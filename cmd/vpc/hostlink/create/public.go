package create

import (
	"fmt"
	"net"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/hostlink"
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
	cmdName        = "create"
	_KeyHostlinkID = config.KeyHostlinkCreateID
)

var Cmd = &command.Command{
	Name: cmdName,

	Cobra: &cobra.Command{
		Use:          cmdName,
		Short:        "create a Hostlink interface",
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

			cons.Write([]byte(fmt.Sprintf("Creating hostlink interface...")))

			id, err := flag.GetID(viper.GetViper(), _KeyHostlinkID)
			if err != nil {
				return errors.Wrap(err, "unable to get VPC Hostlink ID")
			}

			hostlinkCfg := hostlink.Config{
				ID: id,
			}

			hostlinkNIC, err := hostlink.Create(hostlinkCfg)
			if err != nil {
				log.Error().Err(err).Object("hostlink-id", id).Msg("hostlink create failed")
				return errors.Wrap(err, "unable to create Hostlink NIC")
			}
			defer hostlinkNIC.Close()

			if err := hostlinkNIC.Commit(); err != nil {
				log.Error().Err(err).Object("hostlink-id", id).Msg("Hostlink NIC commit failed")
				return errors.Wrap(err, "unable to commit Hostlink NIC")
			}

			cons.Write([]byte("done.\n"))

			var newHostlinkNIC net.Interface
			{ // Get the before/after
				ifacesAfterCreate, err := vpctest.GetAllInterfaces()
				if err != nil {
					return errors.Wrapf(err, "unable to get all interfaces")
				}
				_, newIfaces, _ := existingIfaces.Difference(ifacesAfterCreate)

				var newHostlinkMAC net.HardwareAddr = id.Node[:]
				newHostlinkNIC, err = newIfaces.FindMAC(newHostlinkMAC)
				if err != nil {
					return errors.Wrapf(err, "unable to find new Hostlink NIC with MAC %q", id.Node)
				}
			}

			log.Info().Object("hostlink-id", id).Str("mac", newHostlinkNIC.HardwareAddr.String()).Str("name", newHostlinkNIC.Name).Msg("Hostlink NIC created")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddHostlinkID(self, _KeyHostlinkID, false); err != nil {
			return errors.Wrap(err, "unable to register VPC Hostlink ID flag on VPC Hostlink create")
		}

		return nil
	},
}
