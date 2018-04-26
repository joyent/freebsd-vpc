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

package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"io/ioutil"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/joyent/freebsd-vpc/internal/buildtime"
	"github.com/joyent/freebsd-vpc/internal/logger"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Pool struct {
	pool            *pgx.ConnPool
	stdDriverConfig *stdlib.DriverConfig
	config          Config
}

func New(cfg Config) (*Pool, error) {
	pool := &Pool{
		config: cfg,
	}

	tlsConfig, err := cfg.TLSConfig()
	if err != nil {
		return nil, errors.Wrap(err, "unable to generate a TLS config")
	}

	const keepAliveTimeout = 5 * time.Minute

	// TODO(jen20) do we even want to support password auth vs TLS?
	poolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Logger:   logger.NewPGX(log.Logger),
			Database: cfg.Database,
			User:     cfg.User,
			//Password: cfg.Password,
			Host: cfg.Host,
			Port: cfg.Port,
			Dial: (&net.Dialer{Timeout: cfg.ConnTimeout, KeepAlive: keepAliveTimeout}).Dial,

			UseFallbackTLS: false,
			TLSConfig:      tlsConfig,

			// TODO(seanc): Need to write a zerolog facade that satisfies the pgx logger interface
			// Logger:   log.Logger.With().Str("module", "pgx").Logger(),
			LogLevel: pgx.LogLevelTrace, //pgxLogLevel,
			RuntimeParams: map[string]string{
				"application_name": buildtime.PROGNAME,
			},
		},
	}

	p, err := pgx.NewConnPool(poolConfig)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create a new DB connection pool")
	}
	pool.pool = p

	if err := pool.Ping(); err != nil {
		return nil, errors.Wrap(err, "unable to ping database")
	}

	pool.stdDriverConfig = &stdlib.DriverConfig{
		ConnConfig: pgx.ConnConfig{
			Database:       cfg.Database,
			Dial:           (&net.Dialer{Timeout: cfg.ConnTimeout, KeepAlive: keepAliveTimeout}).Dial,
			Host:           cfg.Host,
			LogLevel:       pgx.LogLevelTrace, //pgxLogLevel,
			Logger:         logger.NewPGX(log.Logger),
			Port:           cfg.Port,
			TLSConfig:      tlsConfig,
			UseFallbackTLS: false,
			User:           cfg.User,
			RuntimeParams: map[string]string{
				"application_name": buildtime.PROGNAME,
			},
		},
		AfterConnect: func(c *pgx.Conn) error {
			return nil
		},
	}
	stdlib.RegisterDriverConfig(pool.stdDriverConfig)

	return pool, nil
}

func (p *Pool) Close() error {
	p.pool.Close()

	return nil
}

func (p *Pool) Ping() error {
	pingCtx, pingCancel := context.WithTimeout(context.Background(), p.config.ConnTimeout)
	defer pingCancel()
	conn, err := p.pool.Acquire()
	if err != nil {
		return errors.Wrap(err, "unable to acquire database connection for ping")
	}
	defer p.pool.Release(conn)

	if err := conn.Ping(pingCtx); err != nil {
		return errors.Wrap(err, "unable to ping database")
	}

	return nil
}

func (p *Pool) Pool() *pgx.ConnPool {
	return p.pool
}

func (p *Pool) STDDB() (*sql.DB, error) {
	username := p.config.User
	hostname := p.config.Host
	port := p.config.Port
	database := p.config.Database

	v := url.Values{}
	if !p.config.InsecureSkipVerify {
		sslMode := "require"

		v.Set("sslmode", sslMode)
		v.Set("sslrootcert", p.config.CAPath)
		v.Set("sslcert", p.config.CertPath)
		v.Set("sslkey", p.config.KeyPath)
	}

	u := url.URL{
		Scheme:   p.config.Scheme,
		User:     url.User(username),
		Host:     net.JoinHostPort(hostname, strconv.Itoa(int(port))),
		Path:     database,
		RawQuery: v.Encode(),
	}

	encodedURI := p.stdDriverConfig.ConnectionString(u.String())
	db, err := sql.Open("pgx", encodedURI)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open a standard database connection")
	}

	return db, nil
}

func (p *Config) TLSConfig() (*tls.Config, error) {
	caCertPool := x509.NewCertPool()
	{
		caPath := p.CAPath
		caCert, err := ioutil.ReadFile(caPath)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to read CA file %q", caPath)
		}

		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, errors.Wrap(err, "unable to add CA to cert pool")
		}
	}

	cert, err := tls.LoadX509KeyPair(p.CertPath, p.KeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read cert")
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify:       p.InsecureSkipVerify,
		ServerName:               p.Host,
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

	return tlsConfig, nil
}
