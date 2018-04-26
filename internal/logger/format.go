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

package logger

import (
	"fmt"
	"strings"

	"github.com/joyent/freebsd-vpc/internal/config"
	"github.com/spf13/viper"
)

type Format uint

const (
	FormatAuto    Format = iota
	FormatZerolog
	FormatHuman
)

func (f Format) String() string {
	switch f {
	case FormatAuto:
		return "auto"
	case FormatZerolog:
		return "zerolog"
	case FormatHuman:
		return "human"
	default:
		panic(fmt.Sprintf("unknown log format: %d", f))
	}
}

func getLogFormat(v *viper.Viper) (Format, error) {
	switch logFormat := strings.ToLower(v.GetString(config.KeyLogFormat)); logFormat {
	case "auto":
		return FormatAuto, nil
	case "json", "zerolog":
		return FormatZerolog, nil
	case "human":
		return FormatHuman, nil
	default:
		return FormatAuto, fmt.Errorf("unsupported log format: %q", logFormat)
	}
}
