package db

import (
	"context"
	"database/sql"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sean-/vpc/agent/logger"
	"github.com/sean-/vpc/buildtime"
	"github.com/sean-/vpc/config"
)

type Pool struct {
	cfg             *config.Config
	pool            *pgx.ConnPool
	stdDriverConfig *stdlib.DriverConfig
}

func New(cfg *config.Config) (*Pool, error) {
	pool := &Pool{
		cfg: cfg,
	}

	p, err := pgx.NewConnPool(cfg.DB.PoolConfig)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create a new DB connection pool")
	}
	pool.pool = p

	if err := pool.Ping(); err != nil {
		return nil, errors.Wrap(err, "unable to ping database")
	}

	tlsConfig, err := cfg.TLSConfig()
	if err != nil {
		return nil, errors.Wrap(err, "unable to generate a TLS config")
	}

	pool.stdDriverConfig = &stdlib.DriverConfig{
		ConnConfig: pgx.ConnConfig{
			Database:       cfg.DB.PoolConfig.ConnConfig.Database,
			Dial:           (&net.Dialer{Timeout: config.DefaultConnTimeout, KeepAlive: 5 * time.Minute}).Dial,
			Host:           cfg.DB.PoolConfig.ConnConfig.Host,
			LogLevel:       pgx.LogLevelTrace, //pgxLogLevel,
			Logger:         logger.NewPGX(log.Logger),
			Port:           cfg.DB.PoolConfig.ConnConfig.Port,
			TLSConfig:      tlsConfig,
			UseFallbackTLS: false,
			User:           cfg.DB.PoolConfig.ConnConfig.User,
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
	pingCtx, pingCancel := context.WithTimeout(context.Background(), config.DefaultConnTimeout)
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
	username := p.cfg.DB.PoolConfig.ConnConfig.User
	hostname := p.cfg.DB.PoolConfig.ConnConfig.Host
	port := p.cfg.DB.PoolConfig.ConnConfig.Port
	database := p.cfg.DB.PoolConfig.ConnConfig.Database
	sslMode := "require"

	v := url.Values{}
	v.Set("sslmode", sslMode)
	v.Set("sslrootcert", p.cfg.DB.CAPath)
	v.Set("sslcert", p.cfg.DB.CertPath)
	v.Set("sslkey", p.cfg.DB.KeyPath)

	u := url.URL{
		Scheme:   p.cfg.DB.Scheme.String(),
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
