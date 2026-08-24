---
authors: Chris K Wensel <cwensel@retrofit.sh>
state: open
cluster: 0002, 0003, 0004, 0005, 0006, 0007, 0008, 0009
labels: intrastate, internal-resolve, refusals, guard-evaluation, rdr-cluster
---

# JDR 0001 What the Resolver Can Trust, and What the User Sees When It Can't

*Go identifiers and CLI transcripts below are **non-normative**; the RDRs own
exact contracts. Source: the `0002-0009` cluster gate —
`docs/rdr/cluster-reconcile/0002-0009/`.*

## Problem statement

Eight RDRs jointly decide one thing: **when the resolver cannot decide, what
happens.** 0002 owns the table format, 0003 the guard vocabulary, 0004 accessor
execution, 0005 the CLI surface, 0006 graph lint. Three later RDRs — 0007
(guard evaluation domain), 0008 (recognized tag key), 0009 (escape-row shape) —
each landed on the `area:internal-resolve` seam between them, two months after
the first five locked.

The structural cause of every finding is one grep:

```
grep -cE 'RDR 0007|RDR 0008|RDR 0009' 0002…md 0003…md 0004…md 0005…md 0006…md
→ 0, 0, 0, 0, 0
```

No older member references any newer one. Meanwhile 0007 names
`Obligation destination (NAMED): RDR 0003's implement stage` three times and
routes A6b to 0004's implement stage. An obligation filed on a document that
never received it binds nobody.

Nothing in `internal/resolve` has a production consumer — one non-test file, no
`flow` verb — so every decision here is settled against test-covered code rather
than a shipped surface. That is why the decisions below prefer the clean shape
over the compatible one.

## Principles

Steering criteria for every decision here. Each is derived from text the members
already lock.

1. **Missing artifact state is never masked** — 0007's reason for existing: the
   domain choice "decides whether missing artifact state can be masked behind an
   escapable refusal class."
2. **The kernel refuses rather than guesses** — 0004: "Indeterminate MUST be a
   refusal-class result, not a false allow and not a false deny."
3. **Every user-visible failure is structured** — 0005 routes all failure
   through `respond.Fail(cmd, *clierr.CLIError)`; the consumer parses JSON.
4. **A refusal names what was missing and who can fix it** — 0007's user
   outcome, plus 0005's exit-code grouping, which encodes remedy.
5. **A green lint means resolution succeeds** — otherwise 0006's blocking
   authority promises nothing.
6. **One home per contract; cite, never restate** — the engine doctrine, and
   why this registry exists.
7. **Pre-release, prefer the clean shape** — no shims for callers that do not
   exist. *(Rests on a project convention, not repo text — confirm before
   leaning on it.)*
8. **The layout makes illegal states unrepresentable** — a shape the strict
   decoder can refuse is never left to a validator (0002's strict-decoding
   clause; resolved by §D7).

---

## D1 — How does the guard seam carry a predicate?

`resolve.Row.Guard` is a `string` (`resolve.go:185`) — an opaque predicate. The
evaluator must reconstruct structure from it and can fail doing so, and 0007
answers that failure with a mandated panic the kernel "MUST NOT recover." 0005
requires every user-visible failure to be a structured envelope. Neither RDR
mentions the other.

The same opacity causes a second defect: 0003's predicate identity is
*per-atom*, so an ordinary two-condition row needs N keys through a one-slot
channel.

- **(a) Panic stands.** A stack trace reaches a JSON consumer. Violates P3.
- **(b) Error channel on `Evaluate`.** 0007 rejected it: a second channel lets
  the evaluator report undecidability outside the verdict taxonomy. The
  objection conflates a *verdict* (cannot decide — modeled as
  `GuardUnevaluable`) with a *defect* (structure is broken), but under (b) the
  distinction is held by discipline rather than structure.
- **(c) Panic at the seam, recover at the CLI.** `panic`/`recover` must cross
  `internal/resolve` → `internal/cli`; recovering only the mapping panic needs a
  typed value and a re-panic default, and every non-CLI consumer inherits the
  obligation. Was the recommendation while "do not reopen Final RDRs" held.
- **(d) Carry the parsed predicate — recommended.** Row carries a slice of atoms
  (key, operator token, literal, block) instead of a string. No reconstruction
  step, so no mapping failure, so no channel question. The kernel carries
  *cardinality*, not *meaning* — it stays as grammar-agnostic as it already is
  for `Match []Tag` and `Writes []Tag`. **Dissolves JD-1.**

  *Cheap:* the only `GuardEvaluator` implementation in the repo is
  `fixtureGuards` in `fixtures_test.go:20`; `resolve.go` is the sole non-test
  file referencing the interface. *Not free:* reopens 0001's `Row` type, which
  0007 declined to touch (A3). `Implemented` has no backward edge, so this lands
  either as 0007 absorbing the change (precedent: it already redefines
  `Row.RequiresOwned`) or a new RDR superseding 0001's row shape.

