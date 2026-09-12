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

### Q3 — `subsumes` is pinned by nothing; (b) stands rather than dissolving

The hoped-for shortcut (an existing record or test already guaranteeing the
`subsumes` property, making A2's condition pre-discharged) DOES NOT EXIST.
Verified directly, not relayed:

- The widening arm is genuinely dead. `indexOf` (`reach.go:249-256`) returns
  a node only where `subsumes` holds; `subsumes` (`:258-278`) requires an
  IDENTICAL key set plus value coverage. So `joinNodes(nodes[j], next)` is
  always key-equal to `nodes[j]`, the `merged.key() == nodes[j].key()` guard
  (`:153-155`) always trips, and `nodes[j] = merged; worklist = append(…)`
  (`:156-157`) cannot fire. The spike's structural claim reproduces at read.
- NO test pins it. `rg indexOf|subsumes` across `internal/` returns zero
  test references to node subsumption — every hit is an unrelated comment
  (row subsumption) or a different symbol (`indexOfStep`). `indexOf`/
  `subsumes` landed in one commit (`5e57b79`, 0006 Stage 8) and were never
  touched again.
- NO Final record states it. 0006 (Implemented) fixes only the intra-call
  join and a fixpoint (`0006:1082-1090`, REQ-108); its D11 pins successor
  identity by presence footprint, not the in-place widening arm.

So A2's "Verified" is genuinely conditional and ruling (b) stands.

COUPLING — a joint decision, not a declare-in-both. The ONLY record leaning
on the property normatively is **0022** (Draft): its C2 anchors on
"`reach.go::indexOf` precedent" (`0022:331-334`) and its A1 is **Pending**,
Method Spike, assuming `subsumes` finds a host node (`0022:142-152`). Draft +
Pending is not a guarantee 0021 may lean on. A pin test for the widening arm
would discharge a dependency shared by `0021:A2` and `0022:A1`/`C2` — which
makes it a joint-decision candidate to home in ONE record and cite from the
other, never to state in both. 0022 and 0021 currently declare mutual
independence (`0022:99-100`, `:510-511`); this shared dependency is the one
real seam between them and is not yet recorded as such.

PRECEDENT for the three shapes ruling (b) needs, all verified in-corpus:
- Qualified Verified is sanctioned and has a house format — 0007
  (Implemented) uses `Verified (<scope>); <what rides where>` and
  `Verified — with a recorded coupling: …` (`0007:450`, `:642`). The closed
  vocabulary is `Verified | Pending | Unverified` (`TEMPLATE.md:178`); there
  is no `Falsified`, so A3's write-up needs the narrowing form, not a new
  status word. A2 gates Phase 1, so it cannot simply stay Pending
  (`04-resolve.prompt.md:114-118`).
- Unreachable-branch pinning has no exhaustive-lemma precedent; the closest
  is 0015's **widening tripwire** — `TestAdvDeadEndExistentialOnMergedNode`
  (`adversarial_0006_test.go:489-506`), a behaviour test whose failure
  message names the re-run obligation. That is the shape to model.
- Narrowing at Stage 4 carries a recorded sweep under an
  `## Amendment sweep (§amendment-sweep)` heading with a per-site table and
  a closing grep statement (`0027/evidence/reconcile/reconcile.md:34-52`,
  `0010/evidence/critique/resolve.md:118-150`).

No contradiction with any Final record. One stale code comment:
`reach.go:124-126` / `:138-144` describe the in-place arm as live, which it
is not — commentary drift, not a record conflict.

## Second grounding round (deeper pass)

Three further analyses. Two CORRECT this file's own earlier entries; both
corrections are recorded here rather than silently overwritten.

### Q1 is forced — but by a different edge than first reported

CORRECTION to the Q1 entry above. That entry leaned on 0029's rulings.md
reasoning (relabelling breaks an inbound `joint-decision-home` edge). That
mechanical premise is WEAKER than 0029's prose claims: `0022:520` spells its
home `cli/0029 §Normative Contracts` C3 — a SECTION anchor with "C3" as
trailing prose — so relabelling 0029's C3 would stale prose, not break a
resolved edge.

The conclusion survives on stronger ground, verified directly in the edge
set: `0029 → cli/0021:C2`, kind `cross-cutting-owner`, `resolved: true`,
`0029.md:1868`. This is ELEMENT-resolved (contrast the bare `mentions` edges
at `:1049`, `:1106`) and originates in a Final, no-amend record.
`cross-cutting-owner` is a real kind (`rdr/tools/rdr/internal/edge/edge.go:74`).
Renaming `0021:C2` would break a resolved inbound edge from a Final record.

So Q1 → override-not-collapse is DETERMINED, not a free choice. The author
still rules, but the two dispositions are not equal-cost.

### The `subsumes` coupling: a CONSTRAINT, not a joint decision

CORRECTION to the Q3 coupling entry above, which called it "a joint
decision, not a declare-in-both". The disposition (one home, cite from the
other, never both) was right; the GRADE was wrong.

