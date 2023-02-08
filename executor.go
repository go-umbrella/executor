package executor

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type (
	Executor interface {
	}

	executor struct {
		config             Config
		taskQueue          chan func()
		workerWG           sync.WaitGroup
		workerStopSignal   chan struct{}
		workerCount        uint64
		workerRunningCount uint64
	}
)

func New(config Config) Executor {
	config.normalize()
	return (&executor{
		config:           config,
		taskQueue:        make(chan func(), config.QueueSize),
		workerStopSignal: make(chan struct{}, runtime.NumCPU()),
	}).initialize()
}

func (e *executor) initialize() *executor {
	e.initializeWorkers()
	return e
}

func (e *executor) initializeWorkers() {
	if !e.config.EagerInitialization {
		return
	}

	for i := uint64(0); i < e.config.Concurrency; i++ {
		e.newWorker()
	}
}

func (e *executor) newWorker() bool {
	if atomic.LoadUint64(&e.workerCount) >= e.config.Concurrency {
		return false
	}

	e.workerWG.Add(1)
	atomic.AddUint64(&e.workerCount, 1)

	go func() {
		defer func() {
			atomic.AddUint64(&e.workerCount, ^uint64(0))
			e.workerWG.Done()
		}()

		for {
			select {
			case task := <-e.taskQueue:
				atomic.AddUint64(&e.workerRunningCount, 1)
				task()
				atomic.AddUint64(&e.workerRunningCount, ^uint64(0))
			case <-e.workerStopSignal:
				return
			}
		}
	}()

	return true
}
