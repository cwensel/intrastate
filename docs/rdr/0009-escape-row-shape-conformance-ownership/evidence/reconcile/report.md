Model: claude-opus-5[1m]
Stage: 6 — Reconcile (RDR 0009)

# Reconcile report — RDR 0009 (escape-row shape conformance ownership)

## Stage 5 completeness preflight

`Profile: foundational` → required lens set is `cove 3amigo critique
repeatability`. All four evidence dirs exist and carry resolved output:

| Lens | Evidence | Resolution |
| --- | --- | --- |
| 3amigo | `consolidation.md`, `persona-1..3` | `resolution.md` — T-1…T-14 |
| critique | `critique.md`, `critique-modelB.md`, `diff.md` | `resolution.md` — M-1…M-15 |
| repeatability | `run-1/2/3.md`, `diff.md` | dispositions absorbed directly into the RDR (no `resolution.md`) |
| cove | `findings-main/a/b.md` | `ledger.md` — L-1…L-15 |

Repeatability ran the **correct variant**: `run-1.md` header reads
`variant: full (profile: foundational)`, matching runs 2 and 3, across three
distinct base models (`claude-opus-5[1m]`, `claude-fable-5`,
`glm-5.2:cloud`) with a `diff.md` authored by a non-run session. No variant
mismatch. Determinacy is a `mid`/`large` trigger and does not fire for
`foundational`; the RDR's Performance Expectations separately record that no
byte-stable hash, canonical serialization, or ordering guarantee is
introduced, so the determinacy checklist does not apply.

Preflight verdict: **Stage 5 complete.**

## Open set (four sources)

1. **Pre-Lock needs-verification lists** — critique/resolution.md carried A2
   (corrected evidence to re-confirm) and A9 (still Pending);
   3amigo/resolution.md raised A9 as net-new.
2. **Pending / Unverified Critical Assumptions** — A9 only. A1–A8 verified at
   Stage 4, A10 verified at Pre-Lock.
3. **Named-but-unrun spikes** — none. The only spike the RDR names is A3
   (`evidence/spikes/a3-fixture-conformance.md`), which has four captured
   `.out` artifacts in `{SPIKE_DIR}`. The `aggregation-probe` reference is a
   citation of RDR 0007's spike, not a spike this RDR owes.
4. **Exactness-word delta (post-mutation only)** — the rounds introduced or
   touched: "reports **every** breaching row", "**exactly** one level deep",
   "**verbatim**, never fmt.Errorf-wrapped", "exhaustive over the five kinds",
   "**identical** reported row sequence", "the table MUST be **mixed**", and
   A9's "**sole** wire cause surface". Each is bound by a Normative Contract
   clause plus a Testing Strategy scenario (9, 9, 4, 10, 10, and A9's
   evidence respectively); no exactness word entered without a record.

Absorption audit (delegated, read-only over the four lens dirs) confirmed the
rounds were folded in: 3amigo 14/14 ledgered findings absorbed, critique
15/15, cove 15/15, repeatability 7/8 (the eighth being A9, whose correct
disposition *was* "carry to Stage 6"). It surfaced two pre-disposition
residue items, dispositioned below as R-2 and R-3.

## Dispositions

| Item | Source | Disposition | Evidence / plan |
| --- | --- | --- | --- |
| **A9** — new `omitempty` field on `clierr.CLIError` carrying row identities is additive to RDR 0005's envelope | 1, 2 (3amigo T-3; critique M-6; repeatability C-1; cove L-10) | **VERIFIED** (with named amendment obligation) | Source search over `internal/cli/clierr/clierr.go`, `internal/cli/version_test.go`, `go list` import set. Four plan items closed — see below. Written into A9's Evidence Record. |
| **A2** — Go error surfaces via `GroupInternal`/exit 2 with no RDR 0005 contract change | 1 (critique: corrected evidence carried forward) | **VERIFIED** (re-confirmed) | A9's search closed the additive question; `ExitCodeFor`'s refusal-to-exit-code mapping untouched. Corrected exit-2 reading (`ExecuteAndEmit` → `cobraErrorToCLIError`) stands. Note added to A2. |
| **R-2** — References line still cited `ExitCodeFor` non-`CLIError` → exit 1 | absorption audit (critique M-8 tail) | **VERIFIED** (corrected in place) | The default *is* exit 1 in isolation but is unreachable for a verb error; References now states both halves. No contract change — the four load-bearing sites (A2, CLI clause, Prerequisites, sc.11) were already correct. |
| **R-3** — `CheckValid` on a zero-row / nil `Rows` table unstated and un-scenario'd | absorption audit (3amigo persona-3-qa, dropped pre-disposition) | **ACCEPTED** (Design Decision, pinned normatively) | Vacuous conformance: a predicate quantified over rows cannot be breached by no rows, and `errors.Join` of nothing is nil. Added to the predicate-export Normative Contract so a nil return there is never read as skipped validation. No spike owed — the property is definitional, not empirical. |
| **A1, A3–A8, A10** | 2 | already terminal (Verified) | Untouched by the rounds' edits; A8's *wording* was clarified by 3amigo T-1, not its verdict. |

