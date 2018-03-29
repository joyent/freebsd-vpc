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
	"bytes"
	"encoding/binary"
	"net"

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
	_OpMuxListenAddrGet      = vpc.Op(8)
)

// Template commands that can be passed to vpc.Ctl() with a valid VPC Mux
// Handle.
const (
	_MuxListenCmd             _MuxCmd = _MuxCmd(vpc.InBit|vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeMux)<<16)) | _MuxCmd(_OpMuxListen)
	_MuxListenAddrCmd         _MuxCmd = _MuxCmd(vpc.OutBit|(vpc.Cmd(vpc.ObjTypeMux)<<16)) | _MuxCmd(_OpMuxListenAddrGet)
	_MuxFTESetCmd             _MuxCmd = _MuxCmd(vpc.InBit|vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeMux)<<16)) | _MuxCmd(_OpMuxFTESet)
	_MuxFTEDelCmd             _MuxCmd = _MuxCmd(vpc.InBit|vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeMux)<<16)) | _MuxCmd(_OpMuxFTEDel)
	_MuxFTEListCmd            _MuxCmd = _MuxCmd(vpc.InBit|(vpc.Cmd(vpc.ObjTypeMux)<<16)) | _MuxCmd(_OpMuxFTEList)
	_MuxUnderlayConnectCmd    _MuxCmd = _MuxCmd(vpc.InBit|vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeMux)<<16)) | _MuxCmd(_OpMuxUnderlayConnect)
	_MuxUnderlayDisconnectCmd _MuxCmd = _MuxCmd(vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeMux)<<16)) | _MuxCmd(_OpMuxUnderlayDisconnect)
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

// ConnectedID returns the VPC ID of the connected interface to this VPC Mux.
func (m *Mux) ConnectedID() (id vpc.ID, err error) {
	// TODO(seanc@): Test to see make sure the descriptor has the mutate bit set.

	out := make([]byte, vpc.IDSize)
	if err := vpc.Ctl(m.h, vpc.Cmd(_MuxConnectedIDGetCmd), nil, out); err != nil {
		return vpc.ID{}, errors.Wrap(err, "unable to get VPC Mux's connected ID")
	}

	buf := bytes.NewReader(out[:])
	if err = binary.Read(buf, binary.LittleEndian, &id); err != nil {
		return vpc.ID{}, errors.Wrap(err, "failed to read VPC ID")
	}

	return id, nil
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

// Disconnect a VPC Mux from a VPC Interface.
func (m *Mux) Disconnect() error {
	// TODO(seanc@): Test to see make sure the descriptor has the mutate bit set.

	if err := vpc.Ctl(m.h, vpc.Cmd(_MuxUnderlayDisconnectCmd), nil, nil); err != nil {
		return errors.Wrap(err, "unable to disconnect VPC Mux to to VPC Interface")
	}

	return nil
}

// Listen instructs the VPC Mux to listen at the given address (host:port) for
// VPC Mux'ed traffic (RFC 7348 VXLAN encapsulated).
func (m *Mux) Listen(addr string) error {
	// TODO(seanc@): Test to see make sure the descriptor has the mutate bit set.

	if err := vpc.Ctl(m.h, vpc.Cmd(_MuxListenCmd), []byte(addr), nil); err != nil {
		return errors.Wrap(err, "unable to listen for VPC Mux traffic")
	}

	return nil
}

// ListenAddr returns the IP and port being used by this VPC Mux to listen for
// muxed traffic (RFC 7348 VXLAN encapsulated). If the Mux is not listening, it
// will return an empty string for the host and port.
func (m *Mux) ListenAddr() (host, port string, err error) {
	// TODO(seanc@): Test to see make sure the descriptor has the mutate bit set.
	const maxListenAddrSize = 128
	out := make([]byte, maxListenAddrSize)
	if err := vpc.Ctl(m.h, vpc.Cmd(_MuxListenAddrCmd), nil, out); err != nil {
		return "", "", errors.Wrap(err, "unable to get the listening address from the VPC Mux")
	}

	host, port, err = net.SplitHostPort(string(out))
	if err != nil {
		return "", "", errors.Wrap(err, "unable to find host/port in listening address")
	}

	return host, port, nil
}
