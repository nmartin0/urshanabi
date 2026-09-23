// Package startup checks at startup what configuration alone cannot
// tell: whether what this service depends on is actually there
// (roadmap R-61). Elysium validated its configuration when it loaded
// it, and learned only on a first failing query that a source was
// unreachable or a credential wrong. Each dependency is named in the
// report, so an operator reads which one, not that something failed.
package startup

import (
	"context"
	"log/slog"
	"time"
)

// Dependency is one thing a service needs, and how to ask whether it is
// there. The question must be cheap and read nothing.
type Dependency struct {
	Name  string
	Where string
	Ask   func(context.Context) error
}

// Report asks each dependency whether it is there, within bound, and
// says so by name. It returns how many could not be reached: a service
// may still start, so that what does not need the missing dependency
// keeps working (R-83), but an operator is told at startup rather than
// by the first request that fails.
func Report(ctx context.Context, log *slog.Logger, bound time.Duration, dependencies ...Dependency) int {
	missing := 0
	for _, dependency := range dependencies {
		asked, cancel := context.WithTimeout(ctx, bound)
		err := dependency.Ask(asked)
		cancel()
		if err != nil {
			missing++
			log.Error("a dependency is not reachable at startup",
				"dependency", dependency.Name, "address", dependency.Where)
			continue
		}
		log.Info("a dependency is reachable", "dependency", dependency.Name, "address", dependency.Where)
	}
	return missing
}
