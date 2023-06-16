package maps

import (
	"encoding/json"
)

// BidiMap is a bidirectional map, which efines a map that allows bidirectional lookup between key and values.
// This map  represents a mapping where a key may lookup a value and a value may lookup a key with equal ease.
// Its key and value types must be comparable, and it enforces the restriction that there is a 1:1 relation between keys and values,
// meaning that multiple keys cannot map to the same value.
type BidiMap[K, V comparable] struct {
	kvMap map[K]V
	vkMap map[V]K
}

// NewBidiMap returns a new bidirectional map with the given capacity.
func NewBidiMap[K, V comparable](capacity int) *BidiMap[K, V] {
	return &BidiMap[K, V]{
		kvMap: make(map[K]V, capacity),
		vkMap: make(map[V]K, capacity),
	}
}

// Get returns the value associated with the given key, and whether it existed or not.
func (bm BidiMap[K, V]) Get(key K) (value V, found bool) {
	value, found = bm.kvMap[key]
	return
}

// GetKey returns the key associated with the given value, and whether it existed or not.
func (bm BidiMap[K, V]) GetKey(value V) (key K, found bool) {
	key, found = bm.vkMap[value]
	return
}

// Put associates the given key with the given value.
// If the key or value already existed, it will be overwritten.
func (bm *BidiMap[K, V]) Put(key K, value V) {
	if v, found := bm.kvMap[key]; found {
		delete(bm.vkMap, v)
	}
	bm.kvMap[key] = value
	bm.vkMap[value] = key
}

// Remove removes the given key and its associated value.
// If the key does not exist, it does nothing.
func (bm *BidiMap[K, V]) Remove(key K) {
	if value, found := bm.kvMap[key]; found {
		delete(bm.kvMap, key)
		delete(bm.vkMap, value)
	}
}

// RemoveValue removes the given value and its associated key.
// If the value does not exist, it does nothing.
func (bm *BidiMap[K, V]) RemoveValue(value V) {
	if key, found := bm.vkMap[value]; found {
		delete(bm.vkMap, value)
		delete(bm.kvMap, key)
	}
}

// Clear removes all the keys and values.
func (bm *BidiMap[K, V]) Clear() {
	bm.kvMap = make(map[K]V)
	bm.vkMap = make(map[V]K)
}

// Len returns the number of keys.
func (bm BidiMap[K, V]) Len() int {
	return len(bm.kvMap)
}

// Keys returns a slice of all the keys.
func (bm BidiMap[K, V]) Keys() []K {
	keys := make([]K, 0, len(bm.kvMap))
	for key := range bm.kvMap {
		keys = append(keys, key)
	}

	return keys
}

// Values returns a slice of all the values.
func (bm BidiMap[K, V]) Values() []V {
	values := make([]V, 0, len(bm.vkMap))
	for value := range bm.vkMap {
		values = append(values, value)
	}

	return values
}

// Range calls the given function for each key/value pair.
// If the function returns false, it stops the iteration.
func (bm BidiMap[K, V]) Range(f func(key K, value V) bool) {
	for key, value := range bm.kvMap {
		if !f(key, value) {
			return
		}
	}
}

// Clone returns a shallow copy of the map.
func (bm *BidiMap[K, V]) Clone() *BidiMap[K, V] {
	if bm == nil {
		return nil
	}

	result := NewBidiMap[K, V](bm.Len())
	bm.Range(func(key K, value V) bool {
		result.Put(key, value)

		return true
	})

	return result
}

// AddBidiMap adds all the key/value pairs from the given map.
func (bm *BidiMap[K, V]) AddBidiMap(m *BidiMap[K, V]) {
	for key, value := range m.kvMap {
		bm.Put(key, value)
	}
}

// AddMap adds all the key/value pairs from the given map.
// You must guarantee that the map has a 1:1 relation between keys and values.
func (bm *BidiMap[K, V]) AddMap(m map[K]V) {
	for key, value := range m {
		bm.Put(key, value)
	}
}

// Equals returns true if the given map is equal to this map.
func (bm BidiMap[K, V]) Equals(m *BidiMap[K, V]) bool {
	if bm.Len() != m.Len() {
		return false
	}

	for key, value := range bm.kvMap {
		if mValue, found := m.kvMap[key]; !found || mValue != value {
			return false
		}
	}

	return true
}

// MarshalJSON marshals the map into JSON.
func (m *BidiMap[K, V]) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.kvMap)
}

// UnmarshalJSON unmarshals the map from JSON.
func (m *BidiMap[K, V]) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &m.kvMap)
	if err != nil {
		return err
	}

	m.vkMap = make(map[V]K, len(m.kvMap))
	for k, v := range m.kvMap {
		m.vkMap[v] = k
	}

	return nil
}
