package sync_test

import (
	"bytes"
	"fmt"
	"sync"
	"testing"
	"time"

	syncx "github.com/smallnest/exp/sync"
	"github.com/stretchr/testify/assert"
)

func ExampleExchanger() {
	buf1 := bytes.NewBuffer(make([]byte, 1024))
	buf2 := bytes.NewBuffer(make([]byte, 1024))

	exchanger := syncx.NewExchanger[*bytes.Buffer]()

	var wg sync.WaitGroup
	wg.Add(2)

	expect := 0
	go func() {
		defer wg.Done()

		buf := buf1
		for i := 0; i < 10; i++ {
			for j := 0; j < 1024; j++ {
				buf.WriteByte(byte(j / 256))
				expect += j / 256
			}

			buf = exchanger.Exchange(buf)
		}
	}()

	var got int
	go func() {
		defer wg.Done()

		buf := buf2
		for i := 0; i < 10; i++ {
			buf = exchanger.Exchange(buf)
			for _, b := range buf.Bytes() {
				got += int(b)
			}
			buf.Reset()
		}
	}()

	wg.Wait()

	fmt.Println(got)
	fmt.Println(expect == got)

	// Output:
	// 15360
	// true
}

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
	v, sent, exchanged := exchanger.ExchangeTimeout(1, 10*time.Millisecond)
	assert.True(t, sent)
	assert.False(t, exchanged)
	assert.Equal(t, 0, v)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		exchanger.Exchange(2) // success
		wg.Done()
	}()
	wg.Wait()

	// the first goroutine has sent but not received yet bucasue of timeout.
	// so we recv the value from the second goroutine.
	v = exchanger.Recv()
	assert.Equal(t, 2, v)
}
