Model: claude-opus-5[1m]
Variant: single-model fallback (two fresh contexts, diffed) — no alt model
reachable; `resources.md` records "Alt-model roster: omitted (single model)".
Recorded as a fallback per `pre-lock/2-critique.md`, not passed off as a full
dual-model pass. RDR Profile is `foundational`, for which dual-model is
*required to converge* or the recorded single-model fallback stands in.

# Critique lens — pass A ↔ pass B diff (RDR 0008, iteration 1)

Ledgers reconciled **by passage anchor**, not by ID (`C-N` IDs are per-file and
do not correspond across passes). Pass A = `critique.md` (15 rows), pass B =
`critique-passB.md` (16 rows).

## Convergent — both passes hit the same passage independently

Agreement across two independent contexts is the strongest signal the fallback
produces; these are the rows the resolve pass treats as highest-confidence.

| # | Passage | A | B | Grounded verdict |
| --- | --- | --- | --- | --- |
| D-1 | A4 Evidence / block 2 first clause | C-1 | C-3, C-4 | **CONFIRMED** — 3/3 committed fixtures declare a `recognized`-provenance tag under a non-`recognized` name: `0002/evidence/spikes/rdr-fixture.toml:25` + `kata-fixture.toml:22` (`[tags.outcome]`), `0003/evidence/spikes/guard-fixture.toml:36` (`[tags.rewind_target]`). A4 quotes these three files as proof nothing breaks. |
| D-2 | A4 "If wrong" — "mechanical rename" | C-2 | C-4 | **CONFIRMED** — five `[rule.match.outcome]` sites (`rdr-fixture.toml:59,74,86`; `kata-fixture.toml:47,57`). `0002…md:272` makes `[rule.match.<tag>]` a *tag* predicate, so references follow the declaration. |
| D-3 | Block 2 reference-position fall-through | C-3 (match position) | C-5 (write position) | **CONFIRMED, two faces of one gap** — the disjunction covers declared-and-referenced + undeclared; it misses declared-but-referencing-a-different-declared-tag (A) and `[rule.write] recognized` on a legally-declared key (B). |
| D-4 | A8 blast-radius sweep method | C-5 | C-9 | **CONFIRMED** — `resolve.go:1-11` package doc: `Input.Owned` is "an accessor-produced owned tag snapshot". A literal grep over `*_test.go` cannot clear a data-derived channel. |
| D-5 | Block 4 "unconditional on `Input.Recognized`" | C-7 | C-8 | **CONFIRMED** — `resolve.go:151` gates injection on non-empty, so the empty-outcome case has nothing to shadow; scenario 6 pins the error as expected. |
| D-6 | Testing Strategy Done split / MVV | C-9 | C-12 | **CONFIRMED** — every scenario delivering the Problem Statement outcome (2, 4, 5, 7, MVV lint half) is carried into 0002's implementation. |
| D-7 | "told why" asserted at data level only | C-10 | C-13 | **CONFIRMED** — already charted as `3amigo/charted.md` C-2; both passes independently re-derive that no verb renders the category. |
| D-8 | Block 5 unenforceable MUST | C-14 | C-14 | **CONFIRMED** — no authored source form, no lint half, scenario 9 explicitly has "no lint half", carried on `Pending` A7. |
| D-9 | Identity rule blesses invisible near-miss | C-13 (case) | C-10 (whitespace) | **CONFIRMED** — both pinned as correct by scenario 4; the whitespace form is the more dangerous half (invisible in diff). |

## Divergent — raised by one pass only

Disagreement is the other half of the signal. Each is grounded on its own
merits rather than discounted for being single-source.

