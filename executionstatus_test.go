package executor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExecutionStatus(t *testing.T) {
	testCases := []struct {
		name            string
		executionStatus ExecutionStatus
		expected        func(status ExecutionStatus) bool
	}{
		{
			name:            "WAITING",
			executionStatus: ExecutionWaitingStatus,
			expected: func(status ExecutionStatus) bool {
				return status.Waiting()
			},
		},
		{
			name:            "REJECTED",
			executionStatus: ExecutionRejectedStatus,
			expected: func(status ExecutionStatus) bool {
				return status.Rejected()
			},
		},
		{
			name:            "RUNNING",
			executionStatus: ExecutionRunningStatus,
			expected: func(status ExecutionStatus) bool {
				return status.Running()
			},
		},
		{
			name:            "CANCELLED",
			executionStatus: ExecutionCancelledStatus,
			expected: func(status ExecutionStatus) bool {
				return status.Cancelled()
			},
		},
		{
			name:            "DONE",
			executionStatus: ExecutionDoneStatus,
			expected: func(status ExecutionStatus) bool {
				return status.Done()
			},
		},
	}

	statusList := []ExecutionStatus{
		ExecutionWaitingStatus,
		ExecutionRejectedStatus,
		ExecutionRunningStatus,
		ExecutionCancelledStatus,
		ExecutionDoneStatus,
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			for _, status := range statusList {
				assert.Equal(t, status == testCase.executionStatus, testCase.expected(status))
			}
		})
	}
}
