Model: claude-opus-5[1m]

# Critique Diff — RDR 0002, iter-2 (A: opus-5 15 rows / B: sonnet-5 10 rows)

Reconciled by RDR passage anchor, not by `C-N` id. Grounding verdicts are
independent, checked against the draft text and, where cited, against
`evidence/spikes/main.go` and `output.txt`.

## Merged Ledger

| M | Sources | RDR passage anchor | Failure mode (one line) | AGREEMENT | GROUNDING |
|---|---------|--------------------|--------------------------|-----------|-----------|
| M-1 | A:C-1 | §Technical Design "The normalized candidate row is the kernel row… predicate atoms (key, operator token, literal, block)"; §Prerequisites "`Row.Guard` is a `string`, `Row.Match` is a flat `[]Tag`" | No clause says how the unified atom set splits back into `Row.Match` vs `Row.Guard`; the normalizer's central step is unspecified. | A-only | grounded |
| M-2 | A:C-2 | §Normative Contracts, outcome-binding: "its combined predicate set (local match block plus inherited contexts) MUST contain exactly one atom on `recognized`" vs guard clause "combine both into one candidate-row predicate set" | "Combined predicate set" is defined twice with different extents; the only implementation (`liftOutcome`) lifts from the post-merge map with no block filter, so `unless.recognized` inverts intent. | A-only | grounded |
| M-3 | A:C-3 | §Normative Contracts, "Shared contexts MAY inherit… MUST normalize to an explicit predicate set" | No cardinality/override rule for same-key atoms; spike's `map[block\x00key\x00op]` makes inheritance last-write-wins, order-dependent on `use = [...]`, silently. | A-only | grounded |
| M-4 | A:C-5, B:C-2 | §Normative Contracts, "`[model]` MUST contain `id` and `version`. Version `1` is the only version this RDR accepts; any other version MUST be refused before normalization." | Version gate refuses but specifies nothing else. **A** attacks *ordering* (strict decode runs first, so a v2 file fails `unknown schema field`, not `unsupported version`); **B** attacks *migration* (no v2 story, no coexistence, no CLI verb). | both | A's ordering half: **grounded**. B's migration half: **scope** |
| M-5 | A:C-6, B:C-4 | §A8 residual, "The recognized-totality clause below is the sole barrier" + §Normative Contracts recognized-total clause | The alphabet ban is a load-time check in a package `internal/resolve` does not import, so it is a convention, not a barrier. **A** frames it as any non-normalizer `Table` producer; **B** frames it as the CLI passing empty `--recognized`. | both | grounded |
| M-6 | A:C-7 | §Normative Contracts, `RequiresOwned`: "the sorted, duplicate-free set of tag keys named by the rule's write block and clear list… does not add guard-read keys to it" | Excluding guard-read keys inverts the kernel's deliberate `missingOwned`-before-undecidable diagnosis ordering; author sees `guard_unevaluable` where `owned_state_unavailable` is the true cause. | A-only | grounded |
| M-7 | A:C-8 | §Round-Trip, "with the source locator's optional line/column detail excluded from the comparison" | Invariant includes the locator then excludes its only varying part; residue is `(model id, rule id)`, already in the identity tuple, so a locator-dropping normalizer passes. | A-only | grounded |
| M-8 | A:C-9 | §Normative Contracts, next-state/writes: "the writes… — the same rendered set" | Two kernel fields contractually forced identical; `output.txt` shows `next=` byte-equal to `write=` on all seven rows. A contract that mandates equality has not modeled a distinction. | A-only | grounded |
| M-9 | A:C-10 | §Normative Contracts, tag-key identity: "no case folding, trimming, or namespace rewriting" | Byte-exact identity is specified for keys and left unspecified for literals; `renderLiteral` space-joins sorted set members, so `in = ["needs work"]` and `in = ["needs","work"]` are one literal. | A-only | grounded |
| M-10 | A:C-11 | §MVV "two sibling candidate rows binding the same outcome" vs §Load-Bearing Decisions "The RDR and kata spike fixtures are the canonical examples implementation tests must promote" + §trace "no witness possible in the current fixture" | Two clauses of one document require contradictory things of one file; nothing forces the canonical fixture to grow sibling rows, so scenario 4 ships untested. | A-only | grounded |
| M-11 | A:C-12, B:C-10 | §Normative Contracts escape clause + §Illustrative Code `draft-no-match-escape` | Escape rescue is per-outcome (kernel filters escapes on `row.Outcome != in.Recognized`), so an N-outcome flow needs N escape rows. **A**: never stated, fixture rescues 1 of 4. **B**: never stated whether an escape row may use `in` to expand across outcomes. | both | grounded (see Contradictions) |
| M-12 | A:C-13 | §Existing Infrastructure Audit "Add stable parse/lint refusal codes later — these are the first of their kind" + §A5 | ~17 stable data-level categories mandated with no CLI code mapping and no Final owner (RDR 0005 Pending); stable categories with no surface are internal constants. | A-only | scope |
| M-13 | A:C-14 | §Normative Contracts dump clause, "`[dump]` settings MAY reorder the rendered columns; they MUST NOT omit a field" | A render preference lives inside the semantic source file, and the round-trip invariant excludes dump text, so `[dump].order` churn is unconstrained. | A-only | grounded |
| M-14 | A:C-15 | §Normative Contracts, "an absent expansion suffix sorts before any present one" + §Round-Trip "the expansion suffix has no reserved separator" | Rule id `a#x` (unexpanded) and rule `a` expanding on `x` produce one identity tuple; the duplicate check compares rule ids, not identities, so the "total" comparator returns 0 for distinct rows. | A-only | grounded |
| M-15 | A:C-4 | §trace table, `expand in` row: "the single-member case is **unwitnessed** (spike suffixes on operator, not expansion count)" | A claims the spike contradicts the contract and the RDR books it as needs-verification rather than a defect. | A-only | **refuted** (see Refuted) |
| M-16 | B:C-1 | §Round-Trip "Rendering is a review surface here, and the normalized value is the contract" vs §Problem Statement "one reviewable artifact" | Two incompatible mental models of the dump: sold as the review surface, then disclaimed as having no grammar with three named lossy sites (`<clear>`, separators, suffix). | B-only | grounded |
| M-17 | B:C-3 | §Normative Contracts dump clause, "The identity tuple… is total **over one model's rows**" | Nothing constrains ordering across separate model files a real repo will have; a reviewer wanting "the table" spanning RDR + kata flows has no promised order. | B-only | **refuted** (see Refuted) |
| M-18 | B:C-5 | §Normative Contracts guard clause, "Normalization MUST combine both into one candidate-row predicate set" | Set semantics undefined for the same atom authored in both `all` and `unless`, or twice in one block; no load category covers "duplicate atom" or "atom in both blocks", so a self-contradictory guard loads and silently prunes. | B-only | grounded |
| M-19 | B:C-6 | §Normative Contracts `RequiresOwned`, "This is a **producer obligation, not a kernel guarantee**… the normalizer is the only enforcement point" | Single point of failure documented with no defense in depth, and the only adversarial coverage cited is the kernel's `TestAdv2b`, not this RDR's own MVV; scenario 3 enumerates the write-block escape variant but not the clear-list-only variant. | B-only | grounded (narrow half) |
| M-20 | B:C-7 | §New Dependencies, "Use `github.com/pelletier/go-toml/v2` as the TOML parser candidate… No production dependency is added until implementation" vs §MVV "the decoder must be configured to reject unmapped keys" | Field layout and strict-decode behavior are locked as normative while the library they depend on is called a non-final "candidate". | B-only | grounded (minor) |
| M-21 | B:C-8 | §MVV / §Testing Strategy scenario 4, "Separately, run RDR 0006's lint over a deliberately overlapping variant" + §Capability Dependencies (RDR 0006 = Pending) | This RDR's own MVV has a hard dependency on a Pending peer, so "MVV executed during implementation, not deferred" is unsatisfiable as written. | B-only | grounded |
| M-22 | B:C-9 | §Prerequisites, "Phases 2 and 3 therefore sequence behind the reshape in their entirety… This item gates implementation sequencing, not lock." | B claims the RDR locks a design whose Phase 2/3 work cannot be attempted, with no described interim state. | B-only | **refuted** on the lock claim; **grounded** on the narrower "no specified interim shim" residue |

