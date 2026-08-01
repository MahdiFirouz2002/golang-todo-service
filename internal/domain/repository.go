package domain

import "context"

// ListFilter defines optional query constraints for listing tasks.
type ListFilter struct {
	Status   *TaskStatus
	Assignee *string
	Page     int
	PageSize int
}

// ListResult is a paginated collection of tasks.
type ListResult struct {
	Items      []*Task `json:"items"`
	Total      int64   `json:"total"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	TotalPages int     `json:"total_pages"`
}

// TaskRepository defines persistence operations for Task entities.
// Implementations live in the infrastructure layer.
type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id string) (*Task, error)
	List(ctx context.Context, filter ListFilter) (*ListResult, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}
