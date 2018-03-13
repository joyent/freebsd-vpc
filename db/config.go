package db

import (
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Config struct {
	Scheme             string        `mapstructure:"scheme"`
	User               string        `mapstructure:"user"`
	Password           string        `mapstructure:"password"`
	Host               string        `mapstructure:"host"`
	Port               uint16        `mapstructure:"port"`
	Database           string        `mapstructure:"database"`
	UseTLSClientAuth   bool          `mapstructure:"use_tls_client_auth"`
	CAPath             string        `mapstructure:"ca_path"`
	CertPath           string        `mapstructure:"cert_path"`
	KeyPath            string        `mapstructure:"key_path"`
	ConnTimeout        time.Duration `mapstructure:"conn_timeout"`
	InsecureSkipVerify bool          `mapstructure:"insecure_skip_verify"`
}

func SetDefaultViperOptions() error {
	viper.SetDefault("db.scheme", "crdb")
	viper.SetDefault("db.user", "root")
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", 26257)
	viper.SetDefault("db.database", "triton")
	viper.SetDefault("db.conn_timeout", 10*time.Second)
	viper.SetDefault("db.insecure_skip_verify", false)
	viper.SetDefault("db.use_tls_client_auth", true)

	// Note: these are the default certificate paths for CockroachDB, as used
	// by the interactive `cockroach sql` command.
	caPath, err := homedir.Expand("~/.cockroach-certs/ca.crt")
	if err != nil {
		return errors.Wrap(err, "error expanding home directory")
	}
	viper.SetDefault("db.ca_path", caPath)

	certPath, err := homedir.Expand("~/.cockroach-certs/client.root.crt")
	if err != nil {
		return errors.Wrap(err, "error expanding home directory")
	}
	viper.SetDefault("db.cert_path", certPath)

	keyPath, err := homedir.Expand("~/.cockroach-certs/client.root.key")
	if err != nil {
		return errors.Wrap(err, "error expanding home directory")
	}
	viper.SetDefault("db.key_path", keyPath)

	return nil
}
