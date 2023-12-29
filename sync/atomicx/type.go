package atomicx

import (
	"sync"
	"sync/atomic"
)

// Pointer is an atomic boolean type with Wait/Notify support.
type Pointer[T any] struct {
	atomic.Pointer[T]

	mu      sync.Mutex
	condvar *sync.Cond
}

// Wait blocks until the pointer is not equal to the given value.
func (ap *Pointer[T]) Wait() {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	if ap.condvar == nil {
		ap.condvar = sync.NewCond(&ap.mu)
	}

	old := ap.Load()

	for old == ap.Load() {
		ap.condvar.Wait()
	}
}

// Broadcast wakes all goroutines waiting on the pointer.
func (ap *Pointer[T]) Broadcast() {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	if ap.condvar == nil {
		ap.condvar = sync.NewCond(&ap.mu)
	}

	ap.condvar.Broadcast()
}

// Signal wakes one goroutine waiting on the pointer.
func (ap *Pointer[T]) Signal() {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	if ap.condvar == nil {
		ap.condvar = sync.NewCond(&ap.mu)
	}

	ap.condvar.Signal()
}

// AtomicBool is an atomic boolean type with Wait/Notify support.
type Bool struct {
	atomic.Bool

	mu      sync.Mutex
	condvar *sync.Cond
}

// Wait blocks until the boolean is not equal to the given value.
func (ab *Bool) Wait() {
	ab.mu.Lock()
	defer ab.mu.Unlock()

	if ab.condvar == nil {
		ab.condvar = sync.NewCond(&ab.mu)
	}

	v := ab.Load()
	for ab.Load() == v {
		ab.condvar.Wait()
	}
}

// Broadcast wakes all goroutines waiting on the boolean.
func (ab *Bool) Broadcast() {
	ab.mu.Lock()
	defer ab.mu.Unlock()

	if ab.condvar == nil {
		ab.condvar = sync.NewCond(&ab.mu)
	}

	ab.condvar.Broadcast()
}

// Signal wakes one goroutine waiting on the boolean.
func (ab *Bool) Signal() {
	ab.mu.Lock()
	defer ab.mu.Unlock()

	if ab.condvar == nil {
		ab.condvar = sync.NewCond(&ab.mu)
	}

	ab.condvar.Signal()
}

// Int32 is an atomic int32 type with Wait/Notify support.
type Int32 struct {
	atomic.Int32

	mu      sync.Mutex
	condvar *sync.Cond
}

// Wait blocks until the int32 is not equal to the given value.
func (ai *Int32) Wait() {
	ai.mu.Lock()
	defer ai.mu.Unlock()

	if ai.condvar == nil {
		ai.condvar = sync.NewCond(&ai.mu)
	}

	v := ai.Load()
	for ai.Load() == v {
		ai.condvar.Wait()
	}
}

// Broadcast wakes all goroutines waiting on the int32.
func (ai *Int32) Broadcast() {
	ai.mu.Lock()
	defer ai.mu.Unlock()

	if ai.condvar == nil {
		ai.condvar = sync.NewCond(&ai.mu)
	}

	ai.condvar.Broadcast()
}

// Signal wakes one goroutine waiting on the int32.
func (ai *Int32) Signal() {
	ai.mu.Lock()
	defer ai.mu.Unlock()

	if ai.condvar == nil {
		ai.condvar = sync.NewCond(&ai.mu)
	}

	ai.condvar.Signal()
}

// Int64 is an atomic int64 type with Wait/Notify support.
type Int64 struct {
	atomic.Int64

	mu      sync.Mutex
	condvar *sync.Cond
}

// Wait blocks until the int64 is not equal to the given value.
func (ai *Int64) Wait() {
	ai.mu.Lock()
	defer ai.mu.Unlock()

	if ai.condvar == nil {
		ai.condvar = sync.NewCond(&ai.mu)
	}

	v := ai.Load()
	for ai.Load() == v {
		ai.condvar.Wait()
	}
}

// Broadcast wakes all goroutines waiting on the int64.
func (ai *Int64) Broadcast() {
	ai.mu.Lock()
	defer ai.mu.Unlock()

	if ai.condvar == nil {
		ai.condvar = sync.NewCond(&ai.mu)
	}

	ai.condvar.Broadcast()
}

// Signal wakes one goroutine waiting on the int64.
func (ai *Int64) Signal() {
	ai.mu.Lock()
	defer ai.mu.Unlock()

	if ai.condvar == nil {
		ai.condvar = sync.NewCond(&ai.mu)
	}

	ai.condvar.Signal()
}

