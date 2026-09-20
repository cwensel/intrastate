Model: claude-fable-5-1

# COVE pre-lock lens — RDR 0012 (declared-kind carrier at the guard value seam)

Source tree read: `main` @ 090c539. Record read through the projector
(`recs inspect`, sections batched to scratch). All 45 `source-anchor`
edges project `resolved: true`; none absent, none false — so the sweep
below is over the UNANCHORED claims (call-flow, "only place", "sibling
already does", "output unchanged") plus a re-read of every anchor the
post-propose rewrite touched (351+/131- across 37 hunks since 43c1e8b).

## Step 0 — grounding

### Confirmed anchors (compact)

- `internal/guard/grammar.go::Evaluator` — `type Evaluator struct{}` (l.97); `eq` arm `value == atom.Literal` (l.107); `in` arm `parseSetLiteral` + `slices.Contains` (l.109-113); `lt..gte` parse both sides (l.114-123); `contains`/`parseHeldSet` nil-check (l.179-196).
- `internal/table/load.go::conformKind` (l.1853) — int arm is `strconv.Atoi` only; bool arm is token equality. `valueMembers` (l.1748) renders `float64` via `strconv.FormatFloat(t,'g',-1,64)` before any kind check.
- `internal/table/load.go::ConformValue` (l.1813) → `conformKind` + `conformDomain`; called from `internal/cli/flow_input.go::canonicalValue` (l.751, l.769) and the `[emit]` site (load.go l.506). So a tightened int arm reaches `flow-tag-invalid`/`flow-write-invalid` and `emit_value_out_of_domain` as C5 says.
- `internal/table/normalize.go::atomsFromBlock` (l.23) — `CatUnknownTag` on undeclared key, block-agnostic (A1 holds). `normalize.go::atom` (l.43) calls `conform(decl, operator, members)` (l.172) for every block incl. match.
- `internal/table/normalize.go::renderWrites` (l.580) — `conform(decl,"eq",members)` at l.624 under `CatMalformedTagDeclaration`; constructs no evaluator on main (C4 dependency note holds).
- `internal/accessor/model.go::ReadResult.OwnedSnapshot` (l.424) — `resolve.Tag{Key: v.Key, Value: v.Value}` raw; `internal/resolve/resolve.go::assemble` (l.220) / `merge` (l.236) store `t.Value` with no conform. `rg ConformValue` over non-test sources: only flow_input.go + load.go — the owned door is unconformed (Problem Statement holds).
- `internal/resolve/resolve.go::TagSet.matches` (l.175) — `tv.value != w.Value` byte compare (Key Discovery holds for the KERNEL match path).
- `internal/table/model.go::KernelTable` (l.540) returns `resolve.Table{Revision, Outcomes, Rows}` — no declarations (ALT1 rejection holds).
- `internal/guard/declaration.go::DeclarationOf` (l.47) — `declarationOf(m.Tags[key])`; takes `*table.Model` (audit row 1 holds).
- `internal/guard/product.go::valueSatisfies` (l.541) — `seam := Evaluator{}`; literal = `strings.Join(atom.Literal,"")` (authored bytes) or `renderSet` for `in`/`contains`. Called from `Denotation(m *table.Model, …)` (l.437) one frame up; `Denotation` returns the unprojectable set on `GuardUnevaluable` and refuses `eq`/`in` over non-single-valued keys (l.481-483).
- `internal/guard/assignment.go::valueAssignments` (l.273) int arm renders `strconv.Itoa(n)` over `[Min..Max]` (A4 held side holds). NOTE the record's anchor spells this `internal/guard/product.go::valueAssignments`; the symbol is DEFINED in assignment.go:273 and only CALLED from product.go:194/336/450/477 — cosmetic misattribution, see F-9.
- `internal/graphlint/reach.go::atomAdmitsValue` (l.502) — `guard.Evaluator{}.Evaluate(...)`; two frames below `matchSatisfiable(m *table.Model, …)` (l.449) via `ownedAtomSatisfiable` (l.468) (A3 frame count holds).
- `internal/cli/flow_resolve.go::guardSeam` (l.561) — zero-arg, returns `guard.Evaluator{}`; called at flow_resolve.go:258 (req.model in scope) and :445 (`escapeClassOf`, `req.model` in scope, resolves over a `probe` derived from `req.model.KernelTable()`).
- `rg 'Evaluator\{\}' --type go` non-test: exactly product.go:542, reach.go:507, flow_resolve.go:561 (A3/C4 enumeration holds).
- `internal/resolve/guardcontract.go::TestGuardEvaluatorContract` (l.14) — 12 cases, all `Key: "subject"`, `Block: BlockAll`.
- `internal/guard/guard_evaluator_0003_test.go::TestReq34_…` — `NumIn() != 3` (l.36), TagSet-typed-param scan (l.40-46), `NumField() != 0` (l.47) (A2 proxy claim holds).
- `internal/resolve` non-test imports of `internal/table`: zero (A5 holds). `internal/guard` imports both `table` and `resolve`.
- `resolve.go:60` `KindGuardUnevaluable RefusalKind = "guard_unevaluable"`; `resolve/guard.go:82` `ReasonUncomparable = "uncomparable"`; CLI code minted at `internal/cli/flow_exec.go:718` as `"flow-" + strings.ReplaceAll(string(kind), "_", "-")`.
- Peer quotes: `0007:A17` Evidence contains verbatim "The declared kind needed to *parse* the value is known to the evaluator from the table it was built for, not from the kernel"; `0007:A18` and `0007:C1` "MAY answer unevaluable for a present value it cannot compare (A18)" verbatim; `0003:C8` "MUST be rejected before resolution" verbatim; `0003:C9` contains neither "stateless" nor "zero fields".
- Corpus (own scan, tomllib over `models/` + `internal/**/testdata`, 123 files / 121 parseable): `eq`/`in` atoms over an `int`-declared key in guard OR match blocks: 0. Non-canonical int spellings across guard/match/`[initial]`: 0. `[rule.write]` int values: 97, non-canonical: 1 — `internal/table/testdata/neg/neg-write-value-wrong-kind.toml` writes `iter = "two"`, a neg fixture already refused as not-an-int. S5 and C5's "zero violators" hold.

### REFUTED / NOT-FOUND (full)

1. **REFUTED — §Key Discoveries "lint and the runtime share the seam for GUARD atoms only."** `internal/graphlint/reach.go::matchSatisfiable` (l.449-465) iterates `row.Atoms` keeping ONLY `a.Block == table.BlockMatch` over owned-provenance keys, and its doc says "no guard atom is consulted at all"; `ownedAtomSatisfiable` → `atomAdmitsValue` (l.502) then runs each MATCH atom through `guard.Evaluator{}.Evaluate`. So the third C4 site types MATCH atoms, not guard atoms. See F-1.
2. **REFUTED (precision) — A3 Evidence "`guardSeam()` is zero-arg but its callers hold `req.model`."** `guardSeam()` has a third non-test call frame the record never names: `internal/cli/flow_next.go:492` inside `probeRow(row table.Row, view map[string]string, owned []resolve.Tag, observed []resolveTag, all bool)` (l.468-474) — no model in that frame; `req` is one frame up at flow_next.go:259. See F-3.
3. **NOT-FOUND (cosmetic) — A4 Evidence anchor `internal/guard/product.go::valueAssignments`.** Defined at `internal/guard/assignment.go:273`; product.go only calls it. See F-9.
4. **Imprecise — C5 "the write-block and emit-value categories."** There is no write-block category: `[rule.write]` kind/domain conformance refuses under `CatMalformedTagDeclaration` = `"malformed_tag_declaration"` (normalize.go:624-625; category.go:28), by deviation D3. Emit is `CatEmitValueOutOfDomain` = `"emit_value_out_of_domain"` (load.go:508). See F-8.

### Inverse sibling check

New rules: (i) a declared-kind lookup at the value seam; (ii) canonical-int-spelling refusal at load.

- (i) Searched `rg 'func .*[Cc]onform|func Conforms'` non-test. The only adjacent decider is `internal/guard/declaration.go::Conforms(m *table.Model, v View) error` (l.280) — the 0003 view-conformance premise. It checks presence of required keys and single-valuedness (`heldValueCount`) ONLY, does NOT check kind, and has ZERO non-test callers. It is the natural venue JDR 0001 §JD-18 names and the closest existing "is this held value legal for its declaration" sibling; the record's Sibling-path check names only `DeclarationOf`. Not a duplicate of the new rule; see F-6.
- (ii) Searched load.go/normalize.go/flow_input.go for any `Itoa`-round-trip or leading-zero check: none exists. `conformKind` (`Atoi` only) and `conformDomain` (`Atoi` again + bounds) are the only int checks. **Searched, none exists** — C5 is genuinely new.
- Second suite implementer (relevant to C3): `internal/resolve/guard_mvv_test.go::conformingContractSeam` (drives `TestGuardEvaluatorContract` at l.462 and, wrapped in `recordingSeam`, at l.257). Its `eq` arm is `value == atom.Literal`. See F-2.

### Spans widened past, and why

- Widened from the three A3 sites to every CALLER of `guardSeam()` (flow_resolve.go:258/:445, flow_next.go:492) because C4's "one model, one evaluator" is a claim about the frames that own the model, not the constructor wrapper.
- Widened from `atomAdmitsValue` to `matchSatisfiable` (reach.go:449) because the Key Discovery claims which atom population the third site evaluates — decidable only at the caller.
- Widened from `guardcontract.go` to every caller of `TestGuardEvaluatorContract` (guard_evaluator_0003_test.go:67, guard_mvv_test.go:257/:462) because C3's extension binds every implementer, and the record names only one.
- Widened from `conformKind` to its full caller set (`conform`, `ConformValue`, normalize.go:172/:624, load.go:506/:1701, flow_input.go:751/:769) to check C5's "binds every authoring site" and its category list.
- Widened from the spike's `models/`-only grep to the full 123-file corpus (the spike's "ZERO eq/in atoms over an int tag" is stated over 123 models but its supporting `rg` ran over `models/` alone, where only 1 of the corpus's 100 int-declaring files lives). Result confirms the claim.
- Read `0007:A17/A18/C1` and `0003:C8/C9` through the projector to check the quoted fences rather than trusting the record.

## Step 1–2 — verification questions (answered independently)

1. **[code] Are there exactly three non-test `Evaluator{}` construction sites on `main`?** Yes. `rg 'Evaluator\{\}' --type go` non-test → `internal/guard/product.go:542`, `internal/graphlint/reach.go:507`, `internal/cli/flow_resolve.go:561`. (Test files also use the zero-value form `var ev guard.Evaluator` at guard_grammar_0003_test.go:31/177/426, guard_scope_0003_test.go:297/353, guard_declaration_0003_test.go:267.)
2. **[code] Does every immediate caller of `guardSeam()` hold a `*table.Model`?** No. `internal/cli/flow_next.go:492` is inside `probeRow(row table.Row, view, owned, observed, all)` — no model parameter; the model is at the caller's frame (flow_next.go:259, `req`). flow_resolve.go:258 and :445 do hold `req.model`.
3. **[code] Does graphlint's `atomAdmitsValue` evaluate GUARD atoms?** No. `internal/graphlint/reach.go::matchSatisfiable` (l.449): `if a.Block != table.BlockMatch { continue }` and doc "no guard atom is consulted at all". The seam call at l.507 decides MATCH atoms over owned keys.
4. **[code] Does the kernel's match path byte-compare?** Yes. `internal/resolve/resolve.go::TagSet.matches` l.178: `tv.value != w.Value`.
5. **[code] Does `conformKind`'s int arm admit `"00"`, `"+1"`, `"-0"`, and is it reached from `--tag`/`--write`?** Yes and yes. load.go:1856 `strconv.Atoi(member)` only (Go's `Atoi` accepts sign and leading zeros); `flow_input.go::canonicalValue` l.751/l.769 → `table.ConformValue` l.1814 → `conformKind`.
6. **[code] Is any owned reader value conformed before it reaches the seam?** No. `accessor/model.go::OwnedSnapshot` l.435 `resolve.Tag{Key: v.Key, Value: v.Value}`; `cli/flow_exec.go:332-337` appends them unchanged; `resolve.go::assemble` l.230 `view.merge(in.Owned, ProvenanceOwned)`; `merge` l.238 stores `t.Value`. No `ConformValue`/`conformKind` caller exists outside flow_input.go and load.go.
7. **[code] Does the conformance suite have an implementer other than `guard.Evaluator`?** Yes. `internal/resolve/guard_mvv_test.go::conformingContractSeam` (struct{}, `eq` arm `value == atom.Literal`) is driven at l.462 and via `recordingSeam` at l.257 by 0007's REQ-71 tests.
8. **[code] Do the suite's cases carry per-kind keys?** No. `guardcontract.go` l.50-55 builds every atom with `Key: "subject"`; kind is nowhere in the case struct (l.17-23).
9. **[code] Does `internal/resolve` (non-test) import `internal/table`?** No. `rg 'intrastate/internal/table' internal/resolve -g '!*_test.go'` → 0 hits. `internal/table/model.go` imports `internal/resolve` (returns `resolve.Table`).
10. **[code] Does the committed corpus contain an `eq`/`in` atom over an `int` key or a non-canonical int spelling?** No / no. Own tomllib scan over 121 parseable models: 0 `eq`/`in` over int (guard or match); 0 non-canonical spellings in guard/match/`[initial]`; `[rule.write]` has 97 int writes, the single non-canonical one being the neg fixture `neg-write-value-wrong-kind.toml` (`iter = "two"`, already refused).
11. **[RDR] Does the record address `0007:A17`'s HEADLINE — the seam "needs neither the view, nor provenance, nor sibling atoms, nor tag declarations at evaluation time"?** RDR is silent. §Approach and Key Discoveries quote only A17's Evidence sentence ("known to the evaluator from the table it was built for") and conclude "A17 froze the signature, not the constructor". C2's arm reads `ev.kinds[atom.Key]` at evaluation time. The record treats the 0003 zero-field test as a recorded deviation but records no deviation against A17's headline.
12. **[RDR] Does C5 name the actual refusal category for `[rule.write]` values?** No — it says "the write-block and emit-value categories". Source: `malformed_tag_declaration` (normalize.go:624-625, deviation D3) and `emit_value_out_of_domain` (load.go:508).
13. **[RDR] Does the record name `flow next` as a seam consumer whose user-visible output changes?** RDR is silent. `flow_next.go:545` documents the `guard_unevaluable` payload as "the ONLY source of `uncomparable`" in the `unknown` pairs; a typed seam changes which rows `flow next` reports as `uncomparable`. Consequences/F1 name only `flow resolve`.

