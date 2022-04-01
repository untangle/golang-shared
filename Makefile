# default to vendor mod, since our minimal supported version of Go is
# 1.11
GOFLAGS ?= "-mod=vendor"
GO111MODULE ?= "on"
GOPRIVATE ?= GOPRIVATE=github.com/untangle/golang-shared

# logging
NC := "\033[0m" # no color
YELLOW := "\033[1;33m"
ifneq ($(DEV),false)
  GREEN := "\033[1;32m"
else
  GREEN :=
endif
LOG_FUNCTION = @/bin/echo -e $(shell date +%T.%3N) $(GREEN)$(1)$(NC)
WARN_FUNCTION = @/bin/echo -e $(shell date +%T.%3N) $(YELLOW)$(1)$(NC)

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

lint:
	$(call LOG_FUNCTION,"Running golang linter...")
	cd /tmp; GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.23.8
	$(shell go env GOPATH)/bin/golangci-lint --version
	$(shell go env GOPATH)/bin/golangci-lint run

.PHONY: build lint environment
