Model: claude-fable-5

Model stamp: single-model pass (claude-fable-5), same base model as the
iteration-2 critique (claude-fable-5); iteration 1 was claude-opus-5. Per
`2-critique.md` §Re-entry this is a same-model repeat and is recorded as such —
the dual-model signal for this cluster remains iter-1 (opus) vs iter-2 (fable).
Fresh context; no prior answers carried except the ledger the brief mandates.

# Hostile Set Critique — RDR 0002–0009, iteration 3

Assignment: whole-set premortem over the eight Final-and-unimplemented members
(0002 `526c481`, 0003 `526c481`, 0004 `526c481`, 0005 `76a9e67`, 0006
`526c481`, 0007 `526c481`, 0008 `4581223`, 0009 `526c481`) plus the home JDR
0001 (`4f9639c`), read in their 2026-08-24 text. Postmortem files ignored. Delta
scope per the brief: 0008's Stage 4 re-entry and re-lock (2026-08-23), and the
home's §D5/§D6/§D7 answers (2026-08-24) closing JD-15/16/17.

**Verdict up front.** Iteration 2 said the members had locked against answers
the home never gave. The home has now given three of them — and every one of
the three names, in its *Lands in* clause, a substantive edit to a document
that was locked the day before and has not moved. §D6 falsifies a fenced
routing clause in Final 0002. §D7 falsifies Final 0002's fenced layout
enumeration, its fenced type-model clause, its canonical fixtures, and a fenced
diagnostic route in Final 0003. §D5 falsifies 0002's Round-Trip invariant. One
of §D7's own sub-rules contradicts §JD-9 at the same home. And 0008, re-locked
on 2026-08-23 to repair two assumptions verified against a superseded 0002,
carries a third (A10) verified against the `[accessors.<id>]` layout §D7
replaced twenty-four hours later. The cluster has reached the state its own
gate defines as SPEC-DEFECT on its producer, at the last iteration the gate
permits a demotion. Meanwhile `internal/resolve/resolve.go:185` still reads
`Guard string`, no commit has touched `internal/` since 2026-08-09, and the
design corpus is 69,271 lines against 3,971 lines of Go (17.4:1, up from 16:1).

---

## Findings Ledger

Every iteration-2 row is carried with a status and a passage anchor in the
current text; rows raised this pass are `N-*`. Status vocabulary: **open** (the
defect is in current text), **closed** (the home answered *and* the member text
is consistent), **superseded** (the defect mutated — the successor row is
named), **new**.

