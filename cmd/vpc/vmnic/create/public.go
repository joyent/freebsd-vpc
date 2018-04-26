// Copyright (c) 2018 Joyent, Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.

package create

import (
	"fmt"
	"net"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/vmnic"
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/vpctest"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/joyent/freebsd-vpc/internal/command/flag"
	"github.com/joyent/freebsd-vpc/internal/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/conswriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	cmdName     = "create"
	keyVMNICID  = config.KeyVMNICCreateID
	keyVMNICMAC = config.KeyVMNICCreateMAC
)

var Cmd = &command.Command{
	Name: cmdName,

	Cobra: &cobra.Command{
		Use:          cmdName,
		Short:        "create a VM NIC",
		SilenceUsage: true,
		// TraverseChildren: true,
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cons := conswriter.GetTerminal()

			existingIfaces, err := vpctest.GetAllInterfaces()
			if err != nil {
				return errors.Wrapf(err, "unable to get all interfaces")
			}

			cons.Write([]byte(fmt.Sprintf("Creating VM NIC...")))

			id, err := flag.GetID(viper.GetViper(), keyVMNICID)
			if err != nil {
				return errors.Wrap(err, "unable to get VPC ID")
			}

			mac, err := flag.GetMAC(viper.GetViper(), keyVMNICMAC, &id)
			if err != nil {
				return errors.Wrap(err, "unable to get MAC address")
			}

			vmnicCfg := vmnic.Config{
				ID:  id,
				MAC: mac,
			}

			vmNIC, err := vmnic.Create(vmnicCfg)
			if err != nil {
				log.Error().Err(err).Object("vmnic-id", id).Msg("vmnic create failed")
				return errors.Wrap(err, "unable to create VM NIC")
			}
			defer vmNIC.Close()

			if err := vmNIC.Commit(); err != nil {
				log.Error().Err(err).Object("vmnic-id", id).Msg("VM NIC commit failed")
				return errors.Wrap(err, "unable to commit VM NIC")
			}

			cons.Write([]byte("done.\n"))

			var newVMNIC net.Interface
			{ // Get the before/after
				ifacesAfterCreate, err := vpctest.GetAllInterfaces()
				if err != nil {
					return errors.Wrapf(err, "unable to get all interfaces")
				}
				_, newIfaces, _ := existingIfaces.Difference(ifacesAfterCreate)

				var newVMNICMAC net.HardwareAddr = id.Node[:]
				newVMNIC, err = newIfaces.FindMAC(newVMNICMAC)
				if err != nil {
					return errors.Wrapf(err, "unable to find new VM NIC with MAC %q", id.Node)
				}
			}

			log.Info().Object("vmnic-id", id).Str("mac", newVMNIC.HardwareAddr.String()).Str("name", newVMNIC.Name).Msg("VM NIC created")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddVMNICID(self, keyVMNICID, false); err != nil {
			return errors.Wrap(err, "unable to register VM NIC ID flag on VM NIC create")
		}

		if err := flag.AddMAC(self, keyVMNICMAC, false); err != nil {
			return errors.Wrap(err, "unable to register MAC flag on VM NIC create")
		}

		return nil
	},
}
