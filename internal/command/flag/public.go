package flag

import (
	"net"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// AddEthLinkID adds the VM NIC ID flag to a given command.
func AddEthLinkID(cmd *command.Command, keyName string, required bool) error {
	key := keyName
	const (
		longName     = "ethlink-id"
		shortName    = "E"
		defaultValue = ""
		description  = "Specify the EthLink ID"
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

// AddHostifID adds the Hostif ID to a given command.
func AddHostifID(cmd *command.Command, keyName string, required bool) error {
	key := keyName
	const (
		longName     = "hostif-id"
		shortName    = "H"
		defaultValue = ""
		description  = "Specify the VPC Hostif ID"
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

// AddInterfaceID adds the Interface ID to a given command.
func AddInterfaceID(cmd *command.Command, keyName string, required bool) error {
	key := keyName
	const (
		longName     = "interface-id"
		shortName    = "I"
		defaultValue = ""
		description  = "Specify the VPC Interface ID"
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

// AddMuxID adds the Mux ID flag to a given command.
func AddMuxID(cmd *command.Command, keyName string, required bool) error {
	key := keyName
	const (
		longName     = "mux-id"
		shortName    = "M"
		defaultValue = ""
		description  = "Specify the VPC Mux ID"
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

// AddPortID adds the Port ID flag to a given command.
func AddPortID(cmd *command.Command, keyName string, required bool) error {
	key := keyName
	const (
		longName     = "port-id"
		shortName    = ""
		defaultValue = ""
		description  = "Specify the VPC Port ID"
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

// AddSwitchID adds the Switch ID flag to a given command.
func AddSwitchID(cmd *command.Command, keyName string, required bool) error {
	key := keyName
	const (
		longName     = "switch-id"
		shortName    = ""
		defaultValue = ""
		description  = "Specify the VPC Switch ID"
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

// AddVMNICID adds the VM NIC ID flag to a given command.
func AddVMNICID(cmd *command.Command, keyName string, required bool) error {
	key := keyName
	const (
		longName     = "vmnic-id"
		shortName    = "N"
		defaultValue = ""
		description  = "Specify the VM NIC ID"
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

// GetID returns the VPC ID address found in the Viper key.
func GetID(v *viper.Viper, key string) (id vpc.ID, err error) {
	if idStr := v.GetString(key); idStr != "" {
		if id, err = vpc.ParseID(idStr); err != nil {
			return vpc.ID{}, errors.Wrapf(err, "unable to parse VPC ID %q", idStr)
		}

		return id, nil
	}

	return vpc.ID{}, errors.Wrapf(err, "unable to lookup VPC ID %q", key)
}

// GetMAC returns the MAC address found in the Viper key.  If id is not nil,
// GetMAC falls back to id.Node for the default value.  If GetMAC uses id.Node,
// it sets the viper key accordingly for future callers.  If id is nil and no
// MAC is found, an error is returned.
func GetMAC(v *viper.Viper, key string, id *vpc.ID) (mac net.HardwareAddr, err error) {
	switch macStr := v.GetString(key); macStr {
	case "":
		if id == nil {
			return net.HardwareAddr{}, errors.Wrapf(err, "missing MAC address")
		}

		mac = id.Node[:]
		v.Set(key, mac.String())
	default:
		if mac, err = net.ParseMAC(macStr); err != nil {
			return net.HardwareAddr{}, errors.Wrapf(err, "unable to parse MAC %q", macStr)
		}
	}

	return mac, nil
}

// GetMuxID returns the VPC Mux ID found in the Viper key.
func GetMuxID(v *viper.Viper, key string) (id vpc.ID, err error) {
	muxIDStr := v.GetString(key)
	if muxIDStr == "" {
		return vpc.ID{}, errors.Wrap(err, "missing VPC Mux ID")
	}

	if id, err = vpc.ParseID(muxIDStr); err != nil {
		return vpc.ID{}, errors.Wrapf(err, "unable to parse VPC Mux ID %q", muxIDStr)
	}

	return id, nil
}

// GetPortID returns the VPC ID found in the Viper key.
func GetPortID(v *viper.Viper, key string) (id vpc.ID, err error) {
	portIDStr := v.GetString(key)
	if portIDStr == "" {
		return vpc.ID{}, errors.Wrap(err, "missing VPC Port ID")
	}

	if id, err = vpc.ParseID(portIDStr); err != nil {
		return vpc.ID{}, errors.Wrapf(err, "unable to parse VPC Port ID %q", portIDStr)
	}

	return id, nil
}

// GetSwitchID returns the VPC ID found in the Viper key.
func GetSwitchID(v *viper.Viper, key string) (id vpc.ID, err error) {
	switchIDStr := v.GetString(key)
	if switchIDStr == "" {
		return vpc.ID{}, errors.Wrap(err, "missing VPC Switch ID")
	}

	if id, err = vpc.ParseID(switchIDStr); err != nil {
		return vpc.ID{}, errors.Wrapf(err, "unable to parse VPC ID %q", switchIDStr)
	}

	return id, nil
}

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
