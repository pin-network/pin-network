# PiN — Pi Integrated Network

> **Are you IN?**

PiN is a free, open-source, decentralized web hosting network built on Raspberry Pi computers and any device that wants to participate. The more nodes that join, the stronger the network becomes. Every node that serves traffic earns **Hashes** — a built-in incentive token that rewards contribution.

No central servers. No subscriptions. No gatekeepers. Just a network owned by everyone in it.

---

## What is PiN?

PiN turns idle devices into a self-sustaining web infrastructure. A Raspberry Pi sitting on your desk, a phone charging overnight, a laptop that sleeps at 11pm — all of these can become part of a global mesh that hosts websites, handles automation, and delivers content to anyone running the PiN browser or app.

It is designed for:

- **Static web hosting** — HTML, CSS, JavaScript sites with no backend requirements
- **Automation and webhooks** — lightweight HTTP triggers for IoT, home automation, and business workflows
- **File and asset delivery** — CDN-style distribution across geographically diverse nodes
- **Off-grid and rural connectivity** — works over WiFi, cellular, point-to-point, and LoRa links
- **Anyone, anywhere** — if a device can connect to a network, it can join PiN

---

## How it works

Every PiN node runs a small daemon called `meshd`. This daemon handles peer discovery, request routing, and logging of served traffic. Nodes find each other automatically using a distributed hash table (DHT), so there is no central directory that can be taken down.

When a user requests a page through the PiN browser, the smart router finds the closest available node that holds the content and serves it directly. If a node goes offline, content is re-routed automatically.

Traffic served earns Hashes. Hashes are recorded in a local proof-of-service ledger and periodically reconciled across the network. They can be spent on priority hosting, extra bandwidth allocation, or transferred to other users.

---

## Node tiers

| Tier | Device | Best for | Hash rate |
|------|--------|----------|-----------|
| 1 | Raspberry Pi, always-on SBC | Static hosting, webhooks, asset delivery | Base |
| 2 | Phone, tablet, laptop | Caching, relay, small APIs — opt-in schedule | 2× |
| 3 | Desktop PC, mini PC, NAS | Heavy compute, large storage, dynamic apps | 4× |

Any device can participate at any tier. You control exactly how much CPU, RAM, storage, and time you contribute.

---

## The Hash incentive

PiN uses a **proof-of-service** model. Hashes are not mined by solving arbitrary puzzles — they are earned by doing useful work: serving web requests, storing content, and staying online.

Your Hash earnings are proportional to:
- Bytes of traffic served
- Uptime percentage during your active window
- Storage pledged to the network
- Tier multiplier of your device

Hashes can be used to boost your own hosted content's priority, purchase additional bandwidth allocation, or gifted to other users. The ledger is distributed — no single authority controls it.

---

## The PiN browser

The PiN browser is the single install that does everything. It is a lightweight browser shell built on open web standards that adds:

- Native `.pin` domain resolution through the mesh DHT — no extension needed
- A built-in co-host daemon that activates when you opt in
- A Hash wallet visible in the toolbar
- A resource scheduler so you choose exactly when and how much your device contributes
- Works on Windows, macOS, Linux, Android, and iOS

Regular users just browse. Contributors flip the hosting switch. Same app, same install.

---

## Connectivity

PiN is designed to work wherever there is any kind of link:

- **WiFi** — primary mode for home and office nodes
- **Ethernet** — for always-on dedicated nodes
- **Cellular** — falls back gracefully; works with any existing data plan including MVNO and low-cost SIM providers
- **Point-to-point** — directional WiFi bridges for line-of-sight rural links
- **LoRa / Meshtastic** — low-bandwidth mesh for remote and off-grid environments (Phase 4)

PiN does not require a fast connection. A 5 Mbps upload is more than sufficient for Tier 1 hosting.

---

## Roadmap

### Phase 0 — Launch (March 14, 2025 — PI Day)
- [x] Public GitHub repository
- [x] Architecture specification
- [x] Community announcement
- [X] Discord server open

### Phase 1 — Core daemon (March → April 2025)
- [X] `meshd` daemon in Go — static file serving, peer discovery, DHT routing
- [X] Raspberry Pi OS image — single flash setup
- [X] CLI tools for node management
- [X] Local Hash ledger (SQLite)

### Phase 2 — Soft nodes and scheduler (May → June 2025)
- [X] Desktop tray app (Tauri) for Windows, macOS, Linux
- [X] Resource control UI — CPU cap, RAM limit, storage share, schedule
- [X] Battery and network rules for mobile devices
- [X] Hash ledger v1 sync across nodes

