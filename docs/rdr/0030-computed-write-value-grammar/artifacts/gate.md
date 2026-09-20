# Finalization Gate — cli/0030 computed-write-value-grammar

- **Record**: `0030-computed-write-value-grammar.md`
- **Date**: 2026-09-20
- **Verdict**: **PASS** — locked to Final
- **Profile**: `foundational` (re-validated this pass; see item 5)
- **Mechanical pre-sweep**: PASS
  (`evidence/tooling-pass/tooling-pass.md`; `recs lint --locking 0030`
  exit 0, `blocking=0 resolution=0`) — after two in-pass C9 relocations,
  recorded there.
- **Joint-decision fence**: clear (`op = none`, `fence-clear`) —
  `overlap_uncited=0`, `rulings_open=0`, `clustered=false`,
  `joint_check_home=homed` over all three joint checks.
- **Repeatability-lite**: `emit.next = none` (`repeatability-full-complete`;
  variant `full`, runs 1-3 and the diff written, `Determinacy: fired`).

Item 4 (Cross-Cutting Concerns) is authored in the record at
`cli/0030:G-cross-cutting` and retained there at lock, because peer RDRs
cite it. It is deliberately not copied here.

## 1. Contradiction Check

No contradictions found between research findings, design principles, and
proposed solution.

Compared: Research Findings (Investigation, Key Discoveries) against
Approach, Problem Statement, the three Normative Contracts C1-C3, Load-
Bearing Decisions, Decision Rationale and Consequences.

Each of the four in-repo shipped facts the research established maps
forward into the design cleanly: `normalize.go::expand` already mints
suffixed rows per member → C1's step choice point; coverage is proved over
guard atoms only → C1's emitted `guard.all eq = <cell>` and its
exclude-on-undecidable arm; a literal write outside its declaration already
refuses under `malformed_tag_declaration` → C2 reuses that category and
mints none; `TagDecl.Domain` is carried in authored order → C1's enum step
order and Consequences' "domain order acquires meaning".

**Three deliberate divergences, each stated in-record with its reason** —
recorded here as divergences, not conflicts, because in every case the
record says it is diverging and says why:

1. *Prior art vs. chosen mechanism.* Investigation found that every peer
   with extended state answers "count in an action" with a host expression
   language, and that no peer carries a data-only computed write or decides
   the bound declaratively. The record diverges deliberately: ALT3 /
   Approach 4 is rejected in the Decision Rationale on rows 1 and 3
   ("undecidable in general"; "models come to depend on it") as what
   RDR 0002 chartered against, and the record states plainly that the bound
   disposition "rests on in-repo principle (`0002:C3` via the
   `valueMembers` comment) rather than external prior art". A design that
   knows it is unlike its prior art and argues the point is not in conflict
   with the finding that surfaced it.
2. *`0002:C13`'s suffix-iff.* C13 states a suffix is non-empty exactly when
   the rule produced more than one row. C1 NARROWS that iff to the `in`
   expansion it was written about, appending the suffix element at every
   admitted cell including a single one. The reason is in the clause: C13's
   iff exists because `eq = "x"` and `in = ["x"]` are two spellings of one
   edge, whereas a stepped row is never the authored row. `D-identity`
   repeats the fence, and the narrowing is declared on the record's
   Overrides field (`0002-transition-table-as-reviewable-data:C13`), so it
   is a recorded override rather than an unannounced contradiction.
3. *`lint`'s per-row `rule` slot.* C3 narrows that slot's published meaning
   from "the authored rule id" to "the row's identity", calls it a gap this
   record closes rather than inherits, and fences the one in-package
   read-back (`coverage.go::groupHasOverlap`, joining on the authored id
   recovered by truncating at the first `#`). S11 pins the invariance on
   the export and the CLI previews. This is the disposition ruling Q1
   settled, and it is the reason A8 stands `Refuted` rather than being
   quietly dropped: the refutation is absorbed into C3 and the Decision
   Rationale rather than left to contradict them.

**Two residuals, stated by the record and not papered over** — neither is a
contradiction because the record concedes each in the same breath as the
benefit. Consequences concedes the enum arm "reduces … but does not
eliminate" the Problem Statement's "every change to the tier domain means
re-typing" — the write is said once, the cap is not. And it concedes an
author converting a hand-unrolled ladder has "no supported way to confirm
the conversion was faithful", marked "Accepted, not closed", with the MVV's
identical-findings comparison named as a test oracle rather than a product
surface. Per-row identity in a flow payload is likewise disclosed in the
Decision Rationale as a possible successor record — "a disclosed
non-obligation, not an owed fix".

