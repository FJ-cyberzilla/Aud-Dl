package strategy

import (
	"context"
	"sync"
	"time"
)

type GracefulLimiter struct {
	sem        chan struct{} // Concurrency slot worker limiter
	rateTicker *time.Ticker
	mu         sync.Mutex
}

func NewGracefulLimiter(maxConcurrent int, requestsPerSec int) *GracefulLimiter {
	limiter := &GracefulLimiter{
		sem:        make(chan struct{}, maxConcurrent),
		rateTicker: time.NewTicker(time.Second / time.Duration(requestsPerSec)),
	}
	return limiter
}

// AcquireSlot blocks until both concurrency limit and rate-limit ticks are satisfied
func (gl *GracefulLimiter) AcquireSlot(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case gl.sem <- struct{}{}: // Acquire concurrency worker token
	}

	select {
	case <-ctx.Done():
		<-gl.sem // Release token if context canceled
		return ctx.Err()
	case <-gl.rateTicker.C: // Pace request frequency
		return nil
	}
}

// ReleaseSlot frees a worker concurrency slot after download finishes
func (gl *GracefulLimiter) ReleaseSlot() {
	<-gl.sem
}
