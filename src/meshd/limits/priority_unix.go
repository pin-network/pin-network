//go:build linux || android

// Priority setting for Linux and Android.
// Covers: amd64, arm64, arm (Raspberry Pi), riscv64, mips, mipsle,
// mips64, mips64le, loong64, ppc64, ppc64le, s390x.
// Android uses the Linux kernel so the same syscall works.
package limits

import (
	"log"
	"syscall"
)

// SetLowPriority sets the current process to nice value 10 on Linux/Android.
// Nice 10 is well below normal (0) but not the lowest (19).
// This leaves headroom for other background processes while keeping
// meshd from competing with foreground apps.
// No elevated permissions required for processes to lower their own priority.
func SetLowPriority() {
	// PRIO_PROCESS = 0, pid = 0 means current process, niceness = 10
	err := syscall.Setpriority(syscall.PRIO_PROCESS, 0, 10)
	if err != nil {
		log.Printf("limits: could not set process priority: %v", err)
		return
	}
	log.Println("limits: process nice value set to 10 (low priority)")
}
