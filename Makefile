# Makefile for Portfolio Website - Static Site Generator
# Usage: `make <target>`, e.g. `make install`, `make build`

.PHONY: help install build build-css dev preview clean stop-dev test fmt

start:
	$(MAKE) install && $(MAKE) build && $(MAKE) preview
help:
	@echo "Portfolio Website - Static Site Generator"
	@echo ""
	@echo "Available commands:"
	@echo "  help        - Show this help message"
	@echo "  install     - Install Node.js dependencies"
	@echo "  build       - Build static site (CSS + HTML)"
	@echo "  build-css   - Build Tailwind CSS only"
	@echo "  preview     - Preview static site locally (port 8080)"
	@echo "  dev         - Development mode (watch CSS + preview)"
	@echo "  test        - Run Go tests"
	@echo "  fmt         - Format Go code"
	@echo "  clean       - Remove generated files (dist/, CSS)"
	@echo "  stop-dev    - Stop Tailwind watcher"
	@echo ""
	@echo "Quick start:"
	@echo "  make install  # Install dependencies"
	@echo "  make build    # Build static site"
	@echo "  make preview  # Preview locally"

install:
	@echo "📦 Installing Node.js dependencies..."
	npm install

build: build-css
	@echo "🔨 Building static site..."
	go run build.go
	@echo "✅ Static site generated in dist/"

build-css:
	@echo "🎨 Building Tailwind CSS..."
	npm run build:css

preview:
	@echo "🌐 Starting preview server on http://localhost:8080"
	@echo "   Press Ctrl+C to stop"
	cd dist && python3 -m http.server 8080

dev:
	@echo "🚀 Starting development mode..."
	@echo "   - Tailwind CSS watcher (background)"
	@echo "   - Preview server on http://localhost:8080"
	@echo ""
	@echo "   Make changes to templates/CSS, then run 'make build' in another terminal"
	@echo "   Press Ctrl+C to stop preview, then run 'make stop-dev' to stop CSS watcher"
	@echo ""
	npm run dev:css & \
	sleep 2; \
	$(MAKE) preview

test:
	@echo "🧪 Running Go tests..."
	go test ./...

fmt:
	@echo "✨ Formatting Go code..."
	go fmt ./...

clean:
	@echo "🧹 Cleaning generated files..."
	rm -rf dist/
	rm -f web/static/dist.css
	@echo "✅ Cleaned dist/ and dist.css"

stop-dev:
	@echo "🛑 Stopping Tailwind CSS watcher..."
	-pkill -f "tailwindcss -i ./web/tailwind.css" || true
	@echo "✅ Stopped CSS watcher"
