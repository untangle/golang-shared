# default to vendor mod, since our minimal supported version of Go is
# 1.11
GOFLAGS ?= "-mod=vendor"
GO111MODULE ?= "on"
GOPRIVATE ?= GOPRIVATE=github.com/untangle/golang-shared

all: environment lint build-discoverd
build-discoverd:
	cd cmd/discoverd ; \
	export GO111MODULE=$(GO111MODULE) ; \
	$(GOPRIVATE) go build $(GOFLAGS)

environment:
	export $(GOPRIVATE)
	mkdir -p ~/.ssh/
	ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts


lint:
 	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.43.0 ; \
 	./bin/golangci-lint run

.PHONY: build lint environment
