model: claude-opus-5[1m]

# Repeatability diff — RDR 0019 (owned-state initialization semantics)

Neutral pass. Author of none of the three runs. Runs compared:

| Run | Base model | Variant header |
| --- | --- | --- |
| run-1 | claude-opus-5[1m] | `full (profile: foundational)` |
| run-2 | claude-sonnet-5 | `full (profile: foundational)` |
| run-3 | claude-fable-5-1 | `full (profile: foundational)` |

All three are expected full-coverage runs (run-1's `variant:` line records
the intended variant for the set). run-3 is the alt-model run; a split
along run-3 vs. {run-1, run-2} is flagged as MODEL-BOUNDARY below.

RDR read through the projector only (`inspect --json --filter
elements,outline`, then `--select` on `0019:C1`, `0019:MVV`,
`0019:D-identity`, `0019:D-wire-byte-format`, `0019:D-naming`,
`0019:D-selection-predicate`, `0019:RT1`–`RT4`, `0019:S1`, `S4`–`S7`,
`S9`–`S12`, `0019:F2`–`F5`, `0019:A6`, `0019:A8`, and the
`§mini-checks`, `§illustrative-code`, `§approach`, `§step-1-verb-skeleton`,
`§step-2-predicate-and-plan`, `§step-3-apply-and-report` sections).

---

## 1. Disagreements

### D1 (3-way) — Exit group of a read-back mismatch — `0019:C1`

| Run | Rendering |
| --- | --- |
| run-1 | `non-zero (2)`, plus a SECOND, separate row: "read-back incomplete (unreachable locator) → exit 3, write possibly applied, inherited `set-state` semantics, distinct from the mismatch class" |
| run-2 | `2 (existing read-back-mismatch family)`, plus a distinct row asserting an UNREADABLE read-back is "inherited `set-state` semantics (exit 3 class)" but caveating that a sealed store is caught by the predicate first as a no-op success |
| run-3 | `non-zero (GUESS: exit 3, the read-back family's group)` — listed in its GUESS ledger item 3 |

