Model: claude-opus-5[1m]

# A3 reshape RE-SPIKE — RDR 0007 (A3 / A23 / A24 / A26)

The prior spike (`a3-reshape-probe.md`) was judged INSUFFICIENT because it
diverged from the specified design on three axes. This re-spike redoes the
reshape **as the Normative Contracts specify**:

| Axis | Prior spike | This re-spike |
| --- | --- | --- |
| (a) `Refusal.Guard` | RETAINED, retyped to `[]GuardAtom`, with `UndecidedAtoms` beside it | **DELETED.** `Refusal.Undecided []UndecidedRow` replaces it |
| (b) Seam | `EvaluateAtom(atom, value, view)` — still passed the view | `Evaluate(atom GuardAtom, value string) GuardResult` — **never sees the view** |
| (c) Existence atom | `{Key:"guard.missing", Op: OpExists}` with NO literal | **Explicit literal, both polarities**, plus the fail-closed leg |

**Overall: A3 VERIFIED (with the re-decide the RDR predicted), A23 VERIFIED,
A24 VERIFIED, A26 VERIFIED.**

---

## 0. Method and hygiene

Worked in a scratch copy; `internal/` in the repo was never edited.

```
$ S=<scratchpad>/respike
$ cp -R internal/resolve "$S/internal/resolve"
$ printf 'module github.com/cwensel/intrastate\n\ngo 1.26.3\n' > "$S/go.mod"
```

The module skeleton uses the **same module prefix** so the external test
package's own import path resolves and `boundary_test.go`'s
`kernelModulePath` constant still matches.

Baseline, before any edit:

```
$ go test ./internal/resolve/... -count=1 -v
PASS lines: 154
FAIL lines: 0
top-level RUN: 154
ok  	github.com/cwensel/intrastate/internal/resolve	0.174s
```

154 PASS / 0 FAIL — matches the prior spike's baseline exactly.

## 1. The reshape as specified

Applied to `resolve.go`:

1. **`Row.Guard string` → `Guard []GuardAtom`.** `GuardAtom` carries the four
   §D1 fields spelled as the contract fixes them: `Key`, `Operator`,
   `Literal`, `Block`.
2. **`Block` is an exported named STRING type** with exactly two constants,
   `BlockAll Block = "all"` and `BlockUnless Block = "unless"`.
3. **Seam narrowed to the per-atom VALUE seam that never sees the view:**

   ```go
   type GuardEvaluator interface {
   	// Evaluate decides atom against the present value bound to atom.Key.
   	Evaluate(atom GuardAtom, value string) GuardResult
   }
   ```

   The view parameter is **gone**, not merely unused. This is the SEAM
   clause's "It never sees the view" enforced by the signature.
4. **`evaluateAtom`** carries the kernel-owned presence + existence rules;
   the seam is unreachable for an existence atom or an absent-key atom.
5. **Strong-Kleene combination in `evaluateGuard`**, against an explicit
   truth-order rank (`F < U < T`), **never** `min` over the raw constants
   (whose order is `F < T < U`). Row verdict is
   `all_result ∧ ¬(unless_conj)`, with the negated term dropping out when
   the `unless` block is empty or omitted.
6. **`Refusal.Guard` DELETED.** Added `Refusal.Undecided []UndecidedRow`,
   `UndecidedRow{RuleID, SourceLocator, Atoms}`,
   `UndecidedAtom{Key, Block, Operator, Literal, Reason}`, the closed
   `Reason` type with `ReasonAbsent` / `ReasonUncomparable`, and the
   `Reasons() []Reason` enumerator mirroring `RefusalKinds()`.
7. **Exported constants:** `OpExists = "exists"`, `LiteralTrue = "true"`,
   `LiteralFalse = "false"`, `BlockAll = "all"`, `BlockUnless = "unless"`.

`slices.MinFunc` over the undecidable set is **retired**, not re-typed —
every undecidable row appears in the payload, so there is no representative
to choose.

