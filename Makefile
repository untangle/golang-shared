# default to vendor mod, since our minimal supported version of Go is
# 1.11
GOFLAGS ?= "-mod=vendor"
GO111MODULE ?= "on"
GOPRIVATE ?= GOPRIVATE=github.com/untangle/golang-shared
EXTRA_TEST_FLAGS ?=
GOTEST_COVERAGE ?= yes
GO_COVERPROFILE ?= /tmp/packetd_coverage.out
COVERAGE_HTML ?= /tmp/packetd_coverage.html
BROWSER ?= open

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

all: build lint
build:
	$(call LOG_FUNCTION,"Compiling protocol buffers...")
	rm -rf structs/protocolbuffers/*
	protoc --proto_path=protobuffersrc --go_out=. --go_opt=module=github.com/untangle/golang-shared protobuffersrc/*

environment:
	$(call LOG_FUNCTION,"Setting up environment...")
	export $(GOPRIVATE)
	mkdir -p ~/.ssh/
	ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts
	git config --global url.ssh://git@github.com/.insteadOf https://github.com/

modules: environment
	$(call LOG_FUNCTION,"Vendoring modules...")
	$(GOPRIVATE) go mod vendor

lint: modules
	$(call LOG_FUNCTION,"Running golang linter...")
	cd /tmp; GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2
	$(shell go env GOPATH)/bin/golangci-lint --version
# IMPORTANT --issues-exit-code 0 will let the build continue without failing lint checks - this should be removed eventually
	$(shell go env GOPATH)/bin/golangci-lint run --issues-exit-code 0

test: build
	if [ $(GOTEST_COVERAGE) = "yes" ]; \
	then \
		go test -vet=off $(EXTRA_TEST_FLAGS) -coverprofile=$(GO_COVERPROFILE) ./...; \
	else \
		go test -vet=off $(EXTRA_TEST_FLAGS) ./...; \
	fi

racetest: EXTRA_TEST_FLAGS=-race
racetest: test

browsecoverage: test
	go tool cover -html=$(GO_COVERPROFILE) -o $(COVERAGE_HTML)
	$(BROWSER) $(COVERAGE_HTML)


funccoverage: test
	go tool cover -func $(GO_COVERPROFILE)


funccoverage: test
	go tool cover -func $(GO_COVERPROFILE)
.PHONY: build lint environment
