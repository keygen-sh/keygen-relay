PACKAGE_NAME    := github.com/keygen-sh/keygen-relay
PACKAGE_VERSION := $(shell cat VERSION)
PACKAGE_MAJOR   := $(shell cat VERSION | cut -d. -f1)
PACKAGE_MINOR   := $(shell cat VERSION | cut -d. -f2)
PACKAGE_PATCH   := $(shell cat VERSION | cut -d. -f3 | cut -d- -f1)
PACKAGE_PRE     := $(shell cat VERSION | grep -oP '(?<=-)[^+]*')
PACKAGE_BUILD   := $(shell cat VERSION | grep -oP '(?<=\+).*')
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

BUILD_LDFLAGS := -s -w -X $(PACKAGE_NAME)/cli.Version=$(PACKAGE_VERSION)
BUILD_FLAGS    = -v

ifdef BUILD_NODE_LOCKED
	BUILD_LDFLAGS += -X $(PACKAGE_NAME)/internal/locker.Fingerprint=$(BUILD_NODE_LOCKED_FINGERPRINT)
	BUILD_LDFLAGS += -X $(PACKAGE_NAME)/internal/locker.PublicKey=$(BUILD_NODE_LOCKED_PUBLIC_KEY)

	ifdef BUILD_NODE_LOCKED_PLATFORM
		BUILD_LDFLAGS += -X $(PACKAGE_NAME)/internal/locker.Platform=$(BUILD_NODE_LOCKED_PLATFORM)
	endif

	ifdef BUILD_NODE_LOCKED_HOSTNAME
		BUILD_LDFLAGS += -X $(PACKAGE_NAME)/internal/locker.Hostname=$(BUILD_NODE_LOCKED_HOSTNAME)
	endif

	ifdef BUILD_NODE_LOCKED_IP
		BUILD_LDFLAGS += -X $(PACKAGE_NAME)/internal/locker.IP=$(BUILD_NODE_LOCKED_IP)
	endif

	ifdef BUILD_NODE_LOCKED_ADDR
		BUILD_LDFLAGS += -X $(PACKAGE_NAME)/internal/locker.Addr=$(BUILD_NODE_LOCKED_ADDR)
	endif

	ifdef BUILD_NODE_LOCKED_PORT
		BUILD_LDFLAGS += -X $(PACKAGE_NAME)/internal/locker.Port=$(BUILD_NODE_LOCKED_PORT)
	endif
endif

ifdef DEBUG
	BUILD_FLAGS += -x
endif

.PHONY: generate
generate:
	go generate ./...
	sqlc generate

.PHONY: build
build: clean generate
	go build $(BUILD_FLAGS) -ldflags "$(BUILD_LDFLAGS)" -o bin/relay ./cmd/relay

.PHONY: build-linux-386
build-linux-386:
	CGO_ENABLED=1 GOOS=linux GOARCH=386 CC="zig cc -target x86-linux" go build $(BUILD_FLAGS) -ldflags "$(BUILD_LDFLAGS)" -o dist/relay-linux-386-$(PACKAGE_VERSION) ./cmd/relay

.PHONY: build-linux-amd64
build-linux-amd64:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC="zig cc -target x86_64-linux" go build $(BUILD_FLAGS) -ldflags "$(BUILD_LDFLAGS)" -o dist/relay-linux-amd64-$(PACKAGE_VERSION) ./cmd/relay

.PHONY: build-linux-arm
build-linux-arm:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 CC="zig cc -target arm-linux-gnueabihf" go build $(BUILD_FLAGS) -ldflags "$(BUILD_LDFLAGS)" -o dist/relay-linux-arm-$(PACKAGE_VERSION) ./cmd/relay

.PHONY: build-linux-arm64
build-linux-arm64:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC="zig cc -target aarch64-linux" go build $(BUILD_FLAGS) -ldflags "$(BUILD_LDFLAGS)" -o dist/relay-linux-arm64-$(PACKAGE_VERSION) ./cmd/relay

# .PHONY: build-linux-mips
# build-linux-mips:
# 	CGO_ENABLED=1 GOOS=linux GOARCH=mips GOMIPS=softfloat CC="zig cc -target mips-linux -mfloat-abi=soft" go build $(BUILD_FLAGS) -ldflags "$(BUILD_LDFLAGS)" -o dist/relay-linux-mips-$(PACKAGE_VERSION) ./cmd/relay

# .PHONY: build-linux-mipsle
# build-linux-mipsle:
# 	CGO_ENABLED=1 GOOS=linux GOARCH=mipsle GOMIPS=softfloat CC="zig cc -target mipsel-linux -mfloat-abi=soft" go build $(BUILD_FLAGS) -ldflags "$(BUILD_LDFLAGS)" -o dist/relay-linux-mipsle-$(PACKAGE_VERSION) ./cmd/relay

