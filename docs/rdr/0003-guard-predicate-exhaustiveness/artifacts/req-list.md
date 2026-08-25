# REQ List — RDR 0003 Guard Predicate Exhaustiveness

Phase 0 audit artifact. Every clause below is a testable obligation drawn from
`docs/rdr/0003-guard-predicate-exhaustiveness.md`. Quotes are verbatim, taken
from the projector (`rdr inspect --select <id>`) for fenced elements and read
from the record for prose outside the fences.

Element ids (`0003:C4`, `0003:MVV`) are carried wherever a REQ derives from a
labelled contract, so a later stage can trace the REQ back to its contract.
The record has 26 `C` elements but more than 26 REQs: several fences carry
multiple independent obligations, and testable prose lives outside the fences.

Source sections are abbreviated:

- `NC` = Proposed Solution / Technical Design / Normative Contracts (fenced)
- `TD` = Proposed Solution / Technical Design (unfenced prose)
- `AP` = Proposed Solution / Approach
- `MTX` = Proposed Solution / Approach, operator/kind matrix table
- `AUTH` = Technical Design / `authority` — source-authority census
- `DISP` = Technical Design / `disposition` — input class × outcome
- `LBD` = Technical Design / Load-Bearing Decisions
- `RT` = Technical Design / Round-Trip / Inverse Invariants
- `IC` = Technical Design / Illustrative Code
- `FM` = Trade-offs / Failure Modes
- `CONS` = Trade-offs / Consequences
- `TS` = Validation / Testing Strategy
- `MVV` = Implementation Plan / Minimum Viable Validation
- `IP` = Implementation Plan (phases, prerequisites, new dependencies)
- `PE` = Validation / Performance Expectations
- `CA` = Research Findings / Critical Assumptions

---

## Atom grammar — what a guard is

- [REQ-1] "A guard predicate MUST be a symbolic atom over a declared tag, not a host language callback and not a free-form expression string." — (NC, `0003:C1`)

- [REQ-2] "The initial operator vocabulary MUST be closed and typed: equality, membership, bounded integer comparison, existence, and set containment. Unknown operators MUST be rejected during parse or lint before resolution." — (NC, `0003:C7`)

- [REQ-3] "Each operator MUST declare which tag value kinds it accepts. A predicate whose literal cannot be parsed as the declared tag kind MUST be rejected before resolution." — (NC, `0003:C8`)

- [REQ-4] The operator/kind matrix fixes acceptance per operator: `eq` over `enum`, `bool`, `int`, `scalar` with "one typed scalar"; `in` over the same kinds with a "non-empty typed scalar set"; "`lt`, `lte`, `gt`, `gte`" over `int` with "one typed integer"; `exists` over "any kind, provided the tag is declared optional" with a boolean literal; `contains` over "`set` with a declared element universe" with a "non-empty typed element set". — (MTX)

- [REQ-5] "a key declared always-present contributes no `{absent}` assignment, so an `exists` atom over it is well-formed but vacuous — lint reports it as such rather than rejecting it" — (MTX, `exists` row)

- [REQ-6] "`Key` must resolve to a declared tag, `Operator` must be allowed by the operator/kind matrix above, and `Literal` must parse to the operator's literal shape." — (TD)

- [REQ-7] "**Source identity is not an atom field** — it is carried by the enclosing normalized row (`RuleID`, `SourceLocator`, per RDR 0002 and `internal/resolve/resolve.go::Row`), which is why this RDR's atom identity tuple is the row's identity joined to the atom's own four fields rather than an identity stored on the atom." — (TD)

- [REQ-8] "a guard atom is identified by the total tuple `(RuleID, SourceLocator, key, block, operator, literal)`, never by its index within `all` or `unless`." — (LBD, Identity)

- [REQ-9] "Semantic equality is `(tag, operator, literal)`." — (LBD, Identity)

- [REQ-10] "**Which equality each operation uses is fixed, not left to the implementer**: diagnostics, deduplication, and any \"same atom\" claim use the **identity tuple**; only domain computation — deciding what subset of the product an atom denotes — uses **semantic equality**" — (LBD, Identity)

- [REQ-11] "Overlap detection is therefore *not* an equality test at all: it intersects the rows' accepted assignment sets, so two byte-identical guards in two different rule ids correctly produce an overlap finding naming both rule ids, rather than being deduplicated into one row." — (LBD, Identity)

- [REQ-12] "equality compares a tag value to one typed literal; membership checks a scalar tag against a typed literal set; bounded integer comparison uses `lt`, `lte`, `gt`, and `gte`; existence checks presence of an optional tag value; set containment checks declared set-valued tags against a typed element set." — (LBD, Operator semantics)

## Tag declaration model (this RDR owns it)

- [REQ-13] "This RDR owns the **tag declaration model** — the typed alphabet every guard atom is written against. A tag declaration MUST carry a value kind, and — as its kind admits, per the agreement clause below — MAY carry a finite domain, an optionality marker, a single-valued marker, and (for set-valued kinds) an element universe." — (NC, `0003:C2`)

- [REQ-14] "The value kinds are exactly five, spelled with these tokens wherever a kind is named: in a declaration, in the operator/kind matrix, and in a diagnostic — `enum`, `bool`, `int`, `set`, and `scalar`." — (NC, `0003:C2`)

- [REQ-15] "The `scalar` kind is the opaque scalar: a typed value compared only by equality and membership. It carries **no** finite domain and can never bear an exhaustiveness claim, so a guard dimension over a `scalar` takes the blocking inability-to-prove outcome, and a `scalar` declaration carrying a finite domain, an element universe, or a single-valued marker is a declaration error." — (NC, `0003:C2`)

- [REQ-16] "RDR 0002 owns where a declaration is authored and how it is carried through normalization; this RDR owns what a declaration means. Neither document restates the other (JDR 0001 P6)." — (NC, `0003:C2`)

- [REQ-17] "A finite domain MUST be declarable for any kind an exhaustiveness claim can range over: an `enum` declares its value set, a `bool` is finite by construction, an `int` declares a `{min..max}` bound … and a `set` declares the element universe its members are drawn from." — (NC, `0003:C3`)

- [REQ-18] "this RDR fixes their meaning — both endpoints are **inclusive**, so `{0..3}` has cardinality 4" — (NC, `0003:C3`)

