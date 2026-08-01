package task_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/domain"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/mocks"
	"github.com/MahdiFirouz2002/golang-todo-service/internal/usecase/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_GetByID_InvalidID(t *testing.T) {
	repo := mocks.NewTaskRepository(t)
	svc := task.NewService(repo)

	_, err := svc.GetByID(context.Background(), "not-a-uuid")
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestService_GetByID_NotFound(t *testing.T) {
	repo := mocks.NewTaskRepository(t)
	repo.On("GetByID", mock.Anything, "11111111-1111-1111-1111-111111111111").
		Return((*domain.Task)(nil), domain.ErrNotFound)

	svc := task.NewService(repo)
	_, err := svc.GetByID(context.Background(), "11111111-1111-1111-1111-111111111111")
	require.ErrorIs(t, err, domain.ErrNotFound)
}

func TestService_List(t *testing.T) {
	repo := mocks.NewTaskRepository(t)
	expected := []*domain.Task{{ID: "11111111-1111-1111-1111-111111111111", Title: "A"}}
	repo.On("List", mock.Anything, mock.Anything).Return(&domain.ListResult{Items: expected}, nil)

	svc := task.NewService(repo)
	result, err := svc.List(context.Background(), task.ListInput{})
	require.NoError(t, err)
	assert.Equal(t, expected, result.Items)
}

func TestService_Create_InvalidStatus(t *testing.T) {
	repo := mocks.NewTaskRepository(t)
	svc := task.NewService(repo)

	_, err := svc.Create(context.Background(), task.CreateInput{
		Title:  "Task",
		Status: domain.TaskStatus("invalid"),
	})
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestService_Update_NoFields(t *testing.T) {
	repo := mocks.NewTaskRepository(t)
	svc := task.NewService(repo)

	_, err := svc.Update(context.Background(), "11111111-1111-1111-1111-111111111111", task.UpdateInput{})
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestService_Update_Success(t *testing.T) {
	repo := mocks.NewTaskRepository(t)
	existing := &domain.Task{
		ID:        "11111111-1111-1111-1111-111111111111",
		Title:     "Old",
		Status:    domain.StatusTodo,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	repo.On("GetByID", mock.Anything, existing.ID).Return(existing, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Task")).Return(nil)

	svc := task.NewService(repo)
	title := "New"
	updated, err := svc.Update(context.Background(), existing.ID, task.UpdateInput{Title: &title})
	require.NoError(t, err)
	assert.Equal(t, "New", updated.Title)
}

func TestService_Delete_InvalidID(t *testing.T) {
	repo := mocks.NewTaskRepository(t)
	svc := task.NewService(repo)

	err := svc.Delete(context.Background(), "bad-id")
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}

func TestService_Create_RepoError(t *testing.T) {
	repo := mocks.NewTaskRepository(t)
	repo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Task")).
		Return(errors.New("db down"))

	svc := task.NewService(repo)
	_, err := svc.Create(context.Background(), task.CreateInput{Title: "Task"})
	require.Error(t, err)
}
