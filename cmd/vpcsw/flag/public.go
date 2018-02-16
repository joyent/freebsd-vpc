package flag

import (
	"github.com/sean-/vpc/internal/command"
	"github.com/sean-/vpc/internal/config"
	"github.com/spf13/viper"
)

type VNICfg struct {
	Required bool
}

// AddFlagVNI adds the VNI flag to a given command.
func AddVNI(cmd *command.Command, cfg VNICfg) error {
	const (
		key          = config.KeySWVNI
		longName     = "vni"
		shortName    = ""
		defaultValue = 0
		description  = "Specify the VNI for a switch operation"
	)

	flags := cmd.Cobra.Flags()
	flags.UintP(longName, shortName, defaultValue, description)
	cmd.Cobra.MarkFlagRequired(longName)

	viper.BindPFlag(key, flags.Lookup(longName))
	viper.SetDefault(key, defaultValue)

	return nil
}
