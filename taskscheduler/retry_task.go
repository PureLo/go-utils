package taskscheduler

import (
	"context"
	"fmt"
	"time"
)

type WithRetryTask struct {
	task    Task
	max     int           // max retry times
	backoff time.Duration // backoff time
}

func (r *WithRetryTask) ID() string {
	return r.task.ID()
}

func (r *WithRetryTask) Execute(ctx context.Context) (any, error) {
	var lastErr error
	for i := range r.max {
		if i > 0 {
			select {
			case <-time.After(r.backoff):
				// continue
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		output, err := r.task.Execute(ctx)
		if err == nil {
			return output, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("after %d attempts: %w", r.max, lastErr)
}
