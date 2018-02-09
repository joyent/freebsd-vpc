package ping

import (
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/agent"
	"github.com/sean-/vpc/config"
	"github.com/sean-/vpc/internal/buildtime"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Use:          "ping",
		Short:        "ping the database to ensure connectivity",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("command", "ping").Msg("")

			// 1. Parse config
			cfg, err := config.New()
			if err != nil {
				return errors.Wrap(err, "unable to load configuration")
			}

			if err := cfg.Load(); err != nil {
				return errors.Wrapf(err, "unable to load %s config", buildtime.PROGNAME)
			}

			// 2. Run agent
			a, err := agent.New(cfg)
			if err != nil {
				return errors.Wrapf(err, "unable to create a new %s agent", buildtime.PROGNAME)
			}

			// 3. Connect to the database to verify database credentials
			start := time.Now()
			if err := a.Pool().Ping(); err != nil {
				return errors.Wrap(err, "unable to ping database")
			}
			elapsed := time.Now().Sub(start)
			log.Info().Str("duration", elapsed.String()).Msg("ping")

			return nil
		},
	},

	Setup: func(parent *command.Command) error {
		return nil
	},
}
