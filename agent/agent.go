package agent

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/joyent/freebsd-vpc/db"
	"github.com/pkg/errors"
	log "github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

type Agent struct {
	config Config
	logger *log.Logger

	dbPool *db.Pool

	rpcListener net.Listener
	rpcServer   *http.Server
}

func New(config Config) (agent *Agent, err error) {
	dbPool, err := db.New(config.DBConfig)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create database pool")
	}

	rpcListener, err := net.Listen("unix", config.AgentConfig.Addresses.Internal)
	if err != nil {
		return nil, errors.Wrap(err, "error creating RPC listener")
	}

	rpcServer := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hlog.FromRequest(r).Info().Msg("got request")

			w.WriteHeader(200)
			w.Write([]byte("Hello World"))
		}),
	}

	return &Agent{
		dbPool:      dbPool,
		rpcListener: rpcListener,
		rpcServer:   rpcServer,
	}, nil
}

func (a *Agent) Start() error {
	if err := a.dbPool.Ping(); err != nil {
		return errors.Wrap(err, "unable to ping database")
	}

	go a.rpcServer.Serve(a.rpcListener)

	return nil
}

func (a *Agent) Shutdown() error {
	if err := a.rpcServer.Shutdown(context.Background()); err != nil {
		a.logger.Warn().Err(err).Msg("error during RPC server shutdown")
	}

	if err := a.rpcListener.Close(); err != nil {
		a.logger.Warn().Err(err).Msg("error during RPC listener shutdown")
	}

	if err := a.dbPool.Close(); err != nil {
		a.logger.Warn().Err(err).Msg("error closing database pool")
	}

	if err := os.Remove(a.config.AgentConfig.Addresses.Internal); err != nil {
		a.logger.Warn().Err(err).Msg("error removing domain socket file")
	}

	a.logger.Info().Msg("graceful shutdown complete")

	return nil
}
