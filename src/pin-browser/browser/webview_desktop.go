//go:build (windows || linux || darwin) && !headless

// WebView launcher for desktop platforms.
// Uses the webview/webview_go library which wraps:
//   - Windows: WebView2 (Chromium-based)
//   - Linux: WebKitGTK
//   - macOS: WKWebView
package browser

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/webview/webview_go"
)

// launchWebView starts the native WebView window on desktop platforms.
func (b *Browser) launchWebView(ctx context.Context) error {
	// Start a local HTTP server that the WebView will talk to
	// This server intercepts .pin requests and proxies them through meshd
	proxyAddr := "127.0.0.1:7071"
	go b.startInterceptServer(proxyAddr)

	// Create the WebView
	w := webview.New(true)
	defer w.Destroy()

	w.SetTitle("PiN Browser")
	w.SetSize(1024, 768, webview.HintNone)

	// Inject the .pin URL handler into the WebView
	// When the user types a .pin URL, we intercept and reroute
	w.Bind("pinNavigate", func(url string) {
		resp, err := b.resolver.Resolve(url)
		if err != nil {
			log.Printf("browser: resolve error: %v", err)
			w.Navigate(fmt.Sprintf("data:text/html,<h1>Mesh Error</h1><p>%v</p>", err))
			return
		}
		if resp != nil {
			// Serve through local intercept server
			w.Navigate(fmt.Sprintf("http://%s/pin/%s", proxyAddr, resp.CID))
		}
	})

	// Navigate to the PiN home page
	w.Navigate(fmt.Sprintf("http://%s/home", proxyAddr))

	// Handle context cancellation
	go func() {
		<-ctx.Done()
		w.Terminate()
	}()

	w.Run()
	return nil
}

// startInterceptServer runs a local HTTP server that the WebView uses
// to serve mesh content. This avoids CORS issues and lets us serve
// binary content (images, files) through the WebView cleanly.
func (b *Browser) startInterceptServer(addr string) {
	mux := http.NewServeMux()

	// Serve mesh content by CID
	mux.HandleFunc("/pin/", func(w http.ResponseWriter, r *http.Request) {
		cid := r.URL.Path[len("/pin/"):]
		resp, err := b.resolver.ResolveCID(cid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", resp.ContentType)
		w.Write(resp.Body)
	})

	// Serve the PiN browser home page
	mux.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(homePage))
	})

	log.Printf("browser: intercept server on %s", addr)
	http.ListenAndServe(addr, mux)
}

// homePage is the default PiN browser start page.
const homePage = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>PiN Browser</title>
<style>
  body { font-family: system-ui, sans-serif; background: #0f0f1a; color: #e0e0ff; 
         display: flex; flex-direction: column; align-items: center; 
         justify-content: center; height: 100vh; margin: 0; }
  h1 { font-size: 3rem; color: #7F77DD; margin-bottom: 0.5rem; }
  p { color: #9090bb; margin-bottom: 2rem; }
  #bar { display: flex; gap: 0.5rem; width: 600px; }
  input { flex: 1; padding: 0.75rem 1rem; border-radius: 8px; border: 1px solid #333; 
          background: #1a1a2e; color: #e0e0ff; font-size: 1rem; }
  button { padding: 0.75rem 1.5rem; border-radius: 8px; border: none; 
           background: #7F77DD; color: white; font-size: 1rem; cursor: pointer; }
  button:hover { background: #6a62cc; }
  .status { margin-top: 2rem; font-size: 0.8rem; color: #555; }
</style>
</head>
<body>
  <h1>PiN</h1>
  <p>Pi Integrated Network — Are you IN?</p>
  <div id="bar">
    <input id="url" type="text" placeholder="Enter a .pin domain or web address..." 
           onkeydown="if(event.key==='Enter') navigate()">
    <button onclick="navigate()">Go</button>
  </div>
  <div class="status" id="status">Connecting to mesh...</div>
  <script>
    function navigate() {
      const url = document.getElementById('url').value.trim();
      if (!url) return;
      if (url.endsWith('.pin') || url.includes('.pin/')) {
        window.pinNavigate(url);
      } else {
        window.location.href = url.startsWith('http') ? url : 'https://' + url;
      }
    }
    // Check mesh status
    fetch('/pin-status').then(r => r.json()).then(s => {
      document.getElementById('status').textContent = 
        s.mesh_healthy ? 'Connected to mesh' : 'Mesh offline — web browsing only';
    }).catch(() => {
      document.getElementById('status').textContent = 'Status unknown';
    });
  </script>
</body>
</html>`
