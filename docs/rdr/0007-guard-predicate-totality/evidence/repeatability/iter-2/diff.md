Model: claude-opus-5

# Repeatability diff — RDR 0007, iteration 2

Runs compared (all `variant: full`, profile foundational):

- **run-1** — `claude-opus-5`
- **run-2** — `claude-fable-5` (alt-model)
- **run-3** — `claude-sonnet-5`

Each finding is labeled **SILENCE** (the RDR never says), **AMBIGUITY** (the
RDR says it in a way that reads two ways), or **RUN ERROR** (the RDR pins it
clearly and one run got it wrong — not an RDR defect, no rewrite owed).
Findings are ordered by how many runs disagree.

---

## 1. Disagreements

### R-1 — `GuardResult` constant ordering vs. the K3 `min` ordering — 3-way — AMBIGUITY

**Contract**: the concrete ordering of the `GuardResult` enum, on which the
normative combination rule ("Conjunction is `min` under `F < U < T`") depends
if implemented as an integer `min`.

- **run-1**: `GuardTrue = iota`, then `GuardFalse`, then `GuardUnevaluable` —
  i.e. `T < F < U`, which inverts the normative order.
- **run-2**: `GuardTrue = iota`, `GuardFalse`, `GuardUnevaluable`, and marks
  the ordering itself GUESS.
- **run-3**: `GuardFalse = iota`, `GuardTrue`, `GuardUnevaluable` — a third,
  distinct order.

**RDR passage**: Normative Contracts, combination clause, states the lattice
order `F < U < T` and the K3 table but says nothing about the Go constant
ordering; Capability Dependencies calls `GuardResult` "Reused unchanged" from
RDR 0001 without reproducing its declaration. The RDR states `min` over a
named order in one place and treats the type as shipped-and-fixed in another,
so a reader cannot tell whether `min` is literal integer comparison over the
declared constants (in which case declaration order is normative) or a
table-driven operator over an abstract lattice. All three runs answered
differently and none matched `F < U < T`.

**Lands in**: Technical Design → **Normative Contracts**, combination clause.

---

### R-2 — where per-atom unevaluability reasons are produced and carried — 3-way — SILENCE

**Contract**: whether the verdict pass returns the `UndecidedAtom` records
alongside the row verdict, or the payload is minted by a second walk over the
atoms after `gate` knows the row is undecidable.

