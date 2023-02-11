package executor

import (
	"context"
	"github.com/go-umbrella/executor/errors"
	"github.com/go-umbrella/executor/options/tasks"
	"sync/atomic"
)

type (
	Execution interface {
		Wait() Execution
		Get() (interface{}, error)
		Done() <-chan struct{}
	}

	execution struct {
		ctx    context.Context
		task   Task
		args   []interface{}
		result atomic.Value
		error  atomic.Value
		done   chan struct{}
	}
)

func newExecution(ctx context.Context, task Task, opts ...tasks.Option) *execution {
	execution := new(execution)
	execution.ctx = ctx
	execution.task = task
	execution.done = make(chan struct{})
	execution.processOptions(opts...)
	return execution
}

func (e *execution) Wait() Execution {
	<-e.done
	return e
}

func (e *execution) Get() (interface{}, error) {
	if err := e.error.Load(); err != nil {
		return e.result.Load(), err.(error)
	}

	return e.result.Load(), nil
}

func (e *execution) Done() <-chan struct{} {
	return e.done
}

func (e *execution) start() {
	defer func() {
		if value := recover(); value != nil {
			e.setResult(nil, errors.NewRecoveredPanicError(1, value))
		}

		e.stop()
	}()

	e.setResult(e.task(e.ctx, e.args))
}

func (e *execution) setResult(result interface{}, error error) {
	if result != nil {
		e.result.Store(result)
	}

	if error != nil {
		e.error.Store(error)
	}
}

func (e *execution) stop() {
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