# .PHONY: build-linux-mips64
# build-linux-mips64:
# 	CGO_ENABLED=1 GOOS=linux GOARCH=mips64 GOMIPS=softfloat CC="zig cc -target mips64-linux -mfloat-abi=soft" go build $(BUILD_FLAGS) -ldflags "$(BUILD_LDFLAGS)" -o dist/relay-linux-mips64-$(PACKAGE_VERSION) ./cmd/relay

# .PHONY: build-linux-mips64le
# build-linux-mips64le:
# 	CGO_ENABLED=1 GOOS=linux GOARCH=mips64le GOMIPS=softfloat CC="zig cc -target mips64el-linux -mfloat-abi=soft" go build $(BUILD_FLAGS) -ldflags "$(BUILD_LDFLAGS)" -o dist/relay-linux-mips64le-$(PACKAGE_VERSION) ./cmd/relay

.PHONY: build-linux-s390x
build-linux-s390x:
	CGO_ENABLED=1 GOOS=linux GOARCH=s390x CC="zig cc -target s390x-linux" go build $(BUILD_FLAGS) -ldflags "$(BUILD_LDFLAGS)" -o dist/relay-linux-s390x-$(PACKAGE_VERSION) ./cmd/relay

.PHONY: build-darwin-amd64
build-darwin-amd64:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 CC="zig cc -target x86_64-macos -g0 --sysroot=./sdk -I./sdk/usr/include -L./sdk/usr/lib -F./sdk/System/Library/Frameworks" go build $(BUILD_FLAGS) -ldflags "$(BUILD_LDFLAGS)" -o dist/relay-darwin-amd64-$(PACKAGE_VERSION) ./cmd/relay

.PHONY: build-darwin-arm64
build-darwin-arm64:
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 CC="zig cc -target aarch64-macos -g0 --sysroot=./sdk -I./sdk/usr/include -L./sdk/usr/lib -F./sdk/System/Library/Frameworks" go build $(BUILD_FLAGS) -ldflags "$(BUILD_LDFLAGS)" -o dist/relay-darwin-arm64-$(PACKAGE_VERSION) ./cmd/relay

.PHONY: build-windows-386
build-windows-386:
	CGO_ENABLED=1 GOOS=windows GOARCH=386 CC="zig cc -target x86-windows" go build $(BUILD_FLAGS) -ldflags "$(BUILD_LDFLAGS)" -o dist/relay-windows-386-$(PACKAGE_VERSION).exe ./cmd/relay

.PHONY: build-windows-amd64
build-windows-amd64:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC="zig cc -target x86_64-windows" go build $(BUILD_FLAGS) -ldflags "$(BUILD_LDFLAGS)" -o dist/relay-windows-amd64-$(PACKAGE_VERSION).exe ./cmd/relay

.PHONY: build-windows-arm64
build-windows-arm64:
	CGO_ENABLED=1 GOOS=windows GOARCH=arm64 CC="zig cc -target aarch64-windows" go build $(BUILD_FLAGS) -ldflags "$(BUILD_LDFLAGS)" -o dist/relay-windows-arm64-$(PACKAGE_VERSION).exe ./cmd/relay

.PHONY: build-installer
build-installer:
	cp scripts/install.sh dist/install.sh

.PHONY: build-version
build-version:
	cp VERSION dist/version

.PHONY: build-image
build-image:
	docker buildx build --platform "linux/amd64,linux/arm64" --output type=oci,dest=- . > dist/relay-$(PACKAGE_VERSION).tar

.PHONY: build-all
build-all: clean generate build-linux-386 build-linux-amd64 build-linux-arm build-linux-arm64 build-linux-s390x \
	build-windows-386 build-windows-amd64 build-windows-arm64 build-installer build-version build-image

# TODO(ezekg) should we refactor into relay-cli and relay-oci packages?
.PHONY: release-new
release-new:
	keygen new --name "Keygen Relay v$(PACKAGE_VERSION)" --channel $(PACKAGE_CHANNEL) --version $(PACKAGE_VERSION)

.PHONY: release-linux-386
release-linux-386:
	keygen upload dist/relay-linux-386-$(PACKAGE_VERSION) --filename relay_linux_386 --release $(PACKAGE_VERSION) --platform linux --arch 386

.PHONY: release-linux-amd64
release-linux-amd64:
	keygen upload dist/relay-linux-amd64-$(PACKAGE_VERSION) --filename relay_linux_amd64 --release $(PACKAGE_VERSION) --platform linux --arch amd64

.PHONY: release-linux-arm
release-linux-arm:
	keygen upload dist/relay-linux-arm-$(PACKAGE_VERSION) --filename relay_linux_arm --release $(PACKAGE_VERSION) --platform linux --arch arm

.PHONY: release-linux-arm64
release-linux-arm64:
	keygen upload dist/relay-linux-arm64-$(PACKAGE_VERSION) --filename relay_linux_arm64 --release $(PACKAGE_VERSION) --platform linux --arch arm64

