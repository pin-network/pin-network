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

## Progress

### What is built and working

The core daemon (`meshd`) and the browser resolver (`pin-browser`) are both functional. The following components exist in `src/`:

| Component | Location | Status |
|-----------|----------|--------|
| `meshd` daemon | `src/meshd/` | ✅ Working |
| libp2p peer mesh + Kad DHT | `src/meshd/node/` | ✅ Working |
| Local HTTP API | `src/meshd/server/` | ✅ Working |
| Content-addressed store | `src/meshd/store/` | ✅ Working |
| `.pin` site manifest format | `src/meshd/manifest/` | ✅ Working |
| Domain registration + resolution | `src/meshd/server/` | ✅ Working |
| Site publish endpoint (`POST /api/v1/publish`) | `src/meshd/server/` | ✅ Working |
| Hash ledger (SQLite, epoch calculator) | `src/meshd/ledger/` | ✅ Working |
| Resource limits (CPU, RAM, bandwidth) | `src/meshd/limits/` | ✅ Working |
| Activity scheduler (active hours, idle ramp) | `src/meshd/scheduler/` | ✅ Working |
| YAML config with defaults | `src/meshd/config/` | ✅ Working |
| `pin-browser` shell + `.pin` resolver | `src/pin-browser/` | ✅ Working |
| Headless proxy mode | `src/pin-browser/browser/` | ✅ Working |

### API endpoints (meshd v0.1.0-dev)

`meshd` exposes a local REST API on port **4002** by default:

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/status` | Node status, peer count, uptime |
| `GET` | `/api/v1/peers` | Connected peers |
| `GET` | `/api/v1/ledger` | Hash balance and epoch history |
| `POST` | `/api/v1/content` | Store a blob, returns CID |
| `GET` | `/api/v1/content/{cid}` | Retrieve a blob by CID (with peer fallback) |
| `GET` | `/api/v1/domain` | List registered `.pin` domains |
| `GET/PUT/DELETE` | `/api/v1/domain/{name}` | Get, register, or delete a domain |
| `POST` | `/api/v1/publish` | Publish a full `.pin` site from JSON payload |

---

## Roadmap

### Phase 0 — Launch (March 14, 2025 — PI Day)
- [x] Public GitHub repository
- [x] Architecture specification
- [x] Community announcement
- [x] Discord server open

### Phase 1 — Core daemon (March → April 2025)
- [x] `meshd` daemon in Go — static file serving, peer discovery, DHT routing
- [x] CLI flags for node init, dev mode, and version
- [x] Local Hash ledger (SQLite) with epoch-based reward calculation
- [x] Content-addressed blob store (SHA-256 CID)
- [x] `.pin` site manifest format with CID integrity
- [x] Domain registration and resolution API
- [x] Site publish endpoint — full `.pin` sites uploadable via API
- [x] Resource limits (CPU, RAM, bandwidth caps)
- [x] Activity scheduler — active-hours windows and idle ramp

### Phase 2 — Soft nodes and scheduler (May → June 2025)
- [x] `pin-browser` shell — WebView wrapper with built-in `.pin` resolver
- [x] Headless proxy mode for displays-free nodes (RPi without monitor)
- [x] Battery and network rules for mobile devices
- [x] Hash ledger v1 sync across nodes

### Phase 3 — PiN browser alpha (July → September 2025)
- [ ] Full browser UI with Hash wallet toolbar
- [ ] `.pin` TLD support visible to end users
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

## Developer quick-start

### Prerequisites

| Tool | Version | Required for |
|------|---------|-------------|
| Go | 1.22+ | `meshd`, `pin-browser` |
| Git | any | cloning the repo |
| GCC / build tools | any | CGo (SQLite) |
| Raspberry Pi OS Lite | Bookworm | RPi hardware testing |

On Debian / Ubuntu / Raspberry Pi OS:

```bash
sudo apt update && sudo apt install -y golang git build-essential
```

On macOS (Homebrew):

```bash
brew install go git
```

On Windows, install [Go](https://go.dev/dl/) and [Git for Windows](https://gitforwindows.org/). CGo requires a C toolchain — install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) or enable WSL.

---

### Clone the repository

```bash
git clone https://github.com/pin-network/pin-network.git
cd pin-network
```

---

### Build and run meshd (node daemon)

```bash
cd src/meshd

# Download dependencies
go mod download

# Build the daemon binary
go build -o meshd .

# Initialise a new node identity (creates ~/.pin/config.yaml and a keypair)
./meshd --init

# Start the daemon in development mode (verbose logging)
./meshd --dev
```

The daemon starts two listeners:

- **Mesh port `4001`** — libp2p peer-to-peer connections
- **API port `4002`** — local REST API used by the browser and CLI tools

Check that it is running:

```bash
curl http://localhost:4002/api/v1/status
```

#### Optional: Windows PowerShell helper

A convenience script is included that sets `CGO_ENABLED=1` and kills any leftover port bindings before launching in dev mode:

```powershell
cd src/meshd
.\run.ps1
```

---

### Build and run pin-browser

The `pin-browser` requires a running `meshd` instance. Start `meshd` first, then in a second terminal:

```bash
cd src/pin-browser

