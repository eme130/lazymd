BINARY      := lm
DESKTOP_BIN := lm-desktop
GO          := /usr/local/go/bin/go
GOFLAGS     ?=
DESKTOP_DIR := cmd/lm-desktop
FRONTEND    := $(DESKTOP_DIR)/frontend

.PHONY: build run test clean fmt vet lint \
        desktop desktop-dev desktop-frontend desktop-install \
        debug debug-desktop help

## --- TUI ---

build: ## Build TUI binary
	$(GO) build $(GOFLAGS) -o bin/$(BINARY) ./cmd/lm

run: ## Run TUI
	$(GO) run ./cmd/lm $(ARGS)

## --- Desktop (Wails) ---

desktop-install: ## Install desktop frontend deps
	cd $(FRONTEND) && npm install

desktop-frontend: ## Build desktop frontend (Vite)
	cd $(FRONTEND) && npm run build

desktop: desktop-frontend ## Build desktop binary
	cd $(DESKTOP_DIR) && $(GO) build $(GOFLAGS) -o ../../bin/$(DESKTOP_BIN) .

desktop-dev: ## Run desktop in Wails dev mode
	cd $(DESKTOP_DIR) && wails dev

## --- Quality ---

test: ## Run all tests
	$(GO) test ./...

fmt: ## Format Go code
	$(GO) fmt ./...

vet: ## Run go vet
	$(GO) vet ./...

lint: vet fmt ## Run all linters

## --- Debug (requires delve) ---

debug: ## Debug TUI with delve
	$(GO) build -gcflags='all=-N -l' -o bin/$(BINARY)-debug ./cmd/lm
	dlv exec bin/$(BINARY)-debug $(ARGS)

debug-desktop: desktop-frontend ## Debug desktop with delve
	cd $(DESKTOP_DIR) && $(GO) build -gcflags='all=-N -l' -o ../../bin/$(DESKTOP_BIN)-debug .
	dlv exec bin/$(DESKTOP_BIN)-debug $(ARGS)

## --- MCP ---

mcp: build ## Run as MCP server
	bin/$(BINARY) --mcp-server $(ARGS)

## --- Misc ---

clean: ## Remove build artifacts
	rm -rf bin/

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
