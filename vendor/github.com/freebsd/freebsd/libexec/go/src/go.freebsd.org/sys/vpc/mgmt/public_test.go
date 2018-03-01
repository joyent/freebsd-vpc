// Test VPC Management handles.
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

package mgmt_test

import (
	"testing"

	"github.com/sean-/seed"
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc"
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/mgmt"
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/vpctest"
)

func init() {
	seed.MustInit()
}

// TestMgmt_OpenClose is intended to verify the basic lifecycle of a management
// handle.
func TestMgmt_OpenClose(t *testing.T) {
	genIDPtr := func() *vpc.ID {
		id := vpc.GenID()
		return &id
	}

	tests := []struct {
		cfg mgmt.Config
	}{
		{ // Empty test
		},
		{
			cfg: mgmt.Config{
				ID:        genIDPtr(),
				Writeable: false,
			},
		},
		{
			cfg: mgmt.Config{
				ID:        genIDPtr(),
				Writeable: true,
			},
		},
	}

	existingIfaces, err := vpctest.GetAllInterfaces()
	if err != nil {
		t.Fatalf("unable to get all interfaces")
	}

	for i, test := range tests {
		// anon func used to trigger execution of defers
		func() { // New/Close
			mgr, err := mgmt.New(&test.cfg)
			if err != nil {
				t.Fatalf("[%d] unable to create new VPC Management Handle: %v", i, err)
			}
			defer mgr.Close()

			{ // Get the before/after the New
				ifacesAfterCreate, err := vpctest.GetAllInterfaces()
				if err != nil {
					t.Fatalf("[%d] unable to get all interfaces")
				}
				oldIfaces, newIfaces, _ := existingIfaces.Difference(ifacesAfterCreate)
				if len(oldIfaces) != 0 || len(newIfaces) != 0 {
					t.Fatalf("[%d] no interfaces should have been created", i)
				}
			}

			if err := mgr.Close(); err != nil {
				t.Fatalf("[%d] unable to close VPC Management Handle: %v", i, err)
			}
		}() // anon func
	} // for
}

// TestMgmt_CountTypes verifies that we can count the number of objects of
// various types with a VPC management descriptor.
func TestMgmt_CountTypes(t *testing.T) {
	tests := []struct {
		objType vpc.ObjType
		count   int64
	}{
		{
			objType: vpc.ObjTypeSwitch,
			count:   -1,
		},
	}

	for i, test := range tests {
		test := test
		t.Run(test.objType.String(), func(t *testing.T) {
			mgr, err := mgmt.New(nil)
			if err != nil {
				t.Fatalf("[%d] unable to create new VPC Management handle: %v", i, err)
			}
			defer mgr.Close()

			count, err := mgr.CountType(test.objType)
			if err != nil {
				t.Fatalf("[%d] unable to get a count of %s VPC objects: %v", i, test.objType, err)
			}

			switch {
			case test.count == -1 && count >= 0:
			case test.count >= 0 && int64(count) != test.count:
				t.Errorf("[%d] wrong number of %s VPC objects", i, test.objType)
			}
		})
	}
}

// TestMgmt_GetAllIDs verifies that we can get the IDs for a given type.
func TestMgmt_GetAllIDs(t *testing.T) {
	for i, objType := range vpc.ObjTypes() {
		objType := objType
		t.Run(objType.String(), func(t *testing.T) {
			mgr, err := mgmt.New(nil)
			if err != nil {
				t.Fatalf("[%d] unable to create new VPC Management handle: %v", i, err)
			}
			defer mgr.Close()

			ids, err := mgr.GetAllIDs(objType)
			if err != nil {
				t.Fatalf("[%d] unable to get IDs for %s VPC objects: %v", i, objType, err)
			}

			_ = ids // Not sure what to do with these UUIDs
		})
	}
}
