// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heap

import (
	stdheap "container/heap"
	"testing"
)

func BenchmarkExpHeap(b *testing.B) {
	h := new(myHeap)
	h.verify(b, 0)

	for k := 0; k < b.N; k++ {
		for i := 20; i > 10; i-- {
			h.Push(i)
		}
		Init[int](h)
		h.verify(b, 0)

		for i := 10; i > 0; i-- {
			Push(h, i)
			h.verify(b, 0)
		}

		for i := 1; h.Len() > 0; i++ {
			x := Pop[int](h)
			if i < 20 {
				Push(h, 20+i)
			}
			h.verify(b, 0)
			if x != i {
				b.Errorf("%d.th pop got %d; want %d", i, x, i)
			}
		}
	}
}

func BenchmarkExpHeapDup(b *testing.B) {
	const n = 10000
	h := make(myHeap, 0, n)
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			Push(&h, 0) // all elements are the same
		}
		for h.Len() > 0 {
			Pop[int](&h)
		}
	}
}

type myStdHeap []int

func (h *myStdHeap) Less(i, j int) bool {
	return (*h)[i] < (*h)[j]
}

func (h *myStdHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *myStdHeap) Len() int {
	return len(*h)
}

func (h *myStdHeap) Pop() (v any) {
	*h, v = (*h)[:h.Len()-1], (*h)[h.Len()-1]
	return
}

func (h *myStdHeap) Push(v any) {
	*h = append(*h, v.(int))
}

func (h myStdHeap) verify(b testing.TB, i int) {
	b.Helper()
	n := h.Len()
	j1 := 2*i + 1
	j2 := 2*i + 2
	if j1 < n {
		if h.Less(j1, i) {
			b.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d", i, h[i], j1, h[j1])
			return
		}
		h.verify(b, j1)
	}
	if j2 < n {
		if h.Less(j2, i) {
			b.Errorf("heap invariant invalidated [%d] = %d > [%d] = %d", i, h[i], j1, h[j2])
			return
		}
		h.verify(b, j2)
	}
}

func BenchmarkStdHeap(b *testing.B) {
	h := new(myStdHeap)
	h.verify(b, 0)

	for k := 0; k < b.N; k++ {
		for i := 20; i > 10; i-- {
			h.Push(i)
		}
		stdheap.Init(h)
		h.verify(b, 0)

		for i := 10; i > 0; i-- {
			stdheap.Push(h, i)
			h.verify(b, 0)
		}

		for i := 1; h.Len() > 0; i++ {
			x := stdheap.Pop(h)
			if i < 20 {
				stdheap.Push(h, 20+i)
			}
			h.verify(b, 0)
			if x != i {
				b.Errorf("%d.th pop got %d; want %d", i, x, i)
			}
		}
	}
}

func BenchmarkStdHeapDup(b *testing.B) {
	const n = 10000
	h := make(myStdHeap, 0, n)
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			stdheap.Push(&h, 0) // all elements are the same
		}
		for h.Len() > 0 {
			stdheap.Pop(&h)
		}
	}
}
