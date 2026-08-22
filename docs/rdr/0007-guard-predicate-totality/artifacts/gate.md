Model: claude-opus-5[1m]

# Finalization Gate — RDR 0007

- **RDR**: 0007-guard-predicate-totality — Guard predicate totality over an
  incomplete evaluation view
- **Date**: 2026-08-21
- **Profile**: foundational
- **Lock**: re-lock. First locked 2026-08-11; demoted 2026-08-12 by JDR 0001
  (`Draft [revised from Final — re-verify A3, A16, A17, A18, A19, A21, A22]`),
  then re-proposed, re-refined, re-resolved, re-lensed (all four
  `foundational` lenses, `*/iter-2/`) and reconciled three times.
- **Mechanical pre-sweep**: PASS
  (`evidence/tooling-pass/iter-2/tooling-pass.md`; five mechanical findings
  fixed in-pass, sweep re-run clean, loop-breaker checked against run 1 —
  no finding re-reported)
- **Verdict**: **READY — Gate PASS**

## 1. Contradiction Check

No contradiction survives between Research Findings and the Proposed Solution.
Four tensions were examined rather than assumed away; each is conceded or
resolved in the text, not papered over.

**Research direction vs. the rule.** The decisive external citation (SCXML
§5.9.1) folds an unevaluable condition to `false` — the *opposite* of this
RDR's rule. *Investigation* concedes this in the open: SCXML mandates the fold
*only* alongside an observable `error.execution` on a second channel, and
"SCXML continues after the error where this RDR halts; the citation supports
'never silently false,' nothing more." The narrow claim drawn is the one the
source actually supports; the rest rests on the in-repo house rule (A8, RDR
0004's "Indeterminate MUST be a refusal-class result") and the accepted-cost
argument. Concession, not contradiction.

**Enforcement site vs. the seam's job.** SQL:2003's `RETURNS NULL ON NULL
INPUT` — "the function itself is not invoked" — is used as instance precedent
for the *host* enforcing a partial domain rather than trusting the routine.
That maps onto the kernel deciding absent-key atoms and never calling the
evaluator for them, which is exactly what the `authority` mini-check table
records ("Value-atom verdict over an ABSENT key … seam is NOT consulted"). The
research and the design agree on the enforcement site.

**Solution vs. shipped behavior.** Where a spike contradicted the draft, the
draft moved. The Stage 4 aggregation probe refuted the escape half of the
original aggregation clause, and the text was corrected to shipped kernel
behavior rather than proposing a kernel change. The Stage 6 re-spike went
further and corrected the RDR against itself twice: Fixup-1d is a **re-decide**
for a deeper reason than the RDR had stated (the fixture's guard names a key
the view does not carry, so the nil seam is never consulted and the test would
go green while testing nothing), and A26's banned lists are **exact**-match,
not substring — so the RDR's claim that REQ-37 forecloses A22's fallback "under
any name containing `Canonicalize`" was overstated and is corrected at A26.
Both corrections run against the author's convenience.

**Stated principle vs. planned feature.** The principle is "never collapse an
undecided third value," and the Normative Contracts apply it reflexively: a
mapping failure the evaluator cannot parse MUST NOT be reported as
`GuardUnevaluable`, because that re-creates the same conflation one layer up.
Likewise `guard_unevaluable` is excluded from escape lists by RDR 0002's
closure, so the design denies itself the escape hatch it removes from authors.

One residual tension is **recorded rather than resolved**, correctly: D8's
original ordering rationale ("absent owned state is the more precise
diagnosis") is invalidated by this RDR's own `RequiresOwned` narrowing, yet the
precedence is retained because it is frozen kernel behavior. The RDR states the
honest cost — a row failing both ways surfaces the write-dependency problem
first — instead of asserting a superseded rationale still holds.

## 2. Assumption Verification

26 Evidence Records. **22 Verified, 4 DOWNGRADED (A12, A19, A22, A25). None
Pending, none Unverified** — confirmed by grep, not by claim. Every record is
internally consistent: Status, Method and Evidence agree, and every "If wrong"
is non-empty and names what fails and how it surfaces.

- **No `Docs Only` records** — zero occurrences; nothing blocks on that axis.
- **No `Source Search` self-reference** (C3 clean). Every citation into this
  RDR's own artifact dir is made by a Spike or Prior Art record; A27's artifact
  citation points at *peer* RDR 0003's fixture, which is external evidence.
- **Every anchor resolves** (C5). ~25 kernel symbols, 13 test anchors, all
  JDR 0001 sections cited (§D1–§D4, §JD-1/2/3/4/6/7/8/9/12, P5) and every
  peer-RDR quotation. No phantom, renamed or moved symbol. Two load-bearing
  **demonstrated negatives** were re-confirmed on source, not inferred:
  `canonicaliz` appears zero times in RDR 0002 (A22) and `RequiresOwned` zero
  times in RDR 0009 (A21).
- **Status consistency holds** (C6). Two reserved-word collisions (A9, A27 used
  "Pending" for a half their Status line did not mark Pending) were fixed
  in-pass as wording; no disposition changed.

**The four DOWNGRADED assumptions do not block lock**, and the reason is the
same in each case: none is blocked on evidence *this RDR can produce*, each
carries a named plan and a **fail-closed** fallback, and the kernel contract is
sound under either branch.

- **A19** (structured `omitempty` field on `CLIError`) — blocked on JDR 0001
  §D4 being reopened. Reversible in both directions; the `Detail`-flattening
  fallback is written. Blocks RDR 0005's JSON rendering only — "the kernel side
  is unblocked either way: the payload exists on `Refusal` regardless of how
  the CLI renders it."
- **A22** (tag-key canonicalization) — blocked on RDR 0002 re-lock. Kernel-side
  fallback (exact byte equality as an input precondition) verified today; a
  foreign spelling becomes a value atom and fails closed.
- **A25** (0009 fixture collision) — the collision itself is **Verified on
  source**; only the *sequencing* is open, which is a cluster question this RDR
  explicitly does not settle.
- **A12** (two-row absence pattern) — home is JDR 0001 §JD-4, substance
  decided, rides RDR 0003's lock. "Nothing in the kernel's domain rule depends
  on the answer."

**The re-verify list is discharged.** The demotion named seven: A3, A16, A17,
A18, A19, A21, A22. A16/A17/A18/A21 verified; A19/A22 downgraded as above; and
**A3 is Verified by a re-spike run against the *specified* shape** — with
`Refusal.Guard` actually deleted and the seam view-free — 154/154 frozen tests,
180/180 with probes, `go vet` clean, `internal/` never edited. That same
re-spike closed A23, A24 and A26. This matters more than a count: the first
spike proved a *superset* shape, and the re-spike is what makes A3's bound a
demonstrated fact rather than an inference.

One process defect is recorded rather than buried, and is the honest reason
this gate reads the absorption audit carefully: **B-10 was dropped at the
cross-model merge step** — a first-class finding from the alt-model critique
pass that appeared in no queue table and never reached dispositions. Stage 6
caught it, dispositioned it ACCEPTED (Design Decision), and the rejected
alternative it named is now written into the GATE-THEN-COUNT clause. A merge
that silently loses a finding is the one defect the dual-model draw exists to
prevent; it is logged in the reconcile report as a mechanism lesson.

## 3. Scope Verification

**The MVV is in scope, kernel-level, and executed during implementation — not
deferred.** It runs against the real kernel with a real atom, no stub, entirely
inside the surface the re-spike proved.

The specific proof is four named scenarios (Testing Strategy rows 1–4):

1. *masking probe inverts* — one candidate row whose guard is a value atom over
   an **absent** key, whose `RequiresOwned` is **satisfied**, plus a modeled
   `no_match` escape ⇒ `Refusal.Kind == guard_unevaluable`, nil `Plan`, the
   absent key in the payload, the row in `Refusal.Rows`.
2. *D8 preserved* — same table, key present and value decided FALSE ⇒ row
   prunes, escape plans. Guards against over-refusing.
3. *unevaluable-blocks-true-sibling* — row A unevaluable beside row B decided
   TRUE ⇒ refusal naming A, no plan.
4. *sanctioned route from absence* — `exists = false` over the absent key
   decides TRUE and the row is selected.

The satisfied `RequiresOwned` in scenario 1 is load-bearing, not incidental,
and the RDR says why: an unevaluable row is a **survivor**, so the owned sweep
runs over it and reports first (GATE, THEN COUNT). A row carrying both an
absent owned key and an unevaluable atom yields `owned_state_unavailable` — so
putting that pairing in the masking probe "would test the owned sweep while
claiming to test the domain rule." That is the RDR policing its own MVV.

**The MVV is writable now.** This is the question iteration 2 of the reconcile
left explicitly open — the MVV asserts on the payload, and the payload surface
(A23) was then unproven, so the MVV "cannot be written without the field it
asserts on." A23 is now Verified *by a spike that built the payload and
asserted on it*. Resolved by evidence, not adjudicated away. The four
DOWNGRADED items were then checked against the MVV text directly: it references
no CLI envelope (A19), no canonicalization (A22), no fixture sequencing (A25)
and no absence pattern (A12).

This gate also credits the RDR for **not overclaiming**: its own "What the MVV
does and does not prove" paragraph states that scenarios 1–2 are regression
pins on already-frozen mapping and that a green MVV MUST NOT be read as
validating the domain rule. The domain rule's honest gate is the Phase 2
matrix.

## 4. Cross-Cutting Concerns

Only concerns that genuinely apply:

- **Incremental adoption / migration.** The rule lands *before* the first
  evaluator exists, which is why "this RDR Final before RDR 0003's
  implementation begins" is a Prerequisite. The RDR corrects its own earlier
  over-claim: "no installed base" counts implementations, not authored intent
  — the 0002 and 0003 reference fixtures already author guards that would
  become `guard_unevaluable`, and Phase 2 MUST classify each as
  safe-or-migration. Owned here.
- **Determinism.** Owned here, and materially expanded since the first lock.
  The payload MUST be sorted on a **total** key `(RuleID, SourceLocator, key,
  block, operator, literal)` so it is a function of the input, with dedup
  stated as a separate kernel obligation ("the sort orders entries, it does not
  dedupe"), and `MissingOwned` is the deduplicated sorted union across
  survivors. Pinned by Testing rows 17, 20 and 23, and independently guarded by
  the four frozen row-order tests A23 verified. Order-stability of the refusal
  as a whole remains RDR 0001's (ADV-3/ADV-3b).
  This RDR claims **no** byte-identical output, content-addressed identity, or
  replay-stable hash — REQ-37 affirmatively pins that the kernel introduces no
  hash or canonical serialization — so the template's hash/pre-image checklist
  does not apply.
- **Concurrency model** — not applicable in the usual sense, but load-bearing
  once: A4's PostgreSQL citation notes commutativity, which rules out a
  short-circuit reading and makes strong-Kleene dominance *semantic* rather
  than evaluation-order dependent. That matters for an evaluator free to visit
  atoms in any order. Addressed by the normative truth tables.
- **Versioning.** Phase 2's vectors are a named, versioned golden vector suite,
  and the harness must be exported so RDR 0003's build can instantiate it
  (Existing Infrastructure Audit marks it Build).
- **Character encoding / identity.** Tag-key identity is exact byte equality at
  the kernel with no case folding (A22); canonicalization is the normalizer's
  duty, not yet a 0002 clause, and the kernel fails closed on a mismatch.

Concerns owned by peers, which this RDR conforms to rather than restates: the
escapable-class closure (RDR 0002), the refusal taxonomy and selection rule
(RDR 0001), the operator vocabulary and lint rejection of unknown operators and
undeclared tags (RDR 0003 + 0002), the no-collapse principle (RDR 0004 A8,
carried by RDR 0005).

Deliberately out of scope with a named home: `recognized`-key binding → RDR
0008; escape-row shape → RDR 0009; read-completeness → RDR 0004 (A6b, closed
here by JDR 0001 §D3/§JD-7); CLI rendering of the payload → RDR 0005, blocked
on A19.

## 5. Proportionality

**Contract count, not word count.** Eleven `normative` blocks resolve to **one
sole-authored contract**:

- Blocks on value-operator partiality, existence totality, atom scope / empty
  block identities, and strong-Kleene combination are one contract decomposed
  by *scope level* — a single answer to "what does the evaluator answer over an
  incomplete view." No implementer can hold one without the others.
- Provenance-blind presence is **forced, not chosen**: `TagSet.Lookup` is the
  only exported accessor an external evaluator can use (A11, re-verified on
  source). Recording a forced consequence is not authoring a second contract.
- The SEAM block explicitly **defers** the atom grammar to JDR 0001 §D1 and the
  operator semantics to RDR 0003 — it cites, it does not restate.
- GATE-THEN-COUNT and the survivor-membership block **pin shipped kernel
  behavior** (Probe A forced the text to match the kernel, not the reverse),
  and the D8 block ratifies RDR 0001's deviation. Not sole-authored.
- The `RequiresOwned` narrowing is the one candidate a critique named as
  plausibly severable. **It is not severable**: it is the originating kata
  (`xg7p`), it is the *mechanism* by which the domain rule resolves the
  conflation (decidability moves to the guard's referenced-tag set *instead of*
  `RequiresOwned`), and it is what invalidates D8's ordering rationale — a
  consequence this same RDR must then re-pin. Splitting it would leave one RDR
  redefining a field and another silently depending on that redefinition to
  close its masking path.

**Verdict: no split.** One contract, stated at four scope levels, plus forced
consequences and pins on frozen behavior.

**Profile re-validated: `foundational` is correct and stays.** It is carried
not by contract count but by blast radius — a cross-RDR producer (0003's
evaluator implements it, 0001's kernel consumes it, 0002 supplies the escape
closure), spanning modules. Critically, **the lenses that actually ran match
the profile's required row against the current draft**: `cove → 3amigo →
critique → repeatability`, all four present under `*/iter-2/`, critique a
genuine dual-model draw (`claude-opus-5` + `claude-fable-5`, diffed by
`claude-sonnet-5`), repeatability the full variant across three distinct base
models. The latch did not route past the battery — and the Stage 6 preflight
that caught the iteration-2 gap (two lenses had never reviewed the current
design, returning NOT RECONCILED) is the evidence that this check has teeth.
`Seam Lineage` records no prior accretion, so the accretion floor is not
engaged. Form conforms: value plus one clause naming the contract, no
matrix/provenance prose.

**Right-sizing.** ~2,400 lines is long, and length was tested rather than
excused. The mass is concentrated in 26 Evidence Records, four of which exceed
the 30-line advisory budget (A3 44, A21 38, A19 34, A26 34 — `4 fields over
budget, 150 lines`). Each was inspected at the question CHECK 9 actually asks:
is the load-bearing anchor still findable? Yes in all four — each opens with
its verdict and stable anchor, and the bulk below is verification content the
grounding sweep reads (re-spike counts, frozen-test enumerations, the
`CLIError` field census, three banned-list tests). Several records are
**demonstrated negatives** ("specified nowhere" in a named peer), which cannot
be short: they must show the sweep that found nothing. Nothing reads as
change-history narration. **No section is flagged for trimming before lock.**