- [REQ-19] "A declaration carrying no finite domain is well-formed — the tag remains runtime-evaluable — but a guard dimension over it cannot carry an exhaustiveness claim, and lint MUST take the blocking inability-to-prove outcome for that dimension." — (NC, `0003:C3`)

- [REQ-20] "A tag declaration MUST be able to state whether the key may be absent. A declaration carrying **no optionality marker declares the key optional** — the conservative default" — (NC, `0003:C4`)

- [REQ-21] "A key declared always-present MUST NOT be absent from a conforming evaluation view; a key declared optional MAY be." — (NC, `0003:C4`)

- [REQ-22] "Presence is a declared property of the tag, not an observation of one view: lint decides `exists` projection and the \"can refuse\" narrowing from this declaration, never from a runtime trace." — (NC, `0003:C4`)

- [REQ-23] "An evaluation view **conforms** to the declared model when every always-present key is present in it and every single-valued tag holds at most one of its declared domain values. Conformance is the premise every lint claim in this RDR is conditional on: a green exhaustiveness result asserts coverage over conforming views only." — (NC, `0003:C4`)

- [REQ-24] "this RDR MUST NOT be read as promising anything about a view that violates the declarations." — (NC, `0003:C4`)

- [REQ-25] "A tag declaration MUST be able to state that the tag is **single-valued**: at most one of its declared domain values holds in any conforming evaluation view. The marker is meaningful only for a kind carrying a finite domain, and a `set`-valued kind MUST NOT carry it." — (NC, `0003:C5`)

- [REQ-26] "a single-valued tag contributes **one dimension of `|domain|` assignments** (for enum `{a,b,c}`: three), because exactly one value holds per view. Absent the marker, a tag whose values could co-occur contributes one **independent boolean dimension per value** (`2^|domain|`, minus nothing — the model does not assume at least one holds)." — (NC, `0003:C5`)

- [REQ-27] "RDR 0006's grouping-dependent lint findings read this field from the declaration and MUST NOT infer it from a tag's name, its value spelling, or a fixture (JDR 0001 §JD-13)." — (NC, `0003:C5`)

- [REQ-28] "A declared domain MUST agree with its value kind. Each kind admits exactly these fields, and any other combination is a declaration error that MUST be rejected before normalization completes" — (NC, `0003:C6`)

- [REQ-29] The kind/field agreement table: `enum` → value set, no element universe, single-valued allowed; `bool` → finite by construction, no element universe, single-valued allowed; `int` → `{min..max}`, no element universe, single-valued allowed; `set` → element universe required to claim exhaustiveness, "**no** — values co-occur by construction" for the marker; `scalar` → "**no**" domain, no element universe, "**no**" marker. So "an `{min..max}` bound on an `enum`, an element universe on any kind but `set`, or a single-valued marker on a `set` or a `scalar`, is each a declaration error." — (NC, `0003:C6`)

- [REQ-30] "A value appearing in a guard literal that lies outside its tag's declared domain MUST likewise be rejected — this is the predicate semantic kind **literal-outside-declared-domain**, distinct from a literal parse failure (the literal parses fine; it is simply not in the domain) — so an unsatisfiable atom is a load-time error rather than a silently-never-matching row." — (NC, `0003:C6`)

- [REQ-31] "a malformed *declaration* (domain disagreeing with its kind) is rejected by the declaration loader before normalization completes … a *literal-outside-domain* is rejected by guard parsing after the declarations have loaded and before rows are yielded." — (NC, `0003:C6`)

- [REQ-32] "Both are this RDR's rejection rules carried by RDR 0002's load categories (JDR 0001 §D7(iii)): the malformed declaration is a `malformed tag declaration`, the out-of-domain literal a `malformed predicate atom`. RDR 0006 mints nothing for either; RDR 0005 maps both onto the envelope under JDR 0001 §JD-8." — (NC, `0003:C6`)

- [REQ-33] "a tag declaration carries five fields: value kind, finite domain, optionality, single-valuedness, and (for set kinds) an element universe. … any enumeration of this model in this document states all five." — (LBD, Declaration model)

## Evaluator scope and the kernel split

- [REQ-34] "This RDR's guard evaluator MUST decide value semantics over a present value only. Key presence, existence atoms, absent-key unevaluability, and the combination of per-atom verdicts belong to the kernel (JDR 0001 §D4, normative in RDR 0007); the evaluator MUST NOT read the tag view." — (NC, `0003:C9`)

- [REQ-35] "an existence atom over an absent key is decided — not unevaluable" and "a value atom over an absent key is unevaluable rather than false" — (TD)

- [REQ-36] "an author expressing \"this row applies when tag X is absent\" writes an existence atom, not a value atom that happens to miss." — (TD)

- [REQ-37] Per the `authority` census, the canonical writer of a value-atom verdict over a present value is "This RDR's evaluator"; key presence, the existence-atom verdict, and per-atom verdict combination are the kernel's, and this RDR's evaluator is "**explicitly not an arm**; it never reads the tag view". — (AUTH)

## Polarity, ordering, and `unless`

- [REQ-38] "Positive guard atoms MUST live in `all`; negative guard atoms MUST live in `unless`. Successful row matching MUST NOT depend on source order or first-match priority." — (NC, `0003:C10`)

- [REQ-39] "a row's accepted assignments are the intersection of all positive `all` atom domains minus the single conjunctive assignment set matched by the row's full `unless` block. `unless` is not per-atom negation, and it does not create source-order priority." — (TD)

- [REQ-40] "a row qualifies only when every `all` atom is true and the `unless` predicate set is not fully true; if multiple rows qualify, RDR 0001's exact-one resolver refuses instead of choosing by priority." — (LBD, Selection / predicate)

- [REQ-41] "**`unless` is subtracted two-valued only when its atoms are decided.** The excluded-intersection subtraction this RDR states for lint is a set operation over decided atoms, while the runtime computes the row verdict in three-valued Kleene — RDR 0007 fixes the verdict as `all_result ∧ ¬(unless_conj)` with `¬U = U`, so an **unevaluable atom inside `unless` makes the whole row unevaluable**, not merely un-excluded." — (NC, `0003:C17`)

