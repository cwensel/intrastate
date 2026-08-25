# REQ List — RDR 0007 Guard predicate totality over an incomplete evaluation view

Phase 0 audit artifact. Every clause below is a testable obligation drawn from
`docs/rdr/0007-guard-predicate-totality.md`. Quotes are verbatim, copied from
the projector (`rdr inspect --select <id>`), never transcribed.

Element ids are carried where the REQ derives from a labelled contract
(`0007:C1` … `0007:C11`, `0007:MVV`). Source sections are abbreviated:

- `NC` = Proposed Solution / Technical Design / Normative Contracts
- `LBD` = Proposed Solution / Technical Design / Load-Bearing Decisions
- `AP` = Proposed Solution / Approach
- `TD` = Proposed Solution / Technical Design (prose)
- `FM` = Trade-offs / Failure Modes
- `RM` = Trade-offs / Risks and Mitigations
- `TS` = Validation / Testing Strategy
- `MC` = Validation / Mini-checks
- `MVV` = Implementation Plan / Minimum Viable Validation
- `IP` = Implementation Plan (phases)
- `PE` = Validation / Performance Expectations

**Two standing corrections apply to this REQ set and are marked at each REQ
they touch** (RDRs are never amended; the correction lives here and in
`deviations.md` D1):

1. **§D12 (§JD-21) widens `0007:C1`'s fenced cardinality.** The record fences
   `Block` as "exactly two constants" (`0007:1264-1266`). JDR 0001 §D12 lands
   `BlockMatch` as a third constant on the *same* type, **in this RDR's Phase 1**
   — one type, no separate slice on `Row`. REQ-6 states the widened form.
   Recorded as a scoped normative re-entry: `deviations.md` D1.
2. **§D13 (§JD-22) unblocks the `contains` leg.** The record says the `contains`
   contract-test leg and TS scenario 8's present-key half are "Blocked on one
   RDR 0003 declaration" (`0007:2159-2170`) and names 0003 as the declarer
   (`0007:2182`). JDR 0001 §D13 reassigns the declaration to **RDR 0002**:
   *"`Tag.Value` stays `string`; a set crosses as its canonical JSON array —
   members sorted, duplicate-free, compact encoding — and 0002 declares it."*
   Per `docs/rdr/BUILD-ORDER.md` ("The one seam this order strains"), 0007 runs
   before 0002 writes that clause and **MUST write these tests against §D13
   directly**, not defer them. REQ-56 and REQ-57 state them. Recorded at
   `deviations.md` D1.

---

## Seam shape and atom transport (`0007:C1`)

- [REQ-1] "SEAM. A candidate row carries its guard as a slice of parsed atoms — key, operator token, literal, block ∈ {all, unless} — the shape JDR 0001 §D1 fixes; this RDR cites it and does not restate the grammar." — (NC, `0007:C1`)

- [REQ-2] "An empty slice is an unguarded row." — (NC, `0007:C1`)

- [REQ-3] "There is no opaque guard string and no reconstruction step, so no mapping failure exists and no panic or error channel is needed for one." — (NC, `0007:C1`)

- [REQ-4] "`Refusal.Guard` has no referent once guards are atoms" — the shipped `Row.Guard string` and `Refusal.Guard string` fields are DELETED, not re-typed. — (NC, `0007:C8`; TD "`Row.Guard` becomes an atom slice … `Refusal.Guard` is replaced by the per-row, per-atom payload")

- [REQ-5] "The atom's four fields are named `Key`, `Operator`, `Literal`, `Block` — the payload's `UndecidedAtom` mirrors them by name, so divergent spellings would make the mirror a mapping." — (NC, `0007:C1`)

- [REQ-6] "`Block` is an exported named STRING type carrying exactly two constants, `BlockAll = \"all\"` and `BlockUnless = \"unless\"`." — (NC, `0007:C1`) **§D12 WIDENS THIS FENCE.** The type gains a third constant `BlockMatch` in this RDR's Phase 1 — one type, no separate slice on `Row`. The testable obligation is: `Block` is an exported named string type; `BlockAll = "all"`, `BlockUnless = "unless"`, and `BlockMatch` are constants of it; the payload's `Block` is that same exported constant type (REQ-50). See `deviations.md` D1 (§D12/§JD-21) — this is a normative change to a Final RDR, taken as a scoped re-entry, not a silent widening.

