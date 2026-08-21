Model: claude-opus-5[1m]

# Persona 3 — QA / Tester (iteration 2, delta pass)

Scope: the reshape surface — 20-row Testing Strategy matrix, Minimum Viable
Validation, `oracle`/`disposition`/`trace` mini-checks, the strong-Kleene
combination matrix, the refusal payload's asserted shape and sort order
(A19/A3), the existence-token drift test (A16), the nil-seam per-atom
narrowing (row 19).

Findings are severity-ranked. Each names the anchored passage, the test it
prevents, and why the pass/fail criterion is missing or unassertable.

---

## QA-1 — The refusal payload has no named field, so every payload assertion in the matrix and MVV is unwritable — HIGH

**Anchored passage.** Normative Contracts, the refusal clause: "The refusal
MUST name what was missing, per row and per atom: for every undecidable row,
its `(RuleID, SourceLocator)` and, for each of its unevaluable atoms, the
referenced key, the block, and a reason drawn from a closed set". Testing
Strategy row 1 ("key in payload"), row 12 ("payload reason `uncomparable`,
not `absent`"), row 17 ("Payload determinism"), row 20 ("Every undecidable
row appears in the payload"); MVV scenario 1 ("the absent key in the
payload"); `oracle` MVV row 1 ("asserts on refusal KIND + payload CONTENT").

**Test I cannot write.** Any assertion on payload content. The RDR names no
field on `Refusal` and no type for the entry. `internal/resolve/resolve.go::Refusal`
today carries `Kind`, `Revision`, `Flow`, `Recognized`, `Rows []RowRef`,
`MissingOwned []string`, `Guard string` — all flat. The Existing
Infrastructure Audit acknowledges this ("`Refusal` is flat scalars plus
`[]RowRef`/`[]string` … EXTEND `Refusal` with the kernel's first two-level
payload") but does not name the extension. A3's re-verification note names
`undecidedAtoms`/`compareAtoms`/`copyAtoms` as *proposed function* names and
explicitly says the spike's `Refusal.UndecidedAtoms`
(`evidence/spikes/a3-reshape/spike_kleene_test.go.txt`, asserting
`got.Refusal.UndecidedAtoms[0].Block`) is **the wrong shape** — it retained
`Refusal.Guard` beside it, whereas the contract says REPLACES. So the one
concrete field name in evidence is disowned by the RDR that cites it.

**Why the criterion is missing.** Six matrix rows and the MVV assert "the key
is in the payload" without a selector. A test author must invent the field
name, the entry struct, and the reason enum's representation (string?
typed constant? `absent` vs `Absent`?), and any two implementers will invent
differently — which means the frozen tests written in Phase 2 do not pin the
contract this RDR states, they pin whatever Phase 1 happened to name. The
"closed set — `absent` / `uncomparable`" is stated in prose only; there is no
`RefusalKinds()`-style enumerator for it, so the analogue of
`resolve_test.go::TestReq7_RefusalKindSetIsExactlyTheFiveNamedKinds` (the
closure test the kernel already applies to its one other closed set) cannot
be written for the reason set.

---

## QA-2 — Row 17's sort order is under-specified for the atoms within a row, so the determinism test has no expected value — HIGH

**Anchored passage.** Normative Contracts, refusal clause: "MUST sort the
payload by row identity then key, so the payload is a function of the input
tuple (RDR 0001 REQ-1), never of atom or row order." Testing Strategy row 17:
"Payload determinism: sorted by row identity then key … `undecidedAtoms`/
`compareAtoms` (both PROPOSED …)". `trace` step 7: "every undecidable row,
sorted row-then-key".

**Test I cannot write.** The row-17 permutation test in the ADV-3 style —
build the same table with atoms permuted inside a row and assert
`reflect.DeepEqual` on the two payloads (the shape
`adversarial_test.go::TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder`
establishes). "Row identity then key" is a total order only if `(row, key)` is
unique per entry. It is not, under this RDR's own clauses:

- The atom slice is `(key, operator token, literal, block)`, and A5 requires a
  row that carries **two atoms over one key** (`X exists = true` beside
  `X eq v`) — "The conjoined value row — two atoms over one key — is admitted
  by the kernel natively once `Row` carries an atom slice (cardinality, not
  grammar)". Testing Strategy row 9 makes exactly this row a test subject.
- The Failure Modes `unless` idiom does the same in the other block:
  "`legal_hold exists = true` beside `legal_hold eq true` in `unless`".
- The same key can be unevaluable in both `all` and `unless` of one row.

So two payload entries can collide on `(row, key)` and differ only in block,
operator, or reason. The comparator is undefined at that point, and the
permutation test's expected value is whatever the implementation's unstable
tie-break produces. This is precisely the class of defect ADV-3 and ADV-3b
were written to catch, re-introduced one level down.

**Why the criterion is missing.** The clause names two of the four atom fields
as sort inputs and is silent on the tie-break. Note the kernel's shipped
comparators are total by construction —
`resolve.go::compareRefs` falls through `RuleID` then `SourceLocator`,
`resolve.go::missingOwned` sorts a deduplicated `[]string` — so there is no
in-repo precedent to borrow the tie-break from.

---

## QA-3 — Row 15 (escape row raises no `owned_state_unavailable`) is a fixture tautology: nothing under test enforces the premise — HIGH

**Anchored passage.** Testing Strategy row 15: "Escape row raises no
`owned_state_unavailable` of its own (empty `RequiresOwned` by composition) |
A21 | `missingOwned` over an empty slice never enters the loop". A21's own
status: "Verified (composition); the PRODUCER question stays open at §JD-3 …
It binds producers, NOT the kernel boundary — so this composition holds for
rows a conforming producer built, and the kernel enforces no part of it."

**Test I cannot write.** A test that can fail. The kernel is `Resolve`, and
`Resolve` takes `Row` values the test itself constructs. To test row 15 I
build an escape row with `RequiresOwned: nil` and assert no
`owned_state_unavailable` — which asserts only that
`resolve.go::missingOwned`'s `for _, key := range row.RequiresOwned` does not
iterate an empty slice. That is a property of Go's `range`, not of this RDR.
The proposition row 15 states — escape rows *have* empty `RequiresOwned` — is
a producer obligation on RDR 0002 (§JD-3, undelivered per Prerequisites) via
RDR 0009's `Writes`-empty clause, and by A21's own admission the kernel
enforces none of it.

**Why the criterion is missing.** There is no fail condition. The shipped
fixture proves it: `fixtures_test.go::escapeRow` builds escape rows with
`RequiresOwned: []string{"status"}` and non-empty `Writes` — i.e. every escape
row in the frozen suite violates the composition row 15 asserts, and
`fixup_test.go::TestFixup1e_AmbiguousClassEscapeIsGatedLikeAnyOtherCandidate`
("missing owned state must not rescue an ambiguity") *depends* on an escape
row raising `owned_state_unavailable` from a non-empty `RequiresOwned`. Row 15
and that frozen test assert opposite things about the same kernel path; the
only thing separating them is which fixture the author writes. Either row 15
needs a kernel-side rejection to test against (which this RDR declines to
add), or it must be marked as a producer-side test that does not belong in the
`internal/resolve` matrix.

---

## QA-4 — Row 19's nil-seam narrowing is only half-discriminating, and its own stated fixture does not discriminate at all — HIGH

**Anchored passage.** Testing Strategy row 19: "Nil seam is per-atom, not
whole-guard: a row of only `exists`/absent-key atoms DECIDES under a nil seam;
only a value atom over a PRESENT key goes unevaluable for want of a seam".
Normative Contracts, nil-seam clause: "a row whose atoms are all
kernel-decidable resolves under a nil seam exactly as it would under a present
one."

**Test I cannot write.** The discriminating half of row 19 for the absent-key
leg. An `exists` atom under a nil seam decides T or F — that leg discriminates
cleanly against the shipped whole-guard rule. An **absent-key value atom**
under a nil seam yields `GuardUnevaluable` under the new per-atom rule and
`GuardUnevaluable` under the shipped whole-guard rule
(`resolve.go::evaluateGuard`: `if seam == nil { return GuardUnevaluable }`).
Same kind, same `Rows`. The two rules are observationally identical on that
input, so the test asserted by row 19's phrase "a row of only … absent-key
atoms DECIDES under a nil seam" cannot be written: the row does not decide, it
goes unevaluable, and the assertion as worded is false against the RDR's own
combination rule.

The RDR conflates two senses of "decides": the clause says *kernel-decidable*
(the kernel answers without the seam), row 19 says *DECIDES* (a T/F verdict).
Only the `exists` sense supports a test.

**Compounding.** The same passage prescribes the Fixup-1d replacement fixture:
"To keep testing the nil-seam rule the fixture needs a value atom over a
PRESENT key (e.g. `reviews >= 3`)." `reviews` is present in the fixture view
(`fixtures_test.go::legalInput` supplies `Observed: {reviews: "2"}`, and
`noMatchInput()` inherits it), so the key is present — but the RDR never says
whether the payload reason for a nil-seam unevaluable is `absent` or
`uncomparable`. The closed reason set is defined as `absent` ("the key was not
in the view") or `uncomparable` ("the key was present and **the seam** could
not compare the value, A18"). With a nil seam there is no seam to report
`uncomparable`, and the key is present so it is not `absent`. Row 12's
"reason `uncomparable`, not `absent`" therefore has no defined answer on the
nil-seam path, which is the very path Fixup-1d exists to pin.

---

## QA-5 — Row 11's "foreign boolean literal is rejected at load" has no load path in the unit under test, and the existence atom's literal is not required by any assertable rule — HIGH

**Anchored passage.** Testing Strategy row 11: "Foreign existence token fails
CLOSED: an atom carrying a non-exported token is a VALUE atom (unevaluable on
absence, seam on presence); **a foreign boolean literal is rejected at load**".
Normative Contracts, existence clause: "a foreign literal is a producer defect
the normalizer MUST reject at load". `disposition` table: "Foreign boolean
literal | rejected at load by the normalizer | producer defect | n/a | loud,
upstream".

**Test I cannot write.** The literal half of row 11, in `internal/resolve`.
The Testing Strategy's own preamble says "Every scenario is a kernel test in
`internal/resolve` (package `resolve_test`)", but load-time rejection happens
in RDR 0002's normalizer, which does not exist and whose duty the Prerequisites
record as **not yet assigned** ("Two duties are added to RDR 0002's Refinement
Context Direction list, which currently carries neither"). There is no load
path in `internal/resolve` — `resolve.go::Resolve` accepts a `Row` value as
given and validates nothing. So row 11's second clause is a test in a package
that does not exist, owned by an RDR that has not accepted the duty, listed in
a matrix whose preamble says every row is a kernel test.

**The kernel-side gap this exposes.** The RDR gives no rule for what the
*kernel* does with a malformed existence atom. The existence verdict is
`presence == literal` with "exactly two boolean literal forms" — but no clause
says what the kernel answers for an existence atom whose literal is neither
form (empty string, `"1"`, `"TRUE"`). The A3 spike's own probe
(`evidence/spikes/a3-reshape/spike_kleene_test.go.txt`, subtest "presence:
exists atom over an absent key is FALSE, not undecidable") constructs
`resolve.GuardAtom{Key: "guard.missing", Op: resolve.OpExists}` with **no
Literal set** and asserts `no_match` — i.e. the spike read a zero-valued
literal as `false` and decided. Under the RDR's `presence == literal` rule an
empty literal is a foreign literal, which the RDR routes to load rejection the
kernel cannot perform. So the one existing executable probe of the existence
rule contradicts the rule as specified, and there is no pass/fail criterion to
adjudicate between them. This is the fail-closed claim's weakest point: it is
asserted for the token but never for the literal on the kernel side.

---

## QA-6 — Row 5's "full strong-Kleene matrix" names no way to produce a U atom of each provenance, and no fixture seam shape — MEDIUM

**Anchored passage.** Testing Strategy row 5: "Full strong-Kleene matrix
across `all` and `unless`, incl. `¬U = U` on a block-level negation". Phase 2:
"encode the domain-rule matrix as kernel tests (operators × present/absent ×
`all`/`unless` × {T,F,U} combinations …)". Risks: "Phase 2's exhaustive kernel
matrix over {T, F, U}³ across both blocks".

**Test I cannot write.** The exhaustive matrix, as a table-driven test whose
cases are enumerable. Three of the four combinations of (verdict, source) are
producible only through distinct fixture machinery, and the RDR names none of
it:

- **T/F from a present key** requires a value stub. `fixtures_test.go::fixtureGuards`
  is keyed on guard *text* (`Evaluate(guard string, _ resolve.TagSet)`) and
  the Existing Infrastructure Audit says "Replace | Atom-level fixture
  evaluator" — but the replacement's keying is unspecified. Keyed on what:
  the key? the `(key, op, literal)` triple? The A3 spike's replacement keys on
  a synthetic `guard.<name>` tag binding, an invention with no normative
  standing.
- **U from absence** requires the key withheld from the view.
- **U from `uncomparable`** requires the stub to answer `GuardUnevaluable` on
  a *present* value — but per the seam contract "The seam decides true or false
  … and MAY answer unevaluable" (A18, permissive), so a conforming stub need
  not offer that path at all.

Without a specified stub contract, "{T,F,U}³" is not an enumeration a test can
mechanically generate; each cell needs a hand-built fixture whose relationship
to the cell is the author's assertion, not the RDR's.

**Compounding: the mutation obligation has no pass criterion.** "Any future
edit to the combination tables MUST be re-mutation-tested against them" names
no mutation tool, no mutant set, and no threshold (all mutants die? which
mutants?). The one mutant on record is "swapping strong-Kleene FALSE-dominance
for UNEVALUABLE-dominance". A CI gate cannot be written from that.

---

## QA-7 — Row 20 and the `Rows`/payload duplication give the "which row" test two candidate oracles — MEDIUM

**Anchored passage.** Normative Contracts, refusal clause: "`Refusal.Rows`
continues to carry every undecidable row. … Any contract or test
distinguishing WHICH row went unevaluable MUST assert on **the row entries**."
Testing Strategy row 20: "Every undecidable row appears in the payload — no
representative row is chosen".

**Test I cannot write.** An unambiguous row-20 assertion. Two structures now
carry the undecidable row set: `Refusal.Rows []RowRef` (shipped, populated by
`resolve.go::gate`'s `rowRefs(undecidable)`) and the new payload's per-row
entries. "MUST assert on the row entries" does not say which, and the RDR does
not state the invariant between them (must the two sets be equal? must the
payload's row set be a subset? can a row appear in `Rows` with zero payload
entries?). Row 6's "refusal naming A" and MVV scenario 3's "refusal naming A"
have the same ambiguity.

The equality is not obviously true: `Refusal.Rows` is populated for *every*
undecidable row, while the payload carries entries keyed by unevaluable
*atom* — a row could in principle appear in `Rows` and contribute no atom
entry if the combination reached U by a route the per-atom scan does not
attribute (the RDR does not enumerate such a route, but it does not rule one
out either). A consistency test between the two is the obvious regression
guard and cannot be written without the stated invariant.

---

## QA-8 — The `disposition` table assigns `GroupUserEnv` to `owned_state_unavailable`, which A19 says has no code row — MEDIUM

**Anchored passage.** `disposition` mini-check, row "Survivor missing an owned
`RequiresOwned` key | `owned_state_unavailable` | **`GroupUserEnv`** |
`MissingOwned` + `Rows` | loud". Contradicted by A19's own evidence:
"`owned_state_unavailable` has **NO row**, a gap §JD-8 names explicitly ('need
`Code` values')."

**Test I cannot write.** Any exit-group assertion driven off this table. I
confirmed against RDR 0005's stable-code table
(`docs/rdr/0005-skill-integration-cli-contract.md`): `flow-guard-unevaluable`
is present with `GroupUserEnv`; there is no row mapping
`owned_state_unavailable`. So the `disposition` table asserts a mapping that
does not exist in the document that owns exit groups. A test that asserts
"`owned_state_unavailable` exits in `GroupUserEnv`" would encode a mapping
this RDR has no authority to make and RDR 0005 has not made.

**Why it matters for the delta.** The `disposition` mini-check is presented as
the input-class ⇒ outcome oracle. Four of its nine rows name `GroupUserEnv`,
and the exit-group column is unassertable for at least one of them and blocked
on A19's Pending reopening for the payload columns of three more.

---

## QA-9 — Row 8's `contains` case is listed as backed while Phase 3 declares its present-key half unwritable — MEDIUM

**Anchored passage.** Testing Strategy row 8: "`contains` over an ABSENT
set-valued tag ⇒ unevaluable, NOT the empty set | A1, A17 | the clause A1 pins
against RDR 0003's silence". Phase 3: "`resolve.Tag.Value` is a bare `string`,
so a set-valued tag reaches the seam as one opaque string with no stated
element encoding. The `contains` leg of this contract test therefore cannot be
written from any current document, and neither can **Testing Strategy scenario
8's present-key half**."

**Test I cannot write.** The present-key control for row 8. Confirmed in
shipped source: `resolve.go::Tag` is `struct { Key, Value string }` with the
package doc "the kernel never parses their internal structure" — there is no
set representation. So row 8 can assert only the absent half, and its negative
control (present set-valued tag containing / not containing the element) is
blocked on an RDR 0003 declaration that Phase 4 only *requests*.

**Why the criterion is missing.** Per the `oracle` mini-check's own standard —
"each has a sibling that must come out the other way" — a row whose control is
unwritable is a one-sided assertion. Row 8 as written passes trivially: the
kernel returns unevaluable on absence for *every* operator it does not
recognize as `exists`, so the test cannot distinguish "`contains` is correctly
partial" from "`contains` is not implemented". The matrix should mark row 8 as
half-covered and blocked, the way Phase 3 does; the matrix row does not carry
that qualification.

---

## QA-10 — Row 16's observed-tag substitution names no expected verdict — MEDIUM

**Anchored passage.** Testing Strategy row 16: "Guard-path observed-tag
substitution: a caller-supplied observed tag turns `guard_unevaluable` into a
decided verdict — the accepted exposure, pinned so it cannot regress silently |
A11, A13 | no existing test contends the GUARD path (ADV-4/ADV-5 defend only
the owned path)".

**Test I cannot write.** The assertion half. "A decided verdict" is not a
pass/fail criterion — decided TRUE (row selected, plan emitted) and decided
FALSE (row pruned, `no_match` or escape) are opposite dispositions, and which
one occurs depends entirely on the value the test supplies and on the stub's
answer for that value. The two existing guards on this exposure are precise by
comparison: `adversarial_test.go::TestAdv4_ObservedTagsMustNotShadowTheOwnedSnapshot`
and `::TestAdv5_ObservedTagCannotSatisfyAnOwnedStateRequirement` each assert a
specific kind.

The regression this row exists to prevent is "the exposure silently widens" —
but a test asserting only "not `guard_unevaluable`" would also pass if the
kernel started deciding absent keys as FALSE, which is the masking path this
whole RDR closes. So the row as worded does not pin what it claims to pin.
Naming the case (`exists = true` over a caller-supplied key ⇒ selected, vs the
same key withheld ⇒ `guard_unevaluable`) would fix it; the RDR does not.

---

## QA-11 — MVV scenario 2's escape leg cannot be built from a conforming escape row — MEDIUM

**Anchored passage.** Minimum Viable Validation: "the same table with the key
present and the value decided FALSE prunes and escapes (D8 preserved)".
Testing Strategy row 2: "Same table, key PRESENT and value decided FALSE ⇒
row prunes and the escape plans (D8 preserved)". `oracle` MVV row 2: "asserts
a PLAN is produced".

**Test I cannot write.** A row-2 test whose escape row is conformant with this
RDR's own composition. To assert "a PLAN is produced" substantively I need the
escape plan to carry something — the `oracle` standard is that verdict-only
assertions are insufficient ("The spike's mutation result is the standing
evidence that verdict-only assertions are insufficient"). But RDR 0009
requires an escape row's `Writes` to be empty, and A21 derives
`RequiresOwned` empty from it, so a conforming escape row yields a plan with
`Writes: nil` and `RequiresOwned: nil`. The only substantive assertion left is
`Plan.Escaped == true` plus `RuleID`. The RDR does not say which of these the
MVV asserts, and the shipped fixture it would naturally reuse
(`fixtures_test.go::escapeRow`) is non-conforming on both fields — so whichever
the author picks, the MVV either asserts almost nothing or asserts against a
row shape RDR 0009 forbids.

**Why the criterion is missing.** The MVV specifies the *input* precisely (a
satisfied `RequiresOwned`, an absent-key value atom, a modeled `no_match`
escape) and, for scenario 2, specifies the outcome only as "prunes and
escapes". Compare scenario 1, which enumerates four assertions (kind, nil
`Plan`, key in payload, row in `Rows`).

---

## QA-12 — The `trace` mini-check's step 5 witness contradicts the shipped `gate` partition it cites — LOW

**Anchored passage.** `trace` mini-check, step 5: "Owned sweep over SURVIVORS |
unevaluable rows ARE survivors, so their `RequiresOwned` counts;
owned-before-unevaluable | **MVV row 1's `RequiresOwned` is satisfied ⇒ passes
through to step 6**".

**Test I cannot write.** A step-5 regression test that fails when the ordering
inverts. The witness given is a row whose owned state is *satisfied* — which
passes through the owned sweep under **either** ordering, so it witnesses
nothing about precedence. The discriminating input is the both-ways row, which
the RDR correctly relocates to Testing Strategy row 13 ("Owned-before-
unevaluable precedence within ONE survivor set"). So the `trace` table's step-5
witness is non-discriminating by construction, and the mini-check's stated
purpose ("the MVV walked stepwise against every assertion in force") is not met
for the one assertion the same section says was the source of the
contradiction it just fixed ("Step 5 caught it").

This is minor because row 13 covers the real case; it is worth recording
because the `trace` table is presented as the walkthrough that validates the
MVV, and a reader taking step 5 as the test specification would write a
passing test that cannot fail.

---

## QA-13 — "No short-circuit" is stated as a MUST with no observable consequence to assert — LOW

**Anchored passage.** Normative Contracts, refusal clause: "The kernel MUST
evaluate every atom of every survivor — no short-circuit on a decided block".

**Test I cannot write.** A short-circuit detection test. The stated MUST is
about kernel *execution*, and the RDR gives it exactly one observable
consequence — payload completeness for unevaluable atoms — which row 20
already covers. But the clause as written is broader than that consequence:
under `F ∧ U = F` the row is pruned and, per "What the payload does NOT
promise", its absences are *not* reported. So for a pruned row, evaluating
every atom is required by the clause and produces no observable difference
whatsoever. A test can only observe seam call counts, and the seam is a test
stub — asserting on stub invocation counts pins the implementation, not the
contract.

Either the MUST needs an observable (e.g. "the payload names every unevaluable
atom of every *surviving* row", scoping it to survivors and making it row 20),
or it should be marked non-testable guidance rather than a MUST. As written a
reviewer cannot tell whether an implementation that short-circuits a
decided-FALSE block violates the contract, because nothing distinguishes the
two behaviors at the boundary.

---

## QA-14 — [out-of-delta] A3's Pending status leaves the frozen-suite baseline the matrix stands on unproven for the specified shape — LOW

**Anchored passage.** Testing Strategy preamble: "The A3 spike
(`evidence/spikes/a3-reshape-probe.md`) established the baseline this matrix
extends: the frozen 154-test suite passes unchanged under the reshape, so
every row below is ADDED surface, not a re-encoding." Contradicted by A3's own
status: "Pending — the BOUND is verified; the SHAPE the spike proved is not
the shape this RDR specifies. … Its PASS for
`fixup_test.go::TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue` depends
on that retention".

**Test I cannot write.** A Phase 1 exit criterion. The preamble asserts a
154-PASS baseline as established; A3 says the spike that produced it kept
`Refusal.Guard` and that at least one frozen test must be **re-decided**, not
re-encoded, under the specified shape. So the number of frozen tests expected
to pass after Phase 1 is not 154 and is not stated. Confirmed against source:
`fixup_test.go::TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue`'s third
assertion is `if got.Refusal.Guard != escape.Guard`, on a field the contract
deletes; `resolve_test.go::TestReq33_UnavailableOwnedStateAndUnevaluableGuardAreValueRefusals`
also reads it (A3 records both). Phase 1's "re-run it (A3 spike)" therefore has
no pass criterion — the honest one is "154 minus the re-decided set", and the
re-decided set is named for Fixup-1d only, by prose in a different section.
