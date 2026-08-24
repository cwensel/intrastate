Model: claude-fable-5

# Answer-vs-fences check — sibling 0009 (escape-row-shape-conformance-ownership) — iteration 3

Date: 2026-08-24
Sibling rev: 526c481 (status qualifier only since iter-2; body frozen 2026-08-11)
Home: `docs/jdr/0001-resolve-kernel-seam.md` §D5 (222-255), JD-15 entry (534-540)

## Scope

0009's qualifier (0009:9) names three home entries:
`§JD-5, §JD-8, §JD-15: precondition precedence; envelope carrier for row identities; <clear> at the write accessor`.

- **JD-15** — ANSWERED 2026-08-24 by §D5(a). Checked below.
- **JD-5, JD-8** — still open at the home; the qualifier still names both. Not checked
  (no answer to check against). Noted only.

## The answer (JD-15 / §D5)

Home 0001:534-537: "the sentinel is reserved: refused as an authored tag value at
load (0002); a `<clear>` write removes the key, read-back asserts absence, clearing
an absent key succeeds (0004)." §D5(a) at 0001:238-243 adds "a read yielding
`<clear>` as a value is unreadable. Kernel `Row` and 0009 untouched." Lands-in
paragraph 0001:251-255: "**0009** is unchanged".

Rejected branch (b) 0001:244-246 is the one that would have touched 0009: "Separate
`Clears []string` on the row … reopens 0001's `Row` and 0009's fenced Writes-only
predicate". (a) was chosen precisely to avoid that.

## Per-clause check — fenced text

| # | Clause (0009 lines) | Fenced? | Verdict | Evidence |
| --- | --- | --- | --- | --- |
| F1 | 822-834 Writes-only predicate: "a Row with a non-empty Escape list MUST have an empty Writes slice. (Authored clears normalize to `<clear>` writes per RDR 0002, so the Writes predicate carries both "no writes" and "no clears" at the kernel boundary.)" | fenced | **CONSISTENT** | The reduction depends on one fact — a clear is rendered as a `<clear>` entry *in `Writes`* — and §D5(a) keeps that rendering ("Three Final documents and the kernel already carry the sentinel", 0001:248). Reserving the sentinel strengthens the clause: with `<clear>` refused as an authored value, an empty `Writes` now provably carries no clear, where before an authored literal could have hidden in the value space. Meaning unchanged. |
| F1b | 826-834 "The predicate is Writes-only and does NOT extend to NextTags: A4 settled that owned state is reachable only through a write accessor" | fenced | **CONSISTENT (silent)** | §D5 says 0002 renders `<clear>` "in both `Writes` and the next-state tags" (0001:224-225) but decides only what the *write accessor* does with the `Writes` entry. NextTags remain non-mutating per A4; the answer neither touches nor needs NextTags. |
| F2 | 836-852 authored path: "an authored escape rule carrying a write block or clear list MUST be rejected at table load under RDR 0002's "malformed escape declaration" validation class … Rejection keys on the PRESENCE of the block" | fenced | **CONSISTENT (silent)** | Presence-keyed; never inspects a tag value. The new 0002 load category ("refusing `<clear>` as an authored tag value") is a *different* category on a *different* shape (a tag value, not an escape rule) and does not collide. |
| F3 | 854-877 kernel entry precondition; error path, no refusal kind | fenced | **CONSISTENT (silent)** | No mention of clear/sentinel; the answer lands nothing on `Resolve`. |
| F4-F6 | 879-1028 typed error, sentinel, aggregate, ordering, `CheckValid` | fenced | **CONSISTENT (silent)** | Untouched by the answer. |
| F7 | 1030-1077 CLI wrap / envelope carrier | fenced | **CONSISTENT (silent)** | Untouched by the answer (rides JD-8, still open). |
| F8 | 1079-1084 "The shared resolve.Row type keeps its single shape … No ordinary/escape type split and no row-kind field is introduced" | fenced | **CONSISTENT** | §D5 rejected (b) `Clears []string` on the row — the only option that would have added a field to `Row`. (a) states "Kernel `Row` and 0009 untouched" (0001:243). |

No fenced clause is falsified. No CONTRADICTS.

## Per-clause check — unfenced text that mentions the sentinel