## 2. Assumption Verification

Sixteen Critical Assumptions. Ten `Verified`, five `Pending`, one
`Refuted`. Every record is internally consistent — Status, Method and
Evidence agree, and no "If wrong" is empty.

**Method vocabulary**: `ca_off_vocabulary=0`. Nine `Source Search`, five
`MVV Test`, one `Spike` (A7), one `Peer RDR` (A9, citing an element, not a
bare record). No record lacks a Method field.

**No Docs-Only.** No assumption carries `Docs Only`, so the gate's
Docs-Only-on-load-bearing bar is vacuous here rather than waived.

**No self-reference.** None of the nine `Source Search` records resolves
its Evidence to this record or to a path under its artifact directory.
Paths into the RDR's own reconcile evidence are Stage 6 verification
artifacts, which is where they belong.

**Anchors resolve.** All 94 `source-anchor` edges report `resolved: true` —
not ABSENT, so the check genuinely ran against the source root rather than
reporting that nothing looked. Every cited `path::Symbol` is on `main`.

**Two Evidence fields were relocated at this gate, not truncated.**
`lint --locking` marked A12 (32 lines) and A15 (48 lines) BLOCKING under
the `foundational` rule. Judged on the check's own question rather than by
reflex: the mass is real verification content and every load-bearing anchor
is a `path::Symbol`, so truncation was refused and "flagged, accepted" was
not taken as a disposition. Both narratives are reproduced verbatim in
`artifacts/evidence-a12-a15.md`; each field keeps its claim, its anchors
and a pointer. Re-lint after the move: `blocking=0`.

**The five `Pending` records — disposition and status consistency.**
A5, A6, A10, A11 and A13 are all `MVV Test` and all unrunnable before
Phase 2 by construction: each asserts a property of rows that do not exist
until the expansion is built. Stage 6 dispositioned them as pinned, not
deferred, and each is pinned to a named authored scenario that must pass
before the implementation is accepted:

