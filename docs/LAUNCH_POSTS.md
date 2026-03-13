# PiN Launch Post Drafts
## PI Day — March 14, 2025

---

## 1. Reddit — r/raspberry_pi

**Title:**
PI Day project: PiN (Pi Integrated Network) — a decentralized web hosting network where your RPI earns "Hashes" for serving traffic. Free, open source, launching today.

**Body:**
Happy PI Day 🥧

I've been designing a project I want to build with the community, and today felt like the right day to put it out there.

**What is PiN?**

PiN (Pi Integrated Network) is a decentralized web hosting mesh where Raspberry Pis — and any other device that wants to join — serve web content to users and earn proof-of-service tokens called Hashes in return.

The idea is simple: your Pi is probably sitting there running 24/7 anyway. Why not let it serve a few web pages and earn something for it?

**How it works:**

- Flash the PiN image onto your Pi (or install the daemon on any device)
- Your node joins a peer mesh using a DHT — no central servers, no accounts
- You serve static sites, automation webhooks, and assets to other users
- Every byte you serve earns Hashes, recorded in a local proof-of-service ledger
- Hashes can boost your own hosted content or be shared with others

**It's not just for RPIs:**

Phones, laptops, and desktop PCs can all participate as opt-in "soft nodes." You set the hours (great for overnight when your device is idle), the CPU cap, the RAM limit, and whether it runs on battery. Very SETI@home in spirit.

**Why open source and free?**

Because the point is a web owned by its users, not a platform owned by a corporation. MIT licensed. No fees, ever.

**Where it stands:**

This is a PI Day launch of the concept, spec, and repository. Phase 1 (the actual `meshd` daemon and RPI image) is under active development. I'm looking for contributors — especially Go developers, networking folks, and people who want to stress-test this on real hardware.

**Repo:** https://github.com/justj1979/pin-network
**Spec:** https://github.com/justj1979/pin-network/blob/main/docs/SPEC.md

Are you IN?

---

## 2. Reddit — r/selfhosted

**Title:**
PiN — Pi Integrated Network: decentralized mesh hosting with proof-of-service Hash rewards. Open source, launching PI Day.

**Body:**
Self-hosters, this one's for you.

PiN is a decentralized web hosting network I'm building and open sourcing today (PI Day felt right). The pitch: instead of paying for a VPS to serve your static sites and automation endpoints, your self-hosted infrastructure *earns tokens* for the traffic it serves.

**The technical shape:**

- `meshd` — a Go daemon, single binary, runs on ARM (RPI) and x86
- libp2p + Kademlia DHT for peer discovery and routing — same stack as IPFS
- Content-addressed storage — sites are referenced by SHA-256 hash of their content
- Proof-of-service ledger in SQLite — no blockchain, no mining, just signed traffic claims
- Tauri desktop app for resource scheduling and the Hash wallet

**What it serves well:**

Static HTML/CSS/JS sites, automation webhooks, IoT endpoints, JSON APIs, file hosting. Not designed for heavy server-side computation — that's intentional.

**Node tiers:**

- Tier 1: RPI and always-on SBCs (base Hash rate)
- Tier 2: phones, laptops, opt-in nodes (2× Hash rate, schedule-controlled)
- Tier 3: dedicated PCs, NAS boxes (4× Hash rate)

**MIT licensed. No central authority. No fees.**

Repo: https://github.com/justj1979/pin-network

Looking for contributors — especially anyone with experience in libp2p, distributed ledgers, or embedded Linux. Come help build it.

---

## 3. Hackaday.io Project Page

**Tagline:** A decentralized web hosting mesh where your Raspberry Pi earns tokens for serving traffic.

**Description:**

PiN — Pi Integrated Network — is a free, open-source project that turns Raspberry Pis (and any other connected device) into nodes in a self-sustaining web hosting mesh.

Every node runs a lightweight Go daemon called `meshd` that handles peer discovery, content-addressed file serving, and proof-of-service logging. Nodes earn Hashes — a built-in reward token — proportional to the bandwidth they serve, their uptime, and their tier.

The browser client resolves `.pin` domains natively through the mesh DHT, so users never need to configure anything. The same app that browses the network can flip into hosting mode and start earning.

