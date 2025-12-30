# GitHub Pages Deployment (Static Site)

Deploy your portfolio website as a static site using GitHub Pages - completely free hosting with HTTPS and custom domain support.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│               GitHub Pages Static Site Architecture                 │
└─────────────────────────────────────────────────────────────────────┘

                         ┌──────────────────┐
                         │  Developer       │
                         │  Pushes Code     │
                         └────────┬─────────┘
                                  │
                                  ▼
                         ┌──────────────────┐
                         │  GitHub Actions  │
                         │  (CI/CD)         │
                         └────────┬─────────┘
                                  │
                    ┌─────────────┴─────────────────┐
                    │  1. Build Static Site         │
                    │     - Generate HTML pages     │
                    │     - Compile CSS             │
                    │     - Copy static assets      │
                    │     - Create JSON API files   │
                    └─────────────┬─────────────────┘
                                  │
                                  ▼
                    ┌──────────────────────────────┐
                    │   GitHub Repository          │
                    │   (gh-pages branch)          │
                    │                              │
                    │   ├── index.html             │
                    │   ├── projects/              │
                    │   │   └── index.html         │
                    │   ├── speaking/              │
                    │   ├── patents/               │
                    │   ├── contact/               │
                    │   ├── api/                   │
                    │   │   ├── projects.json      │
                    │   │   ├── speaking.json      │
                    │   │   └── patents.json       │
                    │   ├── static/                │
                    │   │   ├── dist.css           │
                    │   │   └── images/            │
                    │   └── 404.html               │
                    └──────────────┬───────────────┘
                                   │
                                   ▼
                    ┌──────────────────────────────┐
                    │   GitHub Pages CDN           │
                    │   (Global Distribution)      │
                    └──────────────┬───────────────┘
                                   │
                                   ▼
                    ┌──────────────────────────────┐
                    │  HTTPS URL                   │
                    │  username.github.io/repo     │
                    │  OR                          │
                    │  custom-domain.com           │
                    └──────────────┬───────────────┘
                                   │
                                   ▼
                         ┌──────────────────┐
                         │  End Users       │
                         │  (Static Files)  │
                         └──────────────────┘
```

## ⚠️ Important: Architecture Changes Required

**This deployment requires converting your Go application to a static site generator.**

### What Needs to Change

| Current (Dynamic) | Changed To (Static) |
|------------------|---------------------|
| Go server renders templates | Build script generates HTML files |
| HTMX loads data from API | JavaScript fetches static JSON files |
| Server-side routing | Client-side routing or separate HTML files |
| Contact form POST handler | External form service (Formspree, etc.) |
| Real-time data loading | Pre-generated content |

## Resources Required

### GitHub Resources (All Free!)
1. **GitHub Repository** - Host code and static files
2. **GitHub Pages** - Free hosting (100 GB bandwidth/month)
3. **GitHub Actions** - Free CI/CD (2,000 minutes/month)

### Optional External Services
- **Formspree** - Contact form handling (free tier: 50 submissions/month)
- **Cloudflare** - CDN and DDoS protection (free)

## Cost Breakdown

### Monthly Cost

| Component | Cost |
|-----------|------|
| GitHub Pages Hosting | **$0** (Free) |
| GitHub Actions CI/CD | **$0** (Free tier) |
| SSL Certificate | **$0** (Included) |
| CDN / Bandwidth | **$0** (100 GB included) |
| Custom Domain (optional) | $12/year (~$1/mo) |
| **Total** | **$0 - $1/month** |

### Limits (Free Tier)
- **Bandwidth**: 100 GB/month
- **Build time**: 2,000 minutes/month
- **Storage**: 1 GB
- **Build size**: 1 GB

**Perfect for**: Personal portfolios, documentation sites, low-medium traffic

## Step-by-Step Conversion

### Phase 1: Create Static Site Generator

Create `build.go` in project root:

```go
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"portfolio/internal/domain"
	"portfolio/internal/util"
)

type PageData struct {
	Title       string
	Projects    []domain.Project
	Patents     []domain.Patent
	Speaking    []domain.SpeakingEngagement
	CurrentPage string
}

