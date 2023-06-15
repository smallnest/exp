package win_test

import (
	"log"
	"math/rand"
	"time"

	"github.com/smallnest/exp/stat/win"
)

type Metric struct {
	TimeStamp int64
	DestIP    string
	Success   int
	Fail      int
}

func ExampleSliding() {
	w, err := win.NewChanSize[int64, Metric](time.Second, time.Second, 5*time.Second, 100)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	// use the sliding window in your code for stat.
	go func() {
		for {
			select {
			case <-done:
				return
			default:
			}

			ts := time.Now().UnixNano()
			key := ts / int64(time.Second)
			w.Add(key, Metric{
				TimeStamp: ts,
				DestIP:    "192.168.1.1",
				Success:   rand.Intn(10000),
				Fail:      rand.Intn(10),
			})

			time.Sleep(time.Millisecond)
		}
	}()

	for b := range w.SlidedChan {
		if b.SlideOut == nil {
			return
		}

		key := b.SlideOut.Key

		var total, fail int
		for _, v := range b.SlideOut.Values() {
			total += v.Success
			fail += v.Fail
		}

		log.Printf("key: %s, total: %d, fail: %d, %d buckets in current window", time.Unix(key, 0).Format(time.DateTime), total, fail, len(b.CurrentWindow))
	}

}
