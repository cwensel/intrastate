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