| ID | RDR | Status | RDR passage (anchor) | Failure mode | Symptom user sees | Origin |
|----|-----|--------|----------------------|--------------|-------------------|--------|
| N-1 | 0002 / JDR | new | 0002:692-717 fenced "The unified predicate set splits across the kernel's two predicate fields **by operator** … atoms whose operator is equality on a declared tag populate `Match`; every other atom — set membership, comparison, existence, and every `unless` atom … belongs to the guard"; 0002 A12 evidence "`status.eq=closed@unless` … routes to the guard despite being an equality operator"; JDR §D6 "**Resolved: (b)** … Every atom under `guard.all` / `guard.unless` is a guard atom regardless of operator, `eq` included … *Lands in **0002** — the routing clause restated as block-keyed*" | The home ratified block-keyed routing on 2026-08-24; Final 0002's fence says operator-keyed. An `eq` atom under `[rule.guard.all]` over an optional key is a `Match` atom under 0002 (silent non-candidate, escapable `no_match` — the masking P1 forbids) and a guard atom under the home (withheld by lint, `guard_unevaluable` at runtime). Fenced; meaning changes | 0003's own fixture `foundational-to-cove` (`profile eq "foundational"` under `guard.all`) lints as withheld and resolves as a silent `no_match` escape; which one the implementer builds depends on which document they opened | §1, §2, AT-8 |
| N-2 | 0002 / 0004 / JDR | new | 0002:430-441 fenced "The source schema MUST use the Resolve spike field layout: root `outcomes`, `[model]`, `[tags.<tag>]`, `[accessors.<id>]`, `[context.<id>]`, `[[rule]]`, and `[dump]` … An `[accessors.<id>]` entry carries `mode` and `path`"; 0002:808-818 fenced "authored — under `[tags.<tag>]`, beside `provenance` and the optional accessor reference"; 0002:883-888 Load-Bearing "The RDR and kata spike fixtures are the canonical examples implementation tests must promote"; JDR §D7(ii) "One table per capability — `[read.<id>]`, `[write.<id>]`, `[gate.<id>]` — replacing `[accessors.<id>]` and its `mode` … the tag-side `accessor` reference is a second copy and is removed" | The home replaces a fenced layout in a Final document and deletes a field a fenced clause names. 0002's "canonical examples" (`[accessors.rdr-status]` with `mode`/`path`; `[tags.status] accessor = "rdr-status"`) are non-conforming under the home. Fenced; meaning changes; supersedes C-8 | The first `internal/model` loader either promotes fixtures the home rejects or invents `[read.x]` tables 0002's strict decoder must refuse as `unknown schema field` | §1, §2, AT-8 |
| N-3 | 0002 / 0006 / JDR | new | 0002:430-441 (same fence) + 0002:443-444 fenced "`[model]` MUST contain `id` and `version`"; 0006 A6 "both declarations are **tag predicates, not state names** … what this RDR requires is that neither is a bare identifier a row references by name"; JDR §D7(i) "root `terminal` is a list of **context ids** … a rule may `use` a terminal context for an escape self-loop" | `[initial]` and root `terminal` have a home answer and no home in any fence (supersedes C-1). The answer's shape — a list of context ids, which rules already reference by name via `use` — is the shape 0006 A6 said it did not want. 0006 is Final; A6 still reads "RDR 0002 is `Draft`" (0006:256, 1421) | `intrastate lint` still fails every table 0002 can load with `graph-dangling-edge` (no root); adding `[initial]` still fails 0002's strict decode | §1, AT-8 |
| N-4 | JDR / 0005 / 0002 | new | JDR §D7(ii) "Readers and writers both declare `keys`; **the loader refuses a key served by zero or two readers**, a written or cleared key not in exactly one writer, and a writer key that is not owned"; JDR §JD-9 "Caller-supplied tags enter as `ProvenanceObserved` … Wire it to `Observed`"; 0005 Technical Design "`--tag name=value` supplies one scalar tag fact"; 0002 canonical fixture `rdr-fixture.toml:13-31` declares `profile`, `stage`, `iter`, `rewind_scope`, `prelock_lens` as owned with **no** accessor, and `recognized` is kernel-injected | Read literally, §D7(ii) refuses every model whose tags are not each served by exactly one reader — which is every caller-supplied observed tag JD-9 sanctions, the kernel-injected `recognized` key, and five of the eight tags in 0002's own canonical fixture. The home contradicts itself and its producer's normative fixture in one sentence | The canonical RDR model is refused at load with a binding error naming `profile`; every `--tag` a skill supplies must also have a declared reader or the model will not load | §1, §3, AT-9 |
| N-5 | 0002 / JDR | new | JDR §D6(b) "Match blocks (`[rule.match.*]`, `[context.*.match.*]`) admit `eq`, and `in` only via expansion into one candidate row per member — the expansion 0002 already defines for `recognized`, generalized"; 0002:626-631 fenced "an `in` atom expands into one candidate row per member, each identified by the rule id plus **the outcome literal** as its expansion suffix"; 0002:751-755 fenced "an expansion suffix is a member of the `outcomes` alphabet, which is non-empty and duplicate-free — so no two distinct rows compare equal"; 0002:597-607 fenced `#` is rejected in rule ids and alphabet members only; 0002 fixture `rdr-fixture.toml:61` `in = ["large", "foundational"]` on a non-`recognized` match key | The home generalizes an expansion whose identity, suffix, totality proof and separator-rejection rule are all fenced in 0002 as *outcome-only*. Two `in` atoms on two match keys produce a cross-product with no defined suffix; a tag literal containing `#` now bleeds into the rendered identity; 0002's "seven deterministic rows" normative fixture (`output.txt`, scenario 2) is wrong under the home because the spike expands only `recognized`. Fenced; meaning changes | Scenario 2's golden fixture fails on first run against a correct implementation, or the implementer keeps the spike's behavior and the home's routing rule is silently not implemented for `in` | §1, §3, AT-10 |
| N-6 | 0003 / 0006 / JDR | new | 0003:1060-1067 fenced "Both surface through the predicate semantic kinds this RDR owns (A4) **onto RDR 0006 findings** and the RDR 0005 envelope"; 0003 disposition rows 1518-1520 "Predicate semantic kind (this RDR) → RDR 0006 finding → RDR 0005 envelope"; JDR §D7(iii) "Two load categories in 0002, rules supplied by 0003 … **0006 mints nothing for either**; 0005 maps them under §JD-8"; 0006 disposition "Model unreadable / not conforming to RDR 0002's schema … refused before normalization; lint never runs; none of this RDR's codes" | 0003's fence routes declaration/literal-domain errors through 0006 findings; 0006 says lint never runs on them; the home sides with 0006. Fenced in 0003; the diagnostic surface changes. Also 0003:958-961 fenced spells `min`/`max` in its own voice while the home says that spelling "becomes a citation to 0002" (restatement — cosmetic) | A `kind = "enum"` with a `min` bound is reported as a lint finding by one implementer and as a load refusal with a 0002 category by another; the skill branching on `graph-lint-failed` never sees it | §1, AT-8 |
| N-7 | 0002 / 0004 / JDR | new | 0002:935-938 Round-Trip "Three sites make the rendered form lossy today … the `<clear>` sentinel is **not reserved** in the tag-value space, so an authored value of `<clear>` renders identically to a cleared tag"; 0002 fidelity row 956 "`<clear>` unreserved in the value space"; 0002:827-839 category list carries no "`<clear>` as authored tag value" entry; JDR §D5 "**Resolved: (a).** … 0002 refuses an authored tag value of `<clear>` at load … 0004 states that a `<clear>` write removes the key, that read-back asserts the key is **absent** … *Lands in **0002** — one load category … and the Round-Trip lossy-site note drops its `<clear>` item; and in **0004** — remove-key / read-back-absent / idempotent-clear / unreadable-on-read clauses plus an MVV scenario*" | The home reserves the sentinel; 0002's Round-Trip/Invariant clause still lists it as unreserved (the brief counts Round-Trip clauses as in scope). 0004 contains the word "clear" exactly once — its Status qualifier — so none of the four clauses or the MVV scenario exists; 0004's fence 329-335 "verify the expected owned-tag values" reads green on a literal `<clear>` string stored by a naive writer, which is the exact hazard §D5 names | A writer stores the string `<clear>`; read-back compares expected `<clear>` to observed `<clear>` and reports success; the tag was never removed and the next resolve sees it present | §1, premortem, AT-8 |
| N-8 | 0008 | new | 0008 A10 evidence (0008:882-918) "the committed fixture shows the shape — `[tags.status]` carries `accessor = \"rdr-status\"` while `[accessors.rdr-status]` carries only `mode` and `path` … `main.go::Accessor` is `{Mode, Path}` with **no key field**"; 0008 Status `Final` (re-locked 2026-08-23, Gate PASS); JDR §D7(ii) (2026-08-24) "`keys` is the binding … the tag-side `accessor` reference is a second copy and is removed" | 0008's Stage 4 re-entry re-verified A9 and A4 against a moving 0002 and re-locked; the next day the home moved the layout A10 cites. A10's *conclusion* (an accessor introduces no key spelling of its own) may survive — under §D7 `keys` names declared tags — but its evidence is the third 0008 assumption verified against a peer artifact that then changed (A6/A11 at iter-1, A9 at iter-2, A10 now). The re-entry was STAGE-SCOPED to A9/A4 and could not have caught it | An implementer reading A10 builds the reserved-key predicate against a `{Mode, Path}` accessor shape that will not exist | §1, AT-8 |
| N-9 | 0007 / 0002 / 0006 | new | 0007:1263-1270 fenced "`Block` is an exported named STRING type carrying **exactly two** constants, `BlockAll = \"all\"` and `BlockUnless = \"unless\"`"; 0002 scenario 2 "each atom reports its authored block (`match`/`all`/`unless`)"; 0006:789-801 "`Block` is RDR 0002 authoring vocabulary (`match`/`all`/`unless`)"; JDR §D6 "match atoms carry block `match` (closing the two-valued/three-valued block vocabulary finding)" | The home closes the vocabulary finding by declaring three values; 0007's fence declares two. The values never collide at runtime only because match atoms are lifted into `Match []Tag` at the handoff — a projection no document states. 0006's `Finding.Block` string and 0007's typed `Block` disagree on the domain they serialize | A lint finding on a match-block atom carries `block: "match"`; a refusal payload can never; a consumer keyed on 0007's two constants rejects the lint value | §3 |
| N-10 | 0006 / JDR | new | JDR §D7(i) "A terminal context matching a non-owned tag is a 0006 blocking finding"; 0006 finding-code table (0006:705-722) closed at ten blocking codes, none for this; 0006 fenced 1004-1011 closes the advisory tier; JDR §D7(iii) "0006 mints nothing for either" | The same decision tells 0006 to mint a code for (i) and to mint nothing for (iii); 0006 is Final and its taxonomy names no terminal-non-owned code; whether it rides `graph-dangling-edge` or a new code is unstated | A terminal context over an observed tag either passes lint silently or trips a code no document names | §3 |
| N-11 | set / JDR | new | JDR §D5/§D6/§D7 *Lands in* clauses naming 0002 (×3), 0003 (×2), 0004 (×2), 0006 (×2); README index: all four `Final`; brief "0002 … Status qualifier only (1 line)"; iter-2 report "N=3 is the last iteration that may demote"; iter-2 SPEC-DEFECT rule as applied to 0008: "Where a member's fences contradict a now-ratified answer, the ordinary path applies" | The mechanism that demoted 0008 at iteration 2 now applies to 0002 on three answers (N-1, N-2/N-3, N-7) and to 0003 on one (N-6). 0002 is `foundational`, so a re-entry re-runs the full lens set, moves its spike fixtures, and re-stales 0008 A10 (N-8), 0006 A6/A10 and 0009's `<clear>` citation again. The gate's iteration cap meets its own demotion rule at the producer | Either the cluster implements 0002 against fences its home has already overruled, or the fourth iteration demotes past the cap | §1, §2, AT-8 |
| N-12 | 0002 / 0003 / JDR | new | JDR §D6 "Scenario 4's unevaluable-sibling witness becomes an `eq` guard atom over an absent optional key"; 0002 MVV "The spike fixtures do not yet carry those siblings … the desk trace records 'no witness possible in the current fixture'"; JDR §D7 "Observation for 0003's next touch: an unmarked `kind = \"enum\"` is a `2^\|domain\|` dimension … every author will read `enum` as single-valued" | The home re-authors 0002's MVV witness by fiat and records, as an "observation" for a document that is Final, the default inversion 0003 A21 calls "a real fork … decided here and carried to RDR 0002 and RDR 0006". A fork the home sees and declines to decide, filed on a next touch that no Final document has | The first table author writes `kind = "enum"` on six tags, gets six `graph-unprovable-coverage / tag-not-single-valued` findings, and the marker key they must add is `single_valued` — a key 0002's strict decoder still rejects (N-2) | §3, premortem |
| C-1 | 0006 / 0002 | superseded → N-3 | 0002:430-441; 0006:1015-1018 fenced "A model that declares no initial owned state MUST be rejected" | Home answered (§D7(i)); 0002's fence unchanged and now contradicts the answer rather than merely lacking it | as N-3 | §1 |
| C-2 | 0006 / 0002 | open | 0006 A10 (0006:370-399) "Flips when RDR 0002 states write-replaces … RDR 0002 is `Draft`"; 0002 fenced 821-824 (clear clause only); JDR §D7(iv) "A clause, not a key: a write assigns a tag's whole value … *Lands in **0002***" | Answered at home; absent from 0002's text; 0006 A10 still Pending and still calls 0002 Draft (0006:380, 391, 1425) | `graph-single-valued-state` never fires on a two-write path | §1 |
| C-3 | 0003 / 0006 | open | 0003 Status (0003:12-15) "A10, A12, A17 and A19 are on RDR 0006's refine, A18 and A20 span RDR 0006's refine and RDR 0007's next touch"; 0003:403 "RDR 0006 is now `Draft`"; 0006 Status "A7 is a route-back on RDR 0003 A18" | Unchanged; six Pending records in Final 0003 still name Final 0006's refine; JD-18 unanswered | 0003 Phase 2 group construction gated on A10 with nothing to wait for | §1 |
| C-4 | 0007 / 0009 | open | 0007 Prereq (0007:2054-2065) "[ ] RDR 0009 re-locks against the loss of `Refusal.Guard`"; 0009 A3 (0009:421) "the asserted properties are `Refusal.Kind` / `Refusal.Guard`"; 0009 Gate PASS 2026-08-11 | 0009 has had no content re-lock; its frozen-suite evidence is against a kernel 0007 deletes | 0009 Phase 1 lands on a `Refusal` with no `Guard` field | §1 |
| C-5 | 0007 / JDR | open | 0007 Prereq (0007:2028-2039) "[ ] JDR 0001 §D4 is reopened and §JD-8 amended … a REVERSAL of that ruling"; JDR §D4 "it reaches the CLI through §JD-8's `Detail`" (unchanged); JDR §JD-8 widened note "0007's Prerequisite asks to reopen §D4" | Recorded at the home as a carrier question, not answered; §D4 still says `Detail` | JSON consumers parse a sorted two-level array out of a prose string | §1 |
| C-6 | 0005 / 0006 / 0007 / 0009 | open | 0005 A5 "not new envelope fields or exit groups"; 0006 fenced 959-978 "`findings` field … MUST NOT be `omitempty`"; 0009 fenced 1030-1055 "a NEW omitempty field … MUST NOT be `resolve.RowRef`"; 0007 A19 "ONE `omitempty` structured field"; JDR §JD-8 "*(blank, except the exit-3 call)* … Widened" | JD-8 widened, still blank; four MUSTs on one shipped struct | Three uncoordinated fields on `CLIError` | §1 |
| C-7 | 0002 / 0003 | superseded → N-6 / N-12 | 0003:1658-1660 "The `single_valued` spelling above is illustrative"; JDR §D7(iii) enumerates seven keys | Answered at home (`kind`/`domain`/`min`/`max`/`elements`/`single_valued`/`required`); 0002's decoder fence still admits none of them | as N-2 | §3 |
| C-8 | 0002 / 0004 | superseded → N-2 / N-4 | 0004 Technical Design "name, capability, artifact role, expected tag keys, timeout policy, read-back"; 0002:437-441 | Answered at home; 0002's fence now contradicts the answer | as N-2 | §3 |
| C-9 | 0005 / 0004 / 0002 / 0007 | open | JDR §D7 "**Blank, not decided here:** where a *rule* references a gate accessor — 0002's layout has no site for it"; 0005 fenced 380-386 "It MAY invoke declared gate accessors … when the matched candidate requires a gate fact" | The home now records the blank rather than deciding it; `[gate.<id>]` tables exist (§D7) with no rule-level binding and no kernel input | `flow next`/`resolve` gate binding is dead surface | §3 |
| C-10 | 0003 / 0007 / 0004 / 0005 | open, widened | 0007 Phase 3 (0007:2157-2168) "Blocked on one RDR 0003 declaration … no stated element encoding"; JDR §D7(iv) "for `set` kind the array literal is the whole new set"; 0005 "duplicate tag names in the same request are refused until the transition-model layer exposes a first-class structured literal for set-valued tags" | The home adds set-valued *writes* while the runtime encoding of a set-valued *value* stays unowned and 0005 cannot author one | `contains` uncontractable; `flow set-state --write labels=…` cannot express a set | §3 |
| C-11 | 0007 / 0005 / JDR | open | 0007 row 16 "`--tag X=v` against `X eq v` ⇒ row selected and a plan"; 0007 fenced "PRESENCE IS PROVENANCE-BLIND"; JDR §JD-9 unchanged | No change | Operator unsticks a refusal from the command line | premortem |
| C-12 | 0008 / 0007 | open — survived the re-lock | 0008 scenario 3 (0008:2390-2413) "through a **new** view-capturing guard seam … The captured guard-side view satisfies `Lookup(\"recognized\") == …`"; 0007 fenced SEAM (0007:1283-1290) "`Evaluate(atom, value) GuardResult` … It never sees the view" | The STAGE-SCOPED re-entry re-verified A9/A4 only; scenario 3 still describes capturing a view at a seam the normative seam is forbidden to receive; the added sentence "the guard-side read of the view key is RDR 0007's domain" concedes the point without changing the test | 0008 Phase 3 cannot deliver scenario 3; A2 stays unpinned | §1 |
| C-13 | 0006 | open | 0006:1068-1076 "two edges reaching the same successor produce one node whose per-tag value sets are the union of theirs" | Unchanged; merge key unstated | `graph-product-too-large` on the first real model | §2 |
| C-14 | 0006 | open | 0006:1363-1368 "if that count is non-zero on a model the maintainers accept, the response is a successor RDR" | Unchanged | Lint blocks the checked-in model; fix is an RDR | §2 |
| C-15 | 0003 / 0006 | open, acknowledged at home | 0003 fenced 1417-1426 defaults; JDR §D7 "Observation for 0003's next touch … every author will read `enum` as single-valued" | The home sees the inversion and files it on a next touch of a Final document (N-12) | Every unmarked guard is a blocking finding | §3 |
| C-16 | 0007 / 0003 | open | 0007 Prereq (0007:2066-2073) "[ ] RDR 0003 states the empty/omitted `unless` identity (A9)"; 0003 fenced 893-896 "minus the single conjunctive assignment set matched by the row's full `unless` block" | Unchanged | Lint gaps on every `unless`-less group, or luck | §1 |
| C-17 | 0003 / 0006 / JDR | closed (home corrected 2026-08-23) | JDR §JD-14 "**Corrected 2026-08-23** … overlap is checked in **two populations**"; 0003:1206-1217, 0006:901-903 fenced two-population | Home now matches both members. Residue, cosmetic: 0003 authority row 1492 still reads "**Contested — A17.** §JD-14 decided for one population" and A17 still books "the §JD-14 correction at the next JDR 0001 touch" — the touch has happened | none | — |
| C-18 | 0003 / 0006 / 0007 | open | 0003 fenced 995-1005 "Conformance currently has no enforcer, and that is a booked gap"; JDR §JD-18 open | Unchanged | Green lint, `guard_unevaluable` at runtime | §3 |
| C-19 | 0004 / 0002 / JDR | closed (home) — residue cosmetic | JDR §JD-3 "**CLOSED 2026-08-23**"; 0002:652-670 fenced producer; 0004:567 still "**Open** … JD-3 records the field appears zero times in 0002"; 0007 Prereq 2022 "[ ] JDR 0001 §JD-3 landed in RDR 0002" | The home and 0002 agree; 0004 and 0007 carry stale citations | none | — |
| C-20 | 0002 / 0008 / JDR | closed | JDR §JD-10 "**ANSWERED 2026-08-23**"; 0008:1400-1418 "Boundary against JDR 0001 §JD-10 (answered 2026-08-23)" | Ratified at home; 0008 re-locked citing the answer | none | — |
| C-21 | 0004 | open | 0004:466-472 "A row that consumes the key through `Row.Match` or a guard string alone is dropped by `TagSet.matches` … yielding `no_match`" | Pre-§D4 control flow still narrated as current in a Final document; unfenced | Scenario 6's guard leg expects `no_match`, kernel refuses `guard_unevaluable` | §1 |
| C-22 | 0004 | open | 0004:769-781 Prerequisites — four boxes, all `[ ]`, incl. "All Critical Assumptions verified"; Gate PASS 2026-08-21 | Verbatim from iteration 1 | Implementer cannot tell whether 0004's preconditions hold | AT-1 |
| C-23 | 0005 | open | 0005:599-614 code table (twelve `flow-*` codes); JDR §JD-8 widened, unanswered | No change; every sibling's failure class still unmapped | Skill branches on codes that do not exist | §1 |
| C-24 | 0007 | open (4 of 15 repaired) | 0007:154 "RDR 0002 … (`Draft`, re-entry at refine)"; 0007:2172 "Hand over, with 0003 in Draft"; 0007:148 "RDR 0003 … (`Draft`, re-entry at refine)" | Iteration 2 repaired four pure status claims; the sentence-entangled remainder stands | Phase 4 handoff addressed to a Draft that does not exist | §1 |
| C-25 | 0008 | open | 0008:2334-2366 "this RDR can reach 'implemented' with its user-facing outcome still unshipped … hold this RDR at Final-unimplemented" | Restated at re-lock; unchanged in substance | A green 0008 row changes nothing an author perceives | §2 |
| C-26 | 0009 | open | 0009 fenced 922-929 "Unwrap() []error MUST be EXACTLY ONE LEVEL deep … Resolve MUST return CheckValid's error VERBATIM" | Unchanged | Rewrite on first real caller | §2 |
| C-27 | 0009 | open | 0009:808-812 "callers MUST check the error before reading the disposition"; scenario 7b | Unchanged | `flow resolve` dereferences a nil `Plan` | §2 |
| C-28 | 0008 / 0009 / JDR | open | JDR §JD-5 open; 0008 fenced 1549-1562 "either error is conforming"; 0008 scenario 6 fourth variant "must NOT assert which of the two errors is reported" | Unchanged; the test JD-5 asks for is forbidden until JD-5 closes | Different error per build on a doubly-breaching table | §1 |
| C-29 | 0002 | open — materialized | 0002:1267-1281 Risks "the producer locks *after* its consumers … **Mitigation**: … `/rdr-cluster-reconcile`" | The delegated mechanism has now returned three answers that contradict 0002's fences (N-1, N-2/N-3, N-7); the risk the RDR named is the event that happened | This critique; N-11 | §1 |
| C-30 | 0006 / 0005 | open | 0006 fenced 974-977 "extending `clierr.EmitText` (failure) and the `respond.OK` text branch"; 0005 Prereq "[ ] Add text success payload rendering through `respond.OK`" | Unchanged | Text rendering implemented twice | §1 |
| C-31 | set | open, worse | `wc -l` docs/rdr+docs/jdr = 69,271; `internal/`+`cmd/` = 3,971 (17.4:1); last commit under `internal/` `1c9c0ca` 2026-08-09; `resolve.go:185` `Guard string`; `root.go:55` registers only `newVersionCmd()` | +4,200 lines of design since iteration 2, zero lines of Go | Binary still answers `version` | §1, premortem |
| C-32 | 0003 | open | 0003:707 "A16 owns the authoring location, still `Pending` on RDR 0002" beside A16 "Status: Verified" | Unchanged | Cosmetic | — |
| C-33 | 0006 | open | 0006:1716-1718 "MUST emit `graph-product-too-large` naming the traversal" beside the same code for a declared product | Unchanged | One code, two remedies | §2 |

