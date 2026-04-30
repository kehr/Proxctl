---
layout: home

hero:
  name: Proxctl
  text: Lightweight proxy operations for small VPS nodes.
  tagline: Deploy, adopt, audit, rotate, and export proxy client configs from one production-oriented CLI.
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: Install
      link: /guide/installation
    - theme: alt
      text: GitHub
      link: https://github.com/kehr/Proxctl

features:
  - title: Xray first, provider-ready
    details: Manage existing Xray nodes today while keeping the command model open for future proxy backends.
  - title: Safe operational workflows
    details: Review high-risk plans before applying credential rotation, target changes, SSH hardening, or firewall updates.
  - title: Client config export
    details: Generate profiles for Shadowrocket, Surge, Mihomo/Stash, sing-box, v2rayN, and v2rayNG.
  - title: Versioned releases
    details: Download platform-specific archives with consistent names and checksum verification.
---

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/kehr/Proxctl/main/scripts/install.sh | sh
```

## First checks

```bash
sudo proxctl adopt xray
sudo proxctl health
sudo proxctl audit --skip-updates
```

## Export a client profile

```bash
proxctl client export shadowrocket \
  --provider xray \
  --server <ip-or-domain> \
  --public-key <reality-public-key> \
  --name my-node
```
