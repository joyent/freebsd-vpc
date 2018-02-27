// Go interface to VPC Switch Port objects.
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

package vpcp

import (
	"net"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.freebsd.org/sys/vpc"
)

// Config is the configuration used to populate a given VPC Swith Port.
type Config struct {
	ID        vpc.ID
	MAC       net.HardwareAddr
	Writeable bool
}

func (c Config) MarshalZerologObject(e *zerolog.Event) {
	e.Str("id", c.ID.String()).
		Str("mac", c.MAC.String())
}

// VPCP is an opaque struct representing a VPC Switch Port.
type VPCP struct {
	h   *vpc.Handle
	ht  vpc.HandleType
	id  vpc.ID
	mac net.HardwareAddr
}

// TODO(seanc@): add vpc_ctl(2) to create a port and an op to a vpc switch to
// attach a port to a switch.
//
// // Create creates a new VM NIC using the Config parameters.  Callers are
// // expected to Close a given VPCP (otherwise a file descriptor would leak).
// func Create(cfg Config) (*VPCP, error) {
// 	ht, err := vpc.NewHandleType(vpc.HandleTypeInput{
// 		Version: 1,
// 		Type:    vpc.ObjTypeNICVM,
// 	})
// 	if err != nil {
// 		return nil, errors.Wrap(err, "unable to create a new VM NIC handle type")
// 	}
//
// 	h, err := vpc.Open(cfg.ID, ht, vpc.FlagCreate|vpc.FlagWrite)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "unable to open VM NIC handle")
// 	}
//
// 	return &VPCP{
// 		h:   h,
// 		ht:  ht,
// 		id:  cfg.ID,
// 		mac: cfg.MAC,
// 	}, nil
// }

// Close closes the VPC descriptor.  The life cycle of a VPC Switch Port is
// attached to the VPC Switch Port and will be destroyed when a VPC Switch is
// destroyed.
func (p *VPCP) Close() error {
	if p.h.FD() <= 0 {
		return nil
	}

	if err := p.h.Close(); err != nil {
		return errors.Wrap(err, "unable to close VPC handle")
	}

	return nil
}

// TODO(seanc@): repurpose Commit() to be the attach operation to add a port to
// a switch.
//
// // Commit increments the refcount of the VM NIC in order to ensure the VM NIC
// // lives beyond the life of the current process and is not automatically cleaned
// // up when the VPCP is closed.
// func (p *VPCP) Commit() error {
// 	if p.h.FD() <= 0 {
// 		return nil
// 	}

// 	if err := p.h.Commit(); err != nil {
// 		return errors.Wrap(err, "unable to commit VM NIC")
// 	}

// 	return nil
// }

// TODO(seanc@): Repurpose Destroy to be the detatch operation to remove a port
// from a switch.
//
// // Destroy decrements the refcount of the VM NIC in destroy the the VM NIC when
// // the VPC Handle is closed.
// func (p *VPCP) Destroy() error {
// 	if p.h.FD() <= 0 {
// 		return nil
// 	}
//
// 	if err := p.h.Destroy(); err != nil {
// 		return errors.Wrap(err, "unable to destroy VM NIC")
// 	}
//
// 	return nil
// }

// Open opens an existing VPC Switch Port using the Config parameters.  Callers
// are expected to Close a given VPCP.
func Open(cfg Config) (*VPCP, error) {
	ht, err := vpc.NewHandleType(vpc.HandleTypeInput{
		Version: 1,
		Type:    vpc.ObjTypeSwitchPort,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to create a new VPC Switch Port handle type")
	}

	flags := vpc.FlagOpen | vpc.FlagRead
	if cfg.Writeable {
		flags |= vpc.FlagWrite
	}

	h, err := vpc.Open(cfg.ID, ht, flags)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open VPC Switch Port handle")
	}

	return &VPCP{
		h:  h,
		ht: ht,
		id: cfg.ID,
	}, nil
}
