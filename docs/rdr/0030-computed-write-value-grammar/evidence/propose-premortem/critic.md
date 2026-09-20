Model: claude-fable-5-1

# Premortem critic — computed write value grammar (`{ step = n }` load-time expansion)

Worked from the brief alone. No repository, record, or documentation was read.

## Findings ledger

| ID | passage/claim | Failure mode | Symptom user sees | Origin |
|---|---|---|---|---|
| P-1 | (b) "takes the current value as its row-identity suffix"; claim 1 | Derived row ids (`retry#2`) leak into every user-facing payload (resolve, next preview, set-state, export) as a new untiered identity surface; consumers keyed on the authored id stop matching | Automation/runbook that matched `row == "retry"` silently never fires again; no lint, no refusal, no warning | hindsight |
| P-2 | claim 1 "every downstream consumer ... is unchanged" (dump) | Dump either emits expanded rows whose `#` ids the loader refuses on reload, or emits the source and cannot name the rows lint/export just reported | `intrastate dump \| intrastate load` fails identity validation, or a lint finding names `retry#3` and dump shows no such row | negation |
| P-3 | claim 1 (kernel resolve unchanged) | Synthesized `eq` atom carries a Go-native int; authored atoms and runtime state values flow through TOML/JSON decoding (int64 vs float64 vs string); lint compares generated-to-generated and passes, kernel compares generated-to-decoded and misses | Lint green; every runtime resolve on the stepped rule refuses with "no row matched" | negation |
| P-4 | claim 8 "at most one atom per key ... in each block"; claim 2 | A `guard.all` int range is written as `gte = 1, lt = 5` — two atoms on one key in one block. A per-key single-predicate filter picks one operator; which one depends on builder iteration order | Either a bound refusal on a rule the author explicitly guarded below the bound, or phantom rows for cells the range excludes | refutation |
| P-5 | (b) "guard it below the bound"; problem statement "say it once" | The cap now lives in two literals: `max = 5` on the tag and `lt = 5` in the guard. Raising `max` without touching the guard leaves the ladder capped at the old value and, because the escalation row (`gte = 5`) claims the new cells, coverage is complete | Author raises the cap to 8; lint green; retries still escalate after 5; nobody is told | negation |
| P-6 | claim 3 "never changes which rows overlap" | Overlap and dead-row diagnostics now name expanded ids; an `unless` atom on the stepped tag yields a contradictory row (`guard.all eq = 4` + `unless eq = 4`) the author never wrote | Lint reports overlap/dead row on `retry#4`; author greps the file, finds no such row | negation |
| P-7 | claim 4 "no loader, normalizer, dump, or export sorts or canonicalizes" domain order | Order becomes load-bearing with nothing marking it so. A reviewer alphabetises `domain = [large, medium, small]`; membership semantics unchanged, step direction inverted. Where the guard never admits the new last member, no refusal fires | Escalation de-escalates; lint green | negation |
| P-8 | claim 5 "suffix cannot collide" | (i) A rule expanded on two tags (`in` on X, `step` on Y) mints `rule#a#1` or `rule#1#a` depending on pass order — nondeterministic ids, unstable fingerprint. (ii) Negative int suffix `rule#-1` and enum members containing chars the identity grammar rejects mint ids that fail the loader's own validator, or bypass it | Fingerprint differs between two loads of one file; or "invalid identity" refusal on an id the author never typed | refutation |
| P-9 | claim 6 "bounded by the domain size ... same ceiling" | Expansion is a loader step; the cardinality ceiling is a lint step. `int` with `min = 0, max = 10_000_000` plus a step rule allocates ten million rows before lint can refuse. Two stepped tags in one rule multiply | `load` hangs or OOMs on a model lint would have refused | negation |
| P-10 | claim 7 "reuse the in-atom expansion mechanism" | Expanded row keeps the authored `guard.all in = [...]` AND gains `eq = v` on the same key. If the normalizer/fingerprint keys atoms by `(block, tag)` the `eq` is clobbered; every expanded row has the same guard | Kernel refuses "N rows matched" on every resolve of the stepped rule; lint may or may not have flagged an overlap storm | refutation |
| P-11 | claim 9 "can only over-admit ... author can fix by moving the bound into guard.all" | An `unless` conjunction across two tags (`unless attempt eq 4 AND outcome eq ok`) has no positive `guard.all` equivalent; the only fix is splitting the rule. A hand-unrolled table that simply had no row for that cell loaded fine; the computed form refuses | Author is told to "move the bound" and cannot | negation |
| P-12 | claim 10 "admitted only on the write-block path" | Shared value decoder: initial-state, predicate literal, `set-state` args, and `scalar`-kind opaque values all see an inline table. A `scalar` tag that today legitimately stores `{ step = 1 }` as an opaque literal changes meaning; a `set` tag whose write grammar is already an inline table (`{ add = [...] }`) now has a colliding key | Existing model refuses to load after upgrade, or a misleading "computed write not allowed on initial state" error on a path that never asked for computation | negation |
| P-13 | (b) refuses "past `max`, past the last enum member" — no rule for the negative direction | `{ step = -1 }` at `min` / first enum member is unspecified. Index arithmetic yields `-1`: Go slice panic at load, or modulo wrap to the last member | `load` panics, or de-escalation from `small` lands on `large` | refutation |
| P-14 | (b) "non-zero int" — the only value rule on `step` | No rule for `step = 1.0` (TOML float), `step` larger than the domain, or `max + step` overflowing int64 (`max = 9223372036854775807`): overflow wraps negative, passes the `> max` check, and a garbage literal is written | Load succeeds; the accessor writes a negative attempt count; read-back equality passes because it reads back the garbage | refutation |
| P-15 | (b) new refusal code, `step` reserved key, derived ids and any `expanded_from` field in export | New machine-readable surfaces added without a stability tier; the write inline-table becomes a closed one-key vocabulary with no rule for unknown keys | A downstream consumer pins to `expanded_from` or `#`-suffix shape with no contract saying it may | refutation |
| P-16 | R1 rejection: "every consumer ... would have to learn a second value shape" | The chosen approach also touches every one of the ten sites — each now emits or consumes derived ids and must decide whether to name the source rule; the bound check does become a per-cell analysis, just executed at load | Reviewer accepts R1's rejection and then discovers ten call-site diffs anyway | negation |
| P-17 | R3 rejection: "the generated file becomes the artifact the author reviews" | Coverage gaps, overlaps, dead rows, reach edges, and export nodes are all reported on the expanded rows — the generated artifact is what the author reviews, it just is not on disk and cannot be diffed | Author asks "which of my rows is `retry#3`" and the answer is "none of them" | negation |
| P-18 | R4 rejection: refuse rather than saturate | Refusal is a load-time hard failure on a data edit: lowering `max` or removing an enum member turns a running model into one that will not load at next start. Single-member domains (`min == max`, one enum member) refuse every step | Service fails to start after a config edit that was reviewed as "narrowing a domain" | negation |
| P-19 | (b) "the mechanism the loader already uses to expand an `in` match atom" | Two expansion passes over the same tag (`match attempt in [1,2,3]` + `step` on attempt) compose differently depending on order: nested suffix `rule#1#1`, or 3x3 rows of which six are contradictory | Nine rows in export for a rule the author expected to yield three; six dead | negation |
| P-20 | R2 rejection | `{ step = n }` is a one-operator expression language. Without a closed-grammar rule (`{ step = 1, wrap = true }` refused as unknown key, not ignored), the next feature request extends it by accretion | An unknown key is silently ignored and the author believes `wrap = true` did something | negation |