*The current fixture is the argument:* given a guard it cannot map it returns
`GuardUnevaluable` — the exact conflation 0007's panic clause exists to forbid.

**Resolved: (d).** `Row` carries parsed atoms; the mandated panic and the
mapping-failure question both disappear. JD-1 is closed by this.

*Lands in **0007**, which absorbs the `Row` change (it already redefines
`Row.RequiresOwned`): the panic clause and A10/A15's mapping obligations drop,
and A3 ("a kernel change this RDR declines to make") is superseded — the decline
was justified by a cost the zero-consumer state does not support. Because the
chosen approach changes, 0007 re-enters the flow at propose rather than at a
later stage, and its `foundational` profile carries the full lens set. **0003**
cites the atom shape rather than restating it.*

## D2 — Does the aggregation veto run before or after match counting?

0002 states escape reachability as a function of the match count over non-escape
rows. 0007 refuses *before* counting: "if any surviving candidate row's guard is
GuardUnevaluable, the resolution MUST refuse `guard_unevaluable`." An
implementer reading 0002 builds two stages; reading 0007, three.

- **(a) Count first (0002's shape).** An unevaluable row is pruned, the
  resolution reports `no_match` — which 0002 makes escapable. Violates P1: the
  missing state is masked behind an escapable class. This is precisely the
  outcome 0007 was written to prevent.
- **(b) Gate first (0007's shape) — recommended.** Evaluate guards, veto on
  unevaluable, then count survivors. `guard_unevaluable` is not a modelable
  escape class, so the refusal reaches the user intact. Cost: 0002's Technical
  Design prose and its Scenario 4 expectation need restating at its re-lock —
  0002 states the counting rule in its own voice, so 0007's deferral to 0001
  does not reach it.

**Resolved: (b).** Gate, then count. `Resolve` evaluates guards over surviving
candidates, refuses `guard_unevaluable` if any is undecidable, and only then
applies exact-one matching and escape reachability. P1 is decisive: counting
first prunes the unevaluable row into an escapable `no_match`, which is the
masking 0007 exists to prevent.

*Lands in **0002**, which restates its resolver flow as gate-then-count and adds
to its Scenario 4 expectation the qualifier that no sibling candidate is
unevaluable; and in **0007**, which states the ordering as the kernel's.*

## D3 — Must a read accessor return a *complete* tag set?

0004 says only: "A read accessor MUST return typed tag values or a typed
refusal." A disjunction with no completeness requirement — a partially-read
artifact may conformantly return what it got. `Input.Owned` has no error
channel, so a truncated snapshot and a genuine absence are the same input.

0007 names this A6b, marks it `Pending`, and routes it to 0004's implement
stage. 0004 mentions 0007 zero times.

This is the floor under every presence/absence contract in the cluster: a
truncated read makes `exists = false` decide TRUE, and a plan is produced
against state that exists. It is also the easiest entry to wave through, because
0004's clause reads fine until you notice what it omits.

- **(a) Status quo.** Partial reads stay conformant and indistinguishable.
  Violates P2 — the kernel guesses.
- **(b) Require completeness — recommended.** A read accessor MUST return the
  complete tag set for the keys it was asked for, or take the refusal branch.
  States the rule 0007 needs at the layer that owns it. Cost: one normative
  clause in 0004 plus an MVV scenario asserting a partial read refuses.
- **(c) Error channel on `Input.Owned`.** Lets the kernel distinguish
  read-failed from absent. More expressive, more surface, and 0004 still has to
  say when to use it — so (b) is a prerequisite either way.

**Resolved: (b).** A read accessor MUST return the complete tag set for the keys
it was asked for, or take the refusal branch; a partial read is a refusal, not a
value. (c) stays available if implementation shows the kernel needs to
distinguish read-failed from absent, but (b) is the prerequisite either way.

*Lands in **0004** as one normative clause plus an MVV scenario asserting a
truncated read refuses — no assumption of its own is disturbed. This closes
0007's A6b, which 0004 now carries as its own obligation rather than receiving
it by reference from an RDR it never cites.*

## D4 — Who enforces the guard domain once the row carries atoms?

§D1 put parsed atoms on the row and said the kernel "carries *cardinality*,
not *meaning*." That left one question open: 0007's domain rule — an atom over
an absent key is unevaluable, never false — can now be applied by the kernel,
which sees every atom's key, or left to 0003's evaluator as before. The prior
0007 Final held it in the evaluator plus a conformance harness and recorded the
residual as UNMITIGATED: nothing compels 0003's build to run the harness.

- **(a) Evaluator-enforced.** Row-level seam `Evaluate(atoms, view)`; the
  evaluator decides presence and combines verdicts. Kernel stays
  grammar-blind. The domain rule is held by discipline; drift at the seam can
  yield a plan from absence. Violates P1 conditionally.
- **(b) Kernel-enforced — recommended.** The kernel decides presence
  (provenance-blind), decides `exists` from presence alone, marks an
  absent-key value atom unevaluable **without calling the evaluator**, hands
  only present-key value atoms to a per-atom seam `Evaluate(atom, value)`,
  and combines under strong Kleene itself. Prior art: SQL:2003 strict
  routines — on a null argument "the function itself is not invoked." Cost:
  the kernel learns one grammar fact (the existence operator token and its
  two boolean literal forms), exported as constants 0002's normalizer MUST
  emit. Drift at that boundary fails closed (unevaluable refusal or load
  rejection), never as a masked plan.
- **(c) Kernel evaluates every operator.** Typed kinds and domains are 0003's
  declarations; the kernel has none. Dissolves 0003's evaluator.

**Resolved: (b).** The kernel enforces the guard domain; the evaluator decides
value comparisons over present keys and never sees the view. The
`guard_unevaluable` refusal names what blocked the verdict per row and per
atom — key, block, reason `absent` | `uncomparable` — replacing the
single-valued `Refusal.Guard` text, which has no referent once guards are
atoms; it reaches the CLI through §JD-8's `Detail`. P1 and P4 decide: the
masking path closes by structure rather than by a harness nobody runs, and the
refusal names the missing key.

*Lands in **0007**, the single normative home of the seam, the domain rule,
and the payload. **0003** cites it: its evaluator's scope narrows to value
semantics over a present value, and its placement MUST is read against
`exists = false` there. **0002** cites it: the normalizer emits the kernel's
existence constants and canonicalizes key spellings before a row exists.
**0009** (Final) cites `Refusal.Guard` as a refusal property; that citation is
stale under this decision and rides to 0009's re-lock — see JD-12.*

## D5 — Does `<clear>` cross to the write accessor as a value, and what does read-back assert?

0002 renders an authored rule-level `clear` as a `<clear>` entry in both
`Writes` and the next-state tags, and states the sentinel is **unreserved** in
the tag-value space. 0009's fenced Writes-only predicate relies on that
reduction. 0004 contains the word "clear" zero times: a literal-minded write
accessor stores the string `<clear>`, and its read-back ("expected owned-tag
values") passes green on a tag that was never removed — a success-shaped result
over unverified state. 0006 cannot re-pair a clear from the normalized value
while an author may also write the literal.

Prior art: a delete is a write of a marker (Cassandra tombstones; JSON Merge
Patch's `null`; `kubectl label k-`) — sound exactly when the marker is
**reserved**.

- **(a) Reserve the sentinel; define it at 0004 — recommended.** 0002 refuses
  an authored tag value of `<clear>` at load; 0004 states that a `<clear>` write
  removes the key, that read-back asserts the key is **absent** (already a value
  on its success branch), that clearing an absent key succeeds, and that a read
  yielding `<clear>` as a value is unreadable. Kernel `Row` and 0009 untouched.
- **(b) Separate `Clears []string` on the row.** The clean shape (P7), but it
  reopens 0001's `Row` and 0009's fenced Writes-only predicate — a second
  re-lock for one field.
- **(c) Typed absent value on `Tag`.** (b)'s cost with a weaker type.

**Resolved: (a).** Three Final documents and the kernel already carry the
sentinel; reserving it closes the ambiguity and the "lossy dump site" 0002
lists. Closes JD-15.

*Lands in **0002** — one load category refusing `<clear>` as an authored tag
value, and the Round-Trip lossy-site note drops its `<clear>` item; and in
**0004** — remove-key / read-back-absent / idempotent-clear / unreadable-on-read
clauses plus an MVV scenario for a clearing rule. **0009** is unchanged;
**0006** reads a `<clear>` write as removal by citation.*

## D6 — Which key routes an atom to `Match` or the guard?

0002 routes by **operator** (`eq` on a declared tag → `Match`; every other
operator and every `unless` atom → guard). 0003 proves coverage by **authored
block** (`[rule.match]` is the grouping context; guard blocks are product
dimensions). The partitions disagree on two input classes: a non-`eq` atom
under `[rule.match]` is green under lint and `guard_unevaluable` at runtime (a
P5 failure); an `eq` atom under `[rule.guard.all]` over an optional key — 0003's
own fixture `foundational-to-cove` — is withheld by lint but becomes a silent
non-candidate in the kernel (escapable `no_match`, the masking P1 forbids).

Prior art: statecharts/SCXML separate trigger match (candidate selection) from
`cond` (gate on candidates), and the author declares which is which; rules
engines and SQL (join key vs `WHERE`) are the same split.

- **(a) Operator-keyed (0002 today).** The author's block is advisory; 0003 must
  re-derive routing from operators; `eq` under a guard silently becomes
  selection.
- **(b) Block-keyed, match restricted to what the equality pattern carries —
  recommended.** Match blocks (`[rule.match.*]`, `[context.*.match.*]`) admit
  `eq`, and `in` only via expansion into one candidate row per member — the
  expansion 0002 already defines for `recognized`, generalized. Comparison,
  existence, and `contains` under a match block are a load refusal. Every atom
  under `guard.all` / `guard.unless` is a guard atom regardless of operator,
  `eq` included. A12's exhaustive/disjoint routing becomes structural.
- **(c) Block-keyed, `in` under match refused.** Forces `profile in […]` into
  the guard, where lint reads it as a dimension the group must *cover* —
  gap findings for values the author meant to exclude. Rejected.

**Resolved: (b).** The block is the author's declared intent about masking:
"select on" versus "gate on" (P1, P2). 0002's worry that an `eq` guard atom
"defers a decidable check to a seam that can report `guard_unevaluable`" names
the wanted behavior for a guard over an optional key, and 0002 already does it
for `unless`. Closes JD-16.

*Lands in **0002** — the routing clause restated as block-keyed; the match-block
operator restriction as a load category; the `in`-expansion generalized to every
match block; match atoms carry block `match` (closing the two-valued/three-valued
block vocabulary finding); A12's MVV; Scenario 4's unevaluable-sibling witness
becomes an `eq` guard atom over an absent optional key. **0003** cites: its
participation and can-refuse clauses already read the block. **0006** reads the
three-valued `Block` it already serializes.*

## D7 — What must the closed TOML layout additionally spell?

0002 locks a strict-decoded layout that never received the declarations its
consumers locked against: 0006's root and stop set (A6) and write-replaces
(A10); 0004's per-accessor capability, role, keys, timeout, read-back; 0003's
five type-model fields; and the rejection sites for a bad declaration or a
literal outside its domain. Pre-release, the layout may take any shape (P7), so
the tiebreaker is readability and maintenance.

Two rules decide every sub-question. **Illegal states are unrepresentable in
the layout** — a shape the strict decoder can refuse is never left to a
validator, so the error names a key and a line. **One home per fact** (P6) — a
binding both 0002 and 0004 need is declared once. Prior art: SCXML/XState/ASL
put start and end in the same document as the transitions; NuSMV declares
typed variables (`{a,b,c}`, `0..3`, boolean), roots as assignments
(`init(x) := a`) and goals as predicates; Kubernetes probes model one-of by
sub-object; Go tool configs spell timeouts as duration strings; Cargo keeps a
strict manifest with a sanctioned `[package.metadata]` extension table;
Cassandra replaces scalars and adds collection `+`/`-` forms later without
breaking replace.

- **(i) Initial / terminal.** Root `[initial]` is a table of owned
  `tag = value` assignments (the root must be a node); root `terminal` is a
  list of **context ids** (the stop set is a predicate, and contexts are already
  the model's named predicates — no new grammar; a rule may `use` a terminal
  context for an escape self-loop). Rejected: `terminal = true` inside a
  context (the first non-predicate field in a context; terminals found only by
  scanning); `[[terminal]]` tables (duplicate the context grammar). A terminal
  context matching a non-owned tag is a 0006 blocking finding.
- **(ii) Accessors.** One table per capability — `[read.<id>]`, `[write.<id>]`,
  `[gate.<id>]` — replacing `[accessors.<id>]` and its `mode`. Each decodes to
  its own shape, so `keys` exists only on readers and `read_back` only on
  writers; 0004's "exactly one capability" holds structurally. **`keys` is the
  binding**: 0004 already mandates the requested key set as validated metadata,
  so the tag-side `accessor` reference is a second copy and is removed. Readers
  and writers both declare `keys`; the loader refuses a key served by zero or
  two readers, a written or cleared key not in exactly one writer, and a
  writer key that is not owned. Fields: `role`, `path`, `keys`, `timeout` (Go
  duration string; missing or non-positive refused), and `read_back = true` on
  writers. The same id in `[read.x]` and `[write.x]` is legal — 0004's fenced
  identity is `(flow, name, capability)`; its unfenced "cannot be rebound"
  sentence is a citation repair.
- **(iii) Type-model keys.** `kind`, `domain`, `min`, `max`, `elements`,
  `single_valued`, `required` — 0003's own illustrative spellings, with
  `required` for the optionality marker (defaults optional, as 0003 mandates).
  Rejected: JSON-Schema alignment (`type`/`enum`/`minimum`) — `kind` is in
  every fixture and in 0003's tokens, and `enum` would name a kind in one place
  and a value list in another. Two load categories in 0002, rules supplied by
  0003: `malformed tag declaration` (kind/domain disagreement; before
  normalization completes) and `malformed predicate atom` (literal outside
  domain; after declarations load, before rows are yielded). 0006 mints nothing
  for either; 0005 maps them under §JD-8.
- **(iv) Write-replaces.** A clause, not a key: a write assigns a tag's whole
  value and supplants what was held; for `set` kind the array literal is the
  whole new set; no accumulate form. 0004's read-back compares the held value
  for **equality**. Closes 0006 A10 and G7 (one value per tag per write).
- **(v) Extension table.** `[model.metadata]` is a free-form table the loader
  carries untouched and never interprets — strictness everywhere else, one
  sanctioned namespace for tooling.

**Resolved: (i)–(v) as stated.** Closes JD-17.

*Lands in **0002** — root `terminal` and `[initial]`; the three capability
tables; `[tags.<tag>].accessor` removed; the seven type-model keys enumerated;
`[model.metadata]`; the write-replaces clause; the two load categories; the
binding validations; fixtures re-authored. **0003** — `min`/`max` spelling
becomes a citation to 0002; supplies the two rejection rules. **0004** — its
definition is the capability-table entry; read-back is equality; writers carry
`keys`; the "cannot be rebound" prose repaired. **0006** — A6/A10 flip by
citation; the terminal-non-owned finding.*

**Blank, not decided here:** where a *rule* references a gate accessor — 0002's
layout has no site for it; owner 0002 with 0004's semantics.

**Observation for 0003's next touch (not this decision):** an unmarked
`kind = "enum"` is a `2^|domain|` dimension under 0003's conservative default;
every author will read `enum` as single-valued. The wire key is right either
way.

---

## Interface record

Cite `JDR 0001 §JD-n`; never restate. Entries marked *(blank)* name an owner,
not a negotiation.

- **JD-1 Guard atom transport.** **Closed by §D1.** `Row` carries parsed atoms,
  so a row's N predicates cross the seam natively; 0003's per-atom identity no
  longer has to survive a one-slot channel. *(0007×0003 F1)*
- **JD-2 Resolver control flow.** **Closed by §D2** — gate, then count.
  *(0007×0002 F1)*
- **JD-3 `RequiresOwned` producer.** **CLOSED 2026-08-23** by the `0002-0009`
  iteration-2 gate: RDR 0002's normalizer is the producer (`0002::Normative
  Contracts`, "`Row.RequiresOwned` has no authored form. The normalizer MUST
  derive it … sorted, duplicate-free set of tag keys named by the rule's write
  block and clear list"; escape rows carry an empty set). 0007 owns the field's
  meaning; 0009 is consistent (escape rows carry neither writes nor a set). One
  implementation-order residue, not a decision: whichever of 0007 Phase 1 /
  0009 Phase 1-2 lands first repairs `fixtures_test.go::escapeRow` on both
  `Writes` and `RequiresOwned`. *(0007×0002 F3, 0007×0009 F2; closed via
  `docs/rdr/cluster-reconcile/0002-0009/iter-2/`)*
- **JD-4 Lint's promise is what gives.** Where 0006's exhaustiveness proof and
  0007's veto disagree, the *promise* narrows — P5 decides the substance.
  **CLOSED 2026-08-22** by the `0003-0006-0007` cluster gate, which is the venue
  0003 A8, 0006 Direction 1, and this entry's own recommended arm all routed to.
  **RDR 0003 is the recording document; RDR 0006 cites it and mints no code.**
  Both halves are now decided:
  - **Recording.** 0003 records the narrowing normatively
    (`0003:781-786`, already written). 0006 carries a citation, never a
    restatement. The two triggers genuinely differ and both must survive: 0006's
    clause fires on a *non-finite dimension*; 0003's on a *fully-finite product
    whose participating row can still refuse*. 0003 is the recording document
    because the clause must name the **refusing atom**, and `atom` occurs once
    in 0006, only to delegate atoms to 0003.
  - **Warning category — no.** Lint gains no new code and no non-blocking tier
    for this class. 0006 reuses `graph-unprovable-coverage`, whose stated
    meaning ("Required finite-domain proof unavailable") already covers a
    withheld claim and which already sits in 0006's *mandatory blocking* table
    (`0006:314`) and its MVV fixture matrix (`0006:639`) — so reuse costs no new
    test scaffolding, only a widened trigger. 0006's existing advisory tier stays
    scoped to "redundant rows or unreachable rules" (`0006:302-304`) and MUST NOT
    absorb the withheld claim. This retires the 0003 clause's dependency on this
    entry (`0003:818-826`, "whether lint gains one at all is open under §JD-4").
  - **Consequent duty on 0006 (not a restatement).** 0006's finding contract
    (`0006:363-367`) carries no atom-level field, so the atom-naming half of
    0003's clause has no producer. 0006 MUST extend the finding record to carry
    the refusing atom. This is a payload extension, discharged at 0006's refine.
  *(0007×0006 F2, 0007×0003 F3; closed via `docs/rdr/cluster-reconcile/0003-0006-0007/`.
  The two 0007×0008 rows iteration 1 routed here — absent `recognized` key,
  `exists` on the reserved key — are answered by §D4, not by this entry: the
  kernel decides presence and `exists` provenance-blind.)*
- **JD-5 Precondition precedence.** 0009's breach check and 0008's reserved-key
  check both land at `Resolve` entry; neither orders itself against the other.
  Either order is defensible — pick one and pin it with a test on a table that
  breaches both. The silence is the defect, not the choice. *(0007×0009 F1,
  0008×0009 F2)*
- **JD-6 Producer-defect surface.** **Closed by §D1** — with no reconstruction
  step there is no mapping failure to surface, so 0007's panic clause has no
  trigger and 0005's structured-envelope contract is unopposed. *(0007×0009 F3)*
- **JD-7 Read completeness.** **Closed by §D3** — a partial read refuses.
  *(0007×0004 F1/F3, 0009×0004 F1)*
- **JD-8 Refusal codes.** `owned_state_unavailable` and `reserved_tag_key` need
  `Code` values; 0009 needs `GroupInternal`, which **already ships** in
  `clierr.go`. 0005's own A-block pre-authorizes this: resolver-specific values
  "only require new `Code` constants or literals, **not new envelope fields or
  exit groups**." `Detail` and `Hint` ship and render in text and JSON. One live
  question: `GroupUserEnv` and `GroupInternal` both exit 2, so the only remedy
  distinction an exit code can carry is exit 3 — and an accessor that could not
  read the artifact is what `GroupEnvUnavailable` already means. *(blank, except
  the exit-3 call)* *(0007×0005 F1/F2/F3, 0009×0005 F1/F2, 0008×0005 F1)*
  **Widened 2026-08-23** (iteration-2 gate): the open question is the whole
  code table, not three codes. 0005's June-locked table (`0005:601-614`) has
  no row for 0004's refusal family (`incomplete_read`, post-mutation timeout,
  gate deny — `read_back_incomplete` currently collapses into
  `flow-write-readback-mismatch`, which 0004 forbids), 0002's ~20 load
  categories, 0003's predicate semantic kinds, or an escape-disposition marker
  on a resolve plan; and 0005 has no path form for the model 0002 says enters
  at the CLI seam. Two carrier questions ride here: 0007's Prerequisite asks to
  reopen §D4 for one `omitempty` structured `CLIError` field instead of
  `Detail` text (`0007:2028-2039`), and 0006's non-omitempty `Findings` field
  changes every `flow-*` envelope. 0005 is the sibling this answer re-walks.
- **JD-9 `--tag` provenance.** Caller-supplied tags enter as
  `ProvenanceObserved` and never satisfy an owned-state dependency. `assemble`
  already pins owned-over-observed precedence so the accessor snapshot is "never
  shadowed by caller-supplied context"; the exposure opens only if the `flow`
  verb wires `--tag` into `Input.Owned`. Wire it to `Observed`.
  *(0007×0005 F4, 0008×0005 F2)*
- **JD-10 Recognized-tag totality.** 0008's name constraint invalidates all
  three of 0002's canonical fixtures, which 0002 declares normative — rename
  them; there are no users to migrate. Still open: whether a declared
  `recognized` tag is total (always present, possibly empty) or partial (absent
  with no outcome in flight), which decides whether an empty-outcome row is
  satisfiable, dead, or a lint error. No normalizer exists yet to observe it.
  *(0008×0002 F1/F2/F3, 0008×0009 F3)*
  **ANSWERED 2026-08-23** — ratified at this home by the iteration-2 gate from
  RDR 0002's re-lock: the `recognized` tag is **total over matching** — the
  kernel refuses `unmodeled_outcome` before any row is consulted, the single
  `recognized` match atom is the rule's mandatory outcome binding and is
  lifted into `Row.Outcome`, a `recognized` atom in a guard block is refused at
  load, the `outcomes` alphabet is non-empty and never contains the empty
  string, and a rule requiring `recognized` absent is dead (0006's
  unreachable-rule finding). Fixtures already renamed. 0008's answer-vs-fences
  check failed (A9 verified on the superseded spike; block 2's
  predicate-position clause; scenario 5's two-failure count vs fail-fast load)
  — a SPEC-DEFECT routed to 0008 at Stage 4, STAGE-SCOPED.

- **JD-12 Guard enforcement site.** **Closed by §D4** — the kernel enforces
  presence, `exists`, and combination; the evaluator seam is per-atom over a
  present value; the normalizer emits the kernel's existence constants; the
  refusal payload is per-atom and replaces `Refusal.Guard`. 0009's citation of
  `Refusal.Guard` as a refusal property is stale and is corrected at its
  re-lock; entry kept so the anchor never dangles. *(0007×0003 F1/F2 residue,
  0007×0002 normalizer obligation, 0007×0009 `Refusal.Guard`)*

- **JD-13 Single-valued grouping has no producer.** 0006's lint input contract
  attributes "single-valued grouping when applicable" to "the RDR 0003 tag
  declaration model" (`0006:261-264`), and gates a *mandatory blocking* code on
  it (`graph-single-valued-state`, `0006:315`, asserted by its MVV at
  `0006:640`). But `single-valued` occurs **zero** times in 0003 and **zero**
  times in 0002 — 0003's declaration model enumerates value kind, finite domain,
  optionality, and set-element universe, and stops there (`0003:670-678`). This
  is the same unhomed-producer defect the 2026-08-21 declaration-model rehoming
  was performed to close, surviving in one field the rehoming did not sweep.
  **Decided 2026-08-22: 0003's tag declaration model gains the field**, beside
  the four it already carries — it is a property of a tag class, it is consumed
  by a finite-domain proof, and every sibling property already lives there.
  0006 cites it like the rest of the model. Discharged at 0003's refine, with
  0006's citation at its own. *(0003×0006 F7 — new at this gate, on no prior
  list)*
- **JD-14 Escape rows in the coverage union and the overlap check.** 0003 states
  normatively that an escape row "participates in the coverage identity like any
  other row", that lint "MUST NOT treat 'an escape row exists' as a separate
  coverage-satisfying fact outside the union", and "MUST NOT exclude escape rows
  from overlap checks" (`0003:771-779`). 0006's invariant 3 grants the opposite
  exemption — two rows may overlap "unless the model explicitly routes to one
  deterministic escape row" (`0006:286-288`) — and its invariant 4 phrases escape
  rows as an alternative disjunct satisfying coverage outside the union
  (`0006:289-291`). The same fixture earns a blocking `graph-overlap` under 0003
  and passes under 0006. **0006 is additionally split against itself**: its own
  Load-Bearing Decision says "explicit escape rows are modeled graph edges, not a
  tie-breaker" (`0006:412-414`), siding with 0003 against its own invariant 3.
  **Decided 2026-08-22: 0003's reading governs** — escape rows are ordinary
  participants in both the union and the overlap check, per P5 (a green lint must
  imply resolution succeeds) and RDR 0001's runtime refusal of ambiguity. 0006
  repairs invariants 3 and 4 to match its own Load-Bearing Decision at its
  refine. *(0003×0006 F8 — new at this gate, on no prior list)*
  **Corrected 2026-08-23** (iteration-2 gate): the *overlap* half above
  overshot. Both members now state, in fenced text, that overlap is checked in
  **two populations** — ordinary rows among themselves, and escape rows among
  themselves per declared failure class — never escape-vs-ordinary, because
  the runtime consults escape rows only to rescue a `no_match`/`ambiguous_match`
  (`0003:1206-1217`, `0006:901-903`). The *coverage* half stands: escape rows
  participate in the union. This entry's decision is amended to that reading;
  0003 A17 closes on it.
- **JD-15 `<clear>` at the write accessor.** **Decided 2026-08-24 by §D5** —
  the sentinel is reserved: refused as an authored tag value at load (0002); a
  `<clear>` write removes the key, read-back asserts absence, clearing an absent
  key succeeds (0004). Restored from iteration 1's lost "JD-12". Siblings:
  0002, 0004, 0009. *(0009×0004 F2 ledger; 0002×0004 F-1, 0002×0006 G3,
  iteration 2)*
- **JD-16 Match/Guard routing key.** **Decided 2026-08-24 by §D6** — the
  authored block routes: match blocks admit `eq` and `in` (expanded per
  member), refuse other operators at load; every guard-block atom is a guard
  atom regardless of operator. Siblings: 0002, 0003. *(0002×0003 F1, iteration
  2)*
- **JD-17 What the closed TOML layout must additionally spell.** **Decided
  2026-08-24 by §D7** — root `[initial]` assignments and `terminal` context-id
  list; per-capability accessor tables with `keys` as the binding; type-model
  keys `kind`/`domain`/`min`/`max`/`elements`/`single_valued`/`required`; a
  write-replaces clause; `[model.metadata]`; two load categories for
  declaration and literal-domain errors. Siblings: 0002, 0003, 0004, 0006.
  *(0002×0006 G1/G2, 0002×0004 F1/F2/F3, 0002×0003 F2/F3, iteration 2)*
- **JD-18 Conforming-view enforcer.** 0003 A18/A20 route the runtime half of
  the declaration-conformance check (an observed always-present key omitted at
  runtime refuses under a green lint) and the atom carrier to "0007's next
  touch"; 0007 is Final and silent; 0006 enforces only the owned half
  (`0006:935-940`). Open: which document states the view-level check and where
  it runs. Siblings: 0003, 0007. *(0007×0003 F1, 0003×0006 A18/A20,
  iteration 2)*

**Withdrawn — JD-11 Escape-row identity across the dump.** Re-triaged as a
single-RDR defect: 0002's round-trip invariant requires the dump to preserve
"row kind … and escape failure classes" while its dump-derivation contract
sixty lines away omits both. Both are `normative` blocks inside 0002, so the
contradiction is visible reading 0002 alone; 0009 only escalates it to
breach-laundering. Routed as a spec defect against 0002, fixable at re-lock with
no round re-runs. Entry
kept so the anchor never dangles. *(0009×0002 F3)*

## Ownership corrections

Recorded here because they change which document a cluster member cites, and a
reader arriving at a stale citation needs the pointer.

- **Tag declaration model → RDR 0003** (2026-08-21, from 0003's Stage 6
  reconcile). The typed alphabet guard atoms are written against — value kind,
  finite domain, per-tag optionality, set-element universe — was booked as a
  producer request against 0002 by 0003's A11/A7/A9. It is **rehomed to 0003**.
  The finding: 0002's only normative tag clause requires *provenance* alone;
  `value kind` is normative nowhere in the cluster (0002 names it in a prose
  schema list, and peers quote that prose as if it were a contract); 0002's
  normative validation categories are entirely structural, with no type or
  value-domain axis; and **0007 already records the same division** — "typed
  literals and value kinds (bounded integers, set universes) are RDR 0003's
  declarations" (its rejected-alternative analysis). The model was homeless, not
  homed elsewhere. 0002 keeps provenance, authoring location under
  `[tags.<tag>]`, and normalization carriage, and cites 0003 for meaning; 0006's
  lint input contract re-points at 0003 for finite-domain metadata. This follows
  §D1/§D4's own pattern — one home per contract, cite never restate, and the
  home is the document that owns the surrounding contract. One field does **not**
  travel: per-atom `block` retention (0003 A14) stays a request on 0002, because
  retention through normalization is carriage, not declaration semantics.

## What this does not decide

- **RDR 0008's re-entry** — A6/A11 verified on a false claim about Final 0009;
  already demoted; re-enters at resolve, scoped to those two assumptions.
- **RDR 0002's dump contradiction** — JD-11 above; a 0002 re-lock.
- **RDR 0004's Prerequisites** — `Final` with every box unchecked, including
  "All Critical Assumptions verified." A finalize-gate question for 0004.
- **RDR 0009's frozen API surface** — `Unwrap() []error` "EXACTLY ONE LEVEL",
  "return CheckValid's error VERBATIM". Internal to 0009.
- **Exact code strings, exit integers, envelope field names** — the RDRs own
  them once the decisions above fix the shape.
