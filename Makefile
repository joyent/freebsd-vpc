# NOTES:
#
# 0. make get-tools && pkg intsall cockroachdb
# 1. make crdb-gen-certs
# 2. In different terminals:
#   a. make crdb-start01
#   b. make crdb-start02
#   c. make crdb-start03
# 3. make crdb-mkdb
# 4. make
# 5. ./vpc migrate
# 6. make crdb-sql

build: generate
	mkdir ./bin
	govvv build -o bin/vpc ./cmd/vpc
	bin/vpc shell autocomplete bash -d docs/bash.d/

install:
	govvv install ./cmd/vpc

get-tools::
	go get -u github.com/ahmetb/govvv
	go get -u github.com/jteeuwen/go-bindata/...

DATA_DIR=`go env GOPATH`/src/github.com/sean-/vpc/crdb
CERT_DIR=`go env GOPATH`/src/github.com/sean-/vpc/crdb/certs
KEY_DIR=$(CERT_DIR)/keys

CRDB_HOST?=127.0.0.1
CRDB_PORT?=26257
CRDB_USER?=root
CRDB_DBNAME?=triton
CRDB_CERT_DSN?="sslmode=verify-ca"

generate:
	sh -c "cd db/migrations && go-bindata -o bindata.go -pkg migrations -ignore '(~|\.bak)$$' -prefix crdb crdb/"

crdb-mkdb:
	cockroach sql \
		--certs-dir="$(CERT_DIR)" \
		--host="$(CRDB_HOST)" \
		--port="$(CRDB_PORT)" \
		--user="$(CRDB_USER)" \
		--execute="CREATE DATABASE IF NOT EXISTS $(CRDB_DBNAME);"

crdb-sql:
	cockroach sql \
		--certs-dir="$(CERT_DIR)" \
		--host="$(CRDB_HOST)" \
		--port="$(CRDB_PORT)" \
		--user="$(CRDB_USER)" \
		--database="$(CRDB_DBNAME)"

crdb-gen-certs:
	mkdir -p "$(CERT_DIR)" "$(KEY_DIR)"
	chmod 0700 "$(KEY_DIR)"
	cockroach cert create-ca --certs-dir="$(CERT_DIR)" --ca-key="$(KEY_DIR)/ca.key"
	cockroach cert create-client root --certs-dir="$(CERT_DIR)" --ca-key="$(KEY_DIR)/ca.key"

	cockroach cert create-node 127.0.0.1 crdb01 --certs-dir="$(CERT_DIR)" --ca-key="$(KEY_DIR)/ca.key" --overwrite
	mkdir -p "$(DATA_DIR)/data-crdb01/certs/"
	mv "$(CERT_DIR)/node".* "$(DATA_DIR)/data-crdb01/certs/"
	cp "$(CERT_DIR)/ca.crt" "$(DATA_DIR)/data-crdb01/certs/"

	cockroach cert create-node 127.0.0.1 crdb02 --certs-dir="$(CERT_DIR)" --ca-key="$(KEY_DIR)/ca.key" --overwrite
	mkdir -p "$(DATA_DIR)/data-crdb02/certs/"
	mv "$(CERT_DIR)/node".* "$(DATA_DIR)/data-crdb02/certs/"
	cp "$(CERT_DIR)/ca.crt" "$(DATA_DIR)/data-crdb02/certs/"

	cockroach cert create-node 127.0.0.1 crdb03 --certs-dir="$(CERT_DIR)" --ca-key="$(KEY_DIR)/ca.key" --overwrite
	mkdir -p "$(DATA_DIR)/data-crdb03/certs/"
	mv "$(CERT_DIR)/node".* "$(DATA_DIR)/data-crdb03/certs/"
	cp "$(CERT_DIR)/ca.crt" "$(DATA_DIR)/data-crdb03/certs/"

crdb-start01:
	cockroach start \
		--certs-dir="$(DATA_DIR)/data-crdb01/certs" \
		--store="$(DATA_DIR)/data-crdb01" \
		--host=127.0.0.1 \
		--port=$(CRDB_PORT) \
		--http-port=8080 \
		--http-host=127.0.0.1 \
		--log-file-verbosity=INFO \
		--verbosity=1

crdb-start02:
	cockroach start \
		--certs-dir="$(DATA_DIR)/data-crdb02/certs" \
		--store="$(DATA_DIR)/data-crdb02" \
		--host=127.0.0.1 \
		--port=26258 \
		--http-port=8081 \
		--http-host=127.0.0.1 \
		--join=127.0.0.1:26257

crdb-start03:
	cockroach start \
		--certs-dir="$(DATA_DIR)/data-crdb01/certs" \
		--store="$(DATA_DIR)/data-crdb03" \
		--host=127.0.0.1 \
		--port=26259 \
		--http-port=8082 \
		--http-host=127.0.0.1 \
		--join=127.0.0.1:26257
