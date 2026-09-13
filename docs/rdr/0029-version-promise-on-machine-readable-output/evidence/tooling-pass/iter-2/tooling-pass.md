Model: claude-opus-5

# Tooling Pass — cli/0029 version-promise-on-machine-readable-output

Date: 2026-09-12 · Iteration 2 · Status at sweep: Draft
[revised from Final 2026-09-12; re-verify none @finalize]
Lint (final): `rdr lint --locking 0029` exit 0 —
`blocking=0 resolution=0 placeholder=0 advisory=1` (saved at `lint.txt`).
Re-entry: `@finalize`, RE-LOCK-ONLY. The re-entry note's listed defects ARE
this stage's content fix — applied in place below, then the note deleted.

## Findings

- **C1 (template coverage)** — PASS. No `template:missing-section`, no
  `placeholder:survived`, no `scaffold:row`, no
  `contract:template-example`. The four placeholder hits the prior
  iteration carried inside `§finalization-gate` are gone: the gate's
  guidance blocks went with the first lock, and this re-entry did not
  restore them. `§cross-cutting-concerns` remains authored (six concerns,
  five omitted as inapplicable) and is retained in the record at lock as
  `cli/0029:G-cross-cutting`.
- **C1 (hollow) — one BLOCK, FIXED IN-PASS.** `## Refinement Context
  (cluster re-entry — delete on re-lock)` survived at 1951-2021, reported
  by lint as `parse:section:unknown-to-template`. Per the Stage 7 prompt
  this is NOT READY *only if* a listed defect is still open; on a
  `@finalize` re-entry the listed defects are this stage's content fix.
  Both were discharged in live text (below), then the note was deleted.
  Mechanical class: the deletion is conformance, fixed here, not routed.
- **C2 (Method vocabulary)** — PASS. Seven rows, every
  `method.off_vocabulary` empty (`ca_off_vocabulary=0`), every row carries
  a Method: A1 `Peer RDR`; A2/A3/A4 `Source Search`; A5/A6/A7 `Spike`.
- **C3 (Source Search self-reference)** — PASS. A2/A3/A4 anchors resolve
  to `internal/version/`, `internal/cli/`, `internal/table/`,
  `internal/accessor/`, `internal/graphlint/`, `internal/guard/` and
  `internal/resolve/` symbols. None resolves to the record itself or to
  anything under its artifact dir.
- **C4 (Docs Only on load-bearing)** — PASS. No `Docs Only` record exists,
  so the case cannot arise.
- **C5 (symbol resolution)** — PASS. 78 edges resolve `true`; zero
  `resolved: false`. The C4 addition below cites `0021:C2`/`C5` and C1 by
  element id, all resolving. No bare `file:line` anchor.
- **C6 (status consistency)** — PASS. All 7 CAs `Verified`
  (`ca=all-terminal`); no `Pending`/`Unverified` property is relied on as
  settled fact. One status inconsistency found and FIXED IN-PASS — it is
  the re-entry's own C-10 defect, recorded under Re-entry note below.
- **C9 (evidence budget)** — ADVISORY, accepted unchanged. One
  `evidence:over-budget` at 272-317 (A4, 46 lines vs soft cap 30). Profile
  is `large`, not `foundational`, so `lint --locking` does not mark it
  blocking. Judgement unchanged from iteration 1: the load-bearing anchors
  stay findable in the field and the balance is the per-surface reasoning
  the grounding sweep reads. No relocation, no truncation.
- **C10 (linking)** — PASS. No `label:contracts` (C1–C4 all labelled), no
  `peer-evidence:no-element`, no `edge:unresolved`, no
  `edge:unresolved-terminal`. `joint-decision-home`
  `0029:JC1 → cli/0029:§normative-contracts` resolves.

## Re-entry note — both listed obligations discharged

The note carried one defect and one standing joint decision. Both are now
closed in live text and the note is deleted.

- **C-10, the false-premise delegation — CORRECTED.** Decision Rationale
  described 0021 twice as an early Draft with every assumption still
  Pending, and justified a prose delegation with "nothing is edited in a
  peer's locked text". Verified against the corpus: `0021`'s Date is
  2026-08-28 and its Status reads `Draft [revised from Final 2026-09-12;
  re-verify A4 @refine]` — it was Final for two weeks before this record's
  date and was demoted only at the 0021-0029 cluster gate. The premise was
  false and the justification exactly inverted. Both sites now state 0021
  was Final at authoring (demoted at the cluster gate to discharge the
  obligation), and the "in its own flow" justification is replaced by the
  tracked obligation that actually applies — 0021 re-enters at refine to
  assign the tiers C4 delegates. C4's delegation rule is unchanged.
- **C-7, two-marker read order — ANSWERED IN C4.** C4 already stated the
  order (envelope major decides whether the consumer can parse at all;
  the document's `schema` decides what the payload means). What was
  missing was the mixed-case behaviour the note asks for. Added as one
  paragraph in C4, composed from clauses already on the record rather than
  minted: unsupported envelope major → rejected under C1 before `data` is
  read, the document marker not consulted; supported envelope with
  unrecognized document marker → the envelope parses and its non-`data`
  members are trustworthy while `data` alone is opaque, since C1 already
  fixes that `schema_version` is never projected into `data` and so says
  nothing about the payload's own evolution. 0021 cites it rather than
  restating it, per the note.

The note's "not defects" (C-5/C-6 disclosed limits, C-12 designed
behaviour, C-4's free re-verification) were re-read and required no edit,
as recorded.

## LOOP-BREAKER

`rdr anchors --record 0029` over the prior pass (`$PRIOR_DIR/tooling-pass.md`)
and this one, intersected with `comm -12`. Non-empty — three anchors shared:

```
0029:G-cross-cutting
0029:JC1
0029:§normative-contracts
```

**Judged: NOT a routing loop; no `stopped:finalize-routing-loop`.** The
rule's subject is a FINDING re-reported after its named stage ran. None of
the three is a finding in either report:

- `0029:G-cross-cutting` — a BLOCK in iteration 1 (hollow, template text),
  fixed in-pass at that iteration and cited here only as a PASS under C1
  ("remains authored… retained in the record at lock"). Closed, not
  re-reported.
- `0029:JC1` and `0029:§normative-contracts` — the joint-decision home
  edge, cited as a resolving PASS under C10 and the fence in BOTH reports.
  A passing check named twice is not a finding at all.

The prior pass's actual findings were the hollow `§cross-cutting-concerns`
and the fence's `overlap-uncited`; both were closed at the first lock and
neither is raised here. This pass's findings — the surviving re-entry note
and the C-10 false premise — are new to the re-entry and were raised by the
7.1 cluster gate, not by a stage that failed to clear them. Corroborating:
`0029:C4`, carried by the prior report, is ABSENT from this one, which is
the opposite of a re-report. The intersection is citation overlap between
two PASS narratives.

## Joint-decision fence

`--outcome fence` over `overlap_uncited,rulings_open,clustered,impact_families`
resolves `op = none` (`rule: fence-clear`) — re-run after the edits above.
`rulings_open=0`, `joint_check_home=homed`, `clustered=false` so no
`impact.md` is owed and the fence does not ask for one.

VERDICT: PASS — no blocking finding; proceed to the Gate's written responses.
