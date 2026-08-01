package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/config"
	httpserver "github.com/MahdiFirouz2002/golang-todo-service/internal/delivery/http"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/delivery/http/handler"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/infrastructure/postgres"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/observability/metrics"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/observability/tracing"
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

	shutdownTracing, err := tracing.Init(ctx, "task-manager")
	if err != nil {
		return err
	}
	defer func() {
		_ = shutdownTracing(context.Background())
	}()

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

	metricsCtx, metricsCancel := context.WithCancel(ctx)
	defer metricsCancel()
	metrics.StartTasksCountPoller(metricsCtx, taskRepo, 15*time.Second)

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
