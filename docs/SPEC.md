# PiN Technical Specification

**Version:** 0.1.0-draft  
**Status:** Open for community review  
**Repository:** https://github.com/justj1979/pin-network

---

## Overview

This document defines the technical architecture of the PiN network — the peer discovery protocol, content addressing scheme, request routing model, proof-of-service ledger, and Hash reward mechanism. It is intended as the authoritative reference for contributors building any part of the stack.

---

## 1. Network Model

PiN is a fully decentralized peer-to-peer overlay network. There are no required central servers. Any node can join or leave at any time without disrupting the network.

The network is structured as a **Kademlia distributed hash table (DHT)** using the libp2p Go implementation. Every node is assigned a 256-bit NodeID derived from the SHA-256 hash of its public key, generated at first boot. NodeIDs are stable across restarts.

### 1.1 Bootstrap

New nodes need at least one existing peer address to join the DHT. PiN ships with a small set of well-known bootstrap node addresses embedded in the daemon binary. Bootstrap nodes are ordinary PiN nodes operated by community volunteers — they hold no privileged position in the network and can be replaced by any operator who publishes their address.

Bootstrap node addresses follow the libp2p multiaddr format:

```
/ip4/203.0.113.1/tcp/4001/p2p/QmBootstrapNodeID...
/dns4/bootstrap.pin.network/tcp/4001/p2p/QmBootstrapNodeID...
```

Any node can opt in to being a bootstrap node by setting `bootstrap: true` in its configuration file. Bootstrap nodes are listed publicly in the repository for community awareness.

### 1.2 Peer Discovery

After connecting to a bootstrap node, `meshd` performs a DHT walk to find nearby peers. Peers are stored in a local routing table using the standard Kademlia k-bucket structure (k=20). Peer records include:

- NodeID
- Multiaddr (IP, port, transport)
- Last seen timestamp
- Reported tier (1, 2, or 3)
- Reported available capacity (bandwidth, storage)

Peers are re-validated every 30 minutes. Nodes that fail to respond after three attempts are evicted from the routing table.

### 1.3 NAT Traversal

Most home nodes operate behind NAT. PiN uses the following strategies in order:

1. **UPnP/NAT-PMP** — automatic port mapping where the router supports it
2. **libp2p AutoNAT** — peers probe each other to determine reachability
3. **libp2p hole-punching** — coordinate simultaneous TCP/UDP connections through NAT
4. **Circuit relay** — if direct connection fails, traffic is relayed through a willing peer node

Relay nodes earn a small Hash bonus for relay traffic served.

---

## 2. Content Addressing

All content on PiN is addressed by the **SHA-256 hash of its contents**, not by a domain name or IP address. This is called a Content Identifier (CID).

A `.pin` site is a directory of content-addressed files described by a **manifest file**. The manifest is itself content-addressed. The manifest CID is what a `.pin` domain name resolves to.

### 2.1 Manifest Format

Manifests are JSON files with the following structure:

```json
{
  "version": 1,
  "name": "mysite",
  "description": "My PiN hosted site",
  "created": "2025-03-14T00:00:00Z",
  "updated": "2025-03-14T00:00:00Z",
  "entrypoint": "sha256:abc123...",
  "files": [
    {
      "path": "index.html",
      "cid": "sha256:abc123...",
      "size": 4096,
      "mime": "text/html"
    },
    {
      "path": "style.css",
      "cid": "sha256:def456...",
      "size": 2048,
      "mime": "text/css"
    }
  ]
}
```

### 2.2 Domain Resolution

`.pin` domain names are stored in the DHT as key-value pairs:

- Key: `SHA-256("pin-domain:" + domainname)`
- Value: signed record containing the current manifest CID and owner's public key

Domain records are signed by the owner's private key. Any node that receives a domain update verifies the signature before storing or forwarding it. Domain ownership is established by first publication — the first signed record for a name wins.

Domain records carry a TTL (default 48 hours) and must be republished by the owner to remain active.

---

## 3. Request Routing

When a PiN browser requests `mysite.pin/index.html`, the following sequence occurs:

1. Browser queries local DNS stub for `mysite.pin`
2. Stub queries the local `meshd` daemon on loopback
3. `meshd` looks up `mysite.pin` in the DHT to retrieve the manifest CID
4. `meshd` finds peers that hold the requested file CID using DHT lookup
5. `meshd` selects the optimal peer based on: latency, available bandwidth, tier, and geographic proximity
6. Content is fetched directly from the selected peer over an encrypted libp2p stream
7. Content is verified against its CID before being passed to the browser
8. Content is cached locally for configurable duration

