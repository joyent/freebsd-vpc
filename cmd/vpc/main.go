package main

import (
	"os"

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
