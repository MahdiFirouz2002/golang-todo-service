package metrics

import (
	"context"
	"log/slog"
	"time"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/domain"
)

// StartTasksCountPoller periodically updates the tasks_count gauge.
func StartTasksCountPoller(ctx context.Context, repo domain.TaskRepository, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		update := func() {
			count, err := repo.Count(ctx)
			if err != nil {
				slog.Warn("failed to update tasks_count metric", "error", err)
				return
			}
			SetTasksCount(float64(count))
		}

		update()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				update()
			}
		}
	}()
}
