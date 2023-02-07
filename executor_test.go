package executor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name          string
		expectedValue Executor
	}{
		{
			name:          "should_return_new_executor",
			expectedValue: new(executor),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expectedValue, New())
		})
	}
}
