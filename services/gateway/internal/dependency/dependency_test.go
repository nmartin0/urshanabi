package dependency

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

// at returns a guard whose clock the test drives.
func at(slots, failures int, cooldown, deadline time.Duration, clock *time.Time) *Guard {
	g := New(slots, failures, cooldown, deadline)
	g.now = func() time.Time { return *clock }
	return g
}

func TestACallCarriesTheDeadline(t *testing.T) {
	clock := time.Now()
	g := at(1, 3, time.Minute, 40*time.Millisecond, &clock)
	err := g.Do(context.Background(), func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("the call ended with %v, not at its deadline", err)
	}
}

func TestACallersSoonerDeadlineIsWhatBinds(t *testing.T) {
	clock := time.Now()
	g := at(1, 3, time.Minute, time.Hour, &clock)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	start := time.Now()
	_ = g.Do(ctx, func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})
	if time.Since(start) > time.Second {
		t.Fatal("the guard's own deadline outlived the caller's")
	}
}

func TestAFailingDependencyIsLeftAlone(t *testing.T) {
	clock := time.Now()
	g := at(4, 2, time.Minute, time.Second, &clock)
	broken := errors.New("no")
	calls := 0
	failing := func(context.Context) error { calls++; return broken }

	for i := 0; i < 2; i++ {
		if err := g.Do(context.Background(), failing); !errors.Is(err, broken) {
			t.Fatalf("call %d: %v", i, err)
		}
	}
	if err := g.Do(context.Background(), failing); !errors.Is(err, Unavailable) {
		t.Fatalf("after two failures the dependency was still called: %v", err)
	}
	if calls != 2 {
		t.Fatalf("the dependency was called %d times, want 2", calls)
	}
	// Once the rest is over, it is tried again, and a success ends it.
	clock = clock.Add(2 * time.Minute)
	if err := g.Do(context.Background(), func(context.Context) error { calls++; return nil }); err != nil {
		t.Fatalf("after the rest: %v", err)
	}
	if calls != 3 {
		t.Fatalf("the dependency was called %d times, want 3", calls)
	}
}

func TestSlotsAreBoundedSoAStallCannotSpread(t *testing.T) {
	clock := time.Now()
	g := at(2, 100, time.Minute, time.Second, &clock)
	stalled := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() { // name-ok
			defer wg.Done()
			_ = g.Do(context.Background(), func(context.Context) error { <-stalled; return nil })
		}()
	}
	// Wait for both slots to be taken.
	for i := 0; i < 100 && len(g.slots) < 2; i++ {
		time.Sleep(time.Millisecond)
	}
	if err := g.Do(context.Background(), func(context.Context) error { return nil }); !errors.Is(err, Unavailable) {
		t.Fatalf("a third call was admitted with both slots in use: %v", err)
	}
	close(stalled)
	wg.Wait()
	if err := g.Do(context.Background(), func(context.Context) error { return nil }); err != nil {
		t.Fatalf("once the slots freed: %v", err)
	}
}
