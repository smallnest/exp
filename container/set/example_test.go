package set

import (
	"fmt"
)

func Example() {
	// Create a new set.
	s := New[int]()
	s.Add(1)
	s.AddSet(Of(2, 3))

	// Iterate through list and print its contents.
	for e := range s {
		fmt.Println(e)
	}

	// Output:
	// 1
	// 2
	// 3
}
