package main

import (
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/buildtime"
	"github.com/sean-/vpc/cmd/migrate"
	"github.com/sean-/vpc/cmd/run"
	"github.com/sean-/vpc/cmd/version"
	"github.com/spf13/cobra"
)

var (
	// These fields are populated by govvv
	Version    = "dev"
	BuildDate  string
	GitCommit  string
	GitBranch  string
	GitState   string
	GitSummary string
)

var rootCmd = &cobra.Command{
	Use:   buildtime.PROGNAME,
	Short: buildtime.PROGNAME + " configures and manages VPCs",
}

func main() {
	exportBuildtimeConsts()

	addCommands()

	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("unable to run")
	}
}

func addCommands() {
	addCommand(run.Cmd)
	addCommand(migrate.Cmd)
	addCommand(version.Cmd)
}

func addCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}

func exportBuildtimeConsts() {
	buildtime.GitCommit = GitCommit
	buildtime.GitBranch = GitBranch
	buildtime.GitState = GitState
	buildtime.GitSummary = GitSummary
	buildtime.BuildDate = BuildDate
	buildtime.Version = Version
}
