package win

import (
	"sync"

	"github.com/smallnest/exp/container/heap"
	"golang.org/x/exp/constraints"
)

type Ordered = constraints.Ordered

type buckets[K Ordered, V any] struct {
	mu sync.RWMutex
	bs queue[K, V]
	m  map[K]*Bucket[K, V]
}

func NewBuckets[K Ordered, V any]() *buckets[K, V] {
	return &buckets[K, V]{
		m: make(map[K]*Bucket[K, V]),
	}
}

type queue[K Ordered, V any] []*Bucket[K, V]

func (q queue[_, _]) Len() int           { return len(q) }
func (q queue[_, _]) Less(i, j int) bool { return q[i].Key < q[j].Key }
func (q queue[_, _]) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

func (h *queue[K, V]) Push(b *Bucket[K, V]) {
	*h = append(*h, b)
}

func (h *queue[K, V]) Pop() *Bucket[K, V] {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Add adds a value to the bucket.
func (h *buckets[K, V]) Add(k K, v V) {
	h.mu.Lock()
	bucket := h.m[k]
	if bucket == nil {
		bucket = &Bucket[K, V]{Key: k, Value: []V{}}
		h.m[k] = bucket
		heap.Push(&h.bs, bucket)
	}
	h.mu.Unlock()

	bucket.Add(v)
}

// Get returns the bucket for the key.
func (h *buckets[K, V]) Last() *Bucket[K, V] {
	h.mu.Lock()
	defer h.mu.Unlock()

	buckets := heap.Pop[*Bucket[K, V]](&h.bs)
	if buckets == nil {
		return nil
	}

	delete(h.m, buckets.Key)

	return buckets
}

// LastN returns the last n buckets but not pop them.
func (h *buckets[K, V]) LastN(n int) []*Bucket[K, V] {
	h.mu.Lock()
	defer h.mu.Unlock()

	if n > len(h.bs) {
		n = len(h.bs)
	}

	buckets := h.bs[:n]

	return buckets
}

// Bucket is a key-value pair.
type Bucket[K Ordered, V any] struct {
	Key K

	Mu    sync.RWMutex
	Value []V
}

// Add adds a value to the bucket.
func (b *Bucket[K, V]) Add(v V) {
	b.Mu.Lock()
	b.Value = append(b.Value, v)
	b.Mu.Unlock()
}

// Values returns the values in the bucket.
func (b *Bucket[K, V]) Values() []V {
	return b.Value
}
