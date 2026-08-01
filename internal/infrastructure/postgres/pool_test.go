package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/config"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/infrastructure/postgres"
	"github.com/stretchr/testify/require"
)

func TestNewPool_InvalidURL(t *testing.T) {
	cfg := &config.Config{
		DatabaseURL:       "://invalid",
		DBMaxConns:        5,
		DBMinConns:        1,
		DBMaxConnLifetime: time.Hour,
	}

	_, err := postgres.NewPool(context.Background(), cfg)
	require.Error(t, err)
}

func TestMigrate_InvalidPool(t *testing.T) {
	err := postgres.Migrate(context.Background(), nil)
	require.Error(t, err)
}
