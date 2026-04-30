# Artifacts

Release artifacts use hyphen-separated names:

```text
proxctl-v0.2.0-linux-amd64.tar.gz
proxctl-v0.2.0-linux-arm64.tar.gz
proxctl-v0.2.0-darwin-amd64.tar.gz
proxctl-v0.2.0-darwin-arm64.tar.gz
proxctl-v0.2.0-windows-amd64.zip
checksums.txt
```

Snapshot artifacts use:

```text
proxctl-snapshot-linux-amd64.tar.gz
```

## Verify checksums

```bash
shasum -a 256 -c checksums.txt
```

Linux systems may use:

```bash
sha256sum -c checksums.txt
```

## Build locally

```bash
make dist
```
