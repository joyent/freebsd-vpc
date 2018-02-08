package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"time"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/internal/buildtime"
	"github.com/sean-/vpc/internal/logger"
)

const DefaultConnTimeout = 10 * time.Second

type DBPoolConfig = pgx.ConnPoolConfig
type DBScheme int

const (
	DBSchemeUnknown DBScheme = iota
	DBSchemePostgreSQL
	DBSchemeCRDB
)

func (t DBScheme) String() string {
	switch t {
	case DBSchemeUnknown:
		return "unknown"
	case DBSchemePostgreSQL:
		return "potsgres"
	case DBSchemeCRDB:
		return "cockroachdb"
	default:
		panic(fmt.Sprintf("unknown type: %v", t))
	}
}

type DB struct {
	Scheme     DBScheme
	ConnConfig pgx.ConnConfig
	PoolConfig DBPoolConfig

	CAPath   string
	CertPath string
	KeyPath  string
}

type Config struct {
	DB DB
}

func New() (*Config, error) {
	cfg := &Config{
		DB: DB{
			Scheme: DBSchemeCRDB,

			CAPath:   "/usr/home/seanc/go/src/github.com/sean-/vpc/crdb/certs/ca.crt",
			CertPath: "/usr/home/seanc/go/src/github.com/sean-/vpc/crdb/certs/client.root.crt",
			KeyPath:  "/usr/home/seanc/go/src/github.com/sean-/vpc/crdb/certs/client.root.key",

			PoolConfig: pgx.ConnPoolConfig{
				MaxConnections: 5,
				AfterConnect:   nil,
				AcquireTimeout: 0,

				ConnConfig: pgx.ConnConfig{
					Logger:   logger.NewPGX(log.Logger),
					Database: "triton",    //viper.GetString(KeyPGDatabase),
					User:     "root",      //viper.GetString(KeyPGUser),
					Password: "tls",       //viper.GetString(KeyPGPassword),
					Host:     "127.0.0.1", //viper.GetString(KeyPGHost),
					Port:     26257,       //cast.ToUint16(viper.GetInt(KeyPGPort)),
					Dial:     (&net.Dialer{Timeout: DefaultConnTimeout, KeepAlive: 5 * time.Minute}).Dial,

					UseFallbackTLS: false,
					TLSConfig:      nil,

					// FIXME(seanc@): Need to write a zerolog facade that satisfies the pgx logger interface
					// Logger:   log.Logger.With().Str("module", "pgx").Logger(),
					LogLevel: pgx.LogLevelTrace, //pgxLogLevel,
					RuntimeParams: map[string]string{
						"application_name": buildtime.PROGNAME,
					},
				},
			},
		},
	}

	tlsConfig, err := cfg.TLSConfig()
	if err != nil {
		return nil, errors.Wrap(err, "unable to load TLS config while initializing Config")
	}

	cfg.DB.PoolConfig.ConnConfig.TLSConfig = tlsConfig

	return cfg, nil
}

func (cfg *Config) Load() error {
	return nil
}

func (cfg *Config) TLSConfig() (*tls.Config, error) {
	caCertPool := x509.NewCertPool()
	{
		caPath := cfg.DB.CAPath
		caCert, err := ioutil.ReadFile(caPath)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to read CA file %q", caPath)
		}

		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, errors.Wrap(err, "unable to add CA to cert pool")
		}
	}

	cert, err := tls.LoadX509KeyPair(cfg.DB.CertPath, cfg.DB.KeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read cert")
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify:       true, // TODO(seanc@): make a tunable
		ServerName:               cfg.DB.PoolConfig.ConnConfig.Host,
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: true,
		Certificates:             []tls.Certificate{cert},
		RootCAs:                  caCertPool,
		ClientCAs:                caCertPool,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}
	// tlsConfig.BuildNameToCertificate()

	return tlsConfig, nil
}
