package interfaces

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"portfolio/internal/domain"
)

// ContactHandler handles contact form-related HTTP requests
type ContactHandler struct {
	templates *template.Template
}

// NewContactHandler creates a new contact handler
func NewContactHandler(templates *template.Template) *ContactHandler {
	return &ContactHandler{
		templates: templates,
	}
}

// ServeHTTP handles the contact page request
func (h *ContactHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Contact Maria Lucena",
	}

	if err := h.templates.ExecuteTemplate(w, "contact.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// ServeAPI handles the contact form submission
func (h *ContactHandler) ServeAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response := domain.APIResponse{
			Success: false,
			Error:   "Method not allowed",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %v", err)
		response := domain.APIResponse{
			Success: false,
			Error:   "Failed to parse form data",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create contact message from form data
	message := domain.ContactMessage{
		Name:             strings.TrimSpace(r.FormValue("name")),
		Email:            strings.TrimSpace(r.FormValue("email")),
		Company:          strings.TrimSpace(r.FormValue("company")),
		Message:          strings.TrimSpace(r.FormValue("message")),
		PreferredContact: r.FormValue("preferred_contact"),
		Topic:            r.FormValue("topic"),
		CreatedAt:        time.Now(),
	}

	// Validate the message
	if err := h.validateContactMessage(&message); err != nil {
		response := domain.APIResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Log the contact submission (in production, you'd save to database or send email)
	log.Printf("Contact form submission: Name=%s, Email=%s, Topic=%s",
		message.Name, message.Email, message.Topic)

	// TODO: In production, implement one of these:
	// 1. Save to database
	// 2. Send email notification
	// 3. Forward to CRM system
	// 4. Send to message queue

	// For now, simulate successful processing
	if err := h.processContactMessage(&message); err != nil {
		log.Printf("Error processing contact message: %v", err)
		response := domain.APIResponse{
			Success: false,
			Error:   "Failed to process your message. Please try again later.",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Return success response
	response := domain.APIResponse{
		Success: true,
		Message: "Thank you for your message! I'll get back to you soon.",
		Data: map[string]interface{}{
			"submitted_at": message.CreatedAt,
			"name":         message.Name,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// validateContactMessage validates the contact form data
func (h *ContactHandler) validateContactMessage(msg *domain.ContactMessage) error {
	// Validate name
	if msg.Name == "" {
		return &ValidationError{Field: "name", Message: "Name is required"}
	}
	if len(msg.Name) < 2 {
		return &ValidationError{Field: "name", Message: "Name must be at least 2 characters"}
	}
	if len(msg.Name) > 100 {
		return &ValidationError{Field: "name", Message: "Name must be less than 100 characters"}
	}

	// Validate email
	if msg.Email == "" {
		return &ValidationError{Field: "email", Message: "Email is required"}
	}
	if !h.isValidEmail(msg.Email) {
		return &ValidationError{Field: "email", Message: "Please enter a valid email address"}
	}

	// Validate message
	if msg.Message == "" {
		return &ValidationError{Field: "message", Message: "Message is required"}
	}
	if len(msg.Message) < 10 {
		return &ValidationError{Field: "message", Message: "Message must be at least 10 characters"}
	}
	if len(msg.Message) > 5000 {
		return &ValidationError{Field: "message", Message: "Message must be less than 5000 characters"}
	}

	// Validate topic if provided
	if msg.Topic != "" {
		validTopics := []string{"Consulting", "Speaking", "Mentorship", "Collaboration", "Other"}
		if !h.contains(validTopics, msg.Topic) {
			return &ValidationError{Field: "topic", Message: "Invalid topic selected"}
		}
	}

	// Validate preferred contact if provided
	if msg.PreferredContact != "" {
		validContacts := []string{"Email", "LinkedIn", "Phone"}
		if !h.contains(validContacts, msg.PreferredContact) {
			return &ValidationError{Field: "preferred_contact", Message: "Invalid contact method selected"}
		}
	}

	// Check for spam patterns (basic implementation)
	if h.isLikelySpam(msg) {
		return &ValidationError{Field: "message", Message: "Your message appears to be spam. Please try again."}
	}

	return nil
}

// isValidEmail validates email format
func (h *ContactHandler) isValidEmail(email string) bool {
	// Simple email regex pattern
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(pattern)
	return re.MatchString(email)
}

// contains checks if a slice contains a string
func (h *ContactHandler) contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// isLikelySpam performs basic spam detection
func (h *ContactHandler) isLikelySpam(msg *domain.ContactMessage) bool {
	// Convert to lowercase for checking
	message := strings.ToLower(msg.Message)
	name := strings.ToLower(msg.Name)

	// Common spam patterns
	spamPatterns := []string{
		"viagra", "cialis", "casino", "lottery", "prize",
		"click here", "buy now", "limited time", "act now",
		"make money fast", "work from home", "miracle cure",
	}

	for _, pattern := range spamPatterns {
		if strings.Contains(message, pattern) || strings.Contains(name, pattern) {
			return true
		}
	}

	// Check for excessive URLs
	urlPattern := regexp.MustCompile(`https?://`)
	urls := urlPattern.FindAllString(msg.Message, -1)
	if len(urls) > 3 {
		return true
	}

	// Check for excessive capitalization
	upperCount := 0
	for _, char := range msg.Message {
		if char >= 'A' && char <= 'Z' {
			upperCount++
		}
	}
	if upperCount > len(msg.Message)/2 && len(msg.Message) > 20 {
		return true
	}

	return false
}

// processContactMessage processes the validated contact message
func (h *ContactHandler) processContactMessage(msg *domain.ContactMessage) error {
	// TODO: Implement actual processing logic
	// Options:
	// 1. Send email notification
	// 2. Save to database
	// 3. Forward to CRM
	// 4. Send to Slack/Discord webhook
	// 5. Add to email marketing list

	// For now, just log it
	log.Printf("Processing contact message from %s (%s)", msg.Name, msg.Email)
	log.Printf("Topic: %s", msg.Topic)
	log.Printf("Message: %s", msg.Message)

	// Example: Send email (implement with your preferred email service)
	// return h.sendEmailNotification(msg)

	return nil
}

// sendEmailNotification sends an email notification (placeholder)
func (h *ContactHandler) sendEmailNotification(msg *domain.ContactMessage) error {
	// TODO: Implement with your email service (SendGrid, AWS SES, Mailgun, etc.)
	/*
		emailBody := fmt.Sprintf(`
			New contact form submission:

			Name: %s
			Email: %s
			Company: %s
			Topic: %s
			Preferred Contact: %s

			Message:
			%s

			Submitted at: %s
		`, msg.Name, msg.Email, msg.Company, msg.Topic, msg.PreferredContact, msg.Message, msg.CreatedAt.Format(time.RFC3339))

		// Send email using your preferred service
		// err := emailService.Send(to, subject, emailBody)
		// return err
	*/

	return nil
}

// ValidationError represents a form validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// GetContactStats returns statistics about contact form submissions
// This would typically query a database
func (h *ContactHandler) GetContactStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// TODO: Query database for actual stats
	// This is placeholder data
	stats := struct {
		TotalSubmissions int            `json:"total_submissions"`
		ByTopic          map[string]int `json:"by_topic"`
		ByMonth          map[string]int `json:"by_month"`
		ResponseRate     float64        `json:"response_rate"`
	}{
		TotalSubmissions: 0,
		ByTopic: map[string]int{
			"Consulting":    0,
			"Speaking":      0,
			"Mentorship":    0,
			"Collaboration": 0,
			"Other":         0,
		},
		ByMonth: map[string]int{
			"2024-12": 0,
		},
		ResponseRate: 0.0,
	}

	response := domain.APIResponse{
		Success: true,
		Data:    stats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
