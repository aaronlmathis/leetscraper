# SPDX-License-Identifier: GPL-3.0-or-later
# LeetScraper Makefile

APP_NAME      := leetscraper
CMD_DIR       := ./cmd/$(APP_NAME)
DIST_DIR      := ./dist
VERSION_DIR   := ./internal/version
VERSION_FILE  := $(VERSION_DIR)/version.go
TEMPLATE_FILE := $(VERSION_FILE).tpl
VERSION       := $(shell git describe --tags --always --dirty)
TEST_SCRIPT   := ./test/test_leetscraper.sh

GO            := go

.PHONY: all
all: help

## Generate version.go from template
$(VERSION_FILE): $(TEMPLATE_FILE)
	@mkdir -p $(VERSION_DIR)
	@sed "s/{{VERSION}}/$(VERSION)/" $(TEMPLATE_FILE) > $(VERSION_FILE)

## Build binary with embedded version
.PHONY: build
build: $(VERSION_FILE)
	@echo "🔨 Building $(APP_NAME) v$(VERSION)..."
	@mkdir -p $(DIST_DIR)
	$(GO) build -ldflags="-X github.com/aaronlmathis/leetscraper/internal/version.Version=$(VERSION)" -o $(DIST_DIR)/$(APP_NAME) $(CMD_DIR)



## Run integration tests
.PHONY: test
test: build
	@echo "🧪 Running integration tests..."
	@bash $(TEST_SCRIPT)

## Clean all build artifacts
.PHONY: clean
clean:
	@echo "🧹 Cleaning up..."
	@rm -rf $(DIST_DIR) $(VERSION_FILE)

## Cross-compile for Linux, macOS, Windows
.PHONY: release
release: $(VERSION_FILE)
	@echo "📦 Building cross-platform binaries for v$(VERSION)..."
	@mkdir -p $(DIST_DIR)
	GOOS=linux   GOARCH=amd64 $(GO) build -ldflags="-X github.com/aaronlmathis/leetscraper/internal/version.Version=$(VERSION)" -o $(DIST_DIR)/leetscraper-linux   $(CMD_DIR)
	GOOS=darwin  GOARCH=arm64 $(GO) build -ldflags="-X github.com/aaronlmathis/leetscraper/internal/version.Version=$(VERSION)" -o $(DIST_DIR)/leetscraper-darwin  $(CMD_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags="-X github.com/aaronlmathis/leetscraper/internal/version.Version=$(VERSION)" -o $(DIST_DIR)/leetscraper.exe     $(CMD_DIR)
	@cd $(DIST_DIR) && sha256sum * > checksums.txt
	@echo "✅ Binaries ready in $(DIST_DIR)/"

## Create distributable archives
.PHONY: package
package: release
	@echo "📦 Packaging release archives..."
	cd $(DIST_DIR) && \
	mv leetscraper-linux leetscraper && \
	tar -czf leetscraper-linux.tar.gz leetscraper && \
	mv leetscraper leetscraper-linux && \
	mv leetscraper-darwin leetscraper && \
	tar -czf leetscraper-darwin.tar.gz leetscraper && \
	mv leetscraper leetscraper-darwin && \
	zip -j leetscraper-windows.zip leetscraper.exe


## Show this help
.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@grep -E '^## ' $(MAKEFILE_LIST) | sed -e 's/^## //'
