package doc

import (
	"github.com/joyent/freebsd-vpc/cmd/vpc/doc/man"
	"github.com/joyent/freebsd-vpc/cmd/vpc/doc/md"
	"github.com/joyent/freebsd-vpc/internal/buildtime"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const _CmdName = "doc"

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:     _CmdName,
		Aliases: []string{"docs", "documentation"},
		Short:   "Documentation for " + buildtime.PROGNAME,
	},

	Setup: func(self *command.Command) error {
		subCommands := command.Commands{
			man.Cmd,
			md.Cmd,
		}

		if err := self.Register(subCommands); err != nil {
			log.Fatal().Err(err).Str("cmd", _CmdName).Msg("unable to register sub-commands")
		}

		return nil
	},
}
