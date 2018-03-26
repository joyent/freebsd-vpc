// Go interface to Hostlink NIC objects.
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

package hostlink

import (
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc"
	"github.com/pkg/errors"
)

// _HostlinkCmd is the encoded type of operations that can be performed on a
// Hostlink NIC.
type _HostlinkCmd vpc.Cmd

// Close closes the VPC Handle descriptor.  Created Hostlinkn NICs will not be
// destroyed when the Hostlink is closed if the Hostlink NIC has been Committed.
func (hl *Hostlink) Close() error {
	if hl.h.FD() <= 0 {
		return nil
	}

	if err := hl.h.Close(); err != nil {
		return errors.Wrap(err, "unable to close VPC Hostlink handle")
	}

	return nil
}

// Commit increments the refcount of the VPC Hostlink NIC in order to ensure the
// Hostlink NIC lives beyond the life of the current process and is not
// automatically cleaned up when the Hostlink is closed.
func (hl *Hostlink) Commit() error {
	if hl.h.FD() <= 0 {
		return nil
	}

	if err := hl.h.Commit(); err != nil {
		return errors.Wrap(err, "unable to commit VPC Hostlink NIC")
	}

	return nil
}

// Destroy decrements the refcount of the VPC Hostlink NIC in destroy the the
// Hostlink NIC when the VPC Handle is closed.
func (hl *Hostlink) Destroy() error {
	if hl.h.FD() <= 0 {
		return nil
	}

	if err := hl.h.Destroy(); err != nil {
		return errors.Wrap(err, "unable to destroy VPC Hostlink NIC")
	}

	return nil
}
