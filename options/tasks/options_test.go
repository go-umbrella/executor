package tasks

import (
	"github.com/go-umbrella/executor/options"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArgs(t *testing.T) {
	testCases := []struct {
		name     string
		args     []interface{}
		expected Option
	}{
		{
			name:     "should_return_nil_args",
			args:     nil,
			expected: options.NewOption("task.args", nil),
		},
		{
			name:     "should_return_args",
			args:     []interface{}{1, "2", 3},
			expected: options.NewOption("task.args", []interface{}{1, "2", 3}),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actual := Args(testCase.args...)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
