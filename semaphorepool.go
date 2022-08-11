package core

import (
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

// SemaphorePool is a worker pool that limits concurrency using a semaphore.
type SemaphorePool struct {
	concurrency int64
	sem         *semaphore.Weighted
	wg          *sync.WaitGroup
	ctx         context.Context
}

func NewSemaphorePool(concurrency int64) *SemaphorePool {
	return &SemaphorePool{
		concurrency,
		semaphore.NewWeighted(int64(concurrency)),
		&sync.WaitGroup{},
		context.Background(),
	}
}

func (s SemaphorePool) Submit(f func()) {
	s.wg.Add(1)
	// FIXME: This actually blocks. But moving the Acquire into the goroutine causes `go test . -race` to crash with
	// "limit on 8128 simultaneously alive goroutines is exceeded, dying".
	err := s.sem.Acquire(s.ctx, 1)
	if err != nil {
		panic(err)
	}
	go func() {
		f()
		s.wg.Done()
		s.sem.Release(1)
	}()
}

func (s SemaphorePool) Wait() {
	s.wg.Wait()
	err := s.sem.Acquire(s.ctx, s.concurrency)
	if err != nil {
		panic(err)
	}
}