RDR passage that let them diverge: C1 fixes only "A read-back mismatch
is a distinct terminal refusal naming the key as PRESENT-AND-UNVERIFIED",
and the `disposition` table's row reads `Read-back mismatch | non-zero |
partial, per-key | key PRESENT-AND-UNVERIFIED`. **No exit group is
stated.** C1 fixes the group explicitly (2) for the two NEW codes and
says shared classes "reuse the existing `flow-*` codes unchanged" —
but it never says which group the shared read-back class carries, and
the record never distinguishes mismatch from incomplete/unreadable
read-back at all. Two of three runs invented that second row from
outside the record. This is the single largest divergence in the set:
three different exit dispositions and two different taxonomies for one
contract row.

**Rewrite candidate**: C1's read-back clause should fix the exit group
for the mismatch arm, and either fold in or explicitly exclude the
unreachable/incomplete read-back arm.

### D2 (3-way) — Position of role-binding validation relative to the emptiness read — `0019:C1`

| Run | Rendering |
| --- | --- |
| run-1 | class → registry → carrier → **plan build + role-binding refusal** → emptiness → apply. Plan validation precedes the predicate. |
| run-2 | class → registry → **role-binding refusal** → carrier → emptiness → plan build → encode → apply. Role-binding precedes the CARRIER gate. |
| run-3 | class → registry → carrier → **role-binding refusal** → emptiness → `planInitial` (encode + group) → apply. Role-binding after carrier, before predicate; plan ENCODING split off after the predicate. |

RDR passages that let them diverge — the record carries two orderings
that do not agree:

- The `trace` table's step sequence is `bind model + artifacts → class
  check → carrier check → emptiness read → plan validation → encode +
  write → read-back → second run`, which puts **plan validation AFTER
  the emptiness read**.
- C1's predicate clause says the opposite: "'Bound' is not a filter
  that can shrink the quantifier's domain: the plan-validation clause
  below requires every role a needed writer names to BE bound, and
  refuses otherwise, so the artifacts quantified over here are exactly
  the ones that clause already established." That requires the
  role-binding half of plan validation to run **BEFORE** the predicate.
- `§step-2-predicate-and-plan` states it a third way: "Read the bound
  store; non-empty → the no-op arm... Empty → plan every `[initial]`
  assignment and validate the WHOLE plan (`writerFor` per key,
  artifact roles bound) before any accessor runs" — predicate first,
  all plan validation after.

Consequence the runs are actually splitting on: whether `init-state`
against a non-empty store with an UNBOUND needed role exits 0 (no-op,
per `trace`/Step 2) or exits 2 (artifact-binding refusal, per C1).
run-2 and run-3 refuse; run-1 refuses only because it builds the plan
before the predicate — but run-1's plan build also encodes, which under
`trace` happens after. S6 does not disambiguate: it asserts the refusal
is "reached at the verb" but does not fix the store's state.

**Rewrite candidate**: C1 (or `§step-2`) must state whether the
role-binding refusal preempts the no-op arm. Three runs, three
orderings — this is under-determined.

### D3 (3-way) — Success payload shape — `0019:C1` / `0019:D-wire-byte-format`

| Run | Rendering |
| --- | --- |
| run-1 | `{"outcome": "seeded"\|"noop", "seeded": [{"role","key","value"}], "absent": [...]}` — seeded entries are OBJECTS carrying role and value. Flags as GUESS whether `absent` appears in the seeded arm and whether values are echoed at all. |
| run-2 | Declines to name fields at all; states only the two required distinctions and that names are deferred. |
| run-3 | `initStatePayload{Arm string; Seeded []string; Absent []string; Artifacts map[string]string}` — seeded is a KEY-NAME list, no values; adds an `artifacts` field echoing bound roles (self-marked GUESS); fixes `Absent: []` on the seeded arm. |

RDR passage: `D-wire-byte-format` — "Field names deferred to
Resolve/Pre-Lock, owned here." C1 fixes only that the payload MUST
distinguish the seeded-all case from the no-op case, that its scope is
exactly the `[initial]` key set, and that the no-op arm reports which
`[initial]` keys the store does not carry.

This is a DECLARED silence, not a defect — but the runs diverge past
the declared boundary on two things the record arguably does fix: (a)
whether the seeded arm echoes VALUES (run-1 yes, run-3 no, run-2
undecided) — RT1 asserts `read-state` reports the values, not the init
payload; and (b) whether `absent` appears on the seeded arm (run-1
GUESS, run-3 asserts `[]`, run-2 silent).

### D4 (3-way) — Emptiness probe: carrier and signature — `0019:C1` / A6

| Run | Rendering |
| --- | --- |
| run-1 | Picks the `flowbind` cardinality probe; `func (r Reader) KeyCount() (int, error)` — count of STORE keys, incl. `sealedKey`. Marked GUESS. |
| run-2 | Refuses to pick. "I mark it UNDECIDED rather than guessing a specific signature." |
| run-3 | Picks the exported `flowbind` cardinality probe; `func (r Reader) Empty(art) (bool, error)` returning `len(store) == 0`; plus a per-role wrapper `storeIsEmpty(ctx, reader, art) (bool, error)` and a `type emptiness struct{ Role string; Empty bool }`. Marked GUESS. |

RDR passage: C1 — "Which carrier serves it — a new `ReadBinding`
capability, an exported `flowbind` cardinality probe, or a
`read-state`-family surface — is UNDECIDED here and is booked as A6".
A6 status is **Pending**.

Contract-acknowledged silence, correctly identified by all three. The
divergence is in DISCIPLINE, not comprehension: two runs guessed a
signature anyway (and produced incompatible ones — count-returning vs.
bool-returning), one declined. Not a rewrite candidate for C1; it is
the A6 resolution's job.

### D5 (2-way, MODEL-BOUNDARY) — Whether `--allow-commands` is on the verb's flag set — `0019:C1`

| Run | Rendering |
| --- | --- |
| run-1 | Omits it from the CLI surface entirely; mentions `--allow-commands` only in the ordering prose about preemption. |
| run-2 | Omits it from the CLI surface entirely; same preemption-only mention. |
| run-3 | **Lists it in the synopsis**: `[--allow-commands] # accepted for uniformity; preempted by the carrier refusal (C1 ORDERING)`, and books it as GUESS ledger item 7: "C1 discusses its preemption, implying the flag exists here; not stated outright." |

