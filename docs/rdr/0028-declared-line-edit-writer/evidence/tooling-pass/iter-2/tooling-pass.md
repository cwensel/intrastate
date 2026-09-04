Model: claude-opus-5[1m]

# Tooling Pass — cli/0028 declared-line-edit-writer (iter-2)

Stage 7 mechanical pre-sweep, run 2026-09-04 immediately before the
Finalization Gate's written responses. Second pass: this is a RE-LOCK after the
07.1 cluster reconcile demoted the record (`Draft [revised from Final
2026-09-03; re-verify A7, A9 @refine]`) and Stages 3/4/6 ran over JDR 0003
§D2 (a).

Basis: `rdr lint --locking 0028` (exit 0, captured at `iter-2/lint.txt`) plus
`inspect --json --filter assumptions,edges,metadata 0028` and the outline. Lint
run once; every CHECK below reads the captured file.

## Findings

- **C1 — Template section coverage**: one finding, non-blocking at this point
  in the pass. No `template:missing-section`; no `placeholder:survived` at all
  (the four Gate sub-sections were replaced by the `gate.md` pointer at the
  first lock and stayed replaced through the re-entry). No `scaffold:row`, no
  `contract:template-example`, no seed-skeleton header. The one hit is
  `parse:section:unknown-to-template` at 1634-1660 — the `## Refinement Context
  (cluster re-entry — delete on re-lock)` note. That is the demote note, and
  deleting it is part of the re-lock, not a spine hole: see the re-entry
  disposition below.
- **C2 — Method label vocabulary**: PASS. 12/12 assumption rows carry a Method
  field; `method.off_vocabulary` empty on every row
  (`ca_off_vocabulary_ids=[]`). Distribution unchanged from iter-1: Source
  Search ×8 (A1, A2, A6, A7, A9, A10, A11, A12), Spike ×3 (A3, A4, A5), Peer
  RDR ×1 (A8).
- **C3 — Source Search self-reference**: PASS. The 8 Source Search rows'
  Evidence anchors resolve into the consumer source tree
  (`internal/accessor/…`, `internal/cli/…`, `internal/table/…`); none resolves
  to the record itself or under its artifact directory. A7 and A9 were rewritten
  this round and both still cite source, not the record.
- **C4 — Docs Only on load-bearing claims**: not applicable — no `Docs Only`
  row.
- **C5 — Symbol resolution of anchors**: PASS. 65/65 `source-anchor` edges
  resolve `true` (up from 60 — A7 and A9's re-verification added anchors on
  `parseTags`/`canonicalValue` and `cmdbind.go::refuse`). None `false`, none
  absent-unlooked. No bare `file:line` anchor missing a symbol.
- **C6 — Status consistency**: PASS. `ca=all-terminal` — 12 Verified, 0
  Pending, 0 Unverified. The re-verify set (`reverify=["0028:A7","0028:A9"]`)
  is closed: both read `Verified` and both Evidence fields answer the amended
  C1.6 specifically, A7 establishing that no upstream guard already refuses a
  `-`-prefixed bound value and A9 that `cmdbind.go::refuse` is the rule C1.6
  mirrors. No checklist-vs-gate disagreement is reachable — the gate responses
  are single-sourced to `gate.md`.
- **C9 — Evidence-field budget (ADVISORY)**: one hit, unchanged from iter-1 —
  `evidence:over-budget` on A2's Evidence field, 35 lines against a soft cap of
  30. Answered at the Gate under Proportionality: the load-bearing anchors are
  still findable and the balance is verification content the grounding sweep
  reads. Not truncated, not relocated. Advisory; does not block.
- **C10 — Linking**: PASS. No `label:contracts` (C1.1–C1.6 labelled), no
  `peer-evidence:no-element` (A8 cites JDR 0003 §D1), no `edge:unresolved`, no
  `edge:unresolved-terminal`. The six `joint-decision-home` edges resolve.

## Note — the two unresolved `mentions` edges persist, still not findings

`0028:MVV → 0022` (line 1131) and `0028:S12 → 0021` (line 1275) report
`resolved: false`, as at iter-1. Both name test FIXTURE files under the RDR
tooling's `testdata/status/records/`, which the untyped `mentions` heuristic
pattern-matches on the `NNNN-slug` shape. `mentions` is not a blocking edge
class; rewriting either to point at a peer record would make it wrong.

## Re-entry note disposition

`## Refinement Context (cluster re-entry — delete on re-lock)` at 1634-1660
names four defects. Each is verified closed in live text, so the note is
deleted at lock rather than blocking it:

1. C1.6 `admission:` gains the argv0 exclusion — present: "POSITION: a
   `{tag.<key>}` element at argv0 is `edit_tag_argv0` at LINT (C1.4)", with the
   category registered in C1.4's list and its precedence stated (it fires on
   `command` argv, disjoint from the five that fire on the `edit` table).
2. C1.6 `binding:` gains the `-`-prefixed value refusal — present: "VALUE: a
   bound value beginning with `-` refuses `execution_failure` BEFORE spawn,
   Detail naming the placeholder, mirroring `cmdbind.go::substitute`'s existing
   `{artifact}` rule". Per-USE-SITE scope stated; CWE-88 named as the class.
   C1.3's `order:` and `precedence:` both carry it as an entry-level
   precondition.
3. Cross-Cutting secret/credential response rewritten — present: the response
   now states C1.6 IS a new caller-supplied value source, explicitly retracts
   the earlier reading of A7, and bounds the exposure by the two C1.6 rules.
4. JDR 0003 §D3 cited at C1.3 for the exit group and carried into S27 —
   present at C1.3 `order:` ("EXIT GROUP: … JDR 0003 §D3 (b), which decides
   this and adds no class") and at S27, which pins the exit-2 group, the
   `flow_exec.go::accessorFailureOf` discriminator and the `findings[]`
   carriage while leaving the code's spelling to Stage 8.

Also cleared this pass: `--outcome repeatability` → `emit.next: none`
(`repeatability-lite-complete`); the joint-decision fence → `emit.op: none`
(`fence-clear`), with `joint_check_home=homed`.

## Verdict

**PASS** — no blocking finding. `lint --locking` exit 0, `blocking=0
resolution=0 placeholder=0`; the only advisory is C9 on A2, answered at the
Gate, and the one parse finding is the re-entry note the lock removes.
Proceed to the Finalization Gate's written responses.
