package maps

import (
	"github.com/smallnest/exp/container/list"
)

// Entry is a key-value pair in OrderedMap.
type Entry[K comparable, V any] struct {
	Key   K
	Value V

	element *list.Element[*Entry[K, V]]
}

// Next returns a pointer to the next entry.
func (e *Entry[K, V]) Next() *Entry[K, V] {
	entry := e.element.Next()
	if entry == nil {
		return nil
	}

	return entry.Value
}

// Prev returns a pointer to the previous entry.
func (e *Entry[K, V]) Prev() *Entry[K, V] {
	entry := e.element.Prev()
	if entry == nil {
		return nil
	}

	return entry.Value
}

// OrderedMap is a map with insert order.
type OrderedMap[K comparable, V any] struct {
	entries map[K]*Entry[K, V]
	list    *list.List[*Entry[K, V]]
}

// NewOrderedMap creates an empty OrderedMap.
// The parameter `capability` is the initial size of the map.
func NewOrderedMap[K comparable, V any](capability int) *OrderedMap[K, V] {
	orderedMap := &OrderedMap[K, V]{
		entries: make(map[K]*Entry[K, V], capability),
		list:    list.New[*Entry[K, V]](),
	}

	return orderedMap
}

// Get returns the value of the given key.
func (m *OrderedMap[K, V]) Get(key K) (val V, existed bool) {
	if entry, existed := m.entries[key]; existed {
		return entry.Value, true
	}

	return
}

// Load returns the value of the given key, alias for Get.
func (m *OrderedMap[K, V]) Load(key K) (V, bool) {
	return m.Get(key)
}

// Set sets the value of the given key.
// It returns the old value if the key existed otherwise it returns the passed new value.
// The second return value is true if the key existed.
func (m *OrderedMap[K, V]) Set(key K, value V) (val V, existed bool) {
	if entry, existed := m.entries[key]; existed {
		oldValue := entry.Value
		entry.Value = value
		return oldValue, true
	}

	entry := &Entry[K, V]{
		Key:   key,
		Value: value,
	}
	entry.element = m.list.PushBack(entry)
	m.entries[key] = entry

	return value, false
}

// Store sets the value of the given key, alias for Set.
func (m *OrderedMap[K, V]) Store(key K, value V) (V, bool) {
	return m.Set(key, value)
}

// AddEntries adds entries to the map.
func (m *OrderedMap[K, V]) AddEntries(entries ...*Entry[K, V]) {
	for _, entry := range entries {
		m.Set(entry.Key, entry.Value)
	}
}

// AddOrderedMap adds entries of the given OrderedMap to the map.
func (m *OrderedMap[K, V]) AddOrderedMap(am OrderedMap[K, V]) {
	for _, entry := range am.entries {
		m.Set(entry.Key, entry.Value)
	}
}

// AddMap adds entries of the given map to the map.
func (m *OrderedMap[K, V]) AddMap(am map[K]V) {
	for key, value := range am {
		m.Set(key, value)
	}
}

// Len returns the length of the map.
func (m *OrderedMap[K, V]) Delete(key K) (val V, existed bool) {
	if entry, existed := m.entries[key]; existed {
		m.list.Remove(entry.element)
		delete(m.entries, key)
		return entry.Value, true
	}
	return
}

// Clear clears the map.
func (m *OrderedMap[K, V]) Clear() {
	m.entries = make(map[K]*Entry[K, V])
	m.list = list.New[*Entry[K, V]]()
}

// Len returns the length of the map.
func (m *OrderedMap[K, V]) Len() int {
	return len(m.entries)
}

// Oldest returns the oldest entry of the map.
func (m *OrderedMap[K, V]) Oldest() *Entry[K, V] {
	if m == nil {
		return nil
	}

	list := m.list
	if list == nil {
		return nil
	}

	e := list.Front()
	if e == nil {
		return nil
	}

	return e.Value
}

// Newest returns the newest entry of the map.
func (m *OrderedMap[K, V]) Newest() *Entry[K, V] {
	if m == nil {
		return nil
	}

	list := m.list
	if list == nil {
		return nil
	}

	e := list.Back()
	if e == nil {
		return nil
	}

	return e.Value
}

// Range calls f sequentially for each key and value in the map as insert order .
func (m *OrderedMap[K, V]) Range(f func(key K, value V) bool) {
	list := m.list
	for e := list.Front(); e != nil; e = e.Next() {
		if e.Value != nil {
			if ok := f(e.Value.Key, e.Value.Value); !ok {
				return
			}
		}
	}
}

// ForEach calls f  for each value in the map as random order like builtin map.
func (m *OrderedMap[K, V]) ForEach(f func(key K, value V) bool) {
	for key, value := range m.entries {
		if ok := f(key, value.Value); !ok {
			return
		}
	}
}

// Keys returns all keys of the map as insert order.
func (m *OrderedMap[K, V]) Keys() []K {
	r := make([]K, 0, len(m.entries))
	m.Range(func(key K, value V) bool {
		r = append(r, key)
		return true
	})

	return r
}

// Values returns all values of the map as insert order.
func (m *OrderedMap[K, V]) Values() []V {
	r := make([]V, 0, len(m.entries))
	m.Range(func(key K, value V) bool {
		r = append(r, value)
		return true
	})

	return r
}

// Clone returns a shallow copy of the map.
func (m *OrderedMap[K, V]) Clone() *OrderedMap[K, V] {
	if m == nil {
		return nil
	}

	result := NewOrderedMap[K, V](m.Len())
	result.AddOrderedMap(*m)

	return result
}

// Equal returns true if the given maps are equal to the map in order.
func Equal[K, V comparable](m1, m2 *OrderedMap[K, V]) bool {
	if len(m1.entries) != len(m2.entries) {
		return false
	}

	list1 := m1.list
	list2 := m1.list

	e1 := list1.Front()
	e2 := list2.Front()
	for e1 != nil && e2 != nil {
		if e1.Value == nil && e2.Value == nil {
			e1 = e2.Next()
			e2 = e2.Next()
			continue
		}

		if e1.Value == nil || e2.Value == nil {
			return false
		}

		if e1.Value.Key != e2.Value.Key || e1.Value.Value != e2.Value.Value {
			return false
		}

		e1 = e2.Next()
		e2 = e2.Next()
	}

	if e1 == nil || e2 == nil {
		return true
	}

	return false
}
