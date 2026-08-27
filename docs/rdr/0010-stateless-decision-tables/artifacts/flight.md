# Flight — batch:rdr-0010

`/kata-flight --label batch:rdr-0010 --drain`, run as the `§land-flight` step of
`/rdr-implement-triage 10 --close-and-flight`.

```
flight: 2 shipped, 0 stopped, 0 skipped [over 1 wave]
  shipped: stn1=e744c24  e8dt=c87a120
  stopped: (none)
  skipped: (none)
```

Wave 1 resolved 2 eligible katas. `bycb` was filtered out by the standing
eligibility rule (`kind:rdr-seed` — it exits to RDR authoring, and kata-ship
gate 4 refuses a seed); it remains open under `batch:rdr-0010` by design.

**Drain terminated after one wave.** No `KATA_PUSH` spin-off was minted by
either ship, so nothing new was tagged onto the standing selector and a
re-sweep could only resolve empty (skip-the-guaranteed-empty-sweep rule).

## Review gate (ran before the ship loop, per default)

Both katas passed `/kata-scope-review` as `IN-SCOPE` with source-grounded plans
embedded. The review produced two corrections that changed the shipped work —
in both cases the kata's own text was wrong and grounding caught it:

1. **`stn1`'s blocking open question was FALSE.** `verification.md`
   §Undecidable claimed no runtime gate-deny witness was authorable on a
   decision-table fixture. Disproved by live probe: `0010:C2` forbids only
   *accessors* keying an **owned** tag, while `load.go::accessorTable` imposes
   no owned requirement on **gates** and `checkAccessorBindings` arity-checks
   readers/writers only. A gate keying an **observed** tag is legal. Deny →
   `flow-gate-denied` exit 2; allow → `"gates":[{"id":"ok","result":"allow"}]`.
   Both halves were therefore witnessed on the decision-table class, and the
   correction was appended to `deviations.md` (the RDR body is never amended).

2. **`e8dt`'s own suggested wording was wrong.** Its body proposed text asserting
   an `unless` row "covers the three cells outside `free-eu`". A two-row fixture
   failed with `leaves 1 of 4 assignments uncovered`: `unless` terms are
   intersected then subtracted (`internal/guard/product.go:624-631`), so
   single-atom `unless.tier=free` removes the whole `tier=free` slab — 2 cells,
   not 3. The *thesis* (rows ≠ cells) held; the arithmetic did not. Corrected
   wording shipped instead, and a second instance of the same conflation at
   line 30 ("rows **are** cells") was folded into the same commit.

## Ships

### stn1 → `e744c24` — `test(0010): tighten the vacuous and quadrant-incomplete oracles`

4 files, +415/−47. Six weak oracles tightened; **production code untouched**.

Because the kata is about oracle *strength*, a green run proves nothing, so each
tightened assertion was **falsified against an injected defect** and reverted
before commit. All seven injections confirmed; two are the direct proof of the
alleged defect — the old oracle **passed** where the tightened one **failed**:

| injection | tightened | old oracle |
|---|---|---|
| `gateResult` serializes `verdict` not `result` | FAILED | **PASSED** |
| fixture loses `gate=["ok"]` (vacuity guard) | FAILED | — |
| deny code → `flow-gate-indeterminate` | FAILED | — |
| escape join takes another escape row's emit | FAILED | — |
| resolve overwrites `dt.toml` in place | FAILED | **PASSED** |
| emit dropped on sm escape rows | FAILED | — |
| clear walk `continue`s on non-owned key | FAILED | — |

Refine: 0 findings (jobs 6328, 6329). `make check` green.

Fixture gotcha now recorded in the test: `clear` is a rule-level key and must
precede the first `[rule.*]` sub-table — appending it after `[rule.guard.all.b]`
authors a *guard key named `clear`* and refuses `malformed_predicate_atom`
instead of `write_to_non_owned_tag`.

### e8dt → `c87a120` — `docs(model): state decision-table completeness in cells, not rows`

1 file, +5/−3, documentation only. Both fixtures re-run against the freshly
built worktree binary as acceptance evidence:

- **three-row** table over the 2×2 product → `{"type":"ok","data":{"findings":[]}}`,
  clean exit 0, no escape row — proving rows can span cells.
- **two-row** table → exit 2, `graph-coverage-gap: leaves 1 of 4 assignments
  ... uncovered` — proving the discarded claim really was wrong.

Refine: 0 findings (job 6331). `make check` green.

## Notes

- **No 0011 diff-guard change was needed by either kata.** `TestReq120`'s
  allow-list already admits `*_0010_test.go` and `docs/rdr/`, and
  `docs/model-authoring.md` is admitted via the D15 doc-path allowance
  (`deviations.md:503-505`). The guard was neither weakened nor extended.
- **Pre-existing harness condition, independently reproduced by both resolve
  agents** (verified by stashing to a clean tree, so not a defect of this work):
  three `internal/cli` diff guards — `TestReq94And119_MVV9…`,
  `TestReq34And107And130_…`, `TestReq120_…` — fail with *"diff base is HEAD"* when
  a worktree branch carries **no** commit. That is `f990bee`'s intended
  fail-closed behaviour; it clears as soon as the branch has one commit. Worth
  knowing these three cannot pass in an uncommitted worktree.

## Still open under this batch

- `bycb` — `kind:rdr-seed`, "REQ-54's non-coverage invariant arms are unexercised
  over the decision-table class". Not flight-drainable by design; it routes to
  RDR authoring (`/rdr-seed-triage` → `/rdr-seed`). Carries an `## Open question`
  on whether the node-ceiling arm is reachable at all for this class (a ∅ root
  means one node), which should be settled before fixtures are authored.
