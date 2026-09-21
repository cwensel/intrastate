# Finalization Gate — cli/0012 Declared-kind carrier at the guard value seam

- RDR: `cli/0012` — `0012-declared-kind-carrier-at-the-guard-seam`
- Date: 2026-09-20
- Verdict: **READY** — Gate PASS
- Mechanical pre-sweep: `evidence/tooling-pass/tooling-pass.md` — PASS
  (one BLOCK, C9 evidence-field budget on A4, fixed in-pass by relocating
  the derivation to `artifacts/a4-evidence.md`).

## 1. Contradiction Check

No contradictions between Research Findings and the Proposed Solution.

The one that existed was found and closed in this pass. §Key Discoveries'
match-path bullet still described the PRE-RULING state — "lint parses it at
the guard seam, the kernel byte-compares it at `TagSet.matches`", closing
with "A4 carries the open question of whether the held leg must be closed
too … or the seam's match arm must byte-compare." A4's reconcile ruling
(option (ii), absorbed 2026-09-20) settled exactly that question in the
second direction, and C4 implements it by removing
`reach.go::atomAdmitsValue` from the typed construction sites. The
discovery has been rewritten to the settled state: the match path does not
diverge, because it is excluded from typing rather than canonicalized, and
what survives is a canonicalization gap identical on both sides, carried by
kata `intrastate#ch99`.

Planned features vs stated principles: the record states one asymmetry —
one declared kind, two comparison rules — and does not hide it. C2 states
it on the comparison side, C4 on the construction side, §Load-Bearing
Decisions argues it against the Non-Redundancy/SPOT objection rather than
dismissing it, and the user-facing docs are required to carry it. The
principle that typing applies "exactly where lint and the kernel share a
decision path" is a structural fact about the code (`matchSatisfiable`
filters `BlockMatch`; `nodeMeetsAll` reads `model.Terminal`), stated and
testable, not a convenience boundary.

`0007:A17`'s headline is narrowed by C2 ("needs no declaration it was not
built over"). That is recorded as an explicit adjudicated deviation in
§Approach, licensed by A17's own Evidence and If-wrong, not left for a
cluster pass to discover as a contradiction.

## 2. Assumption Verification

7 of 7 Critical Assumptions terminal; `ca=all-terminal`, 0 Pending, 0
Unverified. Every record is internally consistent: Status, Method and
Evidence agree, and no "If wrong" is empty.

- Method vocabulary: all sanctioned (Peer RDR, Source Search, Spike,
  Design Decision); `method.off_vocabulary` empty on every row. No
  assumption carries `Docs Only`, so no load-bearing Docs-Only record
  needs a Spike/Source-Search plan.
- No self-reference: the Source-Search rows (A1, A2, A3, A4, A5) anchor on
  `internal/…` symbols and on `evidence/spikes/`, `evidence/reconcile/`
  paths. None resolves to the record itself or to `artifacts/`.
- Symbol resolution: `anchors_total=84`, `anchors_unresolved=0`,
  `anchors_unlooked=0`, `peer_evidence_unresolved=0` against the source
  root — every cited `path::Symbol` resolves on `main`.
- A4 was re-ruled at reconcile (option (ii)) and its Verified stamp
  restates the narrowed claim rather than the retired one; its Evidence
  derivation now lives at `artifacts/a4-evidence.md` with the field
  keeping every load-bearing anchor and the pointer.
- A7 is RETIRED, correctly: option (ii) deleted its subject (the
  three-valued widening of graphlint's atom check), rather than leaving a
  stale Verified stamp over a design that no longer exists.
- Status consistency: no settled-fact prose rests on a Pending or
  Unverified assumption — there are none. The one place prose had drifted
  from a settled assumption was the §Key Discoveries bullet fixed under
  item 1.

## 3. Scope Verification

The Minimum Viable Validation is in scope and executed during
implementation, not deferred. It is five concrete steps over one fixture
(`mvv-int`), sharing the model with fixtures S8 and S9 so every step reads
one guard:

- **The specific test/proof** — step 2: with an OWNED reader returning
  `many` for an `int`-declared `iter` against a row guarded `iter eq 7`,
  `flow resolve --outcome <o>` refuses `flow-guard-unevaluable`, naming the
  guarded row and the `iter` atom with reason `uncomparable`, where today
  the same invocation silently prunes the guarded row and plans the
  fallback. That invocation is the seed defect (kata `intrastate#cq5p`)
  refusing loudly instead of misrouting.
- Step 3 witnesses the parsed-comparison leg (`07` matches where it was
  pruned) with `4` as the control proving it is not blanket-true; this leg
  is on the GUARD path and so is unaffected by C4's match exclusion.
- Step 4 is `TestGuardEvaluatorContract` extended per C3; step 5 is C5's
  load refusal on `n eq "00"` (fixture S7).

The venue is the owned ingress because it is the only door reaching the
seam unconformed; the CLI `--tag` path cannot host it (`--tag iter=many`
refuses `flow-tag-invalid` upstream, measured before and after — fixture
S8). The reader is driven by the persisted artifact, not a test double, so
the MVV exercises the real `OwnedSnapshot` → `resolve::assemble` path.

Blast radius: this record declares no Cluster siblings (`clustered=false`,
`impact_families=none`), so no `impact.md` is owed and no peer's shipped
REQ is predicted to re-cut. The three open peers sharing anchors with this
record (0013, 0022, 0030) are each cross-cited and adjudicated — see item 1
and §Load-Bearing Decisions.

## 5. Proportionality

Right-sized. No section is flagged for trimming before locking.

**Contract count (the split test):** this record authors exactly ONE
independent load-bearing contract — C1, the declared-kind carrier at the
value seam. C2, C3, C4 and C5 each declare `Surface — of C1` and are
projected as such (four `surface-of` edges to `0012:C1`). One seam, not
several locked together; no split is warranted.

**Profile re-validation:** `foundational` is correct and matches the
contracts just counted. C1 is consumed by RDR 0030 and by every
`GuardEvaluator` implementer, it is user-facing (a new terminal
`flow-guard-unevaluable` state for owned-value callers), and it locks
cross-RDR. The lenses that ran agree with the field rather than
contradicting it: 3amigo, critique (two models, `critique_models=differ`),
repeatability (variant `full`, runs 1-3 + diff) and cove all ran, with
`lens_findings_open=0` and `lens_stale=none`. The `floor` row raises
nothing (`seam_lineage=0`, `accretion_disposition=false`). The field's form
is value + one clause naming the contract, with no matrix or provenance
prose left from the template.

**Length:** the record is long, and the length is load-bearing rather than
padding — the bulk sits in C2/C3/C4, which carry the guard/match split
that is this design's single most reversible-looking and most consequential
boundary, and in the Load-Bearing Decisions that argue it against the
Non-Redundancy objection. Trimming those is what would make the split read
as arbitrary to an implementer. Two reductions were nevertheless taken this
pass: A4's 34-line Evidence field is now an anchor-bearing pointer to
`artifacts/a4-evidence.md`, and C4's peer-seam paragraph was reduced to a
one-line pointer with the disposition moved to §Load-Bearing Decisions
where it is citable.

**Advisory carried:** nine `prose:exactness` hits on "canonical". Not
loose — C5 defines it operationally as `strconv.Itoa(n) == authored`, and
it is pinned by named fixtures S6 (123-model corpus diff) and S7 (`n eq
"00"` regression), both in the MVV and the conformance suite.
