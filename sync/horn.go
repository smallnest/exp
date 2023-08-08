package sync

import (
	"errors"
	"fmt"
	"sync"

	"github.com/hashicorp/go-multierror"
)

// ErrClosedChannel is returned when a send is attempted on a closed channel.
var ErrClosedChannel = errors.New("send after close")

type Horn[T any] struct {
	m              sync.RWMutex
	listeners      map[uint64]chan<- T
	nextListenerID uint64
	capacity       int
	closed         bool
}

func NewHorn[T any](n int) *Horn[T] {
	return &Horn[T]{capacity: n}
}

// SendNonblocking will send the value to all listeners. If a listener is not
// ready to receive the value, it will be skipped. If the horn is closed, an
// error will be returned.
func (h *Horn[T]) SendNonblocking(v T) error {
	h.m.Lock()
	defer h.m.Unlock()

	if h.closed {
		return ErrClosedChannel
	}

	var errs *multierror.Error

	for id, l := range h.listeners {
		select {
		case l <- v:
		default:
			err := fmt.Errorf("unable to send to listener '%d'", id)
			errs = multierror.Append(errs, err)
		}
	}

	return errs.ErrorOrNil()
}

// Send will send the value to all listeners. If a listener is not ready to
// receive the value, it will be blocked. If the horn is closed, an error will
// be returned.
func (h *Horn[T]) Send(v T) error {
	h.m.RLock()
	defer h.m.RUnlock()

	if h.closed {
		return ErrClosedChannel
	}

	for _, l := range h.listeners {
		select {
		case l <- v:
		}
	}

	return nil
}

// Close will close the horn and all listeners will be closed. Any subsequent
// calls to Send will return an error.
func (h *Horn[T]) Close() {
	h.m.Lock()
	defer h.m.Unlock()

	if h.closed {
		return
	}
	h.closed = true

	for _, l := range h.listeners {
		close(l)
	}
}

// AddListener will return a listener that can be used to receive values from the
// horn. The listener will be closed when the horn is closed.
func (h *Horn[T]) AddListener() *Listener[T] {
	h.m.Lock()
	defer h.m.Unlock()

	if h.listeners == nil {
		h.listeners = make(map[uint64]chan<- T)
	}
	if h.listeners[h.nextListenerID] != nil {
		h.nextListenerID++
	}

	// create a chanthe nel for new listener
	ch := make(chan T, h.capacity)
	if h.closed {
		close(ch)
	}
	h.listeners[h.nextListenerID] = ch

	return &Listener[T]{
		id: h.nextListenerID,
		ch: ch,
		h:  h,
	}
}

// RemoveListener will remove the listener from the horn.
// it is equivalent to calling l.Stop().
func (h *Horn[T]) RemoveListener(l *Listener[T]) {
	l.Stop()
}

// Listener is a handle to a horn listener.
type Listener[T any] struct {
	id uint64
	ch <-chan T
	h  *Horn[T]
}

// Stop will stop listening to the horn.
func (l *Listener[T]) Stop() {
	l.h.m.Lock()
	defer l.h.m.Unlock()
	delete(l.h.listeners, l.id)
}

// Chan returns the channel that can be used to receive values from the horn.
func (l *Listener[T]) Chan() <-chan T {
	return l.ch
}
