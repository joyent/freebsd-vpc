package db

import (
	"github.com/sean-/vpc/cmd/db/migrate"
	"github.com/sean-/vpc/cmd/db/ping"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Use:     "db",
		Aliases: []string{"database"},
		Short:   "Interaction with the VPC database",
	},

	Setup: func(parent *command.Command) error {
		cmds := []*command.Command{
			migrate.Cmd,
			ping.Cmd,
		}

		for _, cmd := range cmds {
			cmd.Setup(cmd)
			parent.Cobra.AddCommand(cmd.Cobra)
		}

		return nil
	},
}
