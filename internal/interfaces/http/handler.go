package interfaces

import (
	"html/template"
)

// Handler groups http handlers and dependencies.
type Handler struct {
	tmpl *template.Template
	repo any
}

// NewHandler constructs a new Handler.
func NewHandler(repo any, tmpl *template.Template) *Handler {
	return &Handler{tmpl: tmpl, repo: repo}
}

