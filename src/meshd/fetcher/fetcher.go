// Package fetcher handles retrieving content from peer nodes.
package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// Fetcher retrieves content from peer nodes by CID.
type Fetcher struct {
	client *http.Client
}

// New creates a new Fetcher.
func New() *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchFromPeer retrieves content with the given CID from a peer's API.
// peerAddr should be the peer's API address, e.g. "192.168.7.200:4002"
func (f *Fetcher) FetchFromPeer(peerAddr string, cid string) ([]byte, error) {
	url := fmt.Sprintf("http://%s/api/v1/content/%s", peerAddr, cid)

	resp, err := f.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching from peer %s: %w", peerAddr, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("peer %s does not have CID %s", peerAddr, cid)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("peer %s returned status %d", peerAddr, resp.StatusCode)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, 50*1024*1024))
	if err != nil {
		return nil, fmt.Errorf("reading response from peer %s: %w", peerAddr, err)
	}

	return data, nil
}

// FetchFromPeers tries each peer in order until one returns the content.
func (f *Fetcher) FetchFromPeers(peerAddrs []string, cid string) ([]byte, string, error) {
	for _, addr := range peerAddrs {
		data, err := f.FetchFromPeer(addr, cid)
		if err != nil {
			continue
		}
		return data, addr, nil
	}
	return nil, "", fmt.Errorf("content %s not found on any peer", cid)
}
