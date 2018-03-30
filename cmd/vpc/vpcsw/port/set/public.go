package set

import (
	"fmt"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc"
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/vpcp"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/joyent/freebsd-vpc/internal/command/flag"
	"github.com/joyent/freebsd-vpc/internal/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	cmdName    = "set"
	_KeyPortID = config.KeySWPortSetPortID
	_KeySetVNI = config.KeySWPortSetVNI
)

var Cmd = &command.Command{
	Name: cmdName,

	Cobra: &cobra.Command{
		Use:          cmdName,
		Short:        "set VPC Port Information",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			portID, err := flag.GetID(viper.GetViper(), _KeyPortID)
			if err != nil {
				return errors.Wrap(err, "unable to get Port VPC ID")
			}

			portCfg := vpcp.Config{
				ID: portID,
			}
			port, err := vpcp.Open(portCfg)
			if err != nil {
				return errors.Wrap(err, "unable to open VPC Port")
			}
			defer port.Close()

			if vni := viper.GetInt(_KeySetVNI); vni >= 0 {
				if err := port.SetVNI(vpc.VNI(vni)); err != nil {
					return errors.Wrapf(err, "unable to set VPC VNI")
				}
			}

			vni, err := port.GetVNI()
			if err != nil {
				return errors.Wrapf(err, "unable to get the VNI for VPC Port")
			}
			fmt.Printf("VNI: %d\n", vni)

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddPortID(self, _KeyPortID, true); err != nil {
			return errors.Wrap(err, "unable to register VPC PortID flag on VPC Switch Port Set")
		}

		{
			const (
				key          = _KeySetVNI
				longName     = "vni"
				shortName    = "n"
				defaultValue = -1
				description  = "set the VNI of a given VPC Port"
			)

			flags := self.Cobra.Flags()
			flags.IntP(longName, shortName, defaultValue, description)

			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		return nil
	},
}
