package executor

import (
	"context"
	stderrors "errors"
	"github.com/go-umbrella/executor/errors"
	"github.com/go-umbrella/executor/options/tasks"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

func TestExecution_Wait(t *testing.T) {
	delay := 100 * time.Millisecond
	result := atomic.Bool{}
	task := func(ctx context.Context, args []interface{}) (interface{}, error) {
		time.Sleep(delay)
		result.Store(true)
		return nil, nil
	}

	execution := newExecution(context.Background(), task)
	start := time.Now()
	go execution.start()

	execution.Wait()
	assert.True(t, result.Load())
	assert.InDelta(t, delay, time.Now().Sub(start), float64(10*time.Millisecond))
}

func TestExecution_Get(t *testing.T) {
	testCases := []struct {
		name      string
		execution *execution
		result    interface{}
		err       error
		panic     bool
	}{
		{
			name: "should_return_result_and_error_null",
			execution: newExecution(context.Background(), func(ctx context.Context, args []interface{}) (interface{}, error) {
				return nil, nil
			}),
			result: nil,
			err:    nil,
		},
		{
			name: "should_return_error",
			execution: newExecution(context.Background(), func(ctx context.Context, args []interface{}) (interface{}, error) {
				return nil, stderrors.New("my_error")
			}),
			result: nil,
			err:    stderrors.New("my_error"),
		},
		{
			name: "should_return_result",
			execution: newExecution(context.Background(), func(ctx context.Context, args []interface{}) (interface{}, error) {
				return "my_result", nil
			}),
			result: "my_result",
			err:    nil,
		},
		{
			name: "should_return_both",
			execution: newExecution(context.Background(), func(ctx context.Context, args []interface{}) (interface{}, error) {
				return "my_result", stderrors.New("my_error")
			}),
			result: "my_result",
			err:    stderrors.New("my_error"),
		},
		{
			name: "should_execute_with_args",
			execution: newExecution(context.Background(), func(ctx context.Context, args []interface{}) (interface{}, error) {
				return args[0], args[1].(error)
			}, tasks.Args("my_result", stderrors.New("my_error"))),
			result: "my_result",
			err:    stderrors.New("my_error"),
		},
		{
			name: "should_recover_from_panic",
			execution: newExecution(context.Background(), func(ctx context.Context, args []interface{}) (interface{}, error) {
				panic("panic_message")
			}),
			result: nil,
			err:    errors.NewRecoveredPanicError(1, "panic_message"),
			panic:  true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			go testCase.execution.start()
			result, err := testCase.execution.Wait().Get()
			assert.Equal(t, testCase.result, result)

			if testCase.panic {
				assert.Equal(t, testCase.err.(*errors.RecoveredPanicError).Value(), err.(*errors.RecoveredPanicError).Value())
			} else {
				assert.Equal(t, testCase.err, err)
			}
		})
	}
}

func TestExecution_Done(t *testing.T) {
	delay := 100 * time.Millisecond
	result := atomic.Bool{}
	task := func(ctx context.Context, args []interface{}) (interface{}, error) {
		time.Sleep(delay)
		result.Store(true)
		return nil, nil
	}

	execution := newExecution(context.Background(), task)
	start := time.Now()
	go execution.start()

	select {
	case <-execution.Done():
		// done successfully
		assert.True(t, result.Load())
		assert.InDelta(t, delay, time.Now().Sub(start), float64(10*time.Millisecond))
	case <-time.After(delay + 25*time.Millisecond):
		assert.FailNow(t, "task too slow")
	}
}
