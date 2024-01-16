package sync_test

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"

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

func BenchmarkExchanger(b *testing.B) {
	exchanger := syncx.NewExchanger[[]int]()
	var buf1 = make([]int, 0, 1024)
	var buf2 = make([]int, 0, 1024)

	var reuslt = make(chan int, 1024)
	var done atomic.Bool

	b.ResetTimer()
	// producer
	go func() {
		buf := buf1
		for i := 0; ; i++ {
			if done.Load() {
				return
			}
			buf = append(buf, i)
			if i == 1023 {
				buf = exchanger.Exchange(buf)
				buf = buf[:0]
				i = 0
			}
		}
	}()

	// consumer
	go func() {
		buf := buf2
		for {
			if done.Load() {
				return
			}
			buf = exchanger.Exchange(buf)
			for _, n := range buf {
				pow(4) // mock process
				reuslt <- n
			}
		}
	}()

	for i := 0; i < b.N; i++ {
		<-reuslt
	}

	done.Store(true)
}

func BenchmarkExchanger_Channel(b *testing.B) {
	var buf = make(chan int, 1024)
	var reuslt = make(chan int, 1024)

	var done atomic.Bool

	b.ResetTimer()
	// producer
	go func() {
		for i := 0; ; i++ {
			if done.Load() {
				return
			}
			buf <- i
		}
	}()

	// consumer
	go func() {
		for i := 0; ; i++ {
			if done.Load() {
				return
			}
			n := <-buf
			pow(4) // mock process
			reuslt <- n
		}
	}()

	for i := 0; i < b.N; i++ {
		<-reuslt
	}

	done.Store(true)
}

func pow(targetBits int) {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	for {
		data := "hello world " + strconv.Itoa(nonce)
		hash = sha256.Sum256([]byte(data))
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			break
		} else {
			nonce++
		}

		if nonce%100 == 0 {
			runtime.Gosched()
		}
	}
}
