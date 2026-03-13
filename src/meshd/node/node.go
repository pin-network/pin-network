// Package node manages the libp2p host and Kademlia DHT for PiN.
package node

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"

	"meshd/config"
	"meshd/ledger"
)

// Node represents a PiN network node.
type Node struct {
	host host.Host
	dht  *dht.IpfsDHT
	cfg  *config.Config
	db   *ledger.DB
}

// identityFile holds the serialised node keypair.
type identityFile struct {
	PrivKey []byte `json:"priv_key"`
}

// Init creates a new node identity and data directory.
func Init(cfg *config.Config) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("finding home directory: %w", err)
	}

	dataDir := filepath.Join(home, ".pin")
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return fmt.Errorf("creating data directory: %w", err)
	}

	privKey, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return fmt.Errorf("generating keypair: %w", err)
	}

	privBytes, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		return fmt.Errorf("marshalling private key: %w", err)
	}

	identity := identityFile{PrivKey: privBytes}
	data, err := json.Marshal(identity)
	if err != nil {
		return fmt.Errorf("serialising identity: %w", err)
	}

	identityPath := filepath.Join(dataDir, "identity.json")
	if err := os.WriteFile(identityPath, data, 0600); err != nil {
		return fmt.Errorf("writing identity file: %w", err)
	}

	storePath := cfg.StorePath()
	if err := os.MkdirAll(storePath, 0755); err != nil {
		return fmt.Errorf("creating store directory: %w", err)
	}

	fmt.Printf("node identity created at %s\n", identityPath)
	return nil
}

// New creates a new PiN node with a libp2p host and DHT.
func New(ctx context.Context, cfg *config.Config, db *ledger.DB) (*Node, error) {
	privKey, err := loadOrCreateIdentity()
	if err != nil {
		return nil, fmt.Errorf("loading node identity: %w", err)
	}

	listenAddr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", cfg.Network.ListenPort)

	h, err := libp2p.New(
		libp2p.Identity(privKey),
		libp2p.ListenAddrStrings(listenAddr),
		libp2p.NATPortMap(),
		libp2p.EnableRelay(),
		libp2p.EnableNATService(),
	)
	if err != nil {
		return nil, fmt.Errorf("creating libp2p host: %w", err)
	}

	kadDHT, err := dht.New(ctx, h, dht.Mode(dht.ModeServer))
	if err != nil {
		h.Close()
		return nil, fmt.Errorf("creating DHT: %w", err)
	}

	return &Node{
		host: h,
		dht:  kadDHT,
		cfg:  cfg,
		db:   db,
	}, nil
}

// ID returns the node's peer ID as a string.
func (n *Node) ID() string {
	return n.host.ID().String()
}

// Addrs returns the node's listen addresses.
func (n *Node) Addrs() []string {
	addrs := make([]string, 0)
	for _, addr := range n.host.Addrs() {
		addrs = append(addrs, fmt.Sprintf("%s/p2p/%s", addr, n.host.ID()))
	}
	return addrs
}

// Bootstrap connects to bootstrap nodes and joins the DHT.
func (n *Node) Bootstrap(ctx context.Context) error {
	if err := n.dht.Bootstrap(ctx); err != nil {
		return fmt.Errorf("bootstrapping DHT: %w", err)
	}

	connected := 0
	for _, addrStr := range n.cfg.Network.BootstrapNodes {
		addr, err := multiaddr.NewMultiaddr(addrStr)
		if err != nil {
			continue
		}

		peerInfo, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			continue
		}

		n.host.Peerstore().AddAddrs(peerInfo.ID, peerInfo.Addrs, peerstore.PermanentAddrTTL)
		if err := n.host.Connect(ctx, *peerInfo); err != nil {
			continue
		}
		connected++
	}

	if connected == 0 {
		return fmt.Errorf("could not connect to any bootstrap nodes")
	}

	return nil
}

// Close shuts down the node gracefully.
func (n *Node) Close() error {
	if err := n.dht.Close(); err != nil {
		return err
	}
	return n.host.Close()
}

// loadOrCreateIdentity loads the node keypair from disk.
func loadOrCreateIdentity() (crypto.PrivKey, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("finding home directory: %w", err)
	}

	identityPath := filepath.Join(home, ".pin", "identity.json")

	data, err := os.ReadFile(identityPath)
	if err != nil {
		if os.IsNotExist(err) {
			return generateAndSaveIdentity(identityPath)
		}
		return nil, fmt.Errorf("reading identity file: %w", err)
	}

	var identity identityFile
	if err := json.Unmarshal(data, &identity); err != nil {
		return nil, fmt.Errorf("parsing identity file: %w", err)
	}

	privKey, err := crypto.UnmarshalPrivateKey(identity.PrivKey)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling private key: %w", err)
	}

	return privKey, nil
}

// generateAndSaveIdentity creates a new ed25519 keypair and saves it.
func generateAndSaveIdentity(path string) (crypto.PrivKey, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return nil, fmt.Errorf("creating identity directory: %w", err)
	}

	privKey, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generating keypair: %w", err)
	}

	privBytes, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		return nil, fmt.Errorf("marshalling private key: %w", err)
	}

	identity := identityFile{PrivKey: privBytes}
	data, err := json.Marshal(identity)
	if err != nil {
		return nil, fmt.Errorf("serialising identity: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return nil, fmt.Errorf("writing identity file: %w", err)
	}

	return privKey, nil
}
