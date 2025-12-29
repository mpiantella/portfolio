package interfaces

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"portfolio/internal/domain"
)

// PatentsHandler handles patent-related HTTP requests
type PatentsHandler struct {
	templates *template.Template
	dataPath  string
}

// NewPatentsHandler creates a new patents handler
func NewPatentsHandler(templates *template.Template, dataPath string) *PatentsHandler {
	return &PatentsHandler{
		templates: templates,
		dataPath:  dataPath,
	}
}

// ServeHTTP handles the patents page request
func (h *PatentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := h.loadPatents()
	if err != nil {
		log.Printf("Error loading patents: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "patents.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// ServeAPI handles the API request for patents data
func (h *PatentsHandler) ServeAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := h.loadPatents()
	if err != nil {
		log.Printf("Error loading patents: %v", err)
		response := domain.APIResponse{
			Success: false,
			Error:   "Failed to load patents data",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Filter by status if requested
	status := r.URL.Query().Get("status")
	if status != "" {
		filtered := make([]domain.Patent, 0)
		for _, p := range data.Patents {
			if p.Status == status {
				filtered = append(filtered, p)
			}
		}
		data.Patents = filtered
	}

	response := domain.APIResponse{
		Success: true,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetPatentByID returns a single patent by ID
func (h *PatentsHandler) GetPatentByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path (assumes route like /api/patents/{id})
	id := r.URL.Path[len("/api/patents/"):]
	if id == "" {
		http.Error(w, "Patent ID is required", http.StatusBadRequest)
		return
	}

	data, err := h.loadPatents()
	if err != nil {
		log.Printf("Error loading patents: %v", err)
		response := domain.APIResponse{
			Success: false,
			Error:   "Failed to load patents data",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Find patent by ID
	var foundPatent *domain.Patent
	for _, p := range data.Patents {
		if p.ID == id {
			foundPatent = &p
			break
		}
	}

	if foundPatent == nil {
		response := domain.APIResponse{
			Success: false,
			Error:   "Patent not found",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := domain.APIResponse{
		Success: true,
		Data:    foundPatent,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// loadPatents reads and unmarshals the patents JSON file
func (h *PatentsHandler) loadPatents() (*domain.PatentsData, error) {
	filePath := filepath.Join(h.dataPath, "patents.json")

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data domain.PatentsData
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, err
	}

	return &data, nil
}

// GetPatentsStats returns statistics about patents
func (h *PatentsHandler) GetPatentsStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := h.loadPatents()
	if err != nil {
		log.Printf("Error loading patents: %v", err)
		response := domain.APIResponse{
			Success: false,
			Error:   "Failed to load patents data",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Calculate statistics
	stats := struct {
		TotalPatents   int `json:"total_patents"`
		GrantedPatents int `json:"granted_patents"`
		PendingPatents int `json:"pending_patents"`
	}{
		TotalPatents: len(data.Patents),
	}

	for _, p := range data.Patents {
		if p.Status == "Granted" {
			stats.GrantedPatents++
		} else if p.Status == "Pending" {
			stats.PendingPatents++
		}
	}

	response := domain.APIResponse{
		Success: true,
		Data:    stats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
