# Welcome to PiN — Pi Integrated Network 🎉

> **Are you IN?**

Welcome to the official PiN-Network wiki! Whether you just discovered this project or you've been watching it grow from the beginning, you're in the right place. This wiki is your guide to everything PiN — what it is, how it works, where it's headed, and how you can be part of it. Grab a coffee and dive in — this is going to be exciting! ☕🚀

---

## What is PiN?

PiN is a **free, open-source, decentralized web hosting network** built on Raspberry Pi computers and any device that wants to participate. What started as a fun Raspberry Pi day project has grown into something far more meaningful — a thoughtfully designed platform that, when complete, will serve as a genuine **private alternative to the internet**.

The vision is bold: **node-to-node (host-to-host) connections** that give users complete access to the regular internet, while keeping the regular internet from accessing the node itself. Think of it as a **one-way gate** — you can reach out, but outside traffic can't come in unless you invite it. This means node users can freely create **private or public networks** on their own terms.

- **Private networks** are truly secure — only nodes with the correct shared code even know another node exists. To everyone else, the node is completely invisible.
- **Public networks** allow a node to make itself visible to peers and openly share its cached data through a friendly web interface, like a community library anyone can browse.

The bigger the network grows, the more powerful and useful it becomes for everyone involved. 🌐

---

## Who is PiN For?

The use cases for PiN are as vast as the internet itself, limited only by the size of the available node network. Here are just a few exciting possibilities:

- 🏢 **Businesses that value internal security** — Deploy scalable private mesh networks for internal communication without exposing sensitive data to the public internet.
- 🌾 **Rural communities with limited internet access** — Set up simple node networks to use IoT devices and control equipment even when traditional broadband isn't an option.
- 📚 **Public access sites acting as library repositories** — Create secured, open repositories of information that function like digital community libraries.
- 🏠 **Home automation enthusiasts** — Build private smart-home networks that stay truly private.
- 🌍 **Anyone who believes the web should be owned by its users.**

The possibilities are genuinely endless. If you can imagine a use case for the internet, PiN can probably help you build a private or community version of it.

---

## How Does PiN Work?

Every PiN node runs a small daemon called **`meshd`**. This daemon handles peer discovery, request routing, and logging of served traffic. Nodes find each other automatically using a **Distributed Hash Table (DHT)**, so there is no central directory that can be taken down.

When a user requests a page through the PiN browser, the smart router finds the closest available node that holds the content and serves it directly. If a node goes offline, content is re-routed automatically. The network heals itself! 💪

### Node tiers

You control exactly how much CPU, RAM, storage, and time you contribute. Any device can participate:

| Tier | Device | Best for | Hash rate |
|------|--------|----------|-----------|
| 1 | Raspberry Pi, always-on SBC | Static hosting, webhooks, asset delivery | Base |
| 2 | Phone, tablet, laptop | Caching, relay, small APIs — opt-in schedule | 2× |
| 3 | Desktop PC, mini PC, NAS | Heavy compute, large storage, dynamic apps | 4× |

Even a simple Raspberry Pi sitting on your desk or a phone charging overnight can become part of a global mesh that hosts websites, handles automation, and delivers content. Every device counts!

---

## The PiN Browser

The **PiN browser** is the single install that does everything. It's a lightweight browser shell built on open web standards that adds:

- Native `.pin` domain resolution through the mesh DHT — no extension needed
- A built-in co-host daemon that activates when you opt in
- A Hash wallet visible in the toolbar
- A resource scheduler so you choose exactly when and how much your device contributes
- Works on **Windows, macOS, Linux, Android, and iOS**

As a **pseudo-browser**, it allows websites to be built and surfed just like the regular internet. Apps and websites can be built and deployed on PiN — and that is fully encouraged! 🎨🛠️

Regular users just browse. Contributors flip the hosting switch. Same app, same install. Simple.

---

## The Hash Incentive — Earn by Doing Good Work 💰

PiN uses a **proof-of-service** model. Hashes are not mined by solving arbitrary puzzles — they are earned by doing **useful work**: serving web requests, storing content, and staying online.

Your Hash earnings are proportional to:
- Bytes of traffic served
- Uptime percentage during your active window
- Storage pledged to the network
- Tier multiplier of your device

Hashes can be used to boost your hosted content's priority, purchase additional bandwidth allocation, or gifted to other users. The ledger is distributed — no single authority controls it.

### 🪙 A Fringe Benefit: Hash Mining & Crypto

Here's a fun bonus: the secure encryption method underlying PiN is built using **hash mining**. This can be leveraged to feed data into a crypto mining operation without expensive or complex setups. You're already running the network — why not let it work even harder for you?

---

## Open Source & Accessible to Everyone

PiN is **open source and freely available** for anyone to use. It's light enough to run on **16-bit devices** (with appropriately limited capability) using minimal internal memory assigned by the user. You decide how much space and compute you contribute — PiN works within whatever you give it.

