BINARY := proxctl
PKG := ./src/cmd/proxctl
BUILD_DIR := build
DIST_DIR := dist
VERSION ?= $(shell ./scripts/version.sh snapshot)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X github.com/kehr/proxctl/src/internal/cli.Version=$(VERSION) -X github.com/kehr/proxctl/src/internal/cli.Commit=$(COMMIT) -X github.com/kehr/proxctl/src/internal/cli.Date=$(DATE)

.PHONY: all test vet build clean install linux-amd64 snapshot dist docs docs-site verify release release-dry-run

all: test build

test:
	go test ./...

vet:
	go vet ./...

build:
	mkdir -p $(BUILD_DIR)
	go build -trimpath -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY) $(PKG)

linux-amd64:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY)-linux-amd64 $(PKG)

snapshot:
	VERSION="$(VERSION)" ./scripts/build-matrix.sh snapshot

dist:
	VERSION="$$(./scripts/version.sh release)" ./scripts/build-matrix.sh release

docs:
	$(BUILD_DIR)/$(BINARY) docs docs/commands
	git diff --exit-code docs/commands

docs-site:
	./scripts/docs-metadata.sh
	npm --prefix docs ci
	npm --prefix docs run build
	npm --prefix docs audit --omit=dev

verify: test vet build dist docs docs-site
	sh -n scripts/*.sh
	PROXCTL_DRY_RUN=1 ./scripts/install.sh

release:
	./scripts/release.sh

release-dry-run:
	DRY_RUN=1 ./scripts/release.sh

install: build
	install -m 0755 $(BUILD_DIR)/$(BINARY) /usr/local/bin/$(BINARY)

clean:
	rm -rf $(BUILD_DIR) $(DIST_DIR)
