package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/domain"
	"github.com/redis/go-redis/v9"
)

const listKeyPrefix = "tasks:list:"

// ListCache provides cache-aside storage for task list queries.
type ListCache struct {
	client *redis.Client
	ttl    time.Duration
}

// NewListCache creates a Redis-backed list cache.
func NewListCache(client *redis.Client, ttl time.Duration) *ListCache {
	return &ListCache{client: client, ttl: ttl}
}

// Get returns a cached list result when present.
func (c *ListCache) Get(ctx context.Context, filter domain.ListFilter) (*domain.ListResult, bool) {
	key := listKey(filter)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}

	var result domain.ListResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, false
	}

	return &result, true
}

// Set stores a list result in cache.
func (c *ListCache) Set(ctx context.Context, filter domain.ListFilter, result *domain.ListResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal list result: %w", err)
	}

	key := listKey(filter)
	if err := c.client.Set(ctx, key, data, c.ttl).Err(); err != nil {
		return fmt.Errorf("set list cache: %w", err)
	}

	return nil
}

// Invalidate removes all cached list entries.
func (c *ListCache) Invalidate(ctx context.Context) error {
	iter := c.client.Scan(ctx, 0, listKeyPrefix+"*", 100).Iterator()
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("delete cache key: %w", err)
		}
	}
	return iter.Err()
}

func listKey(filter domain.ListFilter) string {
	payload := fmt.Sprintf("%v|%v|%d|%d", filter.Status, filter.Assignee, filter.Page, filter.PageSize)
	sum := sha256.Sum256([]byte(payload))
	return listKeyPrefix + hex.EncodeToString(sum[:])
}