---

## Hotspots

Findings both models reached independently, ranked. These are the passages where
two readers with no shared context converged — the strongest signal in this diff.

**H-1 — M-5: A8's "sole barrier" is not a barrier.** A:C-6 and B:C-4 arrive from
opposite directions and land on the same sentence. A reasons *downward* from the
kernel: `internal/resolve` accepts a `Table` from any caller and validates
nothing, so a load-time check in a package it does not import cannot enforce
anything. B reasons *upward* from the CLI: `Input.Recognized == ""` is still
ungoverned at the call boundary, and the draft itself concedes "JDR 0001 §JD-10…
remains open." Both are right, and they are the same defect at two ends of one
seam. The draft's own phrasing — "**Residual, and the reason this RDR is
load-bearing for it**" — shows it knows the hole exists; what it does not do is
name an enforcement point that the kernel actually reaches. This is the highest-
confidence item in the ledger.

**H-2 — M-4: the version gate.** Both models attack the identical sentence, but
they are *not* the same finding, and the split matters. A's half is an ordering
defect: "MUST be refused before normalization" is satisfied by refusing after
full strict decode, so a v2 file with a renamed section dies as `unknown schema
field` — a category collision the draft's own taxonomy separates. That is a real,
in-scope, one-clause fix. B's half is a migration story, which is net-new scope
for a v1-only format RDR (see Refuted note). Take A's ordering clause; leave B's
migration request to a successor.

