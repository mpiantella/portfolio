# Copilot / AI Agent Instructions for this repository ✅

Short, actionable guidance tailored for AI coding agents working on this lightweight Go portfolio app.

## Big picture
- This is a small Go server that renders HTML templates and serves a Tailwind-based UI.
  - Backend: `./cmd/server` + `internal/*` packages (config, logging, server/router, interfaces/http handlers, persistence)
  - Frontend: HTMX + Tailwind templates in `web/templates` and compiled CSS in `web/static/dist.css`
- Data: simple in-memory repository (`internal/infrastructure/persistence/memory`) seeded with demo projects.
- Purpose: small scaffold for a portfolio site; tests and examples assume in-memory data and local templates.

## Key files & patterns (refer to these when making changes) 🔧
- Router: `internal/infrastructure/server/router.go` — wire new endpoints here and use `interfaces.NewHandler`.
- Handlers: `internal/interfaces/http/*` (use `Handler` struct). Handlers typically accept a repo via `h.repo` and templates via `h.tmpl`.
  - Prefer small adapter interfaces for type safety (see `projectRepo` in `project_handler.go`).
- Domain: `internal/domain/project/*` — define entities and repository interface.
- Persistence: `internal/infrastructure/persistence/memory` — simple seeded repo used by tests and development.
- Templates: `web/templates/*.html` — HTMX fragment `projects_fragment.html` is loaded at `/projects/fragment`.
- Static assets: `web/static/dist.css` is the compiled Tailwind output; source is `web/tailwind.css`.

## Developer workflows & commands (exact) 🧭
- Build & run server: `go run ./cmd/server`
- Run tests: `go test ./...` or `make test`
- Tailwind build: `npm run build:css` → outputs `web/static/dist.css`
- Tailwind dev watch: `npm run dev:css` (or `make dev` to run watcher + server; stop with `make stop-dev` or `pkill -f tailwindcss`)
- Common make targets: `make install`, `make build-css`, `make dev-css`, `make run`, `make dev`, `make test`, `make clean` (`Makefile` documents this)

## Project-specific conventions & gotchas ⚠️
- Tests load templates relative to the test package path (e.g. `template.Must(template.ParseGlob("../../../web/templates/*.html"))` in `internal/interfaces/http/*_test.go`). Preserve relative paths when adding tests.
- Handlers use `any` for injected repo; tests and runtime use adapter interfaces for method access (see `projectRepo` in `project_handler.go`).
- `Contact` handler accepts both `application/json` and `application/x-www-form-urlencoded` (HTMX form submissions). When adding contact-related behavior, respect both input types and response formats.
- Environment vars: `PORT` (default `8080`) and `LOG_LEVEL` (`info` by default) from `internal/infrastructure/config`.
- Logger uses `github.com/rs/zerolog` and is configured in `internal/infrastructure/logging`.

## How to add a new backend route (explicit example) ✍️
1. Add domain interface or method if needed in `internal/domain/*`.
2. Implement persistence logic (or mock) in `internal/infrastructure/persistence/*`.
3. Add a handler method in `internal/interfaces/http/<your_handler>.go` and use a small adapter interface for the repo methods required.
4. Wire the handler into `internal/infrastructure/server/router.go` (add exact-path checks if you need subpaths).
5. Add tests in `internal/interfaces/http` package and load templates with the same relative paths as existing tests.
6. Run `go test ./...` and `npm run build:css` (if you changed templates/CSS) and `go run ./cmd/server` to smoke-test locally.

## Guidance for AI edits & PRs ✅
- Keep edits focused and minimal; run `make test` after changes.
- See `.github/agents/beast-mode.md.agent.md` for an existing autonomous agent config and tool usage recommendations.
- When changing templates/CSS: run `npm run build:css` (or `npm run dev:css` in watch mode) and verify UI fragments (`/projects/fragment`).
- Preserve existing test patterns and seeded data in memory repo (tests depend on seeded items).
- Format Go code with `go fmt ./...` and keep import ordering idiomatic.

## Examples of useful quick tasks for agents
- Add a new API endpoint: follow the steps under "How to add a new backend route" and add a unit test.
- Improve contact handler validation: update `Contact` in `internal/interfaces/http/contact.go` and add tests for both JSON and form inputs.
- Add a new project field: update `internal/domain/project/entity.go`, update memory seed, update templates, and adjust any JSON encoding usage.

---
If anything above is unclear or you'd like more examples (e.g., a sample PR for a small change), tell me which area to expand. Thanks! 🎯
