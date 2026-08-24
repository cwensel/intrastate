Model: claude-opus-5[1m]

# Critique Diff — RDR 0002, iter-3 (A: opus-5 24 rows / B: sonnet-5 11 rows)

Reconciled by RDR passage anchor, not by `C-N` id. Grounding verdicts are
independent: each was checked against the draft text and, where the finding
cites one, against the peer RDR file, `evidence/spikes/iter-2/` source, or
`internal/resolve/resolve.go` on `main`.

**Iteration context.** iter-1 and iter-2 converged against an earlier draft.
The RDR was then demoted (`567600f`, cluster-reconcile STAGE-SCOPED) and
re-authored through refine → resolve → cove iter-2 → 3amigo iter-3 against
JDR 0001 §D7's closed layout. This pass reads the *current* draft, so net-new
rows here are genuinely net-new surface, not critique-on-critique drift. Both
passes were instructed not to read prior critique evidence.

## Merged Ledger

| M | Sources | RDR passage anchor | Failure mode (one line) | AGREEMENT | GROUNDING |
|---|---------|--------------------|--------------------------|-----------|-----------|
| M-1 | A:C-1, B:C-3 | §Normative Contracts dump clause, "(`row_identity`, `source_locator`, … `write` …) — that closed list is the column vocabulary" vs. `iter-2/rdr-fixture.toml:182` `order = ["identity","source",…,"writes",…]` | Both normative fixtures author three column ids outside the RDR's own closed vocabulary, so under its own `malformed dump declaration` rule both fixtures MUST refuse at load. | **both** | **grounded** (verified: draft:1052 vs fixture:182-184; mismatch is exactly threefold, and the pre-§D7 fixture has it too) |
| M-2 | A:C-6 | §Normative Contracts next/writes, "**RDR 0009 A4 fixes that these are distinct fields**" (twice) | Fabricated peer citation: 0009 A4 settles the *opposite* — "the predicate does **not** widen; it stays `Writes`-only". | A-only | **grounded** (verified 0009:454-455, 830-833) |
| M-3 | A:C-7 | §Normative Contracts reserved-value, "**RDR 0004 removes the key on that write and asserts absence on read-back**" | Fabricated peer citation: "clear" appears zero times in 0004's body; its read-back asserts *presence*, exemptions "none". | A-only | **grounded** (verified: 0004's own Status line says "`clear` occurs nowhere in 0004's body"; 0004:215,511 assert presence) |
| M-4 | A:C-8 | §Normative Contracts guard, "a three-valued domain, `match`, `all`, or `unless` … **RDR 0007 spells**" | Fabricated peer citation: 0007 spells a **two**-valued domain. The widening is legitimate; its authority is JDR 0001 §D6, not 0007. | A-only | **grounded** (verified 0007:1252 `block ∈ {all, unless}`) |
| M-5 | A:C-9 | §Normative Contracts escape, "Row kind … a **derived view property, never a row field**" | Contradicts 0006's minimum input contract, which requires row kind as an input; 0009 quotes 0002's now-deleted opposite text, and 0009 A4's Verified note rests on it. | A-only | **grounded** (verified 0006:497 "row kind and the declared list of failure classes"; 0006:1186) |
| M-6 | A:C-10 | §Normative Contracts root/stop-set, "root `terminal` is a list of context ids" | Contradicts 0006 A6: "neither is a bare identifier a row references by name" — 0002 ships exactly that. | A-only | **grounded** (verified 0006:266) |
| M-7 | A:C-11 | §Conditional Mini-Checks `disposition`, "expansion counts per rule \| lint \| RDR 0006" | 0006's advisory tier is normatively **closed at four members**; 0002 invents a fifth. 0006 contains no relevant sense of "expansion". | A-only | **grounded** (verified 0006:1004-1006; `grep -c expansion` = 1, unrelated) |
| M-8 | A:C-12 | §Normative Contracts literals, "MUST NOT derive that key … by joining members into a string" | 0006's mandatory finding payload declares `Literal` a flat `string`, reintroducing the banned joined rendering at the locked diagnostic boundary. | A-only | **grounded** (verified 0006:796) |
| M-9 | A:C-2 | §Normative Contracts dump clause, "An `order` naming an unknown identifier … refused at load" vs. Testing Strategy 3 "the spike decodes `[dump]` and never reads it" | `[dump]` is the only sub-schema with a full normative rule and **zero** executed evidence. Reproduced: `["identity","NOT_A_COLUMN"]` loads clean. | A-only | **grounded** (reproduced by running the spike) |
| M-10 | A:C-3 | §Normative Contracts operator clause, "the guard-side mirror … **which without this was absent**" | Guard operator membership unenforced; `frobnicate = "small"` loads clean. The RDR notices the rule was absent, adds prose, adds no fixture and no implementation. | A-only | **grounded** (reproduced; `main.go:529-537` checks only key-declared and not-`recognized`) |
| M-11 | A:C-4 | §Normative Contracts reserved-value, "`<clear>` MUST be refused … **or a predicate literal**" | Enforced in `checkMatchBlock` only; `eq = "<clear>"` under `unless` loads clean — producing exactly the "silent dead rule" the clause names as the dangerous case. | A-only | **grounded** (reproduced; `main.go:563-565`) |
| M-12 | A:C-5 | §Normative Contracts tag-type, "a literal outside the tag's declared domain" | Domain check is match-block only though the clause is unqualified as to block; `eq = "NOT_A_PROFILE"` loads clean under `unless`. | A-only | **grounded** (reproduced) |
| M-13 | A:C-13, B:C-2, B:C-6 | §Conditional Mini-Checks `trace`, "merge atoms" / "write values" — two **CONTRADICTION** verdicts, A13 `Pending` | The draft locks with two self-declared contradictions between its normative clauses and its only executable evidence, on the arm that "trips **zero** load categories". | **both** | **grounded** |
| M-14 | A:C-16, B:C-1 | §Testing Strategy 2, "`output.txt` … SHA-256 `6ccfe901…`, whose nine RDR rows and two kata rows are the expected value" | The SHA-pinned oracle is rendered text produced by a renderer that joins set literals — so a *correct* implementation of the literals clause fails the pinned hash. **B** frames it as the canonical artifact containing the banned defect; **A** frames it as an oracle-class error contradicting the draft's own "assert over the normalized value" guidance. | **both** | **grounded** (verified `main.go:898` `strings.Join(parts, ",")`; `output.txt` shows `labels=needs work`) |
| M-15 | A:C-15 | §Technical Design "(the iter-2 spike populates `Locator` from `rule.Source` … a spike defect against this clause)" vs. §Normative Contracts "the authored provenance string the locator derives from" | Two incompatible locator origins one clause apart; the SHA-pinned witness has two kata rows colliding on `kata:review`. | A-only | **grounded** |
| M-16 | A:C-14 | §Normative Contracts `RequiresOwned` + 0007 "predicate PLACEMENT decides…" | Match-vs-guard placement — the format's most consequential authoring decision — is delegated by a Final peer to 0003; 0002 owns the authoring surface and offers no guidance, lint, or diagnostic. Missing owned state is laundered into an escape-rescued plan. | A-only | **grounded** |
| M-17 | A:C-17, B:C-4 | §Prerequisites, "Phases 2 and 3 therefore sequence behind the reshape **in their entirety**" | Phases 2/3 are unstartable and no RDR/kata/phase owns 0007's reshape. **A** frames it as locking an unstartable plan; **B** frames it as a locked ring of Final consumers around a Draft producer with a hole in the middle. | **both** | **grounded** on the unstartable/unowned half; **re-raise** on the lock half (see Dismissed) |
| M-18 | A:C-18, B:C-8 | §MVV "deferred to the Phase 5 lint handshake" + §Technical Design "Ambiguous overlap … is therefore a lint finding" | Eight hazard classes route to an unimplemented 0006. **A** adds that the RDR's own fixture `[initial] stage = "propose"` is matched by no rule, so the canonical example ships with an unreachable root. **B** adds that a table with no root and no stop set loads and dumps clean. | **both** | **grounded** |
| M-19 | A:C-21, B:C-9 | §A1 Evidence "covering **18 of the 25**" / Testing Strategy 3 "The seven still owed" → six names → "Two of those **five**" | Three cardinalities in one paragraph for the completeness claim of the RDR's central API surface. **B** adds that two of the seven owed categories were minted *at the 3amigo pass itself* and have neither fixture nor implementation. | **both** | **grounded** |
| M-20 | A:C-20, B:C-11 | §Normative Contracts load clause, "**Load is fail-fast** … order … deliberately **unspecified**" | Fail-fast + unspecified order for a format whose primary journey is hand-authoring; the RDR concedes the fit and defers the remedy to an unscheduled surface. | **both** | **grounded** as a stated-and-accepted cost; **re-raise** as a defect (see Dismissed) |
| M-21 | A:C-19 | §Normative Contracts contexts, "**Inheritance therefore never overrides — it only accumulates**" | The whole case for Alternative 1 over Alternative 2 rests on a factoring mechanism with no override; any rule needing a different value on an inherited key forks the chain — Alternative 2's rejection reason, relocated to contexts where the dump cannot show it. | A-only | **grounded** |
| M-22 | A:C-22 | §Normative Contracts accessor, "every **owned** tag MUST be served by exactly one reader" | Unconditional load refusal on an operational property; a write-only owned tag (audit stamp, one-way flag) is unauthorable. The *observed* arm is carefully provenance-scoped; the owned arm is not. | A-only | **grounded** |
| M-23 | A:C-23 | §Normative Contracts next/writes, "MUST NOT populate one by aliasing the other" + Testing Strategy 2 "not assertable by value comparison" | A MUST NOT whose oracle the RDR itself declares unassertable, falling back to review — the control the draft's own `oracle` mini-check rejects elsewhere. | A-only | **grounded** |
| M-24 | A:C-24 | §Round-Trip, "the rendered form of a **set-valued atom literal** non-recoverable" vs. §Problem Statement "one reviewable artifact … has one answer" | The stated user outcome is one reviewable answer; the review surface is then specified provably ambiguous and the reviewer told to read the source the artifact exists to replace. | A-only | **grounded**; overlaps charted `M-16` from iter-2 |
| M-25 | B:C-5 | §Approach "The model is data, not generated code and not a runtime FSM engine" vs. the eight-part schema + ~14 normative blocks | The format's real conceptual weight approaches the "small DSL" it rejected on parser/tooling cost, without paying that cost down via tooling. | B-only | **grounded** (framing/altitude, not a clause defect) |
| M-26 | B:C-7 | §A5 Evidence, "the first parse-validation codes in the envelope … new design work at implementation" | 25-category taxonomy with no worked CLI/Hint mapping will be redesigned once 0005 fits it to a few exit-code groups. | B-only | **re-raise** — iter-2 charted this to 0005 (M-12) |
| M-27 | B:C-10 | §Normative Contracts `[model.metadata]`, "deliberately unconstrained … no consumer today, and that is the point" | A schema-free nested extension table is the one place in a strict schema where anything goes; it grows undisciplined, drifting content with no lint or reviewer signal. | B-only | **grounded** (narrow half — the drift risk; the reserve-now decision itself is adjudicated) |

