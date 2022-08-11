package core

import (
	"fmt"
	"sync/atomic"
)

// ConcurrencyManager represents a worker pool.
type ConcurrencyManager interface {
	// Submit submits a work item to the worker pool.
	Submit(func())
	// Wait blocks until all work items have been completed.
	Wait()
}

// Params is the parameters for a concurrency verification test.
type Params struct {
	ExpectedConcurrency int
	WorkQuantity        int
}

// VerifyConcurrency verifies that the worker pool runs tasks at the expected concurrency.
// The task must have a non-zero runtime to ensure that we test the concurrency limit of the worker pool.
func VerifyConcurrency(cm ConcurrencyManager, p Params, task func()) error {
	var currentConcurrency int32
	var maxObservedConcurrency int32
	var completedJobs int32

	// Wrap the designated task so we can monitor the concurrency in realtime.
	work := func() {
		// Keep track of how many workers are running simultaneously.
		atomic.AddInt32(&currentConcurrency, 1)
		defer atomic.AddInt32(&currentConcurrency, -1)

		// Update the maximum observed concurrency.
		curr := atomic.LoadInt32(&currentConcurrency)
		max := atomic.LoadInt32(&maxObservedConcurrency)
		if curr > max {
			atomic.StoreInt32(&maxObservedConcurrency, curr)
		}

		// Run the task.
		task()

		// Mark this job as completed.
		atomic.AddInt32(&completedJobs, 1)
	}

	// Submit all the work and wait for it to complete.
	for i := 0; i < p.WorkQuantity; i++ {
		cm.Submit(work)
	}
	cm.Wait()

	// Verify that jobs ran concurrently as expected.
	if completedJobs != int32(p.WorkQuantity) {
		return fmt.Errorf("expected %d jobs to complete, but only %d completed", p.WorkQuantity, completedJobs)
	}
	if currentConcurrency != 0 {
		return fmt.Errorf("expected all jobs to be complete, but %d jobs are outstanding", currentConcurrency)
	}
	if maxObservedConcurrency != int32(p.ExpectedConcurrency) {
		return fmt.Errorf("expected max concurrency of %d, got %d", p.ExpectedConcurrency, maxObservedConcurrency)
	}
	return nil
}
