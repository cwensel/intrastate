Model: claude-opus-5

# Persona 3 — QA / Tester

Owned set read in full: `0029:C1`–`0029:C4`, `0029:S1`–`0029:S9`, `0029:MVV`.
WIDENED to `0029:A1`, `0029:A3`, `0029:A4` (the scenarios cite them as their
backing), to `0029:§performance-expectations`, and to the source tree at
`/Users/cwensel/sandbox/newcoinc/intrastate`. What sent me: four of the nine
scenarios name an *accessor symbol* as the thing under assertion, and a
by-value set assertion is only writable if that symbol enumerates. Checking
whether it does is not answerable inside `elements[]`.

Severity-ranked.

---

## F1 (BLOCKING) — `0029:S6` + `0029:C4`: three of the five newly-frozen vocabularies have no enumerable accessor, so the by-value assertion S6 demands cannot be written

`0029:S6` says: "a by-value assertion per frozen vocabulary in C4 —
including the five this RDR newly assigned", and names
`TestReq79_NoSixthRefusalKindIsMinted` as the pattern to follow. That
pattern (`internal/resolve/guard_atoms_test.go:894`) works because
`internal/resolve/resolve.go:67 RefusalKinds()` returns the set as a slice
— `slices.Equal(RefusalKinds(), want)` is writable.

Three of `0029:C4`'s `frozen` entries have no such accessor:

- **the exit-code classes** — C4 names `internal/cli/clierr::ExitCodeFor`.
  That is `func ExitCodeFor(err error) int` (`internal/cli/clierr/clierr.go:131`),
  a `switch` over `ce.Group` returning ints with a `return 1` default. It
  enumerates nothing. There is no `ExitCodes()` and no `Groups()` accessor
  in `clierr` (the package's only exported funcs are `ErrorCode`,
  `ExitCodeFor`, `EmitJSON`, `WriteJSONLine`, `EmitText`,
  `EmitFindingsText`). **Test prevented:** "the frozen exit-code class set
  is exactly {0,1,2,3,130}" — a tester must either reverse the set out of
  the `switch` by exhaustive `errors.As` fixtures (which asserts behaviour,
  not membership, and silently passes when a sixth Group is added with a
  duplicate code) or add a production accessor this RDR does not
  authorize. C4 does not say which, and no scenario names the shape.
- **the stderr advisory `level` set** — C4 names `note`, `warning`. In
  source these are bare string literals inside `map[string]string{"level":
  "note", …}` at `internal/cli/respond/respond.go:207` and `:217`. No
  `Levels()` accessor exists. **Test prevented:** the by-value frozen
  assertion; the only writable test is "invoke `Note`, observe the string",
  which cannot detect a *third* level being added by a future `Alert`
  emitter — precisely the frozen violation the test exists to catch.
- **the `graph-lint-failed` aggregate code** — `AggregateCode` is a single
  `const` (`internal/graphlint/taxonomy.go:48`), not a set. A "set matches
  its declared members exactly" assertion over a one-member non-set is
  either trivial (`AggregateCode == "graph-lint-failed"`, which asserts a
  spelling, not a closure) or unwritable as specified. C4 calls it "the one
  code that envelope ever carries" — that *closure* claim is the testable
  one and nothing enumerates the emit sites to check it.