Tally over iteration-2's 33 rows: **3 closed** (C-17, C-19, C-20), **3
superseded** (C-1, C-7, C-8 — each answered at the home and each answer now
falsifies a member fence), **27 open** (one widened, one materialized). **12
new rows**, of which N-1, N-2, N-3, N-5, N-6, N-7 are fenced contradictions
between the home and a Final member, N-4 is a home-internal contradiction, and
N-8 is a second-lock stale-verification on 0008.

---

## 1. The inter-RDR failure mode

### The failure

**The set will fail because the home has started answering, and its answers
land nowhere.** Iteration 1's defect was silence (obligations routed to a
peer's implement stage). Iteration 2's was lag (members locked against answers
the home had not given; the cure was to move the home). Iteration 3's is the
inverse of lag: the home moved on 2026-08-24, three entries at once, and each
entry's *Lands in* clause is a substantive edit to a document locked on
2026-08-22/23 whose only change since is a one-line Status qualifier. The
qualifier — `Final [joint decision → JDR 0001 §JD-15, §JD-16, §JD-17]` — was
introduced at iteration 2 as a *tolerance*: a marker that a member accepts a
pending answer. It is now being read as a *license*: that a Final member can
absorb a substantive answer by citation without re-lock. It cannot, because
the answers do not fill silence; they overrule fences.

