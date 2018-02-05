package logger

import (
	stdlog "log"
	"os"

	"github.com/jackc/pgx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	// Use a log format that resembles time.RFC3339Nano but includes all trailing
	// zeros so that we get fixed-width logging.
	LogTimeFormat = "2006-01-02T15:04:05.000000000Z07:00"
)

func init() {
	zerolog.TimeFieldFormat = LogTimeFormat
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// os.Stderr isn't guaranteed to be thread-safe, wrap in a sync writer.  Files
	// are guaranteed to be safe, terminals are not.
	zlog := zerolog.New(zerolog.SyncWriter(os.Stderr)).With().Timestamp().Logger()
	log.Logger = zlog

	stdlog.SetFlags(0)
	stdlog.SetOutput(zlog)
}

type PGX struct {
	l zerolog.Logger
}

func NewPGX(l zerolog.Logger) pgx.Logger {
	return &PGX{l: l}
}

func (l PGX) Log(level pgx.LogLevel, msg string, data map[string]interface{}) {
	switch level {
	case pgx.LogLevelDebug:
		l.l.Debug().Fields(data).Msg(msg)
	case pgx.LogLevelInfo:
		l.l.Info().Fields(data).Msg(msg)
	case pgx.LogLevelWarn:
		l.l.Warn().Fields(data).Msg(msg)
	case pgx.LogLevelError:
		l.l.Error().Fields(data).Msg(msg)
	default:
		l.l.Debug().Fields(data).Str("level", level.String()).Msg(msg)
	}
}
