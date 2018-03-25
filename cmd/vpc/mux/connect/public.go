package connect

import (
	"fmt"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/mux"
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
	_CmdName     = "connect"
	_KeyMuxID    = config.KeyMuxConnectMuxID
	_KeyTargetID = config.KeyMuxConnectInterfaceID
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "connect a VPC Mux to a VPC EthLink",
		Aliases:      []string{"conn", "con"},
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cons := conswriter.GetTerminal()

			cons.Write([]byte(fmt.Sprintf("Connecting VPC Mux to VPC ID...")))

			muxID, err := flag.GetMuxID(viper.GetViper(), _KeyMuxID)
			if err != nil {
				return errors.Wrap(err, "unable to get VPC Mux ID")
			}

			targetID, err := flag.GetID(viper.GetViper(), _KeyTargetID)
			if err != nil {
				return errors.Wrap(err, "unable to get VPC Interface ID")
			}

			muxCfg := mux.Config{
				ID:        muxID,
				Writeable: true,
			}

			vpcMux, err := mux.Open(muxCfg)
			if err != nil {
				log.Error().Err(err).Object("mux-id", muxID).Msg("VPC Mux open failed")
				return errors.Wrap(err, "unable to open VPC Mux")
			}
			defer vpcMux.Close()

			if err = vpcMux.Connect(targetID); err != nil {
				log.Error().Err(err).Object("mux-id", muxID).Object("interface-id", targetID).Msg("vpc mux connect failed")
				return errors.Wrap(err, "unable to connect a VPC Interface to VPC Mux")
			}

			cons.Write([]byte("done.\n"))

			log.Info().Object("mux-id", muxID).Object("interface-id", targetID).Msg("VPC Mux connected to VPC Interface")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddInterfaceID(self, _KeyTargetID, true); err != nil {
			return errors.Wrap(err, "unable to register Interface ID flag on VPC Mux connect")
		}

		if err := flag.AddMuxID(self, _KeyMuxID, true); err != nil {
			return errors.Wrap(err, "unable to register Mux ID flag on VPC Mux connect")
		}

		return nil
	},
}
