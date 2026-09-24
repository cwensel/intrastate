# REQ List — RDR 0030 Counting and stepping a tag without writing out every value by hand

Phase 0 audit artifact. Every clause below is a testable obligation drawn from
`docs/rdr/0030-computed-write-value-grammar.md`. Quotes are verbatim, copied
from the projector (`recs inspect --select <id>`), never transcribed; a line
break inside a quoted source span is rendered as one space and `…` marks an
elision.

Element ids are carried where the REQ derives from a labelled contract
(`0030:C1` … `0030:C3`, `0030:MVV`), a Load-Bearing Decision
(`0030:D-identity`, `0030:D-wire-byte-format`), a Failure Mode
(`0030:F1` … `0030:F6`), the Cross-Cutting block (`0030:G-cross-cutting`), or
a Testing Strategy scenario (`0030:S1` … `0030:S11`). Source sections are
abbreviated:

- `NC` = Proposed Solution / Technical Design / Normative Contracts (fenced)
- `TD` = Technical Design, unfenced prose above the Normative Contracts
- `MC` = Technical Design / Mini-check tables
- `LBD` = Technical Design / Load-Bearing Decisions
- `AUD` = Existing Infrastructure Audit
- `DR` = Decision Rationale
- `CONS` / `RISK` / `FM` = Trade-offs / Consequences, Risks and Mitigations, Failure Modes
- `MVV` = Implementation Plan / Minimum Viable Validation
- `IP` = Implementation Plan (prerequisites and phases)
- `TS` = Validation / Testing Strategy
- `G` = Finalization Gate / Cross-Cutting Concerns

**Two overrides of a predecessor are carried by the record itself** (its
`Overrides` field; never amended here, listed so a later stage does not read
them as contradictions):

1. `0002:C4` — the write-value grammar admitted at load is EXTENDED by the
   `{ step = <n> }` form; C4's "a write replaces" application semantics are
   untouched and hold per expanded row (REQ-25).
