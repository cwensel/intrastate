Model: claude-opus-5[1m] (consolidation; personas stamped in their own files)

# 3amigo consolidation — 0024 declared emit vocabulary

Three isolated persona passes (no cross-persona visibility). Overlap below is a
mechanical set intersection on anchored element ids, not a re-judgement. **A
hotspot marks a passage two or more personas independently landed on; it does
not validate the finding, and a single-persona finding is not thereby weaker.**

Totals: PM 9, Implementer 12, QA 9 — 30 findings, 8 high.

## Hotspots (element-id intersection)

| Anchor | Personas | What each independently landed on |
| --- | --- | --- |
| `0024:C1` | PM, Impl, QA | PM-5 scalar hollowing · I1 carrier ordering, I3 reuse shape, I6 stale docs, I7 illustration domain, I8 int authority, I9 bare `[emit]` · Q-1 union order, Q-6 category witnesses, Q-9 dump rendering |
| `0024:C2` | PM, Impl, QA | PM-4 N-round adoption cost · I2 load-step placement/order, I9, I11 category constants · Q-6 witnesses, Q-8 envelope code |
| `0024:C4` | PM, Impl, QA | PM-1 disposition relocation, PM-9 never-omitted vs `omitempty` · I4 declared-but-unemitted join, I10 TestReq39 edit, I12 text render · Q-3 §D1 registration |
| `0024:C3` | Impl, QA | I1 carrier ordering, I3 `TagDecl` contradiction · Q-1 lossless-vs-union, Q-2 unnamed carrier surface |
| `0024:MVV` | PM, QA | PM-7 all oracles intrastate-internal · Q-8 envelope code slug |

The four contracts are the whole hotspot set — expected for a record whose
Normative Contracts are four facets of one contract, and it means the personas
converged on the contract text rather than on peripheral prose.

**Convergent defect, three routes.** The carrier's type and ordering is the one
passage all three personas reached by different questions: the implementer as a
first-hour blocker (I1, I3), QA as a missing oracle (Q-1, Q-2), the PM not at
all — which is itself signal that it is an implementation-facing defect, not an
outcome one. Independent convergence on `C3`'s negative-only carrier definition
is the strongest signal in this pass.

## Findings ledger (origin ledger for the resolve half)

Severity as the persona filed it. `Ground` records the pre-resolve check already
run against `main`; findings without one are grounded during resolve.

### High

| ID | Anchor | Finding | Ground |
| --- | --- | --- | --- |
| I3 | `C1`,`C3`,`A4` | Three incompatible value-conformance reuse shapes: A4 says reuse `ConformValue(decl TagDecl, member string)`; C3 forbids `TagDecl` as carrier; Phase 1 says reuse unexported `conformKind`/`conformDomain` — which also take `TagDecl`. | CONFIRMED — `load.go:793,833,849` all take `TagDecl`; A4's "no adapter or extraction required" holds only if the carrier IS a `TagDecl`. |
| Q-1 | `C3`,`C1`,`S4` | C1 takes the partitioned domain as an unordered "union"; C3 + the `fidelity` mini-check claim a lossless/total carry. Both cannot hold — the union discards partition grouping and within-partition order. S4 has no defined `want`; a map-backed carrier risks `TestReq146`'s determinism sweep. | CONFIRMED — tag precedent carries one authored array (order inherited free); the partitioned form has no authored array to inherit. |
| I1 | `C1`,`D-identity` | Carried declaration set's ordering unstated; a map carrier breaks the determinism every sibling normalized surface holds. | Same defect as Q-1, reached from identity rather than fidelity. |
| I2 | `C2` | "Beside `loadTags`, before `normalizeRules`" leaves a five-step window and never fixes that declaration load must precede the rule cross-check. | Partially CONFIRMED — `load.go:84,90` show `loadTags` then `normalizeRules`; the mutual order of the two new steps is genuinely unstated. |
| I4 | `C4`,`D-selection-predicate` | Join behavior for a declared key the selected row does not emit is unspecified; C2's "authoring headroom" makes a projected-key-set reading plausible. | Open — C4 says "no selected value carries one" ⇒ `{}`, but is silent on per-key absence within a non-empty map. |
| PM-1 | `§problem-statement`,`C4` | Verbatim-surfaced `dispositions` relocates the string discipline from value prefix to payload field; `§consequences`'s "consumers stop re-implementing" is substitution, not removal. | Outcome-level; adjudicate against the decided text (intrastate generic by design). |
| PM-2 | `§phase-4`,`A5` | A5 grounds the disposition partition on `rdr-status.toml`/`rdr-write.toml`, which are not in this repo; Phase 4 declares the pricing example, whose `plan`/`dpa` keys have no route/stop split — so the sub-table grammar ships demonstrated nowhere. | CONFIRMED — those models live in the sibling `rdr` engine repo; pricing example emits only `plan` (basic/pro) and `dpa` (required/none). |
| PM-3 | `§approach`,`§consequences` | Fully opt-in, zero declarations repo-wide (A3), no adopter scheduled — merge-day outcome is one example declaration plus an empty payload field. No falsifiable post-merge observation. | CONFIRMED — `rg '^\[emit'` over `*.toml` returns zero. |

