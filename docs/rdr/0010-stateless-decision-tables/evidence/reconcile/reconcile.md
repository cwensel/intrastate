Model: claude-opus-5[1m]

# Stage 6 — Reconcile, RDR 0010 (stateless decision tables)

**Verdict: RECONCILED.** All fifteen Critical Assumptions are terminal
(`Verified`); no BLOCKER; no item deferred past lock. Ready for Finalize.

## Stage 5 completeness preflight

`Profile` = **foundational**, so the required lens row is
`cove → 3amigo → critique → repeatability` (rdr-common §lens-row). All four
have completed evidence with distinct `Model:` stamps:

| Lens | Evidence | Stamp(s) |
| --- | --- | --- |
| cove | `findings.md`, `step0-grounding.md` | `claude-opus-5[1m]` |
| 3amigo | 3 persona files + `consolidation.md` + `resolve.md` | `claude-opus-5[1m]` |
| critique | `critique.md`, `critique-modelB.md`, `diff.md`, `resolve.md` | `claude-opus-5` / `claude-sonnet-5` (dual-model, converged) |
| repeatability | `run-1/2/3.md`, `diff.md`, `resolve.md` | `claude-sonnet-5` / `claude-fable-5` / `claude-haiku-4-5` |

Repeatability ran the **full** variant (`variant: full (profile: foundational)`),
correct for this profile — not a lite/full mismatch. The Determinacy add-on does
not apply (it is a `mid`/`large` trigger). **Preflight PASS.**

## The open set — all four sources

| Source | Items |
| --- | --- |
| 1. Pre-Lock needs-(re)verification lists | A6, A12, A13, A14, A15 |
| 2. Still-Pending / Unverified assumptions | the same five (A1–A5, A7–A11 were already Verified) |
| 3. Named-but-unrun spikes | **none** — absorption audit found no residue |
| 4. Exactness-word delta (post-mutation) | **none ungrounded** |

Sources 3 and 4 were built by a delegated absorption audit over the four lens
output dirs. It checked all **49 findings** across the rounds against the
current RDR through the projector: every finding marked `fixed` was verified
present, five dismissals were intentional and cited, and every open obligation
was already booked as one of the five Pending assumptions. **Residue list
empty; ungrounded-exactness list empty.** The one claim that would have been
ungrounded — A6's "a flat string carries the rendered command whole" — is
exactly the one the critique lens caught and flipped to Pending, and it is
re-verified below.

## Dispositions

| Item | Source | Disposition | Evidence pointer |
| --- | --- | --- | --- |
| **A6** — flat string-valued `[rule.emit]` sufficient for the motivating consumer | 1 (critique D-7) | **VERIFIED** | `evidence/spikes/a6-emit-keys.md`; `rdr-status/SKILL.md:168`, `rdr-facts.toml:79,99,144,257,552` |
| **A12** — `graph-dangling-edge` splits by arm without a taxonomy change | 1 (cove) | **VERIFIED** | `evidence/spikes/a12-consumers.md`; 12 consumers enumerated; `go test ./internal/graphlint/...` ok |
| **A13** — zero-dimension group takes `graph-unprovable-coverage` from the `len(dims)==0` branch | 1 (3amigo) | **VERIFIED** | `evidence/spikes/a13-zerodim.md` + `a13-zerodim/` spike pkg + `a13-spike.out` |
| **A14** — `no-participating-dimension` appended to 0006's closed `reason` set | 1 (critique D-1) | **VERIFIED** | `evidence/spikes/a14-reason-append.md`; full suite green on scratch copy |
| **A15** — exported `Model.Class` sufficient for the four class readers | 1 (repeatability) | **VERIFIED** | `evidence/spikes/a15-class-readers.md` |

No DOWNGRADED item, no ACCEPTED-as-design-decision item, no BLOCKER.

## What each verification actually established

**A6 — re-verified against the consumer's *actual* emit keys.** The critique
lens was right that the old supporting reading was false. Checked against
rdr#tmxk itself: the consumer's answer is one routing decision, rendered at
`rdr-status/SKILL.md:168` as `Next: /rdr-prelock 0046 critique`. It decomposes
into stage verb + lens (both closed sets of authoring-time literals) and the
record number, which the **caller** already holds — the user typed it. The
match side agrees: every fact the routing table keys on is a closed enum of
fixed strings. The narrowed claim holds; the widening A6 defers (a template
language in an emit value) is genuinely not needed.

**A12 — twelve consumers, none broken, for a structural reason.** `table.Model`
has no `Class` field today and no `class` layout key exists, so every
checked-in fixture is state-machine class by construction and C5's predicate is
*bit-identical* to today's behaviour over all of them. Root-arm consumers all
assert on **code only** over class-less models; the terminal arm's sole
consumer (`TestReq129`) uses a **rooted** model where the root arm was already
silent. The one rootless fixture, `pos-no-initial.toml`, is loader-only and
carries no lint expectation.
*Implementation trap surfaced:* the root-arm predicate must be the **augmenting
OR** (`class == "decision-table" || len(Initial) > 0`), never class alone —
class alone would silence the root arm for rootless *state-machine* models and
flip `TestReq32`. C5 already mandates exactly this, so this is a confirmation
of the contract, recorded in A12 for the implementer.

