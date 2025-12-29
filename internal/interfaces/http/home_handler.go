package interfaces

import (
	"html/template"
	"net/http"
)

type HomeHandler struct {
	templates *template.Template
}

func NewHomeHandler(templates *template.Template) *HomeHandler {
	return &HomeHandler{
		templates: templates,
	}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h.templates.ExecuteTemplate(w, "home.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
