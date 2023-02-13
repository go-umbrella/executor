package executor

import (
	"context"
	stderrors "errors"
	"github.com/go-umbrella/executor/errors"
	"github.com/go-umbrella/executor/options/tasks"
	"sync/atomic"
	"time"
)

type (
	Execution interface {
		Wait() Execution
		WaitCtx(ctx context.Context) bool
		WaitDeadline(deadline time.Time) bool
		WaitTimeout(duration time.Duration) bool
		Cancel() bool
		Get() (interface{}, error)
		Status() ExecutionStatus
		Done() <-chan struct{}
	}

	execution struct {
		ctx     context.Context
		cancel  context.CancelCauseFunc
		task    Task
		args    []interface{}
		result  atomic.Value
		error   atomic.Value
		status  atomic.Value
		stopped int32
		done    chan struct{}
	}
)

var (
	executionCancelledError = stderrors.New("execution_cancelled")
	executionRejectedError  = stderrors.New("execution_rejected")
)

func newExecution(ctx context.Context, task Task, opts ...tasks.Option) *execution {
	execution := new(execution)
	execution.ctx, execution.cancel = context.WithCancelCause(ctx)
	execution.task = task
	execution.status.Store(ExecutionWaitingStatus)
	execution.done = make(chan struct{})
	execution.processOptions(opts...)
	return execution
}

func (e *execution) Wait() Execution {
	<-e.done
	return e
}

func (e *execution) WaitCtx(ctx context.Context) bool {
	select {
	case <-e.done:
		return true
	case <-ctx.Done():
		return false
	}
}

func (e *execution) WaitDeadline(deadline time.Time) bool {
	return e.WaitTimeout(time.Until(deadline))
}

func (e *execution) WaitTimeout(duration time.Duration) bool {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-e.done:
		return true
	case <-timer.C:
		return false
	}
}

func (e *execution) Cancel() bool {
	if !atomic.CompareAndSwapInt32(&e.stopped, 0, 1) {
		return false
	}

	e.status.Store(ExecutionCancelledStatus)
	e.error.Store(executionCancelledError)
	close(e.done)
	return true
}

func (e *execution) Get() (interface{}, error) {
	return e.result.Load(), e.loadError()
}

func (e *execution) Status() ExecutionStatus {
	return e.status.Load().(ExecutionStatus)
}

func (e *execution) Done() <-chan struct{} {
	return e.done
}

func (e *execution) start() {
	var (
		result interface{}
		err    error
	)

	defer func() {
		if value := recover(); value != nil {
			result, err = nil, errors.NewRecoveredPanicError(1, value)
		}

		e.setResult(result, err)
	}()

	result, err = e.task(e.ctx, e.args)
}

func (e *execution) setResult(result interface{}, error error) bool {
	if !atomic.CompareAndSwapInt32(&e.stopped, 0, 1) {
		return false
	}

	if result != nil {
		e.result.Store(result)
	}

	if error != nil {
		e.error.Store(error)
	}

	close(e.done)
	return true
}

func (e *execution) reject() {
	if !atomic.CompareAndSwapInt32(&e.stopped, 0, 1) {
		return
	}

	e.status.Store(ExecutionWaitingStatus)
	e.error.Store(executionRejectedError)
	close(e.done)
}

func (e *execution) processOptions(opts ...tasks.Option) {
	for _, opt := range opts {
		switch opt.Type() {
		case tasks.ArgsType:
			if opt.Value() != nil {
				e.args = opt.Value().([]interface{})
			}
		}
	}
}

func (e *execution) loadError() error {
	err := e.error.Load()
	if err != nil {
		return err.(error)
	}

	return nil
}
