// Go interface to VPC Mux objects.
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

package mux

import (
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Config is the configuration used to populate a given VPC Mux
type Config struct {
	ID        vpc.ID
	Writeable bool
}

func (c Config) MarshalZerologObject(e *zerolog.Event) {
	e.Str("id", c.ID.String())
}

// Mux is an opaque struct representing a VPC Mux.
type Mux struct {
	h  *vpc.Handle
	ht vpc.HandleType
	id vpc.ID
}

// Create creates a new VPC Mux using the Config parameters.  Callers are
// expected to Close a given VPC Mux (otherwise a file descriptor would leak).
func Create(cfg Config) (*Mux, error) {
	ht, err := vpc.NewHandleType(vpc.HandleTypeInput{
		Version: 1,
		Type:    vpc.ObjTypeMux,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to create a new VPC Mux handle type")
	}

	h, err := vpc.Open(cfg.ID, ht, vpc.FlagCreate|vpc.FlagWrite)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open VPC Mux handle")
	}

	return &Mux{
		h:  h,
		ht: ht,
		id: cfg.ID,
	}, nil
}

// Open opens an existing VPC Mux using the Config parameters.  Callers are
// expected to Close a given Mux.
func Open(cfg Config) (*Mux, error) {
	ht, err := vpc.NewHandleType(vpc.HandleTypeInput{
		Version: 1,
		Type:    vpc.ObjTypeMux,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to create a new VPC Mux handle type")
	}

	flags := vpc.FlagOpen | vpc.FlagRead
	if cfg.Writeable {
		flags |= vpc.FlagWrite
	}

	h, err := vpc.Open(cfg.ID, ht, flags)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open VPC Mux handle")
	}

	return &Mux{
		h:  h,
		ht: ht,
		id: cfg.ID,
	}, nil
}
