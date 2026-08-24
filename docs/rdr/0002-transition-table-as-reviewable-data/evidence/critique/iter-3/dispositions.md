Model: claude-opus-5[1m]

# Critique Dispositions — iter-3 (dual-model, cross-model diff)

Origin ledger: `iter-3/critique.md` (opus-5, C-1..C-24) + `iter-3/critique-modelB.md`
(sonnet-5, C-1..C-11), reconciled by passage anchor in `iter-3/diff.md` as
M-1..M-27.

Both passes were run fresh against the post-demotion draft (refine → resolve →
cove iter-2 → 3amigo iter-3 all landed after critique iter-2), so this is a full
pass, not a delta-scoped re-run of iter-2's ledger. Net-new rows are net-new
surface. Net-new *scope* was still held to the charted path.

## Fixed

- **fixed** — M-1 (hotspot, both models) the dump column vocabulary refuses its
  own normative fixtures. Computed, not argued: draft:1052 fixed the vocabulary
  as snake_cased prose field names (`row_identity`, `source_locator`, `write`);
  `iter-2/rdr-fixture.toml:182-184` authors `identity`, `source`, `writes`. Three
  of ten disagree, and the pre-§D7 fixture has the same spelling, so the fixtures
  are the stable side. Fixed the vocabulary to the fixtures' spelling verbatim
  and said why the fixture spelling governs. **This defect was created by the
  immediately preceding lens** (3amigo iter-3 IMP3-007 minted
  `malformed dump declaration` and closed the vocabulary without checking it
  against the fixtures) — the exact "COMPUTE, DON'T ARGUE" failure the resolve
  gate names. Section: Normative Contracts / dump.

- **fixed** — M-2 fabricated citation: "RDR 0009 A4 fixes that these are distinct
  fields" (twice). Verified 0009:454-455 — A4 settles the **opposite**: "the
  predicate does **not** widen; it stays `Writes`-only." A4 fixes no rule about
  the two fields' distinctness. Restated the empty-both-fields consequence as
  this RDR's own authoring-surface result (stricter than A4's kernel-side
  predicate) and re-attributed field distinctness to JDR 0001 §D2. Sections:
  Normative Contracts / next-writes (both sites).

- **fixed** — M-3 fabricated citation: "RDR 0004 removes the key on that write
  and asserts absence on read-back." Verified: "clear" appears zero times in
  0004's body; its read-back asserts *presence* with "lossy exemptions: none",
  and 0004's own Status line records the gap (§D5/§JD-15 not yet carried).
  Replaced the false settled-fact claim with the accurate owed-state: this RDR
  fixes the sentinel's authoring and rendering contract only, and a clearing rule
  is not end-to-end until 0004 lands §D5. Section: reserved-value clause.

- **fixed** — M-4 fabricated citation: the three-valued block domain "RDR 0007
  spells". Verified 0007:1252 — 0007 spells `block ∈ {all, unless}`, two-valued,
  and 0003 quotes the same. The widening to `match` is legitimate but its
  authority is JDR 0001 §D6, not 0007. Re-attributed, and warned the implementer
  that building the atom type from 0007 alone yields two members and MUST widen
  here. Section: guard clause.

- **fixed** — M-6 root `terminal` as bare context ids contradicts RDR 0006 A6
  ("neither is a bare identifier a row references by name"). 0006 leaves the
  spelling to this RDR but fixes that lint receives *tag predicates* — its
  invariant 1 checks them as keys/values and invariant 2 evaluates a terminal as
  a predicate over a node. Collapsed rather than escalated: added the missing
  dereference — a terminal context id is a reference, normalization MUST
  dereference to the context's explicit predicate set, and the normalized value
  carries predicate sets, never bare ids. Section: root/stop-set clause.

- **fixed** — M-7 the headline risk's mitigation routes to a closed lint tier.
  Verified 0006:1004-1006 — the advisory tier is normatively closed at four
  members, and 0006 has no notion of expansion at all. The mitigation for
  "sparse contexts hide an accidental Cartesian product" therefore could not be
  built as written. Collapsed on evidence: expansion count is **derivable from
  this RDR's own dump** (rows sharing a source rule id, already carried by the
  `identity` and `source` columns), so the mitigation moved onto this RDR's
  deliverable and off 0006. Swept two further sites to match. Sections: Risks
  and Mitigations; `disposition` table; Failure Modes.

