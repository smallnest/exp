package heap

import (
	"math/rand"
	"sort"
	"testing"
)

func TestBinHeap(t *testing.T) {
	h := NewBinHeap[int]()

	h.Push(5)
	h.Push(3)
	h.Push(7)

	if h.Len() != 3 {
		t.Errorf("expected Len to be 3, got %v", h.Len())
	}

	if val := h.Pop(); val != 3 {
		t.Errorf("expected Pop to return 3, got %v", val)
	}

	if val := h.Pop(); val != 5 {
		t.Errorf("expected Pop to return 5, got %v", val)
	}

	if val := h.Pop(); val != 7 {
		t.Errorf("expected Pop to return 7, got %v", val)
	}

}

func TestBinHeap_WithMaxHeap(t *testing.T) {
	h := NewBinHeap[int](WithMaxHeap)
	h.Push(5)
	h.Push(3)
	h.Push(7)

	if h.Len() != 3 {
		t.Errorf("expected Len to be 3, got %v", h.Len())
	}

	if val := h.Pop(); val != 7 {
		t.Errorf("expected Pop to return 7, got %v", val)
	}

	if val := h.Pop(); val != 5 {
		t.Errorf("expected Pop to return 5, got %v", val)
	}

	if val := h.Pop(); val != 3 {
		t.Errorf("expected Pop to return 3, got %v", val)
	}
}

func TestBinHeap_WithCapacity(t *testing.T) {
	h := NewBinHeap[int](WithCapacity[int](3))
	h.Push(5)
	h.Push(3)
	h.Push(7)

	if h.Len() != 3 {
		t.Errorf("expected Len to be 7, got %v", h.Len())
	}

	if val := h.Pop(); val != 3 {
		t.Errorf("expected Pop to return 3, got %v", val)
	}

	if val := h.Pop(); val != 5 {
		t.Errorf("expected Pop to return 5, got %v", val)
	}

	if val := h.Pop(); val != 7 {
		t.Errorf("expected Pop to return 7, got %v", val)
	}
}

func TestBinHeap_Rand(t *testing.T) {
	data := make([]int, 0, 100)
	h := NewBinHeap[int](WithCapacity[int](3))
	for i := 0; i < 100; i++ {
		v := rand.Intn(10000)
		data = append(data, v)
		h.Push(v)
	}

	if h.Len() != 100 {
		t.Errorf("expected Len to be 100, got %v", h.Len())
	}

	sort.Ints(data)
	for i := 0; i < 100; i++ {
		if val := h.Pop(); val != data[i] {
			t.Errorf("expected Pop to return %v, got %v", data[i], val)
		}
	}
}

func TestBinHeap_peek(t *testing.T) {
	h := NewBinHeap[int]()

	h.Push(5)
	peek, ok := h.Peek()
	if peek != 5 || !ok {
		t.Errorf("expected Peek to return 5, got %v", peek)
	}

	h.Push(3)
	peek, ok = h.Peek()
	if peek != 3 || !ok {
		t.Errorf("expected Peek to return 3, got %v", peek)
	}

	h.Push(7)
	peek, ok = h.Peek()
	if peek != 3 || !ok {
		t.Errorf("expected Peek to return 3, got %v", peek)
	}

	v := h.Pop()
	if v != 3 {
		t.Errorf("expected Pop to return 3, got %v", v)
	}
	peek, ok = h.Peek()
	if peek != 5 || !ok {
		t.Errorf("expected Peek to return 5, got %v", peek)
	}

	v = h.Pop()
	if v != 5 {
		t.Errorf("expected Pop to return 5, got %v", v)
	}
	peek, ok = h.Peek()
	if peek != 7 || !ok {
		t.Errorf("expected Peek to return 7, got %v", peek)
	}

	v = h.Pop()
	if v != 7 {
		t.Errorf("expected Pop to return 7, got %v", v)
	}
	peek, ok = h.Peek()
	if peek != 0 || ok {
		t.Errorf("expected Peek to return 0, got %v", peek)
	}

}

func TestBinHeap_NewBinHeapWithInitial(t *testing.T) {
	data := []int{5, 3, 7}
	h := NewBinHeapWithInitial[int](data)

	if h.Len() != 3 {
		t.Errorf("expected Len to be 3, got %v", h.Len())
	}

	if val := h.Pop(); val != 3 {
		t.Errorf("expected Pop to return 3, got %v", val)
	}

	if val := h.Pop(); val != 5 {
		t.Errorf("expected Pop to return 5, got %v", val)
	}

	if val := h.Pop(); val != 7 {
		t.Errorf("expected Pop to return 7, got %v", val)
	}

	data1 := make([]int, 0, 100)
	data2 := make([]int, 0, 100)

	for i := 0; i < 100; i++ {
		v := rand.Intn(10000)
		data1 = append(data1, v)
		data2 = append(data2, v)
	}
	h = NewBinHeapWithInitial[int](data1)
	sort.Ints(data2)

	for i := 0; i < 100; i++ {
		if val := h.Pop(); val != data2[i] {
			t.Errorf("expected Pop to return %v, got %v", data[i], val)
		}
	}
}
