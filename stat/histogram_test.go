package stat

import (
	"reflect"
	"testing"
)

func TestHistAll(t *testing.T) {
	tests := []struct {
		name     string
		bins     uint
		values   []int
		expected map[int]int
	}{
		{
			name:     "simple case",
			bins:     3,
			values:   []int{1, 2, 3, 4, 5, 6},
			expected: map[int]int{0: 2, 2: 2, 4: 2},
		},
		{
			name:     "single bin",
			bins:     1,
			values:   []int{1, 2, 3, 4, 5, 6},
			expected: map[int]int{0: 6},
		},
		{
			name:     "empty histogram",
			bins:     3,
			values:   []int{},
			expected: map[int]int{},
		},
		{
			name:     "one value",
			bins:     3,
			values:   []int{5},
			expected: map[int]int{0: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hist := NewHist[int](tt.bins)
			hist.AddBatch(tt.values...)

			result := make(map[int]int)
			hist.All()(func(bin, count int) bool {
				result[bin] = count
				return true
			})

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestHistAll_Log(t *testing.T) {
	tests := []struct {
		name     string
		bins     uint
		values   []int
		expected map[int]int
	}{
		{
			name:     "simple case",
			bins:     3,
			values:   []int{1, 10, 100, 1000, 10000, 100000},
			expected: map[int]int{1: 2, 100: 2, 10000: 2},
		},
		{
			name:     "single bin",
			bins:     1,
			values:   []int{1, 10, 100, 1000, 10000, 100000},
			expected: map[int]int{1: 6},
		},
		{
			name:     "empty histogram",
			bins:     3,
			values:   []int{},
			expected: map[int]int{},
		},
		{
			name:     "one value",
			bins:     3,
			values:   []int{10000},
			expected: map[int]int{1: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hist := NewLogHist[int](tt.bins)
			hist.AddBatch(tt.values...)

			result := make(map[int]int)
			hist.All()(func(bin, count int) bool {
				result[bin] = count
				return true
			})

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestHistStat(t *testing.T) {
	tests := []struct {
		name     string
		bins     uint
		values   []int
		expected struct {
			min, max, avg, median int
		}
	}{
		{
			name:   "simple case",
			bins:   3,
			values: []int{1, 2, 3, 4, 5, 6},
			expected: struct {
				min, max, avg, median int
			}{min: 1, max: 6, avg: 3, median: 3},
		},
		{
			name:   "single value",
			bins:   3,
			values: []int{5},
			expected: struct {
				min, max, avg, median int
			}{min: 5, max: 5, avg: 5, median: 5},
		},
		{
			name:   "empty histogram",
			bins:   3,
			values: []int{},
			expected: struct {
				min, max, avg, median int
			}{min: 0, max: 0, avg: 0, median: 0},
		},
		{
			name:   "even number of values",
			bins:   3,
			values: []int{1, 2, 3, 4},
			expected: struct {
				min, max, avg, median int
			}{min: 1, max: 4, avg: 2, median: 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hist := NewHist[int](tt.bins)
			hist.AddBatch(tt.values...)

			min, max, avg, median := hist.Stat()

			if min != tt.expected.min || max != tt.expected.max || avg != tt.expected.avg || median != tt.expected.median {
				t.Errorf("expected (min: %v, max: %v, avg: %v, median: %v), got (min: %v, max: %v, avg: %v, median: %v)",
					tt.expected.min, tt.expected.max, tt.expected.avg, tt.expected.median, min, max, avg, median)
			}
		})
	}
}
