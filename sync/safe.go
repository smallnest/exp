package sync

import (
	"fmt"
	"io"
)

// SafeCloseChan closes the channel ch safely.
func SafeCloseChan[T any](ch chan T) (justClosed bool) {
	defer func() {
		if recover() != nil {
			justClosed = false
		}
	}()

	close(ch)

	return true
}

// SafeSendChan sends value to the channel ch safely even if ch has been closed.
func SafeSendChan[T any](ch chan T, value T) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()

	ch <- value

	return true
}

// SafeClose closes the closer safely.
// The closer can be a file, a network connection, etc.
func SafeClose(closer io.Closer) (justClosed bool, err error) {
	defer func() {
		if recover() != nil {
			justClosed = false
			fmt.Println("Panic occurred")
		}
	}()

	err = closer.Close()
	justClosed = true

	return justClosed, err
}
