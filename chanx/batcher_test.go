package chanx

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBatch(t *testing.T) {
	ch := make(chan int, 10)
	for i := 0; i < 10; i++ {
		ch <- i
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go Batch[int](ch, 5, func(batch []int) {
		if len(batch) != 5 {
			assert.Fail(t, "expected batch size 5, got %d", len(batch))
		}
		wg.Done()
	})
	wg.Wait()

	for i := 0; i < 10; i++ {
		ch <- i
	}

	wg.Add(3)
	i := 0
	go Batch[int](ch, 3, func(batch []int) {
		if i < 3 && len(batch) != 3 {
			assert.Fail(t, "expected batch size 5, got %d", len(batch))
		}
		if i == 3 && len(batch) != 1 {
			assert.Fail(t, "expected batch size 1, got %d", len(batch))
		}
		i++

		wg.Done()
	})
	wg.Wait()
}

func TestFlatBatch(t *testing.T) {
	ch := make(chan []int, 10)
	for i := 0; i < 10; i++ {
		ch <- []int{i, i}
	}

	go FlatBatch[int](ch, 5, func(batch []int) {
		assert.NotEmpty(t, batch)
	})
	time.Sleep(time.Second)

	for i := 0; i < 10; i++ {
		ch <- []int{i, i}
	}

	go FlatBatch[int](ch, 3, func(batch []int) {
		assert.NotEmpty(t, batch)
	})
	time.Sleep(time.Second)
}
