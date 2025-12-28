package project

// Repository defines persistence operations for projects.
type Repository interface {
	List() ([]Project, error)
}
