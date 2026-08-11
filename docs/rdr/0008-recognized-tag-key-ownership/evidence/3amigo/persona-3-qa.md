Model: claude-opus-5[1m]
Persona: QA / Tester

# 3amigo — Persona 3: QA / Tester (RDR 0008)

Grounding: `internal/resolve/resolve_test.go` (44 `TestReqNN_*` functions,
kernel-only), `internal/resolve/fixtures_test.go` (fixture builders; the only
guard seam is `fixtureGuards.Evaluate(guard string, _ resolve.TagSet)`, which
discards the view), `internal/resolve/resolve.go` (`recognizedTagKey`,
`assemble`, `Resolve (Result, error)`, `TagSet.Lookup/Len`, `Refusal`).
RDR 0002's normalizer does not exist at HEAD, so scenarios 2, 4, 5, 7, and the
lint half of 9 have no harness at all — the findings below are about criteria
that would still be missing *after* that harness exists.

---

## QA-1 — Scenario 6 has two mutually exclusive Expected clauses and no tiebreaker

**Test it prevents**: `TestRdr0008_ProducerObligation_ObservedTagKeyedRecognized`
— construct `resolve.Input` directly with
`Observed: []resolve.Tag{{Key: "recognized", Value: "failed"}}` and
`Recognized: "successful"`, call `Resolve`, and assert the breach is detected.

**Passage**: Validation / Testing Strategy, scenario 6 —

> **Expected**: the breach is caught at whichever locus A6's
> Enforcement-locus decision settles, on the Go error path
> `internal/resolve/resolve.go::Resolve` reserves for programmer mistakes —
> never a new `RefusalKind`. If the locus resolves to (a) documented-only,
> this scenario instead pins the D3 backstop: the collision resolves
> deterministically to the owned/observed value and surfaces as a typed
> refusal, never a wrong edge.

