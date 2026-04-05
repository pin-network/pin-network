// Package resolver handles .pin domain resolution through the local meshd API.
// It is the core of the PiN browser — intercepting .pin requests and routing
// them through the mesh while passing everything else to normal DNS.
package resolver

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	// PinTLD is the mesh top-level domain.
	PinTLD = ".pin"

	// DefaultAPIAddr is the local meshd API address.
	DefaultAPIAddr = "127.0.0.1:4002"

	// MaxContentSize is the maximum content size we'll fetch (50MB).
	MaxContentSize = 50 * 1024 * 1024
)

// Resolver resolves .pin domains through the local meshd API.
type Resolver struct {
	apiAddr string
	client  *http.Client
}

// Response holds the result of a resolved request.
type Response struct {
	Body        []byte
	ContentType string
	CID         string
	FromMesh    bool
	StatusCode  int
}

// New creates a new Resolver pointing at the local meshd API.
func New(apiAddr string) *Resolver {
	if apiAddr == "" {
		apiAddr = DefaultAPIAddr
	}
	return &Resolver{
		apiAddr: apiAddr,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IsPinURL returns true if the URL targets a .pin domain.
func IsPinURL(url string) bool {
	// Normalize — strip scheme
	u := strings.ToLower(url)
	u = strings.TrimPrefix(u, "http://")
	u = strings.TrimPrefix(u, "https://")
	u = strings.TrimPrefix(u, "pin://")

	// Check if the host ends in .pin
	host := strings.SplitN(u, "/", 2)[0]
	host = strings.SplitN(host, ":", 2)[0] // strip port
	return strings.HasSuffix(host, PinTLD)
}

// Resolve fetches content for a .pin URL from the local meshd node.
// For non-.pin URLs it returns nil and the caller should use normal HTTP.
func (r *Resolver) Resolve(url string) (*Response, error) {
	if !IsPinURL(url) {
		return nil, nil
	}

	// Parse the .pin URL into domain and path
	domain, path := parsePinURL(url)

	// Step 1 — resolve the domain to a CID via meshd
	cid, err := r.resolveDomain(domain)
	if err != nil {
		return nil, fmt.Errorf("resolving .pin domain %s: %w", domain, err)
	}

	// Step 2 — if path is not empty, resolve the specific file CID
	// For now we fetch the manifest CID and serve the entrypoint
	// Full path resolution comes in Phase 3.2
	targetCID := cid
	if path != "" && path != "/" {
		targetCID, err = r.resolvePathCID(cid, path)
		if err != nil {
			// Fall back to root CID
			targetCID = cid
		}
	}

	// Step 3 — fetch the content by CID
	return r.fetchCID(targetCID)
}

// ResolveCID fetches content directly by CID from the local meshd node.
// Used when the browser has a direct CID reference.
func (r *Resolver) ResolveCID(cid string) (*Response, error) {
	return r.fetchCID(cid)
}

// Healthy returns true if the local meshd API is reachable.
func (r *Resolver) Healthy() bool {
	resp, err := r.client.Get(fmt.Sprintf("http://%s/api/v1/status", r.apiAddr))
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// resolveDomain looks up a .pin domain name in the meshd DHT
// and returns the root CID for that domain.
func (r *Resolver) resolveDomain(domain string) (string, error) {
	url := fmt.Sprintf("http://%s/api/v1/domain/%s", r.apiAddr, domain)
	resp, err := r.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("querying meshd for domain %s: %w", domain, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("domain %s not found in mesh", domain)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("meshd returned status %d for domain %s", resp.StatusCode, domain)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024))
	if err != nil {
		return "", fmt.Errorf("reading domain response: %w", err)
	}

	// Response is a plain CID string for now
	// Will be JSON manifest in Phase 3.2
	cid := strings.TrimSpace(string(body))
	if cid == "" {
		return "", fmt.Errorf("empty CID returned for domain %s", domain)
	}

	return cid, nil
}

// resolvePathCID resolves a specific file path within a .pin site.
// Takes the root manifest CID and a path, returns the file's CID.
func (r *Resolver) resolvePathCID(manifestCID, path string) (string, error) {
	// Phase 3.2 — manifest parsing and path resolution
	// For now return the manifest CID as a fallback
	_ = manifestCID
	_ = path
	return manifestCID, nil
}

// fetchCID retrieves content by CID from the local meshd node.
func (r *Resolver) fetchCID(cid string) (*Response, error) {
	url := fmt.Sprintf("http://%s/api/v1/content/%s", r.apiAddr, cid)
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching CID %s: %w", cid, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusServiceUnavailable {
		return nil, fmt.Errorf("mesh node is idle — content serving paused")
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("content %s not found in mesh", cid)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meshd returned status %d for CID %s", resp.StatusCode, cid)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, MaxContentSize))
	if err != nil {
		return nil, fmt.Errorf("reading content for CID %s: %w", cid, err)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = http.DetectContentType(body)
	}

	return &Response{
		Body:        body,
		ContentType: contentType,
		CID:         cid,
		FromMesh:    true,
		StatusCode:  http.StatusOK,
	}, nil
}

// parsePinURL extracts the domain and path from a .pin URL.
func parsePinURL(url string) (domain, path string) {
	// Strip scheme
	u := url
	for _, scheme := range []string{"pin://", "http://", "https://"} {
		u = strings.TrimPrefix(u, scheme)
	}

	// Split host and path
	parts := strings.SplitN(u, "/", 2)
	host := parts[0]
	if len(parts) > 1 {
		path = "/" + parts[1]
	} else {
		path = "/"
	}

	// Strip port from host
	domain = strings.SplitN(host, ":", 2)[0]
	return domain, path
}
