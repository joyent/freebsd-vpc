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
	"net"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.freebsd.org/sys/vpc"
)

// DeviceNamePrefix is the prefix of the device name (i.e. "vmnic0").
const DeviceNamePrefix = "vmnic"

// Config is the configuration used to populate a given VM NIC.
type Config struct {
	ID        vpc.ID
	MAC       net.HardwareAddr
	Writeable bool
}

func (c Config) MarshalZerologObject(e *zerolog.Event) {
	e.Str("id", c.ID.String()).
		Str("mac", c.MAC.String())
}

// VMNIC is an opaque struct representing a VM NIC.
type VMNIC struct {
	h   *vpc.Handle
	ht  vpc.HandleType
	id  vpc.ID
	mac net.HardwareAddr
}

// Create creates a new VM NIC using the Config parameters.  Callers are
// expected to Close a given VMNIC (otherwise a file descriptor would leak).
func Create(cfg Config) (*VMNIC, error) {
	ht, err := vpc.NewHandleType(vpc.HandleTypeInput{
		Version: 1,
		Type:    vpc.ObjTypeNICVM,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to create a new VM NIC handle type")
	}

	h, err := vpc.Open(cfg.ID, ht, vpc.FlagCreate|vpc.FlagWrite)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open VM NIC handle")
	}

	return &VMNIC{
		h:   h,
		ht:  ht,
		id:  cfg.ID,
		mac: cfg.MAC,
	}, nil
}

// Close closes the VPC Handle descriptor.  Created VM NICs will not be
// destroyed when the VMNIC is closed if the VM NIC has been Committed.
func (sw *VMNIC) Close() error {
	if sw.h.FD() <= 0 {
		return nil
	}

	if err := sw.h.Close(); err != nil {
		return errors.Wrap(err, "unable to close VPC handle")
	}

	return nil
}

// Commit increments the refcount of the VM NIC in order to ensure the VM NIC
// lives beyond the life of the current process and is not automatically cleaned
// up when the VMNIC is closed.
func (sw *VMNIC) Commit() error {
	if sw.h.FD() <= 0 {
		return nil
	}

	if err := sw.h.Commit(); err != nil {
		return errors.Wrap(err, "unable to commit VM NIC")
	}

	return nil
}

// Destroy decrements the refcount of the VM NIC in destroy the the VM NIC when
// the VPC Handle is closed.
func (sw *VMNIC) Destroy() error {
	if sw.h.FD() <= 0 {
		return nil
	}

	if err := sw.h.Destroy(); err != nil {
		return errors.Wrap(err, "unable to destroy VM NIC")
	}

	return nil
}

// Open opens an existing VM NIC using the Config parameters.  Callers are
// expected to Close a given VMNIC.
func Open(cfg Config) (*VMNIC, error) {
	ht, err := vpc.NewHandleType(vpc.HandleTypeInput{
		Version: 1,
		Type:    vpc.ObjTypeNICVM,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to create a new VM NIC handle type")
	}

	flags := vpc.FlagOpen | vpc.FlagRead
	if cfg.Writeable {
		flags |= vpc.FlagWrite
	}

	h, err := vpc.Open(cfg.ID, ht, flags)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open VM NIC handle")
	}

	return &VMNIC{
		h:  h,
		ht: ht,
		id: cfg.ID,
	}, nil
}
