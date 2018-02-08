package cmd

import (
	"os"
	"path"

	"github.com/google/gops/agent"
	isatty "github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/cmd/db"
	"github.com/sean-/vpc/cmd/run"
	"github.com/sean-/vpc/cmd/version"
	"github.com/sean-/vpc/internal/buildtime"
	"github.com/sean-/vpc/internal/command"
	"github.com/sean-/vpc/internal/config"
	"github.com/sean-/vpc/internal/conswriter"
	"github.com/sean-/vpc/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var subCommands = []*command.Command{
	run.Cmd,
	db.Cmd,
	version.Cmd,
}

var rootCmd = &command.Command{
	Cobra: &cobra.Command{
		Use:   buildtime.PROGNAME,
		Short: buildtime.PROGNAME + " configures and manages VPCs",
	},

	Setup: func(parent *command.Command) error {
		{
			const (
				key         = config.KeyUsePager
				longName    = "use-pager"
				shortName   = "P"
				description = "Use a pager to read the output (defaults to $PAGER, less(1), or more(1))"
			)
			var defaultValue bool
			if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
				defaultValue = true
			}

			flags := parent.Cobra.PersistentFlags()
			flags.BoolP(longName, shortName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		{
			const (
				key          = config.KeyLogLevel
				longOpt      = "log-level"
				shortOpt     = "l"
				defaultValue = "INFO"
				description  = "Change the log level being sent to stdout"
			)

			flags := parent.Cobra.PersistentFlags()
			flags.StringP(longOpt, shortOpt, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longOpt))
			viper.SetDefault(key, defaultValue)
		}

		{
			const (
				key         = config.KeyLogFormat
				longOpt     = "log-format"
				shortOpt    = "F"
				description = `Specify the log format ("auto", "zerolog", or "human")`
			)
			defaultValue := logger.FormatAuto.String()

			flags := parent.Cobra.PersistentFlags()
			flags.StringP(longOpt, shortOpt, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longOpt))
			viper.SetDefault(key, defaultValue)
		}

		{
			const (
				key         = config.KeyLogTermColor
				longOpt     = "use-color"
				shortOpt    = ""
				description = "Use ASCII colors"
			)
			defaultValue := false
			if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
				defaultValue = true
			}

			flags := parent.Cobra.PersistentFlags()
			flags.BoolP(longOpt, shortOpt, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longOpt))
			viper.SetDefault(key, defaultValue)
		}

		{
			const (
				key          = config.KeyUseUTC
				longName     = "utc"
				shortName    = "Z"
				defaultValue = false
				description  = "Display times in UTC"
			)

			flags := parent.Cobra.PersistentFlags()
			flags.BoolP(longName, shortName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		return nil
	},
}

func Execute() error {
	rootCmd.Setup(rootCmd)
	conswriter.UsePager(viper.GetBool(config.KeyUsePager))

	if err := logger.Setup(); err != nil {
		return err
	}

	// Always enable the agent
	//
	// TODO(seanc@): add if viper.GetBool("debug.enable-agent") {
	if err := agent.Listen(&agent.Options{}); err != nil {
		log.Fatal().Err(err).Msg("unable to start gops agent")
	}

	for _, cmd := range subCommands {
		rootCmd.Cobra.AddCommand(cmd.Cobra)
		cmd.Setup(cmd)
	}

	if err := rootCmd.Cobra.Execute(); err != nil {
		return errors.Wrapf(err, "unable to run %s", buildtime.PROGNAME)
	}

	return nil
}

func init() {
	// Initialize viper in order to be able to read values from a config file.
	viper.SetConfigName(buildtime.PROGNAME)
	viper.AddConfigPath(path.Join("$HOME", ".config", buildtime.PROGNAME))
	viper.AddConfigPath(".")

	cobra.OnInitialize(cobraConfig)
}

// cobraConfig reads in config file and ENV variables, if set.
func cobraConfig() {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Debug().Err(err).Msg("unable to read config file")
		} else {
			log.Warn().Err(err).Msg("unable to read config file")
		}
	}
}
