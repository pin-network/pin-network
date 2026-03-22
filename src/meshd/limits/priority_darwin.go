//go:build darwin || ios

// Priority setting for macOS and iOS.
// Both platforms use the Darwin kernel and support the same
// setpriority syscall. On iOS, process priority adjustment is
// permitted for the app's own process without special entitlements.
package limits

import (
	"log"
	"syscall"
)

// SetLowPriority sets the current process to nice value 10 on macOS/iOS.
// This is identical to the Linux implementation as both use POSIX setpriority.
func SetLowPriority() {
	err := syscall.Setpriority(syscall.PRIO_PROCESS, 0, 10)
	if err != nil {
		log.Printf("limits: could not set process priority: %v", err)
		return
	}
	log.Println("limits: process nice value set to 10 (low priority)")
}
