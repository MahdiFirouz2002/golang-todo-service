package handler

import (
	"fmt"
	"net/http"
	"strconv"

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
	input, err := parseListInput(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	result, err := h.service.List(c.Request.Context(), input)
	if err != nil {
		respondError(c, err)
		return
	}

	if result.Items == nil {
		result.Items = []*domain.Task{}
	}

	c.JSON(http.StatusOK, result)
}

func parseListInput(c *gin.Context) (task.ListInput, error) {
	input := task.ListInput{Page: 1, PageSize: 20}

	if page := c.Query("page"); page != "" {
		value, err := strconv.Atoi(page)
		if err != nil || value < 1 {
			return input, fmt.Errorf("invalid page")
		}
		input.Page = value
	}

	if pageSize := c.Query("page_size"); pageSize != "" {
		value, err := strconv.Atoi(pageSize)
		if err != nil || value < 1 {
			return input, fmt.Errorf("invalid page_size")
		}
		input.PageSize = value
	}

	if status := c.Query("status"); status != "" {
		taskStatus := domain.TaskStatus(status)
		if !taskStatus.Valid() {
			return input, fmt.Errorf("invalid status")
		}
		input.Status = &taskStatus
	}

	if assignee := c.Query("assignee"); assignee != "" {
		input.Assignee = &assignee
	}

	return input, nil
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
