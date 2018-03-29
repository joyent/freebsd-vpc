package create

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
	_CmdName  = "create"
	_KeyMuxID = config.KeyMuxCreateMuxID
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "Create a VPC mux interface",
		Aliases:      []string{"new"},
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			switch {
			case viper.GetString(_KeyMuxID) == "":
				// TODO(seanc@): convert mux-id to constants used by cobra when setting
				// the viper key.
				return errors.Errorf("mux-id is required")
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cons := conswriter.GetTerminal()

			cons.Write([]byte(fmt.Sprintf("Creating VPC Mux...")))

			muxID, err := flag.GetMuxID(viper.GetViper(), _KeyMuxID)
			if err != nil {
				return errors.Wrap(err, "unable to get VPC Mux ID")
			}

			// Create a stack of commit and undo operations to walk through in the
			// event of an error.
			var commit bool
			var commitFuncs, undoFuncs []func() error
			defer func() {
				scopeHandlers := undoFuncs
				modeStr := "undo"
				if commit {
					modeStr = "commit"
					scopeHandlers = commitFuncs
				}

				for i := len(scopeHandlers) - 1; i >= 0; i-- {
					if err := scopeHandlers[i](); err != nil {
						log.Error().Err(err).Msgf("failure during %s", modeStr)
					}
				}
			}()
			commitFuncs = append(commitFuncs, func() error {
				cons.Write([]byte("done.\n"))
				return nil
			})

			muxCfg := mux.Config{
				ID:        muxID,
				Writeable: true,
			}

			vpcMux, err := mux.Create(muxCfg)
			if err != nil {
				log.Error().Err(err).Object("mux-cfg", muxCfg).Msg("mux create failed")
				return errors.Wrap(err, "unable to create VPC Mux")
			}
			commitFuncs = append(commitFuncs, func() error {
				if err := vpcMux.Close(); err != nil {
					log.Error().Err(err).Msg("unable to commit VPC Mux")
					return errors.Wrap(err, "unable to commit VPC Mux during commit operation")
				}

				return nil
			})
			undoFuncs = append(undoFuncs, func() error {
				if err := vpcMux.Close(); err != nil {
					log.Error().Err(err).Msg("unable to clean up VPC Mux during error recovery")
				}

				return nil
			})

			commit = true

			log.Info().Object("mux-id", muxID).Msg("mux created")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddMuxID(self, _KeyMuxID, false); err != nil {
			return errors.Wrap(err, "unable to register Mux ID flag on VPC Mux create")
		}

		return nil
	},
}