Read the three decisions against the fences they name:

- **§D6 vs 0002:692-717.** The home: "Every atom under `guard.all` /
  `guard.unless` is a guard atom regardless of operator, `eq` included." The
  fence: "atoms whose operator is equality on a declared tag populate
  `Match`; every other atom — set membership, comparison, existence, and every
  `unless` atom regardless of operator — belongs to the guard." These are the
  two partitions the home itself says "disagree on two input classes." A
  Final document holds one; the home holds the other (N-1).
- **§D7 vs 0002:430-441 and 808-818.** The home replaces `[accessors.<id>]`
  with three capability tables, deletes the tag-side `accessor` field, adds
  `[initial]`, root `terminal`, `[model.metadata]`, and seven type-model keys.
  The fence enumerates a closed layout that has none of them and a strict
  decoder that "MUST reject unmapped keys." The fixtures the fence calls "the
  canonical examples implementation tests must promote" author
  `[accessors.rdr-status]` with `mode`/`path` and `[tags.status] accessor =
  "rdr-status"` — both forms the home retires (N-2, N-3).
- **§D5 vs 0002:935-938.** The home reserves `<clear>`; 0002's Round-Trip
  clause still says "the `<clear>` sentinel is not reserved in the tag-value
  space" and lists it as a lossy site (N-7). 0004, told to gain four clauses
  and a scenario, contains the word "clear" once — in its Status line.
