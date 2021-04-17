#@IgnoreInspection BashAddShebang

export APP=gurl

export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))

export LDFLAGS="-w -s"

all: format lint build

############################################################
# Build and Run
############################################################

build:
	CGO_ENABLED=1 go build -ldflags $(LDFLAGS)  ./cmd/kenny

install:
	CGO_ENABLED=1 go install -ldflags $(LDFLAGS) ./cmd/kenny

############################################################
# Format and Lint
############################################################

check-formatter:
	which goimports || GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports

format: check-formatter
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R goimports -w R
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R gofmt -s -w R

check-linter:
	which golangci-lint || GO111MODULE=off curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.23.8

lint: check-linter
	golangci-lint run $(ROOT)/...

############################################################
# Test
############################################################

test:
	go test -v -race -p 1 ./...