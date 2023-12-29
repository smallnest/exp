package atomicx_test

import (
	"testing"

	"github.com/smallnest/exp/sync/atomicx"
	"github.com/stretchr/testify/assert"
)

func TestPointer(t *testing.T) {
	type S struct {
		v int
	}

	n := 2
	var ap atomicx.Pointer[S]
	ap.Store(&S{v: n})

	running := make(chan bool, n)
	awake := make(chan bool, n)
	for i := 0; i < n; i++ {
		go func() {
			running <- true
			ap.Wait()
			awake <- true
		}()
	}
	for i := 0; i < n; i++ {
		<-running // Wait for everyone to run.
	}
	for n > 0 {
		select {
		case <-awake:
			t.Fatal("goroutine not asleep")
		default:
		}

		ap.Store(&S{v: n})
		ap.Signal()
		<-awake // Will deadlock if no goroutine wakes up
		select {
		case <-awake:
			t.Fatal("too many goroutines awake")
		default:
		}
		n--
	}

	ap.Store(&S{v: n})
	ap.Signal()
}

func TestPointerBroadcast(t *testing.T) {
	type S struct {
		v int
	}

	n := 200
	var av atomicx.Pointer[S]
	av.Store(&S{v: -1})

	running := make(chan int, n)
	awake := make(chan int, n)

	for i := 0; i < n; i++ {
		go func(g int) {
			running <- g
			av.Wait()
			awake <- g
		}(i)
	}

	for i := 0; i < n; i++ {
		<-running // Will deadlock unless n are running.
	}

	seen := make([]bool, n)
	for i := 0; i < n; i++ {
		select {
		case g := <-awake:
			if seen[g] {
				t.Fatal("goroutine woke up twice")
			}
			seen[g] = true
		default:
			av.Store(&S{v: i})
			av.Broadcast()
		}
	}

	assert.Equal(t, n, len(seen))
}

func TestValue(t *testing.T) {
	type S struct {
		v int
	}

	n := 2
	var av atomicx.Value
	av.Store(&S{v: -1})

	running := make(chan bool, n)
	awake := make(chan bool, n)
	for i := 0; i < n; i++ {
		go func() {
			running <- true
			av.Wait()
			awake <- true
		}()
	}
	for i := 0; i < n; i++ {
		<-running // Wait for everyone to run.
	}
	for n > 0 {
		select {
		case <-awake:
			t.Fatal("goroutine not asleep")
		default:
		}

		av.Store(&S{v: n})
		av.Signal()
		<-awake // Will deadlock if no goroutine wakes up
		select {
		case <-awake:
			t.Fatal("too many goroutines awake")
		default:
		}
		n--
	}

	av.Store(&S{v: n})
	av.Signal()
}

func TestValueBroadcast(t *testing.T) {
	type S struct {
		v int
	}

	n := 200
	var av atomicx.Value
	av.Store(&S{v: -1})

	running := make(chan int, n)
	awake := make(chan int, n)

	for i := 0; i < n; i++ {
		go func(g int) {
			running <- g
			av.Wait()
			awake <- g
		}(i)
	}

	for i := 0; i < n; i++ {
		<-running // Will deadlock unless n are running.
	}

	seen := make([]bool, n)
	for i := 0; i < n; i++ {
		select {
		case g := <-awake:
			if seen[g] {
				t.Fatal("goroutine woke up twice")
			}
			seen[g] = true
		default:
			av.Store(&S{v: i})
			av.Broadcast()
		}
	}

	assert.Equal(t, n, len(seen))
}

func TestBool(t *testing.T) {
	n := 2
	var av atomicx.Bool
	av.Store(true)

	running := make(chan bool, n)
	awake := make(chan bool, n)
	for i := 0; i < n; i++ {
		go func() {
			running <- true
			av.Wait()
			awake <- true
		}()
	}
	for i := 0; i < n; i++ {
		<-running // Wait for everyone to run.
	}
	for n > 0 {
		select {
		case <-awake:
			t.Fatal("goroutine not asleep")
		default:
		}

		av.Store(!av.Load())
		av.Signal()
		<-awake // Will deadlock if no goroutine wakes up

		select {
		case <-awake:
			t.Fatal("too many goroutines awake")
		default:
		}
		n--
	}

	av.Store(!av.Load())
	av.Signal()
}

