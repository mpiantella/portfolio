package server

import (
	"html/template"
	"net/http"

	httpHandlers "portfolio/internal/interfaces/http"
)

func NewRouter(templates *template.Template) *http.ServeMux {
	mux := http.NewServeMux()

	// Static files
	fs := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Page handlers
	homeHandler := httpHandlers.NewHomeHandler(templates)
	projectsHandler := httpHandlers.NewProjectsHandler(templates, "./internal/infrastructure/data")
	patentsHandler := httpHandlers.NewPatentsHandler(templates, "./internal/infrastructure/data")
	speakingHandler := httpHandlers.NewSpeakingHandler(templates, "./internal/infrastructure/data")
	contactHandler := httpHandlers.NewContactHandler(templates)

	// Routes
	mux.Handle("/", homeHandler)
	mux.HandleFunc("/api/health", healthHandler)

	// Projects API routes
	mux.Handle("/projects", projectsHandler)
	mux.HandleFunc("/api/projects", projectsHandler.ServeAPI)

	// Patents API routes
	mux.Handle("/patents", patentsHandler)
	mux.HandleFunc("/api/patents", patentsHandler.ServeAPI)
	mux.HandleFunc("/api/patents/stats", patentsHandler.GetPatentsStats)

	// Speaking API routes
	mux.Handle("/speaking", speakingHandler)
	mux.HandleFunc("/api/speaking", speakingHandler.ServeAPI)
	mux.HandleFunc("/api/speaking/stats", speakingHandler.GetSpeakingStats)
	mux.HandleFunc("/api/speaking/upcoming", speakingHandler.GetUpcomingEngagements)

	// Contact API routes
	mux.Handle("/contact", contactHandler)
	mux.HandleFunc("/api/contact", contactHandler.ServeAPI)
	mux.HandleFunc("/api/contact/stats", contactHandler.GetContactStats)

	return mux
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}
