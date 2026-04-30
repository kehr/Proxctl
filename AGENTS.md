# Repository Guidelines

## Project Structure & Module Organization

Proxctl is a Go CLI for operating proxy services on VPS nodes. The entrypoint is `src/cmd/proxctl/main.go`; most behavior lives under `src/internal/`. Key packages include `cli`, `xray`, `client`, `state`, and `system`. Tests sit beside source files as `*_test.go`. Supporting files are in `configs/`, `deployments/systemd/`, `scripts/`, and `docs/`; generated command docs belong in `docs/commands/`.

## Build, Test, and Development Commands

- `make test`: runs `go test ./...`.
- `make vet`: runs `go vet ./...`.
- `make build`: builds `build/proxctl` with version metadata.
- `make all`: runs tests and builds the local binary.
- `make docs`: regenerates command docs with `build/proxctl` and verifies no uncommitted doc drift.
- `make docs-site`: installs docs dependencies, builds VitePress, and audits dependencies.
- `make verify`: runs pre-release checks: tests, vet, builds, archives, docs, shell syntax, and installer dry-run.
- `npm --prefix docs run dev`: starts the VitePress docs site locally after dependencies are installed.

## Coding Style & Naming Conventions

Use standard Go formatting: run `gofmt` on changed Go files and keep package names short, lowercase, and singular where practical. CLI commands and flags should follow existing Cobra patterns in `src/internal/cli`, with kebab-case flags such as `--config-file` and environment bindings through `PROXCTL_*`. Shell scripts in `scripts/` should remain POSIX-compatible where possible and pass `sh -n scripts/*.sh`.

## Testing Guidelines

Add focused unit tests beside the package being changed using the `*_test.go` convention. Prefer table-driven tests for config parsing, exporters, and command behavior. Run `make test` before submitting changes; run `make verify` for release, installer, docs, or cross-platform changes. If CLI behavior changes, run `make build` then `make docs`.

## Commit & Pull Request Guidelines

Git history uses Conventional Commit-style prefixes, for example `feat: add one-line installer`, `fix: harden operational config changes`, `docs: add VitePress documentation site`, and `chore: add make release workflow`. Keep commits scoped and imperative. Pull requests should include a concise summary, tests run, linked issues when relevant, and screenshots or local URLs for docs-site changes. Note operational impact for changes touching install, firewall, SSH, systemd, release, or backup/restore behavior.

## Security & Configuration Tips

Do not commit secrets, private keys, real VPS credentials, or generated archives from `build/` or `dist/`. Use `PROXCTL_DRY_RUN=1 ./scripts/install.sh` when validating installer behavior. Default configuration examples live in `configs/defaults.yaml`; avoid changing production paths or service assumptions without updating docs.
