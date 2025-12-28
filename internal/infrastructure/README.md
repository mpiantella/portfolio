# Infrastructure layer

Purpose: provide concrete implementations for interfaces (DB, persistence, logging, config, server wiring).

Contains:
- In-memory project repository (for demo and tests)
- Server wiring and router

Guidelines:
- Keep real integrations (DB, external services) behind interfaces used by the usecase layer.
- Configuration and logging helpers live here.

---

## Routes (current)
The server uses the standard library `net/http` with an `http.ServeMux` (see `internal/infrastructure/server/router.go`). Current routes:

- `GET /projects` → Handler: `interfaces.Handler.ListProjects` — renders the projects page (HTML)
- `GET /projects/fragment` → Handler: `interfaces.Handler.ListProjectsFragment` — returns an HTML fragment for HTMX (`projects_fragment.html`)
- `GET /api/projects` → Handler: `interfaces.Handler.ListProjectsJSON` — returns JSON list of projects
- `POST /api/contact` → Handler: `interfaces.Handler.Contact` — accepts JSON or form-encoded contact submissions
- `GET /api/health` → Handler: `interfaces.Handler.Health` — lightweight JSON health check
- `GET /contact` → Handler: `interfaces.Handler.ContactPage` — renders contact page
- `GET /static/*` → served from `web/static` via `http.FileServer`
- `GET /` → root redirects to the projects page (handled by `ListProjects`)

Notes:
- The router enforces exact matching for `/projects` so `/projects/fragment` is handled separately.
- Templates live in `web/templates/*.html` and are loaded with `html/template`.

## How to add a new route
Follow this simple pattern:

1. **Add a handler method**
   - Implement a method on `*interfaces.Handler` in `internal/interfaces/http/` (e.g. `func (h *Handler) Foo(w http.ResponseWriter, r *http.Request)`). Keep handlers thin: validate input, map transport to usecase calls, and handle errors.
2. **Wire the route**
   - Update the router in `internal/infrastructure/server/router.go` and add a new `mux.HandleFunc("/path", h.Foo)` (or use a wrapper function for method checks).
   - If you need only a specific HTTP method, check `r.Method` inside the handler and return `405 Method Not Allowed` with `w.Header().Set("Allow", http.MethodPost)`.
3. **Add a test**
   - Add unit tests under `internal/interfaces/http/*_test.go` that exercise the handler behaviors.
4. **Add templates / static assets** (if needed)
   - Put HTML templates in `web/templates/` and static assets in `web/static/`. Templates are executed with `h.tmpl.ExecuteTemplate`.
5. **Run and verify**
   - `go test ./...` to run tests
   - `go run ./cmd/server` then hit the route in the browser or via `curl`

Example snippet (router):

```go
mux.HandleFunc("/api/foo", h.Foo)
```

Example snippet (handler skeleton):

```go
func (h *Handler) Foo(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        w.Header().Set("Allow", http.MethodPost)
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }
    // parse, validate, call usecase, respond
}
```

## Tech stack & patterns (backend)
- Language: **Go** (module-aware project)
- HTTP: `net/http` using `http.ServeMux` and `http.FileServer` (small and dependency-free)
- Templating: `html/template` (server-side rendered pages + fragments)
- Persistence: in-memory demo repo at `internal/infrastructure/persistence/memory` — replaceable with a production-backed implementation behind the domain repository interfaces
- Architecture: **Clean-ish layered structure**
  - `internal/domain` — domain entities and repository interfaces
  - `internal/usecase` — application business logic
  - `internal/interfaces` — HTTP adapters/handlers
  - `internal/infrastructure` — concrete infra implementations (server wiring, logging, config, persistence)
- Frontend: HTMX for progressive enhancement and dynamic fragments; Tailwind CSS for styling (see `package.json` and `web/tailwind.css`)

## Tips
- Keep handlers small and delegate business logic to the usecase layer.
- Use table-driven tests for handler behavior (both HTML and JSON endpoints are covered in tests).
- To integrate a DB or external API, implement the repository interface in `internal/infrastructure` and update wiring in `server.NewRouter`.

---
