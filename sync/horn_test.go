package sync

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestHorn(t *testing.T) {
	h := NewHorn[int](10)

	var count atomic.Int64
	// create listeners
	var wg sync.WaitGroup
	wg.Add(5)

	var allDone sync.WaitGroup
	allDone.Add(5)

	for i := 0; i < 5; i++ {
		go func() {
			defer allDone.Done()

			l := h.AddListener()
			wg.Done()
			for _ = range l.Chan() {
				// t.Logf("received %d", v)
				count.Add(1)
			}
		}()
	}

	wg.Wait()

	// sender
	go func() {
		for i := 0; i < 10; i++ {
			h.Send(i)
		}

		h.Close()
	}()

	allDone.Wait()
	if count.Load() != 50 {
		t.Errorf("expected 50, got %d", count.Load())
	}
}
