Model: claude-opus-5

# Reconciled critique diff — RDR 0010, stateless decision tables

Two independent dual-model critique passes over
`docs/rdr/0010-stateless-decision-tables.md` (Status: Draft, 1542 lines), diffed by
**RDR passage anchor**. The `C-N` ids in each input file are per-file and do not
correspond across passes; the `D-N` ids below are the reconciled, stable keys.

| Pass | File | Model stamp | Rows |
| --- | --- | --- | --- |
| **A** | `evidence/critique/critique.md` | `claude-opus-5` | 13 (C-1…C-13) |
| **B** | `evidence/critique/critique-modelB.md` | `claude-sonnet-5` | 8 (C-1…C-8) |

Grounding was performed by this pass against shipped `main` at `65cf2ac` in
`the repo root` (`internal/table`, `internal/graphlint`,
`internal/cli`, `internal/guard`) and against the RDR's own decided text via
`rdr inspect`. Pass A's source claims were re-derived from the files, not taken on trust.

---

## Reconciled ledger

| ID | Anchor(s) | Raised by | Defect | Agreement | Grounded verdict | Severity |
| --- | --- | --- | --- | --- | --- | --- |
| **D-1** | `0010:C5` ¶2, `0010:A13` | A (A:C-1) | C5's new zero-dimension emission of `graph-unprovable-coverage` carries no `reason`, but 0006 fixes `reason` as a closed 3-value set and a shipped conformance test asserts every emission of that code carries one from it; the RDR never mentions `reason` normatively | A-ONLY | **REAL** — `taxonomy.go:60-62`, `:88-92`, `:104`; `findings_0006_test.go:517-542` | **BLOCKING** |
| **D-2** | `0010:A4`, `0010:§testing-strategy` Done clause | A (A:C-2) | The `[dump]` census counted anchored `.toml` under `testdata/` only; `dump_test.go` hard-codes the ten-column vocabulary as a Go literal under two `reflect.DeepEqual` assertions, and the Done clause licenses no edit of that shape | A-ONLY | **REAL** — `dump_test.go:298-311`; 5 Go files under `internal/table` carry `[dump]`, none in A4's accounting | **BLOCKING** |
| **D-3** | `0010:A10` (silent bucket), `0010:C5` ¶3 | A (A:C-3) | A10's stated silence mechanisms for `graph-unreachable-rule` and `graph-always-present-owned` are wrong: both functions return early on `len(a.model.Initial) == 0`, not for the reasons A10 gives | A-ONLY | **REAL** — `analysis.go:515-520` (`checkUnreachableRules`), `analysis.go:233-238` (`checkAlwaysPresentOwned`) | SUBSTANTIVE |
| **D-4** | `0010:§technical-design` item 1, `0010:§authority`, `0010:C5` | A (A:C-4) | "every site that today asks `len(Initial) == 0` … reads that accessor, so there is exactly one place the question is answered" is false as written: C5 explicitly *augments* rather than replaces the `len(Initial)` test, and two of the four shipped sites must keep the old key | A-ONLY | **REAL** on the contradiction, **WRONG** on the count — A says "five sites"; there are **four** `len(Initial) == 0` sites: `reach.go:91`, `analysis.go:113`, `:233`, `:515`. `reach.go:99` is a `range`, not a test. A's own AT-4 says "4 in graphlint alone", contradicting its own row | SUBSTANTIVE |
| **D-5** | `0010:C4`, `internal/cli/flow_resolve.go::resolvePayload` | A (A:C-5) | C4 fixes no field position for `emit` in the payload struct, while C3 was scrupulous about the dump column's position ("appended last, after `escape`") with a stated rationale | A-ONLY | **REAL** on the omission, **WRONG** on the enumeration — A says C4 "omits ... `model`, `revision`, `observed`, `outcome`, `gates`, `escaped`, `escape_class`" and calls the struct eleven fields (AT-5); the struct carries **13** fields at `flow_resolve.go:31-45`. C4 does name five (`next`, `writes`, `clear`, `owned`, `readers`); position is genuinely unspecified | SUBSTANTIVE |
| **D-6** | `0010:C5` ¶2, `0010:A13` | A (A:C-6) | The fence lands in the wrong branch: `checkCoverage`'s `len(dims)==0` arm calls `emitCoverageArms` and **returns**, above `emitUnprovableDimensions` / `emitWithholdings`; `emitStructurallyUnprovable` is a loop over `guard.Dimensions` and is unreachable in the zero-dim case. A13 points at a site that cannot fire, and C5 fixes no precedence against `graph-coverage-closed-by-escape` | A-ONLY | **REAL** — `coverage.go:33-43` (zero-dim early return), `:147-149` (loop body), `:65-66`/`:87-88` (both other emitters below the return), `:351` (`bareEscapeFor` inside `emitCoverageArms`) | **BLOCKING** |
| **D-7** | `0010:A6`, `0010:C3` ("values MUST be strings"), `0010:BR3` | **BOTH** (A:C-7 §3, B:C-3 §3) | `emit` is string-only and A6's own cited evidence (`Next: /rdr-prelock 0046 critique`) is a verb-plus-arguments answer; A's sharper form: `0046` is caller-supplied and unknowable at authoring time, so the literal block is unwritable without a template — reintroducing the consumer-invented format the Problem Statement exists to abolish | **CONVERGED** | **UNGROUNDED** (judgment) on sufficiency; **REAL** on the textual facts — A6 at RDR:236-252 quotes the example and files it as a limit; BR3 rejects `any`-typed values at RDR:1022-1023 | SUBSTANTIVE |
| **D-8** | `0010:A12` (Pending), `0010:A13` (Pending), `0010:§assumption-verification`, `0010:§trace` row 2″ | **BOTH** (A:C-8 §2, B:C-8 §2) | Two `Pending` assumptions carry C5's load-bearing consequences, while C5 states the fence as an unconditional `MUST` and §risks/§failure-modes treat it as settled — which the RDR's own Finalization Gate status-consistency clause forbids | **CONVERGED** | **REAL** — A12/A13 `Status: Pending` at RDR:412/432; gate clause "no assumption marked `Pending` or `Unverified` may have settled-fact prose elsewhere in the RDR depending on it" at RDR:1457-1460; Trace row 2″ "**unwitnessed — A13 is Pending**" | **BLOCKING** |
| **D-9** | `0010:C1` vs `0010:§approach`, `0010:§decision-rationale` scorecard, `0010:ALT2` | A (A:C-9) | Flat contradiction: C1 says the agreement check is "**one-directional**"; three narrative passages say bidirectional — including the scorecard cell the decision was made on | A-ONLY | **REAL** — C1 at RDR:531-533 "one-directional"; §approach RDR:457 "in either direction" and RDR:483 "a load failure both ways"; scorecard "Silent misclassification … refused at load both ways"; ALT2 "a load-time refusal in both directions" | **BLOCKING** |
| **D-10** | `0010:C3` (join soundness / `rowByID`), `0010:§testing-strategy` S3–S5 | A (A:C-10) | `expand` builds each row from an explicit seed literal, not a struct copy, so a new `Row` field is dropped by default; C3's join-soundness argument is phrased as already-holding, and no scenario asserts an emit block on an expanding (`in`-atom) rule | A-ONLY | **REAL** on the code, **PARTLY WRONG** on the citation — `normalize.go:646-655` is an 8-field seed literal (`ModelID, RuleID, SourceLocator, Source, Gate, RequiresOwned, Escape, setKeys`); A calls it "nine named fields" then lists eight. And C3 *does* already say `Emit` "MUST be carried through … `expand`, **which builds each row from the seed literal rather than copying `base`**" (RDR:628-630) — the RDR is not silent on the mechanism, only on the test | SUBSTANTIVE |
| **D-11** | `0010:§testing-strategy` Done clause vs `0010:A4` | A (A:C-11) | The licensed-diff rule is line-shaped (` emit=[…]` cell or `"emit":` member) but the 103 required edits are to `order = [...]` lines, which are neither shape; A4 asserts they are licensed without amending the rule | A-ONLY | **REAL** — Done clause at RDR:1295-1303; A4's counter-assertion at RDR:1301-1303; 103 anchored `[dump]` files confirmed by re-running A4's own grep | SUBSTANTIVE |
| **D-12** | `0010:A9` "Consequence carried, not restated", `0010:C5` (mitigation parity) | **BOTH** (A:C-12, B:C-7 + B:C-4) | The N-escape-rows-per-outcome consequence is discovered, then routed to Phase 4 documentation as its only mitigation — no contract clause and no distinguishing lint signal, while the structurally identical match-atom mistake gets a `MUST`-emit fence in C5 | **CONVERGED** | **REAL** on the asymmetry — A9 at RDR:318-322 ends "which the Phase 4 authoring docs must say"; C5 gives the sibling hazard a `MUST`. **UNGROUNDED** on whether lint is the right home | SUBSTANTIVE |
| **D-13** | `0010:§problem-statement`, `0010:§approach`, `0010:§technical-design` data-flow vs `0010:§key-discoveries`, `0010:C5` ¶2 | A (A:C-13) | Three passages, including the sentence stating why the RDR exists, say coverage runs "over **observed** dimensions"; coverage is proven over **guard** dimensions, and provenance is orthogonal to block. The top of the file teaches the exact authoring error C5 fences | A-ONLY | **REAL** — `internal/guard/product.go:130-141` `Dimensions` collects from `guardAtoms(row)` only; §problem-statement RDR:126, §approach RDR:471, data-flow RDR:521. Note C5 itself repeats it ("MUST run unchanged over the authored observed dimensions"), so the count is **four** passages, not three | **BLOCKING** |
| **D-14** | `0010:§technical-design` (guard-vs-match discriminator), `0010:C5`, `0010:§illustrative-code` | B (B:C-1 §1) | Authors will discriminate with `[rule.match.<key>]` because that is how every other rule engine reads; only C5's (unverified) fence stops it, and `graph-unprovable-coverage`'s message names the symptom, not the remedy — no contract requires the message to be actionable | B-ONLY | **UNGROUNDED** (authoring-UX judgment). Its mechanical premise is **REAL** — `coverage.go:33-43`: a zero-dim group closes vacuously today. Overlaps D-1/D-6 on the fence, but this row is about *message actionability*, a distinct passage-level defect neither A row raises | SUBSTANTIVE |
| **D-15** | `0010:§authority` (class read at four sites), `0010:C1` | B (B:C-2 §1) | "One accessor, four readers" is a coding convention with no test, linter rule, or code-level enforcement; nothing stops a fifth site inferring class from `len(owned) == 0` | B-ONLY | **REAL** on the absence — §authority names four readers with a `len(owned) == 0` "sibling arm" column and no enforcement mechanism; no self-lint over the intrastate tree exists in the Testing Strategy. Adjacent to D-4 (which is the *contradiction*); this row is the *enforcement gap* | SUBSTANTIVE |
| **D-16** | `0010:C2`, `0010:A9` (the "otherwise" idiom) | B (B:C-4) | The escape-row "otherwise" idiom is authoring convention documented only in Phase 4; an author writing an overlapping catch-all ordinary row gets `flow-ambiguous-match` naming the symptom, not the fix | B-ONLY | **REAL** on the text — A9's "If wrong" arm at RDR:322-323 names exactly this outcome and accepts it. **UNGROUNDED** on severity | MINOR |
| **D-17** | `0010:C4` ("`--outcome` remains required in both classes"), `0010:BR5` | B (B:C-5) | A decision table has no state to select an outcome from either, so the single-outcome ergonomic gap bites hardest on exactly the class this RDR exists to serve — yet it is parked as "a plain kata" with no tracking commitment | B-ONLY | **REAL** on the text — C4 at RDR:673; BR5 at RDR:1027-1028 "parked as a plain kata at triage", no kata id cited. **UNGROUNDED** on whether that is wrong | MINOR |
| **D-18** | `0010:A5`, `0010:§decision-rationale` | **BOTH** (B:C-6 §3, A §3 runner-up) | Exact-one hit policy is verified against one consumer with a self-declared "⚠ no corpus coverage" flag, and is locked as *the* policy with no versioning or extension point for first-hit/priority | **CONVERGED** | **REAL** on the flag — A5 at RDR:224-235 carries the ⚠ verbatim and scopes itself "sufficient for the motivating consumer, not a general survey". **UNGROUNDED** on the prediction | SUBSTANTIVE |