The doctrinal test is whether two records need one shared ANSWER either
could give differently (`rdr/prompts/gate/pairwise.md:69-71`; the JDR
README's "does not join merely for touching a listed member … only if it
joins the same question", `docs/jdr/README.md:88-90`). Neither record
CHANGES `indexOf`/`subsumes`; the answer was fixed by shipped code
(`5e57b79`). Their dependencies even differ in direction — 0021 needs "no
in-place widening after insertion"; 0022 needs re-anchoring determinism and
already has a total fallback. Weakening `subsumes` breaks 0021 and would
make 0022's A1 EASIER. There is no fork where they could answer differently.

The JDR lifecycle names this grade verbatim (`docs/jdr/README.md:152`):
`constraint` — "The answer is already determined by shipped code … Nothing
to negotiate; record it so nobody implements against it. **Does not
block**." So it is a homeless unpinned constraint, not an open decision, and
it does NOT gate the lock.

It is also INDEPENDENT of Q1: a joint edge leaves the citing JC and lands on
the home, so homing elsewhere gives 0021 an OUTBOUND edge touching no
C-label. Last turn's thesis — that the coupling might determine Q1 — is
withdrawn. The two questions are unrelated.

Why the mechanical joint-check missed it: `0021:502-506` grepped peers for
this record's modify-anchors (`reach.go::Reach`, the `graph` verb) and
contract literals (`--emit`, `intrastate.graph/1`, …). The shared anchor is
`subsumes`/`indexOf`, which neither record names as an anchor — so
`joint_check_home=clear` is a true answer to a question that did not cover
this. Recording the fire needs a hand edit to both Draft records.

Homing candidates, ranked (author's fork; doctrine does not settle it):
(a) JDR 0001 `constraint` entry — JD-23 already sits on this `reach.go`
seam with siblings 0015/0022; caveat, the merge relation is an adjacent
question to the quantifier, so this is a same-seam join, arguable.
(b) 0021 as RDR home — it touches `Reach`, produced the lemma and the test
design, and locks first (no `home-ahead-of-lock` flag).
(c) 0022 — weakest; farther from lock, and it needs only the definition it
already cites.

### F3: the caveat is a real spike bug; the fixture is still approvable

Verified at source, not relayed. `dotQuote` (spike `:95-100`) is correct and
its own comment states the required order. But `render` (`:132`) injects the
separator as Go `"\\n"` — two bytes, backslash+n — BEFORE calling
`dotQuote`, whose first step doubles that backslash. The hostile-input line
(`:270`) shows both on one line: injected separator `\\n` (renders as
literal text) beside a genuine newline `\n` (renders as a break). Rule text
right, bytes wrong, cause is the exact ordering fault the spike warns
against. Node ids and edge triples are unaffected — `dotQuote` is applied to
them directly (`:141`, `:155`).

Scope: label composition is a `label=` attribute, and C2 makes the node/edge
SET and marker normative while excluding styling (`0021:302-308`). The bug
falls OUTSIDE F3 as proposed. Caveat 2 is likewise a non-issue: `render`
reads Go fields, so the set/marker claim is rename-invariant. Precedent for
approving a fixture with an explicit scope limit: `0002:2275` ("the expected
row set and row identities — not the expected bytes. This SHA MUST NOT be
asserted as a golden hash").

The escaping ORDER is sound (walked against a hostile input and checked live
on graphviz 12.2.1); reversing it corrupts by closing the string early. One
spike prose error: its claim that a raw newline in a DOT quoted string is
"invalid/ambiguous" (`:298-300`) is false for 12.2.1 — it parses. Escaping
is still correct; the stated rationale is not.

### Q3(a): house form for the A3 narrowing

Precedent is narrow-in-place, not split: the Status head stays on-vocabulary
(`Verified | Pending | Unverified`, `fields.go:196-202`) with the narrowing
in a free-text qualifier, and the Evidence line quotes the pre-edit reading.
Verbatim precedents: `Verified (labels only — see A12 for the reachability
half)` (`0003:219`), `Verified (kernel half); …` (`0007:530`), `Verified as
narrowed — …` (tooling form). Evidence-line disclosure: "(the claim's
wording was narrowed to match at Resolve)" (`0023:438`).

SPLIT is precedented only where the failed half needs its own Method or
owner (`0003` A6→A12). Here the failed half (edges) is ALREADY the subject
of A2 (`0021:150-156`), so splitting would duplicate it — narrow plus a
cross-cite to A2 matches the corpus.

C2 needs NO field-list change: its "carries, at minimum … edges" predicates
*carries* of the JSON DOCUMENT, not the model value, and the document still
carries edges. The sweep's real target is `0021:521-522` ("since
`newAnalysis` already computes everything the document carries" — now false
for edges), plus a C4 consumer re-read.

Existing vocabulary to reuse rather than mint: 0002 C19's carried-vs-derived
pair (`0002:1298-1302`) and 0023's "pure function of" (`0023:441-449`).
Stage 4's proposed "which the model value fully determines" mints new
wording where "derived from / a pure function of the model value" is house.
