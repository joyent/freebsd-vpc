package agent

import (
	"github.com/joyent/freebsd-vpc/agent"
	"github.com/joyent/freebsd-vpc/cmd/vpc/config"
	"github.com/joyent/freebsd-vpc/db"
	"github.com/joyent/freebsd-vpc/internal/buildtime"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const cmdName = "agent"

var Cmd = &command.Command{
	Name: cmdName,

	Cobra: &cobra.Command{
		Use:          cmdName,
		Short:        "Run " + buildtime.PROGNAME,
		SilenceUsage: true,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("command", "run").Msg("")

			// 1. Parse config and construct agent
			var config config.Config
			err := viper.Unmarshal(&config)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to decode config into struct")
			}

			dbPool, err := db.New(config.DBConfig)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to create database pool")
			}

			// 2. Run agent
			a, err := agent.New(dbPool)
			if err != nil {
				return errors.Wrapf(err, "unable to create a new %s agent", buildtime.PROGNAME)
			}

			if err := a.Start(); err != nil {
				return errors.Wrapf(err, "unable to start %s agent", buildtime.PROGNAME)
			}
			defer a.Stop()

			// 3. Connect to the database to verify database credentials

			// 4. Loop until program exit
			if err := a.Run(); err != nil {
				return errors.Wrapf(err, "unable to run %s agent", buildtime.PROGNAME)
			}

			return nil
		},
	},

	Setup: func(parent *command.Command) error {
		return db.SetDefaultViperOptions()
	},
}
