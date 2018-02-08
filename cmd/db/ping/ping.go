package ping

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/internal/command"
	"github.com/sean-/vpc/internal/conswriter"
	"github.com/spf13/cobra"
)

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Use:          "ping",
		Short:        "ping the database to ensure connectivity",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("here\n")
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			cons := conswriter.GetTerminal()

			fmt.Fprintf(cons, "start\n")
			for i := 0; i < 100; i++ {
				fmt.Fprintf(cons, "output\n")
				log.Info().Msg("output")
			}
			fmt.Fprintf(cons, "end\n")

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

	Setup: func(parent *command.Command) error {
		return nil
	},
}
