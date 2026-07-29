package task

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/domain"
)

type mockTaskRepository struct {
	tasks map[string]*domain.Task
}

func newMockTaskRepository() *mockTaskRepository {
	return &mockTaskRepository{tasks: make(map[string]*domain.Task)}
}

func (m *mockTaskRepository) Create(_ context.Context, task *domain.Task) error {
	if task.ID == "" {
		task.ID = "11111111-1111-1111-1111-111111111111"
	}
	now := time.Now().UTC()
	task.CreatedAt = now
	task.UpdatedAt = now
	copyTask := *task
	m.tasks[task.ID] = &copyTask
	return nil
}

func (m *mockTaskRepository) GetByID(_ context.Context, id string) (*domain.Task, error) {
	task, ok := m.tasks[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	copyTask := *task
	return &copyTask, nil
}

func (m *mockTaskRepository) List(_ context.Context) ([]*domain.Task, error) {
	result := make([]*domain.Task, 0, len(m.tasks))
	for _, task := range m.tasks {
		copyTask := *task
		result = append(result, &copyTask)
	}
	return result, nil
}

func (m *mockTaskRepository) Update(_ context.Context, task *domain.Task) error {
	if _, ok := m.tasks[task.ID]; !ok {
		return domain.ErrNotFound
	}
	copyTask := *task
	copyTask.UpdatedAt = time.Now().UTC()
	m.tasks[task.ID] = &copyTask
	return nil
}

func (m *mockTaskRepository) Delete(_ context.Context, id string) error {
	if _, ok := m.tasks[id]; !ok {
		return domain.ErrNotFound
	}
	delete(m.tasks, id)
	return nil
}

func (m *mockTaskRepository) Count(_ context.Context) (int64, error) {
	return int64(len(m.tasks)), nil
}

func TestService_Create(t *testing.T) {
	svc := NewService(newMockTaskRepository())

	task, err := svc.Create(context.Background(), CreateInput{Title: "Write tests"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if task.Title != "Write tests" {
		t.Errorf("Title = %q, want Write tests", task.Title)
	}
	if task.Status != domain.StatusTodo {
		t.Errorf("Status = %q, want todo", task.Status)
	}
}

func TestService_Create_Validation(t *testing.T) {
	svc := NewService(newMockTaskRepository())

	_, err := svc.Create(context.Background(), CreateInput{Title: "   "})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("Create() error = %v, want ErrInvalidInput", err)
	}
}

func TestService_Update(t *testing.T) {
	repo := newMockTaskRepository()
	svc := NewService(repo)

	created, err := svc.Create(context.Background(), CreateInput{Title: "Initial"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	title := "Updated"
	updated, err := svc.Update(context.Background(), created.ID, UpdateInput{Title: &title})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if updated.Title != "Updated" {
		t.Errorf("Title = %q, want Updated", updated.Title)
	}
}

func TestService_Delete(t *testing.T) {
	repo := newMockTaskRepository()
	svc := NewService(repo)

	created, err := svc.Create(context.Background(), CreateInput{Title: "Remove me"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := svc.Delete(context.Background(), created.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err = svc.GetByID(context.Background(), created.ID)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("GetByID() error = %v, want ErrNotFound", err)
	}
}
