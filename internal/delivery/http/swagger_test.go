package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httpserver "github.com/MahdiFirouz2002/golang-todo-service/internal/delivery/http"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/delivery/http/handler"
	"github.com/stretchr/testify/require"
)

func TestSwaggerRoutes(t *testing.T) {
	router := httpserver.NewRouter(httpserver.Dependencies{
		Health: handler.NewHealthHandler(nil),
	})

	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "openapi:")

	req = httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}