### Medium

| ID | Anchor | Finding | Ground |
| --- | --- | --- | --- |
| Q-2 | `C3`,`S4` | Carrier named only negatively ("NOT `TagDecl`") — no `Model` field, type name, or accessor, so S4 has no read site and the `ConformValue` adapter is an unnamed, untested seam. | CONFIRMED — tag analogue is concrete (`Model.Tags map[string]TagDecl`, `model.go:435`). |
| Q-3 | `S1`–`S5` gap vs Phase 3 | Phase 3 names the JDR 0002 §D1 ECHO/PLAN registration as a mandatory red-test risk; no scenario asserts it. Done criterion is satisfiable with `dispositions` unregistered. | Phase 3 text confirms the obligation and its own framing ("not optional"). |
| Q-4 | `S3` | S3 demands both "the only observable delta is `dispositions: {}`" and "byte-for-byte equivalence" — C4's unconditional append makes byte-identity false; "full suite green" also needs Phase 3's 28 licensed regenerations S3 never mentions. | CONFIRMED — A3's own Evidence carries the correct framing; "byte-for-byte" is the stray word. |
| Q-5 | `S1`–`S5`/`MVV` gap vs Phase 4 | Neither the pricing-example `[emit]` block nor the `cli-output-contract.md` edit has a scenario, though peers 0009/0011 pin that doc with `readRepoFile` assertions. | Phase 4 confirmed to ship two untested artifacts. |
| Q-6 | `C1`,`C2`,`S1`,`S2` | Three new categories owe a `testdata/neg/` witness and a `table.Categories()` entry per the REQ-119 convention; the hand-maintained witness map fails nothing if they are omitted (silent degradation). | Convention cited at `dump_test.go:820-875`. |
| I5 | `C1`,`A7` | No `sourceDoc`/`sourceEmitDecl` shape given; A7's arm-closure claim is Pending with no stated gate (Prerequisites lists only A1–A5). | CONFIRMED — Prerequisites names A1–A5 and A6 only; A7 is absent from it. |
| I6 | `C1`, Phase 4 | `docs/model-authoring.md` states emit keys are what "nothing declares" (twice) and that `plan = "<clear>"` loads and answers literally; both read false under a declaration. Phase 1 amends three stale *code* comments making the same claim but misses these *doc* sentences. | CONFIRMED — both passages read as reported; Phase 4 edits only `cli-output-contract.md`. |
| I7 | `C1`,`§illustrative-code`, Phase 4 | Illustration's `[emit.dpa] domain = ["required","waived"]` refuses the shipped example's `dpa = "none"`. | CONFIRMED — same defect QA filed as Q-7. |
| I8 | `C1`,`A4` | `int` runs `strconv.Atoi` in both `conformKind` and `conformDomain`, the latter then testing `Min`/`Max` C1 cannot set; authority for `int` stated two ways. | Partly pre-answered — A4 already notes `Min`/`Max` never fire; the double-`Atoi` authority split stands. |
| I9 | `C1`,`C2` | A bare `[emit]` with zero sub-tables is unadjudicated; the repo treats presence-vs-length as contract-level. | Open. |
| PM-4 | `C2`,`MVV` | Fail-fast one-finding-per-run makes first-declaration adoption an N-round fix-and-rerun loop (54 emit blocks in the named adopter) with unspecified interleaving; never priced from the author's seat nor listed as a Negative. | CONFIRMED — 54 `.emit]` blocks in `rdr-status.toml`. |
| PM-5 | `C1`,`§risks` | The likeliest bulk-adoption path (`scalar` everywhere) yields "declared" with zero proof; only control is human review, since Meyer's flagging condition is structurally refused (`0006:C17`). | Decided text already adjudicates the flagging condition; the adoption-path cost is the live part. |
| PM-6 | `§phase-4` | `docs/model-authoring.md`, quoted twice in `§decision-rationale` as the RDR's own argument, is in no phase's edit list. | CONFIRMED — same gap as I6, reached from the rationale side. |
| PM-7 | `MVV`,`§testing-strategy` | Every oracle is intrastate-internal; the outcome-closing check ("every declared member is a real command or stop token") is twice deferred to an out-of-scope, unscheduled consumer seam test. | Decided text places this out of scope by design; the *unscheduled* part is the live edge. |

