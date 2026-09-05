Model: claude-opus-5[1m]

# Stage 6 Reconcile — cli/0020 undeclared `--tag` key admission

Preflight: `--outcome lens` → `/rdr-reconcile` (`lens-mid-row-complete`, Profile
mid, grounding → 3amigo row complete). `--outcome critique` → `none`
(`critique-off-row`, not on the mid row). `--outcome repeatability` → `none`
(`repeatability-lite-complete`; run-1 stamped `variant: lite (profile: mid)`,
diff written). Determinacy line reads `fired` and the repeatability row emitted
no chain, so the add-on is satisfied by the completed lite pass. Stage 5 done.

## Open set

Built from the four sources. It is EMPTY of untermina items — recorded here as
the audit trail, since an empty set is a finding, not an absence of work.

**Source 1 — Pre-Lock needs-(re)verification lists.** No assumption was flipped
back to Pending by any round, and no round introduced a net-new A-N claim. The
grounding round corrected the *evidence text* of A1/A2/A3 without disturbing
their verdicts (each correction made the negative stronger, not weaker); the
3amigo and repeatability rounds raised no assumption-level finding.

**Source 2 — Pending/Unverified Critical Assumptions.** None.
`ca_total=5 ca_verified=5 ca_pending=0 ca_unverified=0 ca_placeholder=0
ca_off_vocabulary=0 ca=all-terminal`, `ca_pending_ids=[]`.

**Source 3 — `spikes_unrun`.** `[]`. The HEAD spike ran eight cases (a, a2, b,
c, d, e, f, g) plus the selection diff; all captures are on disk under
`evidence/spikes/`, and the absorption audit found no findings file naming a
spike the RDR does not.

**Source 4 — exactness-word delta.** `rdr lint 0020 | grep prose:exactness`
returns three hits, all "byte-identical" at C1 (366, 382). Both claims are
covered by named normative fixtures read at HEAD, not asserted:
- :366 — the hoisted arm's scalar-shaped message for a carrier is
  byte-identical to HEAD → normative fixture **F4**,
  `evidence/spikes/d-undeclared-empty.txt`, and MVV step 5 pins it.
- :382 — a declared scalar's empty value routes through either site
  byte-identically → `evidence/spikes/e-declared-empty.txt`.
Fixtures d and e carry the identical message template differing only in the key
name, which is exactly what the claims assert. No delta owed.

## Absorption audit (delegated)

One sub-agent over the three lens dirs that exist (`grounding`, `3amigo`,
`repeatability`; `critique` and `cove` are off the mid row). Verdict PASS,
non-blocking, no residue. Per round:

- **grounding** — G1 (A1's "one asserting test" miscount), G2 (A3's "exactly one
  place"), G3 (A2's omitted `conflicting` reader), G4 (A1's fixture described as
  "declared scalar" where `kind = "enum"`) all absorbed into the current CA text
  and the §key-discoveries repeat of the G1 count corrected alongside.
