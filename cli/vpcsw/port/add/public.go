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
	_KeyL2ID     = config.KeySWPortAddL2ID
	_KeyL2Name   = config.KeySWPortAddL2Name
	_KeyPortID   = config.KeySWPortAddID
	_KeyPortMAC  = config.KeySWPortAddMAC
	_KeySwitchID = config.KeySWPortAddSwitchID
	_KeyUplink   = config.KeySWPortAddUplink
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
			case viper.GetString(_KeyL2ID) != "" && viper.GetString(_KeyL2Name) == "":
				// TODO(seanc@): convert l2-id and l2-name to constants used by
				// cobra when setting the viper key.
				return errors.Errorf("l2-id requires an l2-name")
			case viper.GetString(_KeyL2ID) == "" && viper.GetString(_KeyL2Name) != "":
				// TODO(seanc@): convert l2-id and l2-name to constants used by
				// cobra when setting the viper key.
				return errors.Errorf("l2-name requires an l2-id")
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

			l2ID, err := _GetLinkID(viper.GetViper(), _KeyL2ID, true)
			if err != nil {
				return errors.Wrap(err, "unable to get l2-id")
			}

			portMAC, err := flag.GetMAC(viper.GetViper(), _KeyPortMAC, nil)
			if err != nil {
				return errors.Wrap(err, "unable to get MAC address")
			}

			l2Name := viper.GetString(_KeyL2Name)

			// Create a stack of commit and undo operations to walk through in the
			// event of an error.
			var commit bool
			var commitFuncs, undoFuncs []func() error
			defer func() {
				scopeHandlers := undoFuncs
				modeStr := "undo"
				if commit {
					modeStr = "commit"
					scopeHandlers = commitFuncs
				}

				for i := len(scopeHandlers) - 1; i >= 0; i-- {
					if err := scopeHandlers[i](); err != nil {
						log.Error().Err(err).Msgf("failure during %s", modeStr)
					}
				}
			}()
			commitFuncs = append(commitFuncs, func() error {
				cons.Write([]byte("done.\n"))
				return nil
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
			commitFuncs = append(commitFuncs, func() error {
				if err := vpcSwitch.Close(); err != nil {
					log.Error().Err(err).Msg("unable to commit VPC Switch")
					return errors.Wrap(err, "unable to commit VPC switch during operation commit")
				}

				return nil
			})
			undoFuncs = append(undoFuncs, func() error {
				if err := vpcSwitch.Close(); err != nil {
					log.Error().Err(err).Msg("unable to clean up VPC Switch during error recovery")
				}

				return nil
			})

			// If we have an L2 Link, add it to the port
			switch {
			case l2Name == "":
				if err = vpcSwitch.PortAdd(portID, portMAC); err != nil {
					log.Error().Err(err).
						Object("port-id", portID).
						Str("port-mac", portMAC.String()).
						Object("switch-cfg", switchCfg).
						Msg("failed to add VPC Switch Port")
					return errors.Wrap(err, "unable to add a port to VPC Switch")
				}
			case l2Name != "":
				l2Cfg := l2link.Config{
					ID:   l2ID,
					Name: l2Name,
				}
				l2, err := l2link.Create(l2Cfg)
				if err != nil {
					return errors.Wrap(err, "unable to create VPC L2 Link")
				}

				if err := l2.Attach(); err != nil {
					return errors.Wrapf(err, "unable to attach L2 link to device %q", l2Name)
				}
				commitFuncs = append(commitFuncs, func() error {
					if err := l2.Commit(); err != nil {
						log.Error().Err(err).Object("l2", l2).Msg("unable to commit VPC L2 Link")
						return errors.Wrap(err, "unable to commit VPC L2 Link")
					}
					return nil
				})

				if viper.GetBool(_KeyUplink) {
					if err := vpcSwitch.PortUplinkSet(portID, portMAC); err != nil {
						log.Error().Err(err).Object("port-id", portID).Object("switch-cfg", switchCfg).Msg("failed to set VPC Switch Port as an Uplink port")
						return errors.Wrap(err, "unable to create a VPC Switch Port uplink")
					}
				}

				portCfg := vpcp.Config{
					ID:        portID,
					Writeable: true,
				}
				vpcPort, err := vpcp.Open(portCfg)
				if err != nil {
					log.Error().Err(err).Object("port-id", portID).Object("switch-cfg", switchCfg).Msg("failed to connect VPC interface to VPC Switch Port")
					return errors.Wrap(err, "unable to open VPC Switch Port")
				}

				if err := vpcPort.Connect(l2Cfg.ID); err != nil {
					log.Error().Err(err).Object("l2-cfg", l2Cfg).Object("l2", l2).Object("port-id", portID).Object("switch-cfg", switchCfg).Msg("failed to connect VPC interface to VPC Switch Port")
					return errors.Wrap(err, "unable to connect VPC Interface to VPC Port")
				}
			default:
				panic("invalid switch port add logic")
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
				key          = _KeyL2ID
				longName     = "l2-id"
				shortName    = ""
				defaultValue = ""
				description  = "Specify the ID of the VPC L2 Link"
			)

			flags := self.Cobra.Flags()
			flags.StringP(longName, shortName, defaultValue, description)

			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		{
			const (
				key          = _KeyUplink
				longName     = "uplink"
				shortName    = "u"
				defaultValue = false
				description  = "make the port ID an uplink for the switch"
			)

			flags := self.Cobra.Flags()
			flags.BoolP(longName, shortName, defaultValue, description)

			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		return nil
	},
}

func _GetLinkID(v *viper.Viper, key string, optional bool) (id vpc.ID, err error) {
	l2IDStr := v.GetString(key)
	switch {
	case optional && l2IDStr == "":
		return vpc.ID{}, nil
	case !optional && l2IDStr == "":
		return vpc.ID{}, errors.Wrap(err, "missing VPC L2 Link ID")
	}

	if id, err = vpc.ParseID(l2IDStr); err != nil {
		return vpc.ID{}, errors.Wrapf(err, "unable to parse VPC L2 Link ID %q", l2IDStr)
	}

	return id, nil
}