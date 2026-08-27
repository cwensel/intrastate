Model: claude-opus-5[1m]

# A3 reshape probe — RDR 0007 Stage 4 Resolve

**Assumption under test (A3).** Replacing `Row.Guard string` with a parsed
atom slice, narrowing `GuardEvaluator` to a per-atom value seam, and adding
presence + strong-Kleene combination touch `Row`, `GuardEvaluator`,
`evaluateGuard`, `Refusal`, and every test fixture — but leave `gate`'s
prune → owned-state → undecidable order, `escapeOrRefuse`'s uniform gating,
and the frozen ADV/Fixup dispositions unchanged.

**Verdict: CONFIRMED, with one scope correction and one coverage gap.**

---

## 1. Importer sweep (A3's explicit re-verification trigger)

```
$ rg -n -t go --glob '!internal/resolve/**' 'cwensel/intrastate/internal/resolve' .
(no output; exit=1)

$ rg -n -t go 'internal/resolve"' .
./internal/resolve/adversarial_test.go:22:	"github.com/cwensel/intrastate/internal/resolve"
./internal/resolve/resolve_test.go:8:	"github.com/cwensel/intrastate/internal/resolve"
./internal/resolve/fixtures_test.go:4:	"github.com/cwensel/intrastate/internal/resolve"
./internal/resolve/boundary_test.go:26:	return kernelModulePath + "/internal/resolve"
./internal/resolve/fixup_test.go:30:	"github.com/cwensel/intrastate/internal/resolve"
./internal/resolve/mvv_test.go:7:	"github.com/cwensel/intrastate/internal/resolve"

$ for p in $(go list ./...); do go list -f '{{.ImportPath}} -> {{join .Imports " "}}' $p | grep -q 'internal/resolve' && echo "IMPORTER: $p"; done
IMPORTER: github.com/cwensel/intrastate/internal/resolve
```

**Result: NO production importer exists.** Every reference is the package's
own `_test.go` files (all in `package resolve_test`, i.e. the external test
package) plus one string constant in `boundary_test.go`. `boundary_test.go:26`
is a *string literal*, not an import.

`internal/cli/root.go::NewRootCmd` registers exactly one subcommand today:

```go
cmd.AddCommand(newVersionCmd())
```

The CLI does not reach the kernel at all. **The `Row` change is a free
change, not a coordinated migration.** The trigger's negative branch holds.

---

## 2. Kernel anchors (pre-reshape)

| Anchor | Finding |
| --- | --- |
| `internal/resolve/resolve.go::Row` | Has `Guard string` ("Empty means no guard") **and** `RequiresOwned []string`. Also `RuleID`, `SourceLocator`, `Outcome`, `Match`, `NextTags`, `Writes`, `Escape`. |
| `internal/resolve/resolve.go::GuardEvaluator` | `Evaluate(guard string, view TagSet) GuardResult` — a whole-predicate seam. |
| `internal/resolve/resolve.go::GuardResult` | Three values: `GuardFalse` (iota), `GuardTrue`, `GuardUnevaluable`. |
| `internal/resolve/resolve.go::Refusal` | `Kind`, `Revision`, `Flow`, `Recognized`, `Rows []RowRef`, `MissingOwned []string`, `Guard string`. |
| `internal/resolve/resolve.go::evaluateGuard` | `guard == ""` → `GuardTrue` (empty-guard branch); `seam == nil` → `GuardUnevaluable` (nil-seam branch); else delegate. Never evaluates a predicate itself. |
| `internal/resolve/resolve.go::gate` | Partition order confirmed: **prune GuardFalse → missingOwned → undecidable**. |
| `internal/resolve/resolve.go::escapeOrRefuse` | Filters on `rescues(kind)`/`Outcome`/`matches`, then calls the *same* `gate`; `case 1` plans, `case 0` returns the underlying refusal `r`, `default` degrades to `ambiguous_match`. |
| `internal/resolve/resolve.go::assemble` | Merges recognized → observed → owned into one map. Only ever *writes* keys; never deletes or filters. Provenance precedence owned > observed > recognized. |
| `internal/resolve/resolve.go::TagSet` | `Lookup` is provenance-*returning* but selection-blind (no filtering by provenance) and is the only exported accessor besides `Len`. `has`/`matches` are unexported. |
| `internal/resolve/resolve.go::missingOwned` | Over an empty `RequiresOwned` list the inner loop never runs → returns `nil` (len 0), so the row raises no owned-state block. |