### Phase 3 — PiN browser alpha (July → September 2025)
- [ ] Browser shell with built-in mesh resolver
- [ ] `.pin` TLD support
- [ ] Hash wallet UI
- [ ] Android and iOS beta

### Phase 4 — Connectivity and rural (October → December 2025)
- [ ] LoRa and Meshtastic integration
- [ ] Point-to-point link support
- [ ] Offline mesh mode
- [ ] Cellular fallback optimisation

### Phase 5 — Business tier (2026)
- [ ] Subscription automation packages for corporate deployment
- [ ] Managed device fleet tooling
- [ ] Rural and off-grid deployment guides
- [ ] Enterprise SLA tooling

---

## Getting started

> Phase 3 is under active development. Star and watch this repo to be notified when the first release drops.

### For Raspberry Pi (coming Phase 1)

```bash
# Flash the PiN image to your SD card using Raspberry Pi Imager
# Select "PiN OS" from the custom image option
# Boot your Pi — meshd starts automatically
# Open the PiN browser on any device and navigate to your node's .pin address
```

### For desktop / laptop (coming Phase 2)

```bash
# Download the PiN tray app for your OS from the releases page
# Install and launch — the daemon runs in the background
# Set your resource schedule in the tray icon preferences
# Start earning Hashes immediately
```

### Build from source (developers)

```bash
git clone https://github.com/justj1979/pin-network.git
cd pin-network
# Requires Go 1.22+
cd src/meshd
go build -o meshd .
./meshd --init
```

---

## Tech stack

| Component | Technology | Why |
|-----------|-----------|-----|
| Node daemon | Go | Single static binary, cross-compiles to ARM, excellent networking primitives |
| Peer mesh | libp2p | Battle-tested DHT, NAT hole-punching, transport encryption |
| Web serving | nginx (embedded) | 2MB RAM at idle, gold standard performance |
| Local ledger | SQLite | Zero configuration, runs on every platform |
| Desktop app | Tauri (Rust + WebView) | 3–10MB binary vs 150MB Electron, cross-platform |
| RPI image | Raspberry Pi OS Lite + pi-gen | Official toolchain, single-flash experience |

---

## Architecture# PiN — Pi Integrated Network

**Are you IN?**

PiN is a free, open-source, decentralized web hosting network built on Raspberry Pi computers and any device that wants to participate. The more nodes that join, the stronger the network becomes. Every node that serves traffic earns Hashes — a built-in proof-of-service incentive token.

No central servers. No subscriptions. No gatekeepers. A network owned by everyone in it.

---

## What is PiN?

PiN turns idle devices into a self-sustaining web infrastructure. A Raspberry Pi sitting on your desk, a phone charging overnight, a laptop that sleeps at 11pm — all of these can become part of a global mesh that hosts websites, handles automation, and delivers content to anyone running the PiN browser or app.

PiN is also a parallel internet. It is not a VPN, not a dark web, and not a CDN. It is a completely separate network that:

- Runs on edge devices most people already own
- Grows organically as nodes cache and share content
- Pulls from the regular internet when it doesn't have something yet
- Eventually becomes self-sufficient as its content library grows
- Feels identical to the regular web from the user's perspective
- Cannot be accessed from the outside — it reaches out, never in

---

## How it works

Every PiN node runs a daemon called `meshd`. This daemon handles peer discovery via a Kademlia DHT, content-addressed file serving, request routing, and a local proof-of-service ledger. Nodes find each other automatically — no central directory that can be taken down.

When a user requests a `.pin` domain through the PiN browser, the resolver finds the closest available node holding that content and serves it directly. If a node goes offline, content is re-routed automatically.

Traffic served earns Hashes, recorded in a local SQLite ledger and reconciled across the network via gossip.

---

## Node tiers

| Tier | Device | Best for | Hash rate |
|------|--------|----------|-----------|
| 1 | Raspberry Pi, always-on SBC | Static hosting, webhooks, asset delivery | Base |
| 2 | Phone, tablet, laptop | Caching, relay, small APIs — opt-in schedule | 2× |
| 3 | Desktop PC, mini PC, NAS | Heavy compute, large storage, dynamic apps | 4× |

Any device can participate at any tier. You control exactly how much CPU, RAM, storage, and time you contribute.

---

## Current status

> **Active development — Phase 3 in progress**

### ✅ Phase 1 — Core daemon (complete)
- `meshd` daemon in Go — peer discovery, DHT routing, content serving
- libp2p + Kademlia DHT — NAT traversal, encrypted transport
- Content-addressed storage — SHA-256 CID addressing
- Local proof-of-service ledger — SQLite, epoch-based Hash calculation
- HTTP API — status, peers, ledger, content, domain, publish endpoints
- Two nodes connected and serving content — Windows and Raspberry Pi 4B
- Pure Go build — no CGO, runs on all platforms
- Permanent bootstrap node — Pi running as systemd service
- End-to-end content transfer verified over Tailscale mesh

