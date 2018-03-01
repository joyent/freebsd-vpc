package ethlink

import (
	"github.com/pkg/errors"
	"github.com/sean-/vpc/cli/ethlink/destroy"
	"github.com/sean-/vpc/cli/ethlink/list"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/cobra"
)

const _CmdName = "ethlink"

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use: _CmdName,
		// NOTE(seanc@): Funny story: ethlink was called l2link but needed to be
		// renmaed because "bad things happen" when a unit name includes a number.
		// Leave l2link as a historical artifact and small easter egg.
		Aliases: []string{"ethlink", "l2link", "phys"},
		Short:   "VPC EthLink management",
	},

	Setup: func(self *command.Command) error {
		subCommands := command.Commands{
			destroy.Cmd,
			list.Cmd,
		}

		if err := self.Register(subCommands); err != nil {
			return errors.Wrapf(err, "unable to register sub-commands under %s", _CmdName)
		}

		return nil
	},
}