- [REQ-42] "Lint MUST therefore treat a value atom in an `unless` block over a key declared optional exactly as it treats one in `all`: the row **can refuse**, and the group's claim is withheld under the narrowing above." — (NC, `0003:C17`)

- [REQ-43] "**A can-refuse row contributes no assignments to the coverage union**, and the group it sits in has no provable product. … lint MUST NOT credit it with its `all`-intersection unsubtracted" — (NC, `0003:C17`)

- [REQ-44] "a withheld group emits the withholding finding for each refusing row, and MUST NOT additionally emit a `graph-coverage-gap` naming a witness assignment that a refusing row would in fact accept when its key is present. Overlap findings among the group's decidable rows are unaffected and still MUST be emitted." — (NC, `0003:C17`)

## Row groups, participation, and the scoped product

- [REQ-45] "Coverage and overlap checks MUST be scoped to a normalized row group supplied by the transition/lint model, and MUST evaluate the participating guard dimensions as one product rather than as independent one-dimensional checks." — (NC, `0003:C13`)

- [REQ-46] "**The row group is defined here**, not deferred: a scoped row group is the set of normalized candidate rows sharing one selection context — the same source state and the same recognized outcome — since that is exactly the set RDR 0001 resolves exact-one over." — (TD)

- [REQ-47] "RDR 0006 supplies the graph traversal that enumerates which selection contexts are reachable; it does not define the grouping predicate, and this RDR does not read one back from it." — (TD)

- [REQ-48] "A dimension **participates** in a row group when any row in that group carries a **guard** atom over that key, in either `all` or `unless` — not only when the rows constrain it differently. The authored block is what makes an atom a guard atom (JDR 0001 §D6): every atom under `guard.all` or `guard.unless` is a guard atom regardless of operator, `eq` included." — (NC, `0003:C13`)

- [REQ-49] "A key every row constrains identically still bounds the product and still requires a finite declared domain; it MUST NOT be dropped from the product because it does not discriminate." — (NC, `0003:C13`)

- [REQ-50] "**Match keys are not product dimensions.** A row's match pattern selects which group the row belongs to — the selection context this RDR groups by — and is a separate field from its guard" — (NC, `0003:C13`)

- [REQ-51] "**The match/guard split is per atom per group, not per key per model.** The same tag key MAY be a match key in one rule and a guard key in another … a key enters that group's product when some row in *that* group carries a guard atom over it, regardless of how the key is used in any other group." — (NC, `0003:C13`)

- [REQ-52] "An `exists` atom projects onto the scoped product as a **per-key presence dimension**: a two-valued dimension `{present, absent}` for that key, alongside the key's value dimension. An `exists` atom denotes `{present}` or `{absent}` on it" — (NC, `0003:C14`)

- [REQ-53] "a value atom denotes a subset of the key's value dimension **as the operator/kind agreement clause projects it** — one held value narrowed for `eq`/`in`/comparisons, which therefore requires the key be declared single-valued, or a containing subset for `contains` — and, because a value atom over an absent key is unevaluable rather than false, implicitly `{present}`." — (NC, `0003:C14`)

- [REQ-54] "`exists` is unaffected by that requirement: it reads presence alone, so it projects over a key whose values co-occur exactly as it does over a single-valued one." — (NC, `0003:C14`)

- [REQ-55] "A key declared always-present contributes no presence dimension — its `{absent}` assignment is not in the product" — (NC, `0003:C14`)

- [REQ-56] "Lint MUST NOT drop `exists` atoms from the product: a group carrying one stays provable, and certifying it exhaustive while ignoring the presence dimension is the false-green this RDR's narrowing forbids." — (NC, `0003:C14`)

- [REQ-57] "**`eq`, `in`, and the integer comparisons are single-value operators** … so an atom denotes the assignments in which the tag's single held value is the literal (`eq`), is a member of the literal set (`in`), or satisfies the comparison. `contains` is the operator for a tag whose values co-occur, and it projects the other way: the assignments whose held set contains every listed element." — (NC, `0003:C8`)

- [REQ-58] "a value atom over a dimension the model does **not** declare single-valued is not a differently-projecting atom — it is one lint **cannot project at all** … Such an atom MUST take the blocking inability-to-prove outcome for that dimension, exactly as a dimension with no finite declared domain does." — (NC, `0003:C8`)

- [REQ-59] "Lint MUST NOT silently pick a reading — resolving it as \"the literal is among the held values\", \"the held set equals the literal\", or \"the held set is contained in the literal\" yields different union cardinalities and different overlap verdicts on the same model, which the published-bound clause forbids." — (NC, `0003:C8`)

- [REQ-60] "Coverage is `union(row_i accepted assignments) == scoped product` for that row group, and overlap is any non-empty `row_i accepted assignments intersect row_j accepted assignments`." — (CA, A2 Evidence — the derivation the contracts implement)

## Exhaustiveness eligibility, default-on, no opt-out

- [REQ-61] "Lint MAY claim guard exhaustiveness only for finite declared domains: enum values, booleans, declared set element universes, or bounded integer ranges." — (NC, `0003:C11`)

- [REQ-62] "If a guard dimension lacks a finite declared domain, lint MUST refuse or downgrade an exhaustiveness claim for that dimension rather than treating the covered examples as complete." — (NC, `0003:C12`)

- [REQ-63] "The exhaustiveness claim is **default-on for every scoped row group whose participating dimensions are all finitely declared**: an opt-in flag would let the guarantee be silently skipped exactly where it matters." — (TD)

- [REQ-64] "**Leaving a dimension undeclared is not an opt-out.** An undeclared dimension does not remove the row group from proof; it makes the proof unavailable, which is the Loud blocking inability-to-prove outcome … never a silent downgrade to unchecked. No authoring gesture quietly exempts a group: a group is either proved, or it carries a blocking finding naming the dimension that cannot be proved." — (TD)

- [REQ-65] "\"Refuse\" and \"downgrade\" are one outcome, not an author's choice: every case this RDR sends to refuse-or-downgrade — a non-finite dimension, a finite product too large to prove deterministically, and a withheld claim under the narrowing above — MUST produce a blocking inability-to-prove finding. This RDR MUST NOT mint a non-blocking warning category for these cases" — (NC, `0003:C20`)

- [REQ-66] "\"Downgrade\" therefore names the same blocking outcome as \"refuse\" on both sides of the seam, and never a silent or advisory one." — (NC, `0003:C20`)

