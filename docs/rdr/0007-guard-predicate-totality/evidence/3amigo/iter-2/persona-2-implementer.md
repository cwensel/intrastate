Model: claude-opus-5[1m]

# Persona 2 — Implementer (iteration 2, delta pass)

Question: if I started coding this Monday, what would I ask in the first hour?
Every finding is a clarification request anchored to a named RDR passage, with
the implementation decision it blocks. Claims about shipped behavior cite
`path::Symbol`.

Scope: weighted to the §D1/§D4 reshape — atom shape, `Evaluate(atom, value)`,
strong-Kleene in `evaluateGuard`, the payload replacing `Refusal.Guard` (A19),
kernel-exported existence constants (A16), key identity (A22), `RequiresOwned`
narrowing (A21), and the A3 spike-shape gap. Out-of-delta findings are marked.

---

## IMP-1 — HIGH — The `unless`-block truth value for a row with `all` atoms but no `unless` atoms is stated twice in incompatible ways

**Anchored passage.** Normative Contracts, empty-block clause: "an empty or
omitted `unless` block is ABSENT and contributes nothing — NOT a
vacuously-true conjunction, which under `all ∧ ¬(unless_conj)` would disable
every row carrying one." And the K3 clause: "The row verdict is
`all_result ∧ ¬(unless_conj)`, where `unless_conj` is the conjunction of the
`unless` block's atoms".

**Clarification I would ask.** The row verdict is given as a *single closed
formula* over `unless_conj`, but the empty-block clause says an absent
`unless` block has no truth value at all ("contributes nothing"). Those are
two different implementations:

- `unless_conj = T` when empty (the conventional empty conjunction) makes
  `¬(unless_conj) = F`, and `all ∧ F = F` — the disable-everything outcome
  the RDR itself names and forbids.
- `unless_conj = F` when empty makes `¬F = T` and the row verdict reduces to
  `all_result` — which is the behavior the RDR clearly wants, but it is NOT
  "the conventional identity" for a conjunction, and the RDR explicitly
  refuses to call the empty `unless` a conjunction at all.
- A third implementation branches structurally: `if len(unlessAtoms) == 0 {
  return allResult }`, never evaluating the formula.

Which do I write? The RDR names the wrong answer and forbids it, but never
states the right one as a rule I can encode — it states an *anti*-rule
("contributes nothing") that is not a truth value, while the only normative
formula I am given requires one. The K3 clause also says "The tables are
normative; the evaluator holds no part of them" — so I cannot resolve this by
convention; I need the identity written down.

**Decision it blocks.** The `evaluateGuard` combination step's base case for
an empty `unless` block — i.e. the verdict of *every* row in both canonical
fixtures, which per A9's own evidence "omit `unless` on most rules". Getting
it wrong is a total resolution failure, not an edge case.

**Note.** A9's Evidence discusses the same identity and Testing Strategy
scenario 10 asserts "omitted `unless` ⇒ absent, contributes nothing" — but a
test asserting "contributes nothing" still needs the implementation to pick a
representation, so the test does not disambiguate either.

---

## IMP-2 — HIGH — The payload's Go shape is never specified, and the RDR names symbols that do not exist and are not defined

**Anchored passage.** Normative Contracts, refusal-payload clause: "for every
undecidable row, its `(RuleID, SourceLocator)` and, for each of its
unevaluable atoms, the referenced key, the block, and a reason drawn from a
closed set — `absent` … or `uncomparable`". Plus Testing Strategy row 17:
"`undecidedAtoms`/`compareAtoms` (both PROPOSED — modeled on shipped
`resolve.go::compareRefs`/`copyTags`)". Plus Load-Bearing Decisions: "The
absent-key payload is a field on `Refusal`, not a new kind."

**Clarification I would ask.** What is the field called, and what is its type?
The RDR tells me the payload's *content* and its *sort order* but never names
the field or its element struct. `Refusal` today
(`internal/resolve/resolve.go::Refusal`) is flat scalars plus
`Rows []RowRef` and `MissingOwned []string`; the Existing Infrastructure
Audit row acknowledges this is "the kernel's first two-level payload." The
A3 status block mentions the spike's `UndecidedAtoms` only to *reject* the
spike's shape (it kept `Refusal.Guard` beside it), so I cannot take that name
as sanctioned. Concretely I need:

