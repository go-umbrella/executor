package options

import "fmt"

type (
	// Option represents an operation customization.
	Option interface {
		// Type describes the type of the option.
		Type() string

		// Value returns a value used to create this option.
		Value() interface{}

		// String returns a string representation of the option.
		String() string
	}

	option struct {
		_type string
		value interface{}
	}
)

func NewOption(_type string, value interface{}) Option {
	return &option{
		_type: _type,
		value: value,
	}
}

func (o *option) Type() string {
	return o._type
}

func (o *option) Value() interface{} {
	return o.value
}

func (o *option) String() string {
	return fmt.Sprintf("%s: %v", o._type, o.value)
}
