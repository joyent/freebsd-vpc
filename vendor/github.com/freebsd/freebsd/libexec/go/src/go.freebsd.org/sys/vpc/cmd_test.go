// Tests for vpc_ctl(2) commands
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

package vpc_test

import (
	"testing"

	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc"
)

func TestOpFlagType(t *testing.T) {
	tests := []struct {
		cmd       vpc.Cmd
		inBit     bool
		outBit    bool
		mutateBit bool
		privBit   bool
		objType   vpc.ObjType
		op        vpc.Op
	}{
		{
			cmd:       vpc.Cmd(0x00000000),
			outBit:    false,
			inBit:     false,
			mutateBit: false,
			privBit:   false,
			objType:   vpc.ObjType(0x00),
			op:        vpc.Op(0x0000),
		},
		{
			cmd:       vpc.Cmd(0xffffffff),
			outBit:    true,
			inBit:     true,
			mutateBit: true,
			privBit:   true,
			objType:   vpc.ObjType(0xff),
			op:        vpc.Op(0xffff),
		},
		{
			cmd:       vpc.Cmd(0x20010000),
			outBit:    false,
			inBit:     false,
			mutateBit: true,
			privBit:   false,
			objType:   vpc.ObjType(0x01),
			op:        vpc.Op(0x0000),
		},
		{
			cmd:       vpc.Cmd(0x50200000),
			outBit:    true,
			inBit:     false,
			mutateBit: false,
			privBit:   true,
			objType:   vpc.ObjType(0x20),
			op:        vpc.Op(0x0000),
		},
		{
			cmd:       vpc.Cmd(0xa0ff0000),
			outBit:    false,
			inBit:     true,
			mutateBit: true,
			privBit:   false,
			objType:   vpc.ObjType(0xff),
			op:        vpc.Op(0x0000),
		},
	}

	for i, test := range tests {
		if test.cmd.Mutate() != test.mutateBit {
			t.Errorf("[%d] Mutate wrong", i)
		}

		if test.cmd.Out() != test.outBit {
			t.Errorf("[%d] Out wrong", i)
		}

		if test.cmd.In() != test.inBit {
			t.Errorf("[%d] In wrong", i)
		}

		if test.cmd.ObjType() != test.objType {
			t.Errorf("[%d] ObjType wrong: 0x%04x 0x%04x", i, test.cmd.ObjType(), test.objType)
		}

		if test.cmd.Op() != test.op {
			t.Errorf("[%d] Op wrong", i)
		}

	}
}
