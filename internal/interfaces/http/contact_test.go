package interfaces_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	interfacespkg "portfolio/internal/interfaces/http"
	"html/template"
)

func TestContactHandlerForm(t *testing.T) {
	// load templates relative to project root
	tmpl := template.Must(template.ParseGlob("../../../web/templates/*.html"))
	h := interfacespkg.NewHandler(nil, tmpl)

	req := httptest.NewRequest(http.MethodPost, "/api/contact", strings.NewReader("name=Test&email=test@example.com&message=hello"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	h.Contact(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusAccepted {
		t.Fatalf("expected 202 got %d", res.StatusCode)
	}
}

func TestContactHandlerJSON(t *testing.T) {
	// load templates relative to project root
	tmpl := template.Must(template.ParseGlob("../../../web/templates/*.html"))
	h := interfacespkg.NewHandler(nil, tmpl)

	json := `{"name":"Jane","email":"jane@example.com","message":"hi"}`
	req := httptest.NewRequest(http.MethodPost, "/api/contact", strings.NewReader(json))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Contact(w, req)
	res := w.Result()
	if res.StatusCode != http.StatusAccepted {
		t.Fatalf("expected 202 got %d", res.StatusCode)
	}
}
