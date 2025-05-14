package taskscheduler

import (
	"context"
)

type SimpleTask struct {
	id     string
	action func(ctx context.Context) (any, error)
}

func (t *SimpleTask) ID() string {
	return t.id
}

func (t *SimpleTask) Execute(ctx context.Context) (any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return t.action(ctx)
	}
}