- **§D7(iii) vs 0003:1060-1067.** 0003's fence routes declaration and
  literal-domain errors "onto RDR 0006 findings"; the home says "0006 mints
  nothing for either." 0006 agreed with the home before the home did (its
  disposition table: "lint never runs"). 0003 is the odd one out and is Final
  (N-6).

Under the gate's own rule, applied to 0008 at iteration 2 — "where a member's
fences contradict a now-ratified answer, the ordinary path applies" — 0002 is
in SPEC-DEFECT on three entries and 0003 on one. Iteration 2 also wrote "N=3
is the last iteration that may demote." The gate has arrived at the point
where its demotion rule and its iteration cap meet on the producer of the
cluster (N-11).

### Root cause across the RDRs

**Enabler 1 — the tolerance qualifier was allowed to carry substantive
answers.** A tolerance means "this member's text is consistent with every
answer the entry could take." 0002's fences were not consistent with §D6(b)
or §D7(ii): the home *chose* the arm 0002's fence rejects ("(a)
Operator-keyed (0002 today)" is listed and rejected in §D6). A tolerance on
an entry whose answer space includes "rewrite the tolerating member's fence"
is not a tolerance; it is a deferred demotion. Iteration 2 typed JD-16 and
JD-17 `blocks-impl` on 0002 and still left 0002 Final with a qualifier.

**Enabler 2 — the home writes fixtures and MVVs by fiat.** §D6 says
"Scenario 4's unevaluable-sibling witness becomes an `eq` guard atom over an
absent optional key"; §D7 says "fixtures re-authored"; §D5 says "plus an MVV
scenario for a clearing rule." The home is not a document that owns MVVs — the
RDR README says the RDRs own exact contracts and the JDR's transcripts are
non-normative. When the home starts specifying test witnesses for Final
documents, the members' Testing Strategies become stale on the day the home
moves, and nobody runs a gate on the home.

**Enabler 3 — the home's own rules are not checked against each other.**
§D7(ii)'s binding validation ("the loader refuses a key served by zero or two
readers") is written for owned tags read back through accessors and stated
over all keys. §JD-9 sanctions caller-supplied observed tags with no accessor.
The kernel injects `recognized` with no accessor. 0002's canonical fixture
declares five owned tags with no accessor. One sentence at the home refuses
the home's own prior answer, the kernel's shipped behavior, and the producer's
normative fixture (N-4). Nothing in the gate diffs a new §D against the
existing §JD entries — the same failure as the members' stale peer-status
sentences, one level up.