Full kernel source: `a3-reshape/respike_resolve.go.txt`.

## 2. Command log

```
$ go build ./...                                    # exit 0 (kernel compiles)
$ go vet ./...                                      # exit 0 after fixture migration
$ go test ./internal/resolve/... -count=1 -v        # runs below
```

Three successive runs are recorded because the failures between them are
the load-bearing evidence.

### Run 1 — reshape applied, Fixup-1d re-encoded FAITHFULLY

```
exit=1
PASS: 147  FAIL: 7
--- FAIL: TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue
--- FAIL: TestReq15_ExactlyOneMatchAfterGuardsIsTheOnlySuccess
--- FAIL: TestReq23_GuardEvaluationIsDelegatedToTheInjectedSeam
    --- FAIL: .../seam_says_fast_lane_holds
    --- FAIL: .../seam_says_slow_lane_holds
    --- FAIL: .../seam_decides_both_true_so_selection_is_ambiguous
    --- FAIL: .../seam_decides_both_false_so_nothing_matches
```

REQ-15 and REQ-23 are **fixture-wiring artifacts**, not contract failures:
both swap in `twoRowsOneGuardFalseTable()` after `legalInput()` has already
bound the guard keys, so the new guard keys were never bound into the view
and the kernel correctly decided them unevaluable **on absence** before the
seam ran. That is the presence rule working. Fixed by re-deriving the
bindings (`in = withGuardBindings(in)`) at both swap sites.

### Run 2 — bindings fixed; Fixup-1d still faithful

```
exit=1
PASS: 153  FAIL: 1
--- FAIL: TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue
```

Verbatim failure:

```
=== RUN   TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue
    fixup_test.go:76: refusal names undecided [{RuleID:rdr.escape.nilseam
      SourceLocator:flows/rdr.toml:90 Atoms:[{Key:iterations Block:all
      Operator:gte Literal:3 Reason:absent}]}]; want [{RuleID:rdr.escape.nilseam
      SourceLocator:flows/rdr.toml:90 Atoms:[{Key:iterations Block:all
      Operator:gte Literal:3 Reason:uncomparable}]}] for diagnosis
--- FAIL: TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue (0.00s)
```

**This is the RDR's predicted RE-DECIDE, reproduced.** See §3.

### Run 3 — Fixup-1d RE-DECIDED per the RDR's own prescription

```
exit=0
PASS: 154  FAIL: 0
ok  	github.com/cwensel/intrastate/internal/resolve	0.241s
```

### Run 4 — full suite with the re-spike probes added

```
$ go test ./internal/resolve/... -count=1 -v
exit=0
PASS: 180  FAIL: 0     # 154 frozen + 26 probe cases
$ go vet ./...
vet exit=0
```

Frozen-only re-run, to show the 154 are intact and not diluted by the probe:

```
$ go test ./internal/resolve/... -count=1 -v -run 'TestReq|TestAdv|TestFixup|TestMVV'
--- PASS count: 154
```

---

## Q-A3 — bound + frozen-suite disposition

### The touch-list holds, with the same additive names the prior spike found

`escapeOrRefuse`: **ZERO diff.**

```
$ diff <(sed -n '/^func escapeOrRefuse/,/^}/p' <baseline>/resolve.go) \
       <(sed -n '/^func escapeOrRefuse/,/^}/p' <respike>/resolve.go)
escapeOrRefuse: ZERO diff
```

`gate`: the prune → owned-state → undecidable **ORDER is byte-identical**.
The complete diff is payload construction plus the `MinFunc` retirement:

