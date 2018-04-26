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

package show

import (
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/mux"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/joyent/freebsd-vpc/internal/command/flag"
	"github.com/joyent/freebsd-vpc/internal/config"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/conswriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	_CmdName  = "show"
	_KeyMuxID = config.KeyMuxShowMuxID
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Aliases:      []string{"get", "status"},
		Short:        "show VPC Mux information",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			cons := conswriter.GetTerminal()

			table := tablewriter.NewWriter(cons)
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			table.SetHeaderLine(false)
			table.SetAutoFormatHeaders(true)

			table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT})
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("")
			table.SetColumnSeparator("")
			table.SetRowSeparator("")

			table.SetHeader([]string{"name", "id"})

			muxID, err := flag.GetMuxID(viper.GetViper(), _KeyMuxID)
			if err != nil {
				return errors.Wrap(err, "unable to get VPC Mux ID")
			}

			muxCfg := mux.Config{
				ID: muxID,
			}

			vpcMux, err := mux.Open(muxCfg)
			if err != nil {
				log.Error().Err(err).Object("mux-id", muxID).Msg("VPC Mux open failed")
				return errors.Wrap(err, "unable to open VPC Mux")
			}
			defer vpcMux.Close()

			interfaceID, err := vpcMux.ConnectedID()
			if err != nil {
				return errors.Wrapf(err, "unable to get VPC ID for interface connected to VPC Mux")
			}

			host, port, err := vpcMux.ListenAddr()
			if err != nil {
				return errors.Wrapf(err, "unable to get VPC Mux listening address")
			}

			table.Append([]string{"interface-id", interfaceID.String()})
			table.Append([]string{"listen-addr", host})
			table.Append([]string{"listen-port", port})

			table.Render()

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		return nil
	},
}
