package create

import (
	"fmt"
	"net"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/hostif"
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
	cmdName      = "create"
	_KeyHostifID = config.KeyHostifCreateID
)

var Cmd = &command.Command{
	Name: cmdName,

	Cobra: &cobra.Command{
		Use:          cmdName,
		Short:        "create a Hostif interface",
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

			cons.Write([]byte(fmt.Sprintf("Creating hostif interface...")))

			id, err := flag.GetID(viper.GetViper(), _KeyHostifID)
			if err != nil {
				return errors.Wrap(err, "unable to get VPC Hostif ID")
			}

			hostifCfg := hostif.Config{
				ID: id,
			}

			hostifNIC, err := hostif.Create(hostifCfg)
			if err != nil {
				log.Error().Err(err).Object("hostif-id", id).Msg("hostif create failed")
				return errors.Wrap(err, "unable to create Hostif NIC")
			}
			defer hostifNIC.Close()

			if err := hostifNIC.Commit(); err != nil {
				log.Error().Err(err).Object("hostif-id", id).Msg("Hostif NIC commit failed")
				return errors.Wrap(err, "unable to commit Hostif NIC")
			}

			cons.Write([]byte("done.\n"))

			var newHostifNIC net.Interface
			{ // Get the before/after
				ifacesAfterCreate, err := vpctest.GetAllInterfaces()
				if err != nil {
					return errors.Wrapf(err, "unable to get all interfaces")
				}
				_, newIfaces, _ := existingIfaces.Difference(ifacesAfterCreate)

				var newHostifMAC net.HardwareAddr = id.Node[:]
				newHostifNIC, err = newIfaces.FindMAC(newHostifMAC)
				if err != nil {
					return errors.Wrapf(err, "unable to find new Hostif NIC with MAC %q", id.Node)
				}
			}

			log.Info().Object("hostif-id", id).Str("mac", newHostifNIC.HardwareAddr.String()).Str("name", newHostifNIC.Name).Msg("Hostif NIC created")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddHostifID(self, _KeyHostifID, false); err != nil {
			return errors.Wrap(err, "unable to register VPC Hostif ID flag on VPC Hostif create")
		}

		return nil
	},
}
