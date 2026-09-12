Model: claude-sonnet-5
Iteration: 3 (pass 2 — delta-scoped)

## Persona 3 — QA / Tester (delta-scoped second pass)

Scope: rewritten `S1`, `S5`, `S6`, `S8` oracles; rewritten `C2`; new `Determinacy:` line; `§pre-lock-mini-checks` `oracle`/`trace`/`disposition` tables against the rewritten scenario text. Not re-raising pass-1-resolved items per the brief's disposition list.

### Widening note

Widened from `S`/`C`/`MVV` into `§pre-lock-mini-checks` (`oracle`, `trace` tables) and `A6`, because S6's new "DECODE" oracle cannot be judged testable from S6's own text alone — the decode mechanism, if it exists, would have to be named either in a contract clause (C2) or in the mini-check tables that operationalize the MVV steps. Neither holds it.

### Findings

**1. `0021:S6` — "DECODE back to its source tag value" names no decode mechanism; test is unwritable as stated. (HIGH)**

S6's rewritten oracle reads: "each hostile label is also asserted to DECODE back to its source tag value, so a well-formed but mis-escaped label fails." This is the right strengthening in principle (exit 0 alone missed A6's real escaping bug), but as written it substitutes one unrunnable verb for another. Contrast with the other three rewritten oracles in this same pass:

- `S1` names the exact mechanism: `GODEBUG=randmapiter=1`, cross-checked by the mini-check `oracle` table row ("A5 §5a: a map probe emitted under `randmapiter=1`, 25 processes").
- `S5` names the exact comparison: `jq -r` unwrap against `document + "\n"`, cross-checked by the `trace` table row 4 and `A1`.
- `S6` names only the verb "DECODE," with no library, flag, or grammar to invoke it. `A6`'s evidence is the only place any escaping detail appears at all ("the escaping ORDER (backslash, then quote, then newline) is sound"), and that is evidence prose, not a normative clause — `C2` explicitly puts DOT "styling/attributes" outside the normative boundary, and no contract clause pins a DOT string-escaping grammar. There is no `dot -Tplain`/`-Txdot` round-trip, no unescaping library, and no reference decoder named anywhere in the record (checked full-text for "decode", "unescap", "dot -T", "graphviz").
- The `§pre-lock-mini-checks` `oracle` table, row "3 — DOT set equals `reach` block," was NOT updated for this pass: it still reads "A6 hostile fixture: quotes/backslash/newline/non-ASCII survive to `dot -Tsvg` exit 0 with exact set match" — the same exit-0-plus-set-match oracle S6's own prose just said is insufficient ("Exit 0 is NOT sufficient"). The `trace` table's row 3 witness is identical and equally stale. S6 and its own mini-check backing now contradict each other on what the passing bar is.

**Test prevented**: a tester cannot write "assert hostile label X, after DOT round-trip, decodes to source tag value X" without a named decoder. Given only graphviz's `dot -Tsvg`/`-Tplain` output and the RDR's own escaping order, a tester has no specified inverse function to apply to the emitted label bytes, and no golden to diff against for the decoded form specifically (only the SVG-renders-without-error and node/edge-SET-match tests are actually specified in the mini-check tables). The test A6's own spike bug (separator injected as two bytes before `dotQuote` doubles the backslash) would need to catch is exactly a decode-mismatch case — but the record gives no way to detect it other than manual/ad hoc inspection, which is not a pass/fail criterion.

**Remedy direction (not prescribing, flagging the gap)**: either (a) `C2` or a new clause pins the DOT label-escaping grammar as normative (so "decode" means "apply grammar G's inverse"), and the mini-check `oracle`/`trace` rows are updated to cite it, or (b) S6's oracle is restated in terms of a concrete, tool-based round-trip (e.g., parse the emitted DOT with a library that exposes unescaped label strings, or diff against a second golden of decoded labels) and the mini-check tables are updated to match.

### Not findings (checked, no defect)

- `0021:S1`, `0021:S5`: oracles are concrete and mechanically executable as written; mini-check `oracle`/`trace` rows agree with the scenario text.
- `0021:S8`: "no document on stdout" assertion and code spellings match `C1`/`C4`'s disposition table; consistent.
- `0021:C2` field spellings, sort orders (nodes by key, edges by (from,to,rule)), and empty/absent (`[]`/`{}` vs ABSENT) rules: each is exercised by MVV step 2 / `S2`'s byte-identical golden, which necessarily pins spelling and order and is separately witnessed for empty/absent by the A5 `empty_probe: []` fixture cited in the `oracle` table. Sort order has no *dedicated* negative-control fixture (unlike `S1`'s map-order provocation), but a byte-identical golden containing ≥2 reach nodes/edges catches an order violation as a side effect of the existing oracle — not a gap worth blocking on.
- `Determinacy:` line (line ~494): claims "C2 fixes exact field spellings, sort orders, and empty/absent renderings that Scenario 2's goldens pin" — accurate; `S2`'s golden-fixture byte comparison is the mechanism that pins all three, even though sort order isn't separately called out as its own witness.