**Enabler 4 — the 0008 pattern is now a series.** 0008 A6/A11 were verified
against a Draft 0009 that then locked (iter-1 SPEC-DEFECT). 0008 A9 was
verified against 0002's pre-lift spike that then changed (iter-2
SPEC-DEFECT). 0008 A10 was re-verified at the 2026-08-23 re-entry against
`[accessors.<id>] {Mode, Path}` "with no key field," and §D7 replaced that
layout on 2026-08-24 (N-8). Three consecutive locks of one document, each
stamping a peer artifact with no re-verification trigger, each falsified
within days. The STAGE-SCOPED re-entry could not have caught A10 because its
scope was A9/A4 — which is the argument for why STAGE-SCOPED re-entries on a
document whose evidence is mostly peer artifacts do not converge.

### The symptom the user sees

The user is the table author with the first real model. They read 0002 —
Final, the producer, `foundational` — and write `[accessors.rdr-status]` with
`mode` and `path` and `[tags.status] accessor = "rdr-status"`, exactly as its
canonical fixture does. The implementer, having read the home, built
`[read.<id>]` tables with `keys`. Load fails: `unknown schema field:
accessors`. They read the home and rewrite. Load fails: binding error — `profile`
is served by zero readers, because §D7(ii) as written requires every key to
have a reader and `profile` is a `--tag`-supplied dimension (N-4). They add a
reader. Load passes. Lint: `graph-unprovable-coverage / tag-not-single-valued`
on every guard, because `kind = "enum"` is a power set by default (C-15) and
the marker they need is `single_valued`, which the home spelled on 2026-08-24
and 0002's fence still rejects (N-12). They run `flow resolve` on a row with
`profile.eq = "foundational"` under `[rule.guard.all]`. Under 0002's fence it
is a `Match` atom and the row is a silent non-candidate; under the home it is
a guard atom and refuses `guard_unevaluable`. Which they get depends on which
document the implementer opened first (N-1).

That is if any of it exists. `internal/cli/root.go:55` registers
`newVersionCmd()`. The last commit under `internal/` is `1c9c0ca`, 2026-08-09.
The design corpus grew by 4,200 lines between iterations 2 and 3 — every one
of them a decision about code that has not moved in fifteen days (C-31).

---

## 2. The RDR that will be rewritten within six weeks of shipping

**RDR 0002 — Transition Table As Reviewable Data. And not in six weeks; at
this gate, because the gate's own rule requires it.**

Iteration 2 named 0006 and put 0002 as runner-up "because it is the producer
that locked last with every consumer's un-landed request pointed at it, and it
is where C-1, C-2, C-7, C-8 all have to land. It will be reopened first; 0006
will be rewritten." Half of that has happened: the home has now *decided* what
lands in 0002 — §D5 (one load category, one Round-Trip edit), §D6 (the routing
clause, a load category, the generalized expansion, a three-valued block,
A12's MVV, Scenario 4's witness), §D7 (root `terminal` and `[initial]`, three
capability tables, `accessor` removed, seven type-model keys,
`[model.metadata]`, write-replaces, two load categories, the binding
validations, fixtures re-authored). Fourteen named edits, four of them to
fenced text, two to the canonical fixtures, one to the MVV. That is not a
citation; it is a Stage 3 refine and a Stage 4 re-verify of A1, A3, A6, A7
(the spike evidence is `output.txt`, which §D6's generalized expansion and
§D7's layout both change), followed by a Stage 8 re-lock with the full
`foundational` lens set.

**Why 0002 and not 0006 this time.** 0006's defects (C-13, C-14, C-33, N-10)
are internal and unchanged; it can be rewritten *after* its inputs exist.
0002's defects are now *external and ratified*: the home has overruled three
of its fences, and the home is the artifact of record. A document whose fences
the home contradicts cannot be implemented as written — the implementer would
be choosing between the producer's normative text and the cluster's registry.
0002 is also the one member whose rewrite cascades: its fixtures are 0008
A4/A9/A10's evidence, 0006 A6/A10's flip condition, 0003 A14/A16's landing
site, and 0009's `<clear>` citation. Rewriting 0002 re-stales four documents
the day it re-locks — which is precisely why the cluster has avoided doing it,
and precisely why avoiding it does not work.

**What would have to be true for 0002 *not* to be rewritten.** The home would
have to withdraw §D5(a), §D6(b) and §D7(i)–(v) and re-decide each for the arm
0002's fences already hold — operator-keyed routing, `[accessors.<id>]` with
`mode`/`path`, an unreserved sentinel, no `[initial]`. §D6 lists that arm and
rejects it on P1. §D7 lists the sidecar arm and 0006 rejected it. The home
would be choosing the members' text over its own principles to avoid a
re-lock. Nothing in the record suggests that is on offer, and iteration 2's
reasoning ("prefer the clean shape; no callers exist") argues against it.

**The six-week clock is real anyway.** Even after the re-lock, 0002's
generalized `in` expansion (N-5) has no identity rule, its binding validation
(N-4) refuses its own fixture, and its `enum` default (N-12) is a fork the
home has seen and left for "0003's next touch." Those three are the first
three things the implementer meets. The rewrite the gate forces now is not the
last one.

*(Runner-up, for the record: 0003 — one fenced contradiction (N-6), one fence
the home says becomes a citation (`min`/`max`), an authority table that still
says "Contested" on an entry the home corrected (C-17 residue), six Pending
records on a Final peer's refine (C-3), and the `enum` default fork (C-15,
N-12) parked on its "next touch." It is Final at `526c481`.)*

---

## 3. The cross-cutting assumption that will not survive first contact

**The assumption: that a Final RDR can receive a later joint decision by
citation — that a Status qualifier reading `[joint decision → §JD-n]` is a
channel through which the home's eventual answer flows into the member's
contract without the member moving.** Every Final member in the cluster
carries one or more of these qualifiers; the iteration-2 report describes them
as "standing tolerances"; the home's *Lands in* clauses are written as if the
landing is a formality. It is load-bearing everywhere because it is what lets
eight documents stay Final while six joint decisions are open and three have
just closed.

Its instances, now that the home has answered:

- **§D6 → 0002's routing fence.** The tolerance on JD-16 was "which key routes
  an atom." The answer is the arm the fence rejects. There is no citation that
  makes "by operator" mean "by block" (N-1).
- **§D7 → 0002's layout fence.** The tolerance on JD-17 was "what the closed
  layout must additionally spell." The answer removes one field the fence
  names and one table the fence names. A citation cannot delete text (N-2,
  N-3).