```
3a4
>       unevaluable := make([][]UndecidedAtom, 0, len(rows))
5c6
<               verdict := evaluateGuard(seam, row.Guard, view)
---
>               verdict, atoms := evaluateGuard(seam, row.Guard, view)
10a12
>               unevaluable = append(unevaluable, atoms)
21a24
>       var payload []UndecidedRow
27a31,35
>                       payload = append(payload, UndecidedRow{
>                               RuleID:        row.RuleID,
>                               SourceLocator: row.SourceLocator,
>                               Atoms:         unevaluable[i],
>                       })
34,36c42
<               lowest := slices.MinFunc(undecidable, func(a, b Row) int {
<                       return compareRefs(refOf(a), refOf(b))
<               })
---
>               sortUndecided(payload)
38,40c44,46
<                       Kind:  KindGuardUnevaluable,
<                       Guard: lowest.Guard,
<                       Rows:  rowRefs(undecidable),
---
>                       Kind:      KindGuardUnevaluable,
>                       Rows:      rowRefs(undecidable),
>                       Undecided: payload,
```

The prune loop, the `missingOwned` block, and the survivor partition are
unchanged. The `slices.MinFunc` representative selection is retired exactly
as the payload clause requires.

Files touched: `resolve.go`, `fixtures_test.go`, `adversarial_test.go`,
`fixup_test.go`, `resolve_test.go`. **`boundary_test.go` and `mvv_test.go`
are untouched** — the boundary inspection helpers needed no change at all.

The additive kernel names the prior spike flagged are confirmed again:
`GuardAtom`, `evaluateAtom`, plus `truthRank` / `kleeneAnd` / `kleeneNot` /
`sortUndecided` / `compareUndecidedAtoms` (this re-spike's spelling; the
prior spike called the last group `undecidedAtoms`/`compareAtoms`/`copyAtoms`).
Note one genuine simplification the "replace" shape buys: **no `copyAtoms`
is needed**, because the payload no longer echoes the row's guard slice —
it is built fresh from the verdict pass, so it cannot alias the caller's
table.

### The frozen ADV/Fixup dispositions

| Test | Disposition |
| --- | --- |
| `TestAdv1_GuardFalseRowsRequiresOwnedMustNotPoisonAnExactOneMatch` | PASS, re-encoded (atoms + guard-key bindings) |
| `TestAdv1b_AllGuardsFalseIsNoMatchNotOwnedStateUnavailable` | PASS, re-encoded |
| `TestAdv2_EscapeEdgeMustNotBypassTheGuardSeam` (both sub-cases) | PASS, re-encoded |
| `TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder` | PASS, re-encoded. Its assertion is `reflect.DeepEqual(forward.Refusal, reversed.Refusal)` — payload-shape-agnostic. Only its `t.Errorf` **format string** read `Refusal.Guard`; that is a diagnostic message, not an assertion, and was retargeted to `Refusal.Undecided`. |
| `TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue` | **RE-DECIDED** — see below |
| `TestFixupGateIsUniformAcrossOrdinaryAndEscapeCandidates` (3 sub-cases) | PASS, re-encoded |

### The critical question: Fixup-1d is RE-DECIDED, not re-encoded

The frozen test:

```go
func TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue(t *testing.T) {
	escape := escapeRow("rdr.escape.nilseam", "flows/rdr.toml:90", resolve.KindNoMatch)
	escape.Guard = "iterations >= 3"

	in := noMatchInput()
	in.Table.Revision = "rev-fixup-1d"
	in.Table.Rows = append(in.Table.Rows, escape)
	in.Guards = nil // no seam at all, not merely a seam that cannot decide

	got := mustResolve(t, in)

	if !got.Refused() { ... }
	if got.Refusal.Kind != resolve.KindGuardUnevaluable {
		t.Errorf("refusal kind = %q; want %q — an absent seam must refuse, not "+
			"let the escape edge through unevaluated",
			got.Refusal.Kind, resolve.KindGuardUnevaluable)
	}
	if got.Refusal.Guard != escape.Guard {
		t.Errorf("refusal names guard %q; want %q for diagnosis",
			got.Refusal.Guard, escape.Guard)
	}
}
```

**The asserted PROPERTY does not survive.** Split by assertion:

