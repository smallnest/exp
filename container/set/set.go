// Package set defines a Set type that holds a set of Eents.
package set

// discussion at https://github.com/golang/go/discussions/47331

// Set is a set of Eents of some comparable type.
// Sets are implemented using maps, and have similar performance characteristics.
// Like maps, Sets are reference types.
// That is, for Sets s1 = s2 will leave s1 and s2 pointing to the same set of Eents:
// changes to s1 will be reflected in s2 and vice-versa.
type Set[E comparable] map[E]struct{}

// New returns a new set.
func New[E comparable]() Set[E] {
	return make(Set[E])
}

// Of returns a new set containing the listed Eents.
func Of[E comparable](v ...E) Set[E] {
	s := make(Set[E])
	s.Add(v...)
	return s
}

// Add adds Eents to a set.
func (s Set[E]) Add(v ...E) {
	for _, e := range v {
		s[e] = struct{}{}
	}
}

// AddSet adds the Eents of set s2 to s.
func (s Set[E]) AddSet(s2 Set[E]) {
	for e := range s2 {
		s[e] = struct{}{}
	}
}

// Remove removes Eents from a set.
// Eents that are not present are ignored.
func (s Set[E]) Remove(v ...E) {
	for _, e := range v {
		delete(s, e)
	}
}

// RemoveSet removes the Eents of set s2 from s.
// Eents present in s2 but not s are ignored.
func (s Set[E]) RemoveSet(s2 Set[E]) {
	for e := range s2 {
		delete(s, e)
	}
}

// Contains reports whether v is in the set.
func (s Set[E]) Contains(v E) bool {
	_, ok := s[v]
	return ok
}

// ContainsAny reports whether any of the Eents in s2 are in s.
func (s Set[E]) ContainsAny(s2 Set[E]) bool {
	for e := range s2 {
		if s.Contains(e) {
			return true
		}
	}
	return false
}

// ContainsAll reports whether all of the Eents in s2 are in s.
func (s Set[E]) ContainsAll(s2 Set[E]) bool {
	for e := range s2 {
		if !s.Contains(e) {
			return false
		}
	}
	return true
}

// Values returns the Eents in the set s as a slice.
// The values will be in an indeterminate order.
func (s Set[E]) Values() []E {
	v := make([]E, 0, len(s))
	for e := range s {
		v = append(v, e)
	}
	return v
}

// Equal reports whether s and s2 contain the same Eents.
func (s Set[E]) Equal(s2 Set[E]) bool {
	if len(s) != len(s2) {
		return false
	}
	for e := range s {
		if !s2.Contains(e) {
			return false
		}
	}
	return true
}

// Clear removes all Eents from s, leaving it empty.
func (s *Set[E]) Clear() {
	// clears
	*s = make(Set[E])
}

// Clone returns a copy of s.
// The Eents are copied using assignment,
// so this is a shallow clone.
func (s Set[E]) Clone() Set[E] {
	s2 := make(Set[E])
	s2.AddSet(s)
	return s2
}

// Filter deletes any Eents from s for which keep returns false.
func (s Set[E]) Filter(keep func(E) bool) {
	for e := range s {
		if !keep(e) {
			delete(s, e)
		}
	}
}

// Len returns the number of Eents in s.
func (s Set[E]) Len() int {
	return len(s)
}

// Do calls f on every Eent in the set s,
// stopping if f returns false.
// f should not change s.
// f will be called on values in an indeterminate order.
func (s Set[E]) Do(f func(E) bool) {
	for e := range s {
		if !f(e) {
			break
		}
	}
}

// Union constructs a new set containing the union of s1 and s2.
func Union[E comparable](s1, s2 Set[E]) Set[E] {
	s := make(Set[E])
	s.AddSet(s1)
	s.AddSet(s2)
	return s
}

// Intersection constructs a new set containing the intersection of s1 and s2.
func Intersection[E comparable](s1, s2 Set[E]) Set[E] {
	s := make(Set[E])
	for e := range s1 {
		if s2.Contains(e) {
			s.Add(e)
		}
	}
	return s
}

// Difference constructs a new set containing the Eents of s1 that
// are not present in s2.
func Difference[E comparable](s1, s2 Set[E]) Set[E] {
	s := make(Set[E])
	for e := range s1 {
		if !s2.Contains(e) {
			s.Add(e)
		}
	}
	return s
}
