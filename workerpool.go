package core

import "github.com/gammazero/workerpool"

// WorkerPoolWrapper is a worker pool that limits concurrency using the gammazero/workerpool package.
type WorkerPoolWrapper struct {
	*workerpool.WorkerPool
}

func NewWorkerPool(workerCount int) *WorkerPoolWrapper {
	return &WorkerPoolWrapper{workerpool.New(workerCount)}
}

func (w *WorkerPoolWrapper) Submit(f func()) {
	w.WorkerPool.Submit(f)
}

func (w *WorkerPoolWrapper) Wait() {
	w.WorkerPool.StopWait()
}
