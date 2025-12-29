package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
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
	generateCaseStudyPages(templates, outputDir, projects)

	// Copy static assets
	copyDir("web/static", filepath.Join(outputDir, "static"))

	// Generate API JSON files
	generateAPIFiles(outputDir, projects, patents, speaking)

	// Create 404 page
	copy404Page(outputDir)

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

	// Calculate stats
	totalCount := len(patents)
	grantedCount := 0
	pendingCount := 0
	for _, p := range patents {
		if p.Status == "Granted" {
			grantedCount++
		} else if p.Status == "Pending" {
			pendingCount++
		}
	}

	data := struct {
		Patents      []domain.Patent
		TotalCount   int
		GrantedCount int
		PendingCount int
	}{
		Patents:      patents,
		TotalCount:   totalCount,
		GrantedCount: grantedCount,
		PendingCount: pendingCount,
	}

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

	// Calculate stats
	totalAudience := 0
	topicsMap := make(map[string]bool)
	for _, s := range speaking {
		totalAudience += s.AudienceSize
		for _, topic := range s.Topics {
			topicsMap[topic] = true
		}
	}

	data := struct {
		SpeakingEngagements []domain.SpeakingEngagement
		TotalAudience       int
		UniqueTopicsCount   int
	}{
		SpeakingEngagements: speaking,
		TotalAudience:       totalAudience,
		UniqueTopicsCount:   len(topicsMap),
	}

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

func generateCaseStudyPages(tmpl *template.Template, outDir string, projects []domain.Project) {
	// Generate individual case study pages for each project
	for _, project := range projects {
		// Create directory for this project
		projectDir := filepath.Join(outDir, "projects", project.Slug)
		os.MkdirAll(projectDir, 0755)

		// Create the case study page
		f, err := os.Create(filepath.Join(projectDir, "index.html"))
		if err != nil {
			panic(err)
		}
		defer f.Close()

		// Execute template with project data
		if err := tmpl.ExecuteTemplate(f, "case-study.html", project); err != nil {
			panic(err)
		}

		fmt.Printf("✅ Generated case study: /projects/%s\n", project.Slug)
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
	data, err := os.ReadFile("internal/infrastructure/data/projects.json")
	if err != nil {
		panic(err)
	}
	var result struct {
		Projects []domain.Project `json:"projects"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		panic(err)
	}
	return result.Projects
}

func loadPatents() []domain.Patent {
	data, err := os.ReadFile("internal/infrastructure/data/patents.json")
	if err != nil {
		panic(err)
	}
	var result struct {
		Patents []domain.Patent `json:"patents"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		panic(err)
	}
	return result.Patents
}

func loadSpeaking() []domain.SpeakingEngagement {
	data, err := os.ReadFile("internal/infrastructure/data/speaking.json")
	if err != nil {
		panic(err)
	}
	var result struct {
		SpeakingEngagements []domain.SpeakingEngagement `json:"speaking_engagements"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		panic(err)
	}
	return result.SpeakingEngagements
}

func writeJSON(path string, data interface{}) {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		panic(err)
	}
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath)
	})
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func copy404Page(outDir string) {
	src := filepath.Join(outDir, "index.html")
	dst := filepath.Join(outDir, "404.html")
	if err := copyFile(src, dst); err != nil {
		fmt.Printf("Warning: Could not create 404.html: %v\n", err)
	}
}
