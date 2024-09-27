package stat

import (
	"iter"
	"math"
	"sort"
)

// Hist is a histogram.
type Hist struct {
	data  []int
	bins  uint
	log10 bool
}

// NewHist creates a new histogram.
// bins is the number of bins in the histogram.
func NewHist[T comparable](bins uint) *Hist {
	return &Hist{
		data:  make([]int, 0),
		bins:  bins,
		log10: false,
	}
}

// NewLogHist creates a new histogram with log10 bins.
// bins is the number of bins in the histogram.
func NewLogHist[T comparable](bins uint) *Hist {
	return &Hist{
		data:  make([]int, 0),
		bins:  bins,
		log10: true,
	}
}

// Add adds a value to the histogram.
func (h *Hist) Add(value int) {
	h.data = append(h.data, value)
}

// AddBatch adds multiple values to the histogram.
func (h *Hist) AddBatch(values ...int) {
	h.data = append(h.data, values...)
}

// Len returns the number of values in the histogram.
func (h *Hist) Len() int {
	return len(h.data)
}

// Histogram returns the histogram.
// The keys are the bin values and the values are the counts.
func (h *Hist) Histogram() map[int]int {
	if len(h.data) == 0 {
		return make(map[int]int)
	}

	min, max := h.data[0], h.data[0]
	for _, value := range h.data {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
	}

	if h.log10 {
		min = int(math.Log10(float64(min)))
		max = int(math.Log10(float64(max)))
	}

	// calculate the bin size
	binSize := (max-min)/int(h.bins) + 1

	histogram := make(map[int]int)

	// fill the histogram
	for _, value := range h.data {
		if h.log10 {
			value = int(math.Log10(float64(value)))
		}
		bin := (value - min) / binSize
		key := bin * binSize
		if h.log10 {
			key = int(math.Pow10(key))
		}

		histogram[key]++
	}

	return histogram
}

// Stat returns the min, max, avg and median of the histogram.
func (h *Hist) Stat() (min, max, avg, median int) {
	if len(h.data) == 0 {
		return 0, 0, 0, 0
	}

	min, max = h.data[0], h.data[0]
	avg = 0
	for _, value := range h.data {
		if value < min {
			min = value
		}
		if value > max {
			max = value
		}
		avg += value
	}
	avg /= len(h.data)

	// Calculate median
	sortedData := make([]int, len(h.data))
	copy(sortedData, h.data)
	sort.Ints(sortedData)
	if len(sortedData)%2 == 0 {
		median = (sortedData[len(sortedData)/2-1] + sortedData[len(sortedData)/2]) / 2
	} else {
		median = sortedData[len(sortedData)/2]
	}

	return min, max, avg, median
}

// All returns an iterator over the histogram.
// The iterator yields the bin value and the count.
func (h *Hist) All() iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		histogram := h.Histogram()
		keys := make([]int, 0, len(histogram))
		for k := range histogram {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		for _, bin := range keys {
			if !yield(bin, histogram[bin]) {
				return
			}
		}
	}
}
