package taskscheduler

import "context"

type Task interface {
	ID() string
	Execute(ctx context.Context) (any, error)
}
