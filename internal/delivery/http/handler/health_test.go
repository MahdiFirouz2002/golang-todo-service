package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestHealthHandler_Live(t *testing.T) {
	h := handler.NewHealthHandler(nil)
	router := gin.New()
	router.GET("/health/live", h.Live)

	req := httptest.NewRequest(http.MethodGet, "/health/live", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if body["status"] != "alive" {
		t.Errorf("status = %q, want alive", body["status"])
	}
}

func TestHealthHandler_Ready(t *testing.T) {
	t.Run("ready when ping succeeds", func(t *testing.T) {
		h := handler.NewHealthHandler(func(_ context.Context) error { return nil })
		router := gin.New()
		router.GET("/health/ready", h.Ready)

		req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}

		var body map[string]any
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("unmarshal body: %v", err)
		}
		if body["status"] != "ready" {
			t.Errorf("status = %v, want ready", body["status"])
		}
		if _, ok := body["uptime"]; !ok {
			t.Error("expected uptime field in response")
		}
	})

	t.Run("not ready when ping fails", func(t *testing.T) {
		h := handler.NewHealthHandler(func(_ context.Context) error {
			return assertErr("db down")
		})
		router := gin.New()
		router.GET("/health/ready", h.Ready)

		req := httptest.NewRequest(http.MethodGet, "/health/ready", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusServiceUnavailable {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusServiceUnavailable)
		}
	})
}

type assertErr string

func (e assertErr) Error() string { return string(e) }
