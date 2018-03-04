#!/bin/sh

set -o errexit

export ASSUME_ALWAYS_YES=yes

pkg update
pkg install -y \
	databases/cockroach \
	editors/vim-console \
	security/ca_root_nss \
	shells/bash \
	sysutils/tmux \
	sysutils/tree

chsh -s /usr/local/bin/bash vagrant
chsh -s /usr/local/bin/bash root

mkdir -p /secrets
chown vagrant:vagrant /secrets
