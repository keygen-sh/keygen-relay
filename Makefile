.PHONY: build build-386 build-amd64 build-arm build-arm64 build-all test test-integration test-all bench clean vet fmt

PACKAGE_NAME    := github.com/keygen-sh/keygen-relay
PACKAGE_VERSION := $(shell cat VERSION)
LD_FLAGS        := "-s -w -X $(PACKAGE_NAME)/cli.Version=$(PACKAGE_VERSION)"

build:
	go build -ldflags $(LD_FLAGS) -o bin/relay ./cmd/relay

build-386:
	CGO_ENABLED=1 GOOS=linux GOARCH=386 go build -ldflags $(LD_FLAGS) -o dist/relay_linux_386 ./cmd/relay

build-amd64:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags $(LD_FLAGS) -o dist/relay_linux_amd64 ./cmd/relay

build-arm:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 CC=arm-linux-gnueabihf-gcc go build -ldflags $(LD_FLAGS) -o dist/relay_linux_arm ./cmd/relay

build-arm64:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc go build -ldflags $(LD_FLAGS) -o dist/relay_linux_arm64 ./cmd/relay

build-all: clean build-386 build-amd64 build-arm build-arm64

test:
	go test -race ./...

test-integration:
	go test -tags=integrity ./...

test-all: test test-integration

bench:
	go test -bench=. -benchmem -run=^# ./...

clean:
	go clean
	rm -rf dist/*

vet:
	go vet ./...

fmt:
	go fmt ./...
