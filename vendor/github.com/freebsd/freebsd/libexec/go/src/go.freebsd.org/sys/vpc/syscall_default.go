// Go interface to VPC syscalls on non-FreeBSD systems.
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

// +build android darwin dragonfly linux nacl netbsd openbsd plan9 solaris windows

package vpc

import (
	"github.com/pkg/errors"
)

// Open obtains a VPC handle to a given object type.  Obtaining an open Handle
// affords no privilges beyond validating that an ID exists on this system.  In
// all other cases Open returns a handle to a resource.  If the id can not be
// found, Open returns ENOENT unless the Create flag is set in flags.  If the
// Create flag is set and the id is found, Open returns EEXIST.  If an invalid
// Flag is set, Open returns EINVAL.  If the HandleType is out of bounds, Open
// returns EOPNOTSUPP.
func Open(id ID, ht HandleType, flags Flags) (Handle, error) {
	return ErrorHandle, errors.New("not implemented")
}

// Ctl manipulates the Handle based on the args
func Ctl(h Handle, op OpFlag, in []byte, out *[]byte) error {
	return errors.New("not implemented")
}
