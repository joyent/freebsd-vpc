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
