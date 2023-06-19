package maps

import (
	"github.com/smallnest/exp/container/list"
)

// AccessOrderedMap is a map with access order.
type AccessOrderedMap[K comparable, V any] struct {
	*OrderedMap[K, V]
}

// NewOrderedMap creates an empty OrderedMap.
// The parameter `capability` is the initial size of the map.
func NewAccessOrderedMap[K comparable, V any](capability int) *AccessOrderedMap[K, V] {
	orderedMap := &OrderedMap[K, V]{
		entries: make(map[K]*Entry[K, V], capability),
		list:    list.New[*Entry[K, V]](),
	}

	return &AccessOrderedMap[K, V]{orderedMap}
}

// Get returns the value of the given key.
func (m *AccessOrderedMap[K, V]) Get(key K) (val V, existed bool) {
	if entry, existed := m.entries[key]; existed {
		m.list.MoveToFront(entry.element)
		return entry.Value, true
	}

	return
}

// Load returns the value of the given key, alias for Get.
func (m *AccessOrderedMap[K, V]) Load(key K) (V, bool) {
	return m.Get(key)
}
