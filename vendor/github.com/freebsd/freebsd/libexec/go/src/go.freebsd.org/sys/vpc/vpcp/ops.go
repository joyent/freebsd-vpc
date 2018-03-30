// Go interface to VM NIC objects.
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

package vpcp

import (
	"encoding/binary"
	"fmt"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc"
	"github.com/pkg/errors"
)

// _PortCmd is the encoded type of operations that can be performed on a VM
// NIC.
type _PortCmd vpc.Cmd

// _PortCmdSetArgType is the value used by a VM NIC set operation.
type _PortSetOpArgType uint64

// Ops that can be encoded into a vpc.Cmd
const (
	_OpInvalid    = vpc.Op(0)
	_OpConnect    = vpc.Op(1)
	_OpDisconnect = vpc.Op(2)
	_OpVNIGet     = vpc.Op(3)
	_OpVNISet     = vpc.Op(4)
	_OpVLANGet    = vpc.Op(5)
	_OpVLANSet    = vpc.Op(6)
	_             = vpc.Op(7) // Unused
	_             = vpc.Op(8) // Unused
	_OpPeerIDGet  = vpc.Op(9)

	_ConnectCmd    _PortCmd = _PortCmd(vpc.InBit|vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeSwitchPort)<<16)) | _PortCmd(_OpConnect)
	_DisconnectCmd _PortCmd = _PortCmd(vpc.InBit|vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeSwitchPort)<<16)) | _PortCmd(_OpDisconnect)
	_VNIGetCmd     _PortCmd = _PortCmd(vpc.OutBit|(vpc.Cmd(vpc.ObjTypeSwitchPort)<<16)) | _PortCmd(_OpVNIGet)
	_VNISetCmd     _PortCmd = _PortCmd(vpc.InBit|vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeSwitchPort)<<16)) | _PortCmd(_OpVNISet)
)

// Connect a VPC Interface to this VPC Port.  VPC Interfaces include VMNIC, and
// L2Link.
func (port *VPCP) Connect(interfaceID vpc.ID) error {
	// TODO(seanc@): Test to see make sure the descriptor has the mutate bit set.

	if err := vpc.Ctl(port.h, vpc.Cmd(_ConnectCmd), interfaceID.Bytes(), nil); err != nil {
		return errors.Wrap(err, "unable to connect VPC Interface to VPC Switch Port")
	}

	return nil
}

// Disconnect a VPC Interface from this VPC Port.  VPC Interfaces include VMNIC,
// and L2Link.
func (port *VPCP) Disconnect(interfaceID vpc.ID) error {
	// TODO(seanc@): Test to see make sure the descriptor has the mutate bit set.

	if err := vpc.Ctl(port.h, vpc.Cmd(_DisconnectCmd), interfaceID.Bytes(), nil); err != nil {
		return errors.Wrap(err, "unable to disconnect VPC Interface from VPC Switch Port")
	}

	return nil
}

// GetVNI gets the VNI assigned to a VPC Switch Port.  A value of 0 means the
// VPC Switch Port has no VNI assigned.
func (port *VPCP) GetVNI() (vpc.VNI, error) {
	out := make([]byte, binary.MaxVarintLen64)
	if err := vpc.Ctl(port.h, vpc.Cmd(_VNIGetCmd), nil, out); err != nil {
		return 0, errors.Wrap(err, "unable to get the VNI of a VPC Switch Port")
	}

	// binary.LittleEndian.Uint16(uuidRaw[4:])
	vni := vpc.VNI(binary.LittleEndian.Uint32(out))
	if vni >= vpc.VNIMin || vni <= vpc.VNIMax {
		return vpc.VNI(vni), nil
	}

	panic(fmt.Sprintf("invariant: VNI too big for kernel interface output"))
}

// SetVNI sets the VNI on a VPC Switch Port.  A value of 0 unsets the value on a
// VPC Port.
func (port *VPCP) SetVNI(vni vpc.VNI) error {
	in := make([]byte, binary.MaxVarintLen64)
	binary.BigEndian.PutUint32(in[0:4], uint32(vni))

	if err := vpc.Ctl(port.h, vpc.Cmd(_VNISetCmd), in, nil); err != nil {
		return errors.Wrap(err, "unable to set the VNI on VPC Switch Port")
	}

	return nil
}
