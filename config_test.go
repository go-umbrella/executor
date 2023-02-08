package executor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
)

func TestConfig_Normalize(t *testing.T) {
	testCases := []struct {
		name     string
		config   Config
		expected Config
	}{
		{
			name:   "should_normalize_name_and_concurrency",
			config: Config{},
			expected: Config{
				Concurrency:         uint64(runtime.NumCPU()),
				QueueSize:           0,
				EagerInitialization: false,
			},
		},
		{
			name: "should_keep_valid_configurations",
			config: Config{
				Name:                "my-beautiful-executor",
				Concurrency:         uint64(runtime.NumCPU() * 4),
				QueueSize:           16,
				EagerInitialization: true,
			},
			expected: Config{
				Name:                "my-beautiful-executor",
				Concurrency:         uint64(runtime.NumCPU() * 4),
				QueueSize:           16,
				EagerInitialization: true,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.expected.Name == "" {
				testCase.expected.Name = fmt.Sprintf("%p", &testCase.config)
			}

			testCase.config.normalize()
			assert.Equal(t, testCase.expected, testCase.config)
		})
	}
}
