package win

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSliding(t *testing.T) {
	w, err := New[int64, int](time.Second, time.Second, 5*time.Second)
	assert.NoError(t, err)

	go func() {
		for i := 0; ; i++ {
			key := time.Now().Unix()
			w.Add(key, i)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	time.Sleep(6 * time.Second)
	_, _, _, err = w.Last()
	assert.NoError(t, err)
	slided, bucket, window, err := w.Last()
	assert.NoError(t, err)

	require.NotEmpty(t, bucket)

	assert.GreaterOrEqual(t, len(bucket.Values()), 9)
	assert.GreaterOrEqual(t, slided, 1)
	assert.GreaterOrEqual(t, len(window), 0)
}

func TestSlidingChannel(t *testing.T) {
	w, err := NewChanSize[int64, int](time.Second, time.Second, 5*time.Second, 100)
	assert.NoError(t, err)

	go func() {
		for i := 0; ; i++ {
			key := time.Now().Unix()
			w.Add(key, i)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	n := 0
loop:
	for {
		timer := time.NewTimer(6 * time.Second)
		select {
		case b := <-w.SlidedChan:
			n++
			if n == 1 {
				assert.GreaterOrEqual(t, len(b.SlideOut.Values()), 0)
				assert.NotEmpty(t, b.CurrentWindow)
				continue
			}
			timer.Stop()
			assert.GreaterOrEqual(t, len(b.SlideOut.Values()), 9)
			assert.NotEmpty(t, b.CurrentWindow)
			break loop
		case <-timer.C:
			t.Fatal("timeout")
			break loop
		}
	}

}
