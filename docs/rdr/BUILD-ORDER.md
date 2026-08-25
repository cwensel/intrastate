# Implementation build order — the `0002-0009` cluster

**Status**: authored 2026-08-24, after JDR 0001 settled (`5c2b96b`) and the
joint-decision qualifier sweep (`730319c`). Read this before `/rdr-implement`
on any cluster member.

This document encodes **build order** (which Go code compiles first). It is a
different axis from the **lock order** the cluster gate reasoned about (which
document re-locks first), and the two run in opposite directions at the head of
the graph. Conflating them is what produced the contradictory orderings in
`cluster-reconcile/0002-0009/iter-3` and `iter-4`; this file exists so the next
reader does not have to re-derive the distinction from 2,300-line RDR bodies.

---

## The inversion, stated once

The iteration-3 gate wrote: *"0002 is the producer every other member's wire
format and fixtures read, so its re-lock is the practical gate for the set."*

**That is true for locks and for fixtures. It is false for Go types**, and 0002
itself says so in five places:

| Site | What 0002 says |
| --- | --- |
| `0002:2263` | Phase 2's handoff-routing assertion "is unsatisfiable until RDR 0007's reshape lands (Prerequisites)" |
| `0002:2278` | Phase 2's atom constants "compile only once RDR 0007's reshape exports them… the spike mirrors the constants locally because the kernel does not yet export them" |
| `0002:1651` | A Phase-2 scenario "waits on RDR 0007's reshape" |
| `0002:1662` | "a sequencing block, not a contradiction" |
| `0002:1772` | "The constants and the atom-shaped row do not exist in `internal/resolve` yet, so Phases 2–3 sequence behind that reshape" |

**0002 Phase 2 normalizes *to* `internal/resolve::Row`. RDR 0007 Phase 1
*defines* that `Row`.** 0002 Phase 2 cannot compile before 0007 Phase 1.

0002 remains the producer of the **wire format and fixtures** — a real
dependency on a different axis, discharged at its Stage 4 spike
(`0002/evidence/spikes/iter-2/`), not at its Stage 8.

## New information since the gate's ordering

Both prior orderings predate the home's final answers. The iteration-4 gate
predicted JD-21/JD-22 would be *answered by* writing the first Go type. Instead
JDR 0001 answered them first — §D12 and §D13, `d937eec`, settled `5c2b96b` —
and **§D12 assigns its landing site to "0007 Phase 1"** by name. That converts
0007 Phase 1 from "where a question gets answered" into "where a decided
contract lands", which is a stronger reason to run it first, and it was not
available when either ordering was written.

---

## Build order

### 1. RDR 0007 Phase 1 — Row and seam reshape · **head of the graph**

Replace `Row.Guard string` with the §D1 atom slice; narrow `GuardEvaluator` to
the per-atom value seam; export the existence token and literal constants;
migrate the frozen fixture suite to drive verdicts THROUGH atoms.

**Also lands §D12 (§JD-21)**: `Block` gains `BlockMatch` — one type, no
separate slice on `Row`. This widens a fenced cardinality ("exactly two
constants", `0007:1264-1266`), so it is a **normative change to a Final RDR,
not a citation** — see `0007/artifacts/deviations.md` D1. Expect the scoped
re-entry here; do not widen the type silently.

Unblocks: everything. Nothing else in the cluster compiles against the kernel
until this lands.

### 2. RDR 0002 — §D13 set-value encoding clause · **one clause, early**

`Tag.Value` stays `string`; a set crosses as its canonical JSON array —
members sorted, duplicate-free, compact. §D13 is explicit: *"Lands in 0002 —
one normative clause"*. See `0002/artifacts/deviations.md` D5.

Pull this **ahead of the rest of 0002** because two other RDRs cite it:
0007's `contains` contract-test leg (Phase 3) and 0004's read-back equality
both consume the encoding. Landing it late makes them cite something unwritten.

### 3. RDR 0007 Phases 2–4, RDR 0002 Phases 2–5

0007 Phase 2 (kernel domain rule, strong-Kleene combination, absent-key
payload) and 0002 Phase 2 (normalizer emitting to the new `Row`) can now both
proceed. 0007 Phase 3's `contains` leg needs step 2 landed.

### 4. RDR 0003 — guard predicate exhaustiveness

Consumes the atom shape (0007) and the normalized model (0002). Carries the
tag declaration model, so 0006's lint input contract reads it.

### 5. RDR 0004, 0005, 0006

0004 (accessor execution) and 0005 (CLI contract) consume the resolver surface;
0006 (graph lint) consumes 0003's declaration model. 0005 already re-entered on
§D8–§D11 and owes no citation.

### 6. RDR 0008, 0009 — smallest, last

Both are ownership questions over checks at `Resolve` entry.

---

## Open joint decisions — neither blocks the build

| JD | Siblings | Why it does not block |
| --- | --- | --- |
| **JD-5** precondition precedence | 0008, 0009 | Ordering of two checks at `Resolve` entry. The home says outright: *"Either order is defensible — pick one and pin it with a test… The silence is the defect, not the choice."* Sharpened form already names the tiebreak: 0009's fences admit only 0009-reports-first. Resolves when the second check is written, at step 6. |
| **JD-18** conforming-view enforcer | 0003, 0007 | Asks which *document states* the view-level check — prose ownership, not behavior. Runtime is already pinned by 0003 A18/A20 and 0006's owned half. Nobody is blocked from writing the check. |

Both are the kind of question the first Go type answers for free — which is
what the iteration-4 gate correctly said about JD-21/JD-22 before the home
answered them at the document layer instead.

## Standing rule for this run

**Findings route to `fix-now` or `kata-bug`, not to `rdr-seed`**, unless a
finding genuinely contradicts a locked contract. The cluster stands at ~80k
lines of RDR markdown against ~3.4k lines of non-spike Go, with no `internal/`
commit since 2026-08-09; the iteration-4 critique found this itself (C-31,
"Code has not moved") and responded with another iteration. Re-seeding RDRs
from implementation findings restarts that loop on the far side of the build.