2. `0002:C13` — the suffix-iff ("a suffix is non-empty exactly when the rule
   produced more than one row") is NARROWED to the `in` expansion it was
   written about; a step point appends its cell element even at one admitted
   cell (REQ-17). `0002`'s REQ-76 (the `in` arm) stays true as written.

**Post-lock drift the auditor found, carried as ASSUMPTIONS below.** This
record locked 2026-09-20 with RDR 0012 still `Draft`; 0012 has since been
implemented (`7d57cf6`). The record's §Prerequisites and `0030:A4`/`0030:A12`
evidence describe the pre-0012 `Evaluator` (`struct{}`, raw-string `eq`). On
the current base the `Evaluator` carries a `kinds map[string]string`, is
constructed through `guard.NewEvaluator`, answers `GuardUnevaluable` for every
`eq`/`in` atom when built with no mapping, and `graphlint/reach.go::atomAdmitsValue`
byte-compares `eq`/`in` in-site (0012 deviations D1). JDR 0004 (JD-1, JD-3),
which binds both 0012 and 0030, decides how the two fit together; see the
first four ASSUMPTIONS.

---

## Step write grammar and admitted kinds (`0030:C1`)

- [REQ-1] "STEP WRITE. A rule's write block MAY assign an owned tag the inline table `{ step = <n> }`, `n` a non-zero TOML integer, and no other table shape: any other key set, a zero step, or a non-integer step (a TOML float `1.0` included) is refused at load under `malformed_tag_declaration` with a detail naming the one admitted form." — (NC, `0030:C1`; `0030:S5`) With "the TOML inline table `{ step = <n> }`, exactly one key." … "A bare table with the reporter's `add`/`next` keys is rejected as a second spelling of one operation; the detail names `step`." — (LBD, `0030:D-wire-byte-format`)

- [REQ-2] "The form is admitted only on a tag whose declared kind is `int` with both `min` and `max`, or `enum` with a non-empty `domain`; on `bool`, `set`, `scalar`, or an `int` missing a bound it is refused under the same category." — (NC, `0030:C1`; `0030:S5`) With "load refusal, `malformed_tag_declaration`, detail names the admitted kinds." — (FM, `0030:F1`)

- [REQ-3] "An `int` whose declared width is not representable — `max - min` negative, or equal to `math.MaxInt`, so the inclusive `+1` wraps — is refused too" … "That is `internal/guard/declaration.go::intWidth`'s rule, and it must be the SAME rule: a declaration whose width lint cannot compute — it saturates `Cardinality` to the ceiling and refuses the product as over-large — must not be one the loader enumerates." — (NC, `0030:C1`; `0030:S5`, whose negative-width arm is a pre-existing `min exceeds max` refusal restated, and whose `max - min == math.MaxInt` arm is new)

- [REQ-4] "the rule MOVES to the shim in the EXISTING `internal/resolve` package beside the relocated `Evaluator` rather than being copied" … "it is EXPORTED at its new home as `resolve.IntWidth(minV, maxV int) (int, bool)`." … "`::IntDomain(d table.TagDecl) []int` does NOT move" … "`guard`'s in-package re-inline (`assignment.go::valueAssignments` walks the same counted loop) and `::domainSize` both repoint through the moved `resolve.IntWidth`." — (NC, `0030:C1`; MC `authority` rows 2–3) Observable: exactly one width rule in the module — `resolve.IntWidth` with that signature; `internal/guard` declares no `intWidth`; `domainSize`, `IntDomain` and `valueAssignments` reach it.

- [REQ-5] "the loader's emitted `guard.all eq = <cell>` must render byte-identically to the cells `valueAssignments` produces or the atom never satisfies at lint." — (NC, `0030:C1`)

- [REQ-6] "The bound is representability, not size — a wide but representable domain is admitted and merely expensive (A7), and the record takes no load-time row ceiling." — (NC, `0030:C1`; RISK "C1 takes no load-time ceiling"; G "No load-time ceiling is owed.")

- [REQ-7] "A step write over an `enum` whose `domain` repeats a member is refused as well" — (NC, `0030:C1`; `0030:S5`) and "A step write over a domain containing a member with the suffix separator `#` is refused at load." — (NC, `0030:C1`; `0030:S5`; RISK "C1 refuses the step over such a domain at load")

- [REQ-8] "**Both guards are step-scoped, not new declaration rules**: an `enum` whose domain repeats a member or carries a `#` member and that NO rule steps continues to load exactly as it does today (A3 records the `#` case as loadable). The refusal fires on the step write, not on the declaration, because it is the step that indexes positions and mints the suffix; this record changes no behaviour on a model that steps nothing." — (NC, `0030:C1`; `0030:S5` negative controls)

- [REQ-9] "`[initial]` values and predicate literals keep the literal-only grammar; the table shape is admitted on the write-block path alone. Each refusing path keeps its OWN existing category, and this record merges none of them: `[initial]` refuses under `malformed_initial_declaration` (`internal/table/load.go`'s `[initial]` arm), a predicate literal under `malformed_predicate_atom` (`internal/table/normalize.go::loader.atom`'s `badAtom`), and the write-block path under `malformed_tag_declaration` (`::renderWrites`)." — (NC, `0030:C1`; `0030:S4`, which covers the predicate literal in all three atom positions: `match … eq`, `guard.all … eq`, and as an `in` member)

- [REQ-10] "Every refusal THIS clause and its surfaces mint — the kind and grammar refusals, the zero-cell refusal, the `#`-member and duplicate-member refusals, and C2's bound — is on the write-block path and so takes `malformed_tag_declaration`." — (NC, `0030:C1`; MC `disposition` table; LBD Naming "The refusal category is reused, not minted")

## Expansion, admitted cells and per-row shape (`0030:C1`)

- [REQ-11] "The loader expands a rule carrying a step write into literal rows, one per ADMITTED CELL of the stepped tag, composed into the rule's existing expansion product (`0002:C13`) as one more choice point. An admitted cell is a domain member — `{min..max}` for `int`, the authored `domain` order for `enum` — that satisfies the CONJUNCTION of every positive atom the rule AUTHORS on that tag, however many per block (`guard.all` `eq`/`in`/`lt`/`lte`/`gt`/`gte`; `match` `eq`/`in`), each evaluated per member as the runtime evaluator would, against the AUTHORED literal and not against any member a sibling choice point has since chosen; `unless` atoms are not consulted." — (NC, `0030:C1`; RISK "two atoms conjoin; the MVV fixture carries such a rule"; LBD `Selection / predicate`)

- [REQ-12] "on a stepped tag the step point SUBSUMES a match `in` on that tag rather than composing with it — whether the `in` is authored locally or inherited from a context, since `expand` receives the merged predicate set and `0002:C13` expands both alike. One choice point is emitted per stepped tag, never two, and the `in` contributes its members to the conjunction instead of expanding separately." — (NC, `0030:C1`; `0030:S8`)

- [REQ-13] "A subsumed `in` is REWRITTEN per row to `match eq = <cell>`, exactly as `expand`'s existing `case expanding:` arm rewrites an `in` it expands — never retained in its `in` form." … "So this mechanism contributes exactly two atoms to the stepped tag per row — the rewritten `match eq = <cell>` where an `in` was subsumed, and the expansion's own `guard.all eq = <cell>` — and no more." … "every OTHER authored atom on the tag, the `guard.all` bounds among them, is retained unchanged beside these two" — (NC, `0030:C1`; `0030:S8`; MC `trace` step 4)

- [REQ-14] "the shim takes a `resolve.GuardAtom` and the held value as a STRING — the cell, never the declared kind — and returns `resolve.GuardResult`, whose members are `GuardTrue`, `GuardFalse` and `GuardUnevaluable`." — (NC, `0030:C1`) See the first and fourth ASSUMPTIONS for how this meets the post-0012 constructed evaluator.

- [REQ-15] "The runtime evaluator is three-valued; an atom answering neither true nor false at a member does NOT admit that member." — (NC, `0030:C1`; TD "C1 fixes the undecided arm to EXCLUDE the member"; MC `authority` row "The UNDECIDED verdict's disposition … loader EXCLUDES") A soundness fence unreachable through a loaded model; asserted at the filter's own seam with an evaluator answering `GuardUnevaluable`.

- [REQ-16] "Each expanded row carries a `guard.all` atom `eq = <cell>` on the tag, the literal write `cell + n` (`int`) or the domain member `n` positions from the cell (`enum`), and the cell appended to its expansion suffix." — (NC, `0030:C1`; `0030:S3` (3c ii), `0030:S9`, `0030:S10`)

- [REQ-17] "The suffix element is appended at EVERY admitted cell, including when exactly one is admitted: a stepped row is never the authored row, so it always carries its cell." … "the two arms keep their own rules, and the `in` expansion's behaviour is unchanged by this record." — (NC, `0030:C1`; LBD `0030:D-identity`; `0030:S9` (a) and the single-member `in` control (c))

- [REQ-18] "The stepped value is computed without overflow at every admitted cell, whatever the step's magnitude; whether it lands inside the domain is C2's bound, not an admission question." — (NC, `0030:C1`; ordering pinned by REQ-30)

- [REQ-19] "A rule admitting zero cells is refused at load." — (NC, `0030:C1`) With "A rule with zero admitted cells is refused, and that refusal is NEW. The loader does not refuse an unsatisfiable guard today" — (TD) and "A step rule admitting no cell — a NEW load refusal taking the `0003:C6` disposition." — (FM, `0030:F3`; `0030:S5` with the unstepped `lt = 0` control that still loads)

- [REQ-20] "A rule that BOTH steps a key and names it in its `clear` list is refused at load under the same category" … "the collision refusal PRECEDES the per-cell walk: a key both stepped and cleared is refused before any cell is evaluated, so it never reaches C2's bound. Where a rule would trip both, the collision is what the author is told" — (NC, `0030:C1`; `0030:S5`)

- [REQ-21] "**The stepped key is placed in `assignments` like any other written key**, carrying a placeholder the expansion replaces per cell: `0002:C14` derives `RequiresOwned` from `maps.Keys(assignments)`, so a spec tracked beside the map instead would silently drop the stepped key from the kernel's owned-state gate." — (NC, `0030:C1`) Observable: every expanded row's `RequiresOwned` names the stepped key.

- [REQ-22] "`::expand` then does only what it does for `in` today: take a list of members per key and mint one row per combination, carrying the supplied literal and appending the suffix element. It therefore stays TOTAL — it keeps its `[]Row` return and gains no error — because every refusal this record mints has already fired before it is called." — (NC, `0030:C1`; IP Phase 2 "the phase adds no fourth site") Observable: `expand`'s result list is `[]Row` alone.

- [REQ-23] "The per-cell stepped literal is the row's rendered assignment for that key, so it lands in BOTH carriers the assignment feeds: `Row.Writes` and `Row.NextTags` are populated per row from the same stepped value, never by aliasing one to the other (`0002:C15`)." — (NC, `0030:C1`; MC `authority` row "neither — C1 requires BOTH be set per row"; `0030:S10`)

- [REQ-24] "the step expansion rides the same `expand` loop and copies `Emit`, `Gate`, `RequiresOwned` and `Escape` identically to every row it mints, so every row a stepped rule expands to carries the same block" — (AUD, "Join a rule's `emit` block back" row; `0030:S7`)

- [REQ-25] "After normalization no surface distinguishes an expanded row from an authored literal row: both carriers hold literals only, and `0002:C4`'s write-replaces clause applies per row unchanged." — (NC, `0030:C1`) Asserted by REQ-32's invariance and `0030:S2`/`0030:S6`.

## Suffix ordering (`0030:D-identity`)

- [REQ-26] "**Step points order by KEY alone.**" … "So `attempt = { step = 1 }` with `tier = { step = 1 }` yields `retry#<attempt>#<tier>`, ordered by the key names. **The same key comparison places a step point among the OTHER candidates** — unsubsumed `in` points on other keys, and the outcome point — so the whole sorted list is keyed consistently and a rule with a step on `attempt` and an `in` on `mode` yields `retry#<attempt>#<mode>`." … "The element is the bare cell, not tag-qualified" — (LBD, `0030:D-identity`; `0030:S9` (d))

## Bound (`0030:C2`)

- [REQ-27] "BOUND. At every admitted cell the stepped literal MUST conform to the tag's declaration exactly as an authored literal write does (`int` within `min..max`; `enum` within `domain`). A cell whose stepped value leaves the domain in either direction — past `max` or below `min`, past the last member or before the first — is a load refusal under `malformed_tag_declaration` — the category a non-conforming literal write already takes — whose detail names the rule, the tag, the cell, and the result." — (NC, `0030:C2`; `0030:S3`, `0030:S3` 3c; FM `0030:F2`)

- [REQ-28] "for `int` it is the computed literal (`10`), which `internal/table/load.go::conformDomain` already reports as \"is above max 9\" / \"is below min 0\"; for `enum` there IS no stepped value — no member sits `n` positions from the cell — so the slot names the signed offset off the end instead (`1 past the last member`, `2 before the first`), and on the checked-add arm (A13) there is no representable literal either, so the slot names the bound it passed." — (NC, `0030:C2`; `0030:S3`)

- [REQ-29] "All four parts ride the refusal's `Detail` STRING" … "This record adds none" … "Tests therefore assert the `Category` plus the detail text, and the detail names the cell FIRST so the assertion is anchored at a stable position rather than mid-sentence." — (NC, `0030:C2`) Includes: `Failure.Rule`, `Offending` and `Remedy` stay unpopulated on this refusal ("populating it with a rule id would publish that hint on the wire").

- [REQ-30] "Where more than one admitted cell fails, the detail names the FIRST in the tag's own domain order — ascending `min..max` for `int`, the authored `domain` order for `enum` — not the row's suffix order" — (NC, `0030:C2`; `0030:S3` 3b; FM `0030:F3` "the detail names the FIRST such cell in the tag's own domain order")

- [REQ-31] "There is no saturating and no wrapping form. The stepped value is computed as a CHECKED `int` add, and that check is ordered BEFORE the render: a step whose magnitude carries a cell outside `int` range is refused as this bound failure, naming the cell, rather than being rendered and handed to conformance." — (NC, `0030:C2`; `0030:S5` A13 fixture, detail names the cell and the bound, NOT "is not an int")

- [REQ-32] "Because `unless` atoms are not consulted when cells are admitted (C1), the detail also names any `unless` atom the rule authors on the stepped tag and says it did not exclude the cell" — (NC, `0030:C2`) With "surfaces as C2's bound refusal, whose detail names the unconsulted `unless` atom and directs the author to the positive-atom (`guard.all`) rewrite" — (TD)

## Invariance and published identity (`0030:C3`)

- [REQ-33] "INVARIANCE. `flow next`, `flow resolve`, `flow set-state`, `lint`, `dump`, and `graph` read expanded rows through `Row.Writes`, `Row.NextTags` and `Row.Atoms` alone — the two write carriers, not one (`0002:C15`; C1's per-row clause is what makes both correct) — and none learns a computed value shape; the graph document's `intrastate.graph/1` vocabulary gains no member." — (NC, `0030:C3`; `0030:S6`; MC `fidelity` row "decodes under shipped `intrastate.graph/1` with no new member")

- [REQ-34] "No emitted vocabulary — `internal/table::Categories()`, the CLIError codes, the finding codes, the graph document's members — gains a member, so no `0029:C4` stability tier is owed by this record." — (NC, `0030:C3`; `0030:S6` golden enumerations) With "Done = every scenario below green with no new finding code and no new envelope member." — (TS preamble)

- [REQ-35] "The authoring guide (`docs/model-authoring.md`) documents the form, the admitted kinds, the admitted-cell rule, and the bound refusal beside the existing write-block section — and, because none is inferable from the form, the four hazards this record creates for an author: that an authored `enum` `domain`'s ORDER is the step order, so reordering a stepped domain changes the model with no lint to catch it; that `unless` never excludes a cell from a step, so the exclusion must be a positive atom; that on an `enum` the only atom that can exclude the terminal member is an `in` re-listing the domain minus that member" … "and that the cell count is the LINT cost, not the load cost" … "What the author sees is `lint` not returning. The guide says exactly that, and names the shape of the cliff — it is the admitted-cell count, not the declared `min..max` width, that drives it, so narrowing the rule's atoms is the fix" … "`intrastate --help-all` is regenerated." — (NC, `0030:C3`; IP Phase 3)

- [REQ-36] "Every surface that names a ROW publishes the suffixed identity `rule#cell` — the shape `in` expansion already emits there (`dump`'s `identity` column, the graph document's `identity` row field) — and every surface that names a RULE publishes the authored rule id (`flow next`'s and `flow resolve`'s `rule`, the graph document's edge `rule`)." — (NC, `0030:C3`; `0030:S7`, `0030:S9`; FM `0030:F6`)

- [REQ-37] "A per-row graph-lint finding on an EXPANDED row therefore publishes the suffixed identity `rule#cell`, not the bare authored id: lint names a ROW, because its subject is a row. That covers BOTH row-naming slots — `Rule`, and the `Element` slot the two findings that name a SECOND row populate" … "A finding may not suffix one slot and not the other" … "The group-level findings that name a group rather than a row (`coverage.go`'s `firstRuleID` sites) are unchanged — they name no row and gain no suffix." — (NC, `0030:C3`; `0030:S11` (a), (b); MVV item 2) Scope of "EXPANDED row": see QUESTIONS Q1.

- [REQ-38] "The suffix is therefore applied at EMISSION and the in-package read-back joins on the AUTHORED ID RECOVERED from the suffixed form — the published identity truncated at the first `#`" … "No in-package consumer may resolve a per-row finding's `Rule` against authored `[[rules]]` ids by equality." … "`internal/graphlint/engine.go::identityKey` … therefore satisfies the fence unchanged and keeps its ordering guarantee under a suffixed value (A15)." — (NC, `0030:C3`; `0030:S11` (c): the `ambiguous_match` coverage arm still fires)

- [REQ-39] "the field's published meaning narrows from \"the authored rule id\" to \"the row's identity\", and `docs/cli-output-contract.md` describes `rule` as graph-lint attribution, so an out-of-repo consumer joining on it is affected in the same way and `span` is the stable key it should use." — (NC, `0030:C3`; G Versioning "one declared exception") See the ASSUMPTION on the contract document.

## Shared comparison (`IP` Phase 2, `MC`, `AUD`)

- [REQ-40] "Move `internal/guard/grammar.go`'s `Evaluator` to `internal/resolve`, then extract the two-valued core of `internal/guard/product.go::valueSatisfies` — render the atom's literal, call the evaluator, return its three-valued verdict — to `internal/resolve` alongside it, carrying the canonicalizing set renderer." — (IP, Phase 2; AUD "Reuse (extract + move)"; MC `authority` row 1 "the shim") With "The loader's admits filter then calls the same comparison the runtime does, applying C1's exclude collapse at its own call site, rather than minting a fourth copy." — (IP, Phase 2) Observable: the loader's admits filter, `guard::valueSatisfies` and `graphlint::atomAdmitsValue` reach one comparison in `internal/resolve`; the loader holds no reimplementation.

- [REQ-41] "Each caller's undecided-arm disposition stays its own — the shim is the comparison, never the policy — so the repoint is behaviour-preserving in all three." — (IP, Phase 2) With "`guard` passes the verdict through, `graphlint` admits, the loader excludes per C1" — (AUD) Observable: no committed model's lint verdict or resolve plan changes across the move and repoint.

- [REQ-42] "the loader takes the CANONICALIZING form, since a set literal's two spellings are one literal" — (AUD; MC `authority` row "Set-literal rendering for `in`" canonical "the canonicalizing one"; G "the set literal the loader renders takes the CANONICALIZING form")

## Derived outcomes (`CONS`, `FM`, `G`)

- [REQ-43] "A cap cell excluded by the step rule and claimed by no other row — `graph-coverage-gap` at lint, the existing finding." — (FM, `0030:F4`; `0030:S2` 2b) and "A stepped tag never initialised — `graph-owned-before-write` at lint." — (FM, `0030:F5`)

- [REQ-44] "appending `critical` to a domain whose step rule carries `in = [\"small\",\"medium\"]` does NOT refuse at load (the `in` still admits only the cells it names, so no stepped value leaves the domain) and does NOT silently widen the ladder (no cell steps into `critical`). The ladder keeps stopping at `large`, and the new member surfaces as a BLOCKING `graph-coverage-gap` at lint" — (CONS)

- [REQ-45] "INTERIOR over-admitted cell conforms, so it mints a row — but that row is DEAD, not overlapping" … "`graph-overlap` cannot fire; what fires is the advisory `graph-redundant-row`" — (TD; MC `disposition` row "`unless` over-admits an INTERIOR cell"; `0030:S5` 5b, 5c)

- [REQ-46] "No map iteration order reaches the row order." — (G, `0030:G-cross-cutting`) Observable: loading the same stepped model repeatedly yields byte-identical `dump` output and row order.

- [REQ-47] "a model with no `{ step = <n> }` write value loads, lints, resolves and exports exactly as today, because the new arm is intercepted before `load.go::valueMembers` on the write path only." — (G, `0030:G-cross-cutting`) Observable: the committed corpus's `lint`/`dump`/`graph` outputs are unchanged; see Q1 for the one reading under which a per-row lint `rule` on an `in`-expanded row changes.

## Minimum Viable Validation (`0030:MVV`)

- [REQ-MVV] "Author two in-repo fixtures of the Background's ladder over `internal/table/testdata/rdr-fixture.toml`'s `iter`-shaped `int` tag (`min = 0, max = 9`) and an ordered `enum` tier the fixtures DECLARE as `tier` with `domain = [\"small\",\"mid\",\"large\"]`" … "`ladder-literal.toml` (unrolled, as today) and `ladder-step.toml` (`attempt = { step = 1 }`, `tier = { step = 1 }`, each step rule carrying at least one row with two positive atoms on the stepped tag)." … "2. `intrastate lint` on both: identical finding sets under S1's `(code, key, dimension, class, reason)` projection, zero `graph-overlap`, zero `graph-coverage-gap` (A5). Plus the half that projection cannot see: a per-row finding on the stepped ladder names the offending CELL (`retry#3`), not the bare rule (S11, A15)." … "3. `flow resolve` over every (state, outcome) cell of both: identical `writes`/`next`/`clear` (A6)." … "4. `ladder-step-unguarded.toml` (the `lt 5` atom removed): load refuses `malformed_tag_declaration` naming `retry`, `attempt`, cell `9`, value `10`; the enum sibling names `tier`, `large`. Two more refusal fixtures: `[initial] attempt = { step = 1 }` refuses `malformed_initial_declaration` and a predicate literal `{ step = 1 }` refuses `malformed_predicate_atom`" … "5. The dead-row outcome and A5's negative control, which take TWO fixtures rather than one (S5b, S5c): (a) a stepped rule with `unless` on the stepped tag excluding an interior cell a literal row claims — loads, and lints exactly one advisory `graph-redundant-row`; (b) the same exclusion written with a POSITIVE atom pair instead, where every row projects — zero `graph-overlap`." … "6. `graph` export of `ladder-step.toml` decodes under the shipped `intrastate.graph/1` document type with no new member (C3)." … "7. Source legibility, the outcome the Problem Statement names: the `[[rules]]` count of `ladder-step.toml` is ONE PER INTENT — one rule per (ladder, outcome) pair rather than one per cell — and is strictly less than `ladder-literal.toml`'s. Both assertions are derived from the fixtures as authored (a count over the parsed source, not a line count)" … "End-state: the step ladder is the literal ladder to every consumer, it is one row per intent in source, and the unguarded cap is a named load refusal." — (MVV, `0030:MVV`)

  The MVV's oracles and negative controls are the MC `oracle` table's rows;
  each is carried by the scenario REQ it names below.

## Testing Strategy obligations (`0030:S1` … `0030:S11`)

- [REQ-48] "identical finding sets compared as sets over the PROJECTION `(code, key, dimension, class, reason)` of each `findings[]` entry — not over whole finding objects." … "The compared set MUST be non-empty, asserted BEFORE the set comparison" … "(a) the CLEAN pair — `ladder-literal.toml` / `ladder-step.toml`, which every other scenario here reuses — must exit OK and yield a non-empty ADVISORY set, guaranteed by authoring a `graph-idempotent-write` on an UNSTEPPED key (a non-zero step is never idempotent) mirrored on both sides; and (b) a BLOCKING pair, `ladder-literal-blocking.toml` / `ladder-step-blocking.toml`, authored solely for this tier and reused by nothing else" … "The blocking defect is `graph-owned-before-write` on a tag neither ladder initialises" … "Zero `graph-overlap` and zero `graph-coverage-gap` on all four. Backs A5." — (TS, `0030:S1`)

- [REQ-49] "`flow resolve --as json --plan-only` swept over every (state, outcome) cell of the pair, the state driven by rewriting the bound artifact file. The sweep's domain is the FULL cross product of both stepped tags' declared domains × the outcome set" … "**Expected**: `writes`, `next` and `clear` equal cell by cell." … "2b. **Scenario**: unclaimed-cell parity. The F4 shape — a cap cell excluded by the step rule and claimed by no other row — authored both ways, swept as in 2. **Expected**: BOTH sides refuse at that cell, with the same code." — (TS, `0030:S2`)

- [REQ-50] "`ladder-step-unguarded.toml` — the `lt 5` atom removed — loaded. **Expected**: refusal under `malformed_tag_declaration` naming rule `retry`, tag `attempt`, cell `9`, result `10`; the enum sibling names `tier`, `large`, and `1 past the last member`" … "3b. … **Expected**: the detail names cell `9` (domain-first), NOT `10` (suffix-lexical-first); and `mid` (authored-domain-first), not `large`." … "3c. … (i) refuses, naming cell `0` and `min` on the `int` arm and `1 before the first member` on the `enum` arm — the only observation a refusing load affords. (ii) loads, and its rows are `retry#1 … retry#9` each writing `cell - 1`." — (TS, `0030:S3`)

- [REQ-51] "`[initial] attempt = { step = 1 }` and a predicate literal `{ step = 1 }`, each loaded. **Expected**: both refuse, and under the categories C1 names — the `[initial]` one `malformed_initial_declaration`, the predicate one `malformed_predicate_atom`." … "The predicate-literal fixture covers all three atom positions — `match … eq`, `guard.all … eq`, and as a member of an `in` list" — (TS, `0030:S4`)

- [REQ-52] "zero-cell (from an ordered bound, e.g. `lt = 0` over `min=0,max=9` — the shape the loader does NOT refuse today), zero-step, float-step, `#`-in-domain, and duplicate-member-in-domain step models loaded; a step write on `bool`, `set`, `scalar`, an unbounded `int`, and an `int` whose width is not representable; a rule that both steps a key and names it in its `clear` list; and a REPRESENTABLE-width `int` whose step magnitude carries an admitted cell past `math.MaxInt`" … "**Expected**: each a load refusal under `malformed_tag_declaration` whose detail names the admitted form or kinds." … "Each carries its paired negative control" … "5b. … **Expected**: loads (no refusal — the stepped value conforms), and lints EXACTLY ONE advisory `graph-redundant-row` (the dead row)" … "The scenario therefore asserts the projectability precondition explicitly" … "5c. … **Expected**: ZERO `graph-overlap`." — (TS, `0030:S5`)

- [REQ-53] "`graph` export of `ladder-step.toml` decoded against the shipped `intrastate.graph/1` document type. **Expected**: decodes with no new member — asserted by STRICT decode into the shipped document type (unknown fields rejected), the oracle `internal/cli/graph_mvv_0021_test.go` already uses; and `internal/table::Categories()`, the CLIError codes and the finding codes each gain none, asserted as golden enumerations" — (TS, `0030:S6`)

- [REQ-54] "`flow resolve` selecting a NON-FIRST expanded row of a stepped rule that authors an `emit` block" … "**Expected**: the payload's `emit` and `dispositions` are the authored block, and `rule` is the authored id" — (TS, `0030:S7`; MC `oracle` row 7 "the first expanded row, which would pass trivially")

- [REQ-55] "A step rule on `attempt` that ALSO authors `match.attempt in = [1,3]`, and a sibling fixture where that same `in` is inherited from a context rather than authored on the rule. **Expected**: both yield exactly |in ∩ admitted cells| rows — two, not the |members| × |cells| product — each carrying `match eq = <cell>` with NO surviving `in`, and a suffix carrying exactly one element for `attempt`. Asserted CLI-level on `dump`'s `atoms` column" — (TS, `0030:S8`)

- [REQ-56] "(a) publishes `retry#<cell>` in `dump`'s `identity` column and the graph document's row `identity` — the suffix is appended even at one cell; (b) publishes `retry#0 … retry#4` on those same surfaces; (c) mints NO suffix element" … "(d) a rule carrying a step on `attempt` AND an unsubsumed `match in` on `mode`, which publishes `retry#<attempt>#<mode>` — the suffix elements ordered by KEY, not by point kind." — (TS, `0030:S9`)

- [REQ-57] "`graph` export and `dump` of `ladder-step.toml`, read for each expanded row's successor. **Expected**: every row's `next` (graph) and next column (`dump`) carries THAT row's stepped value, not the authored one." — (TS, `0030:S10`)

- [REQ-58] "**Expected**: the finding's `rule` reads `retry#3` — the offending cell — and the literal twin's reads `retry-3`" … "Assert additionally that a GROUP-level finding (`coverage.go`'s `firstRuleID` sites) still reads the bare authored id" … "(b) a `graph-overlap` between a literal row and an expanded one: BOTH `rule` and `element` carry their own row's identity — the expanded side suffixed, the literal side not" … "(c) the same fixture authored so the group's `ambiguous_match` coverage arm is reachable: the arm must still fire." — (TS, `0030:S11`; MC `oracle` rows "2 lint attribution (S11)" and "2 `ambiguous_match` survival (S11c)")

---

## EXCLUDED

- EXCLUDED: "This record does NOT wait on RDR 0012." … "0012's typing work later targets `internal/resolve` instead of `internal/guard` — a relocation 0012 absorbs, not a dependency this record has." — (IP, Prerequisites) — record-state precondition, now overtaken: 0012 is implemented on the base. How the relocation carries 0012's typed evaluator is the first four ASSUMPTIONS.
- EXCLUDED: "RDR 0002 is `Implemented` on `main` (it is)" and "All Critical Assumptions verified" — (IP, Prerequisites) — record-state preconditions, not behaviour.
- EXCLUDED: "Determinacy: fired — C1 (step ordering …), C1 (identity …), C3 (identity …)." — (NC preamble) — a lens annotation; the decisions it points at are REQ-20/REQ-26/REQ-38.
- EXCLUDED: "The comparison the filter calls is the extracted shim, and its surface is what is on `main` today, not a new one" … "stating it here rather than leaving it to §Prerequisites because the held-value type is what fixes 0012's later typing work to `internal/resolve`" — (NC, `0030:C1`) — rationale around REQ-14; the pre-0012 "what is on `main` today" half is overtaken (first ASSUMPTION).
- EXCLUDED: "The two widened sites own different halves" … "`::renderWrites` resolves the step spec into the per-key list of admitted cells" … "So `renderWrites` takes the merged predicates as a parameter. Both widenings are local to `normalize.go` and neither adds a site." — (NC, `0030:C1`) — code-siting; REQ-22 carries the observable half (`expand` stays total), REQ-20 the ordering.
- EXCLUDED: "The two arms keep their own rules … That local is computed per choice point, so the step point computes its own as constant-true rather than changing the `in` arm's" — (NC, `0030:C1`) — code structure; REQ-17 carries the observable.
- EXCLUDED: "That split is the existing design, not a gap: a flow payload names the rule it matched, so this record mints no obligation on it." — (NC, `0030:C3`) — a no-obligation statement; REQ-36's flow half is existing behaviour S7 re-pins.
- EXCLUDED: "Per-row identity IN a flow payload is a possible successor record — disclosed here as a non-obligation, not an owed fix." — (DR) — deferred.
- EXCLUDED: "Accepted, not closed — shipping an equivalence checker is a verb, and this record mints none." and "A `dump`-diff or an equivalence mode is the natural successor if conversions prove common." — (CONS) — deferred.
- EXCLUDED: "Extending the guard grammar with an ordered or `neq` atom over an authored `enum` domain would close it and is out of scope here" — (CONS) — deferred; REQ-44 pins the residual's observable.
- EXCLUDED: "Load and normalize are **linear** in the expanded row count: 0.96 s and 584 MB resident at 100,001 rows, 5.08 s at 500,000" and the lint timings "7.5 s at 100 rows, 116 s at 200, beyond 600 s at 1,000" — (Performance Expectations) — spike measurements (`0030:A7`), not an assertion; REQ-6 carries the no-ceiling observable. "Were a ceiling ever wanted it would be memory-motivated" is hypothetical.
- EXCLUDED: "Accepted scope: this is a package move, larger than a function extraction" and "roughly twenty-five `Evaluator` references across seven `internal/guard/*_test.go` files" — (IP, Phase 2) — scope/census notes; the census is stale post-0012 (see `impact.md`).
- EXCLUDED: "Fixtures live in `internal/table/testdata/` beside `rdr-fixture.toml`" … "the assertions are CLI-level, on the JSON envelopes, in `internal/cli/lint_mvv_0030_test.go` and `internal/cli/flow_mvv_0030_test.go` per the `<topic>_mvv_<rdr>_test.go` precedent." — (TS preamble) — test-siting; followed, not asserted.
- EXCLUDED: "This record claims **no** byte-identical output, content-addressed identity or replay-stable hash of its own" — (G) — a negative with no observable; `Fingerprint` and the export's pre-images are 0021's/0006's.
- EXCLUDED: "Not applicable, and deliberately not N/A-bulleted above: build tool compatibility, licensing, …" — (G) — no observable.
- EXCLUDED: Premortem, Ground-sweep and Joint-check notes (`0030:JC1`–`JC3`, §Decision Rationale) — cross-record coordination and review provenance; no observable. The `0013:C1` coupling is REQ-6's "no load-time ceiling".
- EXCLUDED: Alternatives 1–3 and Briefly Rejected (`0030:ALT1`–`ALT3`, `0030:BR1`–`BR5`) — rejected designs. Their negative consequences are carried positively: no saturate/wrap (REQ-31), no string sentinel (REQ-1), no `order` facet (REQ-26 reads authored `domain` order), no `bool` toggle (REQ-2).
- EXCLUDED: "Surface — of C1" headers and the Critical Assumptions' **If wrong** branches (`0030:A1`–`A16`) — record structure and conditional contingencies; the Pending assumptions' named plans (A5 → S1/S5b/S5c, A6 → S2, A10 → S8, A11 → S10, A13 → S5) are REQ-48…REQ-57.
- EXCLUDED: "A10 … Assert it through `KernelRow`, where a retained `in` would surface as `seamValue`'s first member on every row." — (`0030:A10` Evidence) — superseded by S8's oracle choice: "`KernelRow`/`::seamValue` is the mechanism, not the oracle" (REQ-55 asserts on `dump`'s `atoms`).

