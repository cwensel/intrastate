Model: claude-opus-5[1m]

# Persona 3 — QA / Tester: top 5 tests I cannot write

Lens: how do I test this, and what are the pass/fail criteria? Each finding
names the exact test I would try to write and the specific missing criterion
(expected value, oracle, fixture shape, observable) that stops me. Grounded
against the shipped tree: `internal/resolve/resolve.go`,
`fixtures_test.go`, `adversarial_test.go`, `fixup_test.go`,
`resolve_test.go`, `mvv_test.go`.

Severity scale: **Blocker** = the test cannot be written at all;
**Major** = writable only by inventing the contract myself;
**Minor** = writable but its pass criterion is under-determined.

---

## QA-1 — Scenarios 4, 6, 7, 8 have no fixture shape: `Row.Guard` is an opaque string with no atom representation

**Severity**: Blocker

**Passage** — *Validation / Testing Strategy*, Scenario 4:

> "**Scenario**: strong-Kleene combination matrix — each operator × {tag
> present, tag absent} × verdict pairs across `all`, and the same across
> `unless` (recalling `unless` is ¬(conjunction), not per-atom negation).
> **Expected**: the A2 tables. Specifically `F ∧ U = F` (prunes), `T ∧ U = U`,
> `U ∧ U = U`, and `unless` yielding FALSE only when every `unless` atom is
> TRUE."

And *Existing Infrastructure Audit*:

> "| View-reading evaluator for the domain-rule scenarios | None — no
> `GuardEvaluator` implementation reads the `TagSet` | … | Build |"

**What's missing**: the RDR correctly identifies that no *evaluator* exists —
but the deeper blocker is that no **input representation** exists for an atom,
an operator, an `all` block, or an `unless` block anywhere in the kernel. The
seam is `Evaluate(guard string, view TagSet) GuardResult`
(`internal/resolve/resolve.go:93`), and `Row.Guard` is a bare `string`
(`resolve.go:185`) that the kernel explicitly never parses. RDR 0003's atom
shape lives only in a TOML spike fixture
(`docs/rdr/0003-guard-predicate-exhaustiveness/evidence/spikes/guard-fixture.toml`),
in a format with no Go type, no parser, and no loader in this repo (`find . -name
testdata -type d` returns nothing). So to write Scenario 4 I must first invent:
(a) the wire/Go shape of an atom (`{tag, op, literal}`? nested TOML tables?),
(b) how `all` and `unless` blocks are carried into the single `guard string`,
(c) whether the harness parses TOML or constructs Go structs, and (d) which
package owns those types — this RDR's test tree, RDR 0003's future package, or
`internal/resolve`. Each choice yields a different, mutually incompatible
vector file, and RDR 0003's implement stage must consume the *same* one for the
"pass verbatim" obligation (Phase 2) to mean anything. The RDR calls this
"Build" and stops there.

**Test prevented**: `TestKleeneMatrix_ValueOperatorsOverAbsentAndPresentTags` —
the table-driven `{operator, tag present?, atom verdict} → guard verdict` matrix
for Scenario 4, and by the same blocker Scenario 6 (`exists` total), Scenario 7
(`contains` over an absent set-valued tag), and Scenario 8 (two-row absence
pattern). Four of nine scenarios share one unbuildable fixture.

---

## QA-2 — Scenario 1's stated observable is unreachable: `Escaped:false` cannot be asserted on a refusal

**Severity**: Major

**Passage** — *Minimum Viable Validation*:

> "yields `guard_unevaluable` with `Escaped:false` — while the same table with
> the tag present and decided FALSE still prunes and escapes (D8 preserved)."

repeated verbatim in *Validation / Testing Strategy*, Scenario 1:

> "**Expected**: `guard_unevaluable` with `Escaped:false` — never `no_match`,
> never a plan."

