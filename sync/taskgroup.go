package sync

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var defaultTaskGroupContext = context.Background()

// A TaskGroup is a collection of goroutines working on subtasks
// that are part of the same overall task.
//
// A zero TaskGroup is valid, has no limit on the number of active goroutines,
// and does not cancel on error.
type TaskGroup struct {
	ctx                context.Context
	cancel             func(error)
	cancelOnFirstError bool

	wg sync.WaitGroup

	sema chan struct{}

	errMutex sync.Mutex
	err      []error
}

// finish a task and update sema and wg.
func (g *TaskGroup) done() {
	if g.sema != nil {
		<-g.sema
	}
	g.wg.Done()
}

func (g *TaskGroup) errors() error {
	g.errMutex.Lock()
	defer g.errMutex.Unlock()

	return errors.Join(g.err...)
}

func (g *TaskGroup) appendErr(err error) {
	if err == nil {
		return
	}

	g.errMutex.Lock()
	defer g.errMutex.Unlock()

	if g.cancelOnFirstError {
		if len(g.err) == 0 {
			g.err = []error{err}
		}
		return
	}

	g.err = append(g.err, err)
}

// NewTaskGroup returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error in case of cancelOnErr==true or the first time Wait returns, whichever occurs
// first.
func NewTaskGroup(ctx context.Context) (*TaskGroup, context.Context) {
	ctx, cancel := context.WithCancelCause(ctx)
	return &TaskGroup{ctx: ctx, cancel: cancel}, ctx
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the all error (if any) from them.
func (g *TaskGroup) Wait() error {
	done := make(chan struct{})
	go func() {
		defer close(done)
		g.wg.Wait()
	}()

	ctx := g.ctx
	if ctx == nil {
		ctx = defaultTaskGroupContext
	}

	// context done or all tasks done
	select {
	case <-ctx.Done():
		g.appendErr(ctx.Err())
	case <-done:
	}

	err := g.errors()
	if g.cancel != nil {
		g.cancel(err)
	}
	return err
}

// WaitTimeout blocks until all function calls from the Go method have returned, then
// returns the all error (if any) from them.
// Or it will return context.DeadlineExceeded if timeout.
// Or it will return the first error if cancelOnFirstError is true.
func (g *TaskGroup) WaitTimeout(timeout time.Duration) error {
	done := make(chan struct{})
	go func() {
		defer close(done)
		g.wg.Wait()
	}()

	ctx := g.ctx
	if ctx == nil {
		ctx = defaultTaskGroupContext
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// context done or all tasks done
	select {
	case <-ctx.Done():
		g.appendErr(ctx.Err())
	case <-done:
	}

	err := g.errors()
	if g.cancel != nil {
		g.cancel(err)
	}
	return err
}

func (g *TaskGroup) run(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.done()

		ctx := g.ctx
		if ctx == nil {
			ctx = defaultTaskGroupContext
		}
		select {
		case <-ctx.Done():
			return
		default:
		}

		if err := f(); err != nil {
			g.appendErr(err)
			if g.cancelOnFirstError && g.cancel != nil {
				g.cancel(err)
			}
		}
	}()
}

// Go calls the given function in a new goroutine.
// It blocks until the new goroutine can be added without the number of
// active goroutines in the group exceeding the configured limit.
//
// The first call to return a non-nil error cancels the group's context, if the
// group was created by calling WithContext. The error will be returned by Wait.
func (g *TaskGroup) Go(f func() error) {
	if g.sema != nil {
		g.sema <- struct{}{}
	}

	g.run(f)
}

// TryGo calls the given function in a new goroutine only if the number of
// active goroutines in the group is currently below the configured limit.
//
// The return value reports whether the goroutine was started.
func (g *TaskGroup) TryGo(f func() error) bool {
	if g.sema != nil {
		select {
		case g.sema <- struct{}{}:
		default:
			return false
		}
	}

	g.run(f)

	return true
}

// SetLimit limits the number of active goroutines in this group to at most n.
// A negative value indicates no limit.
//
// Any subsequent call to the Go method will block until it can add an active
// goroutine without exceeding the configured limit.
//
// The limit must not be modified while any goroutines in the group are active.
func (g *TaskGroup) SetLimit(n int) {
	if n < 0 {
		g.sema = nil
		return
	}
	if len(g.sema) != 0 {
		panic(fmt.Errorf("taskgroup: modify limit while %v goroutines in the group are still active", len(g.sema)))
	}
	g.sema = make(chan struct{}, n)
}

// CancelOnFirstError configures the group to cancel its context (if any)
// as soon as any of goutines returns a non-nil error.
//
// Not like errgroup.Group, TaskGroup can decide whether to cancel this group or not immediately
// when any of goutines returns a non-nil error.
func (g *TaskGroup) CancelOnFirstError(flag bool) {
	g.cancelOnFirstError = flag
}
