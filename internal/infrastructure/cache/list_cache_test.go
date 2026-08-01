package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/domain"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/infrastructure/cache"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestListCache_SetGetInvalidate(t *testing.T) {
	server := miniredis.RunT(t)
	client := redis.NewClient(&redis.Options{Addr: server.Addr()})
	listCache := cache.NewListCache(client, time.Minute)

	filter := domain.ListFilter{Page: 1, PageSize: 20}
	result := &domain.ListResult{
		Items:    []*domain.Task{{ID: "1", Title: "Cached"}},
		Total:    1,
		Page:     1,
		PageSize: 20,
	}

	ctx := context.Background()
	require.NoError(t, listCache.Set(ctx, filter, result))

	cached, ok := listCache.Get(ctx, filter)
	require.True(t, ok)
	require.Equal(t, result.Items[0].Title, cached.Items[0].Title)

	require.NoError(t, listCache.Invalidate(ctx))

	_, ok = listCache.Get(ctx, filter)
	require.False(t, ok)
}
