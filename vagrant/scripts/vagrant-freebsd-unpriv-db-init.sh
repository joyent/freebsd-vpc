#!/bin/sh

set -o errexit

node_internal=$(ifconfig vmx1 | grep inet | awk '{ print $2; }')

cockroach init --host="${node_internal}"
cockroach sql --host="${node_internal}" \
	--user root \
	--execute "CREATE DATABASE IF NOT EXISTS triton"