**H-3 — M-11: escape rescue is per-outcome.** A:C-12 and B:C-10 both notice that
the escape mechanism is under-described, and both trace it to the same fixture
(`draft-no-match-escape` binding one outcome out of four). The draft states the
kernel behavior *once*, buried in A8's evidence — "`::escapeOrRefuse` applies the
same filter to escape rows" — and never carries it forward to the escape clause
where an author would look. See Contradictions: the two models want opposite
fixes.

---

## Contradictions

Passages where the models pull in opposite directions. Each is a genuine
underspecification signal — the prose admits both readings.

**X-1 — escape-row expansion (M-11).** A treats per-outcome escape scope as
*under*-specified and wants the RDR to say "an N-outcome flow needs N escape
rows." B treats it as *under*-specified in the opposite direction and asks
whether one escape row may use `recognized.in = ["a","b","c"]` to cover three
outcomes at once — which, if permitted, makes A's N-row requirement unnecessary.
The draft supports both readings: the outcome-binding clause says "Every rule —
**ordinary or escape** — MUST bind exactly one outcome" and that "an `in` atom
expands into one candidate row per member," which grants B's expansion; but no
sentence anywhere tells an author that one authored escape rule is required per
outcome to get full coverage. One clause resolves both: state that escape
expansion is permitted and that escape coverage is per-outcome, so `in` is the
idiomatic way to spell a flow-wide escape.

**X-2 — the dump clause: over- vs under-specified (M-13 / M-16 vs A §2 / B §2).**
Both models independently nominated a dump-adjacent section as "the one section
rewritten within six weeks," and they named opposite sections for opposite
reasons. A nominates the **dump ordering clause** as *over*-specified: 30 lines
fixing sort keys to the byte for an object (`Row.Guard`, next-vs-writes) that is
not yet pinned down. B nominates the **Round-Trip / fidelity section** as
*under*-specified: it names three lossy sites and ships anyway, defending a
defect in the artifact the RDR is titled after. These are the same tension read
from two ends — the render contract is simultaneously too precise about ordering
and too vague about grammar. That is the strongest structural signal in this
pass, and it is not addressable by editing either clause alone.

**X-3 — MVV self-containment (M-21) vs Prerequisites sequencing (M-22).** B
raises both, and they cut against each other: C-8 says the MVV is *not*
self-contained because it needs RDR 0006, while C-9 says the RDR should have
specified an interim shim so work can proceed *without* its dependencies. One
asks for stricter dependency hygiene, the other for a documented way to work
around dependencies. Only the first is grounded as a defect in this draft.

