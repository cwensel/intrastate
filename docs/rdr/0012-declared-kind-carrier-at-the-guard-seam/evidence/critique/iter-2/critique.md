Model: claude-sonnet-5

---

# Delta-scoped critique re-run — RDR 0012, iteration 2

Scope: verify whether the 14 claimed pre-lock rewrites held, and whether any
opened a new gap. Verified against `git diff HEAD -- docs/rdr/0012-*.md`
(confirms the rewrite set matches the 14 claims exactly, no unlisted changes)
and against source (`internal/graphlint/reach.go`, `internal/graphlint/analysis.go`,
`internal/table/load.go`, `internal/accessor/model.go`, `internal/resolve/resolve.go`,
`internal/cli/flow_input.go`, `internal/cli/flow_next.go`, `internal/resolve/guardcontract.go`,
test-site counts under `internal/guard`, `internal/cli`).

---

## HELD

1. **C4 "CONSTRUCTION IS NOT ENOUGH" paragraph** — holds. Matches source exactly:
   `reach.go::atomAdmitsValue` (line 502-513) returns `verdict != resolve.GuardFalse`;
   caller `reach.go` reads it positively (line 484-485), `analysis.go::nodeMeetsAll`
   reads it negated (line 438, `!atomAdmitsValue`). The clause correctly forbids the
   two-valued collapse at any consumer and scopes lint-policy as undecided.

2. **C4 dependency paragraph, `renderWrites` removed** — holds. `renderWrites` does
   not exist on `main` (confirmed by source search); the rewritten text correctly
   says the shim "does NOT exist on `main`" and names no present site for it.

3. **§prerequisites A1–A6 with rationale** — holds. Gate line now reads "verified
   (A1–A6)" with the rationale that A6 is the held-side census gating rollout
   shape, not record tidiness.

4. **A4 Pending note re-scoped to HELD values** — holds. Correctly identifies
   `reach.go::heldValues` (source-confirmed: sorts/dedupes via `canonicalValues`,
   no Atoi round-trip) as the second, unmeasured held-side path, distinct from the
   literal side C5 already closes.

5. **§key-discoveries reversal ("held leg IS observable")** — holds. Confirmed against
   `conformKind`'s bare-Atoi int arm (admits `"007"`) and `canonicalValues`'
   sort/dedupe-only behavior — the divergence between lint's parsed guard-seam read
   and the kernel's byte-compare at `TagSet.matches` is real and unclosed on the
   held leg.

6. **§key-discoveries "sole unconformed door" narrowed to KIND-and-SPELLING** — holds,
   and is consistent with §Problem Statement (untouched, in delta scope), which
   still correctly says the owned door is the sole door with "no kind check
   anywhere on the path" — a KIND claim, not a SPELLING claim. The two statements
   operate on different axes and do not conflict.

7. **C2 key-absent arm split (undeclared vs. declared-empty-Kind)** — holds. Confirmed
   against `conformKind`'s switch (line 1853-1865): no `default` case, so an
   empty-`Kind` tag's members pass unvalidated — exactly as the rewrite states.

8. **C2 bool arm bounds rejected population; int arm scopes Atoi width out** — holds.
   `flow_input.go::canonicalValue` (line 738-756) confirmed to call
   `table.ConformValue` before resolution, so the CLI cannot produce a non-`true`/
   `false` bool value — matches the "reachable only from the owned door" claim.
   No platform-width pin found anywhere in `internal/guard` or the conformance
   suite, matching the int-arm claim.

9. **§consequences scope claim with counts** — holds and is accurate: confirmed
   11 `Evaluator{}` literals + 7 `var ev` forms = 18 sites in `internal/guard`,
   plus 4 `guardSeam()` fixture call sites in `internal/cli` — exact match to the
   text's arithmetic.

10. **S4 typed go/ast check + graphlint-site rationale** — holds on its own terms:
    correctly explains why a text grep/regexp cannot distinguish a composite
    literal / `var` / `new` / struct field from a constructor call, and why the
    backstop must be mechanical given C4's new three-valued-consumption
    obligation. `make check` exists as a real target this can be wired into.
    See NEW GAP R2-1 below — the rewrite fixed the mechanism but not two
    in-scope cross-references to it.

11. **C5 commentary on `malformed_predicate_atom` reuse; §technical-design scoped
    to seam only** — holds. Both additions are internally coherent and consistent
    with C5's own scope (predicate ingress only) and C1/C2 (seam only).

12. **F1 second-order break** — holds. The `[match] n = "007"` / `--tag n=007`
    scenario is a correct second-order consequence of C5 (load-time canonicalization)
    combined with the still-unclosed held/CLI leg named in Key Discoveries — no
    internal contradiction, and it correctly cross-references "the held-leg
    divergence (Key Discoveries)."

13. **D-canonical-spelling-at-load admits rejection (c) applies; routes to A4** —
    holds. Confirmed self-consistent: the added paragraph states plainly that C1's
    kernel-byte-compare / seam-parsed split is the exact two-semantics problem
    rejection (c) objected to, that C5 only makes it unobservable for literals,
    and correctly hands the open choice to A4 rather than silently resolving it.

14. **MVV step 1 (persisted-artifact driven reader) and acceptance paragraph
    de-coupled from "rdr navigator" as a library caller** — holds. Step 1 now
    states steps 2/3 set the owned value in the artifact the accessor reads and
    re-run the command; the acceptance paragraph correctly generalizes to "any
    caller reaching the seam through an OWNED value" and explicitly disclaims
    resting on the navigator specifically.

---

## NOT HELD

None. All 14 claimed rewrites hold on inspection.

---

## NEW GAP

**R2-1 — Stale "S4's grep" cross-references after S4's own mechanism changed.**

Item 10's rewrite changed Validation Scenario 4 (S4) from a text-grep backstop to
a typed `go/ast`/`go/packages` check, explicitly rejecting regexp/grep as
insufficient ("A regexp over source cannot separate those forms from a
constructor call and would pass vacuously"). That rewrite did not propagate to
two other delta-scoped anchors that still name the old mechanism:

- **§failure-modes**, "Zero-value evaluator survives migration" (F6): "Detection
  there falls to **S4's grep** and the A3 site audit."
- **§technical-design / Validation, Scenario 5 (S5, the nil-mapping disposition
  scenario)**: "the guarantee **S4's grep** explicitly defers to, asserted
  directly."

Both statements are still substantively true (a missed C4 site whose guards use
only operator-inferred verbs still evades the mechanical check, whatever it's
implemented as), so this is not a factual defect in what the record asserts —
but it now misdescribes the mechanism it names. A reader who reads S4 first
(typed AST pass, explicitly not a grep) and then hits "S4's grep" in F6 or S5
will read it as a contradiction or as evidence the two passages were never
reconciled, which undermines exactly the trust a pre-lock pass exists to
establish. Two-word fix in each spot ("S4's check" or "S4's `go/ast` pass"),
but it is a genuine gap this rewrite opened and left unclosed.

---

## Ledger

| ID | RDR passage | Failure mode | Symptom user sees | Origin |
| --- | --- | --- | --- | --- |
| R2-1 | `cli/0012:F6`, `cli/0012:S9` (nil-mapping Scenario 5) | Stale cross-reference: both still say "S4's grep" after S4 was rewritten (this iteration) to a typed `go/ast`/`go/packages` check, not a grep | A reviewer reading S4 (typed AST pass) then F6/S5 ("S4's grep") sees an apparent contradiction about what the mechanical backstop actually is | Rewrite (item 10, S4) — introduced by narrowing S4's own mechanism without updating its two callers |
