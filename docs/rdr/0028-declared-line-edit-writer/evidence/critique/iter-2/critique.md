Model: claude-sonnet-5

## Delta-scoped critique re-run — RDR 0028, iteration 2

Scope: A2, A10, A11, A12, C1, MVV, S25, S28, S31, §prerequisites,
§capability-dependencies, §existing-infrastructure-audit,
§risks-and-mitigations, §pre-lock-mini-checks, §performance-expectations,
G-proportionality.

### Per-element check notes

- **A2** (Verified→Pending flip): evidence re-checked against
  `internal/accessor/executor.go::Executor.Write` (lines 302–339). The
  described ordering — `readerFor` resolved before `Apply`, `protectedKeys`
  computed, the pre-`Apply` baseline `invokeRead` gated on `len(protected) !=
  0`, and its failure swallowed into `baselineUnread` while the write
  proceeds — matches the source exactly. The "AUTHORITATIVE DETECTOR" framing
  in C1.3 `read-back:` is consistent with this evidence: the baseline arm only
  fires when the reader declares a key beyond the planned set, so it is
  correctly described as not reachable in the MVV fixture. No gap.
- **A10**: unchanged reasoning, still Pending, no new prose depends on it as
  settled fact. Grounding matches `internal/cli/flowbind/registry.go`'s single
  `Registry(...)` construction site. No gap.
- **A11**: the deadline-arm carve-out is consistently hedged everywhere it is
  cited — C1.3 `order:`, the `disposition` mini-check table's timeout row, and
  the Risks section all mark it Pending and do not assert the executor-minted
  timeout's applied sense is fixed. Cross-checked against source
  (`executor.go` lines 356, 359–366): `appliedDeadline` is set from
  `applyCtx.Err()` without consulting `err`, matching the clause's claim
  exactly. No gap.
- **A12** / re-anchor IDENTITY: cross-checked against C1.3 `select:`'s
  "a deletion never shifts a sibling rule's target" and against C1.5's
  `clear` one-way semantics. `re-anchor:`'s "(or, for a deleted line, zero
  lines)" carve-out correctly composes with a `clear`'d rule. The IDENTITY
  formula (held pre-edit index shifted by preceding sibling deletions) is a
  restatement, not a contradiction, of `select:`'s invariant. No gap.
- **C1.3 `write:`** (atomicity boundary): cross-checked against the MVV's
  per-invocation shape (item 2, item 3) and S20. The "ONE ENTRY over ONE
  file, not transactional across entries" language is consistent with the
  MVV's "one line per file PER INVOCATION" framing — the record and README
  are two independent write entries, each independently atomic. No gap.
- **C1.3 `read-back:`** vs S25: S25's two-armed fixture (reader declaring
  only the planned key vs. a second key) precisely mirrors `protectedKeys`'s
  subtraction logic in source (`executor.go` lines 526–538). Grounded
  correctly. No gap.
- **S31**: cross-checked against the mismatch branch (`executor.go` lines
  458–463) — `applied` is left unset on `ClassReadBackMismatch`, confirmed
  against the explicit 0004:C14-scoping comment at lines 452–457. The fix is
  accurate. No gap.
- **§existing-infrastructure-audit** (`save` row, Extend→Sibling): cross-
  checked against `internal/cli/flowbind/flowbind.go::save` (lines 139–165).
  All three named incompatibilities confirmed verbatim: `MkdirAll(0o700)`,
  `clierr.WriteJSONLine(tmp, s)` taking a `store`, and the fixed
  `Chmod(staged, 0o600)`. No gap.
- **§prerequisites** (new BLOCKING item) and **MVV** (new step 8): cross-
  referenced against each other and against F6 — the unlinked `| NNNN |`
  row shape (0001, 0006) is named identically in both places, and
  Prerequisites' consumer-side normalize item correctly defers to "MVV step
  8 covers the refusal for that shape." No gap.
- **G-proportionality** response vs **§prerequisites** BLOCKING item: the
  "no split" conclusion (C1.6 is a dependent extension, not independently
  deployable) is orthogonal to the MVV-execution blocker (the consumer's
  projector verb doesn't exist yet) — one is a contract-count question, the
  other an execution-readiness question. No tension.
- **§risks-and-mitigations** (three new rows): each new risk (read-back
  independence UNMITIGATED, `clear = "line"` enforcement gap, A3's corpus-
  drift mechanism gap) is self-contained and does not restate or contradict
  any other edited clause. No gap.

### New finding

One genuine new gap, in the Performance Expectations / S28 pairing (item 11
of the resolved list).

C1.4 `precedence:` states: "within one entry, fail-fast in the order above
... across entries and tables, 0025:C5's rules apply unchanged." The 0025:C5
clause it cites for that inheritance says, verbatim: "ACROSS entries in the
SAME table there is no order: `accessorTable` ranges a Go map, so which of
two defective entries is reported is unspecified, **and no test may assert
it**."

The newly added S28 sibling-table case and the newly reworded Performance
Expectations "Map order" bullet both introduce a *new* rule not present in
any C1.4 normative clause: sibling `edit.<key>` tables carrying same-category
defects resolve to a deterministic winner via "tables sorted by key," and
S28 explicitly requires a test to assert this determinism ("fails if the
report varies across runs"). This is the opposite disposition from the one
C1.4 `precedence:` says it inherits unchanged for map-iterated siblings, and
the "sorted by key" rule itself appears nowhere in C1.4's normative contract
text — only in the non-normative Performance Expectations prose and in the
S28 test scenario. It is a new load-bearing behavioral claim with no
normative clause behind it, and it is not clearly reconciled with the
inherited "unspecified, no test may assert it" provision it sits beside.

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|---|---|---|---|---|
| E-1 | 0028:S28 / 0028:§performance-expectations ("Map order") vs 0028:C1 C1.4 `precedence:` | New load-bearing claim ("tables sorted by key" determinism across sibling `edit.<key>` tables) stated only in test/perf prose, not in any C1.4 normative clause, and contradicting the 0025:C5 "unspecified, no test may assert it" rule that C1.4 `precedence:` says it inherits unchanged for map-iterated siblings | A model author cannot find a normative rule stating which of two equally-defective `edit.<key>` tables lint will name; if the implementation does not actually sort by key, S28 is an untestable/flaky requirement asserting behavior the contract never promised | Fix to S28 / Performance Expectations "Map order" bullet (item 11 in the resolved list) |
