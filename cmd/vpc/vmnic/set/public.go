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

package set

import (
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/vmnic"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/joyent/freebsd-vpc/internal/command/flag"
	"github.com/joyent/freebsd-vpc/internal/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	cmdName        = "set"
	keySetFreeze   = config.KeyVMNICSetFreeze
	keySetNQueues  = config.KeyVMNICSetNQueues
	keySetUnfreeze = config.KeyVMNICSetUnfreeze
	keyVMNICID     = config.KeyVMNICSetVMNICID
)

var Cmd = &command.Command{
	Name: cmdName,

	Cobra: &cobra.Command{
		Use:          cmdName,
		Short:        "set VM NIC information",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := flag.GetID(viper.GetViper(), keyVMNICID)
			if err != nil {
				return errors.Wrap(err, "unable to get VM NIC ID")
			}

			vmnicCfg := vmnic.Config{
				ID: id,
			}
			vmn, err := vmnic.Open(vmnicCfg)
			if err != nil {
				return errors.Wrap(err, "unable to open VM NIC")
			}
			defer vmn.Close()

			if viper.GetBool(keySetFreeze) {
				if err := vmn.Freeze(true); err != nil {
					return errors.Wrapf(err, "unable to freeze the VM NIC")
				}
			}

			if numQueues := viper.GetInt(keySetNQueues); numQueues > 0 {
				if err := vmn.NQueuesSet(uint16(numQueues)); err != nil {
					return errors.Wrapf(err, "unable to set the number of hardware queues")
				}
			}

			if viper.GetBool(keySetUnfreeze) {
				if err := vmn.Freeze(false); err != nil {
					return errors.Wrapf(err, "unable to unfreeze the VM NIC")
				}
			}

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddVMNICID(self, keyVMNICID, true); err != nil {
			return errors.Wrap(err, "unable to register VM NIC ID flag on VM NIC set")
		}

		{
			const (
				key          = keySetNQueues
				longName     = "num-queues"
				shortName    = "n"
				defaultValue = 0
				description  = "set the number of hardware queues for a given VM NIC"
			)

			flags := self.Cobra.Flags()
			flags.IntP(longName, shortName, defaultValue, description)

			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		{
			const (
				key          = keySetFreeze
				longName     = "freeze"
				shortName    = "E"
				defaultValue = false
				description  = "freeze the VM NIC configuration"
			)

			flags := self.Cobra.Flags()
			flags.BoolP(longName, shortName, defaultValue, description)

			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		{
			const (
				key          = keySetUnfreeze
				longName     = "unfreeze"
				shortName    = ""
				defaultValue = false
				description  = "freeze the VM NIC configuration"
			)

			flags := self.Cobra.Flags()
			flags.BoolP(longName, shortName, defaultValue, description)

			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		return nil
	},
}
