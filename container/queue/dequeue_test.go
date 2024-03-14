package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeque_PushPop(t *testing.T) {
	d := NewDeque[int](4)
	for i := 0; i < 10; i++ {
		d.PushBottom(i)
		if got, want := d.Size(), uint64(i+1); got != want {
			t.Fatalf("unexpected size after push: got %d, want %d", got, want)
		}
		if i < 4 {
			if got, want := d.Capacity(), uint64(4); got != want {
				t.Fatalf("unexpected capacity after push: got %d, want %d", got, want)
			}
		} else if i < 8 {
			if got, want := d.Capacity(), uint64(8); got != want {
				t.Fatalf("unexpected capacity after push: got %d, want %d", got, want)
			}
		} else if i < 10 {
			if got, want := d.Capacity(), uint64(16); got != want {
				t.Fatalf("unexpected capacity after push: got %d, want %d", got, want)
			}
		}

		if d.Empty() {
			t.Fatal("deque should not be empty after push")
		}
	}
	for i := 9; i >= 0; i-- {
		val, ok := d.PopBottom()
		if !ok {
			t.Fatal("unexpected Pop() failure")
		}
		if val != i {
			t.Fatalf("unexpected value popped: got %d, want %d", val, i)
		}
		if got, want := d.Size(), uint64(i); got != want {
			t.Fatalf("unexpected size after pop: got %d, want %d", got, want)
		}
		if d.Empty() && i != 0 {
			t.Fatal("deque should not be empty")
		}
	}
	_, ok := d.PopBottom()
	if ok {
		t.Fatal("Pop() should have failed on an empty deque")
	}
}

func TestDeque_Steal(t *testing.T) {
	d := NewDeque[int](4)
	for i := 0; i < 4; i++ {
		d.PushBottom(i)
	}
	for i := 0; i < 4; i++ {
		val, ok := d.Steal()
		if !ok {
			t.Fatal("unexpected Steal() failure")
		}
		if val != i {
			t.Fatalf("unexpected value stolen: got %d, want %d", val, i)
		}
		if got, want := d.Size(), uint64(4-i-1); got != want {
			t.Fatalf("unexpected size after steal: got %d, want %d", got, want)
		}
	}
	_, ok := d.Steal()
	if ok {
		t.Fatal("Steal() should have failed on an empty deque")
	}
}

func TestDeque_Size(t *testing.T) {
	d := NewDeque[int](4)
	assert.Equal(t, uint64(0), d.Size())
	assert.True(t, d.Empty())

	d.PushBottom(1)
	assert.Equal(t, uint64(1), d.Size())
	assert.False(t, d.Empty())

	d.PushBottom(2)
	d.PushBottom(3)
	d.PushBottom(4)
	assert.Equal(t, uint64(4), d.Size())
	assert.False(t, d.Empty())

	d.PopBottom()
	d.PopBottom()
	assert.Equal(t, uint64(2), d.Size())
	assert.False(t, d.Empty())

	d.Steal()
	assert.Equal(t, uint64(1), d.Size())
	assert.False(t, d.Empty())

	d.PopBottom()
	assert.Equal(t, uint64(0), d.Size())
	assert.True(t, d.Empty())
}

func TestDeque_Capacity(t *testing.T) {
	d := NewDeque[int](4)
	assert.Equal(t, uint64(4), d.Capacity())

	d.PushBottom(1)
	d.PushBottom(2)
	d.PushBottom(3)
	d.PushBottom(4)
	assert.Equal(t, uint64(4), d.Capacity())

	d.PopBottom()
	d.PopBottom()
	d.PopBottom()
	d.PopBottom()
	assert.Equal(t, uint64(4), d.Capacity())
}

func TestDeque_Empty(t *testing.T) {
	d := NewDeque[int](4)
	assert.True(t, d.Empty())

	d.PushBottom(1)
	assert.False(t, d.Empty())

	d.PopBottom()
	assert.True(t, d.Empty())
}

func BenchmarkDeque(b *testing.B) {
	deque := NewDeque[int](uint64(b.N))

	b.Run("Push", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			deque.PushBottom(i)
		}
	})

	b.Run("Pop", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			deque.PopBottom()
		}
	})

	b.Run("Push2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			deque.PushBottom(i)
		}
	})

	b.Run("Steal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			deque.Steal()
		}
	})
}
