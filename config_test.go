package executor

import (
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
			name:   "should_normalize_concurrency",
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
				Concurrency:         uint64(runtime.NumCPU() * 4),
				QueueSize:           16,
				EagerInitialization: true,
				BlockOnFullQueue:    true,
			},
			expected: Config{
				Concurrency:         uint64(runtime.NumCPU() * 4),
				QueueSize:           16,
				EagerInitialization: true,
				BlockOnFullQueue:    true,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.config.normalize()
			assert.Equal(t, testCase.expected, testCase.config)
		})
	}
}
