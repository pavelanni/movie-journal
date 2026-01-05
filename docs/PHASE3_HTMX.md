# Movie Journal - Phase 3: Forms and data entry

## Overview

Build interactive forms using HTMX - create and edit diary entries with real-time validation and feedback, all without page reloads or custom JavaScript.

## Learning goals

By completing these exercises, you will understand:

- How to submit forms with HTMX (`hx-post`, `hx-put`)
- Server-side validation with inline error feedback
- Swapping form content after successful submission
- Using `hx-trigger` for input events (debounce, changed)
- Building search/autocomplete functionality

## Key HTMX concepts for forms

| Attribute                                | Purpose                        | Example                      |
| ---------------------------------------- | ------------------------------ | ---------------------------- |
| `hx-post`                                | Submit form via POST           | `hx-post="/diary/new"`       |
| `hx-put`                                 | Submit form via PUT            | `hx-put="/diary/1"`          |
| `hx-trigger="submit"`                    | Trigger on form submit         | Default for forms            |
| `hx-trigger="input changed delay:500ms"` | Debounced input                | For search fields            |
| `hx-include`                             | Include additional inputs      | `hx-include="[name='csrf']"` |
| `hx-disabled-elt`                        | Disable element during request | `hx-disabled-elt="this"`     |
| `hx-vals`                                | Add extra values to request    | `hx-vals='{"type":"movie"}'` |

## Implementation tasks

### Task 1: Create new diary entry form

**Goal:** Build a form to log a new movie watch.

**What to build:**

1. Create template `templates/diary_form.templ`:
   - Movie title input
   - Date watched (date picker)
   - Location input
   - Rating selector (1-5 stars)
   - "Watched with" input
   - Notes textarea
   - Submit button

2. Create template `templates/diary_new.templ`:
   - Uses `@Layout("Log a Movie")`
   - Contains the form

3. Create handlers:
   - `GET /diary/new` - renders the empty form page
   - `POST /diary/new` - processes form submission

4. For now, just print the submitted data to console and redirect to home

**Files to create/modify:**

- `templates/diary_form.templ` (create)
- `templates/diary_new.templ` (create)
- `internal/handlers/handlers.go` (add handlers)
- `internal/server/server.go` (add routes)

### Task 2: Add form validation

**Goal:** Show validation errors without page reload.

**What to build:**

1. Modify POST handler to validate:
   - Movie title is required
   - Rating is between 1-5
   - Date is not in the future

2. If validation fails:
   - Return the form HTML with error messages
   - Preserve user's input values
   - Highlight invalid fields (red border)

3. Add HTMX attributes to form:
   - `hx-post="/diary/new"`
   - `hx-target="this"` (replace form with response)
   - `hx-swap="outerHTML"`

4. Create error display component:
   - Inline error messages below each field
   - Or summary at top of form

**Files to modify:**

- `templates/diary_form.templ` (add error display)
- `internal/handlers/handlers.go` (add validation logic)

### Task 3: Success feedback

**Goal:** Show success message and update the page after creating an entry.

**What to build:**

1. On successful form submission, handler should:
   - Add entry to sample data (or just acknowledge)
   - Return success HTML fragment

2. Options for success feedback:
   - Replace form with success message + "Add another" link
   - Redirect to home page using `HX-Redirect` header
   - Show toast/notification (more advanced)

3. Try the `HX-Redirect` header approach:

   ```go
   w.Header().Set("HX-Redirect", "/")
   ```

**Files to modify:**

- `internal/handlers/handlers.go`
- `templates/diary_form.templ` (optional success state)

### Task 4: Rating selector component

**Goal:** Build an interactive star rating selector.

**What to build:**

1. Create a clickable star rating input:
   - 5 stars displayed
   - Clicking a star selects that rating
   - Visual feedback (filled vs empty stars)
   - Hidden input field holds the actual value

2. This can be done with:
   - Pure CSS/HTML (radio buttons styled as stars)
   - Or minimal HTMX for interactivity

3. Start with radio buttons approach:

   ```html
   <input type="radio" name="rating" value="1" id="star1">
   <label for="star1">★</label>
   ```

**Files to modify:**

- `templates/diary_form.templ` (add star selector)
- `input.css` (add star styling)

## Challenges (after completing tasks)

### Challenge 1: Edit existing entry

Build an edit form for existing diary entries.

**Hints:**

- `GET /diary/{id}/edit` - renders form pre-filled with entry data
- `PUT /diary/{id}` - processes the update
- Reuse `diary_form.templ` with an `entry` parameter (nil for new, populated for edit)
- Add "Edit" button to `movie_details.templ`

### Challenge 2: Movie search autocomplete

Add a search field that suggests movies as you type.

**Hints:**

- Create search input with `hx-trigger="input changed delay:300ms"`
- `hx-get="/movies/search?q={value}"` with `hx-target="#suggestions"`
- Handler returns list of matching movies from sample data
- Clicking a suggestion fills in the movie details
- Consider using `hx-include` or `hx-vals` to pass search query

### Challenge 3: Add research moments (lookups)

Allow users to add multiple lookups to a diary entry.

**Hints:**

- "Add Lookup" button that appends a new lookup form fragment
- Use `hx-get="/lookups/new"` with `hx-swap="beforeend"` on container
- Each lookup has: question, answer, category (dropdown)
- Consider using `hx-target="#lookups-list"`

### Challenge 4: Inline editing

Edit movie notes directly on the detail view without a separate form page.

**Hints:**

- Click notes text → transforms into textarea
- `hx-get="/diary/{id}/notes/edit"` returns the edit fragment
- `hx-put="/diary/{id}/notes"` saves and returns the display fragment
- `hx-trigger="blur"` or explicit save button

## Success criteria

After Phase 3 completion:

- [ ] Can create new diary entries via form
- [ ] Form shows validation errors inline
- [ ] Success redirects or shows confirmation
- [ ] Star rating selector works
- [ ] At least one challenge completed
- [ ] No custom JavaScript written

## Reference

- HTMX form examples: https://htmx.org/examples/
- HTMX hx-post: https://htmx.org/attributes/hx-post/
- HTMX response headers: https://htmx.org/reference/#response_headers
- Go form handling: https://pkg.go.dev/net/http#Request.FormValue
