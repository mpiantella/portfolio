package projectusecase

import (
	"portfolio/internal/domain/project"
)

// ListProjects returns all projects from the repository.
func ListProjects(repo project.Repository) ([]project.Project, error) {
	return repo.List()
}
