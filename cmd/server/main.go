package main

import (
	"fmt"
	"html/template"
	"net/http"

	"portfolio/internal/infrastructure/config"
	"portfolio/internal/infrastructure/logging"
	"portfolio/internal/infrastructure/server"
)

func main() {
	cfg := config.Load()
	logger := logging.New(cfg.LogLevel)
	_ = logger

	// build address from port
	addr := fmt.Sprintf(":%s", cfg.Port)
	logger.Info().Str("addr", addr).Msg("starting server")

	templates := template.Must(template.ParseGlob("web/templates/**/*.html"))
	r := server.NewRouter(templates)

	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Fatal().Err(err).Msg("server failed")
	}
}
