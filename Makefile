# go-bidi — all automation is driven through this Makefile.
#
# The GitHub Actions workflow calls the same targets, so `make check` here
# is the full quality gate that CI enforces.

GO            ?= go
GOLANGCI_LINT ?= golangci-lint
GOFUMPT       ?= gofumpt

COVERDIR  ?= coverage
COVERFILE ?= $(COVERDIR)/coverage.out

.PHONY: help
help: ## Print this help
	@grep -E '^[a-zA-Z0-9_.-]+:.*?## ' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  %-14s %s\n", $$1, $$2}'

.PHONY: build
build: ## Build all packages
	$(GO) build ./...

.PHONY: test
test: ## Run the unit tests
	$(GO) test ./...

.PHONY: test-race
test-race: ## Run the unit tests with the race detector
	$(GO) test -race ./...

.PHONY: cover
cover: ## Run the tests and print the coverage report
	@mkdir -p $(COVERDIR)
	$(GO) test -coverprofile=$(COVERFILE) ./...
	$(GO) tool cover -func=$(COVERFILE)

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

.PHONY: fmt
fmt: ## Format the code with gofumpt
	$(GOFUMPT) -l -w .

.PHONY: fmt-check
fmt-check: ## Verify the code is gofumpt-formatted
	@out=$$($(GOFUMPT) -l .); \
	if [ -n "$$out" ]; then \
		echo "gofumpt: unformatted files:"; \
		echo "$$out"; \
		exit 1; \
	fi
	@echo "gofumpt: clean"

.PHONY: lint
lint: ## Run golangci-lint
	$(GOLANGCI_LINT) run ./...

.PHONY: lint-fix
lint-fix: ## Auto-fix lint issues
	$(GOLANGCI_LINT) run --fix ./...

.PHONY: check
check: fmt-check vet lint test ## Full quality gate used by CI

.PHONY: examples
examples: build ## Build the runnable examples
	$(GO) build ./examples/...

.PHONY: mod
mod: ## Tidy and verify the module files
	$(GO) mod tidy
	$(GO) mod verify

.PHONY: clean
clean: ## Remove build artifacts
	rm -rf $(COVERDIR)