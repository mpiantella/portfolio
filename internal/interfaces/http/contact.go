package interfaces

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ContactRequest is the expected JSON payload for contact submissions.
type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

// Contact handles POST /api/contact
func (h *Handler) Contact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req ContactRequest
	ct := r.Header.Get("Content-Type")
	if strings.Contains(ct, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}
	} else {
		// accept form-encoded from HTMX
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}
		req.Name = r.Form.Get("name")
		req.Email = r.Form.Get("email")
		req.Message = r.Form.Get("message")
	}
	// basic validation
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Message = strings.TrimSpace(req.Message)
	if req.Name == "" || req.Email == "" || req.Message == "" || !strings.Contains(req.Email, "@") {
		http.Error(w, "missing or invalid fields", http.StatusBadRequest)
		return
	}

	// For now, simply acknowledge receipt. In production, enqueue or send email.
	if strings.Contains(r.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		w.WriteHeader(http.StatusAccepted)
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`<div class="text-sm text-emerald-700">Thanks! Your message was received.</div>`))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
}
