package options

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOption(t *testing.T) {
	testCases := []struct {
		name     string
		option   Option
		expected *option
	}{
		{
			name:   "should_create_option",
			option: NewOption("task.args", "1"),
			expected: &option{
				_type: "task.args",
				value: "1",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			assert.Equal(t, testCase.expected, testCase.option)
			assert.Equal(t, testCase.expected._type, testCase.option.Type())
			assert.Equal(t, testCase.expected.value, testCase.option.Value())
			assert.Equal(t, fmt.Sprintf("%s: %v", testCase.expected._type, testCase.expected.value), testCase.option.String())
		})
	}
}
