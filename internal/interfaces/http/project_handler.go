package interfaces

import (
	"encoding/json"
	"net/http"
	"portfolio/internal/domain"
)

// For type safety, create a small adapter interface the handler expects.
type projectRepo interface {
	List() ([]domain.Project, error)
}

// ListProjects handles GET /projects (HTML)
func (h *Handler) ListProjects(w http.ResponseWriter, r *http.Request) {
	repo, ok := h.repo.(projectRepo)
	if !ok {
		http.Error(w, "repository not available", http.StatusInternalServerError)
		return
	}

	projects, err := repo.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]any{"Projects": projects}

	w.Header().Set("X-Handler", "projects")
	if err := h.tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ListProjectsJSON handles GET /api/projects and returns JSON
func (h *Handler) ListProjectsJSON(w http.ResponseWriter, r *http.Request) {
	repo, ok := h.repo.(projectRepo)
	if !ok {
		http.Error(w, "repository not available", http.StatusInternalServerError)
		return
	}
	projects, err := repo.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(projects)
}

// ListProjectsFragment returns an HTML fragment suitable for HTMX requests.
func (h *Handler) ListProjectsFragment(w http.ResponseWriter, r *http.Request) {
	repo, ok := h.repo.(projectRepo)
	if !ok {
		http.Error(w, "repository not available", http.StatusInternalServerError)
		return
	}
	projects, err := repo.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := map[string]any{"Projects": projects}
	w.Header().Set("X-Handler", "projects-fragment")
	if err := h.tmpl.ExecuteTemplate(w, "projects_fragment.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
