# Go project best practices

A checklist and guidelines for starting new Go projects based on industry standards and community best practices.

## Core infrastructure

### Version management with Git tags

Embed version information using ldflags at build time.

**Implementation**:

```go
// cmd/version.go
package cmd

var (
    Version   = "dev"
    Commit    = "unknown"
    BuildDate = "unknown"
)

func VersionInfo() string {
    return fmt.Sprintf("app version %s\nBuilt: %s\nCommit: %s",
        Version, BuildDate, Commit)
}
```

**Makefile**:

```makefile
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE?=$(shell date -u '+%Y-%m-%d_%H:%M:%S_UTC')

LDFLAGS=-ldflags "\
    -X 'github.com/user/project/cmd.Version=$(VERSION)' \
    -X 'github.com/user/project/cmd.Commit=$(COMMIT)' \
    -X 'github.com/user/project/cmd.BuildDate=$(BUILD_DATE)'"
```

**Usage**:

```bash
# Tag a release
git tag -a v1.0.0 -m "Release v1.0.0"

# Build with version info
make build
./build/app --version
```

### GoReleaser for multi-platform builds

Use GoReleaser for automated, reproducible releases across multiple platforms.

**Installation**:

```bash
# macOS
brew install goreleaser

# Linux
echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' | sudo tee /etc/apt/sources.list.d/goreleaser.list
sudo apt update
sudo apt install goreleaser
```

**Minimal .goreleaser.yaml**:

```yaml
version: 2

builds:
  - binary: myapp
    main: ./cmd/myapp
    env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
    goarch:
      - amd64
      - arm64
    ldflags:
      - -s -w
      - -X github.com/user/project/cmd.Version={{.Version}}
      - -X github.com/user/project/cmd.Commit={{.Commit}}
      - -X github.com/user/project/cmd.BuildDate={{.Date}}

archives:
  - format: tar.gz
    files:
      - README.md
      - LICENSE

nfpms:
  - formats:
      - deb
      - rpm
    bindir: /usr/bin
    contents:
      - src: README.md
        dst: /usr/share/doc/myapp/README.md

checksum:
  name_template: 'checksums.txt'

changelog:
  sort: asc
```

**Makefile targets**:

```makefile
release:
    goreleaser release --clean

snapshot:
    goreleaser release --snapshot --clean

release-check:
    goreleaser check
```

**Usage**:

```bash
# Check configuration
make release-check

# Test locally without Git tag
make snapshot

# Create release (requires Git tag)
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
make release
```

### Standard project structure

```text
myproject/
├── cmd/                    # Command-line applications
│   └── myapp/             # Main application
│       └── main.go
├── pkg/                    # Public libraries (importable by others)
│   └── mylib/
│       ├── types.go
│       └── lib.go
├── internal/              # Private application code (not importable)
│   └── service/
│       └── handler.go
├── build/                 # Build output directory
├── .github/               # GitHub-specific files
│   └── workflows/
│       └── release.yml
├── .goreleaser.yaml      # GoReleaser configuration
├── .gitignore            # Git ignore patterns
├── go.mod                # Go module definition
├── go.sum                # Go dependencies checksums
├── Makefile              # Build automation
├── README.md             # Project documentation
├── LICENSE               # License file
└── CHANGELOG.md          # Release notes (optional)
```

### Makefile with standard targets

Essential targets for every Go project:

```makefile
.PHONY: build clean install test lint deps run help

BINARY=myapp
BUILD_DIR=build

build:
    @mkdir -p $(BUILD_DIR)
    go build -o $(BUILD_DIR)/$(BINARY) ./cmd/myapp

install: build
    @cp $(BUILD_DIR)/$(BINARY) $(GOPATH)/bin/$(BINARY)

clean:
    @rm -rf $(BUILD_DIR)
    @go clean

test:
    go test -v ./...

test-coverage:
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html

lint:
    golangci-lint run

deps:
    go mod download
    go mod tidy

run: build
    ./$(BUILD_DIR)/$(BINARY)

help:
    @echo "Available targets:"
    @echo "  build    - Build the binary"
    @echo "  install  - Install to GOPATH/bin"
    @echo "  clean    - Remove build artifacts"
    @echo "  test     - Run tests"
    @echo "  lint     - Run linter"
    @echo "  deps     - Manage dependencies"
    @echo "  run      - Build and run"
```

## Development tools

### Linting and code quality

**golangci-lint** - Fast, comprehensive linter aggregator.

**Installation**:

```bash
# macOS
brew install golangci-lint

# Linux
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

**Configuration** (.golangci.yml):

```yaml
run:
  timeout: 5m
  tests: true

