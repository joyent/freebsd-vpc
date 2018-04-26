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

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/ethlink"
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
	cmdName      = "destroy"
	keyEthLinkID = config.KeyEthLinkDestroyID
)

var Cmd = &command.Command{
	Name: cmdName,
	Cobra: &cobra.Command{
		Use:              cmdName,
		Aliases:          []string{"rm", "del", "delete"},
		TraverseChildren: true,
		Short:            "destroy a VPC EthLink",
		SilenceUsage:     true,
		Args:             cobra.NoArgs,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: runE,
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddEthLinkID(self, keyEthLinkID, true); err != nil {
			return errors.Wrap(err, "unable to register EthLink ID flag on EthLink destroy")
		}

		return nil
	},
}

func runE(_ *cobra.Command, _ []string) error {
	cons := conswriter.GetTerminal()

	cons.Write([]byte(fmt.Sprintf("Destroying VPC EthLink...")))

	ethLinkID, err := flag.GetID(viper.GetViper(), keyEthLinkID)
	if err != nil {
		return errors.Wrap(err, "unable to get EthLink VPC ID")
	}

	ethLinkCfg := ethlink.Config{
		ID:        ethLinkID,
		Writeable: true,
	}

	log.Info().Object("cfg", ethLinkCfg).Str("op", "destroy").Msg("vpc_ctl")

	ethLink, err := ethlink.Open(ethLinkCfg)
	if err != nil {
		return errors.Wrap(err, "unable to open VPC EthLink")
	}
	defer ethLink.Close()

	if err := ethLink.Destroy(); err != nil {
		return errors.Wrap(err, "unable to destroy VPC EthLink")
	}

	cons.Write([]byte("done.\n"))

	return nil
}
