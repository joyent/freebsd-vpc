package uplink

import (
	"fmt"
	"net"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/vpcsw"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/joyent/freebsd-vpc/internal/command/flag"
	"github.com/joyent/freebsd-vpc/internal/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/conswriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	_CmdName     = "uplink"
	_KeyPortID   = config.KeySWPortUplinkPortID
	_KeySwitchID = config.KeySWPortUplinkSwitchID
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "Sets a given switch port as the uplink",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cons := conswriter.GetTerminal()

			cons.Write([]byte(fmt.Sprintf("Designating VPC Switch Port as Uplink Port...")))

			portID, err := flag.GetPortID(viper.GetViper(), _KeyPortID)
			if err != nil {
				return errors.Wrap(err, "unable to get switch port ID")
			}

			switchID, err := flag.GetSwitchID(viper.GetViper(), _KeySwitchID)
			if err != nil {
				return errors.Wrap(err, "unable to get switch ID")
			}

			switchCfg := vpcsw.Config{
				ID:        switchID,
				Writeable: true,
			}

			vpcSwitch, err := vpcsw.Open(switchCfg)
			if err != nil {
				log.Error().Err(err).Object("switch-id", switchID).Msg("VPC Switch open failed")
				return errors.Wrap(err, "unable to open VPC Switch")
			}
			defer vpcSwitch.Close()

			var portMAC net.HardwareAddr = portID.Node[:]
			if err := vpcSwitch.PortUplinkSet(portID, portMAC); err != nil {
				log.Error().Err(err).Object("port-id", portID).Object("switch-cfg", switchCfg).Msg("failed to set VPC Switch Port as an Uplink port")
				return errors.Wrap(err, "unable to create a VPC Switch Port uplink")
			}

			cons.Write([]byte("done.\n"))

			log.Info().Object("switch-id", switchID).Object("port-id", portID).Msg("Uplink port for VPC Switch set")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddPortID(self, _KeyPortID, true); err != nil {
			return errors.Wrap(err, "unable to register Port ID flag on VPC Switch Port uplink")
		}

		if err := flag.AddSwitchID(self, _KeySwitchID, true); err != nil {
			return errors.Wrap(err, "unable to register Switch ID flag on VPC Switch Port uplink")
		}

		return nil
	},
}
