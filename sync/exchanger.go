package sync

import (
	"sync/atomic"

	"github.com/smallnest/gid"
)

// Exchanger is a synchronization primitive that allows two goroutines to
// exchange values.
// It is similar to a rendezvous point or barrier at which two goroutines
// swap values and proceed.
// It consists of two channels and each goroutine owns one channel.
// Each goroutine calls Exchange with the value to give to the other goroutine
// and receives the value from the other goroutine.
// It is a rendezvous because both goroutines wait for the other before
// exchanging values.
// It is a barrier because both goroutines block until both have called
// Exchange.
type Exchanger[T any] struct {
	leftGoID, rightGoID int64
	left, right         chan T
}

// NewExchanger creates a new exchanger.
func NewExchanger[T any]() *Exchanger[T] {
	return &Exchanger[T]{
		leftGoID:  -1,
		rightGoID: -1,
		left:      make(chan T, 1),
		right:     make(chan T, 1),
	}
}

// Exchange exchanges value between two goroutines.
// It returns the value received from the other goroutine.
//
// It panics if called from neither left nor right goroutine.
//
// If the other goroutine has not called Exchange yet, it blocks.
func (e *Exchanger[T]) Exchange(value T) T {
	goid := gid.ID()

	// left goroutine
	isLeft := atomic.CompareAndSwapInt64(&e.leftGoID, -1, goid)
	if !isLeft {
		isLeft = atomic.LoadInt64(&e.leftGoID) == goid
	}
	if isLeft {
		e.right <- value // send value to right
		return <-e.left  // wait for value from right
	}

	// right goroutine
	isRight := atomic.CompareAndSwapInt64(&e.rightGoID, -1, goid)
	if !isRight {
		isRight = atomic.LoadInt64(&e.rightGoID) == goid
	}
	if isRight {
		e.left <- value  // send value to left
		return <-e.right // wait for value from left
	}

	// other goroutine
	panic("sync: exchange called from neither left nor right goroutine")
}
