package skiplist

func abs(n uint64) uint64 {
	y := n >> 63
	return (n ^ y) - y
}
