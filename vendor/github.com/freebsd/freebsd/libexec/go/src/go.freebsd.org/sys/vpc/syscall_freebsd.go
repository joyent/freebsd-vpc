// Go interface to VPC syscalls on FreeBSD.
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
	"fmt"
	"syscall"
	"unsafe"

	"github.com/pkg/errors"
)

const (
	// SysVPCOpen is the reserved syscall number for vpc_open(2)
	SysVPCOpen = 580

	// SysVPCCtl is the reserved syscall number for vpc_ctl(2)
	SysVPCCtl = 581
)

// Ctl manipulates the Handle based on the args
func Ctl(h *Handle, cmd Cmd, in []byte, out []byte) error {
	// TODO(seanc@): Potential concurrency optimization if we conditionalize the
	// type of lock based on the bits encoded in Cmd.
	h.lock.Lock()
	defer h.lock.Unlock()

	return ctl(h, cmd, in, out)
}

// Open obtains a VPC handle to a given object type.  Obtaining an open Handle
// affords no privilges beyond validating that an ID exists on this system.  In
// all other cases Open returns a handle to a resource.  If the id can not be
// found, Open returns ENOENT unless the Create flag is set in flags.  If the
// Create flag is set and the id is found, Open returns EEXIST.  If an invalid
// Flag is set, Open returns EINVAL.  If the HandleType is out of bounds, Open
// returns EOPNOTSUPP.  Returned Handles must have their information Commit()'ed
// in order for it to persist beyond the life of the Handle.
func Open(id ID, ht HandleType, flags OpenFlags) (h *Handle, err error) {
	h = &Handle{}

	// 580     AUE_VPC         NOSTD   { int vpc_open(const vpc_id_t *vpc_id, vpc_type_t obj_type, \
	//                                   vpc_flags_t flags); }
	r0, _, e1 := syscall.Syscall(SysVPCOpen, uintptr(unsafe.Pointer(&id)), uintptr(ht), uintptr(flags))
	h.fd = HandleFD(r0)
	if e1 != 0 {
		h.fd = HandleErrorFD
		return h, syscall.Errno(e1)
	}

	return h, nil
}


func ctl(h *Handle, cmd Cmd, in []byte, out []byte) error {
	// Implementation sanity checking
	switch {
	case cmd.In() && len(in) == 0:
		return errors.New("operation requires non-zero length input")
	case cmd.Out() && out == nil:
		return errors.New("operation requires non-nil output")
	}

	// 581     AUE_VPC         NOSTD   { int vpc_ctl(int vpcd, vpc_op_t op, size_t innbyte, \
	//                                     const void *in, size_t *outnbyte, void *out); }
	var r1 uintptr
	var e1 syscall.Errno
	switch {
	case len(in) == 0 && out == nil:
		r1, _, e1 = syscall.Syscall6(SysVPCCtl, uintptr(h.fd), uintptr(cmd),
			uintptr(0), uintptr(0),
			uintptr(0), uintptr(0))
	case len(in) != 0 && out != nil:
		sz := uint64(len(out))
		r1, _, e1 = syscall.Syscall6(SysVPCCtl, uintptr(h.fd), uintptr(cmd),
			uintptr(len(in)), uintptr(unsafe.Pointer(&in[0])),
			uintptr(unsafe.Pointer(&sz)), uintptr(unsafe.Pointer(&out[0])))
		out = out[:sz]
	case len(in) != 0 && out == nil:
		r1, _, e1 = syscall.Syscall6(SysVPCCtl, uintptr(h.fd), uintptr(cmd),
			uintptr(len(in)), uintptr(unsafe.Pointer(&in[0])),
			uintptr(0), uintptr(0))
	case len(in) == 0 && out != nil:
		sz := uint64(len(out))
		r1, _, e1 = syscall.Syscall6(SysVPCCtl, uintptr(h.fd), uintptr(cmd),
			uintptr(0), uintptr(0),
			uintptr(unsafe.Pointer(&sz)), uintptr(unsafe.Pointer(&out[0])))
		out = out[:sz]
	default:
		panic(fmt.Sprintf("invalid args to vpc.Ctl()\ncmd: %x\nin: %q\nout: %v", cmd, in, out))
	}
	if r1 != 0 {
		return syscall.Errno(e1)
	}

	return nil
}