## Step 3 — FINDINGS

### F-1 — Third construction site types MATCH atoms, not guard atoms; Key Discovery refuted
- anchor: 0012:§research-findings (Key Discoveries bullet "lint and the runtime share the seam for GUARD atoms only")
- class: a
- evidence: `internal/graphlint/reach.go::matchSatisfiable` l.449-465 (`a.Block != table.BlockMatch → continue`; doc "no guard atom is consulted at all") → `ownedAtomSatisfiable` → `atomAdmitsValue` l.502-517 (`guard.Evaluator{}.Evaluate`). Kernel match path: `internal/resolve/resolve.go::TagSet.matches` l.178 byte-compare.
- what: The record says the seam is shared for guard atoms only and that the kernel's byte-compared match path is "unchanged" — but graphlint's reachability already runs MATCH atoms through `guard.Evaluator`. Under C2 that site's `eq`/`in` over an `int` key becomes parsed while the kernel's `matches` stays byte-equal: two comparison semantics for one declared kind, in the "lint more permissive than kernel" orientation §Load-Bearing Decisions (c) calls the worst. C5's canonicalization makes the divergence unauthorable from model text (reach's held values come from load-conformed writes/initial), so the design survives, but the discovery is false as stated and C4's "the same decision the runtime evaluator makes" needs the qualification that at `atomAdmitsValue` the runtime counterpart is `TagSet.matches`, not the guard seam.

