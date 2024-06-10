package maps

const (
	RED   = true
	BLACK = false
)

type Node[T any] struct {
	Key   int
	Value T
	Color bool
	Left  *Node[T]
	Right *Node[T]
}

type RBTree[T any] struct {
	Root *Node[T]
	Size int
}

func NewRBTree[T any]() *RBTree[T] {
	return &RBTree[T]{}
}

// Get 方法返回与给定键关联的值。如果键不存在，则返回 nil。
func (tree *RBTree[T]) Get(key int) T {
	node := tree.Root
	for node != nil {
		if key < node.Key {
			node = node.Left
		} else if key > node.Key {
			node = node.Right
		} else {
			return node.Value
		}
	}

	var t T
	return t
}

// InOrder 方法返回一个按键升序的所有节点的切片。
func (tree *RBTree[T]) InOrder() []Node[T] {
	nodes := []Node[T]{}
	stack := []*Node[T]{}

	node := tree.Root
	for node != nil || len(stack) > 0 {
		for node != nil {
			stack = append(stack, node)
			node = node.Left
		}
		node = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		nodes = append(nodes, *node)
		node = node.Right
	}
	return nodes
}

// Min 方法返回树中的最小键和对应的值。如果树为空，则返回 nil。
func (tree *RBTree[T]) Min() (int, interface{}) {
	if tree.Root == nil {
		return 0, nil
	}
	node := tree.Root
	for node.Left != nil {
		node = node.Left
	}
	return node.Key, node.Value
}

// Max 方法返回树中的最大键和对应的值。如果树为空，则返回 nil。
func (tree *RBTree[T]) Max() (int, interface{}) {
	if tree.Root == nil {
		return 0, nil
	}
	node := tree.Root
	for node.Right != nil {
		node = node.Right
	}
	return node.Key, node.Value
}

// ReverseInOrder 方法返回一个按键降序的所有节点的切片。
func (tree *RBTree[T]) ReverseInOrder() []Node[T] {
	nodes := []Node[T]{}
	stack := []*Node[T]{}

	node := tree.Root
	for node != nil || len(stack) > 0 {
		for node != nil {
			stack = append(stack, node)
			node = node.Right
		}
		node = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		nodes = append(nodes, *node)
		node = node.Left
	}
	return nodes
}

// 插入节点
func (tree *RBTree[T]) Insert(key int, value T) {
	tree.Root = tree.insert(tree.Root, key, value)
	tree.Root.Color = BLACK
	tree.Size++
}

func (tree *RBTree[T]) insert(node *Node[T], key int, value T) *Node[T] {
	if node == nil {
		return &Node[T]{Key: key, Value: value, Color: RED}
	}

	if key < node.Key {
		node.Left = tree.insert(node.Left, key, value)
	} else if key > node.Key {
		node.Right = tree.insert(node.Right, key, value)
	} else {
		node.Value = value
	}

	return tree.balance(node)
}

// 左旋
func (tree *RBTree[T]) rotateLeft(node *Node[T]) *Node[T] {
	x := node.Right
	node.Right = x.Left
	x.Left = node
	x.Color = node.Color
	node.Color = RED
	return x
}

// 右旋
func (tree *RBTree[T]) rotateRight(node *Node[T]) *Node[T] {
	x := node.Left
	node.Left = x.Right
	x.Right = node
	x.Color = node.Color
	node.Color = RED
	return x
}

// 保持红黑树的平衡
func (tree *RBTree[T]) balance(node *Node[T]) *Node[T] {
	// 如果当前节点的右子节点是红色，而左子节点是黑色，进行左旋转
	// 这是为了保证红色节点都在左侧，即红黑树的性质4：红色节点的两个子节点都是黑色
	if isRed(node.Right) && !isRed(node.Left) {
		node = tree.rotateLeft(node)
	}

	// 如果当前节点的左子节点和左子节点的左子节点都是红色，进行右旋转
	// 这是为了分解连续的红色节点，即红黑树的性质4：红色节点的两个子节点都是黑色
	if isRed(node.Left) && isRed(node.Left.Left) {
		node = tree.rotateRight(node)
	}

	// 如果当前节点的左子节点和右子节点都是红色，进行颜色翻转
	// 这是为了保证每个节点到其任何后代的所有路径都包含相同数目的黑色节点，即红黑树的性质5
	if isRed(node.Left) && isRed(node.Right) {
		flipColors(node)
	}

	return node
}

func isRed[T any](node *Node[T]) bool {
	if node == nil {
		return false
	}
	return node.Color == RED
}

func flipColors[T any](node *Node[T]) {
	node.Color = !node.Color
	if node.Left != nil {
		node.Left.Color = !node.Left.Color
	}

	if node.Right != nil {
		node.Right.Color = !node.Right.Color
	}
}

// 删除节点
func (tree *RBTree[T]) Delete(key int) {
	var deleted bool
	if !isRed(tree.Root.Left) && !isRed(tree.Root.Right) {
		tree.Root.Color = RED
	}
	tree.Root, deleted = tree.delete(tree.Root, key)
	if tree.Root != nil {
		tree.Root.Color = BLACK
	}
	if deleted {
		tree.Size--
	}
}

func (tree *RBTree[T]) delete(node *Node[T], key int) (*Node[T], bool) {
	if node == nil {
		return nil, false
	}

	var deleted bool
	if key < node.Key {
		node.Left, deleted = tree.delete(node.Left, key)
	} else {
		if isRed(node.Left) {
			node = tree.rotateRight(node)
		}
		if key == node.Key && node.Right == nil {
			return nil, true
		}
		if !isRed(node.Right) && !isRed(node.Right.Left) {
			node = tree.moveRedRight(node)
		}
		if key == node.Key {
			x := tree.min(node.Right)
			node.Key = x.Key
			node.Value = x.Value
			node.Right, _ = tree.deleteMin(node.Right)
			deleted = true
		} else {
			node.Right, deleted = tree.delete(node.Right, key)
		}
	}
	return tree.balance(node), deleted
}

func (tree *RBTree[T]) moveRedLeft(node *Node[T]) *Node[T] {
	flipColors(node)
	if node.Right != nil && isRed(node.Right.Left) {
		node.Right = tree.rotateRight(node.Right)
		node = tree.rotateLeft(node)
		flipColors(node)
	}
	return node
}

func (tree *RBTree[T]) moveRedRight(node *Node[T]) *Node[T] {
	flipColors(node)
	if node.Left != nil && isRed(node.Left.Left) {
		node = tree.rotateRight(node)
		flipColors(node)
	}
	return node
}

func (tree *RBTree[T]) deleteMin(node *Node[T]) (*Node[T], bool) {
	if node.Left == nil {
		return nil, true
	}
	var deleted bool
	if !isRed(node.Left) && !isRed(node.Left.Left) {
		node = tree.moveRedLeft(node)
	}
	node.Left, deleted = tree.deleteMin(node.Left)
	return tree.balance(node), deleted
}

func (tree *RBTree[T]) min(node *Node[T]) *Node[T] {
	if node.Left == nil {
		return node
	}
	return tree.min(node.Left)
}
