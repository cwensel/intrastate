# Run order — the `0002-0009` cluster

**Status**: authored 2026-08-24, after JDR 0001 settled (`5c2b96b`) and the
joint-decision qualifier sweep (`730319c`). Read this before `/rdr-implement`
on any cluster member.

Each RDR is implemented **to completion** in one `/rdr-implement NNNN` run.
This document gives the order of those eight runs. It is not a phase-
interleaving plan: no RDR is partially implemented and returned to later.

```
1.  /rdr-implement 7      ← must be first
2.  /rdr-implement 2      ← must be second
3.  /rdr-implement 3
4.  /rdr-implement 6
5.  /rdr-implement 4
6.  /rdr-implement 5
7.  /rdr-implement 8
8.  /rdr-implement 9
```

Steps 1 and 2 are **forced by compile-level proof**. Steps 3–8 are ordered by
citation direction — reasoned, not proven; a mistake there surfaces as a
compile error in minutes, not as rework.

---

## Why 0007 is first

RDR 0007 Phase 1 defines `internal/resolve::Row` — the §D1 atom slice that
replaces the shipped `Row.Guard string` — and exports the existence token and
literal constants. RDR 0002 Phase 2 normalizes **to** that type.

0002 states the dependency itself, five times:

| Site | What 0002 says |
| --- | --- |
| `0002:2263` | Phase 2's handoff-routing assertion "is unsatisfiable until RDR 0007's reshape lands (Prerequisites)" |
| `0002:2278` | The atom constants "compile only once RDR 0007's reshape exports them… the spike mirrors the constants locally because the kernel does not yet export them" |
| `0002:1651` | A Phase-2 scenario "waits on RDR 0007's reshape" |
| `0002:1662` | "a sequencing block, not a contradiction" |
| `0002:1772` | "The constants and the atom-shaped row do not exist in `internal/resolve` yet, so Phases 2–3 sequence behind that reshape" |

**This inverts the lock order.** The iteration-3 cluster gate wrote *"0002 is
the producer every other member's wire format and fixtures read, so its re-lock
is the practical gate for the set."* That is true for **lock order and
fixtures** and false for **Go types**. 0002 produces the wire format; 0007
produces the type the wire format normalizes into. Conflating the two axes
produced contradictory orderings in `cluster-reconcile/0002-0009/iter-3` and
`iter-4`; this file exists so the distinction does not have to be re-derived
from 2,300-line RDR bodies.

0007 Phase 1 also lands **§D12 (§JD-21)**: `Block` gains `BlockMatch` — one
type, no separate slice on `Row`. This widens a fenced cardinality ("exactly
two constants", `0007:1264-1266`), so it is a **normative change to a Final
RDR, not a citation**. See `0007/artifacts/deviations.md` D1: expect the scoped
re-entry, do not widen the type silently.

## Why 0002 is second — the cycle §D13 broke

This order was **impossible before 2026-08-24**. RDR 0007 Phase 3 reads:

> **Blocked on one RDR 0003 declaration.** `resolve.Tag.Value` is a bare
> `string`, so a set-valued tag reaches the seam as one opaque string with no
> stated element encoding. The `contains` leg of this contract test therefore
> cannot be written from any current document.

And 0003 needs 0007's atom shape (`0003:1919`: *"Phase 1 cannot begin against
the atom-slice shape until that reshape lands"*). That is a true cycle between
two whole RDRs — and under one-RDR-per-run, a cycle means **neither can run**.

JDR 0001 §D13 broke it by reassigning the encoding declaration away from 0003:

> `Tag.Value` stays `string`; a set crosses as its canonical JSON array —
> members sorted, duplicate-free, compact encoding — and **0002 declares it**.
> *Lands in 0002 — one normative clause; 0003, 0004, 0007 by citation; 0007's
> `contains` contract-test leg unblocks.*

With 0002 holding the declaration, the cycle becomes the chain
**0007 → 0002 → 0003**. This is the stronger of the two constraints on 0002:
not merely that it consumes 0007's types, but that 0003 cannot run until 0002
has declared the encoding 0003 cites.

### The one seam this order strains

0007 runs before 0002 has written the §D13 clause, so 0007's Phase 3
`contains` leg and Testing Strategy scenario 8 land citing an encoding that is
**decided but not yet in 0002's markdown**.

**Resolution: write them against §D13 directly.** The clause is unambiguous
about the byte form — sorted, duplicate-free, compact JSON array — and 0002's
later run must match it, not redefine it. Do not defer the two tests to a
second visit to 0007; that would be the partial implementation this order
exists to avoid.

## Steps 3–8

| # | RDR | Consumes | Note |
| --- | --- | --- | --- |
| 3 | **0003** guard predicate exhaustiveness | 0007's atom shape, 0002's normalized model | Owns the tag declaration model that 0006 reads. Instantiates `resolve.TestGuardEvaluatorContract` (0007 Phase 3) against its evaluator. |
| 4 | **0006** graph lint | 0003's declaration model | 0003's model **is** 0006's lint input contract — hence ahead of 0004/0005. |
| 5 | **0004** accessor execution | resolver surface; 0002's §D13 encoding | §D9 settles this RDR's internal split: `gate denied` is a **refusal** (`0004:804`), not the typed non-refusal at `0004:502`. Read-back equality is byte equality over the JSON array. |
| 6 | **0005** CLI contract | 0004's semantics; refusal codes | Already re-entered on §D8–§D11 and owes no citation. Maps refusals and lint findings to the envelope, so it wants 0004 settled first. |
| 7 | **0008** recognized tag key | `Resolve` entry | Reserved-key check. |
| 8 | **0009** escape-row shape | `Resolve` entry | Breach check. |

## Open joint decisions — neither blocks any run

| JD | Siblings | Why it does not block |
| --- | --- | --- |
| **JD-5** precondition precedence | 0008, 0009 | Ordering of two checks at `Resolve` entry. The home says outright: *"Either order is defensible — pick one and pin it with a test… The silence is the defect, not the choice."* The sharpened form names the tiebreak: 0009's fences admit only 0009-reports-first. **Resolves during step 8**, when the second of the two checks is written. |
| **JD-18** conforming-view enforcer | 0003, 0007 | Asks which *document states* the view-level check — prose ownership, not behavior. Runtime is already pinned by 0003 A18/A20 and 0006's owned half. Nobody is blocked from writing the check. |

## Standing rule for this run

**Findings route to `fix-now` or `kata-bug`, not to `rdr-seed`**, unless a
finding genuinely contradicts a locked contract. The cluster stands at ~80k
lines of RDR markdown against ~3.4k lines of non-spike Go, with no `internal/`
commit since 2026-08-09; the iteration-4 critique found this itself (C-31,
"Code has not moved") and responded with another iteration. Re-seeding RDRs
from implementation findings restarts that loop on the far side of the build.
