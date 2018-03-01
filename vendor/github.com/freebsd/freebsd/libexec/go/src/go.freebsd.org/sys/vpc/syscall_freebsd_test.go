// Go interface to OS-independent VPC syscalls.
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
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"runtime"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/rs/zerolog/log"
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc"
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc/vpctest"
)

// TestVPCCreateOpenClose performs a serialized and strict test of the Create,
// Open, Close semantics.  Most per-interface tests should go into
// TestVPCCreateOpenCloseParallel.  Serialization is required because of the
// use of vpctest.GetAllInterfaces().
func TestVPCCreateOpenClose(t *testing.T) {
	vpcsw0ID := vpc.GenID()

	existingIfaces, err := vpctest.GetAllInterfaces()
	if err != nil {
		t.Fatalf("unable to get all interfaces")
	}

	ht, err := vpc.NewHandleType(vpc.HandleTypeInput{
		Version: 1,
		Type:    vpc.ObjTypeSwitch,
	})
	if err != nil {
		t.Fatalf("unable to construct a HandleType: %v", err)
	}

	vpcsw0CHandle, err := vpc.Open(vpcsw0ID, ht, vpc.FlagCreate)
	if err != nil {
		t.Fatalf("vpc_open(2) failed: %v", err)
	}

	// NOTE(seanc@): This could conceivably be 0 if we've closed stdin before this
	// test runs.
	if vpcsw0CHandle.FD() == 0 {
		t.Errorf("vpc_open(2) return an FD of 0")
	}

	// Get the before/after
	ifacesAfterCreate, err := vpctest.GetAllInterfaces()
	if err != nil {
		t.Fatalf("unable to get all interfaces")
	}
	_, newIfaces1, _ := existingIfaces.Difference(ifacesAfterCreate)
	if len(newIfaces1) != 1 {
		t.Fatalf("one interface should have been added")
	}

	// For the sake of testing, call this `vpcsw0` even though the name may be
	// different.
	vpcsw0 := newIfaces1.First()
	log.Debug().Str("name", vpcsw0.Name).Msg("created vpcsw0")
	if bytes.Compare(vpcsw0.HardwareAddr[:], net.HardwareAddr{}[:]) == 0 {
		t.Fatalf("%s hardware address is uninitialized, passed in %q", vpcsw0.Name, vpcsw0ID)
	}

	if bytes.Compare(vpcsw0.HardwareAddr, vpcsw0ID.Node[:]) != 0 {
		t.Fatalf("MAC address doesn't match Node portion of ID:\ngot: %v\nwant: %v\nvpcid: %q", vpcsw0.HardwareAddr, vpcsw0ID.Node[:], vpcsw0ID)
	}

	if diff := pretty.Compare(vpcsw0.HardwareAddr, vpcsw0ID.Node[:]); diff != "" {
		t.Fatalf("MAC address doesn't match: (-got, +want)\n%s\ngot: %q\nwant: %q", diff, vpcsw0.HardwareAddr, vpcsw0ID.Node[:])
	}

	log.Debug().Str("VPC ID", vpcsw0ID.String()).Msg("Opening vpcsw0")
	vpcsw0OHandle, err := vpc.Open(vpcsw0ID, ht, vpc.FlagOpen)
	if err != nil {
		t.Fatalf("vpc_open(2) failed: %v", err)
	}
	defer func() {
		if err := vpcsw0OHandle.Close(); err != nil {
			t.Fatalf("unable to close(2) VPC Handle : %v", err)
		}
	}()

	if vpcsw0OHandle == vpcsw0CHandle {
		t.Errorf("vpc_open(2) open and create FDs are identical")
	}

	{
		// Get a new before/after: there should be no change
		ifacesAfterOpen, err := vpctest.GetAllInterfaces()
		if err != nil {
			t.Fatalf("unable to get all interfaces")
		}
		_, newIfaces1a, _ := ifacesAfterCreate.Difference(ifacesAfterOpen)
		if len(newIfaces1a) != 0 {
			t.Fatalf("no new interfaces should have been added")
		}

	}

	vpcsw1ID := vpc.GenID()

	log.Debug().Msg("creating vpcsw1")
	vpcsw1CreateFD, err := vpc.Open(vpcsw1ID, ht, vpc.FlagCreate)
	if err != nil {
		t.Fatalf("vpc_open(2) failed: %v", err)
	}
	defer func() {
		if err := vpcsw1CreateFD.Close(); err != nil {
			t.Fatalf("unable to close(2) vpcsw1CreateFD VPC Handle : %v", err)
		}
	}()

	log.Debug().Int("vpcsw0CreateFD", int(vpcsw0CHandle.FD())).Msg("closing vpcsw0 create")
	if err := vpcsw0CHandle.Close(); err != nil {
		t.Fatalf("unable to close(2) VPC Handle : %v", err)
	}
	if vpcsw0CHandle.FD() != vpc.HandleClosedFD {
		t.Fatalf("handle set to wrong value in vpc.Close()")
	}

	if err := vpcsw0CHandle.Close(); err != nil {
		t.Fatalf("unable to close(2) VPC Handle : %v", err)
	}

	// TODO(seanc@): programmatically verify that vpcsw0 is still present
	//time.Sleep(30 * time.Second)

	log.Debug().Msg("closing vpcsw0 open")
	if err := vpcsw0OHandle.Close(); err != nil {
		t.Fatalf("unable to close(2) VPC Handle : %v", err)
	}

	// TODO(seanc@): programmatically verify that vpcsw0 disappeared after the
	// openfd was closed
	//time.Sleep(30 * time.Second)

	log.Debug().Msg("closing vpcsw1 open")
	if err := vpcsw1CreateFD.Close(); err != nil {
		t.Fatalf("unable to close(2) VPC Handle : %v", err)
	}

	// TODO(seanc@): programmatically verify that vpcsw1 disappeared
	//time.Sleep(30 * time.Second)
}

