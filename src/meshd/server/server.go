// Package server provides the local HTTP API used by the browser and tray app.
package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"meshd/config"
	"meshd/ledger"
	"meshd/node"
)

// API is the local HTTP API server.
type API struct {
	cfg  *config.Config
	node *node.Node
	db   *ledger.DB
	mux  *http.ServeMux
}

// NewAPI creates a new API server.
func NewAPI(cfg *config.Config, n *node.Node, db *ledger.DB) *API {
	a := &API{
		cfg:  cfg,
		node: n,
		db:   db,
		mux:  http.NewServeMux(),
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
	Balance     float64 `json:"balance_hashes"`
	BytesServed int64   `json:"bytes_served_today"`
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

	writeJSON(w, LedgerResponse{
		Balance:     balance,
		BytesServed: bytesServed,
	})
}

// writeJSON writes a JSON response with correct headers.
func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
	}
}