## 1. Prospective hindsight — the failure as accomplished fact

The grammar shipped in the release after the record locked. The lead model was a retry/escalation table: `attempt` as `int` with `min = 0, max = 5`, `tier` as `enum` with `domain = ["small", "medium", "large"]`. The author replaced twenty-two literal rows with two rules — `retry` carrying `attempt = { step = 1 }` under `guard.all attempt = { lt = 5 }`, and `escalate` carrying `tier = { step = 1 }` under `guard.all tier = { in = ["small", "medium"] }`. Lint was green: coverage complete, no overlaps, reach unchanged, export unchanged in field vocabulary.

Three things then happened, none of them refused and none of them warned.

First, an operations job that consumed the resolve payload matched the selected row id against a list of "retry-class" rows by their authored names. The rows now firing were `retry#0` … `retry#4`. The job's match failed, its fallback branch treated every retry as a fresh incident, and the on-call queue filled over a weekend. Nothing in the export, the help text, or the output contract had said that a row id in a payload might carry a suffix the author never typed. The existing `in`-expansion had the same property, but it had only ever appeared on rows whose author had written an explicit member list, so no consumer had encountered it.

Second, the cap was raised. The tag declaration went from `max = 5` to `max = 8`. The guard `lt = 5` was not touched — the author had been told the computed form let them "say it once", and they said it in the tag. Lint stayed green because `escalate`'s sibling `gte = 5` row claimed cells 5, 6, and 7. Retries continued to stop at five. The gap was noticed three weeks later from a dashboard.

