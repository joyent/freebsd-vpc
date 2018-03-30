package port

import (
	"github.com/joyent/freebsd-vpc/cmd/vpc/vpcsw/port/add"
	"github.com/joyent/freebsd-vpc/cmd/vpc/vpcsw/port/connect"
	"github.com/joyent/freebsd-vpc/cmd/vpc/vpcsw/port/disconnect"
	"github.com/joyent/freebsd-vpc/cmd/vpc/vpcsw/port/remove"
	"github.com/joyent/freebsd-vpc/cmd/vpc/vpcsw/port/set"
	"github.com/joyent/freebsd-vpc/cmd/vpc/vpcsw/port/uplink"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const cmdName = "port"

var Cmd = &command.Command{
	Name: cmdName,

	Cobra: &cobra.Command{
		Use:     cmdName,
		Aliases: []string{"po"},
		Short:   "VPC switch port management",
	},

	Setup: func(self *command.Command) error {
		subCommands := command.Commands{
			add.Cmd,
			connect.Cmd,
			disconnect.Cmd,
			//list.Cmd,
			remove.Cmd,
			set.Cmd,
			uplink.Cmd,
		}

		if err := self.Register(subCommands); err != nil {
			return errors.Wrapf(err, "unable to register sub-commands under %s", cmdName)
		}

		return nil
	},
}