1. `!got.Refused()` and `Kind == KindGuardUnevaluable` — **survive**. The
   verdict is unchanged.
2. `got.Refusal.Guard != escape.Guard` — **has no referent.** `Refusal.Guard`
   is deleted. This is a compile error, not a comparison that can be
   re-typed:

   ```
   $ go vet ./...
   vet: internal/resolve/resolve_test.go:1091:8: r.Guard undefined
        (type *resolve.Refusal has no field or method Guard)
   ```

3. Most importantly: **the proposition the test's NAME states no longer
   holds on this fixture.** The guard `iterations >= 3` references a key
   `legalInput`'s view does not carry (it supplies `status`, `reviews`,
   `recognized`). Under the reshape the KERNEL decides that atom unevaluable
   **on ABSENCE**, and the nil seam is never consulted. Run 2's verbatim
   output proves it: the kernel reported `Reason:absent`, where a nil-seam
   rescue-refusal would report `Reason:uncomparable`. The test would still
   go green on kind alone while testing **nothing about the nil seam** —
   it would pass identically with a fully-wired seam.

This is exactly A3's own "If wrong" branch, and it is exactly what the RDR's
SEAM clause already predicted in prose. The re-spike confirms it empirically
rather than by inspection.

**The RDR's prescribed fix works.** Per the SEAM clause ("To keep testing
the nil-seam rule the fixture needs a value atom over a PRESENT key (e.g.
`reviews >= 3`)"), the atom was re-decided onto the present key `reviews`:

```go
	// RE-DECIDED for the atom-shaped guard (RDR 0007 SEAM clause). ...
	escape.Guard = []resolve.GuardAtom{{
		Key: "reviews", Operator: "gte", Literal: "3", Block: resolve.BlockAll,
	}}
	...
	wantUndecided := []resolve.UndecidedRow{{
		RuleID:        escape.RuleID,
		SourceLocator: escape.SourceLocator,
		Atoms: []resolve.UndecidedAtom{{
			Key: "reviews", Block: resolve.BlockAll,
			Operator: "gte", Literal: "3",
			// The key is PRESENT, so `absent` would be a lie; a nil seam is
			// a wiring fact, not a reason of its own.
			Reason: resolve.ReasonUncomparable,
		}},
	}}
	if !reflect.DeepEqual(got.Refusal.Undecided, wantUndecided) { ... }
```

With that, the test again tests what its name says: a value atom over a
present key **would** have gone to the seam, the seam is nil, so the atom is
`uncomparable` and the escape does not rescue. Run 3: 154 PASS / 0 FAIL.

### Coverage gap: reconfirmed, with a CORRECT existence case

The prior spike's mutant is re-run. Two mutations were tried:

**Mutant 1 — integer `min` over the raw constants** (`if a <= b`), the exact
error the combination clause names:

```
--- frozen suite ONLY, against mutant 1 ---
    resolve_test.go:1160: expected a refusal disposition; got plan {RuleID:rdr.guarded.successful ...}
FAIL	github.com/cwensel/intrastate/internal/resolve	0.249s
```

**KILLED by the frozen suite.** New finding, and a mild correction to the
prior spike: this particular mutation is *not* invisible to the 154, because
`GuardTrue <= GuardUnevaluable` under the shipped constant order makes even a
**single-atom** `U` guard fold to `T` against the `allResult = GuardTrue`
accumulator, which the frozen `guard_unevaluable` fixtures see directly.

**Mutant 2 — swap FALSE-dominance for UNEVALUABLE-dominance** (rank `F`→1,
`U`→0), which only a multi-atom conjunction can observe:

```
--- frozen suite ONLY (154), against mutant 2 ---
ok  	github.com/cwensel/intrastate/internal/resolve	0.264s      # SURVIVES

--- full suite incl. respike probe ---
--- FAIL: TestRespikeStrongKleeneCombination
    --- FAIL: TestRespikeStrongKleeneCombination/F_dominates_U_(F_AND_U_=_F):_the_row_is_pruned,_not_undecidable
FAIL	github.com/cwensel/intrastate/internal/resolve	0.177s      # KILLED
```

