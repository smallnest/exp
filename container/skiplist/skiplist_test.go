package skiplist

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	maxN = 1000000
)

type Element int

func (e Element) ExtractKey() float64 {
	return float64(e)
}
func (e Element) String() string {
	return fmt.Sprintf("%03d", e)
}

type FloatElement float64

func (e FloatElement) ExtractKey() float64 {
	return float64(e)
}
func (e FloatElement) String() string {
	return fmt.Sprintf("%.3f", e)
}

type ComplexElement struct {
	E int
	S string
}

func (e ComplexElement) ExtractKey() float64 {
	return float64(e.E)
}
func (e ComplexElement) String() string {
	return fmt.Sprintf("%03d", e.E)
}

// timeTrack will print out the number of nanoseconds since the start time divided by n
// Useful for printing out how long each iteration took in a benchmark
func timeTrack(start time.Time, n int, name string) {
	loopNS := time.Since(start).Nanoseconds() / int64(n)
	fmt.Printf("%s: %d\n", name, loopNS)
}

func TestInsertAndFind(t *testing.T) {
	var list *SkipList[int, int]

	var listPointer *SkipList[int, int]
	listPointer.Insert(0, 0)
	if _, ok := listPointer.Find(0); ok {
		assert.Fail(t, "listPointer should be nil")
	}

	list = New[int, int]()

	if _, ok := list.Find(0); ok {
		assert.Fail(t, "list should be empty")
	}
	if !list.IsEmpty() {
		assert.Fail(t, "list should be empty")
	}

	// Test at the beginning of the list.
	for i := 0; i < maxN; i++ {
		list.Insert(maxN-i, maxN-i)
	}
	for i := 0; i < maxN; i++ {
		if _, ok := list.Find(maxN - i); !ok {
			assert.Fail(t, "list should contain", strconv.Itoa(maxN-i))
		}
	}

	list = New[int, int]()
	// Test at the end of the list.
	for i := 0; i < maxN; i++ {
		list.Insert(i, i)
	}
	for i := 0; i < maxN; i++ {
		if _, ok := list.Find(i); !ok {
			assert.Fail(t, "list should contain", strconv.Itoa(i))
		}
	}

	list = New[int, int]()
	// Test at random positions in the list.
	rList := rand.Perm(maxN)
	for _, e := range rList {
		list.Insert(e, e)
	}
	for _, e := range rList {
		if _, ok := list.Find(e); !ok {
			assert.Fail(t, "list should contain", strconv.Itoa(e))
		}
	}

}

func TestDelete(t *testing.T) {

	var list *SkipList[int, int]

	// Delete on empty list
	list.Delete(0)

	list = New[int, int]()

	list.Delete(0)
	if !list.IsEmpty() {
		assert.Fail(t, "list should be empty")
	}

	// Delete elements at the beginning of the list.
	for i := 0; i < maxN; i++ {
		list.Insert(i, i)
	}
	for i := 0; i < maxN; i++ {
		list.Delete(i)
	}
	if !list.IsEmpty() {
		assert.Fail(t, "list should be empty")
	}

	list = New[int, int]()
	// Delete elements at the end of the list.
	for i := 0; i < maxN; i++ {
		list.Insert(i, i)
	}
	for i := 0; i < maxN; i++ {
		list.Delete(maxN - i - 1)
	}
	if !list.IsEmpty() {
		assert.Fail(t, "list should be empty")
	}

	list = New[int, int]()
	// Delete elements at random positions in the list.
	rList := rand.Perm(maxN)
	for _, e := range rList {
		list.Insert(e, e)
	}
	for _, e := range rList {
		list.Delete(e)
	}
	if !list.IsEmpty() {
		assert.Fail(t, "list should be empty")
	}
}

func TestPrev(t *testing.T) {
	list := New[int, int]()

	for i := 0; i < maxN; i++ {
		list.Insert(i, i)
	}

	smallest := list.GetSmallestNode()
	largest := list.GetLargestNode()

	lastNode := largest
	node := lastNode
	for node != smallest {
		node = list.Prev(node)
		// Must always be incrementing here!
		if node.value >= lastNode.value {
			assert.GreaterOrEqual(t, node.value, lastNode.value)
		}
		// Next.Prev must always point to itself!
		if list.Prev(list.Next(node)) != node {
			assert.NotEqual(t, list.Prev(list.Next(node)), node)
		}
		lastNode = node
	}

	if list.Prev(smallest) != largest {
		assert.NotEqual(t, list.Prev(smallest), largest)
	}
}

func TestNext(t *testing.T) {
	list := New[int, int]()

	for i := 0; i < maxN; i++ {
		list.Insert(i, i)
	}

	smallest := list.GetSmallestNode()
	largest := list.GetLargestNode()

	lastNode := smallest
	node := lastNode
	for node != largest {
		node = list.Next(node)
		// Must always be incrementing here!
		if node.value <= lastNode.value {
			assert.LessOrEqual(t, node.value, lastNode.value)
		}
		// Next.Prev must always point to itself!
		if list.Next(list.Prev(node)) != node {
			assert.NotEqual(t, list.Next(list.Prev(node)), node)
		}
		lastNode = node
	}

	if list.Next(largest) != smallest {
		assert.NotEqual(t, list.Next(largest), smallest)
	}
}

func TestChangeValue(t *testing.T) {
	list := New[int, ComplexElement]()

	for i := 0; i < maxN; i++ {
		list.Insert(i, ComplexElement{i, strconv.Itoa(i)})
	}

	for i := 0; i < maxN; i++ {
		// The key only looks at the int so the string doesn't matter here!
		_, ok := list.Find(i)
		if !ok {
			assert.Fail(t, "could not find element")
		}
		ok = list.ChangeValue(i, ComplexElement{i, "different value"})
		if !ok {
			assert.Fail(t, "could not change value")
		}
		f2, ok := list.Find(i)
		if !ok {
			assert.Fail(t, "could not find element")
		}
		if f2.S != "different value" {
			assert.NotEqual(t, f2.S, "different value")
		}
	}
}

func TestGetNodeCount(t *testing.T) {
	list := New[int, int]()

	for i := 0; i < maxN; i++ {
		list.Insert(i, i)
	}

	assert.Equal(t, list.GetNodeCount(), maxN)
}

func TestInfiniteLoop(t *testing.T) {
	list := New[int, int]()
	list.Insert(1, 1)

	if _, ok := list.Find(2); ok {
		assert.Fail(t, "list should not contain 2")
	}

	if _, ok := list.FindGreaterOrEqual(2); ok {
		assert.Fail(t, "list should not contain 2")
	}
}
