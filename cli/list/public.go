package list

import (
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/sean-/conswriter"
	"github.com/sean-/vpc/internal/command"
	"github.com/sean-/vpc/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.freebsd.org/sys/vpc"
	"go.freebsd.org/sys/vpc/mgmt"
)

const (
	_CmdName      = "list"
	_KeyObjCounts = config.KeyListObjCounts
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "list counts of each VPC type",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			cons := conswriter.GetTerminal()

			if viper.GetBool(_KeyObjCounts) {
				return listTypeCount(cons)
			}

			return listTypeIDs(cons)
		},
	},

	Setup: func(self *command.Command) error {
		{
			const (
				key          = _KeyObjCounts
				longName     = "obj-counts"
				shortName    = "c"
				defaultValue = false
				description  = "list the number of objects per type"
			)

			flags := self.Cobra.Flags()
			flags.BoolP(longName, shortName, defaultValue, description)

			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		return nil
	},
}

func listTypeCount(cons conswriter.ConsoleWriter) error {
	table := tablewriter.NewWriter(cons)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetAutoFormatHeaders(true)

	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")

	table.SetHeader([]string{"name", "count"})

	mgr, err := mgmt.New(nil)
	if err != nil {
		return errors.Wrapf(err, "unable to open VPC Management handle")
	}
	defer mgr.Close()

	var numTypes int64
	for _, objType := range vpc.ObjTypes() {
		count, err := mgr.CountType(objType)
		if err != nil {
			return errors.Wrapf(err, "unable to count object type %s", objType)
		}

		table.Append([]string{
			objType.String(),
			strconv.FormatInt(int64(count), 10),
		})
		numTypes++
	}

	table.SetFooter([]string{"total", strconv.FormatInt(numTypes, 10)})

	table.Render()

	return nil
}

func listTypeIDs(cons conswriter.ConsoleWriter) error {
	table := tablewriter.NewWriter(cons)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetAutoFormatHeaders(true)

	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_LEFT})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")

	table.SetHeader([]string{"type", "id", "unit name"})

	mgr, err := mgmt.New(nil)
	if err != nil {
		return errors.Wrapf(err, "unable to open VPC Management handle")
	}
	defer mgr.Close()

	var numIDs int64
	for _, objType := range vpc.ObjTypes() {
		objHeaders, err := mgr.GetAllIDs(objType)
		if err != nil {
			return errors.Wrapf(err, "unable to count object type %s", objType)
		}

		for _, hdr := range objHeaders {
			table.Append([]string{
				hdr.ObjType().String(),
				hdr.ID().String(),
				hdr.UnitName(),
			})
			numIDs++
		}
	}

	table.SetFooter([]string{"total", strconv.FormatInt(numIDs, 10), ""})

	table.Render()

	return nil
}
