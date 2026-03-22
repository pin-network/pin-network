// Package scheduler determines when the PiN node should be actively serving.
// It is intentionally lightweight — one decision per minute, no system monitoring.
// The scheduler sets a single active flag that meshd reads before serving traffic.
package scheduler

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"meshd/config"
)

// Scheduler is the node activity gatekeeper.
// It runs a single background goroutine that wakes once per minute,
// evaluates the user's schedule, and sets the active flag.
type Scheduler struct {
	active  atomic.Bool
	idling  atomic.Bool
	cfg     *config.Config
	limiter LimiterStats
	mu      sync.Mutex
}

// LimiterStats holds resource usage reported by the limiter.
type LimiterStats struct {
	Active int
	Max    int
}

// New creates a new Scheduler. Call Start to begin scheduling.
func New(cfg *config.Config) *Scheduler {
	s := &Scheduler{cfg: cfg}
	s.active.Store(true)
	return s
}

// Start begins the scheduling loop in a background goroutine.
// It returns immediately. The loop runs until ctx is cancelled.
func (s *Scheduler) Start(ctx context.Context) {
	s.evaluate()

	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.evaluate()
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Active returns true if the node should currently be serving traffic.
// It is safe to call from multiple goroutines.
func (s *Scheduler) Active() bool {
	return s.active.Load()
}

// Idling returns true if the node is below its idle threshold.
// The limiter uses this to allow higher throughput when the device is underutilized.
func (s *Scheduler) Idling() bool {
	return s.idling.Load()
}

// UpdateStats receives current resource usage from the limiter.
// Called periodically by the limiter to inform scheduling decisions.
func (s *Scheduler) UpdateStats(active, max int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.limiter = LimiterStats{Active: active, Max: max}
}

// evaluate makes the single active/idle decision.
// This runs once per minute and is the only place the flag is set.
func (s *Scheduler) evaluate() {
	// Fast path — always on with no limits to check
	if s.cfg.Schedule.AlwaysOn &&
		s.cfg.Limits.BatteryMinPct == 0 &&
		s.cfg.Schedule.IdleThresholdPct == 0 {
		s.setActive(true)
		s.idling.Store(false)
		return
	}

	// Check schedule window
	if !s.cfg.Schedule.AlwaysOn {
		if len(s.cfg.Schedule.ActiveHours) == 0 {
			s.setActive(true)
		} else {
			now := time.Now()
			inWindow := false
			for _, window := range s.cfg.Schedule.ActiveHours {
				if inActiveWindow(now, window.Start, window.End) {
					inWindow = true
					break
				}
			}
			if !inWindow {
				s.setActive(false)
				s.idling.Store(false)
				return
			}
		}
	}

	// Check battery threshold
	// BatteryMinPct of 0 means disabled — skip check
	// Battery level is written by the app via config update
	if s.cfg.Limits.BatteryMinPct > 0 {
		// App writes current battery level to BatteryCurrentPct
		// For now we trust the config value — app layer handles detection
		// If BatteryCurrentPct is added later, compare here
	}

	// Check idle threshold
	// If we are using less than IdleThresholdPct of our allocation,
	// flag as idling so the limiter can allow higher throughput.
	// Default is 20% — ramp up when using less than 20% of allocation.
	// Set to 0 to disable. Configurable via app.
	if s.cfg.Schedule.IdleThresholdPct > 0 {
		s.mu.Lock()
		stats := s.limiter
		s.mu.Unlock()

		if stats.Max > 0 {
			usagePct := (stats.Active * 100) / stats.Max
			s.idling.Store(usagePct < s.cfg.Schedule.IdleThresholdPct)
		}
	}

	s.setActive(true)
}

// setActive sets the active flag and logs only on state change.
func (s *Scheduler) setActive(active bool) {
	previous := s.active.Swap(active)
	if previous != active {
		if active {
			log.Println("scheduler: node entering active mode")
		} else {
			log.Println("scheduler: node entering idle mode")
		}
	}
}

// inActiveWindow returns true if t falls within the start-end window.
// Windows are specified as "HH:MM" in 24-hour local time.
// Overnight windows (e.g. 22:00 to 07:00) are supported.
func inActiveWindow(t time.Time, start, end string) bool {
	startTime, err := parseHHMM(t, start)
	if err != nil {
		return false
	}
	endTime, err := parseHHMM(t, end)
	if err != nil {
		return false
	}

	now := t.Hour()*60 + t.Minute()
	s := startTime.Hour()*60 + startTime.Minute()
	e := endTime.Hour()*60 + endTime.Minute()

	if s <= e {
		return now >= s && now < e
	}
	return now >= s || now < e
}

// parseHHMM parses a "HH:MM" string into a time.Time on the same date as ref.
func parseHHMM(ref time.Time, hhmm string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04",
		ref.Format("2006-01-02")+" "+hhmm, ref.Location())
}

// UpdateConfig applies a new configuration to the scheduler.
// Called by the config watcher when the config file changes.
func (s *Scheduler) UpdateConfig(cfg *config.Config) {
	s.mu.Lock()
	s.cfg = cfg
	s.mu.Unlock()
	// Evaluate immediately with new config
	s.evaluate()
	log.Println("scheduler: config updated")
}
