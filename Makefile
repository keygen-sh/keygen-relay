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
BUILD_FLAGS    = -v

ifdef DEBUG
	BUILD_FLAGS += -x
endif

build:
	go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o bin/relay ./cmd/relay

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

build-windows-arm64:
	CGO_ENABLED=1 GOOS=windows GOARCH=arm64 CC="zig cc -target aarch64-windows" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay_windows_arm64.exe ./cmd/relay

build-install:
	cp scripts/install.sh dist/install.sh

build-version:
	cp VERSION dist/version

build-all: clean build-linux-386 build-linux-amd64 build-linux-arm build-linux-arm64 build-darwin-amd64 \
	build-darwin-arm64 build-windows-386 build-windows-amd64 build-windows-arm64 build-install \
	build-version

release-new:
	keygen new --name "Keygen Relay v$(PACKAGE_VERSION)" --channel ${PACKAGE_CHANNEL} --version ${PACKAGE_VERSION}

release-linux-386:
	keygen upload build/relay_linux_386 --release ${PACKAGE_VERSION} --platform linux --arch 386

release-linux-amd64:
	keygen upload build/relay_linux_amd64 --release ${PACKAGE_VERSION} --platform linux --arch amd64

release-linux-arm:
	keygen upload build/relay_linux_arm --release ${PACKAGE_VERSION} --platform linux --arch arm

release-linux-arm64:
	keygen upload build/relay_linux_arm64 --release ${PACKAGE_VERSION} --platform linux --arch arm64

release-darwin-amd64:
	keygen upload build/relay_darwin_amd64 --release ${PACKAGE_VERSION} --platform darwin --arch amd64

release-darwin-arm64:
	keygen upload build/relay_darwin_arm64 --release ${PACKAGE_VERSION} --platform darwin --arch arm64

release-windows-386:
	keygen upload build/relay_windows_386.exe --release ${PACKAGE_VERSION} --platform windows --arch 386

release-windows-amd64:
	keygen upload build/relay_windows_amd64.exe --release ${PACKAGE_VERSION} --platform windows --arch amd64

release-windows-arm64:
	keygen upload build/relay_windows_arm64.exe --release ${PACKAGE_VERSION} --platform windows --arch arm64

release-installer:
	keygen upload build/install.sh --release ${PACKAGE_VERSION}

release-version:
	keygen upload build/version --release ${PACKAGE_VERSION} --filetype txt

release-publish:
	keygen publish --release ${PACKAGE_VERSION}

release-tag:
	ifeq ($(PACKAGE_CHANNEL),stable)
		keygen untag --release latest
		keygen tag latest --release ${PACKAGE_VERSION}
	endif

release: release-new release-linux-386 release-linux-amd64 release-linux-arm release-linux-arm64 release-darwin-amd64 \
	release-darwin-arm64 release-windows-386 release-windows-amd64 release-windows-arm64 release-installer \
	release-version release-publish release-tag

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
