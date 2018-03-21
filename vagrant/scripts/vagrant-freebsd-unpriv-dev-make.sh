#!/bin/sh

set -o errexit

export GOPATH="/opt/gopath"
export PATH="/opt/gopath/bin:$PATH"

cd /opt/gopath/src/github.com/joyent/freebsd-vpc
make get-tools
make
bin/vpc db migrate
