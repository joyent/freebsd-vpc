package autocompletion

import (
	"github.com/joyent/freebsd-vpc/cmd/vpc/shell/autocompletion/bash"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const cmdName = "autocomplete"

var Cmd = &command.Command{
	Name: cmdName,

	Cobra: &cobra.Command{
		Use:   cmdName,
		Short: "Autocompletion generation",
	},

	Setup: func(self *command.Command) error {
		subCommands := command.Commands{
			bash.Cmd,
		}

		if err := self.Register(subCommands); err != nil {
			log.Fatal().Str("cmd", cmdName).Err(err).Msg("unable to register sub-commands")
		}

		return nil
	},
}
