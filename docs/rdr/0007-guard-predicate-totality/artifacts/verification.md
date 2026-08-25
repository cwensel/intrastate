# Verification — RDR 0007 Guard predicate totality

Stage-8 independent verification. Entries are appended by phase; each
`FAIL-N` names the REQ it violates, the exact failing input, observed vs
spec-required behaviour, and the quote violated.

## Phase 3a — CoVe (chain-of-verification)

Method: violating inputs derived from `req-list.md`, `0007`'s normative
fences, and JDR 0001 §D1/§D12/§D13 **without reading any `*_test.go`**, then
executed against the real kernel (`resolve.go`, `guard.go`,
`guardcontract.go`) via throwaway `package main` probes outside the worktree.

31 probes over the REQ set. Verified clean: the full 4x4 K3 matrix across
`all` x `unless` including `¬U = U` and `all=U ∧ unless=T -> FALSE`
(REQ-33/34/35); `F ∧ U = F` pruning across two atoms (REQ-37, REQ-67);
no short-circuit and one seam call per present-key atom (REQ-51, REQ-54);
gate ordering owned-before-unevaluable across rows (REQ-39, REQ-43);
survivor membership and the aggregation veto (REQ-64, REQ-66);
`MissingOwned` dedup/sort union (REQ-65); escape-set gating both legs
(REQ-40, REQ-41); per-atom nil-seam rule (REQ-13/14/15/16); foreign
existence token vs foreign literal (REQ-24, REQ-25); provenance-blind
presence with exact key equality (REQ-27, REQ-28, REQ-29); every
undecidable row present with no representative (REQ-52, REQ-53); pruned
rows reporting nothing (REQ-58); §D13 set literals compared as opaque
bytes (REQ-56); observed-tag substitution to a specific verdict (REQ-75).

Two genuine violations follow.

### FAIL-1 — `BlockMatch` atoms are evaluated into the `all` conjunction

**Violates:** REQ-33/REQ-35 (`0007:C6`), REQ-1 (`0007:C1`), and the
`req-list.md` ASSUMPTION "`BlockMatch` atoms are NOT evaluated by the guard
pipeline in this RDR".

**Quote violated** (`0007:C1`): "A candidate row carries its guard as a slice
of parsed atoms — key, operator token, literal, **block ∈ {all, unless}**".
(`0007:C6`): "The row verdict is `all_result ∧ ¬(unless_conj)`". JDR 0001
§D12: "0007's per-atom payload may therefore carry `match`, which no refusal
names — **harmless, since match atoms are never unevaluable**".

**Exact failing input:**
```go
Row{RuleID: "M", SourceLocator: "LM", Outcome: "OUT",
    Writes: []Tag{{Key: "w", Value: "1"}},
    Guard: []GuardAtom{
        {Key: "ok",   Operator: OpExists, Literal: LiteralTrue, Block: BlockAll},
        {Key: "gone", Operator: "eq",     Literal: "v",         Block: BlockMatch},
    }}
// Input{Flow: "f", Recognized: "OUT", Observed: []Tag{{Key: "ok", Value: "1"}}}
```

**Observed:** `REFUSE guard_unevaluable`, payload entry `blk="match"`.
Removing only the `BlockMatch` atom yields a PLAN. A second case — a lone
`BlockMatch` atom the seam decides FALSE — **prunes the row** to `no_match`.

**Spec-required:** the row verdict is a function of the `all` and `unless`
blocks only, so the `BlockMatch` atom must not contribute an operand; the
row should PLAN. §D12 asserts a `match` entry in the payload is "harmless"
precisely because match atoms cannot drive a refusal — here one does.

**Mechanism** (`guard.go::evaluateAtoms`): the block switch routes
`BlockUnless` explicitly and sends **everything else** to `default:`, which
folds into `allResult`. `BlockMatch` — and any unrecognized `Block` value,
e.g. `Block("garbage")`, separately confirmed — therefore joins the
conjunction. The kernel is the fail-closed backstop for a normalizer that
did not enforce §D6's load-time match restriction (the pattern `0007:C3`
establishes for foreign literals), so this path is reachable by contract,
not merely by malformed input.

### FAIL-2 — payload order is a function of table row order, not the input tuple

**Violates:** REQ-54 and REQ-55 (`0007:C8`).