- the field name on `Refusal`;
- whether "reason" is a new exported string-typed constant set on the kernel
  (like `RefusalKind`) or a bare string — this matters because RDR 0005 will
  map it and `resolve.go::RefusalKinds` is the house pattern for closed sets;
- whether the "block" is the same exported constant type as the atom's block
  field, or a separately-spelled payload value.

**Decision it blocks.** The `Refusal` type change in Phase 2 and every
assertion in MVV scenario 1 and Testing Strategy rows 1, 11, 12, 20 — all of
which say "key in payload" without a name to assert on. This also blocks
Phase 3's contract-test export and RDR 0005's renderer, both of which type
against this field.

---

## IMP-3 — HIGH — `GuardAtom`'s field types are never fixed, and the `literal` type decides whether the kernel can compare the existence literal at all

**Anchored passage.** Normative Contracts SEAM clause: "A candidate row
carries its guard as a slice of parsed atoms — key, operator token, literal,
block ∈ {all, unless} — the shape JDR 0001 §D1 fixes; this RDR cites it and
does not restate the grammar."

**Clarification I would ask.** JDR 0001 §D1 says exactly the same four words
("key, operator token, literal, block") and nothing more — it explicitly
disclaims Go identifiers as "non-normative." So "cite, don't restate" cites a
document that also does not fix the shape, and no third document does. I need
the `literal` field's Go type before I can write one line of `evaluateAtom`,
because the existence rule is a *comparison against the literal*:

> "Its verdict is `presence == literal`: `exists = true` decides TRUE when
> the key is present, `exists = false` decides TRUE when it is absent."

If `literal` is a `string` (the shape `resolve.Tag.Value` uses today), the
kernel compares against the two exported boolean literal *forms* — which is
what A16 implies ("its literal by exactly two boolean literal forms … exported
as kernel constants"). If `literal` is a typed `any`/`bool`, the exported
forms are meaningless. If `in`/`contains` carry a *set* literal, a single
scalar `literal` field cannot hold it at all — and RDR 0003's matrix
(`0003-guard-predicate-exhaustiveness.md`, operator/kind matrix) gives `in` a
"non-empty typed scalar set" and `contains` a "non-empty typed element set."
A four-field atom with one scalar `literal` cannot represent those two of the
five operators.

**Decision it blocks.** The `GuardAtom` struct definition — Phase 1's first
commit — and with it every fixture migration. It also blocks Testing Strategy
row 11 ("a foreign boolean literal is rejected at load"), which cannot be
written until I know what a literal *is*.

---

## IMP-4 — HIGH — The payload sort key ("row identity then key") is not a total order when one row has two unevaluable atoms on the same key

**Anchored passage.** Normative Contracts, payload clause: "MUST sort the
payload by row identity then key, so the payload is a function of the input
tuple (RDR 0001 REQ-1), never of atom or row order."

**Clarification I would ask.** A5's sanctioned idiom puts *two atoms over one
key* on one row ("the value row MUST itself carry an `X exists = true` atom
beside `X eq v`"), and the P-5 footgun idiom does the same in `unless`
(`legal_hold exists = true` beside `legal_hold eq true`). Those two atoms
share a row identity and a key. If both go unevaluable, `(row, key)` does not
order them, and the payload is then a function of *atom order* — exactly what
this clause forbids. RDR 0003's Identity decision keys an atom by "its
position within `all` or `unless`", but position is atom order, so using it as
a tiebreak reintroduces the dependence the clause rules out.

Do I (a) extend the sort key to `(row, key, block, operator)`, (b) deduplicate
per `(row, key)` and report one entry, or (c) add a positional field to the
atom and accept that the payload depends on it? Each is observably different
under Testing Strategy row 17.

**Decision it blocks.** `compareAtoms`'s implementation (Testing Strategy row
17 names it as PROPOSED) and the determinism test that pins it. Getting this
wrong reproduces the class of defect frozen
`adversarial_test.go::TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder`
exists to prevent — but one level down, at atoms.

---

## IMP-5 — MEDIUM — Whether `Refusal.Guard` is deleted or retained is contradicted between A3's spike status and the Normative Contracts, and Phase 1 needs the answer on day one

**Anchored passage.** A3 Status: "the BOUND is verified; the SHAPE the spike
proved is not the shape this RDR specifies. The spike retained `Refusal.Guard`
(retyped to `[]GuardAtom`) and added `UndecidedAtoms` BESIDE it, whereas the
Normative Contracts and JDR 0001 §D4 both say the payload REPLACES
`Refusal.Guard`." Plus its Re-verification plan: "re-run the reshape with
`Refusal.Guard` actually removed and record the frozen suite's disposition for
that one test."

**Clarification I would ask.** A3 is `Pending` on precisely the question my
first commit answers, and the Prerequisites checklist says "A3 … verified
(Resolve)" is unchecked. But the RDR *also* tells me what the outcome will be:
the SEAM clause states that
`fixup_test.go::TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue` "keeps
its verdict but changes its REASON, and Phase 1 must re-read it rather than
re-encode it," and Testing Strategy row 19 repeats it. So: am I blocked on the
re-spike, or is the SEAM clause the ruling and the spike a formality? If it is
the ruling, is `TestFixup1d`'s guard-text assertion (shipped:
`fixup_test.go:64`, `if got.Refusal.Guard != escape.Guard`) *deleted*, or
*replaced* by an assertion on the new payload — and if replaced, replaced with
what, given IMP-2 leaves the field unnamed?

