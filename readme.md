# Portfolio Website

[![CI](https://github.com/OWNER/REPO/workflows/CI%20-%20Test%20%26%20Build/badge.svg)](https://github.com/OWNER/REPO/actions/workflows/ci.yml)
[![Docker](https://github.com/OWNER/REPO/workflows/Publish%20Docker%20Image/badge.svg)](https://github.com/OWNER/REPO/actions/workflows/docker-publish.yml)
[![Release](https://github.com/OWNER/REPO/workflows/Release/badge.svg)](https://github.com/OWNER/REPO/actions/workflows/release.yml)

A lightweight Go-based portfolio website featuring projects, patents, and speaking engagements.

**Simplified bootstrapping with GitHub Actions** - No local build tools required! Just fork, push, and deploy using pre-built Docker images or binaries.

## Project Structure

```
portfolio-website/
├── cmd/server/
│   └── main.go                         # Application entry point
├── internal/
│   ├── domain/
│   │   └── project.go                  # Domain models
│   ├── interfaces/http/
│   │   ├── home_handler.go             # Home page handler
│   │   ├── projects_handler.go         # Projects handler
│   │   ├── patents_handler.go          # Patents handler
│   │   ├── speaking_handler.go         # Speaking handler
│   │   └── contact_handler.go          # Contact form handler
│   ├── infrastructure/
│   │   ├── config/
│   │   │   └── config.go               # Configuration loader
│   │   ├── logging/
│   │   │   └── logger.go               # Zerolog setup
│   │   ├── server/
│   │   │   └── router.go               # HTTP router setup
│   │   └── data/
│   │       ├── projects.json           # Project data
│   │       ├── patents.json            # Patents data
│   │       └── speaking.json           # Speaking engagements data
│   ├── usecase/
│   │   └── README.md                   # Use case layer docs
│   └── util/
│       └── funcmap.go                  # Template helper functions
├── web/
│   ├── templates/
│   │   ├── layouts/                    # Base layouts
│   │   ├── components/                 # Reusable components
│   │   └── pages/                      # Page templates
│   ├── static/
│   │   ├── dist.css                    # Generated Tailwind CSS (gitignored)
│   │   └── styles.css                  # Additional custom styles
│   └── tailwind.css                    # Tailwind source
├── Dockerfile                          # Multi-stage production build
├── Makefile                            # Development commands
├── go.mod                              # Go dependencies
└── package.json                        # Node.js dependencies
```

## Prerequisites

- **Go** 1.20 or higher
- **Node.js** 18 or higher (for Tailwind CSS)
- **npm** (comes with Node.js)

## Quick Start

### Local Development

#### Option 1: Using Make (Recommended)

```bash
# Install Node.js dependencies
make install

# Start development mode (Tailwind watcher + Go server)
make dev
```

This starts both the Tailwind CSS watcher and the Go server. The server will be available at `http://localhost:8080` (or the port specified in the `PORT` environment variable).

To stop the background Tailwind watcher:
```bash
make stop-dev
```

#### Option 2: Manual Setup

1. **Install dependencies:**
   ```bash
   npm install
   ```

2. **Build CSS assets:**
   ```bash
   npm run build:css
   ```

3. **Run the server:**
   ```bash
   go run ./cmd/server
   ```

4. **For development with live CSS rebuild:**
   ```bash
   # Terminal 1: Start Tailwind in watch mode
   npm run dev:css

   # Terminal 2: Start Go server
   go run ./cmd/server
   ```

### Using Pre-built Docker Images

The easiest way to run the application is using the pre-built Docker images from GitHub Container Registry:

```bash
# Pull and run the latest image
docker run -p 8080:8080 ghcr.io/OWNER/REPO:latest

# Or run a specific version
docker run -p 8080:8080 ghcr.io/OWNER/REPO:v1.0.0
```

Replace `OWNER/REPO` with your GitHub repository path (e.g., `username/portfolio-website`).

## Available Make Commands

```bash
make help          # Show all available commands
make install       # Install Node.js dependencies
make dev           # Start Tailwind watcher + Go server
make build         # Build production assets and binary
make build-css     # Build minified CSS for production
make build-bin     # Build Go binary to ./bin/server
make run           # Run Go server without rebuilding
make start         # Build and run production binary
make test          # Run Go tests
make fmt           # Format Go code
make clean         # Remove generated files (dist.css, binaries)
make stop-dev      # Stop Tailwind watcher
make docker-build  # Build Docker image
make docker-start  # Build and run Docker container
```

## Environment Variables

The application reads the following environment variables:

- `PORT` - Server port (default: `8080`)
- `LOG_LEVEL` - Logging level: `debug`, `info`, `warn`, `error` (default: `info`)

Example:
```bash
PORT=3000 LOG_LEVEL=debug go run ./cmd/server
```

## API Endpoints

### Health Check
- `GET /api/health` - Returns `{"status":"ok"}`

### Projects
- `GET /projects` - Projects page
- `GET /api/projects` - JSON list of all projects
- `GET /api/projects?featured=true` - JSON list of featured projects only

### Patents
- `GET /patents` - Patents page
- `GET /api/patents` - JSON list of all patents
- `GET /api/patents/stats` - Patent statistics

### Speaking
- `GET /speaking` - Speaking engagements page
- `GET /api/speaking` - JSON list of all speaking engagements
- `GET /api/speaking/stats` - Speaking statistics
- `GET /api/speaking/upcoming` - Upcoming engagements

### Contact
- `GET /contact` - Contact page
- `POST /api/contact` - Submit contact form (HTMX endpoint)
- `GET /api/contact/stats` - Contact form statistics

## Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Using Make
make test
```

## Production Build

### Using Make

```bash
# Build everything (CSS + Go binary)
make build

# Start production server
make start
```

The binary will be created at `./bin/server`.

### Using Docker

```bash
# Build Docker image
make docker-build

# Or manually:
docker build -t portfolio:latest .

# Run container
docker run -p 8080:8080 portfolio:latest

# Or using Make:
make docker-start
```

The Dockerfile uses a multi-stage build:
1. **Assets stage**: Builds Tailwind CSS
2. **Builder stage**: Compiles Go binary
3. **Final stage**: Minimal Alpine image with binary and assets

## Development Workflow

1. **Adding a new page:**
   - Create a handler in `internal/interfaces/http/`
   - Register the route in `internal/infrastructure/server/router.go`
   - Create the template in `web/templates/pages/`
   - Add tests in `internal/interfaces/http/`

2. **Adding data:**
   - Add JSON data files to `internal/infrastructure/data/`
   - Update domain models in `internal/domain/` if needed
   - Reference the data in your handler

3. **Styling:**
   - Use Tailwind utility classes in templates
   - Custom CSS goes in `web/static/styles.css`
   - Source Tailwind directives in `web/tailwind.css`

## Technologies Used

- **Backend**: Go 1.20+
- **Logging**: [zerolog](https://github.com/rs/zerolog)
- **Frontend**: HTML templates with HTMX
- **Styling**: Tailwind CSS 3.4+
- **Deployment**: Docker, multi-stage builds

## GitHub Actions CI/CD

This project includes automated workflows for continuous integration and deployment:

### Available Workflows

1. **CI - Test & Build** ([`.github/workflows/ci.yml`](.github/workflows/ci.yml))
   - Runs on: Pull requests and pushes to main
   - **Jobs:**
     - Runs Go tests with coverage
     - Checks code formatting (`go fmt`)
     - Builds Tailwind CSS assets
     - Compiles Go binary for Linux
     - Builds Docker image (main branch only)
   - All jobs run in parallel for fast feedback

2. **Publish Docker Image** ([`.github/workflows/docker-publish.yml`](.github/workflows/docker-publish.yml))
   - Runs on: Pushes to main and version tags
   - **Publishes to:** GitHub Container Registry (ghcr.io)
   - **Tags generated:**
     - `latest` (for main branch)
     - `main-<sha>` (for main branch commits)
     - `v1.0.0`, `v1.0`, `v1` (for version tags)
   - Multi-platform support: `linux/amd64`, `linux/arm64`

3. **Release** ([`.github/workflows/release.yml`](.github/workflows/release.yml))
   - Runs on: Version tags (e.g., `v1.0.0`)
   - **Builds binaries for:**
     - Linux (amd64, arm64)
     - macOS (amd64, arm64)
     - Windows (amd64)
   - Creates GitHub release with binaries and checksums

4. **Deploy to AWS** ([`.github/workflows/deploy-to-aws.yml`](.github/workflows/deploy-to-aws.yml))
   - Runs on: Pushes to main (manual trigger available)
   - Builds and pushes to Amazon ECR
   - Optionally deploys to ECS (if configured)

### Setting Up GitHub Actions

#### For Docker Publishing (GitHub Container Registry)

No secrets required! The workflow uses `GITHUB_TOKEN` automatically.

1. Ensure GitHub Actions has permission to write packages:
   - Go to repository Settings → Actions → General
   - Under "Workflow permissions", select "Read and write permissions"

2. Push to main or create a tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

3. Images will be available at: `ghcr.io/YOUR_USERNAME/portfolio-website`

#### For AWS Deployment

Add these secrets to your repository (Settings → Secrets and variables → Actions):

| Secret | Description |
|--------|-------------|
| `AWS_ACCESS_KEY_ID` | AWS access key with ECR/ECS permissions |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key |
| `AWS_REGION` | AWS region (e.g., `us-east-1`) |
| `ECR_REPOSITORY` | ECR repository name |
| `ECS_CLUSTER` | (Optional) ECS cluster name for auto-deployment |
| `ECS_SERVICE` | (Optional) ECS service name for auto-deployment |

### Creating a Release

```bash
# Create and push a version tag
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

This will:
1. Trigger the Release workflow to build binaries
2. Trigger the Docker publish workflow
3. Create a GitHub release with downloadable binaries

### Simplified Bootstrapping with GitHub Actions

**For new contributors or deployments:**

Instead of installing Go and Node.js locally, you can:

1. **Fork the repository**
2. **Enable GitHub Actions** (they run automatically)
3. **Download pre-built artifacts:**
   - Docker images from GitHub Container Registry
   - Binaries from GitHub Releases
   - No local build tools required!

**For deployments:**

```bash
# Option 1: Use pre-built Docker image (no build required)
docker pull ghcr.io/OWNER/REPO:latest
docker run -p 8080:8080 ghcr.io/OWNER/REPO:latest

# Option 2: Download binary from releases
wget https://github.com/OWNER/REPO/releases/download/v1.0.0/portfolio-linux-amd64
chmod +x portfolio-linux-amd64
./portfolio-linux-amd64
```

## Notes

- Templates use standard Go `html/template` with HTMX for dynamic interactions
- HTMX is loaded via CDN in the base template
- Tailwind CSS is compiled from source (not CDN) for production optimization
- Static files are served from `./web/static`
- All JSON data files are loaded from `./internal/infrastructure/data`
- GitHub Actions automatically builds and publishes Docker images on every push to main