| # | Passage | Raised by | Grounded verdict |
| --- | --- | --- | --- |
| D-10 | Block 4 + Enforcement locus vs **RDR 0009's identical `Resolve`-entry clause** | B C-1 | **CONFIRMED — the most consequential divergence.** `0009…md:510-517` writes the *same* normative sentence over the *same* entry point: "The conformance predicate MUST be exported by the kernel package … and Resolve's entry precondition MUST be that same function — **one predicate, two call sites**." 0009 further pins "the breach error precedes every modeled disposition, `unmodeled_outcome` included" (`:531-538`). RDR 0008 borrows this shape (A6, Enforcement locus) while *forbidding* any symbol-name binding to 0009 — and neither RDR states ordering or composition for two entry preconditions. Pass A missed this entirely. |
| D-11 | Block 1 parenthetical "unreachable past the outcome-alphabet gate" | B C-7 | **CONFIRMED as stated, low severity.** `Table.models` is `slices.Contains(t.Outcomes, outcome)` (`resolve.go:216-218`) — nothing forbids `""` in the alphabet, and `assemble` runs at `:319` *before* the gate at `:321`. The parenthetical already hedges "for any alphabet that excludes the empty string", so this sharpens an acknowledged scope rather than refuting it. |
| D-12 | Block 3 carries **two** normative discriminators | B C-6 | **CONFIRMED** — block 3 makes both `reserved_tag_key` (category token) and `reserved-tag-key/kernel-owned` (rule id) normative and golden-tested, without stating which a consumer keys on. |
| D-13 | A1 verified the category list, not the `[tags.<tag>]` **naming grammar** | B C-16 | **CONFIRMED** — `0002…md:266-274` is a *separate* normative block from the validation-category block at `:342-348`. A1's "including at minimum" extensibility argument covers the latter only; the name constraint lands on the former. |
| D-14 | A5 backstop rationale stale under the settled locus | B C-11 | **PARTLY CONFIRMED.** Real tension: with the predicate at `Resolve` entry, "input that bypasses the precondition" shrinks. But not empty — scenario 8 explicitly resolves "through the package-internal test path", and A5's own Scope paragraph names hand-constructed input. The rationale needs re-deriving, not deleting. |
| D-15 | 0007(Final) A13 rests on this Draft's sentence | B C-2 | **RE-RAISE** — identical to cove F-1, dismissed-with-cite last pass and already recorded as a cluster-reconcile watch item on the joint-check line. B independently rediscovered it; the standing dismissal governs. |
| D-16 | Contract-count closure defended by assertion | A C-8 | **CONFIRMED as a method critique.** The enumeration was extended twice under review (blocks 1–3 → +4 → +5); D-3 and D-10 each name a further channel. |
| D-17 | "no kernel code changes" vs QOC blast-radius row | A C-11 | **CONFIRMED** — Consequences and Enforcement locus concede the surface; the *deciding* QOC row and the Trade-offs "cheapest branch" bullet were never rewritten against the settled locus. |
| D-18 | First non-nil error path breaks caller contract | A C-6 | **PARTLY REFUTED on symptom.** There are **zero** non-test callers of `Resolve` at HEAD (`rg` over `*.go` minus tests returns only the definition at `:318`). All test call sites route through `mustResolve` (`resolve_test.go:14-21`), which `t.Fatalf`s on non-nil — a tripped predicate is a red test, not a silent nil-deref. The contract observation stands for future CLI wiring; pass A's premortem "panic in a CLI verb" does not hold today. `resolve_test.go:748` pins `Resolve(resolve.Input{})` → nil error. |
| D-19 | Missing **lower bound** on the recognized declaration | A C-4 | **CONFIRMED** — block 2 constrains the name, never requires the declaration to exist. Cove Q10 raised the lower bound; the fix pass added only the upper-bound cardinality paragraph. |
| D-20 | A5 derivation confuses declaration space with data space | A C-12 | **CONFIRMED** — P1/P2 are RDR 0002 *table-declaration* rules; `Input.Owned`/`Input.Observed` are runtime data they do not govern. Same root as D-4. |
| D-21 | A7 `Pending` with five settled-fact dependents | A C-15 | **CONFIRMED** — Approach 5, block 5, Capability Dependencies, Phase 2, scenario 9 all state A7's conclusion as settled. This is the Finalization Gate's own Status-consistency rule. |
| D-22 | Problem Statement re-specifies the user's want | B C-15 | **PARTLY RE-RAISE** — 3amigo PM-1 raised the silent-reframing half and it was fixed (the stance is now explicit). B's *new* half is behavioral evidence: 3/3 fixture authors exercised naming freedom, and the XState analogy compares a struct slot to a bare word in the author's own flat namespace. That half is new. |

## Health verdict (Stage 5 review gate)

**Healthy, both passes.** Concrete named passages, real anchors, specific user
journeys; ledger rows span all six origins (§1/§2/§3/premortem/AT). No generic
advice, no all-`§1` ledger, no over-agreement between the passes — B contradicts
A's framing on the enforcement locus (A says rewrite toward (a); B says the
locus collides with 0009) and found four passages A missed. Independence held:
B did not read A.

**Convergence on the fixture defect (D-1) from two independent contexts, plus
the peer-collision finding (D-10), is the load-bearing output of this lens.**

## Tally

- 9 convergent passage groups, 13 divergent rows.
- 1 full re-raise of a standing dismissal (D-15, cove F-1), 1 partial (D-22).
- 1 partial refutation of a pass-A symptom (D-18) — the grounding gate caught it.
- 0 findings dismissed for failing the code-on-`main` ground; every cited
  symbol and fixture line resolves.
