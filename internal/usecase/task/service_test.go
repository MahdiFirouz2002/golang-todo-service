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

func (m *mockTaskRepository) List(_ context.Context, filter domain.ListFilter) (*domain.ListResult, error) {
	items := make([]*domain.Task, 0, len(m.tasks))
	for _, task := range m.tasks {
		if filter.Status != nil && task.Status != *filter.Status {
			continue
		}
		if filter.Assignee != nil && task.Assignee != *filter.Assignee {
			continue
		}
		copyTask := *task
		items = append(items, &copyTask)
	}

	total := int64(len(items))
	page := filter.Page
	if page < 1 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	start := (page - 1) * pageSize
	if start >= len(items) {
		items = []*domain.Task{}
	} else {
		end := start + pageSize
		if end > len(items) {
			end = len(items)
		}
		items = items[start:end]
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}
	if total == 0 {
		totalPages = 0
	}

	return &domain.ListResult{
		Items:      items,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
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

func TestService_List_InvalidStatus(t *testing.T) {
	repo := newMockTaskRepository()
	svc := NewService(repo)

	status := domain.TaskStatus("bad")
	_, err := svc.List(context.Background(), ListInput{Status: &status})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("List() error = %v, want ErrInvalidInput", err)
	}
}

func TestService_List_PaginationDefaults(t *testing.T) {
	repo := newMockTaskRepository()
	svc := NewService(repo)

	result, err := svc.List(context.Background(), ListInput{})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if result.Page != 1 || result.PageSize != 20 {
		t.Fatalf("page = %d pageSize = %d, want 1 and 20", result.Page, result.PageSize)
	}
}

func BenchmarkService_Create(b *testing.B) {
	repo := newMockTaskRepository()
	svc := NewService(repo)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.Create(ctx, CreateInput{Title: "Benchmark task"})
	}
}

func BenchmarkService_List(b *testing.B) {
	repo := newMockTaskRepository()
	svc := NewService(repo)
	ctx := context.Background()

	_, _ = svc.Create(ctx, CreateInput{Title: "Seed"})
	status := domain.StatusTodo

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = svc.List(ctx, ListInput{Status: &status})
	}
}
