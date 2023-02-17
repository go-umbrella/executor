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
		Shutdown() Executor
		Status() Status
		Done() <-chan struct{}
		Name() string
	}

	executor struct {
		name                 string
		config               Config
		dispatcherQueue      chan *execution
		taskQueue            chan func()
		workerWG             sync.WaitGroup
		workerStopSignal     chan struct{}
		workerCount          atomic.Uint64
		workerRunningCount   atomic.Uint64
		status               atomic.Value
		stopped              atomic.Bool
		shutdownSignal       chan struct{}
		dispatcherDoneSignal chan struct{}
		done                 chan struct{}
	}
)

func New(name string, config Config) Executor {
	config.normalize()
	return (&executor{
		name:                 name,
		config:               config,
		dispatcherQueue:      make(chan *execution, runtime.NumCPU()),
		taskQueue:            make(chan func(), config.QueueSize),
		workerStopSignal:     make(chan struct{}, runtime.NumCPU()),
		shutdownSignal:       make(chan struct{}),
		dispatcherDoneSignal: make(chan struct{}),
		done:                 make(chan struct{}),
	}).initialize()
}

func (e *executor) Go(ctx context.Context, task Task, opts ...tasks.Option) Execution {
	execution := newExecution(ctx, task, opts...)

	select {
	case <-e.shutdownSignal:
		execution.reject()
	default:
		e.dispatcherQueue <- execution
	}

	return execution
}

func (e *executor) Shutdown() Executor {
	if !e.stopped.CompareAndSwap(false, true) {
		return e
	}

	e.status.Store(TerminatingStatus)
	close(e.shutdownSignal)
	close(e.dispatcherQueue)
	return e
}

func (e *executor) Status() Status {
	return e.status.Load().(Status)
}

func (e *executor) Done() <-chan struct{} {
	return e.done
}

func (e *executor) Name() string {
	return e.name
}

func (e *executor) initialize() *executor {
	e.status.Store(RunningStatus)
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
	defer func() {
		close(e.taskQueue)
		close(e.dispatcherDoneSignal)
	}()

	for execution := range e.dispatcherQueue {
		e.dispatch(execution)
	}
}

func (e *executor) dispatch(execution *execution) {
	taskEnqueued := e.tryEnqueueTask(execution)

	// always try to create a new worker up to maximum concurrency if there are no idle workers
	if !e.hasIdleWorker() && e.canCreateNewWorker() {
		e.newWorker()

		// new worker created, so if the queue is full, wait until the new goroutine read the first task.
		if !taskEnqueued {
			e.enqueueTask(execution)
			return
		}
	}

	if taskEnqueued {
		return
	}

	if e.config.BlockOnFullQueue {
		e.enqueueTask(execution)
	} else {
		execution.reject()
	}
}

func (e *executor) newWorker() {
	e.workerWG.Add(1)
	e.workerCount.Add(1)

	go func() {
		defer func() {
			workerCount := e.workerCount.Add(^uint64(0))
			e.workerWG.Done()

			if workerCount == 0 && e.stopped.Load() {
				e.status.Store(TerminatedStatus)
				close(e.done)
			}
		}()

		for {
			select {
			case task := <-e.taskQueue:
				e.work(task)
			case <-e.workerStopSignal:
				return
			case <-e.shutdownSignal:
				e.drainTaskQueue()
				return
			}
		}
	}()
}

func (e *executor) work(task func()) {
	e.workerRunningCount.Add(1)
	task()
	e.workerRunningCount.Add(^uint64(0))
}

func (e *executor) drainTaskQueue() {
	for task := range e.taskQueue {
		e.work(task)
	}
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
	return e.workerCount.Load() < e.config.Concurrency
}

func (e *executor) hasIdleWorker() bool {
	return e.workerRunningCount.Load() < e.workerCount.Load()
}