**Quote violated** (`0007:C8`): "the kernel ... MUST sort the payload so it
is a function of the input tuple (RDR 0001 REQ-1), **never of atom or row
order**." And: "**The sort key MUST be TOTAL over payload entries.** ... Two
entries would then tie ... leaving the order to an unstable tie-break — the
atom-order dependence this clause exists to forbid. The ordering is
therefore the tuple `(RuleID, SourceLocator, key, block, operator token,
literal)`, compared field by field in that order."

**Exact failing input:** two or more undecidable rows sharing
`(RuleID, SourceLocator)` but carrying different atoms.
```go
d := func(k string) Row {
    return Row{RuleID: "SAME", SourceLocator: "SAME", Outcome: "OUT",
        Guard: []GuardAtom{{Key: k, Operator: "eq", Literal: "v", Block: BlockAll}}}
}
// Table A rows: d("k01") … d("k14")   (forward)
// Table B rows: d("k14") … d("k01")   (reverse)
// Input{Flow: "f", Recognized: "OUT"}  -- no seam, all keys absent
```

**Observed:**
```
forward table order -> [k01 k02 k03 … k14]
reverse table order -> [k14 k13 k12 … k01]
payloads equal: false
```
Also reproduced on the **escape set** (same `gate` delegation) with two
escape rows sharing `RuleID`/`SourceLocator`.

**Spec-required:** identical `Refusal.Undecided` for both orderings — the
six-tuple is fenced as total over payload *entries*, so the atom half must
break a row-identity tie.

