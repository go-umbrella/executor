package executor

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type expected struct {
		config             Config
		workerCount        uint64
		workerRunningCount uint64
	}

	testCases := []struct {
		name     string
		executor Executor
		expected expected
	}{
		{
			name: "should_eager_initialize_workers",
			executor: New("my-beautiful-executor", Config{
				Concurrency:         4,
				QueueSize:           16,
				EagerInitialization: true,
			}),
			expected: expected{
				config: Config{
					Concurrency:         4,
					QueueSize:           16,
					EagerInitialization: true,
				},
				workerCount:        4,
				workerRunningCount: 0,
			},
		},
		{
			name: "should_eager_initialize_default_workers",
			executor: New("my-beautiful-executor", Config{
				Concurrency:         0,
				QueueSize:           16,
				EagerInitialization: true,
			}),
			expected: expected{
				config: Config{
					Concurrency:         uint64(runtime.NumCPU()),
					QueueSize:           16,
					EagerInitialization: true,
				},
				workerCount:        uint64(runtime.NumCPU()),
				workerRunningCount: 0,
			},
		},
		{
			name: "should_lazily_initialize_workers",
			executor: New("my-beautiful-executor", Config{
				Concurrency:         4,
				QueueSize:           16,
				EagerInitialization: false,
			}),
			expected: expected{
				config: Config{
					Concurrency:         4,
					QueueSize:           16,
					EagerInitialization: false,
				},
				workerCount:        0,
				workerRunningCount: 0,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			executor := testCase.executor.(*executor)
			assert.Equal(t, testCase.expected.config, executor.config)
			assert.Equal(t, testCase.expected.workerCount, atomic.LoadUint64(&executor.workerCount))
			assert.Equal(t, testCase.expected.workerRunningCount, atomic.LoadUint64(&executor.workerRunningCount))
		})
	}
}

func TestExecutor_WorkerLifeCycle(t *testing.T) {
	concurrency := uint64(4)
	executor := New("test", Config{
		Concurrency:         concurrency,
		QueueSize:           4,
		EagerInitialization: true,
	}).(*executor)

	assert.Equal(t, concurrency, atomic.LoadUint64(&executor.workerCount))
	assert.Equal(t, uint64(0), atomic.LoadUint64(&executor.workerRunningCount))

	taskCounter := uint64(0)
	rounds := 8
	task := func() {
		time.Sleep(100 * time.Millisecond)
		atomic.AddUint64(&taskCounter, 1)
	}

	for i := 0; i < rounds; i++ {
		executor.taskQueue <- task
		time.Sleep(10 * time.Millisecond)
		assert.Equal(t, concurrency, atomic.LoadUint64(&executor.workerCount))
		assert.Equal(t, uint64(math.Min(float64(i+1), float64(concurrency))), atomic.LoadUint64(&executor.workerRunningCount))
	}

	time.Sleep(200 * time.Millisecond)
	assert.Equal(t, uint64(rounds), atomic.LoadUint64(&taskCounter))
	assert.Equal(t, concurrency, atomic.LoadUint64(&executor.workerCount))
	assert.Equal(t, uint64(0), atomic.LoadUint64(&executor.workerRunningCount))

	for i := uint64(0); i < concurrency; i++ {
		executor.workerStopSignal <- struct{}{}
		time.Sleep(10 * time.Millisecond)
		assert.Equal(t, concurrency-(i+1), atomic.LoadUint64(&executor.workerCount))
		assert.Equal(t, uint64(0), atomic.LoadUint64(&executor.workerRunningCount))
	}

	executor.workerWG.Wait()
}

func TestExecutor_ShouldNotCreateNewWorkerWhenMaxWorkersHaveBeenCreated(t *testing.T) {
	concurrency := uint64(4)
	executor := New("test", Config{
		Concurrency:         concurrency,
		QueueSize:           4,
		EagerInitialization: false,
	}).(*executor)

	executions := make([]Execution, 0)
	for i := uint64(0); i < concurrency+4; i++ {
		executions = append(executions, executor.Go(context.Background(), func(ctx TaskContext) (interface{}, error) {
			time.Sleep(100 * time.Millisecond)
			return true, nil
		}))

		time.Sleep(15 * time.Millisecond)
		assert.Equal(t, uint64(math.Min(float64(i+1), float64(concurrency))), atomic.LoadUint64(&executor.workerCount))
		assert.Equal(t, uint64(math.Min(float64(i+1), float64(concurrency))), atomic.LoadUint64(&executor.workerRunningCount))
	}

	for _, execution := range executions {
		result, err := execution.Wait().Get()
		assert.Equal(t, true, result)
		assert.Nil(t, err)
	}
}

func TestExecutor_ShouldCreateWorkerWhenQueueIsFullAndHandleNoIdleWorkerAndQueueIsFull(t *testing.T) {
	concurrency := uint64(1)
	executor := New("test", Config{
		Concurrency:         concurrency,
		QueueSize:           0,
		EagerInitialization: false,
	}).(*executor)

	executions := make([]Execution, 0)
	for i := uint64(0); i < concurrency+4; i++ {
		executions = append(executions, executor.Go(context.Background(), func(ctx TaskContext) (interface{}, error) {
			time.Sleep(100 * time.Millisecond)
			return true, nil
		}))

		time.Sleep(15 * time.Millisecond)
		assert.Equal(t, uint64(math.Min(float64(i+1), float64(concurrency))), atomic.LoadUint64(&executor.workerCount))
		assert.Equal(t, uint64(math.Min(float64(i+1), float64(concurrency))), atomic.LoadUint64(&executor.workerRunningCount))
	}

	for i, execution := range executions {
		result, err := execution.Wait().Get()
		if i == 0 {
			assert.Equal(t, true, result)
			assert.Nil(t, err)
		} else {
			assert.Nil(t, result)
			assert.Equal(t, errors.New("execution_rejected"), err)
		}
	}
}

func TestExecutor_Name(t *testing.T) {
	testCases := []struct {
		name         string
		executorName string
		expectedName string
	}{
		{
			name:         "should_return_executor_pointer",
			executorName: "",
			expectedName: "",
		},
		{
			name:         "should_return_name",
			executorName: "my-beautiful-executor",
			expectedName: "my-beautiful-executor",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			executor := New(testCase.executorName, Config{})

			if testCase.executorName == "" {
				testCase.expectedName = fmt.Sprintf("%p", executor)
			}

			assert.Equal(t, testCase.expectedName, executor.Name())
		})
	}
}
