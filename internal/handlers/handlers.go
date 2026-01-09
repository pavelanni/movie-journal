// Package handlers contains HTTP handlers for Movie Journal.
package handlers

import (
	"log/slog"
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
	h.renderDiaryEntry(w, r, func(entry models.DiaryEntry, w http.ResponseWriter, r *http.Request) error {
		return templates.MovieDetails(entry).Render(r.Context(), w)
	})
}

// GetDiaryEntryShort returns a single diary entry's as MovieCard (HTML fragment for HTMX).
func (h *Handlers) GetDiaryEntryShort(w http.ResponseWriter, r *http.Request) {
	h.renderDiaryEntry(w, r, func(entry models.DiaryEntry, w http.ResponseWriter, r *http.Request) error {
		return templates.MovieCard(entry).Render(r.Context(), w)
	})
}

// renderDiaryEntry is a helper that extracts ID, finds entry, and renders using provided function.
func (h *Handlers) renderDiaryEntry(
	w http.ResponseWriter,
	r *http.Request,
	renderFunc func(models.DiaryEntry, http.ResponseWriter, *http.Request) error,
) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	entries := getSampleEntries()
	found := getEntryByID(id, entries)
	if found == nil {
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	if err := renderFunc(*found, w, r); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}

// GetRecentEntries returns filtered diary entries (HTML fragment for HTMX).
func (h *Handlers) GetRecentEntries(w http.ResponseWriter, r *http.Request) {
	entries := getSampleEntries()

	rating := r.URL.Query().Get("min_rating")
	if rating != "" {
		minRating, err := strconv.Atoi(rating)
		if err == nil {
			filtered := make([]models.DiaryEntry, 0)
			for i := range entries {
				if entries[i].Rating >= minRating {
					filtered = append(filtered, entries[i])
				}
			}
			entries = filtered
		}
	}

	err := templates.RecentEntries(entries, rating).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// NewDiaryEntryForm renders the form to create a new diary entry.
func (h *Handlers) NewDiaryEntryForm(w http.ResponseWriter, r *http.Request) {
	err := templates.DiaryNew().Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// CreateDiaryEntry handles the submission of a new diary entry.
func (h *Handlers) CreateDiaryEntry(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Extract form values (in a real app, validate and save to DB)
	watchedDate := r.FormValue("watched_date")
	movieTitle := r.FormValue("movie_title")
	watchedLocation := r.FormValue("watched_location")
	ratingStr := r.FormValue("rating")
	notes := r.FormValue("notes")
	watchedWith := r.FormValue("watched_with")

	// For now, just log the received data (replace with DB save logic)
	slog.Info("Received new diary entry",
		slog.String("watched_date", watchedDate),
		slog.String("watched_location", watchedLocation),
		slog.String("movie_title", movieTitle),
		slog.String("rating", ratingStr),
		slog.String("notes", notes),
		slog.String("watched_with", watchedWith),
	)

	// After logging, redirect back to home page (in a real app, redirect to the new entry's page)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// EditDiaryEntryForm renders the form to edit an existing diary entry.
func (h *Handlers) EditDiaryEntryForm(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Find entry in sample data (will be replaced with DB query later)
	entries := getSampleEntries()
	found := getEntryByID(id, entries)

	if found == nil {
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	// Render the edit form
	err = templates.DiaryEditForm(found).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// EditDiaryEntry handles the editing of an existing diary entry.
func (h *Handlers) EditDiaryEntry(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	idStr := r.PathValue("id")
	// TODO: id will be used when we start working with DB
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Extract form values (in a real app, validate and update in DB)
	watchedDate := r.FormValue("watched_date")
	movieTitle := r.FormValue("movie_title")
	watchedLocation := r.FormValue("watched_location")
	ratingStr := r.FormValue("rating")
	notes := r.FormValue("notes")
	watchedWith := r.FormValue("watched_with")

	// For now, just log the received data (replace with DB update logic)
	slog.Info("Received edit for diary entry",
		slog.String("watched_date", watchedDate),
		slog.String("watched_location", watchedLocation),
		slog.String("movie_title", movieTitle),
		slog.String("rating", ratingStr),
		slog.String("notes", notes),
		slog.String("watched_with", watchedWith),
	)

	// After logging, return to the Movie Details view (in a real app, fetch updated entry from DB)
	entries := getSampleEntries()
	entry := getEntryByID(id, entries)
	if entry == nil {
		http.Error(w, "Entry not found after edit", http.StatusNotFound)
		return
	}
	err = templates.MovieDetails(*entry).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

// DeleteDiaryEntry deletes a diary entry (for HTMX).
func (h *Handlers) DeleteDiaryEntry(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	idStr := r.PathValue("id")
	// TODO: id will be used when we start working with DB
	_, err := strconv.ParseInt(idStr, 10, 64)
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
				Overview: "A depressed man suffering from insomnia meets a strange soap salesman " +
					"named Tyler Durden and soon finds himself living in his squalid house " +
					"after his perfect apartment is destroyed.",
			},
			WatchedDate:     time.Now().AddDate(0, 0, -2),
			WatchedLocation: "Home",
			Rating:          5,
			Notes:           "First rule of Fight Club...",
			WatchedWith:     "Sarah",
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
				Overview: "Cobb, a skilled thief who commits corporate espionage by infiltrating " +
					"the subconscious of his targets is offered a chance to regain his old life " +
					"as payment for a task considered to be impossible: inception.",
			},
			WatchedDate:     time.Now().AddDate(0, 0, -5),
			WatchedLocation: "Cinema",
			Rating:          3,
			Notes:           "The ending still gets me every time. Is it real or not?",
			WatchedWith:     "",
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
				Overview: "A burger-loving hit man, his philosophical partner, a drug-addled " +
					"gangster's moll and a washed-up boxer converge in this sprawling, " +
					"comedic crime caper.",
			},
			WatchedDate:     time.Now().AddDate(0, 0, -10),
			WatchedLocation: "In-flight",
			Rating:          4,
			Notes:           "A masterpiece of non-linear storytelling.",
			WatchedWith:     "Mike",
			Lookups:         []models.Lookup{},
		},
	}
}

func getEntryByID(id int64, entries []models.DiaryEntry) *models.DiaryEntry {
	for i := range entries {
		if entries[i].ID == id {
			return &entries[i]
		}
	}
	return nil
}