I verified the RDR's premise against source: `legalInput`
(`internal/resolve/fixtures_test.go::legalInput`) supplies `status` (owned)
and `reviews` (observed), and `Resolve` adds `recognized`
(`resolve.go::assemble`); the Fixup-1d escape guard is `iterations >= 3`
(`fixup_test.go:44`), whose key is indeed absent. So the RDR's reading is
correct — but "correct" is not the same as "decided," and A3's own Status
calls this "this assumption's own 'If wrong' branch."

**Decision it blocks.** Whether Phase 1 can start before a re-spike, and the
disposition of one frozen test — i.e. whether I am re-encoding a frozen
contract (allowed) or reopening one (which per the Trade-offs "Negative" bullet
must not happen: "the ADV/Fixup dispositions must be re-encoded, not
re-decided (A3)").

---

## IMP-6 — MEDIUM — The existence-operator constants are specified as three exported constants with no names, types, or package placement

**Anchored passage.** Normative Contracts, existence clause: "The kernel
recognizes the existence operator by one operator token, and its literal by
exactly two boolean literal forms, all three exported as kernel constants that
the normalizer MUST emit for existence atoms (A16)."

**Clarification I would ask.** Three exported constants in `internal/resolve`
with (i) no names, (ii) no stated type, and (iii) no stated *values*. Does the
token constant equal the string `"exists"` (RDR 0003's matrix spells the
operator that way — `0003-guard-predicate-exhaustiveness.md` operator/kind
matrix), and do the literal forms equal `"true"`/`"false"`? The RDR is careful
to say the kernel learns "one grammar fact," which implies the *values* matter
and are the drift boundary — yet the values are the one thing not written
down. A16's Evidence quotes JDR §D4's pricing of the exception but likewise
never gives the strings.

Second half of the question: A16's own Status says "The PRODUCER half is not
yet a duty on 0002," and Prerequisites item (a) is unchecked. Since RDR 0002
has no normalizer implementation yet, is Phase 1 free to *choose* the values
and have 0002 import them later, or must I wait? The Risks section says "the
constant is exported from the kernel and the normalizer imports it," which
reads as "kernel chooses" — but that is a mitigation sentence, not a grant.

**Decision it blocks.** The exported constant block in Phase 1 ("export the
existence token and literal constants") and Testing Strategy row 11's
foreign-token test, which needs a known-good token to contrast against.

---

## IMP-7 — MEDIUM — A22's fallback and its normative clause disagree about whether the kernel may assume canonical keys, and only the fallback is implementable today

**Anchored passage.** PRESENCE IS PROVENANCE-BLIND clause: "Key identity is
exact string equality on the key as assembled; canonicalization of authored
key spellings is the normalizer's (RDR 0002), upstream of the kernel, and the
kernel performs none (A22)." Plus A22 Status: "Pending — kernel half Verified;
the producer half is an UNSTATED INFERENCE against RDR 0002's current text …
Fallback if 0002 declines: narrow A22 to kernel-side only (the kernel performs
no canonicalization and requires canonical keys as an INPUT PRECONDITION)."

**Clarification I would ask.** From the kernel's side these two produce
*identical code* (exact string equality, no normalization) — so is there any
implementation difference at all, or is this purely a documentation question
for the doc comment on `Row`/`GuardAtom`? If it is purely documentation, the
Prerequisites item (b) should not read as gating my work. If it is *not* —
e.g. if the "INPUT PRECONDITION" reading implies the kernel should validate or
reject non-canonical keys at `Resolve` entry — that is a new precondition check
alongside the ones JDR 0001 §JD-5 already leaves unordered (0009's breach
check, 0008's reserved-key check), and §JD-5 is explicitly open. I need to know
which, because "precondition" in this codebase has a specific meaning.

**Decision it blocks.** Whether Phase 1 adds any entry-point validation, and
what the `GuardAtom.Key` doc comment promises. Low code cost either way, but it
determines whether I am touching `Resolve`'s entry, which the Technical Design
says is "untouched in control flow."

---

## IMP-8 — MEDIUM — "The kernel MUST evaluate every atom of every survivor — no short-circuit" collides with the seam being consulted for present-key value atoms

**Anchored passage.** Payload clause: "The kernel MUST evaluate every atom of
every survivor — no short-circuit on a decided block — and MUST sort the
payload by row identity then key".

**Clarification I would ask.** "Survivor" is defined by the SURVIVOR
MEMBERSHIP clause as a post-verdict property ("A row whose guard is GuardTrue
or GuardUnevaluable is a SURVIVOR; only GuardFalse rows are pruned") — but I
cannot know a row's verdict until I have evaluated its atoms. So the
no-short-circuit rule necessarily applies to *every candidate row*, including
ones that turn out to be pruned. That means the seam is called for every
present-key value atom on every candidate row, even after an `F` is found in
the same block.

Is that intended? Two consequences I would want confirmed before writing it:
(a) `Evaluate` is now called on atoms whose verdict cannot affect the row —
the seam is RDR 0003's code and this is an observable call-count contract on
it, which nothing states; (b) the Performance Expectations section says "the
kernel replaces one seam call per row with one presence lookup per atom plus a
seam call per present-key value atom," which is consistent with no
short-circuit — but only if "per candidate row" is meant, and it says nothing
about pruned rows.

The RDR gives the *reason* for no-short-circuit ("Every undecidable row appears
in the payload"), which only requires exhaustive evaluation on rows that end up
`U` — a weaker rule I could implement with short-circuiting on a decided `F`.
The stronger literal reading and the weaker justified reading differ in
observable seam calls.

**Decision it blocks.** The `evaluateGuard` loop structure and whether Phase 3's
exported contract test may assume the seam is called at most once per atom.

---

## IMP-9 — LOW — Testing Strategy row 8 and Phase 3 both say the `contains` present-key case cannot be written; row 8 is still listed as a scenario to implement

**Anchored passage.** Phase 3: "**Blocked on one RDR 0003 declaration.**
`resolve.Tag.Value` is a bare `string`, so a set-valued tag reaches the seam as
one opaque string with no stated element encoding. The `contains` leg of this
contract test therefore cannot be written from any current document, and
neither can Testing Strategy scenario 8's present-key half."

**Clarification I would ask.** Row 8 in the Testing Strategy matrix reads
"`contains` over an ABSENT set-valued tag ⇒ unevaluable, NOT the empty set" —
which is the *absent* half and is fine. Phase 3 says scenario 8's "present-key
half" is blocked. But row 8 as written has no present-key half. Is there a
missing row, or is the Phase 3 sentence referring to something that was cut? I
would want to know whether I am expected to deliver a present-key `contains`
test in Phase 2 and am simply blocked, or whether Phase 2's matrix is complete
as written and only Phase 3 is blocked.

I confirmed the premise: `internal/resolve/resolve.go::Tag` is
`{Key string; Value string}`, so there is genuinely no element encoding.

**Decision it blocks.** Phase 2's exit criteria — whether the domain-rule
matrix is deliverable in full or ships with a known gap that needs recording.

---

## IMP-10 — LOW — `RequiresOwned`'s narrowed meaning is stated as a doc contract, but nothing in the kernel can enforce it, and the RDR does not say that is acceptable

**Anchored passage.** `Row.RequiresOwned` clause: "Row.RequiresOwned names the
owned tag keys the row's post-guard transition depends on — the keys its
`Writes` require … Who populates the field is RDR 0002's (JDR 0001 §JD-3)."
Plus A21 Status: "the PRODUCER question stays open at §JD-3 and is carried as a
Capability Dependency, not as a claim of this RDR."

**Clarification I would ask.** Phase 1 says "narrow `Row.RequiresOwned`'s doc
contract" — a comment change only. So after Phase 1, `RequiresOwned` is
*documented* as post-guard write dependencies while `resolve.go::missingOwned`
continues to consume it identically to before, and no producer exists to
populate it under the new meaning. Is any test expected to pin the narrowed
meaning, or is it documentation-only until RDR 0002 lands a producer? Testing
Strategy row 15 pins only the *escape-row* composition ("Escape row raises no
`owned_state_unavailable` of its own"), which — as A21 notes — holds trivially
because `missingOwned` over an empty slice never enters the loop
(`resolve.go::missingOwned`), so that test passes today and would pass under
either meaning. It is a non-discriminating test for this clause.

**Decision it blocks.** Whether kata `xg7p` closes at Phase 1 (doc change) or
stays open until 0002's producer lands. The Trade-offs section claims "kata
`xg7p` closes against behavior, not against a doc comment" once Phases 1–2
land — but the only behavior change I can identify is the domain rule, not
`RequiresOwned` itself, whose consumption is unchanged.

---

## IMP-11 — LOW — [out-of-delta] The MVV requires the escape row to be reachable, but the RDR's own gate ordering means the escape path is never entered in scenario 1

**Anchored passage.** Minimum Viable Validation: "one candidate row whose guard
is a value atom over an absent key, whose `RequiresOwned` is SATISFIED, and a
modeled `no_match` escape row yield `Refusal.Kind == guard_unevaluable`".

**Clarification I would ask.** Per the GATE, THEN COUNT clause and shipped
`resolve.go::Resolve` (`selected, blocked := gate(...); if blocked != nil {
return refuse(in, *blocked), nil }`), a `guard_unevaluable` from the candidate
set returns *before* `escapeOrRefuse` is ever called — the RDR says so
explicitly ("the gate returns before the escape path is reached"). So the
modeled escape row in MVV scenario 1 is inert: the same assertion passes with
the escape row deleted. Is the escape row there to document the *pre-RDR*
masking path (the Background probe that produced `Escaped:true`), or is it
load-bearing to the assertion?

I would keep it — it makes the probe's inversion legible — but I would want the
MVV to say so, because as a test it currently reads as if the escape row were
part of what is being exercised, and a later reader may "simplify" it away and
lose the regression's meaning. Contrast with the MVV's careful note about why
the satisfied `RequiresOwned` *is* load-bearing; the escape row gets no such
note.

**Decision it blocks.** Nothing structural — it blocks only how I write the MVV
test's comment, and whether I add an explicit "escape row present and not
consulted" assertion.
