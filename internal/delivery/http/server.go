package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/config"
	"github.com/gin-gonic/gin"
)

// Server wraps an http.Server with application configuration.
type Server struct {
	httpServer *http.Server
	cfg        *config.Config
}

// New creates an HTTP server bound to the given Gin engine and config.
func New(cfg *config.Config, engine *gin.Engine) *Server {
	return &Server{
		cfg: cfg,
		httpServer: &http.Server{
			Addr:         cfg.Addr(),
			Handler:      engine,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}
}

// Start begins serving HTTP traffic. It blocks until the server stops.
func (s *Server) Start() error {
	if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server failed: %w", err)
	}
	return nil
}

// Shutdown gracefully stops the server within the configured timeout.
func (s *Server) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, s.cfg.ShutdownTimeout)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}
	return nil
}
