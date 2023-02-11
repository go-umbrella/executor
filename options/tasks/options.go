package tasks

import "github.com/go-umbrella/executor/options"

type Option interface {
	options.Option
}

func Args(args ...interface{}) Option {
	if args == nil {
		return options.NewOption(ArgsType, nil)
	}

	return options.NewOption(ArgsType, args)
}
