Model: claude-opus-5[1m]

# 3amigo Consolidation — iteration 3 (re-entry)

Three isolated persona passes (PM 3, Implementer 11, QA 7 = 21 findings), each
delta-scoped to the re-entry: the rehomed tag declaration model, A8's closure via
JDR 0001 §JD-4, §JD-13's single-valued marker, §JD-14's escape-row confirmation,
A5's re-verification, and the four still-open records (A10, A12, A14, A15).

**Isolation held.** No persona file references another persona's output or
finding IDs (grepped for cross-persona tokens: zero hits in all three files). All
three carry a `Model:` stamp. Consolidation below is **mechanical**: hotspots are
computed by which RDR passage two or more personas independently named, not by
re-judging their findings. Overlap marks a hotspot passage, not a validated
finding; a finding raised by exactly one persona is not thereby weaker.

## Origin ledger (this iteration)

| # | Passage anchor | Personas | Hotspot | Concern |
| --- | --- | --- | --- | --- |
| **1** | The single-valued marker's **authoring location** — `Capability Dependencies` decl-model row (`:1055`); `authority` decl-model row (`:937`); RDR 0002's receiving clause (`0002:342-350`) and schema list (`0002:217-220`) | P1-2, P2-7, P3-7 | **HOTSPOT (3)** | §JD-13 homed the field in this RDR's normative clause, but `single-valued` occurs **zero** times in RDR 0002, whose type-model clause and schema list both enumerate the same four fields. The field has no authoring location, so it cannot be written in a fixture or reach RDR 0006's lint input. The §JD-13 checkbox reads `Done` while a mandatory blocking code (`graph-single-valued-state`) still has an unwritable producer — the same unhomed-producer defect §JD-13 closed, reappearing one layer down. |
| **2** | The document's own **four-field enumerations** — `Capability Dependencies` row (`:1055`), `Scope Verification` (`:1740-1742`), `Load-Bearing Decisions` › Operator semantics (`:1012-1016`) | P1-1, P2-7 | **HOTSPOT (2)** | Five-field sites (Metadata `:17`, Approach `:543`, normative clauses `:687`/`:715`/`:729`, both `authority` rows `:937`/`:940`) disagree with four-field sites on the size of the contract this RDR owns. These are exactly the passages a planner reads to size Phase 1. |
| **3** | `Normative Contracts` optionality clause (`:704-712`) + single-valued clause (`:714-725`) — the term **"conforming evaluation view"** | P3-1, P2-1 (adjacent), P2-2 | **HOTSPOT (2+)** | The term occurs exactly twice cluster-wide, both in these two new clauses, with no definition, no violation outcome, no `disposition` row, and no trace step. Lint's soundness is conditional on a view conformance nothing enforces. Compounding: the declaration clause says a tag **MAY** carry an optionality marker (`:686`) with **no stated default**, and the RDR's own `evidence/spikes/guard-fixture.toml` declares seven tags, none carrying one — so every group is withheld under one default and green under the other. The RDR answers the parallel question for domains ("Leaving a dimension undeclared is not an opt-out") but not for presence. |
| **4** | `Normative Contracts` "can refuse" clause (`:859-865`) vs. the optionality clause (`:704-712`) | P2-1 | single | The optionality clause declares two states (always-present / optional); "can refuse" tests a third, differently-named property — `present-for-every-reachable-predecessor`. Synonym, strictly-stronger graph property, or an undeclared field? Reading 2 makes it an A12-gated graph query, which contradicts A12's own survivability argument and the clause's own "decide this syntactically over declarations so the test is total". |
| **5** | `Normative Contracts` participation clause (`:777-782`) vs. `trace` step 3 (`:975`) | P2-3 | single | Participation is defined over `all`/`unless` atoms only, but step 3 counts `status eq "Draft"` as a participating product dimension — authored under `[rule.match.status]` in the fixture and carried in a **separate `Match []Tag` field** in the shipped kernel (`internal/resolve/resolve.go:178-185`), distinct from `Guard`. The trace's worked example contradicts the clause it is cited to justify, and step 3 is the witness the participation clause rests on. |
| **6** | `Normative Contracts` escape-row clause (`:799-807`), §JD-14 | P2-4 | single | §JD-14 resolved this against RDR 0006 only. RDR 0002 (`0002:328-334`) and the shipped kernel (`Resolve` step 2 — "candidate rows are the **non-escape** rows for that outcome") both partition escape rows *out* of the candidate set, so the overlapping escape/guarded pair this clause calls "the ambiguity RDR 0001 refuses" is not an ambiguity at runtime. Separately, a bare escape row is simultaneously coverage-green (denotes the whole product) and overlap-blocking against every peer. |
| **7** | Single-valued clause (`:714-725`) vs. **RDR 0006 invariant 5** (`0006:295-296`) | P2-5, P3-2 | **HOTSPOT (2)** | Two mismatches with the named consumer: per-**tag** marker vs. per-**tag class** invariant, and constrains the **evaluation view** vs. constrains **writes**. Decides the field's type (bool vs. class id) and whether `graph-single-valued-state` actually gained a producer. Test consequence: the marker got a rejection test (Scenario 4) but no acceptance test, and the clause's only stated effect ("licenses treating the domain as a partition rather than independent dimensions") never says what the scoped product *is* under each reading — for enum `{a,b,c}`, 3 assignments or 2³? — so the two expected verdicts are uncomputable. |
| **8** | Declaration clause value-kind list (`:688-689`) vs. operator/kind matrix (`:571-577`) | P2-6, P3-4 | **HOTSPOT (2)** | Three vocabularies for one axis: `enum`/`bool`/`int`/`set`/opaque scalar (normative), "enum, boolean, integer, string-like scalar" (matrix), and `exists`'s "optional scalar or optional set-valued tag". The fifth kind has no token, no finite-domain rule, and no row in the domain/kind rejection clause; `set` never appears as an accepted kind for `eq`/`in`. Blocks the table-driven operator×kind matrix test and the parser's value-kind enum. Also: `exists` over an always-present key is a rejection under the matrix but legal-and-vacuous under the presence-dimension clause (`:785-797`). |
| **9** | Domain/kind agreement clause, second sentence (`:727-735`) + A4's predicate-semantic-kind enumeration (`:196-199`) | P3-3, P2-9 | **HOTSPOT (2)** | The out-of-domain-literal rejection is normative but untested and unnamed: Scenario 4's set is "unknown tag, unknown operator, unsupported operator/tag-kind pair, literal parse mismatch" — an out-of-domain literal is none of those — and A4's enumeration omits it, so there is no error code to assert on. Compounding: one clause names two different phase boundaries ("before normalization completes" vs. "before resolution"), leaving the owning package and error surface unfixed. |
| **10** | `int` finite-domain clause (`:696`) — "`int` declares a `{min..max}` bound" | P2-8 | single | Notation or wire spelling? The RDR's own fixture uses separate `min`/`max` keys. Bound inclusivity at both ends is unstated. Feeds the cardinality arithmetic A15's Scenario 3 compares products with. |
| **11** | `MVV` Authorability paragraph (`:1451-1462`) + `Phase gating` row 3 (`:1471`) — A14's survivability argument | P3-5 | single | A14's `Pending`-but-survivable disposition rests on "the failure is caught by the first normalization test", and no scenario is that test. No case asserts an atom authored in `unless` still reports `Block == unless` after normalization — the detector the argument depends on does not exist. |
| **12** | `Testing Strategy` Scenario 3 (`:1550-1562`) — the named A15 discharge | P3-6 | single | "Two equal-cardinality products of differing shape" gives neither a cardinality relative to the published bound (unknown at fixture-authoring time) nor a quantification of "few wide vs. many narrow". The RDR's last self-owned open assumption is discharged by a described experiment rather than a test with a computable verdict. |
| **13** | `Phase gating` table Phase 2 row (`:1470`); `Technical Design` default-on reading (`:628-634`) vs. A10 `Pending` | P2-10, P2-11 | single (2 findings, one persona) | Phase 2 is marked startable for "the product/proof core", but ledger entries 4 and 5 are *construction* inputs, not integration. And §JD-14 settled default-on while A10 — what it defaults over — is still owed at RDR 0006's refine, leaving group construction a provisional contract against which MVV fixture investment may be unsafe. |
| **14** | `Normative Contracts` declaration clauses (`:684-735`) + `Illustrative Code` (`:1038-1046`) | P1-3 | single | The re-entry moved a five-field type system in, but the only illustration is `[rule.guard.*]` — the surface this RDR does *not* own the authoring of. The one newly-owned surface has no illustration, while `Risks and Mitigations` names declaration boilerplate as a live risk whose mitigation points back at RDR 0002. |

