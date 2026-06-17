# Local development entry points. CI invokes the same targets so local
# and remote runs share one source of truth.

GO            ?= go
GOLANGCI_LINT ?= $(shell $(GO) env GOPATH)/bin/golangci-lint

GOLANGCI_LINT_VERSION ?= v2.11.0

BIN_DIR     ?= bin
BIN          := $(BIN_DIR)/intrastate
INSTALL_DIR ?= $(HOME)/.local/bin
PKG          := ./cmd/intrastate

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -s -w \
	-X github.com/newcoinc/intrastate/internal/version.version=$(VERSION) \
	-X github.com/newcoinc/intrastate/internal/version.commit=$(COMMIT) \
	-X github.com/newcoinc/intrastate/internal/version.date=$(DATE)

.PHONY: all check fmt fmt-check vet lint test test-ci vuln tools tidy clean \
	build install uninstall hooks

all: check

# Local mirror of the checks CI runs (govulncheck lives in its own CI
# job; run `make vuln` to mirror it locally).
check: fmt-check vet lint test

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build -trimpath -ldflags '$(LDFLAGS)' -o $(BIN) $(PKG)

install: build
	@install -d $(INSTALL_DIR)
	install -m755 $(BIN) $(INSTALL_DIR)/intrastate
	@echo "installed: $(INSTALL_DIR)/intrastate"

uninstall:
	rm -f $(INSTALL_DIR)/intrastate
	@echo "uninstalled: $(INSTALL_DIR)/intrastate"

fmt:
	$(GO) fmt ./...

fmt-check:
	@out=$$(gofmt -l .); \
	if [ -n "$$out" ]; then \
		echo "gofmt needs to run on:"; echo "$$out"; exit 1; \
	fi

vet:
	$(GO) vet ./...

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run ./...

test:
	$(GO) test -race -covermode=atomic -coverprofile=coverage.out ./...

# Non-interactive CI test entry point. Mirrors `test`; split out so CI
# can layer JUnit/coverage artifacts here without disturbing the local
# target.
test-ci:
	$(GO) test -race -covermode=atomic -coverprofile=coverage.out ./...

vuln:
	$(GO) run golang.org/x/vuln/cmd/govulncheck@latest ./...

# Install pinned dev tools into $(GOPATH)/bin.
tools: $(GOLANGCI_LINT)

$(GOLANGCI_LINT):
	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

tidy:
	$(GO) mod tidy

clean:
	rm -rf $(BIN_DIR) coverage.out dist

# One-shot: install a pre-commit hook running gofmt + go vet.
hooks:
	@mkdir -p .githooks
	@printf '#!/bin/sh\nset -e\ngofmt -l . | (! grep .) || { echo "gofmt needed"; exit 1; }\ngo vet ./...\n' > .githooks/pre-commit
	@chmod +x .githooks/pre-commit
	git config core.hooksPath .githooks
	@echo "installed .githooks/pre-commit"
