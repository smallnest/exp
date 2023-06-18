package maps

import (
	"encoding/json"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBidiMap(t *testing.T) {
	// create
	m := NewBidiMap[int, string](10)
	assert.Equal(t, 0, m.Len())

	// add
	for i := 0; i < 10; i++ {
		m.Put(i, strconv.Itoa(i))
	}
	assert.Equal(t, 10, m.Len())

	// get
	for i := 0; i < 10; i++ {
		val, existed := m.Get(i)
		assert.True(t, existed)
		assert.Equal(t, strconv.Itoa(i), val)
	}
	for i := 0; i < 10; i++ {
		val, existed := m.GetKey(strconv.Itoa(i))
		assert.True(t, existed)
		assert.Equal(t, i, val)
	}

	// update
	for i := 0; i < 10; i++ {
		m.Put(i, strconv.Itoa(i*100))
	}
	assert.Equal(t, 10, m.Len())
	for i := 0; i < 10; i++ {
		val, existed := m.Get(i)
		assert.True(t, existed)
		assert.Equal(t, strconv.Itoa(i*100), val)
	}
	for i := 0; i < 10; i++ {
		val, existed := m.GetKey(strconv.Itoa(i * 100))
		assert.True(t, existed)
		assert.Equal(t, i, val)
	}
	for i := 1; i < 10; i++ {
		key, existed := m.GetKey(strconv.Itoa(i))
		assert.False(t, existed, key)
	}

	// remove
	for i := 0; i < 5; i++ {
		m.Remove(i)
	}
	assert.Equal(t, 5, m.Len())
	val, existed := m.Get(0)
	assert.False(t, existed)
	val, existed = m.Get(5)
	assert.True(t, existed)
	assert.Equal(t, "500", val)

	// remove by value
	m.RemoveValue("500")
	_, existed = m.Get(5)
	assert.False(t, existed)

	// clear
	m.Clear()
	for i := 0; i < 10; i++ {
		m.Put(i, strconv.Itoa(i))
	}

	//range
	m.Range(func(key int, value string) bool {
		assert.Equal(t, strconv.Itoa(key), value)
		return true
	})

	//keys
	keys := m.Keys()
	assert.Equal(t, 10, len(keys))
	sort.Ints(keys)
	for i := 0; i < 10; i++ {
		assert.Equal(t, i, keys[i])
	}

	// values
	values := m.Values()
	assert.Equal(t, 10, len(values))
	sort.Strings(values)
	for i := 0; i < 10; i++ {
		assert.Equal(t, strconv.Itoa(i), values[i])
	}

	// Clone
	m2 := m.Clone()
	assert.Equal(t, m.Len(), m2.Len())
	assert.True(t, m.Equals(m2))

}

func TestBidiMap_JSON(t *testing.T) {
	m := NewBidiMap[int, string](10)
	for i := 0; i < 10; i++ {
		m.Put(i, strconv.Itoa(i))
	}

	data, err := json.Marshal(m)
	assert.NoError(t, err)

	m2 := NewBidiMap[int, string](10)
	err = json.Unmarshal(data, m2)
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		val, existed := m.Get(i)
		assert.True(t, existed)
		assert.Equal(t, strconv.Itoa(i), val)
	}
	for i := 0; i < 10; i++ {
		val, existed := m.GetKey(strconv.Itoa(i))
		assert.True(t, existed)
		assert.Equal(t, i, val)
	}
}
