package executor

import (
	"github.com/go-umbrella/executor/internal/env"
	"runtime"
)

func Computation() Executor {
	return computation
}

func IO() Executor {
	return io
}

func Single() Executor {
	return single
}

var (
	computation Executor
	io          Executor
	single      Executor
)

func initExecutors() {
	numCPU := uint64(runtime.NumCPU())

	computation = New("computation-executor", Config{
		Concurrency:         env.Uint64("EXECUTOR_COMPUTATION_CONCURRENCY", numCPU),
		QueueSize:           env.Uint64("EXECUTOR_COMPUTATION_QUEUE_SIZE", numCPU*16),
		EagerInitialization: env.Bool("EXECUTOR_COMPUTATION_EAGER_INITIALIZATION", false),
	})

	io = New("io-executor", Config{
		Concurrency:         env.Uint64("EXECUTOR_IO_CONCURRENCY", numCPU*64),
		QueueSize:           env.Uint64("EXECUTOR_IO_QUEUE_SIZE", numCPU*1024),
		EagerInitialization: env.Bool("EXECUTOR_IO_EAGER_INITIALIZATION", false),
	})

	single = New("single-executor", Config{
		Concurrency:         1,
		QueueSize:           env.Uint64("EXECUTOR_SINGLE_QUEUE_SIZE", 16),
		EagerInitialization: env.Bool("EXECUTOR_SINGLE_EAGER_INITIALIZATION", false),
	})
}
