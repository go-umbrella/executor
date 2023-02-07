package executor

type (
	Executor interface {
	}

	executor struct {
	}
)

func New() Executor {
	return new(executor)
}
