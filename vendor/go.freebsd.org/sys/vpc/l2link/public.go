// Go interface to Layer-2 Network Link objects.
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
	"github.com/rs/zerolog"
	"go.freebsd.org/sys/vpc"
)

// Config is the configuration used to populate a given L2 Link.
type Config struct {
	ID   vpc.ID
	Name string
}

func (c Config) MarshalZerologObject(e *zerolog.Event) {
	e.Str("id", c.ID.String()).
		Str("name", c.Name)
}

// L2Link is an opaque struct representing a VM NIC.
type L2Link struct {
	h    *vpc.Handle
	ht   vpc.HandleType
	id   vpc.ID
	name string
}

// Create VPC facade over an existing L2 link (either physical or cloned
// interface) using the Config parameters.  Callers are expected to Close a
// given L2Link (otherwise a file descriptor would leak).
func Create(cfg Config) (*L2Link, error) {
	ht, err := vpc.NewHandleType(vpc.HandleTypeInput{
		Version: 1,
		Type:    vpc.ObjTypeNICVM,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to create a new L2 Link handle type")
	}

	h, err := vpc.Open(cfg.ID, ht, vpc.FlagCreate|vpc.FlagWrite)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open L2 Link handle")
	}

	return &L2Link{
		h:    h,
		ht:   ht,
		id:   cfg.ID,
		name: cfg.Name,
	}, nil
}

// Open opens an existing L2 Link using the Config parameters.  Callers are
// expected to Close a given L2Link.
func Open(cfg Config) (*L2Link, error) {
	ht, err := vpc.NewHandleType(vpc.HandleTypeInput{
		Version: 1,
		Type:    vpc.ObjTypeLinkL2,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to create a new L2 Link handle type")
	}

	h, err := vpc.Open(cfg.ID, ht, vpc.FlagOpen|vpc.FlagRead|vpc.FlagWrite)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open L2 Link handle")
	}

	return &L2Link{
		h:    h,
		ht:   ht,
		id:   cfg.ID,
		name: cfg.Name,
	}, nil
}