**No CONFLICT rows.** The two passes never assert incompatible things about the same
passage. Where they touch the same anchor they agree on direction and differ only in
depth (B argues from the record's own text; A additionally cites shipped source).

---

## Grounding notes, per finding

All paths relative to `the repo root`; line numbers on `main` at `65cf2ac`.

### D-1 — the closed `reason` set — **REAL**

`internal/graphlint/taxonomy.go:58-62`:

```go
// The closed, append-only `reason` set `graph-unprovable-coverage`
// carries (Technical Design, reason table). It is carried on that code
// alone.
const (
	ReasonDimensionNotFinite = "dimension-not-finite"
	ReasonTagNotSingleValued = "tag-not-single-valued"
	ReasonRowCanRefuse       = "row-can-refuse"
)
```

`taxonomy.go:87-92` binds them into `var reasons`; `taxonomy.go:103-104` exports
`Reasons()`. The conformance test exists and asserts both halves —
`internal/graphlint/findings_0006_test.go:517-542`: `slices.Equal` against the
sorted three-element literal, **and** a walk over every emitted finding asserting a
`CodeUnprovableCoverage` finding's `Reason` is in the set while every other code's
`Reason` is `""`. `0010` never uses the word `reason` normatively (C5 at RDR:688-738
does not mention it). A's claim holds exactly as stated.

### D-2 — the dump census — **REAL**

`internal/table/dump_test.go:298-311`:

```go
func TestReq95_DumpColumnVocabularyIsClosedAndVerbatim(t *testing.T) {
	want := []string{
		"identity", "source", "kind", "outcome",
		"atoms", "next", "writes", "requires_owned", "gate", "escape",
	}
	got := table.DumpColumns()
	if !reflect.DeepEqual(got, want) { … }
	m := mustLoad(t, rdrFixture)
	if !reflect.DeepEqual(m.DumpOrder, want) { … }
```

A's quoted comment is verbatim at `dump_test.go:295-297`. Re-running A4's own census
confirms 103 anchored `.toml` files; running it unanchored across Go source returns
`internal/table/{dump.go, dump_test.go, format_test.go, roundtrip_test.go, source.go}`
— five files absent from A4's accounting (RDR:196-218). A4's arithmetic "176 repo-wide,
less the 103 = 73, all under `docs/rdr/0002-*/evidence/`" no longer holds even in its own
terms: unanchored repo-wide is now 229.

### D-3 — A10's silence mechanisms — **REAL**

`internal/graphlint/analysis.go:514-520`:

```go
func (a *analysis) checkUnreachableRules() {
	if len(a.model.Initial) == 0 {
		// With no root there is no reachable set to filter against …
		return
	}
```

`internal/graphlint/analysis.go:232-238` is the same shape for
`checkAlwaysPresentOwned`. A10's evidence (RDR:350-355) attributes the first to
"C5's single ∅ node satisfies every context vacuously" and the second to filtering "on
`ProvenanceOwned` and find none". Both functions never reach that logic under
`len(Initial) == 0`. A's characterisation — the stated reason is an accident that
happens to agree — is accurate.

### D-4 — "exactly one place" — **REAL** (contradiction), **WRONG** (count)

Exhaustive sweep of `len(…Initial) == 0` across `internal/graphlint`,
`internal/table`, `internal/cli` (non-test) returns **four** sites:

- `internal/graphlint/reach.go:91` — `reach`, the traversal root guard
- `internal/graphlint/analysis.go:113` — `checkDanglingEdge`, missing-root arm
- `internal/graphlint/analysis.go:233` — `checkAlwaysPresentOwned`, early return
- `internal/graphlint/analysis.go:515` — `checkUnreachableRules`, early return

Pass A's ledger row says "**five** … (`reach.go:91`, `analysis.go:113`, `:233`, `:515`,
plus `reach.go:99`'s iteration)". `reach.go:99` is `for _, tv := range m.Initial` — a
range, not a `== 0` test — so the fifth is a miscount. A's own AT-4 says "there are 4 in
graphlint alone", so the row and the AT disagree with each other. The **substance**
holds: §technical-design item 1 (RDR:496-499) claims a single answer point, C5
(RDR:690-696) requires `reach`/`checkDanglingEdge` to *augment* rather than replace
`len(Initial)`, and C5 ¶3 lists `always-present-owned` as "vacuous by construction",
leaving two sites permanently on the old key. §authority's reader list names four
readers, but two of them (`normalizeRule`, `checkCoverage` zero-dim arm) are not
`len(Initial)` sites at all, and the two real ones it omits are the two that must stay.

### D-5 — payload field position — **REAL** (omission), **WRONG** (enumeration)

`internal/cli/flow_resolve.go:30-45` — `resolvePayload` carries **13** fields in this
declaration order, which is Go's JSON emission order:

`Model, Revision, Observed, Owned, Readers, Outcome, Rule, Gates, Next, Writes, Clear, Escaped, EscapeClass`

C4 (RDR:656-683) names `emit`, `next`, `writes`, `clear`, `owned`, `readers` and fixes
no position for `emit`. A's row asserts C4 "enumerates five payload fields and omits the
four the struct actually carries alongside them" and then lists **seven** names; AT-5
says "Eleven fields. C4 names five." Neither count matches the struct. The load-bearing
claim — C3 fixes the dump column's position with a stated rationale (RDR:634-639
"**appended last**, after `escape`") while C4 says nothing about the same question one
seam over — is correct.

