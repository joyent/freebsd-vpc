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

package add

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
	_CmdName     = "add"
	_KeyPortID   = config.KeySWPortAddID
	_KeyPortMAC  = config.KeySWPortAddMAC
	_KeySwitchID = config.KeySWPortAddSwitchID
	_KeyUplink   = config.KeySWPortAddUplink
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "add a port to a VPC Switch",
		Aliases:      []string{"create"},
		SilenceUsage: true,
		// TraverseChildren: true,
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cons := conswriter.GetTerminal()

			cons.Write([]byte(fmt.Sprintf("Adding port to VPC Switch...")))

			switchID, err := flag.GetSwitchID(viper.GetViper(), _KeySwitchID)
			if err != nil {
				return errors.Wrap(err, "unable to get VPC ID")
			}

			portID, err := flag.GetPortID(viper.GetViper(), _KeyPortID)
			if err != nil {
				return errors.Wrap(err, "unable to get VPC Switch Port ID")
			}

			portMAC, err := flag.GetMAC(viper.GetViper(), _KeyPortMAC, nil)
			if err != nil {
				return errors.Wrap(err, "unable to get MAC address")
			}

			// Create a stack of commit and undo operations to walk through in the
			// event of an error.
			var commit bool
			var commitFuncs, undoFuncs []func() error
			defer func() {
				scopeHandlers := undoFuncs
				modeStr := "undo"
				if commit {
					modeStr = "commit"
					scopeHandlers = commitFuncs
				}

				for i := len(scopeHandlers) - 1; i >= 0; i-- {
					if err := scopeHandlers[i](); err != nil {
						log.Error().Err(err).Msgf("failure during %s", modeStr)
					}
				}
			}()
			commitFuncs = append(commitFuncs, func() error {
				cons.Write([]byte("done.\n"))
				return nil
			})

			// 1) Open switch and add a port
			switchCfg := vpcsw.Config{
				ID:        switchID,
				Writeable: true,
			}

			vpcSwitch, err := vpcsw.Open(switchCfg)
			if err != nil {
				log.Error().Err(err).Object("switch-cfg", switchCfg).Msg("vpcsw open failed")
				return errors.Wrap(err, "unable to open VPC Switch")
			}
			commitFuncs = append(commitFuncs, func() error {
				if err := vpcSwitch.Close(); err != nil {
					log.Error().Err(err).Msg("unable to commit VPC Switch")
					return errors.Wrap(err, "unable to commit VPC switch during operation commit")
				}

				return nil
			})
			undoFuncs = append(undoFuncs, func() error {
				if err := vpcSwitch.Close(); err != nil {
					log.Error().Err(err).Msg("unable to clean up VPC Switch during error recovery")
				}

				return nil
			})

			if err = vpcSwitch.PortAdd(portID, portMAC); err != nil {
				log.Error().Err(err).
					Object("port-id", portID).
					Str("port-mac", portMAC.String()).
					Object("switch-cfg", switchCfg).
					Msg("failed to add VPC Switch Port")
				return errors.Wrap(err, "unable to add a port to VPC Switch")
			}

			commit = true

			// log.Info().Str("port-id", portAddCfg.ID.String()).Str("switch-id", switchID.String()).Str("uplink-id", uplinkID.String()). /*.Str("name", newPort.Name)*/ Msg("vpcp created")
			log.Info().Object("port-id", portID).Str("switch-id", switchID.String()).Msg("vpcp created")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddPortID(self, _KeyPortID, false); err != nil {
			return errors.Wrap(err, "unable to register Port ID flag on VPC Switch Port add")
		}

		if err := flag.AddMAC(self, _KeyPortMAC, false); err != nil {
			return errors.Wrap(err, "unable to register MAC flag on VPC Switch Port add")
		}

		if err := flag.AddSwitchID(self, _KeySwitchID, false); err != nil {
			return errors.Wrap(err, "unable to register Switch ID flag for VPC Switch Port add")
		}

		{
			const (
				key          = _KeyUplink
				longName     = "uplink"
				shortName    = "u"
				defaultValue = false
				description  = "make the port ID an uplink for the switch"
			)

			flags := self.Cobra.Flags()
			flags.BoolP(longName, shortName, defaultValue, description)

			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		return nil
	},
}