| # | Lines | Verdict | Evidence |
| --- | --- | --- | --- |
| U1 | 296-299 Key Discovery: "clears normalize into writes (RDR 0002: a rule-level `clear` entry "renders as a `<clear>` write"), so at the kernel boundary "no writes and no clears" reduces to an empty `Writes` slice on escape rows." | CONSISTENT | Same reduction as F1; the cited 0002 sentence still exists verbatim (0002:821-823). |
| U2 | 1660-1665 Phase 3: "one canonical authored-clear case pins the `<clear>`-sentinel-write representation identically for the kernel suite and the normalizer suite" | CONSISTENT | Pins the *representation* of an authored clear (a `<clear>` entry in `Writes`), which (a) preserves. It does not pin what the write accessor does with it, so 0004's new remove-key semantics are not contradicted. Not widened: the new 0002 refusal of an authored `<clear>` *value* is not among 0009's fixtures — that is 0002's landing, not 0009's, so the omission is expected, not a gap. |
| U3 | 1788-1795 Testing scenario 8 (write block / clear list / empty write block rejected at load) | CONSISTENT | Presence-keyed, as F2. |
| U4 | 1919-1922 References: 0002 "clear-as-`<clear>`-write normalization" | CONSISTENT | Still true of 0002 (0002:821-823). |
| U5 | 0009 contains **no** statement that `<clear>` is unreserved, a plain value, or storable | CONSISTENT | `grep -n -i clear` over 0009 returns only the sites listed here; the "unreserved" claim §D5 overturns lives in 0002 (0002:956), not 0009. |

## "Lands in" items named by the home (0001:251-255)

| Home item | Target | Exists in target today? |
| --- | --- | --- |
| "**0009** is unchanged" | 0009 | n/a — nothing to land; 0009 body verified consistent above |
| one load category refusing `<clear>` as an authored tag value | 0002 | **ABSENT** — 0002:956 still reads "`<clear>` unreserved in the value space" (dump-grammar table); 0002:1471 category list carries no such category |
| Round-Trip lossy-site note drops its `<clear>` item | 0002 | **ABSENT** — 0002:936-938 still lists "the `<clear>` … renders identically to a cleared tag" as a lossy site |
| remove-key / read-back-absent / idempotent-clear / unreadable-on-read clauses + MVV clearing-rule scenario | 0004 | **ABSENT** — 0004 mentions "clear" only in its Status qualifier (0004:9); body contains the word zero times, as §D5 itself records (0001:226) |
| **0006** reads a `<clear>` write as removal by citation | 0006 | not in this task's scope (assigned to 0006's checker) |

Reported for the gate; 0002/0004 are not this task's sibling.

## Iter-2 cosmetic items on 0009 — re-check

| Item (iter-2 ledger) | Status now | Evidence |
| --- | --- | --- |
| retired 0002 "row kind `escape`" wording ×3 | **STANDS** (×2 found by literal grep; third is the paraphrase at 1213) | 0009:244 and 0009:492 quote "row kind `escape`, its normal predicate set, source rule id, source locator, and modeled failure class list"; current 0002:1465 reads "escape rows retain their modeled failure class list — from which the derived `escape` kind is computed — and carry neither writes nor next-state tags". Substance identical (no writes, no next-tags); only the phrase is stale. 0009:1213 "reopens RDR 0002's single-shape row rendering" is unaffected. Unfenced; cosmetic; single-RDR (0009). |
| "implemented" at 0009:1184 | **STANDS** | 0009:1184: "Predecessor (RDR 0005 + `internal/cli/clierr`, implemented)". README index: 0005 is "Final [joint decision → JDR 0001 §JD-8, §JD-9]", not Implemented. The `clierr` package is shipped, so the row's *Status: Available* is correct; only the parenthetical mislabels 0005. Unfenced; cosmetic. |
| `Refusal.Guard` at 0009:421 | **STANDS, still true** | `internal/resolve/resolve.go:270-272` declares `Guard string` on `Refusal` at HEAD. Unfenced A3 evidence prose ("the asserted properties are `Refusal.Kind` / `Refusal.Guard`"). Could go stale if §D1's guard-seam change renames the field; cosmetic, rides 0009's next touch. |

All three remain citation-repair items pending 0009's next touch; none is fenced, none changes meaning.

## Result

- JD-15 / §D5 vs 0009: **CONSISTENT** on every fenced clause (F1-F8) and every unfenced sentinel mention (U1-U5). The home's "0009 is unchanged" holds.
- JD-5, JD-8: still open; 0009:9 qualifier still names both — carry.
- No CONTRADICTS. No new finding on 0009.