## The runtime-veto narrowing (this RDR is the recording document)

- [REQ-67] "An exhaustiveness claim MUST NOT be stronger than the runtime it describes: lint MUST NOT certify a row group exhaustive when a participating row can refuse `guard_unevaluable` under RDR 0007's aggregation veto. Where the two disagree the lint promise narrows; the runtime veto MUST NOT be weakened." — (NC, `0003:C16`)

- [REQ-68] "**\"Participating row\" here means every row in the group, escape rows included** — the participation clause's population, not the ordinary-row population the overlap check uses." — (NC, `0003:C16`)

- [REQ-69] "Withholding is therefore decided over the whole group, while overlap remains split into the two populations above … The two-population reading applies to overlap only; reading it into the narrowing would certify green exactly the group the runtime refuses." — (NC, `0003:C16`)

- [REQ-70] "A row \"can refuse\" when the row carries a value atom over a key **declared optional** **in either block** — `all` or `unless` … the same one declared field the optionality clause defines, not a second presence property and not a graph query." — (NC, `0003:C21`)

- [REQ-71] "Lint MUST decide this syntactically over declarations so the test is total; it MUST NOT withhold a claim merely because some assignment in the product is unreached." — (NC, `0003:C21`)

- [REQ-72] "An **owned** tag carries a second, graph-level presence condition … but that condition governs whether a row may *match*, not whether its guard can refuse … Presence for the withholding decision reads exactly one field" — (NC, `0003:C21`)

- [REQ-73] "A withheld exhaustiveness claim MUST be observable, not silent. It takes the same blocking inability-to-prove form an unprovable dimension already takes — RDR 0006's `graph-unprovable-coverage` — and MUST name the participating row and the atom that can refuse, using the source rule/context id every other predicate diagnostic names." — (NC, `0003:C18`)

- [REQ-74] "Emitting nothing MUST NOT satisfy this clause: an exit code alone cannot distinguish a withheld claim from a proved one." — (NC, `0003:C18`)

## Escape rows

- [REQ-75] "A declared escape row participates in the coverage identity like any other row: its accepted assignments are computed from its guard atoms and unioned with its peers'. An escape row carrying no guard atoms denotes the whole scoped product and therefore closes coverage by itself. Lint MUST NOT treat \"an escape row exists\" as a separate coverage-satisfying fact outside the union." — (NC, `0003:C15`)

- [REQ-76] "**An escape row closes coverage only for the failure classes it can actually rescue** — those it declares in its rescue list, and of those only `no_match` and `ambiguous_match`. … An escape row therefore cannot rescue a class it does not declare, nor either blocking class, however bare its guard." — (NC, `0003:C15`)

- [REQ-77] "The coverage union is consequently computed per (scoped row group × declared rescuable class), and a row declaring one class MUST NOT close the group's other arm" — (NC, `0003:C15`)

- [REQ-78] "a bare escape row MUST NOT be read as discharging the narrowing: a group whose ordinary row can refuse `guard_unevaluable` still has its claim withheld, because at runtime that refusal is returned before the escape row is ever consulted." — (NC, `0003:C15`)

- [REQ-79] "Overlap is checked in **two separate populations** … an escape row overlapping a *guarded* row is **not** a runtime ambiguity and MUST NOT be reported as one. What lint MUST still check is overlap **among escape rows for the same failure class** … Escape rows are therefore excluded from the ordinary-row overlap check and subjected to their own; excluding them from **coverage** is what MUST NOT happen." — (NC, `0003:C15`)

- [REQ-80] "lint MUST partition the escape rows into **one population per declared failure class**, place a row declaring several classes in **each** of those populations, and run the pairwise overlap check within each population independently." — (NC, `0003:C15`)

- [REQ-81] "A pair overlapping in more than one class is one finding per class, naming the class, so the author can see which rescue path is disabled." — (NC, `0003:C15`)

- [REQ-82] "**A bare escape row is not a silent opt-out from the guarantee.** … Lint MUST therefore report a bare escape row — one carrying no guard atoms — that closes a group's coverage as an **observable** result: the group's verdict names the escape row that closed it … Emitting a bare green for such a group MUST NOT satisfy this clause." — (NC, `0003:C15`)

## Set literals and set-valued proof

- [REQ-83] "A set literal — the right-hand side of `in` and of `contains` — has one canonical spelling: an unordered set of typed elements, duplicate-free, and compared as a set. Two authored spellings differing only in element order or in repeated elements MUST parse to the same literal and therefore to the same atom under the identity tuple; a repeated element MUST be rejected at parse rather than silently collapsed." — (NC, `0003:C22`)

- [REQ-84] "Implementations MUST canonicalize before the literal enters the identity tuple, so reordering a set literal cannot change a diagnostic — the source-order independence MVV Scenario 6 requires." — (NC, `0003:C22`)

- [REQ-85] "Set-valued guard domains MUST be proved with a deterministic symbolic or bitset-equivalent representation. If the finite product is too large for that proof, lint MUST refuse or downgrade the exhaustiveness claim rather than silently capping enumeration." — (NC, `0003:C23`)

## The published cardinality bound

- [REQ-86] "\"Too large to prove\" MUST be a declared, model-independent bound, not an implementation's incidental limit: the implementation MUST publish the bound it enforces, and the same model MUST receive the same verdict on every conforming implementation. A bound discovered by exhausting memory or wall-clock is not a conforming bound." — (NC, `0003:C24`)

- [REQ-87] "The refusal diagnostic MUST report the product size it computed and the bound it exceeded, so an author can tell an over-large product from an undeclared dimension." — (NC, `0003:C24`)

- [REQ-88] "The quantity both figures report is the **cardinality of the scoped product** — the number of assignments in it, the product of every participating dimension's **assignment count** — not a bitset width, byte size, or row count." — (NC, `0003:C24`)

- [REQ-89] The per-kind assignment-count table: "`enum`, `int` (single-valued)" → `|domain|`; "`bool` (single-valued)" → 2; "`enum`, `bool`, `int` **without** the single-valued marker" → "`2^|domain|` — one independent boolean dimension per value"; "`set`" → "`2^|element universe|` … **never** `|universe|`"; "any optional key" → "multiplied by 2 for its `{present, absent}` presence dimension". — (NC, `0003:C24`)