- [REQ-7] "The kernel evaluates a guard ATOM BY ATOM. For each atom it MUST first decide presence of the referenced key in the assembled view, provenance-blind (the `TagSet.Lookup` `ok` test)." — (NC, `0007:C1`)

- [REQ-8] "an existence atom is decided by the KERNEL from presence alone, and the evaluator MUST NOT be consulted for it" — (NC, `0007:C1`)

- [REQ-9] "a value-comparing atom whose key is ABSENT is marked unevaluable by the KERNEL, and the evaluator MUST NOT be consulted for it" — (NC, `0007:C1`)

- [REQ-10] "a value-comparing atom whose key is PRESENT is handed to the evaluator seam with the atom and the present value: `Evaluate(atom, value) GuardResult` — a single-method INTERFACE named `GuardEvaluator`, not a func type, so the nil-seam rule above is a nil interface value and RDR 0003's evaluator satisfies it by declaring the method." — (NC, `0007:C1`)

- [REQ-11] "The seam decides true or false under RDR 0003's typed operator semantics, and MAY answer unevaluable for a present value it cannot compare (A18). It never sees the view." — (NC, `0007:C1`)

- [REQ-12] "The kernel therefore enforces the domain rule by structure: an evaluator cannot fold absence into false, because it is never asked about an absent key." — (NC, `0007:C1`)

## Nil-seam rule is per-atom (`0007:C1`)

- [REQ-13] "A NIL seam does not make a row unevaluable by itself. The kernel decides existence atoms and absent-key atoms without consulting the seam, so a row whose atoms are all kernel-decidable resolves under a nil seam exactly as it would under a present one." — (NC, `0007:C1`)

- [REQ-14] "A nil seam yields unevaluable only for an atom that WOULD have been handed to it — a value-comparing atom over a present key." — (NC, `0007:C1`)

- [REQ-15] "Such an atom is reported in the payload like any other unevaluable atom, with reason `uncomparable` (its key is present, so `absent` would be a lie); the payload never omits it and never names the nil seam as a reason." — (NC, `0007:C1`)

- [REQ-16] "This narrows the shipped whole-guard behavior (`evaluateGuard` returns GuardUnevaluable for any non-empty guard when `seam == nil`) to a per-atom rule." — (NC, `0007:C1`) Discriminating test: TS row 19 — a pure-`exists` row under a nil seam now yields a PLAN and its `Writes` land; assert both that leg AND that a value atom over a present key still refuses under the same nil seam.

- [REQ-17] "Frozen `TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue` keeps its verdict but changes its REASON, and Phase 1 must re-read it rather than re-encode it … To keep testing the nil-seam rule the fixture needs a value atom over a PRESENT key (e.g. `reviews >= 3`)." — (NC, `0007:C1`; TS preamble marks it RE-DECIDED, not re-encoded)

## Value-operator partiality (`0007:C2`)

- [REQ-18] "Value-comparing guard operators (equality, membership, bounded integer comparison, set containment) are PARTIAL over the assembled evaluation view: an atom whose referenced tag key is absent from the view MUST evaluate to unevaluable — never to false and never to true." — (NC, `0007:C2`)

- [REQ-19] "An absent set-valued tag MUST be treated as unevaluable under set containment, NOT as the empty set." — (NC, `0007:C2`; TS row 8)

## Existence totality (`0007:C3`)

- [REQ-20] "The existence operator is the sole TOTAL operator: it MUST decide true or false from presence or absence of its referenced key alone. Its verdict is `presence == literal`: `exists = true` decides TRUE when the key is present, `exists = false` decides TRUE when it is absent." — (NC, `0007:C3`; TS row 7 requires BOTH polarity channels)

- [REQ-21] "Polarity lives in the literal; `all`/`unless` placement composes with it. Absence tests MUST be expressed as existence atoms; no value-comparing operator may act as an implicit existence test." — (NC, `0007:C3`)

- [REQ-22] "The constants are `OpExists` (the operator token) and `LiteralTrue` / `LiteralFalse` (the two boolean literal forms); the kernel compares an atom's operator and literal against these values verbatim, performing no parsing, case-folding, or coercion of its own." — (NC, `0007:C3`)

- [REQ-23] "`OpExists = \"exists\"`, `LiteralTrue = \"true\"`, `LiteralFalse = \"false\"`, all lower-case." — (NC, `0007:C3`) All three exported as kernel constants (A16).

