# [Portfolio Website](https://mpiantella.github.io/portfolio/)

A static portfolio website featuring projects, patents, and speaking engagements. Built with Go templates and deployed to GitHub Pages.

## Architecture

This is a **static site generator** that:
- Builds HTML pages from Go templates at build time
- Generates static JSON API files for dynamic content
- Deploys to GitHub Pages for free hosting
- Uses Tailwind CSS for styling

## Project Structure

```
portfolio-website/
├── build.go                        # Static site generator
├── internal/
│   ├── domain/
│   │   └── project.go              # Domain models (Project, Patent, etc.)
│   ├── util/
│   │   └── funcmap.go              # Template helper functions
│   └── infrastructure/data/
│       ├── projects.json           # Project data
│       ├── patents.json            # Patents data
│       └── speaking.json           # Speaking engagements data
├── web/
│   ├── templates/
│   │   ├── layouts/                # Base layouts
│   │   ├── components/             # Reusable components
│   │   └── pages/                  # Page templates
│   ├── static/
│   │   └── dist.css                # Generated Tailwind CSS
│   └── tailwind.css                # Tailwind source
├── dist/                           # Generated static site (gitignored)
│   ├── index.html
│   ├── projects/
│   ├── patents/
│   ├── speaking/
│   ├── contact/
│   ├── api/                        # Static JSON files
│   └── static/
├── .github/workflows/
│   └── deploy-github-pages.yml     # GitHub Pages deployment
├── package.json                    # Node.js dependencies
└── go.mod                          # Go dependencies
```

## Prerequisites

- **Go** 1.21 or higher
- **Node.js** 18 or higher (for Tailwind CSS)
- **npm** (comes with Node.js)

## Quick Start

### Local Development

1. **Install dependencies:**
   ```bash
   npm install
   ```

2. **Build the static site:**
   ```bash
   npm run build
   ```
   This runs Tailwind CSS build and generates static HTML files in `dist/`

3. **Preview locally:**
   ```bash
   npm run preview
   ```
   Visit http://localhost:8080

### Using Make

```bash
# Install dependencies
make install

# Build static site
make build

# Preview locally
make preview

# Clean generated files
make clean
```

## Available Commands

### NPM Scripts

```bash
npm run build:css    # Build minified Tailwind CSS
npm run dev:css      # Watch Tailwind CSS for changes
npm run build        # Build CSS + generate static site
npm run preview      # Serve dist/ locally on port 8080
```

### Make Commands

```bash
make help            # Show all available commands
make install         # Install Node.js dependencies
make build           # Build static site (CSS + HTML)
make preview         # Preview site locally
make clean           # Remove generated files
make dev             # Development mode (CSS watch + preview)
```

## How It Works

### Build Process

1. **Tailwind CSS** compiles utility classes into `web/static/dist.css`
2. **build.go** reads templates and data files
3. **Static HTML** pages are generated for each route
4. **JSON API files** are created in `dist/api/`
5. **Static assets** are copied to `dist/static/`

### Generated Output

```
dist/
├── index.html              # Home page
├── projects/
│   └── index.html         # Projects page
├── patents/
│   └── index.html         # Patents page
├── speaking/
│   └── index.html         # Speaking page
├── contact/
│   └── index.html         # Contact page
├── api/
│   ├── projects.json      # All projects
│   ├── projects-featured.json
│   ├── patents.json
│   └── speaking.json
├── static/
│   ├── dist.css           # Compiled CSS
│   └── images/
└── 404.html               # Error page
```

## GitHub Pages Deployment

### Automatic Deployment

Every push to `main` triggers automatic deployment:

1. GitHub Actions runs the workflow
2. Tailwind CSS is built
3. Static site is generated with `build.go`
4. Site is deployed to GitHub Pages

### Setup Instructions

1. **Enable GitHub Pages:**
   - Go to repository Settings → Pages
   - Under "Build and deployment", set Source to **GitHub Actions**
   - Save