**Mechanism:** the six-tuple is split across two independent sorts.
`resolve.go::gate` sorts rows with `compareUndecidedRows`, which compares
`(RuleID, SourceLocator)` **only**; `guard.go::evaluateAtoms` sorts atoms
within each row. `slices.SortFunc` is not stable, so rows tying on identity
retain no defined order and the emitted payload follows `Table.Rows`
position. The spec anticipates exactly this failure mode ("Row identity then
key is NOT total") and fences the six-tuple to close it; splitting the sort
reintroduces the tie one level up, at the row.

**Note on scope:** REQ-55's "two entries equal on all six name the same
atom, which the kernel MUST NOT report twice" makes row identity a
*component* of the entry key, not a partition boundary. Whether RDR 0002
guarantees `(RuleID, SourceLocator)` uniqueness is not stated in 0007, and
`0007:C8` places the obligation on the kernel unconditionally.


## Phase 3b — adversarial failure-mode review (independent)

Independently derived; the reviewer did not read the `FAIL-N` entries above.
Anchored to `## Trade-offs > ### Failure Modes` of the record, read through the
projector. Tests live in `internal/resolve/guard_adversarial_0007_test.go`. No
existing test was modified or weakened — the file is additive.

Method note: the per-clause suite (`guard_totality_test.go`,
`guard_atoms_test.go`, `guard_mvv_test.go`) is green and genuinely strong. The
K3 matrix is exhaustive across `all`, `unless`, and cross-block; payload sorting
survives 300 permuted runs at scale under `-race`; §D13 set bytes cross the seam
verbatim and unparsed; operator-token drift (`"Exists"`, `"EXISTS"`, a future
token) fails closed exactly as `0007:C3` requires. Nine further probe vectors
(escape-path payload completeness, pruned-escape owned obligations, out-of-range
seam verdicts, typed-nil seams, contract-test strength against a folding seam,
duplicate row identity) all passed and are recorded here as closed. The three
findings below are the composition cases that survived.

### ADV-1 — a `match`-block atom is evaluated as a guard operand

- **Failure mode**: `0007:F4` — "The domain rule is scoped to guards, so the
  match pattern still folds absence into non-match … predicate PLACEMENT decides
  whether an absent key refuses non-escapably or escapes."
- **Also implicates**: `0007:C6` (verdict formula fenced over `all` and
  `unless` only), `0007:C8` (payload names what blocked the GUARD verdict).
- **What breaks**: §D12 put `BlockMatch` on the same `Block` type — "one type,
  no separate slice on `Row`" — which makes a match atom *representable* in
  `Row.Guard`. `guard.go::evaluateAtoms` has no `match` case: its `default:`
  arm folds every block it does not name into `allResult`. A match atom over an
  absent key therefore refuses `guard_unevaluable` — the NON-escapable kind —
  where F4 states that match-pattern absence yields the escapable `no_match`.
  The match atom also appears in the guard payload with `Block:match`, naming a
  key that was never a guard input.
- **Test**: `TestAdv0007_1_MatchBlockAtomIsNotAGuardOperand` (2 subtests).
- **Status against current implementation**: **FAILS** (both subtests).
- **Suggested resolution**: give `evaluateAtoms` an explicit `case BlockMatch`
  that contributes no operand (the way `0007:C5` drops the omitted `unless`
  term) and emits no payload entry.

### ADV-2 — the `Block` boundary fails OPEN, reopening the `xg7p` masking path

- **Failure mode**: `0007:F8` — "Evaluator/grammar version skew … the seam
  answers unevaluable" — read against `0007:C3`'s fail-closed drift rule and
  `0007:C11`'s ratification of D8.
- **What breaks**: `0007:C3` makes OPERATOR drift fail closed in both
  directions, and the kernel honours that exactly. The BLOCK boundary has no
  equivalent rule and fails OPEN: `default:` sweeps any unrecognized block into
  `allResult`, so an atom there can be DECIDED, and a decided-FALSE one PRUNES
  its row. `0007:C11` ratifies pruning "conditional on the domain rule: pruning
  is safe exactly because GuardFalse can only arise from decided atoms" — a
  condition stated over *guard* atoms, which an atom in an unfenced block is
  not. The consequence is the Background probe verbatim: FALSE guard + absent
  `RequiresOwned` key + modeled `no_match` escape → `Escaped:true`, missing
  artifact state masked. This is kata `xg7p`, reproduced on the post-RDR kernel
  through a boundary the fail-closed clause never reached.
- **Reachable via**: `BlockMatch` (§D12), the zero value `Block("")` (an unset
  field on a producer-built atom), and any future block token.
- **Test**: `TestAdv0007_2_UnfencedBlockFailsOpenAndReopensTheMaskingPath`
  (3 subtests).
- **Status against current implementation**: **FAILS** (all three legs; each
  emits `Escaped:true` where `owned_state_unavailable` naming `gate` is owed).
- **Suggested resolution**: state the block-boundary drift rule the way
  `0007:C3` states the operator one, and make the kernel fail closed —
  an atom whose block the verdict formula does not name must never DECIDE a
  row. Same edit as ADV-1 if the `default:` arm is replaced by explicit cases
  plus a fail-closed fallback.

### ADV-3 — a duplicated tag key makes the guard verdict positional

- **Failure mode**: `0007:F6` — "Refusal defeated by caller-supplied state
  (A13) … an operator can supply the missing key as an observed tag."
- **Also implicates**: `0007:C8` ("MUST sort the payload so it is a function of
  the input tuple … never of atom or row order") and RDR 0001 REQ-1.
- **What breaks**: `0007:C8` closes atom order and row order. It does not reach
  the third ordering the kernel is exposed to — the caller's TAG slices.
  `resolve.go::assemble` resolves a key repeated within one provenance by "last
  tag wins", and its doc comment then claims "the merge is order-insensitive
  across provenances, which is what value-level replay determinism requires
  (REQ-3)". That claim holds only while no key repeats, and A13 leaves
  caller-supplied `Observed` unconstrained. The two orderings of ONE duplicated
  observed key give opposite guard verdicts, which under D8 is the difference
  between a pruned row and a survivor — and therefore between the masking plan
  and the honest refusal. Observed: `[gate=open, gate=closed]` → `Escaped:true`
  plan routing around absent owned `audit`; `[gate=closed, gate=open]` →
  `owned_state_unavailable`. Same tag multiset, opposite disposition. A
  duplicated OWNED key is equally positional (a producer defect, but the kernel
  is the component that must not let slice order decide a verdict).
  `0007:F6` accepts that an operator can SUPPLY a missing key; it does not
  extend to an operator flipping the disposition by REORDERING tags they
  already supplied.
- **Note on scope**: `assemble` is RDR 0001 code, but the defect surfaces
  through `0007:C4` key identity and `0007:C8` determinism, and its observable
  consequence is this RDR's own masking probe. Recorded here rather than
  edited.
- **Test**: `TestAdv0007_3_DuplicateKeyMakesTheVerdictAFunctionOfSlicePosition`
  (3 subtests).
- **Status against current implementation**: **FAILS** (all three subtests).
- **Suggested resolution**: make a duplicated key resolve by a rule that is a
  function of the tuple — reject it as a malformed input tuple, or make the
  collision itself unevaluable. Either way the two orderings must agree.

### Vectors probed and found CLOSED (no test added)

| Vector | Anchor | Result |
| --- | --- | --- |
| K3 matrix across `all`, `unless`, and cross-block, incl. `¬U = U` | `0007:C6` | correct in all 9×3 cells |
| `F ∧ U = F` with the FALSE arriving from either block | `0007:C6`, `0007:C11` | correct |
| Payload determinism: 300 permuted runs, 6 rows × 5 atoms, `-race`, `-count=3` | `0007:C8` | stable |
| §D13 set bytes — non-canonical value crosses the seam verbatim, unparsed | JDR 0001 §D13 | correct |
| Two `in` atoms differing only in member order are distinct payload entries | `0007:C8` REQ-56 | correct |
| Operator drift: `"Exists"`, `"EXISTS"`, `""`, a future token, over an absent key | `0007:C3` | all fail closed, reason `absent`, seam never consulted |
| Escape-path payload completeness — no short-circuit on a decided block | `0007:C7`, `0007:C8` | correct |
| Pruned escape row's owned obligation does not leak | `0007:C11` D8 | correct |
| Escape row both unevaluable and missing owned state | `0007:C7` | owned reported first, `Rows` correctly scoped |
| Escape-set aggregation veto (one TRUE escape, one UNEVALUABLE) | `0007:C10` | vetoes correctly |
| Escape row declaring `guard_unevaluable` in its `Escape` list | `0007:C8` REQ-44 | inert; cannot rescue |
| Out-of-range `GuardResult` from the seam | `0007:C1` | mapped to unevaluable |
| Cross-provenance key collision — owned wins the value, presence is blind | `0007:C4` | correct |
| `MissingOwned` union, deduped and sorted, under permuted rows | `0007:C10` REQ-65 | correct |
| Duplicate row identity in the table | `0007:C8` | payload consistent with `Rows` |
| `TestGuardEvaluatorContract` against a seam folding unparseable → FALSE | `0007:C1` | contract test correctly fails it |


## Phase 3c — fixup dispositions

Every Phase 3a `FAIL-N` and Phase 3b `ADV-N` above is **RESOLVED**. The
probed-and-CLOSED vector table needed no action and was not revisited.
Full suite green: `go build ./...`, `go test ./...`,
`go test -race ./internal/resolve/`, `golangci-lint run` (0 issues). Each of
the three commits builds independently.

No test was weakened. The three Phase 3b adversarial tests were left
byte-for-byte as written and now pass; two regression files' worth of new
coverage was added for the Phase 3a findings, which had none.

| Finding | Status | Fix | Commit | Regression test |
| --- | --- | --- | --- | --- |
| FAIL-1 | **RESOLVED** | fail-closed block boundary | `433f364` | `TestFix0007Fail1_MatchBlockAtomContributesNoOperand` |
| FAIL-2 | **RESOLVED** | total payload sort | `6047b09` | `TestFix0007Fail2_PayloadSortIsTotalOverRowsTyingOnIdentity` |
| ADV-1 | **RESOLVED** | fail-closed block boundary | `433f364` | `TestAdv0007_1_MatchBlockAtomIsNotAGuardOperand` (as written, now green) |
| ADV-2 | **RESOLVED** | fail-closed block boundary | `433f364` | `TestAdv0007_2_UnfencedBlockFailsOpenAndReopensTheMaskingPath` (as written, now green) |
| ADV-3 | **RESOLVED** | duplicated key is unevaluable | `a87fc9e` | `TestAdv0007_3_DuplicateKeyMakesTheVerdictAFunctionOfSlicePosition` (as written, now green) |

### FAIL-1 + ADV-1 + ADV-2 — RESOLVED (one root cause)

**Fix** (`guard.go::evaluateAtoms`, `guard.go::isGuardBlock`): the unfenced
`default:` arm is replaced by explicit `BlockAll` and `BlockUnless` cases
plus a fail-closed fallback. An atom whose block `0007:C6`'s verdict formula
does not name contributes **no operand** to either conjunction — the way an
omitted `unless` block contributes none (`0007:C5`) — emits **no payload
entry** (`0007:C8`'s payload names the atoms that blocked the GUARD verdict,
and an atom contributing no operand cannot have blocked it), and is **never
handed to the seam**. This is the posture `0007:C3` already sets for
operator-token drift: the operator boundary fails closed, so the block
boundary does too.

**What now happens on each reported vector.** The FAIL-1 vector verbatim
(`exists ok = true` in `all` beside a `BlockMatch eq` atom over absent
`gone`) yields a PLAN naming `M`, `Escaped:false`. A lone seam-decided-FALSE
`BlockMatch` atom no longer prunes its row. All three ADV-2 legs
(`BlockMatch`, `Block("")`, `Block("any")`) refuse
`owned_state_unavailable` naming `gate` — the honest disposition, not the
`Escaped:true` plan. No payload entry carries a non-guard block. Q1 reading
(a) is unchanged: `BlockMatch` is still a constant on the `Block` type and
`Row`'s match pattern is untouched (`TestReq78` still holds).

Recorded as deviation **D7** (SPEC-UNDER) with its evidence: `0007:C6`'s
exhaustive verdict formula, `0007:C5`'s "no operand, not a neutral one",
`0007:C11`'s pruning condition stated over GUARD atoms, `0007:C3`'s
fail-closed posture, and JDR 0001 §D12's own "harmless, since match atoms
are never unevaluable" — a claim that only holds under this fix.

### FAIL-2 — RESOLVED

**Fix** (`guard.go::compareUndecidedRows`): the row comparator breaks an
identity tie on the atom half, entry by entry over each row's already-sorted
atom list, then by atom count — making the sort key total over payload
entries as `0007:C8` fences it, with row identity a COMPONENT of the entry
key rather than a partition boundary. `Reason` is compared alongside each
atom because it is part of the emitted entry, though it is not one of
REQ-55's six ordering fields.

**Verified with permuted-input runs**, not a single ordering: the CoVe
vector verbatim (14 rows sharing `RuleID`/`SourceLocator`, one distinct
absent-key atom each, forward and reversed) now yields identical
`Refusal.Undecided` on the candidate set AND on the escape set, which
reaches the same `gate` by delegation. A third leg permutes rows that tie on
identity AND on atom key, so they differ only in block/operator/literal.

Recorded as deviation **D8** (IMPL-DECISION, mechanical translation).

### ADV-3 — RESOLVED (fixed at the kernel boundary only)

**Fix** (`resolve.go::assemble`, `resolve.go::merge`,
`resolve.go::TagSet.conflicting`, `guard.go::evaluateAtom`): a key supplied
more than once WITHIN one provenance with differing values is marked
**conflicted**. It stays PRESENT — presence is provenance-blind and
non-positional (`0007:C4`), and an existence atom still decides from
presence alone, because a collision is about the VALUE, not about whether
the key arrived — but the view carries no single value for it, so a
value-comparing atom over it is `GuardUnevaluable` with reason
`uncomparable`: "the key was present and its value was not compared to a
verdict" (`0007:C8`). The seam is never handed one of the colliding values.

**The two orderings now agree.** `[gate=open, gate=closed]` and
`[gate=closed, gate=open]` both refuse `owned_state_unavailable` naming the
absent owned `audit` — neither masks it behind an `Escaped:true` plan. The
duplicated-OWNED-key leg agrees the same way.

**RDR 0001's merge semantics are not rewritten.** Identical repeats are not
conflicts; cross-provenance precedence (owned > observed > recognized) is
untouched and is a property of the tuple, not of slice order; and a permuted
NON-duplicate input assembles a byte-identical view, so REQ-3 replay
determinism for the non-duplicate case is unchanged. Verified by direct
probe over `assemble` before it was removed.

**Why unevaluable rather than rejection**, grounded before the choice: RDR
0001 reserves the Go error return for programmer mistakes, and A13 makes a
duplicated observed tag caller INPUT rather than a programmer mistake; a
sixth refusal kind is closed by REQ-7 and by `0007:C7`'s refusal to reopen
the taxonomy for a better-motivated case. Unevaluable is the RDR's own
posture — `0007:C2` ("never to false and never to true"), RDR 0004 A8
("Indeterminate MUST be a refusal-class result"), SCXML §5.9.1 and SQL:2003
strict routines — and `0007:C8`'s closed reason set already carries the
exact word. Full argument and citations on deviation **D9** (SPEC-UNDER).

### New public surface

**None.** Every symbol this phase adds is unexported: `isGuardBlock`,
`TagSet.conflicting`, `TagSet.merge`, and the `taggedValue.conflicted`
field. `TagSet.Lookup`'s exported three-value shape, fenced by `0007:C1`'s
presence test, is unchanged, as is every other exported symbol D6
enumerates. No SPEC-UNDER entry is owed on the ADDITIVE-IS-NOT-EXEMPT gate
for new surface; D7 and D9 are SPEC-UNDER for the contract gaps they close,
not for surface.

### Open author decisions

**None.** Both contract-level gaps (D7, D9) are resolved by the spec's own
evidence base, with the evidence recorded on each entry.