- **§D7 → 0003's diagnostic route.** The tolerance on JD-16/JD-17 did not
  cover 0003's fence on where a declaration error surfaces; the answer moved
  it anyway (N-6).
- **§D5 → 0002's Round-Trip invariant and 0004's read-back.** The tolerance on
  JD-15 was "does `<clear>` cross as a value." The answer reserves the sentinel
  in 0002 (whose invariant says it is unreserved) and specifies four clauses
  and a scenario in 0004 (which has none). 0006 "reads a `<clear>` write as
  removal by citation" — citation of a clause in 0002 that does not exist (N-7).
- **§D7 → 0008 A10.** No tolerance at all; 0008 does not carry JD-17. Its
  evidence cites the layout the answer replaced (N-8).
- **§D7 → 0006's finding taxonomy.** 0006 tolerates JD-17 for A6/A10. The
  answer additionally asks it to mint a terminal-non-owned finding its closed
  code table does not name (N-10).

Every instance is the same shape: the tolerance was scoped to the *question*,
the answer reached *text the question did not name*, and the member is Final.
Iteration 2's cross-cutting assumption was that "authoring location vs meaning"
was a complete assignment; §D7 answered that by spelling the keys — and the
spelling has nowhere to go, because the document that owns authoring location
has a fence that admits no new keys and a status that admits no edits.

### How it breaks

**First contact, the implementer of 0002.** They open the producer. Its
fences are unambiguous and closed. They open the home because the Status line
points there. Its *Lands in* clauses are unambiguous and contradict the
fences. There is no third document. They choose — and the project's own rule
("RDRs are never amended; code is the source of truth") means whichever they
choose becomes the contract, with no RDR describing it. Iteration 2's
premortem called this "RDR 0010, seeded from the implementation and locked in
a week." Iteration 3's version is worse: RDR 0010 would now also have to
record which of two *normative* texts it followed.

**First contact, the gate.** Iteration 4, if it runs, finds 0002 in
SPEC-DEFECT on three entries with the demotion cap spent. The gate can (a)
demote past its own cap, (b) record the contradictions as tolerances again and
let 0002 implement against overruled fences, or (c) declare the home
non-normative for the answers it has given. All three are the gate re-deciding
its own rules under pressure — the failure mode 0007 A25 names for tests,
applied to the process.

**First contact, the home.** §D7(ii) refuses 0002's canonical fixture (N-4).
Nobody ran the home's new binding rule against the one model the cluster has.
The home is now a document with its own defects, no Critical Assumptions, no
MVV, no gate, and eight Final members citing it as the artifact of record.

---

## 4. Premortem

*Written from nine months after the third cluster gate.*

We did what the iteration-3 report told us not to: we implemented over it.
The reasoning was that the home had finally answered JD-15/16/17, that 0002's
Status line pointed at the answers, and that a fourth iteration would only
re-demote a producer everyone was tired of re-locking. So `internal/model`
was built "against 0002 as read through the home."

**Weeks 1–4, the loader.** The implementer built the layout from §D7 —
`[read.<id>]`, `[write.<id>]`, `[gate.<id>]`, `[initial]`, root `terminal`,
`[model.metadata]`, seven type-model keys — and the strict decoder from 0002.
0002's two spike fixtures failed to load on the first run (`unknown schema
field: accessors`), so they were re-authored, which meant 0002's A1, A3, A6
and A7 evidence (`output.txt`, SHA `c4be7447…`) was now describing files that
no longer existed. Nobody re-ran the spike; the RDR was Final. Then §D7(ii)'s
binding validation refused the re-authored RDR fixture: `profile` served by
zero readers. The implementer read §JD-9, understood that `--tag` tags have no
reader by design, and scoped the rule to owned keys in code. The home said
"key"; the code said "owned key"; neither document was edited. Two weeks
later the same rule was found to refuse `recognized`, which the kernel
injects, and gained a second exemption.

**Week 5, routing.** The handoff to `resolve.Row` was written block-keyed,
per §D6. 0002's A12 MVV assertion — "the `unless`-block equality atom … routes
to the guard despite being an equality operator" — passed. 0002's fence
("atoms whose operator is equality on a declared tag populate `Match`") was
now false of the shipped normalizer, and its Load-Bearing "Selection /
predicate" decision still described operator routing. §D6's generalized `in`
expansion was then implemented for every match key, and the identity tuple
broke: a rule with `profile.in = ["large","foundational"]` in a context and
`recognized.in = ["a","b"]` in its match block produced four rows whose
suffixes the `(model id, rule id, expansion suffix)` tuple could not
distinguish. The implementer minted `#profile=large#recognized=a`. 0002's
fence rejects `#` in rule ids and alphabet members only, so a profile value
`large#x` collided. The dump's "total ordering" clause was quietly no longer
total (N-5).

**Week 7, `<clear>`.** 0004's implementer, whose document contains the word
"clear" once, built read-back as "verify the expected owned-tag values are
present." The 0002 implementer emitted `<clear>` writes. The first clearing
rule in the RDR flow (`finalized_at` on rewind) stored the string `<clear>`
in the artifact; read-back compared expected `<clear>` to observed `<clear>`
and reported success. Three weeks later a guard `finalized_at exists = false`
decided FALSE on a tag that was "cleared," and the rewind edge never fired.
The fix was in 0004, which was Final, so it landed as a code comment citing
§D5 (N-7).

**Week 9, lint.** 0006 was built. A6 flipped "by citation" to a 0002 clause
that existed only in the home. The reachability root was read from
`[initial]`; `terminal` was a list of context ids, so invariant 1's "as tag
predicates, not as state names" was satisfied by resolving the id — and a
terminal context over an observed tag was accepted silently, because no
finding code existed for it (N-10). 0003's implementer, meanwhile, had
emitted `malformed tag declaration` as a lint finding per 0003's fence; 0006
never saw it because lint "never runs" on a model 0002 refused. The two
implementations were reconciled in a PR comment.

**Week 11, 0008 again.** 0008 Phase 2 was picked up. A10's evidence described
`{Mode, Path}` accessors and a tag-side `accessor` reference; the shipped
loader had neither. The reserved-key predicate was written against the
`keys` binding instead. It was 0008's fourth divergence between a Verified
assumption and the peer artifact it cited, and the first one that shipped
without a gate noticing, because there was no fourth gate.

