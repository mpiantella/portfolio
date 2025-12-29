# Makefile — lightweight bootstrap helpers for the portfolio app
# Usage: `make <target>`, e.g. `make install`, `make dev`

.PHONY: help install build-css dev-css build test run start dev fmt clean stop-dev stop-all

help:
	@echo "Makefile commands:"
	@echo "  help        - Show this help"
	@echo "  install     - Install Node dev deps"
	@echo "  build-css   - Build Tailwind CSS to web/static/dist.css"
	@echo "  dev-css     - Run Tailwind in watch mode"
	@echo "  build       - Build CSS (alias)"
	@echo "  test        - Run go test ./..."
	@echo "  run         - Run the Go server"
	@echo "  dev         - Start Tailwind watcher in background and run server"
	@echo "  fmt         - Run go fmt"
	@echo "  clean       - Remove generated CSS"
	@echo "  stop-dev    - Stop the Tailwind watcher"

install:
	npm install


dev-css:
	npm run dev:css

# Build both frontend assets and server binary for production
build: build-css build-bin

build-css:
	npm run build:css

build-bin:
	@echo "Building server binary into ./bin/server"
	mkdir -p bin
	GOOS=${GOOS:-$(shell go env GOOS)} GOARCH=${GOARCH:-$(shell go env GOARCH)} go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

start: build
	@echo "Starting server (production) ./bin/server"
	./bin/server

# dev: start Tailwind watcher in background, then run Go server in foreground.
# Note: tailwind watcher will continue running in background; stop it with `make stop-dev`
dev:
	@echo "Starting Tailwind watcher in background..."
	npm run dev:css & \
	sleep 1; \
	go run ./cmd/server

fmt:
	go fmt ./...

test:
	go test ./...

clean:
	rm -f web/static/dist.css bin/server

stop-dev:
	# attempts to stop the Tailwind watcher started by `make dev`
	-pkill -f "tailwindcss -i ./web/tailwind.css" || true

# Docker helpers
docker-build:
	docker build -t portfolio:latest .

# Build and run the container locally
docker-start: docker-build
	@echo "Starting Docker container on port 8080 (maps to host 8080)"
	docker run --rm -p 8080:8080 portfolio:latest

# Stop both frontend (Tailwind watcher) and backend (Go server) processes started by `make dev`
stop-all:
	@echo "Stopping Tailwind watcher (frontend)..."
	-pkill -f "tailwindcss -i ./web/tailwind.css" || true
	@echo "Stopping Go server (backend)..."
	-pkill -f "./cmd/server" || pkill -f "/app/server" || true