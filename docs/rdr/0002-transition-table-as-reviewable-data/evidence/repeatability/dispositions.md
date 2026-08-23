Model: claude-opus-5[1m]

# Repeatability Dispositions — RDR 0002

Origin ledger = `diff.md`'s R-1..R-9 findings and C-1..C-7 GUESS clusters.
Diff health: **healthy** (localized to four contract seams, GUESS markers
clustered, all hard normative clauses reproduced identically by all three runs)
— so the disposition is edit, not rerun on another model.

Runs: run-1 `claude-opus-5[1m]`, run-2 `Claude Sonnet 5`, run-3 `claude-fable-5`
— all `variant: full (profile: foundational)`. Cross-model draw satisfied.

## Per finding

| ID | Grounding verdict | Disposition | Section touched |
|---|---|---|---|
| R-1 | Genuine silence — categories decidable only during normalization are labeled "load-time" | **pin** — "load" = the whole source-to-candidate-rows pipeline; internal decomposition free, category escape forbidden | Normative Contracts (category block) |
| R-2 | Genuine silence — L679 contemplates a multi-model dump, L682 says 0006 lints one model per invocation; no clause names the surface | **pin** — `[model]` singular, so the category is scoped to a multi-document invocation; fixture is a *pair* | Normative Contracts (dump ordering) + Testing Strategy |
| R-3 | Genuine silence — only the version gate is ordered; fail-fast vs accumulate unstated | **pin** — fail-fast; order among independent defects explicitly unspecified | Normative Contracts (version gate block) |
| R-4 | Genuine silence — audit row defers path *binding*, never fixes the entry's I/O boundary | **pin** — entry takes bytes + source id; no file I/O in this package | Existing Infrastructure Audit |
| R-5 | Contradiction — L358/L649 carry atoms as one field; L625 routes them into the kernel's two | **single-source** — split applied at the kernel handoff; normalized value keeps the unified block-retaining set (`Guard string` cannot carry blocks) | Normative Contracts (split clause) + Technical Design L358 |
| R-6 | Contradiction — Load-Bearing Decisions uses "combined predicate set", which the RDR *defines* (L453, L471) as the wider extent, for a lifting rule the contracts fix to the narrow one | **single-source** — bullet restated to the match-block extent, with the wider term explicitly excluded | Load-Bearing Decisions |
| R-7 | Genuine silence — L416 requires a write block on ordinary rules; no category covers its absence | **pin** — added `malformed rule shape` to the category list + a scenario 3 variant | Normative Contracts (category list) + Testing Strategy |
| R-8 | Re-raise on the contract (L620 states the MUST plainly) — run-2 violated a stated clause | **dismiss-with-cite** on the contract; **pin** the real half: the clause is unassertable by value comparison (equal sets; `[]Tag` on `main` also shares a backing array) | Testing Strategy (scenario 2) |
| R-9 | Partly decided — prohibition binds normalization and the kernel row, but the dump contract *requires* a derived kind column | **pin** — a render-time view MAY materialize it; it MUST NOT feed back or be branched on | Normative Contracts (row-kind clause) |
| C-1 | Naming only; no downstream RDR cites an identifier from this package | **leave non-normative** — the structural half is R-1, fixed there | — |
| C-2 | Independence already decided ("before CLI mapping" + A5); no import ban existed | **pin** the boundary — MUST NOT import `internal/cli`; `clierr` is a leaf package on `main` for exactly this reason | Normative Contracts (category block) |
| C-3 | Genuine silence — one member (`reserved_tag_key`, owned by RDR 0008) is identifier-spelled, the rest prose | **pin** — snake_case renderings of the prose names, anchored on `reserved_tag_key` | Normative Contracts (category block) |
| C-4 | **False positive** — all three runs marked GUESS *because* the RDR explicitly declines a dump grammar (L812–814) | **dismiss-with-cite** | — |
| C-5 | Same gap as R-2 | folded into R-2 | — |
| C-6 | Genuine silence in the RDR body — but the spike fixture *does* author `mode` + `path`, and the fixtures are normative-by-promotion | **pin** small — entry carries `mode`/`path`; this RDR validates only id resolution, RDR 0004 owns semantics | Normative Contracts (source schema) |
| C-7 | Same seam as R-1 | folded into R-1 | — |

No finding was charted-to-successor: none was net-new scope. No tiebreaker was
escalated — R-5 was the only genuine fork and the evidence collapsed it (block
retention is required downstream and `Guard string` cannot carry blocks, so the
split must be at the handoff).

## Needs (re)verification — carried to Stage 6

- **A9** (new, `Pending`) — the normalized row can carry one unified
  block-retaining atom set and still produce a conforming `resolve.Row` at the
  handoff. Method: MVV Test. Blocked on RDR 0007's kernel reshape.
- **A10** (new, `Pending`) — `duplicate model id` is decidable only across
  documents. Method: MVV Test (paired-document fixture).
- Prerequisites checkbox "All Critical Assumptions verified" flipped back to
  unchecked; A9/A10 are owed.

## Amendment sweep

- `combined predicate set` — 2 sites; L453 is the wider-extent *definition*
  (correct as-is), L566–568's two-extent disambiguation still stands, the
  Load-Bearing bullet is the amended site.
- `duplicate model id` — normative clause and the Testing Strategy fixture list
  both updated; Failure Modes prose unaffected (names the category, not a scope).
- Predicate split — normative clause and the L358 "is the kernel row" sentence
  aligned in the same pass; dump field list and Round-Trip field list already
  agreed with the chosen reading and needed no edit.
- Category list — `malformed rule shape` added normatively and given a scenario 3
  variant in the same pass.

## Mini-checks

Not the RDR's first lens pass (cove, 3amigo, critique all preceded it), and no
fix here fired a new cue — the `fidelity` / `trace` tables already in the draft
are inherited unchanged.
