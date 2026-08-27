Model: claude-opus-5[1m]

# cove — Chain-of-Verification findings, 0010 (stateless decision tables)

Step 0 (grounding) ran as an isolated sub-agent; its output is
[`step0-grounding.md`](step0-grounding.md) — verdict PASS, no codebase claim
refuted, 42/42 anchors semantically accurate, A4's 103-file census reproduced
exactly, and the inverse check decisive (no sibling class discriminator, no
`emit`, no ∅-root seeding anywhere in the tree).

Steps 1–3 below are the silence/contradiction half, read against the whole
record. Findings are merged across both halves and de-duplicated; where the two
halves converged independently on one passage that is noted, since independent
convergence is the strongest signal in this lens.

## Step 1–2 — questions whose answers produced no finding

Recorded so the "searched, found nothing" half is auditable rather than implied.

- **Is `class` visible at rule-normalization time?** Yes. `normalizeRule` is a
  method on `*loader` holding `l.model`, and `loadModelHeader` is first in
  `run()`'s step list (`internal/table/load.go:78-84`). C2's class-conditioned
  arm is implementable exactly as written. CONFIRMED.
- **Does the dump render from the kernel row?** No — `writeModel` iterates
  `m.Rows` (`table.Row`). C3's "`emit` never crosses the kernel boundary, the
  dump renders it" is coherent. CONFIRMED.
- **Is there a second default column set distinct from `dumpColumns`?** No.
  `internal/table/dump.go::writeRow` sets `order = dumpColumns` when
  `m.DumpOrder` is empty, and `loadDump` sets `DumpOrder = DumpColumns()` when
  `[dump]` is absent. C3's two clauses ("MUST join the vocabulary" / "the
  default set MUST include it") name one surface twice — redundant, not wrong.
  A4 already states this. No finding.
- **Does `checkNodeCeiling` fire over a one-node graph?** No. The ceiling is
  enforced inside `reach` (`internal/graphlint/reach.go:116`) and reported by
  `checkNodeCeiling` via `complete`. Over one node it cannot trip. A10's
  placement of it in "exercised unchanged" holds. CONFIRMED.
- **Does a pre-change binary refuse `[rule.emit]` as well as `class`?** Yes —
  `decodeStrict` disallows unknown fields recursively over the struct graph. A8
  already carries this. CONFIRMED.
- **Is C2's `[initial]` prohibition really "already unauthorable"?** Yes, but
  by a category C2 does not name. Executed, not argued — a pure decision-table
  model (zero owned tags, escape row, no write/clear, no accessors) carrying a
  stray `[initial] status = "Draft"` over a *declared observed* tag refuses
  `malformed_accessor_binding: written tag status is served by 0 writers; want
  exactly one` (spike: `evidence/spikes/cove-dt-initial/`). `loadInitial` itself
  refuses only an *undeclared* key (`CatUnknownTag`,
  `internal/table/load.go:546`) and explicitly defers the non-owned case to the
  writer binding (comment at `load.go:574-577`). The deferral survives on a
  decision table because `[initial]` keys are treated as written tags demanding
  a writer, and a decision table declares none. **Claim holds; candidate finding
  dismissed by execution.**

## Findings

### 1 — `0010:C2` × `0010:C5` contradict each other on the `terminal` prohibition (blocking)

**Both halves of this lens reached this passage independently** (Step 0 finding
2, from the enforcement-point direction; Step 3 below, from the contract-
contradiction direction).

C2 asserts the `terminal` prohibition is not a new refusal because "a terminal
predicate over a non-owned tag is 0006's `graph-dangling-edge` terminal arm" —
i.e. it relies on that arm **firing** over a decision table.

C5 requires the opposite: "... and **declared-terminal handling (7)** are
vacuous by construction over this class and **MUST NOT emit**", and A10 lists
`graph-dangling-edge` wholesale in the "provably silent" bucket, attributing its
silence to "C5's override of `0006:C18`" — which is the *missing-root* arm only.