---

## Hotspots

Findings both models reached independently, ranked. These are the passages where
two readers with no shared context converged.

**H-1 — M-1: the dump column vocabulary refuses its own canonical fixtures.**
A:C-1 and B:C-3 land on the same three-name mismatch from opposite directions —
A from the closed-vocabulary clause outward to the fixture, B from the fixture
inward to the clause — and both computed the diff rather than arguing it. This is
the highest-confidence row in the ledger and the cheapest to fix. It is also
*newly created*: the 3amigo iter-3 pass minted `malformed dump declaration` and
closed the column vocabulary (IMP3-007) without checking the clause against the
fixtures the same document declares normative. That is precisely the
"COMPUTE, DON'T ARGUE" failure the resolve gate names, produced by the immediately
preceding lens. Both models caught it on first read of the result.

**H-2 — M-14: the SHA-pinned oracle contradicts the literals clause.** B reaches
it as "the canonical evidence artifact contains the banned defect"; A reaches it
as "a rendered-text golden is the wrong oracle class for a value contract, and
the draft says so itself." Both are right and they compose into a sharper
statement than either: the pinned SHA does not merely *record* the defect, it
*enforces* it — any correct implementation of the literals clause changes the
bytes and fails a Stage-4-approved hash, so the pin actively punishes the fix.
That is a trap laid for the implementer, not a stale artifact.

