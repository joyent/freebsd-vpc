package set

import (
	"github.com/pkg/errors"
	"github.com/sean-/vpc/internal/command"
	"github.com/sean-/vpc/internal/command/flag"
	"github.com/sean-/vpc/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.freebsd.org/sys/vpc/vmnic"
)

const (
	_CmdName       = "set"
	_KeySetNQueues = config.KeyVMNICSetNQueues
	_KeyVMNICID    = config.KeyVMNICSetVMNICID
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "set VM NIC information",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := flag.GetID(viper.GetViper(), _KeyVMNICID)
			if err != nil {
				return errors.Wrap(err, "unable to get VM NIC ID")
			}

			vmnicCfg := vmnic.Config{
				ID: id,
			}
			vmn, err := vmnic.Open(vmnicCfg)
			if err != nil {
				return errors.Wrap(err, "unable to open VM NIC")
			}
			defer vmn.Close()

			if numQueues := viper.GetInt(_KeySetNQueues); numQueues > 0 {
				if err := vmn.NQueuesSet(uint16(numQueues)); err != nil {
					return errors.Wrapf(err, "unable to set the number of hardware queues")
				}
			}

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddVMNICID(self, _KeyVMNICID, true); err != nil {
			return errors.Wrap(err, "unable to register VM NIC ID flag on VM NIC set")
		}

		{
			const (
				key          = _KeySetNQueues
				longName     = "num-queues"
				shortName    = "n"
				defaultValue = 0
				description  = "set the number of hardware queues for a given VM NIC"
			)

			flags := self.Cobra.Flags()
			flags.IntP(longName, shortName, defaultValue, description)

			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		return nil
	},
}
