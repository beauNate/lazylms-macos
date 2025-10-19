.PHONY: all build clean install test run help

BINARY_NAME=lazylms-macos
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

all: build

help:
	@echo "lazylms-macos - Mac OS 26 Liquid Glass TUI for LM Studio"
	@echo ""
	@echo "Usage:"
	@echo "  make build       Build the binary"
	@echo "  make install     Install to /usr/local/bin"
	@echo "  make clean       Remove build artifacts"
	@echo "  make test        Run tests"
	@echo "  make run         Build and run"
	@echo "  make app         Create Mac OS .app bundle"
	@echo ""

build:
	@echo "Building $(BINARY_NAME)..."
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/lazylms
	@echo "Build complete: bin/$(BINARY_NAME)"

clean:
	@echo "Cleaning..."
	@rm -rf bin/ build/ *.app
	@echo "Clean complete"

install: build
	@echo "Installing to /usr/local/bin..."
	@mkdir -p /usr/local/bin
	@cp bin/$(BINARY_NAME) /usr/local/bin/
	@chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "Installation complete"

test:
	@echo "Running tests..."
	go test -v -race -coverprofile=coverage.out ./...
	@echo "Tests complete"

run: build
	@echo "Running $(BINARY_NAME)..."
	@./bin/$(BINARY_NAME)

# Create Mac OS .app bundle
app: build
	@echo "Creating $(BINARY_NAME).app bundle..."
	@mkdir -p build/$(BINARY_NAME).app/Contents/MacOS
	@mkdir -p build/$(BINARY_NAME).app/Contents/Resources
	@cp bin/$(BINARY_NAME) build/$(BINARY_NAME).app/Contents/MacOS/
	@cp Info.plist build/$(BINARY_NAME).app/Contents/
	@echo "App bundle created: build/$(BINARY_NAME).app"

# Development helpers
dev:
	go run ./cmd/lazylms

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run

# Update dependencies
deps:
	go mod tidy
	go mod verify