RDR passage: C1's CARRIER clause enumerates the flag set as "the
shared selection flags (`--flow`|`--model`), explicit `--artifact
role=path` bindings, and no write grammar" — `--allow-commands` is not
in that list. But C1's ORDERING clause then says the carrier refusal
"PREEMPTS the `--allow-commands` refusal... A command-backed accessor
therefore yields the carrier code whether or not the opt-in was
passed", and S10 says its command-backed rows "are run with
`--allow-commands` UNSET" — both presupposing the flag is acceptable
on this verb.

Splits on the model boundary and the alt-model run is the one that
noticed. The record does presuppose the flag exists on `init-state`
without listing it in the flag enumeration; a reader who takes the
CARRIER enumeration as closed concludes it does not.

### D6 (2-way, MODEL-BOUNDARY) — Whether the class arm and the carrier arm have a fixed relative order — `0019:C1`

| Run | Rendering |
| --- | --- |
| run-1 | **Explicit unresolved GUESS**: "the relative order of the CLASS check and the CARRIER check. `trace` lists class before carrier; C1 fixes only that class refuses 'before any accessor runs' and carrier 'after registry construction'. Both are pre-invocation, so a decision-table model with a command-backed writer could take either code — the record does not resolve which." Listed as RDR silence #5. |
| run-2 | Asserts class-before-carrier as settled (pseudo-code step 2 class, step 4 carrier); no GUESS marked. |
| run-3 | Asserts the full chain as "Ordering fixed by C1: model load → class refusal → registry construction → carrier refusal..."; no GUESS marked. |

RDR passages: C1 says the class refusal fires "before any accessor
runs" and, separately, that the carrier refusal is "decided AFTER
registry construction... but BEFORE any accessor is invoked". Both are
pre-invocation constraints that do not order against each other. Only
`trace` (a Mini-check table, not a Normative Contract) lists class
before carrier. S7's fixture "declares zero owned tags", so it can
declare no writer and cannot trip the carrier arm — the tests do not
resolve it either.

Two of three runs treated a Mini-check row's ordering as normative.
Whether that is legitimate is itself the finding: `§mini-checks` says
"Each table is the decision, not a note about it", which arguably makes
`trace` normative — but it is not in the Normative Contracts
subsection, so a reader following the section's own authority
declaration and a reader following the RDR template disagree.

**Rewrite candidate**: C1 should state the class/carrier relative order
directly (a decision-table model cannot have owned writers, so the
combination may be unconstructible — if so, say so, as C1 already does
for the `[initial]`-carrying decision-table case).

### D7 (2-way, MODEL-BOUNDARY) — Exit disposition for model-selection failure — `0019:C1`

| Run | Rendering |
| --- | --- |
| run-1 | Not enumerated. |
| run-2 | Not enumerated. |
| run-3 | Row present: `Model selection failure | 2 (GUESS: same group as set-state) | zero | existing model-selection flow-* code`; booked as GUESS ledger item 4. |

RDR passage: C1's last clause — "Shared refusal classes (model
selection, artifact binding, writer routing, read-back) reuse the
existing `flow-*` codes unchanged." It names model selection as a
shared class but fixes no group for any of the four. Same root gap as
D1. run-3 is the only run that enumerated it and the only run that
noticed the group is unstated.

### D8 (2-way) — `--flow` operand type and `--as` value set — `0019:C1`

