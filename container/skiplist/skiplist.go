// MIT License
//
// Copyright (c) 2018 Maurice Tollmien (maurice.tollmien@gmail.com)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package skiplist is an implementation of a skiplist to store elements in increasing order.
// It allows finding, insertion and deletion operations in approximately O(n log(n)).
// Additionally, there are methods for retrieving the next and previous element as well as changing the actual value
// without the need for re-insertion (as long as the key stays the same!)
// Skiplist is a fast alternative to a balanced tree.
package skiplist

import (
	"fmt"
	"math/bits"
	"math/rand"
	"time"

	"github.com/dolthub/maphash"
)

const (
	// maxLevel denotes the maximum height of the skiplist. This height will keep the skiplist
	// efficient for up to 34m entries. If there is a need for much more, please adjust this constant accordingly.
	maxLevel = 25
)

// Hasher is a function that calculates a hash for a given key.
type Hasher[V comparable] interface {
	Hash(key V) uint64
}

// SkipListElement represents one actual Node in the skiplist structure.
// It saves the actual element, pointers to the next nodes and a pointer to one previous node.
type SkipListElement[V comparable] struct {
	next  [maxLevel]*SkipListElement[V]
	level int
	value V
	prev  *SkipListElement[V]
}

// key calculates the key for a given element.
func getKey[V comparable](hasher Hasher[V], v V) uint64 {
	return hasher.Hash(v)
}

// SkipList is the actual skiplist representation.
// It saves all nodes accessible from the start and end and keeps track of element count, eps and levels.
type SkipList[V comparable] struct {
	startLevels  [maxLevel]*SkipListElement[V]
	endLevels    [maxLevel]*SkipListElement[V]
	maxNewLevel  int
	maxLevel     int
	elementCount int
	eps          float64

	hasher Hasher[V]
}

// NewSeed returns a new empty, initialized Skiplist.
// Given a seed, a deterministic height/list behaviour can be achieved.
func NewSeed[V comparable](seed int64) *SkipList[V] {
	// Initialize random number generator.
	rand.Seed(seed)
	//fmt.Printf("SkipList seed: %v\n", seed)

	hasher := maphash.NewHasher[V]()
	list := &SkipList[V]{
		startLevels:  [maxLevel]*SkipListElement[V]{},
		endLevels:    [maxLevel]*SkipListElement[V]{},
		maxNewLevel:  maxLevel,
		maxLevel:     0,
		elementCount: 0,
		hasher:       &hasher,
	}

	return list
}

// NewHasher returns a new empty, initialized Skiplist.
func NewHasher[V comparable](hasher Hasher[V]) *SkipList[V] {
	t := New[V]()
	t.hasher = hasher

	return t
}

// NewHashWrapper returns a new empty, initialized Skiplist.
type HashWrapper[V comparable] struct {
	hashFunc func(key V) uint64
}

// Hash returns a new empty, initialized Skiplist.
func (hw *HashWrapper[V]) Hash(key V) uint64 {
	return hw.hashFunc(key)
}

// NewHashFunc returns a new empty, initialized Skiplist.
func NewHashFunc[V comparable](hashFunc func(key V) uint64) *SkipList[V] {
	hasher := &HashWrapper[V]{hashFunc}

	return NewHasher[V](hasher)
}

// New returns a new empty, initialized Skiplist.
func New[V comparable]() *SkipList[V] {
	return NewSeed[V](time.Now().UnixNano())
}

// IsEmpty checks, if the skiplist is empty.
func (t *SkipList[V]) IsEmpty() bool {
	return t.startLevels[0] == nil
}

func (t *SkipList[V]) generateLevel(maxLevel int) int {
	level := maxLevel - 1
	// First we apply some mask which makes sure that we don't get a level
	// above our desired level. Then we find the first set bit.
	var x uint64 = rand.Uint64() & ((1 << uint(maxLevel-1)) - 1)
	zeroes := bits.TrailingZeros64(x)
	if zeroes <= maxLevel {
		level = zeroes
	}

	return level
}

func (t *SkipList[V]) findEntryIndex(key uint64, level int) int {
	// Find good entry point so we don't accidentally skip half the list.
	for i := t.maxLevel; i >= 0; i-- {
		if t.startLevels[i] != nil && getKey(t.hasher, t.startLevels[i].value) <= key || i <= level {
			return i
		}
	}
	return 0
}

