// Package config — config reload watcher.
// Watches the config file for changes and reloads it automatically.
// Uses a simple modification time check every 30 seconds.
// No filesystem watchers or OS-specific APIs required.
package config

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Watcher watches the config file for changes and calls the provided
// callback when the file is modified. It checks every 30 seconds.
type Watcher struct {
	path     string
	lastMod  time.Time
	onChange func(*Config)
}

// NewWatcher creates a new config file watcher.
// onChange is called with the new config whenever the file changes.
func NewWatcher(path string, onChange func(*Config)) *Watcher {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil
		}
		path = filepath.Join(home, ".pin", "config.yaml")
	}

	// Record initial modification time
	var lastMod time.Time
	if info, err := os.Stat(path); err == nil {
		lastMod = info.ModTime()
	}

	return &Watcher{
		path:     path,
		lastMod:  lastMod,
		onChange: onChange,
	}
}

// Start begins watching the config file in a background goroutine.
// It returns immediately. The watcher runs until ctx is cancelled.
func (w *Watcher) Start(ctx context.Context) {
	if w == nil {
		return
	}

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				w.check()
			case <-ctx.Done():
				return
			}
		}
	}()
}

// check looks at the config file modification time and reloads if changed.
func (w *Watcher) check() {
	info, err := os.Stat(w.path)
	if err != nil {
		// File may have been temporarily removed during a write — ignore
		return
	}

	if info.ModTime().Equal(w.lastMod) {
		// No change
		return
	}

	// File has changed — reload
	cfg, err := Load(w.path)
	if err != nil {
		log.Printf("config watcher: failed to reload config: %v", err)
		return
	}

	w.lastMod = info.ModTime()
	log.Println("config watcher: config reloaded")
	w.onChange(cfg)
}
