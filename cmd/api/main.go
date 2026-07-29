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
	"github.com/MahdiFirouz2002/golang-todo-service/internal/infrastructure/postgres"
	taskusecase "github.com/MahdiFirouz2002/golang-todo-service/internal/usecase/task"
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

	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, cfg)
	if err != nil {
		return err
	}
	defer pool.Close()

	if err := postgres.Migrate(ctx, pool); err != nil {
		return err
	}

	taskRepo := postgres.NewTaskRepository(pool)
	taskService := taskusecase.NewService(taskRepo)

	router := httpserver.NewRouter(httpserver.Dependencies{
		Health: handler.NewHealthHandler(func(ctx context.Context) error {
			return postgres.Ping(ctx, pool)
		}),
		Tasks: handler.NewTaskHandler(taskService),
	})

	srv := httpserver.New(cfg, router)

	errCh := make(chan error, 1)
	go func() {
		slog.Info("http server starting", "addr", cfg.Addr(), "env", cfg.AppEnv)
		errCh <- srv.Start()
	}()

	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		return err
	case <-sigCtx.Done():
		slog.Info("shutdown signal received")
		if err := srv.Shutdown(context.Background()); err != nil {
			return err
		}
		slog.Info("server stopped gracefully")
		return nil
	}
}