### F-2 — C3 is silent on the suite's second implementer in the kernel's own tests
- anchor: 0012:C3
- class: d
- evidence: `internal/resolve/guard_mvv_test.go::conformingContractSeam` (struct{}; `eq` arm `value == atom.Literal`), driven by `resolve.TestGuardEvaluatorContract` at l.462 and wrapped in `recordingSeam` at l.257 (REQ-71 tests).
- what: C3 extends the shared suite with int/bool want-Unevaluable legs "that fail against today's raw-string arms". The kernel package's own reference seam has exactly those raw-string arms and is `struct{}` — it cannot construct over the published fixture. Once C3 lands, 0007's REQ-71 tests fail unless `conformingContractSeam` is made kind-aware over the same kernel-owned fixture. The record names "every `GuardEvaluator` implementer" in Profile but enumerates none but `guard.Evaluator`; Phase 2 owes this site.

### F-3 — A3/C4 undercount the CLI frames: `probeRow` calls `guardSeam()` with no model in scope
- anchor: 0012:A3
- class: a
- evidence: `internal/cli/flow_next.go:492` (`Guards: guardSeam()`) inside `probeRow(row table.Row, view map[string]string, owned []resolve.Tag, observed []resolveTag, all bool)` l.468-474; model only at flow_next.go:259 (`req.observed` → `req` in scope). Also `flow_resolve.go:445` (`escapeClassOf`) resolves over a `probe` table.
- what: A3's Evidence says "`guardSeam()` is zero-arg but its callers hold `req.model`" and "each needs one added local parameter". The `flow next` probe frame holds only a `table.Row`; threading the model (or the constructed evaluator) needs the `probeRow` signature widened too, and `escapeClassOf` evaluates a filtered probe of the same model (consistent with C4's one-model rule, but unstated). The three-site enumeration is correct; the frame analysis and the "no cross-package signature surgery" cost estimate are incomplete for the CLI.