func TestBoolBroadcast(t *testing.T) {
	n := 200
	var av atomicx.Bool
	av.Store(true)

	running := make(chan int, n)
	awake := make(chan int, n)

	for i := 0; i < n; i++ {
		go func(g int) {
			running <- g
			av.Wait()
			awake <- g
		}(i)
	}

	for i := 0; i < n; i++ {
		<-running // Will deadlock unless n are running.
	}

	seen := make([]bool, n)
	for i := 0; i < n; i++ {
		select {
		case g := <-awake:
			if seen[g] {
				t.Fatal("goroutine woke up twice")
			}
			seen[g] = true
		default:
			av.Store(false)
			av.Broadcast()
		}
	}

	assert.Equal(t, n, len(seen))
}

func TestInt32(t *testing.T) {
	n := 2
	var av atomicx.Int32
	av.Store(-1)

	running := make(chan bool, n)
	awake := make(chan bool, n)
	for i := 0; i < n; i++ {
		go func() {
			running <- true
			av.Wait()
			awake <- true
		}()
	}
	for i := 0; i < n; i++ {
		<-running // Wait for everyone to run.
	}
	for n > 0 {
		select {
		case <-awake:
			t.Fatal("goroutine not asleep")
		default:
		}

		av.Store(int32(n))
		av.Signal()
		<-awake // Will deadlock if no goroutine wakes up

		select {
		case <-awake:
			t.Fatal("too many goroutines awake")
		default:
		}
		n--
	}

	av.Store(int32(n))
	av.Signal()
}

func TestInt32Broadcast(t *testing.T) {
	n := 200
	var av atomicx.Int32
	av.Store(-1)

	running := make(chan int, n)
	awake := make(chan int, n)

	for i := 0; i < n; i++ {
		go func(g int) {
			running <- g
			av.Wait()
			awake <- g
		}(i)
	}

	for i := 0; i < n; i++ {
		<-running // Will deadlock unless n are running.
	}

	seen := make([]bool, n)
	for i := 0; i < n; i++ {
		select {
		case g := <-awake:
			if seen[g] {
				t.Fatal("goroutine woke up twice")
			}
			seen[g] = true
		default:
			av.Store(int32(i))
			av.Broadcast()
		}
	}

	assert.Equal(t, n, len(seen))
}

func TestInt64(t *testing.T) {
	n := 2
	var av atomicx.Int64
	av.Store(-1)

	running := make(chan bool, n)
	awake := make(chan bool, n)
	for i := 0; i < n; i++ {
		go func() {
			running <- true
			av.Wait()
			awake <- true
		}()
	}
	for i := 0; i < n; i++ {
		<-running // Wait for everyone to run.
	}
	for n > 0 {
		select {
		case <-awake:
			t.Fatal("goroutine not asleep")
		default:
		}

		av.Store(int64(n))
		av.Signal()
		<-awake // Will deadlock if no goroutine wakes up

		select {
		case <-awake:
			t.Fatal("too many goroutines awake")
		default:
		}
		n--
	}

	av.Store(int64(n))
	av.Signal()
}

func TestInt64Broadcast(t *testing.T) {
	n := 200
	var av atomicx.Int64
	av.Store(-1)

	running := make(chan int, n)
	awake := make(chan int, n)

	for i := 0; i < n; i++ {
		go func(g int) {
			running <- g
			av.Wait()
			awake <- g
		}(i)
	}

	for i := 0; i < n; i++ {
		<-running // Will deadlock unless n are running.
	}

	seen := make([]bool, n)
	for i := 0; i < n; i++ {
		select {
		case g := <-awake:
			if seen[g] {
				t.Fatal("goroutine woke up twice")
			}
			seen[g] = true
		default:
			av.Store(int64(i))
			av.Broadcast()
		}
	}

	assert.Equal(t, n, len(seen))
}

func TestUint32(t *testing.T) {
	n := 2
	var av atomicx.Uint32
	av.Store(uint32(n + 1))

	running := make(chan bool, n)
	awake := make(chan bool, n)
	for i := 0; i < n; i++ {
		go func() {
			running <- true
			av.Wait()
			awake <- true
		}()
	}
	for i := 0; i < n; i++ {
		<-running // Wait for everyone to run.
	}
	for n > 0 {
		select {
		case <-awake:
			t.Fatal("goroutine not asleep")
		default:
		}

		av.Store(uint32(n))
		av.Signal()
		<-awake // Will deadlock if no goroutine wakes up

		select {
		case <-awake:
			t.Fatal("too many goroutines awake")
		default:
		}
		n--
	}

	av.Store(uint32(n))
	av.Signal()
}

