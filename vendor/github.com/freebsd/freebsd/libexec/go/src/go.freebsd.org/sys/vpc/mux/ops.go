// Go interface to VPC Mux objects.
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

package mux

import (
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc"
	"github.com/pkg/errors"
)

// _MuxCmd is the encoded type of operations that can be performed on a VPC
// Mux.
type _MuxCmd vpc.Cmd

// _MuxCmdSetArgType is the value used by a VPC Mux set operation.
type _MuxSetOpArgType uint64

const (
	// Bits for input
	_DownBit _MuxSetOpArgType = 0x00000000
	_UpBit   _MuxSetOpArgType = 0x00000001
)

// Ops that can be encoded into a vpc.Cmd
const (
	_OpInvalid               = vpc.Op(0)
	_OpMuxListen             = vpc.Op(1)
	_OpMuxFTESet             = vpc.Op(2)
	_OpMuxFTEDel             = vpc.Op(3)
	_OpMuxFTEList            = vpc.Op(4)
	_OpMuxUnderlayConnect    = vpc.Op(5)
	_OpMuxUnderlayDisconnect = vpc.Op(6)
	_OpMuxConnectedIDGet     = vpc.Op(7)
)

// Template commands that can be passed to vpc.Ctl() with a valid VPC Mux
// Handle.
const (
	_MuxListenCmd             _MuxCmd = _MuxCmd(vpc.InBit|vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeMux)<<16)) | _MuxCmd(_OpMuxListen)
	_MuxFTESetCmd             _MuxCmd = _MuxCmd(vpc.InBit|vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeMux)<<16)) | _MuxCmd(_OpMuxFTESet)
	_MuxFTEDelCmd             _MuxCmd = _MuxCmd(vpc.InBit|vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeMux)<<16)) | _MuxCmd(_OpMuxFTEDel)
	_MuxFTEListCmd            _MuxCmd = _MuxCmd(vpc.InBit|(vpc.Cmd(vpc.ObjTypeMux)<<16)) | _MuxCmd(_OpMuxFTEList)
	_MuxUnderlayConnectCmd    _MuxCmd = _MuxCmd(vpc.InBit|vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeMux)<<16)) | _MuxCmd(_OpMuxUnderlayConnect)
	_MuxUnderlayDisconnectCmd _MuxCmd = _MuxCmd(vpc.InBit|vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeMux)<<16)) | _MuxCmd(_OpMuxUnderlayDisconnect)
	_MuxConnectedIDGetCmd     _MuxCmd = _MuxCmd(vpc.InBit|vpc.OutBit|(vpc.Cmd(vpc.ObjTypeMux)<<16)) | _MuxCmd(_OpMuxConnectedIDGet)
)

// Close closes the VPC Mux Handle descriptor.  VPC Muxes will not be destroyed
// when the Mux is closed if the VPC Mux has been Committed.
func (m *Mux) Close() error {
	if m.h.FD() <= 0 {
		return nil
	}

	if err := m.h.Close(); err != nil {
		return errors.Wrap(err, "unable to close VPC Mux handle")
	}

	return nil
}

// Commit increments the refcount of the VPC Mux in order to ensure the Mux
// lives beyond the life of the current process and is not automatically cleaned
// up when the Mux handle is closed.
func (m *Mux) Commit() error {
	if m.h.FD() <= 0 {
		return nil
	}

	if err := m.h.Commit(); err != nil {
		return errors.Wrap(err, "unable to commit VPC Mux")
	}

	return nil
}

// Connect a VPC Mux to a VPC Interface.
func (m *Mux) Connect(interfaceID vpc.ID) error {
	// TODO(seanc@): Test to see make sure the descriptor has the mutate bit set.

	if err := vpc.Ctl(m.h, vpc.Cmd(_MuxUnderlayConnectCmd), interfaceID.Bytes(), nil); err != nil {
		return errors.Wrap(err, "unable to connect VPC Mux to to VPC Interface")
	}

	return nil
}

// Destroy decrements the refcount of the VPC Mux and destroys the object.  The
// Mux resources are cleaned up when the VPC Handle is closed, however the
// object will stop processing traffic when the destroy command is issued.
func (m *Mux) Destroy() error {
	if m.h.FD() <= 0 {
		return nil
	}

	if err := m.h.Destroy(); err != nil {
		return errors.Wrap(err, "unable to destroy VPC Mux")
	}

	return nil
}
