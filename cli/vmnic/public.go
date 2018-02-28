package vmnic

import (
	"github.com/pkg/errors"
	"github.com/sean-/vpc/cli/vmnic/create"
	"github.com/sean-/vpc/cli/vmnic/destroy"
	"github.com/sean-/vpc/cli/vmnic/get"
	"github.com/sean-/vpc/cli/vmnic/list"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

const _CmdName = "vmnic"

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:     _CmdName,
		Aliases: []string{"nic", "if", "iface"},
		Short:   "VM network interface management",
	},

	Setup: func(self *command.Command) error {
		subCommands := command.Commands{
			create.Cmd,
			destroy.Cmd,
			get.Cmd,
			list.Cmd,
		}

		if err := self.Register(subCommands); err != nil {
			return errors.Wrapf(err, "unable to register sub-commands under %s", _CmdName)
		}

		return nil
	},
}
