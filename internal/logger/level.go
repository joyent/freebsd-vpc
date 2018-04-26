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
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Level int

const (
	LevelBegin Level = iota - 2
	LevelDebug
	LevelInfo   // Default, zero-initialized value
	LevelWarn
	LevelError
	LevelFatal

	LevelEnd
)

func (f Level) String() string {
	switch f {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	default:
		panic(fmt.Sprintf("unknown log level: %d", f))
	}
}

func logLevels() []Level {
	levels := make([]Level, 0, LevelEnd-LevelBegin)
	for i := LevelBegin + 1; i < LevelEnd; i++ {
		levels = append(levels, i)
	}

	return levels
}

func logLevelsStr() []string {
	intLevels := logLevels()
	levels := make([]string, 0, len(intLevels))
	for _, lvl := range intLevels {
		levels = append(levels, lvl.String())
	}
	return levels
}

func setLogLevel(v *viper.Viper) (logLevel Level, err error) {
	switch strLevel := strings.ToLower(v.GetString(config.KeyLogLevel)); strLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		logLevel = LevelDebug
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logLevel = LevelInfo
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
		logLevel = LevelWarn
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		logLevel = LevelError
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
		logLevel = LevelFatal
	default:
		return LevelDebug, fmt.Errorf("unsupported error level: %q (supported levels: %s)", logLevel,
			strings.Join(logLevelsStr(), " "))
	}

	return logLevel, nil
}
