package memory

import "portfolio/internal/domain/project"

// InMemoryProjectRepository is a simple in-memory repository used for demos and tests.
type InMemoryProjectRepository struct {
	projects []project.Project
}

// NewProjectRepository constructs a repository with seeded data.
func NewProjectRepository() *InMemoryProjectRepository {
	return &InMemoryProjectRepository{projects: []project.Project{
		{ID: "1", Title: "Demo Project", Description: "A demo project", URL: "https://example.com"},
		{ID: "2", Title: "Another One", Description: "Another demo", URL: "https://example.org"},
		{ID: "3", Title: "Yet Another One", Description: "Yet another demo", URL: "https://example.net"},
	}}
}

// List returns all projects.
func (r *InMemoryProjectRepository) List() ([]project.Project, error) {
	return r.projects, nil
}
