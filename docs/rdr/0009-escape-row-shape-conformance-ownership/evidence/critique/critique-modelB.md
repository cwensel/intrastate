Model: claude-opus-5

# Hostile critique — RDR 0009, escape-row shape conformance ownership

This RDR is 1831 lines of contract text specifying a two-line predicate
(`len(row.Escape) != 0 && len(row.Writes) != 0`) that guards a code path
the RDR itself proves nobody can reach (A5: no production constructor
exists; there is no `flow` verb; `internal/resolve` has zero production
importers). It ships a permanent, name-pinned export surface —
`ErrEscapeShapeBreach`, `*EscapeShapeBreachError{Ref, Count}`,
`Table.CheckValid`, a flat `errors.Join` aggregate with a per-identity
collapse and pre-collapse counting rule — onto an RDR-0001-locked package,
in exchange for deleting a `Writes` field from a test-fixture builder. The
document's own evidence (A3 spike, scenario 6 Mutant A) states in writing
that after Phase 2, *no test in the frozen suite can detect the removal of
the entire mechanism*. That is not a design; that is a specification
looking for a problem, and its most likely failure is not that it breaks —
it is that it lands, costs three days, damages the evidence base of RDR
0001, and guards nothing.

---

## 6. Findings ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
|----|-------------|--------------|-------------------|--------|
| C-1 | §Technical Design: "`Resolve` returns the zero `Result` (`Plan` and `Refusal` both nil …)"; scenario 7b | The zero-`Result`-plus-error return directly contradicts `resolve.go::Result`'s locked doc contract ("never both and never neither") and REQ-1's cardinality test `TestReq1_DispositionIsExactlyOneOfPlanOrRefusal`. The RDR notices the caller hazard and pins it as a test instead of recognizing it as a contract breach of a Final peer, then routes around it by observing the frozen test happens not to cover the case (`mustResolve` fatals first) | On the first real `flow` verb, a caller written against the documented "exactly one disposition" invariant reads `Refused() == false` on a failed call and treats a broken table as a success with a nil `Plan` → nil-pointer panic, exit 2 with a Go stack instead of the promised `escape-row-shape-breach` envelope | §1, §3, premortem, AT-1 |
| C-2 | §Validation scenario 6 *Mutant A*: "`Table.CheckValid` returns `nil` unconditionally. **Expected**: the frozen suite still passes" | The RDR states as a design property that after Phase 2 the entire mechanism is undetectable by the pre-existing suite, and then declares the tautological self-tests (scenarios 1–4, 7, 9, 10, 10b) excluded from the oracle. Nothing outside the RDR's own tests protects the invariant. Combined with Phase 2 stripping the only breaching fixtures, the RDR *removes* the population that would have exercised its own check | Six months on, a refactor drops or short-circuits `CheckValid` and CI stays green; the invariant silently ceases to exist while the RDR still claims it "cannot be emitted" | §1, §2, premortem, AT-2 |
| C-3 | §Approach: "backstopped by the kernel as an entry precondition"; A5 "No production code path constructs `resolve.Row` values at HEAD" | The backstop's entire addressable population is test fixtures, which Phase 2 conforms away. The RDR knowingly builds a runtime guard for a producer set of size zero and defers the only guard the Problem Statement's user needs (RDR 0002 normalizer, Final but unimplemented) to Phase 3 "binding fixtures" — text, not enforcement | The table author from the Problem Statement gets nothing. They author `escape` + `writes` in TOML, and until RDR 0002 is built there is no load-time diagnostic and no runtime path either, because no verb calls `Resolve` | §1, §3, premortem, AT-3 |
| C-4 | A3 Evidence: "every `Plan.Writes` read on an escape-derived plan … sits in a `t.Fatalf` *format argument* … never in a condition" + Phase 2 "drop the two call-site overrides — `adversarial_test.go:229` … `fixup_test.go:103`" | Phase 2 strips the exact `escape.Writes = {status Escaped}` values that make ADV-2b's and Fixup-1e's failure messages *legible*. ADV-2b exists to prove the kernel does not "describe a write off absent owned state"; after Phase 2 its failure output prints `writes []`, deleting the evidence a future regression-triager needs. A3 measured verdicts, not diagnostic value, and explicitly waves this off as "degrades the failure *message*" | A future kernel regression fires ADV-2b; the maintainer reads `kernel emitted a plan for rule "rdr.escape.needsowned" describing writes []` and cannot tell whether the guessed-persistence class recurred | §1, §2, premortem, AT-4 |
| C-5 | §Normative Contracts (multi-breach clause): "breaching rows sharing one `RowRef` contribute ONE reported error"; "Count is PER-IDENTITY and counts PRE-COLLAPSE ROWS" | Collapse-plus-`Count` is a whole invented aggregation algebra to solve an ordering problem that only exists because the RDR insists on reporting *every* breaching row. It is derived, not observed: no consumer exists (A5, no `flow` verb), so the semantics of `Count` are unfalsifiable and pinned in a locked package's exported struct | A producer author fixing a degenerate table sees `RowRef{"",""} Count=3` and cannot locate any of the three rows — the report's own degenerate case, acknowledged in the RDR and shipped anyway | §1, §2, AT-5 |
| C-6 | A2 "Second obligation (found at Pre-Lock)" + A9 **Status: Pending** + scenario 11 "That choice is **pending A9**" | A load-bearing wire-format decision (a new `omitempty` field on `clierr.CLIError`, contradicting the shipped `Cause` doc comment "the wire-visible cause surface is Detail") is written into a `normative` block while its verifying assumption is still `Pending`. The Finalization Gate's own Status-consistency rule ("no assumption marked `Pending` … may have settled-fact prose elsewhere") is violated by the RDR's own Normative Contracts | Either RDR 0005's envelope contract is reopened after lock, or the identities ride `Detail` as prose — the exact re-parse this RDR forbids everywhere else — so operators parse strings to find the broken row | §1, §2, premortem, AT-6 |
| C-7 | §Load-Bearing Decisions *Placement*: "the call is `Resolve`'s first statement, above the existing `view := assemble(in)`"; Normative Contracts: "precedes every modeled disposition" | Placing the whole-table check before the alphabet check makes a *malformed-but-irrelevant* row outrank `unmodeled_outcome`, changing the disposition of tables where the recognized outcome is not even in the alphabet. The RDR calls placement "a cheapness preference, not an observable contract" in one paragraph and pins precedence over all five kinds as normative in another — these are not the same claim, and the RDR conflates them | A caller resolving an unmodeled outcome against a table with one dormant bad row gets `escape-row-shape-breach` instead of `unmodeled_outcome`, and the pre-existing refusal semantics silently change class | §1, §4, AT-7 |
| C-8 | §Normative Contracts: "Resolve MUST return CheckValid's error VERBATIM, never `fmt.Errorf`-wrapped"; "`Unwrap() []error` MUST be EXACTLY ONE LEVEL deep" | The RDR pins internal error-plumbing shape as normative contract on a package whose only consumers are its own tests. This over-constrains the implementer into brittle mechanics (`errors.Join` returns a wrapper even for one element, so callers must never type-assert) and pins that hazard in prose rather than eliminating it with a purpose-built aggregate type | An implementer writes `err.(*EscapeShapeBreachError)`, it compiles, it never matches, and the CLI reports the generic internal error with no row identity | §1, premortem, AT-8 |
| C-9 | §Metadata **Profile**: "foundational"; §Finalization Gate Proportionality: "the sole author of at most one independent load-bearing contract" | The RDR authors at least four independently load-bearing contracts: the producer obligation, the kernel entry precondition + error surface, the multi-breach aggregation/collapse/`Count` algebra, and the CLI envelope extension (`escape-row-shape-breach` code + new serialized field). By its own split test it should have been split or trimmed; instead the Finalization Gate section is an unfilled template | Reviewers cannot re-derive which clause is the decision and which is implementation detail; the next RDR at this seam cites 0009 and inherits the whole 1831 lines | §2, §4 |
| C-10 | A4 *Carried forward*: "this safety rests entirely on RDR 0004's write-accessor scoping, not on any kernel-local property"; Prerequisites: "RDR 0004's implement stage keeps the write accessor scoped to *planned owned-tag writes*" | The stated user outcome ("no tag value can appear in owned state that no authored rule set") is guaranteed by a peer RDR that is Final-but-unimplemented, not by anything this RDR ships. `planOf` still copies `NextTags` symmetrically with zero escape-awareness. The RDR converts a real risk into a checkbox in another RDR's Prerequisites and calls the risk "Retired at Resolve" | If RDR 0004's implementer reads "apply the plan" broadly, an escape's `NextTags` persist and the exact Problem-Statement symptom appears — with RDR 0009 marked Implemented and its tests green | §3, §4, AT-9 |
| C-11 | §Trade-offs: "a hand-built table carrying a *dormant* malformed row … previously resolved fine and now errors on every `Resolve` (whole-table precondition)" | A deliberate, undiscussed-with-users behavior break justified entirely by A5 ("today on nobody"). A5's protection is temporal, not structural: it expires the moment anyone writes a table producer, and nothing in the RDR re-runs A5 at that point | The first non-TOML producer (a test harness, a migration script, a fixture generator) finds every previously-working table now rejected on rows their code path never touches | §3, AT-10 |
| C-12 | §Validation scenario 5: "The comparison baseline is the A3 spike's **conformed** run (`a3-conformed.out`: 154 PASS, 0 FAIL), restricted to the tests existing at that revision" | The Done criterion is a diff against a spike artifact from a specific HEAD (`9678601`, branch `via-claude`) captured before implementation. Any intervening commit to `internal/resolve` invalidates the oracle, and the RDR gives no procedure for regenerating it | The implementer cannot reproduce the 154-PASS baseline, declares the difference "unrelated," and the one scenario that proves conforming tables did not change disposition is quietly downgraded to a smoke test | §2, AT-11 |

