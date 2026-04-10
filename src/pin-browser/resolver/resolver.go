// Package resolver handles .pin domain resolution through the local meshd API.
package resolver

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	PinTLD         = ".pin"
	DefaultAPIAddr = "127.0.0.1:4002"
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

// New creates a new Resolver.
func New(apiAddr string) *Resolver {
	if apiAddr == "" {
		apiAddr = DefaultAPIAddr
	}
	return &Resolver{
		apiAddr: apiAddr,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// IsPinURL returns true if the URL targets a .pin domain.
func IsPinURL(url string) bool {
	u := strings.ToLower(url)
	for _, scheme := range []string{"http://", "https://", "pin://"} {
		u = strings.TrimPrefix(u, scheme)
	}
	host := strings.SplitN(u, "/", 2)[0]
	host = strings.SplitN(host, ":", 2)[0]
	return strings.HasSuffix(host, PinTLD)
}

// Resolve fetches content for a .pin URL from the local meshd node.
func (r *Resolver) Resolve(url string) (*Response, error) {
	if !IsPinURL(url) {
		return nil, nil
	}

	domain, path := parsePinURL(url)

	manifestCID, err := r.resolveDomain(domain)
	if err != nil {
		return nil, fmt.Errorf("resolving .pin domain %s: %w", domain, err)
	}

	if path == "" || path == "/" {
		return r.fetchManifestEntrypoint(manifestCID)
	}

	return r.fetchManifestPath(manifestCID, path)
}

// ResolveCID fetches content directly by CID.
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

	cid := strings.TrimSpace(string(body))
	if cid == "" {
		return "", fmt.Errorf("empty CID returned for domain %s", domain)
	}
	return cid, nil
}

// fetchManifestEntrypoint fetches a manifest and returns its entrypoint content.
func (r *Resolver) fetchManifestEntrypoint(manifestCID string) (*Response, error) {
	manifestResp, err := r.fetchCID(manifestCID)
	if err != nil {
		return nil, fmt.Errorf("fetching manifest: %w", err)
	}

	// Always try to parse as manifest JSON regardless of content type
	body := string(manifestResp.Body)
	if idx := strings.Index(body, `"entrypoint": "`); idx >= 0 {
		start := idx + len(`"entrypoint": "`)
		end := strings.Index(body[start:], `"`)
		if end >= 0 {
			entrypointCID := body[start : start+end]
			if entrypointCID != "" {
				return r.fetchCID(entrypointCID)
			}
		}
	}

	// Not a manifest or no entrypoint — return as-is
	return manifestResp, nil
}

// fetchManifestPath fetches a specific file path from a manifest.
func (r *Resolver) fetchManifestPath(manifestCID, path string) (*Response, error) {
	manifestResp, err := r.fetchCID(manifestCID)
	if err != nil {
		return nil, fmt.Errorf("fetching manifest: %w", err)
	}

	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	body := string(manifestResp.Body)
	searchKey := fmt.Sprintf(`"path":"%s"`, path)
	if idx := strings.Index(body, searchKey); idx >= 0 {
		cidKey := `"cid":"`
		if cidIdx := strings.Index(body[idx:], cidKey); cidIdx >= 0 {
			start := idx + cidIdx + len(cidKey)
			end := strings.Index(body[start:], `"`)
			if end >= 0 {
				return r.fetchCID(body[start : start+end])
			}
		}
	}

	return r.fetchManifestEntrypoint(manifestCID)
}

func (r *Resolver) fetchCID(cid string) (*Response, error) {
	url := fmt.Sprintf("http://%s/api/v1/content/%s", r.apiAddr, cid)
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching CID %s: %w", cid, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusServiceUnavailable {
		return nil, fmt.Errorf("mesh node is idle")
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

func parsePinURL(url string) (domain, path string) {
	u := url
	for _, scheme := range []string{"pin://", "http://", "https://"} {
		u = strings.TrimPrefix(u, scheme)
	}
	parts := strings.SplitN(u, "/", 2)
	host := parts[0]
	if len(parts) > 1 {
		path = "/" + parts[1]
	} else {
		path = "/"
	}
	domain = strings.SplitN(host, ":", 2)[0]
	return domain, path
}