**Week 14, the first table author.** The RDR flow model. Load: three rounds
of binding errors. Lint: nineteen `tag-not-single-valued` findings, fixed by
`single_valued = true` — a key the author found in the home, not in 0002.
Resolve: `profile.eq = "foundational"` under `[rule.guard.all]` over a kata
whose `.kata.toml` did not carry `profile` — `guard_unevaluable`, as the home
intended and 0002's fence forbade. The author read 0002's fence, filed a bug
saying the row should have been a `Match` non-candidate, and was pointed at
§D6. The runbook gained an entry: "when 0002 and JDR 0001 disagree, the JDR
wins, except §D7(ii), where the code wins."

**What actually killed us.** Not the decisions — §D5, §D6, §D7 are each
defensible, and §D6(b) in particular is right. What killed us was that we let
the home decide substantive edits to Final documents and called the result
"landed." A tolerance qualifier was our way of not re-locking; when the answer
came, the qualifier had no mechanism behind it, so the answer lived in the
registry and the fence lived in the RDR and the code lived in neither. We had
a rule for this — a member whose fences contradict a ratified answer demotes —
and an iteration cap that made the rule unusable at exactly the moment it
applied to the producer. We spent iteration 3 confirming 0008's re-lock and
did not run 0002 through the same check, because 0002 had not "moved." The
answers had moved. That was the check.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

AT-1 through AT-7 from iteration 2 stand unchanged and still fail (C-3, C-22,
C-24 on AT-1; C-4, C-12, C-16 on AT-2; C-5, C-6, C-23, C-28 on AT-3; C-9,
C-10, C-15 on AT-4; C-11 on AT-5; C-13, C-14, C-18 on AT-6; C-25, C-29, C-31
on AT-7). Three are added for the defects this iteration introduced.

### AT-8 — A home answer whose *Lands in* names a Final member demotes that member in the same commit
*Catches N-1, N-2, N-3, N-6, N-7, N-8, N-11; would have caught C-1, C-7, C-8 as they closed*

```gherkin
Given a JDR entry moves from open to decided
And its "Lands in" clause names an RDR whose README status is Final
When the deciding commit is made
Then for each named RDR the gate MUST diff the answer against every
  ```normative fence, Round-Trip clause, MVV, and Critical Assumption
  Evidence the answer touches
And if any fenced or invariant clause is falsified, the same commit MUST
  flip that RDR to Draft with the qualifier naming the entry and the
  assumptions to re-verify — a Status qualifier citing the entry is NOT
  a landing
And if the RDR is foundational, the demotion is FULL-FLOW, not
  STAGE-SCOPED, because its spike evidence is what the answer moves
And the iteration cap MUST NOT prevent this flip: a cap that forbids the
  demotion its own rule requires is a cap on the gate, not on the defect
```
Run today: §D5 → 0002 (Round-Trip 935-938), §D6 → 0002 (692-717), §D7 →
0002 (430-441, 808-818, fixtures) and 0003 (1060-1067) → **0002 and 0003
Draft**; 0004 and 0006 carry absent landings (no fence falsified) → tolerance
with a named re-lock trigger; 0008 A10 → re-verify.

### AT-9 — A load-refusal rule decided at the home must be run against every input channel and the producer's normative fixture
*Catches N-4, N-10*

```gherkin
Given a home decision states a loader MUST refuse some input class
When the entry is marked decided
Then the rule MUST be evaluated against: every tag provenance (owned,
  observed, recognized), every channel that supplies a tag (accessor read,
  CLI --tag per §JD-9, kernel injection per §D1/0008), and every fixture
  a Final member calls canonical
And any fixture or sanctioned channel the rule refuses MUST be recorded in
  the entry as a scoping of the rule or a re-authoring of the fixture —
  never left for the implementer to exempt in code
And every "is a 0006 blocking finding" the home mints MUST name a code in
  0006's table or record that 0006 must mint one
```
Run today: §D7(ii) refuses `profile`, `stage`, `iter`, `rewind_scope`,
`prelock_lens` in `rdr-fixture.toml`, every `--tag`-supplied observed tag, and
`recognized` → one failure. §D7(i)'s terminal-non-owned finding names no code
→ one failure.

### AT-10 — A generalized rule must restate every identity, ordering, and separator clause the special case relied on
*Catches N-5, N-9*

```gherkin
Given a home decision generalizes a mechanism a Final member defines for
  one case (in-expansion on `recognized` → on every match key; block ∈
  {all, unless} → {match, all, unless})
When the entry is marked decided
Then the gate MUST list every fenced clause that quantifies over the
  special case — the identity tuple, the suffix definition, the totality
  argument, the separator-rejection rule, the exported constant set — and
  state for each whether it survives the generalization or what replaces it
And a normative fixture the member names (output.txt) MUST be re-run or
  recorded as superseded
```
Run today: §D6's generalized expansion leaves the suffix defined as "the
outcome literal," totality proven by "member of the alphabet," `#` rejected in
two of three positions, and `output.txt` describing seven rows the rule no
longer produces → failure. §D6's `match` block leaves 0007's two-constant
`Block` fence unreconciled → failure.

---

## 6. Delta vs iteration 2

Of iteration 2's 33 rows: **3 closed** — every one closed at the home on
2026-08-23 (JD-3, JD-10, JD-14 correction) with member text already
consistent; **3 superseded** — C-1, C-7, C-8 were answered by §D7 and each
answer now contradicts the fence the row named, so they reappear as N-2/N-3
and N-6; **27 open** — everything routed to JD-5, JD-8, JD-9, JD-18 (all
unanswered), everything internal to 0005, 0006 and 0009 (unmoved), and 0008's
scenario 3 (C-12), which the STAGE-SCOPED re-entry did not reach.

The twelve new rows have one origin: **the home answered and the members did
not move.** Seven are contradictions between a 2026-08-24 home decision and a
2026-08-22/23 fence (N-1, N-2, N-3, N-5, N-6, N-7, N-9); one is the home
contradicting itself (N-4); one is the home telling a Final member to mint a
code it does not name (N-10); one is 0008's third stale verification (N-8);
one is the home authoring MVVs and parking a fork on a next touch (N-12); and
one is the process finding that the gate's demotion rule and its iteration cap
now meet on the cluster's producer (N-11).

What did *not* change: no `Detail`/structured-field answer (C-5, C-6), no code
table (C-23), no precedence (C-28), no conforming-view enforcer (C-18), no gate
binding (C-9, now explicitly blank at the home), no set-value encoding (C-10,
now with set writes on top), no `--tag` guard-path answer (C-11), and no Go
(C-31: 17.4:1).

---

*Ledger rows: 45 (33 carried, 12 new). Acceptance tests: 3 new, 7 carried.
Sections 1–5 are the argument; the ledger is the routable record. Any defect
not carried in the ledger is unreported by construction.*
