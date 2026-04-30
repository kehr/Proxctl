#!/usr/bin/env sh
set -eu

REPO="${PROXCTL_REPO:-kehr/Proxctl}"
INSTALL_DIR="${PROXCTL_INSTALL_DIR:-/usr/local/bin}"
VERSION="${PROXCTL_VERSION:-latest}"
DRY_RUN="${PROXCTL_DRY_RUN:-0}"

log() {
  printf '%s\n' "$*" >&2
}

fail() {
  log "error: $*"
  exit 1
}

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || fail "required command not found: $1"
}

detect_os() {
  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  case "$os" in
    linux) printf 'linux' ;;
    darwin) printf 'darwin' ;;
    *) fail "unsupported OS: $os" ;;
  esac
}

detect_arch() {
  arch="$(uname -m)"
  case "$arch" in
    x86_64|amd64) printf 'amd64' ;;
    arm64|aarch64) printf 'arm64' ;;
    *) fail "unsupported architecture: $arch" ;;
  esac
}

download() {
  url="$1"
  out="$2"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "$out"
  elif command -v wget >/dev/null 2>&1; then
    wget -q "$url" -O "$out"
  else
    fail "curl or wget is required"
  fi
}

latest_version() {
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" |
      sed -n 's/.*"tag_name":[[:space:]]*"\([^"]*\)".*/\1/p' |
      head -n 1
  elif command -v wget >/dev/null 2>&1; then
    wget -qO- "https://api.github.com/repos/${REPO}/releases/latest" |
      sed -n 's/.*"tag_name":[[:space:]]*"\([^"]*\)".*/\1/p' |
      head -n 1
  else
    fail "curl or wget is required"
  fi
}

sha256_check() {
  checksums="$1"
  archive="$2"
  base="$(basename "$archive")"
  if command -v sha256sum >/dev/null 2>&1; then
    (cd "$(dirname "$archive")" && grep "  ${base}$" "$checksums" | sha256sum -c -)
  elif command -v shasum >/dev/null 2>&1; then
    (cd "$(dirname "$archive")" && grep "  ${base}$" "$checksums" | shasum -a 256 -c -)
  else
    fail "sha256sum or shasum is required"
  fi
}

main() {
  need_cmd uname
  need_cmd sed
  need_cmd grep
  need_cmd tar

  os="$(detect_os)"
  arch="$(detect_arch)"

  version="$VERSION"
  if [ "$version" = "latest" ]; then
    version="$(latest_version)"
  fi
  [ -n "$version" ] || fail "could not resolve latest release version"

  archive="proxctl-${version}-${os}-${arch}.tar.gz"
  base_url="https://github.com/${REPO}/releases/download/${version}"
  url="${base_url}/${archive}"
  checksums_url="${base_url}/checksums.txt"

  log "Proxctl installer"
  log "  repo:    ${REPO}"
  log "  version: ${version}"
  log "  target:  ${os}-${arch}"
  log "  install: ${INSTALL_DIR}/proxctl"

  if [ "$DRY_RUN" = "1" ]; then
    log "dry-run: would download ${url}"
    log "dry-run: would verify ${checksums_url}"
    log "dry-run: would install ${INSTALL_DIR}/proxctl"
    exit 0
  fi

  tmp="$(mktemp -d)"
  trap 'rm -rf "$tmp"' EXIT INT TERM

  download "$url" "$tmp/$archive"
  download "$checksums_url" "$tmp/checksums.txt"
  sha256_check "$tmp/checksums.txt" "$tmp/$archive"

  tar -xzf "$tmp/$archive" -C "$tmp"
  src="$tmp/proxctl-${version}-${os}-${arch}/proxctl"
  [ -f "$src" ] || fail "archive did not contain proxctl binary"

  if [ ! -d "$INSTALL_DIR" ]; then
    mkdir -p "$INSTALL_DIR" 2>/dev/null || sudo mkdir -p "$INSTALL_DIR"
  fi

  if [ -w "$INSTALL_DIR" ]; then
    install -m 0755 "$src" "$INSTALL_DIR/proxctl"
  else
    sudo install -m 0755 "$src" "$INSTALL_DIR/proxctl"
  fi

  "$INSTALL_DIR/proxctl" version
}

main "$@"
