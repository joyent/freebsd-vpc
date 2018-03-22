package vm

import (
	"github.com/joyent/freebsd-vpc/cmd/vpc/vm/create"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const cmdName = "vm"

var Cmd = &command.Command{
	Name: cmdName,
	Cobra: &cobra.Command{
		Use:     cmdName,
		Short:   "Interaction with the VM agent",
	},

	Setup: func(self *command.Command) error {
		subCommands := []*command.Command{
			create.Cmd,
		}

		if err := self.Register(subCommands); err != nil {
			log.Fatal().Err(err).Str("cmd", cmdName).Msg("unable to register sub-commands")
		}

		return nil
	},
}
