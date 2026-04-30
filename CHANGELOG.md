# Changelog

## 0.1.0

- Initial Go implementation.
- Cobra/Viper CLI foundation with shell completion and Markdown docs generation.
- GitHub Actions CI and release workflows with multi-platform artifacts and SHA-256 checksums.
- Version-managed artifacts via the root `VERSION` file; release tags must match it.
- Xray provider with status, audit, health, doctor, config summary, adopt, backup, restore, init plan/apply, install plan/apply, rotate, and Reality target switch flows.
- Client exporters for generic URI, Shadowrocket, Surge, Mihomo/Stash, sing-box, v2rayN, and v2rayNG.
- SSH and firewall hardening plan/apply command surfaces.
- Production-oriented project layout with `src`, `docs`, `configs`, `scripts`, and `deployments`.
