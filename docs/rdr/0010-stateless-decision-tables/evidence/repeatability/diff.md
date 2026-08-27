Model: claude-opus-5[1m]

# Repeatability DIFF — RDR 0010 (stateless decision tables)

Three runs, three distinct models: run-1 `claude-sonnet-5`, run-2
`claude-fable-5`, run-3 `claude-haiku-4-5-20251001`. Each read C1–C5, MVV,
S1–S6 and widened into `§technical-design` lead-in, `§authority`,
`§load-bearing-decisions`, `§illustrative-code`.

All three reconstructions agree on the RDR's core: the `class` key and its
two admitted values, the one-directional agreement check and its position
at/after `loadTags`, `EmitValue {Key, Value string}` as a new row type,
`Emit` carried through `expand`, the kernel row untouched, `emit` appended
last in the dump vocabulary, `emit` immediately after `rule` in the payload,
the augmenting seed predicate, and the zero-dimension
`graph-unprovable-coverage` / `no-participating-dimension` arm emitted inside
`len(dims) == 0` ahead of `emitCoverageArms`. The disagreements below are
narrower than that agreement, and most of them are single-run.

---

## 1. Disagreements

### D1 — `0010:C4`: where gate evaluation sits relative to the emit join, and whether a gate denial is reachable over a decision table (2-way split; **model boundary — run-3 alone**)

**Contract:** C4 fixes the payload's `emit` field, its position, its `{}`
never-`null` shape, and the text rendering. It says nothing about gates.
C2 says "Rule ids, the match-block obligation, shared contexts, guards,
**gates**, escape rules, and outcome binding are unchanged in both classes"
— gates are explicitly retained for the decision-table class, but the RDR
never states where gate evaluation falls in the `flow resolve` sequence
relative to the post-selection `rowByID` emit join, nor whether the success
payload the RDR is amending carries a gate-results field that `emit` must be
positioned against.

- **run-1**: no gate step anywhere in the pseudo-code. Sequence is
  `kernel.Resolve → rowByID → payload`. No gate field in the payload.
- **run-2**: no gate step. Same sequence. No gate field. It does mention
  `guard/gate` blocks only in the normalization contrast (C3's
  "EVEN AN EMPTY ONE" arm).
- **run-3**: inserts an explicit step 14, `runGates(selected_rule, gates) →
  IF gate denied: RETURN flow-gate-denied`, **between** kernel selection
  (step 13) and the `rowByID` emit join (step 15), and adds a `gates: gate
  results` field to the `resolvePayload` literal.

