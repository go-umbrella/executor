package executor

import (
	"fmt"
	"runtime"
)

// Config used in Executor.
type Config struct {
	Name                string // Executor's Name.
	Concurrency         uint64 // Concurrency Is the number of goroutines to be kept in the Pool. If EagerInitialization is true, will be created at startup.
	QueueSize           uint64 // QueueSize is the Task queue size. If the queue is full, there are no idle goroutines and the maximum number of goroutines has been created, the Task will be rejected.
	EagerInitialization bool   // If EagerInitialization is true, the goroutines (Concurrency) will be created at startup. Otherwise, will be created on demand.
}

func (c *Config) normalize() {
	c.normalizeName()
	c.normalizeConcurrency()
}

func (c *Config) normalizeName() {
	if c.Name == "" {
		c.Name = fmt.Sprintf("%p", c)
	}
}

func (c *Config) normalizeConcurrency() {
	if c.Concurrency == 0 {
		c.Concurrency = uint64(runtime.NumCPU())
	}
}
