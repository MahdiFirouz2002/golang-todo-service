//go:build integration

package postgres_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/config"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/domain"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/infrastructure/postgres"
	"github.com/stretchr/testify/require"
)

func TestTaskRepository_CRUD(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		t.Skip("DATABASE_URL not set; skipping integration test")
	}

	cfg := &config.Config{
		DatabaseURL:       databaseURL,
		DBMaxConns:        5,
		DBMinConns:        1,
		DBMaxConnLifetime: time.Hour,
	}

	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, cfg)
	require.NoError(t, err)
	defer pool.Close()

	require.NoError(t, postgres.Migrate(ctx, pool))

	repo := postgres.NewTaskRepository(pool)
	task := &domain.Task{
		Title:       "Integration task",
		Description: "created in test",
		Status:      domain.StatusTodo,
		Assignee:    "tester",
	}

	require.NoError(t, repo.Create(ctx, task))
	require.NotEmpty(t, task.ID)

	found, err := repo.GetByID(ctx, task.ID)
	require.NoError(t, err)
	require.Equal(t, task.Title, found.Title)

	result, err := repo.List(ctx, domain.ListFilter{Page: 1, PageSize: 20})
	require.NoError(t, err)
	require.NotEmpty(t, result.Items)

	found.Title = "Updated"
	require.NoError(t, repo.Update(ctx, found))

	require.NoError(t, repo.Delete(ctx, task.ID))

	_, err = repo.GetByID(ctx, task.ID)
	require.ErrorIs(t, err, domain.ErrNotFound)
}