`internal/graphlint/analysis.go::checkDanglingEdge` has **two** arms emitting one
code: the missing-root arm (`element: "model"`) and the terminal-predicate arm
(`element: "terminal[i]"`). Executed against a pure decision-table shape carrying
a stray `terminal` over an observed tag, **both** fire
(`evidence/spikes/cove-dt-terminal/`):

```
graph-dangling-edge  element=model        "the model declares no initial owned state…"
graph-dangling-edge  element=terminal[0]  key=status  "terminal predicate reads the observed tag
                     \"status\"; a terminal is a predicate over owned tags…"
```

A10's per-code grounding of the silent bucket covers `checkDeadEnd` and
`checkTerminalEscape` ("return early on absent terminals") but never grounds
`checkDanglingEdge`'s terminal arm — the one arm C2 depends on. So an
implementer reading C5/A10 suppresses `graph-dangling-edge` by class and C2's
sole enforcement of the `terminal` prohibition disappears with no diagnostic:
a decision table could declare `terminal` and lint clean.

**Why it matters**: C2's "none of these is a new refusal" is the clause that
keeps this RDR's override surface small. If the terminal arm must stay live for
the class, C2 is right but C5 is wrong as written; if C5 is right, C2 owes a new
refusal it currently disclaims. Either way one of the two contracts must move,
and the Failure Modes' "Silent (guarded)" section does not list this one.

**Anchors**: `0010:C2`, `0010:C5`, `0010:A10`,
`internal/graphlint/analysis.go::checkDanglingEdge`.

---

### 2 — `0010:§approach` mis-anchors `Row.Kind`, the sole prior-art basis for rejecting ALT2

The Approach argues the class must be *declared, not inferred* and grounds that
on the codebase's one existing discriminator: "`Row.Kind` (ordinary vs escape,
dump column `kind`) — is inferred from the *presence* of an `escape` list
(`internal/table/normalize.go::normalizeRule`)".

The claim is substantively right, but the inference does not live at
`normalizeRule`. It lives at `internal/table/model.go::Row.Kind`, a method
deriving the kind from `len(r.Escape) > 0`. `normalizeRule` *validates* escape
shape; it never computes the kind. `internal/table/dump.go::column` calls
`r.Kind()` for the `kind` column, matching the record's parenthetical.

**Why it matters**: this anchor is the only prior-art evidence carrying the
ALT2 rejection ("nothing in the tree infers a class from an absence"), and
`Row.Kind` is the exact contrast case — presence-inferred and per-row, versus
declared and per-model. A reader chasing the cited symbol to check the argument
lands on the wrong function. The anchor resolves `true` (the symbol exists), so
neither lint nor the projector catches this; only reading the body does.

**Anchors**: `0010:§approach`, `internal/table/model.go::Row.Kind`.

---

### 3 — `0010:C3` is silent on a present-but-empty `[rule.emit]`, against an explicit house precedent

C3 fixes the absent case ("an absent block normalizes to an empty sequence") and
says nothing about `[rule.emit]` present with zero keys.

This is not a hypothetical gap: the sibling rule in the very function that will
carry the emit arm keys on the opposite principle, and says so in a comment —
`internal/table/normalize.go:339-340`:

> "An escape rule MUST contain an `escape` list and MUST NOT contain a write
> block, clear list, or gate list, **EVEN AN EMPTY ONE: the loader keys on key
> PRESENCE, not on length** (`0002:C4`, deviations.md D2)."

An implementer following house precedent could reasonably make an empty
`[rule.emit]` a refusal; an implementer following C3's silence normalizes it to
the same empty sequence as absent. The two readings diverge at a load refusal,
and C4's "`emit` … present as `{}` — never `null`" guarantee does not
disambiguate, because it constrains the payload, not the loader.

**Why it matters**: a presence-vs-length distinction is precisely the class of
under-specification `0002:C4` already had to settle once with a recorded
deviation (D2). S3's scenario list covers "unordered keys, duplicate keys, a
non-string value, an absent block" — the empty block is the one arm not
enumerated, so no test would catch the divergence either.

