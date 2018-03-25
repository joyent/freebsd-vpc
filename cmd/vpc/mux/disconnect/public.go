package disconnect

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
	_CmdName  = "disconnect"
	_KeyMuxID = config.KeyMuxDisconnectMuxID
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "disconnect a VPC Mux from its VPC Interface",
		Aliases:      []string{"disco"},
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cons := conswriter.GetTerminal()

			cons.Write([]byte(fmt.Sprintf("Disconnecting VPC Mux from VPC Interface...")))

			muxID, err := flag.GetMuxID(viper.GetViper(), _KeyMuxID)
			if err != nil {
				return errors.Wrap(err, "unable to get switch port ID")
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

			if err = vpcMux.Disconnect(); err != nil {
				log.Error().Err(err).Object("mux-id", muxID).Msg("vpc mux disconnect failed")
				return errors.Wrap(err, "unable to disconnect a VPC Mux from its VPC Interface")
			}

			cons.Write([]byte("done.\n"))

			log.Info().Object("mux-id", muxID).Msg("VPC Mux disconnected from VPC Interface")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddMuxID(self, _KeyMuxID, true); err != nil {
			return errors.Wrap(err, "unable to register Mux ID flag on VPC Mux disconnect")
		}

		return nil
	},
}
