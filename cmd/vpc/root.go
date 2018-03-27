package main

import (
	"os"
	"path"

	gopsagent "github.com/google/gops/agent"
	"github.com/joyent/freebsd-vpc/cmd/vpc/agent"
	"github.com/joyent/freebsd-vpc/cmd/vpc/db"
	"github.com/joyent/freebsd-vpc/cmd/vpc/doc"
	"github.com/joyent/freebsd-vpc/cmd/vpc/ethlink"
	"github.com/joyent/freebsd-vpc/cmd/vpc/hostif"
	"github.com/joyent/freebsd-vpc/cmd/vpc/intf"
	"github.com/joyent/freebsd-vpc/cmd/vpc/list"
	"github.com/joyent/freebsd-vpc/cmd/vpc/shell"
	"github.com/joyent/freebsd-vpc/cmd/vpc/version"
	"github.com/joyent/freebsd-vpc/cmd/vpc/vm"
	"github.com/joyent/freebsd-vpc/cmd/vpc/vmnic"
	"github.com/joyent/freebsd-vpc/cmd/vpc/vpcsw"
	"github.com/joyent/freebsd-vpc/internal/buildtime"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/joyent/freebsd-vpc/internal/config"
	"github.com/joyent/freebsd-vpc/internal/logger"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/conswriter"
	"github.com/sean-/seed"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const cmdName = "root"

var subCommands = command.Commands{
	db.Cmd,
	doc.Cmd,
	ethlink.Cmd,
	intf.Cmd,
	hostif.Cmd,
	list.Cmd,
	agent.Cmd,
	shell.Cmd,
	version.Cmd,
	vm.Cmd,
	vmnic.Cmd,
	vpcsw.Cmd,
}

var rootCmd = &command.Command{
	Name: cmdName,

	Cobra: &cobra.Command{
		Use:   buildtime.PROGNAME,
		Short: buildtime.PROGNAME + " configures and manages VPCs",
		//ValidArgs:  subCommands.ValidArgs(),
		//ArgAliases: subCommands.ArgAliases(),
		Example: `# Perform a setup for a VM NIC
$ doas vpc switch create --vni=123 --switch-id=da64c3f3-095d-91e5-df13-5aabcfc52468
$ doas vpc vmnic create --vmnic-id=07f95a11-6788-2ae7-c3ce-ba95cff1db38
$ doas vpc vmnic set --vmnic-id=07f95a11-6788-2ae7-c3ce-ba95cff1db38 --num-queues=2
$ doas vpc switch port add --switch-id=da64c3f3-095d-91e5-df13-5aabcfc52468 --port-id=935cf569-17aa-11e8-a53f-507b9da3d9d0
$ doas vpc switch port connect --port-id=935cf569-17aa-11e8-a53f-507b9da3d9d0 --interface-id=07f95a11-6788-2ae7-c3ce-ba95cff1db38
$ doas vpc switch port add --switch-id=da64c3f3-095d-91e5-df13-5aabcfc52468 --port-id=ea58b648-203b-a707-cdf6-7a552c8d5295 --uplink --l2-name=em0 --ethlink-id=5c4acd32-1b8d-11e8-b4c7-0cc47a6c7d1e

$ vpc vmnic get --vmnic-id=07f95a11-6788-2ae7-c3ce-ba95cff1db38
$ vpc list

# Perform a tear down of the above
$ doas vpc ethlink destroy --ethlink-id=5c4acd32-1b8d-11e8-b4c7-0cc47a6c7d1e
$ doas vpc switch port disconnect --port-id=935cf569-17aa-11e8-a53f-507b9da3d9d0 --interface-id=07f95a11-6788-2ae7-c3ce-ba95cff1db38
$ doas vpc vmnic destroy --vmnic-id=07f95a11-6788-2ae7-c3ce-ba95cff1db38
$ doas vpc switch port remove --switch-id=da64c3f3-095d-91e5-df13-5aabcfc52468 --port-id=935cf569-17aa-11e8-a53f-507b9da3d9d0
$ doas vpc switch destroy --switch-id=da64c3f3-095d-91e5-df13-5aabcfc52468
$ vpc list
`,
	},

	Setup: func(self *command.Command) error {
		{
			const (
				key          = config.KeyUseGoogleAgent
				longName     = "enable-gops"
				shortName    = ""
				description  = "Enable the Google gops(1) socket"
				defaultValue = true
			)

			flags := self.Cobra.PersistentFlags()
			flags.BoolP(longName, shortName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
			flags.MarkHidden(longName)
		}

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

			flags := self.Cobra.PersistentFlags()
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

			flags := self.Cobra.PersistentFlags()
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

			flags := self.Cobra.PersistentFlags()
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

			flags := self.Cobra.PersistentFlags()
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

			flags := self.Cobra.PersistentFlags()
			flags.BoolP(longName, shortName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		return nil
	},
}

func Execute() error {
	rootCmd.Setup(rootCmd)
	//TODO(jen20) fix this up
	if os.Getenv("USE_PAGER") != "" {
		conswriter.UsePager(viper.GetBool(config.KeyUsePager))
	}

	if err := logger.Setup(viper.GetViper()); err != nil {
		return err
	}

	if !viper.GetBool(config.KeyUseGoogleAgent) {
		log.Debug().Msg("gops(1) agent disabled by request")
	} else {
		go func() {
			log.Debug().Msg("starting gops(1) agent")
			if err := gopsagent.Listen(gopsagent.Options{}); err != nil {
				log.Fatal().Err(err).Msg("unable to start gops agent")
			}
		}()
	}

	if secure, err := seed.Init(); !secure {
		log.Fatal().Err(err).Msg("unable to securely seed RNG")
	}

	if err := rootCmd.Register(subCommands); err != nil {
		log.Fatal().Err(err).Str("cmd", cmdName).Msg("unable to register sub-commands")
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
