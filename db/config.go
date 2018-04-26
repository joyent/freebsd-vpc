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
