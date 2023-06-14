package win

import (
	"errors"
	"sync"
	"time"
)

// Sliding is a sliding window.
type Sliding[K Ordered, V any] struct {
	window      time.Duration // window duration
	granularity time.Duration // granularity of the window
	delay       time.Duration // delay of the window

	buckets *buckets[K, V] // buckets not yet in slided windows

	slidedBucketsMu sync.Mutex
	slidedBuckets   []*Bucket[K, V] // buckets in slided windows

	// SlidedChan is a channel to receive slided buckets.
	SlidedChan chan *Bucket[K, V] // insert slided buckets into channel to inform other watchers

	stopOnce sync.Once
	stopped  chan struct{}
}

// New creates a new sliding window.
func New[K Ordered, V any](window, granularity, delay time.Duration) (*Sliding[K, V], error) {
	s, err := newSliding[K, V](window, granularity, delay)
	if err != nil {
		return nil, err
	}

	// s.SlideChan is nil in default. set it as necessary.

	go s.shift()

	return s, nil
}

func newSliding[K Ordered, V any](window, granularity, delay time.Duration) (*Sliding[K, V], error) {
	if window == 0 {
		return nil, errors.New("sliding window cannot be zero")
	}
	if granularity == 0 {
		return nil, errors.New("granularity cannot be zero")
	}
	if window < granularity || window%granularity != 0 {
		return nil, errors.New("window size has to be a multiplier of granularity size")
	}

	s := &Sliding[K, V]{
		window:      window,
		granularity: granularity,
		delay:       delay,
		buckets:     NewBuckets[K, V](),
		stopped:     make(chan struct{}),
	}

	return s, nil
}

// NewChanSize creates a new sliding window with a channel size.
func NewChanSize[K Ordered, V any](window, granularity, delay time.Duration, chanSize int) (*Sliding[K, V], error) {
	s, err := newSliding[K, V](window, granularity, delay)
	if err != nil {
		return nil, err
	}
	s.SlidedChan = make(chan *Bucket[K, V], chanSize)

	go s.shift()

	return s, nil
}

// shift moves the window forward.
func (s *Sliding[K, V]) shift() {
	ticker := time.NewTicker(s.granularity)

	// delay
	dur := s.delay - s.granularity
	if dur > 0 {
		time.Sleep(dur)
	}

	for {
		select {
		case <-ticker.C:
			s.step()
		case <-s.stopped:
			return
		}
	}
}

func (s *Sliding[K, V]) step() {
	last := s.buckets.Last()
	if last != nil {
		s.slidedBucketsMu.Lock()
		s.slidedBuckets = append(s.slidedBuckets, last)
		s.slidedBucketsMu.Unlock()

		if s.SlidedChan != nil {
			select {
			case s.SlidedChan <- last:
			default:
				// chan is full
			}
		}
	}
}

// Add adds a value to the current window.
func (s *Sliding[K, V]) Add(key K, v V) {
	s.buckets.Add(key, v)
}

// Last returns the last bucket.
func (s *Sliding[K, V]) Last() (slided int, last *Bucket[K, V], err error) {
	s.slidedBucketsMu.Lock()
	defer s.slidedBucketsMu.Unlock()

	slided = len(s.slidedBuckets)
	if slided == 0 {
		return 0, nil, nil
	}

	last = s.slidedBuckets[0]
	s.slidedBuckets = s.slidedBuckets[1:]

	return slided, last, nil
}

// Stop stops the sliding window.
func (s *Sliding[_, _]) Stop() {
	s.stopOnce.Do(func() {
		close(s.stopped)
	})
}
