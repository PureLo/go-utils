package taskscheduler

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

type Scheduler struct {
	workerNum   int
	rateLimiter *rate.Limiter
	taskQueue   chan Task
	resultChan  chan *Result
	errorChan   chan *Error
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

type Result struct {
	TaskID   string
	Output   any
	Attempts int
}

type Error struct {
	TaskID   string
	Err      error
	Attempts int
}

func NewScheduler(workerNum, queueSize, limitSec int) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())

	return &Scheduler{
		workerNum:   workerNum,
		rateLimiter: rate.NewLimiter(rate.Limit(limitSec), int(limitSec)),
		taskQueue:   make(chan Task, queueSize),
		resultChan:  make(chan *Result, queueSize),
		errorChan:   make(chan *Error, queueSize),
		wg:          sync.WaitGroup{},
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (s *Scheduler) Submit(task Task) {
	s.taskQueue <- task
}

func (s *Scheduler) Start() {
	for range s.workerNum {
		s.wg.Add(1)
		go s.worker()
	}
}

func (s *Scheduler) Stop() {
	s.cancel()
	s.wg.Wait()
	close(s.taskQueue)
	close(s.resultChan)
	close(s.errorChan)
}

func (s *Scheduler) worker() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		case task := <-s.taskQueue:
			if s.rateLimiter != nil {
				if err := s.rateLimiter.Wait(s.ctx); err != nil {
					s.errorChan <- &Error{
						TaskID: task.ID(),
						Err:    err}
					continue
				}
				output, err := task.Execute(s.ctx)
				if err != nil {
					s.errorChan <- &Error{
						TaskID: task.ID(),
						Err:    err}
				} else {
					s.resultChan <- &Result{
						TaskID: task.ID(),
						Output: output}
				}
			}

		}
	}
}
