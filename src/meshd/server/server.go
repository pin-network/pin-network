// Package server provides the local HTTP API used by the browser and tray app.
package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"meshd/config"
	"meshd/ledger"
	"meshd/node"
	"meshd/store"
)

// API is the local HTTP API server.
type API struct {
	cfg   *config.Config
	node  *node.Node
	db    *ledger.DB
	store *store.Store
	mux   *http.ServeMux
}

// NewAPI creates a new API server.
func NewAPI(cfg *config.Config, n *node.Node, db *ledger.DB, st *store.Store) *API {
	a := &API{
		cfg:   cfg,
		node:  n,
		db:    db,
		store: st,
		mux:   http.NewServeMux(),
	}
	a.registerRoutes()
	return a
}

// ListenAndServe starts the API server on the configured port.
func (a *API) ListenAndServe() error {
	addr := fmt.Sprintf("127.0.0.1:%d", a.cfg.Network.APIPort)
	return http.ListenAndServe(addr, a.mux)
}

// registerRoutes wires up all API endpoints.
func (a *API) registerRoutes() {
	a.mux.HandleFunc("/api/v1/status", a.handleStatus)
	a.mux.HandleFunc("/api/v1/peers", a.handlePeers)
	a.mux.HandleFunc("/api/v1/ledger", a.handleLedger)
	a.mux.HandleFunc("/api/v1/content", a.handleContent)
	a.mux.HandleFunc("/api/v1/content/", a.handleContentGet)
}

// StatusResponse is the response for GET /api/v1/status.
type StatusResponse struct {
	NodeID  string   `json:"node_id"`
	Tier    int      `json:"tier"`
	Addrs   []string `json:"addrs"`
	Version string   `json:"version"`
	Online  bool     `json:"online"`
}

func (a *API) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, StatusResponse{
		NodeID:  a.node.ID(),
		Tier:    a.cfg.Node.Tier,
		Addrs:   a.node.Addrs(),
		Version: "0.1.0-dev",
		Online:  true,
	})
}

// PeersResponse is the response for GET /api/v1/peers.
type PeersResponse struct {
	Count int      `json:"count"`
	Peers []string `json:"peers"`
}

func (a *API) handlePeers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	peers := a.node.Peers()
	writeJSON(w, PeersResponse{Count: len(peers), Peers: peers})
}

// LedgerResponse is the response for GET /api/v1/ledger.
type LedgerResponse struct {
	Balance       float64 `json:"balance_hashes"`
	BytesServed   int64   `json:"bytes_served_today"`
	UptimeMinutes int64   `json:"uptime_minutes_today"`
	UptimePct     float64 `json:"uptime_pct_today"`
}

func (a *API) handleLedger(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	balance, err := a.db.Balance()
	if err != nil {
		http.Error(w, "ledger error", http.StatusInternalServerError)
		return
	}

	bytesServed, err := a.db.BytesServedToday()
	if err != nil {
		http.Error(w, "ledger error", http.StatusInternalServerError)
		return
	}

	uptimeMinutes, err := a.db.UptimeToday()
	if err != nil {
		http.Error(w, "ledger error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, LedgerResponse{
		Balance:       balance,
		BytesServed:   bytesServed,
		UptimeMinutes: uptimeMinutes,
		UptimePct:     float64(uptimeMinutes) / 14.40,
	})
}

// ContentListResponse is the response for GET /api/v1/content.
type ContentListResponse struct {
	Count int      `json:"count"`
	CIDs  []string `json:"cids"`
	Bytes int64    `json:"total_bytes"`
}

// ContentPutResponse is the response for POST /api/v1/content.
type ContentPutResponse struct {
	CID   string `json:"cid"`
	Bytes int    `json:"bytes"`
}

// handleContent handles GET (list) and POST (store) for /api/v1/content.
func (a *API) handleContent(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		cids, err := a.store.List()
		if err != nil {
			http.Error(w, "store error", http.StatusInternalServerError)
			return
		}
		size, err := a.store.Size()
		if err != nil {
			http.Error(w, "store error", http.StatusInternalServerError)
			return
		}
		writeJSON(w, ContentListResponse{
			Count: len(cids),
			CIDs:  cids,
			Bytes: size,
		})

	case http.MethodPost:
		data, err := io.ReadAll(io.LimitReader(r.Body, 50*1024*1024))
		if err != nil {
			http.Error(w, "reading body", http.StatusBadRequest)
			return
		}
		if len(data) == 0 {
			http.Error(w, "empty body", http.StatusBadRequest)
			return
		}

		cid, err := a.store.Put(data)
		if err != nil {
			http.Error(w, "store error", http.StatusInternalServerError)
			return
		}

		writeJSON(w, ContentPutResponse{CID: cid, Bytes: len(data)})

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleContentGet handles GET /api/v1/content/{cid}
func (a *API) handleContentGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cid := strings.TrimPrefix(r.URL.Path, "/api/v1/content/")
	if cid == "" {
		http.Error(w, "missing CID", http.StatusBadRequest)
		return
	}

	data, err := a.store.Get(cid)
	if err != nil {
		if err == os.ErrNotExist {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "store error", http.StatusInternalServerError)
		return
	}

	contentType := http.DetectContentType(data)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("X-Content-CID", cid)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(data)
}

// writeJSON writes a JSON response with correct headers.
func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
	}
}
