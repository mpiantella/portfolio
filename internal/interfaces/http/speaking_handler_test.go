// internal/interfaces/http/speaking_handler_test.go
package interfaces

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSpeakingHandler_ServeAPI(t *testing.T) {
	templates := template.Must(template.New("test").Parse(""))
	handler := NewSpeakingHandler(templates, "../../infrastructure/data")

	req := httptest.NewRequest("GET", "/api/speaking", nil)
	w := httptest.NewRecorder()

	handler.ServeAPI(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestSpeakingHandler_FilterByType(t *testing.T) {
	templates := template.Must(template.New("test").Parse(""))
	handler := NewSpeakingHandler(templates, "../../infrastructure/data")

	req := httptest.NewRequest("GET", "/api/speaking?type=Conference", nil)
	w := httptest.NewRecorder()

	handler.ServeAPI(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Parse response and verify all items are conferences
}

func TestSpeakingHandler_GetUpcoming(t *testing.T) {
	templates := template.Must(template.New("test").Parse(""))
	handler := NewSpeakingHandler(templates, "../../infrastructure/data")

	req := httptest.NewRequest("GET", "/api/speaking/upcoming", nil)
	w := httptest.NewRecorder()

	handler.GetUpcomingEngagements(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
