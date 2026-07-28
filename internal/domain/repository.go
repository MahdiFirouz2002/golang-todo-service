package domain

import (
	"context"
	"errors"
)

var (
	// ErrNotFound is returned when a requested entity does not exist.
	ErrNotFound = errors.New("not found")

	// ErrInvalidInput is returned when domain validation fails.
	ErrInvalidInput = errors.New("invalid input")
)

// TaskRepository defines persistence operations for Task entities.
// Implementations live in the infrastructure layer.
type TaskRepository interface {
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id string) (*Task, error)
	List(ctx context.Context) ([]*Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}
