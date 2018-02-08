package man

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/internal/buildtime"
	"github.com/sean-/vpc/internal/command"
	"github.com/sean-/vpc/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Use:   "man",
		Short: "Generates and install " + buildtime.PROGNAME + " man(1) pages",
		Long: `
This command automatically generates up-to-date man(1) pages of ` + buildtime.PROGNAME + fmt.Sprintf("(%d)", config.ManSect) + `
command-line interface.  By default, it creates the man page files
in the "docs/man" directory under the current directory.`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			header := &doc.GenManHeader{
				Manual:  buildtime.PROGNAME,
				Section: strconv.Itoa(config.ManSect),
				Source:  strings.Join([]string{buildtime.PROGNAME, buildtime.Version}, " "),
			}

			manDir := viper.GetString(config.KeyDocManDir)

			manSectDir := path.Join(manDir, fmt.Sprintf("man%d", config.ManSect))
			if _, err := os.Stat(manSectDir); os.IsNotExist(err) {
				if err := os.MkdirAll(manSectDir, 0777); err != nil {
					return errors.Wrapf(err, "unable to make mandir %q", manSectDir)
				}
			}

			cmd.Root().DisableAutoGenTag = true
			log.Info().Str("MANDIR", manDir).Int("section", config.ManSect).Msg("Installing man(1) pages")

			err := doc.GenManTree(cmd.Root(), header, manSectDir)
			if err != nil {
				return errors.Wrap(err, "unable to generate man(1) pages")
			}

			log.Info().Msg("Installation completed successfully.")

			return nil
		},
	},

	Setup: func(parent *command.Command) error {
		{
			const (
				key          = config.KeyDocManDir
				longName     = "man-dir"
				shortName    = "m"
				description  = "Specify the MANDIR to use"
				defaultValue = config.DefaultManDir
			)

			flags := parent.Cobra.PersistentFlags()
			flags.StringP(longName, shortName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
			viper.BindEnv(key, "MANDIR")
			viper.SetDefault(key, defaultValue)
		}

		return nil
	},
}
