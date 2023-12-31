package sync

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShard(t *testing.T) {
	var shard = NewShard[atomic.Int64]()

	n := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()

			for j := 0; j < 10_000_000; j++ {
				shard.Get().Add(1)
			}
		}()
	}

	wg.Wait()

	var count int64
	shard.Range(func(v *atomic.Int64) {
		count += v.Load()
	})

	assert.Equal(t, int64(10_000_000*n), count)
}

func BenchmarkShardCounter(b *testing.B) {
	counter := NewShard[atomic.Int64]()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Get().Add(1)
		}
	})
}

func BenchmarkMutexCounter(b *testing.B) {
	var counter int64
	var mu sync.Mutex
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			counter += 1
			mu.Unlock()
		}
	})
}

func BenchmarkAtomicCounter(b *testing.B) {
	var counter atomic.Int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Add(1)
		}
	})
}
