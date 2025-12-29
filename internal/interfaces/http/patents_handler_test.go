// internal/interfaces/http/patents_handler_test.go
package interfaces

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPatentsHandler_ServeAPI(t *testing.T) {
	// Setup
	templates := template.Must(template.New("test").Parse(""))
	handler := NewPatentsHandler(templates, "../../infrastructure/data")

	// Create request
	req := httptest.NewRequest("GET", "/api/patents", nil)
	w := httptest.NewRecorder()

	// Execute
	handler.ServeAPI(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