- [REQ-24] "an existence atom carrying a foreign TOKEN is a value atom to the kernel (unevaluable on absence, handed to the seam on presence)." — (NC, `0007:C3`; TS row 11)

- [REQ-25] "A foreign LITERAL on an `OpExists` atom — any value that is neither `LiteralTrue` nor `LiteralFalse`, the empty literal included — is UNEVALUABLE at the kernel, reason `uncomparable`; the kernel MUST NOT decide such an atom from presence, and MUST NOT treat a missing literal as `LiteralFalse`." — (NC, `0007:C3`; TS row 11)

- [REQ-26] "The normalizer additionally MUST reject a foreign literal at load (RDR 0002's typed validation), so the kernel rule is the fail-closed backstop for a producer that did not, not a duplicate of it" — **NOT a kernel test.** TS row 11: "the load-rejection leg is RDR 0002's and is NOT a kernel test (no load path exists under `internal/resolve`)". Recorded as an out-of-scope boundary, not an obligation on this build.

## Presence is provenance-blind (`0007:C4`)

- [REQ-27] "PRESENCE IS PROVENANCE-BLIND. \"Present in the assembled view\" means the key is in the view under ANY provenance — owned, observed, or recognized. An atom MUST NOT be decided differently according to how a tag reached the view." — (NC, `0007:C4`; TS row 18 — `TagSet.Lookup` `ok`, not `TagSet.has`)

- [REQ-28] "Key identity is exact string equality on the key as assembled; canonicalization of authored key spellings is the normalizer's (RDR 0002), upstream of the kernel, and the kernel performs none (A22)." — (NC, `0007:C4`)

- [REQ-29] "This is deliberately NOT the kernel's owned-state predicate: `missingOwned` tests owned provenance because `owned_state_unavailable` is about the owned snapshot; guard decidability takes the broader test. The two MUST NOT be conflated." — (NC, `0007:C4`)

## Empty atom / empty block identities (`0007:C5`)

- [REQ-30] "A row with no atoms is decided TRUE without consulting the view (shipped: `evaluateGuard`'s empty-guard branch)." — (NC, `0007:C5`; TS row 10)

- [REQ-31] "an empty `all` block is TRUE (empty conjunction)." — (NC, `0007:C5`; TS row 10)

- [REQ-32] "An empty or omitted `unless` block is ABSENT: the `¬(unless_conj)` TERM DROPS OUT of the row verdict, which reduces to `all_result`. It is NOT a vacuously-true conjunction … and it is NOT `unless_conj = F` either" — (NC, `0007:C5`; TS row 10)

## Strong-Kleene combination in the kernel (`0007:C6`)

- [REQ-33] "Atom verdicts combine in the KERNEL under strong-Kleene three-valued logic. Conjunction is `min` under the TRUTH order `F < U < T`; negation is `¬T = F`, `¬F = T`, `¬U = U`" — with the normative table `T∧T=T, T∧F=F, T∧U=U, F∧T=F, F∧F=F, F∧U=F, U∧T=U, U∧F=F, U∧U=U`. — (NC, `0007:C6`; TS row 5 requires the full matrix across `all` and `unless` incl. `¬U = U`)

- [REQ-34] "The kernel MUST implement combination against the tables (or against an explicit truth-order rank), and MUST NOT `min` the raw constant values; the constant order is a shipped fact this RDR does not renumber" — shipped `GuardResult` is `GuardFalse, GuardTrue, GuardUnevaluable` (`iota` 0,1,2), the order `F < T < U`, which is NOT the truth order. — (NC, `0007:C6`; TS row 21, mutation-killing)

- [REQ-35] "The row verdict is `all_result ∧ ¬(unless_conj)`, where `unless_conj` is the conjunction of the `unless` block's atoms (`unless` is block-level negation, NOT per-atom negation)." — (NC, `0007:C6`)

- [REQ-36] "A guard is GuardTrue or GuardFalse only when decided atoms alone determine it; any unresolved dependence on an unevaluable atom yields GuardUnevaluable." — (NC, `0007:C6`)

- [REQ-37] "A present-tag atom decided FALSE still yields GuardFalse beside an unevaluable atom (`F ∧ U = F`): the falsity is witnessed by present state and holds under every resolution of the unevaluable atom — this is what makes D8 sound." — (NC, `0007:C6`; TS rows 3 and 21 — mutation-killing, observable only across TWO atoms)

