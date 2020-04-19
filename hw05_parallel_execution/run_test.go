package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Run("if were errors in first M tasks, then finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return fmt.Errorf("error from task %d", i)
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		result := Run(tasks, workersCount, maxErrorsCount)

		require.Equal(t, ErrErrorsLimitExceeded, result)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		result := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.Nil(t, result)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	// new case
	t.Run("run tasks with M <= 0: all of errors must be ignored", func(t *testing.T) {
		tasksCount := 20
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return fmt.Errorf("error from task %d", i)
			})
		}

		workersCount := 3
		maxErrorsCount := -1

		result := Run(tasks, workersCount, maxErrorsCount)

		require.Nil(t, result)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	// new case
	t.Run("run tasks with M >= tasksCount: all of tasks must be executed, error must be returned", func(t *testing.T) {
		tasksCount := 30
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return fmt.Errorf("error from task %d", i)
			})
		}

		workersCount := 10
		maxErrorsCount := 40

		result := Run(tasks, workersCount, maxErrorsCount)

		require.Equal(t, ErrErrorsLimitExceeded, result)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	// new case
	t.Run("run tasks with N > len(tasks)", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 100
		maxErrorsCount := 1

		result := Run(tasks, workersCount, maxErrorsCount)
		require.Nil(t, result)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	// new case
	t.Run("run tasks with N <= 0 : sequential execution", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := -10
		maxErrorsCount := 1

		start := time.Now()
		result := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.Nil(t, result)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(sumTime), int64(elapsedTime), "parallel execution?")
	})
}
