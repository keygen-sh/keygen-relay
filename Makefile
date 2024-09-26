.PHONY: build dist test bench clean

LD_FLAGS := "-s -w"

# TODO: add needed platforms
PLATFORMS = linux darwin
ARCHITECTURES = amd64

build:
	go build -ldflags $(LD_FLAGS) -o bin/ ./cmd/...

build-all: clean
	@for platform in $(PLATFORMS); do \
		for arch in $(ARCHITECTURES); do \
			output="dist/relay-$$platform-$$arch"; \
			echo "Building for $$platform/$$arch..."; \
			env GOOS=$$platform GOARCH=$$arch go build -ldflags $(LD_FLAGS) -o $$output ./cmd/relay/main.go; \
		done; \
	done

test:
	go test ./...

# Run integration tests
test-integration:
	go test -tags=integrity ./...

bench:
	go test -bench=. -benchmem -run=^# ./...

clean:
	rm -rf bin dist

vet:
	go vet ./...

fmt:
	go fmt ./...
