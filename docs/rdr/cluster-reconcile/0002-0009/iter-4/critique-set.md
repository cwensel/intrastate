Model: claude-fable-5

Model stamp: single-model pass (claude-fable-5), same base model as the
iteration-2 and iteration-3 critiques; iteration 1 was claude-opus-5. Per
`2-critique.md` §Re-entry this is a same-model repeat and is recorded as such —
the dual-model signal for this cluster remains iter-1 (opus) vs iter-2 (fable).
Fresh context; no prior answers carried except the ledger the brief mandates.

# Hostile Set Critique — RDR 0002–0009, iteration 4

Assignment: whole-set premortem over the eight Final-and-unimplemented members
(0002 `365bd66`, 0003 `31bb9e4`, 0004 `ffaf897`, 0005 `76a9e67`, 0006 `567600f`,
0007 `526c481`, 0008 `567600f`, 0009 `567600f`) plus the home JDR 0001
(`567600f`), read in their 2026-08-24 text. Postmortem files ignored. Delta scope
per the brief: the re-locks of 0002 (STAGE-SCOPED), 0003 and 0004 (RE-LOCK-ONLY)
against §D5/§D6/§D7; the home gave no new answer (JD-5, JD-8, JD-9, JD-18 open).
N = 4; the demotion cap is spent, so every row below is evidence, not a
disposition.

**Verdict up front.** The three re-locks did what iteration 3 asked: 0002's
routing, layout, expansion, sentinel, provenance-scoped binding and MVV witness
now match the home; 0003 routes its two rejections onto 0002's categories; 0004
carries clear-as-removal and the capability-table definition. Seven of
iteration 3's twelve new rows close on the text. What the re-locks did not do —
could not do, under STAGE-SCOPED and RE-LOCK-ONLY — is look sideways. Every
contradiction left in the set is now **between two members where the home was
silent or blank**, and three of them are fenced on both sides: 0002 puts a
rule-level `gate` list on the row and requires the executor to honour it before
applying the plan, while 0005's fences put gate evaluation before the kernel
call and forbid `set-state` from running anything but writers, and 0004 — the
owner of gate semantics — never mentions a rule's gate list at all; 0005's
`flow resolve` MUST NOT run read accessors and takes state as `--tag`, JD-9
wires `--tag` to `Observed`, and 0002 derives a non-empty `RequiresOwned` for
every ordinary rule, so `flow resolve` as fenced refuses `owned_state_unavailable`
on every write-bearing row of the canonical model; 0002's normalized set-valued
write "MUST NOT be joined into a single delimiter-separated string, in the
stored value" while the only kernel carrier is `Tag.Value string`. And the
producer's own canonical fixture — re-authored on 2026-08-24 under §D7 — leaves
`iter` (int) and `cluster_ready` (bool) without `single_valued`, so under
0003's fenced projection clause and 0006's fenced citation of it the first lint
of the first model emits two blocking `graph-unprovable-coverage` findings and
0002's scenario 4 overlap oracle is never reached. Meanwhile `internal/` has
not changed since 2026-08-09 (`1c9c0ca`), `resolve.go:185` still reads
`Guard string`, and the design corpus is 75,670 lines against 3,971 lines of Go
(19.1:1, up from 17.4:1).

---

## Findings Ledger

Every iteration-3 row is carried with a status and a passage anchor in the
current text; rows raised this pass are `Q-*`. Status vocabulary: **open** (the
defect is in current text), **closed** (the text now agrees, or the home
answered and the member text is consistent), **superseded** (the defect mutated
— the successor row is named), **deferred** (written to a `deviations.md` at
iteration 3 with a mechanical check; carried, not re-litigated), **new**.

