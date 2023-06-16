package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderedMap(t *testing.T) {
	// create
	m := New[int, int](10)
	assert.Equal(t, 0, m.Len())

	// set
	for i := 0; i < 10; i++ {
		m.Set(i, i)
	}
	assert.Equal(t, 10, m.Len())

	// get
	for i := 0; i < 10; i++ {
		val, existed := m.Get(i)
		assert.True(t, existed)
		assert.Equal(t, i, val)
	}

	// update
	for i := 0; i < 10; i++ {
		m.Set(i, i*100)
	}
	assert.Equal(t, 10, m.Len())
	for i := 0; i < 10; i++ {
		val, existed := m.Get(i)
		assert.True(t, existed)
		assert.Equal(t, i*100, val)
	}

	// delete
	for i := 0; i < 5; i++ {
		m.Delete(i)
	}
	assert.Equal(t, 5, m.Len())
	val, existed := m.Get(0)
	assert.False(t, existed)
	assert.Equal(t, 0, val)
	val, existed = m.Get(5)
	assert.True(t, existed)
	assert.Equal(t, 500, val)

	// clear
	m.Clear()
	for i := 0; i < 10; i++ {
		m.Store(i, i)
	}

	//range
	k := 0
	m.Range(func(key int, value int) {
		assert.Equal(t, k, key)
		assert.Equal(t, k, value)
		k++
	})

	// foreach
	m.ForEach(func(k int, v int) {
		if k < 0 || k > 9 {
			t.Error("key out of range", k)
		}
		if v < 0 || v > 9 {
			t.Error("value out of range", v)
		}
	})

	// oldest
	i := 0
	for entry := m.Oldest(); entry != nil; entry = entry.Next() {
		assert.Equal(t, i, entry.Key)
		assert.Equal(t, i, entry.Value)
		i++
	}

	//newest
	i = 9
	for entry := m.Newest(); entry != nil; entry = entry.Prev() {
		assert.Equal(t, i, entry.Key)
		assert.Equal(t, i, entry.Value)
		i--
	}
}

func TestOrderedMap_Add(t *testing.T) {
	m := New[int, int](10)
	for i := 0; i < 10; i++ {
		m.Set(i, i)
	}

	m2 := New[int, int](10)
	for i := 5; i < 15; i++ {
		m.Set(i, i)
	}

	m.AddOrderedMap(*m2)
	assert.Equal(t, 15, m.Len())

	var entries []*Entry[int, int]
	for i := 10; i < 20; i++ {
		entries = append(entries, &Entry[int, int]{
			Key:   i,
			Value: i,
		})
	}
	m.AddEntries(entries...)
	assert.Equal(t, 20, m.Len())

	m3 := make(map[int]int)
	for i := 20; i < 30; i++ {
		m3[i] = i
	}
	m.AddMap(m3)
	assert.Equal(t, 30, m.Len())
}
