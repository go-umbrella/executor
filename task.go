package executor

type Task func(ctx TaskContext) (interface{}, error)
