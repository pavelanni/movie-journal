// Package models defines the data structures for Movie Journal.
package models

import "time"

// Movie represents a movie from TMDB with cached metadata.
type Movie struct {
	ID        int64  `json:"id"`
	TMDBID    int    `json:"tmdb_id"`
	Title     string `json:"title"`
	Year      int    `json:"year"`
	PosterURL string `json:"poster_url"`
	Director  string `json:"director"`
	Genre     string `json:"genre"`
	Overview  string `json:"overview"`
}

// DiaryEntry represents a movie viewing session.
type DiaryEntry struct {
	ID          int64     `json:"id"`
	MovieID     int64     `json:"movie_id"`
	Movie       *Movie    `json:"movie,omitempty"`
	WatchedAt   time.Time `json:"watched_at"`
	Location    string    `json:"location,omitempty"`
	Rating      int       `json:"rating"` // 1-5 stars
	Notes       string    `json:"notes"`
	WatchedWith string    `json:"watched_with"`
	CreatedAt   time.Time `json:"created_at"`
	Lookups     []Lookup  `json:"lookups,omitempty"`
}

// LookupCategory represents the type of research moment.
type LookupCategory string

// Lookup categories for research moments.
const (
	LookupCategoryActor    LookupCategory = "actor"
	LookupCategoryLocation LookupCategory = "location"
	LookupCategoryTrivia   LookupCategory = "trivia"
	LookupCategoryOther    LookupCategory = "other"
)

// Lookup represents a research moment during viewing.
type Lookup struct {
	ID           int64          `json:"id"`
	DiaryEntryID int64          `json:"diary_entry_id"`
	Question     string         `json:"question"`
	Answer       string         `json:"answer"`
	Category     LookupCategory `json:"category"`
	URL          string         `json:"url,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
}

// DiaryEntryInput is used for creating/updating diary entries.
type DiaryEntryInput struct {
	MovieID     int64     `json:"movie_id"`
	WatchedAt   time.Time `json:"watched_at"`
	Location    string    `json:"location,omitempty"`
	Rating      int       `json:"rating"`
	Notes       string    `json:"notes"`
	WatchedWith string    `json:"watched_with"`
}

// LookupInput is used for creating/updating lookups.
type LookupInput struct {
	DiaryEntryID int64          `json:"diary_entry_id"`
	Question     string         `json:"question"`
	Answer       string         `json:"answer"`
	Category     LookupCategory `json:"category"`
	URL          string         `json:"url"`
}
