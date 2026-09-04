# Stage 6 — Reconcile, iteration 2 (RDR 0028)

Model: claude-opus-5[1m]

Re-entry pass. The record was locked Final 2026-09-03, then re-entered from the
0026-0027-0028 cluster reconcile. Metadata Status:

`Draft [revised from Final 2026-09-03; re-verify A7, A9 @refine — JDR 0003 §D2 decided (a): C1.6 gains an argv0 exclusion and a `-`-prefixed value refusal]`

Iteration 1's report (`../reconcile.md`) covers the pre-lock pass and stands. This
iteration reconciles only what the re-entry disturbed.

## Stage 5 preflight

| Check | Result |
|---|---|
| `--outcome lens` | `emit.next = /rdr-reconcile`, `emit.row = none` — rule `lens-large-row-complete`; `lens_stale=none`, so no folder predates the demote date |
| `--outcome critique` | `none` — rule `critique-large-diffed`; `critique_models=differ`, no single-model fallback |
| `--outcome repeatability` | `none` — rule `repeatability-lite-complete`; run-1's `variant:` header reads lite, run-1 and the diff written |
| Determinacy add-on | no `resolve:determinacy` chain emitted; `determinacy=fired` is written in Normative Contracts, so the trigger is judged, not unwritten |

No caveat to carry.

## Open set

Built from the four sources, scoped to the re-entry.

1. **Pre-Lock needs-(re)verification list** — the re-entry's own list, carried on
   the Status qualifier and echoed by `reverify=["0028:A7","0028:A9"]`: **A7, A9**.
   No lens round of this iteration added an A-N claim.
2. **Still-Pending assumptions** — none. `ca_total=12 ca_verified=12 ca_pending=0
   ca_unverified=0 ca_off_vocabulary=0`; `ca_pending_ids=[]`.
3. **Named-but-unrun spikes** — `spikes_unrun=[]`. No iter-2 findings file names a
   spike the RDR does not.
4. **Exactness-word delta** — `rdr lint 0028 | grep prose:exactness` returns
   nothing (`lint-baseline.txt`: `blocking=0 resolution=0 placeholder=0
   advisory=2`). The amendment introduced no unbacked exactness term; C1.6's two
   new rules are each grounded at a source anchor through A7/A9.

Plus the absorption audit's residue over the re-entry-era rounds (3amigo iter-2,
critique iter-2, repeatability, and the cluster reconcile report that caused the
re-entry).

## Dispositions

| Item | Source | Disposition | Evidence pointer / plan |
|---|---|---|---|
| A7 — `--tag` is context-only, so `{tag.nnnn}` binds without a reader run | 1 | **VERIFIED** | `internal/cli/flow.go::registerTagFlag`; `internal/cli/flow_input.go::parseTags` — enforcement is KEY-scoped only (`canonicalValue` applies no flag-shape check), so C1.6's VALUE rule closes a real gap rather than restating an existing guard |
| A9 — `{tag.<key>}` is one more member at `{artifact}`'s substitution site | 1 | **VERIFIED** | `internal/cli/cmdbind/cmdbind.go::substitute` whole-element equality; `::refuse`'s existing `strings.HasPrefix(art.Path, "-")` branch returns before `resolveArgv0` and sits under `if !substituted`, so C1.6's VALUE rule mirrors a rule already there at the per-USE-SITE scope it claims |
| 3amigo iter-2 — disposition table missing the unreadable-target row | audit | **VERIFIED (absorbed)** | `§pre-lock-mini-checks` disposition row 14: `execution_failure`, Detail carries the OS error, the one C1 refusal with no `edit_*` token |
| 3amigo iter-2 — bare `regexp.Expand` citation | audit | **VERIFIED (absorbed)** | Performance Expectations now carries C1.2's disowning qualifier ("the rule matches `regexp.Expand`'s, but the parse is C1.2's own") |
| critique iter-2 E-1 — "sorted by key" contradicted C1.4's inherited unspecified rule | audit | **VERIFIED (absorbed)** | Performance Expectations "Map order" now states the inherited-unspecified rule; S28/S32 assert the category, not the table |
| repeatability #1–#3 — apply-time ordering, `<clear>` ordering, multiline scan scope | audit | **VERIFIED (absorbed)** | C1.3 `precedence:` (six-step fail-fast order); C1.5 `<clear>` decided on the RAW value before selection and `replace`; C1.2 `value shape:` narrowed to "a bound tag value THIS ENTRY's rules reference", per-USE-SITE |

### The four Direction obligations

The re-entry's Refinement Context named four; all four landed, verified against the
current record text:

1. **C1.6 `admission:` argv0 exclusion** — `edit_tag_argv0` at LINT, registered in
   C1.4's load-time category list and kept disjoint from the five `edit`-table
   categories by carrier. Statically decidable, so refused where it is visible.
2. **C1.6 `binding:` `-`-prefixed refusal** — `execution_failure` before spawn,
   Detail naming the placeholder, mirroring `cmdbind.go::substitute`'s `{artifact}`
   rule verbatim. Class named: argument injection (CWE-88), not shell injection.
3. **Cross-Cutting secret/credential response rewritten** — now states that C1.6 IS
   a new source of caller-supplied value on an execution path, and that A7
   establishes only that `--tag` cannot supply an OWNED value. The earlier draft's
   misreading is named and corrected.
4. **JDR 0003 §D3 cited at C1.3, carried into S27** — C1.3 `order:` cites §D3 (b)
   for the exit-2 group as a citation, not a rewrite; S27 asserts exit 2 under a
   distinct CLI code with the reason in `findings[]`, and forbids pinning the
   code's spelling.

No contradiction with siblings: C1.3 `precedence:` places both C1.6 preconditions
first as entry-level, C1.3 `order:` gives them the name-the-placeholder Detail rule,
the disposition table carries matching rows for both, and C1.6's `no other change`
line keeps 0027:C1's "argv WORDS only" promise true (JDR 0003 §D2 (a)).

## MVV floor

MVV step 2 binds `--tag nnnn=NNNN` on both halves of the pipeline and step 5 is the
gate-off refusal before mutation — exactly what A7 and A9 pin, so neither was
eligible for DOWNGRADE past lock. Both are Verified, so the floor is met rather
than waived.

## Completeness

- No `_Draft placeholder._` and no seed-skeleton header survives in any body section.
- `## References` is authored — peer RDRs, JDR 0003 §D1, source anchors, prior art,
  trackers, evidence paths. No bracketed template text.
- `placeholder=0`: the five Finalization Gate findings iteration 1 carried are gone
  (`gate_written=true`).
- Two advisories carry to Stage 7, neither a Stage 6 item: A2's `evidence:over-budget`
  (35 lines against a soft cap of 30 — the lint's own fix text routes it to the Gate
  and says not to truncate), and `parse:section:unknown-to-template` on the
  Refinement Context section, which is deleted on re-lock by its own heading.

## Post-state

`ca_total=12  ca_verified=12  ca_pending=0  ca_unverified=0`; lint
`blocking=0 resolution=0 placeholder=0 advisory=2`.

## Verdict

**RECONCILED** — every item terminal, no BLOCKER, no refutation. The re-entry's two
re-verify assumptions are Verified against the amended C1.6, and all four Direction
obligations landed. Ready for Finalize.