**A3's coverage-gap finding stands.** Every frozen guard is a single atom, so
the frozen suite cannot specify the combinator. Phase 1 must add multi-atom
combination tests; `a3-reshape/respike_kleene_test.go.txt` is the shape, and
unlike the prior probe **its existence cases are correct against this draft**.

### Verdict Q-A3: **VERIFIED**

The bound holds: `escapeOrRefuse` zero diff, `gate`'s partition order
byte-identical, `boundary_test.go` and `mvv_test.go` untouched, 154/154
frozen dispositions preserved. On the critical sub-question the answer is
the RDR's own claim, now demonstrated: **Fixup-1d must be RE-DECIDED, not
re-encoded.** Its `Refusal.Guard` assertion has no referent under "replace"
(compile error), and — the deeper point — its fixture guard names an absent
key, so the kernel decides it on ABSENCE (`Reason:absent`, captured verbatim
in Run 2) and the nil seam is never consulted. The property the test's name
asserts is gone; moving the atom to the present key `reviews` restores it.

---

## Q-A23 — payload additivity

### Does deleting `Refusal.Guard` / adding `Refusal.Undecided` break a frozen field-set assertion?

Sweeps over the **baseline** suite:

```
$ grep -rn 'resolve\.Refusal{\|Refusal{' internal/resolve/*_test.go
(none)
```

**No test constructs a `Refusal` composite literal at all** — so no test can
be broken by a field-set change through that route. (Composite literals with
unkeyed fields would be; keyed ones would still compile. Neither exists.)

```
$ grep -rn 'DeepEqual(.*Refusal' internal/resolve/*_test.go
adversarial_test.go:313:	if !reflect.DeepEqual(forward.Refusal, reversed.Refusal) {
adversarial_test.go:375:	if !reflect.DeepEqual(forward.Refusal, reversed.Refusal) {
fixup_test.go:200:	if !reflect.DeepEqual(forward.Refusal, reversed.Refusal) {
fixup_test.go:236:	if !reflect.DeepEqual(forwardRefusal, reversedRefusal) {
```

All four are **self-comparisons** (forward vs. reversed row order) — they
compare two `Refusal` values against each other, never against a literal with
a fixed field set. They are structurally agnostic to `Refusal`'s shape and
required **zero** edits. They are also the tests that would catch a
non-deterministic payload, and they pass, which is independent evidence that
`sortUndecided` makes `Undecided` a function of the input tuple.

### `TestReq7_RefusalKindSetIsExactlyTheFiveNamedKinds`

```
$ go test ./internal/resolve/... -count=1 -v -run 'TestReq7_RefusalKindSetIsExactlyTheFiveNamedKinds'
=== RUN   TestReq7_RefusalKindSetIsExactlyTheFiveNamedKinds
--- PASS: TestReq7_RefusalKindSetIsExactlyTheFiveNamedKinds (0.00s)
PASS
```

**PASS, unedited.** As A23 predicted, it pins the KIND set, not the field
set. The new `Reason` closed set is a *separate* type and does not enter
`RefusalKinds()`; the five kinds are untouched, so the payload is "a field,
not a sixth kind" as the Load-Bearing Decision states.

### Tests that DID need edits, and why

| Test | Edit | Why |
| --- | --- | --- |
| `resolve_test.go::TestReq33_...` (`guard unevaluable` sub-case) | `r.Guard == ""` → `len(r.Undecided) == 0` | **Compile-forced** (`r.Guard undefined`). A pure non-emptiness check — "the refusal names *something* undecidable for diagnosis" — so the property re-encodes onto the replacement payload with zero semantic loss. |
| `fixup_test.go::TestFixup1d_...` | `Refusal.Guard` equality → `DeepEqual` over `Refusal.Undecided`, **and** the guard atom moved to a present key | **Compile-forced AND semantically forced.** See Q-A3. This is a re-decide, not a re-encode. |
| `adversarial_test.go::TestAdv3_...` | `t.Errorf` format string retargeted `Refusal.Guard` → `Refusal.Undecided` | Compile-forced, but it is a **failure-message** interpolation, not an assertion. No property changed. |
| `fixtures_test.go`, plus guard-key binding sites in `adversarial_test.go` / `fixup_test.go` / `resolve_test.go` | atoms + `guardBindings` | Representation migration A3 already names ("every test fixture"). |