- [REQ-38] "The tables are normative; the evaluator holds no part of them." — (NC, `0007:C6`; MC `authority` row "Atom combination (K3) … evaluator holds no part of the tables")

## Gate ordering (`0007:C7`)

- [REQ-39] "GATE, THEN COUNT (JDR 0001 §D2, stated as the kernel's ordering). Over the candidate set: guard-FALSE rows prune first (D8); absent owned state among survivors is reported next as `owned_state_unavailable`; then any surviving row whose guard is unevaluable vetoes the resolution as `guard_unevaluable`; only then does RDR 0001's exact-one count run over the rows the gate returned, and only then is escape reachability consulted. `no_match` and `ambiguous_match` sit downstream of all three gate facts." — (NC, `0007:C7`; TS row 13)

- [REQ-40] "The escape set is gated identically, as its own row set (D5) — \"identically\" by DELEGATION, not by a parallel implementation: shipped `escapeOrRefuse` filters the escape rows and then calls the same `gate` function (`gate(escapes, in.Guards, view)`), which is why it needs no change here" — (NC, `0007:C7`)

- [REQ-41] "An unevaluable CANDIDATE is never masked by a decidable escape row (the gate returns before the escape path is reached), and an unevaluable ESCAPE row yields `guard_unevaluable` in place of the candidate-set refusal" — (NC, `0007:C7`; TS row 14, both legs; frozen ADV-2 + Fixup-1d)

- [REQ-42] "The owned-before-unevaluable precedence is pinned to shipped behavior, not to D8's original rationale … Authoring docs MUST NOT repeat the superseded rationale." — (NC, `0007:C7`) Phase 2 obligation: "Rewrite `gate`'s doc comment, which still states D8's superseded rationale verbatim" — (IP, Phase 2)

- [REQ-43] "Combined reporting — one refusal carrying BOTH payloads … — is REJECTED here" — a row failing both ways yields `owned_state_unavailable` only, with the owned payload only. — (NC, `0007:C7`; MC `disposition`; TS row 13)

## Refusal mapping and the per-row/per-atom payload (`0007:C8`)

- [REQ-44] "The kernel maps GuardUnevaluable to `guard_unevaluable`, which RDR 0002 excludes from escape lists. Missing artifact state MUST NOT be maskable behind an escapable refusal class." — (NC, `0007:C8`)

- [REQ-45] "The refusal MUST name what was missing, per row and per atom: for every undecidable row, its `(RuleID, SourceLocator)` and, for each of its unevaluable atoms, the referenced key, the block, and a reason drawn from a closed set — `absent` (the key was not in the view) or `uncomparable` (the key was present and its value was not compared to a verdict)." — (NC, `0007:C8`)

- [REQ-46] "`uncomparable` is defined by the OUTCOME, not by which component produced it, and covers all three ways a present key fails to decide: the seam answered unevaluable (A18); the atom carries a foreign or missing literal on `OpExists`; or the seam is NIL, so no comparison could be attempted." — (NC, `0007:C8`; TS rows 12 and 22)

- [REQ-47] "The closed set stays at two members and every unevaluable atom of an undecidable row therefore carries a reason — the payload's per-atom completeness obligation above admits no third state and no omitted entry. A nil seam is not itself a reason" — (NC, `0007:C8`)

- [REQ-48] "`Refusal` carries the field `Undecided []UndecidedRow`, where an `UndecidedRow` names its `RuleID`, `SourceLocator`, and `Atoms []UndecidedAtom`, and an `UndecidedAtom` names its `Key`, `Block`, `Operator`, `Literal`, and `Reason`." — (NC, `0007:C8`) "The payload is named surface, not shape-by-description; every assertion below and RDR 0005's renderer type against it."

- [REQ-49] "`Reason` is a kernel-owned closed constant set — `ReasonAbsent`, `ReasonUncomparable` — spelled and enumerated the way `resolve.go::RefusalKind` / `RefusalKinds()` already spell the refusal taxonomy: a named STRING type with exported constants plus an exported enumerator `Reasons() []Reason` returning the set in declaration order, mirroring `RefusalKinds()`. The enumerator is required surface, not an implementation choice" — (NC, `0007:C8`)

- [REQ-50] "`Block` is the same exported constant type the atom carries, not a separately-spelled payload value." — (NC, `0007:C8`) Pairs with REQ-6 (widened by §D12).

