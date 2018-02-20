package autocompletion

import (
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/cmd/shell/autocompletion/bash"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

const _CmdName = "autocomplete"

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:   _CmdName,
		Short: "Autocompletion generation",
	},

	Setup: func(self *command.Command) error {
		subCommands := command.Commands{
			bash.Cmd,
		}

		if err := self.Register(subCommands); err != nil {
			log.Fatal().Str("cmd", _CmdName).Err(err).Msg("unable to register sub-commands")
		}

		return nil
	},
}
