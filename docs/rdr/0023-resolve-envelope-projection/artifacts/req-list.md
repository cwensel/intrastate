# REQ List — RDR 0023 Resolve-envelope projection opt-out

Phase 0 audit artifact. Every clause below is a testable obligation drawn from
`docs/rdr/0023-resolve-envelope-projection.md`. Quotes for fenced elements are
verbatim — copied from the projector (`rdr inspect --select <id>`), never
transcribed by hand; quotes for testable prose outside the fences are copied
from the record at the cited line.

Element ids (`0023:C1`, `0023:C2`, `0023:MVV`, `0023:S1`…`S5`, `0023:A1`…`A9`,
`0023:D-*`, `0023:F1`…`F7`, `0023:G-cross-cutting`) are carried wherever a REQ
derives from a labelled element, so a later stage can trace the REQ back to its
contract.

**Why the REQ count far exceeds `counts.elements.C = 2`.** This record has
exactly two `normative` fences (`0023:C1` at 182 lines, `0023:C2` at 135), and
between them they carry roughly two dozen independent MUST-clauses — flag
registration, the projected key set, absent-not-null, byte-identity of carried
values, key order, the report-only differential with four sub-assertions and two
pinned measurands, the projection-site rule, the text-subset rule, the
whole-tree non-registration walk with its two preconditions, and C2's declared
partition, constrained carrier, separate projection function, and always-keep
core. Everything else testable is PROSE and is mined below: the `Overrides`
block (the narrowing of `0005:A6`'s pinned envelope), `Approach`,
`Technical Design`, the `Source-authority census`, the `Disposition` grid, the
`Oracle discriminability` grid, `Existing Infrastructure Audit`,
`Consequences`, `Failure Modes`, `Prerequisites`, `MVV` steps 1–6 and its pass
bar, `Phases 1–3`, `Testing Strategy` S1–S5, the `Desk trace`,
`Performance Expectations`, and `Cross-Cutting Concerns`.

The nine Critical Assumptions (`A1`–`A9`), four Load-Bearing Decisions
(`D-identity`, `D-wire-byte-format`, `D-naming`, `D-selection-predicate`), seven
Failure Modes (`F1`–`F7`), five Testing Strategy scenarios (`S1`–`S5`), the
Joint-check (`JC1`), the five Briefly-Rejected entries (`BR1`–`BR5`) and the
three Alternatives (`ALT1`–`ALT3`) were read as context for what the contracts
must hold. `ALT1`–`ALT3` and `BR1`–`BR4` are rationale, not obligation, and mint
no REQ; `BR5` is mined as a negative REQ (no typed refusal code) because C1
states the same fence. `A1`–`A9` are Verified and mint no implementation REQ of
their own EXCEPT where the record explicitly carries a CONSTRUCTION obligation
out of an assumption into the build (`A5`'s total walker → REQ-38/39/40;
`A9`'s empty-container oracle → REQ-52) — those are named in
`Implementation Plan / Prerequisites` as "oracles to write, not assumptions to
settle".

Source sections are abbreviated:

- `NC` = Proposed Solution / Technical Design / Normative Contracts (fenced, `C1`/`C2`)
- `OVR` = Metadata / Overrides
- `AP` = Proposed Solution / Approach
- `TD` = Proposed Solution / Technical Design (unfenced prose)
- `CEN` = Technical Design / Source-authority census
- `DIS` = Technical Design / Disposition (input-class grid)
- `OD` = Technical Design / Oracle discriminability
- `EIA` = Proposed Solution / Existing Infrastructure Audit
- `LBD` = Technical Design / Load-Bearing Decisions
- `CONS` = Trade-offs / Consequences
- `RM` = Trade-offs / Risks and Mitigations
- `FM` = Trade-offs / Failure Modes
- `PRE` = Implementation Plan / Prerequisites
- `MVV` = Implementation Plan / Minimum Viable Validation
- `PH` = Implementation Plan / Phases 1–3
- `TS` = Validation / Testing Strategy (S1–S5)
- `DT` = Validation / Desk trace
- `PERF` = Validation / Performance Expectations
- `CC` = Finalization Gate / Cross-Cutting Concerns

**Standing note on the override.** This RDR overrides ONE thing in the locked
predecessor 0005: `0005:A6`'s pinned minimum success payload fields for
`resolve`, and only in the new opt-in mode. Every other 0005 obligation on
`flow resolve` — one terminal envelope, refusal codes, exit groups, gate
ordering, both-modes agreement — is unchanged and rides in BOTH widths. REQs
below that restate a predecessor obligation are marked **[0005-carried]**,
**[0010-carried]** or **[0011-carried]**: they are obligations to NOT break, not
licence to re-implement the predecessor. Nothing in 0011 is overridden.

**Standing note on the doctrine home.** `0023:C2` cites JDR 0002 §D1 as
normative THERE and does not restate it. REQs derived from the doctrine are
marked **[JDR-0002:D1]** and cite the JDR line; the verb's instance
(concrete field lists, flag semantics, oracle spellings) is C2's and is the
enforceable surface here.

---

## Flag surface and registration (`0023:C1`)

- [REQ-1] "flow resolve MUST accept a boolean flag --plan-only, default false, long form only (no shorthand)." — (NC, `0023:C1`)

- [REQ-2] "--plan-only MUST NOT be accepted by any other command — satisfied by NON-REGISTRATION (registered on resolve only, never the shared registrars, the flow group's own or persistent set, or the root's persistent set)" — (NC, `0023:C1`)

- [REQ-3] The non-registration is "pinned by a STRUCTURAL oracle that walks the WHOLE command tree from the root and asserts the set of commands registering plan-only is exactly {flow resolve} — closed over verbs added later" — (NC, `0023:C1`)

- [REQ-4] The structural oracle is "vacuity-guarded by requiring the four flow verbs to exist." — (NC, `0023:C1`)

- [REQ-5] "WHOLE means total: the walk descends every child of the root with no name-based skip and no Hidden gate, INCLUDING the auto-generated help and completion commands." — (NC, `0023:C1`)

- [REQ-6] "It therefore MUST NOT reuse internal/cli/help_all.go::walkCommandTree, which skips children named help or completion before recursing and whose docs.go callers further gate on Hidden — that walker asserts a strictly weaker negative than this clause requires." — (NC, `0023:C1`) — the negative is the assertable half: the shipped walker symbol is not called by the new oracle.

- [REQ-7] "The oracle MUST therefore materialize both auto-generated commands on the root before walking (calling cobra's own initializers, never hand-constructing a stand-in), and MUST assert their presence in the walked set as a precondition." — (NC, `0023:C1`)

- [REQ-8] "The control asserting the walk reaches help is NOT sufficient on its own: help is present by default, so it cannot distinguish a total walker from an untested tree — completion is the discriminating case and the control MUST name it." — (NC, `0023:C1`)

- [REQ-9] "if a cobra upgrade renames, relocates or stops materializing the command, the precondition fails LOUDLY and the oracle goes red, instead of the walk quietly covering a smaller tree and passing vacuously." — (NC, `0023:C1`)

- [REQ-10] "The behavioural command-error run is corroboration, never the assertion." — (NC, `0023:C1`) — a behavioural sibling-verb run may not stand in for the structural walk.

- [REQ-11] "This contract mints no new refusal code and no new exit group." — (NC, `0023:C1`); restated at `DIS` line 1045 ("The flag mints no refusal code and no exit group (C1)") and as the disposition of `0023:BR5`.

- [REQ-12] "`--plan-only` on `next`/`read-state`/`set-state` fails the parse — `command-error`, exit 2, flag name only in pflag's free text (the `0011:A14` shared-bucket limit, accepted again here)." — (FM, `0023:F1`) — **[0011-carried]**; the `DIS` grid row states the same: "`--plan-only` on `next`/`read-state`/`set-state` | 2 | `command-error` shared bucket".

- [REQ-13] "The walker yields the WALKED COMMAND SET — every command reached, unfiltered — and the registrant filter is applied by the caller, not inside the walk." — (TS, `0023:S4`)

- [REQ-14] The new walker "is also a distinct symbol from `internal/cli/help_all.go::walkCommandTree`, which C1 forbids reusing — that one is a visitor (`func(root *cobra.Command, fn func(*cobra.Command))`, no return) and skips `help`/`completion`, so neither its name nor its shape carries over." — (TS, `0023:S4`)

- [REQ-15] "This scenario also carries the POSITIVE registration-shape assertion on the single registrant — `plan-only` is boolean, defaults false, no shorthand (C1, `D-naming`) — on the `--all` precedent … the shape is otherwise stated twice in this RDR and asserted nowhere." — (TS, `0023:S4`)

- [REQ-16] "the flag is `--plan-only`, boolean, long form only." — (LBD, `0023:D-naming`) — the spelling is normative; `--select`, `--json <fields>`/`--fields`, `--quiet`, `--no-echo`, `--short`/`--brief`, and an inverted `--explain` default are all rejected spellings.

## Projected payload: partition, key set, absence (`0023:C1`)

- [REQ-17] "Under --plan-only a SUCCESS payload MUST carry exactly the plan group — rule, gates, emit, next, writes, clear, escaped, escape_class, revision, plus every field a successor assigns to the plan group under C2" — (NC, `0023:C1`)

- [REQ-18] The same payload "MUST NOT carry the echo group: model, observed, owned, readers, outcome." — (NC, `0023:C1`)

- [REQ-19] "An omitted key is ABSENT, never null and never an empty placeholder." — (NC, `0023:C1`); restated at `LBD` `0023:D-wire-byte-format`: "Absent means absent — never `null`, `{}`, or `\"\"` stand-ins for a projected key."

- [REQ-20] "Every carried field's rendering MUST be byte-identical to its default-mode rendering (same encoder, HTML escaping disabled)" — (NC, `0023:C1`)

- [REQ-21] "the carried keys keep their default-mode relative order." — (NC, `0023:C1`)

- [REQ-22] "The projected TOP-LEVEL KEY SET is itself normative: an oracle MUST assert it as an explicit literal list, so a later payload field reaches the projected width only by a conscious C2 assignment, never by falling through an omit-list." — (NC, `0023:C1`)

- [REQ-23] The literal list is "exactly `revision`, `rule`, `gates`, `emit`, `next`, `writes`, `clear`, `escaped`, plus `escape_class` exactly when the same request's DEFAULT output carries it (the producer's own presence rule, C1 — not an unconditional \"when escaped\")" — (TS, `0023:S2`)

- [REQ-24] "Presence-rule fields keep their own contracts inside the plan group, unchanged and not restated here: emit stays present as {} (0010:C4), and escape_class keeps whatever presence rule its producer already has (0005:A-3 — omitted when unescaped, and also when an escaped row's class is unprobeable, `flow_resolve.go::escapeClassOf`)." — (NC, `0023:C1`) — **[0010-carried]**, **[0005-carried]**

- [REQ-25] "This projection neither widens nor narrows those rules: a field the default mode omits is omitted under the flag for the same reason, so an oracle MUST assert presence-rule fields by comparison against the SAME run's default output, never against an unconditional literal." — (NC, `0023:C1`)

- [REQ-26] "projection is key deletion, never re-encoding: the projected payload is the default payload minus the echo keys, byte-for-byte on every surviving field, keys in surviving declaration order, same non-HTML-escaping encoder." — (LBD, `0023:D-wire-byte-format`)

- [REQ-27] "Default-mode output is byte-identical to today's: the request echo is the correct default and this RDR does not change it." — (NC, `0023:C1`)

- [REQ-28] The `Overrides` narrowing: "under the opt-in `--plan-only` projection (C1 below) the echo-group fields `model`, `observed`, `owned`, `readers`, `outcome` are absent from the success payload (`revision` stays — plan-side identity slot, C2). Default-mode output is byte-identical to today's; every other 0005 obligation on `flow resolve` — one terminal envelope, refusal codes, exit groups, gate ordering, both-modes agreement — is unchanged." — (OVR, lines 105–115) — **[0005-carried]** for the unchanged half.

- [REQ-29] "No clause of 0011 is overridden: `0011:C2`'s fence is `--all`'s non-registration on the other three verbs, which this RDR leaves intact, and its `--select` rejected-spelling pin is on `flow next`" — (OVR, lines 115–118) — **[0011-carried]**; the shipped 0011 oracles must stay green untouched.

- [REQ-30] "`--all` probes are exact-name `Lookup(\"all\")`; a resolve-local `plan-only` is invisible to them." — (DT, step 5 witness; `0023:A5`) — **[0011-carried]**

## Report-only differential (`0023:C1`, `0023:S1`)

- [REQ-31] "--plan-only is REPORT-ONLY, and report-only is ORACLE-ENFORCED, not aspirational: it MUST NOT change what is decided or how the run fails" — (NC, `0023:C1`)

- [REQ-32] "the build MUST carry a differential oracle asserting, over the same request with and without the flag: identical exit codes" — (NC, `0023:C1`)

- [REQ-33] The same differential asserts "byte-identical refusal envelopes (CLIError, findings — never projected, and carrying no echo group to project)" — (NC, `0023:C1`); `0023:A6` grounds it: "`respond.Fail` renders `clierr.CLIError` (+ `findings`) with no request-echo member, so `--plan-only` cannot change a single refusal byte".

- [REQ-34] The same differential asserts "the projected success payload a strict key-subset of the default payload with byte-identical values on every carried key" — (NC, `0023:C1`)

- [REQ-35] The same differential asserts "an identical invoked-reader set" — (NC, `0023:C1`)

- [REQ-36] "the projected encoding STRICTLY SHORTER than the default. The width assertion is not implied by the subset one — a later plan-group field can grow the projected payload past today's full width with the key-set and partition oracles still green — so the flag's whole reason for existing is itself oracle-enforced." — (NC, `0023:C1`)

- [REQ-37] "A later change that skips work whose only consumer is a projected-away field breaches this clause." — (NC, `0023:C1`) — the reader pass is the reachable surface (`CEN` lines 1000–1013).

- [REQ-38] "STRICTLY SHORTER is measured on the FULL EMITTED LINE — the complete {\"type\":\"ok\",\"data\":{…}} NDJSON record as written, excluding the trailing newline — not on the .data payload alone, and it is asserted in BOTH output modes." — (NC, `0023:C1`)

- [REQ-39] "In text mode the measurand is the total rendered byte count of the emitted lines; a text projection that drops no bytes fails this clause even though the S5 subset assertion would still pass, a subset being satisfied by the equal set." — (NC, `0023:C1`)

- [REQ-40] "The measurand is therefore OBSERVED READER EXECUTION, not a recomputed derivation. Comparing invokedReaders(model, outcome) across the two runs is NOT an admissible form of this assertion" — (NC, `0023:C1`)

- [REQ-41] "The oracle MUST instead observe which readers ACTUALLY RAN in each run — recording execution at the reader invocation site (a counting or recording seam around the reader pass, the run's own side effects, never a re-call of the planning function) — and assert the two observed sets are equal." — (NC, `0023:C1`)

- [REQ-42] "The invocation site the seam wraps is the per-reader `exec.Read` call inside `internal/cli/flow_exec.go::(flowRequest).runReaders`, the one pass that yields `readers` and `owned` together; the seam's FORM is left to the implementer, but its position is not" — (NC, `0023:C1`)

- [REQ-43] "The DEFAULT-mode run of the same request supplies the expected set. The discriminating property is that a projected run which skips the reader pass MUST turn this assertion red; an assertion that cannot distinguish that case does not satisfy this clause." — (NC, `0023:C1`)

- [REQ-44] Reading the reader set "from the projected payload" is forbidden and unwritable: "readers is an ECHO field this flag projects away, and the shipped helper reading it fails hard on absence (internal/cli/flow_fixtures_0011_test.go::readersOf)." — (NC, `0023:C1`)

- [REQ-45] "the same invocation ± the flag selects the same rule, runs the same gates, and refuses identically — only report width differs … an oracle asserts the ±-flag plan-group equality." — (LBD, `0023:D-identity`)

- [REQ-46] "because one call yields both fields, no implementation may satisfy the projection by suppressing `owned` and `readers` independently." — (CEN, lines 1012–1013)

## Projection site and flag-blindness (`0023:C1`, `0023:C2`)

- [REQ-47] "The projection SITE is normative, not an implementation-plan detail. It is applied to the verb-specific result BEFORE respond.OK, never in the respond gateway and never per output mode" — (NC, `0023:C1`) — **[JDR-0002:D1]** (the later-verb rule, JDR line 64–67).

- [REQ-48] "it sits on the SUCCESS path only, AFTER the last respond.Fail return" — (NC, `0023:C1`)

- [REQ-49] "the flag is read at exactly ONE lexical site, that projection branch — never passed as a parameter into payload assembly, and never consulted in gate evaluation, rule selection, or refusal construction." — (NC, `0023:C1`)

- [REQ-50] "A defensive \"if refusing, skip projection\" guard is FORBIDDEN: refusal flag-blindness (A6) must hold structurally, because a refusal returns before the projection is reachable" — (NC, `0023:C1`)

- [REQ-51] "The projection function is likewise SEPARATE from payload assembly: it takes the fully assembled payload and returns the projected one, and the flag is never a parameter to the assembly function." — (NC, `0023:C2`)

- [REQ-52] "Fusing the two — assembling conditionally under the flag — destroys both checkable properties this contract rests on … Assembly is flag-blind; projection is the only flag-aware step." — (NC, `0023:C2`)

- [REQ-53] "Nothing upstream of payload assembly reads the flag: readers, kernel call, gate run, and every refusal path are flag-blind. The respond gateway is untouched" — (TD, lines 601–604)

## Text mode (`0023:C1`, `0023:S5`)

- [REQ-54] "Text mode renders the projected result through the same generic payload renderer as every other success, so the two modes cannot disagree (0005:C1 held by construction)." — (NC, `0023:C1`) — **[0005-carried]**

- [REQ-55] "The projected text lines MUST be a subset of the default-mode text lines, byte-identical per line." — (NC, `0023:C1`)

- [REQ-56] "projected text lines a byte-identical, stable subset of the default-mode lines — set membership, not a subsequence: the wire keeps struct declaration order while text is sorted by the shipped `flatten` (`respond/text.go`), and no clause requires the two orders to agree." — (TS, `0023:S5`)

- [REQ-57] "Additionally the projected text is STRICTLY SHORTER in total rendered bytes (C1's width clause binds both modes); a subset assertion alone is satisfied by the equal set, so a zero-saving text projection would otherwise pass." — (TS, `0023:S5`)

- [REQ-58] "The two measurands are JOINTLY SUFFICIENT only because the subset is asserted per-line and byte-identically: every projected line must appear in the default set unchanged, so the width reduction can only come from DROPPED lines, never from a shortened carried value." — (TS, `0023:S5`)

- [REQ-59] Text mode is deterministic across repeated runs — "`--as=text` ± the flag, repeated runs" with a "stable subset" (TS, `0023:S5`); MVV step 3: "the projected lines are a byte-identical subset of step 1's text lines, stable across repeated runs." — (MVV, `0023:MVV`)

## Partition declaration and reflective oracle (`0023:C2`)

- [REQ-60] "Per JDR 0002 §D1's enforcement rule, this verb's reflective oracle MUST assert every field of the resolve success payload is assigned to exactly one group, so an unassigned new field is a test failure, not a silent default — the projection MUST NOT be implemented as a bare omit-list whose complement is \"whatever else exists\"." — (NC, `0023:C2`) — **[JDR-0002:D1]** (JDR lines 52–55).

- [REQ-61] "the ECHO group is the model reference, the observed tags (--tag, echoed unchanged), the assembled owned view and the invoked reader identities, and the requested outcome (--outcome, echoed unchanged)." — (NC, `0023:C2`)

- [REQ-62] "The PLAN group is rule identity, gate results, authored answers and their interpretations, planned next/writes/clear, the escape disposition, and revision." — (NC, `0023:C2`)

- [REQ-63] The per-field assignment is fixed by the census: `model`, `observed`, `owned`, `readers`, `outcome` = ECHO; `revision` = "PLAN (identity slot)"; `rule`, `gates`, `emit`, `next`, `writes`, `clear`, `escaped`, `escape_class` = PLAN. — (CEN, lines 984–998)

- [REQ-64] "The set is FOURTEEN fields — `reflect.TypeOf(resolvePayload{}).NumField()` is 14, the count `internal/cli/decision_table_0010_test.go` already pins" — (CEN, lines 974–976) — **[0010-carried]**

- [REQ-65] "The oracle reflects the payload struct against an assignment source that is INDEPENDENT of the projection code … an oracle that derives the ECHO set by observing what the projection drops is tautological — it restates the implementation and cannot fail." — (NC, `0023:C2`)

- [REQ-66] "the assignment is DECLARED (a standalone table keyed by field name, one entry per field, carrying `echo` or `plan` — the carrier constrained below)" — (NC, `0023:C2`)

- [REQ-67] "the projection is implemented FROM that declaration as a DISTINCT function taking the assembled payload and returning the projected one" — (NC, `0023:C2`)

- [REQ-68] The oracle asserts three things: "every struct field has exactly one entry; the entry set and the field set are equal (neither a field without an entry nor an entry without a field); and the keys a projected run actually emits equal the declaration's `plan` side." — (NC, `0023:C2`)

- [REQ-69] "A new field with no entry then fails at the first assertion rather than defaulting into either width." — (NC, `0023:C2`)

- [REQ-70] "The carrier is CONSTRAINED, not free: it MUST be a standalone table keyed by field name, in its own declaration site, NOT a per-field marker on `resolvePayload`." — (NC, `0023:C2`)

- [REQ-71] "Independence has to be structural — a separate site a projection edit does not open — or S3 asserts only that the implementation agrees with itself." — (NC, `0023:C2`)

- [REQ-72] "Group membership is the declaration table's (C2), never a Go type's." — (CEN, line 1029–1030) — an implementation or oracle that treats "pointer + omitempty" as the mark of a projectable field is forbidden, because `escape_class` is a plain `string` with `,omitempty` and is always-keep core.

## Always-keep core (`0023:C2`)

- [REQ-73] "This verb's ALWAYS-KEEP core (JDR 0002 §D1) is rule, escaped, escape_class, revision: if the boolean axis ever generalizes to an enum or a field list, no mode may PROJECT them away" — (NC, `0023:C2`) — **[JDR-0002:D1]** (JDR lines 56–60).

- [REQ-74] "Always-keep is projection-invariance, not unconditional presence: a core field whose producer already has a presence rule (escape_class) appears under the flag exactly when it appears by default" — (NC, `0023:C2`)

- [REQ-75] "This RDR therefore does NOT mint a separate always-keep oracle (it would today assert exactly what S2 asserts); the obligation it creates is on the successor: an RDR adding a second projection mode owes an always-keep oracle quantified over MODES" — (NC, `0023:C2`) — negative REQ: no separate always-keep oracle is owed by THIS build; S2's literal is the enforcement.

- [REQ-76] "revision rides the PLAN side deliberately: it is produced by the loader, never restated from the request … Projecting it costs 14 bytes and keeps the wire shape stable for the day 0002 admits the key." — (NC, `0023:C2`) — `revision` is carried under the flag even though it is constant-empty today.

- [REQ-77] "Under JDR 0002 §D1, 0024's dispositions, if that RDR lands, is a PLAN-group field … and its never-omitted clause is untouched by this projection." — (NC, `0023:C2`) — **[JDR-0002:D1]**; `0023:A4` and `0023:JC1` ground the ordering tolerance.

## Oracle discriminability — negative controls (`0023:§oracle-discriminability`)

- [REQ-78] S1's controls: "mutate one carried value under the flag; separately, add a large plan-group field and confirm the width clause goes red where the key-set oracle stays green; separately — the control that discriminates this clause — SKIP the reader pass under the flag and confirm the reader assertion goes red." — (OD, `0023:S1` row)

- [REQ-79] S2's controls: "add a field to the payload without a C2 side; S2 and S3 must both go red. Absence control: render one echo key as `null` (the bare non-pointer `omitempty` hazard A2 reproduced) — S2 must go red where a key-presence-only check would stay green" — (OD, `0023:S2` row)

- [REQ-80] S3's discriminating control (not shared with S2): "add a field, declare it `plan`, but omit it from the projection — S3 must go red on the emitted-keys-equal-`plan` assertion while S2 stays green" — (OD, `0023:S3` row)

- [REQ-81] S4's controls: "register `plan-only` on a second command; S4 must go red. Second control: register it on the auto-generated `completion` command — S4 must go red there too, which is the control a bare `NewRootCmd()` tree cannot run" — (OD, `0023:S4` row)

- [REQ-82] S5's controls: "drop a plan-group line under the flag; S5 must go red. Second control: make the text projection a no-op (project nothing) — the subset assertion stays green and only the width assertion goes red" — (OD, `0023:S5` row)

- [REQ-83] "S4 is the one structural-absence oracle here, and the vacuity guard is what stops it passing on an empty command tree." — (OD, lines 1060–1061)

## Construction obligations carried out of the assumptions (`PRE`)

- [REQ-84] "Two CONSTRUCTION obligations carry into the build — these are oracles to write, not assumptions to settle, and neither gates lock: the S4 total walker itself (Phase 2, with both preconditions asserted — auto-generated commands materialized and `completion` present in the walked set), and A9's empty-container default-mode assertion (Phase 1, BEFORE the pointer conversion lands, riding S2's absence control)." — (PRE, lines 1555–1561)

- [REQ-85] A9's Phase-1 oracle: "a default-mode assertion over a request whose containers are empty (no `--tag`, no owned keys, no readers), asserting the rendered bytes carry `{}`/`[]` and no `null`, at implementation Phase 1 BEFORE the conversion lands." — (`0023:A9` Residual)

- [REQ-86] "the five pointers must be non-nil whenever a success payload is emitted, and that is now an implementation invariant stated here rather than an implicit one." — (`0023:A9`, Existing coverage)

- [REQ-87] "Ordering tolerance with cli/0024 confirmed (A4): either RDR may land first; the second lands with `dispositions` already/newly in the plan group and no contract in either moves." — (PRE, lines 1562–1564)

## Mechanism (`0023:A2`, `PH1`)

- [REQ-88] The conversion scope: "the mechanism applies to the five ECHO fields only, each `T` → `*T` with `T` unchanged from the census's Go-type column (`string`, `map[string]string`, `[]string`) — no other field on `resolvePayload` becomes pointer-valued." — (`0023:A2`, Conversion SCOPE)

- [REQ-89] "`escape_class` stays a plain `string` whose `,omitempty` drops the key on `\"\"`. Its presence rule (`0005:A-3`) is therefore realised by a mechanism DISTINCT from the projection mechanism, and the two must not be conflated" — (CEN, lines 1023–1027) — **[0005-carried]**

- [REQ-90] "The `resolvePayload` struct keeps exactly its current field count: `decision_table_0010_test.go` pins `NumField() == 14` (:435) and the 13-key wire list (:422-426), so A2's pointer conversion changes field TYPES only." — (PH, Phase 1) — **[0010-carried]**

- [REQ-91] "Adding a field here would falsify A3's \"no predecessor oracle moves\" claim, so a field addition is out of scope for this RDR by construction, not by preference." — (PH, Phase 1) — negative REQ.

- [REQ-92] "Register `--plan-only` on `resolve` only and hand `respond.OK` the projected verb result when set — the mechanism A2 verified, applied after payload assembly, before the gateway." — (PH, Phase 1)

- [REQ-93] "Flag registration | `newFlowResolveCmd` (verb-local), `registerSelectionFlags` (shared) … Extend verb-local only | The flag never enters the shared registrar (C1)" — (EIA, line 1091)

- [REQ-94] "Both-modes rendering | `internal/cli/respond` (`OK`, `text.go::flatten`) … Reuse unchanged | Projection lands before `respond.OK`, so no gateway change" — (EIA, line 1090) — negative REQ: no edit to `internal/cli/respond`.

## Minimum Viable Validation (`0023:MVV`)

- [REQ-95] MVV step 1: "Run `flow resolve` with a recognized outcome and discriminating `--tag`s, `--as=json`, without the flag → today's full payload, byte-identical to the pre-change fixture." — (MVV, `0023:MVV`)

- [REQ-96] "The fixture is a GOLDEN CAPTURED FROM THE PRE-CHANGE BINARY and checked in BEFORE Phase 1 edits the struct (Phase 1's first step, `git stash`-clean tree) — this is the only comparison in the whole battery that has a pre-change side." — (MVV, `0023:MVV`)

- [REQ-97] "S1–S5 all compare the new build against itself, so none of them can see a default-mode regression that moves both sides together; without a captured golden, C1's headline \"byte-identical to today's\" is asserted nowhere" — (MVV, `0023:MVV`) — the golden is load-bearing for REQ-27, not optional.

- [REQ-98] MVV step 2: "Re-run the same invocation with `--plan-only` → exactly the normative projected key set (`revision`, `rule`, `gates`, `emit`, `next`, `writes`, `clear`, `escaped` — plus `escape_class` iff step 1 carried it); every carried field byte-identical to step 1's; `model`/`observed`/`owned`/`readers`/`outcome` absent (not null, not empty)." — (MVV, `0023:MVV`)

- [REQ-99] MVV step 3: "Re-run step 2 with `--as=text` → the projected lines are a byte-identical subset of step 1's text lines, stable across repeated runs." — (MVV, `0023:MVV`)

- [REQ-100] MVV step 4: "Re-run with an unrecognized outcome, ± the flag → byte-identical refusal envelopes and exit codes." — (MVV, `0023:MVV`)

- [REQ-101] MVV step 5: "Run `flow next --plan-only` → `command-error`, exit 2; and the whole-tree structural oracle (the set of commands registering the flag is exactly `{flow resolve}`, over a root with the auto-generated commands materialized and present in the walked set) passes." — (MVV, `0023:MVV`)

- [REQ-102] "One sibling verb is run, not all three, deliberately: C1 makes the behavioural run CORROBORATION and the structural walk the assertion" — (MVV, `0023:MVV`) — negative REQ: no obligation to run `read-state`/`set-state` behaviourally.

- [REQ-103] MVV step 6: "Record default vs projected byte counts on the motivating-model shape and one gate/write-heavy fixture, measured on the full emitted line (C1's unit), and compare against A1's baseline in `evidence/spikes/a1-byte-width.md`" — (MVV, `0023:MVV`)

- [REQ-104] The step-6 PASS BAR: "every shape measured saves bytes (the S1 width oracle already enforces this per-run), and the CHECKED-IN fixtures (`models/examples/pricing-decision-table.toml`, `release-grammar.toml`, `review-state-machine.toml`) each save at least 40%" — (MVV, `0023:MVV`)

- [REQ-105] "a landing below 40% there means the plan group is carrying materially more than A1 measured, which is A1's \"saves too little to justify a new surface\" condition and routes back rather than recording a number." — (MVV, `0023:MVV`) — the sub-40% outcome is a route-back, not a recorded deviation.

- [REQ-106] "The 79.4% figure is NOT a pass bar, because its shape is not checked in … To gate on it, the step re-materializes the 48-fact model from the spike's recorded generator first; otherwise the checked-in bar above is the gate." — (MVV, `0023:MVV`) — negative REQ.

- [REQ-107] The MVV runs "Over the checked-in decision-table fixture (`models/examples/pricing-decision-table.toml` or the 0005 fixture family)" — (MVV, `0023:MVV`)

## Normative fixtures (`0023:S1`, `0023:S2`)

- [REQ-108] "Normative fixture: the projected wire record of the pricing 2×2 call, whose exact bytes are the \"Projected reference\" line in `evidence/spikes/a2-encoder-mechanism.md` §Reference output" — (TS, `0023:S1`); S2 pins the same fixture.

- [REQ-109] "the corresponding widths (290 B full line → 152 B projected, envelope included — the same unit this scenario asserts) are row S1 of `evidence/spikes/a1-byte-width.md`." — (TS, `0023:S1`)

- [REQ-110] "Witness: 15 default lines → 9 projected on the normative fixture." — (TS, `0023:S5`); the desk trace names the six dropped lines: "`model`, `observed.region`, `observed.tier`, `outcome`, `owned`, `readers`" — (DT, step 3).

- [REQ-111] The width assertion is "unconditional because the five echo keys always render, `model`/`outcome` being non-`omitempty` and the three containers rendering `{}`/`[]`" — (TS, `0023:S1`)

## Phase 2 oracle enumeration (`PH2`)

- [REQ-112] "Five oracles, one per Testing Strategy scenario S1–S5, and the enumeration here is that list — not a sixth" — (PH, Phase 2) — S1 differential, S2 explicit key set carrying absent-not-null, S3 reflective partition-completeness, S4 whole-tree walk, S5 text subset + width.

- [REQ-113] "the explicit projected-key-set oracle (S2), which CARRIES the absent-not-null assertion — a projected key must be absent, never `null` and never an empty placeholder, the live hazard A2's spike reproduced with bare non-pointer `omitempty`" — (PH, Phase 2)

- [REQ-114] "Every C1 MUST maps to at least one oracle below or an MVV step; A1's byte table is recorded evidence, not a test assertion." — (TS, lines 1701–1704) — negative REQ: no oracle asserts A1's byte percentages.

## Docs and help (`PH3`)

- [REQ-115] "`docs/cli-output-contract.md` gains the projected worked payload beside the full one; `flowResolveExtendedDesc` gains the flag under \"Reading a successful plan\"; the partition statement (C2) lands where the payload fields are documented." — (PH, Phase 3)

- [REQ-116] "`flowResolveExtendedDesc` and the flag's own usage string are inputs to the GENERATED artifacts `docs/cli-reference.md` and `llms.txt`, which `make docs-check` verifies as up to date and which `make check` runs … So this phase regenerates those files and commits them in the same change." — (PH, Phase 3)

- [REQ-117] "because the generated artifacts derive from the flag's registration, the commit that registers the flag MUST also carry the regenerated `docs/cli-reference.md` and `llms.txt`, or `make check` fails on `docs-check` at that commit." — (PH, Phase 3)

- [REQ-118] "The flag's usage string is authored here, once, since it ships into a committed artifact rather than staying an implementation detail." — (PH, Phase 3)

- [REQ-119] "the help line states the omitted group by name" — (RM, the model-reference risk mitigation, line 1482).

- [REQ-120] "Verb help | `flow_resolve.go::flowResolveExtendedDesc` … Extend | The help's plan list is the partition's user-facing statement; gains the flag line" — (EIA, line 1094)

## Determinism and cross-cutting (`0023:G-cross-cutting`, `PERF`)

- [REQ-121] "one encoder for both paths, HTML escaping disabled, the measurand fixed as the full emitted NDJSON line excluding the trailing newline (C1). Carried fields keep their default-mode relative order." — (CC, `0023:G-cross-cutting`)

- [REQ-122] "no map-order path exists in text mode: the gateway flattener sorts keys from the decoded map (`internal/cli/respond/text.go::flatten`, A2). JSON key order is struct declaration order, preserved under the pointer conversion." — (CC, `0023:G-cross-cutting`)

- [REQ-123] "An omitted key is **absent**, never `null` and never an empty placeholder (C1); `emit` stays present as `{}` per `0010:C4`; `observed`/`owned` render as `{}` when empty in default mode, which bare-map `omitempty` would wrongly drop (A2)." — (CC, `0023:G-cross-cutting`) — **[0010-carried]**

- [REQ-124] "*Whitespace and case folding* — unchanged; the projection deletes whole keys and never rewrites a carried value." — (CC, `0023:G-cross-cutting`)

- [REQ-125] "*Version marker* — none is minted, deliberately. The wire shape is identified by the flag the caller passed, not by an in-payload marker" — (CC, `0023:G-cross-cutting`) — negative REQ.

- [REQ-126] "Wall time unchanged: projection is key deletion on the assembled payload before `respond.OK` — no extra I/O, no second marshal path" — (PERF); and "no oracle or scenario asserts a wall-time threshold, and none is specified" — negative REQ: no wall-time budget is owed.

- [REQ-127] "one wire encoder (`internal/cli/clierr/clierr.go::WriteJSONLine` — `json.Encoder`, `SetEscapeHTML(false)`, compact one-line NDJSON)" — (PERF, determinism checklist).

- [REQ-128] "No new encoding surface. Tag literals are refused rather than coerced by the shipped parser (`internal/cli/flow_input.go::parseTags`, A7); the sole rewrite is `canonicalSet`'s re-encode of an already-valid set literal, which is value-preserving." — (CC, `0023:G-cross-cutting`) — **[0005-carried]**; negative REQ: no change to tag parsing.

- [REQ-129] "The flag is additive and opt-in, so no existing caller observes a change; default-mode output is byte-identical to the pre-change binary (C1, guarded by the MVV's pre-change golden)." — (CC, `0023:G-cross-cutting`)

## Failure modes and consequences (assertable residue)

- [REQ-130] "a projected payload missing a plan-group field (mechanism bug) fails the ±-flag equality oracle" — (FM, `0023:F2`)

- [REQ-131] "a projection accidentally applied to a refusal path would change refusal bytes — A6 plus a refusal-unchanged oracle make it a test failure, not a field report." — (FM, `0023:F3`)

- [REQ-132] "Version skew is loud, not silent, in BOTH directions. Forward — a binary predating this RDR refuses `--plan-only` at the parse (`command-error`, exit 2), so a caller can never hold a full-width payload believing it was projected, and an absent `observed` always means projection, not an old binary." — (FM, `0023:F5`) — **[0011-carried]** for the shared-bucket limit.

- [REQ-133] "Diagnosis: `--as=json` ± `--plan-only` over the same request diffs to exactly the echo keys; any other diff is a defect in the projection." — (FM, `0023:F6`)

- [REQ-134] "Recovery: drop the flag — the default width is the full record and is contractually unchanged." — (FM, `0023:F7`)

- [REQ-135] "the default is untouched — no consumer changes, no 0005/0010/0011 test moves (A3), no respond-gateway change." — (CONS, lines 1436–1437; the source wraps `0005/` and `0010/0011` across the line break) — **[0005-carried]**, **[0010-carried]**, **[0011-carried]**; the whole predecessor suite must stay green unedited.

- [REQ-136] "this RDR ships the CAPABILITY, not its adoption. The flag is opt-in and no consumer passes it on landing … Phase 2 green plus the MVV step-6 table is the completion bar" — (CONS, lines 1413–1421) — the done bar; no consumer call site is migrated by this RDR.

- [REQ-137] "`intrastate#srz2` (P1) is the ADOPTION ask, so this RDR landing green does not close it — it unblocks it." — (CONS, lines 1427–1430) — negative REQ: do not close the tracker on a green build.

---

## ASSUMPTIONS

Implicit choices made where the wording is imprecise but one reading is
defensible. Each names the REQ it attaches to and the evidence it rests on.

- **A-1 — attaches to REQ-21/REQ-23/REQ-26. "default-mode relative order" for
  the projected wire is `revision, rule, gates, emit, next, writes, clear,
  escaped[, escape_class]`, i.e. `revision` FIRST.** C1 fixes relative order but
  never lists the projected order; the S2/MVV-step-2 literal lists `revision`
  first, and the shipped struct declaration order on `main` is `Model,
  Revision, Observed, Owned, Readers, Outcome, Rule, Gates, Emit, Next, Writes,
  Clear, Escaped, EscapeClass` (`internal/cli/flow_resolve.go:31-51`), so
  deleting the five echo keys leaves `revision` at the head. The desk trace's
  step-2 witness confirms it byte-for-byte:
  `{"revision":"","rule":"paid-eu","gates":[],...}`. Read as: the projected
  order is the default order minus the deleted keys, never a re-sort.

- **A-2 — attaches to REQ-23/REQ-98. The projected key set is EIGHT keys when
  `escape_class` is absent by its producer's rule, NINE when present.** S2 and
  MVV step 2 both write eight names "plus `escape_class` iff step 1 carried
  it", and the desk trace step 2 says "8 keys, equal to S2's list" on a
  non-escaped fixture. Read as: the explicit literal list the oracle asserts is
  conditional on the same run's default output for `escape_class` alone
  (REQ-25), and unconditional for the other eight.

- **A-3 — attaches to REQ-41/REQ-42. The reader-execution seam is TEST-VISIBLE
  instrumentation at the `exec.Read` call inside `runReaders`, not a production
  behaviour change.** C1 says "the seam's FORM is left to the implementer, but
  its position is not" and lists "a counting or recording seam around the
  reader pass, the run's own side effects" as admissible. Nothing in C1 or C2
  permits the flag to be readable at that site (REQ-49 forbids it), so the seam
  cannot be flag-conditional. Read as: a flag-blind recording hook (or an
  observable side effect of the real reader run) that both the default and the
  projected run drive identically, with the oracle comparing the two recordings.

- **A-4 — attaches to REQ-66/REQ-70. The declared assignment table lives in
  non-test production code, keyed by the Go FIELD name.** C2 says "a standalone
  table keyed by field name, in its own declaration site" and REQ-67 requires
  the projection to be "implemented FROM that declaration", which is production
  behaviour, so the table cannot live in a `_test.go` file. C2's oracle
  assertion "every struct field has exactly one entry" is reflective over
  `resolvePayload`, whose reflection surface is field NAMES (`NumField`,
  `FieldByName` — the same surface `decision_table_0010_test.go` uses per
  `0023:A3`), so "field name" reads as the Go identifier, with the JSON tag
  reachable from it for the emitted-keys assertion. Read as: one exported-or-
  package-level `map[string]group` (or equivalent) in `internal/cli`, in a file
  the projection edit does not open.

- **A-5 — attaches to REQ-36/REQ-38. "STRICTLY SHORTER" is a strict inequality
  in bytes (`projected < default`), asserted per-run inside S1, not an aggregate
  or percentage bar.** C1 names the unit (full emitted line, trailing newline
  excluded) but not the comparator's form. The MVV pass bar separately carries
  the 40% figure and explicitly says "the S1 width oracle already enforces this
  per-run" for the "saves bytes" half, distinguishing the two. Read as: S1
  asserts `len(projectedLine) < len(defaultLine)` with no threshold; the 40%
  bar is MVV-only (REQ-104) and is not an oracle.

- **A-6 — attaches to REQ-7. "cobra's own initializers" are the exported
  `(*cobra.Command).InitDefaultHelpCmd` and `(*cobra.Command).InitDefaultCompletionCmd`
  called directly on the root the oracle builds.** C1 names
  `InitDefaultCompletionCmd` and says "reachable from ExecuteC"; `0023:A5`'s
  spike records both as "EXPORTED and callable directly from a test … so no
  fork or unexported access is needed". Read as: the oracle calls the two cobra
  methods on `NewRootCmd()`'s result before walking, rather than invoking
  `Execute`/`ExecuteC` (which would run a command) or calling the repo's
  `registerHelpAllOnTree`.

- **A-7 — attaches to REQ-33/REQ-100. "byte-identical refusal envelopes" is
  asserted on the full emitted refusal LINE, the same unit as the success-width
  clause.** C1 fixes the full-emitted-line unit for the width assertion in the
  same clause block and says refusals are "never projected"; A6 grounds that no
  refusal byte can move. Read as: the differential compares the complete
  refusal NDJSON record (and exit code) between the two runs, not just the
  `.data`/error object.

- **A-8 — attaches to REQ-96. The pre-change golden is checked in under this
  RDR's own evidence/fixture surface as the FIRST Phase 1 commit step, before
  any struct edit.** The MVV says "checked in BEFORE Phase 1 edits the struct
  (Phase 1's first step, `git stash`-clean tree)" but names no path. Read as:
  the location is the implementer's, the ORDERING is not — the capture commit
  precedes the conversion commit, and a golden captured after any Phase 1 edit
  does not satisfy REQ-96/REQ-97.

- **A-9 — attaches to REQ-112/REQ-84. The S4 walker and the A9 empty-container
  assertion are scheduled outside Phase 2's five-oracle enumeration only in the
  sense of phase timing, not count.**
  Phase 2 says "Five oracles, one per Testing Strategy scenario S1–S5, and the
  enumeration here is that list — not a sixth", while `PRE` names the S4
  walker itself and A9's Phase-1 assertion as construction obligations. Read as:
  A9's empty-container assertion rides S2 (its absence control, per `0023:A9`)
  and is written in Phase 1; the S4 walker IS oracle S4 and is written in Phase
  2 — so the build carries five oracles, one of which (S2) gains the A9 arm
  early and one of which (S4) is a new symbol.

- **A-10 — attaches to REQ-88. The five echo fields' JSON tags gain
  `,omitempty` as part of the pointer conversion.** `0023:A2` names the
  mechanism as "pointer-valued echo fields with `omitempty`" and A9 reasons
  throughout about the nil arm "under `omitempty`", but no clause states the tag
  edit as a MUST. The census pins `T` unchanged and forbids other fields
  becoming pointers (REQ-88), and `escape_class` keeps its existing
  `,omitempty` (REQ-89). Read as: `*T` + `,omitempty` on exactly `model`,
  `observed`, `owned`, `readers`, `outcome`; the JSON key strings themselves are
  unchanged, which is what keeps the three shipped wire oracles green
  (`0023:A9`, Existing coverage).

- **A-11 — attaches to REQ-64/REQ-90. "the 13-key wire list" and "FOURTEEN
  fields" are consistent, not a contradiction, because `escape_class` is
  `omitempty` and absent from the non-escaped reference row.** The census says
  fourteen fields "laid out in thirteen rows below because `writes` and `clear`
  share one `plan.Writes` writer" (a different 13), while Phase 1 cites
  `decision_table_0010_test.go`'s 13-key WIRE list. Read as: 14 struct fields,
  14 rendered keys when `escape_class` is present, 13 when it is not — and the
  projected width is 9 or 8 correspondingly (A-2). No REQ asserts a "13" over
  the struct.

---

## QUESTIONS

None. Two clauses were tested against the evidence base for genuine ambiguity
and both are answered by the record rather than escalated:

- The "STRICTLY SHORTER" comparator (threshold vs strict inequality) is settled
  by the MVV's own split between the per-run S1 oracle and the 40% MVV bar —
  recorded as ASSUMPTION A-5 rather than a question.
- The declared-assignment carrier's location (production vs test) is settled by
  C2's requirement that the projection be implemented FROM the declaration,
  which forces production code — recorded as ASSUMPTION A-4.

One item is a NOTED DEPENDENCY, not a question, because the record already
disposes of it (`0023:A4`, `0023:JC1`, REQ-87): if RDR 0024 lands first,
`decision_table_0010_test.go`'s `NumField() != 14` (:435) and its exact
wire-key-set `slices.Equal` (:427) will already have been moved by 0024's own
additive append. `0023:A4`'s critique-lens scope note states this explicitly —
"That cost is 0024's own additive append … and is not created by this
projection" — so a red assertion there on a 0024-first tree is 0024's, and this
RDR's REQ-64/REQ-90 read against whatever count 0024 left. This RDR alone moves
neither assertion.