**H-3 — M-13 / M-17 / M-18 / M-19 / M-20: five convergences on "locks with
known holes."** Both models independently arrived at the same structural
observation from five different clauses: two live CONTRADICTION verdicts (M-13),
unstartable Phases 2/3 (M-17), eight hazards routed to unimplemented lint (M-18),
three contradictory coverage counts (M-19), and a conceded-poor-fit load model
(M-20). No single row blocks lock; the convergence is the signal. Note the models
*disagree on what to do about it* — see Contradictions X-1.

---

## Contradictions

**X-1 — is "locks with known holes" a defect or the plan?** A treats M-17/M-18/
M-20 as lock-blocking (§2's "one section rewritten in six weeks" nominates the
fail-fast clause; the premortem opens on unstartable phases). B treats the same
rows as accepted, disclosed costs and nominates the *category taxonomy* (M-26)
as the rewrite candidate instead. The draft supports both readings because it
states the costs explicitly and then declines to schedule the remedies. Resolution
is not textual: M-17's "gates implementation sequencing, not lock" was already
adjudicated at iter-2 (R-3) and stands; what is *not* adjudicated is whether an
unowned reshape can gate sequencing indefinitely with no owner named. That
residue is real and narrower than either model's framing.

**X-2 — row kind: derived or a field?** M-5 is a genuine cross-RDR fork, not an
under-specification. 0002 says "never a row field"; 0006's minimum input contract
requires it as an input; 0009 quotes 0002's deleted opposite text and rests a
`Verified` assumption on it. Three Final/Draft documents cannot all be right.
This one cannot be collapsed inside 0002 alone — it is a cluster-reconcile
obligation, and it is the strongest argument in this pass for routing back rather
than resolving in place.