| Run | Rendering |
| --- | --- |
| run-1 | `(--flow <id> \| --model <path>)`, `[--as json\|text]` |
| run-2 | Synopsis shows `--model <path>` only (prose says "the shared `--flow`/`--model` flags"), `[--as json]` |
| run-3 | `(--flow <path> \| --model <path>)`, `[--as json\|text]` |

RDR passage: C1 says only "the shared selection flags
(`--flow`|`--model`)" and `D-wire-byte-format` says the payload "rides
the 0005:C1 envelope". `§illustrative-code` shows `--model flow.toml
--artifact state=state.json --as json` — one example, `--model` only,
`--as json` only. Nothing in 0019 says what `--flow` takes or what
`--as` admits; both are 0005's to fix. Low severity — it lands on the
0005 citation, not on new 0019 surface — but the three runs produced
three different synopses for the verb's own CLI contract.

### D9 (2-way) — Empty-scalar seeding in the data model — `0019:F5`

| Run | Rendering |
| --- | --- |
| run-1 | Present in §3.1: `[initial] note = ""` → `{"note":""}`, reads back `Absent:false`, distinct from cleared; notes the argv asymmetry is routed to RDR 0002. |
| run-2 | **Absent.** Mentions empty scalar only in passing as one of the three loader-only admissions `canonicalValue` would refuse; never states what init-state persists for one. |
| run-3 | Present, one line: "An empty scalar `""` seeds and reads back PRESENT (F5)." |

RDR passage: F5 decides it explicitly — "This RDR takes that as the
DECIDED behavior for the seed route" — and the `disposition` table
carries the row `Empty-scalar [initial] value | 0 | the key, as "" |
seeded-all; reads back present, distinct from cleared (F5)`. The
decision is in the Failure Modes section and a Mini-check table, NOT
in C1. run-2, which explicitly declined to widen into "Trade-offs prose
beyond what A1-A8, BR1-5 already restate", missed it for that reason.

**Rewrite candidate**: this is a DECIDED behavior living outside the
Normative Contracts subsection. A reader who reads C1 alone does not
learn it.

### D10 (2-way) — Whether a sealed store can produce an exit-3 arm — `0019:C1` / `0019:S9`

| Run | Rendering |
| --- | --- |
| run-1 | Sealed store is purely the no-op arm; the exit-3 row it lists is for "read-back incomplete (unreachable locator)" during a WRITE, not for the sealed-store predicate. |
| run-2 | Renders a hybrid row: "Read-back reports a requested key UNREADABLE (e.g. sealed artifact) | inherited `set-state` semantics (exit 3 class) — but note: a sealed, all-cleared store is itself a no-op success (exit 0) at the predicate level, since the seal makes the store non-empty *before* any write is attempted". |
| run-3 | Sealed store is purely the no-op arm; no exit-3 row for it. |

RDR passage: C1 says a sealed artifact "is already exit 3 unreadable
to every reader — the read accessor reports every requested key
UNREADABLE rather than absent (`0004:C7`)" while simultaneously
requiring the predicate to see it as a one-key NON-EMPTY store and
decline. S9 fixes the outcome as no-op success. The tension the record
leaves open — how the emptiness probe gets a key COUNT from a reader
that short-circuits and reports everything unreadable — is exactly
A6's undecided carrier. run-2 is the only run that surfaced the
tension; the other two resolved it silently.

### D11 (2-way) — `Model.Writers` element type — `0019:C1`

| Run | Rendering |
| --- | --- |
| run-1 | Not stated (lists `Model.Initial`, `TagValue`, `resolve.Tag`, `Definition.Binding`, `IsDecisionTable`). |
| run-2 | Names `Model.Writers` as "existing per-role writer declarations" without a type. |
| run-3 | `table.Model.Writers []accessor.Definition // GUESS on element type; per-writer Role`. |