- [REQ-90] "**The table reads the declaration's defaults, not an author's intent.** … An implementation MUST apply these defaults when computing the cardinality — reading an unmarked declaration as single-valued or always-present is the same exponential understatement the `set` row forbids" — (NC, `0003:C24`)

- [REQ-91] "The bound is a single integer constant published by the implementation and reported in the diagnostic beside the computed cardinality; it is not per-model, per-group, or configurable per run" — (NC, `0003:C24`)

- [REQ-92] "Two conforming implementations MAY publish different bounds, but each MUST return the same verdict for the same model and MUST report which bound it applied." — (NC, `0003:C24`)

- [REQ-93] "**A dimension with no finite declared domain has no assignment count**, so a product containing one has no cardinality … The too-large comparison is therefore defined **only over a fully-provable product** — lint MUST compute the cardinality and test the bound after every participating dimension is known finite, and MUST NOT report a computed size for a product carrying an unprovable dimension." — (NC, `0003:C24`)

- [REQ-94] "each unprovable dimension draws its own blocking finding naming that dimension, and the over-large refusal — the one that carries `(computed size, bound)` — is simply not among the findings for such a group" — (NC, `0003:C24`)

## Diagnostics and report-every-defect

- [REQ-95] "Lint MUST report every defect it can decide in one pass over a row group, not the first one it encounters. Withholding a group's exhaustiveness claim MUST NOT suppress overlap findings, coverage findings, or further withholding reasons for that same group: each unprovable dimension, each refusing row, each overlapping pair, and any coverage gap over a provable product is its own finding." — (NC, `0003:C19`)

- [REQ-96] "A group with two refusing rows MUST emit a finding for each, so the emitted set does not depend on row or dimension iteration order — the same source-order independence matching already requires." — (NC, `0003:C19`)

- [REQ-97] "A withheld claim and a coverage or overlap finding MAY be emitted together; what MUST NOT be emitted is a green exhaustiveness result alongside any withholding reason." — (NC, `0003:C19`)

- [REQ-98] "Overlap and coverage diagnostics MUST name the source rule id or context id that contributed each predicate involved in the finding." — (NC, `0003:C25`)

- [REQ-99] "A coverage gap has no contributing predicate — it is an absence — so it is attributed differently: a gap finding MUST name the selection context that scopes the group, every source rule id in that group, and at least one concrete uncovered assignment from the scoped product … Naming the group without a witness assignment MUST NOT satisfy this clause." — (NC, `0003:C25`)

## Provenance and owned tags

- [REQ-100] "Predicate lint MUST distinguish owned, observed, and recognized tags. A row that matches an owned tag MUST be rejected unless every reachable predecessor sets or preserves that tag before the match." — (NC, `0003:C26`)

- [REQ-101] "The reachable-predecessor relation is RDR 0006's owned-state reachability relation … this RDR cites it and defines no second one (A12)." — (NC, `0003:C26`)

- [REQ-102] "recognized tags are fresh event inputs, observed tags are re-read before matching, and owned tags must have a reachable predecessor write before a row may match them." — (TD)

## Disposition table (input class × outcome)

- [REQ-103] "Value atom **inside `unless`** over an absent key" → runtime "`guard_unevaluable` refusal — `¬U = U`, so the row is unevaluable, not merely un-excluded"; lint "Exhaustiveness claim **withheld** for that group; the block is **not** subtracted as if decided"; Loud on both surfaces. — (DISP)

- [REQ-104] "Value atom over an absent key" → runtime `guard_unevaluable`; lint "Exhaustiveness claim **withheld** for that group"; diagnostic "the blocking inability-to-prove finding, naming the row and the refusing atom"; Loud. — (DISP)

- [REQ-105] "Existence atom over an absent key" → "Decided (`presence == literal`) — never unevaluable"; lint "Selects `{absent}` on the presence dimension (A7)"; no diagnostic. — (DISP)

- [REQ-106] "Zero rows qualify" → lint "Coverage gap if the product is provable", `graph-coverage-gap` "naming the selection context, every rule id in the group, and one uncovered assignment". "Two or more rows qualify" → `graph-overlap` "naming both source rule ids"; "RDR 0001 refuses — never first-match". — (DISP)

- [REQ-107] "Guard dimension lacks a finite declared domain" → "Evaluable at runtime"; lint "**Refuse or downgrade** the group's claim, one finding per unprovable dimension — never treat examples as complete, never stop at the first"; `graph-unprovable-coverage` naming the dimension. — (DISP)

- [REQ-108] "Finite product too large to prove deterministically" → lint "**Refuse or downgrade** — never silently cap enumeration"; `graph-product-too-large` "reporting the computed product cardinality and the published bound". — (DISP)

- [REQ-109] "Two escape rows matching one failure class" → "Overlap finding within the escape population", `graph-overlap` "naming both escape rule ids and the failure class". "Escape row overlapping a guarded row" → "**No overlap finding**; the escape row still contributes its assignments to the coverage union"; Silent by design. — (DISP)

- [REQ-110] "Row group is domain-exhaustive but a participating row can refuse" → "Claim withheld — the narrowing; one finding per refusing row, never just the first"; "RDR 0007 payload at runtime **and** a blocking lint finding naming the refusing atom — absence of a green result is not the artifact". — (DISP)

- [REQ-111] "`all` atom decides false" → "Row pruned … Silent by design — a decided false is not a defect"; "Full `unless` block decides true" → "Row disabled … Excluded intersection subtracted … Silent by design". — (DISP)

## Failure modes and non-goals

- [REQ-112] "Visible failures should be typed load or lint failures: unknown operator, operator/tag-kind mismatch, literal parse failure, unknown tag, and the two declaration/literal-domain rejections (RDR 0002 load categories, §D7(iii)) at load; non-exhaustive finite domain, overlapping candidate rows, guard dimension not provable because it lacks a finite domain, or finite scoped product too large to prove at lint." — (FM)

- [REQ-113] "A guard that cannot be decided at runtime surfaces as RDR 0007's `guard_unevaluable` refusal, not as a predicate parse kind owned here." — (FM)

- [REQ-114] "Silent failure would be a false exhaustiveness claim; the recovery path is to keep every exactness claim tied to A2 and the MVV fixture." — (FM)

