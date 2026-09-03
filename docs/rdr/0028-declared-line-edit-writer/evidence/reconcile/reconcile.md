# Stage 6 — Reconcile (RDR 0028)

Model: claude-opus-5[1m]

## Stage 5 preflight

| Check | Result |
|---|---|
| `--outcome lens` | `emit.next = /rdr-reconcile`, `emit.row = none` — the large row (grounding → 3amigo → critique → repeatability) is complete |
| `--outcome critique` | `none` — a second pass was taken and diffed; `critique_models=differ`, no single-model fallback |
| `--outcome repeatability` | `none` — variant `lite`, run-1 and the diff written |
| Determinacy add-on | `Determinacy: fired` is written in Normative Contracts; run 1 executed under the lite variant |

No caveat to carry: both cross-model lenses are stamped and complete.

## Open set

Built from the four sources.

1. **Pre-Lock needs-verification list** — A2, A10, A11, A12 (A12 introduced by the
   pre-lock pass itself when C1.3 `re-anchor:` was strengthened from cardinality
   to identity).
2. **Still-Pending assumptions** — the same four (`ca_pending=4` at entry; A1,
   A3–A9 already Verified).
3. **Named-but-unrun spikes** — `spikes_unrun=[]`. A3, A4 and A5 ran at Stage 4
   with captured output under `evidence/spikes/`; no findings file names a spike
   the RDR does not.
4. **Exactness-word delta** — `rdr lint 0028 | grep prose:exactness` returns
   nothing, before and after this pass.

Plus the absorption audit's residue over the four lens rounds.

## Dispositions

| Item | Source | Disposition | Evidence pointer / plan |
|---|---|---|---|
| A2 — reader resolved before `Apply`; gate reachable there | 1, 2 | **VERIFIED** | `internal/accessor/executor.go::Executor.Write` / `::protectedKeys`; `internal/cli/flowbind/registry.go`; `internal/cli/cmdbind` gate inside `spawn` |
| A10 — gate reaches the executor as a field on `accessor.Registry` | 1, 2 | **VERIFIED** | `internal/cli/flowbind/registry.go` sole production site; all 12 `accessor.Registry` literals keyed; `internal/cli/flow_exec.go` sole non-test caller |
| A11 — deadline arm stays 0004's, unchanged | 1, 2 | **VERIFIED** | `internal/accessor/executor.go` deadline arm; `internal/table/load.go` `accessorTable` refuses absent `timeout` before carrier branching |
| A12 — re-anchor identity is index arithmetic, not a diff | 1, 2 | **VERIFIED** | C1.2 `value shape:` owns the `\n` refusal; no production newline guard exists to reuse (`cmdbind::substitute`, `flowbind` load/save, `table/load.go` all swept) |
| 3amigo iter-2 #11 — unreadable-target absent from Failure Modes | absorption audit | **VERIFIED (absorbed)** | new Failure Modes bullet added; S27b and the `disposition` row already existed |
| 3amigo Charted — `${2}` group-reference consistency | absorption audit | **VERIFIED (no edit owed)** | braced `${N}` used at every site; bare `$1` refused by C1.2 `vocabulary:` and tested by S14 |

### A2 — the one refinement

A2's sub-claim that the gate is "inspectable at the resolution point" holds, but
only through a type assertion to the concrete `cmdbind.Reader`: the abstract
`ReadBinding` interface exposes no gate accessor. That is not a refutation — it
is the reason the SITE is A10's field on `accessor.Registry` rather than an
interrogation of the resolved reader. A2's Evidence and the Decision Rationale
"Documented" bullet now say so; C1.3 `read-back:` already named the registry
field as the site.

### Absorption audit — rounds folded in

All four rounds (grounding, 3amigo ×2, critique ×2, repeatability) are absorbed.
Verified independently against the draft: C1.3's `precedence:` clause, the
`<clear>`-before-selection ordering, the per-use-site value-shape scope, S28's
map-order retraction, and the `regexp.Expand` disownment in Performance
Expectations. Three Charted items remain correctly routed to named successor
RDRs (symlink-resolution generalization, literal-anchor dialect, round-trip
carrier-completeness guard) — deferred by design, not residue.

One item is a live Finalization Gate question, not Stage 6 residue: 3amigo PM-5
asked whether Decision Rationale's premortem/ground-sweep/joint-check ledger
should be trimmed; `G-proportionality` answers contract-count and profile but
does not rule on the ledger. It carries to Stage 7.

## MVV floor

MVV step 5 (the gate-off pipeline refusing before mutation, both files
byte-identical) is exactly what A2 and A10 pin, so neither was eligible for
DOWNGRADE past lock. Both are Verified, so the floor is met rather than waived.
A11 and A12 are not consumed by any MVV step; both verified anyway.

## Completeness

- No `_Draft placeholder._` and no seed-skeleton header survives in any body section.
- `## References` is authored — peer RDRs, JDR 0003 §D1, source anchors, prior
  art, trackers and evidence paths. No bracketed template text.
- The five `placeholder:survived` lint findings all sit inside the Finalization
  Gate (1427–1572), which Stage 7 owns and has not yet answered
  (`gate_written=false`). Not a Stage 6 hollow-body failure.
- New advisory introduced: A2's Evidence field is now 35 lines against a soft cap
  of 30 (`evidence:over-budget`). The lint's own fix text routes this to the Gate
  and says not to truncate — the mass is verification content. Carried to Stage 7.

## Post-state

`ca_total=12  ca_verified=12  ca_pending=0  ca_unverified=0`; lint
`blocking=0 resolution=0`.

## Verdict

**RECONCILED** — every item terminal, no BLOCKER, no refutation. Ready for Finalize.
