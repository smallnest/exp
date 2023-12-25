package chanx

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBatch(t *testing.T) {
	ch := make(chan int, 10)
	for i := 0; i < 10; i++ {
		ch <- i
	}

	var count atomic.Int64
	go Batch[int](context.Background(), ch, 5, func(batch []int) {
		if len(batch) != 5 {
			assert.Fail(t, "expected batch size 5, got %d", len(batch))
		}
		count.Add(int64(len(batch)))
	})
	time.Sleep(time.Second)
	close(ch)
	assert.Equal(t, int64(10), count.Load())

	ch = make(chan int, 10)
	for i := 0; i < 10; i++ {
		ch <- i
	}

	count.Store(0)
	i := 0
	go Batch[int](context.Background(), ch, 3, func(batch []int) {
		if i < 3 && len(batch) != 3 {
			assert.Fail(t, "expected batch size 5, got %d", len(batch))
		}
		if i == 3 && len(batch) != 1 {
			assert.Fail(t, "expected batch size 1, got %d", len(batch))
		}
		i++
		count.Add(int64(len(batch)))
	})
	time.Sleep(time.Second)
	close(ch)
	assert.Equal(t, int64(10), count.Load())
}

func TestBatch_Context(t *testing.T) {
	ch := make(chan int, 10)
	for i := 0; i < 10; i++ {
		ch <- i
	}

	var count atomic.Int64
	ctx, cancel := context.WithCancel(context.Background())
	go Batch[int](ctx, ch, 5, func(batch []int) {
		if len(batch) != 5 {
			assert.Fail(t, "expected batch size 5, got %d", len(batch))
		}
		count.Add(int64(len(batch)))
	})
	time.Sleep(time.Second)
	cancel()
	assert.Equal(t, int64(10), count.Load())
}

func TestFlatBatch(t *testing.T) {
	ch := make(chan []int, 10)
	for i := 0; i < 10; i++ {
		ch <- []int{i, i}
	}

	var count atomic.Int64
	go FlatBatch[int](context.Background(), ch, 5, func(batch []int) {
		assert.NotEmpty(t, batch)
		count.Add(int64(len(batch)))
	})
	time.Sleep(time.Second)
	close(ch)
	assert.Equal(t, int64(20), count.Load())

	ch = make(chan []int, 10)
	for i := 0; i < 10; i++ {
		ch <- []int{i, i}
	}

	count.Store(0)
	go FlatBatch[int](context.Background(), ch, 3, func(batch []int) {
		assert.NotEmpty(t, batch)
		count.Add(int64(len(batch)))
	})
	time.Sleep(time.Second)
	close(ch)
	assert.Equal(t, int64(20), count.Load())
}

func TestFlatBatch_Context(t *testing.T) {
	ch := make(chan []int, 10)
	for i := 0; i < 10; i++ {
		ch <- []int{i, i}
	}

	var count atomic.Int64
	ctx, cancel := context.WithCancel(context.Background())
	go FlatBatch[int](ctx, ch, 5, func(batch []int) {
		assert.NotEmpty(t, batch)
		count.Add(int64(len(batch)))
	})
	time.Sleep(time.Second)
	cancel()
	assert.Equal(t, int64(20), count.Load())
}
