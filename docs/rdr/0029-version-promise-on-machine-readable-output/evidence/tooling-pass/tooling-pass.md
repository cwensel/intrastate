Model: claude-opus-5

# Tooling Pass — cli/0029 version-promise-on-machine-readable-output

Date: 2026-09-12 · Iteration 1 (base) · Status at sweep: Draft, no qualifier
Run twice: an initial sweep that BLOCKED on one hollow section, the
mechanical fixes applied in-pass, then this re-run.
Lint (final): `rdr lint --locking 0029` exit 0 —
`blocking=0 resolution=0 placeholder=4 advisory=6` (saved at `lint.txt`).
LOOP-BREAKER: no `$PRIOR_DIR` — first pass, nothing to diff.

## Findings (re-run)

- **C1 (template coverage)** — PASS. No `template:missing-section`. All
  Required sections present and authored; `## References` carries five
  substantive citation groups. The four remaining `placeholder:survived`
  hits (1871-1874, 1880-1892, 1898-1900, 1982-2006) are all inside
  `§finalization-gate` — the guidance blocks the lock removes when
  responses 1/2/3/5 move to `gate.md`. `gate:inline` (1841-2006) is the
  expected pre-lock shape, not a finding.
- **C1 (hollow) — FIXED IN-PASS.** The first sweep found
  `§cross-cutting-concerns` entirely TEMPLATE.md text: the retained-at-lock
  note, the instruction block, the candidate-concern menu and the
  determinism rider, nothing authored. This is the ONE gate item the lock
  keeps in the record (peers cite it as `cli/0029:G-cross-cutting`), so a
  hollow one cannot survive the lock. Authored from material the RDR
  already carried — six concerns (versioning, incremental adoption, build
  tool compatibility, character encoding, canonical-form/determinism,
  secret lifecycle), five omitted as inapplicable. It now emits three
  resolved `cross-cutting-owner` edges. Mechanical class: fillable from
  the record's own body, fixed here, not routed to Refine.
- **C2 (Method vocabulary)** — PASS. Seven rows, every
  `method.off_vocabulary` empty, every row carries a Method:
  A1 `Peer RDR`; A2/A3/A4 `Source Search`; A5/A6/A7 `Spike`.
- **C3 (Source Search self-reference)** — PASS. A2/A3/A4 resolve to
  `internal/version/`, `internal/cli/`, `internal/table/`,
  `internal/accessor/`, `internal/graphlint/`, `internal/guard/`,
  `internal/resolve/` symbols. No anchor resolves to the record itself or
  to anything under its artifact dir.
- **C4 (Docs Only on load-bearing)** — PASS. No `Docs Only` record exists.
- **C5 (symbol resolution)** — PASS. 58 `source-anchor` edges (53 before
  the fixes), all `resolved: true`; no bare `file:line` anchor. The four
  ABSENT edges are `kind: artifact` (`{SPIKE_DIR}` ×3,
  `{ARTIFACT_DIR}/gate.md`) — path-map placeholders, outside this check.
- **C6 (status consistency)** — PASS. Metadata Status `Draft`, no
  qualifier, no re-entry. All 7 CAs `Verified` (`ca=all-terminal`). One
  internal inconsistency found and corrected in-pass: the joint-decision
  paragraph described this record as "`mid` against 0022's `large`" while
  the Profile field reads `large` — the stale self-description was removed.
- **C9 (evidence budget)** — ADVISORY, accepted. One
  `evidence:over-budget` at 272-317 (A4, 46 lines vs soft cap 30). Profile
  is `large`, not `foundational`, so `lint --locking` does not mark it
  blocking. The field is the fourteen-surface vocabulary census; its
  load-bearing anchors (`graphlint::AggregateCode`,
  `resolve::RefusalKinds`, `accessor::Verdicts`, and six more) stay
  findable in place and the balance is the per-surface reasoning the
  grounding sweep reads. No relocation, no truncation.
- **C10 (linking)** — PASS. No `label:contracts` (C1–C4 all labelled), no
  `peer-evidence:no-element`, no `edge:unresolved`, no
  `edge:unresolved-terminal`. `joint-decision-home`
  `0029:JC1 → cli/0029:§normative-contracts` resolves.

## Re-entry note

No `## Refinement Context (cluster re-entry)` block present. Nothing to clear.

## Joint-decision fence

First run stopped `stopped:overlap-uncited` (`overlap_uncited=1+`): three
in-flight peers shared an anchor or literal with no cross-citation — 0012,
0014, 0021. Each was fired as a joint-decision question against the peer's
own contracts (delegated reads; see `gate.md` §Joint-decision fence):

- **0021 — genuine joint decision.** `0021:C2` mints a required `schema`
  marker versioning the exported document; `0021:C5` embeds that document
  in this envelope's `data`. Settled in `0029:C4`: document-versions-itself,
  envelope versioned by C1, later surfaces tier in the record that adds
  them. The propose-time `clear` ("adds no envelope field") was stale and
  is corrected in Decision Rationale.
- **0014 — coincidental.** `0014:C1/C3` decide gate-oracle mechanics on
  `code` / `graph-lint-failed`; C4 here decides their stability tiers.
  Complementary; `Makefile::docs` is non-normative prior art in both.
  Acknowledged in Decision Rationale.
- **0012 — coincidental.** `internal/resolve` / `internal/table` are
  package-path co-mentions; 0012's contracts add no refusal kind, no
  `Block` member, no table category.

Re-run after the citations: 0 uncited pairs on both intersects; the fence
resolves `op = none` (`rule: fence-clear`).

VERDICT: PASS — no blocking finding; proceed to the Gate's written responses.
