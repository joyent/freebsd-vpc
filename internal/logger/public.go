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
	"io"
	"io/ioutil"
	stdlog "log"
	"os"
	"time"

	"github.com/joyent/freebsd-vpc/internal/config"
	"github.com/mattn/go-isatty"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sean-/conswriter"
	"github.com/spf13/viper"
)

const (
	// Use a log format that resembles time.RFC3339Nano but includes all trailing
	// zeros so that we get fixed-width logging.
	logTimeFormat = "2006-01-02T15:04:05.000000000Z07:00"
)

var stdLogger *stdlog.Logger

func init() {
	// Initialize zerolog with a set set of defaults.  Re-initialization of
	// logging with user-supplied configuration parameters happens in Setup().

	// os.Stderr isn't guaranteed to be thread-safe, wrap in a sync writer.  Files
	// are guaranteed to be safe, terminals are not.
	w := zerolog.ConsoleWriter{
		Out:     os.Stderr,
		NoColor: true,
	}
	zlog := zerolog.New(zerolog.SyncWriter(w)).With().Timestamp().Logger()

	zerolog.DurationFieldUnit = time.Microsecond
	zerolog.DurationFieldInteger = true
	zerolog.TimeFieldFormat = logTimeFormat
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Logger = zlog

	stdlog.SetFlags(0)
	stdlog.SetOutput(zlog)
}

func Setup(v *viper.Viper) error {
	logLevel, err := setLogLevel(v)
	if err != nil {
		return errors.Wrap(err, "unable to set log level")
	}

	var logWriter io.Writer
	if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		logWriter = conswriter.GetTerminal()
	} else {
		logWriter = os.Stderr
	}

	logFmt, err := getLogFormat(v)
	if err != nil {
		return errors.Wrap(err, "unable to parse log format")
	}

	if logFmt == FormatAuto {
		if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
			logFmt = FormatHuman
		} else {
			logFmt = FormatZerolog
		}
	}

	var zlog zerolog.Logger
	switch logFmt {
	case FormatZerolog:
		zlog = zerolog.New(logWriter).With().Timestamp().Logger()
	case FormatHuman:
		useColor := v.GetBool(config.KeyLogTermColor)
		w := zerolog.ConsoleWriter{
			Out:     logWriter,
			NoColor: !useColor,
		}
		zlog = zerolog.New(w).With().Timestamp().Logger()
	default:
		return fmt.Errorf("unsupported log format: %q", logFmt)
	}

	zlog.Hook(closeConWriterHook{})

	log.Logger = zlog

	stdlog.SetFlags(0)
	stdlog.SetOutput(zlog)
	stdLogger = &stdlog.Logger{}

	// In order to prevent random libraries from hooking the standard logger and
	// filling the logger with garbage, discard all log entries.  At debug level,
	// however, let it all through.
	if logLevel != LevelDebug {
		stdLogger.SetOutput(ioutil.Discard)
	} else {
		stdLogger.SetOutput(zlog)
	}

	return nil
}

type closeConWriterHook struct{}

func (h closeConWriterHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level != zerolog.FatalLevel {
		return
	}

	conswriter.GetTerminal().Close()
}
