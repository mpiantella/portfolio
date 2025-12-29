package interfaces

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"portfolio/internal/domain"
)

type ProjectsHandler struct {
	templates *template.Template
	dataPath  string
}

func NewProjectsHandler(templates *template.Template, dataPath string) *ProjectsHandler {
	return &ProjectsHandler{
		templates: templates,
		dataPath:  dataPath,
	}
}

func (h *ProjectsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := h.loadProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "projects.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ProjectsHandler) ServeAPI(w http.ResponseWriter, r *http.Request) {
	data, err := h.loadProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Filter featured projects if requested
	featured := r.URL.Query().Get("featured")
	if featured == "true" {
		filtered := make([]domain.Project, 0)
		for _, p := range data.Projects {
			if p.Featured {
				filtered = append(filtered, p)
			}
		}
		data.Projects = filtered
	}

	// Check if this is an HTMX request
	isHTMX := r.Header.Get("HX-Request") == "true"

	if isHTMX {
		// Return HTML fragment for HTMX
		w.Header().Set("Content-Type", "text/html")
		if err := h.templates.ExecuteTemplate(w, "project-card", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	// Return JSON for regular API requests
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (h *ProjectsHandler) loadProjects() (*struct {
	Projects []domain.Project `json:"projects"`
}, error) {
	file, err := ioutil.ReadFile(h.dataPath + "/projects.json")
	if err != nil {
		return nil, err
	}

	var data struct {
		Projects []domain.Project `json:"projects"`
	}

	if err := json.Unmarshal(file, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
