# Movie Journal - Phase 2: HTMX Interactivity

## Overview

Add dynamic interactivity using HTMX - update parts of the page without full reloads, all without writing JavaScript.

## Learning goals

By completing these exercises, you will understand:

- How HTMX attributes trigger HTTP requests
- How to return HTML fragments instead of full pages
- How `hx-target` and `hx-swap` control where content goes
- How to build interactive features without JavaScript

## Key HTMX concepts

| Attribute      | Purpose                   | Example                         |
| -------------- | ------------------------- | ------------------------------- |
| `hx-get`       | Make GET request on click | `hx-get="/movies/1"`            |
| `hx-post`      | Make POST request         | `hx-post="/diary/new"`          |
| `hx-target`    | Where to put the response | `hx-target="#movie-details"`    |
| `hx-swap`      | How to insert content     | `hx-swap="innerHTML"` (default) |
| `hx-trigger`   | What triggers the request | `hx-trigger="click"` (default)  |
| `hx-indicator` | Show loading state        | `hx-indicator="#spinner"`       |

## Implementation tasks

### Task 1: Click to expand movie details

**Goal:** Click a movie card to load full details without page reload.

**What to build:**

1. Create a partial template `templates/movie_details.templ` that renders:
   - Full movie overview/synopsis
   - Director, genre, year
   - All research moments (lookups) for this entry
   - A "Close" button

2. Create handler `GET /diary/{id}` that:
   - Returns just the HTML fragment (no layout wrapper)
   - Fetches the diary entry by ID

3. Add HTMX attributes to `MovieCard`:
   - `hx-get="/diary/{id}"`
   - `hx-target="#detail-panel"`
   - `hx-swap="innerHTML"`

4. Add a detail panel div to `index.templ`:

   ```html
   <div id="detail-panel" class="..."></div>
   ```

**Files to modify:**

- `templates/movie_details.templ` (create)
- `templates/movie_card.templ` (add hx-* attributes)
- `templates/index.templ` (add detail panel)
- `internal/handlers/handlers.go` (add GetDiaryEntry handler)
- `internal/server/server.go` (add route)

### Task 2: Delete diary entry with confirmation

**Goal:** Delete an entry with HTMX, removing it from the list dynamically.

**What to build:**

1. Add a delete button to `movie_details.templ`:
   - `hx-delete="/diary/{id}"`
   - `hx-target="closest .movie-card"` (or use ID)
   - `hx-swap="outerHTML"` (removes the element)
   - `hx-confirm="Delete this entry?"` (browser confirm dialog)

2. Create handler `DELETE /diary/{id}` that:
   - Deletes from sample data (or returns empty response)
   - Returns empty body (element gets removed)

**Files to modify:**

- `templates/movie_details.templ` (add delete button)
- `internal/handlers/handlers.go` (add DeleteDiaryEntry handler)
- `internal/server/server.go` (add route)

### Task 3: Add loading indicator

**Goal:** Show a spinner while content loads.

**What to build:**

1. Add a loading spinner to the detail panel:

   ```html
   <div id="detail-panel">
     <div id="spinner" class="htmx-indicator">Loading...</div>
   </div>
   ```

2. Add `hx-indicator="#spinner"` to the movie card click

3. Add CSS for `.htmx-indicator` (hidden by default, shown during request)

**Files to modify:**

- `templates/index.templ` (add spinner)
- `templates/movie_card.templ` (add hx-indicator)
- `input.css` (add indicator styles)

### Task 4: Close detail panel

**Goal:** Add a close button that clears the detail panel.

**What to build:**

1. In `movie_details.templ`, add close button:
   - `hx-get="/empty"` or use `hx-on:click` to clear
   - `hx-target="#detail-panel"`
   - `hx-swap="innerHTML"`

2. Create simple handler that returns empty HTML or use HTMX's client-side features

**Alternative approach:** Use `hx-on:click="this.closest('#detail-panel').innerHTML=''"` for pure client-side close.

## Challenges (after completing tasks)

### Challenge 1: Toggle expand/collapse

Instead of a separate detail panel, make the movie card expand in place when clicked, and collapse when clicked again.

Hints:

- Use `hx-swap="afterend"` to insert details after the card
- Track expanded state with a class or data attribute
- Second click removes the expanded content

### Challenge 2: Keyboard shortcut

Add keyboard navigation - press `Escape` to close the detail panel.

Hints:

- HTMX has `hx-trigger="keyup[key=='Escape'] from:body"`
- Can be added to a hidden element that listens for the key

### Challenge 3: Add a rating filter

Add buttons to filter movies by rating (show only 5-star, 4-star, etc.).

Hints:

- `hx-get="/diary?rating=5"`
- `hx-target="#entries-list"`
- Return just the filtered movie cards (partial template)

## Success criteria

After Phase 2 completion:

- [ ] Clicking a movie card loads details without page reload
- [ ] Detail panel shows expanded movie info
- [ ] Loading indicator appears during fetch
- [ ] Close button clears the detail panel
- [ ] Delete button removes entry from list
- [ ] No custom JavaScript written

## Reference

- HTMX docs: https://htmx.org/docs/
- HTMX examples: https://htmx.org/examples/
- HTMX attributes reference: https://htmx.org/reference/
