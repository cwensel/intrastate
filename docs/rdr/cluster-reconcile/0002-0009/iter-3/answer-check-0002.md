Model: claude-fable-5
Date: 2026-08-24
Task: B — scoped answer-vs-fences check, sibling 0002 (transition-table-as-reviewable-data)
Sibling revision read: 526c481 (docs/rdr/0002-transition-table-as-reviewable-data.md, 1548 lines)
Home read: docs/jdr/0001-resolve-kernel-seam.md @ 4f9639c — §D5 (222-255), §D6 (257-298), §D7 (300-378), interface record JD-15/16/17 (534-553)

Evidence only. The gate decides disposition. No RDR, JDR, or README was edited.

## 0002 Status line (0002:9)

> `Final [joint decision → JDR 0001 §JD-15, §JD-16, §JD-17: `<clear>` at the write accessor; Match/Guard routing key; wire keys for initial/terminal, write-replaces, accessor metadata, type-model fields]`

The qualifier names all three answered entries. It is the only line that moved since iter-2 (brief: "Status qualifier only (1 line)"). No fenced text was touched, so every finding below is by construction a frozen-fence vs. answer comparison.

---

## JD-15 / §D5 — `<clear>` reserved; refused as an authored tag value at load (0002); remove-key / read-back-absent (0004)

### Clause checks