---

## 1. The three most likely ways implementation goes wrong

### 1.1 The zero-`Result` return breaks RDR 0001's locked cardinality contract, and the RDR pins the breach as a test instead of noticing it

**Root cause.** The RDR needs `Resolve` to return "nothing" alongside an
error, and Go forces it to return *some* `Result`. It picks the zero value
and reasons about the ergonomics — but `Result`'s own doc contract in
`internal/resolve/resolve.go` reads:

> `Result` is the kernel disposition: exactly one transition plan or
> exactly one typed refusal, **never both and never neither**.

And `resolve_test.go::TestReq1_DispositionIsExactlyOneOfPlanOrRefusal`
enforces it as a frozen boundary test, with the explicit error branch
`"disposition carries neither a plan nor a refusal"`. The RDR asserts in
§Overrides that it "does **not** reopen RDR 0001's closed refusal taxonomy"
and that "Neither peer's normative surface is reopened." That is false. It
reopens REQ-1's cardinality clause. The only reason the frozen test does not
turn red is an accident of the test harness: `TestReq1_...` calls
`mustResolve`, which `t.Fatalf`s on a non-nil error before the cardinality
switch executes. The RDR *observes this accident* — scenario 7b says "The
frozen helpers happen to be safe (`mustResolve` fatals on error first), so
no existing test covers this" — and treats surviving-by-accident as
permission.

