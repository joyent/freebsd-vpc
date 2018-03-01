package destroy

import (
	"fmt"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/ethlink"
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
	_CmdName      = "destroy"
	_KeyEthLinkID = config.KeyEthLinkDestroyID
)

var Cmd = &command.Command{
	Name: _CmdName,
	Cobra: &cobra.Command{
		Use:              _CmdName,
		Aliases:          []string{"rm", "del", "delete"},
		TraverseChildren: true,
		Short:            "destroy a VPC EthLink",
		SilenceUsage:     true,
		Args:             cobra.NoArgs,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: runE,
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddEthLinkID(self, _KeyEthLinkID, true); err != nil {
			return errors.Wrap(err, "unable to register EthLink ID flag on EthLink destroy")
		}

		return nil
	},
}

func runE(cmd *cobra.Command, args []string) error {
	cons := conswriter.GetTerminal()

	cons.Write([]byte(fmt.Sprintf("Destroying VPC EthLink...")))

	ethLinkID, err := flag.GetID(viper.GetViper(), _KeyEthLinkID)
	if err != nil {
		return errors.Wrap(err, "unable to get EthLink VPC ID")
	}

	ethLinkCfg := ethlink.Config{
		ID:        ethLinkID,
		Writeable: true,
	}

	log.Info().Object("cfg", ethLinkCfg).Str("op", "destroy").Msg("vpc_ctl")

	ethLink, err := ethlink.Open(ethLinkCfg)
	if err != nil {
		return errors.Wrap(err, "unable to open VPC EthLink")
	}
	defer ethLink.Close()

	if err := ethLink.Destroy(); err != nil {
		return errors.Wrap(err, "unable to destroy VPC EthLink")
	}

	cons.Write([]byte("done.\n"))

	return nil
}
