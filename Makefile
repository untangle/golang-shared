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
	go build $(GOFLAGS) -ldflags "-X main.Version=$(shell git describe --tags --always --long --dirty)"

lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.43.0 ; \
	GO111MODULE=on GOPRIVATE=github.com/untangle/golang-shared ./bin/golangci-lint run 

.PHONY: build lint
