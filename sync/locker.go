package sync

import (
	"runtime"
	"strconv"
	"sync"
)

var lockers Map[string, *sync.Mutex]

func Lock(fn func()) bool {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return false
	}

	key := file + ":" + strconv.Itoa(line)

	mu, _ := lockers.LoadOrStore(key, &sync.Mutex{})
	defer mu.Unlock()

	fn()

	return true
}
