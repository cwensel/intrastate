Model: claude-opus-5[1m]

# cove — chain-of-verification findings, cli/0024

Step 0 grounding swept all 32 `source-anchor` edges (every one `resolved:true`
at symbol level) **plus** the unanchored prose claims, the peer-RDR citations,
and the inverse check. Symbol existence was never taken as sufficient: each
claim was read against the actual source or the projected peer element.

Findings are ordered highest-value first. Each anchors to the element it lands
on. CONFIRMED claims are not reported (per the prompt), except where a
confirmation is needed to bound a finding.

---

## F1 — `A2`'s licensed-diff count is wrong, and Phase 3's adjacency count is a grep artifact (REFUTED, computed)

**Anchors**: `0024:A2`, Implementation Plan Phase 3.

`A2` states the cost leg as "**31** inline whole-payload assertion sites, 3
`resolveGolden` byte-identity goldens, the `shippedFixtureGolden` sweep".
Phase 3 states "the **28** inline whole-payload assertions in
`internal/cli/flow_demand_0011_test.go`, **7 of which** pin the
`"emit":{…},"next"` adjacency".

Both are wrong, and they disagree with each other inside one record.

Executed against the file (payload literals are Go raw strings concatenated
across source lines, so a line-oriented grep undercounts; the count below
normalizes the `` ` + \n` `` joins first):

```
whole-payload ok literals : 28
"emit": occurrences       : 28
emit,next adjacency (norm): 28
emit,next adjacency (raw) : 7     <- the RDR's "7"
resolveGolden entries     : 3
shippedResolveGoldens     : 25
```

- **28**, not 31, inline whole-payload assertion sites.
- **All 28** pin the `emit`/`next` adjacency, not 7. The "7" is exactly the
  count a line-oriented `rg '"emit":\{[^}]*\},"next"'` returns — the 7 sites
  where `emit` and `next` land on the same physical line. It is a measurement
  artifact, not a source fact.
- 3 + 25 = 28 reconciles: the `resolveGolden` goldens and the
  `shippedResolveGoldens` sweep **are** the 28, not additions to it. Both
  A2 and Phase 3 read as though they were additive.
- `shippedFixtureGolden` (`flow_demand_0011_test.go:667`) is the element
  **struct type**; `shippedResolveGoldens` (`:700`) is the **sweep variable**.
  A2 names the type where it means the sweep.

This is the one number that sizes the diff Phase 3 *authorizes by name*, and
Phase 3 closes with "Any assertion outside this list going red is a signal the
append did more than C4 licenses — investigate, do not update it." An
implementer holding "7" who sees 28 red adjacency assertions will read a
correct, licensed change as an over-broad one — the guard fires backwards.

The file's own comment (`:697`) states the true blast radius: "The goldens are
deliberately brittle: any change to the resolve envelope re-breaks every entry
at once."

Mitigation text in Risks ("the fixed insertion point … confines each inline
edit to one field at a known position") is right about the *shape* of each
edit and wrong about the *count*.

## F2 — `A4`'s evidence names the wrong function; the checker it needs is already exported (REFUTED)

**Anchors**: `0024:A4`, Existing Infrastructure Audit row "Value conformance".

`A4`'s Evidence says the reusable checker is
"`internal/table/load.go::conform` / `::conformKind` / `::conformDomain` …
unexported, receiver-free functions taking `(TagDecl, string)`".

Read against source:

```go
func ConformValue(decl TagDecl, member string) error   // load.go:793  EXPORTED
func conform(decl TagDecl, operator string, members []string) error // :806  THREE args
func conformKind(decl TagDecl, member string) error    // :833
func conformDomain(decl TagDecl, member string) error  // :849
```

- `conform` takes **three** arguments, not `(TagDecl, string)`. A4's stated
  signature does not describe it.
- The function that *does* have signature `(TagDecl, string)` and is
  value-level is **`ConformValue`** — already **exported**, already composing
  `conformKind` + `conformDomain`, and A4 never mentions it.
- Its doc comment is the exact warrant A4 is reaching for: "There is no
  operator here: a caller supplies a VALUE, not a predicate" and "A zero
  TagDecl conforms everything".
- It is already used across the package boundary for precisely the analogous
  case — `internal/cli/flow_input.go:349,367`, RDR 0020's `--tag` admission
  path, which is the same "declare-to-tighten" precedent 0024 cites from 0020.

Consequence: A4's operative conclusion (a value-level lexical check is
implementable without `internal/guard`) is **CONFIRMED and in fact stronger
than stated** — no adapter or extraction is needed. But the Infrastructure
Audit's Decision cell ("REUSE — verified separable … adapt via a narrow
helper") prescribes work that is already done, and points the implementer at
three unexported functions instead of the exported one.

