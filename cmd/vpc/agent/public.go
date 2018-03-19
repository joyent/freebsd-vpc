package agent

import (
	"github.com/joyent/freebsd-vpc/agent"
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
			var config agent.Config
			err := viper.Unmarshal(&config)
			if err != nil {
				log.Fatal().Err(err).Msg("unable to decode config into struct")
			}

			// 2. Run agent
			a, err := agent.New(config)
			if err != nil {
				return errors.Wrapf(err, "unable to create a new %s agent", buildtime.PROGNAME)
			}

			return a.Run()
		},
	},

	Setup: func(parent *command.Command) error {
		return db.SetDefaultViperOptions()
	},
}