| ID | RDR | Status | RDR passage (anchor) | Failure mode | Symptom user sees | Origin |
|----|-----|--------|----------------------|--------------|-------------------|--------|
| Q-1 | 0002 (fixture) / 0003 / 0006 | new | 0002 `evidence/spikes/iter-2/rdr-fixture.toml:40-47` `[tags.iter] kind = "int" min = 0 max = 9 required = true` (no `single_valued`), `:59-62` `[tags.cluster_ready] kind = "bool"` (no marker); `:109-110` `[rule.guard.all.iter] lt = 3`, `:128-129` `[rule.guard.all.cluster_ready] eq = true`; 0003:948-952 fenced "a value atom over a dimension the model does **not** declare single-valued … is one lint **cannot project at all** … MUST take the blocking inability-to-prove outcome"; 0003:1273-1290 fenced assignment-count table (`bool` without the marker = `2^\|domain\|`; "An implementation MUST apply these defaults"); 0006:945-947 fenced "An atom over a tag not declared single-valued has no projection and MUST take `graph-unprovable-coverage` (`0003::A21`)"; 0002:1287-1290 fenced "those fixtures are what implementation tests promote"; 0002:2269-2270 scenario 4 "run RDR 0006's lint over a deliberately overlapping variant … reported by lint as an ambiguous overlap"; 0006:1363-1367 "if that count is non-zero on a model the maintainers accept, the response is a successor RDR" | The producer's canonical model, re-authored the day §D7 landed, fails the consumer's blocking lint on two guard dimensions — one of them a `bool` that the defaults read as able to hold `true` and `false` at once. The `round-clean` group (`continue-prelock` ‖ `continue-prelock-cluster`) has no projection for `iter.lt` or `cluster_ready.eq`, so no coverage union, so (0003 trace step 7) no overlap verdict: 0002's scenario 4 overlap oracle is unreachable on the fixture it names, and 0006's scenario-23 "successor RDR" trigger fires on the first model. 0003 A21 ("requiring `single_valued` leaves the target flows authorable") is refuted by the one author who had the clause open while typing | `intrastate lint` on the RDR model: two `graph-unprovable-coverage` findings naming `iter` and `cluster_ready`; the author adds `single_valued = true` to a `bool` and asks why | §1, §3, premortem, AT-11 |
| Q-2 | 0002 / 0004 / 0005 | new — supersedes C-9 | 0002:709-714 fenced "A rule's `gate` list names the gate accessors whose `allow` the executor requires before applying that rule's plan; the list is carried on the normalized row and is part of its value. When it runs, and what `deny` and `indeterminate` mean, are RDR 0004's … This fills the site §D7 left blank"; 0004 — zero occurrences of "gate list" / rule-level gate; 0004:264-267 "runtime invocation applies read and gate accessors … and classifies … gate denied, and gate-indeterminate outcomes" (no site, no ordering); 0005:380-386 fenced "flow resolve … MAY invoke declared gate accessors … **before the pure kernel call** when the matched candidate requires a gate fact"; 0005:393-397 fenced "flow set-state MUST invoke only declared write accessors"; 0002:1146-1166 kernel row fields (no gate field on `resolve.Row`; 0007's reshape adds none); JDR §D7 "**Blank, not decided here:** where a *rule* references a gate accessor … owner 0002 with 0004's semantics" | Three placements of one check. 0002: before applying the plan (which is `set-state`, whose fence admits only writers). 0005: before the kernel call in `resolve`/`next` — before the matched candidate exists. 0004: nowhere. The gate list lives on 0002's normalized row, which is not the kernel row and is not the `Plan`, so nothing carries it from `resolve` to `set-state` across two process invocations. 0002 declares the blank filled; the semantics owner has not landed a clause. Fenced on both 0002 and 0005; meaning-level | `[gate.rdr-lock]` is declared, validated, carried on the dump, and never executed; or `set-state` runs it in violation of 0005's fence; or `resolve` runs every gate of every candidate before knowing which matched | §1, §3, premortem, AT-13 |
| Q-3 | 0005 / 0002 / JDR §JD-9 | new — widens C-11 | 0005:380-383 fenced "flow resolve MUST return exactly one resolved next tag-set or exactly one CLIError refusal … It MUST NOT discover artifacts, run read/write accessors"; 0005:298-302 `--artifact` accepted by `read-state` (readers), `next`/`resolve` (gates only), `set-state` (writers); 0005:317 "`--tag name=value` supplies one scalar tag fact"; JDR §JD-9 "Caller-supplied tags enter as `ProvenanceObserved` and never satisfy an owned-state dependency … Wire it to `Observed`"; 0002:830-833 fenced "An ordinary transition rule MUST contain a write block"; 0002:1125-1130 fenced "`Row.RequiresOwned` … the sorted, duplicate-free set of tag keys named by the rule's write block and clear list"; 0002:684-691 fenced "the kernel evaluates it against the view (`resolve.go::missingOwned`), so an owned tag written on one transition must be readable on the next or that resolve refuses `owned_state_unavailable`"; `adversarial_test.go::TestAdv2b_EscapeEdgeMustNotBypassTheOwnedStateRequirement` | `flow resolve` receives owned state only through `--tag`; JD-9 makes every `--tag` observed; every ordinary rule derives a non-empty `RequiresOwned`; `missingOwned` checks that set against **owned** provenance. So the verb, as fenced, refuses `owned_state_unavailable` for every non-escape row of every model 0002 can load. `read-state` reads the owned snapshot in a different process and its output is not `Resolve` input. 0005 has not moved since June; its §JD-9 tolerance was on the *classification* arm, and the *provenance* arm the home wrote in contradicts the verb design in meaning. Fenced on 0005 and 0002; the home is the contradicting party | `intrastate flow resolve --flow rdr --tag status=Draft --tag stage=prelock … --outcome round-clean` → `flow-fact-missing`/`owned_state_unavailable` naming `stage`, `iter`, every time, for every rule; the only verb that ever returns a plan is the one nobody can call | §1, §3, premortem, AT-12 |
| Q-4 | 0002 / 0007 | new — supersedes N-9 | 0002:928-934 fenced "the widening to a third `match` member is **this RDR's**, on JDR 0001 §D6's authority … An implementer building the atom type from RDR 0007 alone gets two members and MUST widen to three here"; 0007:1263-1267 fenced "`Block` is an exported named STRING type carrying exactly two constants, `BlockAll = \"all\"` and `BlockUnless = \"unless\"`"; 0007:1266-1268 "The payload clause below requires the payload's `Block` to be 'the same exported constant type the atom carries'"; 0002:1170-1178 fenced "The split is applied when a `resolve.Row` is constructed for the kernel"; 0002:2077-2079 "Implementation MUST NOT open Phase 2 by mirroring the constants locally" | Two fences fix the cardinality of one exported constant set at two values. Either 0002 widens `resolve.Block` (0007's "exactly two" is false, and 0007's payload can carry `match`, which no refusal can name) or 0002 mints its own three-valued block type and the handoff **maps** between two `Block` types — the reconstruction step §D1 exists to remove — while 0002 forbids local mirrors of kernel constants. Neither document says which. Fenced both; the Go type identity is the meaning | `go build` after 0007 Phase 1 and 0002 Phase 2 either fails on a duplicate constant or compiles with a `blockOf(b model.Block) resolve.Block` switch nobody specified and one test cannot reach the `match` arm | §3, AT-13 |
| Q-5 | 0002 / 0003 / 0007 | new — supersedes C-10 | 0002:843-853 fenced "a `set`-kind write value MUST normalize to an ordered, member-sorted sequence and MUST NOT be joined into a single delimiter-separated string, **in the stored value** or in any comparison derived from it"; 0002:1146-1150 fenced "Normalization MUST populate … the writes with the owned-tag writes the accessor layer applies"; `internal/resolve/resolve.go` `Tag{Key, Value string}`, `Row.Writes []Tag`; 0007:2157-2168 "`resolve.Tag.Value` is a bare `string`, so a set-valued tag reaches the seam as one opaque string with no stated element encoding … element representation is part of RDR 0003's typed operator semantics"; 0003 — states no encoding (its set-literal clause 0003:1240-1249 fixes *canonical spelling as a set*, not a byte encoding); 0004:367-370 fenced read-back "compares each planned owned tag for equality against the held value"; 0005:317-319 "duplicate tag names … refused until the transition-model layer exposes a first-class structured literal for set-valued tags" | 0002's normalized write value is a `[]string` that must reach `Row.Writes []Tag` whose `Value` is a `string`. The only ways are to join (which 0002's fence forbids "in the stored value"), to change `Tag` (RDR 0001 is Implemented and 0007's reshape leaves `Tag` alone), or to encode per an element encoding 0007 assigns to 0003 and 0003 does not state. 0004 then compares that string for equality; 0005 cannot author one at all. Fenced in 0002; the carrier is a compile-time fact | The `labels` set write in the kata fixture either never reaches `Row.Writes`, or reaches it as `"x,y"`, and read-back on the artifact's real list value reports `read_back_mismatch` on every set write | §3, AT-13 |
| Q-6 | set / JDR §JD-8 | new | 0002:2057-2081 Prereq "[ ] … Phases 2 and 3 therefore sequence behind the reshape in their entirety … **The reshape has no owner** … RDR 0007 is `Final` and therefore will not itself schedule the work, no kata tracks it, and no phase in this RDR's plan claims it"; 0007:2117-2127 "### Phase 1: Row and seam reshape — Replace `Row.Guard string` with the §D1 atom slice"; 0007:2020-2090 seven `[ ]` prerequisites incl. 2028 "JDR 0001 §D4 is reopened and §JD-8 amended", 2054 "RDR 0009 re-locks against the loss of `Refusal.Guard`", 2066 "RDR 0003 states the empty/omitted `unless` identity", 2040 "Two duties are added to RDR 0002's Refinement Context" (a section that no longer exists); 0003:1914-1920 "[ ] … Phase 1 cannot begin … until that reshape lands"; 0003:1959-1964 Phase 4 "Blocked on … §JD-18"; 0008:2360-2365 "hold this RDR at Final-unimplemented until 0002's work starts and the JDs this RDR depends on close" (JD-5/8/9); 0009 Status §JD-5/§JD-8; 0005 Status §JD-8/§JD-9, 0005:631 "[ ]"; 0006:1415-1425 "[ ] A8 … No such model exists today", A6/A10 deferred on 0002; JDR §JD-8 "Widened 2026-08-23 … 0005 is the sibling this answer re-walks" — no answer since | The implementation order is a directed graph of unchecked boxes with a cycle-adjacent head: 0002 → 0007 Phase 1 → (0009 re-lock, JD-8 reopen of §D4, 0003 `unless` identity) → JD-8 → 0005 re-walk → nobody. 0002 says the reshape has no owner; 0007's Phase 1 *is* the reshape and 0007 says it is Final before 0003 implements. Each document has correctly named its dependency; none is on the critical path's head, which is an unanswered interface-record entry with no sibling scheduled to answer it. Unfenced; no clause's meaning changes; nothing mechanical decides it | Fifteen days without a commit under `internal/` become ninety; the first Stage 8 prompt on any member opens with "first, decide JD-8" | §1, premortem, AT-14 |
| Q-7 | 0002 | new — residue of N-7 / C-24 class | 0002:1053-1058 fenced "The accessor-side half is **owed, not settled**: JDR 0001 §D5 (§JD-15) makes read-back assert *absence* for a `<clear>` write, but RDR 0004 does not yet carry that clause — its read-back prose asserts presence with 'lossy exemptions: none', and its own Status line records the gap … a clearing rule is not end-to-end until RDR 0004 lands §D5. Do not read this clause as a claim that it already has"; 0004:358-364 fenced (landed 2026-08-24); 0002:1229-1233 fenced "Re-checked at Stage 6: RDR 0003 is `Draft [revised from Final 2026-08-24]`, so this pins a peer that is not currently Final"; README rows 0003/0004 `Final`; 0002:2073-2076 "The reshape has no owner … no phase in this RDR's plan claims it" vs 0007 Phase 1 | The producer re-locked hours before its two RE-LOCK-ONLY siblings and fenced their pre-re-lock state as fact. An implementer reading 0002's reserved-value clause is told, inside a normative fence, that the clearing contract is not end-to-end and that 0003 is Draft. No clause meaning changes; both are grep-decidable against the README | The 0002 implementer files "0004 has no clear semantics" against a document that has nine lines of them | §1 |
| Q-8 | JDR / 0004 / 0002 | new — supersedes N-4 | JDR §D7(ii) "Each decodes to its own shape, so `keys` exists only on readers and `read_back` only on writers … Readers and writers both declare `keys`; the loader refuses a key served by zero or two readers"; JDR §JD-17 "**Open at (ii)** … the binding validation needs a provenance scope … owner 0002 at its re-entry"; 0002:676-707 fenced "the **binding validations**, which are provenance-scoped … every **owned** tag MUST be served by exactly one reader; an **observed** tag MAY be served by at most one reader, and zero is legal … This settles the provenance scope §D7(ii) left open at this RDR"; 0002:663-664 fenced every entry — gate entries included — carries `keys`; fixture `[gate.rdr-lock] keys = ["status"]`; 0004:254-257 "RDR 0002 owns the carrier and the binding validations (a key served by zero or two readers, a written or cleared key not in exactly one writer)" | 0002 settled the scope in a fence; the home's §D7(ii) sentence still reads unscoped and contradicts itself twice in one paragraph (`keys` only on readers / readers and writers / 0002 puts it on gates); 0004 restates the home's unscoped mechanism instead of citing 0002's scoped one. What a gate entry's `keys` *binds* is stated nowhere (0002 requires it non-empty; 0004 defines `keys` for readers and writers only). Cosmetic at the home and 0004 (unfenced restatement); the gate-`keys` gap folds into Q-2 | A `--tag`-only observed key loads under 0002 and is refused by an implementer reading 0004's sentence; `[gate.x] keys` is validated for declaration and used for nothing | §3 |
| Q-9 | 0004 / 0008 | new | 0004:273-278 "the executor records the same caller-supplied artifact role's observed and recognized tag values … verifies the pre-write observed and recognized tag values are unchanged"; 0004:371 fenced "and that observed and recognized tag values present before the write are unchanged"; 0004:475-480 Round-Trip "any observed or recognized tag values present before the write"; 0002:678-680 fenced "every key in a `keys` list … MUST NOT be `recognized` — the kernel supplies that key and no accessor reads or writes it (RDR 0008)"; `grep -c 'RDR 0008' 0004-*.md` → 0 | The recognized arm of 0004's non-owned read-back comparison names a tag no accessor can read: `recognized` is injected by `assemble` from `Input.Recognized`, never held in an artifact role. Vacuous rather than contradictory — the compared set is empty by construction — but 0004 re-locked twice without citing the RDR that made it so. Fenced (0004:371); no meaning change | none at runtime; a test author writes a "recognized tag changed under the write" fixture that cannot be constructed | §3 |
| N-1 | 0002 / JDR | closed | 0002:1168-1192 fenced "splits … by authored block (JDR 0001 §D6) … every atom carrying block `all` or `unless` belongs to the guard predicate, **regardless of operator** — an `eq` atom under `guard.all` is a guard atom"; 0002:1508-1513 Load-Bearing "Routing key … by its authored block, never by its operator" | Landed as the home decided | none | — |
| N-2 | 0002 / 0004 / JDR | closed | 0002:606-621 fenced closed layout `[read.<id>]`, `[write.<id>]`, `[gate.<id>]`; 0002:1488-1492 "Rejected: a tag-side `accessor` reference … rejected: `[accessors.<id>]` with a `mode` field"; 0002 A1 "the pre-§D7 fixtures are now themselves refused (`unknown schema field: accessors.*, tags.*.accessor`)"; 0004:249-258, 415-418 by citation | Landed; fixtures re-authored; strict decoder refuses the old shape | none | — |
| N-3 | 0002 / 0006 / JDR | closed (0006 residue deferred → `0006/artifacts/deviations.md` D1) | 0002:716-748 fenced `[initial]` owned assignments; root `terminal` context ids "normalization MUST dereference each to the context's explicit predicate set over owned tags … never the bare ids. This is what satisfies RDR 0006 A6's shape requirement"; 0006 A6/A10 still `Pending`, "RDR 0002 is `Draft`" ×7 (D2) | Landed in 0002 in the shape 0006 A6 asked for; 0006's flip-by-citation is the deferred check | none | — |
| N-4 | JDR / 0005 / 0002 | superseded → Q-8 | 0002:676-707 settles; home §D7(ii) literal text unchanged | 0002 answered the home's open note; the home and 0004 still carry the unscoped sentence | as Q-8 | — |
| N-5 | 0002 / JDR | closed | 0002:1103-1122 fenced "Match-block expansion is general (JDR 0001 §D6) … The **expansion suffix** is the sequence of chosen members, one per atom with more than one member, taken in the atoms' sort order"; 0002:1037-1047 `#` refused in rule ids, alphabet members, **and match-block `in` members**; 0002:1317-1325 sequence comparison; A7 re-verified on `continue-prelock#large` / `#foundational`; `output.txt` nine RDR rows | Identity, suffix, totality, separator and the normative fixture all restated for the generalized rule | none | — |
| N-6 | 0003 / 0006 / JDR | closed | 0003:914-925 fenced "Both are this RDR's rejection rules carried by RDR 0002's load categories (JDR 0001 §D7(iii)) … RDR 0006 mints nothing for either; RDR 0005 maps both onto the envelope"; 0003:806-812 `min`/`max` "notation, not wire spelling: RDR 0002 spells the authored form"; 0003:1387-1389 disposition rows re-routed | Landed as §D7(iii) states | none | — |
| N-7 | 0002 / 0004 / JDR | closed (0002 fenced residue → Q-7) | 0002:1049-1061 fenced "The `<clear>` sentinel is reserved (JDR 0001 §D5) … refused at load … as a `reserved tag value`"; 0002:1553-1554 Round-Trip "(The `<clear>` sentinel is reserved … so those two renderings are unambiguous)"; 0004:358-364 fenced remove-key / read-back-absent / idempotent / unreadable-on-read; 0004 Scenario 9; A11 | Both landings exist; 0002's fence still says 0004's does not | as Q-7 | — |
| N-8 | 0008 | deferred (`0008/artifacts/deviations.md` D2) | 0008:882-918 A10 evidence on `[accessors.<id>] {Mode, Path}`; D2 check: grep 0002's re-locked layout for a `keys` entry naming `recognized` → 0002:678-680 fenced forbids it, so A10's conclusion survives; the evidence text is stale | Check now runnable and passes on the conclusion; the stale evidence prose remains until 0008's next touch | none | — |
| N-9 | 0007 / 0002 / 0006 | superseded → Q-4 | 0002:928-934 now asserts the widening explicitly and assigns it to itself | The disagreement is no longer implicit; it is two fences | as Q-4 | — |
| N-10 | 0006 / JDR | deferred (`0006/artifacts/deviations.md` D1 check 2) | 0002:745-748 fenced "whether a terminal context matches a non-owned tag, are RDR 0006's blocking findings, not load failures"; 0006:705-722 code table still names none | 0002 now routes the case to 0006 by fence; 0006 has no code; the deferred check owns it | as iter-3 | — |
| N-11 | set / JDR | closed | 0002/0003/0004 re-locked 2026-08-24 (`365bd66`, `31bb9e4`, `ffaf897`) within the cap | The demotion the rule required happened at N=3; N=4 has nothing to demote and may not | none | — |
| N-12 | 0002 / 0003 / JDR | superseded → Q-1 (fixture) and C-15 (fork) | 0002:2269-2270 scenario 4 now carries the home-authored witness (`cluster_ready.eq=true@all` over an optional key); the `enum`/`bool` default observation remains an "Observation for 0003's next touch" — 0003 was touched 2026-08-24 and did not take it | The MVV half closed; the fork half materialized in 0002's fixture | as Q-1 | — |
| C-1 | 0006 / 0002 | closed | 0002:716-748; 0006:1015-1018 unchanged and consistent (lint rejects a rootless model; 0002 loads it and leaves it to lint, 0002:745-748) | Landed | none | — |
| C-2 | 0006 / 0002 | closed (0006 residue deferred D1) | 0002:838-841 fenced "**A write replaces.** A write assigns a tag's whole value and supplants whatever was held; for a `set` kind the array literal is the whole new set. There is no accumulate form"; 0006 A10 still `Pending` citing Draft 0002 | 0002 states it; 0006's flip is the deferred check | none | — |
| C-3 | 0003 / 0006 | closed | 0003 Status 9-16 "Four records remain open … A18 and A20 close on JDR 0001 §JD-18, and A15 and A21 are discharged by MVV"; A10, A12, A17, A19 `Verified` against Final 0006 | The six Pending records on a Final peer's refine are gone | none | — |
| C-4 | 0007 / 0009 | open | 0007:2054-2065 "[ ] **RDR 0009 re-locks against the loss of `Refusal.Guard`**"; 0009:421 "the asserted properties are `Refusal.Kind` / `Refusal.Guard`"; 0009 unmoved | Unchanged; 0009's non-vacuity evidence names a field 0007 deletes | 0009 Phase 1 lands on a `Refusal` with no `Guard` | §1 |
| C-5 | 0007 / JDR | open | 0007:2028-2039 "[ ] **JDR 0001 §D4 is reopened and §JD-8 amended** … a REVERSAL"; JDR §D4 "reaches the CLI through §JD-8's `Detail`"; §JD-8 unanswered | Unchanged | JSON consumers parse a two-level array out of prose | §1 |
| C-6 | 0005 / 0006 / 0007 / 0009 | open | 0005:188 A5 "not new envelope fields or exit groups"; 0006:969-971 fenced `findings` "MUST NOT be `omitempty`"; 0009:1040-1048 fenced "a NEW omitempty field … MUST NOT be resolve.RowRef"; 0007:844-847 one `omitempty` structured field; §JD-8 blank | Unchanged; four MUSTs on one shipped struct, no carrier decision | Three uncoordinated fields on `CLIError` | §1 |
| C-7 | 0002 / 0003 | closed | 0002:1384-1398 fenced seven wire keys; 0003:805-812, 1522-1524 cite 0002's spelling | Landed both sides | none | — |
| C-8 | 0002 / 0004 | closed | 0002:661-675 fenced `role`, `path`, `keys`, `timeout`, `read_back`; 0004:249-258, 415-418 consume by citation | Landed | none | — |
| C-9 | 0005 / 0004 / 0002 / 0007 | superseded → Q-2 | 0002 filled the layout site; 0004's semantics and 0005's placement did not move | as Q-2 | as Q-2 | — |
| C-10 | 0003 / 0007 / 0004 / 0005 | superseded → Q-5 | 0002 now fences the sequence-not-string rule on the stored write value | as Q-5 | as Q-5 | — |
| C-11 | 0007 / 0005 / JDR | superseded → Q-3 | 0007 row 16 "`--tag X=v` against `X eq v` ⇒ row selected"; JD-9 "Wire it to `Observed`"; the guard path is provenance-blind and the write-dependency path is not | The unsticking hazard (guard) and the unusability hazard (RequiresOwned) are two faces of the same wiring | as Q-3 | premortem |
| C-12 | 0008 / 0007 | deferred (`0008/artifacts/deviations.md` D1) | 0008:2390-2413 scenario 3 view-capturing seam vs 0007:1282-1290 fenced "It never sees the view" | Carried; Phase-1 test decides | as iter-3 | — |
| C-13 | 0006 | open | 0006:1068-1076 "two edges reaching the same successor produce one node whose per-tag value sets are the union of theirs" | Unchanged; merge key unstated | `graph-product-too-large` on the first real model | §2 |
| C-14 | 0006 | open — now with a concrete trigger (Q-1) | 0006:1363-1367 "if that count is non-zero on a model the maintainers accept, the response is a successor RDR" | 0002's canonical model is that model, and the count is two | Lint blocks the checked-in model; fix is an RDR | §2 |
| C-15 | 0003 / 0006 / JDR | open — materialized | 0003:820-824 fenced optional default; 0003:1281-1290 fenced marker defaults "reading an unmarked declaration as single-valued or always-present is the same exponential understatement"; 0003 A21 `Pending` "If wrong … the declaration model should invert it"; JDR §D7 "Observation for 0003's next touch … every author will read `enum` as single-valued"; 0002 fixture `[tags.cluster_ready] kind = "bool"` unmarked | The home predicted the authoring error; 0003 re-locked without deciding; the producer committed the error the same day, on a `bool` | Every unmarked guard is a blocking finding | §3, AT-11 |
| C-16 | 0007 / 0003 | open | 0007:2066-2073 "[ ] **RDR 0003 states the empty/omitted `unless` identity (A9).** 0003 mentions the case zero times"; `grep -nE 'empty (or omitted )?.unless.|omitted .unless.' 0003-*.md` → no match after the re-lock; 0003:1143-1170 fenced subtraction over "the row's full `unless` block" | 0003 re-locked (RE-LOCK-ONLY) without stating the identity 0007 books against it; read literally the subtraction of an empty conjunction subtracts the whole product | Lint disables every `unless`-less row, or luck | §1 |
| C-17 | 0003 / 0006 / JDR | closed | 0003:1361 authority row "JDR 0001 §JD-14 as corrected 2026-08-23; RDR 0006 cites the clause (A17)"; A17 `Verified` | The "Contested" residue is gone | none | — |
| C-18 | 0003 / 0006 / 0007 | open | 0003:845-861 fenced "Conformance is enforced for owned keys and unhomed for the rest — a booked gap"; A18/A20 `Pending`; §JD-18 open | Unchanged | Green lint, `guard_unevaluable` at runtime | §3 |
| C-19 | 0004 / 0002 / JDR | closed (0007 box residue → C-24) | 0004:628 "RDR 0002 / JDR 0001 §JD-3 (CLOSED)"; 0002:1125-1143 fenced producer; 0007:2022 "[ ] JDR 0001 §JD-3 landed in RDR 0002" unchecked though landed | 0004 repaired; 0007's checklist is stale | none | — |
| C-20 | 0002 / 0008 / JDR | closed | unchanged since iter-3 | — | none | — |
| C-21 | 0004 | open | 0004:523-527 "A row that consumes the key through `Row.Match` or a guard string alone is dropped by `TagSet.matches` … before `missingOwned` runs, yielding `no_match`"; 0004:628 repeats "a key consumed only through `Row.Match` or a guard is outside it by design" | Re-locked verbatim. Under §D4/0007 a guard atom over an absent key is `guard_unevaluable`, not a `no_match` drop; only the `Match` half of the sentence is true. Unfenced; the `RequiresOwned` scope claim beside it is correct | Scenario 6's guard leg expects `no_match`; the kernel refuses `guard_unevaluable` | §1 |
| C-22 | 0004 | open | 0004:841-849 Prerequisites — four boxes, all `[ ]`, incl. "All Critical Assumptions verified (A1-A8 Verified; A9, A10, and A11 Pending …)" and "RDR 0001 keeps the resolver stateless"; Gate PASS 2026-08-24 (third lock) | Third consecutive Final with every prerequisite unchecked; the brief named it for the re-lock gate and the gate passed it | Implementer cannot tell whether 0004's preconditions hold | AT-1 |
| C-23 | 0005 | open | 0005:601-614 twelve `flow-*` codes; no row for 0002's 25 load categories, 0004's ten refusal classes, 0006's fourteen finding codes, 0009's `escape-row-shape-breach`, 0008's `reserved_tag_key`; §JD-8 unanswered | No change; the code table is one-third the size of the failure surface the cluster now specifies | Skill branches on codes that do not exist | §1, §2 |
| C-24 | 0007 | open — worsened | 0007:148-149 "RDR 0003 … (`Draft`, re-entry at refine)"; 0007:154 "RDR 0002 … (`Draft`, re-entry at refine)"; 0007:687, 939, 951, 2040 "0002's Refinement Context Direction list" (0002 has no such section; the two duties it names — existence-constant emission 0002:949-957, exact-key identity 0002:959-973 — landed in fences, the second as *no* canonicalization plus `unknown tag` on any non-declaration spelling, which is 0007's own fallback arm); 0007:2172 "Hand over, with 0003 in Draft"; README: both `Final` | Both peers re-locked; 0007 now misdescribes their status, cites a deleted section, and leaves boxes unchecked whose done-conditions are met | Phase 4 handoff addressed to a Draft that does not exist | §1 |
| C-25 | 0008 | open | 0008:2334-2366 "hold this RDR at Final-unimplemented until 0002's work starts" | Unchanged | A green 0008 row changes nothing an author perceives | §2 |
| C-26 | 0009 | open | 0009:922-929 fenced "Unwrap() []error MUST be EXACTLY ONE LEVEL deep … VERBATIM" | Unchanged | Rewrite on first real caller | §2 |
| C-27 | 0009 | open | 0009:808-812 "callers MUST check the error before reading the disposition" | Unchanged | `flow resolve` dereferences a nil `Plan` | §2 |
| C-28 | 0008 / 0009 / JDR | open | §JD-5 sharpened, unanswered; 0008:1549-1562 "either error is conforming" | Unchanged | Different error per build on a doubly-breaching table | §1 |
| C-29 | 0002 | open — materialized twice | 0002:1998-2004 Risks "locks after several of its consumers … Mitigation: The layout, routing key, and sentinel are fixed by JDR 0001 §D5–§D7 … peers absorb by citation at their own re-lock" | The re-lock happened after every consumer locked, again; the absorption mechanism is now two `deviations.md` files and a 0007 checklist that cites a deleted section | Q-4, Q-7, C-24, N-8, N-10 | §1 |
| C-30 | 0006 / 0005 | open | 0006:974-977 fenced "extending `clierr.EmitText` (failure) and the `respond.OK` text branch"; 0005:631 "[ ] Add text success payload rendering through `respond.OK`" | Unchanged | Text rendering implemented twice | §1 |
| C-31 | set | open, worse | `find docs/rdr docs/jdr -name '*.md' \| xargs cat \| wc -l` = 75,670; `internal/`+`cmd/` Go = 3,971 (19.1:1); last commit under `internal/` `1c9c0ca` 2026-08-09; `resolve.go:185` `Guard string`; `root.go:55` registers only `newVersionCmd()` | +6,400 lines of design since iteration 3, zero lines of Go; the binary still answers `version` | Nothing runs | §1, premortem |
| C-32 | 0003 | closed | `grep -n 'A16 owns' 0003-*.md` → no match | Residue swept by the re-lock | none | — |
| C-33 | 0006 | open | 0006:1716-1718 "MUST emit `graph-product-too-large` naming the traversal" beside the same code for a declared product | Unchanged | One code, two remedies | §2 |

Tally over iteration 3's 45 rows: **16 closed** (N-1, N-2, N-3, N-5, N-6, N-7,
N-11, C-1, C-2, C-3, C-7, C-8, C-17, C-19, C-20, C-32 — C-17 and C-20 were
already closed at iter-3 and are re-confirmed), **6 superseded** (N-4→Q-8,
N-9→Q-4, N-12→Q-1/C-15, C-9→Q-2, C-10→Q-5, C-11→Q-3), **3 deferred** (N-8,
N-10, C-12 — carried on their written checks), **20 open** (three materialized
or worsened: C-14, C-15, C-24; one worse by measurement: C-31).
**9 new rows**, of which Q-2, Q-3, Q-4, Q-5 are fenced contradictions between
two *members* (the home is silent or blank on each), Q-1 is the producer's
normative fixture failing two consumers' fenced clauses, Q-6 is the
implementation-order graph, and Q-7, Q-8, Q-9 are citation residue.

---

## 1. The inter-RDR failure mode

### The failure

**The set will fail because the re-locks reconciled each member to the home
and no one reconciled the members to each other where the home never spoke.**
Iteration 1's defect was silence (obligations filed on documents that never
received them). Iteration 2's was lag (members locked against answers the home
had not given). Iteration 3's was the inverse (the home answered; the fences
did not move). Iteration 4's is the residue of fixing exactly that and only
that: three re-entries scoped to "restate the fences to §D5/§D6/§D7" produced
three documents whose fences agree with §D5/§D6/§D7 and whose seams to the
*unmoved* members — 0005 above all — were not on the scope line.

