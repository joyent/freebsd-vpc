package create

import (
	"context"
	"net"
	"net/http"

	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const cmdName = "create"

var Cmd = &command.Command{
	Name: cmdName,
	Cobra: &cobra.Command{
		Use:          cmdName,
		Short:        "create and run a new virtual machine",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("command", "vm").Msg("")

			client := http.Client{
				Transport: &http.Transport{
					DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
						return net.Dial("unix", "/tmp/vpc-agent.sock")
					},
				},
			}

			resp, err := client.Get("http://unix/test")
			if err != nil {
				log.Fatal().Err(err).Msg("error making request")
			}

			log.Info().Int("code", resp.StatusCode).Msg("got response")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		return nil
	},
}