### The load-bearing check: does `gate` read guard TEXT for any decision?

**No.** `gate` touches `row.Guard` at exactly two sites:

1. `evaluateGuard(seam, row.Guard, view)` — passes it **through** to the seam
   without inspecting it.
2. `Guard: lowest.Guard` — copies it into the `Refusal` **payload**.

Every decision in `gate` is made on the `GuardResult` verdict, never on the
guard string. A3's "partition unchanged" claim is therefore structurally
supported, not merely hoped for: the guard's *representation* is opaque to
the partition logic. Changing `string` to `[]GuardAtom` cannot reach the
prune/owned/undecidable ordering.

---

## 3. `Refusal.Guard` assertion sweep (the RDR's "If wrong" condition)

```
$ rg -n '\.Guard\b' internal/resolve/*_test.go
```

Only **two** sites in the whole suite read `Refusal.Guard`:

| Site | What it asserts | Representation-dependent? |
| --- | --- | --- |
| `fixup_test.go::TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue` | `got.Refusal.Guard != escape.Guard` — equality against the value the test itself set on the row | **No.** It is a round-trip identity check, not a text check. Re-encodes to `reflect.DeepEqual` over `[]GuardAtom` with zero semantic loss. |
| `resolve_test.go:1091 (TestReq33 "guard unevaluable")` | `r.Guard == ""` — "refusal names *no* undecidable guard for diagnosis" | **No.** It is a non-emptiness check. Re-encodes to `len(r.Guard) == 0`. |

The remaining `.Guard` hits are all **writes** (`Guard: "..."` in fixtures,
`row.Guard = "..."` in mutators) — inputs, not assertions. `adversarial_test.go:320-321`
prints `Refusal.Guard` inside a `t.Errorf` **format string** only; the actual
assertion is `reflect.DeepEqual(forward.Refusal, reversed.Refusal)`, which is
representation-agnostic by construction.

**No test asserts on guard text identity.** The RDR's "If wrong" condition
does not fire.

### `fixtures_test.go::fixtureGuards`

