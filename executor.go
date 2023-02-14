package executor

import (
	"context"
	"fmt"
	"github.com/go-umbrella/executor/options/tasks"
	"runtime"
	"sync"
	"sync/atomic"
)

type (
	Executor interface {
		Go(ctx context.Context, task Task, opts ...tasks.Option) Execution
		Name() string
	}

	executor struct {
		name               string
		config             Config
		dispatcherQueue    chan *execution
		taskQueue          chan func()
		workerWG           sync.WaitGroup
		workerStopSignal   chan struct{}
		workerCount        uint64
		workerRunningCount uint64
	}
)

func New(name string, config Config) Executor {
	config.normalize()
	return (&executor{
		name:             name,
		config:           config,
		dispatcherQueue:  make(chan *execution, runtime.NumCPU()),
		taskQueue:        make(chan func(), config.QueueSize),
		workerStopSignal: make(chan struct{}, runtime.NumCPU()),
	}).initialize()
}

func (e *executor) Go(ctx context.Context, task Task, opts ...tasks.Option) Execution {
	execution := newExecution(ctx, task, opts...)
	e.dispatcherQueue <- execution
	return execution
}

func (e *executor) Name() string {
	return e.name
}

func (e *executor) initialize() *executor {
	e.normalizeName()
	e.initializeWorkers()
	go e.dispatcher()
	return e
}

func (e *executor) normalizeName() {
	if e.name == "" {
		e.name = fmt.Sprintf("%p", e)
	}
}

func (e *executor) initializeWorkers() {
	if !e.config.EagerInitialization {
		return
	}

	for i := uint64(0); i < e.config.Concurrency; i++ {
		e.newWorker()
	}
}

func (e *executor) dispatcher() {
	for {
		execution := <-e.dispatcherQueue
		taskEnqueued := e.tryEnqueueTask(execution)

		// always try to create a new worker up to maximum concurrency if there are no idle workers
		if !e.hasIdleWorker() && e.canCreateNewWorker() {
			e.newWorker()

			// new worker created, so if the queue is full, wait until the new goroutine read the first task.
			if !taskEnqueued {
				e.enqueueTask(execution)
				continue
			}
		}

		if taskEnqueued {
			continue
		}

		if e.config.BlockOnFullQueue {
			e.enqueueTask(execution)
		} else {
			execution.reject()
		}
	}
}

func (e *executor) newWorker() {
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
}

func (e *executor) enqueueTask(execution *execution) {
	e.taskQueue <- execution.start
}

func (e *executor) tryEnqueueTask(execution *execution) bool {
	select {
	case e.taskQueue <- execution.start:
		return true
	default:
		// queue is full
		return false
	}
}

func (e *executor) canCreateNewWorker() bool {
	return atomic.LoadUint64(&e.workerCount) < e.config.Concurrency
}

func (e *executor) hasIdleWorker() bool {
	return atomic.LoadUint64(&e.workerRunningCount) < atomic.LoadUint64(&e.workerCount)
}
