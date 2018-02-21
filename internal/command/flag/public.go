package flag

import (
	"net"

	"github.com/pkg/errors"
	"github.com/sean-/vpc/internal/command"
	"github.com/spf13/viper"
	"go.freebsd.org/sys/vpc"
)

// AddID adds the ID flag to a given command.
func AddID(cmd *command.Command, keyName string, required bool) error {
	key := keyName
	const (
		longName     = "id"
		shortName    = "I"
		defaultValue = ""
		description  = "Specify the ID"
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

// AddMAC adds the MAC flag to a given command.
func AddMAC(cmd *command.Command, keyName string, required bool) error {
	key := keyName
	const (
		longName     = "mac"
		shortName    = "m"
		defaultValue = ""
		description  = "Specify the MAC address"
	)

	flags := cmd.Cobra.Flags()
	flags.StringP(longName, shortName, defaultValue, description)
	if required {
		cmd.Cobra.MarkFlagRequired(longName)
	}

	flag := flags.Lookup(longName)
	flag.Hidden = true

	viper.BindPFlag(key, flag)
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
		description  = "Specify the VNI"
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

// GetID returns the VPC ID address found in the Viper key.  GetID falls back to
// generating a random ID if no argument was found.  If GetID sets the viper key
// accordingly for future callers if it needs to generate an ID.
func GetID(v *viper.Viper, key string) (id vpc.ID, err error) {
	switch idStr := v.GetString(key); idStr {
	case "":
		id = vpc.GenID()
		v.Set(key, id.String())
	default:
		if id, err = vpc.ParseID(idStr); err != nil {
			return vpc.ID{}, errors.Wrapf(err, "unable to parse UUID %q", idStr)
		}
	}

	return id, nil
}

// GetMAC returns the MAC address found in the Viper key.  GetMAC falls back to
// id.Node for the default value.  If GetMAC uses id.Node, it sets the viper key
// accordingly for future callers.
func GetMAC(v *viper.Viper, key string, id vpc.ID) (mac net.HardwareAddr, err error) {
	switch macStr := v.GetString(key); macStr {
	case "":
		mac = id.Node[:]
		v.Set(key, mac.String())
	default:
		if mac, err = net.ParseMAC(macStr); err != nil {
			return net.HardwareAddr{}, errors.Wrapf(err, "unable to parse MAC %q", macStr)
		}
	}

	return mac, nil
}
