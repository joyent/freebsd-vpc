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
