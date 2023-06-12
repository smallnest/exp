package list_test

import (
	"fmt"
	"log"

	"github.com/smallnest/exp/container/list"
)

func Example() {
	// Create a new list and put some numbers in it.
	l := list.New[int]()
	e4 := l.PushBack(4)
	log.Println("Pushed 4 to back")
	e1 := l.PushFront(1)
	log.Println("Pushed 1 to front")
	l.InsertBefore(3, e4)
	log.Println("Inserted 3 before 4")
	l.InsertAfter(2, e1)
	log.Println("Inserted 2 after 1")

	// Iterate through list and print its contents.
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}

	// Output:
	// 1
	// 2
	// 3
	// 4
}
