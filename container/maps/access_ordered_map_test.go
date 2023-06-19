package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccessOrderedMap(t *testing.T) {
	m := NewAccessOrderedMap[int, int](10)
	assert.Equal(t, 0, m.Len())

	// set
	for i := 0; i < 10; i++ {
		m.Set(i, i)
	}

	values := m.Values()
	assert.Equal(t, 10, len(values))
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, values)

	m.Get(5)
	m.Load(8)
	values = m.Values()
	assert.Equal(t, 10, len(values))
	assert.Equal(t, []int{8, 5, 0, 1, 2, 3, 4, 6, 7, 9}, values)
}