### Low

| ID | Anchor | Finding |
| --- | --- | --- |
| I10 | Phase 3,`C4` | `TestReq39`'s `NumField() != 14` already matches at HEAD; the adjacent "thirteen to fourteen" message and the 13-key list are the real edits, and the message sits outside the licensed set. |
| I11 | `D-naming`,`C2` | Wire slugs fixed, but not the `Cat*` constants nor whether the `flow-model-invalid` mapping is table-driven; two category-enumerating tests need appending. |
| I12 | `C4` | "Generic payload renderer" plus a named two-line format — confirm `emit` already renders identically so Phase 3 carries zero text-mode diff. |
| Q-8 | `MVV`,`C2`,`F1` | `lint` emits envelope code `model-invalid` (`lint.go:183`) while C2 names only `flow-model-invalid`; the MVV runs `lint`, so a tester asserting from C2 writes the wrong code. |
| Q-9 | `S4`,`C1`, Phase 1 | No scenario pins that `dump` still renders a *declared* enum emit value raw; C1's "never canonicalized downstream" goes untested exactly where Phase 1 warns a wrong change is easy. |
| PM-8 | `G-*` | Four gate sub-sections are unfilled template text. |
| PM-9 | `C4` | `dispositions: {}` never-omitted is the entire shipped user-visible delta plus a 28-golden rewrite; `omitempty` is available and unweighed. |
| Q-7 | `§illustrative-code` | Duplicate of I7 (independent). |

## Persona verdicts

All three returned `NEEDS_DECISION`, `blocking: yes`.

## Verified-clean, explicitly not reported (QA)

S5's wire numbers (13→14 keys, `NumField` 14→15, index 9) match
`decision_table_0010_test.go:416-437`; A2's 28/7 counts reproduce; the
`dispositions: (none)` text rendering matches `respond/text.go::flatten`;
`scalar` falls through both conform arms so "never refused" is mechanically
true; escape rules are `sourceRule` too, so S2's "both rule classes" is
decidable pre-`normalizeRules`; `flow next` has no `Emit` field.

## Review gate

Healthy. Every finding anchors to a named element id with a blocked decision or
prevented test. No persona file references another's output (isolation held).
All three carry model stamps. QA's explicit verified-clean list is the signal of
a pass that read the source rather than the prose.
