package metrics_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/observability/metrics"
	"github.com/gin-gonic/gin"
)

func TestMiddleware_RecordsMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(metrics.Middleware())
	router.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestSetTasksCount(t *testing.T) {
	metrics.SetTasksCount(42)
}
