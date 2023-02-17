package env

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestBool(t *testing.T) {
	testCases := []struct {
		name         string
		value        string
		defaultValue bool
		expected     bool
	}{
		{
			name:         "should_return_default_value_false",
			value:        "",
			defaultValue: false,
			expected:     false,
		},
		{
			name:         "should_return_default_value_true",
			value:        "",
			defaultValue: false,
			expected:     false,
		},
		{
			name:         "should_return_default_value_false_on_error",
			value:        "error",
			defaultValue: false,
			expected:     false,
		},
		{
			name:         "should_return_default_value_true_on_error",
			value:        "error",
			defaultValue: false,
			expected:     false,
		},
		{
			name:         "should_return_valid_value_false",
			value:        "false",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "should_return_valid_value_true",
			value:        "true",
			defaultValue: false,
			expected:     true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			require.NoError(t, os.Setenv("ENV_BOOL_TEST", testCase.value))
			assert.Equal(t, testCase.expected, Bool("ENV_BOOL_TEST", testCase.defaultValue))
		})
	}
}

func TestUint64(t *testing.T) {
	testCases := []struct {
		name         string
		value        string
		defaultValue uint64
		expected     uint64
	}{
		{
			name:         "should_return_default_value",
			value:        "",
			defaultValue: 123,
			expected:     123,
		},
		{
			name:         "should_return_default_value_on_error",
			value:        "error",
			defaultValue: 123,
			expected:     123,
		},
		{
			name:         "should_return_valid_value",
			value:        "456",
			defaultValue: 123,
			expected:     456,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			require.NoError(t, os.Setenv("ENV_UINT64_TEST", testCase.value))
			assert.Equal(t, testCase.expected, Uint64("ENV_UINT64_TEST", testCase.defaultValue))
		})
	}
}
