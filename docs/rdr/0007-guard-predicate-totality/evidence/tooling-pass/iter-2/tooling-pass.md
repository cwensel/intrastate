Model: claude-opus-5[1m]

# Tooling Pass — RDR 0007 guard predicate totality (iteration 2)

Mechanical adherence sweep, run as the Stage 7 pre-step to the Finalization
Gate. Post-mutation regression check against a draft that has been rewritten
substantially since iteration 1.

**Why iteration 2.** Run 1 (`../tooling-pass.md`, 2026-08-11) swept the
pre-demotion draft and passed; the RDR locked Final that day. JDR 0001 then
demoted it (`Status: Draft [revised from Final 2026-08-12 — re-verify A3, A16,
A17, A18, A19, A21, A22]`) and it was re-proposed, re-refined, re-resolved, run
through all four `foundational` lenses a second time (`*/iter-2/`), and
reconciled three times. The record set grew from 16 to 26 and ten of eleven
normative blocks were touched. This is a genuinely different document.

**Loop-breaker check (Stage 7 contract).** Run 1's four fixed findings were
read before writing any return pointer. None is re-reported here: the MVV
scope contradiction and the Risks mitigation resting on it (run 1 findings 1–2)
are both absent from the current text, and the two run-1 quotation fixes (RDR
0002 "renders as", RDR 0003 attribution) still hold. Every finding below is
**new**, created by the re-work — no routing loop.

## CHECK 1 — Template section coverage

Delegated. Every **Required** (spine) section is Present-substantive: Metadata
(incl. `Seam Lineage`, `Predecessors`, `Overrides`), Problem Statement,
Critical Assumptions (26 records, all four fields each), Proposed Solution /
Approach / Technical Design / Normative Contracts, Decision Rationale (both
greppable verdict lines present — `Premortem: hardened`, `Joint-check: fired`),
Alternatives Considered, Context, Research Findings, Trade-offs, Implementation
Plan / Prerequisites / MVV, Validation / Testing Strategy, References.