func main() {
	// Create output directory
	outputDir := "dist"
	os.RemoveAll(outputDir)
	os.MkdirAll(outputDir, 0755)

	// Load templates
	templates := template.Must(
		template.New("").Funcs(util.FuncMap()).ParseGlob("web/templates/layouts/*.html"),
	)
	template.Must(templates.ParseGlob("web/templates/components/*.html"))
	template.Must(templates.ParseGlob("web/templates/pages/*.html"))

	// Load data
	projects := loadProjects()
	patents := loadPatents()
	speaking := loadSpeaking()

	// Generate pages
	generateHomePage(templates, outputDir, projects, patents, speaking)
	generateProjectsPage(templates, outputDir, projects)
	generatePatentsPage(templates, outputDir, patents)
	generateSpeakingPage(templates, outputDir, speaking)
	generateContactPage(templates, outputDir)

	// Copy static assets
	copyDir("web/static", filepath.Join(outputDir, "static"))

	// Generate API JSON files
	generateAPIFiles(outputDir, projects, patents, speaking)

	fmt.Println("✅ Static site generated in dist/")
}

func generateHomePage(tmpl *template.Template, outDir string, projects []domain.Project, patents []domain.Patent, speaking []domain.SpeakingEngagement) {
	data := PageData{
		Title:       "Maria Lucena - Director of Architecture",
		Projects:    projects,
		Patents:     patents,
		Speaking:    speaking,
		CurrentPage: "home",
	}

	f, err := os.Create(filepath.Join(outDir, "index.html"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := tmpl.ExecuteTemplate(f, "home.html", data); err != nil {
		panic(err)
	}
}

func generateProjectsPage(tmpl *template.Template, outDir string, projects []domain.Project) {
	os.MkdirAll(filepath.Join(outDir, "projects"), 0755)

	data := struct {
		Projects []domain.Project
	}{Projects: projects}

	f, err := os.Create(filepath.Join(outDir, "projects", "index.html"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := tmpl.ExecuteTemplate(f, "projects.html", data); err != nil {
		panic(err)
	}
}

func generatePatentsPage(tmpl *template.Template, outDir string, patents []domain.Patent) {
	os.MkdirAll(filepath.Join(outDir, "patents"), 0755)

	data := struct {
		Patents []domain.Patent
	}{Patents: patents}

	f, err := os.Create(filepath.Join(outDir, "patents", "index.html"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := tmpl.ExecuteTemplate(f, "patents.html", data); err != nil {
		panic(err)
	}
}

func generateSpeakingPage(tmpl *template.Template, outDir string, speaking []domain.SpeakingEngagement) {
	os.MkdirAll(filepath.Join(outDir, "speaking"), 0755)

	data := struct {
		SpeakingEngagements []domain.SpeakingEngagement
	}{SpeakingEngagements: speaking}

	f, err := os.Create(filepath.Join(outDir, "speaking", "index.html"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := tmpl.ExecuteTemplate(f, "speaking.html", data); err != nil {
		panic(err)
	}
}

func generateContactPage(tmpl *template.Template, outDir string) {
	os.MkdirAll(filepath.Join(outDir, "contact"), 0755)

	f, err := os.Create(filepath.Join(outDir, "contact", "index.html"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := tmpl.ExecuteTemplate(f, "contact.html", nil); err != nil {
		panic(err)
	}
}

func generateAPIFiles(outDir string, projects []domain.Project, patents []domain.Patent, speaking []domain.SpeakingEngagement) {
	apiDir := filepath.Join(outDir, "api")
	os.MkdirAll(apiDir, 0755)

	// Projects API
	projectsData := struct {
		Projects []domain.Project `json:"projects"`
	}{Projects: projects}

	writeJSON(filepath.Join(apiDir, "projects.json"), projectsData)

	// Featured projects
	featured := []domain.Project{}
	for _, p := range projects {
		if p.Featured {
			featured = append(featured, p)
		}
	}
	writeJSON(filepath.Join(apiDir, "projects-featured.json"), struct {
		Projects []domain.Project `json:"projects"`
	}{Projects: featured})

	// Patents API
	writeJSON(filepath.Join(apiDir, "patents.json"), struct {
		Patents []domain.Patent `json:"patents"`
	}{Patents: patents})

	// Speaking API
	writeJSON(filepath.Join(apiDir, "speaking.json"), struct {
		SpeakingEngagements []domain.SpeakingEngagement `json:"speaking_engagements"`
	}{SpeakingEngagements: speaking})
}

func loadProjects() []domain.Project {
	data, _ := ioutil.ReadFile("internal/infrastructure/data/projects.json")
	var result struct {
		Projects []domain.Project `json:"projects"`
	}
	json.Unmarshal(data, &result)
	return result.Projects
}

func loadPatents() []domain.Patent {
	data, _ := ioutil.ReadFile("internal/infrastructure/data/patents.json")
	var result struct {
		Patents []domain.Patent `json:"patents"`
	}
	json.Unmarshal(data, &result)
	return result.Patents
}

func loadSpeaking() []domain.SpeakingEngagement {
	data, _ := ioutil.ReadFile("internal/infrastructure/data/speaking.json")
	var result struct {
		SpeakingEngagements []domain.SpeakingEngagement `json:"speaking_engagements"`
	}
	json.Unmarshal(data, &result)
	return result.SpeakingEngagements
}

func writeJSON(path string, data interface{}) {
	file, _ := os.Create(path)
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(data)
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(src, path)
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		return ioutil.WriteFile(dstPath, data, info.Mode())
	})
}
```

### Phase 2: Update Templates for Static Paths

Update navigation links to use relative paths:

**web/templates/components/nav.html:**
```html
{{define "nav"}}
<nav class="bg-white shadow-lg sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
            <div class="flex items-center">
                <a href="/" class="text-2xl font-bold text-blue-600">ML</a>
            </div>

            <div class="hidden md:flex items-center space-x-8">
                <a href="/" class="text-gray-700 hover:text-blue-600 transition-colors">Home</a>
                <a href="/projects/" class="text-gray-700 hover:text-blue-600 transition-colors">Projects</a>
                <a href="/patents/" class="text-gray-700 hover:text-blue-600 transition-colors">Patents</a>
                <a href="/speaking/" class="text-gray-700 hover:text-blue-600 transition-colors">Speaking</a>
                <a href="/contact/" class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">Contact</a>
            </div>
        </div>
    </div>
</nav>
{{end}}
```

### Phase 3: Update HTMX to Fetch Static JSON

**web/templates/pages/home.html** (Featured Projects section):

Change from:
```html
<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8"
     hx-get="/api/projects?featured=true"
     hx-trigger="load"
     hx-swap="innerHTML">
```

To:
```html
<div id="featured-projects" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
  <!-- Loading skeleton -->
  <div class="animate-pulse bg-white rounded-xl p-6 shadow-md">
    <div class="h-48 bg-gray-200 rounded-lg mb-4"></div>
    <div class="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
    <div class="h-4 bg-gray-200 rounded w-1/2"></div>
  </div>
</div>

<script>
// Fetch and render featured projects
fetch('/api/projects-featured.json')
  .then(res => res.json())
  .then(data => {
    const container = document.getElementById('featured-projects');
    container.innerHTML = data.projects.map(project => `
      <div class="bg-white rounded-xl shadow-lg overflow-hidden transform transition-all duration-300 hover:scale-105 hover:shadow-2xl">
        <div class="relative h-48 bg-gradient-to-br from-blue-500 to-purple-600">
          ${project.images && project.images[0] ?
            `<img src="${project.images[0]}" alt="${project.title}" class="w-full h-full object-cover opacity-80">` : ''}
          <div class="absolute top-4 right-4">
            <span class="px-3 py-1 bg-white bg-opacity-90 rounded-full text-xs font-semibold text-gray-800">
              ${project.category}
            </span>
          </div>
        </div>
        <div class="p-6">
          <h3 class="text-xl font-bold text-gray-900 mb-2">${project.title}</h3>
          <p class="text-gray-600 mb-4">${project.short_description}</p>
          ${project.impact ? `
            <div class="mb-4 p-3 bg-green-50 rounded-lg">
              <p class="text-sm text-green-800 font-medium">${project.impact.metric}</p>
            </div>
          ` : ''}
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-500">${project.timeline}</span>
            <a href="/projects/#${project.id}" class="text-blue-600 hover:text-blue-700 font-semibold">
              View Details →
            </a>
          </div>
        </div>
      </div>
    `).join('');
  });
</script>
```

### Phase 4: Update Contact Form

Replace server-side handler with Formspree:

**web/templates/pages/contact.html:**

```html
<form action="https://formspree.io/f/YOUR_FORM_ID" method="POST" class="space-y-6">
  <div>
    <label for="name" class="block text-sm font-medium text-gray-700">Name</label>
    <input type="text" name="name" id="name" required
           class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500">
  </div>

  <div>
    <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
    <input type="email" name="email" id="email" required
           class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500">
  </div>

  <div>
    <label for="message" class="block text-sm font-medium text-gray-700">Message</label>
    <textarea name="message" id="message" rows="4" required
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500"></textarea>
  </div>

  <div id="form-status" class="hidden p-4 rounded-md"></div>

  <button type="submit" class="w-full bg-blue-600 text-white py-3 px-6 rounded-lg hover:bg-blue-700 transition-colors">
    Send Message
  </button>
</form>

<script>
document.querySelector('form').addEventListener('submit', async (e) => {
  e.preventDefault();
  const form = e.target;
  const status = document.getElementById('form-status');

  try {
    const response = await fetch(form.action, {
      method: 'POST',
      body: new FormData(form),
      headers: { 'Accept': 'application/json' }
    });

    if (response.ok) {
      status.className = 'p-4 rounded-md bg-green-50 text-green-800';
      status.textContent = 'Thank you! Your message has been sent.';
      form.reset();
    } else {
      throw new Error('Form submission failed');
    }
  } catch (error) {
    status.className = 'p-4 rounded-md bg-red-50 text-red-800';
    status.textContent = 'Oops! There was a problem sending your message.';
  }

  status.classList.remove('hidden');
});
</script>
```

### Phase 5: GitHub Actions Workflow

Create `.github/workflows/deploy-github-pages.yml`:

```yaml
name: Deploy to GitHub Pages

on:
  push:
    branches: [ main ]
  workflow_dispatch: {}

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: "pages"
  cancel-in-progress: false

jobs:
  build:
    name: Build Static Site
    runs-on: ubuntu-latest

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '18'
          cache: 'npm'

      - name: Install Node dependencies
        run: npm ci

      - name: Build Tailwind CSS
        run: npm run build:css

      - name: Build static site
        run: go run build.go

      - name: Create 404 page
        run: cp dist/index.html dist/404.html

      - name: Upload artifact
        uses: actions/upload-pages-artifact@v3
        with:
          path: './dist'

  deploy:
    name: Deploy to GitHub Pages
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    needs: build

    steps:
      - name: Deploy to GitHub Pages
        id: deployment
        uses: actions/deploy-pages@v4
```

## Setup Instructions

### 1. Enable GitHub Pages

1. Go to repository **Settings**
2. Navigate to **Pages** (left sidebar)
3. Under "Build and deployment":
   - **Source**: GitHub Actions
4. Click **Save**

### 2. Configure Formspree (Contact Form)

1. Go to [formspree.io](https://formspree.io)
2. Sign up (free tier: 50 submissions/month)
3. Create a new form
4. Copy your form ID (e.g., `xpzgkjvw`)
5. Update contact.html with your form ID:
   ```html
   <form action="https://formspree.io/f/YOUR_FORM_ID" method="POST">
   ```

### 3. Build and Test Locally

```bash
# Install dependencies
npm install

# Build CSS
npm run build:css

# Generate static site
go run build.go

# Test locally with a simple server
cd dist
python3 -m http.server 8080
# Visit http://localhost:8080
```

### 4. Deploy

```bash
# Commit changes
git add .
git commit -m "Convert to static site for GitHub Pages"
git push origin main

# GitHub Actions will automatically build and deploy
```

### 5. Access Your Site

Your site will be available at:
- `https://USERNAME.github.io/REPOSITORY/`

Or with a custom domain (see below).

## Custom Domain Setup

### 1. Add Custom Domain in GitHub

1. Go to Settings → Pages
2. Under "Custom domain", enter your domain: `portfolio.example.com`
3. Click **Save**
4. GitHub will create a `CNAME` file in your repository

### 2. Configure DNS

Add these records at your domain registrar:

**For apex domain (example.com):**
```
Type: A
Name: @
Value: 185.199.108.153
Value: 185.199.109.153
Value: 185.199.110.153
Value: 185.199.111.153
```

**For subdomain (www.example.com or portfolio.example.com):**
```
Type: CNAME
Name: www (or portfolio)
Value: USERNAME.github.io
```

### 3. Enforce HTTPS

1. Wait for DNS to propagate (5-30 minutes)
2. In Settings → Pages
3. Check **Enforce HTTPS**

## Project Structure After Conversion

```
portfolio-website/
├── build.go                    # Static site generator
├── dist/                       # Generated static site (gitignored)
│   ├── index.html
│   ├── projects/
│   │   └── index.html
│   ├── speaking/
│   ├── patents/
│   ├── contact/
│   ├── api/
│   │   ├── projects.json
│   │   ├── projects-featured.json
│   │   ├── speaking.json
│   │   └── patents.json
│   ├── static/
│   │   ├── dist.css
│   │   └── images/
│   └── 404.html
├── web/
│   ├── templates/
│   └── static/
└── .github/workflows/
    └── deploy-github-pages.yml
```

## Advantages

✅ **Free Hosting** - Zero cost for hosting and SSL
✅ **Fast Performance** - Served via GitHub's global CDN
✅ **Simple Deployment** - Just push to main branch
✅ **HTTPS Included** - Free SSL certificates
✅ **Custom Domain** - Easy to configure
✅ **High Availability** - GitHub's infrastructure
✅ **No Server Management** - Static files only
✅ **Version Control** - Full history in Git

## Limitations

❌ **No Server-Side Logic** - Everything is client-side
❌ **No Real-Time Data** - Content updated on build only
❌ **Build Required** - Must regenerate on every change
❌ **Contact Form** - Requires external service
❌ **100 GB Bandwidth** - Free tier limit
❌ **Public Repository** - Private repos need GitHub Pro

## Comparison with Dynamic Options

| Feature | GitHub Pages | Lightsail | App Runner |
|---------|-------------|-----------|------------|
| **Cost** | $0 | $7/mo | $10-40/mo |
| **Server-side** | ❌ No | ✅ Yes | ✅ Yes |
| **Real-time data** | ❌ No | ✅ Yes | ✅ Yes |
| **Setup** | ⭐ Easy | ⭐ Easy | ⭐⭐ Medium |
| **Performance** | ⭐⭐⭐ Fast | ⭐⭐ Good | ⭐⭐ Good |
| **Bandwidth** | 100 GB | Included | Pay-per-use |
| **Custom domain** | ✅ Free | ✅ Free | ✅ Yes |
| **SSL/HTTPS** | ✅ Free | ✅ Free | ✅ Free |

## When to Use GitHub Pages

Choose GitHub Pages if:
- ✅ Your portfolio doesn't need real-time data
- ✅ You want zero hosting costs
- ✅ You're comfortable with static site generation
- ✅ You don't need server-side processing
- ✅ Traffic is under 100 GB/month
- ✅ You want the simplest possible deployment

**Don't use if:**
- ❌ You need server-side rendering
- ❌ You need real-time data updates
- ❌ You need complex backend logic
- ❌ You need private repository hosting (without Pro)

## Maintenance

### Update Content

```bash
# Make changes to data files or templates
vim internal/infrastructure/data/projects.json

# Rebuild and deploy
npm run build:css
go run build.go

# Test locally
cd dist && python3 -m http.server 8080

# Deploy
git add .
git commit -m "Update projects"
git push origin main
```

### Monitor

- **Traffic**: GitHub repository Insights → Traffic
- **Builds**: Actions tab shows build history
- **Errors**: Check Actions logs for build failures

## Troubleshooting

### Site Not Updating

```bash
# Check Actions tab for build status
# Force rebuild by pushing an empty commit
git commit --allow-empty -m "Trigger rebuild"
git push origin main
```

### 404 Errors

```bash
# Ensure 404.html exists
cp dist/index.html dist/404.html

# Check paths are relative, not absolute
# Use /projects/ not http://localhost/projects
```

### CSS Not Loading

```bash
# Verify dist.css is in dist/static/
# Check path in templates:
# Should be: /static/dist.css
# Not: http://localhost:8080/static/dist.css
```

### Form Not Working

```bash
# Verify Formspree form ID is correct
# Check browser console for errors
# Test form directly on formspree.io
```

## Alternative Static Hosts

If GitHub Pages doesn't meet your needs:

### Netlify
- **Cost**: Free (100 GB bandwidth)
- **Pros**: Better build system, form handling
- **Deploy**: Connect GitHub repo

### Vercel
- **Cost**: Free (100 GB bandwidth)
- **Pros**: Excellent performance, Edge functions
- **Deploy**: Import GitHub repo

### Cloudflare Pages
- **Cost**: Free (Unlimited bandwidth!)
- **Pros**: Best performance, global CDN
- **Deploy**: Connect GitHub repo

## Next Steps

1. [Create static site generator](./github-pages-setup.md)
2. [Test locally before deploying](#build-and-test-locally)
3. [Set up custom domain](#custom-domain-setup)
4. [Configure form service](#configure-formspree-contact-form)
5. [Monitor with analytics](./analytics.md)

---

**Ready to deploy?** Follow the [Setup Instructions](#setup-instructions) to get started with GitHub Pages for free! 🚀
