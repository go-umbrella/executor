package executor

import (
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
			executor: New(Config{
				Name:                "my-beautiful-executor",
				Concurrency:         4,
				QueueSize:           16,
				EagerInitialization: true,
			}),
			expected: expected{
				config: Config{
					Name:                "my-beautiful-executor",
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
			executor: New(Config{
				Name:                "my-beautiful-executor",
				Concurrency:         0,
				QueueSize:           16,
				EagerInitialization: true,
			}),
			expected: expected{
				config: Config{
					Name:                "my-beautiful-executor",
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
			executor: New(Config{
				Name:                "my-beautiful-executor",
				Concurrency:         4,
				QueueSize:           16,
				EagerInitialization: false,
			}),
			expected: expected{
				config: Config{
					Name:                "my-beautiful-executor",
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