**Enabling passage.** §Technical Design:

> When any row carries both a non-empty `Escape` and a non-empty `Writes`,
> `Resolve` returns the zero `Result` (`Plan` and `Refusal` both nil —
> `Result` is a struct, so there is no nil `Result` to return) and a
> non-nil error naming every such row. Because the zero `Result` reports
> `Refused() == false`, callers MUST check the error before reading the
> disposition; a caller that branches on `Refused()` first would read a
> success-shaped value with a nil `Plan`.

A design that has to write "callers MUST check the error before reading the
disposition" in a normative document, about a type whose documented
invariant is "never neither," has designed a trap and then documented the
trap rather than removing it. The RDR had at least two exits — return a
`Result` that is structurally impossible to misread, or amend REQ-1's
cardinality clause honestly as this RDR's second override — and took
neither.

**Symptom.** The first `flow` verb is written by someone reading
`resolve.go`'s type docs, not this RDR. They write the natural Go shape:

```go
res, err := resolve.Resolve(in)
if res.Refused() { return refusalToCLIError(res.Refusal) }
// … res.Plan
```

`Refused()` returns false on a breach. Control falls to the plan branch.
`res.Plan` is nil. The user runs `intrastate flow advance` against a table
built by a not-yet-conforming producer and gets a runtime nil-pointer panic
— exit 2 with a Go stack trace on stderr, no `CLIError` envelope, violating
the AGENTS.md never-silent contract the RDR cites approvingly. The
`escape-row-shape-breach` code, the `Hint`, the `RowRef` identities: none
of them reach the user, because the RDR's whole surfacing story assumes a
caller that checks `err` first, which the type's own documentation does not
teach them to do.

### 1.2 Phase 2 deletes the check's own test population, and the RDR writes down that this makes the mechanism undetectable

**Root cause.** The RDR's implementation plan has an internal contradiction
that its Validation section states out loud and then treats as a virtue.
Phase 1 adds `Table.CheckValid` and wires it into `Resolve`. Phase 2 strips
`Writes` from `fixtures_test.go::escapeRow` (line 257) and the two
overrides at `adversarial_test.go:229` and `fixup_test.go:103`. After Phase
2 no pre-existing fixture breaches. Therefore the frozen suite — the entire
evidence base the RDR keeps invoking as its safety net — cannot observe
whether `CheckValid` exists, runs, or returns anything at all.

**Enabling passage.** §Validation scenario 6, *Mutant A*:

