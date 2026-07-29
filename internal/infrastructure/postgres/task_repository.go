package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TaskRepository persists tasks in PostgreSQL.
type TaskRepository struct {
	pool *pgxpool.Pool
}

// NewTaskRepository creates a PostgreSQL-backed TaskRepository.
func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{pool: pool}
}

func (r *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	if task.ID == "" {
		task.ID = uuid.NewString()
	}

	now := time.Now().UTC()
	if task.CreatedAt.IsZero() {
		task.CreatedAt = now
	}
	if task.UpdatedAt.IsZero() {
		task.UpdatedAt = now
	}
	if task.Status == "" {
		task.Status = domain.StatusTodo
	}

	const query = `
		INSERT INTO tasks (id, title, description, status, assignee, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := r.pool.Exec(ctx, query,
		task.ID,
		task.Title,
		task.Description,
		string(task.Status),
		task.Assignee,
		task.CreatedAt,
		task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert task: %w", err)
	}

	return nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	const query = `
		SELECT id, title, description, status, assignee, created_at, updated_at
		FROM tasks
		WHERE id = $1`

	row := r.pool.QueryRow(ctx, query, id)
	task, err := scanTask(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("select task by id: %w", err)
	}

	return task, nil
}

func (r *TaskRepository) List(ctx context.Context) ([]*domain.Task, error) {
	const query = `
		SELECT id, title, description, status, assignee, created_at, updated_at
		FROM tasks
		ORDER BY created_at DESC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("select tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]*domain.Task, 0)
	for rows.Next() {
		task, err := scanTask(rows)
		if err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	return tasks, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
	task.UpdatedAt = time.Now().UTC()

	const query = `
		UPDATE tasks
		SET title = $2,
		    description = $3,
		    status = $4,
		    assignee = $5,
		    updated_at = $6
		WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query,
		task.ID,
		task.Title,
		task.Description,
		string(task.Status),
		task.Assignee,
		task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM tasks WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *TaskRepository) Count(ctx context.Context) (int64, error) {
	const query = `SELECT COUNT(*) FROM tasks`

	var count int64
	if err := r.pool.QueryRow(ctx, query).Scan(&count); err != nil {
		return 0, fmt.Errorf("count tasks: %w", err)
	}

	return count, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanTask(row rowScanner) (*domain.Task, error) {
	var (
		task   domain.Task
		status string
	)

	if err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&status,
		&task.Assignee,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		return nil, err
	}

	task.Status = domain.TaskStatus(status)
	return &task, nil
}
