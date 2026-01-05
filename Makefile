.PHONY: build clean install test lint release snapshot help dev dev-css templ-generate templ-watch css-build css-watch

# Binary name
BINARY=movie-journal

# Build directory
BUILD_DIR=build

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Version information
# Try to get version from git tag, fallback to "dev"
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
# Get git commit hash
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
# Get build timestamp in UTC
BUILD_DATE?=$(shell date -u '+%Y-%m-%d_%H:%M:%S_UTC')

# Build flags
LDFLAGS=-ldflags "\
	-X 'main.Version=$(VERSION)' \
	-X 'main.Commit=$(COMMIT)' \
	-X 'main.BuildDate=$(BUILD_DATE)'"

all: build

build: templ-generate css-build
	@echo "Building $(BINARY)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) ./cmd/movie-journal

build-go:
	@echo "Building $(BINARY) (Go only)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) ./cmd/movie-journal

install: build
	@echo "Installing $(BINARY)..."
	@cp $(BUILD_DIR)/$(BINARY) $(GOPATH)/bin/$(BINARY)

clean:
	@echo "Cleaning..."
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -rf tmp/
	@rm -f static/css/tailwind.css

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed" && exit 1)
	golangci-lint run

deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Templ targets
templ-generate:
	@echo "Generating templ templates..."
	@which templ > /dev/null || (echo "templ not installed. Run: go install github.com/a-h/templ/cmd/templ@latest" && exit 1)
	templ generate

templ-watch:
	@echo "Watching templ templates..."
	@which templ > /dev/null || (echo "templ not installed. Run: go install github.com/a-h/templ/cmd/templ@latest" && exit 1)
	templ generate --watch

# CSS targets
css-build:
	@echo "Building CSS..."
	@test -f ./tailwindcss || (echo "tailwindcss binary not found. Download from https://github.com/tailwindlabs/tailwindcss/releases" && exit 1)
	@mkdir -p static/css
	./tailwindcss -i input.css -o static/css/tailwind.css --minify

css-watch:
	@echo "Watching CSS..."
	@test -f ./tailwindcss || (echo "tailwindcss binary not found. Download from https://github.com/tailwindlabs/tailwindcss/releases" && exit 1)
	@mkdir -p static/css
	./tailwindcss -i input.css -o static/css/tailwind.css --watch

# Development targets
dev:
	@echo "Starting development server with air..."
	@which air > /dev/null || (echo "air not installed. Run: go install github.com/air-verse/air@latest" && exit 1)
	air

dev-css:
	@echo "Starting CSS watcher..."
	$(MAKE) css-watch

# Run both dev and dev-css (in separate terminals)
dev-all:
	@echo "Run these in separate terminals:"
	@echo "  make dev      # Go + templ watcher"
	@echo "  make dev-css  # Tailwind watcher"

run: build
	./$(BUILD_DIR)/$(BINARY) serve

release:
	@echo "Creating release..."
	@which goreleaser > /dev/null || (echo "goreleaser not installed. Install from https://goreleaser.com/install/" && exit 1)
	goreleaser release --clean

snapshot:
	@echo "Creating snapshot release (no Git tag required)..."
	@which goreleaser > /dev/null || (echo "goreleaser not installed. Install from https://goreleaser.com/install/" && exit 1)
	goreleaser release --snapshot --clean

release-check:
	@echo "Checking release configuration..."
	@which goreleaser > /dev/null || (echo "goreleaser not installed. Install from https://goreleaser.com/install/" && exit 1)
	goreleaser check

help:
	@echo "Available targets:"
	@echo "  build          - Build the binary (includes templ + tailwind)"
	@echo "  build-go       - Build Go binary only (no templ/tailwind)"
	@echo "  install        - Build and install to GOPATH/bin"
	@echo "  clean          - Remove build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  lint           - Run linter"
	@echo "  deps           - Download and tidy dependencies"
	@echo "  templ-generate - Generate templ templates"
	@echo "  templ-watch    - Watch and regenerate templ templates"
	@echo "  css-build      - Build Tailwind CSS"
	@echo "  css-watch      - Watch and rebuild Tailwind CSS"
	@echo "  dev            - Start development server with air"
	@echo "  dev-css        - Start Tailwind CSS watcher"
	@echo "  dev-all        - Instructions for full dev setup"
	@echo "  run            - Build and run the server"
	@echo "  release        - Create a release with GoReleaser (requires Git tag)"
	@echo "  snapshot       - Create a snapshot release (no Git tag required)"
	@echo "  release-check  - Check GoReleaser configuration"
