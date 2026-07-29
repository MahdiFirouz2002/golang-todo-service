package handler

import (
	"net/http"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/delivery/http/dto"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/domain"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/usecase/task"
	"github.com/gin-gonic/gin"
)

// TaskHandler exposes task CRUD endpoints.
type TaskHandler struct {
	service *task.Service
}

// NewTaskHandler creates a TaskHandler.
func NewTaskHandler(service *task.Service) *TaskHandler {
	return &TaskHandler{service: service}
}

// Create handles POST /api/v1/tasks
func (h *TaskHandler) Create(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	created, err := h.service.Create(c.Request.Context(), task.CreateInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Assignee:    req.Assignee,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, created)
}

// List handles GET /api/v1/tasks
func (h *TaskHandler) List(c *gin.Context) {
	tasks, err := h.service.List(c.Request.Context())
	if err != nil {
		respondError(c, err)
		return
	}

	if tasks == nil {
		tasks = []*domain.Task{}
	}

	c.JSON(http.StatusOK, tasks)
}

// Get handles GET /api/v1/tasks/:id
func (h *TaskHandler) Get(c *gin.Context) {
	found, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, found)
}

// Update handles PUT /api/v1/tasks/:id
func (h *TaskHandler) Update(c *gin.Context) {
	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	updated, err := h.service.Update(c.Request.Context(), c.Param("id"), task.UpdateInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Assignee:    req.Assignee,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, updated)
}

// Delete handles DELETE /api/v1/tasks/:id
func (h *TaskHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		respondError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
