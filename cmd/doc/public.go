package doc

import (
	"github.com/sean-/vpc/cmd/doc/man"
	"github.com/sean-/vpc/internal/buildtime"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Use:     "doc",
		Aliases: []string{"docs", "documentation"},
		Short:   "Documentation for " + buildtime.PROGNAME,
	},

	Setup: func(parent *command.Command) error {
		cmds := []*command.Command{
			man.Cmd,
		}

		for _, cmd := range cmds {
			cmd.Setup(cmd)
			parent.Cobra.AddCommand(cmd.Cobra)
		}

		return nil
	},
}
