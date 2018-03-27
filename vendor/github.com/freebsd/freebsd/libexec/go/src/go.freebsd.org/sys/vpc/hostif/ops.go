// Go interface to Hostif NIC objects.
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

package hostif

import (
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc"
	"github.com/pkg/errors"
)

// _HostifCmd is the encoded type of operations that can be performed on a
// Hostif NIC.
type _HostifCmd vpc.Cmd

// Close closes the VPC Handle descriptor.  Created Hostif NICs will not be
// destroyed when the Hostif is closed if the Hostif NIC has been Committed.
func (hl *Hostif) Close() error {
	if hl.h.FD() <= 0 {
		return nil
	}

	if err := hl.h.Close(); err != nil {
		return errors.Wrap(err, "unable to close VPC Hostif handle")
	}

	return nil
}

// Commit increments the refcount of the VPC Hostif NIC in order to ensure the
// Hostif NIC lives beyond the life of the current process and is not
// automatically cleaned up when the Hostif is closed.
func (hl *Hostif) Commit() error {
	if hl.h.FD() <= 0 {
		return nil
	}

	if err := hl.h.Commit(); err != nil {
		return errors.Wrap(err, "unable to commit VPC Hostif NIC")
	}

	return nil
}

// Destroy decrements the refcount of the VPC Hostif NIC in destroy the the
// Hostif NIC when the VPC Handle is closed.
func (hl *Hostif) Destroy() error {
	if hl.h.FD() <= 0 {
		return nil
	}

	if err := hl.h.Destroy(); err != nil {
		return errors.Wrap(err, "unable to destroy VPC Hostif NIC")
	}

	return nil
}
