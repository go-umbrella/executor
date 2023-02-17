package executor

type Status string

const (
	RunningStatus     Status = "RUNNING"
	TerminatingStatus Status = "TERMINATING"
	TerminatedStatus  Status = "TERMINATED"
)

func (s Status) Running() bool {
	return s == RunningStatus
}

func (s Status) Terminating() bool {
	return s == TerminatingStatus
}

func (s Status) Terminated() bool {
	return s == TerminatedStatus
}

func (s Status) String() string {
	return string(s)
}