**Why they diverged:** C2 keeps gates live for the class but the RDR's own
data-flow line (`§technical-design`: "TOML → loader → lint → `flow resolve`
→ kernel exact-one → payload `{rule, emit, next: {}, writes: {}, …}`") and
C4's pseudo-envelope both elide gates entirely. An implementer therefore has
no fixed answer to: does a gate-denied selection over a decision table
produce a refusal (so `emit` is never computed), or a success payload
carrying `emit` alongside a denial record? Both are consistent with the text.
This is a real ordering silence, not naming taste: it decides whether the
emit join runs before or after a step that can abort the success path.

**Aggravating factor:** C4 fixes `emit`'s position as "immediately after
`Rule`" in a **thirteen-field** declaration. run-3's payload literal carries
a `gates` field that runs 1 and 2 do not; run-1 explicitly flags the
thirteen-field arithmetic as a GUESS ("13 → 14 is an inference, not a
directly stated total"). If `gates` is or is not among the thirteen, the
field count assertion S5 and the Testing Strategy license differ.

**Lands on:** `0010:C4` — add a clause stating whether the gate step precedes
the emit join and whether a gate-denied decision-table selection is a refusal
(no payload) or a payload; and either enumerate the thirteen fields or drop
the count in favour of the positional clause alone.

---

### D2 — `0010:C4`: whether `escaped` is part of the payload this RDR pins, and whether the escape row's `emit` is joined on the same path (3-way split)

**Contract:** C4 enumerates the fields it constrains: `emit`, then "`next`,
`writes`, `clear`, `owned`, and `readers` keep their `0005:C1` shapes and are
empty over a decision-table model." `escaped` is named nowhere in C4. A9
(non-normative) states "the rescued plan reports `escaped = true`", and
`§authority` routes the answer as "joined by `Plan.RuleID` after selection".

- **run-1**: payload literal is `{Rule, Emit, Next: {}, Writes: {}, Clear:
  [], Owned: {}, Readers: []}` + "remaining 0005:C1 fields". **`escaped` is
  absent from the payload entirely**, and the escape path is described only
  inside the kernel comment ("escape 'otherwise' row rescue, or no_match")
  — the run never states that an escape-rescued selection joins its own
  `emit`.
- **run-2**: carries `escaped` as a separate wire clause ("Escaped selection:
  rescued plan reports `escaped = true` (A9)") and states explicitly that an
  escape row "may carry its own emit", listing it under the error-mode
  section beside `flow-ambiguous-match` / `no_match`.
- **run-3**: carries `escaped: false` as a **literal field inside the payload
  JSON**, positioned last, and models the escape rescue as an explicit branch
  setting `escaped=true` before the join (step 13).

**Why they diverged:** C3 says the join is by rule id and is "sound because
emit is authored per *rule*", and C4 fixes `emit`'s presence for "the
selected row". Whether "the selected row" on a rescued plan is the escape row
— i.e. whether `Plan.RuleID` on a rescue names the escape rule so `rowByID`
returns the escape row's `Emit` — is stated only in A9's prose, never in a
normative clause. run-1's reconstruction is silent on it precisely because
the C-spans are. S5's fixture requires a selectable row authoring **no**
emit block to exercise the `{}`-never-`null` arm, but no scenario pins the
escape-row-with-emit resolve.

**Lands on:** `0010:C4` — state that a rescued plan's `emit` is the escape
row's own block joined on the same `Plan.RuleID` path, and that `escaped`
sits in the payload's fixed order relative to `emit`.

---

### D3 — `0010:C1` / `§authority`: the class's storage shape on `table.Model` — field vs. accessor, and what a hand-constructed Model reads (3-way split)

**Contract:** C1's normative span says only "The class is model data and is
carried on the loaded model; nothing downstream infers it from the owned
set." The accessor and its zero value appear **only** in the non-normative
`§technical-design` lead-in ("`table.Model` gains the class as a field with
an accessor whose **zero value is `state-machine`**") and in `§authority`'s
"via one accessor whose zero value is `state-machine`".

- **run-1**: `func (m *Model) Class() string` — accessor method returning a
  bare string; marks the method identifier a GUESS.
- **run-2**: `func (m *Model) Class() ModelClass` with `type ModelClass
  string` and named constants `ClassStateMachine` / `ClassDecisionTable`;
  marks the Go identifiers a GUESS but the string values normative.
- **run-3**: **exposes `Class string` as a bare public struct field** on
  `table.Model` (`type Model struct { ID string; Class string; Rows []Row }`)
  and separately asserts "Accessor with zero value `state-machine` (RDR
  silent: assumed…, not explicitly in contracts)". run-3's data model and its
  Public API section contradict each other: the struct literal has no
  accessor.

**Why they diverged:** the zero-value clause is load-bearing (it is what
keeps a hand-constructed `table.Model` a state machine), and `§authority`
makes the accessor "the single source for what class is this" with four named
readers and an explicit warning that "nothing in the build refuses a fifth
site that re-derives the class from `len(owned) == 0`". But the clause that
carries that obligation lives outside the normative spans. run-3, reading
closest to the contracts, produced an exported field — which defeats the
single-writer guarantee `§authority` describes, since an exported field is
directly assignable and bypasses any accessor. All three ended up with a
different, non-substitutable Go surface.

**Admissibility note:** the identifier (`Class()` vs `ModelClass`) is naming
taste and is not the finding. The finding is that field-vs-accessor decides
whether the "one writer, four readers" ownership `§authority` asserts is
mechanically expressible at all, and the zero-value obligation is stated
nowhere a contract reader is required to look.

**Lands on:** `0010:C1` — promote the zero-value/accessor clause from the
Technical Design lead-in into C1's normative span, stating that the class is
reachable only through an accessor whose zero value is `state-machine`.

---

### D4 — `0010:C5` / `0010:C1`: which loader step owns the class-agreement check relative to `loadAccessors` and `checkAccessorBindings` (2-way split; **model boundary — run-3 alone**)

**Contract:** C1 fixes the check as "a step at or after `loadTags`" and fixes
the reason ("`l.model.Tags` is empty while the header loads", "so an
undeclared-tag refusal precedes a class-disagreement refusal"). C1 fixes a
**lower** bound only. C2 separately establishes that `checkAccessorBindings`
is "`run`'s last step" and that it is where the `[initial]`-over-observed
refusal lands.

- **run-1**: places the check "at/after `loadTags`" with no upper bound, and
  the pseudo-code runs it immediately after `loadTags`, before
  `normalizeRule`.
- **run-2**: same — `checkClassAgreement(m)` immediately after `loadTags`,
  before the rule loop, with `loadInitial/loadAccessors/checkAccessorBindings`
  running **after** normalization.
- **run-3**: asserts an **upper** bound the RDR does not state — "Called at
  or after `loadTags`, **before `checkAccessorBindings`** (run's last step),
  so undeclared-tag refusals precede class disagreement" — and orders
  `loadAccessors()` (step 7) *before* `normalizeRules()` (step 8), inverting
  runs 1 and 2's sequence.

**Why they diverged:** this decides observable diagnostics. A decision-table
model that both declares an owned tag *and* declares `[initial]` over that
tag refuses with `malformed model declaration` under run-3's ordering but
could refuse with the accessor-binding/writer-arity diagnostic C2 describes
under a later placement. C2 goes out of its way to note that the arity
diagnostic "names writers rather than `[initial]`, which is survivable but
not self-explanatory" — so the RDR is aware the two refusals compete, but
fixes only the `loadTags` floor, never the ceiling. S1 tests the agreement
check in isolation; no scenario pins the two-refusal collision.

**Lands on:** `0010:C1` — fix the upper bound of the agreement check's
position (before or after `loadAccessors` / `checkAccessorBindings`), so the
refusal a doubly-malformed decision table takes is determinate.

---

### D5 — `0010:C3`: what `emit` normalization does with the pair sequence when the same key would sort equal, and whether `expand`'s carried block is shared or copied (2-way split)

**Contract:** C3 fixes "normalization sorts by key and asserts no dedup pass;
'duplicate-free' is a property of the source grammar, not an obligation on
the loader", and "`Row.Emit` MUST be carried through
`internal/table/normalize.go::expand` for every expanded row, which builds
each row from the seed literal rather than copying `base`."

- **run-1**: `row.Emit = sortByKey(decodeEmitBlock(rule))`, and states
  `expand()` "mints one Row per match-block member, carrying Emit on every
  expanded row (not just the first)". Does not say whether the carried value
  is the same backing slice or a copy.
- **run-2**: "`expand` must **copy** `Row.Emit` into every row it mints"
  — explicitly a copy. Adds a `normalizeEmit(map[string]string) []EmitValue`
  helper (marked GUESS).
- **run-3**: "`expand(rule) → multiple rows, each carrying **same** Emit"
  — explicitly shared, and its data model has all expanded rows referencing
  one block.

**Why they diverged:** C3 fixes only that the block is carried and that every
expanded row "carries the same block". `§authority` reinforces "every row a
rule expands to carries the same block", and S3 asserts "**every** row
expanded from the `in`-atom rule carries the authored block
**byte-for-byte**". Byte-for-byte equality is satisfied by both a shared
slice and a copy. Since `Row` is otherwise value-copied through the pipeline
and the dump and the JSON join both read `Emit`, aliasing is unobservable
today — but nothing in the RDR forbids a later mutation site, and the RDR's
own reason for `EmitValue` (avoiding `TagValue`'s `cloneTagValues` machinery)
is precisely a statement about copy semantics that it then leaves open for
`Emit` itself.

**Admissibility note:** this is borderline. It is admitted because C3
explicitly reasons about clone machinery for the neighbouring type and then
declines to state the aliasing rule for the new one — a field-ownership
silence, not a helper-decomposition preference. It is the weakest finding
here and is ordered last among the disagreements.

**Lands on:** `0010:C3` — state whether `expand` shares or clones the `Emit`
sequence across the rows one rule mints.

---

## 2. GUESS clusters

Contracts two or more runs marked GUESS. These are the candidate rewrites.

### G1 — The class's Go surface on `table.Model` (all three runs; `0010:C1`)

- run-1: GUESS on `table.Model.Class()` method identifier and receiver —
  "RDR specifies only 'an accessor whose zero value is state-machine' on
  `table.Model` — the identifier `Class()` is inferred".
- run-2: GUESS on the class type and constants — "the string values are
  normative; the Go identifiers are mine".
- run-3: GUESS #2 — "Assumed `Model.Class` accessor exists and has zero-value
  semantics (**not explicitly stated in contracts**, but required by C1's
  narrative and standard Go patterns)".

Two of three treated the identifier as the guess; run-3 flagged the harder
thing — that the **zero-value semantics are absent from the contracts**. That
is correct: they are in the `§technical-design` lead-in only. This is the
same silence as D3 and is the single strongest rewrite candidate in this
pass, because it is the clause the whole `§authority` ownership table rests
on.

**Rewrite:** move the accessor + `state-machine`-zero-value clause into
`0010:C1`'s normative span.

### G2 — The class-agreement check's identity as a named loader step (all three runs; `0010:C1`)

- run-1: names it as an inline concern of `load.go::run`, no separate
  function — no GUESS marker, but it declines to name one.
- run-2: GUESS `checkClassAgreement` "(or an arm appended to an existing
  loader step); the RDR fixes its *position*, not its name".
- run-3: GUESS #3 `classOwnershipAgreement` "(Conjectured)" — "RDR specifies
  behavior, not implementation".

The **name** is not a finding. But the cluster is real in one respect all
three circled: the RDR fixes a floor position for a step whose existence as a
distinct step is itself unstated — which is what let run-3 (D4) invent an
upper bound and runs 1/2 not to. Rewriting C1's position clause with both
bounds (D4) dissolves this cluster.

### G3 — The zero-dimension coverage arm as a distinct helper (2 runs; `0010:C5`)

- run-2: implicitly inline — "`coverage.go::checkCoverage` gains the
  class-keyed zero-dimension arm".
- run-3: GUESS #5 `emitUnprovableDimensionOrNone` "(Conjectured)".
- run-1: inline in `checkCoverage`, no helper.

**Not a finding.** C5 fixes the emission site, the branch, the ordering
against `emitCoverageArms`, the `reason` value, and the
`graph-coverage-closed-by-escape` suppression. Everything an implementer needs
is fixed; only the decomposition is free, and decomposition is taste. Dropped.

### G4 — The emit-join helper in `flow_resolve.go` (2 runs; `0010:C3`)

- run-3: GUESS #4 `selectRowWithEmit` "(Conjectured)".
- run-2: GUESS on "the exact name of the dump cell renderer for emit …
  inline in the dump column code rather than a named helper — the RDR fixes
  rendering, not a function".
- run-1: uses `rowByID(model, plan.RuleID)` directly, no helper, no GUESS.

**Not a finding.** C3 names `internal/cli/flow_resolve.go::rowByID`
explicitly and states it "is correct as it stands". A helper wrapping it is
decomposition. Dropped.

---

## 3. False GUESSes (reported here so they are not mistaken for silences)

These are cases where a run marked GUESS but the RDR **does** fix the clause.
Each is dropped, not reported as a finding.

- **run-1, `malformed dump declaration` on an explicit `[dump]` omitting
  `emit`** — marked "GUESS: inferred from S3's expected-refusal language …
  treated here as normative behavior since it is asserted, not merely
  illustrative". C3 fixes it directly: "A `[dump]` column list MUST name it".
  A4 grounds the refusal in `loadDump`'s omitted arm with an exact census
  (103 files). run-1 reached the right answer by the wrong route — the
  contract was truncated in its reading, not silent. **Widened and dropped.**
- **run-1, the exact TOML type of `class` beyond "string"** — C1 says "a
  string whose only admitted values are `"state-machine"` and
  `"decision-table"`". There is nothing further to fix. **Dropped.**
- **run-3, `graph-unprovable-coverage` existing reasons** — run-3 enumerates
  `dimension-not-finite`, `tag-not-single-valued`, `row-can-refuse` without a
  GUESS marker; C5 and A14 fix all three plus the append. Correct, no
  divergence. Noted only because run-3 is the only run that enumerated them.
- **All three, `graph-coverage-closed-by-escape` precedence** — run-2 alone
  states the suppression ("where both would apply,
  `graph-coverage-closed-by-escape` MUST NOT be reported"); runs 1 and 3
  state only the ordering against `emitCoverageArms`. C5 fixes the
  suppression explicitly. Runs 1 and 3 under-read a fixed clause; this is a
  reading miss, **not an RDR silence**. **Dropped.**
- **run-3, `checkDanglingEdge` in `analysis.go` vs run-1/run-2's unqualified
  reference** — file placement, not a contract. **Dropped.**
- **Dump cell separator** — run-1 renders `[next=propose]` (single pair, so
  the separator is unexercised); run-2 `[a=x; b=y]`; run-3
  `[next=propose; key2=value2]`. C3 fixes `bracketed like writes ([a=x; b=y])`.
  No divergence on the contract; run-1's single-pair example simply cannot
  show it. **Dropped.**

---

## 4. Agreement

All three runs rendered these identically, confirming the RDR is determinate
there: the `class` key and its two admitted values with absent reading as
`state-machine`; the one-directional agreement check refusing only
decision-table-with-owned-tags, with detail token `owned=<n>`, and running
at/after `loadTags` never in `loadModelHeader`; the `malformed rule shape`
no-write-block arm conditioned on class; `EmitValue {Key, Value string}` as a
new row-level type explicitly not `TagValue`; `Row.Emit` key-sorted, with
absent and present-but-empty both normalizing to the empty sequence without
refusing; `Emit` carried through `expand` to every expanded row; non-string
and nested emit values refusing `malformed TOML` on the decoder type arm
asserted on category only; emit keys as non-tag keys refusing `unknown tag`
when matched; `emit` appended last after `escape` in `dumpColumns`; `emit`
immediately after `rule` in `resolvePayload` as a string→string object,
`{}` never `null`, never omitted; text mode through 0005's generic renderer
as `emit.<key>: <value>` / `emit: (none)` with no per-verb special case;
`--outcome` required in both classes and a stray `--artifact` ignored;
`flow next` unchanged and carrying no `emit` (BR4); the kernel row untouched
with the join by `Plan.RuleID` via first-match `rowByID`; the reach seed
predicate as `class == decision-table || len(Initial) > 0` augmenting not
replacing; exactly two of the four `len(Initial) == 0` sites becoming
class-aware; `graph-unprovable-coverage` with `reason =
"no-participating-dimension"` emitted inside `len(dims) == 0` ahead of
`emitCoverageArms`; the `graph-dangling-edge` split with the missing-root arm
silenced and the terminal-predicate arm live; and `Fingerprint` excluding
`emit` (A7).
