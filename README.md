# Movie Journal

A personal movie diary application for tracking films watched with notes and research moments.

## Motivation

Every night my wife and I watch movies together. We often pause to look things up - who's that actor, where was this filmed, is this based on a true story? Existing movie trackers (Letterboxd, IMDb lists) only log what you watched, not the curiosity moments that made the viewing memorable.

## Features

- **Movie logging** - Search for movies, auto-populate details from TMDB, add ratings and notes
- **Research moments** - Log what you looked up during viewing (actors, locations, trivia)
- **History and browsing** - View past entries with filters, search your diary

## Tech stack

- **Go** - Backend language
- **Templ** - Type-safe HTML templating
- **HTMX** - Frontend interactivity without JavaScript frameworks
- **SQLite** - Embedded database (modernc.org/sqlite for pure Go)
- **Tailwind CSS** - Styling with templui components
- **TMDB API** - Movie metadata, posters, search

## Installation

### From source

```bash
git clone https://github.com/pavelanni/movie-journal.git
cd movie-journal
make build
./build/movie-journal serve
```

### Using Go

```bash
go install github.com/pavelanni/movie-journal/cmd/movie-journal@latest
movie-journal serve
```

## Usage

```bash
# Start the web server
movie-journal serve

# Start on a custom port
movie-journal serve --port 3000

# Use a custom database path
movie-journal serve --db /path/to/diary.db

# Show version
movie-journal version
```

## Development

### Prerequisites

- Go 1.23 or later
- [templ](https://templ.guide/) - `go install github.com/a-h/templ/cmd/templ@latest`
- [air](https://github.com/air-verse/air) - `go install github.com/air-verse/air@latest`
- [Tailwind CSS standalone](https://github.com/tailwindlabs/tailwindcss/releases) - Download for your platform
- [golangci-lint](https://golangci-lint.run/) - For linting

### Development workflow

```bash
# Terminal 1: Start Go + templ watcher
make dev

# Terminal 2: Start Tailwind watcher
make dev-css
```

### Build commands

```bash
make build          # Full build (templ + tailwind + go)
make build-go       # Go build only
make test           # Run tests
make lint           # Run linter
make clean          # Remove build artifacts
```

## Project structure

```
movie-journal/
├── cmd/movie-journal/    # CLI entry point
├── internal/
│   ├── database/         # SQLite operations
│   ├── handlers/         # HTTP handlers
│   ├── models/           # Data structures
│   └── server/           # HTTP server
├── templates/            # Templ templates
├── static/               # Static assets (CSS, JS)
├── PROJECT.md            # Detailed project specification
└── docs/                 # Documentation
```

## License

MIT License
