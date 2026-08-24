Model: claude-opus-5[1m]

# 3amigo Dispositions — RDR 0002, iteration 3

Origin ledger: `iter-3/consolidation.md` — hotspots HOT3-001…HOT3-005 plus four
single-persona findings, traced back to PM3-001…004 / IMP3-001…008 /
QA3-001…008 (20 findings). One line per finding. Delta-scoped to the §D7
re-authoring and the cove iter-2 edits; the iter-1 and iter-2 ledgers stayed
closed and no persona re-opened them.

## Fixed

- **fixed** — IMP3-001 (HOT3-001 adjacent; closed layout omits `[[rule]].source`
  and `[model].description`). Grounded: both fixtures author `source` on all
  seven rules and `description` on both models; `main.go` maps them
  (`Rule.Source`, `Model.Description`). The clause said "No other root key or
  table is admitted", so the normative text refused its own normative fixtures.
  Admitted both keys explicitly and said why the enumeration must be exhaustive
  over what the fixtures author. Section: Normative Contracts / closed-layout.

- **fixed** — IMP3-002 + QA3-008 (HOT3-004; locator provenance and its
  Round-Trip projection). Grounded: `main.go:670` sets `Locator: rule.Source`,
  and both `kata-fixture.toml` rules author `source = "kata:review"`, so two
  rows carry one locator against a floor requiring model id + rule id. Made the
  locator **derived** from `(model id, rule id)`, demoted authored `source` to a
  provenance annotation, named the spike's behavior a defect against the clause,
  and defined the Round-Trip "rule-identifying part" as the `(model id, rule
  id)` projection. Sections: Technical Design / locator floor; Round-Trip.

- **fixed** — IMP3-003 (HOT3-001; set-valued *write* values joined). Grounded:
  `main.go::renderValue` does `strings.Join(parts, ",")`, so `["a,b","c"]` and
  `["a","b,c"]` collide in a value RDR 0004 compares for equality on read-back.
  The literals clause bound atoms only. Added the write-value sequence rule to
  the write-replaces clause, plus a scenario 2 control asserting value
  inequality. Sections: Normative Contracts / write-replaces; Testing Strategy 2.

