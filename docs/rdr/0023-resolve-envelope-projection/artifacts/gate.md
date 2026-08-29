# Finalization Gate — cli/0023 resolve-envelope-projection

RDR: `0023-resolve-envelope-projection` · Date: 2026-08-29 ·
Verdict: **READY — Gate PASS**

Mechanical pre-sweep: `evidence/tooling-pass/tooling-pass.md` — PASS
(one mechanical §References finding fixed in-pass, sweep re-run clean).

## 1. Contradiction Check

No contradictions between Research Findings, the design principles, and
the Proposed Solution.

Three places where the record could read as self-conflicting were
checked and each is a stated reconciliation, not a conflict:

- **The seed's 1529 B / ~72% versus A1's 708→146 B.** §Problem Statement
  names these as two different calls and retires the seed figure as a
  headline; §Consequences quotes A1's range with its shape and says so
  explicitly. The Key Discovery marked "Assumed" (the seed measurement
  transfers) is the one A1 verified, and A1's verdict replaced it rather
  than sitting beside it.
- **"No shipped oracle blocks this" versus `0011`'s fences.** Research
  Findings establish that `0011:C2`'s absence oracles pin `--all` by
  NAME on the sibling verbs and that the `--select` rejected-spelling pin
  sweeps `flow next` only; the Overrides field states the same scope. The
  solution reopens 0011's rejection *ground*, not a fence, and C1 mints a
  strictly stronger structural negative than 0011 shipped.
- **`revision` carried while `model` is projected away.** C2 does not
  claim `revision` delivers provenance today — it records the field as
  contracted-but-vacant, assigns it on its definition (loader-produced,
  never restated from the request), and marks the always-keep clause
  forward-binding. The unattributability residual is charted to the RDR
  that gives `revision` a value (`evidence/critique/Charted.md`), not
  asserted away here.

The chosen design also honors the two locked principles it touches:
`0005:C1`'s both-modes agreement holds by construction (one projection
before `respond.OK`, one shipped renderer), and the "don't edit a
predecessor's surface" principle holds because the flag is additive and
opt-in (A3).

## 2. Assumption Verification

All nine Critical Assumptions are `Verified`; none is `Pending` or
`Unverified`, so no settled-fact prose anywhere in the record rests on an
open assumption.

Every record is internally consistent — Status, Method and Evidence agree
and "If wrong" is non-empty on all nine. Method labels are all in the
sanctioned vocabulary (`Spike` ×3, `Source Search` ×4, `Peer RDR`,
`Source Search + Spike`); `method.off_vocabulary[]` is empty on every
one. **No `Docs Only` record exists**, so nothing is blocked on that
ground.

No `Verified` stamp is self-referential. The four `Spike`-backed records
(A1, A2, A5, A8) cite spike artifacts under this RDR's own
`evidence/spikes/`, which is where spike output belongs; the five
`Source Search` records cite repo source (`internal/cli/flow_input.go`,
`flow_exec.go`, `flow_resolve.go`, `clierr/clierr.go`,
`internal/resolve/resolve.go`), none of them this record or its artifact
directory. A4's `Peer RDR` Evidence cites `0024:C4` — an element, not a
bare record. All 48 `path::Symbol` source anchors resolve on `main`
(`edges[]`: 48 `source-anchor`, all `resolved: true`, zero false).

Three records earned particular scrutiny because they closed late, at the
Stage 6 reconcile, and each closed on measurement rather than deferral:

- **A5** (whole-tree walker implementable) — the cove sweep refuted the
  general "no name registered twice" claim (`help --all` is the
  counterexample), which is *why* C1 mandates a total walker instead of
  reusing `help_all.go::walkCommandTree`. The spike then demonstrated the
  walker at C1's full scope: 15 commands enumerated including `help`,
  `completion` and the four shell children.
- **A8** (text-mode strict width) — measured on the same three shapes and
  in the same unit as A1: 43.7–54.4% saved, so C1's both-modes clause
  binds unconditionally rather than resting on an unmeasured half.
- **A9** (echo containers never nil) — the enumeration half is
  discharged: `resolvePayload` has exactly one construction site
  repo-wide. This is the one assumption whose failure S1 could not see,
  which is why it carries a Phase-1 oracle obligation ahead of the
  pointer conversion.

Two obligations carry into the build and are correctly typed as
construction work, not open assumptions: the S4 total walker (Phase 2,
both preconditions asserted) and A9's empty-container default-mode
assertion (Phase 1, before the conversion lands). Neither gates lock.

