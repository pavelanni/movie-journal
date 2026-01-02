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
	Use:   "example",
	Short: "Example CLI application",
	Long: `Example CLI application demonstrating Go project best practices.

This includes:
- Version management via ldflags
- Cobra for CLI framework
- Proper project structure
- GoReleaser configuration
- GitHub Actions CI/CD`,
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from example CLI!")
		fmt.Println("Use --help to see available commands")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("example version %s\n", Version)
		fmt.Printf("Built: %s\n", BuildDate)
		fmt.Printf("Commit: %s\n", Commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.SetVersionTemplate(fmt.Sprintf("example version %s\nBuilt: %s\nCommit: %s\n",
		Version, BuildDate, Commit))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
