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

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/mux"
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
	_CmdName  = "disconnect"
	_KeyMuxID = config.KeyMuxDisconnectMuxID
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "disconnect a VPC Mux from its VPC Interface",
		Aliases:      []string{"disco"},
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cons := conswriter.GetTerminal()

			cons.Write([]byte(fmt.Sprintf("Disconnecting VPC Mux from VPC Interface...")))

			muxID, err := flag.GetMuxID(viper.GetViper(), _KeyMuxID)
			if err != nil {
				return errors.Wrap(err, "unable to get switch port ID")
			}

			muxCfg := mux.Config{
				ID:        muxID,
				Writeable: true,
			}

			vpcMux, err := mux.Open(muxCfg)
			if err != nil {
				log.Error().Err(err).Object("mux-id", muxID).Msg("VPC Mux open failed")
				return errors.Wrap(err, "unable to open VPC Mux")
			}
			defer vpcMux.Close()

			if err = vpcMux.Disconnect(); err != nil {
				log.Error().Err(err).Object("mux-id", muxID).Msg("vpc mux disconnect failed")
				return errors.Wrap(err, "unable to disconnect a VPC Mux from its VPC Interface")
			}

			cons.Write([]byte("done.\n"))

			log.Info().Object("mux-id", muxID).Msg("VPC Mux disconnected from VPC Interface")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddMuxID(self, _KeyMuxID, true); err != nil {
			return errors.Wrap(err, "unable to register Mux ID flag on VPC Mux disconnect")
		}

		return nil
	},
}
