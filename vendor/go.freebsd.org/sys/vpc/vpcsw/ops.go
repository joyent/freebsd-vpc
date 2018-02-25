// Go interface to VPC Switch objects.
//
// SPDX-License-Identifier: BSD-2-Clause-FreeBSD
//
// Copyright (C) 2018 Sean Chittenden <seanc@joyent.com>
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

package vpcsw

import (
	"net"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.freebsd.org/sys/vpc"
)

// _SwitchCmd is the encoded type of operations that can be performed on a VPC
// Switch.
type _SwitchCmd vpc.Cmd

// _SwitchCmdSetArgType is the value used by a VPC Switch set operation.
type _SwitchSetOpArgType uint64

const (
	// Bits for input
	_DownBit _SwitchSetOpArgType = 0x00000000
	_UpBit   _SwitchSetOpArgType = 0x00000001
)

// Ops that can be encoded into a vpc.Cmd
const (
	_OpInvalid       = vpc.Op(0)
	_OpPortAdd       = vpc.Op(1)
	_OpPortDel       = vpc.Op(2)
	_OpPortUplinkSet = vpc.Op(3)
	_OpPortUplinkGet = vpc.Op(4)
	_OpStateGet      = vpc.Op(5)
	_OpStateSet      = vpc.Op(6)
	_OpReset         = vpc.Op(7)

	_PortAddCmd       _SwitchCmd = _SwitchCmd(vpc.InBit|vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeSwitch)<<16)) | _SwitchCmd(_OpPortAdd)
	_PortRemoveCmd    _SwitchCmd = _SwitchCmd(vpc.InBit|vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeSwitch)<<16)) | _SwitchCmd(_OpPortDel)
	_PortUplinkSetCmd _SwitchCmd = _SwitchCmd(vpc.InBit|vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeSwitch)<<16)) | _SwitchCmd(_OpPortUplinkSet)
)

// Template commands that can be passed to vpc.Ctl() with a valid VPC Switch
// Handle.
var (
	_PortDelCmd    _SwitchCmd
	_PortUplinkGet _SwitchCmd
	_PortStateSet  _SwitchCmd
	_PortStateGet  _SwitchCmd
	_ResetCmd      _SwitchCmd
)

// Close closes the VPC Handle descriptor.  Created VPC Switches will not be
// destroyed when the VPCSW is closed if the VPC Switch has been Committed.
func (sw *VPCSW) Close() error {
	if sw.h.FD() <= 0 {
		return nil
	}

	if err := sw.h.Close(); err != nil {
		return errors.Wrap(err, "unable to close VPC handle")
	}

	return nil
}

// Commit increments the refcount of the VPC Switch in order to ensure the VPC
// Switch lives beyond the life of the current process and is not automatically
// cleaned up when the VPCSW is closed.
func (sw *VPCSW) Commit() error {
	if sw.h.FD() <= 0 {
		return nil
	}

	if err := sw.h.Commit(); err != nil {
		return errors.Wrap(err, "unable to commit VPC Switch")
	}

	return nil
}

// Destroy decrements the refcount of the VPC Switch in destroy the the VPC
// Switch when the VPC Handle is closed.
func (sw *VPCSW) Destroy() error {
	if sw.h.FD() <= 0 {
		return nil
	}

	if err := sw.h.Destroy(); err != nil {
		return errors.Wrap(err, "unable to destroy VPC Switch")
	}

	return nil
}

// PortAdd adds an existing VPC Port to this VPC Switch.  PortID is VPC ID of
// the existing VPC Port to be added to this switch.
func (sw *VPCSW) PortAdd(portID vpc.ID, mac net.HardwareAddr) error {
	// TODO(seanc@): Test to see make sure the descriptor has the mutate bit set.

	// Create the port
	if err := vpc.Ctl(sw.h, vpc.Cmd(_PortAddCmd), portID.Bytes(), nil); err != nil {
		return errors.Wrap(err, "unable to add a VPC Port to VPC Switch")
	}

	// TODO(seanc@): Set the MAC address of the port

	return nil
}

// PortRemove removes a VPC Port from this VPC Switch.  Uses the PortID member
// of Config.
func (sw *VPCSW) PortRemove(cfg Config) error {
	// TODO(seanc@): Test to see make sure the descriptor has the mutate bit set.

	if err := vpc.Ctl(sw.h, vpc.Cmd(_PortRemoveCmd), cfg.PortID.Bytes(), nil); err != nil {
		log.Error().Err(err).
			Object("cfg", cfg).
			Object("cmd", vpc.Cmd(_PortRemoveCmd)).
			Str("cmd", "port remove").
			Str("obj-type", "switch").
			Msg("failed")
		return errors.Wrap(err, "unable to remove a VPC Port from VPC Switch")
	}

	return nil
}

// Reset resets the VPC Switch.
func (sw *VPCSW) Reset() error {
	if sw.h.FD() <= 0 {
		return nil
	}

	// TODO(seanc@): Test to see make sure the descriptor has the mutate bit set.

	if err := vpc.Ctl(sw.h, vpc.Cmd(_ResetCmd), nil, nil); err != nil {
		return errors.Wrap(err, "unable to reset VPC Switch")
	}

	return nil
}

// UplinkSet designates an existing VPC Port as an uplink port for this VPC
// Switch.
func (sw *VPCSW) PortUplinkSet(portID vpc.ID) error {
	// TODO(seanc@): Test to see make sure the descriptor has the mutate bit set.

	// Create the port
	if err := vpc.Ctl(sw.h, vpc.Cmd(_PortUplinkSetCmd), portID.Bytes(), nil); err != nil {
		return errors.Wrap(err, "unable to set VPC Port as uplink in VPC Switch")
	}

	return nil
}
