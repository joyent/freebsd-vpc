package bash

import (
	"os"

	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/internal/buildtime"
	"github.com/sean-/vpc/internal/command"
	"github.com/sean-/vpc/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Use:   "bash",
		Short: "Generates and install " + buildtime.PROGNAME + " bash autocompletion script",
		Long: `Generates a bash autocompletion script for ` + buildtime.PROGNAME + `

By default, the file is written directly to /etc/bash_completion.d
for convenience, and the command may need superuser rights, e.g.:

	$ sudo vpc shell autocomplete bash

Add ` + "`--bash-autocomplete-dir=/path/to/file`" + `. The default file name 
is ` + buildtime.PROGNAME + `.sh.

Logout and in again to reload the completion scripts,
or just source them in directly:

	$ . /bash_completion.d`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {

			target := viper.GetString(config.KeyBashAutoCompletionTarget)
			if _, err := os.Stat(target); os.IsNotExist(err) {
				if err := os.MkdirAll(target, 0777); err != nil {
					return errors.Wrapf(err, "unable to make bash-autocomplete-target %q", target)
				}
			}
			bashFile := fmt.Sprintf("%s/%s.sh", target, buildtime.PROGNAME)
			err := cmd.Root().GenBashCompletionFile(bashFile)
			if err != nil {
				return err
			}

			log.Info().Msg("Installation completed successfully.")

			return nil
		},
	},

	Setup: func(parent *command.Command) error {
		{
			const (
				key          = config.KeyBashAutoCompletionTarget
				longName     = "bash-autocomplete-dir"
				defaultValue = "/etc/bash_completion.d"
				description  = "autocompletion directory"
			)

			flags := parent.Cobra.PersistentFlags()
			flags.String(longName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
		}

		return nil
	},
}