## 3. Scope Verification

The Minimum Viable Validation is **in scope and executed during
implementation**, not deferred. It is a six-step battery over the
checked-in fixtures (`models/examples/pricing-decision-table.toml`,
`release-grammar.toml`, `review-state-machine.toml`), landing as the S1–S5
oracles plus one deliberately manual step:

1. **The pre-change golden** — default-mode output captured from the
   pre-change binary and checked in as Phase 1's first step, on a clean
   tree, before the struct is edited. This is the only comparison in the
   battery with a pre-change side, and it is what makes C1's headline
   "byte-identical to today's" assertable at all; S1–S5 compare the new
   build against itself and cannot see a regression that moves both sides.
2. **S2** — the projected top-level key set asserted as an explicit
   literal list, with `escape_class` compared against the same run's
   default output rather than an unconditional literal.
3. **S5** — `--as=text` subset plus strict byte reduction, stable across
   repeated runs.
4. **S1** — refusal envelopes and exit codes byte-identical ± the flag,
   with the invoked-reader set asserted by *observed execution* at the
   `runReaders` seam (C1 rejects the recomputation form as vacuous).
5. **S4** — the whole-tree structural walk asserting the registrant set is
   exactly `{flow resolve}`, over a root with both auto-generated commands
   materialized and their presence asserted as a precondition.
6. **The width table** — default vs projected bytes on the motivating
   shape and one gate/write-heavy fixture, measured on the full emitted
   line, compared against A1's baseline. The pass bar is stated on shapes
   the implementer can actually run: ≥40% on each checked-in fixture. The
   79.4% synthetic figure is explicitly **not** a pass bar.

**S3** — the reflective partition-completeness oracle — carries C2's
enforcement: bidirectional equality between the declared assignment table
and the struct's field set, with the carrier constrained to a standalone
declaration site so the oracle cannot become tautological.

## 5. Proportionality

Right-sized on the contract test, which is contract count rather than word
count. **This RDR authors one independent load-bearing contract.** C1 is
the seam: `flow resolve`'s caller-controlled projection surface. C2 is not
a second seam — it is this verb's *instance* of a doctrine normative
elsewhere (JDR 0002 §D1), cited rather than restated, and it exists to give
C1's projection a principled membership rule instead of an omit-list. One
seam, two facets; nothing here to split.

**Profile re-validated: `foundational` is correct** and matches the
contracts just counted. The record adds a user-facing CLI surface on a
locked verb, narrows a predecessor's pinned envelope (`0005:A6`), instances
a joint decision (JDR 0002 §D1), and composes with a concurrent sibling
(`0024:C4`) — a cross-RDR producer, not a local change. The full
foundational lens battery ran and left evidence: `cove`, `3amigo`,
`critique` (two models plus diff), and `repeatability` at full strength
(run-1/2/3 plus diff), followed by the Stage 6 reconcile verdicting
RECONCILED. No lens is missing, so nothing routes back on the latch. The
Profile field's form is correct: value plus one clause naming the contract,
with no matrix or provenance prose left from the template.

**Length is the one thing to flag, and it is flagged rather than trimmed.**
At ~2000 lines for a boolean flag the record is long, and the critique lens
said so (M-24). Two considerations decide against trimming before lock.
First, most of the mass is load-bearing and was *added* by the lenses to
close real defects: C1's reader-measurand paragraphs exist because the
obvious oracle form is vacuous by construction (M-8), C2's carrier
constraint exists because the cheaper carrier defeats the check it carries
(M-14), and C1's walker-totality paragraphs exist because the repo already
contains the registration an enumerated sweep misses (`help --all`). Cutting
those re-opens the findings. Second, the one genuine granularity defect —
C1 carrying five distinct obligations under one id, which cost one
repeatability run its width paragraph — is **charted to a successor**
(`evidence/repeatability/Charted.md`) rather than fixed here, because
re-anchoring the split would touch every citing site (§disposition,
§oracle-discriminability, S1, S5, MVV step 6, the desk trace) at the moment
of lock. That is the right trade: the defect is legibility, not correctness,
and the charted successor covers contract granularity across the
0005/0010/0011/0023 family rather than this record alone.

Scope residue is charted, not carried: adoption of the flag at a real call
site (which is what closes `intrastate#srz2`), the traffic-mix measurement,
the modes-quantified always-keep oracle, and `revision` activation are all
recorded in `evidence/critique/Charted.md` as successor work with named
owners-by-shape.