- **run-1**: `evaluateGuard(row, view, seam) GuardResult` — verdict only; a
  separate `undecidedAtoms(undecided, view, seam)` re-walks the atoms in
  `gate` to mint the payload. Run-1 explicitly flags this as GUESS
  ("whether the verdict pass carries reasons forward or the payload pass
  re-walks the atoms").
- **run-2**: `evaluateGuard(row, view, seam) GuardResult` in the signature,
  but the helper narrative and pseudo-code have `gate` accumulate
  `(row, its U-atoms + reasons)` during the verdict loop — reasons carried
  forward. Marked GUESS.
- **run-3**: changes the signature outright to
  `evaluateGuard(atoms []GuardAtom, view) (GuardResult, []UndecidedAtom)` —
  reasons are a second return value, and the seam is dropped from the
  parameter list entirely.

**RDR passage**: SILENCE. The Existing Infrastructure Audit names
`evaluateGuard` as the "per-row verdict hook" to "Extend" and the payload
clause fixes the payload's *shape* and *sort*, but nothing states the
producer of the reason values or the internal signature. Alternative 1's
aside about a seam returning `(GuardResult, unevaluable atoms)` is about the
rejected row-level seam, not the kernel's internal hook, so it does not
settle it. A re-walk and a carry-forward are behaviourally different when the
seam is impure or expensive (the re-walk calls `Evaluate` twice for every
`uncomparable` atom).

**Lands in**: Technical Design → **Normative Contracts**, payload clause
(add: the verdict pass emits reasons; the payload is not re-derived).

---

### R-3 — how a `Row` declares itself an escape row — 3-way — RUN ERROR (run-3), otherwise SILENCE on element type

**Contract**: the shape of `Row.Escape`.

- **run-1**: `Escape []string` — "GUESS: refusal-class names".
- **run-2**: `Escape []RefusalKind` — GUESS on element type.
- **run-3**: `Escape bool` — "GUESS: representation".

**RDR passage**: The RDR repeatedly says "a non-empty `Escape` list" (A21,
A25, RDR 0009 precondition quote, Failure Modes) and RDR 0002's closure is
quoted as "An `escape` list MUST contain only … `no_match` and
`ambiguous_match`". That it is a **list of refusal classes** is pinned;
run-3's `bool` is a **RUN ERROR** against clear text, not an RDR defect. The
residue — whether the element is `string` or `RefusalKind` — is genuine
SILENCE, but it is a type-spelling detail with no behavioural consequence.

**Lands in**: no rewrite owed for the run error. If anything is added, it
belongs in Technical Design → **Normative Contracts** as a one-clause
statement that `Escape` is a list of refusal kinds — low value.

---

### R-4 — `GuardEvaluator` seam declared as an interface vs. a func type — 2-way, splits on the model boundary — SILENCE ⚑

**Contract**: the Go form of the value seam.

- **run-1**: `type GuardEvaluator interface { Evaluate(atom GuardAtom, value string) GuardResult }`, flagging interface-vs-func as GUESS.
- **run-3**: identical interface form.
- **run-2** (alt-model, `claude-fable-5`): `type GuardEvaluator func(atom GuardAtom, value string) GuardResult` — "GUESS: interface vs func type; func type chosen for a minimal seam."

**RDR passage**: SILENCE. The Normative Contracts SEAM clause writes the seam
as `Evaluate(atom, value) GuardResult` and Load-Bearing Decisions repeats
"per-atom `Evaluate(atom, value) GuardResult`". The method-style spelling
`Evaluate(...)` implies a named method, but the RDR never says
interface-or-func, and run-2's pseudo-code consequently calls
`seam.Evaluate(atom, value)` against a declared func type — an internal
inconsistency the RDR's silence permitted. Flagged: this divergence splits
cleanly along the alt-model boundary.

**Lands in**: Technical Design → **Normative Contracts**, SEAM clause.

---

### R-5 — reason attached to an atom left unevaluable by a NIL seam — 2-way — SILENCE

**Contract**: which member of the closed `Reason` set (`absent` /
`uncomparable`) a value-comparing atom over a **present** key gets when
`seam == nil`, given the kernel returns `GuardUnevaluable` for it.

- **run-1**: `evaluateAtom` step 4 says "A nil seam here yields
  `GuardUnevaluable`" and assigns no reason; the payload clause is left with
  no way to fill `Reason` for that atom.
- **run-3**: same gap — its pseudo-code only appends an `UndecidedAtom` in
  the absent branch and in the branch where the seam *answers* unevaluable;
  a nil seam is not modelled at all in `evaluateGuard`, so no payload entry
  is emitted for a nil-seam atom that nonetheless drives the row to `U`.
- **run-2**: models `elif seam == nil: U` as a distinct branch and states in
  its API summary that "a nil seam is never itself a payload reason" —
  reproducing the RDR's parenthetical but still not naming which reason the
  entry carries.

**RDR passage**: SILENCE with a near-miss. The payload clause closes the
reason set to `absent` ("the key was not in the view") and `uncomparable`
("the key was PRESENT and **the seam could not compare** the value, A18"), and
Testing Strategy row 19 adds "a nil seam is never itself a payload reason".
A nil-seam atom satisfies neither definition: the key *is* present, and no
seam answered. The kernel must therefore either emit an entry with a reason
neither definition covers, or emit no entry for an atom that made the row
undecidable — which contradicts "The refusal MUST name what was missing, per
row and per atom: for every undecidable row … for each of its unevaluable
atoms". Two runs silently dropped the entry; one noticed the hole and left it
open. This is the strongest rewrite candidate in the set.

**Lands in**: Technical Design → **Normative Contracts**, payload clause and
the nil-seam paragraph of the SEAM clause.

---

### R-6 — routing of `ambiguous_match` through the escape path — 2-way — RUN ERROR (run-3)

**Contract**: whether the `len(selected) > 1` branch consults escape
reachability.

- **run-1**: `default: return escapeOrRefuse(in, view, KindAmbiguousMatch)`.
- **run-2**: `default: return escapeOrRefuse(ambiguous_match, ...)`.
- **run-3**: `default: return refuse(ambiguous_match)` — no escape
  consultation, and its `case 0` comment scopes `escapeOrRefuse` to "existing
  no_match path".

**RDR passage**: The RDR pins this. Background quotes RDR 0002's closure —
escape lists may contain "`no_match` and `ambiguous_match`" — and the
GATE-THEN-COUNT clause says "`no_match` and `ambiguous_match` sit downstream
of all three gate facts", with escape reachability consulted after the count.
An `ambiguous_match` that cannot reach the escape path makes the modelled
`ambiguous_match` escape class dead surface. **RUN ERROR** against clear
text; no rewrite owed.

**Lands in**: no rewrite. (If the RDR ever restates the count branch it
belongs in Technical Design → **Normative Contracts**, GATE-THEN-COUNT.)

---

### R-7 — owned-state sweep: first missing key vs. accumulate over all survivors — 2-way — SILENCE

**Contract**: whether `MissingOwned` in an `owned_state_unavailable` refusal
carries the missing keys of one survivor or of every survivor.

- **run-1**: `if missing := missingOwned(survivors, view); missing != nil` —
  one call over the whole survivor set; accumulation semantics unstated.
- **run-2**: `missing = missingOwned(survivors, view)` — same, "owned
  provenance only".
- **run-3**: explicitly accumulates — `for (row,…) in survivors: missing +=
  missingOwned(row, view)` — a per-row call summed across survivors.

**RDR passage**: SILENCE. The RDR pins the *ordering* (owned sweep before the
undecidable check, over survivors) and the *field* (`MissingOwned []string`),
and it fixes payload determinism only for the `Undecided` payload — "The
kernel MUST … sort the payload so it is a function of the input tuple (RDR
0001 REQ-1), never of atom or row order." No parallel statement governs
`MissingOwned`, so run-3's accumulation across rows can produce a
row-order-dependent, duplicate-bearing slice while run-1/run-2's single call
may not. Since REQ-1 purity is claimed for the whole result, the omission is
a real gap rather than a taste difference.

**Lands in**: Technical Design → **Normative Contracts**, GATE-THEN-COUNT /
SURVIVOR MEMBERSHIP clause.

---

### R-8 — whether `escapeOrRefuse` reuses `gate` or is untouched — 2-way — AMBIGUITY

**Contract**: how the escape row set gets gated.

- **run-2**: "Reused identically for the escape row set by `escapeOrRefuse`
  (uniform gating…)" — i.e. `escapeOrRefuse` calls `gate`.
- **run-1**: lists `escapeOrRefuse` under "unchanged … (ZERO diff in the
  spike)" and does not model it as calling `gate`.