**Missing criterion**: the *oracle itself is the open decision*. The two
branches are opposite assertions on the same call: (b)/(c) require
`err != nil && result.Refusal == nil`; (a) requires `err == nil` and a
non-nil `Refusal`. Load-Bearing Decisions / Enforcement locus says "Pick the
final form at Pre-Lock", and A6's Residual concedes "the design choice it
gates is not [resolved]". A test cannot branch on an undecided design; I
cannot write scenario 6 in either direction, and there is no stated default
that a test may encode provisionally. Compounding it, branch (a)'s oracle is
also under-specified: it names no `RefusalKind`. Against `recognizedTagSensitiveTable`
the D3 overwrite yields `no_match` (the row's `{Key: "recognized", Value:
"successful"}` pattern fails once the observed value shadows it), but against
a table whose row omits that pattern the same breach yields a *plan*, not a
refusal — so "surfaces as a typed refusal" is not universally true and the
fixture that makes it true is unnamed. Also unspecified for (b)/(c): whether
the precondition rejects an owned/observed `recognized` key when
`Input.Recognized` is empty (the reservation "holds unconditionally" per
block 1, but `assemble` injects nothing in that case, so the breach is
harmless — is it still an error?).

**Severity**: **High**

---

## QA-2 — Scenario 3's "identical assembled view" has no observable and no named seam to assert through

**Test it prevents**: `TestRdr0008_GuardAndMatcherReadOneView` — a guard seam
that inspects the `TagSet` it is handed and asserts it carries the same
`recognized` binding the matcher used in the same resolve.

**Passage**: Validation / Testing Strategy, scenario 3 —

> a guard predicate and a match pattern both read the recognized tag in one
> resolve, through a guard seam that actually inspects the `TagSet` it is
> handed. **Expected**: both observe the identical assembled view.

**Missing criterion**: "identical" is not a checkable predicate against any
exported surface. `TagSet` has unexported fields (`tags map[string]taggedValue`)
and exports only `Lookup(key) (value, prov, ok)` and `Len() int` — no equality
method, no iteration, no `Keys()`. `reflect.DeepEqual` on two `TagSet` values
from an external `resolve_test` package compares unexported maps, and there is
only *one* view to compare against anyway: the matcher's view is never handed
to the test. So the RDR must say what "identical" reduces to observably. The
candidate predicates differ in strength and the RDR picks none: (i) the guard
sees `Lookup("recognized") == (in.Recognized, ProvenanceRecognized, true)`
— pins only the one key, not the view; (ii) guard-side `Len()` equals the
matcher-implied key count — pins cardinality but no values; (iii) the guard
records the view and the test cross-checks every key the row's `Match` names.
Which one counts as "would fail if the two views diverged" (Done means) is
undecided. The RDR also does not say *who owns the fixture change*: A2's
"Test gap" correctly identifies `fixtureGuards.Evaluate(guard string, _
resolve.TagSet)` as discarding the view, but neither scenario 3 nor Phase 3
states whether the existing shared `fixtureGuards` is modified (which would
touch the ~10 tests already using `allGuardsTrue`, e.g.
`TestReq23_GuardEvaluationIsDelegatedToTheInjectedSeam`) or a new
view-capturing seam type is added alongside it. Without that, the test's blast
radius on the existing suite is unknown.

**Severity**: **High**

---

## QA-3 — Scenario 5's "at most one such declaration" is unfalsifiable as stated, and its scope unit is undefined

**Test it prevents**: `TestRdr0008_DuplicateRecognizedDeclarationRejected` — a
TOML fixture carrying two recognized-provenance tag declarations, asserted to
fail lint in the `reserved tag key` category.

**Passage**: Normative Contracts, block 2 —

> A tag declaration with provenance `recognized` MUST be named `recognized`,
> and a flow MUST carry at most one such declaration.

and Testing Strategy scenario 5 —

> a flow carrying two recognized-provenance declarations. **Expected**:
> `reserved tag key` failure

**Missing criterion**: the fixture cannot be built without contradicting the
first clause of the same sentence. RDR 0002's carrier is
`[tags.<tag>]` (Load-Bearing Decisions / Wire format), a TOML table keyed by
tag name, so two recognized-provenance declarations necessarily have two
*different* names — meaning at least one already violates the MUST-be-named-
`recognized` clause and fails on that rule, not on cardinality. Two
declarations both named `recognized` are a TOML duplicate-key error
(`malformed TOML` category), not a `reserved tag key` failure. So there is no
input that isolates the cardinality rule, and the test would silently pass for
the wrong reason. Second gap: the scope unit is "a flow", but RDR 0002's
schema unit is `[model] id`, and "flow" appears in 0002 only as prose about
authoring (`:20`, `:26`, `:194` "Each flow has a named model") — the RDR never
says whether one TOML file is one flow, so a two-file fixture has no defined
verdict either. Third: the RDR does not say which category wins when an input
breaches two rules at once (wrong-named *and* duplicated), so the assertion
target is ambiguous even if a fixture existed.

**Severity**: **High**

---

## QA-4 — Scenario 7's golden-text assertion names no data surface, no field names, and no consumer to render through

**Test it prevents**: `TestRdr0008_ReservedKeyFailureCarriesGuidance` — assert a
`reserved tag key` validation failure value carries the offending name, the
required name `recognized`, and the rule, and that a category-blind consumer
cannot strip them.

**Passage**: Normative Contracts, block 3 —

> Every `reserved tag key` failure … MUST carry, at the data level, the
> offending name, the required name `recognized`, and the rule (the kernel
> owns this key; declarations conform to it). A category consumer may map it,
> but the guidance travels in the failure data, not the renderer.

and Testing Strategy scenario 7's "Golden-text assertion at the data level,
not the renderer."

**Missing criterion**: three assertion inputs are absent. (1) **No named
surface** — there is no type, field set, or symbol for a data-level validation
failure anywhere in the RDR or at HEAD; `internal/resolve`'s `Refusal` is
explicitly the wrong channel ("never a kernel refusal", block 2), and RDR
0002's contract says only that failures "retain stable data-level categories",
naming no carrier struct. A golden test needs a field to read. (2) **"the
rule" is not a value** — the parenthetical gives prose ("the kernel owns this
key; declarations conform to it") but the RDR does not say whether the golden
oracle is that exact string (byte-comparable, and then it is normative text
this RDR must pin verbatim), a stable rule *identifier*, or free-form
implementer prose. A golden-text test with a free-form oracle is untestable.
(3) **The negative half has no harness** — "rendered by a category-enumerating
consumer that does not know the new category" names no such consumer;
RDR 0005's exit-code map is cited in Risks but not as a test target, and
0005's mapping surface is not implemented at this seam today. Without a named
consumer, "a generic renderer cannot strip the guidance" is a claim about
absent code, not a predicate.

**Severity**: **High**

---

## QA-5 — Scenario 8's "proves the collision branch unreachable" is a restatement, not a checkable predicate

**Test it prevents**: `TestRdr0008_D3BackstopStillDeterministic` — assert that
no conforming table can reach `assemble`'s owned/observed overwrite of the
`recognized` entry, while non-conforming hand-built input still resolves
deterministically under D3.

**Passage**: Validation / Testing Strategy, scenario 8 —

> a conforming normalized table proves the `assemble` collision branch
> unreachable; a hand-constructed non-conforming table proves D3 still
> resolves it deterministically. **Expected**: no conforming input reaches
> the overwrite; non-conforming input keeps RDR 0001 D3's precedence
> (`owned` > `observed` > `recognized`) unchanged.

**Missing criterion**: the first half asserts a *universal negative over an
unbounded input domain* with no stated method. A single conforming fixture
demonstrates one non-collision, not unreachability; A5's proof is a derivation
over RDR 0002 semantics (P1+P2+P3), which is a review artifact, not a runnable
oracle. The RDR does not say which discharges the scenario — an exhaustive
generator over conforming declarations, a property/fuzz test with a stated
domain and iteration count, a coverage assertion that the overwrite line is
never taken, or simply "no test; A5's derivation is the evidence". Second, the
overwrite in `internal/resolve/resolve.go::assemble` (the `Observed` then
`Owned` loops writing over `view.tags[recognizedTagKey]`) has **no
observable** — it is an unexported map write with no counter, no error, and no
distinguishing output; from outside the package the only signal is a changed
disposition, which is indirect and fixture-dependent. So "reaches the
overwrite" is not observable at all through the `resolve_test` package
boundary, and the RDR names no instrumentation. Third, the second half's
"deterministically" is unpinned: it does not say the assertion is
`Lookup("recognized")` returning the owned/observed value with the
owned/observed provenance (the direct D3 observable, available via the exported
`TagSet.Lookup` — but `TagSet` is not returned by `Resolve`), versus a
disposition-level assertion, versus a replay-equality assertion in the style of
`TestReq35_RefusalDispositionsReplayValueIdentically`.

**Severity**: **Medium**

---

## Cross-cutting note (not a numbered finding)

"Done means: every scenario below has a green test" (Testing Strategy) is not
satisfiable at lock: scenarios 2, 4, 5, 7 and half of 9 depend on RDR 0002's
normalizer, which the Prerequisites correctly gate ("RDR 0002 implementation
underway") but which the Done criterion does not carve out. As written, the
RDR's own completion predicate cannot go green from this RDR's implementation
alone, so "Done" has no evaluable boundary for the work this RDR authorizes.
