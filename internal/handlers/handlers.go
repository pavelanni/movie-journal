// Package handlers contains HTTP handlers for Movie Journal.
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/pavelanni/movie-journal/internal/database"
	"github.com/pavelanni/movie-journal/internal/models"
	"github.com/pavelanni/movie-journal/templates"
)

// Handlers contains all HTTP handlers.
type Handlers struct {
	db *database.DB
}

// New creates a new Handlers instance.
func New(db *database.DB) *Handlers {
	return &Handlers{db: db}
}

// Home renders the home page with recent diary entries.
func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	// For now, use sample data until we implement database queries
	entries := getSampleEntries()

	err := templates.Index(entries, "").Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// About renders the about page.
func (h *Handlers) About(w http.ResponseWriter, r *http.Request) {
	err := templates.About().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// GetDiaryEntry returns a single diary entry's details (HTML fragment for HTMX).
func (h *Handlers) GetDiaryEntry(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Find entry in sample data (will be replaced with DB query later)
	entries := getSampleEntries()
	var found *models.DiaryEntry
	for i := range entries {
		if entries[i].ID == id {
			found = &entries[i]
			break
		}
	}

	if found == nil {
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	// Render just the details fragment (no layout wrapper)
	err = templates.MovieDetails(*found).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// GetDiaryEntryShort returns a single diary entry's as MovieCard (HTML fragment for HTMX).
func (h *Handlers) GetDiaryEntryShort(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Find entry in sample data (will be replaced with DB query later)
	entries := getSampleEntries()
	var found *models.DiaryEntry
	for i := range entries {
		if entries[i].ID == id {
			found = &entries[i]
			break
		}
	}

	if found == nil {
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	// Render just the MovieCard fragment (no layout wrapper)
	err = templates.MovieCard(*found).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) GetRecentEntries(w http.ResponseWriter, r *http.Request) {
	// For now, use sample data until we implement database queries
	entries := getSampleEntries()

	rating := r.URL.Query().Get("min_rating")
	if rating != "" {
		minRating, err := strconv.Atoi(rating)
		if err == nil {
			filtered := make([]models.DiaryEntry, 0)
			for _, entry := range entries {
				if entry.Rating >= minRating {
					filtered = append(filtered, entry)
				}
			}
			entries = filtered
		}
	}
	// Render just the RecentEntries fragment (no layout wrapper)
	err := templates.RecentEntries(entries, rating).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// DeleteDiaryEntry deletes a diary entry (for HTMX).
func (h *Handlers) DeleteDiaryEntry(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	idStr := r.PathValue("id")
	_, err := strconv.ParseInt(idStr, 10, 64) // we don't use id for now, but it should be added when we start working with DB
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Here you would add code to delete the entry from the database.
	// Return 200 OK with empty body - HTMX will replace the target with nothing (remove it).
	// Note: 204 No Content doesn't trigger HTMX swaps by default.
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	// Empty response body - with hx-swap="outerHTML", this removes the element
}

// getSampleEntries returns sample diary entries for development.
func getSampleEntries() []models.DiaryEntry {
	return []models.DiaryEntry{
		{
			ID:      1,
			MovieID: 1,
			Movie: &models.Movie{
				ID:        1,
				TMDBID:    550,
				Title:     "Fight Club",
				Year:      1999,
				PosterURL: "https://image.tmdb.org/t/p/w185/pB8BM7pdSp6B6Ih7QZ4DrQ3PmJK.jpg",
				Director:  "David Fincher",
				Genre:     "Drama",
				Overview:  "A depressed man suffering from insomnia meets a strange soap salesman named Tyler Durden and soon finds himself living in his squalid house after his perfect apartment is destroyed.",
			},
			WatchedAt:   time.Now().AddDate(0, 0, -2),
			Location:    "Home",
			Rating:      5,
			Notes:       "First rule of Fight Club...",
			WatchedWith: "Sarah",
			Lookups: []models.Lookup{
				{
					ID:       1,
					Question: "Where was the Paper Street house?",
					Answer:   "The house was located in Wilmington, Delaware",
					Category: models.LookupCategoryLocation,
				},
			},
		},
		{
			ID:      2,
			MovieID: 2,
			Movie: &models.Movie{
				ID:        2,
				TMDBID:    27205,
				Title:     "Inception",
				Year:      2010,
				PosterURL: "https://image.tmdb.org/t/p/w185/oYuLEt3zVCKq57qu2F8dT7NIa6f.jpg",
				Director:  "Christopher Nolan",
				Genre:     "Sci-Fi",
				Overview:  "Cobb, a skilled thief who commits corporate espionage by infiltrating the subconscious of his targets is offered a chance to regain his old life as payment for a task considered to be impossible: inception.",
			},
			WatchedAt:   time.Now().AddDate(0, 0, -5),
			Location:    "Cinema",
			Rating:      3,
			Notes:       "The ending still gets me every time. Is it real or not?",
			WatchedWith: "",
			Lookups: []models.Lookup{
				{
					ID:       2,
					Question: "Who composed the score?",
					Answer:   "Hans Zimmer",
					Category: models.LookupCategoryTrivia,
				},
				{
					ID:       3,
					Question: "Where was the rotating hallway filmed?",
					Answer:   "Cardington Studios in Bedfordshire, UK",
					Category: models.LookupCategoryLocation,
				},
			},
		},
		{
			ID:      3,
			MovieID: 3,
			Movie: &models.Movie{
				ID:        3,
				TMDBID:    680,
				Title:     "Pulp Fiction",
				Year:      1994,
				PosterURL: "https://image.tmdb.org/t/p/w185/d5iIlFn5s0ImszYzBPb8JPIfbXD.jpg",
				Director:  "Quentin Tarantino",
				Genre:     "Crime",
				Overview:  "A burger-loving hit man, his philosophical partner, a drug-addled gangster's moll and a washed-up boxer converge in this sprawling, comedic crime caper.",
			},
			WatchedAt:   time.Now().AddDate(0, 0, -10),
			Location:    "In-flight",
			Rating:      4,
			Notes:       "A masterpiece of non-linear storytelling.",
			WatchedWith: "Mike",
			Lookups:     []models.Lookup{},
		},
	}
}