If the primary peer fails mid-transfer, `meshd` automatically retries with the next candidate peer. The browser sees a seamless response.

### 3.1 Smart Router

The smart router selects nodes using a weighted scoring function:

```
score = (1/latency_ms) * bandwidth_weight * tier_weight * uptime_weight
```

Where:
- `bandwidth_weight` = reported available bandwidth / 10 Mbps (capped at 1.0)
- `tier_weight` = 0.8 (Tier 1) | 1.0 (Tier 2) | 1.2 (Tier 3)
- `uptime_weight` = rolling 7-day uptime percentage

The router maintains a short candidate list (top 5 peers for each CID) and probes them with a lightweight ping before committing to a transfer.

### 3.2 Request Weight Classification

The router classifies incoming requests by weight to determine which tier is appropriate:

| Weight class | Criteria | Preferred tier |
|---|---|---|
| Light | Static files < 1MB, HTML/CSS/JS | Tier 1 or 2 |
| Medium | Files 1–50MB, JSON APIs, small DB queries | Tier 2 or 3 |
| Heavy | Files > 50MB, compute tasks, large DB | Tier 3 only |

Heavy requests are only forwarded to Tier 3 nodes that have declared availability in the current time window.

---

## 4. The meshd Daemon

`meshd` is a single Go binary that runs as a background service. It is the only required software on a PiN node.

### 4.1 Responsibilities

- Maintain DHT peer table and respond to peer queries
- Serve content from local storage over libp2p streams
- Proxy browser requests to remote peers
- Log all served traffic to the local ledger
- Report node availability and capacity to the DHT
- Enforce resource limits set by the node operator

### 4.2 Configuration

`meshd` is configured via a YAML file at `~/.pin/config.yaml` (or `/etc/pin/config.yaml` for system installs):

```yaml
node:
  tier: 1                        # 1, 2, or 3
  storage_path: ~/.pin/store     # where content is stored
  storage_limit_gb: 10           # max storage to pledge
  bandwidth_limit_mbps: 5        # max upload bandwidth

schedule:
  always_on: true                # set false for scheduled operation
  active_hours:                  # ignored if always_on: true
    - start: "22:00"
      end: "07:00"
  heavy_tasks_only_during: "active_hours"

network:
  listen_port: 4001
  enable_upnp: true
  enable_relay: true
  bootstrap_nodes:
    - /dns4/bootstrap.pin.network/tcp/4001/p2p/QmBootstrap...

limits:
  cpu_percent: 25                # max CPU usage
  ram_mb: 256                    # max RAM usage
  battery_min_percent: 30        # pause hosting below this (mobile)
  wifi_only: false               # if true, pause on cellular
```

### 4.3 Storage Layout

```
~/.pin/
  config.yaml          node configuration
  identity             ed25519 keypair (node identity)
  ledger.db            SQLite Hash ledger
  store/               content-addressed file storage
    ab/
      cdef1234...      files stored by first 2 chars of CID
  domains/             domain records this node is authoritative for
  logs/                traffic logs (pre-ledger)
```

---

## 5. Proof-of-Service Ledger

The Hash reward system is based on verifiable claims of work done. No work is fabricated — every Hash earned corresponds to real traffic served to real peers.

### 5.1 Traffic Log Format

Every request served by `meshd` is logged locally in SQLite with the following fields:

```sql
CREATE TABLE traffic_log (
  id          INTEGER PRIMARY KEY,
  timestamp   INTEGER NOT NULL,       -- Unix timestamp
  requester   TEXT NOT NULL,          -- requesting NodeID
  content_cid TEXT NOT NULL,          -- CID served
  bytes       INTEGER NOT NULL,       -- bytes transferred
  duration_ms INTEGER NOT NULL,       -- transfer duration
  verified    BOOLEAN DEFAULT FALSE   -- spot-check verified
);
```

### 5.2 Hash Calculation

At the end of each epoch (24 hours), each node calculates its Hash earnings from the traffic log:

```
epoch_hashes = SUM(bytes_served) / BASE_RATE
             * uptime_factor
             * tier_multiplier
             * verification_factor
```

Where:
- `BASE_RATE` = 1,000,000 bytes per Hash (1 MB per Hash, subject to governance)
- `uptime_factor` = minutes_online / 1440 (fraction of day online)
- `tier_multiplier` = 1.0 (Tier 1) | 2.0 (Tier 2) | 4.0 (Tier 3)
- `verification_factor` = 0.0 to 1.0 based on spot-check results

