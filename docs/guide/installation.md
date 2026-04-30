# Installation

## One-line install

```bash
curl -fsSL https://github.com/kehr/Proxctl/releases/latest/download/install.sh | sh
```

The installer detects the operating system and architecture, downloads the latest matching release asset, verifies `checksums.txt`, and installs `proxctl` to `/usr/local/bin`.

## Install a specific version

```bash
curl -fsSL https://github.com/kehr/Proxctl/releases/latest/download/install.sh | PROXCTL_VERSION=v0.2.0 sh
```

## Change install directory

```bash
curl -fsSL https://github.com/kehr/Proxctl/releases/latest/download/install.sh | PROXCTL_INSTALL_DIR="$HOME/.local/bin" sh
```

## Dry-run the installer

```bash
PROXCTL_DRY_RUN=1 ./scripts/install.sh
```

## Manual install

Download the archive for your platform from GitHub Releases:

```bash
curl -LO https://github.com/kehr/Proxctl/releases/download/v0.2.0/proxctl-v0.2.0-linux-amd64.tar.gz
tar -xzf proxctl-v0.2.0-linux-amd64.tar.gz
sudo install -m 0755 proxctl-v0.2.0-linux-amd64/proxctl /usr/local/bin/proxctl
proxctl version
```

## Supported release platforms

- `linux-amd64`
- `linux-arm64`
- `darwin-amd64`
- `darwin-arm64`
- `windows-amd64`

## Installer requirements

The installer expects common system tools:

- `uname`
- `sed`
- `grep`
- `tar`
- `curl` or `wget`
- `sha256sum` or `shasum`
- `sudo` when `/usr/local/bin` is not writable
