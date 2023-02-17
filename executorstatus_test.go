package executor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatus(t *testing.T) {
	testCases := []struct {
		name     string
		status   Status
		expected func(status Status) bool
	}{
		{
			name:   "RUNNING",
			status: RunningStatus,
			expected: func(status Status) bool {
				return status.Running()
			},
		},
		{
			name:   "TERMINATING",
			status: TerminatingStatus,
			expected: func(status Status) bool {
				return status.Terminating()
			},
		},
		{
			name:   "TERMINATED",
			status: TerminatedStatus,
			expected: func(status Status) bool {
				return status.Terminated()
			},
		},
	}

	statusList := []Status{
		RunningStatus,
		TerminatingStatus,
		TerminatedStatus,
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.name, testCase.status.String())

			for _, status := range statusList {
				assert.Equal(t, status == testCase.status, testCase.expected(status))
			}
		})
	}
}