func TestUint32Broadcast(t *testing.T) {
	n := 200
	var av atomicx.Uint32
	av.Store(uint32(n + 1))

	running := make(chan int, n)
	awake := make(chan int, n)

	for i := 0; i < n; i++ {
		go func(g int) {
			running <- g
			av.Wait()
			awake <- g
		}(i)
	}

	for i := 0; i < n; i++ {
		<-running // Will deadlock unless n are running.
	}

	seen := make([]bool, n)
	for i := 0; i < n; i++ {
		select {
		case g := <-awake:
			if seen[g] {
				t.Fatal("goroutine woke up twice")
			}
			seen[g] = true
		default:
			av.Store(uint32(i))
			av.Broadcast()
		}
	}

	assert.Equal(t, n, len(seen))
}

func TestUint64(t *testing.T) {
	n := 2
	var av atomicx.Uint64
	av.Store(uint64(n + 1))

	running := make(chan bool, n)
	awake := make(chan bool, n)
	for i := 0; i < n; i++ {
		go func() {
			running <- true
			av.Wait()
			awake <- true
		}()
	}
	for i := 0; i < n; i++ {
		<-running // Wait for everyone to run.
	}
	for n > 0 {
		select {
		case <-awake:
			t.Fatal("goroutine not asleep")
		default:
		}

		av.Store(uint64(n))
		av.Signal()
		<-awake // Will deadlock if no goroutine wakes up

		select {
		case <-awake:
			t.Fatal("too many goroutines awake")
		default:
		}
		n--
	}

	av.Store(uint64(n))
	av.Signal()
}

func TestUint64Broadcast(t *testing.T) {
	n := 200
	var av atomicx.Uint64
	av.Store(uint64(n + 1))

	running := make(chan int, n)
	awake := make(chan int, n)

	for i := 0; i < n; i++ {
		go func(g int) {
			running <- g
			av.Wait()
			awake <- g
		}(i)
	}

	for i := 0; i < n; i++ {
		<-running // Will deadlock unless n are running.
	}

	seen := make([]bool, n)
	for i := 0; i < n; i++ {
		select {
		case g := <-awake:
			if seen[g] {
				t.Fatal("goroutine woke up twice")
			}
			seen[g] = true
		default:
			av.Store(uint64(i))
			av.Broadcast()
		}
	}

	assert.Equal(t, n, len(seen))
}

func TestUintptr(t *testing.T) {
	var v uintptr
	n := 2
	var av atomicx.Uintptr
	av.Store(v - 1)

	running := make(chan bool, n)
	awake := make(chan bool, n)
	for i := 0; i < n; i++ {
		go func() {
			running <- true
			av.Wait()
			awake <- true
		}()
	}
	for i := 0; i < n; i++ {
		<-running // Wait for everyone to run.
	}
	for n > 0 {
		select {
		case <-awake:
			t.Fatal("goroutine not asleep")
		default:
		}

		av.Store(v)
		v++
		av.Signal()
		<-awake // Will deadlock if no goroutine wakes up

		select {
		case <-awake:
			t.Fatal("too many goroutines awake")
		default:
		}
		n--
	}

	av.Store(v)
	av.Signal()
}

func TestUintptrBroadcast(t *testing.T) {
	var v uintptr

	n := 200
	var av atomicx.Uintptr
	av.Store(v - 1)

	running := make(chan int, n)
	awake := make(chan int, n)

	for i := 0; i < n; i++ {
		go func(g int) {
			running <- g
			av.Wait()
			awake <- g
		}(i)
	}

	for i := 0; i < n; i++ {
		<-running // Will deadlock unless n are running.
	}

	seen := make([]bool, n)
	for i := 0; i < n; i++ {
		select {
		case g := <-awake:
			if seen[g] {
				t.Fatal("goroutine woke up twice")
			}
			seen[g] = true
		default:
			v++
			av.Store(v)
			av.Broadcast()
		}
	}

	assert.Equal(t, n, len(seen))
}
