// Go interface to L2 Link objects.
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

package l2link

import (
	"github.com/pkg/errors"
	"go.freebsd.org/sys/vpc"
)

// _L2Cmd is the encoded type of operations that can be performed on a L2 Link.
type _L2Cmd vpc.Cmd

// _L2LinkCmdSetArgType is the value used by a L2 Link set operation.
type _L2LinkSetOpArgType uint64

const (
	// Bits for input
	_DownBit _L2LinkSetOpArgType = 0x00000000
	_UpBit   _L2LinkSetOpArgType = 0x00000001
)

// Ops that can be encoded into a vpc.Cmd
const (
	_OpInvalid = vpc.Op(0)
	_OpAttach  = vpc.Op(1)

	// _OpReset         = vpc.Op(7)

	_AttachCmd _L2Cmd = _L2Cmd(vpc.InBit|vpc.PrivBit|vpc.MutateBit|(vpc.Cmd(vpc.ObjTypeLinkL2)<<16)) | _L2Cmd(_OpAttach)
)

// Template commands that can be passed to vpc.Ctl() with a valid VM NIC
// Handle.
// var (
// 	_ResetCmd _L2LinkCmd
// )

// Attach attaches the named physical device or cloned interface to this VPC L2
// Link.  The name of the device must be specified in the L2Link Config and
// passed in at Create time.
func (l2 *L2Link) Attach() error {
	// TODO(seanc@): Test to see make sure the descriptor has the mutate bit set.

	if err := vpc.Ctl(l2.h, vpc.Cmd(_AttachCmd), []byte(l2.name), nil); err != nil {
		return errors.Wrap(err, "unable to attach VPC L2 Link to a physical NIC")
	}

	return nil
}

// Close closes the VPC Handle.  Created L2 Links will not be destroyed when the
// L2Link is closed if the L2 Links has been Committed.
func (l2 *L2Link) Close() error {
	if l2.h.FD() <= 0 {
		return nil
	}

	if err := l2.h.Close(); err != nil {
		return errors.Wrap(err, "unable to close L2 Link VPC handle")
	}

	return nil
}

// Commit increments the refcount of the L2 Link in order to ensure the L2 Link
// lives beyond the life of the current process and is not automatically cleaned
// up when the L2Link is closed.
func (l2 *L2Link) Commit() error {
	if l2.h.FD() <= 0 {
		return errors.Errorf("unable to commit VPC L2 Link handle with an empty descriptor")
	}

	if err := l2.h.Commit(); err != nil {
		return errors.Wrap(err, "unable to commit VPC L2 Link")
	}

	return nil
}

// Destroy decrements the refcount of the L2 Link in destroy the the L2 Link
// when the VPC Handle is closed.
func (l2 *L2Link) Destroy() error {
	if l2.h.FD() <= 0 {
		return nil
	}

	if err := l2.h.Destroy(); err != nil {
		return errors.Wrap(err, "unable to destroy L2 Link")
	}

	return nil
}
