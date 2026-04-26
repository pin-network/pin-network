# Contributing to PiN

First — thank you.  
PiN is built by its community. Every contribution, no matter how small, makes the network stronger.

Whether you’re a beginner, a hobbyist, or a seasoned engineer, you are welcome here.  
If you believe the web should be owned by its users, you’re already one of us.

---

# Ways to Contribute

You do **not** need to be a developer to contribute meaningfully to PiN.  
Choose the path that fits your skills and curiosity.

---

## 🧪 Hardware Testers (Raspberry Pi, SBCs, and more)

PiN is still early in its hardware journey — the official Pi image is not yet released.  
But **early testers accelerate development more than anything else**.

If you have a Raspberry Pi or similar device, you can help by:

- Building `meshd` from source  
- Running it on your hardware  
- Testing connectivity, uptime, and performance  
- Trying different network setups (WiFi, Ethernet, cellular, point‑to‑point)  
- Reporting what works — and what doesn’t  
- Sharing logs, configs, and results  

**Important:**  
Hardware testing is experimental. Things will break. That’s the point.  
Your findings help the entire community avoid duplicate work and move faster.

---

## 💻 Go Developers (meshd + pin-browser)

The core of PiN lives in Go:

- `meshd` — the node daemon  
- `pin-browser` — the local resolver and browser backend  

If you write Go, you can help with:

- DHT routing  
- NAT traversal  
- content‑addressed storage  
- ledger logic  
- HTTP APIs  
- performance tuning  
- cross‑platform builds  

Look for issues labeled:

- `daemon`  
- `browser`  
- `networking`  
- `ledger`  

---

## 🖥️ Rust + Web Developers (Tauri desktop app)

The PiN desktop app is built with:

- Rust (Tauri backend)  
- HTML/CSS/JS (frontend)  

You can help with:

- UI/UX  
- resource scheduler UI  
- Hash wallet UI  
- node status dashboard  
- `.pin` browsing shell  

Look for issues labeled:

- `tray-app`  
- `browser`  
- `ui`  

---

## ✍️ Documentation Writers

Clear documentation is as important as code.

You can help by improving:

- README  
- CONTRIBUTING  
- SPEC.md  
- tutorials  
- diagrams  
- examples  

Look for issues labeled `docs`.

---

## 🎨 Designers

PiN needs:

- UI/UX design  
- iconography  
- layout systems  
- accessibility improvements  
- visual diagrams  

Look for issues labeled `design`.

---

## 🌍 Translators

PiN aims to be accessible worldwide.  
If you speak another language, you can help translate:

- documentation  
- UI strings  
- onboarding materials  

Look for issues labeled `i18n`.

---

# Getting Started

## 1. Star and Watch the Repo  
You’ll be notified of new issues, releases, and discussions.

## 2. Join the Community  
Discord link will be posted in the README when public.  
Introduce yourself and share what you’re interested in.

## 3. Read the Technical Spec  
The full architecture lives in:  
👉 `docs/SPEC.md`

Understanding the system before writing code saves everyone time.

## 4. Find an Issue  
Look for:

- `good first issue` — great for beginners  
- `help wanted` — needs an owner  
- subsystem labels (`daemon`, `browser`, `tray-app`, etc.)

## 5. Claim It  
Comment on the issue so others know it’s taken.  
This prevents duplicate work.

---

# Development Setup

## Requirements

- Go 1.22+  
- Git  
- Node.js 20+ (for Tauri app)  
- Rust (for Tauri backend)  
- A Raspberry Pi (optional but helpful)

---

# Building the Subsystems

## meshd (node daemon)

```bash
git clone https://github.com/justj1979/pin-network.git
cd pin-network/src/meshd
go mod tidy
go build -o meshd .
./meshd --dev
