package database

import (
	"context"
	"fmt"
	"log/slog"
)

// schemaVersion is the current database schema version.
const schemaVersion = 1

// Migrate runs database migrations to bring the schema up to date.
func (db *DB) Migrate(ctx context.Context) error {
	// Create migrations table if it doesn't exist
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("creating migrations table: %w", err)
	}

	// Get current version
	var currentVersion int
	err = db.QueryRowContext(ctx, "SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&currentVersion)
	if err != nil {
		return fmt.Errorf("getting current version: %w", err)
	}

	slog.Info("Database migration check",
		slog.Int("current_version", currentVersion),
		slog.Int("target_version", schemaVersion),
	)

	// Run migrations
	for v := currentVersion + 1; v <= schemaVersion; v++ {
		if err := db.runMigration(ctx, v); err != nil {
			return fmt.Errorf("running migration %d: %w", v, err)
		}
		slog.Info("Applied migration", slog.Int("version", v))
	}

	return nil
}

func (db *DB) runMigration(ctx context.Context, version int) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var migration string
	switch version {
	case 1:
		migration = migrationV1
	default:
		return fmt.Errorf("unknown migration version: %d", version)
	}

	if _, err := tx.ExecContext(ctx, migration); err != nil {
		return fmt.Errorf("executing migration: %w", err)
	}

	if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations (version) VALUES (?)", version); err != nil {
		return fmt.Errorf("recording migration: %w", err)
	}

	return tx.Commit()
}

// migrationV1 creates the initial schema.
const migrationV1 = `
-- Movies table: cached movie metadata from TMDB
CREATE TABLE IF NOT EXISTS movies (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	tmdb_id INTEGER UNIQUE NOT NULL,
	title TEXT NOT NULL,
	year INTEGER,
	poster_url TEXT,
	director TEXT,
	genre TEXT,
	overview TEXT
);

CREATE INDEX IF NOT EXISTS idx_movies_tmdb_id ON movies(tmdb_id);
CREATE INDEX IF NOT EXISTS idx_movies_title ON movies(title);

-- Diary entries: individual movie viewing sessions
CREATE TABLE IF NOT EXISTS diary_entries (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	movie_id INTEGER NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
	watched_at DATE NOT NULL,
	rating INTEGER CHECK (rating >= 1 AND rating <= 5),
	notes TEXT,
	watched_with TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_diary_entries_movie_id ON diary_entries(movie_id);
CREATE INDEX IF NOT EXISTS idx_diary_entries_watched_at ON diary_entries(watched_at DESC);

-- Lookups: research moments during viewing
CREATE TABLE IF NOT EXISTS lookups (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	diary_entry_id INTEGER NOT NULL REFERENCES diary_entries(id) ON DELETE CASCADE,
	question TEXT NOT NULL,
	answer TEXT,
	category TEXT CHECK (category IN ('actor', 'location', 'trivia', 'other')) DEFAULT 'other',
	url TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_lookups_diary_entry_id ON lookups(diary_entry_id);
CREATE INDEX IF NOT EXISTS idx_lookups_category ON lookups(category);
`
