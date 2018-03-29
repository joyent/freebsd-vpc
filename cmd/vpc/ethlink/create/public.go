package create

import (
	"fmt"
	"net"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/ethlink"
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
	cmdName       = "create"
	_KeyEthLinkID = config.KeyEthLinkCreateID
)

var Cmd = &command.Command{
	Name: cmdName,

	Cobra: &cobra.Command{
		Use:          cmdName,
		Short:        "create an EthLink interface",
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

			cons.Write([]byte(fmt.Sprintf("Creating ethlink interface...")))

			id, err := flag.GetID(viper.GetViper(), _KeyEthLinkID)
			if err != nil {
				return errors.Wrap(err, "unable to get VPC EthLink ID")
			}

			ethlinkCfg := ethlink.Config{
				ID: id,
			}

			ethlinkNIC, err := ethlink.Create(ethlinkCfg)
			if err != nil {
				log.Error().Err(err).Object("ethlink-id", id).Msg("ethlink create failed")
				return errors.Wrap(err, "unable to create EthLink NIC")
			}
			defer ethlinkNIC.Close()

			if err := ethlinkNIC.Commit(); err != nil {
				log.Error().Err(err).Object("ethlink-id", id).Msg("EthLink NIC commit failed")
				return errors.Wrap(err, "unable to commit EthLink NIC")
			}

			cons.Write([]byte("done.\n"))

			var newEthLinkNIC net.Interface
			{ // Get the before/after
				ifacesAfterCreate, err := vpctest.GetAllInterfaces()
				if err != nil {
					return errors.Wrapf(err, "unable to get all interfaces")
				}
				_, newIfaces, _ := existingIfaces.Difference(ifacesAfterCreate)

				var newEthLinkMAC net.HardwareAddr = id.Node[:]
				newEthLinkNIC, err = newIfaces.FindMAC(newEthLinkMAC)
				if err != nil {
					return errors.Wrapf(err, "unable to find new EthLink NIC with MAC %q", id.Node)
				}
			}

			log.Info().Object("ethlink-id", id).Str("mac", newEthLinkNIC.HardwareAddr.String()).Str("name", newEthLinkNIC.Name).Msg("EthLink NIC created")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddEthLinkID(self, _KeyEthLinkID, false); err != nil {
			return errors.Wrap(err, "unable to register VPC EthLink ID flag on VPC EthLink create")
		}

		return nil
	},
}