- **A5** (no new overlap or coverage gap) → MVV items 2 and 5; S1 at both
  tiers ("`lint --as json` over `ladder-literal.toml` and
  `ladder-step.toml`, the same ladder authored unrolled and stepped"),
  S5b (`unless` excluding an interior cell another row claims literally)
  and S5c, the positive-atom overlap control that exists because the
  record itself observed S5b's zero-overlap half is vacuous.
- **A6** (plan parity) → MVV item 3; S2 ("`flow resolve --as json
  --plan-only` swept over every (state, outcome) cell of the pair"), with
  S2b for unclaimed-cell parity.
- **A10** (subsumption) → S8 ("subsumption. A step rule on `attempt` that
  ALSO authors `match.attempt in = [1,3]`, and a sibling fixture where that
  same `in` is inherited from a context rather than authored on the rule").
- **A11** (both write carriers) → S10 ("both-carriers. `graph` export and
  `dump` of `ladder-step.toml`, read for each expanded row's successor"),
  named as the discriminating oracle: a `Writes`-only implementation passes
  every `Fingerprint` and plan assertion and fails here.
- **A13** (bound before render) → MVV item 5; S5's representable-width
  `int` whose step magnitude carries an admitted cell past `math.MaxInt`,
  whose expected detail "names the cell and the bound it passed, NOT 'is
  not an int' — that is the ordering A13 exists to pin".

A10 and A11 pin to Testing Strategy scenarios rather than to one of MVV's
seven numbered items. That is the correct shape, not a gap: both scenarios
are in-scope Phase 3 work that lands the same fixture pair, and the
assumptions name them as the plan.

**Status consistency (the gate's own bar, checked against that
disposition).** A delegated sweep of C1, C2, C3, Load-Bearing Decisions,
Consequences, Failure Modes and Decision Rationale confirms that no prose
anywhere in the record relies on any of the five as a settled fact. Every
site stating one of those properties states it normatively — what the
implementation MUST do — and carries its forward pin. The strongest
candidate reading, Approach's "everything after normalization … sees only
literal writes and is unchanged", is explicitly hedged in Consequences: "a
`step` ladder and a hand-unrolled ladder are two authorings of one
behaviour but not of one dump … the MVV compares plans and findings, not
dumps", and that comparison "is a TEST oracle, not a product surface".
A11's underlying consumer census is source fact already recorded in its own
Evidence; no site claims the built implementation does it. So "unrunnable
before Phase 2" is consistent with the record's prose, and nothing
downstream depends on any of the five as already true.

**The one `Refuted` record.** A8 is terminal and correctly so: ruling Q1
took disposition (b), and the refutation is absorbed into C3 (narrowed to
claim the dump/graph `identity` precedent only) and into the Decision
Rationale, with the Failure Modes line the ruling required. It is not a
stale `Pending` and does not block.

## 3. Scope Verification

The Minimum Viable Validation is **in scope and will be executed during
implementation**, not deferred. Phase 3 ("Authoring surface and proof")
states it directly: "land the MVV fixture pair and the refusal fixtures,
and assert C3's invariance on the export and the CLI previews."
Prerequisites adds no MVV deferral, and all five Pending assumption Status
notes read "pinned, not deferred" and "MUST pass before the implementation
is accepted".

**The specific proof.** An authored fixture pair — `ladder-literal.toml`
and `ladder-step.toml`, the same ladder written unrolled and stepped —
plus the refusal fixtures `ladder-literal-blocking.toml`,
`ladder-step-blocking.toml` and `ladder-step-unguarded.toml`, over seven
numbered items:

1. author the two fixtures;
2. `intrastate lint` on both — identical finding sets under S1's
   `(code, key, dimension, class, reason)` projection, zero
   `graph-overlap`, zero `graph-coverage-gap`, plus the per-row finding
   naming the CELL (`retry#3`);
3. `flow resolve` over every (state, outcome) cell — identical `writes`,
   `next` and `clear`;
4. `ladder-step-unguarded.toml` refusal, plus the `[initial]` and
   predicate-literal refusal fixtures;
5. the dead-row outcome and A5's negative control, two fixtures
   (S5b, S5c);
6. `graph` export decodes under `intrastate.graph/1` with no new member;
7. source legibility — `ladder-step.toml`'s `[[rules]]` count is one per
   intent and strictly less than the literal ladder's.

End state as the record states it: the step ladder is the literal ladder to
every consumer, it is one row per intent in source, and the unguarded cap
is a named load refusal. Testing Strategy carries the full scenario set
S1-S11 (with S2b, S3b, S3c, S5b, S5c) behind those items.

**Blast radius**: not owed. `clustered=false` — this record declares no
Cluster siblings, so no `impact.md` is required and the fence did not ask
for one (`fence-clear`, not `stopped:impact-unwritten`). The two Overrides
entries against `0002:C4` and `0002:C13` are the cross-record surface, both
declared on the record and both resolving.

## 5. Proportionality

Right-sized. Nothing flagged for trimming before lock.

**Contract count, not word count.** The record is the sole author of
exactly one independent load-bearing contract. `contracts=3`, but
`contracts_durable=1` and `contracts_surface=2`: C2 and C3 are both
`surface-of` C1 (the projector resolves both edges), so the record owns one
seam — the step write and its expansion into literal rows — with C2 (the
bound refusal) and C3 (payload invariance) as surfaces of it rather than
independent seams. No split is warranted.

**Profile re-validated, not re-affirmed.** `rdr-write --outcome profile`
with this clause's own dispositions (`user_facing=yes`, `locks=cross-rdr`)
emits `foundational` under rule `profile-cross-rdr` — a match, not a stop.
The accretion floor resolves `emit.floor: none` (`floor-below-two`: fewer
than two prior point-fixes at the locus), so no raise applies and the
recorded tier is the emitted one. The field's form is already correct:
value, one clause naming the contract, and both dispositions, with no
matrix or provenance prose left from the template. The lenses that ran
agree with the tier rather than contradicting it — the full battery is
complete (3amigo, critique with differing models, repeatability at variant
`full` with three runs and a diff, cove), which is what `foundational`
demands; only `grounding` is absent, and it is absent with no findings
owed. Nothing routed past the lens battery on a wrong Profile.

**Length.** The record runs long, and the length is where the work is: 16
Critical Assumptions with structured Evidence Records, a scored decision
matrix over four approaches, an eleven-row Existing Infrastructure Audit
that turned up the extract-and-move of `valueSatisfies` (four copies
avoided, the first of them in the loader where a divergence is a load
refusal rather than a lint finding), and eleven test scenarios. The two
places bulk was genuinely disproportionate were the A12 and A15 Evidence
fields, and those were relocated to `artifacts/evidence-a12-a15.md` at the
mechanical sweep rather than carried. What remains is proportionate to a
`foundational` record that overrides two clauses of a predecessor and
narrows a shipped payload's published meaning.
