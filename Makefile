.PHONY: build build-386 build-amd64 build-arm build-arm64 build-all test test-integration \
	test-all bench clean vet fmt

PACKAGE_NAME    := github.com/keygen-sh/keygen-relay
PACKAGE_VERSION := $(shell cat VERSION)
PACKAGE_CHANNEL  = stable

ifeq ($(findstring -rc.,$(PACKAGE_VERSION)), -rc.)
	PACKAGE_CHANNEL = rc
else ifeq ($(findstring -beta.,$(PACKAGE_VERSION)), -beta.)
	PACKAGE_CHANNEL = beta
else ifeq ($(findstring -alpha.,$(PACKAGE_VERSION)), -alpha.)
	PACKAGE_CHANNEL = alpha
else ifeq ($(findstring -dev.,$(PACKAGE_VERSION)), -dev.)
	PACKAGE_CHANNEL = dev
endif

BUILD_LDFLAGS := "-s -w -X $(PACKAGE_NAME)/cli.Version=$(PACKAGE_VERSION)"
BUILD_FLAGS    =

ifdef DEBUG
	BUILD_FLAGS += -x
endif

build:
	go build -ldflags $(BUILD_LDFLAGS) -o bin/relay ./cmd/relay

build-linux-386:
	CGO_ENABLED=1 GOOS=linux GOARCH=386 CC="zig cc -target x86-linux" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay_linux_386 ./cmd/relay

build-linux-amd64:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC="zig cc -target x86_64-linux" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay_linux_amd64 ./cmd/relay

build-linux-arm:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 CC="zig cc -target arm-linux-gnueabihf" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay_linux_arm ./cmd/relay

build-linux-arm64:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC="zig cc -target aarch64-linux" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay_linux_arm64 ./cmd/relay

build-darwin-amd64:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 CC="zig cc -target x86_64-macos" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay_darwin_amd64 ./cmd/relay

build-darwin-arm64:
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 CC="zig cc -target aarch64-macos" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay_darwin_arm64 ./cmd/relay

build-windows-386:
	CGO_ENABLED=1 GOOS=windows GOARCH=386 CC="zig cc -target x86-windows" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay_windows_386.exe ./cmd/relay

build-windows-amd64:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC="zig cc -target x86_64-windows" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay_windows_amd64.exe ./cmd/relay

build-windows-arm:
	CGO_ENABLED=1 GOOS=windows GOARCH=arm GOARM=7 CC="zig cc -target arm-windows-gnueabihf" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay_windows_arm.exe ./cmd/relay

build-windows-arm64:
	CGO_ENABLED=1 GOOS=windows GOARCH=arm64 CC="zig cc -target aarch64-windows" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay_windows_arm64.exe ./cmd/relay

build-install:
	cp scripts/install.sh dist/install.sh

build-version:
	cp VERSION dist/version

build-all: clean build-linux-386 build-linux-amd64 build-linux-arm build-linux-arm64 build-darwin-amd64 \
	build-darwin-arm64 build-windows-386 build-windows-amd64 build-windows-arm build-windows-arm64 \
	build-install build-version

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
