# Movie Journal

A personal movie diary application for tracking films watched with notes and research moments.

## Motivation

Every night my wife and I watch movies together. We often pause to look things up - who's that
actor, where was this filmed, is this based on a true story? Existing movie trackers (Letterboxd,
IMDb lists) only log what you watched, not the curiosity moments that made the viewing memorable.

## Goals

1. **Learn HTMX and Templ** - This is primarily a learning project
1. **Actually use it** - Build something we'll use daily, not a throwaway exercise
1. **Capture research moments** - The unique "we looked up..." feature

## Tech stack

- **Go** - Backend language
- **Templ** - Type-safe HTML templating for Go
- **HTMX** - Frontend interactivity without JavaScript frameworks
- **SQLite** - Embedded database (modernc.org/sqlite for pure Go)
- **Tailwind CSS** - Styling (or templui components)
- **TMDB API** - Movie metadata, posters, search

## Core features

### Movie logging

- Search for movies by title (live autocomplete via TMDB)
- Auto-populate: poster, year, director, cast, genre
- Add personal rating (1-5 stars)
- Add notes about the viewing
- Track who you watched with

### Research moments (unique feature)

During or after watching, log what you paused to look up:

- "Who plays the lead?" → Actor name + link
- "Where was this filmed?" → Location + Google Maps link
- "Is this based on a true story?" → Answer + Wikipedia link
- Custom questions and discoveries

### History and browsing

- View past entries with filters (date, rating, genre)
- Search your diary
- Infinite scroll for long history
- Quick stats (movies this month, favorite genres)

## Data model

```
Movie
├── ID (internal)
├── TMDBID (external reference)
├── Title
├── Year
├── PosterURL
├── Director
├── Genre
└── Overview

DiaryEntry
├── ID
├── MovieID (foreign key)
├── WatchedAt (date)
├── Rating (1-5)
├── Notes (text)
├── WatchedWith (string, e.g., "Anna")
└── CreatedAt

Lookup (research moments)
├── ID
├── DiaryEntryID (foreign key)
├── Question ("Who plays the detective?")
├── Answer ("Ralph Fiennes")
├── Category (actor, location, trivia, other)
├── URL (optional source link)
└── CreatedAt
```

## HTMX patterns to learn

| Feature | HTMX pattern |
|---------|--------------|
| Movie search autocomplete | `hx-get` with `hx-trigger="keyup changed delay:300ms"` |
| Auto-populate movie details | `hx-get` on selection, swap into form |
| Inline rating change | `hx-post` with `hx-swap="outerHTML"` |
| Add lookup note | `hx-post` with `hx-swap="beforeend"` |
| Delete with confirmation | `hx-delete` with `hx-confirm` |
| Filter diary list | `hx-get` updating table body |
| Infinite scroll | `hx-get` with `hx-trigger="revealed"` |
| Modal for details | `hx-get` into modal target |
| Toast notifications | `hx-swap-oob` for out-of-band updates |

## Feature phases

### Phase 1: Foundation

- [ ] Project setup (Go module, Templ, HTMX, Tailwind)
- [ ] SQLite database with schema
- [ ] Basic movie entry form (manual, no API yet)
- [ ] List view of diary entries
- [ ] Simple HTMX: delete entry, inline edit rating

### Phase 2: Movie search

- [ ] TMDB API integration
- [ ] Live search autocomplete
- [ ] Auto-populate movie details on selection
- [ ] Display movie posters

### Phase 3: Research moments

- [ ] Add lookup notes to entries
- [ ] Categorize lookups (actor, location, trivia)
- [ ] Quick links generation (IMDb, Wikipedia, Google Maps)

### Phase 4: Polish

- [ ] Filter and search diary
- [ ] Infinite scroll
- [ ] Statistics dashboard
- [ ] templui components for consistent design
- [ ] Mobile-responsive layout

## Resources

- [HTMX documentation](https://htmx.org/docs/)
- [Templ documentation](https://templ.guide/)
- [templui components](https://templui.io/)
- [TMDB API](https://developer.themoviedb.org/docs)
- [Laracasts HTMX course](https://laracasts.com/series/crafting-web-applications-with-htmx)

## Development commands (to be added)

```bash
# Run development server
go run .

# Generate Templ templates
templ generate

# Watch Templ files
templ generate --watch

# Build CSS
npx tailwindcss -i input.css -o static/css/tailwind.css --watch

# Run tests
go test ./...
```
