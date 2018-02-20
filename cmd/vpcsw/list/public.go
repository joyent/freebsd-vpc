package list

import (
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Use:          "list",
		Short:        "list interfaces",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			// mib := []int{unix.CTL_KERN, unix.KERN_OSTYPE}
			// buf := [256]byte{}
			// n := unsafe.Sizeof(buf)
			// if err := unix.SysctlRaw("kern.hostname", mib, &uname.Sysname[0], &n, nil, 0); err != nil {
			// 	return err
			// }

			// log.Info().Str("hostname", h).Msg("list")

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
