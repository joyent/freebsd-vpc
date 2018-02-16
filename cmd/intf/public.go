package intf

import (
	"github.com/sean-/vpc/cmd/intf/list"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Use:     "interface",
		Aliases: []string{"int", "intf"},
		Short:   "VPC interface management",
	},

	Setup: func(parent *command.Command) error {
		cmds := []*command.Command{
			list.Cmd,
		}

		for _, cmd := range cmds {
			cmd.Setup(cmd)
			parent.Cobra.AddCommand(cmd.Cobra)
		}

		return nil
	},
}
