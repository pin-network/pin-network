// Package manifest handles .pin site manifest creation and parsing.
// A manifest describes a complete .pin website — its files, their CIDs,
// and metadata needed to serve it from the mesh.
package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	// Version is the current manifest format version.
	Version = 1

	// DefaultTTLHours is the default domain TTL (48 hours).
	DefaultTTLHours = 48
)

// Manifest describes a complete .pin website.
type Manifest struct {
	Version     int         `json:"version"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	Created     time.Time   `json:"created"`
	Updated     time.Time   `json:"updated"`
	Entrypoint  string      `json:"entrypoint"`
	Files       []FileEntry `json:"files"`
}

// FileEntry describes a single file in a .pin site.
type FileEntry struct {
	Path string `json:"path"`
	CID  string `json:"cid"`
	Size int64  `json:"size"`
	MIME string `json:"mime"`
}

// New creates a new empty manifest for the given domain name.
func New(name string) *Manifest {
	now := time.Now().UTC()
	return &Manifest{
		Version: Version,
		Name:    name,
		Created: now,
		Updated: now,
		Files:   []FileEntry{},
	}
}

// AddFile adds a file entry to the manifest.
func (m *Manifest) AddFile(path, cid string, size int64, mime string) {
	// Detect MIME if not provided
	if mime == "" {
		mime = mimeForPath(path)
	}

	// Set entrypoint to the first index.html found
	if m.Entrypoint == "" && (path == "index.html" || path == "/index.html") {
		m.Entrypoint = cid
	}

	m.Files = append(m.Files, FileEntry{
		Path: path,
		CID:  cid,
		Size: size,
		MIME: mime,
	})
	m.Updated = time.Now().UTC()
}

// CID returns the content identifier for this manifest.
// The manifest CID is the SHA-256 hash of its canonical JSON representation.
func (m *Manifest) CID() (string, error) {
	data, err := m.Marshal()
	if err != nil {
		return "", fmt.Errorf("marshalling manifest: %w", err)
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

// Marshal returns the canonical JSON representation of the manifest.
func (m *Manifest) Marshal() ([]byte, error) {
	return json.MarshalIndent(m, "", "  ")
}

// Parse parses a manifest from JSON bytes.
func Parse(data []byte) (*Manifest, error) {
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing manifest: %w", err)
	}
	if m.Version != Version {
		return nil, fmt.Errorf("unsupported manifest version %d", m.Version)
	}
	if m.Name == "" {
		return nil, fmt.Errorf("manifest missing name")
	}
	return &m, nil
}

// FindFile returns the FileEntry for the given path, or nil if not found.
func (m *Manifest) FindFile(path string) *FileEntry {
	// Normalize path
	if path == "" || path == "/" {
		path = "index.html"
	}
	// Strip leading slash
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	for i, f := range m.Files {
		if f.Path == path {
			return &m.Files[i]
		}
	}
	return nil
}

// EntrypointCID returns the CID of the site's entrypoint file.
// If no explicit entrypoint is set, returns the CID of index.html.
func (m *Manifest) EntrypointCID() string {
	if m.Entrypoint != "" {
		return m.Entrypoint
	}
	f := m.FindFile("index.html")
	if f != nil {
		return f.CID
	}
	// Fall back to first file
	if len(m.Files) > 0 {
		return m.Files[0].CID
	}
	return ""
}

// PublishRequest is the API request for publishing a .pin site.
type PublishRequest struct {
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Files       []PublishFile `json:"files"`
	TTLHours    int           `json:"ttl_hours,omitempty"`
}

// PublishFile is a single file in a publish request.
type PublishFile struct {
	Path    string `json:"path"`
	Content []byte `json:"content"`
	MIME    string `json:"mime,omitempty"`
}

// PublishResponse is the API response for a successful publish.
type PublishResponse struct {
	Name        string      `json:"name"`
	ManifestCID string      `json:"manifest_cid"`
	Files       []FileEntry `json:"files"`
	TTLHours    int         `json:"ttl_hours"`
}

// mimeForPath returns a MIME type based on file extension.
func mimeForPath(path string) string {
	// Find extension
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			ext := path[i:]
			switch ext {
			case ".html", ".htm":
				return "text/html"
			case ".css":
				return "text/css"
			case ".js":
				return "application/javascript"
			case ".json":
				return "application/json"
			case ".png":
				return "image/png"
			case ".jpg", ".jpeg":
				return "image/jpeg"
			case ".gif":
				return "image/gif"
			case ".svg":
				return "image/svg+xml"
			case ".ico":
				return "image/x-icon"
			case ".woff":
				return "font/woff"
			case ".woff2":
				return "font/woff2"
			case ".ttf":
				return "font/ttf"
			case ".txt":
				return "text/plain"
			case ".md":
				return "text/markdown"
			case ".xml":
				return "application/xml"
			case ".pdf":
				return "application/pdf"
			case ".zip":
				return "application/zip"
			case ".mp4":
				return "video/mp4"
			case ".mp3":
				return "audio/mpeg"
			case ".webp":
				return "image/webp"
			}
			break
		}
	}
	return http.DetectContentType([]byte{})
}
