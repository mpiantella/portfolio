Bootstrapping with Make ✅
Use make to simplify common tasks:

* make install — install Node dev deps
* make build-css — build Tailwind CSS to dist.css
* make dev-css — watch and rebuild CSS during development
* make run — run the Go server (go run ./cmd/server)
* make dev — starts Tailwind watcher in background and runs the server (stop watcher with make stop-dev)
* make test — run go test
* make clean — remove generated CSS

Note: make dev backgrounds the Tailwind process; use make stop-dev or pkill -f tailwindcss to stop it.



# Portfolio (scaffold)

Lightweight Go project scaffold for a simple portfolio.

Run server:

```bash
# builds and runs the Go server (reads PORT and LOG_LEVEL from env)
go run ./cmd/server
```

Frontend assets (Tailwind):

```bash
# install dev dependencies (requires node/npm)
npm install
# build the Tailwind CSS into `web/static/dist.css`
npm run build:css
```

Run tests:

```bash
go test ./...
```

Notes:
- Templates use HTMX and Tailwind (CDN HTMX included). The contact form uses HTMX to POST to `/api/contact` and shows an inline response.
- Health endpoint: `GET /api/health` returns JSON `{ "status": "ok" }`.
- Projects API: `GET /api/projects` returns JSON list of projects.

## Developer quickstart
- Start Tailwind in watch mode (auto rebuild CSS during development):
  ```bash
  npm run dev:css
  ```
- In another shell, run the Go server:
  ```bash
  go run ./cmd/server
  ```
- Edit templates under `web/templates` and static assets in `web/static`. HTMX fragments are loaded from `/projects/fragment`.
- Add backend routes by: (1) implementing a handler method in `internal/interfaces/http`, (2) wiring it in `internal/infrastructure/server/router.go`, and (3) adding tests in the `internal/interfaces/http` package.
