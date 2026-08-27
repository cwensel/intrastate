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

.PHONY: docs docs-check all check fmt fmt-check vet lint graph-lint test test-ci vuln tools \
	tidy clean build install uninstall hooks

all: check

# Local mirror of the checks CI runs (govulncheck lives in its own CI
# job; run `make vuln` to mirror it locally).
#
# `build` is a prerequisite because `graph-lint` runs the BUILT command,
# the same one the CI graph-lint job runs — local parity means the same
# binary over the same model, not a second code path.
check: fmt-check vet lint build graph-lint docs-check test

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

# The graph-lint acceptance gate (RDR 0006), local half. It is a separate
# target from `lint` above, which is golangci-lint: the two gates check
# different things and neither is reusable for the other.
MODEL ?= models/rdr.toml

# The worked examples docs/model-authoring.md walks through: one per
# model class, plus the grammar-surface model every snippet in that
# document's "The grammar" section is copied from. They are linted here
# for the same reason the subject above is: a documented model that the
# loader or the analysis has since stopped accepting is worse than no
# example, and prose cannot catch that. All three are expected to lint at
# exit 0 — the two state machines carry `graph-coverage-closed-by-escape`
# advisories, which do not fail.
EXAMPLE_MODELS = \
	models/examples/pricing-decision-table.toml \
	models/examples/review-state-machine.toml \
	models/examples/release-grammar.toml

graph-lint: build
	$(BIN) lint --model $(MODEL) --as=json
	@for m in $(EXAMPLE_MODELS); do \
		echo "lint $$m"; \
		$(BIN) lint --model $$m --as=json || exit 1; \
	done

# Reference docs are GENERATED from the command tree (`internal/cli/docs.go`),
# never hand-edited: the binary already carries the flag grammar, the finding
# taxonomy, and the refusal codes, derived from the same constants the wire is
# emitted from. A hand-maintained copy is a second description that rots.
docs: build
	$(BIN) docs --dir .

# The staleness gate: regenerate into a scratch directory and compare with
# the committed copies. It deliberately does NOT diff the working tree
# against HEAD — that rejected correctly-regenerated-but-uncommitted files,
# so the documented `make docs` then `make check` sequence always failed and
# told you to stash the very files you had just correctly regenerated.
# Comparing against a fresh render answers the only question that matters:
# do the committed files match what this binary emits?
DOCS_FILES = docs/cli-reference.md llms.txt

docs-check: build
	@tmp=$$(mktemp -d) && trap 'rm -rf "$$tmp"' EXIT; \
	$(BIN) docs --dir "$$tmp" >/dev/null; \
	status=0; \
	for f in $(DOCS_FILES); do \
		if ! diff -q "$$f" "$$tmp/$$f" >/dev/null 2>&1; then \
			echo "error: $$f is stale. Run: make docs"; \
			diff -u "$$f" "$$tmp/$$f" | head -40; \
			status=1; \
		fi; \
	done; \
	if [ $$status -eq 0 ]; then echo "docs up to date"; fi; \
	exit $$status

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

# One-shot: point git at the checked-in hooks in .githooks/
# (pre-commit: gofmt + vet; commit-msg: Conventional Commits).
hooks:
	git config core.hooksPath .githooks
	@echo "hooks enabled: core.hooksPath -> .githooks"
