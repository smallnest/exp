package heap

import (
	"cmp"
)

type binHeap[V cmp.Ordered] []V

func (h binHeap[V]) Len() int           { return len(h) }
func (h binHeap[V]) Less(i, j int) bool { return h[i] < h[j] }
func (h binHeap[V]) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *binHeap[V]) Push(x V) {
	*h = append(*h, x)
}

func (h *binHeap[V]) Pop() V {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// BinHeap is a binary heap.
type BinHeap[V cmp.Ordered] struct {
	maxHeap bool
	binHeap binHeap[V]
}

// NewBinHeap returns a new binary heap.
func NewBinHeap[V cmp.Ordered](opts ...BinHeapOption[V]) *BinHeap[V] {
	h := &BinHeap[V]{}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

// NewBinHeapWithInitial returns a new binary heap with the given initial slice.
func NewBinHeapWithInitial[V cmp.Ordered](s []V, opts ...BinHeapOption[V]) *BinHeap[V] {
	h := &BinHeap[V]{}
	h.binHeap = binHeap[V](s)

	for _, opt := range opts {
		opt(h)
	}

	n := len(s)
	for i := n/2 - 1; i >= 0; i-- {
		sift_down[V](&h.binHeap, i, n, h.maxHeap)
	}

	return h
}

// BinHeapOption is a function that configures a binary heap.
type BinHeapOption[V cmp.Ordered] func(*BinHeap[V]) *BinHeap[V]

// WithMaxHeap returns a BinHeapOption that configures a binary heap to be a max heap.
func WithMaxHeap[V cmp.Ordered](h *BinHeap[V]) *BinHeap[V] {
	h.maxHeap = true
	return h
}

// WithMinHeap returns a BinHeapOption that configures a binary heap to be a min heap.
func WithCapacity[V cmp.Ordered](n int) BinHeapOption[V] {
	return func(h *BinHeap[V]) *BinHeap[V] {
		if h.binHeap == nil {
			h.binHeap = make(binHeap[V], 0, n)
		}
		return h
	}
}

// Len returns the number of elements in the heap.
func (h *BinHeap[V]) Len() int {
	return h.binHeap.Len()
}

// Len returns the number of elements in the heap.
func (h *BinHeap[V]) Push(x V) {
	h.binHeap.Push(x)
	sift_up[V](&h.binHeap, h.binHeap.Len()-1, h.maxHeap)
}

// Peek returns the element at the top of the heap without removing it.
func (h *BinHeap[V]) Peek() (V, bool) {
	var v V
	if h.Len() == 0 {
		return v, false
	}

	return h.binHeap[0], true
}

// Push pushes the element x onto the heap.
func (h *BinHeap[V]) Pop() V {
	n := h.binHeap.Len() - 1
	h.binHeap.Swap(0, n)
	sift_down[V](&h.binHeap, 0, n, h.maxHeap)
	return h.binHeap.Pop()
}

func sift_up[V cmp.Ordered](h *binHeap[V], j int, maxHeap bool) {
	less := h.Less
	if maxHeap {
		less = func(i, j int) bool { return !h.Less(i, j) }
	}
	for {
		i := (j - 1) / 2 // parent
		if i == j || !less(j, i) {
			break
		}
		h.Swap(i, j)
		j = i
	}
}

func sift_down[V cmp.Ordered](h *binHeap[V], i0, n int, maxHeap bool) bool {
	less := h.Less
	if maxHeap {
		less = func(i, j int) bool { return !h.Less(i, j) }
	}

	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !less(j, i) {
			break
		}
		h.Swap(i, j)
		i = j
	}
	return i > i0
}
