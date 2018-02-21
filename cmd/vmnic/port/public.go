package port

import (
	"github.com/pkg/errors"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

const _CmdName = "port"

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:     _CmdName,
		Aliases: []string{"sw"},
		Short:   "VPC switch management",
	},

	Setup: func(self *command.Command) error {
		subCommands := command.Commands{
		//create.Cmd,
		//destroy.Cmd,
		//list.Cmd,
		}

		if err := self.Register(subCommands); err != nil {
			return errors.Wrapf(err, "unable to register sub-commands under %s", _CmdName)
		}

		return nil
	},
}
