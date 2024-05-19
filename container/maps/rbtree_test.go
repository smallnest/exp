package maps

import (
	"testing"
)

func TestRBTree(t *testing.T) {
	tree := &RBTree{}

	t.Run("Insert and Size", func(t *testing.T) {
		tree.Insert(5, "five")
		tree.Insert(3, "three")
		tree.Insert(7, "seven")

		if tree.Size != 3 {
			t.Errorf("expected Size to be 3, got %v", tree.Size)
		}

		if val := tree.Root.Value; val != "five" {
			t.Errorf("expected Root Value to be 'five', got %v", val)
		}

		if val := tree.Root.Left.Value; val != "three" {
			t.Errorf("expected Left Value to be 'three', got %v", val)
		}

		if val := tree.Root.Right.Value; val != "seven" {
			t.Errorf("expected Right Value to be 'seven', got %v", val)
		}
	})
}

func TestInOrder(t *testing.T) {
	tree := NewRBTree()
	tree.Insert(3, "three")
	tree.Insert(1, "one")
	tree.Insert(2, "two")
	tree.Insert(4, "four")

	nodes := tree.InOrder()
	if len(nodes) != 4 {
		t.Errorf("Expected 4 nodes, got %d", len(nodes))
	}

	expectedKeys := []int{1, 2, 3, 4}
	for i, node := range nodes {
		if node.Key != expectedKeys[i] {
			t.Errorf("Expected key %d, got %d", expectedKeys[i], node.Key)
		}
	}
}

func TestReverseInOrder(t *testing.T) {
	tree := NewRBTree()
	tree.Insert(3, "three")
	tree.Insert(1, "one")
	tree.Insert(2, "two")
	tree.Insert(4, "four")

	nodes := tree.ReverseInOrder()
	if len(nodes) != 4 {
		t.Errorf("Expected 4 nodes, got %d", len(nodes))
	}

	expectedKeys := []int{4, 3, 2, 1}
	for i, node := range nodes {
		if node.Key != expectedKeys[i] {
			t.Errorf("Expected key %d, got %d", expectedKeys[i], node.Key)
		}
	}
}

func TestGet(t *testing.T) {
	tree := NewRBTree()
	tree.Insert(3, "three")
	tree.Insert(1, "one")
	tree.Insert(2, "two")
	tree.Insert(4, "four")

	if value := tree.Get(1); value != "one" {
		t.Errorf("Expected 'one', got %v", value)
	}

	if value := tree.Get(2); value != "two" {
		t.Errorf("Expected 'two', got %v", value)
	}

	if value := tree.Get(3); value != "three" {
		t.Errorf("Expected 'three', got %v", value)
	}

	if value := tree.Get(4); value != "four" {
		t.Errorf("Expected 'four', got %v", value)
	}

	if value := tree.Get(5); value != nil {
		t.Errorf("Expected nil, got %v", value)
	}
}

func TestInsert(t *testing.T) {
	tree := NewRBTree()
	tree.Insert(3, "three")
	tree.Insert(1, "one")
	tree.Insert(2, "two")
	tree.Insert(4, "four")

	if value := tree.Get(1); value != "one" {
		t.Errorf("Expected 'one', got %v", value)
	}

	if value := tree.Get(2); value != "two" {
		t.Errorf("Expected 'two', got %v", value)
	}

	if value := tree.Get(3); value != "three" {
		t.Errorf("Expected 'three', got %v", value)
	}

	if value := tree.Get(4); value != "four" {
		t.Errorf("Expected 'four', got %v", value)
	}

	if value := tree.Get(5); value != nil {
		t.Errorf("Expected nil, got %v", value)
	}
}

func TestDelete(t *testing.T) {
	tree := NewRBTree()
	tree.Insert(3, "three")
	tree.Insert(1, "one")
	tree.Insert(2, "two")
	tree.Insert(4, "four")

	tree.Delete(1)
	if value := tree.Get(1); value != nil {
		t.Errorf("Expected nil, got %v", value)
	}

	tree.Delete(2)
	if value := tree.Get(2); value != nil {
		t.Errorf("Expected nil, got %v", value)
	}

	tree.Delete(3)
	if value := tree.Get(3); value != nil {
		t.Errorf("Expected nil, got %v", value)
	}

	tree.Delete(4)
	if value := tree.Get(4); value != nil {
		t.Errorf("Expected nil, got %v", value)
	}
}

func FuzzInsert(f *testing.F) {
	f.Add(3, "three")
	f.Add(1, "one")
	f.Add(2, "two")
	f.Add(4, "four")

	f.Fuzz(func(t *testing.T, key int, value string) {
		tree := NewRBTree()
		tree.Insert(key, value)
		if val := tree.Get(key); val != value {
			t.Errorf("Expected %v, got %v", value, val)
		}
	})
}

func FuzzDelete(f *testing.F) {
	f.Add(3, "three")
	f.Add(1, "one")
	f.Add(2, "two")
	f.Add(4, "four")

	f.Fuzz(func(t *testing.T, key int, value string) {
		tree := NewRBTree()
		tree.Insert(key, value)
		tree.Delete(key)
		if val := tree.Get(key); val != nil {
			t.Errorf("Expected nil, got %v", val)
		}
	})
}
