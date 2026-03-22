//go:build !windows && !linux && !android && !darwin && !ios

// Priority setting fallback for all other platforms.
// Covers: FreeBSD, OpenBSD, NetBSD, Plan 9, Solaris, wasip1, js/wasm.
// These platforms either don't support setpriority or have
// platform-specific mechanisms we don't need to support yet.
// The no-op is safe — meshd will run at normal priority on these platforms.
package limits

import "log"

// SetLowPriority is a no-op on unsupported platforms.
// meshd will run at normal OS priority on this platform.
func SetLowPriority() {
	log.Println("limits: priority setting not supported on this platform (running at normal priority)")
}