---

## Refuted / Dismissed

**R-1 — M-17's lock half (A:C-17, B:C-4).** Re-raise of an iter-2 adjudication.
The draft says verbatim: "This item gates implementation sequencing, not lock,"
and iter-2 dismissed the identical claim with the identical cite (R-3). The
*narrow* residue — no RDR, kata, or phase owns 0007's reshape — survives and is
new; the lock claim does not.

**R-2 — M-20 as a defect (A:C-20, B:C-11).** The fail-fast choice is stated,
its cost conceded in the draft's own words, and a batch mode explicitly reserved
to 0006/0005 by the 3amigo iter-3 pass (PM3-003). Re-litigating it is the
"already-decided" ground the resolve gate forbids. What is *not* decided is
whether the reserved batch mode has an owner — same residue shape as R-1.

**R-3 — M-26 (B:C-7).** iter-2 charted the category→CLI mapping to RDR 0005
(`Charted.md` M-12). Unchanged since; a re-raise, not a new finding.

**R-4 — M-25 (B:C-5).** The data-vs-DSL fork is adjudicated in §Alternatives
with a written rationale. B offers no new evidence, only a re-weighing of the
same trade-off. Dismissed as a re-raise; the tooling-gap observation inside it
is real but is 0006/0005 scope, already charted.

**R-5 — M-27's decision half (B:C-10).** Reserving `[model.metadata]`
pre-consumer is decided (JDR 0001 §D7(v)) and was re-affirmed with a testability
fix at 3amigo iter-3 (PM3-004/QA3-002). The drift-risk half is not adjudicated
and survives as a narrow, cheap addition.

---

## Coverage note

A found 24, B found 11; overlap is 7 (M-1, M-13, M-14, M-17, M-18, M-19, M-20).
The models split cleanly by method. **A executed** — it ran the spike binary
against mutated fixtures and reproduced four under-enforcement defects
(M-9…M-12) that no transcript records, because `gen-cases.py` never mutates a
guard block; and it opened every Final peer and diffed the citations, producing
the M-2…M-8 cluster (seven cross-RDR defects, three of them fabricated
citations) that two prior lenses missed. **B read the document as a whole** —
its unique rows are altitude and disclosure judgments (M-25, M-27) plus the
fixture-vs-clause computation it shares with A.

The asymmetry is instructive: 11 of A's 24 rows required *leaving the document*
(running code or opening a peer), and none of those would be visible to a reader
who stayed inside the draft, however careful. The union is the useful artifact,
and the cheapest durable fix this pass suggests is mechanical — a citation-check
and a fixture-conformance check, both scriptable, both currently absent.

Of 27 merged findings: 22 grounded, 4 re-raises of adjudicated decisions
(one with a surviving narrow residue), 1 grounded-in-part.
