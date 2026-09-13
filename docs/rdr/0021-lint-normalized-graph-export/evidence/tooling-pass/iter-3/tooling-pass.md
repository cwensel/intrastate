model: claude-opus-5

# Tooling Pass — RDR 0021 (lint's normalized-graph export), iter-3

Run: Stage 7 mechanical pre-sweep, immediately before the Finalization Gate's
written responses. Source: `rdr lint --locking 0021` at `./lint.txt`, plus
`inspect --json --filter assumptions,edges,metadata`. Findings only; the one
RDR edit this pass made (the 0017 joint-check citation) was made by the Gate,
not by this sweep.

Re-entry context: this is the re-lock pass after the 7.1 cluster reconcile
(`0021-0029`) demoted this record `Draft [revised from Final 2026-09-12;
re-verify A4 @refine]`. The qualifier is already cleared (`status_form=none`,
`status_reentry=false`).

## CHECK 1 — Template section coverage

- One `parse:section:unknown-to-template` (1226-1327): the `## Refinement
  Context (cluster re-entry — delete on re-lock)` block. This is the 7.1
  re-entry note, and Stage 7's rule is that it MUST be gone by re-lock with
  its defects folded into live text. All six of its items verify closed in
  live text (evidence below), so the note is discharged and deleted in this
  pass. Not a spine hole.
- No `template:missing-section`. No `placeholder:survived`, no `scaffold:row`,
  no `contract:template-example` — the four gate responses this lock moves to
  `gate.md` are authored, not bracketed.
- No surviving `this is a seed skeleton` header.

Re-entry note discharge, item by item:

| Note item | Closed in live text | Where |
| --- | --- | --- |
| Defect 1 — tiers never assigned (PW-1, C-1/2/3) | yes | C1 tiers the `--emit` set `append-only` with a named `internal/cli` enumeration accessor, and inherits `seam: none (prose-only)` for `graph-export-too-large` from `0029:C4`'s CLIError row; C2's STABILITY paragraph tiers the field spellings `append-only` with the json-tags seam |
| PW-2 — no enumeration seam | yes | same two seams named above |
| Defect 2 — C2 asserts a forbidden cardinality (C-14) | yes | C2 now reads "carries the dump's field vocabulary rather than restating a cardinality" and adds the explicit MUST-NOT on cardinality/ordinal/tail; the "closed 11-member list … appended last" wording is gone |
| Defect 3 — joint-check never saw 0029 (PW-3, C-9) | yes | Decision Rationale's Joint-check line reads "fired → 0029 … re-run 2026-09-12 on the post-0029 peer set"; JC1 exists with a resolved `joint-decision-home` edge to `cli/0029:§normative-contracts` |
| C-11 — re-verify A4 post-0029 | yes | A4 carries an explicit "Re-verified 2026-09-12 against the post-0029 envelope (C-11)" paragraph, discharging both asks (premise unaffected; no exact envelope key set asserted) |
| C-8 — scope S2's goldens to the document | yes | S2 reads "capture the DOCUMENT, not the envelope — `--as=text`, or `jq .data`" |
| C-7 — standing joint decision (two-marker read order) | homed | `joint_check_home=homed`; home `cli/0029 §Normative Contracts` C4, edge resolved; this record cites rather than restates |

## CHECK 2 — Method label vocabulary

PASS. Eight rows, every `method.off_vocabulary[]` empty
(`ca_off_vocabulary=0`). Members: `Spike` (A1, A2, A5, A6), `Source Search`
(A3, A7, A8), `Peer RDR` (A4) — all sanctioned. No row lacks a Method field.

## CHECK 3 — Source Search self-reference

PASS. The `Source Search` rows are A3, A7, A8; every Evidence anchor resolves
into the product tree (`graphlint/reach.go::Reach`, `::Node`, `::reach`,
`guard/declaration.go::AssignmentCount`, `guard/product.go::Groups`,
`graphlint/analysis.go::newAnalysis`, `graphlint/engine.go::Run`,
`graphlint/taxonomy.go::CodeProductTooLarge`). None resolves to this record or
its artifact dir. No regression from the rounds.

## CHECK 4 — Docs Only on load-bearing claims

PASS, vacuously. No assumption carries `method.members == ["Docs Only"]`.

## CHECK 5 — Symbol resolution of Source Search / Spike anchors

PASS. Every `kind == "source-anchor"` edge reports `resolved: true` — none
false, none ABSENT, so `--repo` was supplied and the lookups genuinely ran. No
blocking `edge:unresolved`. A5's `clierr.go:174::WriteJSONLine` carries a
symbol that resolves; the stale line component is a documented NON-finding.

## CHECK 6 — Status consistency

PASS. Status `Draft`, no qualifier. All eight assumptions `Verified`
(`ca=all-terminal`), so no Pending/Unverified property is relied on as settled
fact. The gate responses move to `gate.md` in this same lock, so there is no
second copy to disagree with the record.

## CHECK 9 — Evidence-field budget

PASS. No `evidence:over-budget`. Profile is `large`, not `foundational`, so
this check is advisory here regardless — and there is nothing to report.

## CHECK 10 — Linking

PASS. C1-C5 labelled (no `label:contracts`); A4's `Method: Peer RDR` cites an
element, not a bare record (no `peer-evidence:no-element`); no
`edge:unresolved`, no `edge:unresolved-terminal`.

## Joint-decision fence — one blocker, fired and cleared in this pass

The fence first resolved `stopped:overlap-uncited` (rule `fence-uncited`):
`index --literal-intersect --record 0021` reported `0017 0021 UNCITED 1
shared: code`. Fired the check rather than syncing copies, grounded per
§ground-before-ask (`--outcome ground` → `apply`, settled in source):

- `0017:C1` decides what `findings[i].code` names for a multi-subject
  producer — a per-finding producer-local discriminator, never inherited from
  the envelope code.
- `0021:C4` forbids the verb to emit findings at all, and its one new code
  `graph-export-too-large` is a scalar ENVELOPE code (GroupUserEnv, exit 2).
- `0021:C1`'s `model-invalid` arm mirrors `lint.go::runLint` verbatim, which
  reaches the shared `internal/cli/flow_input.go::loadFindings` — already
  populating each entry's `code` with the REQ-24 load-category slug
  (`Code: string(category)`, flow_input.go:240), never the envelope code.
  `lint.go:207` and this verb call the one loader.

So 0017's rule is satisfied by the shipped loader this record reuses, and
0021 decides nothing in 0017's domain. Recorded as a citation on the
Joint-check line naming 0017. Post-edit the pair reads `cited`.

## LOOP-BREAKER

`rdr anchors --record 0021` over iter-2's report returns empty (that pass was
a clean PASS carrying no anchored findings), so the `comm -12` intersection
with this pass is empty. No finding re-reported after its named stage ran; no
routing loop.

## Verdict

PASS — no blocking finding; proceed to the Gate's written responses.
