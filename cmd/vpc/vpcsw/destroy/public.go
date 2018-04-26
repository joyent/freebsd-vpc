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
	_CmdName     = "destroy"
	_KeySwitchID = config.KeySWDestroySwitchID
)

var Cmd = &command.Command{
	Name: _CmdName,
	Cobra: &cobra.Command{
		Use:              _CmdName,
		Aliases:          []string{"rm", "del", "delete"},
		TraverseChildren: true,
		Short:            "destroy a VPC switch",
		SilenceUsage:     true,
		Args:             cobra.NoArgs,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: runE,
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddSwitchID(self, _KeySwitchID, true); err != nil {
			return errors.Wrap(err, "unable to register ID flag on VPC Switch destroy")
		}

		return nil
	},
}

func runE(cmd *cobra.Command, args []string) error {
	cons := conswriter.GetTerminal()

	cons.Write([]byte(fmt.Sprintf("Destroying VPC Switch...")))

	id, err := flag.GetID(viper.GetViper(), _KeySwitchID)
	if err != nil {
		return errors.Wrap(err, "unable to get VPC ID")
	}

	switchCfg := vpcsw.Config{
		ID:        id,
		Writeable: true,
	}

	log.Info().Object("cfg", switchCfg).Str("op", "destroy").Msg("vpc_ctl")

	vpcSwitch, err := vpcsw.Open(switchCfg)
	if err != nil {
		return errors.Wrap(err, "unable to open VPC Switch")
	}
	defer vpcSwitch.Close()

	if err := vpcSwitch.Destroy(); err != nil {
		return errors.Wrap(err, "unable to destroy VPC Switch")
	}

	cons.Write([]byte("done.\n"))

	return nil
}
