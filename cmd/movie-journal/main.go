// Movie Journal - Personal movie diary with research moments
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pavelanni/movie-journal/internal/database"
	"github.com/pavelanni/movie-journal/internal/server"
	"github.com/spf13/cobra"
)

var (
	// Version is the semantic version (set via ldflags).
	Version = "dev"
	// Commit is the git commit hash (set via ldflags).
	Commit = "unknown"
	// BuildDate is the build timestamp (set via ldflags).
	BuildDate = "unknown"
)

var (
	port   int
	dbPath string
)

var rootCmd = &cobra.Command{
	Use:   "movie-journal",
	Short: "Personal movie diary application",
	Long: `Movie Journal is a personal movie diary for tracking films watched
with notes and research moments.

Capture what you watch, rate it, add notes, and log the things
you looked up during viewing - who's that actor, where was this
filmed, is this based on a true story?`,
	Version: Version,
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	Long:  `Start the Movie Journal web server to access the diary through your browser.`,
	RunE:  runServe,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("movie-journal version %s\n", Version)
		fmt.Printf("Built: %s\n", BuildDate)
		fmt.Printf("Commit: %s\n", Commit)
	},
}

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to listen on")
	serveCmd.Flags().StringVarP(&dbPath, "db", "d", "movie-journal.db", "Path to SQLite database file")

	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.SetVersionTemplate(fmt.Sprintf("movie-journal version %s\nBuilt: %s\nCommit: %s\n",
		Version, BuildDate, Commit))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runServe(_ *cobra.Command, _ []string) error {
	// Setup logging
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	slog.Info("Starting Movie Journal",
		slog.String("version", Version),
		slog.Int("port", port),
		slog.String("database", dbPath),
	)

	// Open database
	db, err := database.Open(dbPath)
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer func() { _ = db.Close() }()

	// Create server
	srv := server.New(server.Config{
		Port: port,
		DB:   db,
	})

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		fmt.Printf("\nMovie Journal server running at http://localhost:%d\n", port)
		fmt.Println("Press Ctrl+C to stop")
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		slog.Info("Received shutdown signal")
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	}

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutting down server: %w", err)
	}

	slog.Info("Server stopped gracefully")
	return nil
}
