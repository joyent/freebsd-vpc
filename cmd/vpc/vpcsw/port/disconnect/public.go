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

package disconnect

import (
	"fmt"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/vpcp"
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
	_CmdName        = "disconnect"
	_KeyPortID      = config.KeySWPortDisconnectPortID
	_KeyInterfaceID = config.KeySWPortDisconnectInterfaceID
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "disconnect a VPC Interface from a VPC Switch Port",
		Aliases:      []string{"disco"},
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cons := conswriter.GetTerminal()

			cons.Write([]byte(fmt.Sprintf("Disconnecting VPC Interface from VPC Switch Port...")))

			interfaceID, err := flag.GetID(viper.GetViper(), _KeyInterfaceID)
			if err != nil {
				return errors.Wrap(err, "unable to get switch port ID")
			}

			portID, err := flag.GetPortID(viper.GetViper(), _KeyPortID)
			if err != nil {
				return errors.Wrap(err, "unable to get switch port ID")
			}

			portCfg := vpcp.Config{
				ID:        portID,
				Writeable: true,
			}

			vpcPort, err := vpcp.Open(portCfg)
			if err != nil {
				log.Error().Err(err).Str("port-id", portID.String()).Msg("VPC Switch Port open failed")
				return errors.Wrap(err, "unable to open VPC Switch Port")
			}
			defer vpcPort.Close()

			if err = vpcPort.Disconnect(interfaceID); err != nil {
				log.Error().Err(err).Object("port-id", portID).Object("interface-id", interfaceID).Msg("vpc switch port disconnect failed")
				return errors.Wrap(err, "unable to disconnect a VPC Interface from VPC Switch Port")
			}

			cons.Write([]byte("done.\n"))

			log.Info().Object("port-id", portID).Object("interface-id", interfaceID).Msg("VPC Interface disconnected from VPC Switch Port")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddInterfaceID(self, _KeyInterfaceID, true); err != nil {
			return errors.Wrap(err, "unable to register Interface ID flag on VPC Switch Port disconnect")
		}

		if err := flag.AddPortID(self, _KeyPortID, true); err != nil {
			return errors.Wrap(err, "unable to register Port ID flag on VPC Switch Port disconnect")
		}

		return nil
	},
}
