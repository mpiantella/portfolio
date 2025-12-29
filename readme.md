# Portfolio (scaffold)

Lightweight Go project scaffold for a simple portfolio.

```bash
portfolio-go/
├── cmd/server/
│   └── main.go                    # Entry point (modify to add new routes)
├── internal/
│   ├── domain/
│   │   └── project.go             # ADD: Domain models (provided)
│   ├── interfaces/http/
│   │   ├── home_handler.go        # ADD: Home page handler
│   │   ├── projects_handler.go    # ADD: Projects handler
│   │   ├── patents_handler.go     # ADD: Patents handler
│   │   ├── speaking_handler.go    # ADD: Speaking handler
│   │   └── contact_handler.go     # MODIFY: Existing handler
│   └── infrastructure/
│       ├── data/
│       │   ├── projects.json      # ADD: Project data
│       │   ├── patents.json       # ADD: Patents data
│       │   └── speaking.json      # ADD: Speaking data
│       └── server/
│           └── router.go          # MODIFY: Add new routes
├── web/
│   ├── templates/
│   │   ├── layouts/
│   │   │   └── base.html          # MODIFY: Add navigation
│   │   ├── components/
│   │   │   ├── nav.html           # ADD: Navigation component
│   │   │   └── footer.html        # ADD: Footer component
│   │   └── pages/
│   │       ├── home.html          # ADD: Landing page (provided)
│   │       ├── projects.html      # ADD: Projects listing (provided)
│   │       ├── patents.html       # ADD: Patents page
│   │       └── speaking.html      # ADD: Speaking page
│   └── static/
│       ├── css/
│       │   └── input.css          # MODIFY: Add custom animations
│       └── images/                # ADD: Your images
└── data/                          # ADD: JSON data files
```

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
