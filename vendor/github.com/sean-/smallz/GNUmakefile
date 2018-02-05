BIN ?= manta
GOCOVER_TMPFILE?=       $(GOCOVER_FILE).tmp
GOCOVER_FILE?=  .cover.out
GOCOVERHTML?=   coverage.html
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
FIND=`/usr/bin/which 2> /dev/null gfind find | /usr/bin/grep -v ^no | /usr/bin/head -n 1`
XARGS=`/usr/bin/which 2> /dev/null gxargs xargs | /usr/bin/grep -v ^no | /usr/bin/head -n 1`

default:: help

.PHONY: build
build:: $(BIN) ## 10 Build $(BIN) binary

.PHONY: manta
manta::
	go build -ldflags "-X main.commit=`git describe --tags --always` -X main.date=`date +%Y-%m-%d_%H:%d`" -o $@ main.go

.PHONY: check
check:: ## 10 Run go test
	go test -v ./...

cover:: coverage_report ## 10 Generate a coverage report

$(GOCOVER_FILE)::
	@${FIND} . -type d ! -path '*vendor*' ! -path '*.git*' ! -path '*pgdata*' -print0 | ${XARGS} -0 -P4 -I % sh -ec "cd % && rm -f $(GOCOVER_TMPFILE) && go test -coverprofile=$(GOCOVER_TMPFILE) || true"

	@echo 'mode: set' > $(GOCOVER_FILE)
	@${FIND} . -type f ! -path '*vendor*' ! -path '*.git*' ! -path '*pgdata*' -name "$(GOCOVER_TMPFILE)" -print0 | ${XARGS} -0 -n1 cat $(GOCOVER_TMPFILE) | grep -v '^mode: ' >> ${PWD}/$(GOCOVER_FILE)

$(GOCOVERHTML): $(GOCOVER_FILE)
	go tool cover -html=$(GOCOVER_FILE) -o $(GOCOVERHTML)

coverage_report:: $(GOCOVER_FILE)
	go tool cover -html=$(GOCOVER_FILE)

install_audit_tools:: ## 10 Install static analysis tools
	@go get -u github.com/golang/lint/golint && echo "Installed golint"
	@go get -u github.com/fzipp/gocyclo && echo "Installed gocyclo"
	@go get -u github.com/remyoudompheng/go-misc/deadcode && echo "Installed deadcode"
	@go get -u github.com/client9/misspell/cmd/misspell && echo "Installed misspell"
	@go get -u github.com/gordonklaus/ineffassign && echo "Installed ineffassig:"
	@go get -u golang.org/x/tools/cover && echo "Installed cover"

audit:: ## 10 Run static analysis tools over the code
	deadcode
	go tool vet -all *.go
	go tool vet -shadow=true *.go
	golint *.go
	ineffassign .
	gocyclo -over 65 *.go
	misspell *.go


.PHONY: vet
vet:: ## 10 vet the binary (excluding dependencies)
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

.PHONY: fmt
fmt: ## 10 fmt and simplify the code
	gofmt -s -w $(GOFMT_FILES)

.PHONY: vendor-status
vendor-status: ## 10 Display the vendor/ status
	@dep status

.PHONY: release
release: ## 10 Build a release
	#goreleaser --release-notes=release_notes.md
	@goreleaser

.PHONY: release-snapshot
release-snapshot: ## 10 Build a snapshot release
	@goreleaser --snapshot --skip-validate --rm-dist

.PHONY: clean
clean:: cleandb-shard ## 90 Clean target
	find . -name "$(GOCOVER_FILE).tmp" -o -name "$(GOCOVER_FILE)" -delete

.PHONY: help
help:: ## 99 This help message
	@echo "$(BIN) make(1) targets:"
	@grep -E '^[a-zA-Z\_\-]+:[:]?.*?## [0-9]+ .*$$' $(MAKEFILE_LIST) | \
		sort -n -t '#' -k3,1 | awk '				\
BEGIN { FS = ":[:]?.*?## "; section = 10; };				\
{									\
	newSect = int($$2);						\
	if (section != newSect) {					\
		section = newSect;					\
		printf "\n";						\
	}								\
	sub("^[0-9]+", "",  $$2);					\
	printf "\033[36m%-15s\033[0m %s\n", $$1, $$2;			\
}'
