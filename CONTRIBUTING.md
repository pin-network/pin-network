# Contributing to PiN

First — thank you. PiN is built by its community. Every contribution, no matter how small, makes the network stronger.

---

## Ways to contribute

You do not need to be a developer to contribute meaningfully to PiN.

**If you have a Raspberry Pi:** Flash the image when Phase 1 drops, run a node, and report what breaks. Real-world testing on real hardware is invaluable.

**If you write Go:** The `meshd` daemon is where the core work happens. See the issues labelled `daemon` for current priorities.

**If you write Rust or web frontend:** The Tauri desktop app and browser shell need builders. Issues labelled `browser` and `tray-app` are yours.

**If you do design:** The logo and color system are established. UI/UX design for the browser and tray app is wide open.

**If you write documentation:** Every component needs clear docs. Issues labelled `docs` welcome writers of all technical levels.

**If you know networking:** The DHT routing, NAT traversal, and LoRa integration need people who know these domains deeply.

**If you speak other languages:** Internationalisation of the browser and documentation reaches more of the world.

---

## Getting started

1. **Star and watch** the repository so you are notified of new issues and releases.

2. **Join the Discord** — link posted in the README on PI Day (March 14). Introduce yourself in `#introductions` and describe what you are interested in working on.

3. **Read the spec** — [docs/SPEC.md](docs/SPEC.md) is the authoritative technical reference. Understanding the architecture before writing code saves everyone time.

4. **Find an issue** — look for issues labelled `good first issue` if you are new to the project, or `help wanted` for open work that needs an owner.

5. **Claim it** — comment on the issue to let others know you are working on it. This prevents duplicate effort.

---

## Development setup

### Requirements

- Go 1.22 or later
- Git
- A Raspberry Pi (for hardware testing, not required for daemon development)
- Node.js 20+ and Rust (for the Tauri app, Phase 2+)

### Clone and build

```bash
git clone https://github.com/justj1979/pin-network.git
cd pin-network

# Build the daemon
cd src/meshd
go build -o meshd .

# Run tests
go test ./...

# Run a local development node
./meshd --dev --config dev-config.yaml
```

### Running two nodes locally

The best way to test peer discovery is to run two daemon instances with different ports:

```bash
# Terminal 1
./meshd --dev --port 4001 --api-port 4002 --data-dir /tmp/pin-node-a

# Terminal 2 — points to node A as bootstrap
./meshd --dev --port 4011 --api-port 4012 --data-dir /tmp/pin-node-b \
  --bootstrap /ip4/127.0.0.1/tcp/4001/p2p/$(cat /tmp/pin-node-a/identity.pub)
```

---

## Submitting changes

PiN uses a standard GitHub fork-and-pull-request workflow.

1. Fork the repository to your own GitHub account.
2. Create a branch named for your change: `feature/dht-routing`, `fix/nat-traversal`, `docs/spec-clarification`.
3. Make your changes. Write tests for new functionality.
4. Run the test suite: `go test ./...`
5. Commit with a clear message that explains *why*, not just *what*.
6. Open a pull request against the `main` branch. Fill in the PR template.

PRs are reviewed by maintainers within 72 hours where possible. Small, focused PRs are merged faster than large ones. If you are building something significant, open a Discussion first to align on design before writing code.

---

## Code standards

- **Go** — standard `gofmt` formatting, `golint` clean, meaningful variable names. No external dependencies without a Discussion first.
- **Comments** — public functions and types must have doc comments. Complex logic must have inline comments explaining intent.
- **Tests** — new features require tests. Bug fixes require a test that would have caught the bug.
- **No breaking changes** — the daemon API and DHT protocol are public interfaces. Changes that break compatibility require a major version bump and a deprecation period.

---

## Issue labels

| Label | Meaning |
|-------|---------|
| `good first issue` | Suitable for new contributors |
| `help wanted` | Needs an owner, any experience level |
| `daemon` | meshd core work |
| `browser` | PiN browser shell |
| `tray-app` | Desktop tray application |
| `ledger` | Hash ledger and proof-of-service |
| `networking` | DHT, NAT, routing |
| `docs` | Documentation |
| `hardware` | Raspberry Pi and embedded |
| `rural` | Off-grid and connectivity |
| `bug` | Something is broken |
| `discussion` | Design decision needed |

---

## Code of conduct

PiN is for everyone. Contributors are expected to be respectful, constructive, and welcoming to people of all backgrounds and experience levels.

Harassment, gatekeeping, and personal attacks will result in removal from the project. When in doubt, be kind.

---

## Questions?

Open a GitHub Discussion or ask in the Discord `#general` channel. No question is too basic. We were all beginners once.

---

*Are you IN? Let's build this together.*
