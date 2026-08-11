Model: claude-opus-5[1m]

# Finalization Gate — RDR 0007

- **RDR**: 0007-guard-predicate-totality — Guard predicate totality over an
  incomplete evaluation view
- **Date**: 2026-08-11
- **Profile**: foundational
- **Mechanical pre-sweep**: PASS
  (`evidence/tooling-pass/tooling-pass.md`; four mechanical findings fixed
  in-pass, sweep re-run clean)
- **Verdict**: **READY — Gate PASS**

## 1. Contradiction Check

No contradictions remain between Research Findings and the Proposed Solution.
Three tensions were examined rather than assumed away:

**Research vs. solution direction.** The decisive external citation (SCXML
§5.9.1) folds an unevaluable condition to `false` — the *opposite* of this
RDR's rule. The RDR does not paper over this: *Investigation* carries an
explicit "What this citation does and does not authorize" paragraph conceding
that SCXML's disposition is to continue, that the citation does not underwrite
the resolution-level veto, and that it cuts *against* omitting the absent tag
from the refusal. The narrow claim the RDR draws — an undecidable condition
MUST NOT silently become false, and the error must be observable — is what
SCXML actually supports, and the RDR rests the rest on the in-repo house rule
(A8) and the accepted-cost argument. Concession, not contradiction.

