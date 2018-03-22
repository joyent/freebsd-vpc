package vmnic

import (
	"github.com/joyent/freebsd-vpc/cmd/vpc/vmnic/create"
	"github.com/joyent/freebsd-vpc/cmd/vpc/vmnic/destroy"
	"github.com/joyent/freebsd-vpc/cmd/vpc/vmnic/genmac"
	"github.com/joyent/freebsd-vpc/cmd/vpc/vmnic/get"
	"github.com/joyent/freebsd-vpc/cmd/vpc/vmnic/list"
	"github.com/joyent/freebsd-vpc/cmd/vpc/vmnic/set"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/pkg/errors"
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
			genmac.Cmd,
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
