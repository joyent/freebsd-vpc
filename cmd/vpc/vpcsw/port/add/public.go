package add

import (
	"fmt"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc"
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/ethlink"
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/vpcp"
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/vpcsw"
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/vpctest"
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
	_CmdName      = "add"
	_KeyEthLinkID = config.KeySWPortAddEthLinkID
	_KeyL2Name    = config.KeySWPortAddL2Name
	_KeyPortID    = config.KeySWPortAddID
	_KeyPortMAC   = config.KeySWPortAddMAC
	_KeySwitchID  = config.KeySWPortAddSwitchID
	_KeyUplink    = config.KeySWPortAddUplink
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
			case viper.GetString(_KeyEthLinkID) != "" && viper.GetString(_KeyL2Name) == "":
				// TODO(seanc@): convert ethlink-id and l2-name to constants used by
				// cobra when setting the viper key.
				return errors.Errorf("ethlink-id requires an l2-name")
			case viper.GetString(_KeyEthLinkID) == "" && viper.GetString(_KeyL2Name) != "":
				// TODO(seanc@): convert ethlink-id and l2-name to constants used by
				// cobra when setting the viper key.
				return errors.Errorf("l2-name requires an ethlink-id")
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

			ethLinkID, err := _GetLinkID(viper.GetViper(), _KeyEthLinkID, true)
			if err != nil {
				return errors.Wrap(err, "unable to get ethlink-id")
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

			// If we have an EthLink, add it to the port
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
				ethLinkCfg := ethlink.Config{
					ID:   ethLinkID,
					Name: l2Name,
				}
				el, err := ethlink.Create(ethLinkCfg)
				if err != nil {
					return errors.Wrap(err, "unable to create VPC EthLink")
				}

				if err := el.Connect(); err != nil {
					return errors.Wrapf(err, "unable to connect L2 link to device %q", l2Name)
				}
				commitFuncs = append(commitFuncs, func() error {
					if err := el.Commit(); err != nil {
						log.Error().Err(err).Object("ethlink", el).Msg("unable to commit VPC EthLink")
						return errors.Wrap(err, "unable to commit VPC EthLink")
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

				if err := vpcPort.Connect(ethLinkCfg.ID); err != nil {
					log.Error().Err(err).Object("ethlink-cfg", ethLinkCfg).Object("ethlink", el).Object("port-id", portID).Object("switch-cfg", switchCfg).Msg("failed to connect VPC interface to VPC Switch Port")
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
				description  = "Name of the underlying L2 interface to be wrapped as a VPC EthLink and used as the uplink in the VPC Switch"
			)

			flags := self.Cobra.Flags()
			flags.StringP(longName, shortName, defaultValue, description)

			viper.BindPFlag(key, flags.Lookup(longName))
			viper.SetDefault(key, defaultValue)
		}

		{
			const (
				key          = _KeyEthLinkID
				longName     = "ethlink-id"
				shortName    = ""
				defaultValue = ""
				description  = "Specify the ID of the VPC EthLink"
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
	ethLinkIDStr := v.GetString(key)
	switch {
	case optional && ethLinkIDStr == "":
		return vpc.ID{}, nil
	case !optional && ethLinkIDStr == "":
		return vpc.ID{}, errors.Wrap(err, "missing VPC EthLink ID")
	}

	if id, err = vpc.ParseID(ethLinkIDStr); err != nil {
		return vpc.ID{}, errors.Wrapf(err, "unable to parse VPC EthLink ID %q", ethLinkIDStr)
	}

	return id, nil
}