### D-6 — the emission site — **REAL**

`internal/graphlint/coverage.go:33-43`:

```go
func (a *analysis) checkCoverage(g guard.Group) {
	dims := guard.Dimensions(a.model, g)
	if len(dims) == 0 {
		// … the two scans are skipped rather than run over nothing.
		a.emitCoverageArms(g)
		return
	}
```

`emitStructurallyUnprovable` at `coverage.go:147-149` is `for _, key := range
guard.Dimensions(a.model, g)` — zero iterations when `len(dims) == 0` — and is called
only at `coverage.go:65`, inside the product-bound branch, below the early return.
`emitUnprovableDimensions` (`:87`, def `:110`) and `emitWithholdings` (`:66`, `:83`,
`:88`, def `:261`) are likewise all below it. `bareEscapeFor` (`:505`) is reached from
inside `emitCoverageArms` at `:351`. So A13's stated verification target
("whether `emitStructurallyUnprovable`'s existing per-dimension shape admits a
group-level emission", RDR:441-444) is unanswerable as posed, and the precedence
question against `graph-coverage-closed-by-escape` is real and unaddressed by C5.

### D-7 — `emit` string-only — textual facts **REAL**, sufficiency **UNGROUNDED**

A6 (RDR:236-252) quotes `Next: /rdr-prelock 0046 critique` and resolves it as "the
consumer treating the command as one opaque string"; its "If wrong" concedes "the wire
shape needs an `any`-typed value later (a widening, not a break)". BR3 (RDR:1022-1023)
rejects `any`-typed values because "a dump column needs a deterministic rendering".
C3 (RDR:598-600) types `sourceRule.Emit` as `map[string]string`. All quoted text
verifies. A's added observation — `0046` is per-invocation caller data, so no literal
emit block can be written — is a reading of the consumer's requirement, not a code
claim; it is the sharpest form of the finding but remains judgment.

### D-8 — Pending assumptions under a `MUST` — **REAL**

A12 `Status: Pending` (RDR:412), A13 `Status: Pending` (RDR:432), both `Evidence: to
verify`. §metadata `Status: Draft` (RDR:14). The Finalization Gate's own clause
(RDR:1457-1460) reads "**Status consistency:** no assumption marked `Pending` or
`Unverified` may have settled-fact prose elsewhere in the RDR depending on it." C5
(RDR:702-706) states the zero-dim fence as an unconditional `MUST`. The Trace table
row 2″ (RDR:1345) reads "**unwitnessed — A13 is Pending** … this is the one Trace row
with no witness". The record fails its own gate as written. Both passes reach this
independently.

### D-9 — one-directional vs both ways — **REAL**

- C1, RDR:531-533: "The agreement check is **one-directional**: a `"decision-table"`
  model MUST declare zero tags of provenance `owned` … A `"state-machine"` model MUST
  NOT be refused for declaring zero owned tags".
- §approach, RDR:457: "the loader refuses a class that disagrees with the owned set in
  either direction"; RDR:483: "declaring it makes the mismatch a load failure both ways."
- §decision-rationale scorecard, "Silent misclassification" row, column C: "refused at
  load both ways" — and the prose immediately below says "C wins on the misclassification
  and overrides rows", i.e. the decision was made on that cell.
- ALT2 reason-for-rejection: "one declared line buys a load-time refusal in both
  directions and keeps 0006's root rule intact."

C1 is right on the merits (its own parenthetical explains that a bidirectional check
would make every rootless state machine unloadable and MVV step 6's class-omitted
negative control unwritable). Three narrative passages, one of them the scorecard cell,
contradict it.

### D-10 — emit through `expand` — code **REAL**, citation **PARTLY WRONG**

`internal/table/normalize.go:646-655`:

```go
	rows := []Row{{
		ModelID:       base.ModelID,
		RuleID:        base.RuleID,
		SourceLocator: base.SourceLocator,
		Source:        base.Source,
		Gate:          base.Gate,
		RequiresOwned: base.RequiresOwned,
		Escape:        base.Escape,
		setKeys:       base.setKeys,
	}}
```

Eight named fields, no struct copy — a new `Row` field is silently dropped. A says "nine
named fields" and then enumerates eight. `rowByID` (`internal/cli/flow_resolve.go:186-193`)
is first-match, as C3 asserts. But A's framing that C3 "is stated as if it already holds"
overstates: C3 at RDR:628-630 says `Row.Emit` "MUST be carried through
`internal/table/normalize.go::expand` for every expanded row, **which builds each row from
the seed literal rather than copying `base`**" — the RDR names the hazard explicitly. The
surviving defect is narrower and still real: no Testing Strategy scenario asserts an emit
block on a rule with a multi-member match `in` atom.

### D-11 — the licensed-diff rule — **REAL**

Done clause, RDR:1295-1303: "a changed line differs only by the addition of a trailing
` emit=[…]` dump cell or an `"emit":` payload member. Any other changed line is a
regression, not a licensed diff … and the 103 `[dump]`-declaring fixtures under
`internal/table/testdata/` (A4) are licensed diffs under the rule, not violations of it."
The required edit to those 103 files is an added member inside `order = [ … ]`, which is
neither admitted shape. The clause asserts the conclusion instead of admitting the shape.
Combined with D-2, a third shape (a Go `want := []string{…}` literal) is also required
and also unadmitted.

### D-13 — observed vs guard dimensions — **REAL**, count corrected upward

`internal/guard/product.go:130-141`:

```go
func Dimensions(m *table.Model, g Group) []string {
	var out []string
	for _, row := range g.Rows {
		for _, a := range guardAtoms(row) {
```

Dimensions are collected from guard atoms alone; `provenance` (observed/owned) is a
separate `[tags.<tag>]` field. Occurrences of the wrong phrasing:

- §problem-statement, RDR:126: "Exhaustiveness/coverage lint over observed dimensions must keep working unchanged"
- §approach, RDR:471: "overlap/coverage run unchanged over the observed dimensions"
- §technical-design data-flow, RDR:521: "lint (root ∅, groups over observed dimensions)"
- **and C5 itself**, RDR:697-698: "MUST run unchanged over the authored observed dimensions"

A named three; there are four, and the fourth is inside a normative contract — which
makes the defect worse than A states, since C5 ¶1 and C5 ¶2 disagree with each other on
the same page.

### D-14 / D-15 / D-16 / D-17 / D-18 — B's rows

D-14's mechanical premise verifies at `coverage.go:33-43` (a zero-dimension group closes
vacuously on any one row today) — the finding itself is about whether the fence's *message*
must be actionable, which no contract requires; judgment. D-15 verifies as an absence:
§authority (RDR:745-755) names a `len(owned) == 0` "sibling arm" and offers no enforcement,
and the Testing Strategy carries no self-lint over the intrastate tree. D-16's premise is
A9's own "If wrong" arm (RDR:322-323), which names precisely that outcome and accepts it.
D-17: C4 at RDR:673 "`--outcome` remains required in both classes"; BR5 at RDR:1027-1028
parks it "as a plain kata at triage" with no id. D-18: A5 at RDR:224-226 carries
"⚠ no corpus coverage for decision-table hit policies" verbatim and self-scopes to the
motivating consumer.