- **fixed** — PM3-001 (HOT3-001; A13's `If wrong` contradicts the desk trace).
  A13 said the failure carries duplicates "without changing match semantics";
  the trace's `merge atoms` row records a live expanding row where the author
  wrote a dead rule. Prerequisites rested the lock on that text. Split `If wrong`
  into its narrow-key (benign) and wide-key (semantics-changing, zero-category)
  arms, and rewrote the Prerequisites survivability line to make A13 the named
  exception. Sections: Critical Assumptions A13; Prerequisites.

- **fixed** — IMP3-008 (HOT3-002; fail-fast vs the two-pass gate). Reaching
  `[model]` requires a parse, so a malformed **and** v2 document must refuse
  `malformed TOML` while the gate's precedence read as absolute. Added a
  "Syntax precedes the gate" paragraph fixing the order `malformed TOML` →
  version gate → strict decoding → the rest, and scoped A14's oracle to the
  gate-vs-`unknown schema field` discrimination. Section: version-gate clause.

- **fixed** — PM3-003 (HOT3-002; fail-fast forecloses batch diagnostics
  unweighed). The Decision Rationale names authoring ergonomics as the RDR's
  purpose; the clause fixed single-category refusal over 25 categories without
  weighing it or reserving an accumulating mode. Reserved a batch mode
  explicitly as a lint-side (RDR 0006) or repeated-load RDR 0005 surface, and
  narrowed the prohibition to the single `load` entry point's return type.
  Section: Normative Contracts / fail-fast.

- **fixed** — QA3-007 (HOT3-002; absorbed-defect class has no substitute
  oracle). The carve-out named the class and routed coverage to the
  merge/literals clauses, whose control is the one QA3-001 shows is owed.
  Stated that this class is asserted **positively** (a count over the normalized
  value), never by refusal, and made that a standing requirement for any future
  absorption-preventing clause. Section: fail-fast clause.

- **fixed** — PM3-004 + QA3-002 (HOT3-003; `[model.metadata]` has no consumer
  and no oracle). Grounded: no peer RDR references metadata; JDR 0001 §D7(v)
  names no tool. Stated that reserving the namespace pre-consumer is the point,
  then made "carried through untouched" testable: the normalized model MUST
  expose the field (Phase 2 deliverable) and scenario 1 MUST assert **deep
  value/nesting equality**, since a key-name oracle passes a value-discarding
  loader. Sections: `[model.metadata]` clause; Testing Strategy 1.

- **fixed** — IMP3-005 (HOT3-005; `[model]` well-formedness has no category).
  Grounded: the spike invents `malformed model: missing id` off-contract, and an
  absent `version` yields the zero value so it refuses as `unsupported version
  0` — naming a version the author never wrote, the confusion A14's gate exists
  to prevent. Added `malformed model declaration` to the floor covering absent
  `[model]`, `id`, `version`, and `[tags.recognized]`, with an explicit ban on
  folding absent-`version` into `unsupported version`. Section: category floor.

- **fixed** — IMP3-007 (HOT3-005; `[dump]` vocabulary and malformed forms
  unspecified). Grounded: the spike decodes `Dump{Order []string}` and never
  reads it, so both fixtures' ten-element lists are dead text. Fixed `order` as
  the single key, closed the column vocabulary to the snake_cased field names
  already enumerated in the dump clause, and added `malformed dump declaration`
  to the floor for unknown/repeated/omitted columns. Sections: dump clause;
  category floor.

- **fixed** — IMP3-006 (HOT3-005; accessor entry states a rule only for
  `timeout`). Grounded: the spike never checks `role` or `path`, checks
  `timeout` only for `== ""` (so `"0s"` and `"banana"` both load despite the
  clause's stated non-positive rule), and never checks `read_back`. Made all
  four fields required with a per-field refusal rule, and bound `read_back` on
  write entries. Section: Accessor tables clause.

- **fixed** — IMP3-004 (HOT3-005; guard-block operators unconstrained while
  `unknown operator` is a floor category). Grounded: RDR 0003 declares the
  closed set `{eq,in,lt,lte,gt,gte,exists,contains}` but this RDR nowhere cited
  it by member list, and the spike takes any TOML key as an operator token.
  Cited the set by members so `unknown operator` is decidable here, kept
  semantics with RDR 0003, and named this the amendment site should 0003
  re-lock differently. Section: match-block operator clause.

- **fixed** — QA3-001 (HOT3-001; the merge control has no determinate arm).
  Both arms failed: "the implementation's own joining delimiter" is unreadable
  black-box, and "dead rule for lint" is an RDR 0006 verdict the MVV refuses to
  count (the overlap item). Rewrote the control to assert the **atom count** —
  a property of this RDR's own normalized value — parameterized over the
  plausible delimiters, and explicitly barred stating it as a lint verdict.
  Section: Testing Strategy 2.

- **fixed** — QA3-003 (HOT3-005; `[initial]`-undeclared-key control assigned the
  wrong category). Computed against the transcript: `neg-initial-unknown`
  refuses `unknown tag "nosuchtag" written`, not `malformed initial
  declaration`. Scoped that category to its value arm only and reassigned the
  undeclared-key variant to `unknown tag`, with the reason. Section: scenario 3.

- **fixed** — QA3-004 (HOT3-005; three escape-shape controls share one
  category). Computed: `neg-escape-with-write/clear/gate` emit one identical
  string, so the category oracle cannot separate them. Kept one category —
  the three shapes are one authoring error and RDR 0009 binds it as one
  obligation — and made the discriminator the one-mutation construction rather
  than message text. Added the `gate` shape, which the prose had omitted.
  Section: scenario 3.

- **fixed** — QA3-005 (witness counts stated three ways). Computed from
  `negative-cases.txt`: 37 refusal lines, of which 2 are `probe-*` enumeration
  witnesses → **35 category controls**; 3 positive loads + 2 `-strict-only`
  checks. Scenario 3's "35" was right; A1's "32 negative + 2 positive" was
  stale. Corrected A1 and stated the probe exclusion. Section: A1 Evidence.

- **fixed** — QA3-006 (probe fixtures cited by non-existent paths). Verified the
  real paths are `neg/probe-a-initial-observed-no-writer.toml` and
  `neg/probe-b-initial-observed-with-writer.toml`. Corrected. Section:
  scenario 3 probe paragraph.

- **fixed** — PM3-002 (the review surface is lossy and charted to a phantom
  successor). The `fidelity` read-back row charted to "a successor dump-format
  RDR" named nowhere in Capability Dependencies, References, or the phase plan.
  Answered the blocked decision in place — Phase 2's dump is **not** required to
  carry an escaping grammar, because that would make the dump a re-readable
  format whose inverse this RDR does not claim — while requiring set-valued
  fields to render visibly as sets so the ambiguity is not silent. Corrected the
  `fidelity` row to stop naming a dependency that does not exist. The successor
  is charted, not assumed. Sections: Round-Trip; `fidelity`.

## Charted to successor

Recorded durably in `iter-3/Charted.md`:

- **A re-readable expanded-table dump grammar** (PM3-002 residue). Out of scope:
  a new format with its own grammar, escaping rules, and round-trip invariant —
  a contract surface, not a clause. Substantial enough to warrant `/rdr-seed`;
  surfaced in the close packet.
- **Spike artifacts carry a stray `Model:` evidence stamp** (QA persona
  observation, recorded as a non-finding). Evidence hygiene touching no contract
  here; strip at promotion or file a kata.

## Dismissed with cite

None. All 20 findings passed the grounding gate. Every code-level claim was
confirmed at source or against the committed spike before it was actioned —
`Rule.Source`/`Model.Description` in the fixtures, `renderValue`'s join,
`Locator: rule.Source`, the unread `Dump.Order`, the unchecked accessor fields,
the operator token passthrough, and the three counts recomputed from the
transcript rather than taken from any file's prose. None re-litigates an
adjudicated decision.

## Needs (re)verification — carried to Stage 6

1. **A13 is widened** (still `Pending`) to cover set-valued *write* values, not
   only atom literals. Its subject is now "a member sequence is never compared
   as a joined string". Two owed witnesses, same shape: the delimiter-bearing
   atom control (Testing Strategy 2) and the new write-value control. The spike
   fails both as committed.
2. **The derived-locator rule is a new load-bearing claim** (IMP3-002). It
   changes the locator from authored to derived and defines the Round-Trip
   projection. Needs checking against RDR 0006, which consumes normalized rows
   and may key diagnostics on the locator, and against RDR 0004 if any accessor
   diagnostic cites it. The spike contradicts it and must be corrected before
   the fixtures promote.
3. **Two new load categories** — `malformed model declaration` and `malformed
   dump declaration` — are new API-surface identifiers. Both are owed a fixture
   *and* an implementation; the witnessed count moved to 18 of 25, with seven
   owed. `gen-cases.py` should mint their single-mutation cases.
4. **The guard-operator membership citation** pins RDR 0003's closed set
   `{eq,in,lt,lte,gt,gte,exists,contains}` into this RDR. RDR 0003 is currently
   `Draft [revised from Final 2026-08-24]`, so this citation must be re-checked
   when 0003 re-locks; the clause names itself the amendment site.
5. **The accessor per-field refusal rules** (IMP3-006) add four obligations the
   spike does not implement (`role`, `path`, `timeout` parsing, `read_back`).
   Owed fixtures under `malformed accessor declaration`.
6. **The `[model.metadata]` Phase 2 field + deep-equality assertion** (PM3-004,
   QA3-002) is a new deliverable and a new oracle.
7. **The dump's visible-set rendering requirement** (PM3-002) is a new rendering
   obligation; it has no control yet and the spike renders joined.

No previously `Verified` assumption was invalidated this pass. A13 was already
`Pending` from cove iter-2 and was widened, not demoted anew.

## Tiebreakers

None escalated. Two forks were collapsed on evidence rather than surfaced:

- **QA3-004** — split `malformed escape declaration` into three categories, or
  keep one? Collapsed against the draft's own arity logic and RDR 0009, which
  binds the write-free obligation as a single kernel-boundary obligation: the
  three shapes are one authoring error, so the discriminator belongs in the
  fixture construction, not the category vocabulary.
- **PM3-002** — require a dump escaping grammar in Phase 2, or accept the lossy
  render? Collapsed on the RDR's own scope: defining a grammar creates the
  inverse the RDR explicitly declines to claim, and every consumer needing
  recoverability holds the normalized value. The residual visibility duty is
  cheap and was added; the grammar is charted.

## Mini-checks

No new cue fired. This pass edited four tables the cove pass had already
established — `fidelity` (dump row: grammar-by-decision), `disposition` (two new
load categories), `oracle` (three new controls; the absorbed-defect oracle
rule), and `trace` (a second CONTRADICTION row for write values) — none of which
introduces a mini-check the draft was not already carrying. The Stage 5
mini-check cue read was discharged by this RDR's first lens pass this iteration
(cove iter-2); tables persist in the draft.
