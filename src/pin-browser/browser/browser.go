// Package browser provides the PiN browser shell.
// It wraps the OS native WebView and intercepts .pin domain requests,
// routing them through the local meshd node while passing everything
// else to normal HTTP/DNS.
//
// Tier detection:
//   - Full display available: runs WebView with full UI
//   - Headless (no display): runs as a local HTTP proxy on port 7070
//     that other devices can use as a gateway to the mesh
package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"pin-browser/resolver"
)

const (
	// ProxyPort is the port the headless browser proxy listens on.
	ProxyPort = 7070

	// Version is the PiN browser version.
	Version = "0.1.0-dev"
)

// Browser is the PiN browser shell.
type Browser struct {
	resolver *resolver.Resolver
	headless bool
	apiAddr  string
}

// Config holds browser configuration.
type Config struct {
	APIAddr  string // meshd API address (default: 127.0.0.1:4002)
	Headless bool   // force headless mode
	HomePage string // default home page
}

// New creates a new Browser instance.
func New(cfg Config) *Browser {
	if cfg.APIAddr == "" {
		cfg.APIAddr = resolver.DefaultAPIAddr
	}
	return &Browser{
		resolver: resolver.New(cfg.APIAddr),
		headless: cfg.Headless,
		apiAddr:  cfg.APIAddr,
	}
}

// Start starts the browser.
// In headless mode it starts the proxy server.
// In full mode it launches the WebView window.
func (b *Browser) Start(ctx context.Context) error {
	// Check if meshd is running
	if !b.resolver.Healthy() {
		log.Println("browser: warning — meshd not reachable at", b.apiAddr)
		log.Println("browser: .pin domains will not resolve until meshd is running")
	}

	if b.headless {
		return b.startProxy(ctx)
	}
	return b.startWebView(ctx)
}

// startProxy starts the headless HTTP proxy server.
// Other devices on the local network can use this as a gateway
// to access .pin content through the mesh.
func (b *Browser) startProxy(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", b.handleProxyRequest)
	mux.HandleFunc("/pin-status", b.handleStatus)

	addr := fmt.Sprintf("0.0.0.0:%d", ProxyPort)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("browser: headless proxy starting on %s", addr)
	log.Printf("browser: .pin domains accessible via http://[device-ip]:%d/[domain.pin]/[path]", ProxyPort)

	go func() {
		<-ctx.Done()
		server.Close()
	}()

	return server.ListenAndServe()
}

// handleProxyRequest handles incoming proxy requests.
// Intercepts .pin domains and routes them through meshd.
// Passes everything else through to normal HTTP.
func (b *Browser) handleProxyRequest(w http.ResponseWriter, r *http.Request) {
	// Reconstruct the requested URL
	host := r.Host
	path := r.URL.Path
	if r.URL.RawQuery != "" {
		path += "?" + r.URL.RawQuery
	}

	url := fmt.Sprintf("http://%s%s", host, path)

	if resolver.IsPinURL(url) {
		b.servePinContent(w, url)
		return
	}

	// For non-.pin requests in proxy mode, return a helpful message
	// In Phase 3.2 we'll add full HTTP proxying for internet fallback
	http.Error(w, fmt.Sprintf(
		"PiN Browser Proxy\n\nTo access mesh content, request a .pin domain.\nExample: http://[proxy-ip]:%d/mysite.pin/\n\nFor internet access, configure your device to use a standard HTTP proxy.",
		ProxyPort,
	), http.StatusBadGateway)
}

// servePinContent resolves and serves .pin content.
func (b *Browser) servePinContent(w http.ResponseWriter, url string) {
	resp, err := b.resolver.Resolve(url)
	if err != nil {
		log.Printf("browser: failed to resolve %s: %v", url, err)
		http.Error(w, fmt.Sprintf("Failed to resolve mesh content: %v", err), http.StatusBadGateway)
		return
	}
	if resp == nil {
		http.Error(w, "Not a .pin URL", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", resp.ContentType)
	w.Header().Set("X-PiN-CID", resp.CID)
	w.Header().Set("X-PiN-Version", Version)
	w.WriteHeader(resp.StatusCode)
	w.Write(resp.Body)
}

// handleStatus returns browser and mesh status as JSON.
func (b *Browser) handleStatus(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"version":      Version,
		"headless":     b.headless,
		"mesh_healthy": b.resolver.Healthy(),
		"api_addr":     b.apiAddr,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// startWebView launches the native WebView browser window.
// The actual WebView integration is platform-specific and handled
// by the webview_*.go files in this package.
func (b *Browser) startWebView(ctx context.Context) error {
	log.Println("browser: starting WebView")
	return b.launchWebView(ctx)
}

// NavigateTo navigates the browser to the given URL.
// Handles both .pin and regular URLs.
func (b *Browser) NavigateTo(url string) {
	if strings.HasPrefix(url, "pin://") || resolver.IsPinURL(url) {
		log.Printf("browser: navigating to mesh URL: %s", url)
	} else {
		log.Printf("browser: navigating to web URL: %s", url)
	}
}
