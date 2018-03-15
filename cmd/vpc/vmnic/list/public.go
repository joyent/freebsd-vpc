package list

import (
	"strconv"
	"strings"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/vpctest"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/sean-/conswriter"
	"github.com/spf13/cobra"
)

const (
	cmdName = "list"
)

var Cmd = &command.Command{
	Name: cmdName,

	Cobra: &cobra.Command{
		Use:          cmdName,
		Aliases:      []string{"ls"},
		Short:        "list VM NICs",
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
				if !strings.HasPrefix(iface.Name, "vmnic") {
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
