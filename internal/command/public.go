package command

import (
	"github.com/spf13/cobra"
)

type SetupFunc func(parent *Command) error

type Command struct {
	Setup      SetupFunc
	ArgAliases []string
	ValidArgs  []string
	Cobra      *cobra.Command
}

type Commands []*Command

func (cs Commands) ArgAliases() []string {
	args := make([]string, 0, len(cs))
	for _, c := range cs {
		args = append(args, c.ArgAliases...)
	}
	return args
}

func (cs Commands) ValidArgs() []string {
	args := make([]string, 0, len(cs))
	for _, c := range cs {
		args = append(args, c.ValidArgs...)
	}
	return args
}
