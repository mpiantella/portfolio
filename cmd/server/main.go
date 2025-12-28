package main

import (
	"fmt"
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
	
	r := server.NewRouter()
	
	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Fatal().Err(err).Msg("server failed")
	}
}