- [REQ-51] "The reason values are PRODUCED BY THE VERDICT PASS, not re-derived: the per-atom evaluation that decides an atom unevaluable emits that atom's payload entry at the same step, and the kernel carries the entries forward with the row verdict. A second walk over the atoms after the row is known undecidable is FORBIDDEN — it would consult the seam twice for every present-key atom, and the seam is not required to be pure or cheap." — (NC, `0007:C8`)

- [REQ-52] "This payload replaces the single-valued `Refusal.Guard` text — and with it the shipped selection of ONE representative row (`gate`'s `slices.MinFunc` over the undecidable set …). Every undecidable row appears in the payload, so there is no representative to choose; the `MinFunc` call is retired, not re-typed." — (NC, `0007:C8`; TS row 20)

- [REQ-53] "`Refusal.Rows` continues to carry every undecidable row." — (NC, `0007:C8`)

- [REQ-54] "The kernel MUST evaluate every atom of every survivor — no short-circuit on a decided block — and MUST sort the payload so it is a function of the input tuple (RDR 0001 REQ-1), never of atom or row order." — (NC, `0007:C8`)

- [REQ-55] "The ordering is therefore the tuple `(RuleID, SourceLocator, key, block, operator token, literal)`, compared field by field in that order. Those are the row's identity plus the atom's four §D1 fields, so the tuple is total by construction: two entries equal on all six name the same atom, which the kernel MUST NOT report twice." — (NC, `0007:C8`; TS row 17, which requires the tie case: one row carrying two atoms over ONE key, permuted in the slice, must yield equal payloads)

- [REQ-56] "the kernel MUST compare the literal as the exact bytes the normalizer emitted" — the set-literal byte form is **§D13's canonical JSON array: members sorted, duplicate-free, compact encoding**, declared by RDR 0002. The record states this condition as unmet ("Until 0003 declares it, two `in` atoms on one key differing only by literal-set ORDER are not distinguishable by this tuple"); §D13 now supplies it, so the totality claim holds and the sort key needs no interim caveat. — (NC, `0007:C8` + JDR 0001 §D13; see standing correction 2 and `deviations.md` D1)

- [REQ-57] "\"MUST NOT report twice\" is a kernel obligation, not a property of the sort: the kernel emits one entry per unevaluable atom, and an atom is identified by the six-tuple, so duplicate suppression is by construction of the emit loop — the sort orders entries, it does not dedupe them." — (NC, `0007:C8`)

- [REQ-58] "What the payload does NOT promise: a row pruned as GuardFalse under `F ∧ U = F` is not a survivor, so its absent keys are not reported. The guarantee is \"never a plan from absence,\" not \"every absence reported\"; an absence is named exactly when it blocked a verdict." — (NC, `0007:C8`; FM "A pruned row reports nothing")

## `Row.RequiresOwned` narrowed (`0007:C9`)

- [REQ-59] "Row.RequiresOwned names the owned tag keys the row's post-guard transition depends on — the keys its `Writes` require, which per RDR 0002 includes an authored clear (normalization renders it as a `<clear>` write)." — (NC, `0007:C9`) Phase 1: "narrow `Row.RequiresOwned`'s doc contract" — (IP, Phase 1)

- [REQ-60] "Guard decidability is not RequiresOwned's job: guard-input coverage is enforced by the domain rule above." — (NC, `0007:C9`; MC `authority` row "Guard-input coverage … **domain rule**, NOT `RequiresOwned`")

- [REQ-61] "Listing a guard-read owned key in RequiresOwned remains legal and yields the more precise owned_state_unavailable diagnosis among survivors — and is the ONLY way a guard over an owned tag is protected from a caller-supplied observed tag satisfying presence (A13): the guard's key set is NOT required to be a subset of RequiresOwned, because guards legitimately read observed and recognized tags." — (NC, `0007:C9`; TS row 16)

- [REQ-62] "On an escape row, whose `Writes` MUST be empty (RDR 0009), RequiresOwned is therefore empty (A21)." — (NC, `0007:C9`; TS row 15 — the test MUST build the escape row explicitly with empty `Writes` and empty `RequiresOwned`; shipped `fixtures_test.go::escapeRow` is non-conforming on BOTH fields and must not be reused as-is)

- [REQ-63] "Who populates the field is RDR 0002's (JDR 0001 §JD-3)." — (NC, `0007:C9`) Out of scope for this build; recorded as a boundary.

## Survivor membership and the aggregation veto (`0007:C10`)

