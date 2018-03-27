package destroy

import (
	"fmt"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/hostif"
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
	cmdName         = "destroy"
	_KeyInterfaceID = config.KeyHostifDestroyID
)

var Cmd = &command.Command{
	Name: cmdName,
	Cobra: &cobra.Command{
		Use:              cmdName,
		Aliases:          []string{"rm", "del", "delete"},
		TraverseChildren: true,
		Short:            "destroy a Hostif NIC",
		SilenceUsage:     true,
		Args:             cobra.NoArgs,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: runE,
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddHostifID(self, _KeyInterfaceID, true); err != nil {
			return errors.Wrap(err, "unable to register VPC Hostif ID flag on VPC Hostif destroy")
		}

		return nil
	},
}

func runE(cmd *cobra.Command, args []string) error {
	cons := conswriter.GetTerminal()

	cons.Write([]byte(fmt.Sprintf("Destroying VPC Hostif...")))

	id, err := flag.GetID(viper.GetViper(), _KeyInterfaceID)
	if err != nil {
		return errors.Wrap(err, "unable to get Hostif VPC ID")
	}

	hostifCfg := hostif.Config{
		ID:        id,
		Writeable: true,
	}

	// TODO(seanc@): Go back and add hostif/vmnic/vpcsw to other commands
	log.Info().Object("cfg", hostifCfg).Str("op", "destroy").Msg("vpc_ctl")

	hostifNIC, err := hostif.Open(hostifCfg)
	if err != nil {
		return errors.Wrap(err, "unable to open VPC Hostif NIC")
	}
	defer hostifNIC.Close()

	if err := hostifNIC.Destroy(); err != nil {
		return errors.Wrap(err, "unable to destroy VPC Hostif NIC")
	}

	cons.Write([]byte("done.\n"))

	return nil
}
