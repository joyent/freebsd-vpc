// Copyright (c) 2018 Joyent, Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.

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