**No frozen test asserts `Refusal`'s field set exhaustively.** The only two
sites that read `Refusal.Guard` at all are the two the prior spike already
identified, and neither asserts guard-text identity.

### Verdict Q-A23: **VERIFIED**

`Refusal.Undecided` lands as a bare field, `UndecidedRow` / `UndecidedAtom`
/ `Reason` / `Reasons()` land as new top-level surface, and no frozen
contract had to be reopened. The one test whose *property* changed
(Fixup-1d) changed for the SEAM reason in Q-A3, not because `Refusal`'s
shape is frozen. A23's "If wrong" branch (payload lands behind an accessor)
does **not** fire.

---

## Q-A24 — existence literal, both polarities

### The clause asserted against

The existence clause, quoted verbatim:

> The existence operator is the sole TOTAL operator: it MUST decide true or
> false from presence or absence of its referenced key alone. Its verdict is
> `presence == literal`: `exists = true` decides TRUE when the key is
> present, `exists = false` decides TRUE when it is absent.

And the fail-closed leg:

> A foreign LITERAL on an `OpExists` atom — any value that is neither
> `LiteralTrue` nor `LiteralFalse`, the empty literal included — is
> UNEVALUABLE at the kernel, reason `uncomparable`; the kernel MUST NOT
> decide such an atom from presence, and MUST NOT treat a missing literal as
> `LiteralFalse`.

Note this **corrects the task brief's sketch** in one respect, and the
assertions follow the RDR, not the sketch: the brief called `exists=false`
over an absent key "the one sanctioned route from absence to a verdict".
That phrase is A5's, about the *authoring idiom*; the clause itself states
the rule as `presence == literal`, which is what is asserted here. The
polarities the brief predicted are the ones the clause yields, so the
assertions coincide — but they are anchored on `presence == literal`.

### Both polarities, with an EXPLICIT literal

```
--- PASS: TestRespikeExistenceBothPolarities (0.00s)
    --- PASS: .../exists=true_over_an_ABSENT_key_decides_FALSE (0.00s)
    --- PASS: .../exists=false_over_an_ABSENT_key_decides_TRUE (0.00s)
    --- PASS: .../exists=true_over_a_PRESENT_key_decides_TRUE (0.00s)
    --- PASS: .../exists=false_over_a_PRESENT_key_decides_FALSE (0.00s)
```

- `exists = LiteralTrue` over ABSENT `guard.missing` → decides **FALSE** →
  the row is pruned → zero candidates → `no_match`. The row does not apply.
- `exists = LiteralFalse` over ABSENT `guard.missing` → decides **TRUE** →
  the row is **selected** (a plan for `rdr.probe`). Absence reached a
  verdict, by the only route that may.
- The two present-key mirrors both hold.

Each case also asserts the seam was **never consulted** (`len(calls) != 0`
fails the test) — the contract's "the evaluator MUST NOT be consulted for
it" — and that a DECIDED existence atom reports **no** undecided payload.

### The fail-closed leg

```
--- PASS: TestRespikeExistenceForeignLiteralFailsClosed (0.00s)
    --- PASS: .../MISSING_literal_(zero_value)_/_key=reviews (0.00s)
    --- PASS: .../MISSING_literal_(zero_value)_/_key=guard.missing (0.00s)
    --- PASS: .../malformed_literal_/_key=reviews (0.00s)
    --- PASS: .../malformed_literal_/_key=guard.missing (0.00s)
    --- PASS: .../foreign_literal_/_key=reviews (0.00s)
    --- PASS: .../foreign_literal_/_key=guard.missing (0.00s)
```