// TestVPCCreateOpenCloseParallel tests the creation of switches in parallel.
// This is the preferred way of testing the characteristics of individual
// interfaces.
func TestVPCCreateOpenCloseParallel(t *testing.T) {
	maxNumOpensPerSwitch := 32
	cpuScalingFactor := 16
	if testing.Short() {
		cpuScalingFactor = 1
		maxNumOpensPerSwitch = 1
	}

	type _TestCase struct {
		name string
		id   vpc.ID
	}
	testCases := make([]_TestCase, runtime.NumCPU()*cpuScalingFactor)
	for i := range testCases {
		testCases[i] = _TestCase{
			name: fmt.Sprintf("id%d", i),
			id:   vpc.GenID(),
		}
	}

	ht, err := vpc.NewHandleType(vpc.HandleTypeInput{
		Version: 1,
		Type:    vpc.ObjTypeSwitch,
	})
	if err != nil {
		t.Fatalf("unable to construct a HandleType: %v", err)
	}

	log.Debug().Int("num test cases", len(testCases)).Msg("")
	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			vpcswID := vpc.GenID()

			vpcswCHandle, err := vpc.Open(vpcswID, ht, vpc.FlagCreate)
			if err != nil {
				t.Fatalf("vpc_open(2) failed: %v", err)
			}
			defer func(h *vpc.Handle) {
				if err := h.Close(); err != nil {
					t.Fatalf("unable to close(2) VPC Handle : %v c%v", err, h)
				}
			}(vpcswCHandle)

			// NOTE(seanc@): This could conceivably be 0 if we've closed stdin before
			// this test runs.
			if vpcswCHandle.FD() == 0 {
				t.Fatalf("vpc_open(2) return an FD of 0")
			}

			allInterfaces, err := vpctest.GetAllInterfaces()
			if err != nil {
				t.Fatalf("unable to get all interfaces: %v", err)
			}

			// For the sake of testing, call this `vpcsw0` even though the name may be
			// different.
			vpcsw0, err := allInterfaces.FindMAC(net.HardwareAddr(vpcswID.Node[:]))
			if err != nil {
				t.Fatalf("unable to find self: %v", err)
			}

			if bytes.Compare(vpcsw0.HardwareAddr[:], net.HardwareAddr{}[:]) == 0 {
				t.Fatalf("%s hardware address is uninitialized, passed in %q", vpcsw0.Name, vpcswID)
			}

			if bytes.Compare(vpcsw0.HardwareAddr, vpcswID.Node[:]) != 0 {
				t.Fatalf("MAC address doesn't match Node portion of ID:\ngot: %v\nwant: %v\nvpcid: %q", vpcsw0.HardwareAddr, vpcswID.Node[:], vpcswID)
			}

			if diff := pretty.Compare(vpcsw0.HardwareAddr, vpcswID.Node[:]); diff != "" {
				t.Fatalf("MAC address doesn't match: (-got, +want)\n%s\ngot: %q\nwant: %q", diff, vpcsw0.HardwareAddr, vpcswID.Node[:])
			}

			numVPCSwitchOpens := rand.Intn(maxNumOpensPerSwitch)

			openHandles := make(vpcHandleSlice, numVPCSwitchOpens, numVPCSwitchOpens+1)
			for i := 0; i < numVPCSwitchOpens; i++ {
				// Open a handle to the same switch and add it to the list
				vpcswOpenFD, err := vpc.Open(vpcswID, ht, vpc.FlagOpen)
				if err != nil {
					t.Fatalf("vpc_open(2) failed: %v", err)
				}
				defer func(h *vpc.Handle) {
					if err := h.Close(); err != nil {
						t.Fatalf("unable to close(2) VPC Handle : %v o%v", err, h)
					}
				}(vpcswOpenFD)

				openHandles[i] = vpcswOpenFD
			}

			openHandles = append(openHandles, vpcswCHandle)
			openHandles.Shuffle()

			for i := range openHandles {
				if err := openHandles[i].Close(); err != nil {
					t.Fatalf("unable to close(2) VPC Handle : %v b%v", err, openHandles[i])
				}
			}
		})
	}
}