**What's missing**: `Escaped` is a field on `Plan`
(`internal/resolve/resolve.go:250`), not on `Refusal` (`resolve.go:255-273`),
and `Resolve` returns exactly one of the two — `Result.Plan` is nil on every
refusal, enforced by `mvv_test.go:102` ("refusal %q also carries a transition
plan"). So `got.Plan.Escaped == false` is a nil dereference, not an assertion.
The intent is presumably "no plan at all", but that is a different and weaker
claim than the one written: it does not distinguish "the escape row was
correctly never reached" from "the escape row was reached, gated, and blocked
for an unrelated reason." Scenario 1's *actual* discriminating criterion — the
one that would catch a regression — is which row's guard text lands in
`Refusal.Guard`, and the RDR never states it for this scenario. Scenario 5(b)
does state it ("carrying the *candidate's* guard"), which makes the omission in
Scenario 1 look accidental rather than deliberate.

**Test prevented**: `TestMVV0007_UnevaluableNotNoMatch` — I can assert
`Kind == KindGuardUnevaluable`, but I cannot write the pass criterion the RDR
actually names, and I have no stated expected value for `Refusal.Guard`,
`Refusal.Rows`, or `Refusal.MissingOwned` on the masking-probe input.

---

## QA-3 — Scenario 2's cited oracle does not test what the scenario describes; the D8-preservation control has no fixture

**Severity**: Major

**Passage** — *Validation / Testing Strategy*, Scenario 2:

> "**Scenario**: D8 preserved — the same table with the guard's tag *present*
> and the predicate decided FALSE. **Expected**: the row prunes and the modeled
> escape rescues, exactly as
> `adversarial_test.go::TestAdv1b_AllGuardsFalseIsNoMatchNotOwnedStateUnavailable`
> already freezes."

**What's missing**: the cited test does not exercise the property. Reading
`adversarial_test.go:114-152`, `TestAdv1b` uses
`fixtureGuards{decided: map[string]bool{"never": false}}` — a verdict handed in
by *guard text*, with the `TagSet` discarded at `fixtures_test.go:20`. There is
no tag present, no tag absent, and no view read anywhere in it; its owned
snapshot is a fixed `status:Draft`. It freezes "a GuardFalse verdict prunes and
lets an escape rescue," which is the kernel half. Scenario 2 claims something
strictly stronger: that *the same table* differing only in whether the guard's
referenced tag is in the view flips between Scenario 1's refusal and this
rescue. That A/B pair — one input with the tag present, one with it absent,
same table, same revision — is the whole point of the control ("Guards against a
fix to scenario 1 that over-refuses"), and it requires a view-reading evaluator
(QA-1) plus a stated table. "Exactly as TestAdv1b already freezes" reads as if
the coverage exists; it does not. Note the RDR is scrupulous about this
elsewhere — Scenario 3 says "**No shipped test covers this**" — which makes
Scenario 2's claim the outlier.

**Test prevented**: `TestD8Preserved_PresentTagDecidedFalseStillPrunesAndEscapes`
— the paired A/B control for the masking probe. I cannot write it against the
cited oracle (wrong mechanism), and I have no fixture for the tag-present arm.

---

## QA-4 — Scenario 5's two halves assert the same observable, so the test cannot distinguish pass from a specific failure

**Severity**: Major

**Passage** — *Validation / Testing Strategy*, Scenario 5:

> "**Expected**: (a) `guard_unevaluable` carrying the escape row's guard … (b)
> `guard_unevaluable` carrying the *candidate's* guard — the escape path is
> never reached."

read with *Normative Contracts* (aggregation block):

> "An unevaluable CANDIDATE row MUST NOT be masked by a decidable escape row.
> An unevaluable ESCAPE row yields guard_unevaluable in place of the
> candidate-set refusal"

**What's missing**: both halves produce `Kind == KindGuardUnevaluable`, so the
only discriminator is `Refusal.Guard`'s string — and the RDR gives no expected
literal for either. That matters because the discrimination is fragile in a way
the RDR does not address: `gate` picks the guard text of the row lowest by
`compareRefs` over `{RuleID, SourceLocator}` (`resolve.go:409-416`), so the
expected value depends on the fixture's rule IDs, not on which row set was
gated. If a vector author names the escape row `rdr.escape…` and the candidate
`rdr.cand…`, and both were somehow gated together, the assertion would still
pass for the wrong reason. The stronger available observable, `Refusal.Rows`
(which names *which* rows are implicated and would separate the two halves
unambiguously), is never mentioned. Compounding this: 5(b) asserts "the escape
path is never reached" — pure negative-space with no observable at all, since
`escapeOrRefuse` is unexported and leaves no trace in `Result` when it doesn't
run.

**Test prevented**: `TestEscapeSetScoping_UnevaluableEscapeVsUnevaluableCandidate`
— specifically the 5(b) leg. I cannot write an assertion that fails if a future
refactor lets a decidable escape row mask an unevaluable candidate, because the
only stated expected value is a kind both legs share plus an unspecified string.

---

## QA-5 — Scenario 9 depends on assumption A9, which is `Pending`, so both halves have contradictory candidate oracles

**Severity**: Minor

**Passage** — *Validation / Testing Strategy*, Scenario 9:

> "(b) a row whose `unless` block is omitted. **Expected**: … (b) the row is
> *not* disabled — an omitted `unless` block must not read as a vacuously-true
> conjunction. Guards the whole-table failure A9 names."

against *Prerequisites*:

> "- [ ] **A9 verified** (empty-`unless` identity against RDR 0002's
> normalization) … Must close before lock"

and A9 itself: "**Status**: Pending … **Evidence**: Needed."

**What's missing**: A9's own "If wrong" branch is *"the identity belongs to RDR
0002's normalization (not here)"* — i.e. the unresolved question is not merely
what the answer is but **which RDR's test tree owns the assertion**. Until that
closes I cannot tell whether Scenario 9(b) is a vector in this RDR's conformance
suite (asserted against a `GuardEvaluator`), a normalization test in RDR 0002's
tree (asserted against a parser this repo does not have), or neither. Nor can I
state the fixture: "a row whose `unless` block is omitted" is a claim about the
authored TOML, but the kernel sees only `Guard string` — so "omitted" versus
"present but empty" is a distinction with no representation at this seam at all,
which is precisely the distinction the scenario turns on. Scenario 9(a) by
contrast is cleanly writable today (`Guard: ""` → `GuardTrue`, pinned by
`resolve.go:426-428`), so the two halves of one numbered scenario have very
different readiness and the RDR marks them alike.

**Test prevented**: `TestEmptyUnlessBlockDoesNotDisableRow` — I cannot pick a
package to put it in, cannot construct an input that distinguishes omitted from
empty, and have no settled expected value until A9 is verified.

---

## Coverage note (not a numbered finding)

The RDR is honest about two gaps and both are real, verified against source:
the "Owned-state ordering gap" (no shipped test puts a missing-owned row beside
an unevaluable-guard row in one survivor set — confirmed by reading `gate`'s
sequencing at `resolve.go:388` vs `408` and finding no contending test), and
Scenario 3's "**No shipped test covers this**" (`TestAdv3` at
`adversarial_test.go:266` does pair two rows, but both use unknown predicates
against `fixtureGuards{}` and it asserts only payload order-stability). Both
statements check out; they are correctly disclosed rather than findings against
the RDR.
