package version

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/internal/buildtime"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

const cmdName = "version"

var Cmd = &command.Command{
	Name: cmdName,

	Cobra: &cobra.Command{
		Use:          cmdName,
		Short:        "Version " + buildtime.PROGNAME + " schema",
		SilenceUsage: true,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			log.Debug().
				Str("commit", buildtime.GitCommit).
				Str("branch", buildtime.GitBranch).
				Str("state", buildtime.GitState).
				Str("summary", buildtime.GitSummary).
				Str("build-date", buildtime.BuildDate).
				Str("version", buildtime.Version).
				Msg("version")

			fmt.Printf("Version: %s\n", buildtime.Version)

			{
				commit := buildtime.GitCommit
				if buildtime.GitState != "clean" {
					commit += "+" + buildtime.GitState
				}

				fmt.Printf("Commit: %s\n", commit)
			}

			fmt.Printf("Build Date: %s\n", buildtime.BuildDate)

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		return nil
	},
}
