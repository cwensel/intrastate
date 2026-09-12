Model: claude-sonnet-5

# A5 verification spike — enumeration seams for the two non-mechanical owed rows

Scope: RDR `0029:A5`. Verifies whether a tier assignment can be backed by an
enumeration seam in every package that owns a tiered vocabulary, without
breaking a caller or crossing a package boundary the layering forbids.
Focuses on the two rows C4 flags as not mechanical: the CLIError `code`
vocabulary and the `flow next` unknown-`reason` union.

All commands run from `/Users/cwensel/sandbox/newcoinc/intrastate`.

---

## 1. Grounding — existing seams and current shape

Confirmed the four established seams:

```
internal/resolve/resolve.go:67:func RefusalKinds() []RefusalKind
internal/table/category.go:90:func Categories() []Category
internal/accessor/model.go:309:func Verdicts() []Verdict { return slices.Clone(verdicts) }
internal/graphlint/taxonomy.go:112:func Reasons() []string { return slices.Clone(reasons) }
```

Also found `internal/resolve/guard.go::Reasons` (the `graph-unprovable-coverage`
row's seam — already carried, matches C4's census).

### clierr layering (confirmed)

`internal/cli/clierr/clierr.go` doc comment: "Package clierr is the leaf
home of the CLI's structured-error type. It lives in its own package so
other internal packages (config, ...) can construct CLIErrors without
importing internal/cli and forming an import cycle." Imports are stdlib
only (`encoding/json`, `errors`, `fmt`, `io`, `strconv`, `strings`).

Reverse check — `internal/graphlint` already imports `internal/cli/clierr`
(`analysis.go`, `coverage.go`, `engine.go`, `groups.go`) to construct
`clierr.Finding` / `clierr.CLIError` values. `internal/guard` and
`internal/accessor` do not import `clierr` directly. `clierr` itself does
not import `internal/cli`, `internal/graphlint`, `internal/guard`, or
`internal/accessor`.

This fixes the dependency direction: `clierr` must stay upstream of
`graphlint` (and, if either package ever needs to construct a `CLIError`,
of `guard`/`accessor` too). `clierr` importing any of them back is a cycle.

### CLIError `code` raise-site shape

`Code*` / `ValidationCode` constants are minted independently, package by
package, with no shared type tying them to `clierr.CLIError.Code` (a bare
`string` field):

```
internal/graphlint/taxonomy.go:25-43   10 CodeXxx = "..." constants (untyped)
internal/guard/lint.go:21-27            7 CodeXxx Code = "..." constants (guard.Code type)
internal/accessor/model.go:100-117      7 CodeXxx ValidationCode = "..." constants (accessor.ValidationCode type)
```

Raise sites (`Code:` field assignments feeding `clierr.CLIError`/
`clierr.Finding` construction, non-test files): `internal/graphlint` 20,
`internal/guard` 9, `internal/accessor` 7, `internal/cli` 35 — 71 total
`Code:` assignment sites tree-wide (non-test), consistent with the RDR's
"~61 raise sites... ~10 inline literals" order of magnitude; several of
`internal/cli`'s 35 are bare string literals (`"command-error"`,
`"docs-write-failed"`, `"flag-mutually-exclusive"`, etc.) rather than named
constants, confirming the "part constant part inline literal" shape C4
describes.

### `flow next` reason union

`internal/cli/flow_next.go:91-95`:

```go
const (
	reasonAbsent       = string(resolve.ReasonAbsent)
	reasonUncomparable = string(resolve.ReasonUncomparable)
	reasonNotEvaluated = "not-evaluated"
)
```

Confirmed: re-spells `resolve.Reason*` (`internal/resolve/guard.go::ReasonAbsent`,
`::ReasonUncomparable`) plus exactly one local token
(`reasonNotEvaluated`). `internal/cli` already imports `internal/resolve`
directly (`flow_next.go` import block) — no layering boundary between them.

