package sync_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/smallnest/exp/sync"
)

var (
	Web   = fakeSearch("web")
	Image = fakeSearch("image")
	Video = fakeSearch("video")
)

type Result string
type Search func(ctx context.Context, query string) (Result, error)

func fakeSearch(kind string) Search {
	return func(_ context.Context, query string) (Result, error) {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q", kind, query)), nil
	}
}

func compareErr(e1, e2 error) bool {
	if e1 == nil && e2 == nil {
		return true
	}

	if e1 == nil || e2 == nil {
		return false
	}
	return e1.Error() == e2.Error()
}

func TestZeroGroup(t *testing.T) {
	err1 := errors.New("errgroup_test: 1")
	err2 := errors.New("errgroup_test: 2")

	cases := []struct {
		errs []error
	}{
		{errs: []error{}},
		{errs: []error{nil}},
		{errs: []error{err1}},
		{errs: []error{err1, nil}},
		{errs: []error{err1, nil, err2}},
	}

	for _, tc := range cases {
		g := new(sync.TaskGroup)

		var errs []error

		for i, err := range tc.errs {
			err := err
			g.Go(func() error { return err })

			if err != nil {
				errs = append(errs, err)
			}

			if gErr := g.Wait(); !compareErr(errors.Join(errs...), gErr) {
				t.Errorf("after %T.Go(func() error { return err }) for err in #%d\n"+
					"g.Wait() = %v; want %v",
					g, i, gErr, errors.Join(errs...))
			}
		}
	}
}

func TestNewTaskGroup(t *testing.T) {
	errDoom := errors.New("group_test: doomed")

	cases := []struct {
		errs []error
		want error
	}{
		{want: nil},
		{errs: []error{nil}, want: nil},
		{errs: []error{errDoom}, want: errDoom},
		{errs: []error{errDoom, nil}, want: errDoom},
	}

	for _, tc := range cases {
		g, ctx := sync.NewTaskGroup(context.Background())

		for _, err := range tc.errs {
			err := err
			g.Go(func() error { return err })
		}

		if err := g.Wait(); !compareErr(err, tc.want) {
			t.Errorf("after %T.Go(func() error { return err }) for err in %v\n"+
				"g.Wait() = %v; want %v",
				g, tc.errs, err, tc.want)
		}

		canceled := false
		select {
		case <-ctx.Done():
			canceled = true
		default:
		}
		if !canceled {
			t.Errorf("after %T.Go(func() error { return err }) for err in %v\n"+
				"ctx.Done() was not closed",
				g, tc.errs)
		}
	}
}

func TestTryGo(t *testing.T) {
	g := &sync.TaskGroup{}
	n := 42
	g.SetLimit(42)
	ch := make(chan struct{})
	fn := func() error {
		ch <- struct{}{}
		return nil
	}
	for i := 0; i < n; i++ {
		if !g.TryGo(fn) {
			t.Fatalf("TryGo should succeed but got fail at %d-th call.", i)
		}
	}
	if g.TryGo(fn) {
		t.Fatalf("TryGo is expected to fail but succeeded.")
	}
	go func() {
		for i := 0; i < n; i++ {
			<-ch
		}
	}()
	g.Wait()

	if !g.TryGo(fn) {
		t.Fatalf("TryGo should success but got fail after all goroutines.")
	}
	go func() { <-ch }()
	g.Wait()

	// Switch limit.
	g.SetLimit(1)
	if !g.TryGo(fn) {
		t.Fatalf("TryGo should success but got failed.")
	}
	if g.TryGo(fn) {
		t.Fatalf("TryGo should fail but succeeded.")
	}
	go func() { <-ch }()
	g.Wait()

	// Block all calls.
	g.SetLimit(0)
	for i := 0; i < 1<<10; i++ {
		if g.TryGo(fn) {
			t.Fatalf("TryGo should fail but got succeded.")
		}
	}
	g.Wait()
}

func TestGoLimit(t *testing.T) {
	const limit = 10

	g := &sync.TaskGroup{}
	g.SetLimit(limit)
	var active int32
	for i := 0; i <= 1<<10; i++ {
		g.Go(func() error {
			n := atomic.AddInt32(&active, 1)
			if n > limit {
				return fmt.Errorf("saw %d active goroutines; want ≤ %d", n, limit)
			}
			time.Sleep(1 * time.Microsecond) // Give other goroutines a chance to increment active.
			atomic.AddInt32(&active, -1)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkGo(b *testing.B) {
	fn := func() {}
	g := &sync.TaskGroup{}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		g.Go(func() error { fn(); return nil })
	}
	g.Wait()
}

func TestCancelOnFirstError(t *testing.T) {
	err1 := errors.New("errgroup_test: 1")
	err2 := errors.New("errgroup_test: 2")

	cases := []struct {
		errs []error
	}{
		{errs: []error{}},
		{errs: []error{nil}},
		{errs: []error{err1}},
		{errs: []error{err1, nil}},
		{errs: []error{err1, nil, err2}},
	}

	for _, tc := range cases {
		g := new(sync.TaskGroup)
		g.CancelOnFirstError(true)

		var firstErr error
		for i, err := range tc.errs {
			err := err
			g.Go(func() error { return err })

			if firstErr == nil && err != nil {
				firstErr = err
			}

			if gErr := g.Wait(); !compareErr(gErr, firstErr) {
				t.Errorf("after %T.Go(func() error { return err }) for err in %v\n"+
					"g.Wait() = %v; want %v",
					g, tc.errs[:i+1], err, firstErr)
			}
		}
	}

}

func TestWithContext(t *testing.T) {
	errDoom := errors.New("group_test: doomed")

	cases := []struct {
		errs []error
		want error
	}{
		{want: nil},
		{errs: []error{nil}, want: nil},
		{errs: []error{errDoom}, want: errDoom},
		{errs: []error{errDoom, nil}, want: errDoom},
	}

	parent, cancel := context.WithCancel(context.Background())
	cancel() // nolint

	for _, tc := range cases {
		g, _ := sync.NewTaskGroup(parent)

		for _, err := range tc.errs {
			err := err
			g.Go(func() error {
				time.Sleep(time.Second)
				return err
			})
		}

		if err := g.Wait(); !errors.Is(err, context.Canceled) {
			t.Errorf("after %T.Go(func() error { return err }) for err in %v\n"+
				"g.Wait() = %v; want %v",
				g, tc.errs, err, context.Canceled)
		}
	}
}

func TestWaitTimeout(t *testing.T) {
	errDoom := errors.New("group_test: doomed")

	cases := []struct {
		errs []error
		want error
	}{
		{errs: []error{nil}, want: nil},
		{errs: []error{errDoom}, want: errDoom},
		{errs: []error{errDoom, nil}, want: errDoom},
	}

	for _, tc := range cases {
		g, _ := sync.NewTaskGroup(context.Background())

		for _, err := range tc.errs {
			err := err
			g.Go(func() error {
				time.Sleep(time.Second)
				return err
			})
		}

		if err := g.WaitTimeout(time.Millisecond); !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("after %T.Go(func() error { return err }) for err in %v\n"+
				"g.Wait() = %v; want %v",
				g, tc.errs, err, context.Canceled)
		}
	}
}