- [REQ-64] "SURVIVOR MEMBERSHIP. A row whose guard is GuardTrue or GuardUnevaluable is a SURVIVOR; only GuardFalse rows are pruned. The owned-state scan runs over survivors, so an unevaluable row's RequiresOwned keys DO raise owned_state_unavailable, and the undecidable check runs over the same set (shipped partition in `gate`)." — (NC, `0007:C10`)

- [REQ-65] "The owned scan runs ONCE over the whole survivor set, and `MissingOwned` is the deduplicated, sorted union of the missing keys across survivors" — (NC, `0007:C10`; TS row 23: two survivor rows requiring the same absent owned key yield ONE entry, and permuting the rows yields an equal slice)

- [REQ-66] "Aggregation is resolution-level: if any surviving candidate row's guard is GuardUnevaluable, the resolution MUST refuse `guard_unevaluable`; a decided-GuardTrue sibling MUST NOT be selected while an unevaluable candidate exists." — (NC, `0007:C10`; TS row 6; MVV scenario 3)

## D8 ratification (`0007:C11`)

- [REQ-67] "Deviation D8 (guard-FALSE prunes first; a pruned row contributes neither candidacy nor an owned-state obligation) is ratified as normative, conditional on the domain rule: pruning is safe exactly because GuardFalse can only arise from decided atoms — value comparisons over present keys, or existence atoms — never from absence folding into a value comparison." — (NC, `0007:C11`; TS row 2)

## Minimum Viable Validation (`0007:MVV`)

- [REQ-MVV] "Against the real kernel with a real atom — no stub — the masking probe inverts: one candidate row whose guard is a value atom over an absent key, whose `RequiresOwned` is SATISFIED, and a modeled `no_match` escape row yield `Refusal.Kind == guard_unevaluable`, a nil `Plan`, the absent key in the payload, and the row in `Refusal.Rows`; the same table with the key present and the value decided FALSE prunes and escapes (D8 preserved). … Third scenario: *unevaluable-blocks-true-sibling* — row A unevaluable beside row B decided TRUE → refusal naming A, no plan. Fourth: an `exists = false` atom over the absent key decides TRUE and the row is selected — the one sanctioned route from absence to a verdict." — (MVV, `0007:MVV`)

  The satisfied `RequiresOwned` in scenario 1 is load-bearing: "A row carrying BOTH an absent owned key and an unevaluable atom yields `owned_state_unavailable`, not `guard_unevaluable` — that pairing is scenario 13's subject, and putting it in the masking probe would test the owned sweep while claiming to test the domain rule."

## Phase obligations (testable, outside `normative` fences)

- [REQ-68] "migrate every fixture so the frozen suite drives verdicts THROUGH atoms and a value stub — never by injecting row verdicts, which would leave the kernel combinator off the tested path (premortem P-20) — and re-run it (A3 spike)." — (IP, Phase 1; TS preamble repeats the same constraint for every scenario)

- [REQ-69] "Phase 1's exit condition is that count with the field actually removed — not the 154 the superseded superset spike … returned." The frozen suite is expected to pass re-encoded EXCEPT Fixup-1d, which is RE-DECIDED (REQ-17). — (TS preamble)

- [REQ-70] "encode the domain-rule matrix as kernel tests (operators × present/absent × `all`/`unless` × {T,F,U} combinations, empty blocks, two-row pattern, escape-set scoping, foreign existence token); classify each authored guard in the reference fixtures as safe-or-migration." — (IP, Phase 2)

- [REQ-71] "Export a small contract-test function for the value seam (present value × literal × operator; unparseable value → unevaluable, never false) that RDR 0003's implement stage instantiates against its evaluator. It is `resolve.TestGuardEvaluatorContract(t *testing.T, seam GuardEvaluator)` — named here because it is a CROSS-RDR surface 0003 must call by name" — (IP, Phase 3) The exact name and signature are normative surface.

- [REQ-72] The `contains` leg of REQ-71 and TS scenario 8's present-key half are **writable now, against §D13**, not deferred: a set crosses in `Tag.Value` as "its canonical JSON array — members sorted, duplicate-free, compact encoding". The record's "Blocked on one RDR 0003 declaration" (`0007:2159-2170`) and its naming of 0003 as declarer (`0007:2182`) are superseded by JDR 0001 §D13, which assigns the declaration to RDR 0002. Per BUILD-ORDER, 0007 writes these against §D13 and 0002's later run must match, not redefine. See standing correction 2 and `deviations.md` D1.

