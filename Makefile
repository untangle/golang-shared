# default to vendor mod, since our minimal supported version of Go is
# 1.11
GOFLAGS ?= "-mod=vendor"
GO111MODULE ?= "on"
GOPRIVATE ?= GOPRIVATE=github.com/untangle/golang-shared

all: environment modules lint build-discoverd
build-discoverd:
	$(call LOG_FUNCTION,"Building discoverd...")
	cd cmd/discoverd ; \
	export GO111MODULE=$(GO111MODULE) ; \
	$(GOPRIVATE) go build $(GOFLAGS)

environment:
	$(call LOG_FUNCTION,"Setting up environment...")
	export $(GOPRIVATE)
	mkdir -p ~/.ssh/
	ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts

modules:
	$(call LOG_FUNCTION,"Vendoring modules...")
	$(GOPRIVATE) go mod vendor
#	$(GOPRIVATE) go mod tidy

lint:
	$(call LOG_FUNCTION,"Running golang linter...")
 	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.43.0 ; \
 	./bin/golangci-lint run

.PHONY: build lint environment