### F-4 — C3's fixture is keyed, but every suite case hardcodes `Key: "subject"`
- anchor: 0012:C3
- class: d
- evidence: `internal/resolve/guardcontract.go` l.17-23 (case struct has no key/kind), l.50-55 (`Key: "subject"` for all 12 cases).
- what: A `map[string]string` fixture types comparisons BY KEY, so one key carries one kind. The existing `eq`/`in` cases need an enum/scalar key, `gte` an int key, `contains` a set key, and the new int/bool legs their own keys. C3 must say the case table gains a per-case key drawn from the published fixture (and that the existing cases are re-keyed), otherwise the "at least one discriminating case per kind token" meta-check is unsatisfiable against a single-key table.

### F-5 — Record does not acknowledge `0007:A17`'s headline ("nor tag declarations at evaluation time")
- anchor: 0012:§approach
- class: c
- evidence: `0007:A17` headline (projected): "the seam `Evaluate(atom, value)` needs neither the view, nor provenance, nor sibling atoms, nor tag declarations at evaluation time." 0012 §Approach: "A17 verified the *signature* needs no kind precisely because construction was always the intended delivery path"; C2 illustrative arm `switch ev.kinds[atom.Key]`.
- what: 0012 grounds the chosen carrier on A17's Evidence sentence while the verified assumption's headline says the opposite of what C2 does (a declaration lookup at evaluation time, keyed by the atom). The reading "A17 froze the signature, not the constructor" may be right, but it is a reinterpretation of a Verified peer assumption, not its text. Record it the way A2 records the 0003 zero-field proxy: a named deviation against `0007:A17`'s headline, so a 7.1 cluster pass does not surface it as a cross-record contradiction.