Third, a reviewer tidying the model alphabetised the tier domain. Membership did not change; the step direction did. Because `escalate` guarded `in = ["small", "medium"]` and `small` was now last, load refused — but the reviewer read the refusal ("stepped value past last member at cell tier=small") as a bug in the guard, changed it to `in = ["large", "medium"]`, load succeeded, and escalation de-escalated in production.

The follow-up audit found the guard.all two-atom range case (P-4) latent in a second model that had not yet been promoted, and the int64 overflow (P-14) in a test fixture that used `max = math.MaxInt64` as "unbounded".

## 2. Obstacle negation

### Claim 1 — "semantically identical to hand-written; every downstream consumer unchanged"

Negated. Hand-written rows have authored ids; expanded rows have derived ids. Every consumer that emits, stores, or matches a row id — resolve payload, next preview, set-state, export, fingerprint, the idempotent-write advisory, dump — now emits an identity the author cannot find in the file (P-1). Dump in particular cannot be unchanged: either it emits the expanded rows (ids contain `#`, which the loader refuses on the way back in, so dump stops round-tripping) or it emits the source form (and every lint finding that names `retry#3` points at a row dump does not show) (P-2). Separately, the synthesized `eq` atom's value originates in Go from the domain enumeration, while authored atoms and runtime state pass through decoders; if the kernel compares by dynamic type, lint (generated vs generated) passes and the kernel (generated vs decoded) misses (P-3).

### Claim 2 — "guard it below the bound is sufficient; the bound cell is reported by coverage"

Negated twice. The bound must now be written in two places, `max` on the tag and `lt` in the guard, and nothing ties them; raising `max` while a sibling row claims the new cells leaves coverage complete and the ladder silently short (P-5). And "guard it below the bound" presumes the filter reads the guard correctly, which claim 8 does not deliver (P-4).

### Claim 3 — "eq atoms never change overlap beyond the hand-unrolled table"

Negated at the reporting layer. The set of overlapping cells may be the same, but the rows reported are not: every overlap and dead-row finding now names an expanded id. And the `unless`-ignored filter produces rows with self-contradictory guards (`eq = 4` and `unless eq = 4`) that a hand-unrolled author would never have written; a dead-row or unreachable-row lint flags them (P-6).

### Claim 4 — "domain order is stable and meaningful; nothing canonicalises it"

Negated on obligation, not on mechanism. Even if no code sorts it today, order is now load-bearing and nothing says so: no lint marks a domain as "stepped; order matters", no export field publishes the order as a contract, and a membership-preserving reorder is a semantic change with no diff at the rule site. Where the guard does not admit the new last member, no refusal fires (P-7). Also unstated: duplicate members in `domain` (refused, or deduped with which position kept?).

### Claim 5 — "the suffix cannot collide"

Negated on determinism and grammar rather than on `#`. A rule expanded on two tags mints a two-component suffix whose order depends on which pass runs first; unless the suffix is tag-qualified and ordered by a stated rule, ids and therefore fingerprints are unstable across loads (P-8). Negative ints (`rule#-1`) and enum members containing characters the identity grammar rejects mint ids the loader's own validator refuses — or, worse, the derived path skips validation the authored path enforces (Seed 2 class).

### Claim 6 — "row growth bounded; under the ceiling lint already enforces"

Negated on ordering. The ceiling is a lint-stage check on the guard product; expansion is a loader-stage allocation. A wide int domain allocates before the ceiling runs. Two stepped tags in one rule multiply, and combined with `in`-expansion on a third tag the per-rule count is a product the brief never bounds (P-9).

