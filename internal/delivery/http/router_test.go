package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httpserver "github.com/MahdiFirouz2002/golang-todo-service/internal/delivery/http"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/delivery/http/handler"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/domain"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/mocks"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/usecase/task"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestNewRouter_HealthRoutes(t *testing.T) {
	router := httpserver.NewRouter(httpserver.Dependencies{
		Health: handler.NewHealthHandler(nil),
	})

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("live status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestNewRouter_TaskRoutesRegistered(t *testing.T) {
	repo := mocks.NewTaskRepository(t)
	repo.On("List", mock.Anything).Return([]*domain.Task{}, nil)

	svc := task.NewService(repo)
	router := httpserver.NewRouter(httpserver.Dependencies{
		Health: handler.NewHealthHandler(nil),
		Tasks:  handler.NewTaskHandler(svc),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", rec.Code, http.StatusOK)
	}
}