There are no central servers, no subscriptions, and no gatekeepers. Just a network owned by everyone in it. The community builds it together, and the community benefits from it together. That's the PiN spirit. 🤝

---

## Connectivity — Works Wherever You Are

PiN is designed to work wherever there is any kind of network link:

- **WiFi** — primary mode for home and office nodes
- **Ethernet** — for always-on dedicated nodes
- **Cellular** — falls back gracefully; works with any existing data plan
- **Point-to-point** — directional WiFi bridges for line-of-sight rural links
- **LoRa / Meshtastic** — low-bandwidth mesh for remote and off-grid environments *(Phase 4)*

PiN does not require a fast connection. A 5 Mbps upload is more than sufficient for Tier 1 hosting. Got a slow rural connection? PiN was built with you in mind. 🌄

---

## Tech Stack

| Component | Technology | Why |
|-----------|-----------|-----|
| Node daemon | Go | Single static binary, cross-compiles to ARM, excellent networking primitives |
| Peer mesh | libp2p | Battle-tested DHT, NAT hole-punching, transport encryption |
| Web serving | nginx (embedded) | 2MB RAM at idle, gold standard performance |
| Local ledger | SQLite | Zero configuration, runs on every platform |
| Desktop app | Tauri (Rust + WebView) | 3–10MB binary vs 150MB Electron, cross-platform |
| RPI image | Raspberry Pi OS Lite + pi-gen | Official toolchain, single-flash experience |

---

## Roadmap

The project has already made incredible progress and the momentum is only building! Here's where things stand:

### ✅ Phase 0 — Launch (PI Day, March 14, 2025)
- Public GitHub repository live
- Architecture specification published
- Community announcement made
- Discord server open

### ✅ Phase 1 — Core Daemon
- `meshd` daemon in Go — static file serving, peer discovery, DHT routing
- Raspberry Pi OS image — single flash setup
- CLI tools for node management
- Local Hash ledger (SQLite)

### ✅ Phase 2 — Soft Nodes & Scheduler
- Desktop tray app (Tauri) for Windows, macOS, Linux
- Resource control UI — CPU cap, RAM limit, storage share, schedule
- Battery and network rules for mobile devices
- Hash ledger v1 sync across nodes

### 🚧 Phase 3 — PiN Browser Alpha *(Coming July–September 2025)*
- Browser shell with built-in mesh resolver
- `.pin` TLD support
- Hash wallet UI
- Android and iOS beta

### 🔭 Phase 4 — Connectivity & Rural *(October–December 2025)*
- LoRa and Meshtastic integration
- Point-to-point link support
- Offline mesh mode
- Cellular fallback optimisation

### 🏢 Phase 5 — Business Tier *(2026)*
- Subscription automation packages for corporate deployment
- Managed device fleet tooling
- Rural and off-grid deployment guides
- Enterprise SLA tooling

---

## Getting Started

> Phase 1 is under active development. **Star and Watch** this repo to be notified the moment the first release drops!

### For Raspberry Pi *(coming Phase 1)*

```bash
# Flash the PiN image to your SD card using Raspberry Pi Imager
# Select "PiN OS" from the custom image option
# Boot your Pi — meshd starts automatically
# Open the PiN browser on any device and navigate to your node's .pin address
```

### For Desktop / Laptop *(coming Phase 2)*

```bash
# Download the PiN tray app for your OS from the releases page
# Install and launch — the daemon runs in the background
# Set your resource schedule in the tray icon preferences
# Start earning Hashes immediately
```

### Build from Source *(for developers)*

```bash
git clone https://github.com/pin-network/pin-network.git
cd pin-network
# Requires Go 1.22+
cd src/meshd
go build -o meshd .
./meshd --init
```

---

## Architecture & Specifications

Want to go deep? See [SPEC.md](../docs/SPEC.md) for the full technical specification including the DHT protocol design, proof-of-service ledger format, Hash reward model, and node communication protocols.

---

## Contributing 🙌

PiN is built by its community, and **every skill level is welcome** — from flashing Pi images to designing protocols to writing documentation. Your contribution matters, no matter how small.

See [CONTRIBUTING.md](../CONTRIBUTING.md) for how to get involved, claim work, and submit changes. We'd love to have you!

---

## Community & Staying Connected

Come say hello and join the conversation!

- 💬 **GitHub Discussions** — questions, ideas, and architecture debate
- 🎮 **Discord** — real-time community chat
- 🟠 **Reddit** — [r/pinnetwork](https://reddit.com/r/pinnetwork)

---

## License

**MIT License** — free for personal and commercial use. See [LICENSE](../LICENSE) for details.

The PiN name, logo, and "Are you IN?" mark are project identifiers. The software itself is fully open and unrestricted.

---

*Built with ♥ for the Raspberry Pi community and everyone who believes the web should be owned by its users.*

*The project has grown. The vision is clear. The network is yours — are you IN?* 🍓🌐
