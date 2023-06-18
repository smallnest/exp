// Package set defines a Set type that holds a set of elements.
package set

import (
	"github.com/smallnest/exp/container/maps"
)

// Set is a set of elements of some comparable type in insert order.
type SortedSet[E comparable] maps.OrderedMap[E, struct{}]

// NewSortedSet returns a new SortedSet.
func NewSortedSet[E comparable](capability int) *SortedSet[E] {
	return (*SortedSet[E])(maps.NewOrderedMap[E, struct{}](capability))
}

// OfSortedSet returns a new SortedSet containing the listed elements.
func OfSortedSet[E comparable](v ...E) *SortedSet[E] {
	s := NewSortedSet[E](len(v))
	s.Add(v...)
	return s
}

// Add adds elements to a set.
func (s *SortedSet[E]) Add(v ...E) {
	m := (*maps.OrderedMap[E, struct{}])(s)
	for _, e := range v {
		m.Set(e, struct{}{})
	}
}

// AddSet adds the elements of set s2 to s.
func (s *SortedSet[E]) AddSet(s2 *SortedSet[E]) {
	m := (*maps.OrderedMap[E, struct{}])(s2)
	m.Range(func(key E, value struct{}) bool {
		s.Add(key)
		return true
	})
}

// Remove removes elements from a set.
// elements that are not present are ignored.
func (s *SortedSet[E]) Remove(v ...E) {
	m := (*maps.OrderedMap[E, struct{}])(s)
	for _, e := range v {
		m.Delete(e)
	}
}

// RemoveSet removes the elements of set s2 from s.
// elements present in s2 but not s are ignored.
func (s *SortedSet[E]) RemoveSet(s2 *SortedSet[E]) {
	m := (*maps.OrderedMap[E, struct{}])(s2)

	m.Range(func(key E, value struct{}) bool {
		if s.Contains(key) {
			s.Remove(key)
		}
		return true
	})
}

// Contains reports whether v is in the set.
func (s *SortedSet[E]) Contains(v E) bool {
	m := (*maps.OrderedMap[E, struct{}])(s)
	_, ok := m.Get(v)
	return ok
}

// ContainsAny reports whether any of the elements in s2 are in s.
func (s *SortedSet[E]) ContainsAny(s2 *SortedSet[E]) bool {
	m := (*maps.OrderedMap[E, struct{}])(s2)

	existed := false
	m.Range(func(key E, value struct{}) bool {
		if s.Contains(key) {
			existed = true
			return true
		}
		return false
	})

	return existed
}

// ContainsAll reports whether all of the elements in s2 are in s.
func (s *SortedSet[E]) ContainsAll(s2 *SortedSet[E]) bool {
	m := (*maps.OrderedMap[E, struct{}])(s2)

	existed := true
	m.Range(func(key E, value struct{}) bool {
		if s.Contains(key) {
			return true
		}

		existed = false
		return false
	})

	return existed
}

// Values returns the elements in the set s as a slice.
// The values will be in an indeterminate order.
func (s *SortedSet[E]) Values() []E {
	m := (*maps.OrderedMap[E, struct{}])(s)

	v := make([]E, 0, m.Len())
	m.Range(func(key E, value struct{}) bool {
		v = append(v, key)
		return true
	})

	return v
}

// Equal reports whether s and s2 contain the same elements.
func (s *SortedSet[E]) Equal(s2 *SortedSet[E]) bool {
	m := (*maps.OrderedMap[E, struct{}])(s)
	m2 := (*maps.OrderedMap[E, struct{}])(s2)

	if m.Len() != m2.Len() {
		return false
	}

	notEqual := true
	values1 := s.Values()
	values2 := s2.Values()
	for i, v := range values1 {
		if v != values2[i] {
			notEqual = false
			break
		}
	}

	return notEqual
}

// Clear removes all elements from s, leaving it empty.
func (s *SortedSet[E]) Clear() {
	m := (*maps.OrderedMap[E, struct{}])(s)
	m.Clear()
}

// Clone returns a copy of s.
// The elements are copied using assignment,
// so this is a shallow clone.
func (s *SortedSet[E]) Clone() *SortedSet[E] {
	s2 := NewSortedSet[E](s.Len())
	s2.AddSet(s)

	return s2

}

// Filter deletes any elements from s for which keep returns false.
func (s *SortedSet[E]) Filter(keep func(E) bool) {
	m := (*maps.OrderedMap[E, struct{}])(s)
	m.Range(func(key E, value struct{}) bool {
		if !keep(key) {
			m.Delete(key)
		}
		return true
	})
}

// Len returns the number of elements in s.
func (s *SortedSet[E]) Len() int {
	m := (*maps.OrderedMap[E, struct{}])(s)
	return m.Len()
}

// Do calls f on every Eent in the set s,
// stopping if f returns false.
// f should not change s.
// f will be called on values in an indeterminate order.
func (s *SortedSet[E]) Do(f func(E) bool) {
	m := (*maps.OrderedMap[E, struct{}])(s)
	m.Range(func(key E, value struct{}) bool {
		return f(key)
	})
}

// Oldest returns the oldest element in the set.
func (s *SortedSet[E]) Oldest() (E, bool) {
	m := (*maps.OrderedMap[E, struct{}])(s)
	entry := m.Oldest()
	if entry == nil {
		var e E
		return e, false
	}

	return entry.Key, true
}

// Newest returns the newest element in the set.
func (s *SortedSet[E]) Newest() (E, bool) {
	m := (*maps.OrderedMap[E, struct{}])(s)
	entry := m.Newest()
	if entry == nil {
		var e E
		return e, false
	}

	return entry.Key, true
}

// UnionSortedSet constructs a new set containing the union of s1 and s2.
func UnionSortedSet[E comparable](s1, s2 *SortedSet[E]) *SortedSet[E] {
	s := NewSortedSet[E](0)
	s.AddSet(s1)
	s.AddSet(s2)

	return s
}

// IntersectionSortedSet constructs a new set containing the intersection of s1 and s2.
func IntersectionSortedSet[E comparable](s1, s2 *SortedSet[E]) *SortedSet[E] {
	s := NewSortedSet[E](0)

	m := (*maps.OrderedMap[E, struct{}])(s1)
	m.Range(func(key E, value struct{}) bool {
		if s2.Contains(key) {
			s.Add(key)
		}
		return true
	})

	return s
}

// DifferenceSortedSet constructs a new set containing the elements of s1 that
// are not present in s2.
func DifferenceSortedSet[E comparable](s1, s2 *SortedSet[E]) *SortedSet[E] {
	s := NewSortedSet[E](0)

	m := (*maps.OrderedMap[E, struct{}])(s1)
	m.Range(func(key E, value struct{}) bool {
		if !s2.Contains(key) {
			s.Add(key)
		}
		return true
	})

	return s
}