- [REQ-115] "This RDR introduces no encode/decode pair. Parse/render fidelity for the sparse TOML source belongs to RDR 0002." — (RT) — **no Round-Trip / Inverse Invariant is declared.**

- [REQ-116] "There is no embedded host predicate and no free-form expression grammar." — (AP)

- [REQ-117] "No new third-party dependency is proposed. The predicate grammar and finite domain checks should be implemented with local Go code unless Resolve proves a small parsing or set library is necessary." — (IP, New Dependencies)

- [REQ-118] "This RDR creates no persistent resource." — (IP, Day 2 Operations)

- [REQ-119] "no callback invocation, expression parser, or external engine is part of the hot path. If implementation later indexes predicates for speed, the optimization must preserve the normalized atom semantics, scoped product proof, and exact-one refusal behavior." — (PE)

- [REQ-120] "the canonical name is \"guard predicate\"; rejected alternatives are \"condition callback\" and \"guard expression\" because both invite opaque host logic." — (LBD, Naming)

- [REQ-121] "RDR 0002 owns the TOML container; this RDR owns the guard atom grammar embedded in that container." — (LBD, Wire / byte format)

## Implementation phases

- [REQ-122] "Define the tag-kind/operator compatibility matrix and normalized predicate atom shape used by resolver and lint." — (IP, Phase 1)

- [REQ-123] "Define how enum, boolean, set-universe, and bounded-int domains are converted into scoped row-group coverage and overlap checks, including refusal/downgrade behavior for unbounded dimensions and finite products too large to prove deterministically." — (IP, Phase 2)

- [REQ-124] "Build the MVV fixture against representative RDR and kata guards and use the result to confirm or adjust the initial operator set. … the implementation MVV must add at least one `contains` predicate over a declared set-valued tag before the full operator vocabulary is accepted." — (IP, Phase 3)

- [REQ-125] "Connect predicate diagnostics to RDR 0002 source identities, RDR 0001 exact-one selection, RDR 0005 CLI output, and RDR 0006 lint authority." — (IP, Phase 4)

- [REQ-126] "Phase 3 must therefore declare the markers **and** close the routing, or author a separate group for Scenario 2's complete partition — reusing it unchanged makes Scenario 2 fail on first run against a correct lint" — (TS)

- [REQ-127] Phase 3's fixture must "declare `single_valued` on `profile`, `prelock_iterations`, `cluster_eligible`, `stage`, `lens` and `rewind_target`" and, "If any target-flow dimension turns out to need co-occurring values, it is a `set` tag under `contains`, not an unmarked enum — record which, since that changes RDR 0002's authoring request." — (CA, A21 Plan)

## Testing strategy scenarios

- [REQ-128] "**Scenario**: Evaluate representative RDR and kata rows that use equality, membership, set containment, bounded integer comparison, existence, and mixed `all`/`unless` guards. **Expected**: Exactly one qualifying row resolves for the legal input; a row whose `all` predicates match is disabled when its full conjunctive `unless` block also matches; zero or multiple qualifying rows become typed refusals." — (TS, 1)

- [REQ-129] "**Scenario**: Lint finite enum, boolean, set-universe, and bounded-int tag domains inside one normalized row group with one complete partition, one intentional gap visible only in the multi-dimensional product, and one intentional overlap visible only in the multi-dimensional product. **Expected**: Complete partitions pass; product-level gaps and overlaps fail with source rule/context ids. The gap and the overlap are both reported from one run over the group — detecting only the first is a failure. The gap finding names the selection context, every rule id in the group, and one concrete uncovered assignment; the overlap names both contributing rule ids." — (TS, 2)

- [REQ-130] "**Scenario**: Lint an otherwise valid guard over an unbounded integer or undeclared finite domain, and lint a declared finite product too large for the deterministic proof representation. Include a group with two separately unprovable dimensions, and two equal-cardinality products of differing shape." — (TS, 3)

- [REQ-131] "**The shape pair is constructed relative to the implementation's published bound B, not to an absolute size**: build two products of equal cardinality C with C > B — one from few wide dimensions (e.g. two dimensions of ~sqrt(C) values each), one from many narrow ones (e.g. log2(C) boolean dimensions) — so both must refuse if cardinality alone gates the verdict. Repeat the pair just under B, where both must prove." — (TS, 3)

- [REQ-132] "**Expected**: Runtime evaluation remains available, but lint refuses or downgrades the exhaustiveness claim for that dimension/product. The two-dimension group emits two findings, not one. The over-large refusal reports both the computed product cardinality and the published bound. The two equal-cardinality products receive the same verdict — a divergence refutes A15 and the bound clause is restated (a spec defect routed back to this RDR)." — (TS, 3)

- [REQ-133] "for each shape, record whether the proof representation actually completes within the implementation's own resource budget, and compare that to the verdict the bound gave. A15 is refuted when the bound says provable and the representation does not complete, or vice versa" — (TS, 3)

- [REQ-134] "**Scenario**: Parse malformed guard atoms: unknown tag, unknown operator, unsupported operator/tag-kind pair, literal parse mismatch, and a **literal outside its tag's declared domain** … Include the malformed *declarations* the domain/kind agreement clause rejects: a `{min..max}` bound on an `enum`, an element universe on a scalar kind, and a single-valued marker on a `set` kind or on a kind carrying no finite domain." — (TS, 4)

- [REQ-135] "**Expected**: Each failure is rejected before resolution with a predicate semantic kind … The declaration errors are rejected before normalization completes, so a consumer reading the model — including RDR 0006's `graph-single-valued-state` — never sees a marker its kind cannot carry." — (TS, 4)

- [REQ-136] "**Scenario (single-valued acceptance)**: Lint one row group over a finite-domain tag **with** and **without** the single-valued marker, holding every other input fixed. **Expected**: the scoped product differs by construction … so a row set that closes coverage under the marker leaves an uncovered assignment without it. The two runs must therefore return different verdicts on the same rows; identical verdicts mean the marker is not reaching the product and `graph-single-valued-state` has no effective producer. Assert the marker is read from the declaration and never inferred from the tag's name or value spelling." — (TS, 5)

- [REQ-137] "**Scenario**: Reorder authored rows and guard atoms without changing their semantics, including reordering the elements inside an `in` set literal. **Expected**: Successful matching and lint findings are unchanged because source order is not a selection mechanism; an ambiguous pair remains a multiple-match refusal instead of becoming a first-match success." — (TS, 6)

