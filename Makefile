BINARY := proxctl
PKG := ./src/cmd/proxctl
BUILD_DIR := build

.PHONY: all test build clean install linux-amd64

all: test build

test:
	go test ./...

build:
	mkdir -p $(BUILD_DIR)
	go build -trimpath -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY) $(PKG)

linux-amd64:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY)-linux-amd64 $(PKG)

install: build
	install -m 0755 $(BUILD_DIR)/$(BINARY) /usr/local/bin/$(BINARY)

clean:
	rm -rf $(BUILD_DIR)

