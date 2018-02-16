package create

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sean-/conswriter"
	"github.com/sean-/vpc/cmd/vpcsw/flag"
	"github.com/sean-/vpc/internal/command"
	"github.com/sean-/vpc/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.freebsd.org/sys/vpc"
	"go.freebsd.org/sys/vpc/vpcsw"
)

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Use:          "create",
		Short:        "create a VPC switch",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cons := conswriter.GetTerminal()

			cons.Write([]byte(fmt.Sprintf("Creating VPC Switch...")))
			id := vpc.GenID()

			switchCfg := vpcsw.Config{
				ID:  id,
				VNI: vpc.VNI(viper.GetInt(config.KeySWVNI)),
			}

			vpcSwitch, err := vpcsw.New(switchCfg)
			if err != nil {
				return errors.Wrap(err, "invalid VPC Switch configuration")
			}
			defer vpcSwitch.Close()

			if err := vpcSwitch.Create(); err != nil {
				return errors.Wrap(err, "unable to create VPC Switch")
			}

			cons.Write([]byte("done.\n"))

			return nil
		},
	},
	Setup: func(parent *command.Command) error {
		if err := flag.AddVNI(parent, flag.VNICfg{Required: true}); err != nil {
			return errors.Wrap(err, "unable to register VNI flag")
		}

		return nil
	},
}
