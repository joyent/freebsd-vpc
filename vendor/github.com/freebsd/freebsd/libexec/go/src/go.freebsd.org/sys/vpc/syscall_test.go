// Test framework for VPC syscalls.
//
// SPDX-License-Identifier: BSD-2-Clause-FreeBSD
//
// Copyright (C) 2018 Sean Chittenden <seanc@joyent.com>
// Copyright (c) 2018 Joyent, Inc.
// All rights reserved.
//
// This software was developed by Sean Chittenden <seanc@FreeBSD.org> at Joyent,
// Inc.
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
	"encoding/binary"
	"reflect"
	"testing"
	"unsafe"

	"github.com/kylelemons/godebug/pretty"
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc"
)

func TestVPCIDBytes(t *testing.T) {
	tests := []struct {
		bin vpc.ID
		str string
	}{
		{
			bin: vpc.ID{},
			str: "00000000-0000-0000-0000-000000000000",
		},
		{
			bin: vpc.ID{
				TimeLow:     binary.LittleEndian.Uint32([]byte{0xa1, 0x0a, 0xbf, 0xcf}),
				TimeMid:     binary.LittleEndian.Uint16([]byte{0x1a, 0x6f}),
				TimeHi:      binary.LittleEndian.Uint16([]byte{0x11, 0xe8}),
				ClockSeqHi:  uint8(binary.LittleEndian.Uint16([]byte{0x81, 0x00})),
				ClockSeqLow: uint8(binary.LittleEndian.Uint16([]byte{0x77, 0x00})),
				Node:        [6]byte{0x0c, 0xc4, 0x7a, 0x6c, 0x7d, 0x1e},
			},
			str: "a10abfcf-1a6f-11e8-8177-0cc47a6c7d1e",
		},
	}

	for i, test := range tests {
		s := test.bin.String()
		if diff := pretty.Compare(s, test.str); diff != "" {
			t.Errorf("[%d] String VPC ID diff: (-got +want)\n%s", i, diff)
		}

		o, err := vpc.ParseID(s)
		if err != nil {
			t.Errorf("[%d] ParseID failed: %v", err)
		}

		if diff := pretty.Compare(o.String(), test.str); diff != "" {
			t.Errorf("[%d] round-trip VPC ID diff: (-got +want)\n%s", i, diff)
		}
	}
}

func Test_VPC_ID(t *testing.T) {
	origID := vpc.GenID()
	origIDStr := origID.String()
	if len(origIDStr) != 36 {
		t.Fatalf("ID wrong len")
	}

	parseID, err := vpc.ParseID(origIDStr)
	if err != nil {
		t.Fatalf("unable to parse %q: %v", origIDStr, err)
	}

	if !reflect.DeepEqual(origID, parseID) {
		t.Fatalf("parsed bytes don't match: %v %v", origID, parseID)
	}

	if origID.String() != parseID.String() {
		t.Fatalf("string IDs don't match: %q %q", origID.String(), parseID.String())
	}
}

func Test_VPC_ParseID(t *testing.T) {
	tests := []struct {
		idStr string
		ok    bool
	}{
		{
			idStr: "183dddcc-2f8a-85d7-2d7c-3c6a14e22d5d",
			ok:    true,
		},
		{
			idStr: "183dddcc-2f8a-85d7-2d7c-316a14e22d5d",
			ok:    false,
		},
		{
			idStr: vpc.GenID().String(),
			ok:    true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.idStr, func(t *testing.T) {
			t.Parallel()

			if len(test.idStr) != 36 {
				t.Fatalf("ID wrong len")
			}

			parseID, err := vpc.ParseID(test.idStr)
			switch {
			case err == nil && test.ok == true:
			case err != nil && test.ok == false:
				return // ok
			case err != nil && test.ok == true:
				t.Fatalf("unable to parse %q: %v", test.idStr, err)
			case err == nil && test.ok == false:
				t.Fatalf("expected failure for %q", test.idStr)
			}

			if test.idStr != parseID.String() {
				t.Fatalf("string IDs don't match: %q %q", test.idStr, parseID.String())
			}
		})
	}
}

func TestOpenFlags(t *testing.T) {
	if unsafe.Sizeof(vpc.OpenFlags(0)) != 8 {
		t.Fatal("open flags must be 8 bytes")
	}
}

func TestObjType(t *testing.T) {
	if unsafe.Sizeof(vpc.ObjType(0)) != 1 {
		t.Fatal("open flags must be 1 byte")
	}
}

func TestSizeofID(t *testing.T) {
	if vpc.IDSize != 16 {
		t.Errorf("size of vpc.ID is a UUID and expected to be 16B,not %d", vpc.IDSize)
	}

	dynSize := unsafe.Sizeof(vpc.ID{})
	if dynSize != vpc.IDSize {
		t.Errorf("size of vpc.ID changed from %d to %d, ABI mismatch with the kernel guaranteed", vpc.IDSize, dynSize)
	}
}
