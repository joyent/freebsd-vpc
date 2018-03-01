package connect

import (
	"fmt"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/vpcp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/conswriter"
	"github.com/sean-/vpc/internal/command"
	"github.com/sean-/vpc/internal/command/flag"
	"github.com/sean-/vpc/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	_CmdName        = "connect"
	_KeyPortID      = config.KeySWPortConnectPortID
	_KeyInterfaceID = config.KeySWPortConnectInterfaceID
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "connect a VPC Interface to a VPC Switch Port",
		Aliases:      []string{"conn"},
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cons := conswriter.GetTerminal()

			cons.Write([]byte(fmt.Sprintf("Connecting VPC Interface to VPC Switch Port...")))

			interfaceID, err := flag.GetID(viper.GetViper(), _KeyInterfaceID)
			if err != nil {
				return errors.Wrap(err, "unable to get switch port ID")
			}

			portID, err := flag.GetPortID(viper.GetViper(), _KeyPortID)
			if err != nil {
				return errors.Wrap(err, "unable to get switch port ID")
			}

			portCfg := vpcp.Config{
				ID:        portID,
				Writeable: true,
			}

			vpcPort, err := vpcp.Open(portCfg)
			if err != nil {
				log.Error().Err(err).Str("port-id", portID.String()).Msg("VPC Switch Port open failed")
				return errors.Wrap(err, "unable to open VPC Switch Port")
			}
			defer vpcPort.Close()

			if err = vpcPort.Connect(interfaceID); err != nil {
				log.Error().Err(err).Object("port-id", portID).Object("interface-id", interfaceID).Msg("vpc switch port connect failed")
				return errors.Wrap(err, "unable to connect a VPC Interface to VPC Switch Port")
			}

			cons.Write([]byte("done.\n"))

			log.Info().Object("port-id", portID).Object("interface-id", interfaceID).Msg("VPC Interface connected to VPC port")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddInterfaceID(self, _KeyInterfaceID, true); err != nil {
			return errors.Wrap(err, "unable to register Interface ID flag on VPC Switch Port connect")
		}

		if err := flag.AddPortID(self, _KeyPortID, true); err != nil {
			return errors.Wrap(err, "unable to register Port ID flag on VPC Switch Port connect")
		}

		return nil
	},
}
