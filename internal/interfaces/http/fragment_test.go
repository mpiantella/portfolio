package interfaces_test

import (
	"net/http/httptest"
	"testing"
	"io"
	"strings"

	mem "portfolio/internal/infrastructure/persistence/memory"
	interfacespkg "portfolio/internal/interfaces/http"
	"html/template"
)

func TestListProjectsFragment(t *testing.T) {
	tmpl := template.Must(template.ParseGlob("../../../web/templates/*.html"))
	repo := mem.NewProjectRepository()
	h := interfacespkg.NewHandler(repo, tmpl)

	req := httptest.NewRequest("GET", "/projects/fragment", nil)
	w := httptest.NewRecorder()

	h.ListProjectsFragment(w, req)
	res := w.Result()
	if res.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}
	body, _ := io.ReadAll(res.Body)
	if !strings.Contains(string(body), "Demo Project") {
		t.Fatalf("expected demo project in fragment")
	}
}