- **3amigo** — P2-M1 and P3-L1 fixed in D-selection-predicate / the illustrative
  code / S1; P2-M2, P2-L2, P3-L2, P3-I1 dismissed-with-cite; P2-L1 charted to a
  successor kata (a `codeInvalidFor` comment defect outside 0020's seam).
- **repeatability** — Finding 1 resolved. See below.

## Repeatability Finding 1 — the one round finding worth its own row

The lite run flagged a genuine internal fork: for the concrete input
`--tag labels=` against a model declaring `labels` as a set, C1's prose plus an
unconditional `value == ""` guard in the illustrative code pointed at the
scalar-shaped message, while the Disposition table's separate
"declared set key, empty value" row and HEAD fixture g pointed at the
set-specific conformance message. The spike report had already flagged the same
limit as a LOUD FLAG ("declaration-independent for scalars but NOT across
kinds"), and rulings.md Q1 ruled the disposition at Stage 4.

The current text adjudicates it in four mutually consistent places:
- C1: "The hoisted arm is therefore NOT unconditional on `value == \"\"`: it
  fires for a key that is undeclared, or declared with a NON-set kind, and it
  must not intercept a declared SET key" — citing normative fixture **G**.
- Illustrative code: `if value == "" && (!declared || !decl.IsSet())`.
- Disposition table: the two rows carry distinct messages.
- Desk trace row 5b: "resolved by C1's guard", naming the unconditional test as
  the regression it would cause.
- S5 is the pinning witness ("An unconditional `value == \"\"` test ahead of the
  declared lookup fails this scenario, which is what makes it the guard's
  witness"), against fixture `evidence/spikes/g-declared-set-empty.txt`.
Failure Modes carries the same carve-out. No drift against the Q1 ruling.

## Completeness check (grep, not judgment)

- No `_Draft placeholder._` and no `this is a seed skeleton` header survives.
- `## References` is fully authored — source paths, `docs/cli-output-contract.md`,
  peer RDRs 0003/0005/0008/0010/0012/0024, spike artifacts, katas. No template
  brackets.
- The five `placeholder:survived` lint hits are all inside the Finalization Gate
  (949–1039), which Stage 7 owns; no body section this stage owes is hollow.
- `rdr lint 0020`: PASS, blocking=0 resolution=0.

## MVV floor

Nothing is deferred that the MVV depends on. The one leg the MVV rests on that
HEAD cannot witness — an undeclared ARRAY carrier leaving selection unchanged —
is disclosed as such in A2's Evidence, in S1, and in the spike report, and is
covered by MVV 2/3 at implementation. That is not a deferral of a pre-lock
prerequisite: it is unverifiable-before-the-change by construction (the array
refuses at HEAD), and the scalar leg's runtime proof is captured
(`selection-diff-a2-vs-b.txt`, empty).

## Table

| item | source | disposition | evidence pointer / plan |
| --- | --- | --- | --- |
| A1 no caller load-bears on the undeclared-array refusal | 1 | VERIFIED | Source Search; `internal/cli/flow_input.go::canonicalValue` sole producer, zero asserting test repo-wide; grounding G1/G4 corrections absorbed |
| A2 an undeclared value cannot influence resolution | 1 | VERIFIED | Source Search + Spike; five `CatUnknownTag` sites, four kernel-view readers incl. `conflicting`; `evidence/spikes/selection-diff-a2-vs-b.txt` empty |
| A3 no byte-equality obligation reads an undeclared value | 1 | VERIFIED | Source Search; two comparison sites, both in the write-plan/re-read domain (G2 absorbed) |
| A4 external prior art aligns with the carrier arm | 2 | VERIFIED | Prior Art; three opened citations, `evidence/research/resolve-prior-art.md` |
| A5 0005's code table tolerates an unreachable arm | 2 | VERIFIED | Peer RDR; `0005:§failure-modes` header scopes to spellings + exit group |
| spikes named under `{SPIKE_DIR}` | 3 | VERIFIED | `spikes_unrun=[]`; eight captures + diff on disk, `evidence/spikes/spike-report.md` |
| "byte-identical" C1:366 (carrier empty-value message) | 4 | VERIFIED | normative fixture F4, `evidence/spikes/d-undeclared-empty.txt`; MVV 5 |
| "byte-identical" C1:382 (declared scalar empty value) | 4 | VERIFIED | `evidence/spikes/e-declared-empty.txt` |
| repeatability Finding 1 — declared-SET + empty value fork | 1 | VERIFIED (absorbed) | C1 guard + illustrative code + Disposition rows + desk-trace 5b + S5; fixture `evidence/spikes/g-declared-set-empty.txt`; rulings.md Q1 |

No DOWNGRADED items, no ACCEPTED-as-design-decision items owed, no BLOCKER.

## Verdict

**RECONCILED.** All items terminal, no refutation, no MVV-critical deferral.
No in-place edit to Critical Assumptions was owed: every disposition was already
written into the RDR by the round that produced it, and this pass verified the
RDR and this report agree row for row.
