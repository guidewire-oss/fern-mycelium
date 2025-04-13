# Makefile for fern-mycelium with Dagger integration

APP_NAME := fern-mycelium
PKG := github.com/guidewire-oss/fern-mycelium
GOFILES := $(shell find . -name '*.go' -not -path "*/vendor/*")

.DEFAULT_GOAL := help

## ---------- DEV TASKS ---------- ##

.PHONY: tidy
 tidy: ## Run go mod tidy
	@echo ✅ Tidying go.mod...
	go mod tidy

.PHONY: fmt
fmt: ## Run go fmt
	@echo ✅ Formatting source...
	gofmt -s -w $(GOFILES)

.PHONY: lint
lint: ## Run Dagger lint
	dagger call lint --src .

.PHONY: generate
generate: ## Run go generate
	@echo ✅ Running code generation...
	go generate ./...

.PHONY: build
build: ## Build the Go binary
	dagger call build --src .

.PHONY: test
test: ## Run unit tests via Dagger
	dagger call test --src .

.PHONY: acceptance
acceptance: ## Run acceptance tests via Dagger
	dagger call acceptance --src .

.PHONY: scan
scan: ## Run Trivy filesystem scan on container
	dagger call scan --src .

.PHONY: scorecard
scorecard: ## Run OpenSSF Scorecard
	dagger call check-open-ssf --repo $(PKG)

.PHONY: publish
publish: ## Publish container to ttl.sh
	dagger call publish --src .

.PHONY: pipeline
pipeline: tidy fmt lint test build scan publish ## Run full pipeline with Dagger
## pipeline: tidy fmt generate lint test build scan scorecard publish ## Run full pipeline with Dagger
	@echo ✅ All stages complete.

## ---------- HELP ---------- ##

.PHONY: help
help: ## Display this help
	@echo "\nUsage: make <target>\n"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'
	@echo ""

