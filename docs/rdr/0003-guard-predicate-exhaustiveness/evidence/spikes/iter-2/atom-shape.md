Model: claude-opus-5[1m]

# Stage 4 Resolve (scoped re-entry, iter-2) — A5 against as-built kernel

Date: 2026-08-21
Scope: A5 (named in Status qualifier) + anchors the demotion edit touched (A4,
       two new normative clauses citing JDR 0001 §D4 / §JD-4).

## Reuse audit (RDR_ENV reuse-audit paths + internal/)

```sh
$ ls internal/
cli
resolve
version

$ rg -l -i "guard|predicate|resolver|transition" internal/ cmd/
internal/resolve/adversarial_test.go
internal/resolve/resolve_test.go
internal/resolve/mvv_test.go
internal/resolve/fixup_test.go
internal/resolve/resolve.go
internal/resolve/boundary_test.go
internal/resolve/fixtures_test.go
```

FINDING: RDR 0007 has LANDED. internal/resolve/ exists and implements the
kernel seam. RDR 0003's Investigation still claims "there is no implemented
resolver, transition table, or guard package under internal/" — stale.
A5's domain is therefore project source, not peer-RDR prose alone.

## The as-built guard seam

```sh
$ sed -n "77,96p" internal/resolve/resolve.go
```
type GuardResult int

const (
	// GuardFalse: the predicate is decided and does not hold.
	GuardFalse GuardResult = iota
	// GuardTrue: the predicate is decided and holds.
	GuardTrue
	// GuardUnevaluable: the seam cannot decide the predicate. The kernel
	// answers with KindGuardUnevaluable.
	GuardUnevaluable
)

// GuardEvaluator is the delegated guard-evaluation seam owned by RDR
// 0003. The kernel calls it; it does not implement operator semantics.
type GuardEvaluator interface {
	// Evaluate decides guard over the assembled tag-set view.
	Evaluate(guard string, view TagSet) GuardResult
}

// TagSet is the assembled evaluation view: owned, observed, and freshly

```sh
$ sed -n "183,187p" internal/resolve/resolve.go   # Row.Guard
```
	// Guard is the predicate handed to the guard seam. Empty means no
	// guard.
	Guard string

	// NextTags is the next state the row transitions to.

```sh
$ sed -n "425,437p" internal/resolve/resolve.go   # evaluateGuard
```
func evaluateGuard(seam GuardEvaluator, guard string, view TagSet) GuardResult {
	if guard == "" {
		return GuardTrue
	}
	if seam == nil {
		return GuardUnevaluable
	}
	return seam.Evaluate(guard, view)
}

// missingOwned returns the owned tag keys any candidate requires that the
// accessor-produced snapshot does not carry, sorted by key.
//

## Implementation status (the decisive distinction)

`docs/rdr/README.md` index:

| 0001 | Resolution kernel | **Implemented** |
| 0007 | Guard predicate totality | **Final** |

So `internal/resolve/resolve.go` as it stands is **RDR 0001's** kernel.
RDR 0007 is locked-but-unimplemented: it *specifies* the atom-shape change.
`docs/rdr/0007-guard-predicate-totality.md` A3 states the change as future
work — "replacing `Row.Guard string` with a parsed atom slice, narrowing
`GuardEvaluator` to a per-atom value seam" — and its own Evidence is a
re-spike "against the SPECIFIED shape", not against merged code.

### Consequence for A5

A5's citation is sound **as a specification claim** and false **as a
statement about code that exists today**:

| RDR 0003 text | As-built (RDR 0001) | RDR 0007 spec |
| --- | --- | --- |
| "a row carries parsed atoms … not an opaque predicate string" | `Row.Guard string` (resolve.go:185) | parsed atom slice (0007 §A3) |
| "landed in RDR 0007" | not landed — 0007 is Final, not Implemented | — |
| evaluator "MUST NOT read the tag view" | `Evaluate(guard string, view TagSet)` receives the whole view (resolve.go:93) | view-free `Evaluate(atom, value)` (0007:697, 1283) |
| "no reconstruction step that could lose identity" | true for RuleID/SourceLocator (resolve.go:170-172), which the kernel already carries per row | same |

Source identity itself — the actual subject of A5 — **is** carried as built:
`Row.RuleID` and `Row.SourceLocator` (resolve.go::Row), commented "the source
identity RDR 0002 requires every normalized row to retain".

Verdict: A5 holds, but its Evidence must be re-worded to cite RDR 0007 as the
*specified* shape (Peer RDR), not as landed code, and may cite
`internal/resolve/resolve.go::Row` for the source-identity fields that exist now.