linters:
  enable:
    - gofmt
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - ineffassign
    - typecheck
    - misspell
    - gocritic
    - gocyclo
    - goimports

linters-settings:
  gocyclo:
    min-complexity: 15
  gocritic:
    enabled-tags:
      - diagnostic
      - performance
      - style
```

**Usage**:

```bash
# Run linter
make lint

# Auto-fix issues
golangci-lint run --fix
```

### Testing setup

**Test file naming**: `*_test.go`

**Test organization**:

```go
package mylib_test  // External package testing

import (
    "testing"
    "github.com/user/project/pkg/mylib"
)

func TestFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:  "valid input",
            input: "test",
            want:  "expected",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := mylib.Function(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("Function() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("Function() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

**Test commands**:

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Verbose output
go test -v ./...

# Run specific test
go test -run TestFunction ./pkg/mylib
```

### Git configuration

**.gitignore**:

```gitignore
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out

# Build directories
build/
dist/

# Go workspace
go.work
go.work.sum

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db

# Coverage
coverage.out
coverage.html

# Config with secrets
*.local
.env
```

## Documentation

### README structure

Essential sections for a good README:

```markdown
# Project Name

Brief description (1-2 sentences).

## Features

- Feature 1
- Feature 2
- Feature 3

## Installation

### From source
bash
git clone https://github.com/user/project.git
cd project
make build


### Using Go
bash
go install github.com/user/project/cmd/myapp@latest


## Quick start

bash
# Basic usage
myapp command


## Commands

- `command1` - Description
- `command2` - Description

## Configuration

Configuration options and examples.

## Development

bash
make build
make test
make lint


## License

Apache 2.0 (or your chosen license)
```

### Code documentation

Follow Go documentation conventions:

```go
// Package mylib provides functionality for X.
//
// This package is designed for Y use cases.
package mylib

// Function does X and returns Y.
// It accepts parameter Z which must be non-empty.
//
// Example:
//
//     result, err := Function("input")
//     if err != nil {
//         log.Fatal(err)
//     }
//
// Returns an error if the input is invalid.
func Function(input string) (string, error) {
    // Implementation
}
```

## CI/CD

### GitHub Actions for releases

**.github/workflows/release.yml**:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

permissions:
  contents: write

jobs:
  goreleaser:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v6
        with:
          distribution: goreleaser
          version: latest
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### GitHub Actions for CI

**.github/workflows/ci.yml**:

```yaml
name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - name: Run tests
        run: go test -v -race -coverprofile=coverage.out ./...
      - name: Upload coverage
        uses: codecov/codecov-action@v4
        with:
          file: ./coverage.out

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest
```

## Dependency management

### Go modules best practices

```bash
# Initialize module
go mod init github.com/user/project

# Add dependencies (automatically updates go.mod)
go get github.com/spf13/cobra@latest

# Update dependencies
go get -u ./...

# Tidy up (remove unused, add missing)
go mod tidy

# Verify dependencies
go mod verify

# Vendor dependencies (optional)
go mod vendor
```

**go.mod example**:

```go
module github.com/user/project

go 1.23

require (
    github.com/spf13/cobra v1.8.0
    github.com/stretchr/testify v1.9.0
)
```

### Prefer standard library over third-party dependencies

**Philosophy**: Minimize external dependencies by preferring Go's standard library whenever possible.

**Benefits**:

- **Stability**: Standard library is well-maintained and versioned with Go
- **Security**: Fewer dependencies = smaller attack surface
- **Simplicity**: No version conflicts or dependency hell
- **Performance**: Standard library is highly optimized
- **Maintenance**: Less code to update and fewer breaking changes

**Guidelines**:

1. **Always check standard library first** before adding external packages
1. **Use latest versions** of packages unless there's a specific reason not to
1. **Document why** if you must use an older version (e.g., compatibility, known bugs)

**Common examples**:

| Use case              | ❌ Avoid                                           | ✅ Prefer                      |
| --------------------- | ------------------------------------------------- | ----------------------------- |
| Structured logging    | `github.com/sirupsen/logrus`<br>`go.uber.org/zap` | `log/slog` (stdlib, Go 1.21+) |
| HTTP server           | External frameworks                               | `net/http` (stdlib)           |
| HTTP client           | `github.com/go-resty/resty`                       | `net/http` (stdlib)           |
| JSON parsing          | `github.com/json-iterator/go`                     | `encoding/json` (stdlib)      |
| Command line flags    | (basic needs)                                     | `flag` (stdlib)               |
| Environment variables | (basic needs)                                     | `os` (stdlib)                 |
| File I/O              | External wrappers                                 | `os`, `io`, `bufio` (stdlib)  |
| Testing               | (basic assertions)                                | `testing` (stdlib)            |
| Context management    | Third-party packages                              | `context` (stdlib)            |

**Example: Structured logging**

```go
// ❌ Avoid adding dependencies for logging
import "go.uber.org/zap"
import "github.com/sirupsen/logrus"

// ✅ Use standard library slog (Go 1.21+)
import "log/slog"

func main() {
    // JSON handler for production
    logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))

    // Text handler for development
    logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))

    // Structured logging with key-value pairs
    logger.Info("Server started",
        slog.String("addr", ":8080"),
        slog.Int("pid", os.Getpid()),
    )
}
```

**When external dependencies are acceptable**:

Use third-party packages when they provide significant value that stdlib doesn't:

- **CLI frameworks**: `github.com/spf13/cobra` for complex CLIs with subcommands
- **Configuration**: `github.com/spf13/viper` for advanced config management
- **Validation**: `github.com/go-playground/validator` for struct validation
- **Testing utilities**: `github.com/stretchr/testify` for advanced assertions
- **Protocol implementations**: gRPC, GraphQL, message queues
- **Specialized algorithms**: Crypto libraries, compression, etc.

**Version management**:

```bash
# ✅ Use @latest for new dependencies
go get github.com/spf13/cobra@latest