Keys verdicts on the **whole guard string**: `decided map[string]bool`, and
`Evaluate` does `v, ok := g.decided[guard]` — unknown predicate → `GuardUnevaluable`,
known → `GuardTrue`/`GuardFalse`. Empty guard short-circuits to `GuardTrue`
(a redundant mirror of the kernel's own empty branch). Guard strings are
opaque labels ("never", "always", "is-fast-lane", "reviews >= quorum") that
are never parsed — which is exactly why they re-key onto atoms mechanically.

---

## 4. The spike

Scratchpad: `<scratchpad>/spike-a3/`
The probe file is kept at `evidence/spikes/a3-reshape/`. The working
diffs (kernel, fixtures, tests) were scratchpad-only and have been
removed; the substance of each is quoted inline below, and the reshape is
re-derivable from the section that follows.

### Commands

```
$ cp -R internal/resolve <scratch>/spike-a3/internal/resolve
$ printf 'module github.com/cwensel/intrastate\n\ngo 1.26.3\n' > <scratch>/spike-a3/go.mod
$ cd <scratch>/spike-a3 && go test ./... -run . -v -count=1     # baseline
   ... apply reshape ...
$ go vet ./... && go test ./... -run . -v -count=1              # post-reshape
$ go test ./... -count=1 -coverprofile=cov.out                  # combinator coverage
   ... mutate strong-Kleene FALSE-dominance ...
$ go test ./... -count=1                                        # mutation check
```

The module skeleton is the minimal thing that compiles: an empty `go.mod`
with the **same module prefix** so the tests' own import path resolves and
`boundary_test.go`'s `kernelModulePath` constant still matches.

### Baseline (unmodified frozen suite)

```
=== RUN counts: 154
--- PASS: 154
--- FAIL: 0
PASS
ok  	github.com/cwensel/intrastate/internal/resolve	0.182s
```

### The reshape applied

Kernel (165 diff lines, 5 hunks):

- `Row.Guard string` → `Guard []GuardAtom`.
- New `GuardAtom{Key, Op, Literal, Block string}`.
- `GuardEvaluator.Evaluate(guard string, view)` → `EvaluateAtom(atom GuardAtom, value string, view) GuardResult` — a per-atom **value** seam.
- `Refusal.Guard string` → `Guard []GuardAtom` plus `UndecidedAtoms []GuardAtom` (sorted, so the payload stays a function of the tuple).
- `evaluateGuard` now: empty → TRUE; nil seam → UNEVALUABLE; else fold atoms
  with **strong-Kleene conjunction** (any decided FALSE returns FALSE even
  beside an undecided sibling; otherwise one undecided atom makes the whole
  predicate undecidable).
- New `evaluateAtom` does the **kernel-owned presence step** before the seam:
  absent key + `exists` op → decided FALSE; absent key + any other op →
  UNEVALUABLE; present key + `exists` → TRUE; otherwise delegate the value
  comparison. **The seam is never called for an absent key.**

Fixtures: `fixtureGuards` migrated to the
atom seam — it now implements `EvaluateAtom` and decides by comparing the
**bound value** to `atom.Literal`, keyed on `atom.Block`. A `guardAtoms(names...)`
helper re-encodes each old predicate label as a single-atom conjunction over
key `guard.<name>`, and `guardObserved(rows)` binds a value for every guard
key the table mentions so the presence step passes and the atom actually
reaches the value seam. **No row verdict is injected anywhere** — every
verdict is driven through an atom and a value comparison, so the combinator
is on the tested path.

### Post-reshape run

```
=== RUN counts: 160 (154 frozen + 6 spike-only probe cases)
--- PASS: 160
--- FAIL: 0
PASS
ok  	github.com/cwensel/intrastate/internal/resolve	0.250s
```

All 154 frozen tests pass unchanged in disposition.

### Intermediate failure (recorded, and its diagnosis)

The first post-reshape run was **141 PASS / 13 FAIL**. Every failure was
class (a), a mechanical re-encoding artifact of my own fixture wiring: I had
re-keyed the guards onto atom keys but had not yet bound those keys into the
evaluation view, so the kernel's new presence step correctly decided them
UNEVALUABLE before the seam ran. Symptom: `TestAdv2` reported
`guard_unevaluable` where `no_match` was wanted — a guard-FALSE row turning
undecidable. Adding `guardObserved` to bind the keys took it to 154/154.

Note the artifact is itself informative: it is the presence rule **working**.
A guard naming a key the view does not carry is now decided by the kernel
rather than passed to the seam, which is the totality property RDR 0007 wants.

---

## 5. Per-test disposition table (the six frozen tests)

| Test | Depends on guard REPRESENTATION or only VERDICT? | Post-reshape |
| --- | --- | --- |
| `adversarial_test.go::TestAdv1_GuardFalseRowsRequiresOwnedMustNotPoisonAnExactOneMatch` | **Verdict only.** Asserts `!Refused()` and `Plan.RuleID == "rdr.live"`. Guard text is input-side only. | PASS |
| `adversarial_test.go::TestAdv1b_AllGuardsFalseIsNoMatchNotOwnedStateUnavailable` | **Verdict only.** Asserts the escape plan is selected; never reads `Refusal.Guard`. | PASS |
| `adversarial_test.go::TestAdv2_EscapeEdgeMustNotBypassTheGuardSeam` | **Verdict only.** Both sub-cases assert `Refusal.Kind` (`no_match`, `guard_unevaluable`). | PASS (both sub-cases) |
| `adversarial_test.go::TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder` | **Verdict only**, and *structurally* so: the assertion is `reflect.DeepEqual(forward.Refusal, reversed.Refusal)` — a whole-payload comparison that is agnostic to the payload's type. `Refusal.Guard` appears only in the failure-message format string. | PASS |
| `fixup_test.go::TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue` | **Representation-adjacent but not text-dependent.** It compares `Refusal.Guard` to the value the test set on the row. Re-encoded `!=` → `reflect.DeepEqual` over `[]GuardAtom`. The *contract* ("the refusal names the guard it could not decide") is preserved verbatim. | PASS |
| `fixup_test.go::TestFixupGateIsUniformAcrossOrdinaryAndEscapeCandidates` | **Verdict only.** Compares `ordinaryKind` and `escapeKind` against `tc.want` across all three blocking conditions. Guard text lives in `mutate` closures (input side). | PASS (all 3 sub-cases) |

**No frozen assertion had to be reopened.** The single edit to a frozen
assertion (Fixup-1d) is a *type* re-encoding of the same predicate, not a
change in what is asserted. Whether that counts as "unchanged" is a judgment
call the RDR should make explicitly — see the scope correction below.

---

## 6. Did `gate` or `escapeOrRefuse` need changing?

**`escapeOrRefuse`: NO change at all.** Zero diff lines. It never mentions
`Guard`; it delegates entirely to `gate`. Its uniform-gating property is
preserved for free — which is precisely why `TestFixupGateIsUniform...` and
`TestAdv2` pass without being touched.

**`gate`: partition order NO change; refusal payload construction only.**
The prune loop, the `missingOwned` block, and the undecidable loop are
byte-identical. The one hunk inside `gate` builds the refusal:

```go
 		return nil, &Refusal{
-			Kind:  KindGuardUnevaluable,
-			Guard: lowest.Guard,
-			Rows:  rowRefs(undecidable),
+			Kind:           KindGuardUnevaluable,
+			Guard:          copyAtoms(lowest.Guard),
+			UndecidedAtoms: undecidedAtoms(lowest.Guard, seam, view),
+			Rows:           rowRefs(undecidable),
 		}
```

This is payload shape, not partition logic. `gate` still reads guard only to
(1) pass through to `evaluateGuard` and (2) fill the diagnosis payload —
exactly the two roles it had before.

Coverage after the reshape confirms both stay on the tested path:

```
resolve.go:391:	gate			100.0%
resolve.go:441:	evaluateGuard		100.0%
resolve.go:469:	evaluateAtom		 75.0%
resolve.go:562:	escapeOrRefuse		 92.9%
total:					 90.3%
```

---

## 7. Coverage gap found (this is a real finding, not a re-encoding artifact)

A mutation test on the strong-Kleene combinator — swapping FALSE-dominance
for UNEVALUABLE-dominance in `evaluateGuard` — **survived the entire frozen
154-test suite**:

```
$ (mutate evaluateGuard: FALSE no longer short-circuits, UNEVALUABLE does)
$ go test ./... -count=1
ok  	github.com/cwensel/intrastate/internal/resolve	0.247s      # MUTANT SURVIVES
```

The reason: **every guard in the frozen suite is a single atom.** A one-atom
conjunction cannot observe combination order, so the frozen tests pin the
seam's verdicts but say nothing about how multiple atoms combine. The frozen
suite is *sufficient* to prove A3's "dispositions unchanged" claim and
*insufficient* to specify the new combinator.

A 6-case spike-only probe (`evidence/spikes/a3-reshape/spike_kleene_test.go.txt`,
84 lines, NOT part of the frozen suite) closes it, and kills the mutant:

```
$ (same mutation, with the probe present)
--- FAIL: TestSpikeStrongKleeneAndPresence/FALSE_dominates_an_undecided_sibling
FAIL	github.com/cwensel/intrastate/internal/resolve	0.249s      # MUTANT KILLED
```

The probe pins: FALSE dominates an undecided sibling; an undecided atom beside
all-TRUE siblings yields `guard_unevaluable` with the right `UndecidedAtoms`;
an all-TRUE conjunction selects the row; an absent guard key is kernel-decided
without consulting the seam; and an `exists` atom over an absent key is decided
FALSE (pruning the row) rather than undecidable.

**Implication for Phase 1:** A3 is bounded as written, but "leave the frozen
dispositions unchanged" must not be read as "the frozen suite covers the new
behavior". Phase 1 must ADD multi-atom combination tests. That is new test
surface A3's touch-list does not currently name.

---

## 8. Scope correction to A3's touch-list

A3 names `Row`, `GuardEvaluator`, `evaluateGuard`, `Refusal`, and the
fixtures. The spike shows the reshape also requires, inside the kernel:

- a new exported `GuardAtom` type (unavoidable — it is the new seam vocabulary,
  and `boundary_test.go::exportedKernelSymbols` will see it);
- a new `evaluateAtom` helper carrying the kernel-owned presence/existence
  rule — this is where "presence" actually lands, and it is a *new function*,
  not a change to `evaluateGuard`;
- `undecidedAtoms` + `compareAtoms` + `copyAtoms` support functions, needed
  because a slice-valued refusal payload must be sorted and copied to keep
  the existing REQ-1/REQ-10 tuple-determinism property that `rowRefs` and
  `missingOwned` already maintain for their payloads.

None of these disturb `gate` or `escapeOrRefuse`. They are additive. But the
RDR's touch-list should name them so Phase 1 does not read `evaluateGuard`
as the only new logic.

Also worth noting: the reshape does **not** need an exported presence
operator constant to be spelled `exists` — the spike used `OpExists` as a
placeholder. RDR 0002 owns the parsed vocabulary, so the operator alphabet is
a coordination point with 0002, not a free kernel choice.

---

## 9. Verdict

**A3 is VERIFIED as bounded, with the touch-list under-counted and one
coverage caveat.**

The importer sweep — A3's own explicit re-verification trigger — comes back
empty: no production code outside `internal/resolve` imports the kernel, and
`internal/cli/root.go` registers only `newVersionCmd()`. The `Row` change is
therefore a **free change, not a coordinated migration**. The trigger's
conditional does not fire.

The reshape's blast radius is exactly what A3 claims for the parts it names,
and the structural reason is now established rather than assumed: `gate`
never reads guard **text** for any decision — it passes the guard through to
the seam and copies it into the diagnosis payload, and every branch it takes
is on a `GuardResult` verdict. Representation is opaque to the partition. The
spike bears this out empirically: `escapeOrRefuse` has **zero** diff lines,
`gate`'s prune → owned-state → undecidable ordering is byte-identical, and
the sole `gate` hunk is refusal-payload construction. All 154 frozen tests
pass post-reshape with identical dispositions, and **no frozen assertion had
to be reopened** — the one frozen assertion that was edited (Fixup-1d's
`Refusal.Guard` comparison) is a `!=` → `reflect.DeepEqual` type re-encoding
of the identical contract, and the only other `Refusal.Guard` reader
(`resolve_test.go:1091`) is a `== ""` → `len() == 0` non-emptiness check.
No test anywhere asserts on guard-text identity, so the RDR's named "If
wrong" condition does not fire.

Two corrections the RDR should absorb. First, the touch-list is
under-counted: the reshape also adds `GuardAtom`, a new `evaluateAtom`
carrying the presence/existence rule, and `undecidedAtoms`/`compareAtoms`/
`copyAtoms` to keep the slice-valued payload tuple-deterministic under
REQ-1/REQ-10. These are additive and disturb nothing, but "touches
`evaluateGuard`" understates where the new logic lands. Second, and more
consequential for Phase 1: the frozen suite **cannot** validate the new
combinator. A mutation swapping strong-Kleene FALSE-dominance for
UNEVALUABLE-dominance survives all 154 frozen tests, because every frozen
guard is a single atom and a one-atom conjunction cannot observe combination
order. "Dispositions unchanged" is true and must not be misread as "covered".
Phase 1 must add multi-atom combination tests as new surface; a 6-case probe
is sufficient to kill the mutant and is included as evidence.

Nothing here refutes A3. The kernel change is bounded, the importer sweep is
clean, and the frozen contracts survive re-encoding intact.

---

## 8. Post-run housekeeping (main context, not the spike agent)

The probe file was written into the project tree as
`spike_kleene_test.go`, which put it in the module's build: `go build ./...`
passed but `go vet ./...` FAILED on it, because it references the proposed
`resolve.GuardAtom` shape that does not exist in the kernel yet.

```
$ go vet ./...
# .../evidence/spikes/a3-reshape_test
vet: .../a3-reshape/spike_kleene_test.go:13:31: undefined: resolve.GuardAtom
```

Renamed to `spike_kleene_test.go.txt` so it stays readable as evidence
without entering the build. Re-verified from the repo root:

```
$ go vet ./...          → exit 0
$ go test ./...         → internal/cli ok, internal/resolve ok
$ git status --short internal/   → (empty)
```

Pre-existing, out of scope for this stage: `go test ./...` also lists
`docs/rdr/0004-accessor-execution-safety-model/evidence/spikes` as a live
package in the build — the same leftover class, from RDR 0004's spike.
Flagged, not fixed here.
