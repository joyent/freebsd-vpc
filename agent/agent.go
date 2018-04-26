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

package agent

import (
	"context"
	"net"
	"net/http"

	"github.com/joyent/freebsd-vpc/db"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Agent struct {
	config Config

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
	if err := a.rpcListener.Close(); err != nil {
		log.Warn().Err(err).Msg("error during RPC listener shutdown")
	}

	if err := a.rpcServer.Shutdown(context.Background()); err != nil {
		log.Warn().Err(err).Msg("error during RPC server shutdown")
	}

	if err := a.dbPool.Close(); err != nil {
		log.Warn().Err(err).Msg("error closing database pool")
	}

	return nil
}
