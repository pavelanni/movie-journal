// Package server provides the HTTP server for Movie Journal.
package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/pavelanni/movie-journal/internal/database"
	"github.com/pavelanni/movie-journal/internal/handlers"
)

// Config holds server configuration.
type Config struct {
	Port int
	DB   *database.DB
}

// Server is the Movie Journal HTTP server.
type Server struct {
	config     Config
	httpServer *http.Server
	mux        *http.ServeMux
	handlers   *handlers.Handlers
}

// New creates a new server with the given configuration.
func New(cfg Config) *Server {
	mux := http.NewServeMux()

	s := &Server{
		config:   cfg,
		mux:      mux,
		handlers: handlers.New(cfg.DB),
		httpServer: &http.Server{
			Addr:         fmt.Sprintf(":%d", cfg.Port),
			Handler:      mux,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
	}

	s.setupRoutes()

	return s
}

// setupRoutes configures all HTTP routes.
func (s *Server) setupRoutes() {
	// Static files
	fs := http.FileServer(http.Dir("static"))
	s.mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// Health check
	s.mux.HandleFunc("GET /health", s.handleHealth)

	// Home page
	s.mux.HandleFunc("GET /", s.handlers.Home)

	// About page
	s.mux.HandleFunc("GET /about", s.handlers.About)

	// HTMX endpoints
	s.mux.HandleFunc("GET /diary/{id}", s.handlers.GetDiaryEntry)
	s.mux.HandleFunc("DELETE /diary/{id}", s.handlers.DeleteDiaryEntry)
	s.mux.HandleFunc("GET /diary-short/{id}", s.handlers.GetDiaryEntryShort)
	s.mux.HandleFunc("GET /recent-entries", s.handlers.GetRecentEntries)
	s.mux.HandleFunc("GET /diary/new", s.handlers.NewDiaryEntryForm)
	s.mux.HandleFunc("POST /diary/new", s.handlers.CreateDiaryEntry)

}

// Start starts the HTTP server.
func (s *Server) Start() error {
	slog.Info("Starting server",
		slog.String("addr", s.httpServer.Addr),
	)
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown(ctx context.Context) error {
	slog.Info("Shutting down server")
	return s.httpServer.Shutdown(ctx)
}

// handleHealth returns server health status.
func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, `{"status":"ok"}`)
}
