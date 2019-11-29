package domain

import "context"

type TaskRepository interface {
	Get(ctx context.Context, id string) (*Task, error)
}
