// Go interface to VPC EthLink objects.
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

package ethlink

import (
	"encoding/binary"
	"fmt"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc"
	"github.com/pkg/errors"
)

// _EthLinkCmd is the encoded type of operations that can be performed on a VPC
// EthLink.
type _EthLinkCmd vpc.Cmd

// _EthLinkCmdSetArgType is the value used by a VPC EthLink set operation.
type _EthLinkSetOpArgType uint64

const (
	// Bits for input
	_DownBit _EthLinkSetOpArgType = 0x00000000
	_UpBit   _EthLinkSetOpArgType = 0x00000001
)

// Ops that can be encoded into a vpc.Cmd
const (
	_OpInvalid          = vpc.Op(0)
	_OpConnect          = vpc.Op(1)
	_OpCloneAttach      = vpc.Op(2)
	_OpDevCtl           = vpc.Op(3)
	_OpDisconnect       = vpc.Op(4)
	_OpConnectedNameGet = vpc.Op(5)
	_OpVTagGet          = vpc.Op(6)
	_OpVTagSet          = vpc.Op(7)

	_ConnectCmd _EthLinkCmd = _EthLinkCmd(vpc.InBit|vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeLinkEth)<<16)) | _EthLinkCmd(_OpConnect)
	_VTagGetCmd _EthLinkCmd = _EthLinkCmd(vpc.OutBit|(vpc.Cmd(vpc.ObjTypeLinkEth)<<16)) | _EthLinkCmd(_OpVTagGet)
	_VTagSetCmd _EthLinkCmd = _EthLinkCmd(vpc.InBit|vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeLinkEth)<<16)) | _EthLinkCmd(_OpVTagSet)
)

// Close closes the VPC Handle.  Created EthLink will not be destroyed when the
// EthLink is closed if the EthLink has been Committed.
func (el *EthLink) Close() error {
	if el.h.FD() <= 0 {
		return nil
	}

	if err := el.h.Close(); err != nil {
		return errors.Wrap(err, "unable to close VPC EthLink handle")
	}

	return nil
}

// Commit increments the refcount of the EthLink in order to ensure the EthLink
// lives beyond the life of the current process and is not automatically cleaned
// up when the EthLink is closed.
func (el *EthLink) Commit() error {
	if el.h.FD() <= 0 {
		return errors.Errorf("unable to commit VPC EthLink handle with an empty descriptor")
	}

	if err := el.h.Commit(); err != nil {
		return errors.Wrap(err, "unable to commit VPC EthLink")
	}

	return nil
}

// Connect attaches the named physical device or cloned interface to this VPC
// EthLink.  The name of the device must be specified in the EthLink Config and
// passed in at Create time.
func (el *EthLink) Connect() error {
	// TODO(seanc@): Test to see make sure the descriptor has the mutate bit set.

	if err := vpc.Ctl(el.h, vpc.Cmd(_ConnectCmd), []byte(el.name), nil); err != nil {
		return errors.Wrap(err, "unable to connect VPC EthLink to a physical NIC")
	}

	return nil
}

// Destroy decrements the refcount of the VPC EthLink.  This EthLlink will be
// cleaned up when this VPC Handle is closed, however the object is destroyed
// before this call returns.  Some operations may still be performed on the open
// - and now invalidated - EthLink handle.
func (el *EthLink) Destroy() error {
	if el.h.FD() <= 0 {
		return nil
	}

	if err := el.h.Destroy(); err != nil {
		return errors.Wrap(err, "unable to destroy VPC EthLink")
	}

	return nil
}

// VTagGet returns the VTag (VLAN ID) associated with this EthLink interface.
func (el *EthLink) VTagGet() (vpc.VTag, error) {
	out := make([]byte, binary.MaxVarintLen64)
	if err := vpc.Ctl(el.h, vpc.Cmd(_VTagGetCmd), nil, out); err != nil {
		return 0, errors.Wrap(err, "unable to get VTag from EthLink")
	}

	vtagID, n := binary.Uvarint(out)
	if n > 0 && n <= 2 {
		switch {
		case vtagID < vpc.VTagMin:
			return 0, errors.Errorf("vtag value less than min: %d < %d", vtagID, vpc.VTagMin)
		case vtagID > vpc.VTagMax:
			return 0, errors.Errorf("vtag value greater than max: %d > %d", vtagID, vpc.VTagMax)
		default:
			return vpc.VTag(vtagID), nil
		}
	}

	panic(fmt.Sprintf("invariant: num vtag bytes read too big for kernel interface output (want/got: 2/%d", n))
}

// VTagSet sets the VTag (VLAN ID) on a given EthLink interface.  Setting the
// value to 0 clears the VTag.
func (el *EthLink) VTagSet(vtagID vpc.VTag) error {
	in := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(in, uint64(vtagID))
	if n < 2 {
		in = in[:2]
	} else {
		panic(fmt.Sprintf("invariant: vtag size too big for kernel interface input (want/got: 2/%d", n))
	}

	if err := vpc.Ctl(el.h, vpc.Cmd(_VTagSetCmd), in, nil); err != nil {
		return errors.Wrap(err, "unable to set the VTag for EthLink NIC")
	}

	return nil
}
