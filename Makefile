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

.PHONY: build
build:
	go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o bin/relay ./cmd/relay

.PHONY: build-linux-386
build-linux-386:
	CGO_ENABLED=1 GOOS=linux GOARCH=386 CC="zig cc -target x86-linux" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay-linux-386-$(PACKAGE_VERSION) ./cmd/relay

.PHONY: build-linux-amd64
build-linux-amd64:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC="zig cc -target x86_64-linux" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay-linux-amd64-$(PACKAGE_VERSION) ./cmd/relay

.PHONY: build-linux-arm
build-linux-arm:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 CC="zig cc -target arm-linux-gnueabihf" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay-linux-arm-$(PACKAGE_VERSION) ./cmd/relay

.PHONY: build-linux-arm64
build-linux-arm64:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC="zig cc -target aarch64-linux" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay-linux-arm64-$(PACKAGE_VERSION) ./cmd/relay

# .PHONY: build-linux-mips
# build-linux-mips:
# 	CGO_ENABLED=1 GOOS=linux GOARCH=mips GOMIPS=softfloat CC="zig cc -target mips-linux -mfloat-abi=soft" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay-linux-mips-$(PACKAGE_VERSION) ./cmd/relay

# .PHONY: build-linux-mipsle
# build-linux-mipsle:
# 	CGO_ENABLED=1 GOOS=linux GOARCH=mipsle GOMIPS=softfloat CC="zig cc -target mipsel-linux -mfloat-abi=soft" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay-linux-mipsle-$(PACKAGE_VERSION) ./cmd/relay

# .PHONY: build-linux-mips64
# build-linux-mips64:
# 	CGO_ENABLED=1 GOOS=linux GOARCH=mips64 GOMIPS=softfloat CC="zig cc -target mips64-linux -mfloat-abi=soft" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay-linux-mips64-$(PACKAGE_VERSION) ./cmd/relay

# .PHONY: build-linux-mips64le
# build-linux-mips64le:
# 	CGO_ENABLED=1 GOOS=linux GOARCH=mips64le GOMIPS=softfloat CC="zig cc -target mips64el-linux -mfloat-abi=soft" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay-linux-mips64le-$(PACKAGE_VERSION) ./cmd/relay

.PHONY: build-linux-s390x
build-linux-s390x:
	CGO_ENABLED=1 GOOS=linux GOARCH=s390x CC="zig cc -target s390x-linux" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay-linux-s390x-$(PACKAGE_VERSION) ./cmd/relay

.PHONY: build-darwin-amd64
build-darwin-amd64:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 CC="zig cc -target x86_64-macos" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay-darwin-amd64-$(PACKAGE_VERSION) ./cmd/relay

.PHONY: build-darwin-arm64
build-darwin-arm64:
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 CC="zig cc -target aarch64-macos" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay-darwin-arm64-$(PACKAGE_VERSION) ./cmd/relay

.PHONY: build-windows-386
build-windows-386:
	CGO_ENABLED=1 GOOS=windows GOARCH=386 CC="zig cc -target x86-windows" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay-windows-386-$(PACKAGE_VERSION).exe ./cmd/relay

.PHONY: build-windows-amd64
build-windows-amd64:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC="zig cc -target x86_64-windows" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay-windows-amd64-$(PACKAGE_VERSION).exe ./cmd/relay

.PHONY: build-windows-arm64
build-windows-arm64:
	CGO_ENABLED=1 GOOS=windows GOARCH=arm64 CC="zig cc -target aarch64-windows" go build $(BUILD_FLAGS) -ldflags $(BUILD_LDFLAGS) -o dist/relay-windows-arm64-$(PACKAGE_VERSION).exe ./cmd/relay

.PHONY: build-install
build-install:
	cp scripts/install.sh dist/install.sh

.PHONY: build-version
build-version:
	cp VERSION dist/version

.PHONY: build-all
build-all: clean build-linux-386 build-linux-amd64 build-linux-arm build-linux-arm64 build-darwin-amd64 \
	build-darwin-arm64 build-windows-386 build-windows-amd64 build-windows-arm64 build-install \
	build-version

.PHONY: release-new
release-new:
	keygen new --name "Keygen Relay v$(PACKAGE_VERSION)" --channel ${PACKAGE_CHANNEL} --version ${PACKAGE_VERSION}

.PHONY: release-linux-386
release-linux-386:
	keygen upload build/relay-linux-386-$(PACKAGE_VERSION) --filename relay_linux_386 --release ${PACKAGE_VERSION} --platform linux --arch 386

