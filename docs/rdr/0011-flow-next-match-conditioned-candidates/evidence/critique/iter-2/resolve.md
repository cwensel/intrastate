Model: claude-opus-5[1m] (dispatcher; pass/diff contexts stamped in their own files)

# Critique resolve — RDR cli/0011, iter-2

Origin ledger: the reconciled table in `diff.md` (18 distinct defects over
pass A's 12 rows and pass B's 9, reconciled by passage anchor). Dispositions
below trace to that table's numbering.

**Outcome: ROUTE-BACK.** Two findings exceed what this half may absorb by
editing — each refutes a load-bearing clause of `C1` on shipped code, and
fixing either re-opens the approach rather than the wording. Per
rdr-common §punt-ledger both are recorded in
`docs/rdr/0011-flow-next-match-conditioned-candidates-postmortem.md` at the
moment of the route-back, before any refine can collapse the history.

## Blocking (force the route-back)

- **#1 `C1` per-key filter is unlicensed** — `needs-propose`. C1 permits its
  filter because "a single probe row holds at most one match tag per key"
  (`:140`, citing `0002:C13`'s row-level `in`-expansion). Refuted on `main`:
  `internal/table/normalize.go::atomsFromBlock` nests a loop over OPERATORS
  inside its loop over keys, and `:78` admits both `eq` and `in` under a match
  block, so one authored `[rule.match.<key>]` carrying both emits TWO match
  tags on that key. Independently reproduced by the diff context (own probe
  model → `table.Load` accepts → `KernelRow().Match` carries two tags on one
  key). Consequence: C1's per-key filter keeps such a row a candidate when the
  key is ABSENT and excludes it unreported when the key is PRESENT at any
  value — supplying state REMOVES rows, the non-monotonicity this RDR exists
  to fix. Touches `C1`, `A12` (was `Verified`), `A3`, `S3`, `S4`, `F1`.
  Flagged in place at the `C1` boundary; `A12` flipped to `Pending`.
- **#2 `C1`'s "flow resolve unchanged" is false at the verb** —
  `needs-resolve`. `internal/cli/flow_exec.go::invokedReaders` has two callers,
  `flow_next.go:97` and `flow_resolve.go:98` (verified directly), and C1 writes
  the demand-set extension as a property of the MODEL, not the mode.
  `runReaders` aborts the verb on an unbound artifact role (exit 2) or a reader
  refusal (exit 3), so a model with a match-only owned key served by a refusing
  or unbound reader turns a passing `flow resolve` into a refusal. C1's "This
  contract changes nothing in `internal/resolve`: … `flow resolve`'s
  dispositions are unchanged in every case" (`:182`) conflates a true claim
  about the PACKAGE with a false one about the VERB. Compounded by #3.
  `A15`(iii) corrected in place — it did name the second caller but framed the
  uniformity as benign without reaching `runReaders`' refusal arms.

## Status-consistency (rides the route-back)

- **#3 `A15` `Pending` while six passages state it as settled** — `fixed`
  (flagged, not closed). C1's MUST body, `JC1`, the `authority` mini-check,
  `S8`, `MVV` 8 and Phase 1 all assert the demand-set extension as fact. This
  is the breach `§Assumption Verification` forbids. Not closable here: #2
  changes what A15 must claim. Recorded on the Status qualifier.

## Confirmed, in-draft absorbable on the way back through

Not edited now — the re-propose owns the clauses they land in, and editing
them against an approach that is itself being re-decided would be churn.

- **#4 `A6`'s 21→3 witness reproduces the declared alphabet** (A C-6,
  CONFIRMED, generalizes to every non-terminal stage). The three survivors are
  one rule per declared outcome — the alphabet `outcomes[]` already carried, so
  the motivating caller is no better able to choose. A problem-statement
  honesty edit, and it bears on whether the Approach's headline win is real.
- **#5 `--all` prescribed as the diagnostic but forbidden from naming why**
  (A C-8/C-10, B C-4). Already stated honestly in `F1`; the finding is that the
  RDR predicts its own amendment.
- **#6 `A11` true by call-graph coincidence, no oracle pins it** (B C-3) —
  the best unique contribution of pass B.
- **#7 filter-after-`KernelRow()` rule lives only in prose** (B C-6); no
  set-kinded fixture discriminates the two implementations.
- **#8 no multi-key mixed-state fixture** (B C-8) — adjacent to #1 and
  partially subsumed by it.
- **#9 rename bundled with the predicate change** (A C-10, B C-5).
- **#10 `TestReq44` survives by fixture accident** (A C-4 first half).
- **#11 ~15 `unresolved` comment/message sites uninventoried** (A C-5 second
  half) — C3's census counts test READS only, which it says; the finding is
  that the production sweep needs the rest.

## Dismissed-with-cite

- **#12 "C3 counts `:103` among the three needing ok-bool"** (A C-5 first
  half) — **REFUTED** by the diff context: C3 names its three explicitly and
  `:103` is not among them. Not edited.
- **#13 "A4's BREAKS 0 cited as suite clearance"** (A C-4 second half) —
  **REFUTED**: fenced three times (A4's own Evidence, C3, S6). Not edited.
- **#14 `A3` states the invariant one arm short** (A C-9) — **PARTIAL**: the
  0-reader arm does fire, but A3's `!= 1` already covers it. Editorial at most.
- **#17 permanent `next`/`resolve` disagreement on absent-match rows**
  (B C-9) — **RE-RAISE**. Decided in `§Decision Rationale` (why `0007:REQ-78`
  stays deferred) and recorded in `§Consequences`. The fix half must not edit
  the draft to satisfy it (grounding gate 3).
- **#16 record size vs. change size** (A C-11, B C-2) — soft; B's framing
  unanchored. `Profile: large` is set by the contract count and already
  justified in `§Proportionality`.
- **#15 `excluded()` refactor bigger than framed** (B C-7) — PARTIAL; `S7`
  exists to pin it. Rides the re-propose since #2 changes that function's
  contract anyway.
- **#18 text mode doubles via `flatten`** (A C-7) — not independently
  re-verified by the diff; already open as 3amigo iter-3 L7
  (`D-undecided-reporting-shape`), so it is not net-new.

## Needs (re)verification (Stage 6 closes these)

- `A12` — flipped `Verified` → `Pending`. Owes: which loader tier (if any)
  refuses an `eq`+`in` pairing on one match key, and what the predicate does
  when a key carries several atoms.
- `A15` — stays `Pending`, scope corrected. Owes: the exit-2/3 arms over a
  match-only-owned-key fixture, and the scoped-vs-uniform decision for the
  demand-set term.
- `C1`, `JC1` — named on the Status qualifier; both rest on clauses #1/#2
  refute.
- `A3` — not flipped, but #1 lands in the same reachability territory; re-read
  it when C1's predicate is re-decided.

## Charted to successor

None net-new this pass. Every confirmed finding lands inside a clause this
RDR already owns; `Charted.md` from the previous critique cycle still holds
its own item and is untouched.
