// Package database provides SQLite database operations for Movie Journal.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "modernc.org/sqlite" // SQLite driver
)

// DB wraps the SQL database connection with Movie Journal operations.
type DB struct {
	*sql.DB
}

// Open opens a SQLite database at the given path.
// It creates the database file if it doesn't exist and runs migrations.
func Open(path string) (*DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	ctx := context.Background()

	// Enable foreign keys
	if _, err := db.ExecContext(ctx, "PRAGMA foreign_keys = ON"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("enabling foreign keys: %w", err)
	}

	// Enable WAL mode for better concurrency
	if _, err := db.ExecContext(ctx, "PRAGMA journal_mode = WAL"); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("enabling WAL mode: %w", err)
	}

	wrapped := &DB{DB: db}

	// Run migrations
	if err := wrapped.Migrate(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	slog.Info("Database opened successfully", slog.String("path", path))
	return wrapped, nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	return db.DB.Close()
}

// WithTimeout returns a context with the given timeout.
func WithTimeout(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
