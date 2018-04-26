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

package bash

import (
	"os"
	"path"

	"github.com/joyent/freebsd-vpc/internal/buildtime"
	"github.com/joyent/freebsd-vpc/internal/command"
	"github.com/joyent/freebsd-vpc/internal/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const cmdName = "bash"

var Cmd = &command.Command{
	Name: cmdName,

	Cobra: &cobra.Command{
		Use:   cmdName,
		Short: "Generates and install " + buildtime.PROGNAME + " bash autocompletion script",
		Long: `Generates a bash autocompletion script for ` + buildtime.PROGNAME + `

By default, the file is written directly to ` + config.DefaultBashAutoCompletionDir + `
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
			bashDir := viper.GetString(config.KeyShellAutoCompBashDir)
			if _, err := os.Stat(bashDir); os.IsNotExist(err) {
				if err := os.MkdirAll(bashDir, 0777); err != nil {
					return errors.Wrapf(err, "unable to create bash autocomplete directory %q", bashDir)
				}
			}

			bashFile := path.Join(bashDir, buildtime.PROGNAME+".sh")
			err := cmd.Root().GenBashCompletionFile(bashFile)
			if err != nil {
				return errors.Wrap(err, "unable to generate bash completion")
			}

			log.Info().Msg("Installation completed successfully.")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		{
			const (
				key               = config.KeyShellAutoCompBashDir
				shortOpt, longOpt = "d", "dir"
				defaultValue      = config.DefaultBashAutoCompletionDir
				description       = "autocompletion directory"
			)

			flags := self.Cobra.Flags()
			flags.StringP(longOpt, shortOpt, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longOpt))
		}

		return nil
	},
}