# ❌ Avoid pinning to old versions without reason
go get github.com/spf13/cobra@v1.0.0

# ✅ Document if you must use a specific version
# In go.mod or documentation:
# Using v1.0.0 because v1.1.0 has breaking changes in our use case
# See: https://github.com/project/issues/123
```

**Updating dependencies**:

```bash
# Check for updates
go list -m -u all

# Update all dependencies to latest
go get -u ./...

# Update specific dependency
go get -u github.com/spf13/cobra@latest

# Test after updates
go test ./...
make lint
```

## Modern Go practices

Use modern Go language features and idioms introduced in recent versions. The Go compiler often suggests these improvements.

### Use `any` instead of `interface{}`

**Since**: Go 1.18

The `any` type alias is clearer and more concise than `interface{}`.

```go
// ❌ Old style
func Process(data interface{}) error {
    return nil
}

var result interface{}

// ✅ Modern style
func Process(data any) error {
    return nil
}

var result any
```

### Use generics instead of `interface{}`

**Since**: Go 1.18

When you need type safety, use generics instead of `any`/`interface{}` with type assertions.

```go
// ❌ Old style - runtime type assertions
func Max(a, b interface{}) interface{} {
    if a.(int) > b.(int) {
        return a
    }
    return b
}

// ✅ Modern style - compile-time type safety
func Max[T cmp.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}
```

### Use built-in `min`, `max`, `clear`

**Since**: Go 1.21

```go
// ❌ Old style
func getMax(a, b int) int {
    if a > b {
        return a
    }
    return b
}

// Clear map
for k := range myMap {
    delete(myMap, k)
}

// ✅ Modern style
result := max(a, b)
minVal := min(x, y, z)  // Works with multiple values

clear(myMap)  // Clear map
clear(slice)  // Zero out slice elements
```

### Use `slices` package

**Since**: Go 1.21

The `slices` package provides generic slice operations.

```go
import "slices"

// ❌ Old style - manual operations
func contains(slice []string, item string) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}

func sortStrings(s []string) {
    sort.Strings(s)
}

// ✅ Modern style - use slices package
found := slices.Contains(items, "search")
slices.Sort(items)
slices.Reverse(items)
index := slices.Index(items, "search")
filtered := slices.DeleteFunc(items, func(s string) bool {
    return s == ""
})
```

### Use `maps` package

**Since**: Go 1.21

```go
import "maps"

// ❌ Old style
func copyMap(m map[string]int) map[string]int {
    result := make(map[string]int, len(m))
    for k, v := range m {
        result[k] = v
    }
    return result
}

// ✅ Modern style
copied := maps.Clone(original)
maps.Copy(dst, src)
equal := maps.Equal(map1, map2)
keys := slices.Collect(maps.Keys(myMap))  // Go 1.23+
```

### Use `strings.Cut` for string splitting

**Since**: Go 1.18

`strings.Cut` is clearer and more efficient than `Index` + slicing.

```go
// ❌ Old style
idx := strings.Index(s, "=")
if idx >= 0 {
    key := s[:idx]
    value := s[idx+1:]
}

