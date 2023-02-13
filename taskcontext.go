package executor

import (
	"context"
	"time"
)

type (
	TaskContext interface {
		Args() []interface{}
		context.Context
	}

	taskContext struct {
		ctx    context.Context
		cancel context.CancelCauseFunc
		args   []interface{}
	}
)

func newTaskContext(ctx context.Context, cancel context.CancelCauseFunc, args []interface{}) *taskContext {
	return &taskContext{
		ctx:    ctx,
		cancel: cancel,
		args:   args,
	}
}

func (c *taskContext) Args() []interface{} {
	return c.args
}

func (c *taskContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *taskContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *taskContext) Err() error {
	return c.ctx.Err()
}

func (c *taskContext) Value(key any) any {
	return c.ctx.Value(key)
}
