package hostlink

import (
	"github.com/joyent/freebsd-vpc/cmd/vpc/hostlink/create"
	"github.com/joyent/freebsd-vpc/cmd/vpc/hostlink/destroy"
	"github.com/joyent/freebsd-vpc/cmd/vpc/hostlink/genmac"
	"github.com/joyent/freebsd-vpc/cmd/vpc/hostlink/list"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const cmdName = "hostlink"

var Cmd = &command.Command{
	Name: cmdName,

	Cobra: &cobra.Command{
		Use:     cmdName,
		Aliases: []string{"host"},
		Short:   "Host network stack interface",
	},

	Setup: func(self *command.Command) error {
		subCommands := command.Commands{
			create.Cmd,
			destroy.Cmd,
			genmac.Cmd,
			list.Cmd,
		}

		if err := self.Register(subCommands); err != nil {
			return errors.Wrapf(err, "unable to register sub-commands under %s", cmdName)
		}

		return nil
	},
}
