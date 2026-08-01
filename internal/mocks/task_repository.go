package mocks

import (
	"context"
	"testing"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/domain"
	"github.com/stretchr/testify/mock"
)

// TaskRepository is a testify mock for domain.TaskRepository.
type TaskRepository struct {
	mock.Mock
}

// NewTaskRepository creates a mock repository that asserts expectations on cleanup.
func NewTaskRepository(t *testing.T) *TaskRepository {
	t.Helper()
	repo := &TaskRepository{}
	repo.Test(t)
	return repo
}

func (m *TaskRepository) Create(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *TaskRepository) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	args := m.Called(ctx, id)
	task, _ := args.Get(0).(*domain.Task)
	return task, args.Error(1)
}

func (m *TaskRepository) List(ctx context.Context) ([]*domain.Task, error) {
	args := m.Called(ctx)
	tasks, _ := args.Get(0).([]*domain.Task)
	return tasks, args.Error(1)
}

func (m *TaskRepository) Update(ctx context.Context, task *domain.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *TaskRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *TaskRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}
