# default to vendor mod, since our minimal supported version of Go is
# 1.11
GOFLAGS ?= "-mod=vendor"
GO111MODULE ?= "on"

all: build-discoverd
build-%:
	cd cmd/$* ; \
	export GO111MODULE=$(GO111MODULE) ; \
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.43.0 ; \
	./bin/golangci-lint run ; \
	export GOPRIVATE=github.com/untangle/golang-shared ; \
	go build $(GOFLAGS) -ldflags "-X main.Version=$(shell git describe --tags --always --long --dirty)"

lint:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	GO111MODULE=on GOPRIVATE=github.com/untangle/golang-shared $(shell go env GOPATH)/bin/golint -set_exit_status $(shell go list $(GOFLAGS) ./...)

.PHONY: build lint
