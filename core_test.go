package core_test

import (
	"testing"
	"time"

	. "github.com/mplewis/go-worker-pools"
)

// Parameters for validation.
const delay = time.Millisecond * 10 // Time taken by each job to simulate work being done
const workers = 100                 // Worker count
const work = 10000                  // Number of jobs to queue

var params = Params{ExpectedConcurrency: workers, WorkQuantity: work}

// Task sleeps for the specified delay.
func Task() {
	time.Sleep(delay)
}

func TestWorkerPool(t *testing.T) {
	wp := NewWorkerPool(workers)
	err := VerifyConcurrency(wp, params, Task)
	if err != nil {
		t.Error(err)
	}
}

func TestChannelPool(t *testing.T) {
	chp := NewChannelPool(workers)
	err := VerifyConcurrency(chp, params, Task)
	if err != nil {
		t.Error(err)
	}
}

func TestSemaphorePool(t *testing.T) {
	sp := NewSemaphorePool(workers)
	err := VerifyConcurrency(sp, params, Task)
	if err != nil {
		t.Error(err)
	}
}