Built on libp2p (the same peer stack as IPFS), SQLite, nginx, and Tauri. Designed to run on a Raspberry Pi Zero 2W with 512MB RAM.

**Why this matters for the maker community:**

The web is increasingly centralized. A handful of cloud providers serve the majority of internet traffic. PiN is a bet that distributed hardware — the RPIs and old laptops and phones that the maker community already has — can form a viable alternative infrastructure.

Phase 1 is under development. Contributions welcome.

**Links:**
- GitHub: https://github.com/justj1979/pin-network
- Technical spec: https://github.com/justj1979/pin-network/blob/main/docs/SPEC.md

---

## 4. Hacker News — Show HN

**Title:** Show HN: PiN – Pi Integrated Network, a decentralized web hosting mesh with proof-of-service rewards

**Body:**

PiN is an open-source project I'm launching today (PI Day) that builds a decentralized web hosting network on top of Raspberry Pis and any other device that wants to participate.

The core idea: instead of hosting content on central servers, content is served from a peer mesh. Nodes earn proof-of-service tokens (Hashes) proportional to the traffic they serve and their uptime. No mining, no puzzles — just useful work.

**Technical stack:**
- Go daemon (meshd) using libp2p + Kademlia DHT for peer discovery and routing
- Content-addressed storage (SHA-256 CIDs, similar to IPFS but purpose-built for web serving)
- Proof-of-service ledger in SQLite with gossip-based sync between nodes
- Tauri app for desktop/mobile resource scheduling
- nginx embedded for actual HTTP serving

**Design decisions I'd love feedback on:**
- Using libp2p rather than rolling our own DHT — is the overhead acceptable on a Pi Zero?
- Gossip-based ledger sync vs a lightweight BFT consensus layer at scale
- Content moderation — the spec acknowledges this is an open question
- Whether 24-hour epochs for Hash calculation are the right granularity

**What exists today:** architecture spec, repository structure, and Go project skeleton. Phase 1 (working daemon + RPI image) is the immediate build target.

Repo: https://github.com/justj1979/pin-network
Spec: https://github.com/justj1979/pin-network/blob/main/docs/SPEC.md

Looking especially for contributors with libp2p experience and distributed systems background.

---

## 5. Hackster.io Project

**Title:** PiN — Pi Integrated Network: Turn Your Raspberry Pi Into a Rewarded Web Hosting Node

**Intro paragraph:**

What if your Raspberry Pi earned tokens just for being on? PiN (Pi Integrated Network) is an open-source project that connects RPIs — and any other always-on device — into a decentralized web hosting mesh. Every node that serves traffic earns Hashes, a proof-of-service reward built into the protocol.

**Things used in this project:**

Hardware: Raspberry Pi 4, Raspberry Pi 3B+, Raspberry Pi Zero 2W, USB storage drives (optional), PoE HAT (optional for clean always-on setup)

Software: PiN OS image (coming Phase 1), Go 1.22+, libp2p, nginx, SQLite, Tauri

**The story:**

[Describe the project vision, the PI Day launch, and invite Hackster readers to follow for updates as Phase 1 development progresses. Link to GitHub for the spec and to join as contributors.]

---

## 6. Twitter / X Thread

**Tweet 1:**
Launched PiN today on PI Day 🥧

Pi Integrated Network — a free, open-source decentralized web hosting mesh where your Raspberry Pi (or any device) earns tokens for serving traffic.

Are you IN? → github.com/justj1979/pin-network

**Tweet 2:**
How it works:
→ Flash PiN onto your Pi
→ Node joins a peer mesh (libp2p DHT, no central servers)  
→ You serve static sites, webhooks, automation
→ Every byte earns Hashes (proof-of-service, not proof-of-work)
→ Hashes boost your own hosted content

**Tweet 3:**
Not just for RPIs.

Phones, laptops, desktops can all opt in as soft nodes. Set your hours, CPU cap, RAM limit. Run heavy tasks overnight when your device is idle.

Very SETI@home in spirit 🔭

**Tweet 4:**
MIT licensed. No fees. No central authority. No ads.

The web should be owned by the people in it.

Spec: github.com/justj1979/pin-network/blob/main/docs/SPEC.md

Go devs, networking folks, RPI tinkerers — come build this with us.

#RaspberryPi #OpenSource #PIDay #PiN #decentralized
