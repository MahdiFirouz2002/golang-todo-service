package domain_test

import (
	"testing"

	"github.com/MahdiFirouz2002/golang-todo-service/internal/domain"
)

func TestTaskStatus_Valid(t *testing.T) {
	tests := []struct {
		name   string
		status domain.TaskStatus
		want   bool
	}{
		{name: "todo", status: domain.StatusTodo, want: true},
		{name: "in_progress", status: domain.StatusInProgress, want: true},
		{name: "done", status: domain.StatusDone, want: true},
		{name: "empty", status: "", want: false},
		{name: "unknown", status: "archived", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
