package utilities

import (
	"fmt"
	"os"

	"github.com/sean-/conswriter"
	"github.com/sean-/smallz/buildtime"
	"github.com/sean-/vpc/internal/command"
	"github.com/sean-/vpc/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

var InstallManCommand = &command.Command{
	Cobra: &cobra.Command{
		Use:   "manpages",
		Short: "Generates and installs triton cli man pages",
		Long: `This command automatically generates up-to-date man pages of Triton CLI
command-line interface.  By default, it creates the man page files
in the "docs/man" directory under the current directory.`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cons := conswriter.GetTerminal()

			header := &doc.GenManHeader{
				Manual:  "VPC",
				Section: "3",
				Source:  fmt.Sprintf("VPC %s", buildtime.VERSION),
			}

			location := viper.GetString(config.KeyManPageDirectory)
			if location == "" {
				location = "./docs/man"
			}
			if _, err := os.Stat(location); os.IsNotExist(err) {
				os.Mkdir(location, 0777)
			}

			cmd.Root().DisableAutoGenTag = true
			cons.Write([]byte(fmt.Sprintf("Generating manpages to %s", location)))

			err := doc.GenManTree(cmd.Root(), header, location)
			if err != nil {
				return err
			}

			cons.Write([]byte("\nManpage generation complete"))

			return nil
		},
	},
	Setup: func(parent *command.Command) error {
		return nil
	},
}
