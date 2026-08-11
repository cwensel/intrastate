Model: claude-opus-5[1m]
Lens: 3amigo (persona 3 — QA / tester)

# Persona 3 — QA / Tester: top 5 untestable promises

Grounding: read the full RDR (1576 lines) and the frozen suite at HEAD —
`internal/resolve/{resolve,mvv,adversarial,fixup,boundary,fixtures}_test.go`
(2656 test lines) plus `internal/resolve/resolve.go` and
`internal/cli/clierr/clierr.go`. Suite assertion style, verified: every
disposition test routes through `resolve_test.go::mustResolve`, which
**`t.Fatalf`s on any non-nil error** ("Resolve returned a Go error for a
modeled disposition"); assertions are `reflect.DeepEqual` on whole values
(`Refusal`, `Result`) or field equality on `Refusal.Kind`/`.Guard`/`.Rows`;
`Writes` assertions are all four length- or DeepEqual-based
(`resolve_test.go:668`, `:672`, `:709`, `mvv_test.go:110`), none nil-sensitive;
zero tests assert on error *strings*; `boundary_test.go` does AST-level
package inspection (`kernelPackageImports`, `exportedKernelSymbols`).

---

## QA-1 — Scenario 6 (mutation check) has no stated mutant and no oracle

**Passage** — *Validation → Testing Strategy*, scenario 6:

> **Scenario**: Mutation check — revert the entry check and re-run.
> **Expected**: At least one test fails, proving the conformed fixtures did
> not leave the check untested.

Reinforced in *Phase 2*: "Because no remaining fixture breaches, the
mutation-style step … is load-bearing here rather than a formality — it is
what proves the new check is exercised at all."

**Severity**: high

**The gap**: As written this scenario is **vacuous by construction and
cannot fail**. Phase 2 conforms *all three* breaching fixture sites, so after
Phase 2 no pre-existing test constructs a breaching table. The only tests that
can fail when the entry check is reverted are the RDR's *own* new tests
(scenarios 1, 2, 4, 7, 9, 10, 10b) — which trivially fail, because they were
written to assert exactly that check. "At least one test fails" is therefore
satisfied by a tautology and proves nothing about coverage. The RDR states
the intent ("proves the new check is exercised at all") but never names:
(a) which mutant is applied (delete the `Resolve` entry call? make
`Table.CheckValid` return nil unconditionally? invert the predicate to
`len(Writes)==0`?), (b) which tests are *in scope* as the surviving oracle
(the A3 mutation precedent used a *specific* mutant — `viable := escapes` —
and a *specific* count, "still fails 11 assertions across ADV-2, ADV-2b,
Fixup-1d, Fixup-1e, and TestFixupGateIsUniform…", per
`evidence/spikes/a3-fixture-conformance.md`), and (c) whether excluding the
RDR's own new tests is required for the check to mean anything. Without
(a)–(c) the pass criterion is unfalsifiable.

Contrast: A3's own spike **did** state mutant + expected failure count. The
RDR's scenario 6 drops both, so it is strictly weaker than the evidence it
cites.

**Test it prevents**: `TestEscapeShapeCheckIsNonVacuous` — a mutation-harness
test (or a documented manual mutation procedure) that reverts a *named*
mutant and asserts a *named, non-self-referential* set of assertions fails.
As written I cannot decide what to revert or what must go red.

---

## QA-2 — Scenario 5's "still discriminates" has an observable with no expected value

**Passage** — *Validation → Testing Strategy*, scenario 5:

> **Scenario**: The full frozen `internal/resolve` suite (ADV/MVV/boundary)
> run against conformed escape fixtures.
> **Expected**: Every assertion still passes and still discriminates — no
> test silently becomes vacuous (A3).

**Severity**: high

**The gap**: "Still discriminates" and "no test silently becomes vacuous" are
**not checkable by running the suite**. `go test` reports pass/fail; it does
not report discrimination. The RDR gives no oracle for the second half of the
conjunction. A3 established discrimination *by a mutation run* — a separate
procedure with its own mutant — but scenario 5 does not say to run one, does
not name the mutant, and does not restate A3's "11 assertions" figure as the
expected value. Meanwhile the first half is checkable but **near-tautological
at HEAD**: A3 already ran it (154 PASS / 0 FAIL / identical sorted outcome
set, `a3-conformed.out`), so re-running post-Phase-2 asserts a fact already in
evidence.

Compounding: the RDR nowhere pins *which* run is the baseline for
comparison. `a3-baseline.out` is the **pre-conformance** tree; the
post-Phase-2 tree will additionally contain this RDR's new tests, so the PASS
count will not be 154 and the "sorted outcome set" will not be identical to
`a3-baseline.out`. Any test or CI check written literally against "identical
to the A3 baseline" fails for a reason unrelated to the contract. The RDR
does not say whether the comparison excludes new tests, and there is no
stated substitute expected value.

**Test it prevents**: `TestFrozenSuiteStillDiscriminatesUnderConformedFixtures`
— a CI gate comparing the post-Phase-2 run against a pinned baseline. I
cannot write it because the baseline artifact is stale-by-design and no
replacement expected value is given, and the "discriminates" half has no
oracle at all.

---

## QA-3 — Scenarios 9 / 10 / 10b: the report's asserted shape is contradicted between clauses (aggregate vs. traversal vs. `errors.As`)

**Passage** — three clauses that must be read together.

*Normative Contracts* (typed-error block):

> Because errors.Join returns a wrapper even for a single error, callers MUST
> classify and extract with errors.Is / errors.As / errors.AsType rather than
> a direct type assertion or equality against the sentinel.

*Failure Modes → Visible break (programmatic producer)*:

> A single `errors.As`/`AsType` call reports only the FIRST breach in the
> chain, so it is the wrong instrument for reading a multi-breach report — a
> test asserting on all offending rows must traverse the aggregate.

*Testing Strategy*, scenario 9:

> each per-row error is individually recoverable by traversing the aggregate's
> `Unwrap() []error` — a single `errors.As` finds only the first match, so the
> assertion must walk the slice.

*MVV fixture 1*:

> the error yields that row's `RowRef{"rdr.escape.needsowned",
> "flows/rdr.toml:90"}` through `errors.As`/`AsType` … — both through the
> `errors.Join` aggregate, which wraps even this single-breach case.

**Severity**: high

**The gap**: The extraction contract is stated three ways and they do not
compose into one assertion. Concretely, the RDR never fixes the **traversal
contract** the test must code against:

1. Is the aggregate's `Unwrap() []error` guaranteed to be **exactly one level
   deep** (each element being the per-row typed error), or may an element
   itself be a joined wrapper? A test that walks the slice and type-asserts
   each element is correct only under the flat guarantee, which is never
   stated.
2. `errors.Join` is named as the combinator, but nothing says whether
   `Table.CheckValid`'s single-breach return is *also* joined. MVV fixture 1
   says yes ("wraps even this single-breach case") and the Normative
   Contracts say yes ("MUST NOT return the bare per-row error in the
   one-breach case"), yet *Load-Bearing Decisions → Reporting* justifies the
   choice on the grounds that `errors.Join` "degrades to the single error's
   own message for one breach" — a *message* property, not a *structure*
   property. A test asserting `len(err.(interface{Unwrap() []error}).Unwrap()) == 1`
   for the single-breach case is required by the first two, and neither
   confirmed nor denied by the third.
3. Whether the *sentinel* is wrapped by each per-row error, by the aggregate,
   or by both, is unstated. `errors.Is(err, ErrEscapeShapeBreach)` passing on
   the aggregate does not tell me whether it also passes on each unwrapped
   element — which is exactly what scenario 9's per-row assertion needs.

Contrast with HEAD: the frozen suite's diagnostic assertions are whole-value
`reflect.DeepEqual` on `Refusal` (`fixup_test.go:200`, `:236`) precisely
because the payload shape is pinned. The breach report has no equivalent
pinned shape, so no `DeepEqual` is writable and the piecewise traversal is
under-specified.

**Test it prevents**: `TestMultiBreachReportNamesEveryRow` (scenario 9) —
I cannot write the traversal loop, because the number of unwrap levels, the
per-element type, and per-element `errors.Is` behavior are all indeterminate.

---

## QA-4 — Scenario 10b: "the breach count is stated alongside" has no assertable carrier

**Passage** — *Testing Strategy*, scenario 10b:

> **Expected**: Identical report under both permutations — equal identities
> collapse to one reported entry, and the breach count is stated alongside.

and *Normative Contracts* (multi-breach block):

> Rows carrying no source identity collapse to a single RowRef{"",""} entry;
> that is a degenerate producer, and the error MUST remain diagnostic in that
> case by stating the breach count alongside the identities.

**Severity**: medium

**The gap**: "Stated alongside" names no carrier and therefore no assertion.
Two readings produce **contradictory tests**:

- *Prose reading* — the count appears in the formatted message. But this
  directly contradicts the same RDR's own binding rule two blocks earlier:
  "Row identity MUST NOT be recoverable only from formatted prose" and
  scenario 1's "no assertion parses message text." If the count lives only in
  `Error()`, the only assertion available is a string match — the one
  instrument the RDR forbids everywhere else, and which appears **zero times**
  in the frozen suite (verified: no test in `internal/resolve` asserts on
  error text).
- *Structured reading* — the count is a field on the typed error (e.g.
  `Count int`). But the Normative Contracts specify the typed error's payload
  as "the offending row identity as the kernel's existing RowRef value" and
  nothing else; A8's "If wrong" explicitly frames a count field as the
  *rejected/contingent* branch ("If a consumer needs a per-row breach count
  rather than per-identity, the report must carry a count field instead of
  collapsing"), i.e. as an alternative to collapsing, not a companion to it.

So the RDR simultaneously mandates collapsing *and* a count, while its only
discussion of a count field treats it as the alternative to collapsing. The
expected value of `count` is also unstated in the collapsed case: for three
breaching rows sharing `RowRef{"",""}`, is the stated count 3 (pre-collapse
rows) or 1 (post-collapse entries)? Both are defensible from the text.

**Test it prevents**:
`TestCollapsedBreachReportStatesCountForDegenerateProducer` — I cannot name
the observable (field vs. message) and cannot pick between the two candidate
expected values (3 vs 1).

---

## QA-5 — Scenario 11 is unexecutable at HEAD and the RDR does not say so in the Testing Strategy

**Passage** — *Testing Strategy*, scenario 11:

> **Scenario** (binds the verb that first calls `Resolve`): a breach surfaced
> through the CLI.
> **Expected**: exit 2 via `CLIError{Group: GroupInternal}` carrying the
> stable code, the offending row identity, and the remedy `Hint` … an
> assertion that reads identities off the JSON output is what proves the wire
> carrier exists.

**Severity**: medium

**The gap**: There is **no `flow` verb at HEAD** — the RDR itself verifies
this twice (A2: "no `flow` verb exists at HEAD
(`internal/cli/root.go::NewRootCmd` registers only `newVersionCmd()`)"; A5:
same). So scenario 11 cannot be executed by this RDR's implementation at all.
The RDR *does* record the dependency, but only in **Prerequisites** (an
unchecked box) and in A2's "Obligation, not a gap" — **not in the Testing
Strategy**, where scenario 11 is listed flush with the eight executable
scenarios and carries no deferral marker. Scenario 8 gets the parenthetical
"(binds RDR 0002's build)"; scenario 11's parenthetical "(binds the verb that
first calls `Resolve`)" reads as scoping, not as "cannot run in this RDR."
The reader of the Testing Strategy cannot tell which scenarios are Done
criteria for *this* RDR versus obligations levied on a future one — and the
section opens "Done means: every producer path is closed…", implying all
eleven are in scope.

Second, unspecified observable: "the stable code" is never given a value.
`clierr.CLIError.Code` is a string (`internal/cli/clierr/clierr.go`), and the
shipped precedent the RDR cites uses a literal (`config.Load` →
`"config-read-error"`). The RDR names the breach informally ("escape-row
shape breach") and explicitly defers spelling ("the exact type/sentinel
spelling is sharpened at Pre-Lock"), but the *CLI code* string is not
enumerated anywhere. A test asserting `got.Code == ???` has no right-hand
side. Same for the `Hint` text — "a Hint stating the remedy" gives no
assertable value, and Failure Modes' remedy prose ("fix the producer — the
table value is broken; there is nothing to retry") is not offered as the
literal.

Third, the choice of wire carrier is left open in the same breath as the
assertion is mandated: "either the existing `Detail` (rendered from the row
identities …) or a new `omitempty` field". A JSON assertion must pick a key.
`{"detail": "..."}` requires parsing rendered prose out of a serialized
field — which the RDR forbids consumers from doing ("never re-parsed by any
consumer"), leaving the test in the same bind as QA-4.

**Test it prevents**:
`TestFlowVerbSurfacesEscapeShapeBreachAsExit2` — blocked three ways: no verb
to invoke, no expected `Code` literal, and no committed JSON key to read the
row identities from.

---

## Additional observations (below the top 5)

- **Promised in prose, never a scenario — precedence over *every* modeled
  disposition.** *Load-Bearing Decisions → Precedence* says "the breach error
  precedes every modeled disposition, `unmodeled_outcome` included," and the
  Normative Contracts say "it precedes every modeled disposition." Scenario 4
  covers **only** `unmodeled_outcome`. The other four kinds
  (`no_match`, `ambiguous_match`, `owned_state_unavailable`,
  `guard_unevaluable`) have no precedence scenario. `unmodeled_outcome` is
  the *easiest* case (it is the first branch of `Resolve`, at
  `resolve.go:321`); the sharp cases are the ones reached via
  `escapeOrRefuse`, where `gate` may return a `blocked` refusal from the
  breaching row itself. The MVV's five-kinds precedent
  (`mvv_test.go` leg 2, which enumerates all five and asserts
  `len(covered) != 5`) shows the house style is exhaustive-over-kinds; this
  RDR's precedence coverage is 1-of-5 with no stated reason.

- **Scenario 3's expected value is anchored to a moving baseline.**
  "resolves exactly as the A3 baseline" — but the A3 baseline
  (`a3-baseline.out`) is the *unconformed* tree, and scenario 3's fixture
  (`Writes: []Tag{}`) exists in *no* run captured by A3. The three A3 runs
  were: original, `Writes` stripped (i.e. **nil**), and `Writes`+`NextTags`
  stripped. The empty-non-nil variant was never executed. So "identically to
  the A3 baseline" cites evidence that does not cover the case, and the only
  independently assertable clause is `len(Plan.Writes) == 0` — which is
  satisfied by the nil case too and therefore does **not** discriminate
  empty-from-nil. The scenario meant to pin the nil-vs-empty decision cannot
  distinguish the two states it exists to separate. (Verified against
  `copyTags` at `resolve.go:569`: nil in → nil out; `[]Tag{}` in → `[]Tag{}`
  out, both `len == 0`. A test wanting to prove the *empty* case conforms must
  assert on the input row, not the plan — the RDR does not say to.)

- **Scenario 7's "verdicts identical" is under-specified for the error case.**
  "Verdicts identical to `Resolve`'s entry check — one predicate, two call
  sites, no drift." For the conforming tables that is `nil == nil`, fine. For
  the breaching ones, "identical" between a bare
  `Table.CheckValid() error` and `Resolve`'s second return value needs a
  comparison relation, and none is named. `reflect.DeepEqual` on two
  `errors.Join` aggregates is not obviously well-defined (joinError holds a
  slice; equality depends on element identity), and `errors.Is`-equivalence
  is much weaker than "identical." Given the frozen suite's habit of
  whole-value `DeepEqual` (`fixup_test.go:200`), a reader will reach for
  `DeepEqual` and may get a false green or a spurious red depending on
  implementation details the RDR does not pin. **Unverified**: I did not run
  a probe to establish whether `reflect.DeepEqual` over two independently
  constructed `errors.Join` results is stable here.

- **`Resolve`'s new caller contract has no scenario.** *Technical Design*
  states: "Because the zero `Result` reports `Refused() == false`, callers
  MUST check the error before reading the disposition; a caller that branches
  on `Refused()` first would read a success-shaped value with a nil `Plan`."
  This is a genuine, newly-introduced hazard — and at HEAD it is *already*
  live in the test helpers: `resolve_test.go::planOf` calls
  `got.Plan == nil → t.Fatal`, and `mustResolve` fatals on error first, so
  the frozen helpers happen to be safe. But no scenario asserts the zero
  `Result` shape (`Plan == nil && Refusal == nil && Refused() == false`) as a
  property. MVV fixture 1 asserts `Plan: nil, Refusal: nil` in passing; the
  `Refused() == false` consequence — the part that is actually novel and
  trap-shaped — is never an expected value anywhere.

- **No scenario for `Table.CheckValid` on a table with zero rows / nil Rows.**
  The predicate is a whole-table scan; `errors.Join` of nothing is nil, so
  presumably nil. Not stated, not scenario'd. Low severity — the frozen suite
  has an empty-input test (`resolve_test.go` REQ-20 region) so the shape is
  reachable.
