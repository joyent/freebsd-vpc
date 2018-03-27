package genmac

import (
	"net"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/olekukonko/tablewriter"
	"github.com/sean-/conswriter"
	"github.com/spf13/cobra"
)

const (
	cmdName = "genmac"
)

var Cmd = &command.Command{
	Name: cmdName,

	Cobra: &cobra.Command{
		Use:          cmdName,
		Short:        "generate a random VPC ID and MAC address for a VPC Hostif Interface",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			cons := conswriter.GetTerminal()

			table := tablewriter.NewWriter(cons)
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			table.SetHeaderLine(false)
			table.SetAutoFormatHeaders(true)

			table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT})
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("")
			table.SetColumnSeparator("")
			table.SetRowSeparator("")

			table.SetHeader([]string{"id", "mac"})

			id := vpc.GenID(vpc.ObjTypeHostif)
			var macAddr net.HardwareAddr = id.Node[:]

			table.Append([]string{
				id.String(),
				macAddr.String(),
			})

			table.Render()

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		return nil
	},
}
