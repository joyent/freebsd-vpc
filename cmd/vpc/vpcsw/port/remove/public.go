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

package remove

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
	_CmdName     = "remove"
	_KeyPortID   = config.KeySWPortRemovePortID
	_KeySwitchID = config.KeySWPortRemoveSwitchID
)

var Cmd = &command.Command{
	Name: _CmdName,
	Cobra: &cobra.Command{
		Use:              _CmdName,
		Aliases:          []string{"rm", "del", "delete"},
		TraverseChildren: true,
		Short:            "remove a port from a VPC switch",
		SilenceUsage:     true,
		Args:             cobra.NoArgs,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			// TODO(seanc@): Verify that a given port-id belongs to a switch
			return nil
		},

		RunE: runE,
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddPortID(self, _KeyPortID, true); err != nil {
			return errors.Wrap(err, "unable to register Port ID flag on VPC Switch Port remove")
		}

		if err := flag.AddSwitchID(self, _KeySwitchID, false); err != nil {
			return errors.Wrap(err, "unable to register Switch ID flag for VPC Switch Port add")
		}

		return nil
	},
}

func runE(cmd *cobra.Command, args []string) error {
	cons := conswriter.GetTerminal()

	cons.Write([]byte(fmt.Sprintf("Removing Port from VPC Switch...")))

	// 1) get switch ID
	switchID, err := flag.GetSwitchID(viper.GetViper(), _KeySwitchID)
	if err != nil {
		return errors.Wrap(err, "unable to get VPC Switch ID")
	}

	// 2) get port id
	portID, err := flag.GetPortID(viper.GetViper(), _KeyPortID)
	if err != nil {
		return errors.Wrap(err, "unable to get VPC Switch Port ID")
	}

	// 3) open switch
	switchCfg := vpcsw.Config{
		ID:        switchID,
		Writeable: true,
	}

	sw, err := vpcsw.Open(switchCfg)
	if err != nil {
		log.Error().Err(err).Str("switch-id", switchID.String()).Object("switch-cfg", switchCfg).Msg("vpc_open() failed on switch")
		return errors.Wrap(err, "unable to vpc_open(2) switch")
	}
	defer sw.Close()

	// 4) send op to remove
	if err = sw.PortRemove(portID); err != nil {
		log.Error().Err(err).Object("switch-cfg", switchCfg).Object("port-id", portID).Msg("vpc_ctl(2): switch port destroy failed")
		return errors.Wrap(err, "unable to remove VPC Switch Port")
	}

	// 5) close switch
	if err := sw.Close(); err != nil {
		return errors.Wrap(err, "unable to close VPC Switch")
	}

	cons.Write([]byte("done.\n"))

	return nil
}
