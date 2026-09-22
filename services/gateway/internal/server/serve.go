package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"
)

// Serve serves web on lis until ctx is done, then stops taking requests
// and waits up to grace for those in flight to finish. It returns only
// once they have, so whatever they use can safely be closed afterwards.
func Serve(ctx context.Context, web *http.Server, lis net.Listener, grace time.Duration) error {
	stopped := make(chan error, 1)
	go func() { stopped <- web.Serve(lis) }() // name-ok
	select {
	case err := <-stopped:
		return err
	case <-ctx.Done():
	}
	shutdown, cancel := context.WithTimeout(context.Background(), grace)
	defer cancel()
	if err := web.Shutdown(shutdown); err != nil {
		return err
	}
	if err := <-stopped; !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