- [REQ-138] "**Scenario (block retention)**: Normalize a row carrying atoms over the **same key in both blocks** — one in `all`, one in `unless` — and inspect the normalized atoms. **Expected**: each atom still reports the block it was authored in (`Block == all` / `Block == unless`), and the two remain distinguishable and separately identifiable under the identity tuple." — (TS, 7)

- [REQ-139] "**Scenario**: Evaluate a row group whose guards are exhaustive over their declared domains but where one participating row carries a value atom over a key that can be absent. **Expected**: The value atom is unevaluable rather than false, resolution refuses `guard_unevaluable` under RDR 0007's veto, and lint emits the blocking inability-to-prove finding naming that row and atom — not merely an absent green result … Paired negative control: an otherwise identical group over always-present keys certifies green. An existence atom over the same absent key decides instead of refusing, and this RDR's evaluator is never consulted for either." — (TS, 8)

- [REQ-140] "**Both blocks, and the escape path.** Run the same group a second time with the optional-key value atom moved from `all` into `unless`: the verdict must be identical — claim withheld … Then add a bare escape row to the group: the claim must **still** be withheld, since `resolve.go::Resolve` returns the gate's refusal before `escapeOrRefuse` is reachable, so the escape row cannot rescue `guard_unevaluable`. A run that goes green once the escape row is present is the false-green A19 guards." — (TS, 8)

- [REQ-141] "The MVV should become production tests that exercise both runtime predicate evaluation and lint-time finite-domain reasoning. Done means the same normalized predicate atoms drive exact-one row selection, overlap detection, and exhaustiveness proof/refusal without host callbacks or source-order priority." — (TS)

---

## REQ-MVV

- [REQ-MVV] "Encode one RDR flow slice and one kata flow slice as normalized candidate rows, including equality, enum membership, set containment, bounded integer comparison, existence, and mixed `all`/`unless` guards. Lint must prove one exhaustive and mutually exclusive scoped row group, then detect one intentional gap and one intentional overlap that only appear in the multi-dimensional product, with source rule/context ids in the diagnostic. The fixture must also include one row group that is domain-exhaustive yet contains a possibly-absent guard key, and assert **positively**: lint emits the blocking inability-to-prove finding naming that row and the refusing atom. Asserting only the absence of a green result does not discharge this — a run that emits nothing must fail the test. The MVV must also carry a **negative control**: a row group whose keys are all declared always-present, which lint certifies green, proving the narrowing is tight rather than blanket." — (MVV, `0003:MVV`)

REQ-MVV is the gating validation. It decomposes into five obligations that must
all be present in one production test suite (Phase 3, `Target-Flow Fixture`):

1. a fixture encoding one RDR flow slice and one kata flow slice as normalized
   candidate rows, covering `eq`, `in`, `contains`, bounded integer comparison,
   `exists`, and mixed `all`/`unless` (REQ-128, REQ-124);
2. one scoped row group proved exhaustive *and* mutually exclusive — the MVV's
   only passing exhaustiveness verdict (REQ-129, REQ-60);
3. one intentional gap and one intentional overlap that appear only in the
   multi-dimensional product, both reported from a single run, with source
   rule/context ids and a concrete uncovered assignment (REQ-95, REQ-98,
   REQ-99);
4. a **positive** assertion on the narrowing: a domain-exhaustive group carrying
   a possibly-absent guard key emits the blocking inability-to-prove finding
   naming the row and the refusing atom. "A run that emits nothing must fail
   the test" — asserting only the absence of a green result does not discharge
   it (REQ-73, REQ-74, REQ-110);
5. a **negative control** — an otherwise identical group over always-present
   keys that lint certifies green, proving the narrowing is tight rather than
   blanket (REQ-55, REQ-139).

Authorability is closed by the record: every case is writable against this
RDR's own declaration model, and the `single_valued` marker the positive half
depends on is writable because RDR 0002's schema carries it (A16, `Verified`).

The remaining dependency is **sequencing, not authorability**: "RDR 0007's
kernel reshape must land before Phase 1 builds against the atom-slice shape"
(MVV, Authorability), and the shipped kernel still carries `Row.Guard` as a
`string`.

---

## ASSUMPTIONS

- ASSUMPTION: the wire spellings of the five declaration fields are RDR 0002's
  `kind`, `domain`, `min`, `max`, `elements`, `single_valued`, and `required`,
  and this RDR's REQs are tested against those keys rather than inventing new
  ones. Grounded: `0003:C3` says `{min..max}` is "notation, not wire spelling:
  RDR 0002 spells the authored form as the two wire keys `min` and `max`";
  the Illustrative Code block writes `kind`/`domain`/`single_valued`/`required`
  and says "The key spellings above are RDR 0002's wire keys (JDR 0001
  §D7(iii); A16)"; A16 quotes RDR 0002's normative type-model clause listing
  all seven keys.

- ASSUMPTION: "optionality marker" (REQ-20) is the `required` key, where
  `required = true` means always-present and an omitted `required` means
  optional. Grounded: the Illustrative Code comment `required = true #
  always-present; contributes no {absent} assignment`, and the following line
  "`required` defaults to optional when omitted, as the optionality clause
  mandates."

