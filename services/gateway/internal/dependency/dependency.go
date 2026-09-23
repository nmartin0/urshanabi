// Package dependency bounds what one failing service can cost this one
// (roadmap R-83). A synchronous chain of calls is a known cause of
// cascading failure, and the established remedies are all here: every
// call carries a deadline, a dependency that keeps failing is left alone
// until it recovers, and each dependency has its own slots, so one
// stalled service cannot consume the resources another's callers need.
package dependency

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Unavailable is returned without calling: either the dependency is
// being left alone to recover, or its slots are all in use. A caller
// learns only that it is unavailable (R-22).
var Unavailable = errors.New("the dependency is unavailable")

// Guard bounds the calls to one dependency.
type Guard struct {
	deadline time.Duration
	failures int
	cooldown time.Duration
	slots    chan struct{}
	now      func() time.Time

	mu       sync.Mutex
	failed   int
	restsTil time.Time
}

// New returns a guard allowing slots calls at once, leaving the
// dependency alone for cooldown once failures calls in a row have
// failed, and bounding every call by deadline.
func New(slots, failures int, cooldown, deadline time.Duration) *Guard {
	return &Guard{
		deadline: deadline,
		failures: failures,
		cooldown: cooldown,
		slots:    make(chan struct{}, slots),
		now:      time.Now,
	}
}

// Do runs call against the dependency, with the deadline, unless the
// dependency is resting or its slots are full.
func (g *Guard) Do(ctx context.Context, call func(context.Context) error) error {
	if g.resting() {
		return Unavailable
	}
	select {
	case g.slots <- struct{}{}:
		defer func() { <-g.slots }()
	default:
		// Every slot is in use: failing now keeps this service
		// answering, where queueing would sink it with the dependency.
		return Unavailable
	}
	// The deadline travels with the call and shrinks at each hop: a
	// caller's own deadline, if sooner, is what binds.
	ctx, cancel := context.WithTimeout(ctx, g.deadline)
	defer cancel()
	err := call(ctx)
	g.record(err == nil)
	return err
}

// resting says whether the dependency is being left alone to recover.
func (g *Guard) resting() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.now().Before(g.restsTil)
}

// record counts a call towards resting the dependency, or ends the rest.
func (g *Guard) record(worked bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if worked {
		g.failed = 0
		g.restsTil = time.Time{}
		return
	}
	g.failed++
	if g.failed >= g.failures {
		g.restsTil = g.now().Add(g.cooldown)
		g.failed = 0
	}
}
