package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortedSet(t *testing.T) {
	// Create a set
	s := OfSortedSet(1, 2, 3)

	// Test Add
	s.Add(4)
	if !s.Contains(4) {
		t.Errorf("Expected set to contain element 4, but it does not")
	}

	// Test AddSet
	s2 := OfSortedSet(5, 6, 7)
	s.AddSet(s2)
	if !s.ContainsAll(s2) {
		t.Errorf("Expected set to contain all elements from s2, but it does not: %v", s.Values())
	}

	// Test Remove
	s.Remove(1)
	if s.Contains(1) {
		t.Errorf("Expected set to not contain element 1, but it does")
	}

	// Test RemoveSet
	s.RemoveSet(OfSortedSet(6))
	if s.Contains(6) {
		t.Errorf("Expected set to not contain element 6, but it does: %v", s.Values())
	}

	// Test ContainsAny
	if !s.ContainsAny(OfSortedSet(7, 8, 9)) {
		t.Errorf("Expected set to contain at least one element from the provided set, but it does not")
	}

	// Test ContainsAll
	if !s.ContainsAll(OfSortedSet(2, 3)) {
		t.Errorf("Expected set to contain all elements from the provided set, but it does not")
	}

	// Test Values
	values := s.Values()
	expectedValues := []int{2, 3, 4, 5, 7}
	if !s.Equal(OfSortedSet(expectedValues...)) {
		t.Errorf("Expected values %v, but got %v", expectedValues, values)
	}

	// Test Equal
	s3 := OfSortedSet(2, 3, 4, 5, 7)
	if !s.Equal(s3) {
		t.Errorf("Expected sets to be equal, but they are not")
	}

	// Test Clear
	s.Clear()
	if s.Len() != 0 {
		t.Errorf("Expected set length to be 0, but got %d", s.Len())
	}

	// Test Clone
	s.Add(1)
	s4 := s.Clone()
	if !s4.Equal(s) {
		t.Errorf("Expected cloned set to be equal to the original set, but they are not")
	}

	// Test Filter
	s.Filter(func(e int) bool {
		return e%2 == 0
	})
	if s.Contains(1) {
		t.Errorf("Expected set to not contain odd elements, but it does")
	}

	// Test Len
	if s.Len() != 0 {
		t.Errorf("Expected set length to be 0, but got %d", s.Len())
	}

	// Test Do
	sum := 0
	s = OfSortedSet(1, 2, 3)
	s.Do(func(e int) bool {
		sum += e
		return true
	})
	if sum != 6 {
		t.Errorf("Expected sum to be 6, but got %d", sum)
	}

	// Test Union
	s5 := UnionSortedSet(OfSortedSet(1, 2, 3), OfSortedSet(3, 4, 5))
	if !s5.Equal(OfSortedSet(1, 2, 3, 4, 5)) {
		t.Errorf("Expected union set to be {1, 2, 3, 4, 5}, but got %v", s5.Values())
	}

	// Test Intersection
	s6 := IntersectionSortedSet(OfSortedSet(1, 2, 3), OfSortedSet(3, 4, 5))
	if !s6.Equal(OfSortedSet(3)) {
		t.Errorf("Expected intersection set to be {3}, but got %v", s6.Values())
	}

	// Test Difference
	s7 := DifferenceSortedSet(OfSortedSet(1, 2, 3, 4, 5), OfSortedSet(3, 4))
	if !s7.Equal(OfSortedSet(1, 2, 5)) {
		t.Errorf("Expected difference set to be {1, 2, 5}, but got %v", s7.Values())
	}
}

func TestSortedSet_Order(t *testing.T) {
	s := OfSortedSet(1, 2, 3)
	s.Add(4)
	s.Add(5)
	s.Add(6)
	s.Remove(3)
	s.Remove(5)
	s.AddSet(OfSortedSet(7, 8, 9))
	s.Add(3)

	expectedValues := []int{1, 2, 4, 6, 7, 8, 9, 3}
	assert.Equal(t, expectedValues, s.Values())

	e, ok := s.Oldest()
	assert.True(t, ok)
	assert.Equal(t, 1, e)

	e, ok = s.Newest()
	assert.True(t, ok)
	assert.Equal(t, 3, e)
}