> `Table.CheckValid` returns `nil` unconditionally. **Expected**: the frozen
> suite still passes. This is the honest result, and it is the point: after
> Phase 2 no pre-existing fixture breaches, so the frozen suite alone cannot
> detect the check's removal. Recording this refutes the weaker "at least one
> test fails" claim rather than hiding behind it.

Honesty about a defect is not a mitigation of the defect. And the RDR
compounds it: scenario 6 *excludes* this RDR's own scenarios 1–4, 7, 9, 10,
10b from the oracle on tautology grounds. So the RDR simultaneously
declares (a) the frozen suite cannot detect the mechanism's removal, and
(b) the only tests that can are excluded from the non-vacuity argument. The
net protection against silent deletion of the invariant is zero. Mutant B
("restore `Writes` on the builder, suite fails") proves only that the check
fires when handed the fixture Phase 2 just removed — it proves the check is
wired, not that anything continues to protect it.

**Symptom.** Not immediate. Six weeks to six months out, someone
benchmarking, or simplifying, or resolving a merge conflict, removes or
short-circuits the `CheckValid` call at `Resolve`'s first statement. CI is
green. The RDR is still marked Implemented. The Trade-offs section still
reads "no write-bearing escape plan can be emitted by any producer path."
The user-visible symptom arrives much later, in exactly the form the
Problem Statement opens with — an owned tag with a value no authored rule
set — and the RDR that claims to prevent it will be cited as evidence that
it cannot be happening.

### 1.3 The RDR ships a permanent exported surface and an invented error algebra for a consumer population of zero

**Root cause.** The RDR pins four exported names, two struct fields, an
`Unwrap() []error` depth constraint, a verbatim-return prohibition, an
identity-collapse rule, and pre-collapse `Count` semantics — onto
`internal/resolve`, a package RDR 0001 locked, with zero production
importers (A5, verified) and no `flow` verb (A2, verified). Every one of
those decisions is derived from imagined consumers. §Normative Contracts
justifies the pinning as necessary:

> The surface is fixed here, not left to the implementer, because it lands
> permanently on RDR 0001's locked package surface and every test below
> asserts through it

That is precisely backwards. Landing permanently on a locked surface is the
argument for shipping the *minimum* — one unexported check and one
`error` — until a consumer exists and tells you what it needs. Instead the
RDR reasons its way to `Count int` on a per-identity collapsed error, a
value whose only stated use case is a degenerate table of identity-less
hand-built rows that the RDR itself calls "a degenerate producer."

**Enabling passage.** §Normative Contracts, multi-breach clause:

> Count is PER-IDENTITY and counts PRE-COLLAPSE ROWS: three breaching rows
> sharing `RowRef{"",""}` yield ONE reported error with `Ref ==
> RowRef{"",""}` and `Count == 3`. A single-row breach carries `Count == 1`,
> so the field is uniform rather than present only in the degenerate case.

Read the user journey this serves. A programmer hand-builds a table in a
test, forgets to set `RuleID`/`SourceLocator` on three rows, and puts
`Writes` on the escape rows. They get back an error saying: one identity,
`{"", ""}`, count 3. They now know three rows are wrong and have no way to
find any of them, because the report deliberately discarded the only
information that could locate them — position — on REQ-2 grounds. The
RDR's answer ("Identity is what the producer fixes by") is false for the
exact population the RDR names as this precondition's "whole target
population (A5)." Hand-built rows *do not have* identity. The RDR argues
itself into a report that is maximally uninformative precisely where it is
most likely to fire.

**Symptom.** The implementer builds all of it — sentinel, typed error, flat
join, collapse, count, ordering — and it costs a day. Then the first real
consumer (a `flow` verb, a year later) needs something the RDR did not
imagine: per-row entries, or the source position, or a single first error
for a log line. Now they are modifying a name-pinned, normatively-frozen
surface in a locked package, and every change is an RDR amendment.

---

## 2. The one section that will be rewritten within 6 weeks of shipping

**§Proposed Solution → Normative Contracts, the multi-breach reporting
clause** (lines 818–875: `errors.Join` aggregation, `compareRefs` ordering,
identity collapse, per-identity pre-collapse `Count`), together with its
dependent §Normative Contracts error-surface clause (lines 765–816: the
four pinned spellings, the one-level-deep `Unwrap()`, the verbatim-return
rule).

**Why it gets rewritten.**

1. **It is the only section with no ground truth.** Every other clause is
   anchored in something that exists: `rescues`, `len(Escape) != 0`, the
   reserved error return, `GroupInternal`. The aggregation algebra is
   anchored in nothing — no consumer, no test outside this RDR's own, no
   analogous surface elsewhere in the repo. It is pure derivation from
   REQ-2/REQ-10, and derived contracts survive contact with reality worst.

