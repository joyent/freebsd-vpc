package destroy

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/conswriter"
	"github.com/sean-/vpc/internal/command"
	"github.com/sean-/vpc/internal/command/flag"
	"github.com/sean-/vpc/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.freebsd.org/sys/vpc/vmnic"
)

const (
	_CmdName    = "destroy"
	_KeyVMNICID = config.KeyVMNICDestroyID
)

var Cmd = &command.Command{
	Name: _CmdName,
	Cobra: &cobra.Command{
		Use:              _CmdName,
		Aliases:          []string{"rm", "del", "delete"},
		TraverseChildren: true,
		Short:            "destroy a VM NIC",
		SilenceUsage:     true,
		Args:             cobra.NoArgs,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: runE,
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddVMNICID(self, _KeyVMNICID, true); err != nil {
			return errors.Wrap(err, "unable to register VM NIC ID flag on VPC Switch destroy")
		}

		return nil
	},
}

func runE(cmd *cobra.Command, args []string) error {
	cons := conswriter.GetTerminal()

	cons.Write([]byte(fmt.Sprintf("Destroying VM NIC...")))

	id, err := flag.GetID(viper.GetViper(), _KeyVMNICID)
	if err != nil {
		return errors.Wrap(err, "unable to get VPC ID")
	}

	vmnicCfg := vmnic.Config{
		ID:        id,
		Writeable: true,
	}

	// TODO(seanc@): Go back and add vmnic/vpcsw to other commands
	log.Info().Object("cfg", vmnicCfg).Str("op", "destroy").Msg("vpc_ctl")

	vmNIC, err := vmnic.Open(vmnicCfg)
	if err != nil {
		return errors.Wrap(err, "unable to open VM NIC")
	}
	defer vmNIC.Close()

	if err := vmNIC.Destroy(); err != nil {
		return errors.Wrap(err, "unable to destroy VM NIC")
	}

	cons.Write([]byte("done.\n"))

	return nil
}
