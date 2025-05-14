package taskscheduler

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"
	"time"
)

func TestTaskScheduler(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// 10 goroutine, 100 task queue, 10 task per second
	scheduler := NewScheduler(10, 15, 10)

	scheduler.Start()

	// select results
	go func() {
		for result := range scheduler.resultChan {
			fmt.Printf("Task %s completed: %v\n", result.TaskID, result.Output)
		}
	}()

	// select errors
	go func() {
		for err := range scheduler.errorChan {
			fmt.Printf("Task %s failed: %v (attempts: %d)\n", err.TaskID, err.Err, err.Attempts)
		}
	}()

	for i := range 20 {
		actionFn := func(ctx context.Context) (any, error) {
			// Simulate a task that takes some time to complete
			// And there is a certain probability of failure
			// 10% chance of failure
			if rand.IntN(10) == 0 {
				return nil, fmt.Errorf("task failed")
			} else {
				return fmt.Sprintf("task-%d completed", i), nil
			}
		}

		retryTask := WithRetryTask{
			task: &SimpleTask{
				id:     fmt.Sprintf("task-%d", i),
				action: actionFn,
			},
			max:     3,
			backoff: 200 * time.Millisecond,
		}

		scheduler.Submit(&retryTask)
	}

	time.Sleep(time.Second * 5)
	scheduler.Stop()
}
