# cove — resolve pass (iteration 1)

Origin ledger: the union of `findings.md` (claude-fable-5-1) and
`findings-modelB.md` (claude-sonnet-5), keyed by element id. Pass B's F-1
and pass A's F-9 are the same defect on `0012:A4`; counted once.

## Dispositions

| Ledger | Anchor | Disposition | Section touched |
| --- | --- | --- | --- |
| L1 (A/F-1) Key Discovery "guard atoms only" REFUTED | 0012:§research-findings | fixed | Key Discoveries |
| L2 (A/F-2) second suite implementer unnamed | 0012:C3 | fixed | Normative Contracts C3 |
| L3 (A/F-3) probeRow holds no model | 0012:A3 | fixed | Critical Assumptions A3 Evidence |
| L4 (A/F-4) suite cases hardcode Key "subject" | 0012:C3 | fixed | Normative Contracts C3 |
| L5 (A/F-5) 0007:A17 headline deviation unrecorded | 0012:§approach | fixed | Approach + Key Discoveries |
| L6 (A/F-6) guard.Conforms sibling omitted | 0012:§research-findings | fixed | Sibling-path check |
| L7 (A/F-7) flow next unnamed seam consumer | 0012:F1 | fixed | Failure Modes |
| L8 (A/F-8) "write-block category" does not exist | 0012:C5 | fixed | Normative Contracts C5 |
| L9 (A/F-9 + B/F-1) valueAssignments misattributed | 0012:A4 | fixed | Critical Assumptions A4 Evidence |
| L10 (A/F-10) Evaluator{} grep narrower than the universal | 0012:S4 | fixed | Testing Strategy S4 + C4 |
| L11 (B/F-2) renderWrites phrasing could read as "doesn't exist" | 0012:C4 | fixed | Normative Contracts C4 |
| L12 (B/F-3) Failure Modes silent on A1's If-wrong | 0012:F1 | fixed | Failure Modes |

B/F-4 (S7 fixture string) was recorded by its author as informational-CONFIRMED,
not a defect: no disposition owed.

Net-new scope: none. Every row traced to a ledger entry; nothing charted.

## Grounding gate

Each finding was confirmed against `main` before it edited the draft:

- L1 — `internal/graphlint/reach.go::matchSatisfiable` keeps only
  `a.Block == table.BlockMatch` ("no guard atom is consulted at all"), and
  `atomAdmitsValue` calls `guard.Evaluator{}`. CONFIRMED; the Key Discovery
  was false as written.
- L2 — `internal/resolve/guard_mvv_test.go::conformingContractSeam` is driven
  at :462 and via `recordingSeam` at :256. CONFIRMED.
- L3 — `internal/cli/flow_next.go::probeRow(row, view, owned, observed, all)`
  calls `guardSeam()` with no model parameter. CONFIRMED.
- L4 — `internal/resolve/guardcontract.go:51` is the single `Key: "subject"`
  construction for all cases. CONFIRMED.
- L6 — `internal/guard/declaration.go::Conforms` checks required-key presence
  and single-valuedness only, zero non-test callers. CONFIRMED.
- L8 — no write-block category exists; `normalize.go::renderWrites` files under
  `CatMalformedTagDeclaration` (category.go:28). CONFIRMED.
- L9 — `internal/guard/assignment.go:273 func valueAssignments`; product.go
  only calls it. CONFIRMED.
- L10 — non-test zero-value forms are the three literals; `var`/`new`/embedded
  forms are admitted by Go and would yield the same nil mapping. CONFIRMED.

No finding failed the gate, so nothing was dismissed-with-cite.

L5 was collapsed rather than escalated: `0007:A17`'s own Evidence ("known to
the evaluator from the table it was built for") and its If-wrong ("the seam
takes whatever RDR 0003 shows the evaluator needs") license the narrowing, so
this is a recorded deviation against the headline's wording, not a contradiction
needing a route-back.

## Mini-checks

All five cues fired; five compact tables were written into Technical Design
(`#### Pre-lock mini-checks`): `authority`, `disposition`, `fidelity`, `oracle`,
`trace`. The desk trace walked the MVV stepwise against C1/C2/C3/C4/C5 and
produced NO CONTRADICTION row.