func TestVPCCreateCommitDestroyClose(t *testing.T) {
	vpcsw0ID := vpc.GenID()

	ht, err := vpc.NewHandleType(vpc.HandleTypeInput{
		Version: 1,
		Type:    vpc.ObjTypeSwitch,
	})
	if err != nil {
		t.Fatalf("unable to construct a HandleType: %v", err)
	}

	log.Debug().Msg("creating vpcsw0")
	vpcsw0CHandle, err := vpc.Open(vpcsw0ID, ht, vpc.FlagCreate|vpc.FlagWrite)
	if err != nil {
		t.Fatalf("vpc_open(2) failed: %v", err)
	}

	// NOTE(seanc@): This could conceivably be 0 if we've closed stdin before this
	// test runs.
	if vpcsw0CHandle.FD() == 0 {
		t.Errorf("vpc_open(2) return an FD of 0")
	}

	if err := vpcsw0CHandle.Commit(); err != nil {
		t.Fatalf("vpc_open(2) unable to commit: %v", err)
	}

	if err := vpcsw0CHandle.Destroy(); err != nil {
		t.Fatalf("vpc_open(2) unable to destroy: %v", err)
	}

	log.Debug().Int("vpcsw0CreateFD", int(vpcsw0CHandle.FD())).Msg("closing vpcsw0 create")
	if err := vpcsw0CHandle.Close(); err != nil {
		t.Fatalf("unable to close(2) VPC Handle : %v", err)
	}
	if vpcsw0CHandle.FD() != vpc.HandleClosedFD {
		t.Fatalf("handle set to wrong value in vpc.Close()")
	}
}

func BenchmarkVPCCreateOpenClose(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		vpcsw0ID := vpc.GenID()

		ht, err := vpc.NewHandleType(vpc.HandleTypeInput{
			Version: 1,
			Type:    vpc.ObjTypeSwitch,
		})
		if err != nil {
			b.Fatalf("unable to construct a HandleType: %v", err)
		}

		for pb.Next() {
			vpcsw0CreateFD, err := vpc.Open(vpcsw0ID, ht, vpc.FlagCreate)
			if err != nil {
				b.Fatalf("vpc_open(2) failed: %v", err)
			}

			if err := vpcsw0CreateFD.Close(); err != nil {
				b.Fatalf("unable to close(2) VPC Handle : %v", err)
			}
		}
	})
}
