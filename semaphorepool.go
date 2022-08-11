package core

import (
	"context"

	"golang.org/x/sync/semaphore"
)

// SemaphorePool is a worker pool that limits concurrency using a semaphore.
type SemaphorePool struct {
	concurrency int64
	sem         *semaphore.Weighted
	ctx         context.Context
}

func NewSemaphorePool(concurrency int64) *SemaphorePool {
	return &SemaphorePool{concurrency, semaphore.NewWeighted(int64(concurrency)), context.Background()}
}

func (s SemaphorePool) Submit(f func()) {
	err := s.sem.Acquire(s.ctx, 1)
	if err != nil {
		panic(err)
	}
	go func() {
		f()
		s.sem.Release(1)
	}()
}

func (s SemaphorePool) Wait() {
	err := s.sem.Acquire(s.ctx, s.concurrency)
	if err != nil {
		panic(err)
	}
}
