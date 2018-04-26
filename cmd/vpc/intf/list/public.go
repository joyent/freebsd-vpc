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
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/spf13/cobra"
)

const _CmdName = "list"

var Cmd = &command.Command{
	Name: _CmdName,
	Cobra: &cobra.Command{
		Use:          _CmdName,
		Aliases:      []string{"ls"},
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
