package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
// If m <= 0 then all of errors will be ignored.
func Run(tasks []Task, n int, m int) error {
	wg := &sync.WaitGroup{}
	var errorsCount, tasksRunning int32

	if n <= 0 {
		n = 1
	}

	i := 0
	for i < len(tasks) {
		if atomic.LoadInt32(&tasksRunning) < int32(n) { // start new task
			wg.Add(1)
			atomic.AddInt32(&tasksRunning, 1)
			go func(tIndex int) {
				defer wg.Done()
				if err := tasks[tIndex](); err != nil {
					atomic.AddInt32(&errorsCount, 1)
				}
				atomic.AddInt32(&tasksRunning, -1) // one task completed
			}(i)
			if m > 0 && atomic.LoadInt32(&errorsCount) >= int32(m) { // wait all running tasks and terminate, errors count could be greather than m
				wg.Wait()
				return ErrErrorsLimitExceeded
			}
			i++
		}
	}

	wg.Wait()

	if m >= len(tasks) && int(errorsCount) == len(tasks) { // when m >= len(tasks)
		return ErrErrorsLimitExceeded
	}
	return nil
}
