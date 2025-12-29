package interfaces

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"time"

	"portfolio/internal/domain"
)

// SpeakingHandler handles speaking engagement-related HTTP requests
type SpeakingHandler struct {
	templates *template.Template
	dataPath  string
}

// NewSpeakingHandler creates a new speaking handler
func NewSpeakingHandler(templates *template.Template, dataPath string) *SpeakingHandler {
	return &SpeakingHandler{
		templates: templates,
		dataPath:  dataPath,
	}
}

// ServeHTTP handles the speaking page request
func (h *SpeakingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := h.loadSpeakingEngagements()
	if err != nil {
		log.Printf("Error loading speaking engagements: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Sort by date (most recent first)
	sort.Slice(data.SpeakingEngagements, func(i, j int) bool {
		return data.SpeakingEngagements[i].Date.After(data.SpeakingEngagements[j].Date)
	})

	// Calculate stats for template
	stats := h.calculatePageStats(data)

	templateData := struct {
		SpeakingEngagements []domain.SpeakingEngagement
		TotalTalks          int
		ConferenceCount     int
		TotalAudience       int
		UniqueTopicsCount   int
	}{
		SpeakingEngagements: data.SpeakingEngagements,
		TotalTalks:          stats.TotalTalks,
		ConferenceCount:     stats.ConferenceCount,
		TotalAudience:       stats.TotalAudience,
		UniqueTopicsCount:   stats.UniqueTopicsCount,
	}

	if err := h.templates.ExecuteTemplate(w, "speaking.html", templateData); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// calculatePageStats calculates stats for the speaking page
func (h *SpeakingHandler) calculatePageStats(data *domain.SpeakingData) struct {
	TotalTalks        int
	ConferenceCount   int
	TotalAudience     int
	UniqueTopicsCount int
} {
	stats := struct {
		TotalTalks        int
		ConferenceCount   int
		TotalAudience     int
		UniqueTopicsCount int
	}{
		TotalTalks: len(data.SpeakingEngagements),
	}

	topicsMap := make(map[string]bool)

	for _, s := range data.SpeakingEngagements {
		if s.Type == "Conference" {
			stats.ConferenceCount++
		}
		stats.TotalAudience += s.AudienceSize

		for _, topic := range s.Topics {
			topicsMap[topic] = true
		}
	}

	stats.UniqueTopicsCount = len(topicsMap)

	return stats
}

// ServeAPI handles the API request for speaking engagements data
func (h *SpeakingHandler) ServeAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := h.loadSpeakingEngagements()
	if err != nil {
		log.Printf("Error loading speaking engagements: %v", err)
		response := domain.APIResponse{
			Success: false,
			Error:   "Failed to load speaking engagements data",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Filter by type if requested
	eventType := r.URL.Query().Get("type")
	if eventType != "" {
		filtered := make([]domain.SpeakingEngagement, 0)
		for _, s := range data.SpeakingEngagements {
			if s.Type == eventType {
				filtered = append(filtered, s)
			}
		}
		data.SpeakingEngagements = filtered
	}

	// Filter by year if requested
	year := r.URL.Query().Get("year")
	if year != "" {
		filtered := make([]domain.SpeakingEngagement, 0)
		for _, s := range data.SpeakingEngagements {
			if s.Date.Format("2006") == year {
				filtered = append(filtered, s)
			}
		}
		data.SpeakingEngagements = filtered
	}

	// Filter future/past events
	timeFilter := r.URL.Query().Get("time")
	if timeFilter != "" {
		filtered := make([]domain.SpeakingEngagement, 0)
		now := time.Now()
		for _, s := range data.SpeakingEngagements {
			if timeFilter == "upcoming" && s.Date.After(now) {
				filtered = append(filtered, s)
			} else if timeFilter == "past" && s.Date.Before(now) {
				filtered = append(filtered, s)
			}
		}
		data.SpeakingEngagements = filtered
	}

	// Sort by date (most recent first)
	sort.Slice(data.SpeakingEngagements, func(i, j int) bool {
		return data.SpeakingEngagements[i].Date.After(data.SpeakingEngagements[j].Date)
	})

	response := domain.APIResponse{
		Success: true,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSpeakingByID returns a single speaking engagement by ID
func (h *SpeakingHandler) GetSpeakingByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path (assumes route like /api/speaking/{id})
	id := r.URL.Path[len("/api/speaking/"):]
	if id == "" {
		http.Error(w, "Speaking engagement ID is required", http.StatusBadRequest)
		return
	}

	data, err := h.loadSpeakingEngagements()
	if err != nil {
		log.Printf("Error loading speaking engagements: %v", err)
		response := domain.APIResponse{
			Success: false,
			Error:   "Failed to load speaking engagements data",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Find speaking engagement by ID
	var foundEngagement *domain.SpeakingEngagement
	for _, s := range data.SpeakingEngagements {
		if s.ID == id {
			foundEngagement = &s
			break
		}
	}

	if foundEngagement == nil {
		response := domain.APIResponse{
			Success: false,
			Error:   "Speaking engagement not found",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := domain.APIResponse{
		Success: true,
		Data:    foundEngagement,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSpeakingStats returns statistics about speaking engagements
func (h *SpeakingHandler) GetSpeakingStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := h.loadSpeakingEngagements()
	if err != nil {
		log.Printf("Error loading speaking engagements: %v", err)
		response := domain.APIResponse{
			Success: false,
			Error:   "Failed to load speaking engagements data",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Calculate statistics
	stats := struct {
		TotalEngagements  int            `json:"total_engagements"`
		ConferenceCount   int            `json:"conference_count"`
		WebinarCount      int            `json:"webinar_count"`
		InternalCount     int            `json:"internal_count"`
		TotalAudience     int            `json:"total_audience"`
		UpcomingCount     int            `json:"upcoming_count"`
		EngagementsByYear map[string]int `json:"engagements_by_year"`
		UniqueTopics      []string       `json:"unique_topics"`
	}{
		TotalEngagements:  len(data.SpeakingEngagements),
		EngagementsByYear: make(map[string]int),
	}

	topicsMap := make(map[string]bool)
	now := time.Now()

	for _, s := range data.SpeakingEngagements {
		// Count by type
		switch s.Type {
		case "Conference":
			stats.ConferenceCount++
		case "Webinar":
			stats.WebinarCount++
		case "Internal":
			stats.InternalCount++
		}

		// Total audience
		stats.TotalAudience += s.AudienceSize

		// Upcoming events
		if s.Date.After(now) {
			stats.UpcomingCount++
		}

		// By year
		year := s.Date.Format("2006")
		stats.EngagementsByYear[year]++

		// Unique topics
		for _, topic := range s.Topics {
			topicsMap[topic] = true
		}
	}

	// Convert topics map to slice
	for topic := range topicsMap {
		stats.UniqueTopics = append(stats.UniqueTopics, topic)
	}
	sort.Strings(stats.UniqueTopics)

	response := domain.APIResponse{
		Success: true,
		Data:    stats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetUpcomingEngagements returns upcoming speaking engagements
func (h *SpeakingHandler) GetUpcomingEngagements(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := h.loadSpeakingEngagements()
	if err != nil {
		log.Printf("Error loading speaking engagements: %v", err)
		response := domain.APIResponse{
			Success: false,
			Error:   "Failed to load speaking engagements data",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Filter upcoming events
	upcoming := make([]domain.SpeakingEngagement, 0)
	now := time.Now()
	for _, s := range data.SpeakingEngagements {
		if s.Date.After(now) {
			upcoming = append(upcoming, s)
		}
	}

	// Sort by date (soonest first)
	sort.Slice(upcoming, func(i, j int) bool {
		return upcoming[i].Date.Before(upcoming[j].Date)
	})

	response := domain.APIResponse{
		Success: true,
		Data: struct {
			SpeakingEngagements []domain.SpeakingEngagement `json:"speaking_engagements"`
		}{
			SpeakingEngagements: upcoming,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetEngagementsByEvent returns all engagements for a specific event
func (h *SpeakingHandler) GetEngagementsByEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	event := r.URL.Query().Get("event")
	if event == "" {
		http.Error(w, "Event name is required", http.StatusBadRequest)
		return
	}

	data, err := h.loadSpeakingEngagements()
	if err != nil {
		log.Printf("Error loading speaking engagements: %v", err)
		response := domain.APIResponse{
			Success: false,
			Error:   "Failed to load speaking engagements data",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Filter by event
	filtered := make([]domain.SpeakingEngagement, 0)
	for _, s := range data.SpeakingEngagements {
		if s.Event == event {
			filtered = append(filtered, s)
		}
	}

	response := domain.APIResponse{
		Success: true,
		Data: struct {
			SpeakingEngagements []domain.SpeakingEngagement `json:"speaking_engagements"`
			EventName           string                      `json:"event_name"`
		}{
			SpeakingEngagements: filtered,
			EventName:           event,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// loadSpeakingEngagements reads and unmarshals the speaking engagements JSON file
func (h *SpeakingHandler) loadSpeakingEngagements() (*domain.SpeakingData, error) {
	filePath := filepath.Join(h.dataPath, "speaking.json")

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data domain.SpeakingData
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, err
	}

	return &data, nil
}
