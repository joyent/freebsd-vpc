package db

import (
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/cli/db/migrate"
	"github.com/sean-/vpc/cli/db/ping"
	"github.com/sean-/vpc/internal/command"
	"github.com/sean-/vpc/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const _CmdName = "db"

var Cmd = &command.Command{
	Name: _CmdName,
	Cobra: &cobra.Command{
		Use:     _CmdName,
		Aliases: []string{"database"},
		Short:   "Interaction with the VPC database",
	},

	Setup: func(self *command.Command) error {
		subCommands := []*command.Command{
			migrate.Cmd,
			ping.Cmd,
		}

		{
			const (
				key          = config.KeyPGDatabase
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
				key          = config.KeyPGUser
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
				key          = config.KeyPGPassword
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
				key          = config.KeyPGHost
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
				key          = config.KeyPGPort
				longName     = "db-port"
				description  = "Database port"
				defaultValue = 26257
			)

			flags := self.Cobra.PersistentFlags()
			flags.Uint(longName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		if err := self.Register(subCommands); err != nil {
			log.Fatal().Err(err).Str("cmd", _CmdName).Msg("unable to register sub-commands")
		}

		return nil
	},
}