| # | 0002 clause | Fenced? | Verdict | Meaning changes? |
| --- | --- | --- | --- | --- |
| 15.1 | 0002:820-824 `Clearing a tag MUST be represented by an explicit rule-level `clear` entry that normalization renders as a `<clear>` write. Absence from both the write block and the clear list MUST NOT imply deletion.` | fenced | **CONSISTENT** — §D5(a) keeps the sentinel and the rule-level `clear` list ("Kernel `Row` and 0009 untouched"); the answer reserves the value, it does not change how a clear is authored or rendered. | no |
| 15.2 | 0002:672-680 next-state tags and writes "the same rendered set, including `<clear>` entries" | fenced | **CONSISTENT** — reservation does not alter carriage; §D5 explicitly leaves the kernel `Row` shape alone. | no |
| 15.3 | 0002:726-729 dump "MUST carry … writes including `<clear>` entries" | fenced | **CONSISTENT** | no |
| 15.4 | 0002:827-839 load category list (data-level categories "including at minimum …") | fenced | **CONSISTENT-by-silence** — the list is open ("at minimum") so a new category `reserved tag value` / `<clear>` refusal can be added without falsifying it; but it is **absent** today (see Lands-in table). No existing category covers an authored `<clear>` literal: `malformed predicate atom` is operator/literal-form only (0002:371-373), `write to non-owned tag` is provenance only. | no (additive) |
| 15.5 | 0002:931-942 Round-Trip prose: "Three sites make the rendered form lossy today … the `<clear>` sentinel is not reserved in the tag-value space, so an authored value of `<clear>` renders identically to a cleared tag" | **unfenced** (Round-Trip / Inverse Invariants section prose, no ```normative fence) | **CONTRADICTS** — §D5 resolves "(a) Reserve the sentinel"; the JD-15 entry reads "the sentinel is reserved: refused as an authored tag value at load (0002)". The frozen text asserts the opposite fact ("not reserved"). | Yes in fact, no in invariant: the round-trip invariant itself (0002:915-922, over the normalized value) is unaffected; only the lossy-site enumeration becomes stale (three sites → two). |
| 15.6 | 0002:956 `fidelity` table row `dump`: "**no grammar**: `<clear>` unreserved in the value space; separators unescaped; suffix separator unreserved" | unfenced (Conditional Mini-Checks table) | **CONTRADICTS** — same fact as 15.5, restated in the table. | same as 15.5 |
| 15.7 | 0002:949-950 "the rendered form is lossy at three named sites and carries no inverse" | unfenced | **CONTRADICTS** (count) — after §D5 the count is two; the "no inverse" half stands. | count only |
| 15.8 | 0002:101 A1 Evidence "multi-tag writes, and explicit clears"; 0002:165-168 A7 "writes including `<clear>`" | Critical Assumptions (unfenced) | **CONSISTENT** — witness of rendering; unaffected by reservation. | no |
| 15.9 | 0002:1388-1390 MVV permutation control "carries the escape row, the `in`-expansion, `<clear>`, and inherited contexts" | MVV (unfenced) | **CONSISTENT** | no |
| 15.10 | 0002:1463 Scenario 1 "explicit clears … decode"; 1471-1472 Scenario 3 category list | Testing Strategy (unfenced) | **CONSISTENT / silent** — no fixture for the new refusal; owed if the category lands (see Lands-in). | no |

**Fenced verdict for JD-15: no CONTRADICTS on fenced text.** The two contradictions (15.5, 15.6, 15.7) are unfenced Round-Trip / mini-check prose that state the pre-decision fact the home explicitly overturns ("Reserving it closes the ambiguity and the 'lossy dump site' 0002 lists", home:248-249).

### "Lands in 0002" items (§D5, home:251-255)

| Item | In 0002? | Anchor / note |
| --- | --- | --- |
| "one load category refusing `<clear>` as an authored tag value" | **absent** | 0002:827-839 category fence has no such member; 0002:370-376 Technical Design load list has none; Scenario 3 (0002:1471) has no fixture. Nothing contradicts adding it (list is "at minimum"). |
| "the Round-Trip lossy-site note drops its `<clear>` item" | **contradicted** | 0002:936-938 still lists it as lossy site one; 0002:956 and 0002:949-950 repeat it. Unfenced. |

---

## JD-16 / §D6 — routing is block-keyed; match blocks admit `eq` and `in` (expanded per member), refuse other operators at load; every guard-block atom is a guard atom regardless of operator

### Clause checks

| # | 0002 clause | Fenced? | Verdict | Meaning changes? |
| --- | --- | --- | --- | --- |
| 16.1 | 0002:693-695 `**The unified predicate set splits across the kernel's two predicate fields by operator, and this RDR owns the split — at the kernel handoff, not in the normalized value.**` | **fenced** | **CONTRADICTS** — §D6 resolves "(b) Block-keyed" and names "(a) Operator-keyed (0002 today)" as the rejected option. The bolded clause title names the rejected key. | **Yes** — the routing key changes from operator to authored block. |
| 16.2 | 0002:704-708 `the handoff MUST route each atom to exactly one of them: atoms whose operator is equality on a declared tag populate `Match`; every other atom — set membership, comparison, existence, and every `unless` atom regardless of operator — belongs to the guard predicate.` | **fenced** | **CONTRADICTS** — under §D6(b): an `eq` atom under `guard.all` is a guard atom (frozen text sends it to `Match`); an `in` atom under a match block expands into per-member `Match` rows (frozen text sends `in` to the guard). Two input classes flip. Note the frozen text already routes `unless`-`eq` to the guard, which §D6 cites as "0002 already does it for `unless`" — so the fence is a hybrid the answer generalizes. | **Yes** |
| 16.3 | 0002:710-713 `Routing an atom into `Match` that the kernel cannot evaluate as an equality tag would make it silently fail to match rather than reach the guard seam, and routing an equality atom into the guard defers a decidable check to a seam that can report `guard_unevaluable`.` | **fenced** | **CONTRADICTS** (rationale) — the second half is the argument §D6 rejects: home:293-296 says that sentence "names the wanted behavior for a guard over an optional key". The first half survives as the reason for the match-block operator restriction. | Yes (rationale inverts; the load refusal it motivates is the replacement) |
| 16.4 | 0002:713-716 "this clause fixes only which atoms are the guard's, which is stable across that reshape" | fenced | **CONTRADICTS** (by consequence) — the "which atoms" is precisely what §D6 changes. | yes |
| 16.5 | 0002:539-545 `Guard predicates MUST be represented as positive `all` predicates and negative `unless` predicates … each atom in that set MUST retain … the block (`all` or `unless`) it was authored in` | **fenced** | **CONTRADICTS** (vocabulary) — the fence enumerates the block as two-valued (`all`/`unless`); §D6 lands "match atoms carry block `match` (closing the two-valued/three-valued block vocabulary finding)". Scenario 2 (0002:1465) already says three-valued "(`match`/`all`/`unless`)" and the spike emits `@match` (output.txt line 4 `status.eq=Draft@match`), so the fence is the odd one out. | Yes — block domain gains `match`. Under block-keyed routing the block is load-bearing, not merely carried. |
| 16.6 | 0002:398-399 (Technical Design prose) "each atom still records the block it was authored in"; 0002:406-408 "the kernel's `Match`/`Guard` routing is applied at the handoff … not stored in the normalized value" | unfenced | **CONSISTENT** — handoff placement and block retention both survive §D6; only the key changes. §D6 makes the retained block the routing key, which this prose supports. | no |
| 16.7 | 0002:623-631 outcome binding `**Outcome binding reads the match blocks only** … exactly one atom on `recognized`, using `eq` or `in` … an `in` atom expands into one candidate row per member` | fenced | **CONSISTENT** — this is the "expansion 0002 already defines for `recognized`" that §D6(b) generalizes to every match-block `in` atom. The fence is narrower (recognized only) but not contradicted; generalization is additive. | no |
| 16.8 | 0002:636-641 `recognized` atom under `guard.all`/`guard.unless` refused as malformed outcome binding | fenced | **CONSISTENT** — compatible with block-keyed routing (a guard-block atom is never selection). | no |
| 16.9 | 0002:431-435 layout fence "rule predicates live under `[rule.match.<tag>]`, `[rule.guard.all.<tag>]`, and `[rule.guard.unless.<tag>]`" | fenced | **CONSISTENT** — the blocks §D6 keys on already exist in the layout. | no |
| 16.10 | 0002:827-839 load category list | fenced | **CONSISTENT-by-silence** — "at minimum"; the match-block operator restriction category is **absent** (no `malformed predicate atom` sub-rule for comparison/existence/`contains` under a match block; 0002:371-373 defines that category by operator/literal form only). | additive |
| 16.11 | A12 (0002:267-282): "the union of the constructed row's `Match` and guard atoms equals the normalized atom set and the intersection is empty — including the `unless`-block equality atom (`status.eq=closed@unless` in the kata fixture), which routes to the guard despite being an equality operator." | Critical Assumption (unfenced) | **CONSISTENT** on the totality property; the discriminating case is still valid under §D6 (an `unless` atom is a guard atom). But §D6 lands "A12's MVV" restated — the property "exhaustive/disjoint routing becomes structural" and the discriminating witness should become an `eq` under `guard.all`, which is the class the operator-keyed rule got wrong. Present but keyed on the old rationale ("despite being an equality operator"). | partially |
| 16.12 | A9 (0002:220-234) "the operator-based `Match`/`Guard` split … pinned to the handoff" | Critical Assumption (unfenced) | **CONTRADICTS** (wording) — names the split "operator-based". The verification (one normalized value feeds both scenario 2 and scenario 4) survives. | wording only |
| 16.13 | Scenario 2 A12 paragraph (0002:1470): "The kata fixture's `status.eq=closed@unless` atom is the discriminating case — an equality operator that must route to the guard because its block is `unless`." | Testing Strategy (unfenced) | **CONSISTENT** — literally a block-keyed justification ("because its block is `unless`"). §D6 wants the sharper witness (an `eq` under `guard.all`); not present. | no |
| 16.14 | Scenario 4 (0002:1474-1475): "one tag-set in which a sibling candidate's guard is unevaluable" — witness unspecified | Testing Strategy (unfenced) | **silent** — §D6 lands "Scenario 4's unevaluable-sibling witness becomes an `eq` guard atom over an absent optional key". The scenario does not name the witness; the MVV fixture requirement (0002:1360-1362) requires "one guard that actually carries a predicate" without specifying its operator. | no (absent, not contradicted) |
| 16.15 | Desk trace (0002:988) "`continue-prelock … atoms=[…profile.eq=small@unless…]`" | unfenced | **CONSISTENT** | no |
| 16.16 | Illustrative TOML (0002:1016-1020) `[rule.guard.all] iter.lt = 3` / `[rule.guard.unless] profile.eq = "small"` | illustrative | **CONSISTENT** — both atoms are guard atoms under either key. | no |

**Fenced verdict for JD-16: CONTRADICTS on fenced text** — 0002:693-717 (the whole handoff-routing fence, at 16.1/16.2/16.3/16.4) and the two-valued block enumeration at 0002:539-545 (16.5). Meaning changes: the routing key flips from operator to authored block; two input classes (`eq` under a guard block; `in` under a match block) move sides; a new load refusal (comparison/existence/`contains` under a match block) has no home.

### "Lands in 0002" items (§D6, home:298-303)

| Item | In 0002? | Anchor / note |
| --- | --- | --- |
| "the routing clause restated as block-keyed" | **contradicted** | 0002:693-717 fence is operator-keyed, titled as such. |
| "the match-block operator restriction as a load category" | **absent** | No category in 0002:827-839; no scenario 3 fixture (0002:1471). |
| "the `in`-expansion generalized to every match block" | **absent** (recognized-only present) | 0002:623-631 fence and Load-Bearing Identity (0002:876-879: "the outcome literal when an `in` atom on `recognized` expands the rule") both scope expansion to `recognized`. The expansion-suffix contract (0002:643-648, 0002:751-753 "an expansion suffix is a member of the `outcomes` alphabet") is **implicitly contradicted** by generalization: a non-`recognized` `in` expansion would need a suffix that is not an alphabet member, and the total-order argument at 0002:751-755 rests on the suffix being an alphabet member. This is a consequence the home does not name; flagged as evidence. |
| "match atoms carry block `match` (closing the two-valued/three-valued block vocabulary finding)" | **contradicted in the fence / present elsewhere** | 0002:542 fence says "(`all` or `unless`)"; 0002:1465 scenario 2 says "(`match`/`all`/`unless`)"; spike output emits `@match`. |
| "A12's MVV" (restated block-keyed) | **present, old rationale** | 0002:267-282 and 0002:1470 — property present; discriminating witness still the `unless`-`eq` case, phrased "despite being an equality operator". |
| "Scenario 4's unevaluable-sibling witness becomes an `eq` guard atom over an absent optional key" | **absent** | 0002:1474-1475 names no witness operator. |

---

## JD-17 / §D7 — root `[initial]` and `terminal`; `[read.<id>]`/`[write.<id>]`/`[gate.<id>]` with `keys` as the binding; seven type-model keys; write-replaces clause; `[model.metadata]`; two load categories

### Clause checks

| # | 0002 clause | Fenced? | Verdict | Meaning changes? |
| --- | --- | --- | --- | --- |
| 17.1 | 0002:430-441 `The source schema MUST use the Resolve spike field layout: root `outcomes`, `[model]`, `[tags.<tag>]`, `[accessors.<id>]`, `[context.<id>]`, `[[rule]]`, and `[dump]`.` | **fenced** | **CONTRADICTS** — §D7(ii) replaces `[accessors.<id>]` with `[read.<id>]`, `[write.<id>]`, `[gate.<id>]`; §D7(i) adds root `[initial]` and `terminal`; §D7(v) adds `[model.metadata]`. The fence is a closed enumeration ("MUST use the Resolve spike field layout") coupled to the strict-decoder obligation (0002:465-470: unmapped keys refused), so a document authored to §D7 is refused `unknown schema field` under the frozen text. | **Yes** — layout membership. |
| 17.2 | 0002:437-440 `An `[accessors.<id>]` entry carries `mode` and `path`, as the spike fixtures author it; this RDR validates only that a referenced accessor id resolves to a declared entry (`unknown accessor` otherwise)` | **fenced** | **CONTRADICTS** — §D7(ii): no `mode` (capability is the table name); fields are `role`, `path`, `keys`, `timeout`, `read_back`; and 0002's validation grows beyond id-resolution to the binding validations ("the loader refuses a key served by zero or two readers, a written or cleared key not in exactly one writer, and a writer key that is not owned"). Note the spike fixture authors `mode = "read-write"` on one entry (rdr-fixture.toml:42-44), which §D7(ii)'s one-capability-per-table shape cannot express as a single entry. | **Yes** — accessor table shape and 0002's validation scope. |
| 17.3 | 0002:813-814 fence: declaration authored "under `[tags.<tag>]`, beside `provenance` and the optional accessor reference"; 0002:348-351 Technical Design item 2 "optional accessor reference for observed or owned read-back" | fenced (813-814) / unfenced (348-351) | **CONTRADICTS** — §D7(ii): "the tag-side `accessor` reference is a second copy and is removed"; `keys` on the capability table is the binding. Spike fixture authors `accessor = "rdr-status"` under `[tags.status]` (rdr-fixture.toml:11). | **Yes** — binding direction (tag→accessor becomes accessor→keys). |
| 17.4 | 0002:883-888 Load-Bearing "Wire / byte format — … `[accessors.<id>]` … The RDR and kata spike fixtures are the canonical examples implementation tests must promote." | Load-Bearing Decisions (unfenced) | **CONTRADICTS** — same layout; and "fixtures re-authored" (home:365-366) conflicts with "canonical examples … must promote" until the fixtures are re-authored. 0002:1376-1378 ("a fixture may not be narrowed on promotion, only extended") is compatible with re-authoring that adds/renames tables only if renames count as extension — evidence, not a call. | yes |
| 17.5 | 0002:807-818 fence: type model "value kind, and optionally a finite domain, an optionality marker, a set-element universe, and a single-valued marker. RDR 0003 is the normative home … MUST NOT be restated here." | fenced | **CONSISTENT** — the five concepts match §D7(iii)'s seven keys (`kind`, `domain`, `min`/`max` for the numeric domain, `elements`, `single_valued`, `required`). The fence names concepts, not wire keys; §D7 fills the spelling. The "MUST NOT be restated here" governs semantics, and §D7 assigns 0002 only the spelling (0003's own illustrative spellings). Wire keys are **absent** from 0002 except `kind` in the spike fixture. | no (fills silence) |
| 17.6 | 0002:370-371 "Load-time validation rejects malformed tags …" (prose); 0002:827-839 category fence | prose / fenced | **CONSISTENT-by-silence** — `malformed tag declaration` and `malformed predicate atom` (literal outside domain) per §D7(iii): the fence's `malformed predicate atom` is defined (0002:371-373) as "unknown operator, or a literal ill-formed for its operator" — domain-membership is not in that definition, so §D7 **extends** the category's meaning; the fence's "at minimum" admits `malformed tag declaration`, which is absent from the fence (only the prose "malformed tags" at 0002:370 gestures at it). | additive; `malformed predicate atom` definition widens |
| 17.7 | 0002:820-824 clear fence; 0002:672-690 next-state/writes fence | fenced | **CONSISTENT** with §D7(iv) write-replaces — "a write assigns a tag's whole value and supplants what was held" is not stated in 0002 and nothing in 0002 states an accumulate form. **Absent**, not contradicted. | no |
| 17.8 | 0002:1057 Infrastructure Audit: load entry "takes already-read bytes plus a source id … performs no file I/O" | unfenced | **CONSISTENT** — §D7 adds no I/O. | no |
| 17.9 | Terminal / initial: 0002 has no root-node or stop-set concept anywhere (grep `initial`/`terminal` hits only the Status line, `terminal-archive` rule id, and "A1–A8 are terminal"). | — | **absent** — silent; §D7(i) fills. But 17.1's closed layout fence refuses the keys. | additive once layout opens |
| 17.10 | `[model]` fence 0002:443-445 "`[model]` MUST contain `id` and `version`" | fenced | **CONSISTENT** — `[model.metadata]` is a sub-table; "MUST contain" is a lower bound. Strict decoding (0002:466) would still refuse `metadata` until mapped; §D7(v) "carries untouched and never interprets" is compatible with a mapped free-form table. | no |
| 17.11 | 0002:1044-1045 Capability table rows for 0004 "This RDR may reference accessors but does not execute them" | unfenced | **CONSISTENT** — reference vs execute split survives; the reference site moves. | no |
| 17.12 | Scenario 1 (0002:1463) "accessor references … decode"; Scenario 3 (0002:1471) "unknown accessors" | unfenced | **CONSISTENT** for `unknown accessor` (an id referenced by nothing after `accessor` removal? — under §D7 the tag-side reference is gone, so `unknown accessor` has **no remaining trigger site** in 0002's layout except a rule-level gate reference, which the home leaves **blank** at home:369-371). Flagged: the category may become unreachable. | category reach |

**Fenced verdict for JD-17: CONTRADICTS on fenced text** — 0002:430-441 (layout enumeration and `[accessors.<id>]` shape, 17.1/17.2) and 0002:813-814 (tag-side accessor reference, 17.3). Meaning changes: layout membership, accessor table shape, binding direction, and 0002's accessor-validation scope.

### "Lands in 0002" items (§D7, home:362-366)

| Item | In 0002? | Anchor / note |
| --- | --- | --- |
| root `terminal` and `[initial]` | **absent** (and refused by 17.1's closed layout + strict decode) | no mention anywhere in 0002 |
| the three capability tables `[read.<id>]`/`[write.<id>]`/`[gate.<id>]` | **contradicted** | 0002:431-432, 437-440, 884-885 all fix `[accessors.<id>]` with `mode` |
| `[tags.<tag>].accessor` removed | **contradicted** | 0002:813-814 fence, 0002:350-351, fixture line 11/40 |
| the seven type-model keys enumerated | **absent** | 0002:808-810 names five concepts, no wire keys; fixture uses `kind` only |
| `[model.metadata]` | **absent** | — |
| the write-replaces clause | **absent** | nothing in 0002 states replace vs accumulate; nothing contradicts |
| the two load categories (`malformed tag declaration`, `malformed predicate atom` literal-outside-domain) | **absent / partially present** | `malformed predicate atom` present (0002:833) but defined narrower (0002:371-373); `malformed tag declaration` absent from fence (prose "malformed tags" 0002:370 only) |
| the binding validations (key served by zero/two readers; written or cleared key not in exactly one writer; writer key not owned) | **absent** | 0002:438-440 fence limits 0002's accessor validation to id resolution |
| fixtures re-authored | **absent / contradicted** | rdr-fixture.toml still `[accessors.*]` + `mode` + tag-side `accessor`; 0002:886-888 mandates promoting them verbatim as "canonical examples" |

**Blank the home leaves for 0002** (home:369-371): "where a rule references a gate accessor — 0002's layout has no site for it; owner 0002 with 0004's semantics." Confirmed absent in 0002; no fence contradicts a future site.

---

## Cross-entry observations (evidence, not disposition)

1. All CONTRADICTS on fenced text are in JD-16 and JD-17. JD-15's contradictions are confined to unfenced Round-Trip / mini-check prose (0002:936-938, 949-950, 956).
2. The strict-decoder fence (0002:465-470) turns every JD-17 "absent" layout key into an active refusal under the frozen text — absence is not neutral here.
3. The JD-16 `in`-generalization interacts with a fence the home does not name: expansion-suffix totality (0002:751-753) rests on the suffix being an `outcomes` member; a match-block `in` on a non-`recognized` tag would need a different suffix rule and a different total-order argument.
4. Nothing in 0002's Status qualifier or elsewhere claims the answers have been absorbed; the qualifier only points at the home. RDRs are never amended (user convention), so the gate's question is whether these fenced contradictions route to a re-lock of 0002 or to a successor.
