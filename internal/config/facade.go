package config

import (
	"github.com/freebsd/freebsd/libexec/go/src/go.freebsd.org/sys/vpc"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func GetUUID(v *viper.Viper, key string) (vpc.ID, error) {
	uuidStr := v.GetString(key)

	id, err := vpc.ParseID(uuidStr)
	if err != nil {
		return vpc.ID{}, errors.Wrapf(err, "unable to parse UUID: %q", uuidStr)
	}

	return id, nil
}
