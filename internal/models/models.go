// Package models defines the data structures for Movie Journal.
package models

import "time"

// Movie represents a movie from TMDB with cached metadata.
type Movie struct {
	Title     string `json:"title"`
	PosterURL string `json:"poster_url"`
	Director  string `json:"director"`
	Genre     string `json:"genre"`
	Overview  string `json:"overview"`
	ID        int64  `json:"id"`
	TMDBID    int    `json:"tmdb_id"`
	Year      int    `json:"year"`
}

// DiaryEntry represents a movie viewing session.
type DiaryEntry struct {
	WatchedDate     time.Time `json:"watched_date"`
	CreatedAt       time.Time `json:"created_at"`
	Movie           *Movie    `json:"movie,omitempty"`
	WatchedLocation string    `json:"watched_location,omitempty"`
	WatchedWith     string    `json:"watched_with"`
	Notes           string    `json:"notes"`
	Lookups         []Lookup  `json:"lookups,omitempty"`
	ID              int64     `json:"id"`
	MovieID         int64     `json:"movie_id"`
	Rating          int       `json:"rating"`
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
	CreatedAt    time.Time      `json:"created_at"`
	Question     string         `json:"question"`
	Answer       string         `json:"answer"`
	Category     LookupCategory `json:"category"`
	URL          string         `json:"url,omitempty"`
	ID           int64          `json:"id"`
	DiaryEntryID int64          `json:"diary_entry_id"`
}

// DiaryEntryInput is used for creating/updating diary entries.
type DiaryEntryInput struct {
	WatchedAt   time.Time `json:"watched_at"`
	Location    string    `json:"location,omitempty"`
	Notes       string    `json:"notes"`
	WatchedWith string    `json:"watched_with"`
	MovieID     int64     `json:"movie_id"`
	Rating      int       `json:"rating"`
}

// LookupInput is used for creating/updating lookups.
type LookupInput struct {
	Question     string         `json:"question"`
	Answer       string         `json:"answer"`
	Category     LookupCategory `json:"category"`
	URL          string         `json:"url"`
	DiaryEntryID int64          `json:"diary_entry_id"`
}
