package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/config"
	httpserver "github.com/MahdiFirouz2002/golang-todo-service/internal/delivery/http"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	if err := run(); err != nil {
		slog.Error("application stopped with error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := httpserver.NewRouter(httpserver.Dependencies{
		Health: handler.NewHealthHandler(),
	})

	srv := httpserver.New(cfg, router)

	errCh := make(chan error, 1)
	go func() {
		slog.Info("http server starting", "addr", cfg.Addr(), "env", cfg.AppEnv)
		errCh <- srv.Start()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		slog.Info("shutdown signal received")
		if err := srv.Shutdown(context.Background()); err != nil {
			return err
		}
		slog.Info("server stopped gracefully")
		return nil
	}
}