# .PHONY: release-linux-mips
# release-linux-mips:
# 	keygen upload dist/relay-linux-mips-$(PACKAGE_VERSION) --filename relay_linux_mips --release $(PACKAGE_VERSION) --platform linux --arch mips

# .PHONY: release-linux-mipsle
# release-linux-mipsle:
# 	keygen upload dist/relay-linux-mipsle-$(PACKAGE_VERSION) --filename relay_linux_mipsle --release $(PACKAGE_VERSION) --platform linux --arch mipsle

# .PHONY: release-linux-mips64
# release-linux-mips64:
# 	keygen upload dist/relay-linux-mips64-$(PACKAGE_VERSION) --filename relay_linux_mips64 --release $(PACKAGE_VERSION) --platform linux --arch mips64

# .PHONY: release-linux-mips64le
# release-linux-mips64le:
# 	keygen upload dist/relay-linux-mips64le-$(PACKAGE_VERSION) --filename relay_linux_mips64le --release $(PACKAGE_VERSION) --platform linux --arch mips64le

.PHONY: release-linux-s390x
release-linux-s390x:
	keygen upload dist/relay-linux-s390x-$(PACKAGE_VERSION) --filename relay_linux_s390x --release $(PACKAGE_VERSION) --platform linux --arch s390x

.PHONY: release-darwin-amd64
release-darwin-amd64:
	keygen upload dist/relay-darwin-amd64-$(PACKAGE_VERSION) --filename relay_darwin_amd64 --release $(PACKAGE_VERSION) --platform darwin --arch amd64

.PHONY: release-darwin-arm64
release-darwin-arm64:
	keygen upload dist/relay-darwin-arm64-$(PACKAGE_VERSION) --filename relay_darwin_arm64 --release $(PACKAGE_VERSION) --platform darwin --arch arm64

.PHONY: release-windows-386
release-windows-386:
	keygen upload dist/relay-windows-386-$(PACKAGE_VERSION).exe --filename relay_windows_386.exe --release $(PACKAGE_VERSION) --platform windows --arch 386

.PHONY: release-windows-amd64
release-windows-amd64:
	keygen upload dist/relay-windows-amd64-$(PACKAGE_VERSION).exe --filename relay_windows_amd64.exe --release $(PACKAGE_VERSION) --platform windows --arch amd64

.PHONY: release-windows-arm64
release-windows-arm64:
	keygen upload dist/relay-windows-arm64-$(PACKAGE_VERSION).exe --filename relay_windows_arm64.exe --release $(PACKAGE_VERSION) --platform windows --arch arm64

.PHONY: release-installer
release-installer:
	keygen upload dist/install.sh --release $(PACKAGE_VERSION)

.PHONY: release-version
release-version:
	keygen upload dist/version --release $(PACKAGE_VERSION) --filetype txt

.PHONY: release-publish
release-publish:
	keygen publish --release $(PACKAGE_VERSION)

.PHONY: release-tag
release-tag:
ifeq ($(PACKAGE_CHANNEL),stable)
	keygen untag --release latest
	keygen tag latest --release $(PACKAGE_VERSION)
endif

.PHONY: release-image-new
release-image-new:
	keygen new --name "Keygen Relay v$(PACKAGE_VERSION)" --channel $(PACKAGE_CHANNEL) --version $(PACKAGE_VERSION) --package relay

.PHONY: release-image-tarball
release-image-tarball:
	keygen upload dist/relay-$(PACKAGE_VERSION).tar --filename relay.tar --release $(PACKAGE_VERSION) --package relay

.PHONY: release-image-tag
release-image-tag:
ifeq ($(PACKAGE_CHANNEL),stable)
	keygen untag latest v$(PACKAGE_MAJOR) v$(PACKAGE_MAJOR).$(PACKAGE_MINOR) --release latest --package relay
	keygen tag latest v$(PACKAGE_MAJOR) v$(PACKAGE_MAJOR).$(PACKAGE_MINOR) --release $(PACKAGE_VERSION) --package relay
endif

.PHONY: release-image
release-image: release-image-new release-image-tarball release-image-tag

# FIXME(ezekg) refactor into release-cli and release-oci recipes
.PHONY: release
release: release-new release-linux-386 release-linux-amd64 release-linux-arm release-linux-arm64 \
	release-linux-s390x release-windows-386 release-windows-amd64 release-windows-arm64 \
	release-installer release-version release-publish release-tag release-image

.PHONY: test
test:
	go test -v -race ./...

.PHONY: test-integration
test-integration:
	go test -v -tags=integration ./...

.PHONY: test-all
test-all: test test-integration

.PHONY: bench
bench:
	go test -bench=. -benchmem -run=^# ./...

.PHONY: clean
clean:
	go clean
	rm -rf dist/* bin/*
	mkdir -p dist bin

.PHONY: vet
vet:
	@go vet ./...

.PHONY: fmt
fmt:
	@go fmt ./...
