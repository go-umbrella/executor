package executor

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTaskContext_Args(t *testing.T) {
	args := []interface{}{1, 2, 3}
	ctx := newTaskContext(context.Background(), args)
	assert.Equal(t, args, ctx.Args())
}

func TestTaskContext_Deadline(t *testing.T) {
	deadlineCtx, cancel := context.WithDeadline(context.Background(), time.Now().Add(100*time.Millisecond))
	defer cancel()

	ctx := newTaskContext(deadlineCtx, nil)

	deadline, ok := ctx.Deadline()
	deadlineCtxDeadline, deadlineCtxOk := deadlineCtx.Deadline()
	assert.Equal(t, deadlineCtxDeadline, deadline)
	assert.Equal(t, ok, deadlineCtxOk)
}

func TestTaskContext_DoneAndErr(t *testing.T) {
	cancelCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctx := newTaskContext(cancelCtx, nil)

	select {
	case <-ctx.Done():
		assert.FailNow(t, "should not done")
	default:
		// continue
	}

	cancel()

	select {
	case <-ctx.Done():
		// continue
	default:
		assert.FailNow(t, "should not done")
	}

	assert.NotNil(t, cancelCtx.Err())
	assert.Equal(t, cancelCtx.Err(), ctx.Err())
}

func TestTaskContext_Value(t *testing.T) {
	ctx := newTaskContext(context.WithValue(context.Background(), "key", "value"), nil)
	assert.Equal(t, "value", ctx.Value("key"))
}