### 5.3 Spot Verification

To prevent fraudulent claims, the network performs random spot checks. When node A claims to have served content CID X to node B, a third node C may query node B to confirm it received CID X from a node matching A's description within the claimed time window.

Nodes that fail spot checks have their `verification_factor` reduced for that epoch. Repeated failures result in temporary exclusion from Hash rewards.

### 5.4 Ledger Synchronisation

Epoch Hash totals are broadcast to neighboring peers as signed claim records. The network maintains a rolling consensus by gossiping claim records and cross-validating signatures. No global consensus is required — nodes maintain local views that converge over time through gossip.

```json
{
  "version": 1,
  "node_id": "QmNodeID...",
  "epoch": 20250314,
  "hashes_earned": 1847,
  "bytes_served": 1847000000,
  "uptime_minutes": 1380,
  "tier": 1,
  "signature": "ed25519:..."
}
```

---

## 6. Security Model

### 6.1 Transport Security

All peer-to-peer communication is encrypted using the **Noise protocol** via libp2p's default transport security. Every connection is mutually authenticated using the node's ed25519 identity key.

### 6.2 Content Integrity

Every piece of content is verified against its CID before serving or caching. A node that serves corrupted content will fail spot checks and lose Hash rewards. There is no way to tamper with content without changing its CID, which would break all references to it.

### 6.3 Domain Ownership

Domain records are signed by the owner's private key. A domain cannot be hijacked without access to the owner's private key. Key rotation is supported via a signed revocation + re-registration process.

### 6.4 Sybil Resistance

Hash rewards are tied to verifiable traffic, not to node count. Creating many fake nodes earns nothing because fake nodes serve no real traffic and will fail spot checks. Resource costs (storage, bandwidth, uptime) act as natural Sybil deterrents.

---

## 7. API Reference

`meshd` exposes a local HTTP API on `127.0.0.1:4002` for use by the browser and tray app.

```
GET  /api/v1/status              Node status and stats
GET  /api/v1/peers               Current peer table
GET  /api/v1/ledger              Hash balance and recent earnings
GET  /api/v1/content             List locally stored content
POST /api/v1/publish             Publish content to the network
POST /api/v1/domain/register     Register a .pin domain
PUT  /api/v1/config              Update node configuration
POST /api/v1/schedule            Update hosting schedule
```

---

## 8. Connectivity Modes

### 8.1 Standard (WiFi / Ethernet)
Default mode. No special configuration required.

### 8.2 Cellular
`meshd` detects cellular connections and applies conservative defaults: reduced bandwidth limit (1 Mbps), no large file serving, and increased cache TTLs to minimise data usage. Operators can override these defaults.

### 8.3 Point-to-Point
For rural deployments using directional WiFi bridges, PiN treats point-to-point links as standard network interfaces. No special configuration is required if the link presents as a normal network adapter.

### 8.4 LoRa / Meshtastic (Phase 4)
Low-bandwidth overlay for off-grid environments. LoRa nodes operate in a degraded mode: they can resolve domain names and serve cached content locally, but cannot serve large files to remote peers. Synchronisation of content and ledger records occurs when a LoRa node comes within range of a standard WiFi node.

---

## 9. Open Questions (Community Input Welcome)

The following design decisions are not finalised and are open for community discussion in GitHub Discussions:

- **Hash governance** — who can adjust `BASE_RATE` and tier multipliers, and how?
- **Ledger convergence** — is gossip sufficient at scale, or do we need a lightweight BFT consensus layer?
- **Domain expiry** — 48-hour TTL is conservative; should active sites auto-renew via the daemon?
- **Content moderation** — how does the network handle illegal content? Node operators bear responsibility for what they store and serve.
- **Mobile battery optimisation** — what is the right default battery threshold for Android vs iOS?

---

## Appendix A — Glossary

| Term | Definition |
|------|-----------|
| CID | Content Identifier — SHA-256 hash of a file's contents |
| DHT | Distributed Hash Table — decentralised key-value store |
| Hash | PiN's proof-of-service reward token |
| meshd | The PiN node daemon |
| NodeID | 256-bit identifier derived from a node's public key |
| Epoch | 24-hour period for Hash calculation |
| Manifest | JSON file describing a `.pin` site's content |
| Tier | Node classification (1=RPI, 2=soft, 3=power) |

---

*This specification is a living document. Submit issues or pull requests to propose changes.*
