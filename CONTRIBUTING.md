# Contributing to Portfolio Website

Thank you for your interest in contributing! This guide will help you get started.

## Quick Start for Contributors

You have two options for getting started:

### Option 1: Simplified Bootstrapping (No Local Build Tools)

Perfect for quick contributions or testing:

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/YOUR_USERNAME/portfolio-website.git
   cd portfolio-website
   ```

2. **Use pre-built Docker image from GitHub Container Registry**
   ```bash
   # Pull the latest image built by GitHub Actions
   docker pull ghcr.io/ORIGINAL_OWNER/portfolio-website:latest

   # Run it locally
   docker run -p 8080:8080 ghcr.io/ORIGINAL_OWNER/portfolio-website:latest
   ```

3. **Make your changes** to code files

4. **Push and create a PR**
   - GitHub Actions will automatically build and test your changes
   - No need to build locally!

### Option 2: Full Local Development

For active development with live reloading:

1. **Prerequisites:**
   - Go 1.20+
   - Node.js 18+
   - Make (optional but recommended)

2. **Setup:**
   ```bash
   # Install dependencies
   make install

   # Start development mode (auto-reload CSS)
   make dev
   ```

3. **Make changes and test**
   ```bash
   # Run tests
   make test

   # Format code
   make fmt
   ```

## Development Workflow

### Making Changes

1. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes:**
   - Code changes go in `internal/`, `cmd/`, or `web/`
   - Update tests in `*_test.go` files
   - Update documentation if needed

3. **Test your changes:**
   ```bash
   # Run tests locally
   make test

   # Or let GitHub Actions test for you by pushing
   git push origin feature/your-feature-name
   ```

4. **Create a Pull Request:**
   - GitHub Actions will automatically:
     - Run all tests
     - Check code formatting
     - Build Docker image
     - Report results on your PR

### What GitHub Actions Does For You

Every time you push code, GitHub Actions automatically:

✅ Runs all Go tests with coverage
✅ Checks code formatting (`go fmt`)
✅ Builds Tailwind CSS assets
✅ Compiles the Go binary
✅ Builds Docker image (for main branch)
✅ Reports all results on your PR

**You don't need to build anything locally!** Just push your changes and the CI will validate them.

## Project Structure

```
portfolio-website/
├── cmd/server/              # Application entry point
├── internal/
│   ├── domain/              # Domain models and business logic
│   ├── interfaces/http/     # HTTP handlers
│   ├── infrastructure/      # Config, logging, data, server setup
│   └── util/                # Helper utilities
├── web/
│   ├── templates/           # HTML templates
│   └── static/              # Static assets
├── .github/workflows/       # CI/CD automation
└── Makefile                 # Development commands
```

## Common Tasks

### Adding a New Page

1. **Create handler** in `internal/interfaces/http/`
2. **Register route** in `internal/infrastructure/server/router.go`
3. **Create template** in `web/templates/pages/`
4. **Add tests** in `internal/interfaces/http/*_test.go`

Example:
```bash
# 1. Create handler
touch internal/interfaces/http/blog_handler.go

# 2. Edit router to add route
# (Edit internal/infrastructure/server/router.go)

# 3. Create template
touch web/templates/pages/blog.html

# 4. Add tests
touch internal/interfaces/http/blog_handler_test.go
```

### Adding Data

1. Add JSON file to `internal/infrastructure/data/`
2. Update domain models in `internal/domain/` if needed
3. Reference in your handler

### Updating Styles

1. **Edit templates** with Tailwind utility classes
2. **Custom CSS** goes in `web/static/styles.css`
3. **Tailwind config** is in `tailwind.config.js`

If running locally with `make dev`, CSS rebuilds automatically.

## Testing

### Local Testing
```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -v ./internal/interfaces/http -run TestHomeHandler
```

### CI Testing
Just push your branch - GitHub Actions runs all tests automatically!

## Code Style

- Follow standard Go formatting: `make fmt` or `go fmt ./...`
- Write tests for new features
- Keep handlers focused and simple
- Use descriptive variable names

The CI will check formatting automatically.

## Pull Request Process

1. **Fork** the repository
2. **Create a branch** from `main`
3. **Make your changes**
4. **Push to your fork**
5. **Create a Pull Request** to the main repository
6. **Wait for CI** to complete (automatic)
7. **Address review feedback** if any
8. **Merge** once approved!

### PR Checklist

- [ ] Code follows Go conventions (`go fmt`)
- [ ] Tests added/updated for new features
- [ ] All tests pass (checked by CI)
- [ ] Documentation updated if needed
- [ ] No sensitive data in commits

## Getting Help

- **Questions?** Open a GitHub Discussion
- **Bugs?** Open a GitHub Issue
- **Feature ideas?** Open a GitHub Issue with the "enhancement" label

## Using Pre-built Artifacts

### Docker Images

All commits to `main` automatically build and publish Docker images:

```bash
# Latest from main branch
docker pull ghcr.io/OWNER/REPO:latest

# Specific commit
docker pull ghcr.io/OWNER/REPO:main-abc123

# Tagged release
docker pull ghcr.io/OWNER/REPO:v1.0.0
```

### Release Binaries

Download pre-built binaries from [GitHub Releases](https://github.com/OWNER/REPO/releases):

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

No compilation required!

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.
