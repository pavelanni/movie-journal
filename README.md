# Go Project Template

A comprehensive template for Go CLI applications following modern best practices. Includes everything you need to start a production-ready Go project with proper tooling, CI/CD, and release automation.

## Features

- 🏗️ **Modern Project Structure** - Standard Go project layout (cmd/, pkg/, internal/)
- 📦 **GoReleaser Integration** - Multi-platform builds (Linux/macOS, ARM64/AMD64)
- 🔍 **golangci-lint Configuration** - Comprehensive linting with 20+ linters
- ⚙️ **GitHub Actions Workflows** - CI/CD for testing, linting, and releases
- 📋 **Makefile** - Common targets for build, test, lint, release
- 🏷️ **Version Management** - Git tag-based versioning via ldflags
- 📝 **Best Practices Documentation** - Comprehensive guide for Go development
- 🎨 **Cobra CLI Framework** - Full-featured command-line interface

## Quick start

### Option 1: Use GitHub Template (Recommended)

1. Click "Use this template" button on GitHub
2. Create your new repository
3. Clone and start coding!

### Option 2: Bootstrap Script

```bash
# Clone this repository
git clone https://github.com/pavelanni/go-project-template.git
cd go-project-template

# Create a new project
./scripts/new-project.sh myapp github.com/yourname/myapp

# Start working on your new project
cd myapp
make build
./build/myapp --version
```

### Option 3: Manual Setup

1. Copy the template files to your new project
2. Customize placeholders in:
   - `go.mod` (module path)
   - `Makefile` (binary name, module path)
   - `.goreleaser.yaml` (GitHub username, project name)
   - `README.md` (project details)

## What's Included

### Project Structure

```text
your-project/
├── cmd/
│   └── yourapp/
│       └── main.go              # Application entry point
├── pkg/                          # Public libraries
├── internal/                     # Private application code
├── build/                        # Build output directory
├── .github/
│   └── workflows/
│       ├── ci.yml               # CI workflow (build, test, lint)
│       └── release.yml          # Release workflow (GoReleaser)
├── .goreleaser.yaml             # Multi-platform build configuration
├── .golangci.yml                # Linter configuration
├── .gitignore                   # Go-specific gitignore
├── Makefile                     # Build automation
├── go.mod                       # Go module definition
└── README.md                    # Project documentation
```

### Configuration Files

#### `.goreleaser.yaml`

- Builds for Linux (AMD64, ARM64) and macOS (AMD64, ARM64)
- Creates DEB and RPM packages
- Generates checksums
- Auto-generates changelogs from conventional commits
- Optional Homebrew tap support

#### `.golangci.yml`

Enabled linters:

- **Core**: errcheck, govet, ineffassign, staticcheck, unused
- **Code Quality**: revive, gocritic, gocyclo, dupl
- **Security**: gosec
- **Style**: misspell, whitespace, goconst, godot
- **Errors**: errorlint, errname
- **Formatters**: gofmt, goimports

#### `Makefile`

Standard targets:

```bash
make build          # Build the binary
make test           # Run tests
make test-coverage  # Generate coverage report
make lint           # Run linter
make clean          # Remove build artifacts
make install        # Install to $GOPATH/bin
make release        # Create release (requires Git tag)
make snapshot       # Test release locally
make help           # Show all targets
```

### GitHub Actions

#### CI Workflow (`.github/workflows/ci.yml`)

Runs on every push and PR:

- Build the project
- Run tests
- Run linter (golangci-lint)

#### Release Workflow (`.github/workflows/release.yml`)

Triggers on Git tags (`v*`):

- Builds for all platforms
- Creates GitHub release
- Uploads binaries and packages
- Generates changelog

## Documentation

### Best Practices Guide

See [`docs/BEST_PRACTICES.md`](docs/BEST_PRACTICES.md) for a comprehensive guide covering:

- Project structure standards
- Version management
- Testing patterns
- Linting configuration
- CI/CD setup
- Dependency management
- Release process
- New project checklist

## Usage Examples

### Creating a New Project

```bash
# Create a CLI tool
./scripts/new-project.sh mycli github.com/yourname/mycli

# Create in specific directory
./scripts/new-project.sh mycli github.com/yourname/mycli ~/projects
```

### Building and Testing

```bash
cd your-project

# Build
make build

# Test
make test

# Lint
make lint

# Test release locally (no Git tag required)
make snapshot
```

### Creating a Release

```bash
# Commit your changes
git add .
git commit -m "feat: Add awesome feature"

# Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# GitHub Actions will automatically:
# - Build binaries for all platforms
# - Create DEB and RPM packages
# - Generate changelog
# - Create GitHub release
```

## Integration with Claude Code

### Global Context (Recommended)

Copy the best practices guide to Claude's configuration directory:

```bash
mkdir -p ~/.claude
cp docs/BEST_PRACTICES.md ~/.claude/GO_BEST_PRACTICES.md
```

Now Claude Code will automatically reference it in all your Go projects!

### Per-Project Context

```bash
# In your project
mkdir -p .claude
curl -o .claude/go-best-practices.md \
  https://raw.githubusercontent.com/pavelanni/go-project-template/main/docs/BEST_PRACTICES.md
```

### Usage with AI

Tell Claude:

```text
"Create a new Go CLI tool called 'myapp' following the best practices
from my go-project-template repository."
```

Claude will:

1. Create proper project structure
2. Set up version management
3. Configure GoReleaser and linting
4. Add GitHub Actions workflows
5. Initialize with working example code

## Prerequisites

### Required

- Go 1.23 or later
- Git
- Make

### Optional (for full feature set)

- [golangci-lint](https://golangci-lint.run/welcome/install/) - For linting
- [goreleaser](https://goreleaser.com/install/) - For releases
- [GitHub CLI](https://cli.github.com/) - For creating repositories

### Installation

```bash
# macOS
brew install golangci-lint goreleaser gh

# Linux
# See respective project documentation for installation
```

## Customization

### Project-Specific Changes

After creating a new project, customize:

1. **README.md** - Add your project description
2. **LICENSE** - Choose and add a license file
3. **.goreleaser.yaml** - Update GitHub username and repository name
4. **cmd/yourapp/main.go** - Implement your CLI logic

### Template Changes

To improve this template:

1. Fork this repository
2. Make your changes
3. Update [`docs/BEST_PRACTICES.md`](docs/BEST_PRACTICES.md)
4. Submit a pull request

## Examples

This template is used by:

- [aistor-snapshot](https://github.com/pavelanni/aistor-snapshot) - MinIO cluster configuration management tool

## Contributing

Contributions are welcome! Areas for improvement:

- Additional linters or quality checks
- More comprehensive examples
- Integration with other tools (Docker, Kubernetes, etc.)
- Additional GitHub Actions workflows
- Language-specific templates (Python, Rust, etc.)

## License

MIT License - Feel free to use this template for any purpose.

## Resources

- [Effective Go](https://go.dev/doc/effective_go)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [GoReleaser Documentation](https://goreleaser.com/)
- [golangci-lint Documentation](https://golangci-lint.run/)
- [Cobra CLI Framework](https://github.com/spf13/cobra)

## Support

- 📖 [Best Practices Guide](docs/BEST_PRACTICES.md)
- 🐛 [Report Issues](https://github.com/pavelanni/go-project-template/issues)
- 💬 [Discussions](https://github.com/pavelanni/go-project-template/discussions)