2. **It solves a self-inflicted problem.** Aggregate-not-fail-fast forces
   an ordering question; ordering forces the `compareRefs` tiebreak
   question; `compareRefs` not being total forces the collapse; collapse
   losing multiplicity forces `Count`. Four layers of contract, each
   necessitated by the previous one, all descending from a single
   discretionary choice (report every breach) that serves an audience — "a
   table author fixing rows" — that this RDR's Phase 1/2 scope does not
   have, because table authors are on the RDR 0002 path (Phase 3, deferred,
   fixtures only).

3. **The moment a real caller appears, at least one of its pins breaks.**
   `Count` on a collapsed identity is not what a CLI wants to render; the
   RDR admits the wire path for identities is unsolved and `Pending` (A9).
   The flat-`Unwrap()` and verbatim-return rules will collide with the
   first `fmt.Errorf("resolving %s: %w", flow, err)` an implementer writes
   at the verb boundary, which is idiomatic Go and which this clause
   forbids at one level and does not govern at the other.

4. **It is the section a maintainer most wants to delete.** It is the
   longest and least load-bearing part of the RDR, and it exists at all
   because §Load-Bearing Decisions chose aggregate reporting on a prior-art
   analogy (`fstest.TestFS`) to an audience that does not exist yet. The
   cheapest correct version — return the first breach, or return an
   unexported error with the row identity, and widen when someone asks — is
   two lines and no contract.

Runner-up: **§Validation scenario 6**, which will be rewritten because
Mutant A's honest result ("the frozen suite still passes") makes the
scenario an admission rather than a test, and the first person to read it
during a regression investigation will replace it with a real guard.

---

## 3. The one assumption that will not survive first contact with a real user

**A5**: *"No production code path constructs `resolve.Row` values at HEAD —
the only constructors are test fixtures — so the entry precondition changes
no shipped behavior and needs no migration."*

A5 is **verified and correct today, and load-bearing for a future it does
not govern.** It is not an assumption about the design; it is an
observation about a repository that has not been built yet. The RDR leans
on it structurally, in at least five places:

- §Trade-offs: the dormant-row behavior break is acceptable because "today
  on nobody (A5: no production constructor)."
- §Trade-offs: "A5's protection expires exactly when authored tables
  arrive, which is also when the load-time diagnostic that makes the
  refusal actionable arrives, so the strictness never lands without its
  remedy."
- §Normative Contracts multi-breach clause: "Hand-built rows are this
  precondition's whole target population (A5)."
- §Decision Rationale: A is rejected because it "leaves open the only path
  that can [misbehave]" — the hand-built path A5 documents.
- §Context: "Reachability today is narrow."

That third quotation is the load-bearing claim, and it is the one that
breaks. A5 says the target population is hand-built rows. Phase 2 then
*conforms every hand-built row in the repository*. So at RDR completion the
target population is empty. The RDR has built a runtime guard whose
addressable set it deliberately emptied in the same change.

The second quotation is the one that will not survive first contact.
"A5's protection expires exactly when authored tables arrive, which is also
when the load-time diagnostic arrives" assumes the two events are
simultaneous. They are not, and the RDR's own Capability Dependencies table
says so: RDR 0002 is `Final, unimplemented` / `Deferred`. The realistic
sequence is: someone builds a `flow` verb, or a fixture generator, or a
migration harness, or a benchmark table builder — any non-TOML producer —
*before* the RDR 0002 normalizer exists, because the normalizer is a much
larger job. At that instant A5 is false, there is no load-time diagnostic,
and the RDR's promise that "the strictness never lands without its remedy"
fails. The user gets the whole-table strictness (C-11: dormant rows in
tables their code path never touches now error) with none of the
authored-path diagnosis, delivered through the caller trap of C-1.

The RDR knows A5 is temporal — the "If wrong" clause says "An existing
production producer could begin erroring at `Resolve` entry; a migration
step and an RDR 0005 surfacing review become prerequisites" — but nothing
in Prerequisites re-runs A5, and nothing gates the strictness on RDR 0002's
existence. The assumption is verified at a point in time and consumed as if
it were an invariant.

---

## 4. Premortem

*Written at draft time, as if the failure has already happened.*

