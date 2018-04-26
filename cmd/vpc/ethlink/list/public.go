// Copyright (c) 2018 Joyent, Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.

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
	_CmdName   = "list"
	_KeySortBy = config.KeyEthLinkListSortBy
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Aliases:      []string{"ls"},
		Short:        "list VPC EthLink interfaces",
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

			table.SetHeader([]string{"name", "id"})

			mgr, err := mgmt.New(nil)
			if err != nil {
				return errors.Wrapf(err, "unable to open VPC Management handle")
			}
			defer mgr.Close()

			objHeaders, err := mgr.GetAllIDs(vpc.ObjTypeLinkEth)
			if err != nil {
				return errors.Wrapf(err, "unable to count %s VPC objects", vpc.ObjTypeLinkEth)
			}

			sortBy := viper.GetString(_KeySortBy)
			switch k := strings.ToLower(viper.GetString(_KeySortBy)); k {
			case "id":
				sort.SliceStable(objHeaders, func(i, j int) bool { return bytes.Compare(objHeaders[i].ID().Bytes(), objHeaders[j].ID().Bytes()) < 0 })
			case "name":
				sort.SliceStable(objHeaders, func(i, j int) bool { return objHeaders[i].UnitName() < objHeaders[j].UnitName() })
			default:
				return errors.Errorf("unsupported sort option: %q", sortBy)
			}

			for _, hdr := range objHeaders {
				table.Append([]string{
					hdr.UnitName(),
					hdr.ID().String(),
				})
			}

			table.SetFooter([]string{"total", strconv.FormatInt(int64(len(objHeaders)), 10), "", "", ""})

			table.Render()

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		{
			const (
				key          = _KeySortBy
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

		return nil
	},
}
