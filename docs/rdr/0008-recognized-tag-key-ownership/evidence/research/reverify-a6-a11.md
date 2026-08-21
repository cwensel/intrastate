# RDR 0008 — Stage 4 scoped re-entry: A6 / A11 re-verification

Model: claude-opus-5[1m]
Date: 2026-08-21
Scope: `Draft [revised from Final 2026-08-11; re-verify A6, A11]` — the
`re-verify` set is exactly {A6, A11} plus the anchors the demotion touched.
Every other assumption's lock-time `Verified` stamp carries forward undisturbed
and was NOT re-derived. Stage-2 prior-art cache (`prior-art.md`) reused, not
re-run — no assumption in scope depends on it.

## The defect, confirmed

Both assumptions rested on a negative existential: "the string `Input` appears
**zero** times in RDR 0009." Re-counted against Final 0009 (case-sensitive,
whole file): **5 occurrences** at raw lines 424, 509, 736, 741, 742, 857 — of
which 424 is the unrelated token `legalInput`, leaving 5 true references.
Line 857 sits inside 0009's own Normative Contracts.

Causal chain: the count was true of the **Draft** 0009 this RDR read. 0009 then
gained **A10** (raised by its repeatability lens, diff finding D-1, which had
reconstructed a two-parameter `Resolve`), pinning
`func Resolve(in Input) (Result, error)` with the table reached as `in.Table`.
That is what introduced the `Input` references, after which 0009 locked to Final.

## What falsifies A6 as written

0009's contract is **both** a producer obligation **and** a `Resolve`-entry
precondition — not one or the other:

- producer half — "a PRODUCER obligation on every constructor of resolve.Row
  values: a Row with a non-empty Escape list MUST have an empty Writes slice"
- kernel half — "The kernel enforces the same obligation as an entry
  precondition of Resolve. Resolve's signature is UNCHANGED —
  `func Resolve(in Input) (Result, error)`, one parameter — and the checked
  table is the one already reached through the input, `in.Table`"
- two-site duality — "the conformance predicate MUST be exported by the kernel
  package as a construction-time check callable by any table producer, **and
  Resolve's entry precondition MUST be that same function** — one predicate,
  two call sites"
- scope — "evaluated over every row of the supplied table at `Resolve` entry,
  not only rows the resolution touches"

So the boundary A6 concluded this RDR would have to establish alone already
carries a peer precondition. The conclusion "needs an enforcement locus of its
own" does not survive in that form.

## What survives, and why the design does not change

