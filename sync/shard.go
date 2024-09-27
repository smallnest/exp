package sync

import (
	"runtime"

	"github.com/smallnest/gid"
	"golang.org/x/sys/cpu"
)

// Shard is a container of values of the same type
// have n data of type T. Each P has a shard of n data.
type Shard[T any] struct {
	values []value[T]
}

// NewShard creates a new Shard and initializes it with runtime.GOMAXPROCS.
func NewShard[T any]() *Shard[T] {
	n := runtime.GOMAXPROCS(0)

	return &Shard[T]{
		values: make([]value[T], n),
	}
}

type value[T any] struct {
	_ cpu.CacheLinePad // prevent false sharing
	v T
	_ cpu.CacheLinePad // prevent false sharing
}

// Get gets the shard of current goroutine's P.
func (s *Shard[T]) Get() *T {
	if len(s.values) == 0 {
		panic("sync: Sharded is empty and has not been initialized")
	}

	return &s.values[int(gid.PID())%len(s.values)].v
}

// Range calls f for all data of type T in Shard.
//
// It is not goroutine-safe to modify the Shard while iterating.
func (s *Shard[T]) Range(f func(*T)) {
	if len(s.values) == 0 {
		panic("sync: Sharded is empty and has not been initialized")
	}

	for i := range s.values {
		f(&s.values[i].v)
	}
}