**Downstream consequence C1 does not take**: reusing `ConformValue` *decides*
the question C1 defers. `conformKind`'s `int` arm is `strconv.Atoi`, which
admits `03`, `+5`, `-0`. C1 says "Whether a non-canonical `int` literal (`03`)
is admitted is a Resolve-level detail of the lexical rule" — under reuse it is
already answered (admitted), and answered by a function whose behavior 0003
owns. `conformDomain`'s `int` arm additionally applies `Min`/`Max`, which emit
declarations do not carry (C1 excludes `min/max`), so the reuse is partial and
the RDR does not say which arms it takes.

## F3 — C1's dual-form `domain` opens a strictness hole the rest of the schema does not have; the malformed-declaration list does not cover it (SILENCE, spiked)

**Anchor**: `0024:C1` (GRAMMAR, the `domain` clause and its refusal list).

C1 makes **two** `domain` forms normative — flat `domain = [...]` and the
sub-table `[emit.<key>.domain]` whose keys are disposition tokens. Every other
field in `sourceDoc` is concretely typed, and strictness
(`DisallowUnknownFields`) is what catches authoring typos for free.

A single struct field cannot decode both a TOML array and a TOML table, so the
emit decl's `domain` must be `any` (or carry a custom unmarshaller). Spiked
against the repo's actual decoder (`pelletier/go-toml/v2 v2.2.4`):

```
flat        err=<nil>  domain=[]interface{}{"a","b"}
subtable    err=<nil>  domain=map[string]interface{}{"route":[...], "stop":[...]}
```

So the dual form **is** decodable — C1's grammar is implementable, and the
Risks entry that worries it "proves awkward to … decode against real TOML
tooling" is over-pessimistic on that point.

