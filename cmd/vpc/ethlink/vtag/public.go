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

package vtag

import (
	"fmt"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc"
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
	cmdName            = "vtag"
	_KeyEthLinkID      = config.KeyEthLinkVTagID
	_KeyEthLinkGetVTag = config.KeyEthLinkGetVTag
	_KeyEthLinkSetVTag = config.KeyEthLinkSetVTag
)

var Cmd = &command.Command{
	Name: cmdName,
	Cobra: &cobra.Command{
		Use:              cmdName,
		Aliases:          []string{"vtag", "vlan"},
		TraverseChildren: true,
		Short:            "Get or set the VTag on a VPC EthLink",
		SilenceUsage:     true,
		Args:             cobra.NoArgs,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: runE,
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddEthLinkID(self, _KeyEthLinkID, true); err != nil {
			return errors.Wrap(err, "unable to register EthLink ID flag on EthLink vtag")
		}

		{
			const (
				key          = _KeyEthLinkGetVTag
				longName     = "get-vtag"
				shortName    = "g"
				defaultValue = true
				description  = "get the VTag for a given VPC EthLink NIC"
			)

			flags := self.Cobra.Flags()
			flags.BoolP(longName, shortName, defaultValue, description)

			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		{
			const (
				key          = _KeyEthLinkSetVTag
				longName     = "set-vtag"
				shortName    = "s"
				defaultValue = -1
				description  = "set the VTag for a given VPC EthLink NIC"
			)

			flags := self.Cobra.Flags()
			flags.IntP(longName, shortName, defaultValue, description)

			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		return nil
	},
}

func runE(_ *cobra.Command, _ []string) error {
	cons := conswriter.GetTerminal()

	ethLinkID, err := flag.GetID(viper.GetViper(), _KeyEthLinkID)
	if err != nil {
		return errors.Wrap(err, "unable to get EthLink VPC ID")
	}

	ethLinkCfg := ethlink.Config{
		ID: ethLinkID,
	}
	if vtagID := viper.GetInt(_KeyEthLinkSetVTag); vtagID >= vpc.VTagMin {
		ethLinkCfg.Writeable = true
	}

	log.Info().Object("cfg", ethLinkCfg).Str("op", "vtag").Msg("vpc_ctl")

	ethLink, err := ethlink.Open(ethLinkCfg)
	if err != nil {
		return errors.Wrap(err, "unable to open VPC EthLink")
	}
	defer ethLink.Close()

	if vtagID := viper.GetInt(_KeyEthLinkSetVTag); vtagID >= vpc.VTagMin {
		cons.Write([]byte(fmt.Sprintf("Setting VPC EthLink VTag...")))
		if err := setVTag(cons, ethLink, vpc.VTag(vtagID)); err != nil {
			return errors.Wrap(err, "unable to set EthLink VTag")
		}
		cons.Write([]byte("done.\n"))
	}

	if viper.GetBool(_KeyEthLinkGetVTag) {
		cons.Write([]byte(fmt.Sprintf("Getting VPC EthLink VTag...")))
		if err := getVTag(cons, ethLink); err != nil {
			return errors.Wrap(err, "unable to get EthLink VTag")
		}
		cons.Write([]byte("done.\n"))
	}

	return nil
}

func getVTag(cons conswriter.ConsoleWriter, ethLink *ethlink.EthLink) error {
	vtagID, err := ethLink.VTagGet()
	if err != nil {
		return errors.Wrap(err, "unable to get VTag from EthLink")
	}

	cons.Write([]byte(fmt.Sprintf("VTag: %d\n", vtagID)))

	return nil
}

func setVTag(cons conswriter.ConsoleWriter, ethLink *ethlink.EthLink, vtagID vpc.VTag) error {
	if err := ethLink.VTagSet(vtagID); err != nil {
		return errors.Wrap(err, "unable to set VTag on EthLink")
	}

	return nil
}