- **fixed** — M-8 the literals clause's join ban vs. RDR 0006's flat-`string`
  `Literal` field. Verified 0006:796 — the string typing is deliberate, to keep
  `clierr` free of this RDR's vocabulary. Narrower than the finding framed it:
  that is a **transport/serialization** boundary, not an identity one. Stated the
  distinction — display is not identity; rendering into a diagnostic field is
  permitted, deriving merge keys/dedup/sort/round-trip from a rendering is not —
  and asked producers to render sets visibly. Section: literals clause.

- **fixed** — M-10, M-11, M-12 (one clause, one root cause) atom-level validation
  is enforced in match blocks only. All three reproduced by **running the spike**
  against hand-mutated guard blocks: `frobnicate = "small"` (unknown operator),
  `eq = "<clear>"` (reserved value), and `eq = "NOT_A_PROFILE"` (outside declared
  domain) each load clean under `[rule.guard.unless.profile]`; `main.go:529-537`
  checks only key-declared and not-`recognized`. Added a block-agnostic normative
  clause covering all three blocks for every atom rule, named the spike's
  behavior a defect against it, and recorded why no transcript witnesses the gap
  (`gen-cases.py` never mutates a guard block, so the hole is invisible to
  `negative-cases.txt` by construction — which is how two prior lenses passed over
  it). Booked as **A15 (Pending)**. Sections: new normative clause after the
  operator clause; A15; Testing Strategy 3.

- **fixed** — M-14 (hotspot, both models) the SHA-pinned oracle enforces the
  banned defect. The pin makes `output.txt` "the expected value", but its bytes
  were produced by `renderValue`'s `strings.Join(parts, ",")` — so a *correct*
  implementation of the literals clause changes them and fails a Stage-4-approved
  hash. That is an anti-oracle: it teaches the implementer to revert the fix.
  Demoted the SHA to provenance-only with an explicit MUST NOT assert-as-golden,
  restated the expectation as the row set and identities, and pointed at the
  draft's own "assert over the normalized value" rule. Checked the other two SHA
  sites and **left them**: they assert run-to-run and permutation byte-identity,
  which holds under any renderer including a corrected one. Section: Testing
  Strategy 2.

- **fixed** — M-15 two locator origins one clause apart. The floor clause is
  explicit and correct ("derived, not authored"; composed from `(model id, rule
  id)`); only the schema line contradicted it by calling `source` "the authored
  provenance string the locator derives from". Fixed the schema line to match the
  floor. Section: closed-layout schema.

- **fixed** — M-16 match-vs-guard placement has no authoring guidance in the
  document that owns the authoring surface. Grounded on `main`: `RequiresOwned`
  is checked by `resolve.go::missingOwned`, and a match atom on a failed owned
  read drops every candidate, so `escapeOrRefuse` returns a **plan** — missing
  owned state laundered into a successful transition, which JDR 0001 §D2 exists
  to forbid. RDR 0007 routes this guidance to RDR 0003; repeated it here because
  this is the document open while the author types the block, and a rule stated
  only where the author is not looking is not a control. Section: block-intent
  clause.

- **fixed** — M-19 (both models) three cardinalities in one paragraph. Recomputed
  from the transcript rather than from any file's prose: `negative-cases.txt`
  has 37 `refused:` lines, 2 of them `probe-*`, so **35 category controls** — the
  draft's 35/18-of-25/2-probe figures are all correct. The only defect was
  internal: "seven still owed" followed by six names and then "two of those
  five". Named the arithmetic explicitly (five pre-existing + two added = seven).
  Section: Testing Strategy 3.

- **fixed** — M-21 inheritance never overrides, so contexts fork. Grounded: the
  no-override rule is correct and load-bearing (an override reintroduces the
  order dependence the full-atom-identity merge key exists to prevent), but the
  cost was undisclosed and it is Alternative 2's rejection reason relocated from
  rows to contexts — where it is *less* visible, since the dump shows expanded
  rows but not context repetition. Stated the cost, the reason it is accepted,
  and the successor trigger. Section: contexts/merge clause.

