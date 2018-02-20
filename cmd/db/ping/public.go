package ping

import (
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/agent"
	"github.com/sean-/vpc/config"
	"github.com/sean-/vpc/internal/buildtime"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

const _CmdName = "ping"

var Cmd = &command.Command{
	Name: _CmdName,
	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "ping the database to ensure connectivity",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info().Str("command", "ping").Msg("")

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

			// 3. Connect to the database to verify database credentials
			start := time.Now()
			if err := a.Pool().Ping(); err != nil {
				return errors.Wrap(err, "unable to ping database")
			}
			elapsed := time.Now().Sub(start)
			log.Info().Str("duration", elapsed.String()).Msg("ping")

			return nil
			// tritonClientConfig, err := api.InitConfig()
			// if err != nil {
			// 	return err
			// }

			// client, err := tritonClientConfig.GetComputeClient()
			// if err != nil {
			// 	return err
			// }

			// instances, err := client.Instances().List(context.Background(), &compute.ListInstancesInput{})
			// if err != nil {
			// 	return err
			// }

			// table := tablewriter.NewWriter(cons)
			// table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			// table.SetHeaderLine(false)
			// table.SetAutoFormatHeaders(true)

			// table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT})
			// table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			// table.SetCenterSeparator("")
			// table.SetColumnSeparator("")
			// table.SetRowSeparator("")

			// table.SetHeader([]string{"id", "name", "image", "package"})

			// var numInstances uint
			// for _, instance := range instances {
			// 	table.Append([]string{instance.ID, instance.Name, instance.Image, instance.Package})
			// 	numInstances++
			// }

			// table.Render()

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		return nil
	},
}
