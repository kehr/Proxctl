# Artifacts

Release artifacts use hyphen-separated names:

<LatestVersion />

<ArtifactNames />

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
