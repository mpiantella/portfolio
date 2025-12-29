package interfaces

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"portfolio/internal/domain"
)

func TestContactHandler_ValidateEmail(t *testing.T) {
	handler := &ContactHandler{}

	tests := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"user+tag@domain.co.uk", true},
		{"invalid.email", false},
		{"@example.com", false},
		{"user@", false},
	}

	for _, tt := range tests {
		result := handler.isValidEmail(tt.email)
		if result != tt.valid {
			t.Errorf("isValidEmail(%s) = %v, want %v", tt.email, result, tt.valid)
		}
	}
}

func TestContactHandler_SpamDetection(t *testing.T) {
	handler := &ContactHandler{}

	spamMessage := &domain.ContactMessage{
		Name:    "Spammer",
		Email:   "spam@example.com",
		Message: "BUY VIAGRA NOW! Click here for amazing deals!",
	}

	if !handler.isLikelySpam(spamMessage) {
		t.Error("Expected spam message to be detected")
	}

	validMessage := &domain.ContactMessage{
		Name:    "John Doe",
		Email:   "john@example.com",
		Message: "I'm interested in speaking opportunities.",
	}

	if handler.isLikelySpam(validMessage) {
		t.Error("Valid message incorrectly marked as spam")
	}
}

func TestContactHandler_FormSubmission(t *testing.T) {
	templates := template.Must(template.New("test").Parse(""))
	handler := NewContactHandler(templates)

	formData := url.Values{
		"name":              {"John Doe"},
		"email":             {"john@example.com"},
		"company":           {"Acme Corp"},
		"message":           {"I'd like to discuss a speaking opportunity."},
		"topic":             {"Speaking"},
		"preferred_contact": {"Email"},
	}

	req := httptest.NewRequest("POST", "/api/contact", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	handler.ServeAPI(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
