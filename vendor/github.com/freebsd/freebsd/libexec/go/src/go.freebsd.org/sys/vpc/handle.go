// Go interface for VPC Handles.
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

package vpc

import (
	"encoding/binary"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// HandleFD is the descriptor number associated with an opened VPC Object.
type HandleFD int

// Handle is a handle to the actual descriptor
type Handle struct {
	lock sync.RWMutex
	fd   HandleFD
}

func (h Handle) MarshalZerologObject(e *zerolog.Event) {
	e.Int("fd", int(h.fd))
}

const (
	// HandleErrorFD is the value returned when an error occurrs during a call to
	// Open.
	HandleErrorFD HandleFD = -1

	// HandleClosedFD is the value used to indicate a Handle has been closed.
	HandleClosedFD HandleFD = -2

	errVersion HandleType = 0x1
)

// HandleVersion is the version number of the VPC API and controls the ABI used
// to talk with a VPC Handle.
type HandleVersion uint64

func (v HandleVersion) MarshalZerologObject(e *zerolog.Event) {
	e.Int64("version", int64(v))
}

// HandleTypeInput is passed to the constructor NewHandleType
type HandleTypeInput struct {
	Version  HandleVersion
	Type     ObjType
	Writable bool
}

// HandleType is the Object Type.  In sys/amd64/vmm/net/vmmnet.c this is
// defined as:
//
//    typedef struct {
//      uint64_t vht_version:4;
//      uint64_t vht_pad1:4;
//      uint64_t vht_obj_type:8;
//      uint64_t vht_pad2:48;
//    } vpc_handle_type_t;
type HandleType uint64

// NewHandleType constructs a new HandleType
func NewHandleType(cfg HandleTypeInput) (ht HandleType, err error) {
	if ht, err = ht.SetVersion(cfg.Version); err != nil {
		return errVersion, err
	}

	if ht, err = ht.SetObjType(cfg.Type); err != nil {
		return errVersion, err
	}

	return ht, err
}

func (ht HandleType) MarshalZerologObject(e *zerolog.Event) {
	e.Object("version", ht.Version()).
		Object("obj-type", ht.ObjType())
}

const (
	objTypeMask HandleType = 0x00ff000000000000
	versionMask HandleType = 0xf000000000000000
)

// Version returns the HandleVersion being opened
func (t HandleType) Version() HandleVersion {
	return HandleVersion(t >> (64 - 4))
}

// SetVersion returns a new HandleType with the version encoded in the result.
func (t HandleType) SetVersion(ver HandleVersion) (HandleType, error) {
	switch {
	case ver > ((2 << 4) - 1):
		return errVersion, errors.New("API version too large")
	}

	// clear version
	tu := uint64(t)
	tu = tu &^ uint64(versionMask)

	// set version
	uVer := uint64(ver)
	uVer = uVer << (64 - 4)
	return HandleType(tu | uVer), nil
}

// ObjType returns the ObjType from a given HandleType
func (t HandleType) ObjType() ObjType {
	t &= objTypeMask
	t = t >> (64 - 8 - 8)
	return ObjType(t)
}

// SetObjType encodes the ObjType into a copy of the HandleType receiver and
// returns a new HandleType with the ObjType encoded.
func (t HandleType) SetObjType(objType ObjType) (HandleType, error) {
	// clear version
	tu := uint64(t)
	tu = tu &^ uint64(objTypeMask)

	// set ObjType
	uVer := uint64(objType)
	uVer = uVer << (64 - 8 - 8)
	return HandleType(tu | uVer), nil
}

// Operations that can be applied to all VPC Object types
const (
	// Note: sorted by op number.  These commands are prefixed by the type of
	// descriptor required to run the op.  Meta ops can be run by any descriptor.
	// The Mgmt commands can only be run by the mgmt command.

	_MetaDestroyOp = Op(0x0001)
	_MetaTypeGetOp = Op(0x0002)
	_MetaCommitOp  = Op(0x0003)
	_MetaMACSetOp  = Op(0x0004)
	_MetaMACGetOp  = Op(0x0005)
	_MetaMTUSetOp  = Op(0x0006)
	_MetaMTUGetOp  = Op(0x0007)
	_MetaGetIDOp   = Op(0x0008)
	_MgmtCountOp   = Op(0x0009)

	// Management Commands
	_CountCmd = PrivBit | (Cmd(ObjTypeMgmt) << 16) | Cmd(_MgmtCountOp)

	// Meta commands
	_CommitCmd  = PrivBit | MutateBit | (Cmd(ObjTypeMeta) << 16) | Cmd(_MetaCommitOp)
	_DestroyCmd = PrivBit | MutateBit | (Cmd(ObjTypeMeta) << 16) | Cmd(_MetaDestroyOp)
	_GetIDCmd   = (Cmd(ObjTypeMeta) << 16) | Cmd(_MetaGetIDOp)
	_MACGetCmd  = (Cmd(ObjTypeMeta) << 16) | Cmd(_MetaMACGetOp)
	_MACSetCmd  = PrivBit | MutateBit | (Cmd(ObjTypeMeta) << 16) | Cmd(_MetaMACSetOp)
	_MTUGetCmd  = (Cmd(ObjTypeMeta) << 16) | Cmd(_MetaMTUGetOp)
	_MTUSetCmd  = PrivBit | MutateBit | (Cmd(ObjTypeMeta) << 16) | Cmd(_MetaMTUSetOp)
	_TypeCmd    = OutBit | (Cmd(ObjTypeMeta) << 16) | Cmd(_MetaTypeGetOp)
)

// Commit increments the refcount on the object referrenced by this VPC Handle.
// Commit is used to ensure that the life of the referred VPC object outlives
// the current process with the open VPC Handle.
func (h *Handle) Commit() error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if err := ctl(h, _CommitCmd, nil, nil); err != nil {
		return errors.Wrap(err, "unable to commit VPC object")
	}

	return nil
}

// Destroy decrements the refcount on the object referrenced by this VPC Handle.
// Destroy is used to terminate the life of the referred VPC object so that the
// VPC Object's resources are cleaned up when the Handle is closed.
func (h *Handle) Destroy() error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if err := ctl(h, _DestroyCmd, nil, nil); err != nil {
		return errors.Wrap(err, "unable to destroy VPC object")
	}

	return nil
}

// FD returns the integer Unix file descriptor referencing the open file. The
// file descriptor is valid only until h.Close is called.
func (h *Handle) FD() HandleFD {
	h.lock.RLock()
	defer h.lock.RUnlock()

	return h.fd
}

// Type returns the VPC Object Type used by this handle.
func (h *Handle) Type() (ObjType, error) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	out := make([]byte, binary.MaxVarintLen64)
	if err := ctl(h, _TypeCmd, nil, nil); err != nil {
		return ObjTypeInvalid, errors.Wrap(err, "unable to destroy VPC object")
	}

	objType, n := binary.Uvarint(out)
	if n > 0 && n <= 2 {
		return ObjType(objType), nil
	}

	panic(fmt.Sprintf("invariant: obj type too big for kernel interface output (want/got: 2/%d", n))
}
