package core_test

import (
	"testing"

	. "github.com/mplewis/go-worker-pools"
)

// Noop does nothing.
func Noop() {}

// RunBenchOn runs a benchmark on the given ConcurrencyManager using a builder.
func RunBenchOn(b *testing.B, cmBuilder func() ConcurrencyManager) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		cm := cmBuilder()
		b.StartTimer()
		for i := 0; i < work; i++ {
			cm.Submit(Noop)
		}
		cm.Wait()
	}
}

func BenchmarkWorkerPool(b *testing.B) {
	RunBenchOn(b, func() ConcurrencyManager {
		return NewWorkerPool(workers)
	})
}

func BenchmarkChannelPool(b *testing.B) {
	RunBenchOn(b, func() ConcurrencyManager {
		return NewChannelPool(workers)
	})
}

func BenchmarkSemaphorePool(b *testing.B) {
	RunBenchOn(b, func() ConcurrencyManager {
		return NewSemaphorePool(workers)
	})
}
