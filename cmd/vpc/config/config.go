package config

import (
	"github.com/joyent/freebsd-vpc/db"
)

type Config struct {
	DBConfig db.Config `mapstructure:"db"`
	MetadataConfig map[string]struct {
		InstanceID string `mapstructure:"instance_id"`
	} `mapstructure:"mdata"`
}