RDR passage: C1 cites `internal/table/model.go::Model.Writers` and
notes `::runFlowSetState` "already iterates per writer, checking
`req.artifacts[role]` for each" — the field is named but never typed.
Minor; the reconstruction is not blocked by it.

---

## 2. GUESS clusters (contracts 2+ runs marked GUESS)

Ordered by cluster size, then by how load-bearing the contract is.

### G1 — Payload field names (3/3 runs, all explicit) — `0019:D-wire-byte-format`
run-1: "every payload field NAME below is a GUESS by the record's own
admission". run-2: "explicitly deferred by the record itself to
Resolve/Pre-Lock". run-3: GUESS ledger item 1. All three cite
`D-wire-byte-format`'s "Field names deferred to Resolve/Pre-Lock, owned
here" verbatim. **Declared, not a defect — but see D3: the runs also
diverged on what the payload CARRIES (values vs. key names, `absent` on
the seeded arm), which is content, not naming, and which C1 arguably
does own.**

### G2 — Emptiness probe carrier and signature (3/3) — `0019:C1` / A6
run-1 §1.4 "GUESS: one of C1's three named options". run-2 "genuinely
undecided in the record; I did not invent a specific API for it".
run-3 GUESS ledger item 2. Contract-acknowledged (A6 Pending). Not a
C1 rewrite; a Resolve obligation.

### G3 — Internal helper names and signatures (3/3) — `0019:C1`
run-1: `allBoundStoresEmpty`, `checkFileBackedCarrier`, `canonicalSeed`
— all "GUESS: name". run-2: `isStoreEmpty`/`allBoundArtifactsEmpty`,
`checkAccessorCarrier`, `encodeSeed` — all "unnamed in the record".
run-3: `initCarrierCheck`, `storeIsEmpty`, `planInitial` — "Names are
GUESSES". Three runs, nine names, zero overlap. The RDR names
`newFlowInitStateCmd` (Step 1), `groupByWriter`, `writerFor`,
`canonicalSet`, `IsDecisionTable`, `sealedKey` — real shipped symbols —
but nothing for the three new helpers. Expected for an unimplemented
verb; NOT a rewrite candidate (the RDR should not fix private helper
names).

### G4 — Literal spellings of the two new exit-2 codes (3/3) — `0019:C1`
All three correctly render `flow-init-class-unsupported` and
`flow-init-carrier-unsupported` AND all three correctly flag the
spellings as non-normative on the `0028:C1.3::codeWriteEditRefused`
carve-out, noting no test may pin the string. **Unanimous and
correct** — C1's carve-out paragraph is doing its job. Listed here only
because it is a 3/3 GUESS-marked contract; it needs no rewrite.

### G5 — Go/CLI framework identity (2/3) — `0019:C1` / `§step-1-verb-skeleton`
run-1: "`*cobra.Command` is a GUESS — the record never names the CLI
framework." run-3: GUESS ledger item 6, "cobra as the command framework
(inferred from 'beside the four shipped verbs')". run-2 does not
render a Go constructor signature at all. Both guessing runs guessed
cobra and both flagged it. Environmental, not contractual.

### G6 — Exit group of the shared refusal classes (2/3) — `0019:C1`
run-3 marks the read-back group (item 3) and the model-selection group
(item 4) as GUESSes. run-1 marks the read-back group implicitly by
writing "non-zero (2)" with the parenthetical hedge and by listing a
second exit-3 arm the record does not carry. run-2 asserts 2 without
hedging. **This is the strongest rewrite candidate in the set**: it is
the same gap as D1 and D7, it was reached independently by two runs,
and unlike G1/G2 the record does NOT declare it as deferred — C1 says
these classes "reuse the existing `flow-*` codes unchanged", which
reads as settled while fixing no group.

