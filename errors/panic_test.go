package errors

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewRecoveredPanicError(t *testing.T) {
	err := NewRecoveredPanicError(1, "panic error")
	require.NotNil(t, err)
	assert.Equal(t, "panic error", err.Value())
	assert.NotEmpty(t, err.Callers())
	assert.NotEmpty(t, err.Stack())
	assert.NotNil(t, err.Error())
	assert.Equal(t, err.Error(), err.String())
	assert.Nil(t, err.Unwrap())

	err = NewRecoveredPanicError(1, errors.New("panic error"))
	require.NotNil(t, err)
	assert.Equal(t, errors.New("panic error"), err.Value())
	assert.NotEmpty(t, err.Callers())
	assert.NotEmpty(t, err.Stack())
	assert.NotNil(t, err.Error())
	assert.Equal(t, err.Error(), err.String())
	assert.Equal(t, errors.New("panic error"), err.Unwrap())
}
