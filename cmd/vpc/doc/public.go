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

package doc

import (
	"github.com/joyent/freebsd-vpc/cmd/vpc/doc/man"
	"github.com/joyent/freebsd-vpc/cmd/vpc/doc/md"
	"github.com/joyent/freebsd-vpc/internal/buildtime"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const _CmdName = "doc"

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:     _CmdName,
		Aliases: []string{"docs", "documentation"},
		Short:   "Documentation for " + buildtime.PROGNAME,
	},

	Setup: func(self *command.Command) error {
		subCommands := command.Commands{
			man.Cmd,
			md.Cmd,
		}

		if err := self.Register(subCommands); err != nil {
			log.Fatal().Err(err).Str("cmd", _CmdName).Msg("unable to register sub-commands")
		}

		return nil
	},
}
