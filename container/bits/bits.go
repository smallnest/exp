package bits

import (
	"errors"
	"math/bits"
)

// ErrNotPowerOfTwo is returned when the size of the bitset is not a power of two.
var ErrNotPowerOfTwo = errors.New("size must be a power of two")

// Bits is a simple bitset implementation, which is a fixed-size array of bits.
// It supports setting, clearing, and checking the value of a bit at a given position.
// It also supports counting the number of 1s in the bitset and shifting all bits to the left or right.
type Bits struct {
	data []uint64
	size int
}

// NewBits creates a new bitset with the given size.
func NewBits(size int) (*Bits, error) {
	if !isPowerOfTwo(size) {
		return nil, ErrNotPowerOfTwo
	}

	numElements := (size + 63) / 64
	return &Bits{
		data: make([]uint64, numElements),
		size: size,
	}, nil
}

// isPowerOfTwo checks if the given number is a power of two.
func isPowerOfTwo(n int) bool {
	if n <= 0 {
		return false
	}
	return (n & (n - 1)) == 0
}

// SetBit sets the bit at the given position to 1.
func (b *Bits) SetBit(pos int) {
	if pos >= b.size {
		return
	}
	elementIndex := pos / 64
	bitIndex := pos % 64

	b.data[elementIndex] |= 1 << bitIndex
}

// ClearBit clears the bit at the given position to 0.
func (b *Bits) ClearBit(pos int) {
	if pos >= b.size {
		return
	}
	elementIndex := pos / 64
	bitIndex := pos % 64
	b.data[elementIndex] &^= 1 << bitIndex
}

// IsBitSet checks if the bit at the given position is set to 1.
func (b *Bits) IsBitSet(pos int) bool {
	if pos >= b.size {
		return false
	}
	elementIndex := pos / 64
	bitIndex := pos % 64
	return (b.data[elementIndex] & (1 << bitIndex)) != 0
}

// CountOnes calculates the number of bits set to 1 in the bitset.
func (b *Bits) CountOnes() int {
	count := 0
	for i, value := range b.data {
		if i == len(b.data)-1 && b.size%64 != 0 {
			value &= (1 << (b.size % 64)) - 1
		}
		count += bits.OnesCount64(value)
	}
	return count
}

// LeftShift left shifts the entire bitset by n positions.
func (b *Bits) LeftShift(n int) {
	if n == 0 {
		return
	}

	// if n is greater than or equal to the size of the bitset, set all bits to 0.
	if n >= b.size {
		for i := 0; i < len(b.data); i++ {
			b.data[i] = 0
		}
		return
	}

	b.slowLeftShift(n)
}

func (b *Bits) slowLeftShift(n int) {
	// calculate the number of elements and bits to shift.
	shiftElements := n / 64
	// the number of bits to shift is the remainder of n divided by 64.
	shiftBits := n % 64

	for i := len(b.data) - 1; i >= shiftElements; i-- {
		b.data[i] = b.data[i-shiftElements] << shiftBits
		if i > shiftElements && shiftBits != 0 {
			b.data[i] |= b.data[i-shiftElements-1] >> (64 - shiftBits)
		}
	}
	for i := 0; i < shiftElements; i++ {
		b.data[i] = 0
	}
	if b.size%64 != 0 {
		b.data[len(b.data)-1] &= (1 << (b.size % 64)) - 1
	}
}

// RightShift right shifts the entire bitset by n positions.
func (b *Bits) RightShift(n int) {
	if n == 0 {
		return
	}

	// if n is greater than or equal to the size of the bitset, set all bits to 0.
	if n >= b.size {
		for i := 0; i < len(b.data); i++ {
			b.data[i] = 0
		}
		return
	}

	b.slowRightShift(n)
}

func (b *Bits) slowRightShift(n int) {
	// calculate the number of elements and bits to shift.
	shiftElements := n / 64
	// the number of bits to shift is the remainder of n divided by 64.
	shiftBits := n % 64
	for i := 0; i < len(b.data)-shiftElements; i++ {
		b.data[i] = b.data[i+shiftElements] >> shiftBits
		if i < len(b.data)-shiftElements-1 && shiftBits != 0 {
			b.data[i] |= b.data[i+shiftElements+1] << (64 - shiftBits)
		}
	}
	for i := len(b.data) - shiftElements; i < len(b.data); i++ {
		b.data[i] = 0
	}
}