---

## Convergence

Three passages drew independent findings from both models. Independent convergence at the
same anchor is the strongest signal in this diff:

| ID | Anchor | What both passes saw |
| --- | --- | --- |
| **D-7** | `0010:A6` / `0010:C3` / `0010:BR3` | Both named A6 as *the* assumption that will not survive first contact, both quoting the same `Next: /rdr-prelock 0046 critique` evidence line back at the RDR, and both reaching the same conclusion: the RDR's self-awareness about the limit is the tell, not the mitigation. A adds the decisive detail (`0046` is caller data, so no literal block exists). |
| **D-8** | `0010:A12` / `0010:A13` / `0010:§assumption-verification` | Both flagged two `Pending` assumptions under C5's `MUST`, both cited the Trace table's own unwitnessed row 2″, and both wrote an acceptance test that is simply the RDR's own Finalization Gate clause turned on the RDR. B additionally named C5 as "the one section rewritten within 6 weeks" for this reason; A named C5 for the same reason plus D-1 and D-6. |
| **D-12** | `0010:A9` "Consequence carried, not restated" | Both saw the N-escape-rows-per-outcome consequence discovered and then routed to documentation, and both framed it as a *mitigation-parity* failure against C5's fence for the structurally identical mistake. |
| **D-18** | `0010:A5` | Both flagged the ⚠ single-consumer hit-policy verification; B as a ledger row, A as its §3 runner-up. |

