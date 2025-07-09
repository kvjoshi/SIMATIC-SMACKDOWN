BINARY_NAME=simatic_smackdown
VERSION=v1
GO=go

default: help

help: ## List Makefile targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

all: build

fmt: ## Format Go files
	gofumpt -w .

build: ## Build
	env $(if $(GOOS),GOOS=$(GOOS)) $(if $(GOARCH),GOARCH=$(GOARCH)) $(GO) build -o build/$(BINARY_NAME)

clean: ## Clean up build artifacts
	$(GO) clean
	rm ./build/$(BINARY_NAME)

run: build ## Run SIMATIC-SMACKDOWN (network scan mode)
	./build/$(BINARY_NAME)

run-targets: build ## Run with specific target IPs (example)
	./build/$(BINARY_NAME) 192.168.1.50 192.168.1.51 192.168.1.52

test-single: build ## Test with a single target IP
	./build/$(BINARY_NAME) 192.168.1.50
