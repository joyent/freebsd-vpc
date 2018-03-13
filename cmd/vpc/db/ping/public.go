package ping

import (
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/cmd/vpc/config"
	"github.com/sean-/vpc/db"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const cmdName = "ping"

var Cmd = &command.Command{
	Name: cmdName,
	Cobra: &cobra.Command{
		Use:          cmdName,
		Short:        "ping the database to ensure connectivity",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("command", "ping").Msg("")

			var config config.Config
			err := viper.Unmarshal(&config)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to decode config into struct")
			}

			dbPool, err := db.New(config.DBConfig)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to create database pool")
			}

			start := time.Now()
			if err := dbPool.Ping(); err != nil {
				return errors.Wrap(err, "unable to ping database")
			}
			elapsed := time.Now().Sub(start)
			log.Info().Str("duration", elapsed.String()).Msg("ping")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		return db.SetDefaultViperOptions()
	},
}
