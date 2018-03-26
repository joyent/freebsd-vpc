package destroy

import (
	"fmt"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/hostlink"
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
	_KeyInterfaceID = config.KeyHostlinkDestroyID
)

var Cmd = &command.Command{
	Name: cmdName,
	Cobra: &cobra.Command{
		Use:              cmdName,
		Aliases:          []string{"rm", "del", "delete"},
		TraverseChildren: true,
		Short:            "destroy a Hostlink NIC",
		SilenceUsage:     true,
		Args:             cobra.NoArgs,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: runE,
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddHostlinkID(self, _KeyInterfaceID, true); err != nil {
			return errors.Wrap(err, "unable to register VPC Hostlink ID flag on VPC Hostlink destroy")
		}

		return nil
	},
}

func runE(cmd *cobra.Command, args []string) error {
	cons := conswriter.GetTerminal()

	cons.Write([]byte(fmt.Sprintf("Destroying VPC Hostlink...")))

	id, err := flag.GetID(viper.GetViper(), _KeyInterfaceID)
	if err != nil {
		return errors.Wrap(err, "unable to get Hostlink VPC ID")
	}

	hostlinkCfg := hostlink.Config{
		ID:        id,
		Writeable: true,
	}

	// TODO(seanc@): Go back and add hostlink/vmnic/vpcsw to other commands
	log.Info().Object("cfg", hostlinkCfg).Str("op", "destroy").Msg("vpc_ctl")

	hostlinkNIC, err := hostlink.Open(hostlinkCfg)
	if err != nil {
		return errors.Wrap(err, "unable to open VPC Hostlink NIC")
	}
	defer hostlinkNIC.Close()

	if err := hostlinkNIC.Destroy(); err != nil {
		return errors.Wrap(err, "unable to destroy VPC Hostlink NIC")
	}

	cons.Write([]byte("done.\n"))

	return nil
}
