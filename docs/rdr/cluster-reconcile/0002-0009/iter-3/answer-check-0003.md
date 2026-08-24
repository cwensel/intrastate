Model: claude-fable-5

# Answer-vs-fences check — sibling 0003 (guard-predicate-exhaustiveness), iteration 3

Date: 2026-08-24. Read-only. Home: `docs/jdr/0001-resolve-kernel-seam.md` §D6 (JD-16), §D7 (JD-17), interface record JD-16/17/18. Sibling: `docs/rdr/0003-guard-predicate-exhaustiveness.md` (2370 lines, read in full). Line anchors below are current-file lines.

## Status line qualifier (0003:9)

> `Final [joint decision → JDR 0001 §JD-16, §JD-18: Match/Guard routing key; conforming-view enforcer] [locked 2026-08-22 — Gate PASS. …]`

- JD-16 qualifier **still present** although §D6 answered it on 2026-08-24 (stale as a "pending" marker; the home says 0003's landing for D6 is a citation).
- JD-18 qualifier **still present**; JD-18 remains open at the home, and A18 (0003:606-638) is still `Pending`. Correct.
- JD-17 is **not** in 0003's qualifier, though the home names 0003 a JD-17 sibling ("0002, 0003, 0004, 0006") and lists two "Lands in 0003" items for it. Note only — a qualifier omission, not a fenced-text issue.

## JD-16 (§D6 — block-keyed routing, resolved (b))

Answer: the authored block routes. Match blocks admit `eq` and `in` (expanded per member); comparison/existence/`contains` under a match block are a load refusal (0002 category). Every atom under `guard.all`/`guard.unless` is a guard atom regardless of operator, `eq` included. "Lands in 0003 — cites: its participation and can-refuse clauses already read the block."

| Clause | Fenced? | Verdict | Evidence |
| --- | --- | --- | --- |
| Participation clause | fenced (0003:1132-1136) | CONSISTENT | "A dimension **participates** in a row group when any row in that group carries a **guard** atom over that key, in either `all` or `unless`" — block-keyed by construction; no operator test. |
| Match keys are not product dimensions | fenced (0003:1138-1147) | CONSISTENT | "a separate field from its guard both in RDR 0002's authored shape (`[rule.match.*]` vs. `[rule.guard.*]`)" — reads the authored block as the split. |
| Match/guard split per atom per group | fenced (0003:1149-1158) | CONSISTENT | "Participation is therefore decided by where each atom is authored **within the group being proved**" — exactly D6(b)'s "the block is the author's declared intent". |
| Can-refuse clause | fenced (0003:1360-1366) | CONSISTENT | "a value atom over a key **declared optional** **in either block** — `all` or `unless`" — reads block + declared optionality; no operator routing. |
| Presence-dimension clause | fenced (0003:1162-1177) | CONSISTENT | Value atoms (incl. `eq`) project on the guard product; D6 confirms an `eq` guard atom stays a guard atom. |
| Operator/kind matrix | unfenced (0003:823-829) | CONSISTENT | Lists `eq` as a guard operator over `enum`/`bool`/`int`/`scalar`; under D6(b) `eq` under a guard block is a guard atom, so the matrix is the guard vocabulary D6 presumes. 0003 defines no match-block grammar — silent, D6 fills it in 0002. |
| Desk trace steps 2-3 | unfenced (0003:1534-1535) | CONSISTENT | "`status eq "Draft"` atom is authored under `[rule.match.status]`, so it is the grouping key … not a guard dimension" — block-keyed routing witnessed on the fixture. |
| Fixture `foundational-to-cove` (`profile eq "foundational"` under guard) | unfenced (0003:1536) | CONSISTENT | D6 names this fixture as the input class where 0003's withhold reading is the wanted behavior; 0003 already treats the `eq` guard atom as a guard dimension. |
| Illustrative `[rule.guard.all] profile.in = […]` | unfenced (0003:1651-1655) | CONSISTENT | `in` under a guard block is a guard atom (D6(b)); the D6(c) worry ("gap findings for values the author meant to exclude") is exactly the reading 0003 gives a guard-block `in`, which D6(b) endorses as the author's declared intent. |
| Row-group definition (source state × recognized outcome) | unfenced (0003:869-878) | CONSISTENT (silent) | 0003 says nothing about match-block `in` expansion; each expanded row is one candidate row and groups by its own selection context. Nothing to falsify. |

No fenced or unfenced clause in 0003 routes by operator. Search for an operator-keyed routing statement found none.

**"Lands in 0003" items (JD-16):**

| Item | Status | Evidence |
| --- | --- | --- |
| Participation clause reads the block | present | 0003:1132-1136, 1149-1158 (fenced). |
| Can-refuse clause reads the block | present | 0003:1360-1366 (fenced). |
| Citation of §D6 / JD-16 as decided | absent | The only reference is the Status-line qualifier listing §JD-16 as a pending joint decision (0003:9). No clause cites §D6. Substance already matches; only the anchor is missing. |

## JD-17 (§D7 — closed TOML layout, resolved (i)-(v))

Answer items touching 0003: (iii) type-model keys `kind`/`domain`/`min`/`max`/`elements`/`single_valued`/`required` — 0003's own spellings, `required` for the optionality marker, defaults optional; two load categories in 0002 with rules supplied by 0003 — `malformed tag declaration` (kind/domain disagreement; before normalization completes) and `malformed predicate atom` (literal outside domain; after declarations load, before rows are yielded); "0006 mints nothing for either; 0005 maps them under §JD-8". (ii) `[tags.<tag>].accessor` removed. "Lands in 0003 — `min`/`max` spelling becomes a citation to 0002; supplies the two rejection rules."

| Clause | Fenced? | Verdict | Evidence |
| --- | --- | --- | --- |
| Declaration-model ownership | fenced (0003:937-954) | CONSISTENT | "RDR 0002 owns where a declaration is authored and how it is carried … this RDR owns what a declaration means." D7 keeps the wire keys in 0002 and the meaning here. |
| Five kind tokens `enum`/`bool`/`int`/`set`/`scalar` | fenced (0003:943-945) | CONSISTENT | D7 rejects JSON-Schema `type`/`enum` precisely because "`kind` is in every fixture and in 0003's tokens". |
| Finite-domain clause — `{min..max}` wire spelling | fenced (0003:958-966) | CONSISTENT (content) | "the authored form is two fields, `min` and `max`, and both endpoints are **inclusive**" — same keys D7(iii) fixes. 0003 asserts the spelling rather than citing 0002 (see Lands-in table). No meaning change. |
| Optionality default | fenced (0003:970-972) | CONSISTENT | "A declaration carrying **no optionality marker declares the key optional**" — D7: "`required` for the optionality marker (defaults optional, as 0003 mandates)". |
| Single-valued marker clause | fenced (0003:1009-1035) | CONSISTENT | Marker semantics unchanged; D7 spells the key `single_valued`, which 0003 already uses (0003:1520, 1639). |
| Domain/kind agreement — declaration rejection rule | fenced (0003:1039-1052, 1060-1063) | CONSISTENT | "a declaration error that MUST be rejected before normalization completes … rejected by the declaration loader before normalization completes" — matches D7's `malformed tag declaration` (kind/domain disagreement; before normalization completes). |
| Literal-outside-declared-domain rejection rule | fenced (0003:1053-1058, 1063-1065) | CONSISTENT | "rejected by guard parsing after normalization has supplied the declaration to check against, and before resolution" — D7's `malformed predicate atom` window ("after declarations load, before rows are yielded") sits inside 0003's "before resolution" window; D7 pins the earlier bound, 0003 is not falsified. |
| Rejection-rule surfacing route | **fenced (0003:1065-1067)** | **CONTRADICTS (narrow)** | 0003: "Both surface through the predicate semantic kinds this RDR owns (A4) onto **RDR 0006 findings** and the RDR 0005 envelope." D7: "Two load categories in **0002** … **0006 mints nothing for either**; 0005 maps them under §JD-8." 0002's own bucket table (0002:964-968) puts `malformed predicate atom` in the `load` bucket owned by 0002, lint (0006) only for cross-row classes. The phrase "onto RDR 0006 findings" names a reporting surface D7 says does not exist for these two classes. **Meaning of the rules (what is rejected, at which phase, by whom) does not change**; only the emitting document for the artifact does. 0003's own unfenced text already leans D7's way: "a load-time error" (0003:1057), disposition rows "Rejected at load" (0003:1518-1520), and A4 (0003:199) scopes 0006 to "graph-lint finding codes for non-exhaustive and overlapping row groups". Candidate citation repair rather than a design conflict — the gate decides. |
| Disposition table rows for the two rejections | unfenced (0003:1519-1520) | CONSISTENT / same route note | Phase and owner agree with D7; the "→ RDR 0006 finding → RDR 0005 envelope" arrow repeats the route noted above (unfenced). |
| Assignment-count defaults (unmarked enum = `2^|domain|`) | fenced (0003:1417-1426) | CONSISTENT | D7's closing "Observation for 0003's next touch" restates this default as-is and explicitly defers it ("not this decision"). |
| Illustrative declaration `optional = false` | unfenced (0003:1634-1646, 1658-1660) | CONSISTENT (spelling differs) | D7's key is `required`; the illustrative block spells the same marker as `optional = false` and is labelled "Illustrative only — RDR 0002 owns the authored container". Same meaning (always-present); spelling only. Other keys in the block — `kind`, `domain`, `single_valued`, `elements` — match D7 exactly. |
| A16 evidence quoting 0002's "optional accessor reference" | unfenced (0003:519) | CONSISTENT (stale quote) | Quotes 0002 pre-D7 text ("beside `provenance` and the optional accessor reference"); D7(ii) removes `[tags.<tag>].accessor`. A quotation of a peer's superseded sentence inside a Critical Assumption evidence line — no 0003 clause depends on the accessor reference. |
| MVV Scenario 4 (malformed declarations / literal outside domain) | unfenced (0003:2277-2290) | CONSISTENT | "rejected before resolution with a predicate semantic kind that RDR 0006 can map to a lint finding and RDR 0005 can map to the structured CLI gateway" — "can map" is weaker than the fenced route phrase; same note applies. Phase ("before normalization completes") matches D7. |
| D7(i) initial/terminal, (iv) write-replaces, (v) `[model.metadata]` | — | CONSISTENT (silent) | 0003 says nothing about any of these; the answer merely fills. |

**"Lands in 0003" items (JD-17):**

| Item | Status | Evidence |
| --- | --- | --- |
| `min`/`max` spelling becomes a citation to 0002 | present as assertion; citation form absent | 0003:958-966 (fenced) states the authored form itself ("two fields, `min` and `max`"); it does not cite 0002 for the spelling. Spelling agrees with D7, so no contradiction — the form is what differs. |
| Supplies rejection rule 1 — `malformed tag declaration` (kind/domain disagreement, before normalization completes) | present | 0003:1039-1052, 1060-1063 (fenced). |
| Supplies rejection rule 2 — `malformed predicate atom` (literal outside domain, after declarations load, before rows yielded) | present | 0003:1053-1058, 1063-1065 (fenced); 0003 spells the semantic kind `literal-outside-declared-domain` and the window "before resolution" (superset of D7's). |
| Rules routed as 0002 load categories, 0006 mints nothing | contradicted (narrow) | 0003:1065-1067 (fenced) routes both "onto RDR 0006 findings"; see the CONTRADICTS row. |

## JD-18 (open — note only)

Qualifier present on 0003:9 (`§JD-18 … conforming-view enforcer`); A18 `Pending` (0003:606-638); conformance clause still says "Conformance currently has no enforcer, and that is a booked gap" (0003:995-1005). Unchanged since iter-2; nothing to check against an answer.

## Summary

- JD-16: every fenced clause CONSISTENT; both Lands-in clauses present; the §D6 citation anchor is absent and the Status qualifier still lists JD-16 as pending.
- JD-17: one narrow CONTRADICTS on fenced text — the surfacing route "onto RDR 0006 findings" (0003:1065-1067) versus D7's "0006 mints nothing for either"; rule substance, phase, and ownership all CONSISTENT. `min`/`max` spelling agrees but is stated, not cited. Two unfenced stale spellings/quotes (`optional = false`; "optional accessor reference"). JD-17 missing from 0003's qualifier.
