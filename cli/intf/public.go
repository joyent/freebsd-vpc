package intf

import (
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/cli/intf/list"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

const _CmdName = "interface"

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:     _CmdName,
		Aliases: []string{"int", "intf"},
		Short:   "VPC interface management",
	},

	Setup: func(self *command.Command) error {
		subCommands := command.Commands{
			list.Cmd,
		}

		if err := self.Register(subCommands); err != nil {
			log.Fatal().Err(err).Str("cmd", _CmdName).Msg("unable to register sub-commands")
		}

		return nil
	},
}