func (t *SkipList[V]) findExtended(key uint64, findGreaterOrEqual bool) (foundElem *SkipListElement[V], ok bool) {

	foundElem = nil
	ok = false

	if t.IsEmpty() {
		return
	}

	index := t.findEntryIndex(key, 0)
	var currentNode *SkipListElement[V]

	currentNode = t.startLevels[index]
	nextNode := currentNode

	// In case, that our first element is already greater-or-equal!
	if findGreaterOrEqual && getKey(t.hasher, currentNode.value) > key {
		foundElem = currentNode
		ok = true
		return
	}

	for {
		if getKey(t.hasher, currentNode.value) == key {
			foundElem = currentNode
			ok = true
			return
		}

		nextNode = currentNode.next[index]

		// Which direction are we continuing next time?
		if nextNode != nil && getKey(t.hasher, nextNode.value) <= key {
			// Go right
			currentNode = nextNode
		} else {
			if index > 0 {

				// Early exit
				if currentNode.next[0] != nil && getKey(t.hasher, currentNode.next[0].value) == key {
					foundElem = currentNode.next[0]
					ok = true
					return
				}
				// Go down
				index--
			} else {
				// Element is not found and we reached the bottom.
				if findGreaterOrEqual {
					foundElem = nextNode
					ok = nextNode != nil
				}

				return
			}
		}
	}
}

// Find tries to find an element in the skiplist based on the key from the given ListElement.
// elem can be used, if ok is true.
// Find runs in approx. O(log(n))
func (t *SkipList[V]) Find(e V) (elem *SkipListElement[V], ok bool) {

	if t == nil {
		return
	}

	elem, ok = t.findExtended(getKey(t.hasher, e), false)
	return
}

// FindGreaterOrEqual finds the first element, that is greater or equal to the given ListElement e.
// The comparison is done on the keys (So on ExtractKey()).
// FindGreaterOrEqual runs in approx. O(log(n))
func (t *SkipList[V]) FindGreaterOrEqual(e V) (elem *SkipListElement[V], ok bool) {

	if t == nil {
		return
	}

	elem, ok = t.findExtended(getKey(t.hasher, e), true)
	return
}

// Delete removes an element equal to e from the skiplist, if there is one.
// If there are multiple entries with the same value, Delete will remove one of them
// (Which one will change based on the actual skiplist layout)
// Delete runs in approx. O(log(n))
func (t *SkipList[V]) Delete(e V) {

	if t == nil || t.IsEmpty() {
		return
	}

	key := getKey(t.hasher, e)

	index := t.findEntryIndex(key, 0)

	var currentNode *SkipListElement[V]
	nextNode := currentNode

	for {

		if currentNode == nil {
			nextNode = t.startLevels[index]
		} else {
			nextNode = currentNode.next[index]
		}

		// Found and remove!
		if nextNode != nil && getKey(t.hasher, nextNode.value) == key {

			if currentNode != nil {
				currentNode.next[index] = nextNode.next[index]
			}

			if index == 0 {
				if nextNode.next[index] != nil {
					nextNode.next[index].prev = currentNode
				}
				t.elementCount--
			}

			// Link from start needs readjustments.
			if t.startLevels[index] == nextNode {
				t.startLevels[index] = nextNode.next[index]
				// This was our currently highest node!
				if t.startLevels[index] == nil {
					t.maxLevel = index - 1
				}
			}

			// Link from end needs readjustments.
			if nextNode.next[index] == nil {
				t.endLevels[index] = currentNode
			}
			nextNode.next[index] = nil
		}

		if nextNode != nil && getKey(t.hasher, nextNode.value) < key {
			// Go right
			currentNode = nextNode
		} else {
			// Go down
			index--
			if index < 0 {
				break
			}
		}
	}

}

// Insert inserts the given ListElement into the skiplist.
// Insert runs in approx. O(log(n))
func (t *SkipList[V]) Insert(e V) {

	if t == nil {
		return
	}

	level := t.generateLevel(t.maxNewLevel)

	// Only grow the height of the skiplist by one at a time!
	if level > t.maxLevel {
		level = t.maxLevel + 1
		t.maxLevel = level
	}

	elem := &SkipListElement[V]{
		next:  [maxLevel]*SkipListElement[V]{},
		level: level,
		value: e,
	}

	t.elementCount++

	newFirst := true
	newLast := true

	elemKey := getKey(t.hasher, elem.value)
	if !t.IsEmpty() {
		newFirst = elemKey < getKey(t.hasher, t.startLevels[0].value)
		newLast = elemKey > getKey(t.hasher, t.endLevels[0].value)
	}

	normallyInserted := false
	if !newFirst && !newLast {

		normallyInserted = true

		index := t.findEntryIndex(elemKey, level)

		var currentNode *SkipListElement[V]
		nextNode := t.startLevels[index]

		for {

			if currentNode == nil {
				nextNode = t.startLevels[index]
			} else {
				nextNode = currentNode.next[index]
			}

			// Connect node to next
			if index <= level && (nextNode == nil || getKey(t.hasher, nextNode.value) > elemKey) {
				elem.next[index] = nextNode
				if currentNode != nil {
					currentNode.next[index] = elem
				}
				if index == 0 {
					elem.prev = currentNode
					if nextNode != nil {
						nextNode.prev = elem
					}
				}
			}

			if nextNode != nil && getKey(t.hasher, nextNode.value) <= elemKey {
				// Go right
				currentNode = nextNode
			} else {
				// Go down
				index--
				if index < 0 {
					break
				}
			}
		}
	}

	// Where we have a left-most position that needs to be referenced!
	for i := level; i >= 0; i-- {

		didSomething := false

		if newFirst || normallyInserted {

			if t.startLevels[i] == nil || getKey(t.hasher, t.startLevels[i].value) > elemKey {
				if i == 0 && t.startLevels[i] != nil {
					t.startLevels[i].prev = elem
				}
				elem.next[i] = t.startLevels[i]
				t.startLevels[i] = elem
			}

			// link the endLevels to this element!
			if elem.next[i] == nil {
				t.endLevels[i] = elem
			}

			didSomething = true
		}

		if newLast {
			// Places the element after the very last element on this level!
			// This is very important, so we are not linking the very first element (newFirst AND newLast) to itself!
			if !newFirst {
				if t.endLevels[i] != nil {
					t.endLevels[i].next[i] = elem
				}
				if i == 0 {
					elem.prev = t.endLevels[i]
				}
				t.endLevels[i] = elem
			}

			// Link the startLevels to this element!
			if t.startLevels[i] == nil || getKey(t.hasher, t.startLevels[i].value) > elemKey {
				t.startLevels[i] = elem
			}

			didSomething = true
		}

		if !didSomething {
			break
		}
	}
}