### A9 — the four plan items

- **(a) No consumer depends on `Detail` being the sole wire cause surface.**
  `internal/cli/` has exactly one test file (`version_test.go`); no
  `testdata/`, goldens, or snapshots exist in the CLI tree. Its only JSON
  assertion (`TestVersion_JSON`) unmarshals into a partial anonymous struct —
  unknown fields ignored — and tests the `respond.Success` envelope, not
  `CLIError`. The two error-path tests assert on `clierr.ErrorCode` /
  `ExitCodeFor` and never touch serialized bytes. The break-shapes
  (whole-JSON comparison, `DeepEqual` over an unmarshalled map,
  `DisallowUnknownFields`) exist nowhere in the repo; every `reflect.DeepEqual`
  hit is in `internal/resolve/*_test.go` over kernel types. No non-Go consumer
  parses the envelope. **Breaks nothing, populated or not.**
- **(b) The allowance is real and conditional.** `clierr.go::CLIError` doc,
  verbatim: "Extend with new optional fields as needed — keep them
  `omitempty` so the envelope stays append-only and stable for tools." The
  `omitempty` condition is what the RDR's clause already requires. Tags
  confirmed: `Code`/`Message` unconditional; `Param`/`Detail`/`Hint` already
  `omitempty`; `Group` and `Cause` both `json:"-"`. `EmitJSON` is a plain
  `json.Marshal` with no custom `MarshalJSON`, so a tagged field reaches the
  wire with no other code change.
- **(c) The `Cause` doc comment contradicts the addition → amendment
  obligation.** Verbatim: "Cause preserves the underlying Go error for
  errors.Is/errors.As traversal. Not serialized — the wire-visible cause
  surface is Detail." The definite article is a uniqueness claim, and row
  identities are cause information, so the sentence goes false when the field
  lands. Documentary, not behavioral: nothing branches on it, and `Cause`
  stays `json:"-"`. **Recorded as a Prerequisite** — reword the second clause
  and update `docs/cli-output-contract.md`'s field list in the same change.
  The comment's normative home is RDR 0005, but amending a stale *code
  comment* to match an additive field is not reopening 0005's envelope
  contract, which is what A2's verdict turns on.
- **(d) `clierr` stays a leaf; the two open sub-questions settled.** `go list`
  gives the complete import set as `encoding/json errors fmt io` — four
  stdlib packages, no intra-repo import — and the package doc records the
  leaf property as deliberate (avoiding an import cycle). Nothing forces a
  `resolve` import. *Settled*: the carrier is a **clierr-local
  representation** (plain strings or a small row-identity struct declared in
  `clierr`), never `resolve.RowRef`; the CLI verb layer, importing both, does
  the conversion. *Settled*: the per-identity `Count` **does serialize** —
  carrying it structurally in the kernel and dropping it at the wire would
  force exactly the prose re-parse this RDR forbids, for the degenerate
  `RowRef{"",""}` case the multi-breach clause calls out.
  Caveat for the implementer: `internal/resolve/resolve_test.go::TestReq11_KernelImportsNoCLIOutputOrPersistenceFacility`
  guards one direction only (kernel must not import `clierr`), so the leaf
  property is preserved by design intent, not by a test.

## Hard-rule checks

- **Refutation** — none. No spike or source search refuted an assumption the
  RDR relies on. A9's search *confirmed* its claim; the one friction it found
  (the stale `Cause` comment) is a documentation defect with a same-change
  remedy, not a design refutation. No return-to-stage is triggered.
- **MVV dependency** — nothing deferred past lock that the MVV depends on. The
  MVV's three normative fixtures exercise `Resolve`, `Table.CheckValid`, the
  typed error, and the aggregate — all kernel-side, all covered by scenarios
  1–7b, 9, 10, 10b. A9's subject (the CLI wire carrier) is scenario 11, which
  the RDR already marks **DEFERRED: not executable by this RDR** because no
  `flow` verb exists at HEAD (A5) — a deferral of *execution*, not of
  verification: the assumption itself is now Verified and the obligation is
  bound in Prerequisites.

## Completeness check

- No `_Draft placeholder._` survives (grep: 0 hits).
- No `this is a seed skeleton` header remains.
- `## References` carries no bracketed template placeholders; it is fully
  populated from citations the RDR already held.
- Every `**Method**:` label is in vocabulary (Source Search ×7, Spike ×1,
  Peer RDR ×1, Prior Art ×1). **No `Docs Only`** on any load-bearing claim.
- Every spike named in the RDR has captured output in `{SPIKE_DIR}`.

## Verdict

**RECONCILED.** All ten Critical Assumptions are terminal (A1–A8, A10
Verified before this stage; A9 Verified here). No BLOCKER. The two
pre-disposition residue items found by the absorption audit are dispositioned
in place. Ready for Finalize.
