package main

import (
	"os"

	"github.com/joyent/freebsd-vpc/internal/buildtime"
	"github.com/rs/zerolog/log"
	"github.com/sean-/conswriter"
	"github.com/sean-/sysexits"
)

var (
	// Variables populated by govvv(1).
	Version    = "dev"
	BuildDate  string
	DocsDate   string
	GitCommit  string
	GitBranch  string
	GitState   string
	GitSummary string
)

func realmain() int {
	exportBuildtimeConsts()

	defer func() {
		p := conswriter.GetTerminal()
		p.Wait()
	}()

	if err := Execute(); err != nil {
		log.Error().Err(err).Msg("unable to run")
		os.Exit(sysexits.Software)
	}

	return sysexits.OK
}

func main() {
	os.Exit(realmain())
}

func exportBuildtimeConsts() {
	buildtime.GitCommit = GitCommit
	buildtime.GitBranch = GitBranch
	buildtime.GitState = GitState
	buildtime.GitSummary = GitSummary
	buildtime.BuildDate = BuildDate
	if DocsDate != "" {
		buildtime.DocsDate = DocsDate
	} else {
		buildtime.DocsDate = BuildDate
	}
	buildtime.Version = Version
}
