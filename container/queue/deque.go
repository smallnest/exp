package queue

import (
	"sync/atomic"
	"unsafe"
)

// RingBuffer is a basic ring buffer, the capacity must be a power of 2.
type RingBuffer[T any] struct {
	cap  uint64
	mask uint64
	buff []T
}

// NewRingBuffer creates a new ring buffer with the given capacity.
func NewRingBuffer[T any](cap uint64) *RingBuffer[T] {
	if cap == 0 || cap&(cap-1) != 0 {
		panic("capacity must be power of 2")
	}
	return &RingBuffer[T]{
		cap:  cap,
		mask: cap - 1,
		buff: make([]T, cap),
	}
}

// Capacity returns the capacity of the ring buffer.
func (r *RingBuffer[T]) Capacity() uint64 {
	return r.cap
}

// Store stores the value x at the given index i.
func (r *RingBuffer[T]) Store(i uint64, x T) {
	r.buff[i&r.mask] = x
}

// Load loads the value at the given index i.
func (r *RingBuffer[T]) Load(i uint64) T {
	return r.buff[i&r.mask]
}

// Grow rows the ring buffer to the given bounds.
func (r *RingBuffer[T]) Grow(b, t uint64) *RingBuffer[T] {
	newRing := NewRingBuffer[T](2 * r.cap)
	for i := t; i != b; i++ {
		newRing.Store(i, r.Load(i))
	}
	return newRing
}

// Deque is a lock-free single-producer multi-consumer deque.
// Deque is a bleeding-edge lock-free, single-producer multi-consumer, Chase-Lev work stealing deque as presented
// in the paper ["Dynamic Circular Work-Stealing Deque"](https://dl.acm.org/doi/10.1145/1073970.1073974) and further improved in the follow up paper:
// ["Correct and Efficient Work-Stealing for Weak Memory Models"](https://dl.acm.org/doi/10.1145/2442516.2442524).
//
// Still has data race issues.
type Deque[T any] struct {
	top    uint64
	bottom uint64
	buffer unsafe.Pointer
	// garbage []*RingBuffer[T]
}

// NewDeque creates a new deque with the given capacity.
func NewDeque[T any](cap uint64) *Deque[T] {
	return &Deque[T]{buffer: unsafe.Pointer(NewRingBuffer[T](cap))}
}

// Size returns the number of elements in the deque.
func (d *Deque[T]) Size() uint64 {
	b := atomic.LoadUint64(&d.bottom)
	t := atomic.LoadUint64(&d.top)
	if b >= t {
		return b - t
	}
	return 0
}

// Capacity returns the capacity of the deque.
func (d *Deque[T]) Capacity() uint64 {
	return (*RingBuffer[T])(d.buffer).Capacity()
}

// Empty returns true if the deque is empty.
func (d *Deque[T]) Empty() bool {
	return d.Size() == 0
}

// PushBottom pushes the value x to the deque.
// The value x is pushed to the bottom of the deque.
func (d *Deque[T]) PushBottom(x T) {
	b := atomic.LoadUint64(&d.bottom)
	t := atomic.LoadUint64(&d.top)
	buf := (*RingBuffer[T])(atomic.LoadPointer(&d.buffer))
	if buf.Capacity() < b-t+1 {
		// d.garbage = append(d.garbage, buf)
		buf = buf.Grow(b, t)
		atomic.StorePointer(&d.buffer, unsafe.Pointer(buf))
	}
	buf.Store(b, x)
	atomic.StoreUint64(&d.bottom, b+1)
}

// PopBottom pops the value from the deque.
// The value x is popped from the bottom of the deque.
func (d *Deque[T]) PopBottom() (x T, ok bool) {
	b := atomic.AddUint64(&d.bottom, ^uint64(0))
	buf := (*RingBuffer[T])(atomic.LoadPointer(&d.buffer))
	t := atomic.LoadUint64(&d.top)
	if t <= b {
		if t == b {
			if !atomic.CompareAndSwapUint64(&d.top, t, t+1) {
				atomic.StoreUint64(&d.bottom, b+1)
				return x, false
			}
			atomic.StoreUint64(&d.bottom, b+1)
		}
		x = buf.Load(b)
		return x, true
	}
	atomic.StoreUint64(&d.bottom, b+1)
	return x, false
}

// PeekBottom peeks at the value from the deque without removing it.
// The value x is peeked from the top of the deque without removing it.
func (d *Deque[T]) PeekBottom() (x T, ok bool) {
	b := atomic.AddUint64(&d.bottom, ^uint64(0))
	buf := (*RingBuffer[T])(atomic.LoadPointer(&d.buffer))
	t := atomic.LoadUint64(&d.top)
	if t <= b {
		x = buf.Load(b)
		return x, true
	}

	return x, false
}

// StealTop steals a value from the deque.
// The value x is stolen from the top of the deque.
// ok indicates whether the value was stolen, false if the deque is empty or other goroutines are stealing.
func (d *Deque[T]) Steal() (x T, ok bool) {
	t := atomic.LoadUint64(&d.top)
	b := atomic.LoadUint64(&d.bottom)
	if t < b {
		buf := (*RingBuffer[T])(atomic.LoadPointer(&d.buffer))
		x = buf.Load(t)
		if !atomic.CompareAndSwapUint64(&d.top, t, t+1) {
			return x, false
		}
		return x, true
	}
	return x, false
}

// PeekTop peeks at a value from the deque without removing it.
// The value x is peeked from the top of the deque without removing it.
func (d *Deque[T]) PeekTop() (x T, ok bool) {
	t := atomic.LoadUint64(&d.top)
	b := atomic.LoadUint64(&d.bottom)
	if t < b {
		buf := (*RingBuffer[T])(atomic.LoadPointer(&d.buffer))
		x = buf.Load(t)
		return x, true
	}
	return x, false
}
