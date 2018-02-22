package add

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
	"go.freebsd.org/sys/vpc/vpcp"
	"go.freebsd.org/sys/vpc/vpcsw"
	"go.freebsd.org/sys/vpc/vpctest"
)

const (
	_CmdName     = "add"
	_KeyPortID   = config.KeySWPortAddID
	_KeyPortMAC  = config.KeySWPortAddMAC
	_KeySwitchID = config.KeySWPortAddSwitchID
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "add a port to a VPC Switch",
		Aliases:      []string{"create"},
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

			cons.Write([]byte(fmt.Sprintf("Adding port to VPC Switch...")))

			switchID, err := flag.GetSwitchID(viper.GetViper(), _KeySwitchID)
			if err != nil {
				return errors.Wrap(err, "unable to get VPC ID")
			}

			portID, err := flag.GetPortID(viper.GetViper(), _KeyPortID)
			if err != nil {
				return errors.Wrap(err, "unable to get VPC Switch Port ID")
			}

			mac, err := flag.GetMAC(viper.GetViper(), _KeyPortMAC, nil)
			if err != nil {
				return errors.Wrap(err, "unable to get MAC address")
			}

			switchCfg := vpcsw.Config{
				ID:        switchID,
				Writeable: true,
			}

			vpcSwitch, err := vpcsw.Open(switchCfg)
			if err != nil {
				log.Error().Err(err).Str("switch-id", switchID.String()).Msg("vpcsw open failed")
				return errors.Wrap(err, "unable to open VPC Switch")
			}
			defer vpcSwitch.Close()

			portAddCfg := vpcsw.Config{
				PortID: portID,
				MAC:    mac,
			}
			if err = vpcSwitch.PortAdd(portAddCfg); err != nil {
				log.Error().Err(err).Str("port-id", portAddCfg.PortID.String()).Msg("vpc switch port add failed")
				return errors.Wrap(err, "unable to add a port to VPC Switch")
			}

			portCfg := vpcp.Config{
				ID: portID,
			}
			swPort, err := vpcp.Open(portCfg)
			if err != nil {
				log.Error().Err(err).Str("port-id", portCfg.ID.String()).Msg("vpcp open failed")
				return errors.Wrap(err, "unable to open VPC Switch port")
			}
			defer swPort.Close()

			cons.Write([]byte("done.\n"))

			var newPort net.Interface
			{ // Get the before/after
				ifacesAfterCreate, err := vpctest.GetAllInterfaces()
				if err != nil {
					return errors.Wrapf(err, "unable to get all interfaces")
				}
				_, newIfaces, _ := existingIfaces.Difference(ifacesAfterCreate)

				var newPortMAC net.HardwareAddr = portID.Node[:]
				newPort, err = newIfaces.FindMAC(newPortMAC)
				if err != nil {
					return errors.Wrapf(err, "unable to find new VPC Port on Switch with MAC %q", portID.Node)
				}
			}

			log.Info().Str("port-id", portCfg.ID.String()).Str("switch-id", switchID.String()).Str("mac", newPort.HardwareAddr.String()).Str("name", newPort.Name).Msg("vpcp created")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddID(self, _KeyPortID, false); err != nil {
			return errors.Wrap(err, "unable to register ID flag on VPC Port add")
		}

		if err := flag.AddMAC(self, _KeyPortMAC, false); err != nil {
			return errors.Wrap(err, "unable to register MAC flag on VPC Port add")
		}

		if err := flag.AddSwitchID(self, _KeySwitchID, false); err != nil {
			return errors.Wrap(err, "unable to register Switch ID flag for VPC Port add")
		}

		return nil
	},
}