**Conditional sections cleanly deleted** (PASS per the false-positive guard):
Round-Trip / Inverse Invariants (the Mini-checks block explicitly records that
the cue does not fire — "no inverse invariant is claimed"), Illustrative Code,
Day 2 Operations, New Dependencies, Phase 2: Operational Activation (this
RDR's Phases 2–4 are code phases).

**Conditional sections present and filled**: Load-Bearing Decisions (4),
Capability Dependencies (7 rows), Existing Infrastructure Audit (6 rows),
Alternatives 1–2 + Briefly Rejected, Performance Expectations.

No `_Draft placeholder._`, no seed-skeleton header, no surviving bracketed
template instruction from any template version. The `## Finalization Gate`
section reads `_Not yet run — Stage 7 authors the gate responses._` — the
expected mid-relock state, replaced by the gate.md pointer at lock, not a
defect.

*Advisory, not a finding:* Critical Assumptions is nested as `###` under
Research Findings rather than the top-level `##` TEMPLATE.md orders second.
Checked against the corpus — 0001, 0002, 0005, 0006 and 0009 all do the same.
This is the established house shape; changing 0007 alone would make it the
outlier. No action.

Verdict: **PASS**.

## CHECK 2 — Method label vocabulary

26 Evidence Records. **Two off-vocabulary labels found and fixed in-pass**
(both created during the re-work; run 1 saw neither):

1. **A6b — `Prior Decision`** is not one of the sanctioned eight. The record
   relies on a property fixed in another design document (JDR 0001 §D3, quoted
   in full in its own Evidence line) and carries no obligation of its own —
   which is exactly README §Verifying load-bearing claims' definition of
   **`Peer RDR`** ("relies on a property defined in another RDR. Evidence: RDR
   ID + section"). Relabeled; Evidence already carried document + section, so
   nothing else changed. Note run 1 recorded A6b as `Source Search`/`Pending`;
   the label drifted when Stage 6 closed it against the JDR.
2. **A19 — `Source Search (…) + \`internal/cli/clierr\` source`.** The trailing
   conjunct is a *source pointer*, not a second Method, so the label parsed as
   a compound whose second half is off-vocabulary. Both halves are the same
   method. Folded the pointer into the parenthetical: `Source Search (RDR 0005
   envelope; JDR 0001 §JD-8; \`internal/cli/clierr\` source)`.

After the fixes, the distribution is: Source Search ×17, Spike ×4, Design
Decision ×2 (incl. one sanctioned `Spike + Design Decision` compound),
Derivation ×1, Prior Art ×1, Peer RDR ×1. Every record has all four fields.

Verdict: **PASS** (after the two in-pass fixes).

## CHECK 3 — Source Search self-reference

Delegated. No findings. Every citation into
`docs/rdr/0007-guard-predicate-totality/` is made by a **Spike** record (A3,
A23, A24, A26 → the re-spike artifacts) or a **Prior Art** record (A4 →
`evidence/research/resolve-citations.md`) — never by a Source Search record.
The one Source Search record citing an artifact directory, **A27**, points at
*peer* RDR 0003's `evidence/spikes/guard-fixture.toml`, which is external
evidence; the quoted spellings (`exists = true`, `domain = [true, false]`) were
confirmed in that file.

Verdict: **PASS**.

## CHECK 4 — Docs Only on load-bearing claims

Zero `Docs Only` records (grep: 0 occurrences of the label or its variants
anywhere in the document). Nothing blocks on this axis.

Verdict: **PASS**.

## CHECK 5 — Symbol resolution of Source Search / Spike anchors

Delegated; all anchors resolve **in their cited files** on branch
`via-claude`. **One form finding, fixed in-pass.**

- **A23 cited bare `file:line`** — `adversarial_test.go:313,375`;
  `fixup_test.go:200,236` for the four `reflect.DeepEqual`-over-`Refusal`
  sites. All four lines were verified to be exactly the described comparison,
  so the claim is true and the anchors currently hit — but CHECK 5 makes a bare
  `file:line` with no symbol a finding in its own right. Resolved each to its
  enclosing test symbol and rewrote as `path::Symbol`:
  `::TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder`,
  `::TestAdv3b_MissingOwnedPayloadMustNotDependOnTableRowOrder`,
  `::TestFixup3c_AmbiguousMatchRowsPayloadMustNotDependOnTableRowOrder`,
  `::TestFixup3c_DegradedEscapeAmbiguityPayloadMustNotDependOnRowOrder`.
  No bare `file:line` anchor remains in the document.

Everything else resolved, with the structural facts the normative blocks assert
confirmed rather than assumed:

- Kernel: `resolve.go::assemble`, `::Resolve`, `::gate`, `::evaluateGuard`,
  `::missingOwned`, `::escapeOrRefuse`, `::RefusalKinds`,
  `::KindGuardUnevaluable` (doc string matches A18's quote verbatim),
  `TagSet::Lookup`/`Len`/`has`/`matches`.
- CLI: `clierr.go::CLIError` — A19's quotes verify verbatim, including the
  "Extend with new optional fields … `omitempty`" invitation and `Detail` being
  a `string` documented "May be multi-line"; field set is exactly
  `Code`/`Message`/`Param`/`Detail`/`Hint`. `root.go::NewRootCmd`.
- 13 test anchors resolve (ADV-1/1b/2/3/3b/4/5, Fixup1d, the gate-uniformity
  test, Req7/25/33/36/37).
- Fixtures: `fixtures_test.go::escapeRow`, `::fixtureGuards`,
  `boundary_test.go::exportedKernelSymbols`, `::parseKernelPackage`. **A25's
  collision re-verifies on source** — `escapeRow` sets both `Writes` and
  `RequiresOwned: []string{"status"}` beside a non-empty `Escape`. **A23's
  grounds (1) re-verifies** — zero `Refusal{}` literals in the test files.
- JDR 0001 §D1–§D4, §JD-1/2/3/4/6/7/8/9/12 and principle P5 all resolve; §JD-4
  is confirmed *not* among the Closed entries, exactly as A12/A21 state.
- Peer-RDR quoted clauses resolve across 0001, 0002, 0003, 0004, 0005, 0009,
  including the two load-bearing **demonstrated negatives**: `canonicaliz`
  appears zero times in RDR 0002 (A22) and `RequiresOwned` zero times in RDR
  0009 (A21). RDR 0003's section heading "Phase 1: Predicate Model" exists.

No phantom, never-built, renamed, or moved symbols.

Verdict: **PASS** (after the in-pass anchor rewrite).

## CHECK 6 — Status consistency

Delegated. No `Pending` or `Unverified` Status value exists: all 26 records
read `Verified` (22) or `DOWNGRADED` (4 — A12, A19, A22, A25). So no unsettled
assumption can be leaned on as settled fact by construction.

Each DOWNGRADED assumption's **blocked half** was checked against every prose
site that references it, and all four are hedged consistently:

- **A19** — the structured `omitempty` field is never asserted as granted; the
  ownership paragraph says CLI rendering "is blocked on the A19 envelope-field
  grant — or … A19's `Detail`-flattening fallback", the `authority` table marks
  the row "**unresolved — see A19** / §D4 stands until reopened", Prerequisites
  carries it unchecked and explicitly "Does not block lock."
- **A22** — the `authority` table says "**normalizer** (duty not yet a 0002
  clause)"; the PRESENCE clause states the kernel performs no canonicalization.
- **A25** — Testing row 15 and Prerequisites both say the sequencing "is a
  cluster question this RDR does not settle."
- **A12** — scoped to the authoring story throughout ("Nothing in the kernel's
  domain rule depends on the answer").

**Two status-vocabulary collisions found and fixed in-pass.** Both reused the
reserved word `Pending` inside a Status/Evidence line for a half the Status
line does not mark Pending, which grep-collides with Prerequisites' assertion
that "No assumption remains Pending":

1. **A9** Status read `Verified (kernel half); the agreement with RDR 0003's
   subtractive lint algebra is Pending`. The record's own *Downgraded*
   sub-bullet already treats that half as a downgrade, so the Status line now
   reads `… is DOWNGRADED` — the vocabulary the other four blocked-half records
   use.
2. **A27** Evidence closed "which is why this is Pending rather than blocking"
   while its Status reads `Verified (kernel half)`. Reworded to "the producer
   half is carried as a downgrade rather than a blocker" — same meaning, no
   reserved word.

Both are wording-only; neither changed a disposition. After them, `Pending`
occurs exactly once in the document: the Prerequisites line asserting none
remain, which is now literally true.

The Metadata `Status:` re-verify list (A3, A16, A17, A18, A19, A21, A22) is
consistent with the dispositions: A3/A16/A17/A18/A21 Verified, A19/A22
DOWNGRADED. The qualifier self-clears at re-lock.

Verdict: **PASS** (after the two in-pass fixes).

## CHECK 9 — Evidence-field budget (ADVISORY — never blocks)

Four fields exceed the 30-line budget: **A3** 44, **A21** 38, **A19** 34,
**A26** 34. `4 fields over budget, 150 lines`.

No truncation proposed. The question this check asks — has the prose outgrown
the record, is the load-bearing anchor still findable — answers *yes, still
findable* in all four: each opens with its verdict and its stable anchor, and
the mass below is verification content the grounding sweep reads (A3's
re-spike counts and frozen-test enumeration, A21's §JD-3 composition, A19's
`CLIError` field enumeration, A26's three banned-list tests). Carried to the
Gate as a report, per the check's own instruction.

## Findings fixed in-pass (all MECHANICAL)

Per the Stage 7 contract these are repaired here and the sweep re-run, never
routed to Refine (conformance is outside Refine's contract). All five are
drawn from material already in the RDR — no new evidence, no design call, no
disposition changed:

| # | Check | Record | Fix |
| --- | --- | --- | --- |
| 1 | C2 | A6b | `Prior Decision` → `Peer RDR` (off-vocabulary label) |
| 2 | C2 | A19 | source pointer folded into the parenthetical so the label is one sanctioned token |
| 3 | C5 | A23 | four bare `file:line` anchors → `path::Symbol` |
| 4 | C6 | A9 | Status `Pending` → `DOWNGRADED` (reserved-word collision) |
| 5 | C6 | A27 | Evidence "is Pending" → "carried as a downgrade" |

**Sweep re-run after the fixes**: 26 records, 0 off-vocabulary, 0 incomplete,
0 template brackets, 0 placeholders, 0 `Docs Only`, 0 `Unverified`, 0 bare
`file:line`, 1 `Pending` (the Prerequisites assertion). No regression, no new
finding.

## Verdict

**PASS** — no finding survives; proceed to the Gate's written responses.
