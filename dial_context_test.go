package ch03

import (
	"context"
	"net"
	"syscall"
	"testing"
	"time"
)

// A context is an object that you can use the send cancellation signals to
// your asynchronous processes. It also allows you to send a cancellation signal
// after it reaches a deadline or after its timer expires.

func TestDialContext(t *testing.T) {
	dl := time.Now().Add(5 * time.Second)

	// Create a context with a deadline and its cancel function 
	ctx, cancel := context.WithDeadline(context.Background(), dl)
	defer cancel()

	var d net.Dialer // DialContext is a method on a Dialer

	// Called at the beginning of DialContext?
	d.Control = func(_, _ string, _ syscall.RawConn) error {
		// Sleep long enough to reach the context's deadline
		time.Sleep(5 * time.Second + time.Millisecond)
		return nil
	}
	conn, err := d.DialContext(ctx, "tcp", "10.0.0.0:80")
	if err == nil {
		conn.Close()
		t.Fatal("connection did not time out")
	}
	nErr, ok := err.(net.Error)
	if !ok {
		t.Error(err)
	} else {
		if !nErr.Timeout() {
			t.Errorf("error is not a timeout: %v", err)
		}
	}
	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("expected deadline exceeded; actual: %v", ctx.Err())
	}
}

func TestDialContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sync := make(chan struct{})
	go func() {
		defer func() { sync <- struct{}{} }()

		var d net.Dialer
		d.Control = func(_, _ string, _ syscall.RawConn) error {
			time.Sleep(time.Second)
			return nil
		}
		conn, err := d.DialContext(ctx, "tcp", "10.0.0.1:80")
		if err != nil {
			t.Log(err)
			return
		}

		conn.Close()
		t.Error("connection did not time out")
	}()

	// cancel while dialer is attempting to connect to and handshake with server
	// node
	cancel()
	<-sync

	if ctx.Err() != context.Canceled {
		t.Errorf("expected canceled context; actual: %q", ctx.Err())
	}
}