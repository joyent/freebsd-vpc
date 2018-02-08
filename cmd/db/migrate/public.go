package migrate

import (
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	bindata "github.com/mattes/migrate/source/go-bindata"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/agent"
	"github.com/sean-/vpc/config"
	"github.com/sean-/vpc/db/migrations"
	"github.com/sean-/vpc/internal/buildtime"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Use:          "migrate",
		Short:        "Migrate " + buildtime.PROGNAME + " schema",
		SilenceUsage: true,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("command", "migrate").Msg("")

			cfg, err := config.New()
			if err != nil {
				return errors.Wrap(err, "unable to load configuration")
			}

			if err := cfg.Load(); err != nil {
				return errors.Wrapf(err, "unable to load %s config", buildtime.PROGNAME)
			}

			a, err := agent.New(cfg)
			if err != nil {
				return errors.Wrapf(err, "unable to create a new %s agent", buildtime.PROGNAME)
			}
			defer a.Shutdown()

			// verify db credentials
			if err := a.Pool().Ping(); err != nil {
				return errors.Wrap(err, "unable to ping database")
			}

			// Wrap jackc/pgx in an sql.DB-compatible facade.
			db, err := a.Pool().STDDB()
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
				cfg.DB.PoolConfig.Database, driver)
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

	Setup: func(parent *command.Command) error {
		return nil
	},
}