Three literal shapes — `""` (the missing/zero-value literal the prior spike
wrongly decided on), `"True"` (case drift, the exact `"exists"` vs `"Exists"`
hazard the clause names), and `"yes"` — each over both a present and an
absent key. All six yield `guard_unevaluable` with

```go
UndecidedAtom{Key: <key>, Block: BlockAll, Operator: OpExists,
              Literal: <foreign>, Reason: ReasonUncomparable}
```

asserted by `reflect.DeepEqual` against the full payload. In particular the
absent-key/missing-literal case does **not** decide TRUE (which is what
treating a missing literal as `LiteralFalse` would have produced), and the
present-key/missing-literal case does **not** decide TRUE either. The kernel
performs no parsing, case-folding, or coercion: `switch atom.Literal` over
the two exported constants, `default:` → `(GuardUnevaluable, ReasonUncomparable)`.

### Verdict Q-A24: **VERIFIED**

Verbatim comparison against two exported constants is total over well-formed
input and every other value is `uncomparable`. The kernel needs **no**
literal-parse step, so A24's "If wrong" branch does not fire. The prior
spike's existence case is confirmed WRONG against this draft and is
superseded by `respike_kleene_test.go.txt`.

---

## Q-A26 — exported-name sweep

### The three frozen boundary tests, against the compiled new surface

```
$ go test ./internal/resolve/... -count=1 -v -run 'TestReq25_...|TestReq36_...|TestReq37_...'
=== RUN   TestReq25_NoExportedSymbolImpliesOrchestrationOrPersistence
--- PASS: TestReq25_NoExportedSymbolImpliesOrchestrationOrPersistence (0.00s)
=== RUN   TestReq36_KernelExposesNoEncodeDecodeOrInverseOperation
--- PASS: TestReq36_KernelExposesNoEncodeDecodeOrInverseOperation (0.00s)
=== RUN   TestReq37_KernelIntroducesNoHashOrCanonicalSerialization
--- PASS: TestReq37_KernelIntroducesNoHashOrCanonicalSerialization (0.00s)
PASS
ok  	github.com/cwensel/intrastate/internal/resolve	0.176s
```

**All three PASS.** These run `exportedKernelSymbols(t)`, which AST-walks the
real package, so they saw the actual new declarations: `Block`, `BlockAll`,
`BlockUnless`, `OpExists`, `LiteralTrue`, `LiteralFalse`, `GuardAtom`,
`Reason`, `ReasonAbsent`, `ReasonUncomparable`, `Reasons`, `UndecidedRow`,
`UndecidedAtom`.

### Mechanical sweep of the fourteen names

`respike_sweep_test.go.txt` checks each proposed name against the union of
all three banned lists, under **both** the frozen tests' `==` rule and the
stricter **substring** reading A26 states:

```
--- PASS: TestRespikeNewExportedNamesClearBannedLists (0.00s)
```

Banned set swept: `Run`, `Execute`, `Apply`, `Persist`, `Save`, `Commit`,
`Start`, `Orchestrate` (REQ-25); `Encode`, `Decode`, `Marshal`, `Unmarshal`,
`Import`, `Export`, `Parse`, `Serialize`, `Deserialize`, `Invert` (REQ-36);
`Hash`, `Checksum`, `Canonicalize`, `Fingerprint`, `Digest` (REQ-37).

**No collision, exact or substring**, for any of `GuardAtom`,
`UndecidedRow`, `UndecidedAtom`, `Undecided`, `OpExists`, `LiteralTrue`,
`LiteralFalse`, `ReasonAbsent`, `ReasonUncomparable`, `Reasons`, `Block`,
`BlockAll`, `BlockUnless`, `TestGuardEvaluatorContract`.

### One finding the sweep surfaced

