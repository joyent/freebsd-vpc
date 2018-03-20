package agent

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joyent/freebsd-vpc/agent"
	"github.com/joyent/freebsd-vpc/db"
	"github.com/joyent/freebsd-vpc/internal/buildtime"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sys/unix"
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
				return errors.Wrapf(err,"unable to decode config into struct")
			}

			// 2. Run agent
			a, err := agent.New(config)
			if err != nil {
				return errors.Wrapf(err, "unable to create a new %s agent", buildtime.PROGNAME)
			}

			if err := a.Start(); err != nil {
				return errors.Wrapf(err, "unable to start agent")
			}

			signalCh := make(chan os.Signal, 10)
			signal.Notify(signalCh, os.Interrupt, unix.SIGTERM, unix.SIGPIPE)

			for {
				var sig os.Signal
				select {
				case s := <-signalCh:
					sig = s
				}

				switch sig {
				case syscall.SIGPIPE:
					continue

				default:
					log.Info().Str("signal", sig.String()).Msg("caught signal")

					log.Info().Msg("initiating graceful shutdown of agent")
					gracefulCh := make(chan struct{})
					go func() {
						if err := a.Shutdown(); err != nil {
							log.Fatal().Err(err).Msg("error during agent shutdown")
							return
						}
						close(gracefulCh)
					}()

					gracefulTimeout := 15 * time.Second
					select {
					case <-signalCh:
						log.Info().Str("signal", sig.String()).Msg("caught second signal, exiting")
						os.Exit(1)
					case <-time.After(gracefulTimeout):
						log.Info().Dur("timeout", gracefulTimeout).Msg("timeout on graceful shutdown, exiting")
						os.Exit(1)
					case <-gracefulCh:
						log.Info().Msg("graceful shutdown complete")
						os.Exit(0)
					}
				}
			}
		},
	},

	Setup: func(parent *command.Command) error {
		if err := db.SetDefaultViperOptions(); err != nil {
			return err
		}

		if err := setAgentDefaultViperOptions(); err != nil {
			return err
		}

		return nil
	},
}

func setAgentDefaultViperOptions() error {
	viper.SetDefault("agent.addresses.internal", "/tmp/vpc-agent.sock")

	return nil
}