**Anchors**: `0010:C3`, `0010:S3`,
`internal/table/normalize.go::normalizeRule`.

---

### 4 — `0010:C4` / `0010:S5` leave the never-null mechanism for `emit` undetermined

C4 requires `emit` "present as `{}` — never `null`, never omitted". Step 0 found
`internal/cli/flow_resolve.go::resolvePayload` already enforces never-null by
**two different mechanisms** on different fields — construction-time
initialization for some, a post-hoc nil guard for others — and the RDR does not
say which one `emit` takes.

S5's scenario asserts `emit` "present (`{}` when unauthored, never `null`)" but
its fixture is the MVV table, in which **every** ordinary rule carries a
`[rule.emit]` block (MVV step 1: "three ordinary rules each carrying
`[rule.emit]`"). The `null` arm can only surface on a row authoring *no* emit
block, which no scenario in the Testing Strategy exercises — the assertion is
written but unreachable on its own fixture.

**Why it matters**: a `"emit": null` on the wire is a payload-contract break for
every consumer that maps the field, and it is the single most likely
implementation slip for a newly added map field. The contract states the
requirement and the test names it, but nothing pins the arm that would fail.

**Anchors**: `0010:C4`, `0010:S5`, `0010:MVV`,
`internal/cli/flow_resolve.go::resolvePayload`.

---

### 5 — `0010:A9`'s per-outcome escape consequence reaches no contract and no test

A9's Evidence carries a consequence it explicitly flags as such: "`0002:C5`
scopes rescue per outcome, so 'the otherwise row' is per-outcome — a table over
an N-outcome alphabet needs N escape rows (or one `in`-atom rule expanding to
N), which the Phase 4 authoring docs must say."

That consequence appears nowhere else. C5 and C2 do not mention it; the MVV
authors "one recognized outcome"; Testing Strategy scenario 4's escape variant
inherits that single-outcome fixture. So the N-outcome case is stated once, in
an assumption's evidence prose, assigned to a docs phase, and exercised by no
test.

**Why it matters**: this is the one authoring trap in the class that produces a
*coverage* outcome rather than a load refusal — an author who writes one escape
row for a three-outcome table gets coverage findings on two outcomes they
believed were rescued. A9's own "If wrong" arm names this failure shape. It is a
test-coverage and contract-reach gap rather than an unrecorded silence, so it
ranks below 1–4.

**Anchors**: `0010:A9`, `0010:MVV`, `0010:S4`.

---

### 6 — `0010:A4`'s repo-wide count is off by two (minor, non-load-bearing)

A4 states "the 75 further repo-wide matches are all under
`docs/rdr/0002-*/evidence/`". Re-measured:

```
grep -rl '^\[dump\]' --include='*.toml' . | wc -l                        → 176
grep -rl '^\[dump\]' internal/table/testdata/ | wc -l                    → 103
… outside internal/table/testdata                                        →  73
```

**73**, not 75, and the "all under `docs/rdr/0002-*/evidence/`" attribution is
correct (2 + 10 + 15 + 43 + 3 across the spike subdirs). The load-bearing count —
103 fixtures under `internal/table/testdata/`, all 103 with an explicit
`order =` — reproduces **exactly**.

Step 0 separately notes that those excluded fixtures become non-loadable once
`emit` joins `dumpColumns`, and that leaving 0002's spike evidence stale is
defensible but unstated in C3/S6.

**Why it matters**: little, on its own — the number is decorative and the
decisive count is right. It is recorded because A4's Evidence presents both
counts in the same breath as "**Count confirmed exact**", and one of the two is
not.

**Anchors**: `0010:A4`, `0010:S6`.

## Step 3 roll-up

| # | Class | Lands on | Severity |
|---|---|---|---|
| 1 | (c) internal contradiction | `0010:C2` × `0010:C5` × `0010:A10` | blocking |
| 2 | (a) codebase claim mis-anchored | `0010:§approach` | correctness of cited evidence |
| 3 | (d) silence → missing requirement | `0010:C3`, `0010:S3` | spec gap |
| 4 | (d) silence → missing requirement | `0010:C4`, `0010:S5` | test-reach gap |
| 5 | (d) silence → unreached consequence | `0010:A9` | test-reach gap |
| 6 | (a) count drift | `0010:A4` | minor |

No finding of class (b) — Step 0's inverse check established that no sibling
path already makes any of the four decisions this RDR introduces.

No codebase claim was REFUTED or NOT-FOUND: the record's grounding is sound, and
every finding above is a contract/silence defect rather than a false frame.

---

## Resolution (same pass)

| # | Disposition | Origin ledger entry | Section touched |
|---|---|---|---|
| 1 | **fixed** | finding 1 (C2 × C5 contradiction) | `0010:C5` (override scoped to the root arm; terminal arm explicitly live), `0010:A10` (claim split by arm + executed evidence), `0010:C2` (refusal categories corrected against the spike), `0010:S4` (stray-`terminal` control) |
| 2 | **fixed** | finding 2 (mis-anchored `Row.Kind`) | `0010:§approach`, `0010:ALT2` — re-anchored to `internal/table/model.go::Row.Kind`, per-row nature made explicit |
| 3 | **fixed** | finding 3 (empty `[rule.emit]` silence) | `0010:C3` (present-but-empty pinned, with the reason it diverges from the escape-block precedent), `0010:S3` (arm added) |
| 4 | **fixed** | finding 4 (unreachable never-null assertion) | `0010:S5` — fixture must carry a row authoring no emit block; assertion pinned on that row |
| 5 | **fixed** | finding 5 (A9 per-outcome escape unreached) | `0010:S4` — two-outcome escape control added |
| 6 | **fixed** | finding 6 (count drift) | `0010:A4` — 73, with the 176/103 subtraction written out |
| — | **dismissed-with-cite** | candidate: C2's `[initial]` clause is unenforced | Executed: refuses `malformed_accessor_binding` (`evidence/spikes/cove-dt-initial/result.md`). Claim holds; C2's category list corrected rather than its claim. |

### Fork collapsed (tiebreaker-reduction gate)

Finding 1 presented as an either/or — C2 moves or C5 moves. Collapsed on
evidence rather than escalated: `0006:C3` carries "dangling edge" and "declared
terminal/escape handling" as **distinct** invariant classes, and `0006:C18`
names only the absent root, so C5's override has no basis to reach the terminal
arm. The arm additionally cannot false-positive over a zero-owned model (every
terminal predicate necessarily reads a non-owned tag), so it is not vacuous.
**C5 moves; C2 stands.** No tiebreaker escalated.

### Mini-checks fired

Read per `$RDR_HOME/stages/05-prelock.md`. Cue present → table written into the RDR.

| Mini-check | Cue | Table |
|---|---|---|
| source-authority census | three class-keyed sites that "must agree"; declared-vs-inferred class | `Authority` (Technical Design) |
| test-discriminability | MVV step 3 is an exit-0 / findings-`[]` absence oracle | `Oracle` (Validation) |
| desk trace | five contracts bearing on the MVV end-state | `Trace` (Validation) — surfaced the C2 × C5 CONTRADICTION at step 2′, now fixed |
| round-trip / fidelity | not fired — no import/export or inverse operation | — |
| disposition | not fired — refusals are per-clause, not a disposition over input classes | — |

### Needs (re)verification — carried to Stage 6

- **A12 (new, Pending)** — splitting `graph-dangling-edge` by arm is implementable
  without a taxonomy change, and no consumer keys on the code alone. Method: MVV
  Test. Booked because C5's fix introduces a load-bearing claim about how the
  existing finding set may be re-scoped; the arm split is executed and proven
  (`evidence/spikes/cove-dt-terminal/`), the *consumer* side is not.
- No previously Verified assumption was invalidated. A10's claim changed shape
  but its correction ships with executed evidence, so it stays Verified.

### Charted to successor

None. Every finding landed inside this RDR's existing contract surface; no
net-new scope was absorbed.
