package projectusecase_test

import (
	"testing"

	mem "portfolio/internal/infrastructure/persistence/memory"
	uc "portfolio/internal/usecase/project"
)

func TestListProjects(t *testing.T) {
	repo := mem.NewProjectRepository()
	projects, err := uc.ListProjects(repo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(projects) == 0 {
		t.Fatalf("expected at least one project")
	}
}
