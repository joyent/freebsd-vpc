// Tests for VPC Handles
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
	"math/rand"
	"testing"
	"unsafe"

	"github.com/kylelemons/godebug/pretty"
	"github.com/sean-/seed"
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc"
)

func init() {
	seed.MustInit()
}

type vpcHandleSlice []*vpc.Handle

// Shuffle pseudo-randomizes the order of elements using Fisher-Yates shuffle
// function.
func (s vpcHandleSlice) Shuffle() {
	for i := len(s) - 1; i > 0; i-- {
		j := rand.Int31n(int32(i + 1))
		s[i], s[j] = s[j], s[i]
	}
}

func TestHandle(t *testing.T) {
	if unsafe.Sizeof(vpc.HandleType(0)) != 8 {
		t.Fatal("size must be 8 bytes")
	}
}

func TestHandleTypeVersion(t *testing.T) {
	tests := []struct {
		in         vpc.HandleType
		outVersion vpc.HandleVersion
		outObjType vpc.ObjType
	}{
		{
			in:         0x0000000000000000,
			outVersion: 0x0,
			outObjType: vpc.ObjTypeInvalid,
		},
		{
			in:         0xe001000000000000,
			outVersion: 14,
			outObjType: vpc.ObjTypeSwitch,
		},
		{
			in:         0x3002000000000000,
			outVersion: 3,
			outObjType: vpc.ObjTypeSwitchPort,
		},
		{
			in:         0x4003000000000000,
			outVersion: 4,
			outObjType: vpc.ObjTypeRouter,
		},
		{
			in:         0x5004000000000000,
			outVersion: 5,
			outObjType: vpc.ObjTypeNAT,
		},
		{
			in:         0x6005000000000000,
			outVersion: 6,
			outObjType: vpc.ObjTypeLinkVPC,
		},
		{
			in:         0x7006000000000000,
			outVersion: 7,
			outObjType: vpc.ObjTypeNICVM,
		},
		{
			in:         0x8007000000000000,
			outVersion: 8,
			outObjType: vpc.ObjTypeMgmt,
		},
		{
			in:         0x9008000000000000,
			outVersion: 9,
			outObjType: vpc.ObjTypeLinkEth,
		},
		{
			in:         0xa009000000000000,
			outVersion: 10,
			outObjType: vpc.ObjTypeMeta,
		},
		{
			in:         0xb00a000000000000,
			outVersion: 11,
			outObjType: vpc.ObjTypeAny,
		},
	}

	// Test Version bits
	for i, test := range tests {
		ht := vpc.HandleType(test.in)

		if diff := pretty.Compare(ht.Version(), test.outVersion); diff != "" {
			t.Errorf("[%d] %#v: Version diff: (-got +want)\n%s", i, ht.Version(), diff)
		}

		{
			oldVer := ht.Version()
			newHandle, err := ht.SetVersion(oldVer + 1)
			if err != nil {
				t.Errorf("[%d] %#v: set failed: %v", i, ht, err)
			}
			if oldVer == newHandle.Version() {
				t.Fatalf("[%d] %#v: cmp failed", i, ht)
			}
			ht = newHandle
		}

		if diff := pretty.Compare(ht.Version(), test.outVersion+1); diff != "" {
			t.Errorf("[%d] %#v: Version diff: (-got +want)\n%s", i, ht.Version(), diff)
		}
	}

	// Test ObjType bits
	for i, test := range tests {
		ht := vpc.HandleType(test.in)

		if diff := pretty.Compare(ht.ObjType(), test.outObjType); diff != "" {
			t.Errorf("[%d] %#v: ObjType diff: (-got +want)\n%s", i, ht.ObjType(), diff)
		}

		// Enumerate all known object types to catch their accidental removal in the
		// future.
		objTypes := []vpc.ObjType{
			vpc.ObjTypeInvalid,
			vpc.ObjTypeSwitch,
			vpc.ObjTypeSwitchPort,
			vpc.ObjTypeRouter,
			vpc.ObjTypeNAT,
			vpc.ObjTypeLinkVPC,
			vpc.ObjTypeNICVM,
			vpc.ObjTypeMgmt,
			vpc.ObjTypeLinkEth,
		}

		origHandleType := ht
		for _, objType := range objTypes {
			ht := origHandleType
			ver := ht.Version()
			newType, err := ht.SetObjType(objType)
			if err != nil {
				t.Errorf("[%d] %#v: set objtype failed: %v", i, ht, err)
			}

			if diff := pretty.Compare(newType.ObjType(), objType); diff != "" {
				t.Errorf("[%d] %#v: set objtype failed: %v", i, ht, err)
			}

			if diff := pretty.Compare(newType.Version(), ver); diff != "" {
				t.Errorf("[%d] %#v: set objtype nicked version: %v", i, ht, err)
			}
		}
	}
}
