package redis

import (
	"context"
	"fmt"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/config"
	"github.com/redis/go-redis/v9"
)

// NewClient creates a Redis client from configuration.
func NewClient(cfg *config.Config) (*redis.Client, error) {
	opts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("parse redis url: %w", err)
	}

	client := redis.NewClient(opts)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	return client, nil
}