### Claim 7 — "the in-expansion mechanism can be reused"

Negated. In-expansion rewrites a match atom in place; step-expansion adds a guard atom while retaining the authored one. The expanded row now carries both `in` and `eq` on one key in `guard.all`. If any stage keys atoms by `(block, tag)` — normalizer map, fingerprint canonicalisation, overlap comparator — one atom clobbers the other and every expanded row has an identical guard; the kernel then sees N matching rows and refuses at every resolve (P-10). Composing the two passes on the same tag is also unspecified (P-19).

### Claim 8 — "at most one atom per key per block"

Refuted directly (Seed 1). A `guard.all` int range is written as two atoms on one key (`gte = 1, lt = 5`); the brief says guards admit all six operators. A filter that fetches "the" predicate for the key returns whichever the builder emits first, and behaviour differs by operator order: pick `lt` and cells below `gte` expand into contradictory rows; pick `gte` and the bound cell expands and refuses despite the author's `lt` (P-4). The filter must evaluate the conjunction of every positive atom on the key.

### Claim 9 — "ignoring unless is conservative; move the bound into guard.all"

Negated on the remedy. An `unless` that conjoins the stepped tag with another tag has no positive single-block equivalent; the author must split the rule. And "conservative" is only true relative to under-admission — relative to the hand-unrolled table it is strictly more refusing: a table with no row for cell 4 loaded yesterday; its computed equivalent refuses today (P-11).

### Claim 10 — "initial-state and predicate paths untouched"

Negated where the value decoder is shared. Initial-state, predicate literal, `set-state` argument parsing, and the `scalar` kind (which may already admit an inline table as an opaque literal) all see an inline table with a `step` key; the `set` kind's own inline-table write grammar, if one exists, now has a reserved key it did not choose (P-12). Unstated value rules on `step` itself — float, oversize, int64 overflow (P-14) — and the whole negative direction (P-13) are Seed 2 class: a token admitted with no value rule.

### R1 — symbolic write into the kernel

The rejection reason ("ten sites learn a second value shape") is symmetric: the chosen approach makes those same ten sites emit or consume derived identities, and each must decide whether to name the source rule alongside (P-16). The "bound check becomes interval analysis" objection is also symmetric — the chosen approach performs exactly that analysis, per cell, at load.

### R2 — expression language

The chosen grammar is a one-operator expression language. Its rejection reason ("defeats static coverage") is answered by expansion, but the grammar needs a closed-vocabulary rule so unknown keys refuse rather than pass (P-20).

### R3 — external generator

"The generated file becomes what the author reviews" is true of the chosen approach's diagnostics: every lint finding, reach edge, and export node is reported on generated rows the author cannot diff or open (P-17).

### R4 — saturate instead of refuse

Refusal converts a domain-narrowing data edit into a start-time hard failure of the running service; the idempotent-write advisory would have degraded gracefully. Single-member domains and `min == max` refuse every step with no useful authoring escape (P-18).

## 3. Consumer artifacts — what would have caught each at review