## Hotspot summary

Six of fourteen entries are multi-persona hotspots. They cluster on two roots:

- **§JD-13's field was homed in prose but never swept** — entries **1, 2, 7**.
  Three personas reached it from three different questions (can the author get
  the outcome / what do I code / what do I assert). The field is enumerated
  inconsistently *within* this document (2), has no authoring location in the
  peer that owns authoring (1), and does not match the shape its named consumer
  reads (7).
- **The declaration model's edges are under-specified for a first implementer** —
  entries **3, 8, 9**. Presence has an undefined conformance term and no default;
  the value-kind axis has three spellings; and a normative rejection has neither
  a test nor an error kind.

Entries **5 and 6** are the two places where the draft conflicts with **shipped
code and RDR 0002's normative text**, not merely with an unbuilt peer — both in
delta-scoped passages, and both raised by the persona whose question is "what
would I ask in the first hour".

## Not filed

- PM verified and cleared: A8/§JD-4's recording assignment (JDR 0001 §JD-4
  CLOSED, RDR 0006's Status line agrees — both halves match); §JD-14's "no edit
  owed here"; A5's re-verification (`RuleID`/`SourceLocator` at
  `internal/resolve/resolve.go:173-174`); the A11 rehoming and the 0002/0003
  split (RDR 0002 states the same split independently at `0002:909-919`, RDR 0007
  `Final` assigns it here — the model genuinely was homeless, so the scope growth
  is justified by the outcome); A10/A12 homed on RDR 0006's Draft Status line;
  A14 present on RDR 0002's duty list (`0002:897-902`); A15's Scenario 3 written
  to make the needed comparison.
- QA checked and cleared: the §JD-14 escape-row clause as *written* (`:799-807`),
  A8's narrowing plus withheld-claim observability (`:809-823`, with Scenario 6
  asserting positively and failing an empty run), and A5's code anchors.
- No persona reported an OUT-OF-SCOPE finding.

## Delta coverage

Every element of the re-entry's re-verify list was reached by at least one
persona: the declaration model (1, 2, 3, 8, 9, 10, 14), §JD-13 (1, 2, 7),
A8/§JD-4 (cleared by PM and QA), §JD-14 (6, and cleared as written by QA), A5
(cleared by PM and QA), A10/A12 (4, 13), A14 (11), A15 (12).