// ✅ Modern style
key, value, found := strings.Cut(s, "=")
if found {
    // Use key and value
}
```

### Use `errors.Join` for multiple errors

**Since**: Go 1.20

```go
import "errors"

// ❌ Old style - lose error information
var err error
if err1 := step1(); err1 != nil {
    err = err1
}
if err2 := step2(); err2 != nil {
    err = err2  // Lost err1!
}

// ✅ Modern style - preserve all errors
var errs []error
if err := step1(); err != nil {
    errs = append(errs, err)
}
if err := step2(); err != nil {
    errs = append(errs, err)
}
return errors.Join(errs...)
```

### Use `errors.Is` and `errors.As`

**Since**: Go 1.13

```go
// ❌ Old style
if err == io.EOF {
    // Handle
}

// ✅ Modern style - works with wrapped errors
if errors.Is(err, io.EOF) {
    // Handle
}

// Type assertions
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Println(pathErr.Path)
}
```

### Use structured logging with `log/slog`

**Since**: Go 1.21

See "Prefer standard library over third-party dependencies" section for details.

```go
import "log/slog"

// ✅ Modern style
slog.Info("Request processed",
    slog.String("method", "GET"),
    slog.Int("status", 200),
    slog.Duration("duration", elapsed),
)
```

### Use `context.WithoutCancel`

**Since**: Go 1.21

```go
// ✅ Detach context from parent cancellation
func background(ctx context.Context) {
    // Continue work even if parent is cancelled
    detached := context.WithoutCancel(ctx)
    go cleanup(detached)
}
```

### Minimum Go version recommendation

Set minimum Go version in `go.mod` to enable modern features:

```go
module github.com/user/project

go 1.23  // Or latest stable version

require (
    // dependencies
)
```

**Benefits of staying current**:

- Access to latest language features and standard library improvements
- Better performance and compiler optimizations
- Security fixes and bug fixes
- Cleaner, more maintainable code

## Additional best practices

### Semantic versioning

Follow SemVer (MAJOR.MINOR.PATCH):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes

Example tags:

- `v1.0.0` - First stable release
- `v1.1.0` - Added new feature
- `v1.1.1` - Bug fix
- `v2.0.0` - Breaking changes

### Changelog management

Use conventional commits for automatic changelog generation:

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `chore:` - Maintenance tasks
- `refactor:` - Code refactoring
- `test:` - Test additions/changes
- `perf:` - Performance improvements

GoReleaser can auto-generate changelogs from these commit messages.

### License

Choose and include a LICENSE file:

- **MIT**: Permissive, simple
- **Apache 2.0**: Permissive, includes patent grant
- **GPL**: Copyleft, requires derivative works to be open source

Add license header to files:

```go
// Copyright 2025 Your Name
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
```

## New project checklist

When starting a new Go project:

- [ ] Initialize Git repository and push to GitHub
- [ ] Create go.mod with `go mod init`
- [ ] Set up standard directory structure (cmd/, pkg/, internal/)
- [ ] Create Makefile with essential targets
- [ ] Add version management code and Makefile flags
- [ ] Create .goreleaser.yaml configuration
- [ ] Add .gitignore file
- [ ] Create README.md with project information
- [ ] Add LICENSE file
- [ ] Set up golangci-lint with .golangci.yml
- [ ] Create GitHub Actions workflows (ci.yml, release.yml)
- [ ] Write initial tests
- [ ] Tag first release (v0.1.0 or v1.0.0)
- [ ] Test release process with `make snapshot`

## Useful tools and libraries

**Note**: Always prefer standard library first (see section 13). Only use these when stdlib doesn't meet your needs.

### CLI applications

- **cobra** - CLI framework with commands and flags (when `flag` pkg is insufficient)
- **viper** - Advanced configuration management (for complex config needs)
- **lipgloss** - Terminal styling and colors
- **bubbletea** - Interactive TUI framework

### Testing

- **testify** - Advanced test assertions and mocking (beyond stdlib `testing`)
- **gomock** - Mock generation
- **httptest** - HTTP testing (standard library)

### Logging

- **log/slog** - ✅ Structured logging (standard library, Go 1.21+) - USE THIS
- ~~logrus~~ - ❌ Avoid, use `log/slog` instead
- ~~zap~~ - ❌ Avoid, use `log/slog` instead

### Utilities

- **errors** - Error wrapping (standard library, Go 1.13+)
- **validator** - Struct validation (when needed beyond manual checks)

## References

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [GoReleaser Documentation](https://goreleaser.com/)
- [golangci-lint](https://golangci-lint.run/)
- [Semantic Versioning](https://semver.org/)
