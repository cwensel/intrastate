Model: claude-opus-5[1m] (Opus 5, 1M context)

# Repeatability DIFF — RDR 0012

Runs compared (all `variant: full`, profile foundational):

- run-1 — `claude-opus-5[1m]`
- run-2 — `claude-sonnet-5`
- run-3 — `claude-fable-5-1`  ← alt-model; splits along this boundary are flagged

Coverage is as intended: run-1's `variant:` reads `full`, three runs
across three model families. No mismatch to flag.

---

## 1. Disagreements

Ordered by how many runs diverge.

### D1 — `GuardResult` constant ordering and underlying type — 3-WAY — `0012:C2` (kernel type, referenced not defined)

- run-1: `type GuardResult int`, `GuardTrue = iota, GuardFalse,
  GuardUnevaluable`; explicitly marks the ordering **GUESS**.
- run-2: declines to write the declaration at all — names only the
  three-valued set "`GuardTrue | GuardFalse | GuardUnevaluable`
  (existing type, referenced not redefined here)". No ordering, no
  underlying type.
- run-3: `type GuardResult int` with `GuardFalse = iota, GuardTrue,
  GuardUnevaluable` — the **inverse** of run-1's first two constants;
  marks both the int-ness and the ordering **GUESS**.

Silence: C1 freezes the seam *signature* and C2 names the three
verdicts and forbids a two-valued collapse, but no element in the
record pins the encoding, the iota order, or even that `GuardResult`
is an integer type. Two runs invented opposite orderings and the third
refused to render one — the clearest 3-way split in the set. It is
benign for behavior but not for a reconstruction: an implementer who
picks run-3's order makes the zero value `GuardFalse`, which silently
changes what an unconstructed verdict means.

### D2 — `in`-literal member representation — 3-WAY — `0012:C2`

- run-1: `literalMembers(atom)` helper, "1 for eq, n for in"; marks the
  member separator and the helper **GUESS**, stating the RDR "never
  says how an `in` literal encodes its members".
- run-2: `parseAllMembers(atom.Literal)` plus a separate
  `parseUnderKind(kind, atom.Literal)` returning a typed
  `parsedLiteral` — i.e. the literal is a delimited string parsed
  twice, once for the literal and once for the member list. Not marked
  GUESS at the representation level.
- run-3: names the fork explicitly as a silence — "`[]string` members
  alongside `Literal`, or `Literal` split by a fixed delimiter" — and
  leaves it open, carrying it into its terminal silences list.

Silence: C1 names the atom shape as `Key, Operator, Literal, Block`
(JDR 0001 §D1) and C2 states the poison rule over "each member, for
`in`", but nothing states how a single `Literal` string carries n
members. Every run had to invent an accessor, and run-2's version
additionally implies a second parse path the other two do not have.

### D3 — `bool` arm: is the literal checked before or inside the kind switch — 3-WAY — `0012:C2`

- run-1: checks held value tokens, then loops `literalMembers` checking
  each against `{true,false}`, then token equality. Literal check lives
  *inside* the `bool` arm.
- run-2: hoists a `literalOK, parsedLiteral := parseUnderKind(kind,
  atom.Literal)` check **above** the kind switch, so the unparseable-
  literal arm fires for every kind before dispatch; the `bool` arm then
  only re-checks the held value.
- run-3: same as run-1 — both checks inside the arm, held first then
  members.

Silence: C2 lists "unparseable LITERAL: GuardUnevaluable (defense in
depth)" as a *bullet peer* of the per-kind bullets rather than as a
step in an ordering, so whether it is a pre-dispatch gate or a
per-arm obligation is undetermined. The verdicts agree; the control
flow does not, and run-2's hoist means an unknown kind token would hit
the literal parse before the `default` arm — an order C2's
"dispatch ... is EXHAUSTIVE ... and the fallthrough `default` means
GuardUnevaluable" arguably contradicts.

### D4 — the `default` operator arm / `exists` disposition — 3-WAY — `0012:C2`

- run-1: a trailing `default: return GuardUnevaluable` for an
  unrecognized operator, marked **GUESS**, with `exists` noted as
  never reaching the seam.
- run-2: an explicit `case "exists": panic("unreachable: exists
  resolved upstream")` — a **panic**, not a verdict — and *no* operator
  `default` arm at all (the pseudo-code's switch falls off the end with
  no return).
- run-3: `exists` not present as a case; no operator `default` arm
  rendered.

