package application

import (
	"context"
	"errors"
	"gomock_test/domain"
)

type TaskApplicationService struct {
	taskRepository domain.TaskRepository
}

func NewTaskApplicationService(taskRepository domain.TaskRepository) *TaskApplicationService {
	return &TaskApplicationService{taskRepository: taskRepository}
}

func (s *TaskApplicationService) GetTask(ctx context.Context, id string) (*domain.Task, error) {
	t, err := s.taskRepository.Get(ctx, id)
	if err != nil {
		errors.New("get task failed.")
		return nil, err
	}

	return t, nil
}