### ✅ Phase 2 — Resource scheduler and node management (complete)
- Lightweight scheduler — 60-second tick, active/idle flag
- Schedule windows — HH:MM time ranges, overnight windows supported
- OS process priority — below-normal on Windows, nice 10 on Linux/Android/ARM
- Memory cap — Go runtime soft limit
- Concurrency limiting — max simultaneous requests from CPU% setting
- Bandwidth limiting — token bucket rate limiter
- Battery threshold — config flag, app layer handles detection
- WiFi-only mode — config flag, user defined
- Idle threshold — 20% default, user defined, app-configurable
- Config reload — 30-second file watcher, live updates without restart
- Platform support — Windows, Linux, Android, macOS, iOS, ARM, RISC-V, MIPS, AMD64

### 🔄 Phase 3 — PiN browser (in progress)
- `.pin` domain registration and resolution API ✅
- Site manifest format — multi-file `.pin` site publishing ✅
- `pin-browser` backend — local HTTP proxy server resolving `.pin` domains ✅
- Browser home page UI ✅
- Tauri desktop app — native window wrapper (next)
- Android and iOS (Phase 4)

### ⬜ Phase 4 — Tauri desktop app
- Native window on Windows, macOS, Linux
- Resource allocation UI — sliders for CPU, RAM, bandwidth, storage
- Schedule UI — active hours, idle mode
- Hash wallet UI — balance, earnings history, epoch countdown
- Node status dashboard — peers, uptime, bytes served
- Mobile — Android and iOS

### ⬜ Phase 5 — Meta-crawler and search
- Mesh content indexer
- Internet meta-crawler — seeds mesh with regular web content
- Unified search — mesh first, internet fallback
- Self-perpetuating mesh growth

### ⬜ Phase 6 — Extended connectivity
- LoRa / Meshtastic offline mesh overlay
- Point-to-point rural deployment
- Cellular optimization

### ⬜ Phase 7 — Commercial and ecosystem
- BooFunk Enterprises business tier
- Hash economy — spending, tipping, content monetization
- Developer SDK for `.pin` native apps
- PiN OS — dedicated Raspberry Pi image

---

## Build from source

### Requirements
- Go 1.22+
- Git

### meshd (node daemon)
```bash
git clone https://github.com/pin-network/pin-network.git
cd pin-network/src/meshd
go mod tidy
go build -o meshd .
./meshd --dev
```

### pin-browser (browser backend)
```bash
cd pin-network/src/pin-browser
go mod tidy
go build -o pin-browser .
./pin-browser
# Open http://localhost:7070 in your browser
```

---

## Tech stack

| Component | Technology | Why |
|-----------|-----------|-----|
| Node daemon | Go | Single binary, cross-compiles to ARM, excellent networking |
| Peer mesh | libp2p + Kademlia DHT | Battle-tested, NAT traversal, encrypted transport |
| Local ledger | SQLite (pure Go) | Zero configuration, runs on every platform, no CGO |
| Browser backend | Go HTTP server | No dependencies, works everywhere |
| Desktop app | Tauri (Rust + WebView) | 3–10MB binary, cross-platform native window |

---

## Architecture

See [SPEC.md](docs/SPEC.md) for the full technical specification including the DHT protocol, proof-of-service ledger format, Hash reward model, and node communication protocols.

---

## Contributing

PiN is built by its community. See [CONTRIBUTING.md](CONTRIBUTING.md) for how to get involved.

---

## License

MIT License — free for personal and commercial use.  
Copyright (c) 2025 BooFunk Enterprises and PiN Contributors.

---

*Built with ♥ for the Raspberry Pi community and everyone who believes the web should be owned by its users.*

See [SPEC.md](docs/SPEC.md) for the full technical specification including the DHT protocol design, proof-of-service ledger format, Hash reward model, and node communication protocols.

---

## Contributing

PiN is built by its community. Every skill level is welcome — from flashing Pi images to designing protocols to writing documentation.

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to get involved, claim work, and submit changes.

---

## License

MIT License — free for personal and commercial use. See [LICENSE](LICENSE) for details.

The PiN name, logo, and "Are you IN?" mark are project identifiers. The software itself is fully open and unrestricted.

---

## Community

- **GitHub Discussions** — questions, ideas, and architecture debate
- **Discord** — coming March 14, link posted here on launch day
- **Reddit** — [r/pinnetwork](https://reddit.com/r/pinnetwork) (launching PI Day)

---

*Built with ♥ for the Raspberry Pi community and everyone who believes the web should be owned by its users.*
