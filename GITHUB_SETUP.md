# GitHub Setup Guide
## Getting PiN live at github.com/justj1979/pin-network

---

## Step 1 — Create the repository

1. Go to https://github.com/new
2. Repository name: `pin-network`
3. Description: `PiN — Pi Integrated Network. A free, open-source decentralized web hosting mesh. Are you IN?`
4. Set to **Public**
5. Do NOT initialise with a README (we have our own)
6. Click **Create repository**

---

## Step 2 — Push all files

Open a terminal where you have the files from this package:

```bash
cd pin-network

git init
git add .
git commit -m "PI Day launch: PiN — Pi Integrated Network

Launching on March 14 (PI Day) with:
- Full architecture specification (docs/SPEC.md)
- README with project vision and roadmap
- Contributing guide
- Go project skeleton for meshd daemon
- GitHub Pages landing page
- MIT License

Are you IN?"

git branch -M main
git remote add origin https://github.com/justj1979/pin-network.git
git push -u origin main
```

---

## Step 3 — Enable GitHub Pages

1. Go to your repo → **Settings** → **Pages**
2. Source: **Deploy from a branch**
3. Branch: `main` / folder: `/docs`
4. Click **Save**
5. Your site will be live at: `https://justj1979.github.io/pin-network`

This takes 1–2 minutes to deploy. The URL will appear on the Pages settings page.

---

## Step 4 — Add repo metadata

On your repo main page, click the ⚙️ gear next to "About":

- Description: `PiN — Pi Integrated Network. A free, open-source decentralized web hosting mesh. Are you IN?`
- Website: `https://justj1979.github.io/pin-network`
- Topics: `raspberry-pi`, `decentralized`, `p2p`, `mesh-network`, `go`, `libp2p`, `open-source`, `hosting`, `iot`

---

## Step 5 — Add the logo to your repo

The logo SVG files are in `assets/`. GitHub READMEs display SVGs inline. The README already references them correctly.

---

## Step 6 — Enable GitHub Discussions

1. Go to **Settings** → **General**
2. Scroll to **Features**
3. Check **Discussions**
4. This gives you a community forum right on the repo

---

## Step 7 — Create a Discord server (optional but recommended)

1. Create a new Discord server named "PiN — Pi Integrated Network"
2. Channels to create:
   - `#announcements` (announcements only)
   - `#introductions`
   - `#general`
   - `#daemon-dev` (meshd development)
   - `#browser-dev`
   - `#hardware` (RPI setup, images)
   - `#hashes-ledger`
   - `#ideas`
3. Get the invite link and add it to the README

---

## Step 8 — Post to communities

See `docs/LAUNCH_POSTS.md` for ready-to-post text for:
- r/raspberry_pi
- r/selfhosted
- Hackaday.io
- Hacker News (Show HN)
- Hackster.io
- Twitter/X

Post to r/raspberry_pi first — it's the highest value audience for this project.

---

## Checklist

- [ ] Repo created at github.com/justj1979/pin-network
- [ ] All files pushed
- [ ] GitHub Pages enabled → justj1979.github.io/pin-network
- [ ] Repo description and topics set
- [ ] Discussions enabled
- [ ] Discord created (optional)
- [ ] r/raspberry_pi post submitted
- [ ] Hacker News Show HN submitted
- [ ] Hackaday.io project page created

---

*Happy PI Day. Let's build this.*