**C5 is the hotspot.** Both passes independently nominated `#### Normative Contracts →
C5` as the section that gets rewritten within six weeks, by different routes: B via A13's
Pending status and the guard/match authoring trap, A via the `reason` set (D-1), the
unreachable emission site (D-6), and the class-conditionality argument. Five of the six
BLOCKING rows in this ledger (D-1, D-6, D-8, D-9, D-13) land on or contradict C5.

---

## Ledger bugs

Findings raised in a file's prose that its own ledger does not row:

1. **Pass A, §2, second candidate — `0010:§load-bearing-decisions`.** A's §2 argues that
   the Load-Bearing Decisions section "is four bullets long and says nothing about the
   payload field order (C-5), the expansion carry-through (C-10), or the emit-key/tag-key
   namespace collision risk it delegates entirely to a Risks bullet … For a 'foundational'
   Profile spanning three modules, four bullets is not a Load-Bearing Decisions section."
   That is a distinct passage-level defect against `0010:§load-bearing-decisions`
   (RDR:757-775, confirmed: four D-rows — `D-identity`, `D-wire-byte-format`, `D-naming`,
   `D-selection-predicate`), and **no ledger row carries it**. A ledger bug.

2. **Pass A, §4 premortem week 10 / §1 §3 — the `emit` key namespace.** The premortem and
   §2 both reference "the emit-key/tag-key namespace collision risk"; the ledger has no row
   for it. C3 (RDR:621-625) does address it normatively ("Emit keys are NOT tag keys …
   MUST NOT be matched, guarded, written, or read by any accessor"), so the unrowed
   concern is about *placement* (Risks bullet vs Load-Bearing Decisions), not about a
   missing contract — which is why it folds into ledger bug 1 rather than standing alone.

3. **Pass B, §2 — the state-machine false-positive risk.** B's §2 raises "if it turns out
   some checked-in state-machine fixture has a legitimate zero-dimension group that would
   false-positive under a wrongly-keyed arm (a risk A13's own 'If wrong' names)". This is a
   sweep obligation on `0010:A13` distinct from B:C-8's "A13 is Pending" row, and B's
   ledger does not carry it. Minor ledger bug; A13's own evidence (RDR:444-446) already
   states the verification step, so the RDR is not silent — B's prose is.

No row in either ledger is unanchored, and no row is raised in a ledger without prose
support.

---

## Pass health (per §expected-signal)

**Pass A (`claude-opus-5`) — healthy, with a counting discipline problem.**

Thirteen rows, every one anchored to a named RDR element (`0010:C5` ¶2, `0010:A4`,
`0010:§technical-design` item 1, `0010:§problem-statement`) and, in eleven of thirteen,
to a named `path::Symbol` on `main`. Origins span §1, §2, §3, the premortem, and the AT
block — no row is premortem-only or AT-only. It found five defects that require reading
shipped source to see at all (D-1, D-2, D-3, D-6, D-10), and four of those verify exactly
as stated. The failure mode is arithmetic, not grounding: three rows carry a wrong count
(five `len(Initial)` sites vs four; eleven or "five named, four omitted" payload fields vs
thirteen; nine seed-literal fields vs eight), and two rows disagree with their own
acceptance tests on the same number. Two rows also under-credit the RDR's own text
(D-10's "stated as if it already holds" — C3 names the seed literal explicitly; D-13
misses that C5 ¶1 commits the same error it diagnoses elsewhere, which strengthens rather
than weakens the finding). **Verdict: healthy.** Substance survives every correction;
only the numbers need fixing before any of this reaches the record.

**Pass B (`claude-sonnet-5`) — healthy but shallower.**

Eight rows, all anchored to real elements (`0010:C1`, `0010:A6`, `0010:C4`,
`0010:§authority`), origins spanning §1, §2, §3, the premortem, and the AT block. Every
claim it makes about the RDR's *own text* verifies. Its distinctive contribution is the
authoring/UX and enforcement axis A does not cover: the guard-vs-match discriminator as a
newcomer trap (D-14), the absence of any mechanism enforcing the four-reader invariant
(D-15), the missing-fix-guidance problem in refusal messages (D-16), and the
single-outcome `--outcome` friction (D-17). It makes almost no independent source claims —
every code reference it uses is quoted from the RDR's own evidence rather than re-derived
— so it could not have found D-1, D-2, D-3, D-6, or D-10, and its C-8 is the only row that
would block lock on its own. Its premortem is concrete and named (three tickets, each
traced to a ledger row), not generic advice. **Verdict: healthy.** Independent
convergence with A on D-7, D-8, D-12, and D-18 without shared source access is the
strongest evidence in this diff that those four are real.

**Combined.** Both passes are healthy by §expected-signal: concrete named passages, real
anchors, origins spanning several sections, and no generic advice. Six BLOCKING rows
(D-1, D-2, D-6, D-8, D-9, D-13); five of them land on or contradict C5, and D-9 and D-13
are pure-text contradictions inside the record that a grounding sweep, a cove pass, and a
3amigo pass all missed.