Silence: C2 says only "`exists` never reaches the seam (`0007:C1`)"
and says nothing about an operator the switch does not recognize. One
run makes it a refusal, one makes it a crash, one leaves it
unhandled. This is the disagreement with the widest blast radius: a
panic at a seam the record describes as total and LOUD-not-crashing
(C1's nil-mapping disposition) is a different product.

### D5 — lint's `atomAdmitsValue` post-change return type and disposition — 3-WAY — `0012:C4`

- run-1: assumes `resolve.GuardResult` and that
  `analysis.go::nodeMeetsAll` switches on it; explicitly marks the
  return type **GUESS** and enumerates the live alternatives (a
  three-state lint enum, a `(bool, bool)`).
- run-2: does not render `atomAdmitsValue`'s new signature at all —
  only reproduces the construction/threading rule.
- run-3: renders the return as "three-valued", and additionally
  **guesses the policy**: "lint over-approximates an unevaluable atom
  as admissible AND emits a finding", while noting the record leaves
  the policy undecided.

Silence: C4 forbids the `!= GuardFalse` collapse and requires
GuardUnevaluable be its own case at every consumer, then explicitly
declines to decide "WHAT lint then does with an unevaluable atom". It
does not say what *type* carries the now-three-valued answer. run-3
filled the deliberately-open policy with a concrete both-arms answer;
run-1 correctly left it open. The contract's own deferral is being
read by one run as an invitation to pick.

### D6 — home package and nil-handling of `DeclaredKinds` — 2-WAY + BOUNDARY — `0012:C1`

- run-1: homes it in `internal/guard` but marks the home **GUESS**,
  naming `internal/table` as "equally import-legal"; `DeclaredKinds(nil)`
  returns an empty non-nil map, marked **GUESS**.
- run-2: homes it in `internal/guard` unqualified, no GUESS; says
  nothing about a nil model. Additionally asserts it "reuses the
  existing lookup surface `internal/guard/declaration.go::DeclarationOf`
  internally per the infra audit" — an implementation obligation
  neither other run states.
- run-3: `internal/guard`, no GUESS on the home; `DeclaredKinds(nil)`
  returns an empty non-nil map, marked **GUESS** ("record silent").

Silence: C1 gives the signature `func DeclaredKinds(m *table.Model)
map[string]string` with **no package qualifier** and no nil-input
clause. run-1 and run-3 agree on nil behavior and mark it; run-2 is
silent on nil and instead binds an internal reuse obligation lifted
from the infrastructure audit rather than from a contract. The
audit-vs-contract boundary is doing load-bearing work here that the
normative text does not.

### D7 — the published conformance fixture: identifier and var-vs-func — 2-WAY — `0012:C3`

- run-1: `func ContractKinds() map[string]string`, name marked
  **GUESS**, noting "a package-level `var` is equally likely".
- run-2: does not render the fixture's identifier at all — describes
  it only as a kernel-owned `map[string]string`.
- run-3: `var ContractKinds = map[string]string{...}` — same invented
  name as run-1 but the **opposite form** (var, not func), with both
  the name and the form marked **GUESS**.

Silence: C3 makes the fixture's existence, type, ownership and
kind-coverage normative ("publishes the declaration fixture its cases
assume as a kernel-owned `map[string]string`") but never names it. Two
runs independently converged on `ContractKinds` and then diverged on
whether it is a var or a constructor — a var is mutable shared test
state, a func returns a fresh copy. That is a real difference in a
fixture multiple packages import.

### D8 — the `probeRow` / per-row verdict aggregation — 2-WAY, BOUNDARY-SPLIT (run-3 vs run-1/run-2) — `0012:C4`

- run-1: renders only the threading rule — `probeRow` RECEIVES the
  constructed `resolve.GuardEvaluator` — and does not render the loop
  body or aggregation.
- run-2: same; threading rule only, "NOT rebuilt per row".
- run-3: renders a full `probeRow` body and surfaces a question neither
  other run raises — "**GUESS**: first non-True verdict short-circuits,
  Unevaluable preferred over False? — record silent" — plus a guess
  that `probeRow` passes `""` for an absent key, and a guess at the new
  parameter's position.

Silence: C4 fixes *where* the evaluator is constructed and that it is
threaded, not *how a row with several guard atoms combines verdicts*.
MVV step 2 says the refusal "names the guarded row and the `iter`
atom" — singular — which presumes but does not state a
first-unevaluable-wins rule. Only the alt-model run noticed. This is a
boundary split in the useful direction: run-3 found a genuine silence
the same-family runs rendered past.

### D9 — the C5 canonicality check's factoring and error type — 2-WAY — `0012:C5`

- run-1: `func (l *loader) conformCanonicalInt(decl TagDecl, site
  string, members []string) error`, signature "entirely **GUESS**".
- run-2: **no new function** — the check is "folded into
  `(*loader).atom`" beside the existing `conform` call, with no
  extracted helper rendered.
- run-3: "No new exported function", but guesses a private
  `canonicalIntSpelling(decl, members, site) error`, and separately
  flags the **Go error type** as a silence ("record names the code,
  not the type").

Silence: C5 fixes the rule, the site (`normative.go::(*loader).atom`,
beside `conform(...)` at `normalize.go:172`), the refusal code and the
message text, but neither the decomposition nor the error value's Go
type. Only run-3 named the error-type silence.

### D10 — the S4 zero-value backstop's packaging — 2-WAY — `0012:S4` (Validation scenario 4; no contract id)

- run-1: does not render it in the helper list at all.
- run-2: names it as a near-miss fourth helper and explicitly declines
  to invent a name — "I have not invented one above to avoid
  manufacturing a false API surface."
- run-3: names it as an honourable mention and guesses the packaging —
  "a small `cmd`/`tools` program or a Go test under `internal/guard`
  with a build tag".

Silence: S4 fixes that the check is typed (`go/ast` or
`x/tools/go/packages`), what it matches, and that it is wired into
`make check` — but not where it lives in the tree. The scenario is the
*only* backstop for the graphlint site until C4 lands, so its home
being unstated is load-bearing, and it lands on a validation scenario
rather than a contract, which is why no id names it.

### D11 — `Evaluator` receiver form and constructor return — 2-WAY — `0012:C1`

- run-1: value receiver, value return, with a paragraph of reasoning
  and an explicit **GUESS**; argues a `*Evaluator` return would make
  `Evaluator{}` a differently-shaped hazard.
- run-2: value receiver `(e Evaluator)`, value return; **not** marked
  GUESS. Instead marks a different thing GUESS — "exact type name
  `Evaluator` vs `evaluator` for the unexported struct" — which the
  other two runs do not treat as open (C1 writes `NewEvaluator(...)
  Evaluator` and `Evaluator{}`).
- run-3: value receiver `(ev Evaluator)`, value return; not marked
  GUESS at the receiver, marked GUESS only at the field name.

Silence: C1 writes the constructor's return as bare `Evaluator` and
speaks of the "zero-value `Evaluator{}` construction form", which does
pin the value shape more tightly than run-1 credits — but it never
states the receiver form on `Evaluate`, and the value receiver is what
makes the nil-map read safe (C1's LOUD-not-panic disposition). All
three landed in the same place; only one treated it as inferred.

---

## 2. GUESS clusters

Contracts two or more runs independently marked GUESS. These are the
candidate rewrites, ordered by cluster size.

### G1 — `in`-literal member encoding — 3/3 runs — `0012:C2`
run-1 ("never says how an `in` literal carries its members"), run-2
(invents `parseAllMembers`/`parseUnderKind` without pinning the shape),
run-3 (lists it first among its terminal silences). C2's poison rule is
the RDR's sharpest behavioral claim and it is stated over a member
sequence the atom shape does not visibly contain. Highest-value
rewrite in the set.

### G2 — `GuardResult` encoding / constant ordering — 2/3 explicit GUESS, 3/3 divergent — `0012:C2`
run-1 and run-3 both marked it GUESS and then chose **opposite**
orderings; run-2 declined to render it. A one-clause pin ("the
existing `resolve.GuardResult` encoding is unchanged; this RDR fixes no
ordering") would close it, or a cross-reference to wherever 0007 pins
it.

### G3 — the conformance fixture's identifier and form — 2/3 — `0012:C3`
run-1 (`func ContractKinds()`, GUESS) and run-3 (`var ContractKinds`,
GUESS) converged on the name and split on the form; run-2 refused to
name it. C3 already makes the fixture's type and ownership normative —
naming it is one clause.

### G4 — `DeclaredKinds` on a nil model — 2/3 — `0012:C1`
run-1 and run-3 both marked it GUESS and both chose "empty non-nil
map, by the same LOUD-not-crash reasoning"; run-2 is silent. The two
that noticed agreed, which makes this cheap to close and low-risk
either way.

### G5 — the private dispatch helper's shape inside `Evaluate` — 3/3 — `0012:C2`
run-1 (`kindOf` + `evalEq`/`evalIn` + `parseInt`, all GUESS), run-2
("GUESS at exact internal shape"), run-3 (`compareTyped(kind, op,
held, members)`, GUESS). All three flagged it; all three invented a
different decomposition. Arguably *should* stay open — decomposition
is implementation latitude — but the runs' disagreement on whether
`eq` and `in` share an arm (run-1 splits them, run-3 merges them via
`members`) touches the poison rule, which is not latitude.

### G6 — `GuardAtom.Block`'s type and domain — 2/3 — `0012:C1`
run-1 (GUESS on the spellings `match | guard.all | guard.unless` and
on whether `Operator`/`Block` are named string types), run-3 (GUESS on
`Block`'s type). run-2 renders the four field names without types.
C1 cites JDR 0001 §D1 by reference; the field types never appear.

### G7 — the C5 refusal's Go error type / factoring — 2/3 — `0012:C5`
run-1 (whole signature GUESS), run-3 (private helper GUESS + error
type GUESS). run-2 declines to extract a helper.

### G8 — the `flow-guard-unevaluable` envelope's field names — 3/3 — `0012:C5` adjacent / MVV step 2
All three runs reached the same conclusion by the same route: MVV step
2 states the *semantic* content ("names the guarded row and the `iter`
atom with reason `uncomparable`") and S8 pins only the
*`flow-tag-invalid`* envelope (`code`/`message`/`schema_version`/
`param`). run-1 and run-2 both explicitly widened to S8's shape as the
convention to reuse and both marked the row/atom field keys GUESS;
run-3 called the payload shape "unchanged" without naming keys. The
record has a normative fixture for the control envelope and none for
the new one — that asymmetry is the finding. Lands on the Validation
Scenarios section (no contract owns the unevaluable envelope).

---

## 3. Agreement

All three runs rendered identically, without divergence or GUESS: the
`NewEvaluator(kinds map[string]string) Evaluator` and
`DeclaredKinds(m *table.Model) map[string]string` signatures; the
totality rule with empty-`Kind` keys OMITTED not mapped to `""`; the
five-kind vocabulary and its spelling; the frozen `Evaluate(atom,
value) GuardResult` seam signature; the full `eq`/`in` verdict matrix
including parse-all-then-compare poisoning, `bool` token-only
equality, the shared `enum`/`scalar` arm, `set` → GuardUnevaluable,
key-absent → GuardUnevaluable, and `default` → GuardUnevaluable rather
than string equality; the operator-inferred `lt/lte/gt/gte` and
`contains` arms not consulting the mapping; the widened
`TestGuardEvaluatorContract(t, newSeam func(...) GuardEvaluator)`
signature; the three C4 construction sites and the thread-don't-rebuild
rule at both `probeRow` and `atomAdmitsValue`; the `strconv.Itoa`
round-trip canonicality rule at `(*loader).atom` with
`malformed_predicate_atom` reused and `conformKind` left untouched;
and "nothing new is persisted, no envelope field, no sixth refusal
kind". These confirm the RDR is determinate across its core.

---

## Signal

**HEALTHY.**

The diff localizes to a small, specific set of interfaces rather than
spreading across the record: the `in`-literal member encoding (C2),
the `GuardResult` encoding (C2), the published fixture's identifier
and form (C3), lint's three-valued return type (C4), and the new
refusal envelope's field names (MVV/Validation). Five interfaces, and
every one of them sits in the same two places — C2's comparison
mechanics and the test/consumer surfaces around them.

GUESS markers cluster on the same under-specified contracts across
runs rather than scattering: G1 (member encoding) and G5 (dispatch
decomposition) were flagged by 3/3; G2, G3, G4, G6, G7 by 2/3, with
the two flaggers usually being run-1 and run-3 — different model
families reaching the same silences independently. That convergence is
the health signal.

The failure mode this check exists to catch is absent: the runs are
not identical-but-confidently-wrong. They disagree substantively (D1's
inverted constant ordering, D4's panic-vs-verdict, D5's filled-in lint
policy) and they mark their uncertainty, so the disposition is
**edit the record**, not rerun the lens on another model.

One asymmetry worth naming for the resolver: the alt-model run (run-3)
carries a terminal silences list and surfaced two silences neither
same-family run reached — per-row verdict aggregation (D8) and the C5
refusal's Go error type (D9). run-2 is the most conservative renderer
(it declines to name identifiers it cannot ground) and therefore
contributes the fewest GUESS markers but also the fewest findings;
its divergences are mostly *omissions* rather than alternatives,
except D3's pre-dispatch literal hoist and D4's panic, which are
genuine alternative readings.