// GetValue extracts the ListElement value from a skiplist node.
func (e *SkipListElement[V]) GetValue() V {
	if e == nil {
		var zero V
		return zero
	}
	return e.value
}

// GetSmallestNode returns the very first/smallest node in the skiplist.
// GetSmallestNode runs in O(1)
func (t *SkipList[V]) GetSmallestNode() *SkipListElement[V] {
	return t.startLevels[0]
}

// GetLargestNode returns the very last/largest node in the skiplist.
// GetLargestNode runs in O(1)
func (t *SkipList[V]) GetLargestNode() *SkipListElement[V] {
	return t.endLevels[0]
}

// Next returns the next element based on the given node.
// Next will loop around to the first node, if you call it on the last!
func (t *SkipList[V]) Next(e *SkipListElement[V]) *SkipListElement[V] {
	if e.next[0] == nil {
		return t.startLevels[0]
	}
	return e.next[0]
}

// Prev returns the previous element based on the given node.
// Prev will loop around to the last node, if you call it on the first!
func (t *SkipList[V]) Prev(e *SkipListElement[V]) *SkipListElement[V] {
	if e.prev == nil {
		return t.endLevels[0]
	}
	return e.prev
}

// GetNodeCount returns the number of nodes currently in the skiplist.
func (t *SkipList[V]) GetNodeCount() int {
	return t.elementCount
}

// ChangeValue can be used to change the actual value of a node in the skiplist
// without the need of Deleting and reinserting the node again.
// Be advised, that ChangeValue only works, if the actual key from ExtractKey() will stay the same!
// ok is an indicator, wether the value is actually changed.
func (t *SkipList[V]) ChangeValue(oldValue V, newValue V) (ok bool) {
	e, ok := t.Find(oldValue)
	if !ok || e == nil {
		return false
	}

	e.value = newValue

	return true
}

// String returns a string format of the skiplist. Useful to get a graphical overview and/or debugging.
func (t *SkipList[V]) String() string {
	s := ""

	s += " --> "
	for i, l := range t.startLevels {
		if l == nil {
			break
		}
		if i > 0 {
			s += " -> "
		}
		next := "---"
		if l != nil {
			next = fmt.Sprintf("%+v", l.value)
		}
		s += fmt.Sprintf("[%v]", next)

		if i == 0 {
			s += "    "
		}
	}
	s += "\n"

	node := t.startLevels[0]
	for node != nil {
		s += fmt.Sprintf("%v: ", node.value)
		for i := 0; i <= node.level; i++ {

			l := node.next[i]

			next := "---"
			if l != nil {
				next = fmt.Sprintf("%+v", l.value)
			}

			if i == 0 {
				prev := "---"
				if node.prev != nil {
					prev = fmt.Sprintf("%+v", node.prev.value)
				}
				s += fmt.Sprintf("[%v|%v]", prev, next)
			} else {
				s += fmt.Sprintf("[%v]", next)
			}
			if i < node.level {
				s += " -> "
			}

		}
		s += "\n"
		node = node.next[0]
	}

	s += " --> "
	for i, l := range t.endLevels {
		if l == nil {
			break
		}
		if i > 0 {
			s += " -> "
		}
		next := "---"
		if l != nil {
			next = fmt.Sprintf("%+v", l.value)
		}
		s += fmt.Sprintf("[%v]", next)
		if i == 0 {
			s += "    "
		}
	}
	s += "\n"
	return s
}
