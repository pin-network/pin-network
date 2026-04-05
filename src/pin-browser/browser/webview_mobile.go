//go:build android || ios

// WebView launcher for mobile platforms.
// On mobile the browser shell is embedded in the native app container.
// The Go layer provides the .pin resolver and mesh connectivity.
// The UI is handled by the native app (Phase 4 — Tauri/React Native).
package browser

import (
	"context"
	"log"
)

// launchWebView on mobile starts the resolver service only.
// The native UI layer handles rendering via the platform WebView.
func (b *Browser) launchWebView(ctx context.Context) error {
	log.Println("browser: mobile mode — resolver service running")
	log.Println("browser: UI handled by native app layer")

	// On mobile the Go layer runs as a background service
	// exposing the resolver via the proxy server
	// The native app calls into Go via gomobile bindings
	b.headless = true
	return b.startProxy(ctx)
}
