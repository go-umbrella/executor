package executor

import (
	"github.com/stretchr/testify/assert"
	"os"
	"runtime"
	"testing"
)

func TestComputation(t *testing.T) {
	numCPU := uint64(runtime.NumCPU())

	_ = os.Setenv("EXECUTOR_COMPUTATION_CONCURRENCY", "")
	_ = os.Setenv("EXECUTOR_COMPUTATION_QUEUE_SIZE", "")
	_ = os.Setenv("EXECUTOR_COMPUTATION_EAGER_INITIALIZATION", "")
	_ = os.Setenv("EXECUTOR_COMPUTATION_BLOCK_ON_FULL_QUEUE", "")

	assert.Equal(t, "computation-executor", Computation().(*executor).name)
	assert.Equal(t, Config{
		Concurrency:         numCPU,
		QueueSize:           numCPU * 16,
		EagerInitialization: false,
		BlockOnFullQueue:    false,
	}, Computation().(*executor).config)

	_ = os.Setenv("EXECUTOR_COMPUTATION_CONCURRENCY", "4")
	_ = os.Setenv("EXECUTOR_COMPUTATION_QUEUE_SIZE", "16")
	_ = os.Setenv("EXECUTOR_COMPUTATION_EAGER_INITIALIZATION", "true")
	_ = os.Setenv("EXECUTOR_COMPUTATION_BLOCK_ON_FULL_QUEUE", "true")

	initExecutors()

	assert.Equal(t, "computation-executor", Computation().(*executor).name)
	assert.Equal(t, Config{
		Concurrency:         4,
		QueueSize:           16,
		EagerInitialization: true,
		BlockOnFullQueue:    true,
	}, Computation().(*executor).config)
}

func TestIO(t *testing.T) {
	numCPU := uint64(runtime.NumCPU())

	_ = os.Setenv("EXECUTOR_IO_CONCURRENCY", "")
	_ = os.Setenv("EXECUTOR_IO_QUEUE_SIZE", "")
	_ = os.Setenv("EXECUTOR_IO_EAGER_INITIALIZATION", "")
	_ = os.Setenv("EXECUTOR_IO_BLOCK_ON_FULL_QUEUE", "")

	assert.Equal(t, "io-executor", IO().(*executor).name)
	assert.Equal(t, Config{
		Concurrency:         numCPU * 64,
		QueueSize:           numCPU * 1024,
		EagerInitialization: false,
		BlockOnFullQueue:    false,
	}, IO().(*executor).config)

	_ = os.Setenv("EXECUTOR_IO_CONCURRENCY", "1")
	_ = os.Setenv("EXECUTOR_IO_QUEUE_SIZE", "1")
	_ = os.Setenv("EXECUTOR_IO_EAGER_INITIALIZATION", "true")
	_ = os.Setenv("EXECUTOR_IO_BLOCK_ON_FULL_QUEUE", "true")

	initExecutors()

	assert.Equal(t, "io-executor", IO().(*executor).name)
	assert.Equal(t, Config{
		Concurrency:         1,
		QueueSize:           1,
		EagerInitialization: true,
		BlockOnFullQueue:    true,
	}, IO().(*executor).config)
}

func TestSingle(t *testing.T) {
	_ = os.Setenv("EXECUTOR_SINGLE_QUEUE_SIZE", "")
	_ = os.Setenv("EXECUTOR_SINGLE_EAGER_INITIALIZATION", "")
	_ = os.Setenv("EXECUTOR_SINGLE_BLOCK_ON_FULL_QUEUE", "")

	assert.Equal(t, "single-executor", Single().(*executor).name)
	assert.Equal(t, Config{
		Concurrency:         1,
		QueueSize:           16,
		EagerInitialization: false,
		BlockOnFullQueue:    false,
	}, Single().(*executor).config)

	_ = os.Setenv("EXECUTOR_SINGLE_QUEUE_SIZE", "32")
	_ = os.Setenv("EXECUTOR_SINGLE_EAGER_INITIALIZATION", "true")
	_ = os.Setenv("EXECUTOR_SINGLE_BLOCK_ON_FULL_QUEUE", "true")

	initExecutors()

	assert.Equal(t, "single-executor", Single().(*executor).name)
	assert.Equal(t, Config{
		Concurrency:         1,
		QueueSize:           32,
		EagerInitialization: true,
		BlockOnFullQueue:    true,
	}, Single().(*executor).config)
}
