package sync_test

import (
	"sync"
	"testing"
	"time"

	syncx "github.com/smallnest/exp/sync"
	"github.com/stretchr/testify/assert"
)

func TestExchanger(t *testing.T) {
	exchanger := syncx.NewExchanger[int]()

	var wg sync.WaitGroup
	wg.Add(2)

	var leftReceived []int
	var rightReceived []int

	go func() {
		defer wg.Done()

		for i := 0; i < 5; i++ {
			n := exchanger.Exchange(i*2 + 1) // send odd
			leftReceived = append(leftReceived, n)
		}
	}()

	go func() {
		defer wg.Done()

		for i := 0; i < 5; i++ {
			n := exchanger.Exchange(i*2 + 2) // send even
			rightReceived = append(rightReceived, n)
		}
	}()

	wg.Wait()

	assert.Equal(t, []int{2, 4, 6, 8, 10}, leftReceived)
	assert.Equal(t, []int{1, 3, 5, 7, 9}, rightReceived)
}

func TestExchanger_panic(t *testing.T) {
	exchanger := syncx.NewExchanger[int]()

	var wg sync.WaitGroup
	wg.Add(2)

	var leftReceived []int
	var rightReceived []int

	go func() {
		defer wg.Done()

		for i := 0; i < 5; i++ {
			n := exchanger.Exchange(i*2 + 1) // send odd
			leftReceived = append(leftReceived, n)
			if i == 3 {
				wg.Done()
			}
		}
	}()

	go func() {
		defer wg.Done()

		for i := 0; i < 4; i++ {
			n := exchanger.Exchange(i*2 + 2) // send even
			rightReceived = append(rightReceived, n)
		}
	}()

	wg.Wait()

	assert.Equal(t, []int{2, 4, 6, 8}, leftReceived)
	assert.Equal(t, []int{1, 3, 5, 7}, rightReceived)

	assert.Panics(t, func() { exchanger.Exchange(10) })
}

func TestExchanger_timeout(t *testing.T) {
	exchanger := syncx.NewExchanger[int]()
	v, exchanged := exchanger.ExchangeTimeout(1, 10*time.Millisecond)
	assert.False(t, exchanged)
	assert.Equal(t, 0, v)
}