RDR 0009 shipped in three days. Phase 1 added `Table.CheckValid() error`,
`ErrEscapeShapeBreach`, `*EscapeShapeBreachError{Ref RowRef; Count int}`,
and a `CheckValid()` call as the first statement of
`internal/resolve/resolve.go::Resolve`, above `view := assemble(in)`. Phase
2 removed `Writes:` from `fixtures_test.go::escapeRow` (line 257) and the
two call-site overrides at `adversarial_test.go:229` (ADV-2b) and
`fixup_test.go:103` (Fixup-1e). Phase 3 wrote a fixture file for a
normalizer that does not exist. The suite went green at 154 pre-existing
PASS plus eleven new ones. The RDR was marked Implemented. Everyone
believed the invariant now held.

**Four months later, the `flow` verb landed** as part of the RDR 0005
build. Its author read `resolve.go`, not RDR 0009. They read the `Result`
doc contract — "exactly one transition plan or exactly one typed refusal,
never both and never neither" — and wrote the shape that contract teaches:

```go
res, err := resolve.Resolve(in)
if err != nil { /* handled below */ }
if res.Refused() { return respond.Fail(cmd, refusalToCLIError(res.Refusal)) }
return respond.OK(cmd, planPayload(res.Plan))
```

The `err` branch got written last and got the generic treatment: wrap into
`clierr.CLIError{Code: "resolve-error", Group: GroupInternal}`. Nobody
wired `escape-row-shape-breach`, nobody added the `omitempty` identity
field to `clierr.CLIError` — A9 was still `Pending` at lock, and a Pending
assumption is not an implementation task anyone tracks.

**Two months after that**, the flow-table generator for the kata harness
landed: a small Go program that builds `resolve.Table` values from kata
metadata without going through TOML. Its author copied the shape of
`fixtures_test.go` — including, on one escape row, a `Writes` entry, because
the row's purpose was to mark a flow blocked when nothing matched. Exactly
the Problem Statement's mental model: *the escape names where control goes*
— and the author wrote the write because they wanted the tag set.

The first user ran `intrastate flow advance`. `Resolve` called
`CheckValid`, which found the breach and returned an `errors.Join`
aggregate. `Resolve` returned `Result{}` and that error. In the verb, `err
!= nil` was true — and the `if err != nil` branch was three lines below the
`res.Refused()` branch, because that is how the author had laid it out
while reading the type docs. `res.Refused()` returned `false`. Control fell
through to `planPayload(res.Plan)`. `res.Plan` was nil.

**The user saw a Go panic**, exit 2, a stack trace on stderr, and no JSON
envelope — the never-silent contract in AGENTS.md broken at the exact
boundary RDR 0009 spent nine hundred words specifying. The row identity
that `*EscapeShapeBreachError.Ref` had faithfully carried out of
`CheckValid` never reached anyone: `CLIError.Cause` is `json:"-"`, A9 was
never resolved, and in any case the error value never got wrapped.

**Triage went to `internal/resolve`.** The maintainer ran the frozen suite:
154 + 11 PASS. They ran ADV-2b
(`TestAdv2b_EscapeEdgeMustNotBypassTheOwnedStateRequirement`) specifically,
because "kernel guessing at persistence" was the named class. It passed. To
check whether the escape path had regressed at all, they mutated
`escapeOrRefuse` — `viable := escapes` — and watched eleven assertions
fail, exactly as the A3 spike had recorded. Then they read ADV-2b's failure
message format string:

```go
t.Fatalf("kernel emitted a plan for rule %q describing writes %v; want "+
    "owned_state_unavailable — the escape edge required owned tag %q, "+
    "which the accessor snapshot does not carry",
    got.Plan.RuleID, got.Plan.Writes, "never-present")
```

Phase 2 had stripped `escape.Writes = {status Escaped}` from line 229. The
message that had been the repository's single clearest statement of the
guessed-persistence class now printed `describing writes []`. The A3 spike
had classified this as harmless — "Stripping `Writes` degrades the failure
*message* … but changes no verdict" — because A3 measured pass/fail
outcomes, which is the wrong oracle for a diagnostic.

**Then they checked whether `CheckValid` was even still there.** It was,
but the check took thirty minutes, because nothing in the suite depended on
it: Mutant A had told them in writing that `Table.CheckValid` could return
`nil` unconditionally and the frozen suite would still pass. The RDR's own
Validation section was the reason they could not answer the question by
running tests.

**The RDR 0002 normalizer still did not exist.** The Problem Statement's
table author — the one who "writes an escape rule to say when this failure
class happens, leave by this route" — had never received a single
diagnostic from this RDR, at load time or otherwise, because the authored
path was Phase 3 and Phase 3 shipped fixtures for an unimplemented builder.

**The second incident** came from `planOf`. RDR 0009 had guaranteed
`Writes` empty on escape rows and, per A4, explicitly declined to guard
`NextTags`. `planOf` still copies both symmetrically with no
escape-awareness:

