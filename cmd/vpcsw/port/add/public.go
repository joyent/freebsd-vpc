package add

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
	"go.freebsd.org/sys/vpc"
	"go.freebsd.org/sys/vpc/l2link"
	"go.freebsd.org/sys/vpc/vpcp"
	"go.freebsd.org/sys/vpc/vpcsw"
	"go.freebsd.org/sys/vpc/vpctest"
)

const (
	_CmdName     = "add"
	_KeyL2Name   = config.KeySWPortAddL2Name
	_KeyPortID   = config.KeySWPortAddID
	_KeyPortMAC  = config.KeySWPortAddMAC
	_KeySwitchID = config.KeySWPortAddSwitchID
	_KeyUplinkID = config.KeySWPortAddUplinkID
)

var Cmd = &command.Command{
	Name: _CmdName,

	Cobra: &cobra.Command{
		Use:          _CmdName,
		Short:        "add a port to a VPC Switch",
		Aliases:      []string{"create"},
		SilenceUsage: true,
		// TraverseChildren: true,
		Args: cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			switch {
			case viper.GetString(_KeyUplinkID) != "" && viper.GetString(_KeyL2Name) == "":
				// TODO(seanc@): convert uplink-id and l2-name to constants used by
				// cobra when setting the viper key.
				return errors.Errorf("uplink-id requires an l2-name")
			}

			if l2Name := viper.GetString(_KeyL2Name); l2Name != "" {
				existingIfaces, err := vpctest.GetAllInterfaces()
				if err != nil {
					return errors.Wrapf(err, "unable to get all interfaces")
				}

				var found bool
				for _, iface := range existingIfaces {
					if l2Name == iface.Name {
						found = true
						break
					}
				}

				if !found {
					return errors.Errorf("unable to find interface %q", l2Name)
				}
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cons := conswriter.GetTerminal()

			cons.Write([]byte(fmt.Sprintf("Adding port to VPC Switch...")))

			switchID, err := flag.GetSwitchID(viper.GetViper(), _KeySwitchID)
			if err != nil {
				return errors.Wrap(err, "unable to get VPC ID")
			}

			portID, err := flag.GetPortID(viper.GetViper(), _KeyPortID)
			if err != nil {
				return errors.Wrap(err, "unable to get VPC Switch Port ID")
			}

			var uplinkID vpc.ID
			if uplinkStr := viper.GetString(_KeyUplinkID); uplinkStr != "" {
				uplinkID, err = vpc.ParseID(uplinkStr)
				if err != nil {
					// TODO(seanc@): convert uplink-id to a constant usable within this
					// package.
					return errors.Wrapf(err, "unable to parse uplink-id %q", viper.GetString(_KeyUplinkID))
				}
			}

			portMAC, err := flag.GetMAC(viper.GetViper(), _KeyPortMAC, nil)
			if err != nil {
				return errors.Wrap(err, "unable to get MAC address")
			}

			l2Name := viper.GetString(_KeyL2Name)

			// Create a stack of commit and undo operations to walk through in the
			// event of an error.
			var commit bool
			var commitFuncs, undoFuncs []func()
			defer func() {
				if commit {
					for _, fn := range commitFuncs {
						fn()
					}
				} else {
					for _, fn := range undoFuncs {
						fn()
					}
				}
			}()
			commitFuncs = append(commitFuncs, func() {
				cons.Write([]byte("done.\n"))
			})

			// 1) Open switch and add a port
			switchCfg := vpcsw.Config{
				ID:        switchID,
				Writeable: true,
			}

			vpcSwitch, err := vpcsw.Open(switchCfg)
			if err != nil {
				log.Error().Err(err).Object("switch-cfg", switchCfg).Msg("vpcsw open failed")
				return errors.Wrap(err, "unable to open VPC Switch")
			}
			commitFuncs = append(commitFuncs, func() {
				if err := vpcSwitch.Close(); err != nil {
					log.Error().Err(err).Msg("unable to commit VPC Switch")
				}
			})
			undoFuncs = append(undoFuncs, func() {
				if err := vpcSwitch.Close(); err != nil {
					log.Error().Err(err).Msg("unable to clean up VPC Switch during error recovery")
				}
			})

			if err = vpcSwitch.PortAdd(portID, portMAC); err != nil {
				log.Error().Err(err).
					Object("port-id", portID).
					Str("port-mac", portMAC.String()).
					Object("switch-cfg", switchCfg).
					Msg("failed to add VPC Switch Port")
				return errors.Wrap(err, "unable to add a port to VPC Switch")
			}

			// If we have an L2 Link, add it to the port
			if l2Name != "" {
				l2Cfg := l2link.Config{
					ID:   uplinkID,
					Name: l2Name,
				}
				l2, err := l2link.Create(l2Cfg)
				if err != nil {
					return errors.Wrap(err, "unable to create VPC L2 Link")
				}
				commitFuncs = append(commitFuncs, func() {
					if err := l2.Commit(); err != nil {
						log.Error().Err(err).Msg("unable to commit VPC L2 Link")
					}
				})
				undoFuncs = append(undoFuncs, func() {
					if err := l2.Destroy(); err != nil {
						log.Error().Err(err).Msg("unable to clean up VPC L2 Link during error recovery")
					}

					if err := l2.Close(); err != nil {
						log.Error().Err(err).Msg("unable to clean up VPC L2 Link during error recovery")
					}
				})

				if err := l2.Attach(); err != nil {
					l2.Close()
					return errors.Wrapf(err, "unable to attach L2 link to device %q", l2Name)
				}

				// if err := vpcSwitch.PortAdd(portID); err != nil {
				// 	log.Error().Err(err).Object("port-id", portID).Object("switch-cfg", switchCfg).Msg("failed to add VPC Switch Port")
				// }

				if err := vpcSwitch.PortUplinkSet(portID); err != nil {
					log.Error().Err(err).Object("port-id", portID).Object("switch-cfg", switchCfg).Msg("failed to set VPC Switch Port as an Uplink port")
				}

				portCfg := vpcp.Config{
					ID: portID,
				}
				vpcPort, err := vpcp.Open(portCfg)
				if err != nil {
					log.Error().Err(err).Object("port-id", portID).Object("switch-cfg", switchCfg).Msg("failed to connect VPC interface to VPC Switch Port")
				}

				if err := vpcPort.Connect(l2Cfg.ID); err != nil {
					log.Error().Err(err).Object("port-id", portID).Object("interface", switchCfg).Msg("failed to connect VPC interface to VPC Switch Port")
				}
			}

			commit = true

			// log.Info().Str("port-id", portAddCfg.ID.String()).Str("switch-id", switchID.String()).Str("uplink-id", uplinkID.String()). /*.Str("name", newPort.Name)*/ Msg("vpcp created")
			log.Info().Object("port-id", portID).Str("switch-id", switchID.String()).Msg("vpcp created")

			return nil
		},
	},

	Setup: func(self *command.Command) error {
		if err := flag.AddPortID(self, _KeyPortID, false); err != nil {
			return errors.Wrap(err, "unable to register Port ID flag on VPC Switch Port add")
		}

		if err := flag.AddMAC(self, _KeyPortMAC, false); err != nil {
			return errors.Wrap(err, "unable to register MAC flag on VPC Switch Port add")
		}

		if err := flag.AddSwitchID(self, _KeySwitchID, false); err != nil {
			return errors.Wrap(err, "unable to register Switch ID flag for VPC Switch Port add")
		}

		{
			const (
				key          = _KeyL2Name
				longName     = "l2-name"
				shortName    = "n"
				defaultValue = ""
				description  = "Name of the L2 interface to use as an uplink in the VPC Switch"
			)

			flags := self.Cobra.Flags()
			flags.StringP(longName, shortName, defaultValue, description)

			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		{
			const (
				key          = _KeyUplinkID
				longName     = "uplink-id"
				shortName    = "u"
				defaultValue = ""
				description  = "Specify the ID of the VPC Switch's uplink port"
			)

			flags := self.Cobra.Flags()
			flags.StringP(longName, shortName, defaultValue, description)

			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		return nil
	},
}

func _GetUplinkID(v *viper.Viper, key string) (id vpc.ID, err error) {
	uplinkIDStr := v.GetString(key)
	if uplinkIDStr == "" {
		return vpc.ID{}, errors.Wrap(err, "missing VPC Uplink ID")
	}

	if id, err = vpc.ParseID(uplinkIDStr); err != nil {
		return vpc.ID{}, errors.Wrapf(err, "unable to parse VPC Uplink ID %q", uplinkIDStr)
	}

	return id, nil
}
