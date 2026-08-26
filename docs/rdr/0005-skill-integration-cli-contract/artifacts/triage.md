# Triage — RDR 0005 Skill Integration CLI Contract

Single-pass roborev triage of the implementation launch.
Window frozen once: `BASE=7f2bdfd` .. `HEAD=2e16e52`. Batch: `batch:rdr-0005`.

18 findings over 7 in-window jobs (6230, 6232, 6234, 6235, 6236, 6238, 6239).
Each was grounded by a sub-agent against the RDR, `req-list.md`,
`deviations.md`, and `verification.md` — the context roborev's repo-sandboxed
prompt never sees — then verified by the orchestrator before execution.

## Verdicts

| # | Job | Location | Verdict | odc-type / trigger | Evidence | Outcome |
|---|---|---|---|---|---|---|
| F-A | 6234 | `flowbind.go:266` (**High**) | DROP superseded-at-HEAD | algorithm / logic-flow | `sealReadBack` deleted by `e76b3a5`; HEAD sets the seal into the already-mutated map before ONE atomic save. Verified: no such symbol at HEAD. | drop |
| F-B | 6236 | `flowbind.go:264` | KATA-BUG | algorithm / rare-situation | No `else` clears `sealedKey`; the seal is monotonic. No REQ governs seal lifetime (DEV-6 is an IMPL-DECISION). | kata `dkcf` ★ |
| F-C | 6236 | `flowbind.go:270` | FIX-NOW | test-oracle / test-coverage | The CLI oracle asserts code+exit+detail but never opens the artifact; the package had NO unit tests. | `5c43098` |
| F-D | 6232 | `flowbind.go:246` | DROP superseded-at-HEAD | interface / design-conformance | `readBackReader` absent at HEAD (`e76b3a5`); the seal design is DEV-6(c)-adjudicated. | drop |
| F-E | 6232 | `flow_next.go:147` | DROP rdr-adjudicated | interface / design-conformance | DEV-4 decided this exact question; HEAD already clones the declared alphabet. | drop |
| F-F | 6232 | `flow_next.go:251` | DROP superseded-at-HEAD | algorithm / logic-flow | `71f5762` routes exclusion through the kernel probe (ADV-1). | drop |
| F-G | 6232 | `flow_exec.go:491` | DROP superseded-at-HEAD | function / logic-flow | `00f6088` passes the full result slice (FAIL-1). | drop |
| F-H | 6238 | `flow_exec.go:616` | KATA-BUG | interface / design-conformance | Live. RDR names no per-gate code, but REQ-24's convention and the function's own doc comment are falsified by its default arm. | kata `sd16` ★ |
| F-I | 6239 | `flow_exec.go:89` | FIX-NOW | algorithm / boundary | Live regression introduced by `79b88fc`'s own demand-set widening. Kernel skips differing-outcome escape rows (`resolve.go:563`), so their keys are not demands. | `d587e1f` |
| F-J | 6232 | `flow_input.go:236` | KATA-BUG | checking / boundary | `canonicalValue` never receives the `TagDecl`; an out-of-domain value persists and read-back reports success. | kata `98n6` ★ |
| F-K | 6232 | `flow_state.go:310` | KATA-BUG | checking / boundary | Loader counts writers only for keys in `written`; a key no rule names escapes both checks. Phase 3b's unreachability argument is incomplete. | kata `8dg3` |
| F-L | 6232 | `flow.go:55` | KATA-BUG | interface / design-conformance | Bare `flow` exits 0. Downgraded from FIX-NOW: no REQ governs it and a fix needs a code the closed 25-row table lacks. | kata `3cwy` ★ |
| F-M | 6235 | `cli-output-contract.md:38` | FIX-NOW | documentation / doc-code-drift | Carrier rule false for one-element aggregates, verified at `flow_exec.go:301` and `flow_input.go:146`. | `5c43098` |
| F-N | 6235 | `cli-output-contract.md:63` | FIX-NOW | documentation / doc-code-drift | `Fingerprint` is serialized but absent from the field table. | `5c43098` |
| F-O | 6235 | `cli-output-contract.md:82` | FIX-NOW | documentation / doc-code-drift | `WriteJSONLine` neither sorts nor dedupes; `canonicalSet` is the real encoder. | `5c43098` |
| F-P | 6230 | `finding_0005_test.go:188` | DROP superseded-at-HEAD | n/a-not-a-defect | Fixed by `77f190a` (DEV-2). Deliberately-RED Phase-1 commit. | drop |
| F-Q | 6230 | `flow_exit_0005_test.go:41` | KATA-BUG | test-oracle / test-coverage | SURVIVES at HEAD — not pre-filtered. The exit-3 branch is dead; the test proves only the exit-2 half. | kata `g6vn` |
| F-R | 6230 | `flow_next_0005_test.go:283` | KATA-BUG | test-oracle / test-coverage | SURVIVES at HEAD. Any nonempty verdict passes, including a coerced `allow`. | kata `r2ba` |

★ = the kata body carries an `## Open question`.

## Fix-commit sweep (bounded, one round)

`d587e1f` reviewed clean. `5c43098` raised one Medium finding — my own F-M edit
added the family-based carrier rule but left the superseded cardinality block
below it, so the section contradicted itself. Because that defect was
introduced BY this triage rather than pre-existing, it was corrected in place
(`aee38b0`) instead of deferred. No further round was run.

## Verification the orchestrator performed directly

Grounding verdicts were not taken on trust:

- **F-A** (the only High): confirmed `sealReadBack` and `readBackReader` are
  absent at HEAD and the current `Apply` preserves the mutation — a High
  severity DROP needs positive evidence, and it has it.
- **F-I**: confirmed both legs — `invokedReaders`' escape conjunct and the
  kernel's `row.Outcome != in.Recognized` skip at `resolve.go:563`.
- **F-C** and **F-I** regressions were **mutation-tested**: reintroducing each
  defect fails the new oracle, restoring the fix passes it. Neither is a
  tautological pass.
- **F-L** was downgraded from the sub-agent's FIX-NOW after checking that REQ-5
  binds only the four verbs and the refusal table names no usage code —
  shipping one would violate Stage 8's additive-is-not-exempt rule.

## Distribution (excludes `n/a-not-a-defect`)

`algorithm=3 checking=2 documentation=3 function=1 interface=4 test-oracle=3`