```go
NextTags: copyTags(row.NextTags),
Writes:   copyTags(row.Writes),
```

The RDR 0004 accessor implementer, reading "apply the plan," applied next
state as well as planned writes. The escape row's
`NextTags: {status Blocked}` — still set on `fixtures_test.go::escapeRow`
after Phase 2, because "NextTags stays on the builder: A4 did not widen the
predicate" — persisted. An owned tag took a value no authored rule set. The
exact symptom of the Problem Statement, occurring in a repository where RDR
0009 was marked Implemented, its tests green, and its Trade-offs section
reading "no write-bearing escape plan can be emitted by any producer path."
The mitigation had been recorded as a checkbox in *another RDR's*
Prerequisites, which no one read while implementing that other RDR.

**Post-hoc accounting.** Three days of implementation. Four names
permanently on RDR 0001's locked package surface. A collapse-and-`Count`
algebra with zero consumers. Two frozen adversarial tests with their
diagnostic payloads deleted. One REQ-1 cardinality contract quietly
violated. Zero user-visible breaches prevented — because at the moment the
mechanism could have prevented one, the mechanism's caller could not read
its output.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

These are RDR-review gates, not implementation tests. Each is written to be
runnable against the *document* plus the existing source tree.

### AT-1 — the zero-`Result` return contradicts a frozen contract (C-1)

```gherkin
Given the RDR specifies a return value for Resolve on the breach path
When I extract the (Result, error) pair the RDR specifies
And I evaluate it against internal/resolve/resolve.go::Result's doc contract
  ("never both and never neither")
And against resolve_test.go::TestReq1_DispositionIsExactlyOneOfPlanOrRefusal
Then the RDR MUST either
  (a) specify a Result satisfying the cardinality invariant, or
  (b) list REQ-1's cardinality clause in §Overrides as reopened
And the RDR MUST NOT justify the return by observing that the frozen test
  happens not to reach the assertion (mustResolve fatals on error first)
```

Fails today. §Technical Design chooses (c): document the trap in prose and
add scenario 7b. Neither (a) nor (b) appears in §Overrides, which claims
"Neither peer's normative surface is reopened."

### AT-2 — the invariant must be protected by something outside this RDR's own tests (C-2)

```gherkin
Given the RDR proposes a runtime check
When I read §Validation scenario 6 Mutant A
Then the RDR MUST NOT state that the pre-existing suite still passes with
  the check neutralized
Or, if it does, the RDR MUST specify a durable guard the check's removal
  would trip that is not one of the RDR's own new scenarios
  (e.g. a retained breaching fixture, an exported-surface presence test in
   boundary_test.go, or a Phase 2 that conforms fixtures but keeps one
   deliberately-breaching table under a named test)
```

Fails today. Mutant A states the honest result; no durable guard is
specified; scenarios 1–4, 7, 9, 10, 10b are explicitly excluded from the
oracle as tautological.

### AT-3 — the guard must have a non-empty addressable population at completion (C-3)

```gherkin
Given A5 verifies no production constructor of resolve.Row exists at HEAD
And Phase 2 conforms every test-fixture constructor
When I enumerate the producers the kernel entry precondition can catch
  at the moment this RDR is Done (Phases 1–3 landed)
Then the set MUST be non-empty
Or the RDR MUST state that the check ships dormant and name the event that
  makes it live, with that event gated in Prerequisites
```

Fails today. The set is empty. §Trade-offs asserts the opposite — "at this
RDR's own completion the invariant holds everywhere a row can be
constructed" — which is true only in the vacuous sense that no row can be
constructed wrongly because no row is constructed at all.

### AT-4 — fixture conformance must not degrade adversarial diagnostics (C-4)

```gherkin
Given Phase 2 removes escape.Writes from adversarial_test.go:229 and
  fixup_test.go:103
When I read each t.Fatalf format string that consumed those values
Then for each, the RDR MUST show the failure message remains sufficient to
  identify the regression class the test names
  (ADV-2b: "the kernel guessing at persistence")
And the A3 spike MUST evaluate diagnostic sufficiency, not only pass/fail
  verdict equality
```

Fails today. A3 §a explicitly scopes to verdicts and dismisses the message
change: "Stripping `Writes` degrades the failure *message* (it would now
print `[]` instead of `[{status Escaped}]` if the kernel regressed) but
changes no verdict." ADV-2b's message becomes `describing writes []`, which
states the opposite of the class it was written to name.

### AT-5 — a diagnostic must locate the defect for its stated target population (C-5)

