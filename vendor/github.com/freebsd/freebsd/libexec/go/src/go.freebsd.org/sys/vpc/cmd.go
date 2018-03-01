// Go interface to vpc_ctl(2) commands
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

	"github.com/rs/zerolog"
)

// Cmd is the operation being performed.  A Cmd encodes:
//
// 1. The direction of the arguments (in or out)
// 2. The object type (ObjType)
// 3. The per-object type operation being performed (Op).
type Cmd uint32

// Op is the action being performed on a given object type.
type Op uint16

// String returns a hex-encded string representing the op.  Individual Object
// Types should have their own Op type that satisfies the stringer method.
func (op Op) String() string {
	return fmt.Sprintf("op-0x%04x", uint16(op))
}

// Constants used to test or extract information from a command.
const (
	// Constants taken from: sys/sys/ioccom.h and sys/net/if_vpc.h
	//
	// PrivBit and MutateBit are capsicum rights encoded in the command.  If the
	// rights in the command don't match up with the rights stored in the Handle,
	// the operation will fail.  Encoding rights into the command and handle
	// allows privileges to be dropped.

	PrivBit   Cmd = 0x10000000
	MutateBit Cmd = 0x20000000
	OutBit    Cmd = 0x40000000
	InBit     Cmd = 0x80000000

	ObjTypeMask Cmd = 0xFF00FFFF
	OpMask      Cmd = 0xFFFF0000
)

// In returns true if the command requires input when passed to vpc.Ctl()
func (cmd Cmd) In() bool {
	return cmd&InBit != 0
}

// Mutate returns true if the command indicates that operation has a side effect
// and will change the state of the object attached to the VPC Handle.
func (cmd Cmd) Mutate() bool {
	return cmd&MutateBit != 0
}

// ObjType returns the encoded ObjType in the command.
func (cmd Cmd) ObjType() ObjType {
	return ObjType((cmd &^ ObjTypeMask) >> 16)
}

// Op returns the operation encoded in the command.
func (cmd Cmd) Op() Op {
	return Op(cmd &^ OpMask)
}

// Out returns true if the command requires an output argument when passed to
// vpc.Ctl().
func (cmd Cmd) Out() bool {
	return cmd&OutBit != 0
}

// Privileged returns true if the command indicates that the operation is
// privileged.
func (cmd Cmd) Privileged() bool {
	return cmd&PrivBit != 0
}

func (cmd Cmd) MarshalZerologObject(e *zerolog.Event) {
	e.Bool("in", cmd.In())
	e.Bool("out", cmd.Out())
	e.Bool("priv", cmd.Privileged())
	e.Bool("mutate", cmd.Mutate())
	e.Str("obj-type", cmd.ObjType().String())
	e.Str("obj-op", cmd.Op().String())
}
