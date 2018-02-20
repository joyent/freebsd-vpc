package list

import (
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/sean-/conswriter"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
	"go.freebsd.org/sys/vpc/vpctest"
)

const (
	_CmdName = "list"
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "list interfaces",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			cons := conswriter.GetTerminal()

			existingIfaces, err := vpctest.GetAllInterfaces()
			if err != nil {
				return errors.Wrapf(err, "unable to get all interfaces")
			}

			table := tablewriter.NewWriter(cons)
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			table.SetHeaderLine(false)
			table.SetAutoFormatHeaders(true)

			table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT})
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("")
			table.SetColumnSeparator("")
			table.SetRowSeparator("")

			table.SetHeader([]string{"name", "index", "mtu", "mac", "flags"})

			var numInterfaces int64
			for _, iface := range existingIfaces {
				if !strings.HasPrefix(iface.Name, "vpcsw") {
					continue
				}

				table.Append([]string{
					iface.Name,
					strconv.FormatInt(int64(iface.Index), 10),
					strconv.FormatInt(int64(iface.MTU), 10),
					iface.HardwareAddr.String(),
					iface.Flags.String(),
				})
				numInterfaces++
			}

			table.SetFooter([]string{"total", strconv.FormatInt(numInterfaces, 10), "", "", ""})

			table.Render()

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		return nil
	},
}
