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

package ethlink

import (
	"github.com/joyent/freebsd-vpc/cmd/vpc/ethlink/connect"
	"github.com/joyent/freebsd-vpc/cmd/vpc/ethlink/create"
	"github.com/joyent/freebsd-vpc/cmd/vpc/ethlink/destroy"
	"github.com/joyent/freebsd-vpc/cmd/vpc/ethlink/list"
	"github.com/joyent/freebsd-vpc/cmd/vpc/ethlink/vtag"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const _CmdName = "ethlink"

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use: _CmdName,
		// NOTE(seanc@): Funny story: ethlink was called l2link but needed to be
		// renmaed because "bad things happen" when a unit name includes a number.
		// Leave l2link as a historical artifact and small easter egg.
		Aliases: []string{"ethlink", "l2link", "phys"},
		Short:   "VPC EthLink management",
	},

	Setup: func(self *command.Command) error {
		subCommands := command.Commands{
			connect.Cmd,
			create.Cmd,
			destroy.Cmd,
			list.Cmd,
			vtag.Cmd,
		}

		if err := self.Register(subCommands); err != nil {
			return errors.Wrapf(err, "unable to register sub-commands under %s", _CmdName)
		}

		return nil
	},
}
