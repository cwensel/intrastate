Model: claude-opus-5[1m]

# Stage 6 reconcile — cli/0023 resolve-envelope projection

Date: 2026-08-29

Verdict: **RECONCILED** — all nine Critical Assumptions terminal
(`ca=all-terminal`, 9 Verified / 0 Pending), no BLOCKER, no route-back.

## Stage 5 preflight

The routing model reports the `foundational` row complete — cove →
3amigo → critique → repeatability all present, `emit.next:
/rdr-reconcile`. The two completion outcomes a folder cannot answer
both return `none`: critique has two differing base-model stamps
(dual-model obligation discharged), repeatability ran the `full`
variant with run-1/2/3 and the diff written. No caveat to carry into
Caveats — no single-model fallback, no unstamped pass, no variant
mismatch. Determinacy is not at issue: the trigger is a `mid`/`large`
judgement and this record is `foundational`, where `repeatability` is
already a row entry.

## Open set

Built from the four sources. Sources 1–3 converge on the same three
assumptions; source 4 contributed the M-11 type-inspection claim.

| # | Item | Source | Disposition | Evidence / plan |
|---|------|--------|-------------|-----------------|
| 1 | A5 whole-tree registration oracle implementable | 1, 2 (cove flipped it Verified→Pending) | **VERIFIED** | `evidence/spikes/a5-total-walker.md` — total walker enumerates 15 commands incl. `completion` + 4 shell children; partial walker reaches 9 |
| 2 | A8 strict width satisfiable in TEXT mode | 1, 2 (3amigo booked it new) | **VERIFIED** | `evidence/spikes/a8-text-width.md` — 267→130 B (51.3%), 338→154 B (54.4%), 368→207 B (43.7%) |
| 3 | A9 echo containers never nil under pointer conversion | 1, 2 (critique booked it from M-9) | **VERIFIED** | Single construction site `flow_resolve.go::runFlowResolve`; three branch-free `make` writers; oracle half carried to Phase 1 as construction |
| 4 | M-11 A3's stamp covers count/keys but not `reflect.Kind` | 4 (claim introduced by the A2 fix pass) | **ACCEPTED — non-finding** | Sole `reflect.TypeOf(resolvePayload{})` site asserts `NumField`/`FieldByName`/adjacency, never `.Type`/`.Kind()`/`.Tag`; A3 holds unqualified |
| 5 | M-17 version skew analysed one-directionally | 4 | **ACCEPTED — recorded** | Rollback arm added to F5: same parse refusal fires; `0011:A14` message limit inherited, not minted |
| 6 | M-22 joint-check names 0012 but not 0018 | 4 | **ACCEPTED — recorded** | 0018's anchors are `kernelResolveFailure`/`escapeShapeBreaches`, both refusal-path; disjointness is structural via A6 + C1 |

### Why A5/A8/A9 were closed rather than downgraded

Each pinned something C1 asserts as a MUST or the MVV consumes, so
the no-defer-past-lock rule applied to all three rather than the
survivable-downgrade rule:

- **A5** — C1 states the whole-tree scope as a MUST and carries its
  own pre-lock branch ("if that spike cannot reach this scope, C1's
  structural negative needs its own form"). Deferring would have
  locked a contract whose satisfiability was unknown. The crux was
  narrow and checkable: whether `InitDefaultCompletionCmd` is
  reachable outside `ExecuteC`. It is exported, so the oracle is
  writable to full scope and the contract stands at full strength.
- **A8** — MVV step 3 and C1's both-modes width clause consume it.
  It was one measurement away from settled, and the measurement
  confirmed it with margin (text reduction exceeds JSON's on every
  shape).
- **A9** — guards C1's byte-identity clause on DEFAULT-mode output,
  i.e. callers who never opt in. The enumeration half was a bounded
  source search over a single construction site.

## Residual construction obligations

Not deferred assumptions — oracles to write, neither gating lock:

1. **S4 total walker** (Phase 2), with both preconditions asserted:
   auto-generated commands materialized via cobra's own initializers,
   and `completion` present in the walked set.
2. **A9 empty-container assertion** (Phase 1, BEFORE the pointer
   conversion lands): a default-mode oracle over a request with no
   `--tag`, no owned keys, no readers, asserting `{}`/`[]` and no
   `null`. Rides S2's absence control.

The A9 invariant is additionally guarded on landing by three shipped
wire oracles requiring the five echo keys present in default mode
(`decision_table_0010_test.go::wireKeyOrder`, `::TestReq41…`,
`flow_resolve_0005_test.go::TestReq76…`). They do not discharge the
Phase 1 oracle — none is written as an empty-container assertion,
which is the arm A9 names — but the conversion cannot land silently
broken.

## Charted to successors (not this RDR's residue)

Confirmed still correctly out of scope, no action at reconcile:
M-3 adoption/migration (kata `intrastate#srz2` unblocked, not
closed); M-18 traffic mix unsized; M-5 `revision` activation
residual; D8 C1 contract granularity; the `help --all` second
registration (pre-existing `0011` condition this RDR cites but
neither creates nor widens); the `0005` multi-class `escape_class`
gap, left as found.

M-24 (record length / proportionality) is Stage 7 gate work — the
`§proportionality` response must count contracts and re-validate the
`foundational` Profile. Flagged here, not answered.

## Amendment sweep

A5 and A8 flipping to Verified changed conditional wording in the C1
normative block and downstream. Sites updated in the same pass:

- C1 width clause — "rests on A8, presently Pending and unmeasured"
  → both halves measured, clause binds unconditionally.
- C1 implementability paragraph — "A5, presently Pending… written for
  the first time at Phase 2" → verified by spike before lock,
  construction remains Phase 2.
- §prerequisites — the "ONE carried exception" bullet replaced; all
  nine Verified before Phase 1, with the two construction
  obligations listed separately.
- §desk-trace steps 3 and 5 — witnesses upgraded from projected to
  measured (A8's bytes; the 15-vs-9 walk).

Anchor correction made in the same pass: the 3amigo note cited
`help_all.go::InitDefaultHelpCmd`, which is cobra's method, not a
repo symbol. The repo's force-init is
`internal/cli/help_all.go::registerHelpAllOnTree`. The behavior the
note described is unchanged and was re-confirmed by the bare-root
probe.

## Completeness check

No `_Draft placeholder._` and no seed-skeleton header survives in any
body section. `## References` is authored — JDR/peer-RDR citations,
the output contract, langref peer-CLI prior art, literature, and all
four spike artifacts.

`rdr lint` PASSes. Remaining findings are Stage 7's own: `gate:inline`
and five `placeholder:survived` blocks confined to the Finalization
Gate section, which Stage 7 authors. The five `prose:exactness`
`byte-identical` hits each have an Evidence Record — A2's normative
fixture (the "Projected reference" line) plus MVV step 1's pre-change
golden — and are covered by S1/S2/S5's Expected clauses; they are not
a post-mutation delta.

## Verdict

**RECONCILED.** Ready for Finalize.