- [REQ-73] "Every scenario is a kernel test in `internal/resolve` (package `resolve_test`)" — (TS preamble)

- [REQ-74] "Any future edit to the combination tables MUST be re-mutation-tested against" scenarios 3–5, which kill the FALSE-vs-UNEVALUABLE dominance mutant that "survives all 154 frozen tests". — (TS, coverage floor)

- [REQ-75] "Guard-path observed-tag substitution: a caller-supplied observed tag turns `guard_unevaluable` into a SPECIFIC verdict — assert the exact disposition (`--tag X=v` against `X eq v` ⇒ row selected and a plan; against `X eq w` ⇒ row pruned, no plan), never merely \"not `guard_unevaluable`\", which the masking path would also satisfy" — (TS row 16)

- [REQ-76] "Two-row absence pattern (`X exists = false` row ‖ `X exists = true` + `X eq v` row), plus a bare-value-row negative control that must refuse" — (TS row 9)

## Scope boundaries (negative contract — asserted as non-obligations)

- [REQ-77] "`assemble`, `gate`, `missingOwned`, `escapeOrRefuse`, and `Resolve` are untouched in control flow (the payload is a data change `gate` populates, not a reordering)." — (TD) The change "is confined to `internal/resolve/resolve.go`" plus test fixtures.

- [REQ-78] "The domain rule is scoped to guards, so the match pattern still folds absence into non-match. `resolve.go::TagSet.matches` returns false on `!ok` … This is deliberate … and NOT closed by this RDR". — (FM) The match-tag path MUST NOT be changed by this build.

- [REQ-79] "verdicts keep `GuardUnevaluable` / `guard_unevaluable`; rejected: a sixth refusal kind (reopens RDR 0001's closed taxonomy). The absent-key payload is a field on `Refusal`, not a new kind." — (LBD, Naming) RDR 0001 REQ-7 pins the kind set at exactly five; this build mints no sixth `RefusalKind`.

- [REQ-80] "Evaluation stays linear in the atom count per candidate row … and the payload sort is over undecidable rows only, on the refusal path. `Resolve` is pure and in-process; no new I/O." — (PE) No performance dimension is claimed; this is a shape constraint, not a benchmark.

---

## ASSUMPTIONS

- ASSUMPTION: `BlockMatch`'s literal string value is `"match"`, by symmetry with
  `BlockAll = "all"` and `BlockUnless = "unless"` (REQ-6). Grounded: §D12 names
  only the constant, and `0007:C1` fixes the value of every other member of the
  type verbatim; the type is an "exported named STRING type" so a value is
  required. If §D12's home later spells it otherwise, this is a one-line change
  with no behavioral consequence in this RDR — no 0007 clause reads the
  `match` block's bytes.

- ASSUMPTION: `BlockMatch` atoms are NOT evaluated by the guard pipeline in this
  RDR. Grounded: every 0007 clause that names blocks names `all` and `unless`
  only — `0007:C5` ("empty `all` block", "empty or omitted `unless` block"),
  `0007:C6` ("The row verdict is `all_result ∧ ¬(unless_conj)`"), `0007:C1`
  ("block ∈ {all, unless}"). §D12 adds the constant to the type so one type
  serves both; it does not extend the K3 verdict formula, which `0007:C6` fences
  normatively. Match-tag semantics remain RDR 0001's closed-world
  `TagSet.matches` (REQ-78).

- ASSUMPTION: "unevaluable" as a payload-emitting state applies only to atoms of
  rows that end up SURVIVORS and undecidable. Grounded: REQ-58's stated
  non-promise (pruned rows report nothing) plus REQ-45's scoping ("for every
  undecidable row"). The emit-at-verdict-time rule (REQ-51) therefore produces
  entries that are discarded for pruned rows; only the emit ORDER is fixed, not
  the retention.