---

## 2. Spike — writing the two accessors

### 2a. `flow next` unknown-`reason` seam (mechanical case)

Added to `internal/cli/flow_next.go`, immediately after the existing const
block, matching `graphlint::Reasons`' shape:

```go
var flowNextUnknownReasons = []string{reasonAbsent, reasonUncomparable, reasonNotEvaluated}

func UnknownReasons() []string { return slices.Clone(flowNextUnknownReasons) }
```

(`slices` was already imported in this file.)

```
$ go build ./...
(exit 0, no output)
```

Result: clean. No import needed beyond what the file already carries,
because `internal/cli` already sits downstream of `internal/resolve`. This
confirms A5's framing — this row is mechanical, not one of the two hard
cases.

### 2b. CLIError `code` seam — attempt 1: natural registry shape

Added to `internal/cli/clierr/clierr.go`, importing the three producing
packages to enumerate their real constants:

```go
import (
	...
	"github.com/cwensel/intrastate/internal/accessor"
	"github.com/cwensel/intrastate/internal/graphlint"
	"github.com/cwensel/intrastate/internal/guard"
)

var cliErrorCodes = []string{
	graphlint.CodeDanglingEdge,
	string(guard.CodeUnprovableCoverage),
	string(accessor.CodeMissingAccessor),
}

func Codes() []string { ... }
```

```
$ go build ./...
package github.com/cwensel/intrastate/cmd/intrastate
	imports github.com/cwensel/intrastate/internal/cli from main.go
	imports github.com/cwensel/intrastate/internal/cli/clierr from docs.go
	imports github.com/cwensel/intrastate/internal/graphlint from clierr.go
	imports github.com/cwensel/intrastate/internal/cli/clierr from analysis.go: import cycle not allowed
(exit 1)

$ go vet ./...
(same cycle, exit 1)
```

**Import cycle confirmed, exactly as A5 predicted**:
`clierr` → `graphlint` → `clierr` (via `graphlint/analysis.go` et al.,
which already import `clierr`). This is not a hypothetical — `go build`
fails hard the moment `clierr` imports any of the three producing
packages.

### 2c. CLIError `code` seam — attempt 2: leaf-preserving registry

Reverted the import and instead hand-copied a subset of the literal code
strings directly into `clierr.go`, with no import of the producing
packages:

```go
var cliErrorCodes = []string{
	"graph-dangling-edge",        // mirrors graphlint.CodeDanglingEdge
	"graph-unprovable-coverage",  // mirrors guard.CodeUnprovableCoverage
	"missing_accessor",           // mirrors accessor.CodeMissingAccessor
}

func Codes() []string { ... }
```

```
$ go build ./...
(exit 0, no output)

$ go vet ./...
(exit 0, no output)
```

Result: builds clean, no cycle, `clierr` stays a leaf. **But this is not an
enumeration seam in the sense C4 requires.** It is a second, independently
maintained copy of a subset of the code literals with no compiler or
runtime tie back to the actual raise sites. Nothing prevents a raise site
from minting a new code, or changing an existing one, without this list
being updated — the four established seams (`RefusalKinds`, `Categories`,
`Verdicts`, `Reasons`) all return the literal backing slice the emit sites
actually use; this shape returns a stand-in slice the emit sites never
touch. It would not support C2's by-value tier assertion (`append-only`
membership/uniqueness checking) against the real vocabulary — only against
itself.

---

## 3. Revert

```
$ git status --porcelain   (before revert of my two files)
 M internal/cli/clierr/clierr.go
 M internal/cli/flow_next.go
 (plus unrelated foreign-session changes already present in the tree before this spike started — see below)

$ git checkout -- internal/cli/flow_next.go internal/cli/clierr/clierr.go

$ git diff --quiet internal/cli/clierr/clierr.go && echo clean
clean
$ git diff --quiet internal/cli/flow_next.go && echo clean
clean
```

