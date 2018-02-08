package md

import (
	"fmt"
	"os"

	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/internal/buildtime"
	"github.com/sean-/vpc/internal/command"
	"github.com/sean-/vpc/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/spf13/viper"
)

const template = `---
date: %s
title: "%s"
slug: %s
url: %s
---
`

var Cmd = &command.Command{
	Cobra: &cobra.Command{
		Use:   "md",
		Short: "Generates and install " + buildtime.PROGNAME + " markdown pages",
		Long: `
Generate Markdown documentation for the ` + buildtime.PROGNAME + `
It creates one Markdown file per command `,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			mdDir := viper.GetString(config.KeyDocMdDir)

			if _, err := os.Stat(mdDir); os.IsNotExist(err) {
				if err := os.MkdirAll(mdDir, 0777); err != nil {
					return errors.Wrapf(err, "unable to make mddir %q", mdDir)
				}
			}

			log.Info().Str("MDDIR", mdDir).Int("section", config.ManSect).Msg("Installing markdown pages")

			now := time.Now().UTC().Format(time.RFC3339)
			prepender := func(filename string) string {
				name := filepath.Base(filename)
				base := strings.TrimSuffix(name, path.Ext(name))
				url := "/commands/" + strings.ToLower(base) + "/"
				return fmt.Sprintf(template, now, strings.Replace(base, "_", " ", -1), base, url)
			}

			linkHandler := func(name string) string {
				base := strings.TrimSuffix(name, path.Ext(name))
				return "/commands/" + strings.ToLower(base) + "/"
			}

			doc.GenMarkdownTreeCustom(cmd.Root(), mdDir, prepender, linkHandler)

			log.Info().Msg("Installation completed successfully.")

			return nil
		},
	},

	Setup: func(parent *command.Command) error {
		{
			const (
				key          = config.KeyDocMdDir
				longName     = "md-dir"
				description  = "Specify the MDDIR to use"
				defaultValue = config.DefaultMdDir
			)

			flags := parent.Cobra.PersistentFlags()
			flags.String(longName, defaultValue, description)
			viper.BindPFlag(key, flags.Lookup(longName))
			viper.BindEnv(key, "MANDIR")
			viper.SetDefault(key, defaultValue)
		}

		return nil
	},
}
