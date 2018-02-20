package flag

import (
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/viper"
)

// AddID adds the ID flag to a given command.
func AddID(cmd *command.Command, keyName string, required bool) error {
	key := keyName
	const (
		longName     = "id"
		shortName    = "I"
		defaultValue = ""
		description  = "Specify the ID for a VPC Switch operation"
	)

	flags := cmd.Cobra.Flags()
	flags.StringP(longName, shortName, defaultValue, description)
	if required {
		cmd.Cobra.MarkFlagRequired(longName)
	}

	viper.BindPFlag(key, flags.Lookup(longName))
	viper.SetDefault(key, defaultValue)

	return nil
}

type VNICfg struct {
	Name     string
	Required bool
}

// AddFlagVNI adds the VNI flag to a given command.
func AddVNI(cmd *command.Command, cfg VNICfg) error {
	key := cfg.Name
	const (
		longName     = "vni"
		shortName    = ""
		defaultValue = 0
		description  = "Specify the VNI for a VPC Switch operation"
	)

	flags := cmd.Cobra.Flags()
	flags.UintP(longName, shortName, defaultValue, description)
	if cfg.Required {
		cmd.Cobra.MarkFlagRequired(longName)
	}

	viper.BindPFlag(key, flags.Lookup(longName))
	viper.SetDefault(key, defaultValue)

	return nil
}
