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

package command

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SetupFunc func(parent *Command) error
type SetViperFunc func(v *viper.Viper)

type Command struct {
	Name       string
	Setup      SetupFunc
	ArgAliases []string
	ValidArgs  []string
	Cobra      *cobra.Command
}

func (cmd *Command) Register(subCommands Commands) (err error) {
	for i, subCommand := range subCommands {
		// Validate sub-command's configuration
		switch {
		case subCommand.Name == "":
			return fmt.Errorf("%q[%d] has no name registered", cmd.Name, i)
		case strings.ContainsAny(subCommand.Name, "."):
			return errors.Errorf("%q[%d] contains an invalid character: %q", cmd.Name, i, subCommand.Name)
		case subCommand.Cobra == nil:
			return errors.Errorf("%q[%d].%q is missing a Cobra instance", cmd.Name, i, subCommand.Name)
		}

		if subCommand.Setup != nil {
			if err := subCommand.Setup(subCommand); err != nil {
				return errors.Wrapf(err, "unable to register %q subcommands", subCommand.Name)
			}
		}

		cmd.Cobra.DisableAutoGenTag = true
		cmd.Cobra.AddCommand(subCommand.Cobra)
	}

	return nil
}

type Commands []*Command

func (cs Commands) ArgAliases() []string {
	args := make([]string, 0, len(cs))
	for _, cmd := range cs {
		args = append(args, cmd.ArgAliases...)
	}
	return args
}

func (cs Commands) ValidArgs() []string {
	args := make([]string, 0, len(cs))
	for _, cmd := range cs {
		args = append(args, cmd.ValidArgs...)
	}
	return args
}
