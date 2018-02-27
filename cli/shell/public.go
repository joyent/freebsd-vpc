package shell

import (
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/cli/shell/autocompletion"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

const _CmdName = "shell"

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:   _CmdName,
		Short: "shell commands",
	},

	Setup: func(self *command.Command) error {
		subCommands := command.Commands{
			autocompletion.Cmd,
		}

		if err := self.Register(subCommands); err != nil {
			log.Fatal().Err(err).Str("cmd", _CmdName).Msg("unable to register sub-commands")
		}

		return nil
	},
}
