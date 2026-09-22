// Package logging writes the service's logs as JSON lines, with every
// timestamp in UTC (roadmap R-128).
package logging

import (
	"io"
	"log/slog"
)

// New returns a logger writing to w.
func New(w io.Writer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey && len(groups) == 0 {
				a.Value = slog.TimeValue(a.Value.Time().UTC())
			}
			return a
		},
	}))
}
