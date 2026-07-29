package dto

import "github.com/MahdiFirouz2002/golang-todo-service/internal/domain"

// CreateTaskRequest is the payload for creating a task.
type CreateTaskRequest struct {
	Title       string            `json:"title" binding:"required"`
	Description string            `json:"description"`
	Status      domain.TaskStatus `json:"status"`
	Assignee    string            `json:"assignee"`
}

// UpdateTaskRequest is the payload for updating a task.
type UpdateTaskRequest struct {
	Title       *string            `json:"title"`
	Description *string            `json:"description"`
	Status      *domain.TaskStatus `json:"status"`
	Assignee    *string            `json:"assignee"`
}

// ErrorResponse is a standard API error body.
type ErrorResponse struct {
	Error string `json:"error"`
}
