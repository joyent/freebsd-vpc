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

package vmnic

import (
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
	"go.freebsd.org/sys/vpc"
)

// _VMNICCmd is the encoded type of operations that can be performed on a VM
// NIC.
type _VMNICCmd vpc.Cmd

// _VMNICCmdSetArgType is the value used by a VM NIC set operation.
type _VMNICSetOpArgType uint64

// Ops that can be encoded into a vpc.Cmd
const (
	_OpInvalid    = vpc.Op(0)
	_OpNQueuesGet = vpc.Op(1)
	_OpNQueuesSet = vpc.Op(2)
	_             = vpc.Op(3) // unused
	_             = vpc.Op(4) // unused
	_             = vpc.Op(5) // unused
	_             = vpc.Op(6) // unused
	// _OpAttach     = vpc.Op(7) // bhyve SPI
	// _OpMSIX       = vpc.Op(8) // kvirtio SPI
	_OpFreeze   = vpc.Op(9)
	_OpUnfreeze = vpc.Op(10)
)

// Cmds that can be sent to vpc.Ctl()
const (
	_NQueuesGetCmd _VMNICCmd = _VMNICCmd(vpc.OutBit|(vpc.Cmd(vpc.ObjTypeNICVM)<<16)) | _VMNICCmd(_OpNQueuesGet)
	_NQueuesSetCmd _VMNICCmd = _VMNICCmd(vpc.InBit|vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeNICVM)<<16)) | _VMNICCmd(_OpNQueuesSet)
	_FreezeCmd     _VMNICCmd = _VMNICCmd(vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeNICVM)<<16)) | _VMNICCmd(_OpFreeze)
	_UnfreezeCmd   _VMNICCmd = _VMNICCmd(vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeNICVM)<<16)) | _VMNICCmd(_OpUnfreeze)
)

// Close closes the VPC Handle descriptor.  Created VM NICs will not be
// destroyed when the VMNIC is closed if the VM NIC has been Committed.
func (vmn *VMNIC) Close() error {
	if vmn.h.FD() <= 0 {
		return nil
	}

	if err := vmn.h.Close(); err != nil {
		return errors.Wrap(err, "unable to close VPC handle")
	}

	return nil
}

// Commit increments the refcount of the VM NIC in order to ensure the VM NIC
// lives beyond the life of the current process and is not automatically cleaned
// up when the VMNIC is closed.
func (vmn *VMNIC) Commit() error {
	if vmn.h.FD() <= 0 {
		return nil
	}

	if err := vmn.h.Commit(); err != nil {
		return errors.Wrap(err, "unable to commit VM NIC")
	}

	return nil
}

// Destroy decrements the refcount of the VM NIC in destroy the the VM NIC when
// the VPC Handle is closed.
func (vmn *VMNIC) Destroy() error {
	if vmn.h.FD() <= 0 {
		return nil
	}

	if err := vmn.h.Destroy(); err != nil {
		return errors.Wrap(err, "unable to destroy VM NIC")
	}

	return nil
}

// Freeze freezes the VMNIC so it can be plugged into a VPC Switch Port.
func (vmn *VMNIC) Freeze(enable bool) error {
	cmd := _UnfreezeCmd
	cmdStr := "unfreeze"
	if enable {
		cmd = _FreezeCmd
		cmdStr = "freeze"
	}

	if err := vpc.Ctl(vmn.h, vpc.Cmd(cmd), nil, nil); err != nil {
		return errors.Wrapf(err, "unable to %s VM NIC", cmdStr)
	}

	return nil
}

// NQueuesGet returns the number of queues assigned to this VMNIC.
func (vmn *VMNIC) NQueuesGet() (uint16, error) {
	out := make([]byte, binary.MaxVarintLen64)
	if err := vpc.Ctl(vmn.h, vpc.Cmd(_NQueuesGetCmd), nil, out); err != nil {
		return 0, errors.Wrap(err, "unable to get the number of hardware queues from VMNIC")
	}

	numQueues, n := binary.Uvarint(out)
	if n > 0 && n <= 2 {
		return uint16(numQueues), nil
	}

	panic(fmt.Sprintf("invariant: num queues too big for kernel interface output (want/got: 2/%d", n))
}

// NQueuesSet sets the number of queues for this VMNIC.
func (vmn *VMNIC) NQueuesSet(numQueues uint16) error {
	in := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(in, uint64(numQueues))
	if n < 2 {
		in = in[:2]
	} else {
		panic(fmt.Sprintf("invariant: num queuese size too big for kernel interface input (want/got: 2/%d", n))
	}

	if err := vpc.Ctl(vmn.h, vpc.Cmd(_NQueuesSetCmd), in, nil); err != nil {
		return errors.Wrap(err, "unable to set the number of hardware queues for VMNIC")
	}

	return nil
}
