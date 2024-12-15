package sync

import "sync"

type Broadcaster[T any] struct {
	mu       *sync.Mutex
	cond     *sync.Cond
	signaled bool

	v T
}

// Broadcast sends a signal with a value to all goroutines waiting on the broadcaster.
func NewBroadcaster[T any]() *Broadcaster[T] {
	var mu sync.Mutex
	return &Broadcaster[T]{
		mu:       &mu,
		cond:     sync.NewCond(&mu),
		signaled: false,
	}
}

// Go waits until something is broadcasted, and runs the given function in a new
// goroutine with the value that was broadcasted.
func (b *Broadcaster[T]) Go(fn func(v T)) {
	go func() {
		b.cond.L.Lock()
		defer b.cond.L.Unlock()

		for !b.signaled {
			b.cond.Wait()
		}
		fn(b.v)
	}()
}

// Broadcast broadcasts a signal to all
// waiting function and unblocks them.
func (b *Broadcaster[T]) Broadcast(v T) {
	b.cond.L.Lock()
	b.v = v
	b.signaled = true
	b.cond.L.Unlock()

	b.cond.Broadcast()
}

func (b *Broadcaster[T]) Reset() {
	var t T
	b.cond.L.Lock()
	b.signaled = false
	b.v = t
	b.cond.L.Unlock()
}
