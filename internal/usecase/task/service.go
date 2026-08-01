package task

import (
	"context"
	"fmt"
	"strings"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/domain"
	"github.com/google/uuid"
)

// Service implements task business rules on top of a TaskRepository.
type Service struct {
	repo domain.TaskRepository
}

// NewService creates a task Service.
func NewService(repo domain.TaskRepository) *Service {
	return &Service{repo: repo}
}

// CreateInput holds data required to create a task.
type CreateInput struct {
	Title       string
	Description string
	Status      domain.TaskStatus
	Assignee    string
}

// UpdateInput holds optional fields for updating a task.
type UpdateInput struct {
	Title       *string
	Description *string
	Status      *domain.TaskStatus
	Assignee    *string
}

// Create persists a new task.
func (s *Service) Create(ctx context.Context, input CreateInput) (*domain.Task, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, fmt.Errorf("%w: title is required", domain.ErrInvalidInput)
	}

	status := input.Status
	if status == "" {
		status = domain.StatusTodo
	}
	if !status.Valid() {
		return nil, fmt.Errorf("%w: invalid status", domain.ErrInvalidInput)
	}

	task := &domain.Task{
		Title:       title,
		Description: strings.TrimSpace(input.Description),
		Status:      status,
		Assignee:    strings.TrimSpace(input.Assignee),
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// GetByID returns a single task.
func (s *Service) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, id)
}

// ListInput holds list query parameters.
type ListInput struct {
	Status   *domain.TaskStatus
	Assignee *string
	Page     int
	PageSize int
}

// List returns tasks matching optional filters with pagination.
func (s *Service) List(ctx context.Context, input ListInput) (*domain.ListResult, error) {
	page := input.Page
	if page < 1 {
		page = 1
	}

	pageSize := input.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	if input.Status != nil && !input.Status.Valid() {
		return nil, fmt.Errorf("%w: invalid status filter", domain.ErrInvalidInput)
	}

	return s.repo.List(ctx, domain.ListFilter{
		Status:   input.Status,
		Assignee: input.Assignee,
		Page:     page,
		PageSize: pageSize,
	})
}

// Update modifies an existing task.
func (s *Service) Update(ctx context.Context, id string, input UpdateInput) (*domain.Task, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}

	if input.Title == nil && input.Description == nil && input.Status == nil && input.Assignee == nil {
		return nil, fmt.Errorf("%w: at least one field must be provided", domain.ErrInvalidInput)
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.Title != nil {
		title := strings.TrimSpace(*input.Title)
		if title == "" {
			return nil, fmt.Errorf("%w: title cannot be empty", domain.ErrInvalidInput)
		}
		task.Title = title
	}

	if input.Description != nil {
		task.Description = strings.TrimSpace(*input.Description)
	}

	if input.Status != nil {
		if !input.Status.Valid() {
			return nil, fmt.Errorf("%w: invalid status", domain.ErrInvalidInput)
		}
		task.Status = *input.Status
	}

	if input.Assignee != nil {
		task.Assignee = strings.TrimSpace(*input.Assignee)
	}

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// Delete removes a task by ID.
func (s *Service) Delete(ctx context.Context, id string) error {
	if err := validateID(id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func validateID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: invalid task id", domain.ErrInvalidInput)
	}
	return nil
}
