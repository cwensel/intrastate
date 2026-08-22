Model: claude-opus-5[1m]

# Repeatability resolve — disposition ledger (iteration 2)

Origin ledger = `diff.md`'s R-1..R-9 and G-1..G-8. Every entry exits once.
Dispositions follow the REPEATABILITY DIFF clause: **pin** / **cut** /
**single-source** / **leave non-normative** / **tiebreaker**, plus the standard
`dismissed-with-cite` for a finding that fails the grounding gate.

Grounding for the source-facing entries was delegated; all seven questions
returned verbatim from `internal/resolve/resolve.go` on the current branch.
Findings the diff itself labeled RUN ERROR were not treated as RDR defects.

## Disagreements

| ID | Disposition | Origin | Section touched |
| --- | --- | --- | --- |
| R-1 | **pin** (fixed) — grounding turned this from "ambiguity" into a live implementation trap. Shipped `GuardResult` is `GuardFalse, GuardTrue, GuardUnevaluable` = `iota` 0,1,2 — the order `F < T < U`, NOT the truth order `F < U < T` the clause names. An integer `min` over the raw constants returns `T ∧ U = T` where the table requires `U`, and passes every other cell. Computed rather than argued: enumerating all 9 cells of `min` over the shipped `iota` values against the RDR's table returns 7 matches and exactly 2 mismatches — the commutative pair `T ∧ U`, which raw `min` gives as `T` where the table requires `U`. That is the worst possible shape for a bug: it passes every other cell. Clause now states the tables are normative, `min` names the lattice only, and the kernel MUST NOT `min` the constant values; the constants are not renumbered (a shipped fact, and renumbering would silently change any existing comparison). | R-1 | Normative Contracts, combination clause |
| R-2 | **pin** (fixed) — reason values are produced by the verdict pass and carried forward; a second walk after the row is known undecidable is FORBIDDEN, because it consults the seam twice per present-key atom and the seam is not required to be pure or cheap. Internal data flow only; exported surface unchanged. | R-2 | Normative Contracts, payload clause |
| R-3 | **dismissed-with-cite** (run error) + **pin** on the residue. `Row.Escape` is `[]RefusalKind` on shipped source; run-3's `bool` contradicts "a non-empty `Escape` list" in four places. The element-type residue is closed incidentally by the G-6 fix, which now names the five kinds. No clause added for the `bool` error. | R-3 | none (Context, via G-6) |
| R-4 | **pin** (fixed) — `GuardEvaluator` is a single-method INTERFACE, not a func type. Split cleanly on the alt-model boundary (run-2 declared a func type, then called `seam.Evaluate(...)` against it — an internal inconsistency the silence permitted). Interface also makes "nil seam" a nil interface value, which the nil-seam rule already assumes. | R-4 | Normative Contracts, SEAM clause |
| R-5 | **pin** (fixed) — the strongest finding in the set, and a genuine hole rather than a spelling gap. A value atom over a PRESENT key under a nil seam is unevaluable, but the closed reason set covered neither case: `absent` is false (the key is present) and `uncomparable` was glossed as "the seam could not compare" (no seam answered). Two runs silently dropped the payload entry — contradicting "for every undecidable row … for each of its unevaluable atoms" — and one noticed and left it open. `uncomparable` is now defined by OUTCOME, not producer, and explicitly covers all three routes (seam answered unevaluable; foreign/missing `OpExists` literal; nil seam). Set stays closed at two; no entry may be omitted. | R-5 | Normative Contracts, payload clause + SEAM nil-seam paragraph |
| R-6 | **dismissed-with-cite** (run error) — the RDR pins `ambiguous_match` downstream of the gate facts with escape reachability consulted after the count, and RDR 0002's closure admits `ambiguous_match` to escape lists. Grounding confirms shipped `escapeOrRefuse` returns `KindAmbiguousMatch` on `len(viable) > 1`. Run-3's non-escape routing would make the modeled `ambiguous_match` escape class dead surface. No rewrite owed. | R-6 | none |
| R-7 | **pin** (fixed) — shipped `missingOwned` already takes `[]Row`, dedupes via a `seen` map and `slices.Sort`s, so run-3's per-row accumulation is wrong against source. But the RDR's silence is what permitted it, and REQ-1 purity is claimed for the whole result while determinism was pinned only for `Undecided`. Clause now states the single sweep and the deduplicated, sorted union — shipped behavior restated, not a change. | R-7 | Normative Contracts, SURVIVOR MEMBERSHIP |
| R-8 | **single-source** (fixed) — "gated identically" and "untouched in control flow" are jointly true only because `escapeOrRefuse` already delegates to `gate`; the RDR asserted both without stating the mechanism, leaving a reader to choose between "inherits the payload for free" and "has its own gating to update". Grounding shows the delegation verbatim (`gate(escapes, in.Guards, view)`). Mechanism now stated once, in the clause that makes the claim. | R-8 | Normative Contracts, GATE-THEN-COUNT |
| R-9 | **leave non-normative** — `Result`'s shape is RDR 0001's surface and this RDR touches only `Refusal`; all three runs correctly declined to invent it. Specifying it here would be the over-spec trap and would fork a peer RDR's type. The `Escaped:true` behaviour the Background asserts is already cited to RDR 0001. | R-9 | none |

## GUESS clusters

