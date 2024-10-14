package gods

import (
	"sync"
	"testing"
)

func BenchmarkPoolDequeue(b *testing.B) {
	const size = 1024
	pd := NewPoolDequeue(size)
	var wg sync.WaitGroup

	// Producer
	go func() {
		for i := 0; i < b.N; i++ {
			pd.PushHead(i)
		}
		wg.Done()
	}()

	// Consumers
	numConsumers := 10
	wg.Add(numConsumers + 1)
	for i := 0; i < numConsumers; i++ {
		go func() {
			for {
				if _, ok := pd.PopTail(); !ok {
					break
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func BenchmarkPoolChain(b *testing.B) {
	pc := NewPoolChain()
	var wg sync.WaitGroup

	// Producer
	go func() {
		for i := 0; i < b.N; i++ {
			pc.PushHead(i)
		}
		wg.Done()
	}()

	// Consumers
	numConsumers := 10
	wg.Add(numConsumers + 1)
	for i := 0; i < numConsumers; i++ {
		go func() {
			for {
				if _, ok := pc.PopTail(); !ok {
					break
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

func BenchmarkChannel(b *testing.B) {
	ch := make(chan interface{}, 1024)
	var wg sync.WaitGroup

	// Producer
	go func() {
		for i := 0; i < b.N; i++ {
			ch <- i
		}
		close(ch)
		wg.Done()
	}()

	// Consumers
	numConsumers := 10
	wg.Add(numConsumers + 1)
	for i := 0; i < numConsumers; i++ {
		go func() {
			for range ch {
			}
			wg.Done()
		}()
	}

	wg.Wait()
}
