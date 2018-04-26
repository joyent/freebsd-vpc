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

package destroy

import (
	"fmt"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/hostif"
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
	cmdName         = "destroy"
	_KeyInterfaceID = config.KeyHostifDestroyID
)

var Cmd = &command.Command{
	Name: cmdName,
	Cobra: &cobra.Command{
		Use:              cmdName,
		Aliases:          []string{"rm", "del", "delete"},
		TraverseChildren: true,
		Short:            "destroy a Hostif NIC",
		SilenceUsage:     true,
		Args:             cobra.NoArgs,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: runE,
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddHostifID(self, _KeyInterfaceID, true); err != nil {
			return errors.Wrap(err, "unable to register VPC Hostif ID flag on VPC Hostif destroy")
		}

		return nil
	},
}

func runE(cmd *cobra.Command, args []string) error {
	cons := conswriter.GetTerminal()

	cons.Write([]byte(fmt.Sprintf("Destroying VPC Hostif...")))

	id, err := flag.GetID(viper.GetViper(), _KeyInterfaceID)
	if err != nil {
		return errors.Wrap(err, "unable to get Hostif VPC ID")
	}

	hostifCfg := hostif.Config{
		ID:        id,
		Writeable: true,
	}

	// TODO(seanc@): Go back and add hostif/vmnic/vpcsw to other commands
	log.Info().Object("cfg", hostifCfg).Str("op", "destroy").Msg("vpc_ctl")

	hostifNIC, err := hostif.Open(hostifCfg)
	if err != nil {
		return errors.Wrap(err, "unable to open VPC Hostif NIC")
	}
	defer hostifNIC.Close()

	if err := hostifNIC.Destroy(); err != nil {
		return errors.Wrap(err, "unable to destroy VPC Hostif NIC")
	}

	cons.Write([]byte("done.\n"))

	return nil
}
