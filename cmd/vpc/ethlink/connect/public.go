package connect

import (
	"fmt"

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
	_CmdName          = "connect"
	_KeyEthLinkID     = config.KeyEthLinkConnectID
	_KeyEthLinkL2Name = config.KeyEthLinkConnectL2Name
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "connect a VPC EthLink interface to a physical or cloned interface",
		Long:         "`vpc ethlink connect` is used to create a VPC interface that wraps a cloned or physical interface.  The cloned or physical interface is typically the interface used for the underlay network to route between different hosts.",
		Aliases:      []string{"conn"},
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cons := conswriter.GetTerminal()

			cons.Write([]byte(fmt.Sprintf("Connecting VPC EthLink to a physical or cloned interface...")))

			ethLinkID, err := flag.GetID(viper.GetViper(), _KeyEthLinkID)
			if err != nil {
				return errors.Wrap(err, "unable to get EthLink's VPC ID")
			}

			l2Name := viper.GetString(_KeyEthLinkL2Name)
			{
				if l2Name == "" {
					return errors.Wrap(err, "unable to get EthLink's physical or cloned interface name")
				}

				existingIfaces, err := vpctest.GetAllInterfaces()
				if err != nil {
					return errors.Wrapf(err, "unable to get all interfaces")
				}

				var found bool
				for _, iface := range existingIfaces {
					if l2Name == iface.Name {
						found = true
						break
					}
				}

				if !found {
					return errors.Errorf("unable to find interface %q", l2Name)
				}

			}

			ethLinkCfg := ethlink.Config{
				ID:        ethLinkID,
				Writeable: true,
			}

			ethLinkNIC, err := ethlink.Open(ethLinkCfg)
			if err != nil {
				log.Error().Err(err).Object("ethlink-id", ethLinkCfg).Msg("VPC EthLink open failed")
				return errors.Wrap(err, "unable to open VPC EthLink")
			}
			defer ethLinkNIC.Close()

			if err = ethLinkNIC.Connect(l2Name); err != nil {
				log.Error().Err(err).Object("ethlink-id", ethLinkID).Str("l2-name", l2Name).Msg("vpc ethlink connect failed")
				return errors.Wrap(err, "unable to connect a VPC EthLink to physical or cloned interface")
			}

			cons.Write([]byte("done.\n"))

			log.Info().Object("ethlink-id", ethLinkID).Str("l2-name", l2Name).Msg("VPC EthLink connected to physical or cloned interface")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddEthLinkID(self, _KeyEthLinkID, false); err != nil {
			return errors.Wrap(err, "unable to register VPC EthLink ID flag on VPC EthLink connect")
		}

		{
			const (
				key          = _KeyEthLinkL2Name
				longName     = "l2-name"
				shortName    = "n"
				defaultValue = ""
				description  = "Name of the underlay L2 interface to be wrapped by VPC EthLink"
			)

			flags := self.Cobra.Flags()
			flags.StringP(longName, shortName, defaultValue, description)

			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		return nil
	},
}
