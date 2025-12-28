package interfaces

import "net/http"

// ContactPage renders the contact form page.
func (h *Handler) ContactPage(w http.ResponseWriter, r *http.Request) {
	if err := h.tmpl.ExecuteTemplate(w, "layout.html", map[string]any{"Content": "contact"}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
