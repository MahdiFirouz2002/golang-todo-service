package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/delivery/http/handler"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/domain"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/usecase/task"
	"github.com/gin-gonic/gin"
)

type stubTaskRepository struct {
	task *domain.Task
}

func (s *stubTaskRepository) Create(_ context.Context, task *domain.Task) error {
	if task.ID == "" {
		task.ID = "11111111-1111-1111-1111-111111111111"
	}
	copyTask := *task
	s.task = &copyTask
	return nil
}

func (s *stubTaskRepository) GetByID(_ context.Context, id string) (*domain.Task, error) {
	if s.task == nil || s.task.ID != id {
		return nil, domain.ErrNotFound
	}
	copyTask := *s.task
	return &copyTask, nil
}

func (s *stubTaskRepository) List(_ context.Context) ([]*domain.Task, error) {
	if s.task == nil {
		return []*domain.Task{}, nil
	}
	copyTask := *s.task
	return []*domain.Task{&copyTask}, nil
}

func (s *stubTaskRepository) Update(_ context.Context, task *domain.Task) error {
	if s.task == nil || s.task.ID != task.ID {
		return domain.ErrNotFound
	}
	copyTask := *task
	s.task = &copyTask
	return nil
}

func (s *stubTaskRepository) Delete(_ context.Context, id string) error {
	if s.task == nil || s.task.ID != id {
		return domain.ErrNotFound
	}
	s.task = nil
	return nil
}

func (s *stubTaskRepository) Count(_ context.Context) (int64, error) {
	if s.task == nil {
		return 0, nil
	}
	return 1, nil
}

func newTaskRouter(t *testing.T) (*gin.Engine, *stubTaskRepository) {
	t.Helper()

	repo := &stubTaskRepository{}
	svc := task.NewService(repo)
	h := handler.NewTaskHandler(svc)

	router := gin.New()
	router.POST("/api/v1/tasks", h.Create)
	router.GET("/api/v1/tasks", h.List)
	router.GET("/api/v1/tasks/:id", h.Get)
	router.PUT("/api/v1/tasks/:id", h.Update)
	router.DELETE("/api/v1/tasks/:id", h.Delete)

	return router, repo
}

func TestTaskHandler_CreateAndGet(t *testing.T) {
	router, _ := newTaskRouter(t)

	body := bytes.NewBufferString(`{"title":"Ship feature","status":"todo"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d", rec.Code, http.StatusCreated)
	}

	var created domain.Task
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal create response: %v", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+created.ID, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestTaskHandler_ListEmpty(t *testing.T) {
	router, _ := newTaskRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", rec.Code, http.StatusOK)
	}

	var tasks []domain.Task
	if err := json.Unmarshal(rec.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("unmarshal list response: %v", err)
	}
	if len(tasks) != 0 {
		t.Fatalf("len(tasks) = %d, want 0", len(tasks))
	}
}

func TestTaskHandler_UpdateAndDelete(t *testing.T) {
	router, repo := newTaskRouter(t)

	repo.task = &domain.Task{
		ID:        "22222222-2222-2222-2222-222222222222",
		Title:     "Old title",
		Status:    domain.StatusTodo,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	body := bytes.NewBufferString(`{"title":"New title"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/tasks/"+repo.task.ID, body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("update status = %d, want %d", rec.Code, http.StatusOK)
	}

	req = httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/"+repo.task.ID, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d", rec.Code, http.StatusNoContent)
	}
}
