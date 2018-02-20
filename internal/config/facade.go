package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.freebsd.org/sys/vpc"
)

func GetUUID(v *viper.Viper, key string) (vpc.ID, error) {
	uuidStr := v.GetString(key)

	id, err := vpc.ParseID(uuidStr)
	if err != nil {
		return vpc.ID{}, errors.Wrapf(err, "unable to parse UUID: %q", uuidStr)
	}

	return id, nil
}
