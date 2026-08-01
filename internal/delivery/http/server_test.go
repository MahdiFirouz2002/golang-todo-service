package httpserver_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/config"
	httpserver "github.com/MahdiFirouz2002/golang-todo-service/internal/delivery/http"
	"github.com/gin-gonic/gin"
)

func TestServer_Shutdown(t *testing.T) {
	cfg := &config.Config{
		HTTPPort:        "0",
		ReadTimeout:     time.Second,
		WriteTimeout:    time.Second,
		IdleTimeout:     time.Second,
		ShutdownTimeout: time.Second,
	}

	engine := gin.New()
	srv := httpserver.New(cfg, engine)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	time.Sleep(100 * time.Millisecond)

	if err := srv.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown() error = %v", err)
	}

	select {
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			t.Fatalf("Start() error = %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("server did not stop in time")
	}
}