2. **Push to main:**
   ```bash
   git add .
   git commit -m "Deploy to GitHub Pages"
   git push origin main
   ```

3. **Your site will be live at:**
   - `https://USERNAME.github.io/REPOSITORY/`

### Custom Domain (Optional)

1. Add custom domain in Settings → Pages
2. Configure DNS at your registrar:
   ```
   Type: CNAME
   Name: www (or subdomain)
   Value: USERNAME.github.io
   ```
3. Enable "Enforce HTTPS"

## Content Management

### Adding/Updating Content

Content is stored in JSON files under `internal/infrastructure/data/`:

- `projects.json` - Portfolio projects
- `patents.json` - Patents and innovations
- `speaking.json` - Speaking engagements

After updating JSON files:
```bash
npm run build        # Rebuild static site
git commit -am "Update content"
git push origin main # Auto-deploys to GitHub Pages
```

### Adding New Pages

1. Create template in `web/templates/pages/`
2. Add generation function in `build.go`
3. Call function in `main()`
4. Rebuild and deploy

## Technologies Used

- **Static Generator**: Go 1.21+
- **Templates**: Go `html/template`
- **Styling**: Tailwind CSS 3.4+
- **Hosting**: GitHub Pages (free)
- **CI/CD**: GitHub Actions

## GitHub Actions Workflows

### 1. Deploy to GitHub Pages
- **File**: `.github/workflows/deploy-github-pages.yml`
- **Trigger**: Push to main
- **Actions**: Build CSS → Generate static site → Deploy to GitHub Pages

### 2. Release Binaries
- **File**: `.github/workflows/release.yml`
- **Trigger**: Version tags (e.g., `v1.0.0`)
- **Actions**: Build cross-platform binaries for self-hosting

## Deployment Options

This portfolio supports multiple deployment strategies:

### GitHub Pages (Current - Recommended)
- ✅ **Free** hosting with SSL
- ✅ **Automatic** deployment on push
- ✅ **Fast** CDN delivery
- ✅ **Simple** - no server management
- 📖 [Setup Guide](./docs/deployment/github-pages.md)

### Self-Hosting (Alternative)
Download pre-built binaries from [Releases](https://github.com/OWNER/REPO/releases) for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

Run the server binary:
```bash
./portfolio-linux-amd64
```

## Development Workflow

### Making Changes

1. **Update content**: Edit JSON files in `internal/infrastructure/data/`
2. **Update templates**: Modify HTML in `web/templates/`
3. **Update styles**: Add Tailwind classes or edit `web/tailwind.css`
4. **Build**: `npm run build`
5. **Preview**: `npm run preview`
6. **Deploy**: Push to main

### Development Mode

For live reloading during development:

```bash
# Terminal 1: Watch CSS
npm run dev:css

# Terminal 2: Rebuild on changes
# (Manual rebuild required after template/data changes)
npm run build && npm run preview
```

Or use Make:
```bash
make dev
```

## Cost & Performance

| Metric | Value |
|--------|-------|
| **Hosting Cost** | $0 (GitHub Pages free tier) |
| **SSL Certificate** | Free (included) |
| **Bandwidth** | 100 GB/month (free) |
| **Build Time** | ~10-15 seconds |
| **Deploy Time** | ~30 seconds |
| **Page Load** | <500ms (global CDN) |

## Troubleshooting

### Build fails locally

```bash
# Clear and rebuild
make clean
npm install
npm run build
```

### CSS not updating

```bash
# Rebuild Tailwind CSS
npm run build:css
```

### Site not deploying

1. Check GitHub Actions tab for errors
2. Verify GitHub Pages is enabled (Settings → Pages)
3. Ensure workflow has write permissions (Settings → Actions → General)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes and test locally
4. Push and create a pull request

## License

This project structure is open source. Content and personal information are © Maria Lucena.

---

**Live Site**: https://USERNAME.github.io/REPOSITORY/
**Source Code**: https://github.com/OWNER/REPO
