Model: claude-opus-5[1m]

# Repeatability diff — RDR 0007 guard predicate totality

Three full-variant reconstructions, three distinct base models:

| run | model | variant |
| --- | --- | --- |
| run-1 | `claude-opus-5[1m]` | full (profile: foundational) |
| run-2 | `claude-fable-5` | full (profile: foundational) |
| run-3 | `glm-5.2:cloud` | full (profile: foundational) |

All three are cross-model draws, so every 2-vs-1 split is a model-boundary
split; each finding below names the dissenting model.

## 1. Disagreements

### D-1 — Does an unevaluable row join `survivors` for the `missingOwned` scan? (3-way)

**Contract**: `gate`'s row bookkeeping — whether a `GuardUnevaluable` row is a
member of the survivor set that `missingOwned` scans, or only of the
`undecidable` list.

- **run-1** (`opus-5`): appends the row to **both** `survivors` and
  `undecidable`, and flags the choice as an explicit GUESS: "whether an
  unevaluable row also joins `survivors` for the `missingOwned` scan, or only
  `undecidable`. The RDR says the owned check runs 'among surviving rows' … —
  which only bites if it does."
- **run-2** (`fable-5`): `survivors ← rows minus decided-GuardFalse rows` —
  unevaluable rows are survivors by construction (subtractive definition). Not
  marked GUESS.
- **run-3** (`glm-5.2`): `survivors := [r for r in rows if Evaluate(...) !=
  GuardFalse]` — same subtractive definition, then `if any survivor
  GuardUnevaluable`. Not marked GUESS.

Net semantics coincide, but only run-1 recognized this as underdetermined; runs
2 and 3 derived it silently from a subtractive reading. The RDR never states the
membership rule constructively.

