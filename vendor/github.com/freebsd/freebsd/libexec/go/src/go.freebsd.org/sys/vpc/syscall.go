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

package vpc

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
)

// VNI is the type for VXLAN Network Identifiers ("VNI")
type VNI int32

const (
	// IDSize is the sizeof(ID)
	IDSize = 16

	// VNIMax is the largest permitted VNI
	VNIMax VNI = (1 << 24) - 1

	// VNIMin is the smallest permitted VNI.  NOTE: a VNI of 0 implies
	// un-encapsulated frames.
	VNIMin VNI = 0
)

// Byter ensures objects can be converted to binary
type Byter interface {
	// Bytes returns a byte representation of the receiver
	Bytes() []byte
}

// ID is the globally unique identifier for a VPC object.
type ID struct {
	TimeLow    uint32
	TimeMid    uint16
	TimeHi     uint16
	ClockSeqHi uint8
	ObjType    ObjType
	Node       [6]byte
}

func (id ID) MarshalZerologObject(e *zerolog.Event) {
	e.Str("id", id.String())
}

// Bytes returns a id as a little endian byte slice
func (id ID) Bytes() []byte {
	var binBuf bytes.Buffer
	binBuf.Grow(16)
	binary.Write(&binBuf, binary.LittleEndian, id)
	return binBuf.Bytes()
}

// GenID randomly generates a new UUID
func GenID(objType ObjType) ID {
	randUint8 := func() uint8 {
		var b [1]byte
		if _, err := rand.Read(b[:]); err != nil {
			panic("bad")
		}
		return uint8(b[0])
	}

	randUint16 := func() uint16 {
		var b [2]byte
		if _, err := rand.Read(b[:]); err != nil {
			panic("bad")
		}
		return uint16(binary.LittleEndian.Uint16(b[:]))
	}

	randUint32 := func() uint32 {
		var b [4]byte
		if _, err := rand.Read(b[:]); err != nil {
			panic("bad")
		}
		return uint32(binary.LittleEndian.Uint32(b[:]))
	}

	randNode := func() [6]byte {
		var b [6]byte
		if _, err := rand.Read(b[:]); err != nil {
			panic("bad")
		}
		// #define    ETHER_IS_MULTICAST(addr) (*(addr) & 0x01) /* is address mcast/bcast? */
		b[0] = b[0] &^ 0x01
		return b
	}

	// FIXME(seanc@): I took the bruteforce way of populating a struct with random
	// data vs just populating a [16]byte slice w/ random data and casting it to
	// an ID because I didn't want to fight with the language, but this should be
	// done better and differently.
	return ID{
		TimeLow:    randUint32(),
		TimeMid:    randUint16(),
		TimeHi:     randUint16(),
		ClockSeqHi: randUint8(),
		ObjType:    objType,
		Node:       randNode(),
	}
}

// ParseID parses a UUID string and converts it into an ID.  ParseID will return
// an error if the UUID is malformed or if the Node portion of the UUID has its
// multicast/broadcast bit set.
func ParseID(idStr string) (ID, error) {
	uuidRaw, err := uuid.FromString(idStr)
	if err != nil {
		return ID{}, errors.Wrapf(err, "unable to parse UUID: %q", idStr)
	}

	// #define    ETHER_IS_MULTICAST(addr) (*(addr) & 0x01) /* is address mcast/bcast? */
	if uuidRaw[10]&0x01 == 1 {
		return ID{}, errors.New("broadcast bit set in Node portion of UUID")
	}

	id := ID{
		TimeLow:    binary.LittleEndian.Uint32(uuidRaw[0:]),
		TimeMid:    binary.LittleEndian.Uint16(uuidRaw[4:]),
		TimeHi:     binary.LittleEndian.Uint16(uuidRaw[6:]),
		ClockSeqHi: uint8(binary.LittleEndian.Uint16(uuidRaw[8:])),
		ObjType:    ObjType(binary.LittleEndian.Uint16(uuidRaw[9:])),
	}

	buf := bytes.NewReader(uuidRaw[10:])
	err = binary.Read(buf, binary.LittleEndian, id.Node[0:])
	if err != nil {
		return ID{}, errors.Wrap(err, "unable to read bytes")
	}

	return id, nil
}

func (id ID) String() string {
	var binBuf bytes.Buffer
	binBuf.Grow(16)
	binary.Write(&binBuf, binary.LittleEndian, id)
	uuid := binBuf.Bytes()

	var buf [36]byte
	hex.Encode(buf[:], uuid[:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], uuid[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], uuid[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], uuid[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], uuid[10:])

	return string(buf[:])
}

// OpenFlags is the flags passed to Open
type OpenFlags uint64

const (
	// FlagCreate is used to create a new VPC object.
	FlagCreate OpenFlags = 1 << 0

	// FlagOpen is used to open an existing VPC object.
	FlagOpen OpenFlags = 1 << 1

	// FlagRead is used to open an existing VPC object for reading.
	FlagRead OpenFlags = 1 << 2

	// FlagWrite is used to open a VPC object for writes (including commit).
	FlagWrite OpenFlags = 1 << 3
)

// ObjType distinguishes the different types of supported VPC Object Types.
type ObjType uint8

func (objType ObjType) Bytes() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(buf, uint64(objType))
	return buf[:n]
}

func (objType ObjType) MarshalZerologObject(e *zerolog.Event) {
	e.Str("type", objType.String())
}

// Exported enumerated types of available VPC objects.
const (
	ObjTypeInvalid    ObjType = 0
	ObjTypeSwitch     ObjType = 1
	ObjTypeSwitchPort ObjType = 2
	ObjTypeRouter     ObjType = 3
	ObjTypeNAT        ObjType = 4
	ObjTypeMux        ObjType = 5
	ObjTypeNICVM      ObjType = 6
	ObjTypeMgmt       ObjType = 7
	ObjTypeLinkEth    ObjType = 8
	ObjTypeMeta       ObjType = 9
	ObjTypeAny        ObjType = 10
)

// ObjTypes returns a lits of supported Object Types
func ObjTypes() []ObjType {
	return []ObjType{
		// ObjTypeInvalid,
		ObjTypeSwitch,
		ObjTypeSwitchPort,
		ObjTypeRouter,
		ObjTypeNAT,
		ObjTypeMux,
		ObjTypeNICVM,
		ObjTypeMgmt,
		ObjTypeLinkEth,
		// ObjTypeMeta, // Not a queriable type
		// ObjTypeAny, // Not a queriable type
	}
}

// String returns the string representation of a given object
func (obj ObjType) String() string {
	switch obj {
	case ObjTypeInvalid:
		return "invalid"
	case ObjTypeSwitch:
		return "vpcsw"
	case ObjTypeSwitchPort:
		return "vpcp"
	case ObjTypeRouter:
		return "vpcrtr"
	case ObjTypeNAT:
		return "vpcnat"
	case ObjTypeMux:
		return "vpcmux"
	case ObjTypeNICVM:
		return "vmnic"
	case ObjTypeMgmt:
		return "mgmt"
	case ObjTypeLinkEth:
		return "ethlink"
	case ObjTypeMeta:
		return "meta"
	case ObjTypeAny:
		return "any"
	default:
		panic(fmt.Sprintf("unsupported object type: 0x%02x", uint8(obj)))
	}
}

// Close closes a VPC Handle.  Closing a VPC Handle does not destroy any
// resources.
func (h *Handle) Close() error {
	h.lock.Lock()
	defer h.lock.Unlock()

	if h.fd == HandleClosedFD {
		return nil
	}

	return h.closeHandle()
}