### F-6 — Sibling-path check omits `guard.Conforms`, the existing (kind-blind, uncalled) view-conformance seam
- anchor: 0012:§research-findings (Sibling-path check)
- class: b
- evidence: `internal/guard/declaration.go::Conforms(m *table.Model, v View) error` l.280-300 — checks required-key presence and single-valuedness only; zero non-test callers (`rg '\bConforms\('` non-test → definition only).
- what: The record says "searched, no second kind-lookup path exists" and names `DeclarationOf` as the only sibling. `guard.Conforms` is the 0003 view-conformance premise (its doc: "this RDR promises nothing about a view that violates the declarations") and the natural venue for JDR 0001 §JD-18 — it already has the model and the view, returns an error, and mints no refusal kind, so it is an O4 variant the matrix never weighed. It does not decide kind today, so it does not duplicate C1/C2 and does not reopen the choice; but the Sibling-path check and the JD-18 note should name it so a later §JD-18 answer lands there rather than in a fresh path.

### F-7 — `flow next` is an unnamed seam consumer whose `unknown` output changes
- anchor: 0012:F1
- class: d
- evidence: `internal/cli/flow_next.go:492` (`guardSeam()` in `probeRow`); flow_next.go:545 "This is the ONLY source of `uncomparable`"; flow_next.go:435 "same `guardSeam()` `flow resolve` hands the kernel".
- what: Consequences and F1 describe the visible break for `flow resolve` only. `flow next` derives its `{key, uncomparable}` pairs from the probe's `guard_unevaluable` payload, so a malformed owned int now produces `uncomparable` there too (and a `"07"` owned value flips a row from excluded to candidate under parsed comparison). That is the intended honesty, but it is a user-visible output change on a second verb governed by RDR 0011's closed-reason contract and the record should name it (Consequences, F1, and S8's venue).

### F-8 — C5 names a "write-block category" that does not exist
- anchor: 0012:C5
- class: a
- evidence: `internal/table/normalize.go::renderWrites` l.624-625 refuses under `CatMalformedTagDeclaration` (`"malformed_tag_declaration"`, category.go:28, deviation D3); `[emit]` refuses under `CatEmitValueOutOfDomain` (`"emit_value_out_of_domain"`, load.go:508).
- what: C5's "reuses the refusal code each ingress already carries" is right in substance, but its list says "the write-block and emit-value categories". The write-block ingress files kind/domain refusals under the tag-DECLARATION category; a reader implementing C5 from the fence would look for a write category that is not in `category.go`. Spell the two codes as the wire does, per the project's help-body convention.

### F-9 — A4 anchor misattributes `valueAssignments` to product.go
- anchor: 0012:A4
- class: a
- evidence: `internal/guard/assignment.go:273 func valueAssignments(d table.TagDecl)`; product.go:194/336/450/477 are call sites (`valueAssignments(m.Tags[key])`).
- what: The Evidence text reads "`internal/guard/product.go::valueAssignments` ranges over `m.Tags[key]` and `internal/guard/assignment.go`'s `int` arm renders …" — the function is defined in assignment.go and takes a `TagDecl`; product.go passes `m.Tags[key]`. The projector resolved the anchor (symbol exists in the package), but the `path::Symbol` is wrong. Cosmetic; the int-arm claim itself is confirmed (l.308 `strconv.Itoa(n)`).

### F-10 — S4's mechanical enforcer is narrower than the universal it enforces
- anchor: 0012:S4
- class: d
- evidence: zero-value forms on `main` other than the literal: `var ev guard.Evaluator` (guard_grammar_0003_test.go:31/177/426, guard_scope_0003_test.go:297/353, guard_declaration_0003_test.go:267); Go also admits `new(guard.Evaluator)` and a zero `Evaluator` field in a struct literal.
- what: S4 states "no `Evaluator{}` literal survives in non-test code" as "the mechanical enforcer for C4's universal". A `var`/`new`/embedded zero value is not that literal and constructs a nil-mapping evaluator just the same. F5 makes the miss loud, so this is not silent, but the scenario should either name the grep pattern it actually runs (any zero-value construction of `guard.Evaluator` outside `NewEvaluator`) or say plainly that the loud disposition, not the grep, is the guarantee.