But strictness then stops descending, exactly as `sourceModel.Metadata`
documents at `source.go:39-42` ("a free-form map so strict decoding descends
no further"). Probed:

```
typo at decl level   (domaim = [...])              err=strict mode: ... REFUSED
scalar under domain  (route = "not-an-array")      err=<nil>   ADMITTED
deep junk under domain ([emit.next.domain.route.oops]) err=<nil>   ADMITTED
```

C1's `malformed_emit_declaration` list enumerates: unknown `kind`; `enum` with
no domain; empty domain; empty-string member; duplicate member; `domain` on a
non-enum kind; empty-string disposition token. It does **not** cover:

- a disposition whose value is not an array of strings (`route = "x"`,
  `route = 3`, `route = {…}`);
- nesting below the disposition level (`[emit.next.domain.route.oops]`);
- a disposition sub-table that is empty (`[emit.next.domain]` with no keys) —
  distinct from "an empty domain", which C1 does list for the flat form.

These are the cases the decoder catches for free everywhere else in the schema
and cannot catch here. The RDR is silent on the trade, so the whole-model
strictness C2 sells ("a typo'd KEY is indistinguishable from an intentionally
undeclared one") is weaker one level down than a reader would expect: a typo'd
**disposition sub-key shape** is admitted silently.

This is the highest-value silence: it lands on the grammar C1 locks, and the
refusal list is where it must be fixed.

## F4 — C2's fail-fast argument imports two peer citations from a channel their own text distinguishes (REFUTED as support; conclusion survives)

**Anchor**: `0024:C2` ("Reporting is FAIL-FAST …" and "Why fail-fast is the
correct tier behavior").

C2 grounds fail-fast on `0008:C4` and on `0018:C1`/`0018:ALT1`. Both quotes are
verbatim. Both are scoped to a channel this RDR is not on.

- **`0008:C4`** is scoped to **one kernel-entry predicate** over
  `resolve.Input` (the reserved-`recognized`-key check across three sequences),
  on the Go-error path, where C4 itself says "A breach is a producer programmer
  mistake — not table data". C2 introduces it as "the load tier's cardinality
  is already contracted: `0008:C4` fixes it". It does not: it fixes one
  predicate's cardinality on the programmer-mistake channel.
- **`0018:ALT1`** draws the line explicitly, and draws it *against* C2's use:
  "The prior art misfits the channel: **aggregation serves user-facing
  declarative validation; this is the programmer-mistake path**, where
  `0008:C4` already rules a producer 'is not owed an exhaustive list'."

An emit-declaration defect is an **author** mistake in model data, surfaced to
a human through `intrastate lint` — i.e. user-facing declarative validation,
the side of 0018's own distinction where it says aggregation *does* belong. So
the citation is not merely stretched; read straight it is mild counter-evidence.

- **`0018:C1`** is also cited for "the house answer is PRECEDENCE, not
  aggregation", immediately before C2 declines to fix any precedence ("order
  stays unspecified"). 0018:C1's entire content is that the winner is *named
  and observable* ("reordering the entry checks is a breaking change to this
  contract"). C2 takes 0018's rejection-of-aggregation half and drops the half
  that gives precedence meaning.

**The conclusion is not refuted and should not change.** C2's own second
ground is independently sufficient and verified verbatim in source —
`internal/table/load.go`'s `Load` doc (`:16-29`): "the refusal is singular:
load is fail-fast and returns one categorized error, never a list" and "the
order in which independent defects are checked is deliberately unspecified".
Together with the absence of a comparator (the real load/lint asymmetry C2
correctly identifies), that carries the whole argument. The fix is to lead
with `load.go`'s contract and demote 0008/0018 to channel-analogous colour, or
drop them.

Related, minor: C2 cites the load.go text as "`internal/table/load.go`'s
package contract". It is the **`Load` function** doc; the package doc lives in
`internal/table/doc.go`.

## F5 — `0020` is Draft; "the house already adjudicated this shape once" is not available from it (REFUTED)

**Anchor**: Decision Rationale, "House precedent (verified at Resolve)".

The paragraph's quotes from `0020:C1` and `0020:BR1` are verbatim and
`0020:BR1` does name "the 0024 shape" as compatible layered work. But
`0020`'s `Status` is **`Draft`** (projector-confirmed), not Final.

The paragraph is headed "House precedent" and asserts "the house already
adjudicated this shape once" and "`0020` decided the sibling question for
`--tag` keys". A Draft record has adjudicated and decided nothing; it is a
proposal that may still move. 0024 tracks peer status carefully elsewhere
(Phase 3: "0023 is Final and unimplemented"), which makes the omission here
read as an oversight rather than a convention.

The paragraph's load-bearing evidence — the **loader census** ("tag keys, tag
kinds, operators, the outcome alphabet, dump columns, reserved keys, and
schema fields all refuse unknown members at load; `[rule.emit]` is the sole
admitted-unconditionally vocabulary") — is independent of 0020's status and
survives intact, as does the `docs/model-authoring.md` quote. Only the
"already adjudicated" framing needs to weaken to "0020 proposes the same shape
for the sibling `--tag` question and names this RDR as compatible".

Also in this paragraph: 0010 is described in the record elsewhere as a locked
Final contract; its actual `Status` is **`Implemented`** (stronger, not
weaker — worth correcting only for precision).

## F6 — a 0010 spec test the RDR does not mention constrains the emit refusal arm (SILENCE)

**Anchor**: Implementation Plan Phase 1 (the amended-comments paragraph, which
names the 0010 tests it keeps passing).

Phase 1 names `internal/table/emit_0010_test.go::TestReq27_AnEmitKeyIsNotATagKey`
as the 0010 spec test that "keeps passing unchanged". It does not mention
`internal/table/emit_shape_0010_test.go::TestReq21_NoHandWrittenTypeCheckDiscriminatesTheEmitRefusal`,
whose failure message is:

> "a wrong-typed emit value refuses %q while another wrong-typed scalar
> refuses %q; **the emit arm is a hand-written type check rather than the
> decoder's**"

It asserts a wrong-typed emit value refuses with `CatMalformedTOML`,
indistinguishable from any other wrong-typed scalar. 0024 adds exactly a
hand-written kind check over emit values.

**They do not actually collide** — verified: `sourceRule.Emit` is
`map[string]string` (`source.go:102`), so `verdict = 42` refuses at the decoder
before any 0024 check is reached; 0024's check is lexical over an
already-decoded *string*. TestReq21 keeps passing.

But the RDR should say so, for two reasons. First, the test's name and message
are a standing tripwire aimed at this exact move, and an implementer will hit
it; Phase 1 pre-clears the *other* 0010 test by name and leaves this one
unexplained. Second, the interaction is genuinely counter-intuitive at the
authoring surface: under `kind = "int"`, an author who writes `count = 42`
gets `malformed_toml` from the decoder, never `emit_value_out_of_domain`, and
must write `count = "42"`. C1 states every kind check is "a LEXICAL check on
the authored string" but never says the authored string is *required to be a
TOML string* — which is where `kind = "int"` and `kind = "bool"` become
surprising.

## F7 — C2 does not say which loader step the rule-side checks occupy (SPEC-UNDER)

**Anchor**: `0024:C2` ("MUST refuse, before yielding rows").

The loader is a fixed 11-step slice (`internal/table/load.go::run`, `:78-98`),
returning on the first step error:

```
loadModelHeader, loadOutcomes, loadTags, loadAccessors, loadDump,
loadContexts, loadInitial, loadTerminal, normalizeRules,
checkClassAgreement, checkAccessorBindings
```

"Before yielding rows" is satisfiable — rows are minted in `normalizeRules`
(step 9) — so C2's clause is achievable, CONFIRMED.

But the declaration-side check (`malformed_emit_declaration`) and the
rule-side checks (`unknown_emit_key`, `emit_value_out_of_domain`) have
different natural homes, and C2 fixes neither:

- placed as a step beside `loadTags` (before step 9), the rule-side checks
  read `l.src.Rule[].Emit` — the **source** rules;
- placed after `normalizeRules`, they read normalized rows — but then they run
  after rows are yielded, contradicting C2's own clause.

Two steps already sit after `normalizeRules` precisely because they need
normalized rows (`checkClassAgreement`, `checkAccessorBindings`), so the
pattern for both placements exists and the choice is real. The RDR names the
pipeline but not the position, and the position determines which representation
the check reads. Given C2's "before yielding rows" the source-rule reading is
forced — worth stating rather than leaving to the implementer.

## F8 — three citation-precision defects (minor, but each is a claim about a peer)

**Anchors**: References; `0024:C2`; Phase 3.

1. **`0006:C13` is a dead citation.** It appears in References ("`0006:C3`,
   `0006:C13`, `0006:C17` (finding tiers)") and is used nowhere in the
   record's argument. C3 and C17 both do real work; C13 (the blocking-finding
   payload contract) does not.
2. **References omits every peer the argument actually leans on.** Cited
   in-argument but absent from the References list: `0003` (C19),
   `0008` (C4), `0009` (C5), `0011`, `0018` (C1, ALT1), `0020` (C1, BR1),
   `0023` (C2, A4), and `JDR 0002 §D1`. The list names only 0002/0005/0006/0010.
3. **The `0009` citation understates what made the order total.** C2 says
   0009's report "aggregates because `internal/resolve/resolve.go::compareRefs`
   orders `RowRef` lexicographically on `(RuleID, SourceLocator)`". `0009:C5`
   says explicitly that `compareRefs` on that pair leaves distinct rows
   comparing **equal**, and the total order comes from *collapsing* equal
   identities into one reported error with a `Count`. Since C2's whole argument
   turns on "the project accumulates findings precisely where a TOTAL ORDER
   exists", presenting a partial order as the total one weakens the strongest
   part of the argument. `compareRefs` at `resolve.go:663` confirms the
   two-field comparison.

Also: `0011` is characterized in A2 and Phase 3 as "RDR 0011's guard against
exactly this class of change". 0011 decides `flow next` candidate
match-conditioning; the full-payload goldens in `flow_demand_0011_test.go` are
0011's *own adversarial oracle* for its demand-set change, not a payload-shape
contract 0011 owns. Practical effect is favorable (updating them is a fixture
edit, not a peer-spec amendment) and Phase 3's conservative "authorized HERE by
name" handling is right — only the characterization overstates.

## F9 — `A5`'s sub-counts do not match the models it cites (minor; conclusion CONFIRMED)

**Anchor**: `0024:A5`.

Re-counted against the live consumer models:

| A5 claims | actual |
|---|---|
| `rdr-status.toml`: 54 `[rule.emit]` blocks | **54** ✓ |
| `next`: 22 distinct values, all fixed literals | **22** ✓, all literal ✓ |
| "16 route-like and 6 `stopped:`-like" | 14 `/rdr-*` + 6 `stopped:` + `none` + `resolve:lens` |
| `rdr-write.toml`: 30 `[rule.emit]` blocks | **29** |
| `op`: 10 distinct values | **10** ✓ (claim, demote, lock, readme-add, readme-flip / none, stopped:×4) |

A5's load-bearing content is **correct**: the vocabularies are closed and
enumerable, and the partition is "three-valued in practice, not binary" —
which A5 states correctly two sentences later, and which is the live evidence
against `ALT1`. The "16 route-like" sub-count is what silently folds `none`
and `resolve:lens` into the routes, i.e. it commits on a small scale the exact
two-valued error A5 is arguing against. Block count is 29, not 30.

Not a refutation; the numbers should match the models they cite.

---

## Step 0 — CONFIRMED, reported only where a finding depends on them

- All 32 `source-anchor` edges resolve `true`; each was additionally read.
- `sourceDoc` declares 11 fields, **no `emit`** — `A1` CONFIRMED verbatim,
  including the `decodeStrict` → `CatUnknownSchemaField` = `"unknown_schema_field"`
  refusal leg and its unconditional call from pass-2.
- **Inverse check** (the RDR adds a new discriminator over emit): searched,
  **none exists**. `internal/graphlint` never reads `Row.Emit` (zero
  non-test matches). The only four non-test `.Emit` sites are carry-through
  or render (`normalize.go:457,697`, `dump.go:101`, `flow_resolve.go:239`).
  No emit category in `category.go`, no emit advisory in `advisory.go`. The
  sole existing constraint is shape-not-domain via `map[string]string`.
- `resolvePayload` has **14 fields / 13 wire keys**, one `omitempty`
  (`escape_class`) — `A2`'s and Testing Strategy's 14→15 / 13→14 arithmetic
  CONFIRMED. `TestReq39`'s `want` list matches the Testing Strategy's stated
  pre-change order exactly.
- `TestReq43` asserts **relative** indices — Phase 3's "needs no edit"
  CONFIRMED.
- `emitMap` (`flow_resolve.go:285`) + `rowByID` (`:294`) — single row, by
  `Plan.RuleID`, no default-row and no merge path. C4's claim CONFIRMED.
  Unconditional allocation yields `{}` not `null`, as C4 requires.
- `flow next`'s `nextPayload` carries no `emit` — C4's "carries no
  `dispositions`" CONFIRMED with 0010:C4's stated reason.
- Text mode: `respond/text.go::flatten` is fully generic (recursive
  `map[string]any` arm + empty-map `(none)` arm, no key-name logic), and
  `TestReq45` already attests the behavior for `emit`. C4's
  `dispositions.<key>: <token>` / `dispositions: (none)` rendering needs **no**
  renderer change. CONFIRMED.
- Key ordering: `flow_resolve.go:282` already documents "the payload encoder
  emits a map's keys in byte order" — C4's byte-order clause is inherited,
  not new.
- Advisory tier: `internal/graphlint/taxonomy.go:37` — "The four advisory
  finding codes. The tier is CLOSED at these four". `0006:C17` CONFIRMED in
  code as well as in the record, so the declined Meyer condition is grounded.
- `internal/cli/lint.go:176-191` loads via `LoadWithAdvisories` and returns
  `respond.Fail` with `loadFindings(...)` before `graphlint.Run` at `:193` —
  C2's lint-surfacing claim CONFIRMED. Two precision notes: `loadFindings` is
  defined in `internal/cli/flow_input.go:185` (lint.go is a caller), and
  "blocking" is structural (the `Fail` path), not a severity field.
- Emit exists only on rules (`sourceContext` carries none), so C2's "any
  `[rule.emit]` key of any rule, ordinary or escape" is complete coverage.
- `emitSequence` key-sorts and `expand` copies `base.Emit` — C2's "proving the
  AUTHORED form is proving the executed form" CONFIRMED.
- `declaredKinds` = `["enum","bool","int","set","scalar"]`. C1 says it reuses
  0003's spellings "verbatim" and lists them `enum | bool | int | scalar`
  (excluding `set`, correctly and with a stated reason). Same tokens; only the
  ordering differs from the source literal. Not a finding.
- Comment amendments (Phase 1): `model.go::EmitValue` and
  `normalize.go::emitSequence` quotes CONFIRMED verbatim. `dump.go::renderEmit`
  reads "An emit key has neither: it is undeclared, so there is no kind to
  consult" — CONFIRMED, but "neither" refers back to `renderValue`'s "declared
  kind and member count", and that antecedent stays true under 0024 (an emit
  key still has no *tag* kind and its value is still one string). This is the
  least stale of the three comments; Phase 1's narrowing is defensible but
  overstates its staleness.
- `JDR 0002 §D1` **does name `dispositions` by name** and assigns it PLAN
  ("in **0024** — `dispositions` is a PLAN-group field"), corroborated by
  `0023:C2` and `0023:A4` (Verified). 0023's Status is **Final**. The forward
  reference in Phase 3 CONFIRMED — and `0023:A4` independently names the same
  two `decision_table_0010_test.go` assertions 0024's A2 names, so the two
  records agree on the licensed diff.
- `0002:C22`, `0002:C24` (all three legs), `0005:C1` (including 0024's own
  citation correction, which is right), `0006:C3`, `0006:C17`, all of
  `0010:C3`, all of `0010:C4`, `0010:A6`(a)–(d), `0003:C19`/`0006:C16` and
  `engine.go::sortFindings`, `0020:C1`/`BR1` quoted text — all CONFIRMED
  verbatim.
- `0006:C1`'s quoted span is verbatim; its lead-in paraphrases "a lint result"
  as "a design-time proof", a mild widening that the ALT3 inference survives.

---

## Dispositions

| # | Disposition | Origin ledger entry | Section touched |
| --- | --- | --- | --- |
| F1 | fixed | F1 (counting refutation, executed) | `A2`; Phase 3; Risks |
| F2 | fixed | F2 (`ConformValue` misattribution) | `A4`; `C1`; Existing Infrastructure Audit |
| F3 | fixed | F3 (strictness hole, spiked) | `C1` refusal list; Risks; Testing Strategy scenario 1; new `A7` |
| F4 | fixed | F4 (channel-mismatched citations) | `C2` both fail-fast paragraphs |
| F5 | fixed | F5 (0020 is Draft) | Decision Rationale, house-precedent paragraph |
| F6 | fixed | F6 (TestReq21 unmentioned) | Phase 1 |
| F7 | fixed | F7 (loader step unspecified) | `C2` opening clause |
| F8 | fixed | F8 (citation precision) | References; `C2`'s 0009 citation; MVV step 2 |
| F9 | fixed | F9 (`A5` sub-counts) | `A5` |

No entry dismissed, none charted to a successor: every finding landed inside
this RDR's existing scope — its own contracts, assumptions, and citations. No
finding required a tiebreaker; each collapsed on the evidence once the source
or the projected peer element was read.

Ledger closed: 9 of 9 dispositioned, 0 open. Converged at iteration 1.

## Needs (re)verification — carried to Stage 6

1. **`A7` (new, Status `Pending`, Method Spike)** — booked by F3. The dual-form
   decode and the bounded strictness loss are spiked
   (`evidence/spikes/c1-dual-form-domain.md`); what is unverified is the
   CLOSURE claim, that C1's five added arms are the complete set of shapes the
   decoder stops refusing. The spike probed three of them; the rest were
   derived from the type structure. Verification plan is on the assumption.
2. **`A4` — no status change, evidence corrected.** F2 refuted its evidence
   POINTER, not its verdict: the operative claim (implementable without
   importing `internal/guard`) is now better supported, since the exported
   `ConformValue` has exactly the needed signature and already crosses the
   package boundary for RDR 0020's analogous `--tag` path. Stays `Verified`.
3. **New normative claim in `C1`, grounded not booked** — that an emit value
   is always a TOML string, so `kind = "int"` takes `count = "42"` and a bare
   `count = 42` refuses `malformed_toml` at the decoder. Read directly from
   `internal/table/source.go:102` (`Emit map[string]string`) and from
   `0010:C3`; no assumption booked because the source is dispositive.
4. **New normative consequence in `C1`, grounded not booked** — `03`, `+5`,
   `-0` are admitted under `kind = "int"`, inherited from `conformKind`'s
   `strconv.Atoi`. Read from `internal/table/load.go:833`. This CLOSES a
   question C1 previously deferred to Resolve; it is now decided in the draft.

No assumption was flipped from `Verified` back to `Pending`; no fix invalidated
a previously verified claim.

## Mini-checks

Cue read performed (this RDR's first lens pass). Three of five fired and their
tables are written into the draft under Technical Design → Mini-checks:
`disposition`, `fidelity`, `trace`. The desk trace walked the MVV stepwise
against all four contracts and returned **no CONTRADICTION row**.

Not fired: `source-authority census` (no fallback, derived, or propagated path;
no source-of-truth contest — the declaration is the sole authority and C3
carries it to exactly one reader) and `test-discriminability` (every MVV oracle
asserts a named value — category slugs, a 14-key wire order, `{}` vs `null` —
and step 1 is the negative control).