.PHONY: release-linux-amd64
release-linux-amd64:
	keygen upload build/relay-linux-amd64-$(PACKAGE_VERSION) --filename relay_linux_amd64 --release ${PACKAGE_VERSION} --platform linux --arch amd64

.PHONY: release-linux-arm
release-linux-arm:
	keygen upload build/relay-linux-arm-$(PACKAGE_VERSION) --filename relay_linux_arm --release ${PACKAGE_VERSION} --platform linux --arch arm

.PHONY: release-linux-arm64
release-linux-arm64:
	keygen upload build/relay-linux-arm64-$(PACKAGE_VERSION) --filename relay_linux_arm64 --release ${PACKAGE_VERSION} --platform linux --arch arm64

# .PHONY: release-linux-mips
# release-linux-mips:
# 	keygen upload build/relay-linux-mips-$(PACKAGE_VERSION) --filename relay_linux_mips --release ${PACKAGE_VERSION} --platform linux --arch mips

# .PHONY: release-linux-mipsle
# release-linux-mipsle:
# 	keygen upload build/relay-linux-mipsle-$(PACKAGE_VERSION) --filename relay_linux_mipsle --release ${PACKAGE_VERSION} --platform linux --arch mipsle

# .PHONY: release-linux-mips64
# release-linux-mips64:
# 	keygen upload build/relay-linux-mips64-$(PACKAGE_VERSION) --filename relay_linux_mips64 --release ${PACKAGE_VERSION} --platform linux --arch mips64

# .PHONY: release-linux-mips64le
# release-linux-mips64le:
# 	keygen upload build/relay-linux-mips64le-$(PACKAGE_VERSION) --filename relay_linux_mips64le --release ${PACKAGE_VERSION} --platform linux --arch mips64le

.PHONY: release-linux-s390x
release-linux-s390x:
	keygen upload build/relay-linux-s390x-$(PACKAGE_VERSION) --filename relay_linux_s390x --release ${PACKAGE_VERSION} --platform linux --arch s390x

.PHONY: release-darwin-amd64
release-darwin-amd64:
	keygen upload build/relay-darwin-amd64-$(PACKAGE_VERSION) --filename relay_darwin_amd64 --release ${PACKAGE_VERSION} --platform darwin --arch amd64

.PHONY: release-darwin-arm64
release-darwin-arm64:
	keygen upload build/relay-darwin-arm64-$(PACKAGE_VERSION) --filename relay_darwin_arm64 --release ${PACKAGE_VERSION} --platform darwin --arch arm64

.PHONY: release-windows-386
release-windows-386:
	keygen upload build/relay-windows-386-$(PACKAGE_VERSION).exe --filename relay_windows_386.exe --release ${PACKAGE_VERSION} --platform windows --arch 386

.PHONY: release-windows-amd64
release-windows-amd64:
	keygen upload build/relay-windows-amd64-$(PACKAGE_VERSION).exe --filename relay_windows_amd64.exe --release ${PACKAGE_VERSION} --platform windows --arch amd64

.PHONY: release-windows-arm64
release-windows-arm64:
	keygen upload build/relay-windows-arm64-$(PACKAGE_VERSION).exe --filename relay_windows_arm64.exe --release ${PACKAGE_VERSION} --platform windows --arch arm64

.PHONY: release-installer
release-installer:
	keygen upload build/install.sh --release ${PACKAGE_VERSION}

.PHONY: release-version
release-version:
	keygen upload build/version --release ${PACKAGE_VERSION} --filetype txt

.PHONY: release-publish
release-publish:
	keygen publish --release ${PACKAGE_VERSION}

.PHONY: release-tag
release-tag:
	ifeq ($(PACKAGE_CHANNEL),stable)
		keygen untag --release latest
		keygen tag latest --release ${PACKAGE_VERSION}
	endif

.PHONY: release
release: release-new release-linux-386 release-linux-amd64 release-linux-arm release-linux-arm64 release-darwin-amd64 \
	release-darwin-arm64 release-windows-386 release-windows-amd64 release-windows-arm64 release-installer \
	release-version release-publish release-tag

.PHONY: test
test:
	go test -race ./...

.PHONY: test-integration
test-integration:
	go test -tags=integrity ./...

.PHONY: test-all
test-all: test test-integration

.PHONY: bench
bench:
	go test -bench=. -benchmem -run=^# ./...

.PHONY: clean
clean:
	go clean
	rm -rf dist/*

.PHONY: vet
vet:
	@go vet ./...

.PHONY: fmt
fmt:
	@go fmt ./...
