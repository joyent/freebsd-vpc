// Copyright (c) 2018 Joyent, Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR AND CONTRIBUTORS ``AS IS'' AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
// OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
// HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
// OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
// SUCH DAMAGE.

package man

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/joyent/freebsd-vpc/internal/buildtime"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/joyent/freebsd-vpc/internal/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

const _CmdName = "man"

var Cmd = &command.Command{
	Name: _CmdName,
	Cobra: &cobra.Command{
		Use:   _CmdName,
		Short: "Generates and install " + buildtime.PROGNAME + " man(1) pages",
		Long: `This command automatically generates up-to-date man(1) pages of ` + buildtime.PROGNAME + fmt.Sprintf("(%d)", config.ManSect) + `
command-line interface.  By default, it creates the man page files
in the "docs/man" directory under the current directory.`,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			now, err := time.Parse(time.RFC3339, buildtime.DocsDate)
			if err != nil {
				log.Warn().Err(err).Msg("unable to parse docsdate")
				now = time.Now()
			}

			header := &doc.GenManHeader{
				Manual:  buildtime.PROGNAME,
				Section: strconv.Itoa(config.ManSect),
				Source:  strings.Join([]string{buildtime.PROGNAME, buildtime.Version}, " "),
				Date:    &now,
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

			err = doc.GenManTree(cmd.Root(), header, manSectDir)
			if err != nil {
				return errors.Wrap(err, "unable to generate man(1) pages")
			}

			log.Info().Msg("Installation completed successfully.")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		{
			const (
				key          = config.KeyDocManDir
				longName     = "man-dir"
				shortName    = "m"
				description  = "Specify the MANDIR to use"
				defaultValue = config.DefaultManDir
			)

			flags := self.Cobra.PersistentFlags()
			flags.StringP(longName, shortName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
			viper.BindEnv(key, "MANDIR")
			viper.SetDefault(key, defaultValue)
		}

		return nil
	},
}
