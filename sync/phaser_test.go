package sync

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPhaser(t *testing.T) {
	// 示例用法
	phaser := NewPhaser(0)

	var wg sync.WaitGroup
	wg.Add(3)

	for i := 0; i < 3; i++ {
		id := i
		phaser.Join()
		go func() {
			defer wg.Done()

			t.Logf("goroutine %d started\n", id)

			// phase == 0
			phase := phaser.ArriveAndWait()
			assert.Equal(t, int32(1), phase)
			t.Logf("goroutine %d finished phase 0\n", id)

			time.Sleep(time.Duration(rand.Intn(100) * int(time.Millisecond)))

			// phase == 1
			assert.Equal(t, int32(1), phaser.Phase())
			if id == 1 {
				phase = phaser.ArriveAndLeave()
				t.Logf("goroutine %d finished phase %d and deregistered\n", phase, id)
				return
			} else {
				phase = phaser.ArriveAndWait()
			}

			assert.Equal(t, int32(2), phase)
			t.Logf("goroutine %d finished phase 1\n", id)

			time.Sleep(time.Duration(rand.Intn(100) * int(time.Millisecond)))

			phaser.ArriveAndWait()
			assert.Equal(t, int32(2), phase)
			t.Logf("goroutine %d finished phase 2\n", id)
		}()
	}

	wg.Wait()

	// assert.Equal(t, 2, phaser.Phase())

	t.Logf("phaser terminated")
}

func TestPhaser_phase(t *testing.T) {
	// 3 phases and one exit after phase==0

	phaser := NewPhaser(0)

	var wg sync.WaitGroup
	wg.Add(3)

	for i := 0; i < 3; i++ {
		phaser.Join()
		go func(id int) {
			defer wg.Done()

			t.Logf("goroutine %d started\n", id)

			// phase == 0
			phase := phaser.ArriveAndWait()
			assert.Equal(t, int32(1), phase)
			t.Logf("goroutine %d finished phase 0\n", id)

			if id == 1 {
				phaser.Leave()
				t.Logf("goroutine %d exit after phase 0\n", id)
				return
			}

			time.Sleep(time.Duration(rand.Intn(100) * int(time.Millisecond)))

			// phase == 1
			assert.Equal(t, int32(1), phaser.Phase())
			phase = phaser.ArriveAndWait()
			assert.Equal(t, int32(2), phase)
			t.Logf("goroutine %d finished phase 1\n", id)

			time.Sleep(time.Duration(rand.Intn(100) * int(time.Millisecond)))

			phaser.ArriveAndWait()
			assert.Equal(t, int32(2), phase)
			t.Logf("goroutine %d finished phase 2 and deregistered\n", id)
		}(i)
	}

	wg.Wait()

	// assert.Equal(t, 2, phaser.Phase())

	t.Logf("phaser terminated")
}