The four places it shows are all seams the home explicitly declined to decide
or never saw:

- **The gate list (Q-2).** §D7 wrote "Blank, not decided here: where a rule
  references a gate accessor — 0002's layout has no site for it; owner 0002
  with 0004's semantics." 0002 added the site and a fenced sentence about when
  it is honoured ("before applying that rule's plan") and declared the blank
  filled. 0004 — RE-LOCK-ONLY, re-verify none — added nothing: the string
  "gate list" does not occur in it. 0005, untouched since June, has two fenced
  sentences that put the gate somewhere else ("before the pure kernel call")
  and forbid the verb that applies plans from running gates at all. The row
  carries the list; the kernel row does not; the `Plan` does not; two CLI
  invocations separate the resolve from the write. The blank is now three
  answers.
- **The `--tag` wiring (Q-3).** JD-9 decided the provenance half in the entry
  text ("Wire it to `Observed`") and left the classification half open; 0005's
  qualifier tolerates the open half. Nobody diffed the decided half against
  0005's fences. Under it, `flow resolve` — which "MUST NOT … run read/write
  accessors" — can obtain owned state only as observed tags, and 0002's
  normalizer gives every ordinary rule a non-empty `RequiresOwned` that
  `missingOwned` evaluates against owned provenance. The verb refuses every
  write-bearing row. The kernel's own adversarial test
  (`TestAdv2b_EscapeEdgeMustNotBypassTheOwnedStateRequirement`) is the proof
  that this is by design.
