<div align="center">

# 🚀 `executor`: elegant concurrency for golang 🚀

![Build status](https://github.com/go-umbrella/executor/actions/workflows/build.yml/badge.svg)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/go-umbrella/executor)
<a title="Release" target="_blank" href="https://github.com/go-umbrella/executor/releases"><img src="https://img.shields.io/github/v/release/go-umbrella/executor"></a>
<a title="Codecov" target="_blank" href="https://codecov.io/gh/go-umbrella/executor"><img src="https://codecov.io/gh/go-umbrella/executor/branch/main/graph/badge.svg"/></a>
<a title="Go Report Card" target="_blank" href="https://goreportcard.com/report/github.com/go-umbrella/executor"><img src="https://goreportcard.com/badge/github.com/go-umbrella/executor"/></a>
[![GoDoc](https://pkg.go.dev/badge/github.com/go-umbrella/executor)](https://pkg.go.dev/github.com/go-umbrella/executor)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/go-umbrella/executor/blob/main/LICENSE)

</div>

---

<p align="center">
 • <a href="#-about-executor">About Executor</a> •
 <a href="#-installation">Installation</a> •
 <a href="#-examples">Examples</a> •
</p>

---

## 💻 About Executor

`executor` is your toolkit for simple, reliable and elegant concurrency in go.

## 🛠 Installation

To install `executor`, run:

```
go get -u github.com/go-umbrella/executor
```

## 📚 TL;DR
[`Executor`](https://pkg.go.dev/github.com/go-umbrella/executor#Executor) is a specialized workerpool to run async tasks, then:
- Use [`executor.Computation()`](https://pkg.go.dev/github.com/go-umbrella/executor#Computation) if you want to run **computation** intensive tasks (math operations, etc).
- Use [`executor.IO()`](https://pkg.go.dev/github.com/go-umbrella/executor#IO) if you want to run **IO** intensive tasks  (http calls, database connections, etc).
- Use [`executor.Single()`](https://pkg.go.dev/github.com/go-umbrella/executor#Single) if you want to run only **single** async task at a time.
- Use [`executor.New(name, config)`](https://pkg.go.dev/github.com/go-umbrella/executor#New) if you want to create a **new custom executor**.

After `executor` was configured, then:
- [`execution := e.Go(ctx, task, opts)`](https://pkg.go.dev/github.com/go-umbrella/executor#Executor.Go) to submit a task to run async.
  - [`Execution`](https://pkg.go.dev/github.com/go-umbrella/executor#Execution) control async execution.
  - [`result, err := execution.Wait().Get()`](https://pkg.go.dev/github.com/go-umbrella/executor#Execution.Get) wait task to be finished and get result.
  - [`cancelled := execution.Cancel()`](https://pkg.go.dev/github.com/go-umbrella/executor#Execution.Cancel) cancel long-running tasks.
  - [`status := execution.Status()`](https://pkg.go.dev/github.com/go-umbrella/executor#Execution.Status) to see execution status.
  - [`<-execution.Done()`](https://pkg.go.dev/github.com/go-umbrella/executor#Execution.Done) to read done channel.

## 📚 Examples

### 👨🏻‍🏭 Executor (Workerpool)

Create a custom `executor` (Workerpool):
```go
import "github.com/go-umbrella/executor"

var taskExecutor = executor.New("task-executor", executor.Config{
	Concurrency:         16,
	QueueSize:           256,
	EagerInitialization: true,
})
```

Or use pre-defined executors:

```go
import "github.com/go-umbrella/executor"

executor.Computation()  // Optimized for intensive computation tasks (math operations, etc).
executor.IO()           // Optimized for intensive IO operations (http calls, database connections, etc)
executor.Single()       // To run only single task at a time
```

### 📝 Task

Create a task:
```go
func longDatabaseQuery(ctx executor.TaskContext) (interface{}, error) {
    time.Sleep(100 * time.Millisecond)
    return nil, nil
}
```

Submit task and get the result:

```go
execution := taskExecutor.Go(context.Background(), longDatabaseQuery)
result, err := execution.Wait().Get()
```