- ASSUMPTION: this RDR ships no lint finding *codes* of its own — it states the
  semantics and RDR 0006 mints the codes (`graph-unprovable-coverage`,
  `graph-product-too-large`, `graph-coverage-gap`, `graph-overlap`,
  `graph-coverage-closed-by-escape`, `graph-owned-before-write`,
  `graph-single-valued-state`). What this RDR owns is the *predicate semantic
  kind* (unknown operator, type mismatch, literal parse failure,
  literal-outside-declared-domain, declaration/kind disagreement, unsupported
  operator/tag-kind pairing). Grounded: A4's Evidence; REQ-32; REQ-65's "This
  RDR MUST NOT mint a non-blocking warning category"; the `authority` census
  row for the withheld-claim artifact ("RDR 0006 widens
  `graph-unprovable-coverage`'s trigger and mints no code").

- ASSUMPTION: "blocking inability-to-prove outcome" (REQ-15, REQ-19, REQ-58,
  REQ-65) names one outcome with one carrier, not a family — it is
  `graph-unprovable-coverage` for the non-finite dimension, the unprojectable
  atom and the withheld claim, and `graph-product-too-large` for the over-large
  product. Grounded: REQ-65/`0003:C20` enumerates exactly the three cases and
  names both codes as the carriers, both blocking.

- ASSUMPTION: the guard evaluator this RDR owns is a Go seam satisfying
  `internal/resolve`'s `GuardEvaluator`, reshaped per RDR 0007 to take a parsed
  atom rather than a `string`. Grounded: Research Findings — "The kernel calls
  the guard seam this RDR owns (`resolve.go::GuardEvaluator`) but supplies no
  implementation of it"; the Prerequisites item naming
  `resolve.go::GuardEvaluator.Evaluate(guard string, view TagSet)` as the
  pre-reshape form; the `authority` census row "Guard atom shape | Kernel seam
  (JDR 0001 §D1, normative in RDR 0007) … the shipped string form is
  superseded, not an arm". The cluster deviations file records this RDR
  "Instantiates `resolve.TestGuardEvaluatorContract` (0007 Phase 3)".

- ASSUMPTION: `scalar` accepts `eq` and `in` only. Grounded: the operator/kind
  matrix lists `enum, bool, int, scalar` for `eq` and `in` and omits `scalar`
  from `lt/lte/gt/gte` and `contains`; `0003:C2` says the `scalar` kind is
  "compared only by equality and membership". `exists` also applies to it via
  "any kind, provided the tag is declared optional".

- ASSUMPTION: "the tag remains runtime-evaluable" (REQ-19) means the guard
  evaluator still decides such an atom at runtime; only the *lint*
  exhaustiveness claim is withheld. Grounded: REQ-107's disposition row
  ("Evaluable at runtime" / lint refuses), REQ-62, and the `Consequences`
  bullet "Unbounded integer or free-form string dimensions cannot receive
  silent exhaustiveness claims".

- ASSUMPTION: the "selection context" a row group is keyed by is the pair
  (source state, recognized outcome), and a match key is any atom authored
  under `[rule.match.*]`. Grounded: REQ-46 ("the same source state and the same
  recognized outcome"); REQ-50 ("`[rule.match.*]` vs. `[rule.guard.*]`, routed
  by block per JDR 0001 §D6"); desk-trace step 3 uses `status eq "Draft"` under
  `[rule.match.status]` as the grouping key rather than a dimension.

- ASSUMPTION: the four REQ-derived `Pending` records do not gate the REQ set —
  A15 and A21 are discharged *by* the MVV (REQ-133, REQ-127), and A18/A20 are
  routed to JDR 0001 §JD-18, which is open. Where a REQ depends on A20's
  runtime atom carrier (the runtime half of REQ-73), the record itself books
  the shortfall: "Booked as **A20** (JDR 0001 §JD-18) rather than left standing
  as a MUST no producer can meet" (the `Carriers` note after `0003:C18`). So
  REQ-73's atom-naming obligation is testable on the **lint** surface today and
  its runtime half is a known-open dependency, not a REQ this stage can close.

- ASSUMPTION: the three cluster-gate deviation entries (D1 empty/omitted
  `unless` identity, D2 stale peer claims, D3 §D13 set-value encoding citation)
  are dispositioned by running their named checks, not by editing contracts.
  Grounded: `artifacts/deviations.md` — "Stage 8 opens this file; running each
  named check IS the entry's disposition. An entry escalates only if its check
  contradicts a contract." D1's check is REQ-129/MVV Scenario 2 over a group
  with no `unless` block; the identity (`unless = ∅` subtracts nothing) is
  RDR 0007's fenced clause, not restated here.

- ASSUMPTION: a set tag's value crosses the kernel seam as RDR 0002's canonical
  JSON array (members sorted, duplicate-free, compact) while this RDR owns the
  set *literal* spelling (REQ-83) and the element universe. Grounded: JDR 0001
  §D13 as recorded in `artifacts/deviations.md` D3 — "0003 … declares spelling
  and universe only"; "Encoding is carriage, not declaration semantics". The
  two normal forms agree (unordered/duplicate-free literal vs.
  sorted/duplicate-free carriage), so no REQ here constrains `Tag.Value` bytes.

---

## QUESTIONS

None.

Four clauses were examined as ambiguity candidates and each resolved against
the record's own text, the deviations file, or the predecessor artifacts, so
none is recorded as a QUESTION:

1. **`exists` over an always-present key** — REQ-5 calls such an atom "well-
   formed but vacuous — lint reports it as such rather than rejecting it",
   while the matrix's acceptance column reads "any kind, provided the tag is
   declared optional", which could be read as a rejection. The two are
   reconciled by REQ-55 (an always-present key "contributes no presence
   dimension") and by the matrix's own parenthetical, which is the more
   specific statement and explicitly says "rather than rejecting it". Resolved:
   accepted, reported as vacuous, no `{absent}` assignment contributed.
   Recorded as REQ-5.

2. **Whether a group with no `unless` block is silent or covered** — the
   cluster gate flagged this (deviations D1). It is not ambiguous *here*: this
   RDR's `0003:C17` is scoped to decided atoms and is silent, and the identity
   (`unless = ∅` subtracts nothing) is fenced in RDR 0007. Filling the silence
   with RDR 0007's identity "changes no 0003 clause" (deviations D1, condition
   (c)). Resolved by deferring to the peer's fenced clause; the Stage 8 check
   is MVV Scenario 2 over a group with no `unless`.

3. **Whether withholding suppresses the coverage-gap finding** — REQ-95
   requires every decidable defect to be reported, while REQ-44 forbids
   emitting `graph-coverage-gap` for a withheld group. These are not in
   conflict: REQ-95 scopes the gap finding to "any coverage gap over a provable
   product", and a withheld group has no provable product (REQ-43). Overlap
   findings among decidable rows still MUST be emitted (REQ-44). Resolved
   against the record's own scoping language.

4. **Whether "refuse or downgrade" is one outcome or two** — the phrase recurs
   in REQ-62, REQ-85, and the disposition table, reading as a choice.
   `0003:C20` closes it normatively: "one outcome, not an author's choice",
   both carriers blocking, no non-blocking tier. Recorded as REQ-65/REQ-66.