**RDR passage**: Normative Contracts define the survivor set only by what it
excludes — "Only decided-GuardFalse rows are pruned from consideration" — and
state the precedence as "Among surviving rows, absent owned state is reported
BEFORE an undecidable guard: a survivor set carrying both MUST refuse
`owned_state_unavailable`." The clause presupposes an unevaluable row is a
survivor without saying so. A3's Evidence describes the shipped mechanism
(`gate` calls `missingOwned(survivors, view)`, "both fail if `rows` replaced
`survivors`") but never says which verdicts populate `survivors`.

**Lands in**: Technical Design → Normative Contracts (the aggregation /
ordering block). State constructively: a row whose guard is GuardTrue or
GuardUnevaluable is a survivor; only GuardFalse rows are pruned.

### D-2 — `Refusal.Rows` element type (3-way)

**Contract**: the type of the discriminator field the RDR makes normatively
load-bearing.

- **run-1** (`opus-5`): `Rows []Row`, marked GUESS — "could be `[]Row`,
  `[]RuleID`, or `[]RowRef`."
- **run-2** (`fable-5`): `Rows []RowRef`, marked GUESS — "row ids? `*Row`?"
- **run-3** (`glm-5.2`): `Rows []RowID` with an invented `RowID` struct
  carrying `RuleID` + `SourceLocator`, marked GUESS.

Three different types, three different GUESS notes. Same for the sibling
`MissingOwned` (run-1 `[]string` GUESS, run-2 `[]string`, run-3 `[]TagKey`).

**RDR passage**: "`Refusal.Rows` carries every undecidable row. `Refusal.Guard`
is single-valued — the kernel selects the lowest row by `(RuleID,
SourceLocator)` — so it is a diagnostic convenience, NOT the discriminator: any
contract or test distinguishing *which* row went unevaluable MUST assert on
`Rows`." A normative MUST-assert-on-`Rows` with no stated element type; Testing
Strategy scenario 5 compounds it ("give the two rows distinct `RuleID`s and
assert on `Rows`") without saying what `Rows` holds.

**Lands in**: Technical Design → Normative Contracts (the refusal-payload
block). Naming the element type is required for the Phase 2 vectors to assert
on `Rows` at all.

### D-3 — Whether guard text that cannot be mapped to atoms is a reachable case, and what it yields (3-way)

**Contract**: `Evaluate`'s behavior when `recoverAtoms` fails.

- **run-1** (`opus-5`): no failure branch in the pseudo-code; raises the
  question only as an error-mode GUESS — "whether `Evaluate` may return a Go
  `error` at all (e.g. for a guard string the evaluator cannot parse) … the RDR
  never says so."
- **run-2** (`fable-5`): explicit branch `if atoms unrecoverable: return
  GuardUnevaluable`, flagged GUESS that the case exists at all, with the
  reasoning "nothing may silently decide, so U."
- **run-3** (`glm-5.2`): affirmatively rules the case out — "the RDR pushes all
  parse/lint failure to RDR 0003's load-time … so at Evaluate time the guard is
  well-formed and an `error` return has no defined trigger."

This is the sharpest 3-way: one run left it open, one invented a normative
verdict for it, one closed it as unreachable. `fable-5`'s invented
`GuardUnevaluable` fallback is the dangerous rendering — it makes an evaluator
bug indistinguishable from missing artifact state, which is precisely the
conflation the RDR exists to prevent.

**RDR passage**: The atom-vocabulary scope clause requires the mapping be
"total and lossless" and says "an evaluator that cannot recover which tag keys
an atom references cannot implement the domain rule at all" — but states no
runtime disposition for a mapping failure. A7/A1 route undeclared keys and
unknown operators to load-time rejection; neither covers a well-loaded guard
the evaluator cannot parse.

**Lands in**: Technical Design → Normative Contracts (atom-vocabulary scope
clause). Either state that mapping failure is unreachable by construction (and
why), or fix its verdict — and if the verdict is `GuardUnevaluable`, say
explicitly that it is not the same signal as absent state.

### D-4 — `Evaluate`'s return shape (3-way in confidence, not in value)

**Contract**: bare `GuardResult` vs `(GuardResult, error)`.

- **run-1** (`opus-5`): reconstructs bare `GuardResult`, GUESS inline on the
  signature and again in the ranked guess inventory (#4).
- **run-2** (`fable-5`): bare `GuardResult`, GUESS, with an argued lean ("the
  verdict is the only observable") but explicitly "never rules out
  `(GuardResult, error)`."
- **run-3** (`glm-5.2`): bare `GuardResult`, GUESS, and additionally asserts in
  §1.6 as settled fact: "There is no Go `error` channel for guard evaluation."

All three land on the same value; `glm-5.2` alone hardens it into a stated
contract. The RDR quotes only the parameter list.

**RDR passage**: The RDR quotes `Evaluate(guard string, view TagSet)` (via
`fixtureGuards`' `Evaluate(guard string, _ resolve.TagSet)`) and never the
return. Related but not decisive: "the verdict is the only observable" appears
in Research Findings about the *kernel's* channels, not about this signature.

**Lands in**: Technical Design → Normative Contracts (atom-vocabulary scope
clause, alongside D-3 — same seam, same silence).

### D-5 — Whether `GuardEvaluator` is an interface (2-vs-1 — `fable-5` dissents)

**Contract**: the Go shape of the seam.

- **run-1** (`opus-5`) and **run-3** (`glm-5.2`): `type GuardEvaluator
  interface { Evaluate(...) }`, asserted without a GUESS marker.
- **run-2** (`fable-5`): marks it GUESS — "a function type `type GuardEvaluator
  func(string, TagSet) GuardResult` would also satisfy every quoted sentence."

`fable-5` is right that the RDR never states it; the other two read the
interface in from the fact that a stub implements it.

**RDR passage**: RDR is silent on the declaration form. It names
`resolve.go::GuardEvaluator` as "the seam" and the Phase 2 harness constraint
"an exported function taking the seam" — which reads either way.

**Lands in**: Technical Design → Normative Contracts (atom-vocabulary scope
clause). Low-severity relative to D-1..D-4, but it is the parameter type of
the Phase 2 harness signature the RDR *does* fix normatively.

### D-6 — `TagSet` struct vs interface (2-vs-1 — `opus-5` dissents)

**Contract**: whether the evaluation view is a concrete type or an interface
the evaluator programs against.

- **run-1** (`opus-5`): renders `type TagSet interface { Lookup(...) }` and
  flags the tension itself as GUESS — "interface vs struct is not stated; the
  RDR calls `TagSet.Lookup` 'the only exported accessor', and elsewhere writes
  `view.tags` / `view.has(key, ProvenanceOwned)` as unexported internals, which
  reads more like a struct with methods."
- **run-2** (`fable-5`) and **run-3** (`glm-5.2`): both render `type TagSet
  struct { tags map[...]taggedValue }` with unexported `has`/`matches`/`Len`,
  no GUESS on the struct-vs-interface question.

Consequential: if `TagSet` is a concrete struct with unexported fields, an
out-of-package evaluator (RDR 0003's, in a separate package) can construct a
view only through kernel-provided constructors — which bears directly on
whether the Phase 2 harness can build the presence/absence views scenarios
4/6/7/8 require.

**RDR passage**: A11's Evidence enumerates "exactly four readers" with
`TagSet.Lookup` "the only *exported* accessor — the one an evaluator outside
the package can use at all," and the provenance-blind clause cites
`view.has(key, ProvenanceOwned)` as kernel-internal. The RDR never states the
type's declaration form, nor how a test constructs a `TagSet` with chosen
contents.

**Lands in**: Technical Design → Normative Contracts (the provenance-blind
presence clause) plus Implementation Plan → Phase 2 (view construction for the
harness).

### D-7 — `no_match` / `ambiguous_match` determination point relative to the refusal checks (2-vs-1 — `opus-5` dissents)

**Contract**: where zero/multiple-candidate counting sits in `gate`'s order.

- **run-1** (`opus-5`): puts the counting in the *caller* (`resolve`), after
  `gate` returns survivors — `if len(candidates) == 1 … == 0 … else
  ambiguous_match` — and GUESSes both the escape-exhaustion mapping and the
  `>1 → ambiguous_match` mapping ("inferred from the closed escape classes").
- **run-2** (`fable-5`): puts the counting inside `gate`, after the two refusal
  checks, and marks the placement GUESS explicitly: "relative order of
  `no_match`/`ambiguous_match` determination vs the two refusal checks shown —
  the RDR fixes only: prune first (D8), owned-before-undecidable among
  survivors, and the undecidable veto discarding `selected`."
- **run-3** (`glm-5.2`): does not model it at all — `gate` ends at `return
  selected`; no zero/multiple mapping appears anywhere in the reconstruction.

Three different structural placements for the same closed taxonomy.

**RDR passage**: The RDR fixes only three ordering facts — D8 prune-first,
"absent owned state is reported BEFORE an undecidable guard", and the
aggregation veto. RDR 0001's selection rule is quoted in Technical Environment
("the only successful selection is exactly one matching edge after guard
evaluation. Zero, multiple, unavailable, or unevaluable candidates are
refusals…") but this RDR never restates where in `gate` the zero/multiple
branch sits or which function owns it.

**Lands in**: Technical Design → Normative Contracts (aggregation block).
Arguably RDR 0001's to own — in which case this RDR should cite rather than
leave the reader to infer.

### D-8 — The fifth refusal kind (3-way, all GUESS)

**Contract**: the closed five-kind taxonomy the RDR repeatedly leans on.

- **run-1** (`opus-5`): names four, "the fifth kind is never named in this RDR"
  (ranked #5 in its guess inventory).
- **run-2** (`fable-5`): names four, then *invents* a candidate — "Guess:
  something like an input/table-validity refusal (`invalid_input` or similar),
  never reachable from the guard seam" — and additionally guesses the Go
  constant spellings `KindNoMatch` / `KindAmbiguousMatch` /
  `KindOwnedStateUnavailable`.
- **run-3** (`glm-5.2`): names four and explicitly declines — "The fifth kind
  is named only by test, not by this RDR; I do not reconstruct it."

`fable-5`'s invented fifth kind is the notable divergence: nothing in the RDR
supports it, and it is the sort of confabulation an implementer could carry
forward.

**RDR passage**: The taxonomy is called closed in three places ("the closed
refusal taxonomy", "the five-kind set is pinned closed by
`resolve_test.go::TestReq7_RefusalKindSetIsExactlyTheFiveNamedKinds`", "reopens
RDR 0001's closed five-kind taxonomy") and the fifth member is never named.

**Lands in**: Technical Design → Normative Contracts, or Technical Environment
→ RDR 0001 summary. Cheapest fix in the set: name all five once where the
taxonomy is first called closed.

### D-9 — The Phase 2 view-reading evaluator's exported shape (2-vs-1 — `fable-5` is the only run to model it)

**Contract**: whether the Phase 2 reference evaluator is exported, and how it
sidesteps A10.

- **run-2** (`fable-5`): gives it its own section (§1.5) and two GUESSes —
  "unexported reference evaluator beside the harness, since the *normative*
  evaluator is RDR 0003's future build and exporting a second one would invite
  production use of a stopgap"; and that it must define "a private literal
  guard encoding for vector purposes only … without claiming that encoding for
  RDR 0003."
- **run-1** (`opus-5`): reconstructs only the harness (`RunGuardVectors`) as
  new public surface; the view-reading evaluator appears nowhere in its API
  section.
- **run-3** (`glm-5.2`): does not model Phase 2's harness or evaluator at all —
  no `RunGuardVectors`, no vector record.

`fable-5` alone surfaced the real bootstrapping problem: the Phase 2 evaluator
must encode guard structure while A10 is open, so it necessarily invents an
encoding the RDR does not sanction, and nothing says that encoding is
non-normative.

**RDR passage**: Phase 2 requires "A view-reading evaluator … the domain rule's
first executable expression" and fixes the *harness* shape
("`resolveconform.RunGuardVectors(t, eval)`", non-`_test` package, "caller
supplies the evaluator"), but says nothing about the evaluator's visibility or
about how it represents guards while A10 is Pending. Testing Strategy states
the blocker ("no Go representation of any of that exists") without resolving
the tension that Phase 2 must nonetheless build one.

**Lands in**: Implementation Plan → Phase 2. State that the Phase 2 evaluator's
guard encoding is vector-local and non-normative, and fix its visibility.

### D-10 — Whether `exists` takes a literal (2-vs-1 — `glm-5.2` dissents)

**Contract**: `exists` atom arity.

- **run-1** (`opus-5`) and **run-2** (`fable-5`): treat `exists` as carrying a
  polarity/literal. run-1's atom struct is uniform (`Op`, `TagKey`, `Literal`)
  and `atomVerdict` returns `boolVerdict(present)`; run-2's `evalAtom` is
  explicit — `return bool3(present == atom.wantPresent)`, i.e. an authored
  expected-presence value.
- **run-3** (`glm-5.2`): "I GUESS `exists` is unary over a key with no literal,
  since its 'whole job is deciding presence'" — `return ok ? TRUE : FALSE`.

This changes what a negative existence test *is*: under run-2's reading
`exists = false` is expressible as one atom; under run-3's it is expressible
only via `unless … exists`, which is the A5 blessed-pattern route. The two
readings disagree about whether A5's two-row pattern is the only way to author
absence.

**RDR passage**: The RDR uses both forms without reconciling them. A11 and the
provenance-blind clause quote the fixture as `exists = true` (a literal);
Testing Strategy scenario 6 says "a positive existence atom … and a negative
existence atom (`unless … exists`)" (block placement carries polarity); A5
describes "one row guarded on absence via `unless` existence". Whether polarity
lives in the literal or in the block is never fixed.

**Lands in**: Technical Design → Normative Contracts (the existence-operator
clause). Materially affects A5/A12, since the two-row pattern's necessity
depends on it.

## 2. GUESS clusters

Contracts two or more runs marked GUESS. These are the candidate RDR rewrites,
ordered by cluster size and by how much normative text hangs on them.

### G-1 — The `Row.Guard` ↔ atom-structure mapping (3/3 GUESS; all three rank it first)

All three runs name A10 as the hinge and mark the recovery step GUESS
(`recoverStructure` / `recoverAtoms` / `recoverAtoms`). run-1 ranks it #1 in
its guess inventory ("the entire normative block quantifies over a structure
the RDR explicitly does not define"); run-2 calls it "the hinge"; run-3 calls
it "the blocking prerequisite for every atom-level clause". Downstream: D-3,
D-4, D-9, D-10 all derive from this silence.

**Status in RDR**: A10 Pending, and Prerequisites already require it verified
before lock. Not a hidden gap — but the three runs confirm that *every*
atom-level normative clause is unrenderable until it closes, which is stronger
than "a Pending assumption."

**Lands in**: Technical Design → Normative Contracts (atom-vocabulary scope
clause) once A10 resolves.

### G-2 — `TagValue` / typed value representation (3/3 GUESS)

run-1 ranks it #2 and calls it "load-bearing for three of the five operators
(`contains` absence semantics is *normative*)"; run-2 GUESSes `Value` inside
`taggedValue`; run-3 "reconstruct a sum type; the real type is RDR 0003's."

The RDR's own `contains` clause is normative — "an absent set-valued tag MUST
be treated as unevaluable under set containment, NOT as the empty set" — and
quantifies over a scalar/set distinction the RDR never models. Bounded integer
comparison implies a third typing.

**Lands in**: Technical Design → Normative Contracts (the partial-operator
clause) — at minimum, state that scalar/set/bounded-int typing is RDR 0003's
and that this RDR's clause binds regardless of representation.

### G-3 — `Refusal.Rows` element type and `Row` field types (3/3 GUESS)

See D-2. All three GUESS the element type and diverge on it; run-1 additionally
GUESSes `Row.RuleID`, `Row.SourceLocator`, `Row.RequiresOwned`, and `Row.Writes`
types; run-2 GUESSes `Row`'s match-pattern and escape-marking fields entirely
("the RDR references 'matched the outcome and the tag pattern' and 'modeled
escape edge' without showing the fields"); run-3 marks the selection fields
GUESS.

**Lands in**: Technical Design → Normative Contracts (refusal-payload block,
and the `Row.RequiresOwned` block which is the one field this RDR *does* fix).

### G-4 — `Evaluate`'s return shape (3/3 GUESS)

See D-4. Unanimous value, unanimous GUESS marking, one run (`glm-5.2`)
hardening it into prose anyway.

**Lands in**: Technical Design → Normative Contracts (atom-vocabulary scope
clause).

### G-5 — The fifth refusal kind (3/3 GUESS)

See D-8. Two decline to name it, one invents `invalid_input`.

**Lands in**: Technical Environment (RDR 0001 summary) or Normative Contracts.

### G-6 — Conformance harness: package name, reporter interface, vector encoding (2/3 GUESS)

run-1 (#6 in its inventory) and run-2 both GUESS the package name, the
`TestingT` method set, and whether vectors are Go literals or `testdata/`
golden files. run-3 does not model the harness at all — itself a signal that
Phase 2's surface is thin enough to be skipped in a reconstruction.

The RDR marks the encoding *deliberately deferred* ("the exact package name and
vector encoding are Phase 2's to fix") while three adjacent constraints are
normative (importability, exported-function shape, caller-supplies-evaluator).
The deferral is legitimate; both runs read it as such. Recorded here because it
is the largest deliberate silence, not because it needs a rewrite.

**Lands in**: Implementation Plan → Phase 2 (no change required unless the
reporter interface is to be pinned).

### G-7 — `ProvenanceRecognized` constant name (2/3 GUESS)

run-1 GUESSes it explicitly (#10) — "prose only ('recognized'), constant name
never given"; run-2 GUESSes the whole `Provenance` type including the third
constant. run-3 lists the three provenances in a comment without marking it.

The provenance-blind clause is normative and quantifies over "owned, observed,
or recognized"; two of the three constants are quoted from source
(`ProvenanceOwned`, `ProvenanceObserved`), the third is not.

**Lands in**: Technical Design → Normative Contracts (provenance-blind clause).

### G-8 — The strong-Kleene negation table for `unless` (2/3 GUESS)

run-1 ranks it #7 — "stated by reference, never written out." run-2 GUESSes the
`bool3`/`min` framing ("the RDR fixes the truth tables (A2) but not that the
implementation literalizes them as a min-lattice"). run-3 writes `¬T=F, ¬F=T,
¬U=U` as settled, correctly quoting A2's Evidence.

A2 does state the negation clause (`¬T=F, ¬F=T, ¬U=U`) but only in its Evidence
block; the ```normative``` combination clause references strong Kleene without
carrying the tables. run-3 found the negation table in A2; the other two did
not, which is the real finding — the tables are load-bearing but live outside
the normative block.

**Lands in**: Technical Design → Normative Contracts (combination clause) —
lift the conjunction and negation tables from A2's Evidence into the normative
block, or cite A2 explicitly from it.

## 3. Agreement

All three runs rendered these identically, confirming the RDR is determinate
there: the three-valued `GuardResult` names and their role (`GuardFalse`
prunes, `GuardUnevaluable` never folded); the per-atom domain rule (four
value-comparing operators partial, absent key ⇒ unevaluable; `exists` the sole
total operator; `contains` over an absent set-valued tag unevaluable and *not*
the empty set); provenance-blind presence via `TagSet.Lookup`'s `ok`, expressly
distinguished from `missingOwned`'s `has(key, ProvenanceOwned)`; the
strong-Kleene combination shape `all ∧ ¬(unless_conj)` with `F ∧ U = F` as
witnessed falsity making D8 sound; empty `all` ⇒ TRUE and empty/omitted
`unless` ⇒ absent-not-vacuously-true (three different encodings, one identical
net semantics — an encoding-taste difference, not a finding); `gate`'s ordering
(D8 prune-first → `missingOwned` among survivors → undecidable loop) and the
resolution-level veto discarding accumulated GuardTrue rows; D5 escape-set
scoping with an unevaluable escape row replacing the candidate-set refusal;
`guard_unevaluable` non-escapable by RDR 0002's closure; the `Row.RequiresOwned`
narrowing to post-guard write-dependency keys with guard-read keys still legal
there; `Refusal.Guard` as single-valued diagnostic convenience and `Refusal.Rows`
as the discriminator; no missing-tag payload on `GuardResult`; and A6b's
read-failed-vs-absent conflation carried open. Helper decomposition and naming
(`atomVerdict`/`evalAtom`, `combineKleene`/`combine`/`k3And`) varied freely and
is not reported as a finding.