`Undecided` is a **FIELD of `Refusal`**, not a top-level declaration.
`boundary_test.go::exportedKernelSymbols` walks top-level decls only
(`*ast.FuncDecl`, and `*ast.TypeSpec` / `*ast.ValueSpec` inside
`*ast.GenDecl`) — it never descends into struct fields. So the three frozen
boundary tests are **structurally incapable of seeing `Undecided`** at all.

This does not refute A26 — the name is clear under the substring reading
anyway, and this re-spike verified that separately — but it is a scope fact
worth recording: A26's claim for `Undecided` is established by the standalone
sweep, not by REQ-25/36/37, which cannot reach it. The same is true of every
field name on the new payload types (`Key`, `Block`, `Operator`, `Literal`,
`Reason`, `Atoms`, `RuleID`, `SourceLocator`).

`TestGuardEvaluatorContract` is Phase 3 surface not declared by this
re-spike; it is swept as a name only. It also lives in a `_test.go` file,
which `parseKernelPackage` filters out, so it is doubly out of the frozen
tests' reach.

### Verdict Q-A26: **VERIFIED**

The surface compiles and all three frozen boundary tests pass against it.
No proposed name collides with a banned substring. A26's "If wrong" branch
does not fire; no normative name needs renaming and no RDR 0001 contract
needs reopening.

---

## Summary of verdicts

| Question | Verdict | Reason |
| --- | --- | --- |
| **Q-A3** | **VERIFIED** | Bound holds (`escapeOrRefuse` zero diff, `gate` order byte-identical, `boundary_test.go`/`mvv_test.go` untouched, 154/154). Fixup-1d **must be RE-DECIDED**: its `Refusal.Guard` assertion has no referent under "replace", and its fixture guard names an ABSENT key, so the kernel decides it `Reason:absent` and the nil seam is never consulted — the property its name asserts is gone. The RDR's own prescribed fix (a value atom over present `reviews`) restores it. |
| **Q-A23** | **VERIFIED** | No frozen test constructs a `Refusal` literal; all four `DeepEqual`-over-`Refusal` sites are self-comparisons and needed zero edits; `TestReq7` PASSES unedited (it pins the kind set, not the field set). Only two tests needed edits, both compile-forced by the `Refusal.Guard` deletion. |
| **Q-A24** | **VERIFIED** | Both polarities hold on an explicit literal (`exists=true`/ABSENT → FALSE → pruned; `exists=false`/ABSENT → TRUE → selected), the seam is never consulted for an existence atom, and all six foreign/missing-literal cases yield `uncomparable` rather than a verdict. No literal-parse step needed. |
| **Q-A26** | **VERIFIED** | REQ-25 / REQ-36 / REQ-37 all PASS against the compiled surface; no name collides exact or substring. Scope note: `Undecided` is a struct field, which `exportedKernelSymbols` cannot see, so its clearance rests on the standalone sweep. |

**No question is REFUTED.** The specified design compiles and no frozen test
genuinely breaks in a way the RDR has not already anticipated and prescribed
a fix for.

---

## Artifacts

- `a3-reshape/respike_resolve.go.txt` — the reshaped kernel (the shape proved)
- `a3-reshape/respike_kleene_test.go.txt` — Kleene combination, both existence
  polarities, the fail-closed leg, nil-seam per-atom rule, `unless` block
- `a3-reshape/respike_sweep_test.go.txt` — the A26 banned-name sweep
- `a3-reshape/respike_fixtures_test.go.txt` — the migrated fixture builders
  (per-atom value seam, `guardAtoms`, `guardBindings`, `withGuardBindings`)

## Housekeeping

All work was done in a scratch copy. The repository's `internal/` tree was
never edited; the probe files were never placed inside the repo module (the
prior spike's `go vet` breakage class is avoided by construction). Verified:

```
$ git status --porcelain
(only the new evidence files under docs/rdr/0007-.../evidence/spikes/)
```
