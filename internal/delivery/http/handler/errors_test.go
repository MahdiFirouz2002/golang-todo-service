package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/delivery/http/handler"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/domain"
	"github.com/gin-gonic/gin"
)

func TestRespondError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "not found", err: domain.ErrNotFound, wantStatus: http.StatusNotFound},
		{name: "invalid input", err: domain.ErrInvalidInput, wantStatus: http.StatusBadRequest},
		{name: "internal", err: errors.New("boom"), wantStatus: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/err", func(c *gin.Context) {
				handler.RespondError(c, tt.err)
			})

			req := httptest.NewRequest(http.MethodGet, "/err", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}
