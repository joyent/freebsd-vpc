package autocompletion

import (
	"github.com/sean-/vpc/cmd/shell/autocompletion/bash"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Use:   "autocomplete",
		Short: "Autocompletion generation",
	},

	Setup: func(parent *command.Command) error {
		cmds := []*command.Command{
			bash.Cmd,
		}

		for _, cmd := range cmds {
			cmd.Setup(cmd)
			parent.Cobra.AddCommand(cmd.Cobra)
		}

		return nil
	},
}
