package executor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"sync/atomic"
	"testing"
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