---

## ASSUMPTIONS

- ASSUMPTION: the admitted-cell filter compares through an evaluator CONSTRUCTED over the loader's own declarations — key → declared `Kind` for every tag in `l.model.Tags`, a key with an empty `Kind` omitted — never through the zero-value form. A zero-value evaluator on the current base answers `GuardUnevaluable` for every `eq`/`in` atom (`0012:C1`), which REQ-15's EXCLUDE arm would turn into zero admitted cells and REQ-19's zero-cell refusal on a correct model. Grounded: JDR 0004 JD-3 ("The loader-side admitted-cell shim constructs through `NewEvaluator` with the loader's own declarations"), which binds 0030 and says it "Lands in 0030: the Phase 2 extraction step names the construction"; C1's "evaluated per member as the runtime evaluator would … (`0012:C2`)". For the kinds C1 admits the typed comparison and a byte comparison agree on every cell (canonical `strconv.Itoa` cells; predicate `int` literals are canonical by `0012:C5`; `enum` compares exact strings), so this choice is behaviour-neutral for admission and only avoids the zero-value trap.
- ASSUMPTION: the relocation moves the `Evaluator` type, its `Evaluate` method and its private helpers to `internal/resolve` (JDR 0004 JD-1: "`internal/resolve` may host the value-comparison seam's implementation, declaration state included"), with a constructor reachable from `internal/table`. `internal/guard` keeps `Evaluator` as a re-export (type alias) and keeps `NewEvaluator(kinds map[string]string) Evaluator` and `DeclaredKinds(m *table.Model) map[string]string` declared there, so 0012's shipped pins (`TestReq1_0012_…`, `TestReq2_0012_…`, `TestReq10_0012_…`, `TestReq35_0012_…`, the REQ-48 zero-value check, which matches the type object through aliases per 0012 deviations D4) keep holding where they can. The loader-side mapping has one body. `guard.DeclaredKinds` may delegate to a `table`-side helper the loader also calls, following C1's own `IntWidth`/`IntDomain` pattern ("one width rule, two unpackings"). That helper is not named `DeclaredKinds` (0012 REQ-2 checks for a second producer by name). Any 0012 or 0003 test pin that the relocation still makes false (a home asserted by package directory, or the zero-value allow-list keyed on the `internal/guard/` prefix) is migrated and recorded in `deviations.md`, on the precedent of 0012 deviations D3. `impact.md` predicts 168 such tests across 11 families, all in `internal/guard` and `internal/resolve`. Grounded: C1 "the relocated `Evaluator`"; IP Phase 2 "`guard_evaluator_0003_test.go`'s field-count reflection assertion … follows the type to `resolve` rather than being rewritten".
- ASSUMPTION: the `graphlint::atomAdmitsValue` repoint (REQ-40/REQ-41) touches only its render-then-`Evaluate` fall-through for operators other than `eq`/`in`. Its in-site byte comparison of `eq`/`in` (0012 deviations D1), its name, home, `bool` return, `guard.Evaluator{}` construction and `!= resolve.GuardFalse` collapse all stay (0012 REQ-35/REQ-48). Grounded: REQ-41 requires the repoint to be behaviour-preserving, and 0012 D1 records that routing `eq`/`in` through a zero-value evaluator flips 15 committed models' lint verdicts.
- ASSUMPTION: REQ-14's shim surface ("takes a `resolve.GuardAtom` and the held value as a STRING … returns `resolve.GuardResult`") is the constructed evaluator's `Evaluate`, reached through the shared comparison. The render step is the canonicalizing set renderer exported beside it in `internal/resolve` (REQ-40, REQ-42), and each of the three callers renders its own `table.Atom` into a `resolve.GuardAtom` with it at the call site. Grounded: IP Phase 2 "each of the three callers renders its own `table.Atom` into that form at the call site"; "carrying the canonicalizing set renderer". `guard::valueSatisfies` now takes its seam as a parameter (0012), and it keeps doing so. Switching it from `renderSet` (as-authored) to the canonicalizing renderer preserves behaviour: the two differ only in member order, duplicates and the nil list. `in`/`contains` membership does not depend on order or duplicates. The loader refuses an empty `in`/`contains` list ("`in` takes a non-empty member set"), so the nil case never reaches the comparison from a loaded model.
- ASSUMPTION: the step interception runs inside `renderWrites`' write loop AFTER the existing undeclared-key (`unknown_tag`) and non-owned (`write_to_non_owned_tag`) checks, so a `{ step = n }` on an undeclared or non-owned tag keeps those categories. Grounded: C1 "MAY assign an owned tag"; REQ-9 "Each refusing path keeps its OWN existing category"; REQ-10 scopes `malformed_tag_declaration` to refusals "THIS clause and its surfaces mint".
- ASSUMPTION: "the detail names the cell FIRST" (REQ-29) means first within the text this record adds, after `renderWrites`' existing `rule <id> write <key>: ` prefix. Example: `rule retry write attempt: cell 9 steps to 10, which is above max 9`. Grounded: S3 "`renderWrites` already prefixes rule and key, so the cell and the enum offset are the new text".
- ASSUMPTION: C2's multi-failure rule (REQ-30) and the step/clear ordering (REQ-20) mean the per-cell walk runs over the tag's domain order and reports the first failure. The row set's lexical suffix order plays no part.
- ASSUMPTION: the `enum` sibling of the unguarded fixture (REQ-50, MVV item 4) is a separate fixture file, because a refusing load reports one refusal. Likewise every S5 refusal is its own fixture with its own paired control.
- ASSUMPTION: REQ-35's guide text follows C3's fenced wording for the wide-domain lint cost ("What the author sees is `lint` not returning"). It does not follow the unfenced RISK / Performance prose ("degrading the verdict (`complete=false`, reported as a finding)"). C3 states that nothing in-band signals in the band between a few hundred cells and the 2048/4096 bounds. The RISK/Performance prose describes what happens past those bounds. The guide may state both halves: a hang in the band, a degraded verdict past the bounds.
- ASSUMPTION: REQ-39 obliges an edit to `docs/cli-output-contract.md`. Its description of a graph-lint finding's `rule` gains the narrowed meaning (the row's identity; `rule#cell` on an expanded row) and names `span` as the stable key for out-of-repo consumers. Grounded: C3 "Stating this is part of the clause, not an implementation note: the field's published meaning narrows". `docs/cli-reference.md` and `llms.txt` change only through `make docs` (AGENTS.md §Docs).
- ASSUMPTION: the suffixed identity a lint finding publishes is `<rule id>#<suffix element>#…` — the authored rule id joined with the row's suffix elements by `#`, without the model-id prefix `model.go::Identity` carries. Grounded: S11 expects `retry#3`, MVV item 2 names `retry#3`, and the Decision Rationale ground-sweep notes "`Identity` joins the model id with `.`, and only the suffix elements take `#`". `dump`'s `identity` column and the graph document's row `identity` keep whatever they publish today (REQ-36 cites them as the precedent, unchanged).
- ASSUMPTION: REQ-45's "INTERIOR over-admitted cell" row keeps the authored `unless` atom, so its accepted set is empty. Grounded: TD "the authored `unless` is retained on it"; REQ-13 "every OTHER authored atom on the tag … is retained unchanged".
- ASSUMPTION (Q1, grounded by Phase 3d → reading (a)): C3's per-row lint suffixing applies to EVERY row carrying a suffix, `in`-minted rows included. A per-row finding's `Rule`/`Element` publish the authored rule id joined with the row's full `Suffix` (the bare id for an unexpanded row). Grounded: `internal/table/model.go:279` (`Row.Suffix` is one undifferentiated list of `in` and step elements, with no origin marker); `internal/table/normalize.go:805-814` (the `in` arm appends to that same list); C3 ("Every surface that names a ROW publishes the suffixed identity … the shape `in` expansion already emits there"); D-identity (a mixed row is `retry#<attempt>#<mode>`); G Versioning (the `rule` slot narrows to "the row's identity" with no step condition); A15 Verified plus `evidence-a12-a15.md:76-84` ("reachable without new state"). A committed golden that asserts a bare id on an `in`-expanded per-row finding is re-cut under REQ-39's narrowing, and that re-cut is logged in `deviations.md`.

