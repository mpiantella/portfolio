# Makefile — lightweight bootstrap helpers for the portfolio app
# Usage: `make <target>`, e.g. `make install`, `make dev`

.PHONY: install build-css dev-css build test run start dev fmt clean stop-dev

install:
	npm install

build-css:
	npm run build:css

dev-css:
	npm run dev:css

build: build-css

test:
	go test ./...

run:
	go run ./cmd/server

start: build run

# dev: start Tailwind watcher in background, then run Go server in foreground.
# Note: tailwind watcher will continue running in background; stop it with `make stop-dev`
dev:
	@echo "Starting Tailwind watcher in background..."
	npm run dev:css & \
	sleep 1; \
	go run ./cmd/server

fmt:
	go fmt ./...

clean:
	rm -f web/static/dist.css

stop-dev:
	# attempts to stop the Tailwind watcher started by `make dev`
	-pkill -f "tailwindcss -i ./web/tailwind.css" || true