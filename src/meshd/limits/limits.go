// Package limits enforces resource caps for the PiN node.
// It is cross-platform and uses no OS-specific system calls.
//
// Three caps are enforced:
//   - Memory: soft cap via Go runtime, triggers aggressive GC near the limit
//   - Concurrency: max simultaneous requests, indirectly caps CPU usage
//   - Bandwidth: token bucket rate limiter on bytes served per second
package limits

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"meshd/config"
)

// Limiter enforces resource caps for the node.
type Limiter struct {
	// concurrency controls max simultaneous content requests
	concurrency chan struct{}

	// bandwidth is a token bucket for bytes per second
	bandwidth *tokenBucket

	cfg *config.Config
}

// New creates a Limiter from the given config and applies the memory cap.
func New(cfg *config.Config) *Limiter {
	l := &Limiter{cfg: cfg}

	// Apply memory limit via Go runtime
	// This tells the GC to be aggressive when approaching the limit
	// rather than letting the process grow unbounded.
	if cfg.Limits.RAMMB > 0 {
		memLimit := int64(cfg.Limits.RAMMB) * 1024 * 1024
		debug.SetMemoryLimit(memLimit)
	}

	// Set up concurrency limiter
	// We derive max concurrent requests from CPU percent:
	// 100% = 16 concurrent, 25% = 4 concurrent, 10% = 2 concurrent
	maxConcurrent := (cfg.Limits.CPUPercent * 16) / 100
	if maxConcurrent < 1 {
		maxConcurrent = 1
	}
	l.concurrency = make(chan struct{}, maxConcurrent)

	// Set up bandwidth token bucket
	// Convert Mbps to bytes per second
	bytesPerSec := int64(cfg.Limits.BandwidthMbps) * 125000 // 1 Mbps = 125000 bytes/sec
	if bytesPerSec > 0 {
		l.bandwidth = newTokenBucket(bytesPerSec)
	}

	return l
}

// Acquire acquires a concurrency slot. Call Release when done.
// Blocks if the max concurrent requests limit is reached.
// Returns an error if ctx is cancelled while waiting.
func (l *Limiter) Acquire(ctx context.Context) error {
	select {
	case l.concurrency <- struct{}{}:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("request cancelled while waiting for slot")
	}
}

// Release releases a concurrency slot acquired by Acquire.
func (l *Limiter) Release() {
	<-l.concurrency
}

// Wait blocks until the token bucket allows the given number of bytes.
// This is called before serving each chunk of content to enforce bandwidth limits.
func (l *Limiter) Wait(ctx context.Context, bytes int64) error {
	if l.bandwidth == nil {
		return nil // No bandwidth limit configured
	}
	return l.bandwidth.wait(ctx, bytes)
}

// Stats returns current limiter statistics.
func (l *Limiter) Stats() LimiterStats {
	active := len(l.concurrency)
	cap := cap(l.concurrency)
	return LimiterStats{
		ActiveRequests: active,
		MaxRequests:    cap,
		MemLimitMB:     l.cfg.Limits.RAMMB,
		BandwidthMbps:  l.cfg.Limits.BandwidthMbps,
	}
}

// LimiterStats holds current resource usage information.
type LimiterStats struct {
	ActiveRequests int `json:"active_requests"`
	MaxRequests    int `json:"max_requests"`
	MemLimitMB     int `json:"mem_limit_mb"`
	BandwidthMbps  int `json:"bandwidth_mbps"`
}

// tokenBucket is a simple token bucket rate limiter.
// It refills at a constant rate and blocks callers when empty.
type tokenBucket struct {
	mu       sync.Mutex
	tokens   int64
	capacity int64
	refillPS int64 // tokens added per second
	lastFill time.Time
}

// newTokenBucket creates a token bucket with the given bytes-per-second rate.
func newTokenBucket(bytesPerSec int64) *tokenBucket {
	return &tokenBucket{
		tokens:   bytesPerSec, // start full
		capacity: bytesPerSec,
		refillPS: bytesPerSec,
		lastFill: time.Now(),
	}
}

// wait blocks until the bucket has enough tokens for the requested bytes.
func (tb *tokenBucket) wait(ctx context.Context, bytes int64) error {
	for {
		tb.mu.Lock()
		tb.refill()

		if tb.tokens >= bytes {
			tb.tokens -= bytes
			tb.mu.Unlock()
			return nil
		}

		// Not enough tokens — calculate wait time
		needed := bytes - tb.tokens
		waitMS := (needed * 1000) / tb.refillPS
		if waitMS < 1 {
			waitMS = 1
		}
		tb.mu.Unlock()

		// Wait for tokens to refill
		select {
		case <-time.After(time.Duration(waitMS) * time.Millisecond):
			// Try again
		case <-ctx.Done():
			return fmt.Errorf("request cancelled while waiting for bandwidth")
		}
	}
}

// refill adds tokens based on elapsed time since last fill.
// Must be called with mu held.
func (tb *tokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastFill)
	tb.lastFill = now

	newTokens := int64(elapsed.Seconds() * float64(tb.refillPS))
	tb.tokens += newTokens
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
}
