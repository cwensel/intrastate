Model: claude-fable-5-1

# Tooling pass — cli/0019 owned-state-initialization-semantics (2026-09-05, iter 1)

Run first: `rdr lint --locking 0019` → exit 0, `blocking=0 resolution=0
placeholder=5 advisory=12` (captured in `lint.txt`, pre-edit); the
post-edit re-run is the receipt. `inspect --json --filter
outline,assumptions,edges,metadata` read once.

- C1 (coverage) — no `template:missing-section`; every Required section
  present and authored (`## References` carries peers, source paths, docs,
  prior art). HOLLOW: the five `placeholder:survived` hits are all inside
  `## Finalization Gate` (the template's item guidance) — expected pre-lock
  state, resolved in-pass: `### Cross-Cutting Concerns` authored in the
  record (retained at lock), the four judged items written to
  `artifacts/gate.md` and replaced by the pointer at lock. `gate:inline`
  (conformance) is the same finding. No `scaffold:row`, no
  `contract:template-example`, no seed-skeleton header.
- C2 (Method vocabulary) — 8/8 rows carry a `method`; `off_vocabulary`
  empty on every row (Source Search ×4, Spike ×2, Source Search + Spike,
  Peer RDR).
- C3 (self-reference) — Source Search rows A1, A3, A6, A7, A8 anchor
  `internal/...::Symbol` paths only; no anchor resolves to the record or
  its artifact dir. Clean.
- C4 (Docs Only) — no row has `method.members == ["Docs Only"]`. Clean.
- C5 (anchor resolution) — every `source-anchor` edge `resolved: true`
  (repo-resolved via `$RDR_SOURCE_REPO`); none ABSENT; no bare
  `file:line` anchors. Two foreign-codebase anchors (qmuntal/stateless,
  XState) resolve as prior-art citations. Clean.
- C6 (status consistency) — Status `Draft`, canonical form, no qualifier;
  all 8 assumptions `Verified`, none Pending/Unverified, so no settled-fact
  prose can depend on an unverified one. Gate not yet written (no second
  copy to disagree). Clean.
- C9 (evidence budget, advisory) — `evidence:over-budget` ×5: A2 (50
  lines), A4 (43), A6 (60), A7 (31), A8 (58). Gate answer: each field opens
  on its load-bearing anchor; the balance is spike output / read-back
  verification the grounding sweep reads. Kept in the record; not moved,
  not truncated.
- C10 (linking) — no `label:contracts*`, no `peer-evidence:no-element`
  (A4 cites `0005:C1` / `0005:D-naming` as elements), no
  `edge:unresolved`, no `edge:unresolved-terminal`. `prose:exactness`
  (advisory, C1 "canonical"): covered by A2's byte-for-byte spike, RT3/S4
  and MVV step 3.

Also swept: no `## Refinement Context (cluster re-entry)` block
(`status_reentry=false`); repeatability-lite `emit.next: none`
(`repeatability-full-complete`; determinacy `determinacy-foundational`);
LOOP-BREAKER skipped (first pass, no prior dir).

Joint-decision fence (not a sweep check, recorded here for the lock
record): first run `stopped:overlap-uncited` — six open peers shared an
anchor or contract literal with no scanner-readable cross-citation
(0012, 0015, 0016, 0017, 0021, 0022), and cli/0028's fire on this record
(home cli/0019:C1) had no symmetric line here. Propose's three arms re-run
on the written proposal at finalize; every pair disposed cite-don't-restate
(no shared answer) and the `Joint-check:` line rewritten with `cli/NNNN`
citations and the symmetric fire. Re-run: `fence-clear`, `op = none`;
`joint_check_home=homed`, `overlap_uncited=0`.

Verdict: PASS — no blocking finding; proceed to the Gate's written responses.