The two that ARE writable: `findings[].operator`
(`internal/guard/grammar.go:57 Operators()`, and C4 correctly cites the
existing `slices.Equal` test) and the gate `verdict` set
(`internal/accessor/model.go:309 Verdicts()`). So S6 is writable for two of
five and underspecified for three. **Blocks:** Phase 1 Step 3 ("Assert the
tier rules the consumer is promised") cannot be scoped or estimated — the
implementer does not know whether it is three test files or three test
files plus three new production accessors.

---

## F2 (BLOCKING) — `0029:C4` + `0029:S7`: the `append-only` CLIError `code` vocabulary has no enumeration anywhere, so S7's membership-and-uniqueness assertion has no subject

`0029:C4` assigns `append-only` to "the CLIError `code` vocabulary".
`0029:S7` says append-only tiers get "membership-and-uniqueness assertions
only", pointing at
`internal/table/command_carrier_0025_test.go:924 TestReq79_TheCategorySetIsAppendOnlyAndDuplicateFree`
as the pattern — which works because `internal/table/category.go:90
Categories()` returns the set.

There is no analogous accessor for CLIError codes. `grep` over
`internal/cli/clierr/` finds no `Codes()`, no `codes = []string{…}`, no
code-constant block; codes are minted as string literals at the raise
sites across `internal/cli/`. "Membership" and "uniqueness" are both
undefined predicates over a set that is not materialized. **Test
prevented:** `TestTheCLIErrorCodeVocabularyIsAppendOnlyAndDuplicateFree` —
there is nothing to pass to `slices.Contains` or to dedupe. Worse, the
append-only *promise* (no member removed or renamed within a major) has no
mechanical guard at all: a rename is a one-literal edit that no test
observes. **Blocks:** the C2 `append-only` tier is unenforceable for the
single largest vocabulary the RDR tiers, and the implementer cannot tell
whether C4 is asking them to build that enumeration.

Same defect, lesser scale, for two more `append-only` C4 entries: the
`graph-unprovable-coverage` `reason` set and the `flow next` unknown-`reason`
set. The latter is composed at `internal/cli/flow_next.go:92-93` from
`resolve.Reason*` constants plus a locally-declared `reasonNotEvaluated`,
with no accessor over the union. `0029:A4` correctly identifies it as "a
set DISTINCT from the `graph-unprovable-coverage` `reason` set… on a
different field with a different producer" — but a distinct set that is
never materialized is a distinct set no test can range over.

---

## F3 (BLOCKING) — `0029:S9`: the `closed`-retirement grep is scoped to two files while `closed` is live in at least four more, and the scope is the pass/fail criterion

`0029:S9` specifies the test exactly: "no occurrence of `closed` describing
a tier in `internal/graphlint/taxonomy.go` or `internal/table/category.go`.
A grep assertion is sufficient". That is a precise, writable test. It is
also, as scoped, a test that passes while the ambiguity `0029:C2` exists to
remove survives in production comments on vocabularies **this same RDR
tiers**:

- `internal/guard/grammar.go:47,56,59,60,67,145` — six occurrences, all
  describing `Operators()`, which `0029:C4` assigns `frozen`. C4 itself
  quotes `0003:C7`'s "MUST be closed and typed" approvingly.
- `internal/accessor/model.go:308` — "the closed three-member gate verdict
  vocabulary", which C4 assigns `frozen`.
- `internal/resolve/resolve.go:45,66` — "the closed kernel-owned refusal
  kind set", i.e. `data.escape_class`, which C4 assigns `frozen`.
- `internal/cli/flow_input.go:349` — "RDR 0005's refusal table is closed",
  describing the CLIError code vocabulary, which C4 assigns **`append-only`**
  — the one place where a surviving `closed` comment states the *opposite*
  of the tier C4 now assigns.

`0029:C2`'s clause is unqualified: "The term `closed` MUST NOT be used to
describe any of these tiers, in code comments or documentation". S9's
two-file scope does not test C2's clause; it tests a strict subset of it.
**Test prevented:** the real C2-last-clause test. As written the scenario is
green-on-day-one for the flow_input.go:349 comment that actively
contradicts C4. A tester writing S9 verbatim ships a false pass.
Note also that C4 and A1 *themselves* retain the word (C4's `frozen` entry
quotes "MUST be closed and typed"; `0006:C17`'s amendment target is named
`…IsClosedAtExactlyFourMembers`), so the RDR does not say whether a
quotation of a peer contract counts as an occurrence — the tester cannot
decide whether `internal/graphlint/findings_0006_test.go:287`'s renamed
successor may keep the word.

---

## F4 (MAJOR) — `0029:S1` asserts an exact top-level key set that `respond.Success` cannot produce

`0029:S1` **Expected**: "emits top-level keys `data,schema_version,type`…
The pre-change baseline is `data,type` (captured), so this row is the one
added key and nothing else."

`internal/cli/respond/respond.go` `Success` has four marshalled fields, not
two: `Type` (non-omitempty), `Notes []Advisory` (`omitempty`), `Warnings
[]Warning` (`omitempty`), `Data` (`omitempty`). The captured baseline
`data,type` is therefore the key set *of one particular clean-model run
that emitted no notes and no warnings* — not the envelope's shape. S1 states
it as the shape. **Test prevented / mis-specified:** a tester writing S1 as
an exact `sort(keys) == ["data","schema_version","type"]` assertion has
written a test that fails the moment `lint` emits a `Note` — for reasons
that have nothing to do with this RDR. The scenario gives no pass/fail rule
for the omitempty siblings: is the criterion "exactly these three keys", or
"these three keys are present and the only *added* key is `schema_version`"?
Those are different tests with different failure sets, and `0029:MVV` step 1
inherits the same ambiguity ("record the full envelope… It reports
`"schema_version":"0.1"` and no findings"). `0029:S2` handles the refusal
record's optional fields carefully — "`code,detail,message,param,
schema_version` (or with `findings` where the failure aggregates)" — which
shows the RDR knows the distinction and did not apply it to S1. (S2 has its
own smaller version of the gap: `CLIError` also carries `Hint`
(`clierr.go:60`, omitempty), which neither S2's key list nor its
parenthetical mentions.)

---

## F5 (MAJOR) — `0029:C2` requires every tier be "recorded in `docs/cli-output-contract.md` beside that vocabulary", and no scenario tests it; the doc today contains none of the three tier words

`0029:C2`'s first sentence makes the doc the normative register of record:
"the tier is recorded in `docs/cli-output-contract.md` beside that
vocabulary". `0029:C4`'s last paragraph makes it an invariant: "A surface
added later takes a tier assignment in the same document as part of the
change that adds it; an unassigned machine-readable surface is a defect."

The file exists (`/Users/cwensel/sandbox/newcoinc/intrastate/docs/cli-output-contract.md`)
and contains **zero** occurrences of `frozen`, `append-only`, or `growing`.
No scenario S1–S9 covers the doc. `0029:MVV` does not touch it.
Phase 2 Activation Step 1 ("Publish the promise where agents already read")
is the only place it appears, and Phase 2 is by construction not Phase 1.
**Test prevented:** the C4-last-paragraph defect check — "every tiered
vocabulary named in C4 appears in `docs/cli-output-contract.md` with a tier
word beside it, and no machine-readable surface appears there without one".
That is the single test that keeps `0029:A4`'s finding (C4 missed seven of
fourteen surfaces on the first pass, and two more inside `findings[]` on the
grounding pass) from recurring silently. The RDR diagnoses the recurrence in
A4's "If wrong" — "That the gap recurred once inside the fix is why C4 now
states where a vocabulary must be looked for" — and then leaves the
recurrence undetectable, because "where to look" is prose, not an assertion.

---

## F6 (MAJOR) — `0029:MVV` step 4 requires observing `"0.1"` → `"0.2"`, but nothing in the record says who moves the literal or what makes a failure to move it detectable

`0029:MVV` step 4 **Expected end state** includes "C1 moves the minor
(`"0.1"` → `"0.2"`) and never the major", and the `oracle` table's "4 version
movement" row makes this a by-value read: "major moved, or minor did not —
read off the emitted string, not inferred from the change class". Its
negative control is "a no-op release MUST move neither".

`0029:C1` says the minor "increments for backward-compatible additions",
but neither C1 nor `0029:§phase-1-code-implementation` states where the
literal lives or what increments it. If it is a hand-edited constant, the
oracle's stated failure mode — a contributor adds a `growing`-tier member and
forgets the bump — is caught by **no** automated test, only by the human
running the MVV by hand once. The `disposition` table's row "model firing an
`info` code" does not mention the version at all. **Test prevented:**
`TestAddingAnAdvisoryCodeMovesTheSchemaMinor` — there is no mechanical link
between "a vocabulary gained a member" and "the literal changed", so the
assertion has no left-hand side. The `oracle`'s negative control ("a no-op
release MUST move neither") is likewise unwritable: nothing defines what a
release is, in-repo, for a test to range over.

Related, smaller: `0029:C1` says the consumer "MUST reject an envelope
reporting an unsupported major", and the `disposition` table marks both
consumer-side rows "n/a — consumer side". The repo *has* an in-repo consumer
— `internal/cli/flow_input.go:338 planEnvelope`, which `0029:S3` and
`0029:A3` both name. No scenario asks whether `planEnvelope` implements
C1's two consumer MUSTs. **Test prevented:** "the repo's own envelope
consumer rejects an unsupported major" — and the answer today is that it
does not, since `planEnvelope` has no `schema_version` field at all. The RDR
does not say whether that is in scope (making it a Phase 1 obligation) or
out (making C1's consumer rules advisory-to-third-parties only).

---

## F7 (MINOR) — `0029:S5` and `0029:MVV` step 5 are the same assertion, and neither says how the promotion is staged in a test

`0029:S5` ("Promotion to `blocking` IS the verdict-changing event",
expected exit 2 and the record shape flips to bare `CLIError`) restates
`0029:MVV` step 5 without adding a mechanism. Both require a code to exist
at `info` and then at `blocking` within one test binary.
`internal/graphlint/taxonomy.go:106,109` return `slices.Clone` of unexported
package vars, and `severityFor` (`:123`) derives severity from
`slices.Contains(blockingCodes, code)` — so promotion means mutating a
package-level var. Neither scenario says whether the test does that (needs a
seam that does not exist), or whether S5 is a manual MVV step only. **Test
prevented:** an automated S5. `0029:S8` has the same shape but is explicit
about being a source edit ("is updated in the same change"), which is why
S8 is actionable and S5 is not.

---

## F8 (MINOR) — `0029:S3`'s second half states a diff criterion with no named artifact path for the *predecessor*

`0029:S3`'s golden clause: "that golden is re-captured in the same commit…
and the re-captured file differs from its predecessor by exactly the one
`schema_version` key". The path is real
(`docs/rdr/0023-resolve-envelope-projection/artifacts/mvv-step1-default-golden.json`,
read by `internal/cli/flow_mvv_0023_test.go:44`). But "differs from its
predecessor by exactly one key" is a *review* criterion over a git diff, not
a test — the predecessor ceases to exist in the tree once re-captured.
**Test prevented:** none, strictly; but the criterion is stated in the
Testing Strategy as though it were a test row, so an implementer may believe
it is covered by CI when it is covered by a human reading a diff. Say which.

---

## Not findings (checked, clean)

- Every test symbol the RDR names exists at the path it names:
  `TestReq74_TheAdvisoryTierIsClosedAtExactlyFourMembers`
  (`internal/graphlint/findings_0006_test.go:287`),
  `TestReq79_NoSixthRefusalKindIsMinted`
  (`internal/resolve/guard_atoms_test.go:894`),
  `TestReq79_TheCategorySetIsAppendOnlyAndDuplicateFree`
  (`internal/table/command_carrier_0025_test.go:924`), `planEnvelope`
  (`internal/cli/flow_input.go:338`), the 0023 golden. `0029:A1`'s claim
  that the four-element `want` literal is the only blocker on a fifth
  advisory code holds — `severityFor` derives, nothing counts.
- `0029:S4` is writable as specified; `severityFor` deriving from
  `!IsBlocking` is exactly what the scenario says it is.
- `0029:S7`'s refusal to add a byte-golden over the envelope is correct and
  well-argued, and `0029:§performance-expectations` correctly scopes the
  byte-stability non-claim. No finding.
- `0029:C4`'s two deliberately tier-less namespaces (`findings[].class`,
  `data.dispositions`) are consistently handled; `0029:A4` backs them.
