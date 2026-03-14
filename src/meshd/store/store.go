// Package store manages content-addressed file storage for PiN nodes.
// Files are stored by the SHA-256 hash of their contents (CID).
package store

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Store is a content-addressed file store.
type Store struct {
	root string
}

// New creates a new Store rooted at the given directory.
func New(root string) (*Store, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, fmt.Errorf("creating store directory: %w", err)
	}
	return &Store{root: root}, nil
}

// Put stores data and returns its CID (SHA-256 hex string).
func (s *Store) Put(data []byte) (string, error) {
	cid := hashBytes(data)

	path := s.cidPath(cid)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", fmt.Errorf("creating store subdirectory: %w", err)
	}

	// Don't rewrite if already stored
	if _, err := os.Stat(path); err == nil {
		return cid, nil
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("writing content: %w", err)
	}

	return cid, nil
}

// PutFile stores a file from disk and returns its CID.
func (s *Store) PutFile(srcPath string) (string, int64, error) {
	f, err := os.Open(srcPath)
	if err != nil {
		return "", 0, fmt.Errorf("opening source file: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return "", 0, fmt.Errorf("reading source file: %w", err)
	}

	cid, err := s.Put(data)
	if err != nil {
		return "", 0, err
	}

	return cid, int64(len(data)), nil
}

// Get retrieves content by CID. Returns os.ErrNotExist if not found.
func (s *Store) Get(cid string) ([]byte, error) {
	path := s.cidPath(cid)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, os.ErrNotExist
		}
		return nil, fmt.Errorf("reading content: %w", err)
	}

	// Verify integrity
	actual := hashBytes(data)
	if actual != cid {
		// Corrupted — delete and return not found
		os.Remove(path)
		return nil, fmt.Errorf("content integrity check failed for CID %s", cid)
	}

	return data, nil
}

// Has returns true if the store contains the given CID.
func (s *Store) Has(cid string) bool {
	_, err := os.Stat(s.cidPath(cid))
	return err == nil
}

// Delete removes content by CID.
func (s *Store) Delete(cid string) error {
	return os.Remove(s.cidPath(cid))
}

// List returns all CIDs stored locally.
func (s *Store) List() ([]string, error) {
	var cids []string

	err := filepath.WalkDir(s.root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		// Filename is the full CID
		cids = append(cids, d.Name())
		return nil
	})

	return cids, err
}

// Size returns the total bytes stored.
func (s *Store) Size() (int64, error) {
	var total int64
	err := filepath.WalkDir(s.root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		total += info.Size()
		return nil
	})
	return total, err
}

// cidPath returns the filesystem path for a given CID.
// Files are stored in subdirectories based on the first 2 chars of the CID
// to avoid having too many files in a single directory.
func (s *Store) cidPath(cid string) string {
	if len(cid) < 2 {
		return filepath.Join(s.root, cid)
	}
	return filepath.Join(s.root, cid[:2], cid)
}

// hashBytes returns the SHA-256 hex string of data.
func hashBytes(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}