- **The `Block` type (Q-4) and the set-write carrier (Q-5).** Both are Go-type
  facts. 0007 fences a two-member exported type; 0002 fences a third member
  and says an implementer "MUST widen to three here"; the handoff either
  widens 0007's type or maps between two — and 0002 forbids local mirrors of
  kernel constants. 0002 fences that a set-valued write "MUST NOT be joined
  into a single delimiter-separated string, in the stored value"; the stored
  value's only kernel carrier is `Tag.Value string`; 0007 says the element
  encoding is 0003's; 0003 says nothing. These are the two seams where a
  compiler, not a reviewer, will report the contradiction — and the two seams
  no answer-vs-fences check was ever scoped to reach, because the home has no
  entry for either.

Under the gate's own rule — a fenced contradiction between two members with no
home is a joint decision to be homed *before* implementation — Q-2 and Q-3 are
the iteration-2 pattern recurring one level down, on a producer that has
already re-locked once.

### Root cause across the RDRs

**Enabler 1 — re-entry scope was "member × home," never "member × member."**
STAGE-SCOPED and RE-LOCK-ONLY were defined by which *assumptions* to re-verify,
and the answer-vs-fences checks were defined by which *JD entries* the member
tolerated. Neither reaches a fence in 0005, because 0005 tolerates JD-8 and
JD-9 only and both are "open." So 0002 could fence a gate-list semantics and
0005 could keep a contradictory one, and both passed their gates the same day.

