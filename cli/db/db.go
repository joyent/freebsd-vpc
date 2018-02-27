package db

import (
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/cli/db/migrate"
	"github.com/sean-/vpc/cli/db/ping"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

const _CmdName = "db"

var Cmd = &command.Command{
	Name: _CmdName,
	Cobra: &cobra.Command{
		Use:     _CmdName,
		Aliases: []string{"database"},
		Short:   "Interaction with the VPC database",
	},

	Setup: func(self *command.Command) error {
		subCommands := []*command.Command{
			migrate.Cmd,
			ping.Cmd,
		}

		if err := self.Register(subCommands); err != nil {
			log.Fatal().Err(err).Str("cmd", _CmdName).Msg("unable to register sub-commands")
		}

		return nil
	},
}