**Solution vs. shipped behavior.** The Stage 4 aggregation spike *refuted* the
escape half of the original aggregation clause (Probe A: the shipped kernel
returns the escape set's blocking refusal, frozen by three tests and D5). This
was resolved in the honest direction — the RDR text was corrected to shipped
behavior and Phase 1 stayed doc-comments-only — rather than by proposing a
kernel change. A3's Status records the correction inline.

**Stated principle vs. planned feature.** The RDR's own principle is "never
collapse an undecided third value." The Normative Contracts apply it
reflexively: a mapping failure the evaluator cannot parse MUST NOT be reported
as `GuardUnevaluable`, because that would re-create the same conflation one
layer up. The design rules against its own convenience.

One residual tension is *recorded* rather than resolved, correctly: D8's
original ordering rationale ("absent owned state is the more precise
diagnosis") is invalidated by this RDR's own `RequiresOwned` narrowing, yet the
precedence is retained because it is frozen kernel behavior. The RDR states the
honest cost — a row failing both ways surfaces the write-dependency problem
first — instead of asserting the superseded rationale still holds.

## 2. Assumption Verification

16 Evidence Records. **15 Verified, 1 Pending (A6b) by explicit decision.**
Every record is internally consistent: Status, Method, and Evidence agree, and
every "If wrong" is non-empty and specific (each names what fails and how it
surfaces).

- **No `Docs Only` records** — nothing blocks on that axis.
- **No self-referential `Verified` stamp** — CHECK 3 clean; every Source Search
  resolves to peer-RDR text or to `internal/resolve/` source, never to this RDR.
- **Every cited `path::Symbol` resolves on the branch** — CHECK 5 verified 40+
  anchors in their cited files, including the structural facts the normative
  blocks assert (`GuardEvaluator` is an interface with the exact cited
  signature; `RowRef` is `(RuleID, SourceLocator)`; `fixtureGuards` discards the
  `TagSet`). 39 peer-RDR quotations verified; two imprecisions found and fixed
  in-pass.
- **Status consistency holds** — A6b is the only non-Verified record, and all
  nine of its references treat it as open. No settled-fact prose leans on it.

**A6b (Pending) — plan, and why it does not block lock.** The question (is a
*failed* accessor read distinguishable from *genuine* absence at the
accessor→kernel boundary?) is not resolvable inside this RDR: `Input.Owned`
carries no error channel, so a truncated snapshot and a real absence are
byte-identical at the kernel seam. It is not MVV-critical, and the contract is
sound without it — the exposure is that an unevaluable refusal cannot be read
as "retryable", which Failure Modes states as operator guidance. The obligation
is **routed to a destination that exists**: RDR 0004's implement stage, which
owns the accessor executor and must state whether the "typed tag values" branch
carries a completeness guarantee.

**Four assumptions verified as named obligations rather than settled facts**
(A10, A12, A15 → RDR 0003's `Phase 1: Predicate Model`, confirmed to exist and
to already charter "the normalized predicate atom shape used by resolver and
lint"; A12 also → RDR 0006's implement stage). This is the correct disposition
under a never-amend rule: each is *silence in a Final peer*, not refutation, and
each lands inside an existing implement-stage scope with no peer amendment. A13
is verified as a knowingly **accepted exposure**, bounded today because nothing
imports `internal/resolve` and `root.go` registers only `version` — both
independently confirmed.

The honest consequence is stated in the RDR rather than hidden: Testing Strategy
rows 4–8 and the Phase 2 vectors are **un-encodable until RDR 0003 discharges
A10**. That is a declared downstream dependency, not a defect here.

## 3. Scope Verification

**The MVV is in scope and will be executed during implementation, not
deferred.** It is kernel-level and reachable against the tree as it stands.

Specific proof: three named scenarios in Testing Strategy rows 1–3, run against
`internal/resolve/fixtures_test.go::fixtureGuards` —

1. *unevaluable-not-no_match* — guard over an absent tag + absent
   `RequiresOwned` key + modeled `no_match` escape ⇒ `Refusal.Kind ==
   guard_unevaluable`, nil `Plan`. Inverts the masking probe recorded under
   *Background*.
2. *D8 preserved* — same table, tag present and predicate decided FALSE ⇒ row
   prunes, escape rescues. Guards against over-refusing.
3. *unevaluable-blocks-true-sibling* — row A unevaluable beside row B decided
   TRUE ⇒ refusal, no plan, row B not selected.

Scenario 3 is the load-bearing one: **no shipped test covers it**, and the
Stage 4 spike (Probe B) already observed the behavior against the shipped
kernel, so its expected value is captured, not guessed.

This gate credits the RDR for *not* overclaiming here. Its own "What the MVV
does and does not prove" paragraph states that scenarios 1–2 are regression
pins on already-frozen mapping and that a green MVV MUST NOT be read as
validating the domain rule. The domain rule's honest gate is the Phase 2
harness. A residual contradiction on this exact point — the MVV also claiming
mixed-verdict atom combination, which `fixtureGuards` structurally cannot
express — was found by the sweep and **fixed in-pass**, along with a Risks
mitigation that rested on it.

## 4. Cross-Cutting Concerns

Only concerns that genuinely apply:

- **Incremental adoption / migration.** The rule lands *before* the first
  evaluator exists, which is why sequencing (0007 Final before 0003 implement)
  is a Prerequisite. The RDR corrects its own earlier over-claim: "no installed
  base" counts implementations, not authored intent — the 0002 and 0003 spike
  fixtures already author guards that would become `guard_unevaluable`, and
  Phase 2 MUST classify each as safe-or-migration. Owned here.
- **Concurrency model — not applicable in the usual sense, but load-bearing
  once.** A4's PostgreSQL citation notes commutativity, which rules out a
  short-circuit reading and makes strong-Kleene dominance *semantic* rather
  than evaluation-order dependent. That matters for an evaluator free to visit
  atoms in any order. Addressed by the normative truth tables.
- **Determinism.** Refusal payload order-stability is RDR 0001's (frozen by
  ADV-3/ADV-3b). This RDR adds one determinism obligation of its own —
  `Refusal.Guard` is single-valued and selected as the lowest row by
  `(RuleID, SourceLocator)`, so any test discriminating *which* row went
  unevaluable MUST assert on `Rows`. Owned here, and Scenario 5 pins it.
  This RDR claims no byte-identical output, content-addressed identity, or
  replay-stable hash, so the hash/pre-image checklist does not apply.
- **Versioning.** Phase 2's vectors are a "named, versioned golden vector
  suite" — versioning is in the contract, and the harness must be exported so
  RDR 0003's build can instantiate it.

Concerns owned by peers, which this RDR conforms to rather than restates:
the escapable-class closure (RDR 0002), the refusal taxonomy and selection rule
(RDR 0001), the operator vocabulary and lint/parse rejection of unknown
operators and undeclared tags (RDR 0003 + RDR 0002), the no-collapse principle
for an undecided third value (RDR 0004 A8, carried by RDR 0005).

Deliberately out of scope, with a named home: `recognized`-key binding →
RDR 0008 (cove F-8, charted); escape-row shape → RDR 0009; read-completeness →
RDR 0004 implement (A6b); `--tag` channel constraint → seeded against RDR 0005.

## 5. Proportionality

Stage 6 explicitly charted this item here rather than deciding it, because both
critique passes (C-14, B-15) independently counted **five or more** independently
load-bearing contracts against a `Profile` reading "foundational — one
contract". Counted in front of me, with the split test being contract count, not
word count:

Ten `normative` blocks resolve to **one sole-authored contract**, not five:

- Blocks 2, 3, 5, 6 (value-operator partiality, `exists` totality, atom scope /
  empty-block identities, strong-Kleene combination) are one contract
  decomposed by *scope level*. They are a single answer to "what does the
  evaluator answer over an incomplete view"; no implementer can hold one without
  the others, and each is meaningless alone.
- Block 4 (provenance-blind presence) is **forced, not chosen** —
  `TagSet.Lookup` is the only exported accessor an external evaluator can use
  (A11), and RDR 0003's own fixture guards an `observed` tag. Recording a forced
  consequence is not authoring a second contract.
- Blocks 1, 7, 9, 10 are **not sole-authored**: block 1 explicitly defers the
  mapping to RDR 0003; blocks 7 and 9 pin *shipped* kernel behavior (A3 verifies
  no change, and Probe A forced the text to match the kernel rather than the
  reverse); block 10 ratifies RDR 0001's deviation D8.
- Block 8 (`RequiresOwned` narrowing) is the one candidate the critique named as
  most plausibly severable. **It is not severable.** It is the originating kata
  (`xg7p`); it is the *mechanism* by which the domain rule resolves the
  conflation — guard decidability moves to the guard's own referenced-tag set
  *instead of* `RequiresOwned`; and the narrowing is what invalidates D8's
  ordering rationale, a consequence this same RDR must then re-pin. Splitting it
  out would leave one RDR redefining a field and another silently depending on
  that redefinition to close its masking path — a worse seam than the one it
  would relieve.

**Verdict: no split.** One contract, stated at four scope levels, plus forced
consequences and pins on frozen behavior.

**Profile re-validated: `foundational` is correct and stays.** It is not
carried by the contract axis but by blast radius — this RDR is a cross-RDR
producer (RDR 0003's evaluator implements it, RDR 0001's kernel consumes it,
RDR 0002 supplies the escape closure) and it spans modules. The lenses that
actually ran match the profile's required set (`cove 3amigo critique
repeatability`, all four present, critique dual-model, repeatability full
variant across three distinct base models), so the latch did not route past the
battery. `Seam Lineage` records no prior accretion, so the accretion floor is
not engaged; the value stands on the contract axis's cross-RDR trigger alone.

**Form check**: the field is value + one clause naming the contract, with no
matrix/provenance prose left from the template. Conforms.

**Right-sizing.** At ~2,460 lines this is a long document, and length was
tested rather than excused. The volume is concentrated in 16 Evidence Records
(~700 lines) that are load-bearing — six of them are the demonstrated negatives
("specified nowhere") on which the routed obligations rest, and a demonstrated
negative cannot be short: it must show the sweeps that found nothing. Nothing
here reads as change-history narration or bloat, and the Refine stage already
ran. No section is flagged for trimming before lock.
