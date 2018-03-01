package get

import (
	"strconv"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/vmnic"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/sean-/conswriter"
	"github.com/sean-/vpc/internal/command"
	"github.com/sean-/vpc/internal/command/flag"
	"github.com/sean-/vpc/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	_CmdName       = "get"
	_KeyGetNQueues = config.KeyVMNICGetNQueues
	_KeyVMNICID    = config.KeyVMNICGetVMNICID
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "get VMNIC information",
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

			table.SetHeader([]string{"id", "key", "value"})

			id, err := flag.GetID(viper.GetViper(), _KeyVMNICID)
			if err != nil {
				return errors.Wrap(err, "unable to get VM NIC ID")
			}

			vmnicCfg := vmnic.Config{
				ID: id,
			}
			vmn, err := vmnic.Open(vmnicCfg)
			if err != nil {
				return errors.Wrap(err, "unable to open VM NIC")
			}
			defer vmn.Close()

			if viper.GetBool(_KeyGetNQueues) {
				numQueues, err := vmn.NQueuesGet()
				if err != nil {
					return errors.Wrapf(err, "unable to get the number of hardware queues")
				}

				table.Append([]string{
					id.String(),
					"num-queues",
					strconv.FormatInt(int64(numQueues), 10),
				})
			}

			table.Render()

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		{
			const (
				key          = _KeyGetNQueues
				longName     = "num-queues"
				shortName    = "n"
				defaultValue = true
				description  = "get the number of hardware queues for a given VM NIC"
			)

			flags := self.Cobra.Flags()
			flags.BoolP(longName, shortName, defaultValue, description)

			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		if err := flag.AddVMNICID(self, _KeyVMNICID, true); err != nil {
			return errors.Wrap(err, "unable to register VMNIC ID flag on VM NIC get")
		}

		return nil
	},
}
