package errors

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

type RecoveredPanicError struct {
	value   interface{}
	callers []uintptr
	stack   []byte
}

func NewRecoveredPanicError(skip int, value interface{}) *RecoveredPanicError {
	var callers [64]uintptr
	n := runtime.Callers(skip+1, callers[:])
	return &RecoveredPanicError{
		value:   value,
		callers: callers[:n],
		stack:   debug.Stack(),
	}
}

func (e *RecoveredPanicError) Value() interface{} {
	return e.value
}

func (e *RecoveredPanicError) Callers() []uintptr {
	return e.callers
}

func (e *RecoveredPanicError) Stack() []byte {
	return e.stack
}

func (e *RecoveredPanicError) Error() string {
	return e.String()
}

func (e *RecoveredPanicError) Unwrap() error {
	if err, ok := e.value.(error); ok {
		return err
	}

	return nil
}

func (e *RecoveredPanicError) String() string {
	return fmt.Sprintf("panic: %v\nstacktrace:\n%s\n", e.value, e.stack)
}
