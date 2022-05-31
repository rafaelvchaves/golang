package ch03

import (
	"context"
	"io"
	"time"
)

const defaultPingInterval = 30 * time.Second

// Writes ping messages to a given writer at regular intervals.
func Pinger(ctx context.Context, w io.Writer, reset <-chan time.Duration) {
	var interval time.Duration
	select {
	case <-ctx.Done():
		return
	case interval = <-reset: // pull initial interval off reset channel
	default:
	}
	if interval <= 0 {
		interval = defaultPingInterval
	}

	timer := time.NewTimer(interval)
	defer func() {
		// avoid timer channel leak
		if !timer.Stop() {
			<-timer.C
		}
	}()
	for {
		select {
		case <-ctx.Done():
			// context was cancelled, exit
			return
		case newInterval := <-reset:
			// signal to reset timer was received. Stop timer and set interval.
			if !timer.Stop() {
				<-timer.C
			}
			if newInterval > 0 {
				interval = newInterval
			}
		case <-timer.C:
			// timer expired, send ping.
			if _, err := w.Write([]byte("ping")); err != nil {
				// track and act on consecutive timeouts here.
				return
			}
		}

		_ = timer.Reset(interval)
	}
}