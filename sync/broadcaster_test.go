package sync

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBroadcaster(t *testing.T) {
	b := NewBroadcaster[int]()

	var count atomic.Int32
	var sum atomic.Int64

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		b.Go(func(v int) {
			defer wg.Done()

			count.Add(1)
			sum.Add(int64(v))
		})
	}

	b.Broadcast(10)

	wg.Wait()

	assert.Equal(t, int32(10), count.Load())
	assert.Equal(t, int64(100), sum.Load())

}
func TestBroadcasterReset(t *testing.T) {
	b := NewBroadcaster[int]()

	var count atomic.Int32
	var sum atomic.Int64

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		b.Go(func(v int) {
			defer wg.Done()

			count.Add(1)
			sum.Add(int64(v))
		})
	}

	b.Broadcast(10)
	wg.Wait()

	assert.Equal(t, int32(10), count.Load())
	assert.Equal(t, int64(100), sum.Load())

	// Reset the broadcaster
	b.Reset()

	// Reset the counters
	count.Store(0)
	sum.Store(0)

	// Add new goroutines after reset
	for i := 0; i < 10; i++ {
		wg.Add(1)
		b.Go(func(v int) {
			defer wg.Done()

			count.Add(1)
			sum.Add(int64(v))
		})
	}

	b.Broadcast(20)
	wg.Wait()

	assert.Equal(t, int32(10), count.Load())
	assert.Equal(t, int64(200), sum.Load())
}