**Enabler 2 — the home wrote a decision into an "open" entry.** JD-9's text
contains a directive ("Wire it to `Observed`") and a status ("Still open …
the classification arm"). The iteration-3 contract ran the answered-tolerance
check on JD-15/16/17 — entries that flipped to *Decided* — and not on JD-9,
whose decisive half never flipped. A member cannot tell, from the entry, which
sentence it is tolerating.

**Enabler 3 — the producer's fixture was re-authored under one consumer's
rules and never linted under another's.** 0002's A1 re-verification ran the
strict decoder and the normalizer over the §D7 fixtures and stamped
`Verified`. It did not compute 0003's assignment-count table over the same
file — the table 0003 fences "an implementation MUST apply" — which would have
returned `2^2` for a `bool` and `2^10` for `iter` and the question "why is
`cluster_ready` not single-valued?" (Q-1). 0006's fence cites `0003::A21`, a
*Pending* assumption whose "If wrong" inverts the default. The cluster fenced
a rule on a Pending premise and then falsified the premise in its own fixture.

**Enabler 4 — the implementation graph has no head (Q-6).** Every document
names what it waits on, correctly. 0002 waits on 0007 Phase 1 and says the
reshape "has no owner … no phase in this RDR's plan claims it" — 0007's Phase
1 is that reshape by name. 0007 waits on 0009's re-lock, on JD-8 reopening §D4,
on 0003 stating an `unless` identity 0003's re-lock did not state (C-16), and
on a "Refinement Context" section 0002 deleted. 0009 and 0008 wait on JD-5 and
JD-8. JD-8 "re-walks 0005." 0005's last edit is June. The one entry with no
sibling scheduled to answer it is the one every path passes through.

### The symptom the user sees

The user is the first skill author. They write the RDR model exactly as 0002's
canonical fixture does, run `intrastate lint`, and get two blocking
`graph-unprovable-coverage` findings naming `iter` and `cluster_ready`. They
read 0003, learn that a `bool` can hold both values unless told otherwise, add
`single_valued = true` to both, and lint goes green (C-15, Q-1). They run
`flow read-state --artifact rdr=./0010.md` and get
`{status: Draft, stage: prelock, profile: large, iter: 2, …}`. They pass those
to `flow resolve --tag status=Draft --tag stage=prelock … --outcome
round-clean` and get `owned_state_unavailable` naming `stage` and `iter` —
every time, for every rule, because `--tag` is observed and the rule writes
`stage` (Q-3). They read 0005's fence, which says `resolve` runs no readers,
and 0004's, which says owned state comes from readers, and file a bug titled
"resolve cannot resolve." Meanwhile the `[gate.rdr-lock]` they declared, which
0002 validated and rendered on the dump, has never been invoked by any verb
(Q-2).

That is if any of it exists. The last commit under `internal/` is fifteen days
old; the design corpus grew 6,400 lines in the interval (C-31).

---

## 2. The RDR that will be rewritten within six weeks of shipping

**RDR 0005 — Skill Integration CLI Contract. Not touched since 2026-06-19 by
anything but a Status qualifier; the only member never re-entered; and the
document every finding above lands on when it reaches a user.**

Iteration 3 named 0002 and was right for that iteration: 0002 has now been
rewritten — fourteen home landings, both fixtures, the MVV — and re-locked.
The rewrite that is now unavoidable is the one the cluster has been routing
*around*: 0005 is the sibling JD-8 "re-walks," and JD-8 is the entry whose
answer is "the whole code table" (C-23). Count what 0005's twelve `flow-*`
codes must now map: 0002's twenty-five load categories (fenced, snake_case
identifiers, "an API surface, not message text"), 0004's ten refusal classes
(`read_back_incomplete` "MUST NOT be reported as `read_back_mismatch`" —
0005's table has only the mismatch), 0006's fourteen finding codes and a
non-`omitempty` `findings` field, 0009's `escape-row-shape-breach` and its own
`omitempty` field, 0007's per-atom payload and its one `omitempty` field,
0008's `reserved_tag_key` plus a near-miss advisory channel. Twelve codes, one
`Detail` string, and a June-locked A5 that pre-authorizes "not new envelope
fields."

Then the verbs. `flow resolve` as fenced cannot return a plan for any
write-bearing rule under JD-9 (Q-3). `flow next` / `flow resolve` "MAY invoke
declared gate accessors … before the pure kernel call" where 0002 puts the
gate after the plan and 0004 puts it nowhere (Q-2). `set-state` "MUST NOT treat
context-only `--tag` values as writes" and "MUST invoke only declared write
accessors" — so the gate cannot run there either. `--tag name=value` "supplies
one scalar tag fact; duplicate tag names … refused until the transition-model
layer exposes a first-class structured literal for set-valued tags" — the
layer has (0002:843-853), as a member sequence the kernel's `Tag.Value string`
cannot hold (Q-5). And "0002 says the model enters at the CLI seam" with "no
path form" (JD-8 widened note) — 0005 has no `--model` flag and 0002's loader
"takes already-read bytes plus a source id … A path-taking convenience wrapper
MAY live at the CLI seam" (0002:1782). Every one of these is a fenced 0005
clause or a missing one; none can be cited in.

**Why 0005 and not 0002 again.** 0002's residue is cosmetic (Q-7) or a fixture
edit (Q-1, two lines) or a type-identity decision it shares with 0007 (Q-4) or
0003 (Q-5). 0005's residue is its verb contract. 0002 has now been through
every lens twice; 0005 has been through them once, against a cluster that had
five members and no home. It is also the only member whose tolerance
qualifier tolerates an entry (JD-9) that contains a directive it contradicts.

**What would have to be true for 0005 *not* to be rewritten.** JD-8 would have
to be answered "twelve codes suffice; everything else is `Detail`," JD-9's
provenance directive withdrawn so `--tag` enters as owned, the gate placed
before the kernel call as 0005 says (so 0002's "before applying the plan"
sentence is the one that yields), and set-valued tags kept out of the CLI. Each
of those is the arm the home's principles reject: P3 (structured failure), P2
(the kernel refuses rather than guesses — `--tag` as owned is the "shadowing"
`assemble` pins against), P6 (one home per fact). Nothing in the record offers
it.

**The six-week clock.** Even after the rewrite, 0005 owns the exit-code
grouping that "encodes remedy" (P4) and has one bit to spend: `GroupUserEnv`
and `GroupInternal` both exit 2, `GroupEnvUnavailable` exits 3. 0004's
"may have been applied and was not verified" (`read_back_incomplete`,
post-mutation `timeout`) is a remedy class — re-read before acting — that
neither exit code carries. That is the second rewrite.

*(Runner-up: 0007 — seven unchecked prerequisites, two peers described as Draft
that are Final, four citations to a deleted section, a fenced two-member type
its consumer has fenced at three (Q-4), a Phase 3 blocked on an encoding it
assigns to a document that will not state it (Q-5), and the reshape every
member sequences behind. It is Final at `526c481` and has not moved since
iteration 2's status repairs.)*

---

## 3. The cross-cutting assumption that will not survive first contact

**The assumption: that owned state reaches the kernel through an accessor
snapshot, and that some *other* document's verb is the one that reads it.**

Every member is written against `Input.Owned` being populated:

- 0002 declares readers with `keys` and requires every owned tag to be served
  by exactly one (fenced), and derives `RequiresOwned` so the kernel can refuse
  when the snapshot lacks a key (fenced).
- 0004 defines the read that produces the snapshot — complete, typed, absence
  as omission — and the write that read-backs it (fenced), "executed only from
  a successful transition plan."
- 0007 evaluates guards over "the assembled view," provenance-blind for
  presence, and refuses when the view lacks a key (fenced).
- 0003 proves coverage over "conforming evaluation views" and books the
  view-level conformance check on JD-18 (fenced: "the check belongs to the
  kernel's view assembly").
- 0006 promises "a green lint means resolution succeeds" (P5) over the
  owned-state graph — the states the snapshot will hold.
- 0009 and 0008 evaluate their preconditions at `Resolve` entry, over `Input`.
- 0005 — the only document that owns a verb — fences that `resolve` "MUST NOT
  … run read/write accessors," gives `read-state` the readers and `set-state`
  the writers, and supplies `resolve` its state as `--tag`, which JD-9 makes
  observed.

Nobody's verb assembles the owned snapshot and then calls `Resolve` in the
same invocation. `read-state` reads; `resolve` resolves; `set-state` writes;
and the only channel between them is the caller re-typing tag values as
flags that the kernel will not count as owned. The assumption is load-bearing
in every fenced clause that mentions `owned_state_unavailable`, `RequiresOwned`,
`missingOwned`, or "the view," and it is false at the one place a user touches
the system.

### How it breaks

**First contact, the `flow resolve` implementer (Q-3).** They wire `--tag` to
`Observed` as JD-9 says, run the canonical model, and get
`owned_state_unavailable` on every ordinary rule. They have three moves: wire
`--tag` to `Owned` (JD-9 forbids; `assemble`'s owned-over-observed precedence
means a stale `--tag` shadows nothing but a missing reader now shadows
everything); make `resolve` run the readers (0005's fence forbids; 0004's
executor becomes a dependency of a verb that "MUST NOT discover artifacts");
or add a fifth verb that reads-then-resolves (0005's four-verb A2 is
`Verified`). Whichever they choose, no fence describes it.

**First contact, the gate (Q-2).** The same implementer finds `gate` on the
normalized row, no field for it on `resolve.Row` or `Plan`, and three documents
placing its evaluation at three times. They run it in `resolve` before the
kernel call, on every candidate, because that is the only verb with the
`--artifact` binding and the only sentence with a "when." Every gate in the
model runs on every resolve, for rows that will not match. `[gate.rdr-lock]`'s
`keys = ["status"]` binds nothing anyone can name (Q-8).

**First contact, the fixture (Q-1).** 0006's implementer lints 0002's
canonical model as scenario 23 requires and records "count of blocking findings
a model its authors consider correct receives" = 2. 0006's own text says the
response is a successor RDR on guard-aware pruning. The correct response is
two lines in a TOML file and a decision the home already flagged and 0003
declined at its re-lock (C-15). A `bool` that may hold both values is the
signal that the conservative default was never authorable; A21 was never
going to be verified by a Phase 3 fixture because the Stage 4 fixture already
refutes it.

**First contact, the compiler (Q-4, Q-5).** `resolve.Block` has two constants
or three. `Tag.Value` is a string or it is not. These are not review findings;
they are the first two build failures of Phase 2, and the documents that would
decide them are Final on both sides with no home entry between them.

---

## 4. Premortem

*Written from six months after the fourth cluster gate.*

The fourth iteration reported RECONCILED WITH TOLERANCES. It had to: the cap
was spent, the three re-locks had landed what iteration 3 named, and the
findings left were "between members, not against the home," which the
tolerance table had a column for. Nine rows were written into three
`deviations.md` files and two new interface-record entries (JD-19: gate
placement; JD-20: `--tag` provenance vs `RequiresOwned`) were opened with 0005
as the re-walked sibling. 0005 was not re-entered — it had never been, and the
gate could not demote. Implementation started on 0007 Phase 1 because it was
the head of every dependency list.

**Weeks 1–3, the reshape.** `Row.Guard string` became `[]Atom`. `Block` was
exported with two constants per 0007's fence. 0002's Phase 2 then needed a
third; the implementer added `BlockMatch = "match"` to `resolve` — 0007's
"exactly two" was now false in code — because minting a second type and
mapping at the handoff was the reconstruction step §D1 had just spent a
re-entry removing (Q-4). 0007's payload could now carry `match`, which no
refusal ever produces; a test asserting the payload's `Block` is `all` or
`unless` was written with a `default: t.Fatal("unreachable")` arm.

**Week 4, the set write.** The kata fixture's `labels` write reached
`Row.Writes` as `Tag{Key: "labels", Value: "x,y"}` — the join 0002's fence
forbids "in the stored value" — because `Tag.Value` was a string and 0001 was
Implemented and 0007's reshape had left `Tag` alone (Q-5). 0004's read-back
compared `"x,y"` to the artifact's list and reported `read_back_mismatch` on
every set write. The fix was an encoding — `"[\"x\",\"y\"]"` — documented in a
code comment that cited 0003, which had never stated one.

**Week 6, lint.** 0006 lints the RDR model: two `graph-unprovable-coverage`
findings, `iter` and `cluster_ready` (Q-1). Scenario 23's count is non-zero on
a model the maintainers accept; per 0006's text a successor RDR on guard-aware
pruning is seeded. The reviewer adds `single_valued = true` to a `bool`, the
findings clear, and the RDR is closed as "fixture defect." Three weeks later
the first real author writes `kind = "bool"` for a flag, gets the same finding,
and the runbook gains "always mark bools single-valued." The default inversion
0003 A21 named as "a real fork … decided here and carried to RDR 0002 and RDR
0006" was decided in a runbook (C-15).

**Week 8, the CLI.** `flow resolve` shipped with `--tag` wired to `Observed`
per JD-9. Every resolve of every write-bearing rule refused
`owned_state_unavailable` (Q-3). The implementer read 0005's "MUST NOT run
read/write accessors," 0004's "owned state comes from readers," and JD-20's
blank, and added `--owned name=value` — a flag that enters as owned provenance,
which is exactly the `Input.Owned` wiring JD-9 said would open "the exposure."
0007's row 16 hazard ("`--tag X=v` against `X eq v` ⇒ row selected and a
plan") was now `--owned X=v` and the operator could unstick any refusal from
the command line. Nobody had re-derived C-11 under the new flag name.

**Week 9, the gate.** `[gate.rdr-lock]` had been declared, validated, and
dumped for two months without executing. `resolve` ran it "before the pure
kernel call" per 0005, on every candidate row that carried a `gate` list,
because the matched row was not yet known; a gate that returned `deny` for a
row that would not have matched refused the whole resolve (Q-2). The fix
moved gate evaluation after `Resolve` and before returning the plan, which
0002's "before applying that rule's plan" could be read to permit and 0005's
"before the pure kernel call" could not. 0004 still had no sentence about it.
`set-state`, in a separate process, applied the plan without the gate — the
allow had been observed one invocation earlier and carried in the caller's
memory.

**Week 12, 0008 and 0009.** JD-5 was still open. 0009 Phase 1 landed
`CheckValid` at `Resolve` entry; 0008 Phase 1 landed the reserved-key
predicate at the same entry; the order was whichever the implementer wrote
first, and 0008's scenario 6 "must NOT assert which of the two errors is
reported" was the test that pinned nothing (C-28). 0009's non-vacuity evidence
still cited `Refusal.Guard`, which had been deleted in week 2 (C-4).

**What actually killed us.** Not any one contradiction — each was two lines
in code. What killed us was that the fourth gate had a column for "between
members, not against the home" and used it, so that the four seams the home
had never seen were tolerated as if the home had answered and the members
were merely lagging. They were not lagging; they were disagreeing, in fences,
with no venue. And the one member that owned the verbs — the only place the
user meets the cluster — was never re-entered, because the gate's scope was
defined by which documents had moved, and 0005 was the one that never did.

---

## 5. Acceptance tests that would have caught each failure at RDR-review time

AT-1 through AT-10 from iterations 2 and 3 stand; AT-8 (home answer → member
demotion) is now discharged for JD-15/16/17 and would have caught nothing new
this pass, which is the point: the failures below are not home-vs-member. Four
are added.

### AT-11 — The producer's normative fixture is run through every consumer's fenced arithmetic before the producer re-locks
*Catches Q-1, C-15, C-14*

```gherkin
Given a producer RDR names a fixture "canonical" or "what implementation
  tests must promote"
And a consumer RDR fences a rule that computes over declared fields of that
  fixture (0003's assignment-count table; 0006's `graph-unprovable-coverage`
  trigger)
When the producer's Stage 4 re-verifies the fixture
Then the re-verification MUST include the consumer's computation over the
  fixture — every guard dimension's assignment count and projection status —
  and MUST record the lint verdict the consumer's fences yield
And a fixture whose verdict is a blocking finding MUST NOT be stamped
  Verified as "canonical"; either the fixture is repaired or the consumer's
  default is re-decided at the home, never both left standing
```
Run today: `[tags.iter]` and `[tags.cluster_ready]` in
`0002/evidence/spikes/iter-2/rdr-fixture.toml` carry no `single_valued`;
0003:1273-1290 assigns `2^10` and `2^2`; 0003:948-952 gives `iter.lt` and
`cluster_ready.eq` no projection; 0006:945-947 mints two blocking findings →
**failure**. A `bool` requiring the marker is the C-15 fork, materialized.

### AT-12 — Every CLI verb is walked through the kernel's provenance-aware inputs before the verb's fence locks or a JD directive that touches provenance is written
*Catches Q-3, C-11*

```gherkin
Given a verb's fence names which accessors it MAY and MUST NOT run and which
  flags supply state
And the home records a provenance for each flag's tags (JD-9)
When either text is decided
Then for each verb the gate MUST list, per kernel input (`Input.Owned`,
  `Input.Observed`, `Input.Recognized`), the source that populates it in that
  invocation — accessor, flag, or nothing
And any verb that calls `Resolve` with `Input.Owned` populated by nothing
  MUST be recorded as unable to return a plan for any row whose
  `RequiresOwned` is non-empty — with the producer's derivation rule cited —
  and the contradiction homed before the verb's fence stands
```
Run today: `flow resolve` — `Input.Owned` ← nothing (0005:380-383); `--tag` →
`Observed` (JD-9); every ordinary rule's `RequiresOwned` non-empty
(0002:830-833, 1125-1130) → **failure**. `flow next` — same inputs; same
result for any candidate summary that consults `missingOwned`.

### AT-13 — Every field a producer puts on the normalized row has a named carrier to every consumer that reads it, or a recorded "no carrier" with an owner
*Catches Q-2, Q-4, Q-5*

```gherkin
Given a producer fences a field on the normalized row (gate ids; the
  three-valued block; a set-valued write value as a member sequence)
When the producer locks
Then for each consumer that reads the field at a later stage (the executor
  for gate ids; the kernel handoff for block; `Row.Writes` for write values)
  the producer MUST name the Go carrier — the kernel field, the plan field,
  or the CLI channel — and its type
And where the carrier's type cannot hold the field (a `string` for a
  sequence; a two-constant type for three values; no field at all for the
  gate list) the gate MUST route the seam to the home as a joint decision
  with the consumer named, never leave "this fills the site" standing
```
Run today: gate ids → `resolve.Row`/`Plan` have no field; 0005's channel
contradicts 0002's timing → **failure**. Block → `resolve.Block` "exactly two"
(0007) vs "MUST widen to three" (0002) → **failure**. Set write →
`Tag.Value string` vs "MUST NOT be joined … in the stored value" → **failure**.

### AT-14 — The cluster's Prerequisites checklists form a rooted DAG with an owner at the root
*Catches Q-6, C-4, C-5, C-16, C-24*

```gherkin
Given every member's Prerequisites section lists boxes naming peers,
  JD entries, or phases
When the cluster gate runs at N ≥ 2
Then the gate MUST draw the graph: member-box → the document or entry it
  waits on
And every JD entry at a root MUST have a named sibling scheduled to answer
  it at a named stage, or the report MUST say "no owner" for that root
And every box whose done-condition is met in a peer's current text MUST be
  reported as stale (0007:2022, 2024, 2040), and every box citing a section
  that no longer exists MUST be reported as dangling (0007:2040)
And a member that names a peer's phase as "unowned" when the peer's plan
  claims it (0002:2073-2076 vs 0007 Phase 1) MUST be reported as a
  restatement defect
```
Run today: roots are JD-8 (no sibling scheduled; "re-walks 0005," untouched
since June), JD-5, JD-18 → **failure**. Four 0007 boxes are stale or dangling
→ **failure**.

---

## 6. Delta vs iteration 3

Of iteration 3's 45 rows: **16 closed** — every one on the three re-locks
(N-1, N-2, N-3, N-5, N-6, N-7 by 0002/0003/0004's fences; N-11 by the fact of
re-locking within the cap; C-1, C-2, C-7, C-8 by 0002's layout and
write-replaces clauses; C-3, C-17, C-32 by 0003's sweep; C-19 by 0004's
citation) — **6 superseded** into the seam rows (N-4→Q-8, N-9→Q-4,
N-12→Q-1/C-15, C-9→Q-2, C-10→Q-5, C-11→Q-3), **3 deferred** on their written
checks (N-8, N-10, C-12), **20 open**, of which C-14, C-15 and C-24 are worse
in kind (a concrete trigger; a materialized fork; a deleted section cited)
and C-31 worse in measure (19.1:1; the interval added 6,400 lines and no Go).

The nine new rows have one origin: **the re-locks were scoped to the home and
the contradictions that remain are between members.** Four are fenced on
both sides with no home entry (Q-2 gate placement, Q-3 `--tag`/`RequiresOwned`,
Q-4 `Block` cardinality, Q-5 set-write carrier). One is the producer's
normative fixture failing two consumers' fences (Q-1). One is the
implementation-order graph with an unowned root (Q-6). Three are citation
residue of re-locking in the wrong order — 0002 before 0003 and 0004 (Q-7),
the home and 0004 after 0002 settled §D7(ii) (Q-8), and 0004 twice without
citing 0008 (Q-9).

What did *not* change: no code table (C-23), no carrier (C-5, C-6), no
precedence (C-28), no conforming-view enforcer (C-18), no `unless` identity
(C-16), no 0009 re-lock (C-4), no 0005 touch at all, and no Go (C-31).

---

*Ledger rows: 54 (45 carried, 9 new). Acceptance tests: 4 new, 10 carried.
Sections 1–5 are the argument; the ledger is the routable record. Any defect
not carried in the ledger is unreported by construction.*
