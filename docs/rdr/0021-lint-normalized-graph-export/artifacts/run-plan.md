# Run plan — rdr-draft-to-lock (cli/0021)

What this run scheduled and decided. Not a position statement — the RDR's
current stage is derived from on-disk evidence (`/rdr-status 0021`).

```
rdr: 0021-lint-normalized-graph-export   profile: large (as read; Draft = provisional)
posture: upfront=false each=false finalize=true   (rule posture-draft-large)
lenses: grounding,3amigo,critique                 (critique runs --auto)
stages: refine -> resolve -> [lenses] -> reconcile -> finalize   stop-after: finalize
```

Entered at Stage 3 (first run; status Draft, form none).

Re-invoked 2026-09-12. Status still Draft, form `none` (not a routed-back
re-entry). Router re-answers `/rdr-refine` as before. The prior run's park was
put to the author, who ruled that its committed refine PASS satisfies the
advance carve-out (evidence/rulings.md, 2026-09-12); this run resumes at
Stage 4. Lens row and posture re-asked on re-entry, both unchanged.

## Ledger

| Stage | Verdict | Blocking | Note |
| --- | --- | --- | --- |
| refine | PASS | no | 5 contradictions resolved, 1 redundancy collapsed; Validation authored (9 scenarios). Commit aed9efe. |
| — | PARKED | yes | Router names `/rdr-refine` again (rule `locate-draft-refine`, guard `ca=all-pending`); refine leaves no done-signal, so the router cannot see commit aed9efe. Advancing to Stage 4 would override `emit.next`. Fork put to user. |
| — | RULED | no | Fork settled by author: accept inherited PASS, enter at Stage 4. See evidence/rulings.md. |
| resolve | NEEDS_DECISION | yes | A1/A2/A5/A6 verified by spike; A3 falsified in part (no edges carrier on `Reach`); A4 verified, citation re-anchor owed. Reuse audit: no existing export capability. Record left unedited — unapproved fixtures are not Evidence. Author's round at evidence/authors-round.md (Q1 blocking, Q2, Q3, F1–F3). Evidence commit a783b7e. |
| — | PARKED | yes | Packet row `fork` → park, stage `same`. Q1 blocks the `Profile` latch (`--outcome profile` → `stopped:split-signal`), so no lens row can be re-asked and `/rdr-prelock grounding` has no latch to read. |

## Pre-ask grounding (rdr-common §ground-before-ask)

Three read-only precedent sweeps spawned over the RDR corpus and Go source
before putting the round to the author — one-seam-vs-split practice (Q1),
empty-collection wire shape (Q2), conditional-assumption and
unreachable-branch precedent (Q3). Advisory only; they inform the author's
ruling and do not answer it. Not rulings: `rulings.md` carries the author's
words alone.

### Q2 — precedent found, and it diverges from Stage 4's recommendation

An empty-collection convention already exists and is unanimous. Verified
against source, not relayed: 16 declared collection fields across the CLI
render `[]` on empty, each with an explicit nil→empty normalization at its
producer; 0 are left nil-renderable. Fixed normatively by `0006:C14` ("the
empty list is the proof's receipt — so the `findings` field MUST NOT be
`omitempty`"), `0010` (nil-to-`[]string{}` normalization, "never `null`"),
`0023:C1`, `0024:C4`. Pinned by raw-bytes tests: `lint_0006_test.go:578-580`
(literal `[]`; "a `null` means the field carried a nil slice"),
`flow_partition_0023_test.go:871` (no `:null` on the line),
`flow_plan_c3xz_test.go:970-983`.

`omitempty` appears on exactly 4 output fields, every one an optional surface
(failure-envelope `findings` per `0005` REQ-8/16; `Notes`/`Warnings`;
`candidate.Gates` under a flag) — never on a declared collection.
`clierr.go:47-48` states the rule: "keep them `omitempty` so the envelope
stays append-only".

DIVERGENCE: Stage 4 recommended "`[]` for every declared collection,
reserving `null` for genuinely absent optional objects". The first half
matches precedent exactly. The second half has none — this codebase spells
an absent optional as key ABSENCE (`0023:C1`: "Absent means absent — never
`null`, `{}`, or `""` stand-ins"), never as `null`. No writer emits `null`
on any envelope. Adopting Stage 4's wording verbatim would make 0021 the
first record in the corpus to emit one.

Precedent-consistent wording, for the author's Q2 ruling: *declared
collections render `[]`/`{}`; optional members are absent, never `null`.*
C2 currently states neither (`0021.md:278` fixes only strict additivity) —
confirmed gap, so this is an addition to C2 and triggers §amendment-sweep.

### Q1 — corpus practice, and an inbound citation the sweep missed

Corpus practice on multiple `**Cn**` fences is unanimous and settles the
framing: multi-C records are facets of ONE seam, never a split. Zero records
have ever been split on contract count. Final/Implemented records carrying
2+ labelled C-blocks include 0002 (24 blocks), 0003 (26), 0004 (17), 0006
(20), 0007 (11), 0008 (6), 0009 (8), 0010 (5), 0011 (3), 0023 (2), 0024 (4),
0025 (7) — every one judged one seam at its gate. 0010's Profile literally
reads "five contracts … facets of one seam"; 0007's gate overruled a
critique that argued for five independent contracts. 0021's preamble wording
is copied from 0024's, which passed its gate with that exact framing.

So `stopped:split-signal` is not a claim that 0021 should be split. The row
counts labels; the split test is seams (`stages/README.md:234-236`,
`TEMPLATE.md:223-231`). Two post-rule records satisfied it differently:
0028 COLLAPSED (commit dcfe056, C1–C6 → clauses C1.1–C1.6); 0029 OVERRODE by
written ruling, keeping its labels.

CORRECTION to the precedent sweep, verified against primary sources: the
sweep concluded 0021 has no inbound label citation and that collapse is
therefore lower-friction. That is wrong. `0029` (Final) cites `0021:C2` and
`0021:C5` BY LABEL in both its body (`0029.md:1052-1056`, `:1106-1110`) and
its gate, and `0029/artifacts/gate.md:172` records 0021 as a "genuine joint
decision, settled here" — C2 mints the `schema` marker, C5 embeds the
document in the envelope's `data`. 0029's own Profile ruling
(`0029/evidence/rulings.md:5`) blocked ITS collapse for exactly this reason:
"collapsing to one `**C1**` is blocked because `0022:JC2` homes a joint
decision at this section's **C3** specifically, and relabelling would break
that inbound `joint-decision-home` edge."

Consequence for the author's Q1 ruling: collapsing C1–C5 into one `**C1**`
would break live inbound citations from a Final record, and would need those
0029 references re-anchored in the same change. The 0029 override model —
keep the labels, write the ruling, hand-write the Profile — is the path with
a matching constraint and a matching precedent. Still the author's call;
recorded here because the cost asymmetry was not visible in the round.