- **P-1** — Output-contract test `TestResolvePayload_RowIDIsAuthoredOrTiered`: load the retry model, resolve at `attempt = 2`, assert the payload names the row and that the published contract for the `row` field states whether a `#` suffix may appear. User journey "operator matches on row id" in the output contract doc.
- **P-2** — Round-trip test `TestDump_ReloadsIdentically`: `load → dump → load` on a model with a step rule; assert equal fingerprints and that every id lint reports appears in dump output (or that dump names the source rule for each).
- **P-3** — End-to-end test `TestStepRule_ResolveWritesReadBack`: set-state `attempt = 2` through the real accessor, resolve, assert `retry#2` selected and read-back `attempt == 3`; run the same under the JSON-state path if state is persisted as JSON.
- **P-4** — Normalizer test `TestStepFilter_RangeGuardIsConjunction`: rule with `guard.all attempt = { gte = 1, lt = 5 }`, `max = 5`; assert exactly rows for 1..4, no refusal, no row for 0. Second case swaps the authoring order of the two keys and asserts identical output.
- **P-5** — Lint advisory test `TestLint_StepGuardBoundDriftsFromMax`: `max = 8`, `guard lt = 5`, sibling `gte = 5`; assert an advisory names the stepped tag and both literals. User journey "raise the cap" in the record's worked example.
- **P-6** — Golden lint output `TestLint_UnlessOnSteppedTag_NamesSourceRule`: rule with `unless attempt eq 4`; assert any dead-row/overlap finding names the source rule and the cell, not only `retry#4`.
- **P-7** — Test `TestDomainReorder_ChangesFingerprintAndWarns`: reorder enum members without changing membership; assert fingerprint changes and a lint advisory says the domain is stepped and order-bearing. Export test asserts domain order is published.
- **P-8** — Determinism test `TestExpansion_TwoTags_StableIDs`: rule with `in` on X and `step` on Y loaded twice; assert identical id sets. Identity test `TestExpansion_MintedIDPassesIdentityGrammar` over negative ints and enum members with spaces, dots, and `-` prefix.
- **P-9** — Test `TestLoad_WideIntDomainRefusesBeforeExpansion`: `min = 0, max = 10_000_000` with a step rule; assert refusal within a memory budget (testing.AllocsPerRun / timeout).
- **P-10** — Normalizer test `TestExpansion_RetainsInAndEqOnOneKey`: `guard.all attempt = { in = [1,2,3] }` + step; assert each expanded row carries both atoms and that fingerprints of the three rows differ; kernel test asserts exactly one row matches at `attempt = 2`.
- **P-11** — Refusal-message test `TestStepRefusal_UnlessConjunctionRemedy`: `unless` across two tags; assert the refusal text names the cell and offers "split the rule", not "move the bound".
- **P-12** — Table test `TestInlineTable_RefusedOffWritePath` over initial-state, predicate literal, set-state args, `scalar` write, `set` write; each asserts a distinct code and that a pre-existing `scalar` opaque-table fixture still loads.
- **P-13** — Test `TestStep_NegativeAtLowerBoundRefuses` for `int` at `min` and enum at first member with `step = -1`; assert refusal, no panic, no wrap.
- **P-14** — Table test `TestStep_ValueRules`: `step = 1.0`, `step = 0`, `step = 100` on a 6-cell domain, `max = math.MaxInt64` with `step = 1`; each asserts refusal with a named code.
- **P-15** — Docs-check assertion that every new refusal code and export field is listed with a stability tier in the generated reference; test `TestWriteInlineTable_UnknownKeyRefused`.
- **P-16** — Review checklist item: diff touches to the ten named sites enumerated in the record, each stating "names source rule: yes/no".
- **P-17** — User journey "author reads a coverage-gap finding on an expanded row" documented with the exact finding text and how it maps back to the authored rule.
- **P-18** — Test `TestStep_SingleMemberDomainRefusesWithRemedy`; operational note in help text that narrowing a stepped domain is a load-breaking edit.
- **P-19** — Test `TestExpansion_MatchInPlusStepOnSameTag`: assert exactly three rows, ids stated in the record.
- **P-20** — Grammar test `TestWriteInlineTable_ClosedVocabulary`: `{ step = 1, wrap = true }` refuses.

## 4. Refutation targets — the seed classes, found here

- **Seed 1 class (multiple atoms per key; non-monotonic filter)** — present at P-4 (guard.all range is two atoms on one key) and P-10 (expansion adds a second atom on the key it filtered). The brief's claim 8 is the same false license Seed 1 carried. The fix is a conjunction filter over every positive atom on the key, and an explicit statement of what the expanded row carries.
- **Seed 2 class (token admitted with no value rule)** — present at P-13 (negative step direction unspecified), P-14 (`step` float / overflow / oversize), and P-8 (minted ids bypass the identity grammar the authored path enforces).
- **Seed 3 class (new surfaces, no tier; closed list a peer forbids relying on)** — present at P-15 (refusal code, reserved key, derived-id shape, any `expanded_from` field) and P-1 (derived ids in payloads with no contract statement).

## Verdict reasoning

No finding forces a switch to R1–R4: each is answerable inside the chosen design by stating a rule the brief omits (conjunction filter, tag-qualified deterministic suffix, lower-bound and overflow refusals, closed grammar, tiering, dump-from-source with source-rule naming in diagnostics, order-bearing lint, guard/max drift advisory, ceiling-before-expansion). The approach passes with those mitigations folded in.
