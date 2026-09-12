Model: claude-opus-5

# Stage 6 Reconcile — cli/0029

Date: 2026-09-12 · Iteration 1 · Verdict: **RECONCILED**

## Stage 5 preflight

`--outcome lens` → `/rdr-reconcile` (Profile large; grounding → 3amigo →
critique row complete). `--outcome critique` → `none` (two passes, models
differ, diffed). `--outcome repeatability` → `none` (variant lite, run-1 and
diff written), which is the Determinacy add-on's answer too — the
`Determinacy: fired` line is satisfied by the lite variant already on disk.
No caveat to carry into Caveats: no single-model fallback, no unstamped
pass, no variant mismatch.

## Open set

| Item | Source | Disposition | Evidence pointer or plan |
| --- | --- | --- | --- |
| A5 — enumeration seam per tiered vocabulary | 1, 2 | **VERIFIED (split)** — confirmed for `flow next`; **REFUTED** for CLIError `code`, taking C4's pre-authorized fallback | `evidence/spikes/a5-enumeration-seams.md`; `internal/cli/flow_next.go` accessor builds clean; `clierr` → `graphlint` → `clierr` cycle reproduced under `go build ./...` |
| A6 — `schema_version` constant home | 1, 2 | **VERIFIED** | `evidence/spikes/a6-schema-version-home.md`; `go build`/`go vet` clean, no cycle; predicted 0023 golden byte-identity break fired as A3 records |
| A7 — CI snapshot of seams + severities | 1, 2 | **VERIFIED** | `evidence/spikes/a7-ci-snapshot-check.md`; golden-file `go test` in `graphlint_test`, red on `graph-vacuous-atom` `info`→`blocking`, green on revert |
| Spikes named but unrun | 3 | none owed | `spikes_unrun=[]` before and after |
| Exactness-word delta | 4 | none owed | `rdr lint 0029 \| grep prose:exactness` → zero findings |
| 3amigo F4 — §step-2 grep scope | absorption audit | **ALREADY ABSORBED** | current §step-2 and S9 scope the grep to C4-tiered vocabularies and enumerate the sites; S9 states outright that an unscoped grep is not the assertion |

## The A5 refutation, and why it is not a route-back

The spike refuted A5 for the CLIError `code` row: a `clierr`-side registry
drawing on the real `Code*` constants is a hard `go build` cycle, and the
only leaf-preserving shape that compiles is a disconnected shadow list that
cannot back a by-value `append-only` assertion.

This does not reopen the design, because **C4's seam clause authorized both
outcomes in writing before the spike ran**: "If that spike shows either seam
cannot be built at the layering this repo keeps … the tier on that row
STANDS and the row records `seam: none (prose-only)` with the reason … Which
of the two outcomes lands is A5's to report before lock — not Phase 1's to
discover." The spike reported which branch landed; no contract was amended
to accommodate it, no dependent clause went partial, and no peer's locked
text was touched. `§ground-before-ask` resolved `apply` (settled in source,
recommend with the cite) rather than `ask`.

The residual exposure is recorded rather than absorbed: the CLIError `code`
vocabulary is the surface consumers branch on most, and through this RDR its
`append-only` tier is a prose promise with no mechanical assertion behind
it. A7's snapshot cannot cover it, because there is no seam to read.

## Amendment sweep (rdr-common §amendment-sweep)

A5's outcome changed C4's seam obligation; seven sites carried the change:

1. C4 census — CLIError `code` row → `seam: none (prose-only)` with the reason.
2. C4 census — `flow next` row → seam confirmed buildable, shape named.
3. C4 contingency clause — rewritten from conditional to recorded outcome.
4. C3 mechanical backing — names the one vocabulary the check cannot reach.
5. Capability Dependencies — seam row (five of six buildable) and
   release-notes row (A7 Pending → Verified).
6. Risks & Mitigations — A7 residual resolved; A5's prose-only gap stated.
7. Phase 1 Step 3 — "six vocabularies" → the five that admit a seam.
8. S7 — the CLIError `code` set has no seam to assert over; scoped out on
   the positive side, explicitly in scope on the negative.

## Completeness check

No `_Draft placeholder._` survives in any body section; no seed-skeleton
header remains. The five `placeholder:survived` lint findings are all in
§Finalization Gate (1747-1870), which Stage 7 authors — not a hollow body.
`## References` carried the template's four bracketed placeholders and was
filled this pass from citations the RDR already held (SemVer §4/§8, the
in-repo census surfaces, repo documents, peer RDRs 0005/0006/0017/0022/0023/0028
and JDR 0001 §D10/§D12, the external prior-art roster, and this record's own
evidence artifacts). Collection, not research.

Two advisories remain and both are Stage 7's to answer, not Stage 6's:
`evidence:over-budget` on A4 (judged at the Gate) and `gate:inline`
(resolved by the lock pass's `sections` move).

## Verdict

**RECONCILED.** All seven Critical Assumptions are terminal
(`ca_verified=7`, `ca_pending=0`), `spikes_unrun=[]`, no BLOCKER, and no
MVV-critical item deferred past lock. `rdr lint` PASS, blocking=0,
resolution=0.
