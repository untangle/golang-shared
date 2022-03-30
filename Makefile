# currently the only thing we do in this Makefile is run lint, and build the protocol buffer structures
all: lint compile-protobuffs

compile-protobuffs:
	rm -rf structs/protocolbuffers/*
	protoc --proto_path=protobuffersrc --go_out=. --go-grpc_out=protobuffersrc --go-grpc_opt=paths=source_relative --go_opt=module=github.com/untangle/golang-shared protobuffersrc/*

lint:
	GO111MODULE=off go get -u golang.org/x/lint/golint
	$(shell go env GOPATH)/bin/golint -set_exit_status $(shell go list $(GOFLAGS) ./...)