Ownership does not transfer with co-residency. The separation is **structural**:
0009's Normative Contracts require the predicate be "a METHOD ON Table taking no
arguments," so its receiver is `Table` and it cannot read `Input.Owned` /
`Input.Observed` at all — the very fields this RDR's precondition must read.
0009 reaches the table *through* `Input` but checks a property of the table
value alone ("a malformed table is malformed as a value, independent of the
input tuple").

Read-domains, re-derived from each RDR's own contract:

| | this RDR (0008) | RDR 0009 |
| --- | --- | --- |
| reads | `Input.Owned`, `Input.Observed`, `Row.RequiresOwned` | `Row.Escape`, `Row.Writes` |
| subject | a reserved *name* | a row's *shape* |
| receiver | the input tuple | `Table` |

Disjoint, both pure reads. A shape check on `Table` cannot detect a
reserved-key breach over the input's tag snapshots, so riding 0009's contract
was never available: Enforcement-locus candidate (b)+(c) is unchanged. What
changes is the **rationale** — the locus is shared infrastructure this RDR
joins, not virgin boundary it establishes. Strictly less novel surface than the
pre-demotion text claimed, which is the direction A6's own "If wrong" line
anticipated.

Independently re-derived at cluster-reconcile rather than taken from either RDR:
`docs/rdr/cluster-reconcile/0002-0009/pairwise-0008-0009.md` → FINDING 2.

## Source confirmation (RDR 0001 is `Implemented`)

Anchors are durable (`path::Symbol`), per the demotion's instruction that bare
line numbers are what went stale.

- `internal/resolve/resolve.go::Resolve` — declared
  `func Resolve(in Input) (Result, error)`, one parameter. Confirms 0009's A10
  against real code, not just peer text.
- `internal/resolve/resolve.go::Input` — carries `Table Table` alongside `Flow`,
  `Owned`, `Observed`, `Recognized`, `Guards`.
- `internal/resolve/resolve.go::Row` — carries `RequiresOwned []string` (0008's
  field) alongside `Escape []RefusalKind` and `Writes []Tag` (0009's), so both
  predicates are co-resident over one row type without overlapping on any field.
- `internal/resolve/resolve.go::Resolve` doc comment — "the error return is
  reserved for programmer mistakes, not for modeled refusals" (RDR 0001 REQ-6),
  verbatim as A6 quoted it.
- **No non-nil error path exists at HEAD**: every `return` in `Resolve`'s body
  pairs its `Result` with a literal `nil` (5 sites). The reserved channel is
  unexercised surface — which both RDRs' preconditions intend to land on, and is
  precisely why JD-5 exists: both claim the same single `error` slot.
- `internal/resolve/resolve.go::recognizedTagKey` — `const recognizedTagKey =
  "recognized"`, unexported. No `Table.CheckValid` method and no reserved-key
  precondition exist yet; **neither RDR's surface is built**, so nothing in
  source contradicts either RDR's pins.

## Where the precedence question went

Re-entry outcome **(2)**: the co-residency is load-bearing, so the precedence
question belongs to the umbrella, not to this RDR. It has a named home:

> **JD-5 Precondition precedence.** 0009's breach check and 0008's reserved-key
> check both land at `Resolve` entry; neither orders itself against the other.
> Either order is defensible — pick one and pin it with a test on a table that
> breaches both. The silence is the defect, not the choice.
> — `docs/jdr/0001-resolve-kernel-seam.md` → Interface record → JD-5

**JD-5 is OPEN.** Its neighbours carry explicit closures (JD-6 "Closed by §D1",
JD-7 "Closed by §D3", JD-11 "Withdrawn"); JD-5 carries none. JDR 0001's three
resolved decisions (§D1 guard seam, §D2 gate-then-count, §D3 complete reads) do
not reach it.

RDR 0009 carries the matching latch on its live Status:
`Final [joint decision → JDR 0001 §JD-5]`.

JDR 0001 also names this exact re-entry as outside its own scope — "RDR 0008's
re-entry — A6/A11 verified on a false claim about Final 0009; already demoted;
re-enters at resolve, scoped to those two assumptions" — so the two records
agree on the division of labour and neither absorbs the other.

Block 4's interim rule ("either error is conforming; what is *not* conforming is
skipping this check because another fired") is compatible with JD-5's "either
order is defensible" and stands until JD-5 closes. Stale "Stage 7.1" routing in
A11 and block 4 re-pointed to JD-5.

## Citations repaired

| Was | Now |
| --- | --- |
| "`Input` appears zero times in RDR 0009" (A6) | withdrawn; replaced by the quoted producer/precondition clauses |
| "`Input` occurs zero times … re-counted at Reconcile" (A11) | withdrawn; disjointness rests on the field sets |
| `0009…md:536-540` (A11) | RDR 0009 → Technical Design → Normative Contracts (predicate clause) |
| `0009…md:510-517` (A11, and body block 4 rationale) | RDR 0009 → Technical Design → Normative Contracts (two-call-sites clause) |
| "reconciled at Stage 7.1" (A11, block 4) | JDR 0001 §JD-5 |

0009's Normative Contracts blocks carry **no per-contract IDs** — they are eight
positional untagged ```normative``` fences. The durable anchor is therefore the
section heading plus the clause subject, which is what was used.

## Out of scope, recorded not absorbed

**JD-10** — that this RDR's name constraint invalidates all three of RDR 0002's
canonical fixtures, contradicting this RDR's `Overrides` claim to "narrow
nothing in either peer" — is settled at the umbrella per the demotion's "Also
note". JDR 0001 §JD-10 records the disposition ("rename them; there are no users
to migrate") and keeps one sub-question open (whether a declared `recognized`
tag is total or partial). Not widened into this re-entry.

Predecessor RDR 0002 is itself now `Draft [revised from Final 2026-08-12;
re-verify A2, A7]`. This does not disturb the re-verify set: no assumption in
scope depends on 0002's re-verified pair.

## Author's round (Stage 4's single interaction)

Fixtures rendered: **none**. No assumption in the re-verify set turns on an
exact value — A6/A11 are structural claims about which contract owns a boundary
and which fields each predicate reads, not about a wire record, byte layout, or
count. Rendering a fixture here would be invented, not read.

Two questions put, both carrying their grounding:

1. **JD-5 posture** — grounded on JDR 0001 §JD-5 (quoted open, no closure
   marker), RDR 0009's Status latch, and JDR 0001's own "What this does not
   decide" entry scoping this re-entry to A6/A11.
   **Approved: cite JD-5, stay Draft.** A6/A11 verify on the surviving
   substance; the ordering is cited and never restated. The RDR does not re-lock
   to Final until JD-5 closes, because block 4's interim "either error is
   conforming" rule depends on it. No Status latch added — it would compose
   awkwardly with the `re-verify` qualifier Stage 8 expects to self-clear.
   Rejected: proposing a precedence here (exceeds re-entry scope; JDR 0001
   reserves it to the umbrella).

2. **Doubly-breaching test variant** — grounded on JD-5's own instruction ("pin
   it with a test on a table that breaches both") against Testing Strategy
   scenario 6, which had no such case.
   **Approved: keep as authored.** Added as scenario 6's fourth variant,
   asserting only the non-silent-skip property that block 4 licenses, explicitly
   NOT asserting which error is reported, and marked as the test that pins JD-5
   once it closes.

Round tally: 0 fixtures put / 2 questions put / 2 approved / 0 rejected.

## Profile

Unchanged: **`foundational`**, one contract (who owns the *name* of the
recognized-provenance tag key in the assembled view).

Recount, not carried from the Seed estimate: the Normative Contracts section
holds six ```normative``` blocks, but they are clauses of a **single**
independent contract — the reserved-key identity rule and its enforcement — not
six contracts. No second independent contract (no distinct hash, wire format,
taxonomy, or destructive-op policy) is present, so the ≥2 split signal does not
fire.

The `foundational` value comes from the cross-RDR-producer trigger, and the
re-verification **strengthens** it rather than weakening it: 0008 binds RDR
0001's kernel carrier to RDR 0002's declared model, extends 0002's validation
surface, and now additionally co-resides with RDR 0009 at `Resolve` entry with a
precedence question deferred to JDR 0001 §JD-5. That is more cross-RDR coupling
than at lock, never less.

Accretion floor: **not applicable** — `Seam Lineage` reads "`area:internal-resolve`
(recognized-tag key identity at the RDR 0001 ↔ 0002 boundary) — no prior
accretion," i.e. zero closed prior point-fixes, below the ≥2 threshold. The
count and the floor agree, so no disposition was needed. The count could only
have raised the profile here; `foundational` is already the ceiling.
