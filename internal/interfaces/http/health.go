package interfaces

import (
	"encoding/json"
	"net/http"
)

// Health responds with a lightweight JSON status for monitoring.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{"status": "ok"}
	_ = json.NewEncoder(w).Encode(resp)
}
