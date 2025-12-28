package interfaces_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	mem "portfolio/internal/infrastructure/persistence/memory"
	interfacespkg "portfolio/internal/interfaces/http"
	"html/template"
)

func TestListProjectsJSON(t *testing.T) {
	tmpl := template.Must(template.ParseGlob("../../../web/templates/*.html"))
	repo := mem.NewProjectRepository()
	h := interfacespkg.NewHandler(repo, tmpl)

	req := httptest.NewRequest("GET", "/api/projects", nil)
	w := httptest.NewRecorder()

	h.ListProjectsJSON(w, req)
	res := w.Result()
	if res.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}
	var projects []map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&projects); err != nil {
		t.Fatalf("failed decode: %v", err)
	}
	if len(projects) == 0 {
		t.Fatalf("expected projects, got none")
	}
}
