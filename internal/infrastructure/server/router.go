package server

import (
	"html/template"
	"net/http"
	"path/filepath"

	mem "portfolio/internal/infrastructure/persistence/memory"
	interfaces "portfolio/internal/interfaces/http"
)

// NewRouter builds the application's HTTP router.
func NewRouter() http.Handler {
	// simple in-memory repo
	repo := mem.NewProjectRepository()
	tmpl := template.Must(template.ParseGlob(filepath.Join("web", "templates", "*.html")))
	h := interfaces.NewHandler(repo, tmpl)

	mux := http.NewServeMux()
	mux.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		// ensure exact match so subpaths like /projects/fragment map to their own handlers
		if r.URL.Path != "/projects" {
			http.NotFound(w, r)
			return
		}
		h.ListProjects(w, r)
	})
	mux.HandleFunc("/projects/fragment", h.ListProjectsFragment)
	mux.HandleFunc("/api/projects", h.ListProjectsJSON)
	mux.HandleFunc("/api/contact", h.Contact)
	mux.HandleFunc("/api/health", h.Health)
	
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	
	// root redirect
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		h.ListProjects(w, r)
	})
	mux.HandleFunc("/contact", h.ContactPage)
	return mux
}
