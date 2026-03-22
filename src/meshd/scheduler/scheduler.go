// Package scheduler determines when the PiN node should be actively serving.
// It is intentionally lightweight — one decision per minute, no system monitoring.
// The scheduler sets a single active flag that meshd reads before serving traffic.
package scheduler

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"meshd/config"
)

// Scheduler is the node activity gatekeeper.
// It runs a single background goroutine that wakes once per minute,
// evaluates the user's schedule, and sets the active flag.
type Scheduler struct {
	active atomic.Bool
	cfg    *config.Config
}

// New creates a new Scheduler. Call Start to begin scheduling.
func New(cfg *config.Config) *Scheduler {
	s := &Scheduler{cfg: cfg}
	// Default to active so the node serves immediately on startup
	// if no schedule is configured.
	s.active.Store(true)
	return s
}

// Start begins the scheduling loop in a background goroutine.
// It returns immediately. The loop runs until ctx is cancelled.
func (s *Scheduler) Start(ctx context.Context) {
	// Evaluate immediately on startup
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
// This is the single flag that meshd checks before accepting requests.
// It is safe to call from multiple goroutines.
func (s *Scheduler) Active() bool {
	return s.active.Load()
}

// evaluate makes the single active/idle decision.
// This runs once per minute and is the only place the flag is set.
func (s *Scheduler) evaluate() {
	// If always_on is set, never idle — simple fast path
	if s.cfg.Schedule.AlwaysOn {
		s.setActive(true)
		return
	}

	// Check if current time falls within any active window
	if len(s.cfg.Schedule.ActiveHours) == 0 {
		// No schedule defined — default to active
		s.setActive(true)
		return
	}

	now := time.Now()
	inWindow := false

	for _, window := range s.cfg.Schedule.ActiveHours {
		if inActiveWindow(now, window.Start, window.End) {
			inWindow = true
			break
		}
	}

	s.setActive(inWindow)
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
		// Normal window: e.g. 09:00 to 17:00
		return now >= s && now < e
	}
	// Overnight window: e.g. 22:00 to 07:00
	return now >= s || now < e
}

// parseHHMM parses a "HH:MM" string into a time.Time on the same date as ref.
func parseHHMM(ref time.Time, hhmm string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04",
		ref.Format("2006-01-02")+" "+hhmm, ref.Location())
}
