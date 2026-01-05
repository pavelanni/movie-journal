# Movie Journal - HTMX & Templ Learning Plan

## Overview

A guided exploration of the Movie Journal codebase to learn how Templ and HTMX work together to create dynamic web applications without writing JavaScript.

## Learning goals

By completing these exercises, you will understand:

- How Templ templates compile to Go code
- How to create reusable components in Templ
- How HTMX enables dynamic updates without JavaScript
- The request/response cycle between HTMX and Go handlers

## Part 1: Templ fundamentals

### Exercise 1.1: Template structure

**Files to explore:** `templates/layout.templ`, `templates/index.templ`

**Questions:**

1. Open `templates/layout.templ`. What is the function signature of the `Layout` templ? What does `children ...templ.Component` mean?

1. Find where HTMX is loaded in the layout. What HTML tag is used and what is the `src` attribute?

1. Look at `templates/index.templ`. How does it use the `Layout` component? What syntax wraps the page content?

1. Run `templ generate` and look at the generated file `templates/layout_templ.go`. Can you find the Go function that corresponds to your templ? How does it render HTML?

### Exercise 1.2: Components and data

**Files to explore:** `templates/movie_card.templ`, `internal/models/models.go`

**Questions:**

1. What data type does `MovieCard` accept as a parameter? Trace it back to `internal/models/models.go` - what fields does `DiaryEntry` have?

1. In `movie_card.templ`, find the conditional rendering. What syntax does Templ use for `if` statements? How does it differ from Go's normal syntax?

1. Find the `for` loop in `movie_card.templ` (hint: look at `StarRating`). How does Templ handle loops?

1. How does Templ handle dynamic attribute values? Look at lines like `src={ entry.Movie.PosterURL }` - what do the curly braces do?

### Exercise 1.3: Template compilation

**Questions:**

1. What file extension do Templ source files use? What extension do the generated Go files have?

1. In `go.mod`, find the templ import. What package path is it?

1. Look at `.gitignore` - are the generated `*_templ.go` files committed to git? Why or why not?

1. In the `Makefile`, find the `templ-generate` target. When is this run during the build process?

## Part 2: Understanding the request flow

### Exercise 2.1: From URL to template

**Files to explore:** `internal/server/server.go`, `internal/handlers/handlers.go`

**Questions:**

1. In `server.go`, find where routes are registered. What pattern matches the home page (`/`)?

1. Trace the home page request: which handler function is called? Where is it defined?

1. In `handlers.go`, find the `Home` function. How does it render the templ template? What method is called on the template?

1. What two arguments does `.Render()` receive? Why does it need the request context?

### Exercise 2.2: Sample data flow

**Files to explore:** `internal/handlers/handlers.go`

**Questions:**

1. Find `getSampleEntries()` in handlers.go. What data structure does it return?

1. How is this data passed to the `Index` template?

1. In `index.templ`, how does the template access this data? Find the loop that iterates over entries.

1. What would happen if you passed an empty slice? Find the conditional that handles this case.

## Part 3: HTMX basics (preparation for Phase 2)

### Exercise 3.1: HTMX setup

**Files to explore:** `templates/layout.templ`, `static/js/htmx.min.js`

**Questions:**

1. Where is HTMX loaded? Is it loaded before or after the page content?

1. Check the HTMX file size: `ls -lh static/js/htmx.min.js`. How does this compare to typical JavaScript frameworks?

1. HTMX works through HTML attributes. Look at the existing templates - do you see any `hx-*` attributes yet? (Hint: not yet - we'll add them in Phase 2)

### Exercise 3.2: HTMX concepts (theory)

**Research questions** (use https://htmx.org/docs/ as reference):

1. What does `hx-get` do? How is it different from a regular `<a href="...">`?

1. What does `hx-target` do? Why is this powerful?

1. What does `hx-swap` do? What are some common values?

1. What is the main difference between HTMX and traditional JavaScript frameworks like React?

## Part 4: Hands-on challenges

After completing the exploration exercises, try these small modifications:

### Challenge 1: Add a field

Add a `Location` field to display where you watched the movie (e.g., "home", "theater", "friend's house").

- Modify `internal/models/models.go`
- Update sample data in `handlers.go`
- Display it in `movie_card.templ`

### Challenge 2: Conditional styling

Make the star rating show different colors based on the rating:

- 1-2 stars: red/orange
- 3 stars: yellow
- 4-5 stars: green

Hint: Use Templ conditionals to change the Tailwind class.

### Challenge 3: Add a new page

Create a simple "About" page:

- Create `templates/about.templ`
- Add a handler in `handlers.go`
- Add a route in `server.go`
- Add a navigation link in `layout.templ`

## Answer key location

After attempting the exercises, you can ask me to explain any answers you're unsure about.

## Next steps (Phase 2 preview)

Once comfortable with the basics, Phase 2 will add:

- `hx-get` for loading movie details without page refresh
- `hx-post` for submitting new diary entries
- `hx-trigger` for search-as-you-type
- Partial template rendering (returning HTML fragments)