- ASSUMPTION: `Reason` is a named string type (not an int) with exported
  constants whose byte values are `"absent"` and `"uncomparable"`. Grounded:
  REQ-49 fixes the SHAPE explicitly ("a named STRING type with exported
  constants plus an exported enumerator"), and REQ-45 spells the two members as
  `absent` and `uncomparable` in backticks — the same convention `0007:C3` uses
  where it declares byte values normative.

- ASSUMPTION: `UndecidedRow.SourceLocator` and `RuleID` are the same types the
  shipped `resolve.RowRef` already carries for those names. Grounded: REQ-55
  calls the pair "the row's identity" and REQ-53 keeps `Refusal.Rows` (a
  `[]RowRef`) alongside; a second spelling of one identity would make the
  payload's row key a mapping, which `0007:C1`'s (1) forbids for the atom half.

- ASSUMPTION: the exported contract-test function (REQ-71) lives in
  non-`_test.go` source in package `resolve`, since RDR 0003's implement stage
  must import it by name from outside the package. Grounded: "a CROSS-RDR
  surface 0003 must call by name"; Go test files are not importable. The
  Existing Infrastructure Audit marks the harness "Build" and notes "only value-
  operator contract tests are exported".

- ASSUMPTION: "the same exported constant type the atom carries" (REQ-50) means
  the payload field is declared as `Block` (the named type), not `string`.
  Grounded: the clause exists specifically to forbid "a separately-spelled
  payload value", and `0007:C1` states the requirement "is unstatable while the
  type and its members are unnamed".

- ASSUMPTION: `Reasons() []Reason` returns a fresh slice per call, mirroring the
  shipped `RefusalKinds()`. Grounded: REQ-49 says "mirroring `RefusalKinds()`"
  and "spelled and enumerated the way `resolve.go::RefusalKind` /
  `RefusalKinds()` already spell the refusal taxonomy" — the mirror is on the
  shipped implementation, which is the readable source of truth for the idiom.

- ASSUMPTION: the §D13 JSON array (REQ-56, REQ-72) is compared by the kernel as
  opaque bytes — the kernel neither parses nor re-canonicalizes it. Grounded:
  REQ-22 ("performing no parsing, case-folding, or coercion of its own"),
  REQ-28 (kernel performs no canonicalization), and §D13's own framing
  ("Encoding is carriage, not declaration semantics"). Set-member comparison
  for `contains` is the SEAM's, under RDR 0003's typed semantics (REQ-11).

- ASSUMPTION: Phase 4 (authoring/diagnostics handoff to RDR 0003) mints no code
  and no test in this build. Grounded: its text is entirely "Hand over …" and
  "the REQUEST that 0003 declare …"; no clause states a kernel obligation.
  Note that its one substantive request — the set-element encoding — is already
  answered by §D13 (REQ-72), so the handoff carries a citation, not an open ask.

---

## QUESTIONS

- **Q1 — Does `BlockMatch` on the `Block` type imply match tags become atoms on
  `Row` in this phase?** Two readings, materially different. (a) §D12 adds the
  constant so ONE type serves both the guard payload and a future match-atom
  representation, but `Row`'s match pattern stays the shipped `TagSet`-style
  map in this RDR — no behavior change, REQ-78 holds. (b) §D12's "one type, no
  separate slice on `Row`" means match tags travel in the SAME atom slice as
  guard atoms, discriminated by `Block == BlockMatch` — which would put match
  evaluation inside `evaluateGuard` and collide head-on with `0007:C6`'s fenced
  row verdict `all_result ∧ ¬(unless_conj)` and with REQ-78.
  **Proceeding under (a)**, recorded as an ASSUMPTION above: it is the only
  reading that leaves every `0007` normative fence intact, and `deviations.md`
  D1 scopes §D12 to "`Block` gains `BlockMatch`" — a change to the constant set
  of a type, not to the row's shape. If (b) is intended, the widening is far
  larger than D1's scoped re-entry describes and needs its own re-entry.
  No predecessor settles it: RDR 0003 (which owns the match/guard authoring
  split) is unimplemented, and `0001/artifacts/req-list.md` predates atoms
  entirely.

- **Q2 — Is `Refusal.Undecided` populated for an unevaluable ESCAPE row's
  refusal?** REQ-41 says an unevaluable escape row "yields `guard_unevaluable`
  in place of the candidate-set refusal", and REQ-40 says the escape set is
  gated "by DELEGATION" through the same `gate` function, which by REQ-52
  populates the payload — so the escape set's payload rides along. The
  `disposition` mini-check agrees ("Unevaluable ESCAPE row … escape set's
  payload"). **Proceeding on that reading** — it is stated in the mini-check
  table, so this is closer to a confirmed ASSUMPTION than a genuine fork, and
  it is recorded here only because the two normative clauses reach it by
  composition rather than by statement. TS row 14 pins the kind, not the
  payload contents; the build asserts payload contents too.
