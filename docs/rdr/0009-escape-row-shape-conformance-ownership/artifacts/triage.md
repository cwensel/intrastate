# RDR 0009 — roborev triage

- **Window (frozen once):** BASE `0506f21` → HEAD `3b4f082` (16 commits)
- **Batch label:** `batch:rdr-0009`
- **Spine:** 7 in-window per-commit auto-reviews (16 findings)
- **Range net:** 1 review (job #6281), 2 findings — 1 new, 1 duplicate
- **Launch status:** `COMPLETE` (verified before triage)

## Verdicts

| # | Job | Finding | Verdict | ODC type / trigger | Outcome |
|---|-----|---------|---------|--------------------|---------|
| 1 | 6272 | `breachElements` uses `errors.As`, so a `%w`-wrapped child passes REQ-35 | DROP `over-engineering` | test-oracle / design-conformance | The helper re-asserts "no nested join" right after `errors.As` (`fixtures:122-147`) — the clause REQ-35 names. REQ-36 closes the wrapping path upstream with its own test. |
| 2 | 6272 | `compareRefs` ordering passes when sorted by locator alone | DROP `over-engineering` | test-oracle / design-conformance | Input is supplied in REVERSE order and `sameRefs` (`fixtures:178-186`) compares the full `RowRef` tuple position-by-position. REQ-41 pins the identity-tie collapse separately. |
| 3 | 6273, 6276 | REQ-78/79 compare two outcomes from the same implementation | DROP `rdr-adjudicated` | test-oracle / design-conformance | REQ-79 is a build-time procedure per the req-list ASSUMPTION (`req-list.md:768-771`), executed at `verification.md:166-173`: base `05607f5` via `git archive`, 493 outcomes, sorted sets identical both directions. |
| 4 | 6273, 6276 | REQ-74 "value-for-value" compares only 3 fields | **FIX-NOW** | test-oracle / design-conformance | Fixed in `83a13d3`. Whole-`Plan` `DeepEqual`, normalizing only the nil-vs-empty `Writes` difference. Mutation-verified. |
| 5 | 6273 | REQ-99 `n` breaches proves cardinality, not O(rows) | DROP `over-engineering` | test-oracle / design-conformance | REQ-99 sits under "Performance and non-goals"; REQ-102 introduces no ordering/determinism guarantee. The test documents it asserts linearity as a behavioural property. |
| 6 | 6274 | Clear-sentinel assertions loose (`strings.Contains`, any row); `:195` omits rule ID; `:466` fixture not importable | **KATA-BUG** | test-oracle / design-conformance | Filed **`0wc8`**. Real residue, but REQ-90 defers TS scenario 8 enforcement to RDR 0002's build (`req-list.md:793-800`). |
| 7 | 6275, 6278, 6281 | No end-to-end `runFlowResolve` test; CLI could regress to `codeAccessorFailed` | DROP `rdr-adjudicated` / `unreachable` | test-oracle / test-coverage | `normalize.go:344-348` rejects write-bearing escape rules, and `Model.KernelTable()` is the only route into `Resolve`, so no loadable model yields a breaching kernel table. Intended C2/C3 layering (REQ-47); recorded at `verification.md:243-253` and `:497-506`; REQ-91 defers TS 11. The proposed injectable seam is barred by ADDITIVE IS NOT EXEMPT. |
| 8 | 6275, 6278 | No assertion that `CLIError` preserves the kernel chain via `Cause` | DROP `rdr-adjudicated` | test-oracle / test-coverage | Premise false at HEAD: `Cause` is `json:"-"` and this RDR rewrote its doc comment (`clierr.go:79-87`) to "nothing is recoverable from `cause` itself (`0009:C7`)". REQ-59 makes the serialized `findings` carrier the contract. |
| 9 | 6278, 6279 | Multi-breach Rule↔Locator pairing / count round-trip under-asserted | DROP `over-engineering` | test-oracle / design-conformance | Covered by the reverse-ordered full-tuple assertions; added strictness exceeds REQ-40's text. |
| 10 | 6279 | Adversarial tests invent a text-mode `Count` requirement | DROP `superseded-at-HEAD` | test-oracle / design-conformance | Already removed in `30606dc` by the Phase 3c adjudication, reached independently: REQ-44 forbids prose-ONLY carriage, and the REQ-96 ASSUMPTION bars asserting on message text. |
| 11 | 6279 | Wrapped-error scenario unreachable per REQ-36 | DROP `superseded-at-HEAD` | checking / design-conformance | The scenario is unreachable, but the consumer-side fragility was real and was fixed in `52b257a`: a bare type assertion on the outer aggregate would silently truncate an N-breach report. REQ-36 binds the kernel, not this mapper. |
| 12 | 6281 | JD-5 dual-breach precedence untested | **RDR-SEED** | interface / design-conformance | Filed **`6cek`** ★. JDR 0001 §JD-5 is an OPEN joint decision spanning RDR 0008 × 0009 and also requires naming the error wrapping. `resolve.go:432-452` documents placement as "a cheapness preference, not an observable contract"; REQ-22's ASSUMPTION forbids asserting on it; Phase 3c probed and deliberately declined to pin it (`verification.md:475`). |

★ = carries an `## Open question`.

## Fix-commit sweep (bounded, one round)

`83a13d3` triggered a fresh per-commit auto-review. Collected once: **no open
findings**. Not re-fixed, per the one-round rule.

## Katas filed

| short_id | Kind | Severity | Priority | Title |
|---|---|---|---|---|
| `0wc8` | `type:bug` (`lifecycle:queued`) | medium | 2 | Pin the `<clear>` sentinel representation exactly on the `ordinary-clear` row |
| `6cek` | `kind:rdr-seed` (no `lifecycle:*`) | medium | 2 | Decide JDR 0001 JD-5: precondition precedence between RDR 0008 and RDR 0009 |

Both carry `src:roborev`, `batch:rdr-0009`, `release:post-1.0`.

## Distribution (excludes n/a-not-a-defect)

- `test-oracle`: 11 · `checking`: 1 · `interface`: 1
- `design-conformance`: 9 · `test-coverage`: 3

The corpus is dominated by test-oracle/design-conformance findings — reviewer
pressure toward stricter assertions. Nine of twelve were dropped on grounded
evidence, most because the RDR or its launch artifacts had already adjudicated
the point; two of those (the text-`Count` requirement and the wrapped-error
fragility) had already been actioned by Phase 3c before the review landed.
