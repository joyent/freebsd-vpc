package listen

import (
	"fmt"
	"net"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/mux"
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
	_CmdName       = "listen"
	_KeyMuxID      = config.KeyMuxListenMuxID
	_KeyListenAddr = config.KeyMuxListenAddr
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "listen address to use when sending/receiving muxed VPC traffic",
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cons := conswriter.GetTerminal()

			cons.Write([]byte(fmt.Sprintf("Enabling VPC Mux listener...")))

			muxID, err := flag.GetMuxID(viper.GetViper(), _KeyMuxID)
			if err != nil {
				return errors.Wrap(err, "unable to get VPC Mux ID")
			}

			listenAddr := viper.GetString(_KeyListenAddr)
			if listenAddr == "" {
				return errors.Wrap(err, "unable to listening address for VPC Mux")
			}
			{
				host, port, err := net.SplitHostPort(listenAddr)
				if err != nil {
					return errors.Wrap(err, "unable to find host/port")
				}

				ip := net.ParseIP(host)
				if ip == nil {
					return errors.Wrap(err, "invalid IP address")
				}

				if port == "" {
					// FIXME(seanc@): turn this into a const
					port = "4789"
				}

				listenAddr = net.JoinHostPort(ip.String(), port)
				viper.Set(_KeyListenAddr, listenAddr)
			}

			muxCfg := mux.Config{
				ID:        muxID,
				Writeable: true,
			}

			vpcMux, err := mux.Open(muxCfg)
			if err != nil {
				log.Error().Err(err).Object("mux-id", muxID).Msg("VPC Mux open failed")
				return errors.Wrap(err, "unable to open VPC Mux")
			}
			defer vpcMux.Close()

			if err = vpcMux.Listen(listenAddr); err != nil {
				log.Error().Err(err).Object("mux-id", muxID).Str("listen-addr", listenAddr).Msg("vpc mux listen failed")
				return errors.Wrap(err, "unable to setup VPC Mux listener")
			}

			cons.Write([]byte("done.\n"))

			log.Info().Object("mux-id", muxID).Str("listen-addr", listenAddr).Msg("VPC Mux listener started")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		{
			const (
				key          = _KeyListenAddr
				longName     = "listen-addr"
				shortName    = ""
				defaultValue = ""
				description  = "Address and port the VPC Mux will use to listen for traffic on the underlay network"
			)

			flags := self.Cobra.Flags()
			flags.StringP(longName, shortName, defaultValue, description)

			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		if err := flag.AddMuxID(self, _KeyMuxID, true); err != nil {
			return errors.Wrap(err, "unable to register Mux ID flag on VPC Mux connect")
		}

		return nil
	},
}
