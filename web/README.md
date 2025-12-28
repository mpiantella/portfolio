# Frontend (templates & assets)

This folder contains HTML templates, Tailwind input CSS, and static assets used by the Go server.

## Structure
- `templates/` — server-side HTML templates (rendered by `html/template`). Main files:
  - `layout.html` — app shell and HTMX setup
  - `projects.html` and `_projects_fragment.html` — projects full page and fragment
  - `contact.html` — contact form
- `static/` — compiled static files (dist.css) and other assets served at `/static/`
- `tailwind.css` — Tailwind entry file used to build `dist.css`

## Patterns & usage ✅
- **Progressive enhancement with HTMX**: The `layout.html` uses HTMX to load the projects fragment (`hx-get="/projects/fragment" hx-swap="innerHTML"`) — this keeps the initial page server-rendered and enables incremental updates without a full SPA.
- **Server-side templates**: Use `{{ template "name" . }}` and data maps (e.g. `Projects`) passed from handlers.
- **Tailwind CSS**: Tailwind is used for styling; compiled CSS is referenced as `/static/dist.css` in `layout.html`.

## Development (start here) 🔧
1. Install Node dependencies (only dev tooling for Tailwind):
   ```bash
   cd portfolio
   npm install
   ```
2. Dev workflow:
   - Run Tailwind in watch mode to auto-rebuild CSS while you edit classes:
     ```bash
     npm run dev:css
     ```
   - Run the Go server:
     ```bash
     go run ./cmd/server
     ```
   - Open `http://localhost:8080` (or the PORT you set) and make changes to `web/templates` or CSS.

   - Tailwind config: `tailwind.config.js` at the repository root is the canonical Tailwind configuration used when building CSS (it includes content paths for templates and Go files).
3. Build for production (minified CSS):
   ```bash
   npm run build:css
   ```

## How to add UI fragments or pages
1. Add/modify a template in `web/templates/*.html` (create a named template with `{{ define "name" }}`).
2. Add a handler in `internal/interfaces/http` that renders the template using `h.tmpl.ExecuteTemplate`.
3. Wire a new route in `internal/infrastructure/server/router.go` to point to the handler.
4. If you plan to fetch partial HTML with HTMX, return a fragment template that can be swapped into the DOM.

## Debugging & tips 💡
- Use the browser dev tools to inspect HTMX requests and verify responses.
- When changing CSS classes, ensure the Tailwind watcher is running or re-run the build script to regenerate `dist.css`.
- Keep fragments small and focused; they are easiest to test and reuse.

---

If you want, I can add a short example that shows how to add a new HTMX fragment endpoint and template.