// Uint32 is an atomic uint32 type with Wait/Notify support.
type Uint32 struct {
	atomic.Uint32

	mu      sync.Mutex
	condvar *sync.Cond
}

// Wait blocks until the uint32 is not equal to the given value.
func (au *Uint32) Wait() {
	au.mu.Lock()
	defer au.mu.Unlock()

	if au.condvar == nil {
		au.condvar = sync.NewCond(&au.mu)
	}

	v := au.Load()
	for au.Load() == v {
		au.condvar.Wait()
	}
}

// Broadcast wakes all goroutines waiting on the uint32.
func (au *Uint32) Broadcast() {
	au.mu.Lock()
	defer au.mu.Unlock()

	if au.condvar == nil {
		au.condvar = sync.NewCond(&au.mu)
	}

	au.condvar.Broadcast()
}

// Signal wakes one goroutine waiting on the uint32.
func (au *Uint32) Signal() {
	au.mu.Lock()
	defer au.mu.Unlock()

	if au.condvar == nil {
		au.condvar = sync.NewCond(&au.mu)
	}

	au.condvar.Signal()
}

// Uint64 is an atomic uint64 type with Wait/Notify support.
type Uint64 struct {
	atomic.Uint64

	mu      sync.Mutex
	condvar *sync.Cond
}

// Wait blocks until the uint64 is not equal to the given value.
func (au *Uint64) Wait() {
	au.mu.Lock()
	defer au.mu.Unlock()

	if au.condvar == nil {
		au.condvar = sync.NewCond(&au.mu)
	}

	v := au.Load()
	for au.Load() == v {
		au.condvar.Wait()
	}
}

// Broadcast wakes all goroutines waiting on the uint64.
func (au *Uint64) Broadcast() {
	au.mu.Lock()
	defer au.mu.Unlock()

	if au.condvar == nil {
		au.condvar = sync.NewCond(&au.mu)
	}

	au.condvar.Broadcast()
}

// Signal wakes one goroutine waiting on the uint64.
func (au *Uint64) Signal() {
	au.mu.Lock()
	defer au.mu.Unlock()

	if au.condvar == nil {
		au.condvar = sync.NewCond(&au.mu)
	}

	au.condvar.Signal()
}

// Uintptr is an atomic uintptr type with Wait/Notify support.
type Uintptr struct {
	atomic.Uintptr

	mu      sync.Mutex
	condvar *sync.Cond
}

// Wait blocks until the uintptr is not equal to the given value.
func (au *Uintptr) Wait() {
	au.mu.Lock()
	defer au.mu.Unlock()

	if au.condvar == nil {
		au.condvar = sync.NewCond(&au.mu)
	}

	v := au.Load()
	for au.Load() == v {
		au.condvar.Wait()
	}
}

// Broadcast wakes all goroutines waiting on the uintptr.
func (au *Uintptr) Broadcast() {
	au.mu.Lock()
	defer au.mu.Unlock()

	if au.condvar == nil {
		au.condvar = sync.NewCond(&au.mu)
	}

	au.condvar.Broadcast()
}

// Signal wakes one goroutine waiting on the uintptr.
func (au *Uintptr) Signal() {
	au.mu.Lock()
	defer au.mu.Unlock()

	if au.condvar == nil {
		au.condvar = sync.NewCond(&au.mu)
	}

	au.condvar.Signal()
}

// Value is an atomic interface{} type with Wait/Notify support.
type Value struct {
	atomic.Value

	mu      sync.Mutex
	condvar *sync.Cond
}

// Wait blocks until the value is not equal to the given value.
func (av *Value) Wait() {
	av.mu.Lock()
	defer av.mu.Unlock()

	if av.condvar == nil {
		av.condvar = sync.NewCond(&av.mu)
	}

	v := av.Load()
	for av.Load() == v {
		av.condvar.Wait()
	}
}

// Broadcast wakes all goroutines waiting on the value.
func (av *Value) Broadcast() {
	av.mu.Lock()
	defer av.mu.Unlock()

	if av.condvar == nil {
		av.condvar = sync.NewCond(&av.mu)
	}

	av.condvar.Broadcast()
}

// Signal wakes one goroutine waiting on the value.
func (av *Value) Signal() {
	av.mu.Lock()
	defer av.mu.Unlock()

	if av.condvar == nil {
		av.condvar = sync.NewCond(&av.mu)
	}

	av.condvar.Signal()
}
