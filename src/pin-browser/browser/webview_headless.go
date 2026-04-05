//go:build headless || (!windows && !linux && !darwin && !android && !ios)

// WebView launcher for headless and embedded platforms.
// Used on: RPi Zero, RISC-V, MIPS, any device without a display.
// Runs as a local HTTP proxy that other devices can use to access the mesh.
package browser

import (
	"context"
	"log"
)

// launchWebView on headless platforms starts the proxy server.
func (b *Browser) launchWebView(ctx context.Context) error {
	log.Println("browser: headless mode — starting proxy server")
	log.Println("browser: other devices can use this node as a mesh gateway")
	b.headless = true
	return b.startProxy(ctx)
}
