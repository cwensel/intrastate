Model: claude-opus-5[1m]

# Finalization Gate — RDR cli/0002, Transition Table As Reviewable Data

- **Date**: 2026-08-24
- **Verdict**: **PASS — READY. Locked to Final.**
- **Lock generation**: re-lock. The 2026-08-23 lock was undone by
  `/rdr-cluster-reconcile` iter-3 (SPEC-DEFECT, STAGE-SCOPED, re-verify
  A1/A9/A12). This record **overwrites** that one and describes the current
  draft: re-authored to the JDR 0001 §D7 closed layout, 15 assumptions, three
  `Pending`.
- **Mechanical pre-sweep**: `evidence/tooling-pass/iter-2/tooling-pass.md` —
  **PASS**, no blocking findings. One MECHANICAL C1 (stale 2026-08-23 lock
  pointer in `## Finalization Gate`) is resolved by this lock's own prescribed
  steps. C9 advisory: 3 fields over the 30-line budget (A13 60, A15 33, A7 31),
  reviewed and kept.
- **Pre-lock lenses run** (`Profile: foundational` → cove → 3amigo → critique →
  repeatability): all four complete with dispositions — cove iter-2, 3amigo
  iter-3, critique iter-3 (dual-model), repeatability ×3 across three distinct
  model stamps, diff health **Healthy**.
- **Stage 6 reconcile**: `evidence/reconcile/report.md` (iteration 2) —
  RECONCILED, every item terminal, no BLOCKER surviving.

## 1. Contradiction Check

**No contradictions found between research findings, design principles, and
proposed solution.** The four candidates re-examined against the re-authored
draft all resolve.

**Research supports the chosen split.** The Investigation's prior art is
uniformly *source-model-plus-rendered-view*: Sismic's `import_from_dict` /
`export_to_dict`, `transitions`' positive `conditions` / negative `unless` guard
lists, Stateless's `StateGraph` built from machine metadata, Sismic's
`PlantUMLExporter` as a rendered view over model data. The sparse-TOML-source +
expanded-table-dump shape is that same prior. The statechart hit
(hierarchy/extended state as the answer to dimensional explosion) is honored by
inherited contexts, while "not a runtime FSM engine" is honored by keeping the
libraries as vocabulary/validation/visualization only. No research finding
recommends a representation this RDR rejects.