- **run-3**: "its internals are not re-derived here since the RDR says it
  needs 'ZERO changes'".

**RDR passage**: AMBIGUITY. The Technical Design says `escapeOrRefuse` is
"untouched in control flow"; Key Discoveries says "`escapeOrRefuse` gates the
escape set the same way"; the GATE-THEN-COUNT clause says "The escape set is
gated identically, as its own row set (D5)". "Gated identically" and
"untouched" are jointly satisfiable only if `escapeOrRefuse` already calls
`gate` today — which the RDR asserts but never states as the mechanism, so a
reader must choose between "delegates to `gate`, therefore inherits the new
payload for free" and "has its own gating that must be separately updated".
The escape-set payload cell in the `disposition` mini-check ("escape set's
payload") assumes the former without saying so.

**Lands in**: Technical Design → **Normative Contracts**, GATE-THEN-COUNT
clause (escape-set paragraph).

---

### R-9 — `Result` shape and the escape marker — 2-way — SILENCE

**Contract**: what `Resolve` returns on success and on refusal.

- **run-2**: models `Result{Plan *Plan; Refusal *Refusal; Escaped bool}`,
  citing the Background's observed `Escaped:true`, and marks the whole shape
  GUESS.
- **run-1**: does not model `Result` at all beyond `(Result, error)`, and
  explicitly says "the full `Input`/`Result` struct shape is not stated".
- **run-3**: same — signature marked GUESS, no struct.

**RDR passage**: SILENCE, and defensibly so — `Result` is RDR 0001's surface
and this RDR touches only `Refusal`. Recorded because `Escaped` is the only
field of it this RDR's own Background asserts behaviour about ("yields a plan
via the escape, `Escaped:true`"), and that assertion is the masking path the
RDR exists to close. Low-value rewrite; a citation to RDR 0001's `Result`
would close it.

**Lands in**: Technical Design (prose), not Normative Contracts.

---

## 2. GUESS clusters

Contracts two or more runs independently marked **GUESS**. These are the
candidate RDR rewrites, strongest first.

### G-1 — the `GuardAtom` struct's Go field names and types — 3 runs GUESS

run-1 GUESSes `Key`/`Op`/`Literal`/`Block` and the existence of named
`Operator` and `Block` types; run-2 GUESSes `Op string` and `Literal string`;
run-3 invents named types `OperatorToken` and `Literal` and GUESSes all
three. The RDR names the type `GuardAtom` (A26's banned-name sweep) and fixes
the four *fields* semantically ("key, operator token, literal, block ∈ {all,
unless}") and the payload's mirror fields by name (`Key`, `Block`,
`Operator`, `Literal`, `Reason`), but never spells `GuardAtom`'s own field
names. Because the payload clause requires `Block` to be "the same exported
constant type the atom carries", the atom's field types are load-bearing for
the payload, not incidental. → **Normative Contracts**, SEAM clause.

### G-2 — values of `OpExists`, `LiteralTrue`, `LiteralFalse` — 3 runs GUESS

run-1 and run-2 guess `"exists"` / `"true"` / `"false"`; run-3 declines to
guess a spelling at all (`"..."`). The RDR names all three constants and
makes the kernel compare "verbatim, performing no parsing, case-folding, or
coercion", and makes the normalizer's emission of exactly these values a MUST
(A16) — a contract whose whole content is byte equality, stated without the
bytes. → **Normative Contracts**, existence-operator clause.

### G-3 — the `Block` constant names and values — 3 runs GUESS

run-1 `BlockAll`/`BlockUnless` = `"all"`/`"unless"` (GUESS on both); run-2
identical, GUESS on spelling; run-3 names the type `BlockKind` and gives no
constants. A26 sweeps "the `Block` constants" against banned names, so they
are known to exist and to be exported, but they are never enumerated. →
**Normative Contracts**, payload clause (which already pins `Block` as an
exported constant type).

### G-4 — a `Reasons()` enumerator — 2 runs GUESS

run-1 and run-2 both posit `func Reasons() []Reason` from the clause "spelled
and enumerated the way `resolve.go::RefusalKind` / `RefusalKinds()` already
spell the refusal taxonomy, so … a test can pin the set exactly", and both
mark it GUESS; run-3 declares the constants with no enumerator. The clause
mandates the *pattern* and a test that pins the set, which is unsatisfiable
without an enumerator, but never states one. → **Normative Contracts**,
payload clause.

### G-5 — `Reason`'s underlying type and constant values — 3 runs GUESS

All three write `type Reason string` with `"absent"`/`"uncomparable"` and
mark the underlying type (run-1, run-2) or the exact strings (run-3) as
GUESS. Same root as G-2: the RDR fixes the closed set by constant name and
by English gloss, not by value. → **Normative Contracts**, payload clause.

### G-6 — the fifth refusal kind — 3 runs GUESS

All three note that the taxonomy is closed at five kinds and that only four
are named in this document (`no_match`, `ambiguous_match`,
`owned_state_unavailable`, `guard_unevaluable`); run-1 guesses "an
internal/table-error kind", run-2 guesses "`invalid_input`-style", run-3
declines. Not a defect of this RDR — the fifth kind is RDR 0001's and this
RDR reopens no kind — but it is a repeatability cost of the RDR asserting
"exactly five" without citing the fifth. → optional one-word citation in
Technical Design prose; no Normative Contracts change.

### G-7 — the Phase 3 exported value-seam contract-test signature — 2 runs GUESS

run-1: `func TestGuardEvaluatorContract(t *testing.T, mk func() GuardEvaluator)`;
run-2: `func TestGuardEvaluatorContract(t *testing.T, seam GuardEvaluator)`;
both GUESS entirely (run-1: "name and signature entirely"). run-3 does not
model it. The RDR commits to exporting a contract-test function that RDR
0003's implementation instantiates — a cross-RDR surface — without naming it.
→ Implementation Plan, Phase 3 (not Normative Contracts).

### G-8 — CLI transport of the payload — 2 runs GUESS, all 3 flag unresolved

All three correctly report the transport as unresolved between JDR 0001 §D4
(`Detail` flattening) and A19 (one `omitempty` structured field on
`clierr.CLIError`); run-2 and run-3 additionally GUESS the wire field name /
JSON encoding. This is a *deliberate* deferral the RDR states as such, and
all three runs reproduced the deferral faithfully — recorded as a GUESS
cluster only to note it is not a rewrite candidate.

---

## 3. Agreement

All three runs rendered these identically, confirming the RDR is determinate
on them: the `Refusal.Undecided []UndecidedRow` / `UndecidedRow{RuleID,
SourceLocator, Atoms}` / `UndecidedAtom{Key, Block, Operator, Literal,
Reason}` two-level payload with `Refusal.Guard` removed and no sixth refusal
kind; the closed reason set `{ReasonAbsent, ReasonUncomparable}` with
`uncomparable` reserved for a present-key seam answer; the four-way per-atom
dispatch (presence via provenance-blind `TagSet.Lookup`, existence decided by
the kernel as `presence == literal`, absent-key value atom unevaluable with
the seam never consulted, present-key value atom handed to the seam);
fail-closed drift (foreign token ⇒ value atom, foreign or missing literal on
`OpExists` ⇒ unevaluable/`uncomparable`, never presence-decided); the
strong-Kleene combination `all_result ∧ ¬(unless_conj)` with the `unless`
term dropping out on an empty or omitted block, an empty atom slice decided
TRUE without the view, `F ∧ U = F` pruning as witnessed falsity, and no
short-circuit; the gate-then-count order (prune GuardFalse → owned sweep over
survivors → undecidable veto → exact-one count → escape reachability) with
unevaluable rows counted as survivors and the resolution-level veto over a
decided-TRUE sibling; the six-field total sort tuple `(RuleID, SourceLocator,
key, block, operator token, literal)` with `slices.MinFunc` retired and every
undecidable row reported; the narrowed `Row.RequiresOwned` as post-guard
write-dependency keys, empty on escape rows, with the guard key set not
required to be a subset; `guard_unevaluable` as non-escapable; and the
absence of any new panic or Go-error channel from this RDR (the only Go error
being RDR 0009's malformed-table entry precondition).
