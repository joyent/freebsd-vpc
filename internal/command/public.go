package command

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SetupFunc func(parent *Command) error
type SetViperFunc func(v *viper.Viper)

type Command struct {
	Name       string
	Setup      SetupFunc
	ArgAliases []string
	ValidArgs  []string
	Cobra      *cobra.Command
}

func (cmd *Command) Register(subCommands Commands) (err error) {
	for i, subCommand := range subCommands {
		// Validate sub-command's configuration
		switch {
		case subCommand.Name == "":
			return fmt.Errorf("%q[%d] has no name registered", cmd.Name, i)
		case strings.ContainsAny(subCommand.Name, "."):
			return errors.Errorf("%q[%d] contains an invalid character: %q", cmd.Name, i, subCommand.Name)
		case subCommand.Cobra == nil:
			return errors.Errorf("%q[%d].%q is missing a Cobra instance", cmd.Name, i, subCommand.Name)
		}

		if err := subCommand.Setup(subCommand); err != nil {
			return errors.Wrapf(err, "unable to register %q subcommands", subCommand.Name)
		}

		cmd.Cobra.AddCommand(subCommand.Cobra)
	}

	return nil
}

type Commands []*Command

func (cs Commands) ArgAliases() []string {
	args := make([]string, 0, len(cs))
	for _, cmd := range cs {
		args = append(args, cmd.ArgAliases...)
	}
	return args
}

func (cs Commands) ValidArgs() []string {
	args := make([]string, 0, len(cs))
	for _, cmd := range cs {
		args = append(args, cmd.ValidArgs...)
	}
	return args
}
