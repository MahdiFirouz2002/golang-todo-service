package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/delivery/http/handler"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/domain"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/mocks"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/usecase/task"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTaskRouterWithMock(t *testing.T) (*gin.Engine, *mocks.TaskRepository) {
	t.Helper()

	repo := mocks.NewTaskRepository(t)
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

func TestTaskHandler_Create_InvalidJSON(t *testing.T) {
	router, _ := newTaskRouterWithMock(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBufferString(`{`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTaskHandler_Create_ValidationError(t *testing.T) {
	router, _ := newTaskRouterWithMock(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBufferString(`{"title":""}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTaskHandler_Get_NotFound(t *testing.T) {
	router, repo := newTaskRouterWithMock(t)
	id := "11111111-1111-1111-1111-111111111111"
	repo.On("GetByID", mock.Anything, id).Return((*domain.Task)(nil), domain.ErrNotFound)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+id, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestTaskHandler_Get_InvalidID(t *testing.T) {
	router, _ := newTaskRouterWithMock(t)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks/not-a-uuid", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTaskHandler_List_Error(t *testing.T) {
	router, repo := newTaskRouterWithMock(t)
	repo.On("List", mock.Anything, mock.Anything).Return((*domain.ListResult)(nil), errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestTaskHandler_Update_InvalidJSON(t *testing.T) {
	router, _ := newTaskRouterWithMock(t)
	id := "11111111-1111-1111-1111-111111111111"

	req := httptest.NewRequest(http.MethodPut, "/api/v1/tasks/"+id, bytes.NewBufferString(`{`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestTaskHandler_Update_NotFound(t *testing.T) {
	router, repo := newTaskRouterWithMock(t)
	id := "11111111-1111-1111-1111-111111111111"
	repo.On("GetByID", mock.Anything, id).Return((*domain.Task)(nil), domain.ErrNotFound)

	body := bytes.NewBufferString(`{"title":"Updated"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/tasks/"+id, body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestTaskHandler_Delete_NotFound(t *testing.T) {
	router, repo := newTaskRouterWithMock(t)
	id := "11111111-1111-1111-1111-111111111111"
	repo.On("Delete", mock.Anything, id).Return(domain.ErrNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/"+id, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestTaskHandler_CreateAndGet_WithMock(t *testing.T) {
	router, repo := newTaskRouterWithMock(t)

	repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Task")).
		Run(func(args mock.Arguments) {
			taskArg := args.Get(1).(*domain.Task)
			taskArg.ID = "11111111-1111-1111-1111-111111111111"
			taskArg.CreatedAt = time.Now().UTC()
			taskArg.UpdatedAt = time.Now().UTC()
		}).
		Return(nil)

	body := bytes.NewBufferString(`{"title":"Ship feature","status":"todo"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var created domain.Task
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &created))

	repo.On("GetByID", mock.Anything, created.ID).Return(&created, nil)

	req = httptest.NewRequest(http.MethodGet, "/api/v1/tasks/"+created.ID, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestHealthHandler_Ready_WithPing(t *testing.T) {
	h := handler.NewHealthHandler(func(_ context.Context) error { return nil })
	router := gin.New()
	router.GET("/health/ready", h.Ready)

	req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}