**Exact-one vs. first-match is a deliberate, argued divergence.** Decision
Rationale rejects first-match explicitly ("priority order is convenient in code,
but it makes review harder and lets source or dump reordering change
behavior"), and A2 verifies at source that the shipped kernel already implements
gate-then-count. A named divergence from surveyed prior art is not a
contradiction.

**"Hand-authored, not generated" and the dump are compatible.** Background says
the table "must be hand-authored from the legal graph audits, not generated."
The RDR generates the *expanded table*, never the *source model* — Approach:
"The rendered table is review support, not the source authors maintain."

**The producer-locks-after-consumers inversion is booked, not hidden.** The
Risks section names it directly ("This RDR is the cluster's wire-format producer
and locks after several of its consumers"), and its mitigation is that the
layout, routing key, and sentinel are fixed by JDR 0001 §D5–§D7 and *cited*
rather than derived, with cross-RDR drift checked by `/rdr-cluster-reconcile`
before implementation. Critique iter-3 sharpened one concrete instance (M-5:
RDR 0009 quotes RDR 0002 text on row kind that this draft has since deleted and
inverted to "a **derived view property, never a row field**"), and the
disposition is correct: this RDR's own clause is self-consistent, RDR 0002
cannot edit a `Final` peer, and per project doctrine RDRs are never amended, so
the obligation is charted to `/rdr-cluster-reconcile`. Stale peer *fact*, not a
contract conflict. Not a lock blocker — and it is the specific reason the Next
pointer below routes through cluster-reconcile.

**One internal contradiction class was found and closed rather than carried.**
Critique's M-13 named two live CONTRADICTION rows between normative clauses and
the only executable evidence — the spike compared member sequences by a joined
rendering (A13) and enforced atom rules in match blocks only (A15), so the
document's clauses and its own spike disagreed. Stage 6 escalated to
`§strong-consult`, verified both premises at source, and **ran the fix in-pass**
rather than deferring: both rows now hold, both assumptions are `Verified` on
executed controls, and the baseline digest `6ccfe901…` is unchanged. The draft
locks with no clause its evidence contradicts.

## 2. Assumption Verification

**15 Evidence Records (A1–A15). Every record is internally consistent:** Status,
Method, and Evidence agree, and no `If wrong` is empty.

**No `Docs Only` records** — the vocabulary distribution is `Spike` ×8
(A1, A3, A6, A7, A11, A13, A14, A15), `Source Search` ×3 (A2, A5, A8),
`Peer RDR` ×1 (A4), `MVV Test` ×3 (A9, A10, A12). Nothing blocks on the
Docs-Only rule.

**No self-referential `Verified` stamp, and every cited symbol resolves.** The
three `Source Search` records cite `internal/resolve/` and `internal/cli/`
paths; none resolves to this RDR file or its artifact directory. Symbol
resolution was verified against the working tree this pass — `Row`
(`resolve.go:169`), `Resolve` (`:318`), `assemble` (`:148`), `gate` (`:376`),
`escapeOrRefuse` (`:473`), `missingOwned` (`:444`), `Table.models` (`:216`),
`recognizedTagKey` (`:110`), `CLIError` (`clierr.go:48`), `Fail`
(`respond.go:133`), `Load` (`config.go:74`), and
`TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder`
(`adversarial_test.go:266`). No bare `file:line` anchor appears in any Evidence
field.

**Twelve of fifteen are `Verified`, and the promotions since the last lock are
backed by executed evidence, not by re-reading.** A1 was re-verified at Stage 4
by re-authoring both fixtures to the §D7 layout and re-running under strict
decoding; A7, A11, and A14 were promoted by that same run. Three worth naming:

- **A14** now has the discriminating control it lacked. The old
  `unsupported version` fixture was v1-shaped and trips on either check
  ordering, witnessing nothing about precedence. `neg/neg-v2-shaped.toml` is
  `version = 2` plus a v2-only root table the strict decoder would reject; it
  refuses `unsupported version 2`, **not** `unknown schema field`, and the
  paired control (same document, `version = 1`, v2-only key retained) refuses
  `unknown schema field: schema_v2`. One key, two outcomes, only the gate
  ordering separating them.
- **A13** is the strongest record in the set and the reason the delimiter
  parameterization earned its keep twice. The pre-fix spike joined literals on
  `,` and collapsed **only** the comma case — all four other delimiters yielded
  two atoms while the defect was live, so a single-delimiter control chosen
  unluckily would have certified it. The `|` case then caught a residual in the
  *render*, fixed by quoting each member. Fifteen fixtures, both arms, minted by
  `gen-cases.py`.
- **A15** entered at Stage 6 and closed in the same pass. The gap was invisible
  to `negative-cases.txt` by construction (`gen-cases.py` mutated match blocks
  only), which is why two prior lenses passed over it; the fix generalized
  `checkMatchBlock` to `checkAtomBlock` and scoped domain/kind conformance
  **per operator** — a real design refinement the pass forced, since applying
  the clause literally to a `lt` bound breaks the RDR fixture's own
  `iter.lt = 3`.

**A8's residual is stated rather than overclaimed.** The alphabet ban is a
load-time check in the normalizer; the kernel imports no normalizer and
validates no alphabet well-formedness, so the barrier covers exactly one
producer — tables this RDR's loader built. Closing it for all producers is
RDR 0001's to make. Correct shape for a bounded guarantee.

**Three records are `Pending` (A9, A10, A12), each a deliberate Stage 6
DOWNGRADE with a named MVV oracle — down from six at the previous lock.**

| ID | Why Pending | Where the oracle lives |
| --- | --- | --- |
| A9 | Blocked on RDR 0007's kernel reshape — the field this assumption targets does not exist, so the block is structural, not a scheduling preference | Scenario 2 normalization assertions + scenario 4 `Resolve` run against one normalized value |
| A10 | Decidability half source-verified (`Model` is a singular struct field mirroring singular `[model]`); paired-document behavior owed, and no single-file surface can express it | Scenario 3, a **pair** of documents sharing a `model id` |
| A12 | Blocked on the same reshape; discriminating witnesses now authored (`continue-prelock-cluster`'s `cluster_ready.eq=true@all`, kata's `status.eq=closed@unless`) but no field exists to route into | Scenario 2, union of `Match` and guard atoms equals the normalized set, intersection empty |

**The deferral floor was applied, not waived.** Stage 6 distinguished items it
*could* close from items it could not: A13 and A15 each had a wide arm that
changes match semantics while tripping **zero** load categories — the owed
fixture was the only detector, so deferring the fixture deferred the only
detector — and both were run now. A9 and A12 are blocked on work Stage 6 cannot
unblock; A10's remaining half needs a surface that does not exist. That is the
floor working as intended rather than a carve-out.

**Status consistency holds.** No settled-fact prose leans on an unsettled
assumption. The three sites where a Pending assumption's subject appears outside
its own record were read in full: the `Match`/`Guard` routing clause
(`:1170–1185`) is a `MUST`, and its one factual claim — `Guard string` carries
no per-atom block — is source-verified; the three `duplicate model id` sites
(`:1330–1342`, `:1430–1440`, `:1601`) are `MUST` / `MUST NOT` clauses about a
loader that does not exist yet plus the source-verified singular-`[model]` fact;
A12's Testing Strategy entry states its own blocker. Prerequisites names the
actual `Pending` set and its count, with the `[x]` qualified "**except A9, A10,
and A12**" — the run-1-era checkbox mismatch was corrected at Stage 6.

The RDR's `Status: Draft` line carries no 07.1 qualifier while the README row
reads `Draft [revised from Final 2026-08-24; re-verify A1, A9, A12 —
STAGE-SCOPED, re-enter Stage 3]`. Per TEMPLATE.md that qualifier self-clears at
this flip, and its named re-verification is discharged: A1 re-verified at
Stage 4 against §D7; A9 and A12 re-examined at Stage 6 and downgraded on stated,
structural grounds.

## 3. Scope Verification

**The Minimum Viable Validation is in scope and executed during implementation,
not deferred — with one item explicitly and correctly excluded.**

The specific proof: parse the RDR and kata sparse TOML fixtures into typed
source structs; normalize them into candidate rows; dump the expanded table;
validate the full load-time category set; then prove by unit test that one
sample tag-set resolves through `internal/resolve::Resolve` to exactly one
ordinary row, one unmatched tag-set resolves to exactly one modeled escape row,
and one unsupported-version variant is refused before normalization. Testing
Strategy scenarios 1–6 carry the assertions, and **41** load-time controls are
already witnessed one-refusal-per-mutated-fixture in
`evidence/spikes/iter-2/negative-cases.txt`, covering 18 of 25 categories. The
seven still owed are enumerated by name, with two of them (`cyclic context
inheritance`, `malformed tag declaration`) identified as owed only a *fixture*
since the spike already decides them — the cheapest outstanding work, named
rather than left to implementation discovery.

**The one excluded item is properly excluded.** The ambiguous-overlap assertion
is cross-row and therefore RDR 0006's by the arity split, and RDR 0006 is
unimplemented. The RDR states it is "deferred to the Phase 5 lint handshake
rather than counted as satisfied here," names what this RDR still owes before
that handshake (overlapping rows distinguishable by row identity and comparable
by predicate set, which the identity tuple and atom-set contract already fix),
and says plainly that "marking the overlap assertion green before RDR 0006 lands
would be asserting on a stub." A scoped handoff with a named obligation.

**The oracle discipline is the strongest part of this RDR and is not
boilerplate.** "Every oracle must have a failing control" names specific
assertions that *previously passed against a wrong implementation* and replaces
each: the category oracle (was asserting message text), the determinism oracles
(a normalizer emitting a constant passed; now paired with the value assertion
and run on the RDR fixture carrying the escape row, the `in`-expansion,
`<clear>`, and inherited contexts), "the comparator is total" (was satisfied by
the positional tiebreak it existed to exclude; now rule-order permutation), and
the no-alias obligation (not assertable by value comparison, so it is a
**mutation** test — mutate one field in place, assert the other unchanged). The
RDR explicitly refuses "review the normalizer" as an oracle because it names no
wrong implementation that fails it.

**The SHA ban is correctly reasoned, and Stage 6 fixed its rationale.** The
draft previously argued the recorded digest must not be a golden because "a
*correct* implementation necessarily changes these bytes" — and the correct
implementation reproduced them exactly. The ban now stands on the oracle's
*character*: a rendered-text hash is satisfied by any renderer that happens to
agree on these inputs, so it witnesses nothing about the identity rules the
clauses bind. Correcting a rationale the evidence falsified is the same defect
class the rounds were catching, and the RDR records it as such.

**Fixture promotion is `extend-then-promote`, and the two-sibling requirement is
now met.** The pre-§D7 fixtures never had two *ordinary* candidates contending,
which is why the desk trace once recorded "no witness possible in the current
fixture." `continue-prelock-cluster` supplies the MVV's two-sibling requirement
and A12's discriminating witness at once. Those old fixtures are now refused
outright by this RDR's own contracts (`unknown schema field` on `accessors.*`,
`tags.*.accessor`) — the sharpest possible statement of why they could not be
promoted verbatim.

## 4. Cross-Cutting Concerns

**Versioning.** Owned here. `[model].version = 1` is the only accepted version,
and the version gate is normatively ordered **before** strict field decoding —
the only ordering this RDR fixes — so a future v2 document refuses as
`unsupported version` rather than as `unknown schema field` on whichever v2-only
key the decoder reaches first. A14 is now `Verified` by the v2-shaped fixture
and its v1 control, and the two-pass shape it requires was confirmed feasible
against the library rather than assumed.

**Canonical form / determinism.** The core concern and the reason this RDR is
`foundational`. It claims byte-identical dumps, so the sub-clause applies item
by item:

- *Hash function + library* — **deliberately not a contract**; the SHA is spike
  provenance, not an identity commitment, and production goldens must assert the
  normalized value. Nothing downstream content-addresses this.
- *Pre-image byte layout* — the invariant is stated over the **normalized value,
  not the dump text**, because this RDR fixes the field list and row ordering
  but defines no dump grammar.
- *Map iteration order* — explicitly neutralized: rows, atoms, next tags,
  writes, required-owned keys, and escape classes are all sorted slices before
  rendering. Witnessed by three consecutive byte-identical runs despite Go's
  randomized iteration.
- *Ordering* — the identity tuple `(model id, rule id, expansion suffix)`
  compared field by field, byte-lexicographically, with the source locator
  excluded so unrelated source edits cannot churn goldens. The suffix is now a
  **sequence** compared element-wise (`slices.Compare`), since general
  match-block expansion can expand on more than one atom. Byte-lexicographic
  comparison is a Go spec guarantee, so ordering is locale- and
  platform-independent.
- *Whitespace / case folding* — "no case folding, trimming, or namespace
  rewriting is applied at any stage," matching RDR 0008's byte-exact reserved-key
  comparison, with a negative control (`Status` against a declared `status` must
  fail `unknown tag`).
- *Empty / null / absent distinguishability* — the expansion suffix is absent or
  an alphabet member, never the empty string, so the identity tuple is total;
  the empty string as an alphabet member is a load-time refusal; a member or
  rule id containing the suffix separator is its own category, which is what
  keeps `(rule id, suffix)` recoverable.
- *Member-sequence identity* — **new this iteration, and the sharpest
  determinism finding in the set.** A13 established that no member sequence,
  atom literal or set-valued write value, is ever compared by a joined
  rendering. The identity is length-prefixed (`memberKey`, `%d:%s`, injective
  for any member content), and the display clause was raised from SHOULD to
  **MUST** because an unquoted `|` separator spells `["x|y","z"]` and
  `["x","y|z"]` alike.
- *Version marker* — `[model].version`, gated first.

The three named **lossy rendering sites** are disclosed rather than hidden, and
the obligation to close them is charted to a successor dump-format RDR — with
the charting note now stating plainly that no such RDR exists yet, rather than
naming a phantom dependency. Since the contract is the normalized value and the
dump is a review surface, this is a bounded, stated limit.

**Character encoding.** Tag keys and literals are exact byte strings with no
folding, trimming, or namespace rewriting; set-valued literals render with
members sorted so two authored orderings are one literal, and each member is now
delimited unambiguously.

**Incremental adoption.** The table is new data with no existing model to
migrate; `internal/resolve::Row` is reused rather than reshaped by this RDR, and
Phases 2–3 sequence behind RDR 0007's reshape rather than forcing it.

**Build tool compatibility.** `github.com/pelletier/go-toml/v2` (v2.3.1, MIT) is
a candidate only; no production dependency is added until implementation.
Strict decoding is stated as an obligation on the *format*, not a property of
the library — verified at the library's source, where `toml.Unmarshal` is
permissive by default — so a parser swap cannot silently retire the
`unknown schema field` category. The RDR also notes the dependency settles the
open TOML-library `TODO` at `config.go::Load`, and that the two consumers should
share one library.

Peer-owned policies this RDR conforms to rather than authors: predicate operator
grammar and tag type model (RDR 0003), accessor execution and read-back
(RDR 0004), CLI surface and error mapping (RDR 0005), cross-row graph lint
(RDR 0006), guard seam / atom shape / existence constants (RDR 0007), reserved
`recognized` key name (RDR 0008), escape-row shape conformance (RDR 0009).

Concurrency model, secret/credential lifecycle, memory management, licensing,
deployment model, and IDE compatibility do not apply — omitted per the template
rather than N/A-bulleted.

## 5. Proportionality

**Right-sized. Nothing flagged for trimming; lock as written.**

**Contract count — one seam, not several.** This RDR authors exactly one
independent load-bearing contract: *the sparse transition-model format and its
normalization to candidate rows.* The wire schema, the version gate, the
category taxonomy, the identity tuple, and the expanded-table dump order are all
facets of that single seam — none separable, since each is defined in terms of
the normalized candidate row. Every adjacent seam is explicitly delegated
(0003–0009 above). The load/lint **arity split** — single-rule → load here,
cross-row → lint at 0006, evaluation-time → kernel refusal at 0001 — is the
principle that keeps the boundary decidable rather than negotiated case by case,
and `duplicate model id` is correctly called out as the one load-time check that
is cross-*document* rather than single-rule, listed as an owned refusal rather
than smuggled in as single-rule. No split warranted.

**Profile re-validated against the contracts just counted.** `foundational` is
correct and not inflated: RDRs 0003, 0006, 0007, 0008, and 0009 all consume the
wire format, the normalization semantics, the dump's total ordering, or the
validation category taxonomy, and 0008 and 0009 both name this RDR's normalizer
as the enforcement point their own checks land inside. That is the matrix's
cross-RDR-producer trigger on the contract axis, independent of the accretion
axis (`Seam Lineage`: no prior accretion, so no floor applies). The Profile
field carries the value plus one clause naming the contracts, with no matrix or
provenance prose left from the template.

**Lenses match the Profile.** `foundational` owes cove → 3amigo → critique →
repeatability, and all four ran with complete evidence and dispositions at the
current iteration (cove iter-2, 3amigo iter-3, critique iter-3 dual-model). The
repeatability files predate the §D7 re-authoring — the reconcile records this as
a carried caveat rather than hiding it — but §lens-row's completion test is the
run/diff files at the right variant with distinct stamps
(`claude-opus-5[1m]`, `Claude Sonnet 5`, `claude-fable-5`, all `variant: full`),
which is met, and the absorption audit confirmed no clause repeatability touched
was re-opened by the later rounds. Not a Stage-5 return. No `Transient`-marked
contract exists to discount.

**Length is load-bearing, not bloat.** The document is long, and it grew this
iteration — but the growth is where a `foundational` producer must be precise.
Three stretches earn their length by recording a *rejected* reading a reader
would otherwise reconstruct wrongly: the "combined predicate set" exclusion in
outcome binding, the row-kind clause (a render-time view MAY materialize the
derived column but MUST NOT feed it back — the exact clause RDR 0009 now quotes
stalely), and the fail-fast ordering clause (order among independent defects is
*deliberately* unspecified). Each was a repeatability-diff finding where
independent model runs diverged.

**The advisory C9 report, reviewed per assumption as the check asks.** Three
fields over the 30-line budget: A13 (60), A15 (33), A7 (31) — up from one at the
previous lock, and the growth is exactly the Stage 6 close. A13's and A15's mass
is the verification narrative: the defect as found, the fix, the discriminating
witnesses, and — for A13 — why five delimiters rather than one. That last is not
decoration: the pre-fix spike collapsed *only* the comma case, so a
single-delimiter control would have certified the defect, and a pointer-only
field would leave that unrecoverable. Load-bearing anchors
(`main.go::Atom.identity`, `::memberKey`, `::renderValue`, `::checkAtomBlock`,
the `delim/` fixture dir, the preserved repro) remain findable within each field.
Keep as written. Advisory, non-blocking, no action.

**One scope pressure was correctly refused.** 3amigo charted a re-readable dump
grammar to a successor RDR rather than absorbing it — "a new format with its own
grammar, escaping rules, and round-trip invariant — a contract surface, not a
clause." That is the scope-expansion wormhole named and declined.

---

## Verdict

**READY.** No open blockers.

The mechanical sweep passes with no blocking findings. All four `foundational`
lenses ran with dispositions at the current iteration. Stage 6 closed the two
assumptions whose wide arm changes match semantics while tripping zero load
categories — the arms where the owed fixture was the only detector — rather than
deferring them, and the draft locks with no clause its own evidence contradicts.
The three remaining `Pending` assumptions are structural: two blocked on RDR
0007's unimplemented kernel reshape, one on a multi-document surface that does
not exist, each carrying a named MVV oracle and a survivable `If wrong`. The MVV
is in scope with its one cross-row assertion correctly handed to RDR 0006, and
the RDR owns exactly one contract at the `foundational` profile its consumer set
requires.

**Implementation readiness, stated plainly.** Locking does not mean Stage 8 can
open on this RDR alone. Prerequisites names an unowned blocker in its own words:
RDR 0007's kernel reshape is unimplemented, `Final` (so 0007 will not schedule
it), tracked by no kata, and claimed by no phase here — "sequencing behind an
unowned prerequisite is indefinite, not merely ordered." Phases 2–3 sit behind
it in their entirety, and implementation MUST NOT open Phase 2 by mirroring the
constants locally. Landing the reshape, or explicitly assigning it, is the first
implementation action. Separately, critique's M-5 is a live cross-RDR citation
staleness (RDR 0009 quotes deleted RDR 0002 text, and a `Verified` peer
assumption reasons from it) that only `/rdr-cluster-reconcile` can close, and
the existing 0002-0009 iteration-3 gate stamp did not catch it. Both are
implementation-sequencing facts, not lock gates — the RDR states each one — and
both are why the next command is cluster-reconcile rather than implement.

Status → **Final**.
