package migrate

import (
	"github.com/joyent/freebsd-vpc/db"
	"github.com/joyent/freebsd-vpc/db/migrations"
	"github.com/joyent/freebsd-vpc/internal/buildtime"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	"github.com/mattes/migrate/source/go-bindata"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const cmdName = "migrate"

var Cmd = &command.Command{
	Name: cmdName,
	Cobra: &cobra.Command{
		Use:          cmdName,
		Short:        "Migrate " + buildtime.PROGNAME + " schema",
		SilenceUsage: true,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("cmd", cmdName).Msg("")

			var config struct {
				DBConfig db.Config `mapstructure:"db"`
			}
			err := viper.Unmarshal(&config)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to decode config into struct")
			}

			dbPool, err := db.New(config.DBConfig)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to create database pool")
			}

			// verify db credentials
			if err := dbPool.Ping(); err != nil {
				return errors.Wrap(err, "unable to ping database")
			}

			// Wrap jackc/pgx in an sql.DB-compatible facade.
			db, err := dbPool.STDDB()
			if err != nil {
				return errors.Wrap(err, "unable to conjur up sql.DB facade")
			}

			source, err := bindata.WithInstance(
				bindata.Resource(migrations.AssetNames(),
					func(name string) ([]byte, error) {
						return migrations.Asset(name)
					}))
			if err != nil {
				return errors.Wrap(err, "unable to create migration source")
			}

			if err := db.Ping(); err != nil {
				return errors.Wrap(err, "unable to ping with stdlib driver")
			}

			driver, err := postgres.WithInstance(db, &postgres.Config{})
			if err != nil {
				return errors.Wrap(err, "unable to create migration driver")
			}

			m, err := migrate.NewWithInstance("file:///migrations/crdb/", source,
				config.DBConfig.Database, driver)
			if err != nil {
				return errors.Wrap(err, "unable to create migration")
			}

			if err := m.Down(); err != nil && err != migrate.ErrNoChange {
				return errors.Wrap(err, "unable to downgrade schema")
			}

			if err := m.Up(); err != nil {
				return errors.Wrap(err, "unable to upgrade schema")
			}

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		return db.SetDefaultViperOptions()
	},
}
