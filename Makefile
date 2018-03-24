RELEASE_DATE?=2018-02-28T23:59:59Z
RELEASE_EPOCH_TIME?=1519862399

GOPATH?=`go env GOPATH`

PROJECT_NAME?=github.com/joyent/freebsd-vpc
BUILDTIME_PATH?=$(PROJECT_NAME)/internal/buildtime
VPC_CMD_PATH=$(PROJECT_NAME)/cmd/vpc
GOVVV_FLAGS!=govvv -flags -pkg $(BUILDTIME_PATH)
GO_LDFLAGS?=-ldflags="$(VPC_CMD_PATH)=$(GOVVV_FLAGS) -X $(BUILDTIME_PATH).DocsDate=$(RELEASE_DATE)"

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
	mkdir -p ./bin
	go build $(GO_LDFLAGS) -o bin/vpc $(VPC_CMD_PATH)
	bin/vpc shell autocomplete bash -d docs/bash.d/ | cat
	bin/vpc docs man | cat
	bin/vpc docs md | cat

check:
	gometalinter \
		--deadline 10m \
		--vendor \
		--sort="path" \
		--aggregate \
		--enable-gc \
		--disable-all \
		--enable goimports \
		--enable misspell \
		--enable vet \
		--enable deadcode \
		--enable varcheck \
		--enable ineffassign \
		--enable structcheck \
		--enable unconvert \
		--enable gofmt \
		./...

install:
	go install $(GO_LDFLAGS) $(VPC_CMD_PATH)

get-tools:
	go get -u github.com/ahmetb/govvv
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u github.com/golang/dep/cmd/dep
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install

vagrant-box:
	go get -u github.com/sean-/cfgt
	go get -u github.com/hashicorp/packer
	cd vagrant/packer && cfgt --in template.json5 | \
		packer build -

DATA_DIR=`go env GOPATH`/src/github.com/joyent/freebsd-vpc/crdb
CERT_DIR=`go env GOPATH`/src/github.com/joyent/freebsd-vpc/crdb/certs
KEY_DIR=$(CERT_DIR)/keys

CRDB_HOST?=127.0.0.1
CRDB_PORT?=26257
CRDB_USER?=root
CRDB_DBNAME?=triton
CRDB_CERT_DSN?="sslmode=verify-ca"

generate:
	cd db/migrations && \
		$(GOPATH)/bin/go-bindata \
			-o bindata.go \
			-modtime $(RELEASE_EPOCH_TIME) \
			-pkg migrations \
			-ignore '(~|\.bak)$$' \
			-prefix crdb crdb/

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
