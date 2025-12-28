package config

import "os"

// Config holds environment-driven configuration values.
type Config struct {
	Port     string // e.g. "8080"
	LogLevel string // e.g. "debug", "info"
}

// Load loads configuration from environment with sensible defaults.
func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	lvl := os.Getenv("LOG_LEVEL")
	if lvl == "" {
		lvl = "info"
	}
	return Config{Port: port, LogLevel: lvl}
}