---

## QUESTIONS

- **Q1 → RESOLVED (Phase 3d: reading (a); see the Q1 ASSUMPTION above).** **Q1 — Does C3's per-row lint suffixing apply to every row that carries a suffix (including rows minted by a match `in` expansion), or only to rows minted by a step expansion?** The two readings produce materially different lint output on models that step nothing.
  (a) **Uniform.** Every per-row finding publishes its row's identity in `Rule`/`Element`: the authored id plus the row's suffix elements, which equals the bare id for an unexpanded row. On models that step nothing, an `in`-expanded row's per-row finding changes from `retry` to `retry#fast`.
  (b) **Step rows only.** Only rows a step point minted are suffixed. `in`-expanded rows keep the bare id.
  For (a): C3 "lint names a ROW, because its subject is a row"; "Every surface that names a ROW publishes the suffixed identity … the shape `in` expansion already emits there"; "`lint` is the THIRD position, and it is a gap this record closes rather than inherits"; G Versioning, "the per-row finding `rule` slot's meaning narrows from \"the authored rule id\" to \"the row's identity\"". Reading (b) also needs a step-origin marker on the normalized row, which C1's closing sentence ("After normalization no surface distinguishes an expanded row from an authored literal row") argues against.
  For (b): C1 "this record changes no behaviour on a model that steps nothing" (fenced, though stated in the context of the `#`/duplicate guards); G Incremental adoption, "a model with no `{ step = <n> }` write value loads, lints, resolves and exports exactly as today"; C3's motivating case and S11 are step-only.
  No predecessor settles it: 0002/0006 fix `rule` as `row.RuleID`, and 0012 does not touch finding attribution.
  **Proceeding under (a).** It is the reading C3's fenced lint paragraph states, and it needs no new row state. The recovery join (REQ-38) is exact for any number of suffix elements (G Character encoding). Consequence to record at Phase 2 if taken: any committed test or golden that asserts a bare `rule`/`element` on a per-row finding over an `in`-expanded row changes, and REQ-47's "lints … exactly as today" is then true except for that attribution field. That exception is logged as a SPEC-DEFECT deviation against G/C1's "no behaviour on a model that steps nothing".

Candidates examined and resolved without a QUESTION:

1. **The post-lock 0012 landing** (the record's pre-0012 description of `Evaluator`, and its "does NOT wait on RDR 0012"). This is not ambiguous. JDR 0004 JD-1 and JD-3 bind 0030 and decide both the host package and the construction. Resolved as the first four ASSUMPTIONS.
2. **Lint cost wording** (C3 "`lint` not returning" against RISK/Performance "`complete=false`, reported as a finding"). The two describe different bands, and C3 is fenced. Resolved as an ASSUMPTION.
3. **Where "names the cell FIRST" sits** relative to `renderWrites`' existing prefix. S3 itself says the existing prefix stays and the cell is the new text. Resolved as an ASSUMPTION.
4. **Shim signature against the render step** (C1's shim takes a rendered `resolve.GuardAtom`, while Phase 2 says the shim "carr[ies] the canonicalizing set renderer"). Phase 2's own "each of the three callers renders its own `table.Atom` into that form at the call site" reconciles them. Resolved as an ASSUMPTION.