| ID | Disposition | Origin | Section touched |
| --- | --- | --- | --- |
| G-1 | **pin** (fixed) — atom field names `Key`, `Operator`, `Literal`, `Block`. Not a re-decision of §D1's grammar: the payload clause requires `UndecidedAtom` to mirror the atom's fields by name and to carry "the same exported constant type" for `Block`, which is unstatable while the fields are unnamed. Scope explicitly fenced — operator vocabulary, literal typing and the parse stay §D1's. | G-1 | Normative Contracts, SEAM clause |
| G-2 | **pin** (fixed) — `OpExists = "exists"`, `LiteralTrue = "true"`, `LiteralFalse = "false"`, lower-case. This is the sharpest class in the set: a contract whose entire content is byte equality, stated without the bytes, plus a cross-component MUST (A16) telling the normalizer to emit "exactly these values". A producer cannot emit what the contract does not spell. | G-2 | Normative Contracts, existence-operator clause |
| G-3 | **pin** (fixed) — `Block` is an exported named string type with `BlockAll = "all"` and `BlockUnless = "unless"`. A26 already swept "the `Block` constants" as known-to-exist and exported while nothing enumerated them. | G-3 | Normative Contracts, SEAM clause |
| G-4 | **pin** (fixed) — `Reasons() []Reason` is required surface, mirroring `RefusalKinds()`. The clause already mandated the pattern and "a test can pin the set exactly", which is unsatisfiable without an enumerator. Two runs independently invented it from that sentence. | G-4 | Normative Contracts, payload clause |
| G-5 | **pin** (fixed) — `Reason` is a named string type, resolved by the same edit as G-4 (the `RefusalKind` pattern the clause cites is `type RefusalKind string`; grounding confirms). Values `absent` / `uncomparable` were already in-text. | G-5 | Normative Contracts, payload clause |
| G-6 | **pin** (fixed) — the fifth kind is `unmodeled_outcome`. Not a defect of this RDR (the taxonomy is RDR 0001's and no kind is reopened), but asserting "exactly five" while naming four made all three runs guess. Now cited where the taxonomy is introduced. | G-6 | Context, RDR 0001 boundary |
| G-7 | **pin** (fixed) — `resolve.TestGuardEvaluatorContract(t *testing.T, seam GuardEvaluator)`. A cross-RDR surface RDR 0003's implement stage must call BY NAME; a caller cannot instantiate a function whose name and signature the contract leaves open. | G-7 | Implementation Plan, Phase 3 |
| G-8 | **dismissed-with-cite** (not a defect) — the CLI transport is a deliberate, stated deferral (JDR 0001 §D4 vs A19), and all three runs reproduced the deferral faithfully. The diff recorded it only to mark it a non-candidate. Deferral is doing its job. | G-8 | none |

## Agreement (no action)

The diff's agreement paragraph is large and covers everything load-bearing: the
two-level payload with `Refusal.Guard` removed and no sixth kind; the closed
reason set; the four-way per-atom dispatch; fail-closed drift; the strong-Kleene
verdict shape with the `unless` term dropping out; the gate-then-count order with
the resolution-level veto; the six-field total sort tuple with `MinFunc` retired;
the narrowed `RequiresOwned`; `guard_unevaluable` as non-escapable; and no new
panic or error channel. Three independent models reconstructed all of that
identically — the determinacy this lens exists to measure, and it holds where it
matters.

The findings are concentrated in Go SPELLING (G-1..G-5, G-7, R-4) rather than in
rule content, with two genuine rule gaps (R-5's uncovered reason, R-2's
unspecified producer) and one implementation trap the RDR's own phrasing set
(R-1). That distribution is the healthy signature for a draft this far along.

## Lens health

Healthy per the lens's Expected signal. GUESS markers cluster on the same
contracts across runs (24 / 46 / 24 markers; three distinct model stamps
`claude-opus-5`, `claude-fable-5`, `claude-sonnet-5`), which is the "3–5 specific
interfaces" shape rather than the identical-but-confidently-wrong overfit that
would force a model switch. One divergence (R-4) split on the alt-model boundary
and was flagged as the stronger signal it is.

## Needs (re)verification — carried to Stage 6

- **A27 (NEW, Pending)** — the five literal byte values this pass spelled
  (`"exists"`, `"true"`, `"false"`, `"all"`, `"unless"`) are consistent with what
  RDR 0002's normalizer and RDR 0003's grammar already use. Kernel half decided
  here; producer half rides the same unlanded 0002 duty as A16/A22/A24. Fails
  closed on mismatch (foreign token ⇒ value atom; foreign literal ⇒
  `uncomparable`), which is why it is Pending rather than blocking.
- **A26 (already Pending) — enumeration widened.** The banned-name sweep must now
  also cover `Reasons`, `Block`, `BlockAll`, `BlockUnless`, and
  `TestGuardEvaluatorContract`. Still discharged by A3's re-spike.
- **No assumption flipped from Verified to Pending.** R-1, R-7 and R-8 were each
  grounded verbatim against shipped source this pass, so they are Verified facts
  newly STATED rather than new claims; A2's abstract `min` derivation is
  unaffected (it derives over the lattice, not over the Go constants).

## Charted

None. No finding was net-new scope: every entry either pinned a contract this RDR
already owns, restated shipped behavior the RDR was silent on, or was dismissed.
