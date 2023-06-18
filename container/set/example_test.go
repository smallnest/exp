package set

import (
	"fmt"
)

func Example() {
	// Create a new set.
	s := NewSet[int]()
	s.Add(1)
	s.AddSet(OfSet(2, 3))

	// Iterate through list and print its contents.
	for i := 1; i < 5; i++ {
		fmt.Println(s.Contains(i))
	}

	// Output:
	// true
	// true
	// true
	// false
}