### G7 — In-memory `store` type (2/3) — `0019:C1`
run-2: "`map[string]string` — inferred from described behavior
(presence/absence, verbatim values), not a cited type declaration."
run-1 renders the same shape as "a flat JSON object, key → canonical
string value" and explicitly scopes it: "the rest of the format is a
`flowbind` IMPL-DECISION, explicitly not contract". run-3 renders the
JSON shape without typing the Go map. C1 and S9 both say only the two
format properties (key PRESENCE distinguishes empty-set from cleared;
values stored VERBATIM) are contract — the runs agree on that scoping.
Correctly handled; no rewrite.

---

## 3. Agreement

All three runs rendered these identically, confirming the RDR is
determinate there — not findings: the verb is `flow init-state` in the
`flow` group with `--flow`/`--model` selection, repeatable `--artifact
role=path`, and **no write grammar** (`0019:C1` CARRIER, `0019:D-naming`);
the seed encoder is kind-dispatched — `canonicalSet` for `kind="set"`,
`members[0]` verbatim for every scalar kind — via neither
`::canonicalValue` nor a re-conform pass (`0019:C1`); the empty-store
predicate is the **ALL** quantifier over every bound artifact counting
**STORE** keys including `::sealedKey`, never ANY, never per-artifact,
never per-key-absent (`0019:C1`, `0019:D-selection-predicate`, `0019:S9`,
`0019:S11`); the class arm keys on `table.IsDecisionTable` ALONE and
never on `len(owned)` (`0019:C1`); the carrier gate covers BOTH
capabilities, admits only `*flowbind.Writer` (POINTER) and
`flowbind.Reader` (VALUE), runs after registry construction and before
any accessor invocation, and PREEMPTS the `--allow-commands` refusal
(`0019:C1`, `0019:S10`, A7/A8); the `≠1 writer` arm refuses as
`flow-model-invalid` at LOAD and never reaches the verb (`0019:C1`,
`0019:S5`); a non-empty store — torn, post-clear, sealed, or fully
seeded — is a NO-OP SUCCESS at exit 0 with zero writes whose payload
lists the absent `[initial]` keys and which never repairs a torn seed
(`0019:C1`, `0019:F2`, `0019:S12`); both new codes are exit group 2,
mutually distinct, spellings non-normative (`0019:C1`); execution is
the `runFlowSetState` shape — `groupByWriter` → per-writer apply →
read-back as the only commit check, with no cross-writer atomicity
(`0019:C1`, `§approach`, `§step-3-apply-and-report`).

---

## Landing summary

| Finding | Runs split | RDR element |
| --- | --- | --- |
| D1 read-back mismatch exit group | 3-way | `0019:C1` |
| D2 role-binding vs. emptiness order | 3-way | `0019:C1` (vs. `trace`, `§step-2`) |
| D3 payload shape/content | 3-way | `0019:C1`, `0019:D-wire-byte-format` |
| D4 emptiness probe signature | 3-way | `0019:C1` / A6 (declared) |
| D5 `--allow-commands` on flag set | 2-way, model boundary | `0019:C1` |
| D6 class vs. carrier order | 2-way, model boundary | `0019:C1` |
| D7 model-selection exit group | 2-way, model boundary | `0019:C1` |
| D8 `--flow` operand, `--as` values | 2-way | `0019:C1` (0005 citation) |
| D9 empty-scalar seeding | 2-way | `0019:F5` → should land in `0019:C1` |
| D10 sealed store vs. exit-3 | 2-way | `0019:C1`, `0019:S9` |
| D11 `Model.Writers` element type | 2-way | `0019:C1` |
| G6 shared-class exit groups | 2/3 GUESS | `0019:C1` |
| G1, G2, G3, G4, G5, G7 | 2–3/3 GUESS | declared silences / non-contractual |

Four 3-way disagreements, seven 2-way, seven GUESS clusters. Three of
the eleven disagreements split on the model boundary (D5, D6, D7) and
in all three the alt-model run (run-3) is the one that surfaced an
under-determination the other two resolved silently. Nine of the eleven
land on `0019:C1`; the tenth (D9) is a decided behavior living outside
the Normative Contracts subsection.
