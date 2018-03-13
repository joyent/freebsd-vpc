package vmnic

import (
	"github.com/pkg/errors"
	"github.com/sean-/vpc/cmd/vpc/vmnic/create"
	"github.com/sean-/vpc/cmd/vpc/vmnic/destroy"
	"github.com/sean-/vpc/cmd/vpc/vmnic/get"
	"github.com/sean-/vpc/cmd/vpc/vmnic/list"
	"github.com/sean-/vpc/cmd/vpc/vmnic/set"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

const cmdName = "vmnic"

var Cmd = &command.Command{
	Name: cmdName,

	Cobra: &cobra.Command{
		Use:     cmdName,
		Aliases: []string{"nic", "if", "iface"},
		Short:   "VM network interface management",
	},

	Setup: func(self *command.Command) error {
		subCommands := command.Commands{
			create.Cmd,
			destroy.Cmd,
			get.Cmd,
			list.Cmd,
			set.Cmd,
		}

		if err := self.Register(subCommands); err != nil {
			return errors.Wrapf(err, "unable to register sub-commands under %s", cmdName)
		}

		return nil
	},
}
