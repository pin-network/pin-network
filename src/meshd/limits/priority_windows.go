//go:build windows

// Priority setting for Windows.
// Sets the process to BELOW_NORMAL_PRIORITY_CLASS so meshd is
// naturally deprioritized by the OS scheduler without requiring
// any permissions or invasive monitoring.
package limits

import (
	"log"
	"syscall"
)

// Windows priority class constants
const (
	belowNormalPriorityClass = 0x00004000
)

var (
	kernel32          = syscall.NewLazyDLL("kernel32.dll")
	getCurrentProcess = kernel32.NewProc("GetCurrentProcess")
	setPriorityClass  = kernel32.NewProc("SetPriorityClass")
)

// SetLowPriority sets the current process to below-normal priority on Windows.
// This makes meshd invisible to task manager alerts and app monitors.
// No elevated permissions required.
func SetLowPriority() {
	handle, _, err := getCurrentProcess.Call()
	if handle == 0 {
		log.Printf("limits: could not get process handle: %v", err)
		return
	}

	ret, _, err := setPriorityClass.Call(handle, belowNormalPriorityClass)
	if ret == 0 {
		log.Printf("limits: could not set process priority: %v", err)
		return
	}

	log.Println("limits: process priority set to below-normal")
}