Both spiked files confirmed byte-identical to HEAD after revert (`git diff`
empty on both). No commit was made.

Note: the working tree carried modifications NOT made by this spike, both
before and after — `docs/rdr/0029-version-promise-on-machine-readable-output.md`,
`internal/cli/respond/respond.go`, `internal/graphlint/taxonomy.go`
(carrying an inline comment `// SPIKE: undisclosed promotion for A7
red-case demo`, i.e. a *different*, concurrent A7 spike), plus untracked
`internal/cli/a7_snapshot_spike_test.go` and `internal/cli/testdata/`.
These were present at task start (except `taxonomy.go`, which changed
mid-session) and are foreign/parallel work from another session, not
artifacts of this A5 spike. They were left untouched, per instruction to
only revert files this spike touched.

Final `git status --porcelain` (post-revert):
```
 M docs/rdr/0029-version-promise-on-machine-readable-output.md
 M internal/cli/respond/respond.go
 M internal/graphlint/taxonomy.go
?? docs/rdr/0029-version-promise-on-machine-readable-output/evidence/reconcile/
?? internal/cli/a7_snapshot_spike_test.go
?? internal/cli/testdata/
```
— identical in kind to the tree's state before this spike began (the
`taxonomy.go` change and the A7 test/testdata additions appeared mid-session
from the concurrent process, not from this spike).

---

## 4. Verdict

**`flow next` unknown-`reason` union: ACHIEVABLE, mechanical.** An
`UnknownReasons() []string` accessor in `internal/cli/flow_next.go`,
matching `graphlint::Reasons`' shape, builds clean with no cycle and no
new import. A5 stands confirmed for this row.

**CLIError `code` vocabulary: NOT achievable as a real enumeration seam
without inverting the leaf layering or moving the constants.**

- A `clierr`-side registry that imports `graphlint`/`guard`/`accessor` to
  draw on their real `Code*` constants produces a confirmed, reproducible
  `go build` import cycle (`clierr` → `graphlint` → `clierr`), because
  `graphlint` already imports `clierr` to construct `CLIError`/`Finding`
  values. This is not avoidable by reordering the spike's own code — the
  cycle is structural, given `clierr`'s existing role as the leaf
  `graphlint`/`respond`/`cli`/`flowbind` all import.
- The only shape that keeps `clierr` a leaf and still builds
  (`internal/cli/clierr::Codes`, spiked) is a hand-maintained, disconnected
  copy of code-literal strings living in `clierr` itself. This satisfies
  "does not cross the forbidden boundary" but fails the actual bar C4 sets
  for a seam: an accessor returning the vocabulary's real backing set, the
  same set the raise sites emit against, so a tier's by-value assertion
  means something. A drifting shadow list is not that.
- The minimum viable shape that would be a REAL seam is what C4/A5 already
  name as the fallback-avoiding option: move the `Code*` /
  `ValidationCode` constants (or re-declare them canonically) into
  `clierr`, with `graphlint`, `guard`, and `accessor` importing
  `clierr.Code*` instead of minting their own. That is a real migration
  across ~36 constants and ~71 raise sites, not an additive accessor, and
  is out of scope for this spike (and arguably out of scope for this RDR,
  which is about tiering existing vocabularies, not consolidating error
  taxonomies).

**Conclusion for A5**: the assumption is refuted for the CLIError `code`
row as stated ("a tier assignment can be backed by an enumeration seam...
without... crossing a package boundary the layering forbids") — no such
seam is buildable today without either the cycle or the migration. The
`flow next` row is confirmed mechanical. Per C4's own contingency
language, the recommended path is the RDR's fallback: the CLIError `code`
row keeps its `append-only` tier assignment but records `seam: none
(prose-only)` with the reason (leaf-layering conflict), rather than
obliging a seam C4 cannot actually get built. The `flow next` row proceeds
with a seam as C4 specifies.
