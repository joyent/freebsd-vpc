package list

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc"
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/mgmt"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/joyent/freebsd-vpc/internal/config"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/sean-/conswriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	cmdName      = "list"
	keyObjCounts = config.KeyListObjCounts
	keySortBy    = config.KeyListObjSortBy
	keyType      = config.KeyListObjType
)

var Cmd = &command.Command{
	Name: cmdName,

	Cobra: &cobra.Command{
		Use:          cmdName,
		Aliases:      []string{"ls"},
		Short:        "list counts of each VPC type",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		Long: `The list operation of vpc(8) is used to display all VPC objects in the system
and their respective VPC IDs.`,
		Example: `% vpc list
 TYPE     ID                                    UNIT NAME
 ethlink  5c4acd32-1b8d-11e8-b408-0cc47a6c7d1e  ethlink0
 vmnic    07f95a11-6788-2ae7-c306-ba95cff1db38  vmnic0
 vmnic    a774ba3a-1f77-11e8-8006-0cc47a6c7d1e  vmnic1
 vpcp     0ebf50e1-1f79-11e8-8002-0cc47a6c7d1e  vpcp1
 vpcp     ea58b648-203b-a707-cd02-7a552c8d5295  vpcp2
 vpcp     fd436f9c-1f77-11e8-8002-0cc47a6c7d1e  vpcp0
 vpcsw    da64c3f3-095d-91e5-df01-5aabcfc52468  vpcsw0

   TOTAL                    7`,

		RunE: func(cmd *cobra.Command, args []string) error {
			cons := conswriter.GetTerminal()

			if viper.GetBool(keyObjCounts) {
				return listTypeCount(cons)
			}

			return listTypeIDs(cons)
		},
	},

	Setup: func(self *command.Command) error {
		{
			const (
				key          = keyObjCounts
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

		{
			const (
				key          = keySortBy
				longName     = "sort-by"
				shortName    = "s"
				defaultValue = "id"
			)
			sortOptions := []string{"id", "name"}
			sortOptionsStr := strings.Join(sortOptions, ", ")
			description := fmt.Sprintf("Change the sort order within a given type: %s", sortOptionsStr)

			flags := self.Cobra.Flags()
			flags.StringP(longName, shortName, defaultValue, description)

			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		{
			const (
				key          = keyType
				longName     = "obj-type"
				shortName    = "t"
				defaultValue = "all"
			)
			objTypes := vpc.ObjTypes()
			sort.SliceStable(objTypes, func(i, j int) bool { return objTypes[i].String() < objTypes[j].String() })

			objTypesStrs := make([]string, len(objTypes))
			for i := range objTypes {
				objTypesStrs[i] = objTypes[i].String()
			}
			objTypesStr := strings.Join(objTypesStrs, ", ")
			description := fmt.Sprintf("List objects of a given type. Valid types: %s", objTypesStr)

			flags := self.Cobra.Flags()
			flags.StringP(longName, shortName, defaultValue, description)

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

	var objTypes []vpc.ObjType
	{
		objTypes = vpc.ObjTypes()
		sort.SliceStable(objTypes, func(i, j int) bool { return objTypes[i].String() < objTypes[j].String() })

		wantObjTypeStr := viper.GetString(keyType)
		if objTypeStr := strings.ToLower(wantObjTypeStr); objTypeStr != "all" {
			var found bool
			for _, objType := range objTypes {
				if objTypeStr == strings.ToLower(objType.String()) {
					found = true
					objTypes = []vpc.ObjType{objType}
					break
				}
			}

			if !found {
				return errors.Errorf("unsupported VPC Object Type %q", wantObjTypeStr)
			}
		}
	}

	var numIDs int64
	for _, objType := range objTypes {
		objHeaders, err := mgr.GetAllIDs(objType)
		if err != nil {
			return errors.Wrapf(err, "unable to count object type %s", objType)
		}

		sortBy := viper.GetString(keySortBy)
		switch k := strings.ToLower(viper.GetString(keySortBy)); k {
		case "id":
			sort.SliceStable(objHeaders, func(i, j int) bool { return bytes.Compare(objHeaders[i].ID().Bytes(), objHeaders[j].ID().Bytes()) < 0 })
		case "name":
			sort.SliceStable(objHeaders, func(i, j int) bool { return objHeaders[i].UnitName() < objHeaders[j].UnitName() })
		default:
			return errors.Errorf("unsupported sort option: %q", sortBy)
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
