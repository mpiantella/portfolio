You are acting as a senior software architect and Go expert.

Goal:
Create a professional software architect portfolio website with:
- A modern, responsive UI
- A Go backend using Clean Architecture principles
- Clear separation of concerns
- Production-quality structure, naming, and patterns

Tech Stack Requirements:
Frontend:
- Use a simple, modern UI framework with minimal setup
- **Chosen:** HTMX + Tailwind CSS (lightweight, server-driven UX that integrates with Go templates)
- Apply modern design patterns: clean layout, typography-first, subtle animations
- Mobile-first and accessible

Notes:
- Tailwind CLI is used to build `web/static/dist.css`.
- HTMX is included via CDN to progressively enhance interactions (forms, fragments).

Backend:
- Language: Go
- Architecture: Clean Architecture (Domain, Use Cases, Interfaces, Infrastructure)
- HTTP server using net/http or chi (no heavy frameworks)
- RESTful APIs (JSON)
- Clear boundaries:
  - domain (entities, value objects)
  - usecase (application logic)
  - interface adapters (HTTP handlers)
  - infrastructure (DB, config, logging)
- Dependency inversion enforced via interfaces
- No framework-specific logic in domain or usecase layers

Core Features:
- Home page (architect summary)
- Experience / Projects (API-driven)
- Architecture philosophy / blog-ready structure
- Contact form (POST endpoint)
- Health check endpoint

Quality Requirements:
- Idiomatic Go
- Testable components
- Context-aware request handling
- Structured logging
- Environment-based configuration
- Clear folder structure with README per layer

Deliverables:
1. Suggested project folder structure
2. Initial Go module setup
3. Example entity + use case + HTTP handler
4. Frontend integration example
5. Minimal styling using the selected UI framework
6. Comments explaining architectural decisions

Assume this is a long-lived, maintainable codebase intended to demonstrate senior-level architectural thinking.
Prefer clarity and correctness over cleverness.

UI should feel modern and minimal:
- Neutral color palette
- Subtle shadows and rounded corners
- Large readable typography
- Layout inspired by modern SaaS landing pages
- Avoid clutter and unnecessary animations
