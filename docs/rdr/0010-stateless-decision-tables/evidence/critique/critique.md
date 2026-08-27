Model: claude-opus-5

# Critique — RDR 0010, stateless decision tables

## Findings Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | `0010:C5` ¶2 ("MUST take `graph-unprovable-coverage` naming the group"), `0010:A13` ("without a taxonomy change") | The new zero-dimension emission has no `reason`. `0006` REQ-80 fixes `reason` as a **closed, append-only three-value set** (`dimension-not-finite`, `tag-not-single-valued`, `row-can-refuse`) and asserts it byte-exact in `internal/graphlint/findings_0006_test.go::TestReq80_UnprovableCoverageCarriesAReasonFromTheClosedSet`, which additionally asserts *every* emitted finding of that code carries a reason from the set. A group-level emission has no dimension and none of the three reasons applies. The RDR contains the word "reason" nowhere normative. | `make check` fails on a shipped 0006 conformance test the implementer did not know existed. The fix is either a fourth reason value (an append to a closed 0006 set = a route-back 0006 owns) or a `reason`-less emission that fails REQ-80 outright. Implementation stalls mid-Phase-2. | §1, §3, premortem, AT-1 |
| C-2 | `0010:A4` ("Count confirmed exact: **103** files under `internal/table/testdata/` … the 73 further repo-wide matches … are all under `docs/rdr/0002-*/evidence/`"), `0010:§testing-strategy` Done clause ("There is no golden-file tree in this repo") | The census counted `.toml` files only. `internal/table/dump_test.go::TestReq95_DumpColumnVocabularyIsClosedAndVerbatim` hard-codes the ten-column vocabulary as a Go slice literal compared with `reflect.DeepEqual`, and its own comment says "The vocabulary is fixed VERBATIM to the fixtures' spelling — it is the one place the record forbids re-derivation." The Done clause's licensed-diff rule admits only "a trailing ` emit=[…]` dump cell or an `\"emit\":` payload member" — it does not license editing a `want` literal in a Go test, so the *one* file that must change is the one the rule forbids changing. | `go test ./internal/table` fails with `dump column vocabulary = [… escape emit]; want [… escape]`. The implementer must either violate the RDR's own Done clause or route back to 0002. | §1, §2, premortem, AT-2 |
| C-3 | `0010:A10` ("provably silent … `graph-unreachable-rule`"), evidence bullet ("`checkUnreachableRules` is silent because C5's single ∅ node satisfies every context vacuously") | The stated mechanism is wrong. `internal/graphlint/analysis.go::checkUnreachableRules` returns early on `len(a.model.Initial) == 0` — it never consults `a.reachable` for a rootless model. Under C5 a decision table still has `len(Initial) == 0`, so the code is silent for a reason C5 does not control, and the RDR's reason (∅ node satisfies every context) is an accident that happens to agree. The same misattribution applies to `checkAlwaysPresentOwned` (analysis.go:233), also gated on `len(Initial)`, which A10 lists as silent "because it filters on `ProvenanceOwned` and finds none." | Silent today; a latent trap. Any later change that makes `Initial` non-empty for a decision table, or that replaces the `len(Initial)` guards with the class accessor C1 mandates ("every site that today asks `len(Initial) == 0` … reads that accessor"), flips both checks live and floods every decision table with `graph-unreachable-rule` on every row. | §1, premortem, AT-3 |
| C-4 | `0010:§technical-design` item 1 ("every site that today asks `len(Initial) == 0` or would ask `len(owned) == 0` reads that accessor, so there is exactly one place the question is answered") | There are **five** `len(Initial) == 0` sites in shipped source (`reach.go:91`, `analysis.go:113`, `:233`, `:515`, plus `reach.go:99`'s iteration), and C5/§technical-design names only two of them (`reach`, `checkDanglingEdge` root arm). The instruction "every site … reads that accessor" is directly contradicted by C5, which requires `checkAlwaysPresentOwned` and `checkUnreachableRules` to stay `len(Initial)`-keyed (they are listed as "vacuous by construction," not class-keyed). The design's own single-source-of-truth claim is false as written. | An implementer following §technical-design literally converts all five sites and breaks `graph-unreachable-rule`/`graph-always-present-owned` for the state-machine class; an implementer following C5 leaves three sites on the old key and the "exactly one place" invariant is a lie in the record. Either way the four-site-drift risk the RDR names is the risk the RDR causes. | §1, §2, premortem, AT-4 |
| C-5 | `0010:C4` ("The `flow resolve` success payload MUST carry `emit`… `next`, `writes`, `clear`, `owned`, and `readers` keep their `0005:C1` shapes") | C4 enumerates five payload fields and omits the four the struct actually carries alongside them: `model`, `revision`, `observed`, `outcome`, `gates`, `escaped`, `escape_class`. It never says **where** `emit` sits in `internal/cli/flow_resolve.go::resolvePayload`'s field order, which is the JSON key order Go's encoder emits. C3 was scrupulous about the dump ("appended last, after `escape`" — with the reason spelled out) and C4 says nothing at all about the same question one seam over. | Two implementers produce two different key orders. Any consumer or fixture asserting on payload bytes (the repo asserts JSON inline in Go tests) sees an unlicensed diff, and the Done clause cannot adjudicate because it licenses `"emit":` "as a member" without fixing position. | §1, §2, AT-5 |
| C-6 | `0010:C5` ¶2 ("a table that discriminates with `[rule.match.<key>]` atoms contributes no guard dimension … so its scoped product is the empty product") + `0010:§illustrative-code` | The class-conditioned fence lands in the *wrong branch*. `internal/graphlint/coverage.go::checkCoverage`'s `len(dims) == 0` arm calls `emitCoverageArms` and **returns before `emitUnprovableDimensions` and `emitWithholdings` ever run**. Adding the fence there is a brand-new emission site, not the reuse A13 claims ("`emitStructurallyUnprovable`'s existing per-dimension shape" — that function loops `guard.Dimensions`, which is empty, so it can never fire here). Worse: over a decision table, `coverageUnionFor`'s zero-dim arm returns the whole product on the first row of *any* kind, so today a one-row match-only table is green — and the fence must fire *before* that, in a branch A13 does not name. | The implementer adds the arm where A13 points (`emitStructurallyUnprovable`), sees no finding, and either debugs for hours or ships the arm in `emitCoverageArms` where it now fires *after* `bareEscapeFor` has already computed `graph-coverage-closed-by-escape` — yielding both findings on one group. | §1, premortem, AT-6 |
| C-7 | `0010:A6` ("A flat, string-valued `[rule.emit]` table is sufficient for the motivating consumer's answer"), `0010:BR3` ("`any`-typed emit values … widenable later") | The consumer's answer is already structured and A6 admits it in its own evidence: `rdr/skills/rdr-status/SKILL.md` renders "`Next: /rdr-prelock 0046 critique`" — a verb plus two arguments, one of which (`0046`) is *per-invocation data the table cannot know*. A6 resolves this by declaring the command "one opaque string," which means the emit block must carry a template or the consumer must string-concatenate. Neither is specified. | The first real consumer authors `next = "/rdr-prelock {id} {lens}"` or `verb = "/rdr-prelock"` + an out-of-band argument convention — reinventing exactly the "consumer-invented undocumented format" the Problem Statement says this RDR exists to abolish. | §3, premortem, AT-7 |
| C-8 | `0010:A12` ("Status: Pending"), `0010:A13` ("Status: Pending"), `0010:§metadata` **Status: Draft** | Two Pending assumptions carry the two load-bearing consequences of C5, and prose elsewhere in the record treats both as settled: C5 states the zero-dimension fence as a `MUST` with no conditional; §risks-and-mitigations calls it "a contract, not documentation"; the Trace table's row 2″ admits it is "the trace's only unwitnessed row." The Finalization Gate's own Assumption Verification clause forbids exactly this ("no assumption marked `Pending` … may have settled-fact prose elsewhere in the RDR depending on it"). | The RDR locks with the gate's own stated rule violated, and C5 ships a `MUST` whose implementability is unproven. When it fails (C-1), the failure lands in implementation, not review. | §2, §3, AT-8 |
| C-9 | `0010:C1` ("The agreement check is **one-directional**") vs `0010:§approach` ("the loader refuses a class that disagrees with the owned set **in either direction**") and `0010:§decision-rationale` row "Silent misclassification: refused at load both ways" and `0010:ALT2` reason-for-rejection ("one declared line buys a load-time refusal **in both directions**") | Flat contradiction between the normative contract and three separate narrative passages. C1 is right (bidirectional would make rootless state machines unloadable and break 0006:C18); the Approach, the scoring table's winning cell, and the Alternative-2 rejection rationale are all wrong. The scoring table's decisive row — the one the Decision Rationale says C "wins on" — is stated in terms of the property C1 explicitly disclaims. | A reviewer reads the Approach and the winning scorecard cell and implements the bidirectional check. Every existing rootless state-machine model (and the RDR's own class-omitted negative control, MVV step 6) refuses at load instead of lint. The MVV's own control becomes unwritable — which C1 says in so many words, and no one reading the Approach will notice. | §1, §2, premortem, AT-9 |
| C-10 | `0010:C3` ("the CLI joins it back on `Plan.RuleID` … the existing first-match `internal/cli/flow_resolve.go::rowByID` is correct as it stands") | `rowByID` returns the **first** row whose `RuleID` matches, and `normalize.go::expand` mints one row per match-block `in` member. C3's soundness argument is "emit is rule-scoped so every expanded row carries the same block" — but that is only true if `expand` copies `Emit`, and `expand` builds each row from an explicit seed literal (`normalize.go:646`) enumerating `ModelID, RuleID, SourceLocator, Source, Gate, RequiresOwned, Escape, setKeys` — nine named fields, no struct copy. C3 does say `Emit` "MUST be carried through `expand`," but the join-soundness argument is stated as if it already holds, and nothing in the Testing Strategy asserts an emit block on an **expanding** (`in`-atom) rule. | A rule with `[rule.match.k] in = ["a","b"]` and an emit block resolves to the right rule id and an empty `emit`. Silent wrong answer — the payload says `{}` and the consumer takes the no-answer branch. No test in the RDR catches it: S3 tests emit normalization, S4 tests lint, S5 tests resolve on a non-expanding fixture. | §1, premortem, AT-10 |
| C-11 | `0010:§testing-strategy` Done clause ("every diff to a checked-in expectation is **licensed** … a changed line differs only by the addition of a trailing ` emit=[…]` dump cell or an `\"emit\":` payload member") | The licensing rule is a *line-level* rule applied to a *file-level* problem. 103 `[dump]` fixtures must gain `emit` to their `order` list — that is a change to the `order = [...]` line, which adds neither a dump cell nor a payload member. The rule as written classifies all 103 required edits as regressions. A4 asserts the opposite ("the 103 … are licensed diffs under the rule") without amending the rule. | The one mechanical gate protecting a 103-file sweep from hiding a real regression is self-contradictory, so the implementer will disable it or reason around it — and the regression it exists to catch ships. | §2, premortem, AT-11 |
| C-12 | `0010:A9` ("**Consequence carried, not restated:** `0002:C5` scopes rescue per outcome, so 'the otherwise row' is per-outcome — a table over an N-outcome alphabet needs N escape rows … which the Phase 4 authoring docs must say") | The N-escape-rows consequence is discovered, written down, given a test (S4's "two-outcome escape control"), and then routed to *documentation* as its only mitigation. There is no contract clause and no lint finding for "a decision table has an outcome with no rescue row." Compare C5, which was willing to add a fence for the match-atom mistake. The two mistakes are the same shape; one gets a `MUST`, the other gets a doc. | A four-outcome decision table with one "otherwise" row lints green on the rescued outcome and reports a coverage gap on the other three — but the author reads the gap as "add more rows," adds them, and the table is now three escape rows short with no diagnostic saying so. The idiom is undiscoverable from the tool. | §2, §3, AT-12 |
| C-13 | `0010:§problem-statement` ("Exhaustiveness/coverage lint over observed dimensions must keep working unchanged") vs `0010:§key-discoveries` bullet 6 ("coverage is proven over **guard** dimensions") and `0010:§illustrative-code` | The Problem Statement — the sentence stating *why the RDR exists* — says "observed dimensions." Coverage is proven over **guard** dimensions; observed-ness is a `provenance`, orthogonal to whether an atom is a match or guard atom. The Approach repeats the error ("overlap/coverage run unchanged over the observed dimensions"), and the Data-flow line repeats it a third time ("groups over observed dimensions"). The correct statement appears only in Key Discoveries and in C5's second paragraph. | The one paragraph a future reader consults to learn what the RDR is for teaches the exact authoring mistake C5 exists to fence. The record undermines its own contract at the top of the file. | §2, §3, AT-13 |

---

## 1. The three most likely ways implementation goes wrong

### 1.1 Phase 2 stalls on a closed vocabulary the RDR never looked at

**Root cause in the RDR.** C5 and A13 assert, twice and emphatically, that the zero-participating-dimension fence is "not a new code" and needs "no taxonomy change." Both statements are true about the *code* field and false about the *record*. `graph-unprovable-coverage` is the one finding in 0006's taxonomy that carries a second discriminator — `reason` — and 0006 fixes it as a closed, append-only set of exactly three values.

**The passage that enabled it.** `0010:C5` ¶2:

> A `"decision-table"` group whose scoped product has **zero participating dimensions** MUST take `graph-unprovable-coverage` naming the group, and MUST NOT be reported as covered. This is the one fence the class needs and it is not a new code

and `0010:A13`'s title clause, "*without a taxonomy change*."

Ground truth, `internal/graphlint/taxonomy.go`:

```go
ReasonDimensionNotFinite = "dimension-not-finite"
ReasonTagNotSingleValued = "tag-not-single-valued"
ReasonRowCanRefuse       = "row-can-refuse"
```

and the shipped conformance test, `internal/graphlint/findings_0006_test.go::TestReq80_UnprovableCoverageCarriesAReasonFromTheClosedSet`, which asserts `slices.Equal(graphlint.Reasons(), want)` against that literal three-element slice **and** walks every emitted finding asserting that any `CodeUnprovableCoverage` finding carries a reason drawn from it.

Every one of the three reasons is *dimension-scoped* and names a *per-dimension remedy* — "declare the domain," "add the single-valued marker," "the tag is optional so the row can refuse." A group-level finding over zero dimensions has no dimension to name and no such remedy. Its remedy is "author your discriminators as guard atoms," which is not in the set.

So the implementer's options are: (a) emit with `reason: ""` — REQ-80's second half fails; (b) emit with a borrowed reason — the message contradicts the reason and REQ-81's "reads as a withholding, not an error" test is next in line; (c) append a fourth reason — an append to a closed set that 0006 owns, i.e. a route-back, and the exact "route-back" outcome A13's own "If wrong" arm names, arriving after Phase 1 and Phase 2 are already merged.

The RDR contains the word "reason" nowhere in a normative position. Thirteen assumptions, 1543 lines, five contracts, a Trace table, an Oracle table, and the single field that decides whether C5's central fence is implementable is not mentioned once.

**Symptom.** The implementer writes the fence, runs `make check`, and `TestReq80` fails on a test they had no reason to read. The blast radius is a Final predecessor's conformance suite. Phase 2 does not land. C5 is renegotiated in a hallway rather than in the record.

### 1.2 The `[dump]` sweep breaks a test the census could not see

**Root cause.** A4's verification is a `grep -rl '^\[dump\]' internal/table/testdata/ | wc -l` — a count of TOML *files*. The dump vocabulary is also asserted from Go, as a slice literal, in a test whose comment explicitly forbids re-derivation.

**The passage.** `0010:A4`:

> **Count confirmed exact: 103** files under `internal/table/testdata/` declare `[dump]`, all 103 with an explicit `order =` list (`grep -rl '^\[dump\]' internal/table/testdata/ | wc -l`); the 73 further repo-wide matches (176 repo-wide, less the 103) are all under `docs/rdr/0002-*/evidence/`

Ground truth: `grep -rl '\[dump\]'` (no `^`, no extension filter) returns five files under `internal/table/*.go`. The decisive one is `internal/table/dump_test.go:293`:

```go
// The vocabulary is fixed VERBATIM to the fixtures' spelling — it is the
// one place the record forbids re-derivation.
func TestReq95_DumpColumnVocabularyIsClosedAndVerbatim(t *testing.T) {
	want := []string{
		"identity", "source", "kind", "outcome",
		"atoms", "next", "writes", "requires_owned", "gate", "escape",
	}
	got := table.DumpColumns()
	if !reflect.DeepEqual(got, want) { ... }

	m := mustLoad(t, rdrFixture)
	if !reflect.DeepEqual(m.DumpOrder, want) { ... }
```

Two `reflect.DeepEqual` assertions against a hand-written literal, plus a third against a fixture's decoded `DumpOrder`. A4's arithmetic (176 − 103 = 73, all under `docs/`) is not just incomplete — it is *reassuring*, and it is the reassurance that stops anyone looking further. The count was performed with a regex anchored to `^`, on `.toml`, in one directory, and reported as "**Count confirmed exact**."

The Done clause compounds it. `0010:§testing-strategy`:

> a changed line differs only by the addition of a trailing ` emit=[…]` dump cell or an `"emit":` payload member. Any other changed line is a regression, not a licensed diff. There is no golden-file tree in this repo — expectations live inline in Go tests and in `testdata` TOML/JSON

The Done clause *knows* expectations live inline in Go tests, and its licensing rule still admits only two line shapes, neither of which is `want := []string{...}` gaining an eleventh element.

**Symptom.** `go test ./internal/table` fails:
`dump column vocabulary = [identity source kind outcome atoms next writes requires_owned gate escape emit]; want [identity ... escape]`.
The implementer edits `want`, and now the Done clause says their diff is a regression. They will disable the rule, because the alternative is stopping the RDR — and the rule is the only mechanical guard on a 103-file sweep.

### 1.3 The fence is implemented in the wrong branch and lands after the closure

**Root cause.** A13 grounds its mechanism in the right place and its *implementation site* in the wrong one.

**The passage.** `0010:A13`:

> The mechanism is grounded — the zero-dimension branch is explicit at `internal/graphlint/coverage.go::checkCoverage` (`len(dims) == 0` → `emitCoverageArms` …). What is unverified: whether `emitStructurallyUnprovable`'s existing per-dimension shape admits a group-level emission with no dimension to name (its `Dimension` field would be empty)

`emitStructurallyUnprovable` cannot be the site. Ground truth, `coverage.go:147`:

```go
func (a *analysis) emitStructurallyUnprovable(g guard.Group) bool {
	var any bool
	for _, key := range guard.Dimensions(a.model, g) {
```

The function's entire body is a loop over `guard.Dimensions`. When `len(dims) == 0` that loop has zero iterations and the function returns `false` unconditionally. It is not a matter of whether it "admits a group-level emission with no dimension to name" — it is structurally unreachable in the case A13 wants it for. Worse, `emitStructurallyUnprovable` is called only from inside the *product-bound* branch, which `checkCoverage` reaches only after the `len(dims) == 0` early return has already fired:

```go
dims := guard.Dimensions(a.model, g)
if len(dims) == 0 {
	a.emitCoverageArms(g)
	return
}
```

`emitUnprovableDimensions` and `emitWithholdings` are *both* below that return. The zero-dimension path bypasses every existing unprovable-coverage emitter in the file. A13's "verify by adding the class-keyed arm" is not a verification of an existing shape; it is the design of a new emission site, deferred to implementation, on a `MUST` clause.

And the ordering matters concretely. `emitCoverageArms` computes `bareEscapeFor` and may emit `graph-coverage-closed-by-escape`. If the fence is added inside `emitCoverageArms` rather than before it, a match-discriminated decision table with an escape row gets *both* `graph-unprovable-coverage` and `graph-coverage-closed-by-escape` on one group — an unprovable group that lint simultaneously reports as closed. Nothing in C5 forbids that pairing, and the Testing Strategy's match-only control (S4) says "reports exactly `graph-unprovable-coverage`," which such an implementation would fail without telling anyone which of the two is wrong.

**Symptom.** The implementer follows A13 to `emitStructurallyUnprovable`, adds the arm, writes the fixture, gets no finding, and burns a day before discovering the function is unreachable. Then they put it in `emitCoverageArms`, S4 fails on a second finding, and C5 needs a precedence clause it does not have.

---

## 2. The one section that will be rewritten within 6 weeks of shipping

**`#### Normative Contracts` → `C5`.**

Not the Testing Strategy (bad but survivable), not A4 (a count, fixable in an hour). C5, because C5 is the only contract in this RDR that changes *what lint says about a model that is not a decision table's fault*, and it does it with a clause that cannot be implemented as specified.

Three independent forces converge on C5 inside six weeks:

**Force 1 — the `reason` field (C-1).** C5's second paragraph mandates an emission that 0006's REQ-80 conformance test rejects. The fix is a fourth reason value, which is an amendment to a closed set in a Final predecessor, which is a new RDR. C5's text will be superseded the moment that RDR lands.

**Force 2 — the wrong branch (C-6).** C5 says *what* must be emitted and never *where*, and A13's answer to "where" is a function that cannot fire in the case. Whoever implements it discovers the precedence question against `graph-coverage-closed-by-escape` and the mutual-exclusion question against `emitWithholdings` — neither answered in the record. Both answers are contract text, and contract text that is discovered during implementation gets written into the next RDR, not this one.

**Force 3 — C5's clause is class-conditional on a property that is not class-conditional.** C5 justifies binding the fence to `decision-table` alone because "over a state machine that shape is legitimate (a group may genuinely range over no guard dimension)." That justification is about *whether the author meant it*, not about *whether the claim is provable* — and the claim is exactly as vacuous in both classes. The first state-machine author who ships a green lint over a zero-dimension group and later discovers it proved nothing will ask why the fence was class-keyed, and the honest answer ("because we only noticed the hole while writing the decision-table RDR") does not survive the question.

The second candidate is `#### Load-Bearing Decisions`, which is four bullets long and says nothing about the payload field order (C-5), the expansion carry-through (C-10), or the emit-key/tag-key namespace collision risk it delegates entirely to a Risks bullet. For a "foundational" Profile spanning three modules, four bullets is not a Load-Bearing Decisions section; it is a Naming section with three neighbours.

---

## 3. The one assumption that will not survive first contact with a real user

**A6.** "A flat, string-valued `[rule.emit]` table is sufficient for the motivating consumer's answer."

It does not survive because **A6's own evidence already refutes it and the record files the refutation as a limit rather than a failure.**

Read the evidence clause verbatim:

> The consumer's actual answer shape was checked, not assumed: `rdr/skills/rdr-status/SKILL.md` renders "the exact command to run, e.g. `Next: /rdr-prelock 0046 critique`" — a command with arguments, which a flat string carries whole. That is precisely the limit this assumption's "If wrong" already names, so the claim holds *with its limit made concrete*

`/rdr-prelock 0046 critique`. The `0046` is an RDR number. It is not knowable at authoring time — it is the thing the caller is asking *about*. A decision table cannot emit it as a literal. So the consumer's real emit block is one of:

- `next = "/rdr-prelock {id} critique"` — a template. Intrastate has no template semantics, no substitution point, no escaping rule, and C3 says emit values are "undeclared, uninterpreted, compared by exact byte equality." The consumer now owns an interpolation format that intrastate refuses to know about — which is, word for word, the "consumer-invented file format" the Problem Statement gives as the reason this RDR exists.
- `verb = "/rdr-prelock"` + `lens = "critique"` and the caller assembles — which needs the verb and its arguments as separate fields, which A6's "If wrong" identifies as "the widening A6 defers."

The assumption is not that a string suffices. The assumption is that the *first consumer* will accept a string. And the first consumer is rdr#tmxk, whose answer shape A6 read, quoted, and then declared sufficient by redefining sufficiency as "the consumer treats the command as one opaque string" — a property of the consumer, asserted here, verified nowhere, and contradicted by the presence of `0046` in the example the RDR chose.

The self-awareness is the tell. A6 anticipates the failure precisely, in the right words, and then routes it to "a widening, not a break." Widenings are cheap only before a wire format has consumers. `emit` is a dump column and a payload field; both are byte-asserted; and BR3 already rejected `any`-typed values because "a dump column needs a deterministic rendering." The widening A6 defers is the one BR3 rejected on grounds that will not have changed.

Runner-up: **A5** (exact-one is the hit policy), which is `Method: Prior Art` with `⚠ no corpus coverage` and a single consumer's acceptance criterion standing in for a survey. It survives only as long as there is exactly one consumer. The first table that wants "first match wins" gets `flow-ambiguous-match` and an answer of "restructure your table," and A5's "If wrong" concedes the class then "needs a selection rule the kernel does not have."

---

## 4. Premortem

*Written from twelve weeks after lock. RDR 0010 shipped in four phases. It is now being partly reverted.*

**Week 1 — Phase 1 lands late and larger than planned.** The implementer adds `class` to `sourceModel`, `emit` to `sourceRule`, `Row.Emit`, and appends `emit` to `dumpColumns`. `make check` fails immediately in `internal/table`, not in `testdata`: `TestReq95_DumpColumnVocabularyIsClosedAndVerbatim` compares `table.DumpColumns()` against a hand-written ten-element `want` with `reflect.DeepEqual`, and a second assertion compares `mustLoad(t, rdrFixture).DumpOrder` against the same literal. A4 had counted 103 `.toml` files with `grep -rl '^\[dump\]' internal/table/testdata/` and reported "Count confirmed exact"; the Go-side literals were outside the grep. The implementer edits `want`, which the RDR's own Done clause classifies as an unlicensed diff — the rule admits only a trailing ` emit=[…]` dump cell or an `"emit":` payload member. They note the conflict in the PR and disable the check for the sweep. The 103 fixture edits go in unaudited. Two of them silently lose a `requires_owned` cell to a bad `sed`; nobody notices, because the mechanical rule that would have caught it was turned off in week 1 to accommodate a defect in the RDR that wrote the rule.

**Week 3 — Phase 2 hits the closed reason set.** `reach` is class-seeded and `checkDanglingEdge`'s root arm is class-keyed; the A1 spike's arm-C finding set reproduces exactly, as promised. Then the implementer writes C5's zero-dimension fence. Following A13, they open `emitStructurallyUnprovable` and find its whole body is `for _, key := range guard.Dimensions(a.model, g)` — zero iterations in the zero-dimension case, and reachable only from the product-bound branch, which sits *below* `checkCoverage`'s `len(dims) == 0` early return. A13's "verify whether its existing per-dimension shape admits a group-level emission" cannot be answered because the function cannot be reached. They move the emission into `checkCoverage`'s zero-dim branch, ahead of `emitCoverageArms`, and `TestReq80_UnprovableCoverageCarriesAReasonFromTheClosedSet` fails: `graph-unprovable-coverage carries reason "", outside the closed set [dimension-not-finite row-can-refuse tag-not-single-valued]`. The reason vocabulary is closed and append-only and 0006 owns it. C5's "not a new code" was true and irrelevant; the RDR never once mentions `reason`. Phase 2 stops. A hallway decision borrows `dimension-not-finite`. `TestReq81_RowCanRefuseMessageReadsAsAWithholdingNotAnError` passes by luck, and the shipped message now tells authors to "declare the domain" of a dimension that does not exist.

**Week 5 — the first real table returns an empty answer.** rdr#tmxk's navigator is rewritten as `class = "decision-table"`. One rule collapses four statuses with `[rule.match.recognized] in = ["locate","route"]`. `flow resolve` returns `{"rule":"draft-propose","emit":{}}`. `normalize.go::expand` builds each expanded row from an explicit nine-field seed literal — `ModelID, RuleID, SourceLocator, Source, Gate, RequiresOwned, Escape, setKeys` — and `Emit` was never added to it. C3 did say "`Row.Emit` MUST be carried through `expand`," but C3's *join-soundness argument* for keeping `rowByID` first-match is stated as already true ("every row a rule expands to carries the same block"), so the reviewer read it as a fact rather than an obligation. No test covers it: S3 tests normalization of a non-expanding rule, S5 tests resolve on the MVV fixture where every rule is non-expanding. The consumer's `rdr-status` skill takes the no-answer branch and reports `Next: (unknown)` for two days before anyone traces it to `expand`.

**Week 6 — the answer format leaks.** With emit fixed, the navigator needs `Next: /rdr-prelock 0046 critique`. `0046` is the RDR the caller asked about. The consumer authors `next = "/rdr-prelock {id} critique"` and does the substitution in the skill. C3 says emit values are "undeclared, uninterpreted, compared by exact byte equality," so intrastate has no opinion on `{id}` — no escaping rule, no validation, no dump rendering distinction. A second consumer picks `$ID`. The consumer-invented undocumented format the Problem Statement named as the defect is back, one seam over. A6 predicted this in its own evidence and filed it as "a widening, not a break"; the widening is BR3, which was rejected because "a dump column needs a deterministic rendering," and that reason has not changed.

**Week 8 — the bidirectional check goes in and comes back out.** A second engineer picks up a follow-up kata, reads `§approach` ("the loader refuses a class that disagrees with the owned set **in either direction**") and the Decision Rationale's winning scorecard cell ("Silent misclassification: refused at load both ways"), and adds the reverse arm: a `state-machine` model with zero owned tags refuses at load. `make check` goes red across `internal/graphlint` — every rootless-model fixture that 0006 relies on to produce `graph-dangling-edge` now fails at *load*, before lint runs — and MVV step 6's own class-omitted negative control becomes unwritable. C1 says all of this explicitly, in a parenthetical, three paragraphs in. The Approach, the scorecard, and Alternative 2's reason-for-rejection all say the opposite. Revert, plus a `//` comment in the loader that will outlive the RDR.

**Week 10 — a decision table lints green over nothing.** A four-outcome table ships with one "otherwise" escape row. A9 discovered that `0002:C5` scopes rescue per outcome and that an N-outcome table needs N escape rows — and routed the consequence to "the Phase 4 authoring docs must say" it. Phase 4 said it. The author did not read Phase 4. Three outcomes report coverage gaps; the author reads them as missing rows, adds ordinary rows until the gaps close, and ships a table where three of four outcomes have no rescue at all. `flow resolve` returns `flow-no-match` in production on a dimension combination the author believed was covered. Lint had no finding for it because A9's consequence got a doc where C5's got a `MUST`.

**Week 12 — the `len(Initial)` cleanup breaks the state-machine class.** §technical-design promised "every site that today asks `len(Initial) == 0` or would ask `len(owned) == 0` reads that accessor, so there is exactly one place the question is answered." A cleanup kata converts all five shipped sites — `reach.go:91`, `analysis.go:113`, `:233` (`checkAlwaysPresentOwned`), `:515` (`checkUnreachableRules`) — to the class accessor. C5 had only ever named two. `checkUnreachableRules` and `checkAlwaysPresentOwned` were silent over decision tables *because of the `len(Initial)` guard*, not for the reason A10 gave ("C5's single ∅ node satisfies every context vacuously" and "filter on `ProvenanceOwned` and find none"). Class-keyed, both go live: every decision-table row draws `graph-unreachable-rule` and every always-present owned tag check runs against a root that does not exist. The revert restores three `len(Initial)` guards, and the RDR's "exactly one place the question is answered" is now a comment explaining why there are four.

**Post-mortem finding.** Nothing in this record was under-researched. Thirteen assumptions, ten with executed or read-through evidence, an Oracle table that correctly identifies its own absence-of-error hazard, a Trace table that flags its one unwitnessed row. The failures were: two `Pending` assumptions carrying C5's load-bearing consequences past a lock the record's own gate forbids (`no assumption marked Pending … may have settled-fact prose elsewhere`); one census performed with a regex too narrow to see Go-side expectations and reported as "confirmed exact"; and one contract (C1) whose normative text contradicts three narrative passages including the scorecard cell the decision was made on. The document was not too small. It was precise in the places it looked and silent in the two places — `reason`, and `expand`'s seed literal — where the code disagreed with it.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

These are review-time gates, executable against shipped `main` before lock. Each names the finding it kills.

### AT-1 — the closed `reason` set (kills C-1)

```gherkin
Feature: a new graph-lint emission is admissible under 0006's finding record
  Scenario: C5's zero-dimension fence carries a legal reason
    Given RDR 0010 mandates a new emission of "graph-unprovable-coverage"
    When I read internal/graphlint/taxonomy.go::Reasons
    Then the reason set is closed, append-only, and has exactly 3 members
    And internal/graphlint/findings_0006_test.go::TestReq80_... asserts
        slices.Equal over that literal AND that every emitted finding of
        this code carries a reason from it
    Then RDR 0010 MUST name, in C5, which of the 3 reasons the new emission
        carries, and MUST justify why the reason's message contract
        (REQ-81: reads as a withholding, not an authoring error) still holds
    And if none of the 3 fits, A13's "without a taxonomy change" is FALSE
        and the RDR routes back to 0006 before lock
```

Plain step: for any RDR adding an emission of an existing finding code, enumerate every *discriminator field* that code carries (not just the code), check each against its owning RDR's closed set, and require the new RDR to bind a value for each. `grep -c 'reason' 0010.md` returning zero normative hits is the review-time smell.

### AT-2 — the census must cover every expectation medium (kills C-2, C-11)

```gherkin
Feature: a closed-vocabulary extension enumerates all its assertions
  Scenario: the emit column census
    Given A4 claims "Count confirmed exact: 103" from
          grep -rl '^\[dump\]' internal/table/testdata/
    When I re-run the census WITHOUT the ^ anchor, WITHOUT the extension
         filter, and across internal/**/*.go
    Then the result set MUST equal A4's, or A4 is wrong
    And the Done clause's licensed-diff rule MUST admit every file shape
        in the result set — including a Go `want := []string{...}` literal
        gaining a member and a TOML `order = [...]` line gaining a member
    And a rule that admits neither is self-contradictory and blocks lock
```

Plain step: run `grep -rl '\[dump\]' . | grep -v '^./.git'` and diff the directory histogram against A4's claim. `internal/table` (5 Go files) is absent from A4's accounting. Then: for each distinct file shape found, quote the Done clause and decide whether the required edit is licensed. Two shapes are not.

### AT-3 — every "provably silent" claim states the guard that makes it silent (kills C-3)

```gherkin
Feature: A10's silent bucket is grounded in the actual early-return
  Scenario: checkUnreachableRules and checkAlwaysPresentOwned
    Given A10 claims graph-unreachable-rule is silent "because C5's single
          ∅ node satisfies every context vacuously"
    And A10 claims graph-always-present-owned is silent because the check
          "filters on ProvenanceOwned and finds none"
    When I read internal/graphlint/analysis.go
    Then checkUnreachableRules returns early on len(a.model.Initial) == 0
    And checkAlwaysPresentOwned returns early on len(a.model.Initial) == 0
    Then both stated mechanisms are WRONG, and the real guard is a field
         C1 instructs the implementation to stop reading
    And A10 MUST be corrected to name the len(Initial) guard, and C5 MUST
        state that these two sites keep it
```

Plain step: `grep -n 'len(a.model.Initial)\|len(m.Initial)' internal/graphlint/*.go` returns five hits. Cross them against C5's two named sites. Three are unaccounted for. Any RDR that class-keys a discriminator must inventory *every* current reader of the old discriminator, not the ones it plans to change.

### AT-4 — "exactly one place" is falsifiable (kills C-4)

```gherkin
Feature: the single-source-of-truth claim is checkable
  Scenario: the class accessor
    Given §technical-design says "every site that today asks
          len(Initial) == 0 ... reads that accessor, so there is exactly
          one place the question is answered"
    When I enumerate the len(Initial) == 0 sites on main
    Then there are 4 in graphlint alone
    And C5 requires 2 of them to change and 2 to stay
    Then "exactly one place" is FALSE as written and MUST be replaced by an
         explicit per-site table: site | old key | new key | why
```

Plain step: the Authority table has four rows and names `reach`, `checkDanglingEdge` root arm, `normalizeRule`, `checkCoverage` zero-dim arm. It does not name `checkUnreachableRules` or `checkAlwaysPresentOwned`, both of which read the sibling discriminator the table's "Sibling arms" column identifies as the drift risk. The Authority table is the artifact that should have caught this and it under-enumerated.

### AT-5 — the payload field position is specified (kills C-5)

```gherkin
Feature: a new wire field has a fixed position
  Scenario: emit on resolvePayload
    Given C3 fixes the dump column position ("appended last, after escape")
          with an explicit rationale about fixture diff confinement
    When I read C4 for the payload
    Then C4 states no position for emit within resolvePayload's fields
    And Go's encoder emits JSON keys in struct field-declaration order
    Then C4 MUST fix the field's position the way C3 fixes the column's,
         or two implementations produce two byte-different payloads and
         the Done clause cannot adjudicate
```

Plain step: read `internal/cli/flow_resolve.go:31-46`. Eleven fields. C4 names five. Ask: after which one does `emit` go, and what breaks if it goes elsewhere? The RDR has no answer.

### AT-6 — the emission site is reachable (kills C-6)

```gherkin
Feature: a Pending assumption names an implementable site
  Scenario: A13's emitStructurallyUnprovable
    Given A13 says "verify whether emitStructurallyUnprovable's existing
          per-dimension shape admits a group-level emission"
    When I read internal/graphlint/coverage.go::emitStructurallyUnprovable
    Then its entire body is a loop over guard.Dimensions(a.model, g)
    And it is called only from checkCoverage's product-bound branch,
        which is BELOW the len(dims) == 0 early return
    Then the function is UNREACHABLE in the zero-dimension case and A13's
         question is unanswerable as posed
    And C5 MUST specify the emission site and its precedence against
        emitCoverageArms / bareEscapeFor / emitWithholdings
```

Plain step: trace `checkCoverage` top to bottom. The `len(dims) == 0` branch calls `emitCoverageArms` and returns — bypassing `emitUnprovableDimensions`, `emitWithholdings`, and the product-bound branch entirely. Then ask: on a match-discriminated table *with* an escape row, does the group get the fence, the closed-by-escape advisory, or both? Nothing in C5 answers it, and S4 asserts "exactly `graph-unprovable-coverage`."

### AT-7 — the consumer's answer is authorable as a literal (kills C-7)

```gherkin
Feature: A6's sufficiency is demonstrated, not asserted
  Scenario: the navigator's real answer
    Given A6 cites rdr-status rendering "Next: /rdr-prelock 0046 critique"
    When I ask whether "0046" is knowable at model-authoring time
    Then it is not — it is the caller's own request parameter
    Then A6 MUST exhibit the literal [rule.emit] block the consumer would
         author, end to end, with no substitution
    And if that block contains a placeholder of any syntax, A6 is FALSE:
        intrastate has no interpolation contract, C3 forbids interpretation,
        and the consumer-invented format the Problem Statement rejects
        has returned
```

Plain step: write the emit block. `next = "/rdr-prelock ??? critique"`. If you cannot fill the `???` with a literal, the assumption is refuted by its own example. A6 should be downgraded to `Unverified` and the widening decided before lock, not after `emit` has byte-asserted consumers.

### AT-8 — no Pending assumption underwrites a MUST (kills C-8)

```gherkin
Feature: the Finalization Gate's own status-consistency rule
  Scenario: A12 and A13 at lock
    Given the gate says "no assumption marked Pending or Unverified may
          have settled-fact prose elsewhere in the RDR depending on it"
    And A12 (dangling-edge arm split) and A13 (zero-dim fence) are Pending
    When I grep the RDR for prose depending on them
    Then C5 states the zero-dim fence as an unconditional MUST
    And §risks-and-mitigations calls it "a contract, not documentation"
    And §failure-modes lists it as one of the two "now guarded" paths
    And the Trace table's row 2″ admits it is "unwitnessed"
    Then the gate's own rule BLOCKS lock until A12 and A13 are Verified
```

Plain step: this test needs no source access. It is the RDR checked against its own gate text, and it fails. Run it before every lock.

### AT-9 — the contract and the narrative agree on directionality (kills C-9)

```gherkin
Feature: no narrative passage contradicts a normative clause
  Scenario: the class/owned-set agreement check
    Given C1 says "The agreement check is one-directional"
    When I grep the RDR for the same check described elsewhere
    Then §approach says "in either direction"
    And §decision-rationale's winning row says "refused at load both ways"
    And ALT2's reason-for-rejection says "in both directions"
    Then three passages contradict the contract, one of them being the
         scorecard cell the decision was made on
    And the RDR MUST NOT lock until all three read one-directional
```

Plain step: `grep -n 'both ways\|either direction\|both directions'` on the record. Three hits, all wrong, all outside the normative block. The scorecard is the worst of the three: the Decision Rationale says "C wins on the misclassification and overrides rows," and the misclassification row's winning cell states a property C1 disclaims.

### AT-10 — emit survives match-block expansion (kills C-10)

```gherkin
Feature: a row-carried field survives expand
  Scenario: an emit block on an expanding rule
    Given C3 argues rowByID's first-match join is sound because "every row
          a rule expands to carries the same block"
    When I read internal/table/normalize.go::expand
    Then each row is built from an explicit field-by-field seed literal
         (ModelID, RuleID, SourceLocator, Source, Gate, RequiresOwned,
          Escape, setKeys) — no struct copy, so a new field is dropped
         by default
    Then the Testing Strategy MUST carry a scenario with an emit block on a
         rule whose match block has a multi-member `in` atom, asserting the
         emit block is present on EVERY expanded row and on the resolve
         payload for each member
    And S3/S4/S5 as written contain no such case
```

Plain step: read the seed literal at `normalize.go:646`. Count its fields. Ask which of `Row`'s fields are absent from it and why. Then check whether any RDR test asserts an emit on an expanding rule. None does — S5's fixture is the MVV table, every rule non-expanding.

### AT-11 — the licensed-diff rule admits every required edit (kills C-11)

```gherkin
Feature: the mechanical Done rule is not self-contradictory
  Scenario: the 103-fixture sweep
    Given the rule admits "a trailing ` emit=[…]` dump cell or an
          `\"emit\":` payload member" and nothing else
    And A4 requires 103 files to gain `emit` inside `order = [...]`
    When I classify that edit under the rule
    Then it is neither shape, so all 103 required edits are "regressions"
    And A4's claim that they "are licensed diffs under the rule" is not
        supported by the rule's own text
    Then the rule MUST be amended to a third admitted shape before lock,
         or it will be disabled during the sweep and the regression it
         exists to catch will ship
```

### AT-12 — a discovered authoring hazard gets a contract, not a doc (kills C-12)

```gherkin
Feature: mitigation parity between discovered hazards
  Scenario: A9's N-escape-rows consequence vs C5's match-atom fence
    Given both are authoring mistakes that produce an unsound table
    And C5 gives the match-atom mistake a MUST-emit lint fence
    And A9 gives the per-outcome-rescue mistake a Phase 4 doc sentence
    When I ask what lint says about a 4-outcome decision table with one
         escape row
    Then it reports coverage gaps on 3 outcomes with no indication that
         the remedy is 3 more escape rows, not 3 more ordinary rows
    Then A9's consequence MUST either take a contract clause of its own or
         the RDR MUST state why lint is the wrong place for it — the
         asymmetry with C5 is unjustified as written
```

### AT-13 — the Problem Statement names the right dimension kind (kills C-13)

```gherkin
Feature: the "why" paragraph does not teach the mistake the contract fences
  Scenario: observed vs guard dimensions
    Given §problem-statement says coverage runs "over observed dimensions"
    And §approach says "over the observed dimensions"
    And §technical-design's data-flow says "groups over observed dimensions"
    When I read internal/guard/product.go::Dimensions
    Then dimensions are collected from guard.all / guard.unless atoms only,
         and provenance ("observed") is orthogonal to block ("guard")
    And C5 exists precisely to fence a table that discriminates with match
        atoms — which may be over observed tags
    Then all three passages MUST say "guard dimensions", or the record's
         top-level statement of purpose teaches the exact error C5 fences
```

Plain step: substitute "guard" for "observed" in the three passages and confirm the sentences still parse. They do — which means the error was mechanical, undetected across a grounding sweep, a cove pass, and a 3amigo pass, and sits in the first paragraph anyone reads.
