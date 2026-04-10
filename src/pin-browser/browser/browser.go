// Package browser provides the PiN browser backend.
// Runs as a local HTTP proxy server resolving .pin domains through meshd.
// Becomes the Go sidecar for the Tauri app in Phase 4.
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

const Version = "0.1.0-dev"

// Browser is the PiN browser backend.
type Browser struct {
	resolver *resolver.Resolver
	apiAddr  string
	port     int
}

// Config holds browser configuration.
type Config struct {
	APIAddr string
	Port    int
}

// New creates a new Browser instance.
func New(cfg Config) *Browser {
	if cfg.APIAddr == "" {
		cfg.APIAddr = resolver.DefaultAPIAddr
	}
	if cfg.Port == 0 {
		cfg.Port = 7070
	}
	return &Browser{
		resolver: resolver.New(cfg.APIAddr),
		apiAddr:  cfg.APIAddr,
		port:     cfg.Port,
	}
}

// Start starts the browser proxy server.
func (b *Browser) Start(ctx context.Context) error {
	if !b.resolver.Healthy() {
		log.Println("browser: warning — meshd not reachable at", b.apiAddr)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", b.handleRequest)
	mux.HandleFunc("/pin-status", b.handleStatus)
	mux.HandleFunc("/pin-resolve/", b.handleResolve)

	addr := fmt.Sprintf("0.0.0.0:%d", b.port)
	server := &http.Server{Addr: addr, Handler: mux}

	log.Printf("browser: PiN browser started on http://localhost:%d", b.port)
	log.Printf("browser: open http://localhost:%d in your browser", b.port)

	go func() {
		<-ctx.Done()
		server.Close()
	}()

	return server.ListenAndServe()
}

func (b *Browser) handleRequest(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	path := r.URL.Path
	isLocal := strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.0.0.1")

	if isLocal && (path == "/" || path == "") {
		b.serveHomePage(w)
		return
	}

	url := fmt.Sprintf("http://%s%s", host, path)
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}

	if resolver.IsPinURL(url) {
		b.servePinContent(w, url)
		return
	}

	b.serveHomePage(w)
}

func (b *Browser) handleResolve(w http.ResponseWriter, r *http.Request) {
	target := strings.TrimPrefix(r.URL.Path, "/pin-resolve/")
	if target == "" {
		http.Error(w, "missing domain", http.StatusBadRequest)
		return
	}

	resp, err := b.resolver.Resolve("pin://" + target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	if resp == nil {
		http.Error(w, "not a .pin URL", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", resp.ContentType)
	w.Header().Set("X-PiN-CID", resp.CID)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(resp.StatusCode)
	w.Write(resp.Body)
}

func (b *Browser) servePinContent(w http.ResponseWriter, url string) {
	resp, err := b.resolver.Resolve(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Mesh error: %v", err), http.StatusBadGateway)
		return
	}
	if resp == nil {
		http.Error(w, "Not a .pin URL", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", resp.ContentType)
	w.Header().Set("X-PiN-CID", resp.CID)
	w.WriteHeader(resp.StatusCode)
	w.Write(resp.Body)
}

func (b *Browser) handleStatus(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"version":      Version,
		"mesh_healthy": b.resolver.Healthy(),
		"api_addr":     b.apiAddr,
		"port":         b.port,
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(status)
}

func (b *Browser) serveHomePage(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(homePage))
}

const homePage = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>PiN Browser</title>
<style>
  * { box-sizing: border-box; margin: 0; padding: 0; }
  body { font-family: system-ui, sans-serif; background: #0f0f1a; color: #e0e0ff;
         display: flex; flex-direction: column; align-items: center;
         justify-content: center; min-height: 100vh; padding: 2rem; }
  .logo { font-size: 4rem; color: #7F77DD; font-weight: 700; letter-spacing: -2px; }
  .tagline { color: #6060aa; margin: 0.5rem 0 2rem; font-size: 1rem; }
  .bar { display: flex; gap: 0.5rem; width: 100%; max-width: 600px; }
  input { flex: 1; padding: 0.875rem 1.25rem; border-radius: 10px;
          border: 1px solid #2a2a4a; background: #1a1a2e; color: #e0e0ff;
          font-size: 1rem; outline: none; }
  input:focus { border-color: #7F77DD; }
  button { padding: 0.875rem 1.75rem; border-radius: 10px; border: none;
           background: #7F77DD; color: white; font-size: 1rem; cursor: pointer; font-weight: 600; }
  button:hover { background: #6a62cc; }
  .status { margin-top: 2rem; font-size: 0.8rem; color: #404060;
            display: flex; align-items: center; gap: 0.5rem; }
  .dot { width: 8px; height: 8px; border-radius: 50%; background: #404060; }
  .dot.online { background: #44dd88; }
  .examples { margin-top: 3rem; text-align: center; }
  .examples p { color: #404060; font-size: 0.8rem; margin-bottom: 0.75rem; }
  .example { display: inline-block; margin: 0.25rem; padding: 0.4rem 0.8rem;
             background: #1a1a2e; border: 1px solid #2a2a4a; border-radius: 6px;
             color: #7F77DD; font-size: 0.8rem; cursor: pointer; }
  .example:hover { border-color: #7F77DD; }
</style>
</head>
<body>
  <div class="logo">PiN</div>
  <p class="tagline">Pi Integrated Network — Are you IN?</p>
  <div class="bar">
    <input id="url" type="text" placeholder="Enter a .pin domain..."
           onkeydown="if(event.key==='Enter') navigate()">
    <button onclick="navigate()">Go</button>
  </div>
  <div class="status">
    <div class="dot" id="dot"></div>
    <span id="status-text">Checking mesh connection...</span>
  </div>
  <div class="examples">
    <p>Try a .pin domain</p>
    <span class="example" onclick="go('test.pin')">test.pin</span>
    <span class="example" onclick="go('hello.pin')">hello.pin</span>
  </div>
  <script>
    function navigate() {
      const url = document.getElementById('url').value.trim();
      if (!url) return;
      go(url);
    }
    function go(url) {
      if (url.endsWith('.pin') || url.includes('.pin/')) {
        fetch('/pin-resolve/' + url)
          .then(r => r.text())
          .then(html => { document.open(); document.write(html); document.close(); })
          .catch(e => alert('Mesh error: ' + e));
      } else {
        window.location.href = url.startsWith('http') ? url : 'https://' + url;
      }
    }
    fetch('/pin-status').then(r => r.json()).then(s => {
      const dot = document.getElementById('dot');
      const txt = document.getElementById('status-text');
      if (s.mesh_healthy) { dot.classList.add('online'); txt.textContent = 'Connected to mesh'; }
      else { txt.textContent = 'Mesh offline — start meshd to browse .pin sites'; }
    }).catch(() => { document.getElementById('status-text').textContent = 'Status unknown'; });
  </script>
</body>
</html>`