```gherkin
Given the RDR names hand-built rows as "this precondition's whole target
  population (A5)"
And observes that hand-built rows "collide trivially: rows built with no
  source identity all carry RowRef{"",""}"
When I walk the journey of a producer author receiving the report for a
  table with three such rows
Then the report MUST let them locate at least one offending row
```

Fails today. The report is one entry, `Ref == RowRef{"",""}`, `Count == 3`.
The RDR's stated remedy ("Identity is what the producer fixes by") is false
for the population it names as the target.

### AT-6 — no normative clause may depend on a Pending assumption (C-6)

```gherkin
Given A9 has Status: Pending
When I grep the ```normative blocks for claims that depend on A9
Then the set MUST be empty
```

Fails today. The CLI-wrapping normative block states "The chosen carrier is
a NEW omitempty field added under the type's own 'Extend with new optional
fields as needed' allowance — NOT the existing Detail," then adds "subject
to A9." A normative MUST subject to a Pending verification is not a
contract. The RDR's own Finalization Gate template forbids exactly this:
"no assumption marked `Pending` or `Unverified` may have settled-fact prose
elsewhere in the RDR depending on it."

### AT-7 — precedence over existing dispositions must be stated as a behavior change, not a placement preference (C-7)

```gherkin
Given the RDR places the check as Resolve's first statement, above the
  alphabet check that yields unmodeled_outcome
When I compare that placement's observable effect against the RDR's claim
  that placement is "a cheapness preference, not an observable contract"
Then the two statements MUST be reconciled
And the change in disposition for a table with a dormant breaching row and
  an out-of-alphabet recognized outcome MUST be enumerated in
  §Trade-offs → Consequences as a behavior change
```

Fails today. §Load-Bearing Decisions *Placement* calls it unobservable
("`assemble` is pure and allocation-only, so placing the check after it
would produce identical dispositions") — but the relevant comparison is not
against post-`assemble` placement, it is against *not having the check*,
and against the alphabet check that follows. §Normative Contracts pins
precedence over "every modeled disposition" as normative, and scenario 4
tests it per kind. The Consequences section mentions the dormant-row break
but never that `unmodeled_outcome` is now outranked.

### AT-8 — no error-plumbing mechanic may be pinned normatively without a consumer (C-8)

```gherkin
Given the RDR pins Unwrap() []error depth, verbatim-return, and a
  no-direct-type-assertion caller rule
When I identify the consumer whose requirements those pins serve
Then that consumer MUST exist at HEAD or be scheduled in this RDR's scope
```

Fails today. A5 and A2 jointly verify there is no consumer: no production
importer of `internal/resolve`, no `flow` verb. Scenario 11 — the only one
that would exercise the CLI side — is marked "DEFERRED: not executable by
this RDR."

### AT-9 — the stated user outcome must be guaranteed by something this RDR ships (C-10)

```gherkin
Given §Testing Strategy defines Done as "no tag value can appear in owned
  state that no authored rule set"
When I trace that guarantee to its enforcing artifact
Then the artifact MUST be shipped by this RDR
```

Fails today. A4's *Carried forward* clause states the guarantee "rests
entirely on RDR 0004's write-accessor scoping, not on any kernel-local
property." `planOf` still copies `NextTags` symmetrically. The RDR's own
Prerequisites acknowledge this by binding a checkbox onto RDR 0004's
implement stage — which is a promise about another RDR's future, not a
guarantee.

### AT-10 — a temporal assumption may not be consumed as an invariant (C-11)

```gherkin
Given A5 is an observation about HEAD, not a structural property
When I list every place the RDR relies on A5 to license a behavior break
  or a strictness decision
Then each MUST carry a re-verification trigger tied to the event that
  falsifies A5 (the first production resolve.Row constructor)
```

Fails today. Four reliance sites (Trade-offs ×2, Normative Contracts,
Decision Rationale); zero triggers. The claim "A5's protection expires
exactly when authored tables arrive, which is also when the load-time
diagnostic arrives" assumes RDR 0002 precedes every other producer, which
the Capability Dependencies table (`Deferred`) gives no reason to believe.

### AT-11 — the Done oracle must be reproducible after intervening commits (C-12)

```gherkin
Given scenario 5's baseline is a3-conformed.out captured at HEAD 9678601
When any commit touches internal/resolve between spike and implementation
Then the RDR MUST specify how the baseline is regenerated and what
  "restricted to the tests existing at that revision" means operationally
```

Fails today. No procedure is given. The spike doc pins branch `via-claude`,
HEAD `9678601`; the RDR says the oracle is "the sorted outcome set over
that pre-existing test set" without saying how to reconstruct that set.
