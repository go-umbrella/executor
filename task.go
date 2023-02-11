package executor

import "context"

type Task func(ctx context.Context, args []interface{}) (interface{}, error)
