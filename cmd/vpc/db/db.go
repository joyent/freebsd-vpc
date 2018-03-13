package db

import (
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/cmd/vpc/db/migrate"
	"github.com/sean-/vpc/cmd/vpc/db/ping"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

const cmdName = "db"

var Cmd = &command.Command{
	Name: cmdName,
	Cobra: &cobra.Command{
		Use:     cmdName,
		Aliases: []string{"database"},
		Short:   "Interaction with the VPC database",
	},

	Setup: func(self *command.Command) error {
		subCommands := []*command.Command{
			migrate.Cmd,
			ping.Cmd,
		}

		if err := self.Register(subCommands); err != nil {
			log.Fatal().Err(err).Str("cmd", cmdName).Msg("unable to register sub-commands")
		}

		return nil
	},
}
