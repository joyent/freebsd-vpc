package create

import (
	"fmt"
	"net"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/vpcsw"
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/vpctest"
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
	_CmdName      = "create"
	_KeySwitchID  = config.KeySWCreateSwitchID
	_KeySwitchMAC = config.KeySWCreateSwitchMAC
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "create a VPC switch",
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

			cons.Write([]byte(fmt.Sprintf("Creating VPC Switch...")))

			id, err := flag.GetID(viper.GetViper(), _KeySwitchID)
			if err != nil {
				return errors.Wrap(err, "unable to get VPC ID")
			}

			mac, err := flag.GetMAC(viper.GetViper(), _KeySwitchMAC, &id)
			if err != nil {
				return errors.Wrap(err, "unable to get MAC address")
			}

			switchCfg := vpcsw.Config{
				ID:  id,
				MAC: mac,
			}

			vpcSwitch, err := vpcsw.Create(switchCfg)
			if err != nil {
				log.Error().Err(err).Str("id", id.String()).Msg("vpcsw create failed")
				return errors.Wrap(err, "unable to create VPC Switch")
			}
			defer vpcSwitch.Close()

			if err := vpcSwitch.Commit(); err != nil {
				log.Error().Err(err).Str("id", id.String()).Msg("vpcsw commit failed")
				return errors.Wrap(err, "unable to commit VPC Switch")
			}

			cons.Write([]byte("done.\n"))

			var newSwitch net.Interface
			{ // Get the before/after
				ifacesAfterCreate, err := vpctest.GetAllInterfaces()
				if err != nil {
					return errors.Wrapf(err, "unable to get all interfaces")
				}
				_, newIfaces, _ := existingIfaces.Difference(ifacesAfterCreate)

				var newSwitchMAC net.HardwareAddr = id.Node[:]
				newSwitch, err = newIfaces.FindMAC(newSwitchMAC)
				if err != nil {
					return errors.Wrapf(err, "unable to find new VPC Switch with MAC %q", id.Node)
				}
			}

			log.Info().Str("id", id.String()).Str("mac", newSwitch.HardwareAddr.String()).Str("name", newSwitch.Name).Msg("vpcsw created")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddSwitchID(self, _KeySwitchID, false); err != nil {
			return errors.Wrap(err, "unable to register ID flag on VPC Switch create")
		}

		if err := flag.AddMAC(self, _KeySwitchMAC, false); err != nil {
			return errors.Wrap(err, "unable to register MAC flag on VPC Switch create")
		}

		return nil
	},
}
