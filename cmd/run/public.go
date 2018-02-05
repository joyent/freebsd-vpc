package run

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/smallz/buildtime"
	"github.com/sean-/vpc/agent"
	"github.com/sean-/vpc/config"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:          "run",
	Short:        "Run " + buildtime.PROGNAME,
	SilenceUsage: true,

	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info().Str("command", "run").Msg("")

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

		if err := a.Start(); err != nil {
			return errors.Wrapf(err, "unable to start %s agent", buildtime.PROGNAME)
		}
		defer a.Stop()

		// 3. Connect to the database to verify database credentials
		if err := a.Pool().Ping(); err != nil {
			return errors.Wrap(err, "unable to ping database")
		}

		// 4. Loop until program exit
		if err := a.Run(); err != nil {
			return errors.Wrapf(err, "unable to run %s agent", buildtime.PROGNAME)
		}

		return nil
	},
}
