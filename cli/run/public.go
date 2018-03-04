package run

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/agent"
	"github.com/sean-/vpc/config"
	"github.com/sean-/vpc/internal/buildtime"
	"github.com/sean-/vpc/internal/command"
	internal_config "github.com/sean-/vpc/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const _CmdName = "run"

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
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
	},

	Setup: func(self *command.Command) error {
		{
			const (
				key          = internal_config.KeyPGDatabase
				longName     = "db-name"
				description  = "Database name"
				defaultValue = "triton"
			)

			flags := self.Cobra.PersistentFlags()
			flags.String(longName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		{
			const (
				key          = internal_config.KeyPGUser
				longName     = "db-username"
				description  = "Database username"
				defaultValue = "root"
			)

			flags := self.Cobra.PersistentFlags()
			flags.String(longName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		{
			const (
				key          = internal_config.KeyPGPassword
				longName     = "db-password"
				description  = "Database password"
				defaultValue = "tls"
			)

			flags := self.Cobra.PersistentFlags()
			flags.String(longName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		{
			const (
				key          = internal_config.KeyPGHost
				longName     = "db-host"
				description  = "Database server address"
				defaultValue = "127.0.0.1"
			)

			flags := self.Cobra.PersistentFlags()
			flags.String(longName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		{
			const (
				key          = internal_config.KeyPGPort
				longName     = "db-port"
				description  = "Database port"
				defaultValue = 26257
			)

			flags := self.Cobra.PersistentFlags()
			flags.Uint(longName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		return nil
	},
}
