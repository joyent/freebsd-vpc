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

package uplink

import (
	"fmt"
	"net"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/vpcsw"
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
	_CmdName     = "uplink"
	_KeyPortID   = config.KeySWPortUplinkPortID
	_KeySwitchID = config.KeySWPortUplinkSwitchID
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "Sets a given switch port as the uplink",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cons := conswriter.GetTerminal()

			cons.Write([]byte(fmt.Sprintf("Designating VPC Switch Port as Uplink Port...")))

			portID, err := flag.GetPortID(viper.GetViper(), _KeyPortID)
			if err != nil {
				return errors.Wrap(err, "unable to get switch port ID")
			}

			switchID, err := flag.GetSwitchID(viper.GetViper(), _KeySwitchID)
			if err != nil {
				return errors.Wrap(err, "unable to get switch ID")
			}

			switchCfg := vpcsw.Config{
				ID:        switchID,
				Writeable: true,
			}

			vpcSwitch, err := vpcsw.Open(switchCfg)
			if err != nil {
				log.Error().Err(err).Object("switch-id", switchID).Msg("VPC Switch open failed")
				return errors.Wrap(err, "unable to open VPC Switch")
			}
			defer vpcSwitch.Close()

			var portMAC net.HardwareAddr = portID.Node[:]
			if err := vpcSwitch.PortUplinkSet(portID, portMAC); err != nil {
				log.Error().Err(err).Object("port-id", portID).Object("switch-cfg", switchCfg).Msg("failed to set VPC Switch Port as an Uplink port")
				return errors.Wrap(err, "unable to create a VPC Switch Port uplink")
			}

			cons.Write([]byte("done.\n"))

			log.Info().Object("switch-id", switchID).Object("port-id", portID).Msg("Uplink port for VPC Switch set")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddPortID(self, _KeyPortID, true); err != nil {
			return errors.Wrap(err, "unable to register Port ID flag on VPC Switch Port uplink")
		}

		if err := flag.AddSwitchID(self, _KeySwitchID, true); err != nil {
			return errors.Wrap(err, "unable to register Switch ID flag on VPC Switch Port uplink")
		}

		return nil
	},
}
