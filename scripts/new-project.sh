#!/bin/bash
# Create a new Go project from the template
# Usage: ./new-project.sh <project-name> <module-path> [output-dir]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Check arguments
if [ $# -lt 2 ]; then
    print_error "Insufficient arguments"
    echo ""
    echo "Usage: $0 <project-name> <module-path> [output-dir]"
    echo ""
    echo "Arguments:"
    echo "  project-name  Name of the binary/project (e.g., myapp)"
    echo "  module-path   Go module path (e.g., github.com/user/myapp)"
    echo "  output-dir    Optional output directory (default: current directory)"
    echo ""
    echo "Example:"
    echo "  $0 myapp github.com/user/myapp"
    echo "  $0 myapp github.com/user/myapp ~/projects"
    exit 1
fi

PROJECT_NAME=$1
MODULE_PATH=$2
OUTPUT_DIR=${3:-$(pwd)}

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TEMPLATE_DIR="$(dirname "$SCRIPT_DIR")"

print_info "Creating new Go project: $PROJECT_NAME"
echo "  Module path: $MODULE_PATH"
echo "  Output directory: $OUTPUT_DIR"
echo ""

# Create project directory
PROJECT_DIR="$OUTPUT_DIR/$PROJECT_NAME"
if [ -d "$PROJECT_DIR" ]; then
    print_error "Directory $PROJECT_DIR already exists"
    exit 1
fi

print_info "Creating project structure..."
mkdir -p "$PROJECT_DIR"
cd "$PROJECT_DIR"

# Create directory structure
mkdir -p cmd/$PROJECT_NAME
mkdir -p pkg
mkdir -p internal
mkdir -p build
mkdir -p .github/workflows

print_success "Created directory structure"

# Copy and customize templates
print_info "Copying template files..."

# Copy .goreleaser.yaml
if [ -f "$TEMPLATE_DIR/templates/.goreleaser.template.yaml" ]; then
    sed -e "s|github.com/pavelanni/aistor-snapshot|$MODULE_PATH|g" \
        -e "s|aissnap|$PROJECT_NAME|g" \
        -e "s|aistor-snapshot|$PROJECT_NAME|g" \
        -e "s|Pavel Anni|$(git config user.name 2>/dev/null || echo 'Your Name')|g" \
        -e "s|pavelanni@gmail.com|$(git config user.email 2>/dev/null || echo 'your@email.com')|g" \
        -e "s|pavelanni|$(echo $MODULE_PATH | cut -d'/' -f2)|g" \
        "$TEMPLATE_DIR/templates/.goreleaser.template.yaml" > .goreleaser.yaml
    print_success "Created .goreleaser.yaml"
fi

# Copy .golangci.yml
if [ -f "$TEMPLATE_DIR/templates/.golangci.yml" ]; then
    cp "$TEMPLATE_DIR/templates/.golangci.yml" .golangci.yml
    print_success "Created .golangci.yml"
fi

# Copy .gitignore
if [ -f "$TEMPLATE_DIR/templates/.gitignore" ]; then
    cp "$TEMPLATE_DIR/templates/.gitignore" .gitignore
    print_success "Created .gitignore"
fi

# Copy GitHub workflows
if [ -d "$TEMPLATE_DIR/templates/.github/workflows" ]; then
    cp -r "$TEMPLATE_DIR/templates/.github/workflows/"* .github/workflows/
    print_success "Created GitHub workflows"
fi

# Create and customize Makefile
if [ -f "$TEMPLATE_DIR/templates/Makefile.template" ]; then
    sed -e "s|{{BINARY_NAME}}|$PROJECT_NAME|g" \
        -e "s|{{MODULE_PATH}}|$MODULE_PATH|g" \
        "$TEMPLATE_DIR/templates/Makefile.template" > Makefile
    print_success "Created Makefile"
fi

# Create and customize README
if [ -f "$TEMPLATE_DIR/templates/README.template.md" ]; then
    # Extract owner/repo from module path
    OWNER_REPO=$(echo $MODULE_PATH | cut -d'/' -f2-)
    REPO_URL="https://$(echo $MODULE_PATH | cut -d'/' -f1)/$OWNER_REPO"

    sed -e "s|{{PROJECT_NAME}}|$PROJECT_NAME|g" \
        -e "s|{{MODULE_PATH}}|$MODULE_PATH|g" \
        -e "s|{{REPO_URL}}|$REPO_URL|g" \
        "$TEMPLATE_DIR/templates/README.template.md" > README.md
    print_success "Created README.md"
fi

# Create main.go from example
print_info "Creating main application file..."
cat > "cmd/$PROJECT_NAME/main.go" << 'EOF'
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version is the semantic version (set via ldflags)
	Version = "dev"
	// Commit is the git commit hash (set via ldflags)
	Commit = "unknown"
	// BuildDate is the build timestamp (set via ldflags)
	BuildDate = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "PROJECT_NAME",
	Short: "Short description of PROJECT_NAME",
	Long:  `Longer description of what PROJECT_NAME does.`,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from PROJECT_NAME!")
		fmt.Println("Use --help to see available commands")
	},
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("PROJECT_NAME version %s\nBuilt: %s\nCommit: %s\n",
		Version, BuildDate, Commit))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
EOF

# Replace PROJECT_NAME placeholder
sed -i.bak "s/PROJECT_NAME/$PROJECT_NAME/g" "cmd/$PROJECT_NAME/main.go"
rm "cmd/$PROJECT_NAME/main.go.bak"
print_success "Created cmd/$PROJECT_NAME/main.go"

# Initialize Go module
print_info "Initializing Go module..."
go mod init "$MODULE_PATH"
print_success "Initialized Go module"

# Add Cobra dependency
print_info "Adding dependencies..."
go get github.com/spf13/cobra@latest
go mod tidy
print_success "Added dependencies"

# Initialize git repository
print_info "Initializing Git repository..."
git init -q
git add .
git commit -q -m "chore: Initial commit from go-project-template"
print_success "Initialized Git repository"

echo ""
print_success "Project $PROJECT_NAME created successfully!"
echo ""
echo "Next steps:"
echo "  1. cd $PROJECT_NAME"
echo "  2. make build"
echo "  3. ./build/$PROJECT_NAME --version"
echo ""
echo "Optional:"
echo "  - Edit README.md with your project details"
echo "  - Choose a LICENSE file"
echo "  - Update .goreleaser.yaml with your GitHub username"
echo "  - Create GitHub repository and push:"
echo "      gh repo create $PROJECT_NAME --public --source=. --remote=origin"
echo "      git push -u origin main"
echo ""
