package intf

import (
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/cmd/vpc/intf/list"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

const cmdName = "interface"

var Cmd = &command.Command{
	Name: cmdName,

	Cobra: &cobra.Command{
		Use:     cmdName,
		Aliases: []string{"int", "intf"},
		Short:   "VPC interface management",
	},

	Setup: func(self *command.Command) error {
		subCommands := command.Commands{
			list.Cmd,
		}

		if err := self.Register(subCommands); err != nil {
			log.Fatal().Err(err).Str("cmd", cmdName).Msg("unable to register sub-commands")
		}

		return nil
	},
}
