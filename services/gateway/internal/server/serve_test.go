package server

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestShutdownLetsInFlightRequestsFinish(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	started := make(chan struct{})
	web := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		close(started)
		time.Sleep(300 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})}
	ctx, stop := context.WithCancel(context.Background())
	served := make(chan error, 1)
	go func() { served <- Serve(ctx, web, lis, 5*time.Second) }() // name-ok

	answered := make(chan int, 1)
	go func() { // name-ok
		resp, err := http.Get("http://" + lis.Addr().String())
		if err != nil {
			answered <- 0
			return
		}
		_ = resp.Body.Close()
		answered <- resp.StatusCode
	}()
	<-started
	stop()
	if code := <-answered; code != http.StatusOK {
		t.Fatalf("the in-flight request got %d, not 200", code)
	}
	if err := <-served; err != nil {
		t.Fatalf("Serve returned %v", err)
	}
}