- **fixed** — M-22 write-only owned tag unauthorable. Collapsed on the kernel
  rather than escalated: `RequiresOwned` carries every written/cleared key and
  `missingOwned` evaluates it against the view, so an owned tag written on one
  transition must be readable on the next or that resolve refuses
  `owned_state_unavailable`. The one-reader rule is load-bearing, and the load
  refusal is the early diagnosable form of a resolve-time failure. Stated the
  reason (the observed arm was justified; the owned arm was not) rather than
  relaxing the rule. Section: accessor binding clause.

- **fixed** — M-23 a MUST NOT whose oracle the RDR declares unassertable. The
  draft named the real control in its own last clause (mutate one field, assert
  the other unchanged — both are `[]Tag`, so an alias shares a backing array) but
  led with "enforcement is by review". Promoted the mutation test to the control
  and cut the review fallback, per this RDR's own `oracle` rule that an assertion
  is not accepted until a wrong implementation that fails it is named — "review
  the normalizer" names none. Section: Testing Strategy 2.

- **fixed** — M-24 the review surface is specified ambiguous against a Problem
  Statement promising "one answer". The 3amigo pass had already fixed the
  mechanism (visible-set rendering, grammar-by-decision); what was missing was
  reconciling it with the promise. Stated the narrowing precisely: the artifact
  answers which rows exist, what each binds, and in what order — all recoverable;
  what is not recoverable from rendered text alone is a set literal's member
  decomposition, a strictly narrower question. Named the two rejected
  alternatives and why the residue is cheaper. Section: Round-Trip.

- **fixed** — M-27 (narrow half) `[model.metadata]` drift. The reserve-now
  decision is adjudicated (JDR 0001 §D7(v)) and was made testable at 3amigo
  iter-3; the drift risk was not addressed. Since constraining the namespace
  would make it a schema, the control is scope, not validation: stated that
  metadata is for data no consumer reads, that anything acted on is a tag, and
  that the first consumer is the trigger to revisit. Section: `[model.metadata]`
  clause.

- **fixed** — M-17 / M-20 residues (the halves that survived the re-raise
  dismissals below). M-17: RDR 0007's reshape gates Phases 2/3 and **has no
  owner** — 0007 is Final so will not schedule it, no kata tracks it, no phase
  claims it; sequencing behind an unowned prerequisite is indefinite, not merely
  ordered. Stated that, and banned the local-constant-mirroring workaround the
  Risks section already names as drift. M-20: neither reserved batch-diagnostic
  surface is scheduled, so the authoring cost is real and unmitigated on arrival;
  labelled it a known debt and the strongest follow-up candidate. Sections:
  Prerequisites; fail-fast clause.

## Dismissed with cite

- **dismissed-with-cite** — M-5 row kind "never a row field" contradicts RDR
  0006/0009. **The finding's direction is wrong, and checking it inverted the
  defect.** 0002's clause is self-consistent: it forbids a kind field in the
  *normalized value and kernel row* while explicitly permitting a render-time
  view column ("a render-time view struct MAY materialize the computed column").
  0006 consumes a *graph view*, so its minimum-input "row kind" is satisfied by
  the permitted derived column. What is actually broken is **RDR 0009**, which
  quotes 0002's deleted text at 0009:244 and 0009:492 and rests an A4
  `Consistency note` marked `Verified` on it. That is a stale citation in a
  `Final` peer, not a defect in this draft, and it is not this RDR's to fix —
  **charted to cluster-reconcile**. Note both 0006 and 0009 carry "checked
  consistent at the 0002-0009 iteration-3 gate", so the gate missed it.

- **dismissed-with-cite** — M-17 (lock half) / M-20 (defect half) / M-26 / M-25 /
  M-27 (decision half). Five re-raises of adjudicated decisions:
  - M-17 lock half: the draft says verbatim "This item gates implementation
    sequencing, not lock", and iter-2 dismissed the identical claim with the
    identical cite (R-3). The narrow ownership residue was fixed above.
  - M-20 as a defect: fail-fast is chosen, its cost conceded in the draft's own
    words, and a batch mode reserved by 3amigo iter-3 (PM3-003). The scheduling
    residue was fixed above.
  - M-26 category→CLI mapping: charted to RDR 0005 at iter-2 (`Charted.md` M-12);
    unchanged since.
  - M-25 data-vs-DSL weight: adjudicated in §Alternatives with a written
    rationale; B re-weighs the same trade-off with no new evidence. The tooling
    residue inside it is 0006/0005 scope, already charted.
  - M-27 decision half: reserving the namespace pre-consumer is decided
    (JDR 0001 §D7(v)) and re-affirmed at 3amigo iter-3.

