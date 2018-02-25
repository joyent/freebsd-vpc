package attach

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
	_CmdName   = "attach"
	_KeyPortID = config.KeySWPortAttachID
	_KeyL2Name = config.KeySWPortAttachL2Name
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "add an L2 link to a VPC Switch",
		Aliases:      []string{"connect"},
		SilenceUsage: true,
		Args:         cobra.NoArgs,
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

			portID, err := flag.GetPortID(viper.GetViper(), _KeyPortID)
			if err != nil {
				return errors.Wrap(err, "unable to get port ID")
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
				Uplink: viper.GetBool(_KeyUplink),
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

		{
			const (
				key          = _KeyL2Name
				longName     = "l2-name"
				shortName    = "n"
				defaultValue = ""
				description  = "Name of the Layer-2 link Tag the port as an uplink port in the VPC Switch"
			)

			flags := self.Cobra.Flags()
			flags.BoolP(longName, shortName, defaultValue, description)

			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		return nil
	},
}