---

## Refuted

Findings whose cited passage does not say what the finding claims, or that the
draft already decides explicitly. Deciding sentence quoted in each case.

**R-1 — M-15 (A:C-4), "the spike contradicts the expansion-suffix contract."**
A asserts the spike "suffixes on `len(outcomes) > 1`, which is the *right* rule"
while the RDR trace says otherwise, and books this as the RDR ducking a defect.
The spike source refutes the framing. `main.go::normalize`:

```go
suffix := ""
if len(outcomes) > 1 {
    suffix = outcome
}
```

That is *exactly* the contract: "A suffix is present exactly when the rule
produced more than one row." Spike and contract agree. What is actually wrong is
the trace table's parenthetical — "(spike suffixes on operator, not expansion
count)" — which is a stale note describing an earlier spike revision. The finding
as written (contract contradicted by its only implementation) is refuted; the
residue is a one-line stale trace note, not a design defect. Note this
*strengthens* the draft: the single-member `in` case is witnessed by construction,
not "unwitnessed."

**R-2 — M-17 (B:C-3), "nothing constrains cross-model ordering."** The draft
decides this explicitly, in the same clause B cites:

> "A dump spanning more than one model MUST carry a model-unique `model id`,
> which is what makes the leading tuple field discriminating; two models sharing
> an id is a `duplicate model id` load failure. RDR 0006 lints exactly one model
> per invocation."

The identity tuple leads with `model id`, so a multi-model dump is ordered by
construction, and `duplicate model id` is already in the load-category list.
B's symptom ("an unspecified concatenation order") is directly answered.

**R-3 — M-22 (B:C-9), "the RDR locks a design whose Phase 2/3 cannot be
attempted."** The draft decides the lock question verbatim in the sentence B
quotes:

> "This item gates implementation sequencing, not lock."

The Prerequisites bullet is unchecked *on purpose* and says so. B's broader claim
— that this blocks lock — is refuted. The narrow residue (no normative interim
shim; only a Testing Strategy aside says "the spike mirrors the constants locally
because the kernel does not yet export them") is a real but minor sequencing
convenience, not a lock blocker.

**R-4 — M-4, B's half (B:C-2), "the version gate lacks a migration path."**
Migration is not in scope for a format RDR that accepts exactly one version. The
draft states the acceptance set and the refusal, which is the complete v1
contract:

> "`[model]` MUST contain `id` and `version`. Version `1` is the only version
> this RDR accepts; any other version MUST be refused before normalization."

A v2 migration story requires a v2 schema to migrate *to*, which does not exist
and is not this RDR's to invent. Classified `scope`. A's ordering half of the
same passage (M-4) survives and is grounded.

**R-5 — not raised, checked as instructed: ambiguous overlap as a load failure.**
Neither model raised this. Confirmed the draft decides it explicitly:

> "Load and lint split on **arity, not severity**… Ambiguous overlap between
> candidate rows is therefore a lint finding, not a load failure."

**R-6 — not raised, checked as instructed: `[model].version` inert.**
Neither model claimed the version field is inert. Confirmed present and load-
bearing in the draft (version clause, `unsupported version` in the category list,
MVV "one unsupported-version variant is refused before normalization", and
witnessed in `negative-cases.txt`). The iter-1 fix holds.

---

## Coverage note

A found 15, B found 10, overlap is 3 (M-4, M-5, M-11). Overlap is low because the
models read at different altitudes: A reads *downward into the kernel* — nine of
its fifteen rows are grounded in `internal/resolve` or spike source and would not
be visible to a reader who stayed inside the document. B reads *across the
document* — its unique rows are internal-consistency and sequencing defects
(MVV vs Capability Dependencies, dependency hedging, dump framing vs Problem
Statement) that require no source reading. Neither pass subsumes the other, and
the union is the useful artifact.

Of 22 merged findings: 17 grounded, 3 refuted, 2 scope.