- **dismissed-with-cite** — M-9 `[dump]` has a normative rule and zero evidence.
  Not dismissed on the merits — it is **already booked**: Testing Strategy 3
  names `malformed dump declaration` as owed both a fixture and an
  implementation. M-1's fix corrects what those fixtures must assert against, and
  the new fixture-conformance control is the check that would have caught it.
  No separate edit owed.

- **dismissed-with-cite** — M-18 eight hazards route to an unimplemented RDR
  0006. The load/lint split is adjudicated on **arity, not severity** (a
  Load-Bearing Decision, re-affirmed at iter-2), and cross-row properties are
  not decidable at load by construction. B's added observation — a table with no
  root and no stop set loads and dumps clean — is the same decision seen from the
  other side and is explicitly the documented behavior ("the loader accepts the
  absence and lint refuses to certify"). A's added observation about the
  fixture's unreachable `[initial]` is real but is an RDR 0006 scenario-4
  obligation, already named.

## Charted to successors

Recorded in `iter-3/Charted.md`: M-5 (RDR 0009's stale quotes of deleted 0002
text, one carrying a `Verified` A4 consistency note → **cluster-reconcile**,
the sharpest item this pass produced), and a **citation-conformance check** plus
a **fixture-conformance check** as cheap mechanical guards → tooling-pass.

## Needs verification (Stage 6)

New or amended load-bearing claims from this pass, all **Pending**:

1. **A15 — atom-level validation is block-agnostic** (new assumption). The
   current evidence *contradicts* it: three reproduced cases load clean under
   `guard.unless`. Method: `gen-cases.py` must mutate guard blocks; one negative
   control per atom rule per block. Until then the coverage figures overstate
   what is proven — an atom rule counted "witnessed" is witnessed in one of three
   blocks.
2. **The terminal dereference rule** (M-6) is a new normative obligation on
   normalization and a new claim about what RDR 0006 receives. Needs checking
   against 0006's invariants 1 and 2 before lock.
3. **Expansion count is derivable from the dump** (M-7). Asserted, not
   witnessed — no control computes rows-per-source-rule-id today.
4. **The dump column vocabulary** (M-1) changed value. Every clause and fixture
   naming a column must agree; the new fixture-conformance control is the
   oracle, and it does not exist yet.
5. **`output.txt`'s hash is no longer a golden** (M-14). The re-rendered fixture
   and its hash are owed at implementation; nothing pins them now.
6. **The no-alias mutation test** (M-23) is a new control with no fixture.
7. **The `<clear>` end-to-end path is owed on RDR 0004** (M-3), not on this RDR.
   Cross-RDR: a clearing rule cannot be validated end-to-end until 0004 carries
   §D5.

No previously `Verified` assumption was invalidated. A15 is net-new and Pending
from birth.

## Tiebreakers

None escalated. Three forks were collapsed on evidence rather than surfaced:

- **M-5 row kind** — appeared to be a three-document design contradiction
  requiring redesign or route-back. Reading 0002's own clause (view column
  permitted), 0006's input contract (graph view, not kernel row), and 0009's
  quotes showed 0002 is correct and 0009 is stale. Collapsed to a charted
  cluster-reconcile item.
- **M-22 write-only owned tag** — relax the one-reader rule, or keep it? Collapsed
  against the kernel: `missingOwned` makes a write-only owned tag unusable, not
  merely unbound, so the load refusal is the early form of a resolve-time
  failure. Kept the rule, stated the reason.
- **M-8 `Literal` as string** — is the diagnostic boundary a violation of the join
  ban? Collapsed on 0006's own rationale: the string typing is a transport
  decision to keep `clierr` dependency-free. Display is not identity; both
  clauses stand.

## Mini-checks

No new cue fired. This pass edited three tables the earlier lenses established —
`disposition` (expansion-count routing corrected off RDR 0006), `oracle` (the
no-alias mutation control promoted; the SHA demoted from golden), and `fidelity`
(the Problem-Statement narrowing stated) — none introducing a mini-check the
draft was not already carrying. The Stage 5 mini-check cue read was discharged by
this RDR's first lens pass this iteration (cove iter-2); tables persist in the
draft.
