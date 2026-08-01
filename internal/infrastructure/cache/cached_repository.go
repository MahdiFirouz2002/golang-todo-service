package cache

import (
	"context"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/domain"
)

// CachedTaskRepository decorates a TaskRepository with cache-aside list reads.
type CachedTaskRepository struct {
	repo  domain.TaskRepository
	cache *ListCache
}

// NewCachedTaskRepository wraps a repository with list caching.
func NewCachedTaskRepository(repo domain.TaskRepository, cache *ListCache) *CachedTaskRepository {
	return &CachedTaskRepository{repo: repo, cache: cache}
}

func (c *CachedTaskRepository) Create(ctx context.Context, task *domain.Task) error {
	if err := c.repo.Create(ctx, task); err != nil {
		return err
	}
	_ = c.cache.Invalidate(ctx)
	return nil
}

func (c *CachedTaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	return c.repo.GetByID(ctx, id)
}

func (c *CachedTaskRepository) List(ctx context.Context, filter domain.ListFilter) (*domain.ListResult, error) {
	if cached, ok := c.cache.Get(ctx, filter); ok {
		return cached, nil
	}

	result, err := c.repo.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	_ = c.cache.Set(ctx, filter, result)
	return result, nil
}

func (c *CachedTaskRepository) Update(ctx context.Context, task *domain.Task) error {
	if err := c.repo.Update(ctx, task); err != nil {
		return err
	}
	_ = c.cache.Invalidate(ctx)
	return nil
}

func (c *CachedTaskRepository) Delete(ctx context.Context, id string) error {
	if err := c.repo.Delete(ctx, id); err != nil {
		return err
	}
	_ = c.cache.Invalidate(ctx)
	return nil
}

func (c *CachedTaskRepository) Count(ctx context.Context) (int64, error) {
	return c.repo.Count(ctx)
}
