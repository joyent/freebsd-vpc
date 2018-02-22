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
	"go.freebsd.org/sys/vpc"
	"go.freebsd.org/sys/vpc/vpcsw"
	"go.freebsd.org/sys/vpc/vpctest"
)

const (
	_CmdName = "create"
	_KeyID   = config.KeySWCreateID
	_KeyMAC  = config.KeySWCreateMAC
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "create a VPC switch",
		SilenceUsage: true,
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

			cons.Write([]byte(fmt.Sprintf("Creating VPC Switch...")))

			id, err := flag.GetID(viper.GetViper(), _KeyID)
			if err != nil {
				return errors.Wrap(err, "unable to get VPC ID")
			}

			mac, err := flag.GetMAC(viper.GetViper(), _KeyMAC, &id)
			if err != nil {
				return errors.Wrap(err, "unable to get MAC address")
			}

			switchCfg := vpcsw.Config{
				ID:  id,
				MAC: mac,
				VNI: vpc.VNI(viper.GetInt(config.KeySWCreateVNI)),
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
		if err := flag.AddID(self, _KeyID, false); err != nil {
			return errors.Wrap(err, "unable to register ID flag on VPC Switch create")
		}

		if err := flag.AddMAC(self, _KeyMAC, false); err != nil {
			return errors.Wrap(err, "unable to register MAC flag on VPC Switch create")
		}

		if err := flag.AddVNI(self, flag.VNICfg{Name: config.KeySWCreateVNI, Required: true}); err != nil {
			return errors.Wrap(err, "unable to register VNI flag on VPC Switch create")
		}

		return nil
	},
}
