package concurrency

import (
	"fmt"
)

// ThreadPool structure to handle tasks.
type ThreadPool struct {
	tasks     chan interface{}
	isRunning bool
	poolSize  int
}

// Callback structure to handle callbacks
type Callback[T any] struct {
	Error   func(errorType error)
	Success func(T)
}

// NewThreadPool constructor for creating a new thread pool with a given pool size.
func NewThreadPool(poolSize int, bufferSize int) *ThreadPool {
	return &ThreadPool{
		poolSize: poolSize,
		tasks:    make(chan interface{}, bufferSize),
	}
}

// Start initializes the thread pool and starts worker goroutines.
func (pool *ThreadPool) Start() {
	for i := 0; i < pool.poolSize; i++ {
		go runThread(pool)
	}

	pool.isRunning = true
}

// runThread continuously fetches and executes tasks from the tasks channel.
func runThread(pool *ThreadPool) {
	for task := range pool.tasks {
		switch t := task.(type) {
		case func():
			t()
		case func(any):
			t(nil)
		default:
			fmt.Println("Unknown task type")
		}
	}
}

// Submit adds a new task to the thread pool.
func (pool *ThreadPool) Submit(task interface{}) {
	if !pool.isRunning {
		return
	}

	//Failsafe go-func TODO: remove this find another way
	go func() {
		pool.tasks <- task
	}()
}
