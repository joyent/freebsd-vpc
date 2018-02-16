package vpcsw

import (
	"github.com/sean-/vpc/cmd/vpcsw/create"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

var Cmd = &command.Command{
	ValidArgs: []string{"switch", "sw"},
	Cobra: &cobra.Command{
		Use:     "switch",
		Aliases: []string{"sw"},
		Short:   "VPC switch management",
	},

	Setup: func(parent *command.Command) error {
		cmds := []*command.Command{
			create.Cmd,
			// list.Cmd,
		}

		for _, cmd := range cmds {
			parent.Cobra.AddCommand(cmd.Cobra)
			cmd.Setup(cmd)
		}

		return nil
	},
}