# Download dependencies
go mod download

# Build the browser binary
go build -o pin-browser .

# Launch with a UI (requires a display)
./pin-browser

# Launch in headless proxy mode (no display needed — useful on RPi without monitor)
# Starts an HTTP proxy on port 7070 that routes .pin requests through the local node
./pin-browser --headless
```

Point the browser at your local node's API if it is not on the default address:

```bash
./pin-browser --api 127.0.0.1:4002
```

---

### Raspberry Pi setup

The following instructions get `meshd` running on a freshly flashed Raspberry Pi.

#### 1. Flash Raspberry Pi OS Lite

Use [Raspberry Pi Imager](https://www.raspberrypi.com/software/) to flash **Raspberry Pi OS Lite (64-bit, Bookworm)** to an SD card. In the Imager settings, configure your WiFi credentials and enable SSH before writing.

#### 2. Boot and connect

```bash
ssh pi@<your-pi-ip>
```

#### 3. Install Go and build tools

```bash
sudo apt update && sudo apt install -y build-essential

# Download and install Go (adjust version/arch as needed)
wget https://go.dev/dl/go1.22.0.linux-arm64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-arm64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version
```

#### 4. Clone and build meshd

```bash
git clone https://github.com/pin-network/pin-network.git
cd pin-network/src/meshd
go mod download
go build -o meshd .
```

#### 5. Initialise and start the node

```bash
./meshd --init
./meshd
```

#### 6. (Optional) Run as a systemd service so it starts on boot

```bash
sudo tee /etc/systemd/system/meshd.service > /dev/null <<EOF
[Unit]
Description=PiN meshd node daemon
After=network-online.target
Wants=network-online.target

[Service]
ExecStart=/home/pi/pin-network/src/meshd/meshd
Restart=on-failure
User=pi
WorkingDirectory=/home/pi

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable --now meshd
sudo systemctl status meshd
```

#### 7. Run pin-browser in headless proxy mode

On a Pi without a monitor, run the browser as a proxy that other devices on your network can use to reach `.pin` sites:

```bash
cd ~/pin-network/src/pin-browser
go build -o pin-browser .
./pin-browser --headless
# Now listening on :7070
```

---

### Configuration reference

`meshd` reads `~/.pin/config.yaml` on start. If the file does not exist, defaults are used. Run `./meshd --init` to create it.

```yaml
node:
  tier: 1                        # 1 = RPi/SBC, 2 = laptop/phone, 3 = desktop/NAS
  storage_path: ~/.pin/store     # where content blobs are stored
  storage_limit_gb: 10           # maximum storage pledged to the network
  bandwidth_limit_mbps: 5        # maximum upload bandwidth

schedule:
  always_on: true                # set false to use active_hours windows
  active_hours:
    - start: "08:00"
      end: "22:00"
  idle_threshold_pct: 20         # ramp up resources when below this utilisation %

network:
  listen_port: 4001              # libp2p peer port
  api_port: 4002                 # local REST API port
  enable_upnp: true
  enable_relay: true
  bootstrap_nodes:
    - /dns4/bootstrap.pin.network/tcp/4001/p2p/QmPlaceholderBootstrapID
  peer_apis: []                  # optional list of known peer API URLs for content fallback

limits:
  cpu_percent: 25                # maximum CPU the daemon may use
  ram_mb: 256                    # maximum RAM the daemon may use
  bandwidth_mbps: 5              # maximum bandwidth the daemon may use
  battery_min_percent: 30        # pause when battery drops below this level
  wifi_only: false               # only run when connected to WiFi
```

---

## Tech stack

| Component | Technology | Why |
|-----------|-----------|-----|
| Node daemon | Go | Single static binary, cross-compiles to ARM, excellent networking primitives |
| Peer mesh | libp2p + Kad DHT | Battle-tested DHT, NAT hole-punching, transport encryption |
| Content store | SHA-256 CID (custom) | Minimal, no external dependencies, content-integrity by default |
| Local ledger | SQLite (modernc, pure Go) | Zero configuration, runs on every platform without CGo issues |
| Browser shell | Go + WebView (`webview_go`) | Lightweight native window, zero Electron overhead |
| Config | YAML | Human-readable, easy to hand-edit on a Pi |

---

## Architecture

See [SPEC.md](docs/SPEC.md) for the full technical specification including the DHT protocol design, proof-of-service ledger format, Hash reward model, and node communication protocols.

---

## Contributing

PiN is built by its community. Every skill level is welcome — from running nodes on real hardware to designing protocols to writing documentation.

See [CONTRIBUTING.md](CONTRIBUTING.md) for how to get involved, claim work, and submit changes.

---

## License

MIT License — free for personal and commercial use. See [LICENSE](LICENSE) for details.

The PiN name, logo, and "Are you IN?" mark are project identifiers. The software itself is fully open and unrestricted.

---

## Community

- **GitHub Discussions** — questions, ideas, and architecture debate
- **Discord** — link in the repository description
- **Reddit** — [r/pinnetwork](https://reddit.com/r/pinnetwork)

---

*Built with ♥ for the Raspberry Pi community and everyone who believes the web should be owned by its users.*
