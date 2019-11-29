package application

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"gomock_test/domain"
	"testing"
)

func TestTaskApplicationService_GetTask(t *testing.T) {
	ctrl := gomock.NewController(t)

	taskRepository := domain.NewMockTaskRepository(ctrl)

	s := NewTaskApplicationService(taskRepository)

	t.Run("ok", func(t *testing.T) {
		task := &domain.Task{ID: "TEST_TASK"}
		taskRepository.EXPECT().Get(gomock.Any(), "TEST_TASK").Return(task, nil)

		ctx := context.Background()
		id := "TEST_TASK"

		task, err := s.GetTask(ctx, id)

		if err != nil {
			t.Fatal("expected no error")
		}
	})

	t.Run("fail", func(t *testing.T) {
		taskRepository.EXPECT().Get(gomock.Any(), "TEST_TASK").Return(nil, errors.New("something failed"))

		ctx := context.Background()
		id := "TEST_TASK"

		task, err := s.GetTask(ctx, id)

		if task != nil || err == nil {
			t.Fatal("expected error")
		}
	})
}
