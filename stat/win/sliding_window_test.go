package win

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSlidingChannel(t *testing.T) {
	w, err := NewChanSize[int64, int](time.Second, time.Second, 5*time.Second, 100)
	assert.NoError(t, err)
	w.Stop()

	w, err = New[int64, int](time.Second, time.Second, 5*time.Second)
	assert.NoError(t, err)
	w.SlidedChan = make(chan Result[int64, int], 100)

	defer w.Stop()

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
		timer := time.NewTimer(10 * time.Second)
		select {
		case b := <-w.SlidedChan:
			n++
			if n == 1 {
				assert.GreaterOrEqual(t, len(b.SlideOut.Values()), 0)
				assert.NotEmpty(t, b.CurrentWindow)
				continue
			}
			timer.Stop()
			assert.GreaterOrEqual(t, len(b.SlideOut.Values()), 8)
			assert.NotEmpty(t, b.CurrentWindow)

			assert.NotNil(t, w.Last())

			w.ForceForward()

			break loop
		case <-timer.C:
			t.Fatal("timeout")
			break loop
		}
	}

}

func TestSlidingChannel_New(t *testing.T) {
	w, err := NewChanSize[int64, int](-time.Second, time.Second, 5*time.Second, 100)
	assert.Error(t, err)
	assert.Nil(t, w)

	w, err = NewChanSize[int64, int](time.Second, -time.Second, 5*time.Second, 100)
	assert.Error(t, err)
	assert.Nil(t, w)

	w, err = NewChanSize[int64, int](time.Second, time.Second, -5*time.Second, 100)
	assert.Error(t, err)
	assert.Nil(t, w)

	w, err = NewChanSize[int64, int](time.Second, 5*time.Second, -5*time.Second, 100)
	assert.Error(t, err)
	assert.Nil(t, w)
}
