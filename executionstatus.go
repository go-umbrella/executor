package executor

type ExecutionStatus string

const (
	ExecutionWaitingStatus   ExecutionStatus = "WAITING"
	ExecutionRejectedStatus  ExecutionStatus = "REJECTED"
	ExecutionRunningStatus   ExecutionStatus = "RUNNING"
	ExecutionCancelledStatus ExecutionStatus = "CANCELLED"
	ExecutionDoneStatus      ExecutionStatus = "DONE"
)

func (s ExecutionStatus) Waiting() bool {
	return s == ExecutionWaitingStatus
}

func (s ExecutionStatus) Rejected() bool {
	return s == ExecutionRejectedStatus
}

func (s ExecutionStatus) Running() bool {
	return s == ExecutionRunningStatus
}

func (s ExecutionStatus) Cancelled() bool {
	return s == ExecutionCancelledStatus
}

func (s ExecutionStatus) Done() bool {
	return s == ExecutionDoneStatus
}

func (s ExecutionStatus) String() string {
	return string(s)
}