## Needs (re)verification — carried to Stage 6

- **A3** — stays Verified; its headline ("every non-test construction site has
  the loaded model in scope") is unchanged, but its Evidence now names a fourth
  CLI frame (`flow_next.go::probeRow`) where the model sits one frame up rather
  than at the call. Stage 6 should confirm the widened cost statement, not the
  verdict.
- **C3** — now binds a second implementer (`conformingContractSeam`) and
  re-keyed suite cases. Both are claims about test code that does not exist yet;
  they are Phase 2 obligations, verified at implementation, not now.
- No assumption flipped to Pending; no new A-N was added. The pass corrected
  and narrowed existing claims rather than introducing load-bearing new ones.

---

# cove — resolve pass (iteration 2, delta-scoped)

The iter-1 fixes were substantial, so the loop row returned `rerun`. The
second pass was delta-scoped to the 11 anchors iter-1 touched plus the five
newly-authored mini-check tables. It CONFIRMED every new source claim iter-1
introduced (probeRow's frame, guardSeam's call sites, conformingContractSeam,
renderWrites, Conforms, matchSatisfiable, valueAssignments, the `-0.0` path)
and found four defects the rewrite itself created or exposed.

| Ledger | Anchor | Disposition | Section touched |
| --- | --- | --- | --- |
| I2-F1 trace walks the MVV against two guard literals (`eq 3` vs `eq 7`) | 0012:MVV, mini-checks `trace` | fixed | MVV step 1/3, trace rows 1 and 3b |
| I2-F2 C5's venue pinned two ways (load-only vs every `conformKind` ingress) | 0012:C5, 0012:F1, `authority`/`disposition` | fixed | C5, Failure Modes, both tables |
| I2-F3 inline F4/F5 labels collide with projector ids after the A1 bullet | 0012:§failure-modes | fixed | Failure Modes + 4 citations |
| I2-F4 blanket re-keying puts `gte`/`contains` on matrix-illegal keys | 0012:C3 | fixed | Normative Contracts C3 |

Net-new scope: none — all four land inside the anchors iter-1 edited.

## The two that matter

**I2-F1 was self-inflicted, and is the `COMPUTE, DON'T ARGUE` failure.** The
iter-1 `trace` table asserted the MVV was jointly satisfiable instead of
walking it: MVV step 1 authored `iter eq 3` while step 3 and fixtures S7/S8
used `eq 7`, so under the guard the table itself named, step 3's `07` witness
parses to 7 ≠ 3 and prunes — the parsed-comparison leg never fires, and step
3b's `4 != 3` contradicts step 3's `7 == 7`. The desk trace is the owner of
joint satisfiability within one RDR, and it caught the defect only when it
was actually walked. Pinned to `eq 7` (S7/S8's `mvv-int` literal).

**I2-F2 was a real under-specification, collapsed on evidence, not escalated.**
`internal/cli/flow_input.go::canonicalValue` reaches `conformKind` through
`table.ConformValue`, so a round-trip placed in `conformKind` binds the CLI
BY CONSTRUCTION — "load-only" would require deliberately excluding it. The
record's own spike (`evidence/spikes/a4-typed-compare.md`) measures the cost
of that exclusion: `--tag iter=07` flips GuardFalse → GuardTrue under C2's
parsed comparison, a silent verdict change where the refusal is loud. Q1's
ruling ("reject non-canonical spellings at every authoring site", naming
`flow-tag-invalid`/`flow-write-invalid` among the reused codes) already
settled the direction. C5 now states the ingress rule and its reach; F1
carries the third visible break with the trade named.

## Needs (re)verification — carried to Stage 6

Unchanged from iteration 1 (A3's widened cost statement; C3's two
implementation-time obligations), plus:

- **C5's CLI reach** — C5 now claims `--tag`/`--write` inherit the round-trip
  through `ConformValue`. Grounded on `flow_input.go` calling `ConformValue`
  (confirmed) and on the measured flip in `a4-typed-compare.md`; the claim
  that the refusal diagnostic is identical at both venues is an
  implementation obligation, verified when C5 lands.
- No assumption flipped to Pending; no new A-N added.
