package create

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/conswriter"
	"github.com/sean-/vpc/internal/command"
	"github.com/sean-/vpc/internal/command/flag"
	"github.com/sean-/vpc/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.freebsd.org/sys/vpc/vmnic"
	"go.freebsd.org/sys/vpc/vpctest"
)

const (
	_CmdName     = "create"
	_KeyVMNICID  = config.KeyVMNICCreateID
	_KeyVMNICMAC = config.KeyVMNICCreateMAC
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "create a VM NIC",
		SilenceUsage: true,
		// TraverseChildren: true,
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cons := conswriter.GetTerminal()

			existingIfaces, err := vpctest.GetAllInterfaces()
			if err != nil {
				return errors.Wrapf(err, "unable to get all interfaces")
			}

			cons.Write([]byte(fmt.Sprintf("Creating VM NIC...")))

			id, err := flag.GetID(viper.GetViper(), _KeyVMNICID)
			if err != nil {
				return errors.Wrap(err, "unable to get VPC ID")
			}

			mac, err := flag.GetMAC(viper.GetViper(), _KeyVMNICMAC, &id)
			if err != nil {
				return errors.Wrap(err, "unable to get MAC address")
			}

			vmnicCfg := vmnic.Config{
				ID:  id,
				MAC: mac,
			}

			vmNIC, err := vmnic.Create(vmnicCfg)
			if err != nil {
				log.Error().Err(err).Object("vmnic-id", id).Msg("vmnic create failed")
				return errors.Wrap(err, "unable to create VM NIC")
			}
			defer vmNIC.Close()

			if err := vmNIC.Commit(); err != nil {
				log.Error().Err(err).Object("vmnic-id", id).Msg("VM NIC commit failed")
				return errors.Wrap(err, "unable to commit VM NIC")
			}

			cons.Write([]byte("done.\n"))

			var newVMNIC net.Interface
			{ // Get the before/after
				ifacesAfterCreate, err := vpctest.GetAllInterfaces()
				if err != nil {
					return errors.Wrapf(err, "unable to get all interfaces")
				}
				_, newIfaces, _ := existingIfaces.Difference(ifacesAfterCreate)

				var newVMNICMAC net.HardwareAddr = id.Node[:]
				newVMNIC, err = newIfaces.FindMAC(newVMNICMAC)
				if err != nil {
					return errors.Wrapf(err, "unable to find new VM NIC with MAC %q", id.Node)
				}
			}

			log.Info().Object("vmnic-id", id).Str("mac", newVMNIC.HardwareAddr.String()).Str("name", newVMNIC.Name).Msg("VM NIC created")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddVMNICID(self, _KeyVMNICID, false); err != nil {
			return errors.Wrap(err, "unable to register VM NIC ID flag on VM NIC create")
		}

		if err := flag.AddMAC(self, _KeyVMNICMAC, false); err != nil {
			return errors.Wrap(err, "unable to register MAC flag on VM NIC create")
		}

		return nil
	},
}
