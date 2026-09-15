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
	-X github.com/cwensel/intrastate/internal/version.version=$(VERSION) \
	-X github.com/cwensel/intrastate/internal/version.commit=$(COMMIT) \
	-X github.com/cwensel/intrastate/internal/version.date=$(DATE)

.PHONY: docs docs-check all check fmt fmt-check vet lint graph-lint test test-ci vuln tools \
	tidy clean build install uninstall hooks snapshot release-check

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
# document's "The grammar" section is copied from, plus the routing table
# docs/cli-output-contract.md documents `dispositions` against (RDR 0024
# `0024:PH4`). They are linted here for the same reason the subject above
# is: a documented model that the loader or the analysis has since stopped
# accepting is worse than no example, and prose cannot catch that. All examples
# are expected to lint at exit 0 — the state machine carries
# `graph-coverage-closed-by-escape` advisories, which do not fail.
EXAMPLE_MODELS = \
	models/examples/markdown-review.toml \
	models/examples/pricing-decision-table.toml \
	models/examples/routing-decision-table.toml \
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

# Build the full release matrix locally, exactly as CI does, without
# publishing anything. The local mirror of the `snapshot` CI job — run it
# before cutting a tag to see the real artifacts.
snapshot:
	$(GO) run github.com/goreleaser/goreleaser/v2@latest release --snapshot --clean

# The distribution gate: assert a released artifact was stamped BY LDFLAGS.
# The `-X` paths in .goreleaser.yaml name the module path, and if they stop
# matching go.mod nothing fails to build — the binary silently falls back to
# internal/version's VCS stamps. A shipped binary that cannot name its
# release is the defect this catches.
#
# The oracle is the commit width, not the presence of text. Checking for
# "dev"/"none"/"unknown" does NOT work: inside a git repo the VCS fallback
# fills those in with real-looking values, so an unstamped build passes such
# a check (verified). The two paths are distinguishable because the fallback
# truncates to internal/version's vcsRevisionLen (12) and may append
# "-dirty", while goreleaser passes {{ .Commit }} — the full 40-char SHA.
# So: exactly 40 hex chars means ldflags won, which is the property under
# test. Keep this in sync with vcsRevisionLen if that constant changes.
#
# Runs against the freshly built dist/, so it must follow a goreleaser run;
# CI invokes it directly after one.
# The binary under test must be the one for THIS host: dist/ holds every
# target, and running a foreign one dies with "Exec format error", whose
# empty output reads as drift and reports the wrong cause. goreleaser names
# the per-target dirs `intrastate_<goos>_<goarch>...`, so ask the toolchain
# what host it is and match that prefix.
release-check:
	@goos=$$($(GO) env GOOS); goarch=$$($(GO) env GOARCH); \
	bin=$$(find dist -type f -name intrastate -path "*_$${goos}_$${goarch}*" 2>/dev/null | head -1); \
	if [ -z "$$bin" ]; then \
		echo "error: no $${goos}/$${goarch} binary under dist/ — run 'make snapshot' first" >&2; \
		exit 1; \
	fi; \
	if ! out=$$("$$bin" version 2>&1); then \
		echo "error: could not execute $$bin: $$out" >&2; \
		exit 1; \
	fi; \
	echo "$$out"; \
	commit=$$(printf '%s' "$$out" | sed -n 's/.*commit \([0-9a-f]*\).*/\1/p'); \
	if [ $${#commit} -ne 40 ]; then \
		echo "error: built binary was not stamped by ldflags (commit '$$commit' is $${#commit} chars, want 40)" >&2; \
		echo "it fell back to VCS stamps: the -X paths in .goreleaser.yaml no longer match the module path in go.mod" >&2; \
		exit 1; \
	fi; \
	echo "release-check: build identity stamped by ldflags ($${goos}/$${goarch})"

clean:
	rm -rf $(BIN_DIR) coverage.out dist

# One-shot: point git at the checked-in hooks in .githooks/
# (pre-commit: gofmt + vet; commit-msg: Conventional Commits).
hooks:
	git config core.hooksPath .githooks
	@echo "hooks enabled: core.hooksPath -> .githooks"