**A13 — executed against the real loader and engine**, from a spike package
under the evidence dir, `internal/` untouched.
- *Q1:* a naive insertion **would double-report**. Over a zero-dimension group
  the product is the empty product and IS projectable, `coverageUnionFor` takes
  its own zero-dimension arm where membership alone decides, so every such
  group closes every arm and `bareEscapeFor` emits
  `graph-coverage-closed-by-escape`. Emitting and falling through yields both
  findings — the pair C5 forbids. The **early `return`** delivers C5's
  precedence, and forfeits nothing: `TestQ1b` establishes the closure advisory
  is the *only* invariant-4 output reachable on that path.
- *Q2:* **151 zero-dimension groups across 37 of 37 loadable fixtures**,
  including production `models/rdr.toml`; none declares a class. This confirms
  "legitimate and stays silent" empirically and makes the class-keying a hard
  requirement with a measured cost — an unconditional arm would fire a
  *blocking* code on 151 groups and redden every fixture at once.

**A14 — the potential BLOCKER, and it does not block.** All three halves pass.
- *Mechanical:* exactly **one** append-sensitive consumer in the tree
  (`TestReq80`'s sorted literal). Every production site assigns rather than
  enumerates, or switches **with a `default`**. ⚠ `internal/resolve/guard.go:81
  Reasons()` is a same-named but entirely **unrelated** set — a false positive
  for anyone grepping.
- *Test:* the update is a two-line source append plus one sorted-position
  literal edit; `go build ./...` and `go test ./...` fully green across all
  nine packages on a scratch copy (committed `internal/` never modified).
- *Owner assent — structural, no reopen.* The decisive reading is 0006's own
  contrasting vocabulary: where it means a set that may **not** grow it writes
  "closed **at** …" and names the members (`0006:C17`); where it means a set
  that may grow it writes **append-only** (`0006:C14`,
  `0006:D-wire-byte-format`). REQ-80 declares this set "closed, append-only",
  so `closed` = no value outside the set is legal and `append-only` = it may be
  extended. **The append is the sanctioned evolution, pre-authorized by 0006.**
  No 0006 sentence becomes false. Peer precedent: 0008 (added a validation
  category on an "including at minimum" list) and 0011 (added a reason token
  while explicitly declining to widen the adjacent closed kernel set) — the
  project already distinguishes an authorized append from a prohibited one.
- *Constraint fixed into C5:* the new arm MUST supply its **own message**
  rather than routing through `coverage.go::unprovableMessage`, whose `default`
  would phrase the remedy as "declare the domain" — wrong for a dimension that
  does not exist.

**A15 — the projection risk does not bite.** All four readers reach the model
where they run (`normalizeRule` via `*loader.model`, `reach` via its own
`*table.Model` param, both graph-lint checks via `analysis.model`). The
load-bearing question was whether graphlint sees only the `KernelTable`
projection, which would have made an exported field insufficient: it does not —
`graphlint.Request` carries the **full `*table.Model`**. Ordering holds
(`loadModelHeader` is `run`'s first step; graph lint runs after load
completes), and the C1 window has **six** valid insertion points.

## MVV-critical check (hard rule)

The MVV depends on C5's fence behaving as specified — step 3 asserts
`findings exactly []` with no taxonomy code present, and step 2″ exercises the
match-discriminated fence. A13 and A14 are therefore **pre-lock
prerequisites, not deferrable**, which is how they were treated: both were run
now rather than downgraded. Nothing MVV-critical is deferred past lock.

## Refutation check (hard rule)

**No spike or source-search refuted an assumption the RDR relies on.** The one
refutation in this RDR's history — the critique lens against A6's supporting
reading — was caught at Stage 5, flipped A6 to Pending, and is discharged here
by re-verification against the consumer's real emit keys. No route-back is
owed, so no `§punt-ledger` row and no `§strong-consult` were triggered at this
stage.

## RDR edits made (dispositions written into the record)

Stage 7 reads the RDR, not this report, so every disposition landed in the
record:

- **A6, A12, A13, A14, A15** — `Status: Pending` → `Verified`, each with its
  concrete evidence pointer and the substance of what was established.
- **C5** — two clause amendments: the precedence paragraph now states that the
  early **`return`** is what delivers it (falling through double-reports,
  executed), and the "conditional on A13 and A14" paragraph now records both as
  Verified plus the own-message constraint.
- **§prerequisites** — both open boxes checked, with the structural-assent
  reasoning for the 0006 append.
- **§risks-and-mitigations** — "Carried, not yet demonstrated" → demonstrated,
  carrying the 151-group measurement.
- **§trace** — step 2″ re-witnessed; the closing paragraph now separates the
  witnessed 2″ from 4′, which is an implementation obligation by design.

Per rdr-common §amendment-sweep, the amended C5 clauses were swept: all
`zero-dimension` / `zero participating` sites were re-read for agreement, and
the two that had gone stale (the Trace row and the Risks mitigation) were
updated in the same pass.

## Post-edit checks

- `rdr lint 0010` → **PASS** (the one `conformance gate:inline` advisory is
  Stage 7's cross-file move of the gate responses to `artifacts/gate.md`, not a
  reconcile item).
- No `_Draft placeholder._`, no seed-skeleton header, no surviving template
  bracket.
- `## References` is authored, not template text.
- All 15 assumptions `Verified`; zero `Pending`/`Unverified`.
