package destroy

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/conswriter"
	"github.com/sean-/vpc/cmd/vpcsw/flag"
	"github.com/sean-/vpc/internal/command"
	"github.com/sean-/vpc/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.freebsd.org/sys/vpc/vpcsw"
)

const (
	_CmdName = "destroy"
	_KeyID   = config.KeySWDestroyID
)

var Cmd = &command.Command{
	Name: _CmdName,
	Cobra: &cobra.Command{
		Use:              _CmdName,
		Aliases:          []string{"rm", "del", "delete"},
		TraverseChildren: true,
		Short:            "destroy a VPC switch",
		SilenceUsage:     true,
		Args:             cobra.NoArgs,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: runE,
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddID(self, _KeyID, true); err != nil {
			return errors.Wrap(err, "unable to register ID flag on VPC Switch destroy")
		}

		return nil
	},
}

func runE(cmd *cobra.Command, args []string) error {
	cons := conswriter.GetTerminal()

	cons.Write([]byte(fmt.Sprintf("Destroying VPC Switch...")))

	id, err := config.GetUUID(viper.GetViper(), _KeyID)
	if err != nil {
		return errors.Wrapf(err, "unable to parse UUID %q", viper.GetString(_KeyID))
	}

	switchCfg := vpcsw.Config{
		ID: id,
	}

	log.Info().Object("cfg", switchCfg).Str("op", "destroy").Msg("vpc_ctl")

	vpcSwitch, err := vpcsw.Open(switchCfg)
	if err != nil {
		return errors.Wrap(err, "unable to open VPC Switch")
	}
	defer vpcSwitch.Close()

	if err := vpcSwitch.Destroy(); err != nil {
		return errors.Wrap(err, "unable to destroy VPC Switch")
	}

	cons.Write([]byte("done.\n"))

	return nil
}
