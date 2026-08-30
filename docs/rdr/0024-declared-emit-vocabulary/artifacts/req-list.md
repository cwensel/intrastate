# REQ List — RDR 0024 Declared Emit Vocabulary

Phase 0 audit artifact. Every clause below is a testable obligation drawn from
`docs/rdr/0024-declared-emit-vocabulary.md`. Quotes are verbatim — copied from
the projector (`rdr inspect --select <id>`) for fenced elements and read from
the record for testable prose outside the fences — never transcribed by hand.

Element ids (`0024:C1` … `0024:C4`, `0024:MVV`, `0024:S1`…`0024:S9`,
`0024:F1`…`0024:F3`, `0024:D-naming`, `0024:G-cross-cutting`) are carried
wherever a REQ derives from a labelled element, so a later stage can trace the
REQ back to its contract.

**Shape.** `counts.elements` gives `C=4, MVV=1, A=9, F=3, D=3, S=9, G=1, BR=3`.
The four `normative` fences (C1–C4) are the spine and carry the majority of the
obligations, but each fence carries many independent clauses, and a material
remainder is UNFENCED prose: `Approach` (the enforcement-point argument and the
"declaration constrains authoring, never evaluation" clause), `Technical
Design` (the data-flow / single-reader clause), `Load-Bearing Decisions`
(identity, naming, per-member selection), the `disposition` and `fidelity`
mini-check tables, `Trade-offs / Failure Modes` (`0024:F1`–`F3`), the four
`Implementation Plan` phases (which carry named, non-optional edits: the
comment/doc amendments, the licensed golden diff, the second worked example),
the nine `Testing Strategy` scenarios, and `Cross-Cutting Concerns`.

The `Critical Assumptions` (A1–A9) are read as EVIDENCE, not as REQs, except
where an assumption fixes an implementation obligation a contract delegates to
it — A2 (the enumerated licensed diff), A4 (`ConformValue` reuse), and A9 (the
locator technique) are in that class and are cited at the REQs they bind.
A6 is Pending and explicitly NOT load-bearing; A7 is Pending on its closure leg
only, discharged by `0024:S1`; A8 is **Refuted**, and its consequence (no
`TestReq146` edit; scenario 4 is the sole determinism oracle) is a REQ.

Source sections are abbreviated:

- `NC` = Proposed Solution / Technical Design / Normative Contracts (fenced)
- `AP` = Proposed Solution / Approach
- `TD` = Proposed Solution / Technical Design (unfenced prose)
- `LBD` = Technical Design / Load-Bearing Decisions
- `MC` = Technical Design / Mini-checks (`disposition`, `fidelity`, `trace`)
- `IC` = Technical Design / Illustrative Code (**non-normative** — "Illustrative
  only — shapes, not fixtures"; Phase 4 warns against copying it)
- `AUDIT` = Proposed Solution / Existing Infrastructure Audit
- `CONS` = Trade-offs / Consequences
- `RISK` = Trade-offs / Risks and Mitigations
- `FM` = Trade-offs / Failure Modes (`0024:F1`–`0024:F3`)
- `PRE` = Implementation Plan / Prerequisites
- `MVV` = Implementation Plan / Minimum Viable Validation (`0024:MVV`)
- `PH1`…`PH4` = Implementation Plan / Phase 1…4
- `TS` = Validation / Testing Strategy (preamble)
- `SC-n` = Validation / Testing Strategy, numbered scenario *n* (`0024:Sn`)
- `XC` = Finalization Gate / Cross-Cutting Concerns (`0024:G-cross-cutting`)
- `CA` = Critical Assumptions

**Standing note on the surface at HEAD.** Verified while writing this list:
`internal/table/source.go::sourceDoc` declares no `emit` field (A1 holds);
`internal/table/load.go`'s step slice is `[]func() error` holding bound method
values in the order `…, l.loadTags, l.loadAccessors, l.loadDump,
l.loadContexts, l.loadInitial, l.loadTerminal, l.normalizeRules, …`;
`internal/cli/flow_resolve.go::resolvePayload` has **14 struct fields** and
**13 wire keys** with `Emit` immediately after `Gates`. The record's "fourteen
→ fifteen fields, thirteen → fourteen keys" arithmetic is therefore correct
against HEAD.

---

## A. Grammar: the `[emit]` declaration table (C1)

- [REQ-1] `0024:C1` "GRAMMAR. A model MAY declare its emit vocabulary in a top-level `[emit]` TOML table, one sub-table per emit key: `[emit.<key>]`." — (NC)

- [REQ-2] `0024:C1` "**The opt-in trigger is the COUNT of declared keys, never the presence of the `[emit]` table** (a bare table with zero sub-tables is a zero-declaration model; the tail of this contract gives the reasoning)." — (NC)

- [REQ-3] `0024:C1` "`kind` (required): one of `enum | bool | int | scalar` — RDR 0003's token spellings reused verbatim" … "`set` is excluded because an emit value is one authored string, `0010:C3`" — (NC) — the elided clause names the adjacent vocabulary, `internal/table/model.go::declaredKinds` (the record breaks that symbol across a line inside its backticks). `set` is NOT an admitted emit kind.

- [REQ-4] `0024:C1` "a domain, `enum` only, authored in either of two TOML spellings of the SAME `domain` key" — `domain = [ ... ]` — "a flat member array, no dispositions; or" `[emit.<key>.domain]` — "a sub-table whose keys are MODEL-AUTHORED disposition tokens (e.g. `route`, `stop`, `terminal` — intrastate fixes no vocabulary) and whose values are member arrays. The key's domain is the union; each member carries the one disposition it is listed under." — (NC)

- [REQ-5] `0024:C1` "**These are two spellings of one key, not two coexisting fields, so \"both present\" is not a refusal arm this contract owns** — authoring both is a duplicate-key error the TOML decoder raises before any check here runs. An implementer must not add a hand-written arm for it" — (NC) — **negative REQ**: no `both-domain-forms` arm may be implemented.

- [REQ-6] `0024:C1` "`bool` fixes the implicit domain `true | false`; `int` constrains the value to a base-10 integer literal; `scalar` is the declared-but-unvalidated escape hatch — the key is admitted, the value unconstrained. None of the three takes a `domain`, and dispositions attach only to declared enum members." — (NC)

- [REQ-7] `0024:C1` "Every kind check is a LEXICAL check on the authored string: no value is parsed into a typed representation, canonicalized, or converted anywhere downstream — evaluation and the payload carry the authored bytes (`0010:C3`)." — (NC) — the negative control is `0024:S9` (REQ-64).

- [REQ-8] `0024:C1` "The authored value is always a TOML **string**: `sourceRule.Emit` is `map[string]string` (`0010:C3`), so under `kind = \"int\"` an author writes `count = \"42\"`, and a bare `count = 42` refuses as `malformed_toml` from the decoder — not as `emit_value_out_of_domain` — exactly as it does today." — (NC; restated `MC` disposition row "Non-string TOML emit value (`verdict = 42`)", and `PH1`'s TestReq21 pre-clearing, REQ-72)

- [REQ-9] `0024:C1` "Reusing `ConformValue` (A4) settles non-canonical literals in the permissive direction: `03`, `+5` and `-0` are admitted for `int`, because `strconv.Atoi` admits them." — (NC; `XC` restates it; A4 is the evidence) — these three literals MUST NOT refuse under `kind = "int"`.

- [REQ-10] `0024:C1` "Because `domain` is two-shaped, its decoded field is the ONE place in the source schema that cannot be concretely typed: it is declared **`any`** on the source-schema struct, mirroring `internal/table/source.go`'s `sourceModel.Metadata`" — (NC) — the source-schema `domain` field MUST be `any`-typed.

- [REQ-11] `0024:C1` "**The carve-out is that field type, NOT a decoder option** — `pelletier/go-toml/v2` has no per-subtree strictness switch, and `DisallowUnknownFields` stays set for the whole document" — (NC) — **negative REQ**: no decoder-strictness toggle may be introduced.

- [REQ-12] `0024:C1` "**That declaration-level refusal is the decoder's, not this grammar's**: it arrives as `unknown_schema_field` (A1's `decodeStrict` arm), an existing `0002:C24` category and NOT one of the three this RDR registers. A test covering a `domaim` typo asserts the decoder's slug, and an implementer must not add a hand-written arm to re-report it as `malformed_emit_declaration`." — (NC; `MC` disposition row "Declaration-level key typo (`domaim = [...]`)")

- [REQ-13] `0024:C1` "Refused as `malformed_emit_declaration`: an unknown `kind` token" — (NC) — refusal arm 1.

- [REQ-14] `0024:C1` "an `enum` carrying no usable domain" — (NC) — refusal arm 2; see REQ-19 for its one-arm/three-fixtures ruling.

- [REQ-15] `0024:C1` "an empty-string member, or a duplicate member (duplicates across disposition lists included — one member, one disposition)" — (NC) — refusal arms 3 and 4.

- [REQ-16] `0024:C1` "a `domain` on a non-enum kind" — (NC) — refusal arm 5.

- [REQ-17] `0024:C1` "a disposition token that is the empty string" — (NC) — refusal arm 6.

- [REQ-18] `0024:C1` "the arms that exist only because strictness cannot reach them — a `domain` that is neither a flat array of strings nor a table of disposition keys, a disposition whose value is not an array of strings, any nesting below the disposition level, and a non-string member." — (NC; `MC` disposition row "Malformed shape INSIDE `domain`") — refusal arms 7–10. These four are the arms A7's Pending closure leg covers, discharged by `0024:S1` (REQ-53).

- [REQ-19] `0024:C1` "**\"No usable domain\" is ONE arm, not two, because the decoder cannot tell the two authorings apart.** … an `enum` declaring no `domain` key at all and an `[emit.<key>.domain]` sub-table carrying zero keys BOTH decode to `Domain == nil`, with no error and nothing to discriminate on. … An empty flat `domain = []` IS distinguishable (it decodes to a non-nil empty slice) and is covered by the same arm" — (NC; `SC-1` fixes the three-fixtures-one-expectation test shape, REQ-54) — all three authorings refuse as one category, and no test may assert distinct messages for the first two.

- [REQ-20] `0024:C1` "Emit declarations are NOT tag declarations: no provenance, no accessor reference, no `min/max/elements/single_valued/required`, and an emit key remains barred from match, guard, write, and accessor use (`0010:C3`)." — (NC; `AUDIT` "Extend the pattern, not the type") — **negative REQ**: the declaration grammar admits none of those keys, and `internal/table/emit_0010_test.go::TestReq27_AnEmitKeyIsNotATagKey` keeps passing unchanged (`PH1`).

- [REQ-21] `0024:C1` "A **bare `[emit]` table with zero sub-tables is a zero-declaration model**, identical in every observable to omitting the table: it takes C2's opt-in leg, no refusal becomes reachable, and the payload carries `dispositions: {}`." — (NC)

- [REQ-22] `0024:C1` "A declaration is **author-owned and unversioned**: a domain may be widened or narrowed by editing the model, and intrastate holds no history to check the edit against — every check in this contract reads the model file against itself at one point in time." — (NC; `XC` Versioning; `F1` "every check reads only the model file against itself (never the consumer's command surface)") — **negative REQ**: no narrowing detection, no history, no external command-surface check is implemented.

## B. Proof: the load-pipeline refusals (C2)

- [REQ-23] `0024:C2` "PROOF (opt-in, whole-model). With ZERO `[emit.*]` declarations the load pipeline is byte-for-byte today's: no new refusal is reachable." — (NC; `SC-3`; `CONS`; `XC` Incremental adoption)

- [REQ-24] `0024:C2` "With ONE OR MORE declarations present, the source-to-candidate-rows pipeline (`0002:C24`'s \"load\") MUST refuse, before yielding rows. \"Before yielding rows\" fixes the position: rows are minted in `internal/table/load.go`'s `normalizeRules` step, so both checks run as a step ahead of it, beside `loadTags`, and therefore read the SOURCE rules (`sourceRule.Emit`) rather than normalized rows." — (NC)

- [REQ-25] `0024:C2` "**Two steps, in this order, at this position.** The work is `loadEmitDecls` (C1's grammar checks over `[emit]`, building C3's carrier) then `checkRuleEmit` (this contract's two cross-checks over every `sourceRule.Emit`). `loadEmitDecls` MUST precede `checkRuleEmit` … and both are inserted immediately after `loadTags` in `load.go::run`'s step slice." — (NC)

- [REQ-26] `0024:C2` "**Both are methods on `*loader` taking no arguments and returning `error`** — `func (l *loader) loadEmitDecls() error` and `func (l *loader) checkRuleEmit() error`" — (NC) — the signatures are normative, because the step slice holds bound method values.

- [REQ-27] `0024:C2` "**The zero-declaration opt-in gate is evaluated INSIDE each step**, not by the caller: the step slice is uniform and has no room for a conditional call, so both steps run unconditionally and return `nil` early when no key is declared." — (NC)

- [REQ-28] `0024:C2` "`unknown_emit_key` — any `[rule.emit]` key of any rule, ordinary or escape, not declared under `[emit]`. Strictness is whole-model, not per-key" — (NC; `SC-2` asserts both rule classes)

- [REQ-29] `0024:C2` "`emit_value_out_of_domain` — an authored value that is not a member of its key's declared enum domain, not a `bool` token, or not an `int` literal, per C1's kinds (`scalar` values are never refused)." — (NC; `MC` disposition row "Value under `kind = \"scalar\"` … never refused")

- [REQ-30] `0024:C2` "**All three categories MUST carry a source line**, stamped through `internal/table/category.go::atLine` the way `loadTags` stamps a tag refusal." — (NC; `F1` "offending file in `locator`"; `MVV` step 2 "each carrying the offending block's SOURCE LINE, not `:1`"; A9 Verified)

- [REQ-31] `0024:C2` "**Two keying strategies are required, not one**: a declaration defect keys on the top-level `[emit.<key>]` header, which `tagHeaderLine`'s technique reaches unchanged" — (NC; `PH1`)

- [REQ-32] `0024:C2` "a rule-side defect (`unknown_emit_key`, `emit_value_out_of_domain`) keys on the offending RULE ID, **not on its `[rule.emit]` header**. … The rule-side locator anchors on the rule's `id = \"<ruleID>\"` line and scans forward to the emit block … The anchor must match the id's VALUE, since `[model]` also carries an `id` key." — (NC; `PH1` names the helper `emitHeaderLine(src []byte, ruleID string) int` as "a rule-id-anchored, block-bounded forward scan"; A9)

- [REQ-33] `0024:C2` "The decoder is not an alternative route — `pelletier/go-toml/v2` exposes position only on `DecodeError` with unexported fields" — (NC) — **negative REQ**: no decoder-position plumbing.

- [REQ-34] `0024:C2` "These two categories and C1's `malformed_emit_declaration` join `0002:C24`'s data-level set … and therefore map to `flow-model-invalid` under `0005:C1`, and to a blocking finding under `intrastate lint` via its load-refusal arm (`internal/cli/lint.go`)." — (NC; `AP` "The Problem Statement's predicate holds verbatim")

- [REQ-35] `0024:C2` "**The envelope code differs by surface and neither is this RDR's to change**: the `flow` verbs emit `flow-model-invalid` … while `intrastate lint` emits `model-invalid` … The category slug travels in the inner finding's `code` on both" — (NC) — an assertion on the OUTER envelope code picks the one belonging to the surface under test.

- [REQ-36] `0024:C2` "(`lint.go`'s help text naming `codeModelInvalid` for a branch it does not emit is pre-existing drift, noted so an implementer does not \"fix\" the emitted value to match the prose and break a consumer.)" — (NC) — **negative REQ**: do not touch that drift.

- [REQ-37] `0024:C2` "All three categories are refusals — nonzero exit, never advisory, under every surface that loads the model." — (NC; `MC` disposition table; the `0006:C17` closed-advisory-tier argument in `AP` and `RISK`) — **negative REQ**: no advisory-tier finding may be added by this RDR.

- [REQ-38] `0024:C2` "**Reporting is FAIL-FAST, one refusal per run — inherited from the load tier, not decided here.** … `internal/table/reserved_key_0008_test.go` actively forbids the alternative, failing any load whose error carries `Unwrap() []error`" — (NC; `MVV` step 2 "Do NOT expect both in one run"; `SC-1`/`SC-2` "one blocking finding per run")

- [REQ-39] `0024:C2` "This RDR fixes no precedence among the three emit categories and the existing ones, and takes none: order stays unspecified." — (NC; `SC-2` "the test does NOT assert WHICH") — **negative REQ**: no test may pin which of two coexisting defects is reported.

- [REQ-40] `0024:C2` "`0005:C1`'s \"one findings[] entry per category hit\" is satisfied vacuously, because no load path yields more than one hit to map — this RDR introduces no accumulation and takes no `Overrides` entry against it." — (NC) — **negative REQ**: no `Overrides` entry, no accumulation.

- [REQ-41] `0024:C2` "Proving the AUTHORED form is proving the executed form: `0010:C3` fixes that normalization carries the emit block through `expand` unmutated … Nothing is checked at resolve time, and an emit value's evaluation semantics — uninterpreted, byte-compared — are unchanged (`0010:C3`)." — (NC; `AP` "a declaration constrains authoring, never evaluation"; `ALT3` rejected) — **negative REQ**: no resolve-time validation.

- [REQ-42] `0024:C2` "A declared key no rule emits, and a declared member no rule authors, are NOT findings of any tier (authoring headroom; the advisory tier is closed, `0006:C17`)" — (NC; `MC` disposition row; `CONS`) — **negative REQ**.

- [REQ-43] `AP` "The enforcement point is the load pipeline, not lint … No graphlint code changes; the advisory tier stays closed (`0006:C17`) and the blocking tier's \"at least\" floor (`0006:C3`) is not touched." — (AP; `TD` "Technical Environment" lists `internal/graphlint/` as unchanged) — **negative REQ**: `internal/graphlint` gains no code change.

## C. Carry: the normalized declaration carrier (C3)

- [REQ-44] `0024:C3` "CARRY. Every declared field — key, kind, domain members, and each member's disposition — is carried through normalization onto the normalized model, the same clause `0002:C22` states for tag declarations" — (NC; `MC` fidelity row 1)

- [REQ-45] `0024:C3` "The carried form is a new model-level type, deliberately NOT `TagDecl` (the mirror of `EmitValue` not being `TagValue`, `0010:C3`). It is carried at `Model.EmitDecls map[string]EmitDecl`, keyed by emit key — the sibling of `Model.Tags map[string]TagDecl`." — (NC)

- [REQ-46] `0024:C3` "`EmitDecl` is an exported struct with exactly three exported fields: `Kind string` (one of `enum | bool | int | scalar`), `Domain []string` (enum only; the bytewise-sorted union; `nil` for every other kind), and `Dispositions map[string]string` (member → its disposition token; `nil` when the domain carries none)." — (NC) — exactly three exported fields.

- [REQ-47] `0024:C3` "**`Model.EmitDecls` is always non-nil after a successful load** — empty for a zero-declaration model, never `nil` — matching the `Model.Tags` convention" — (NC; `XC` Canonical-form)

- [REQ-48] `0024:C3` "`Domain` is `nil` (not an empty non-nil slice) for `bool`, `int`, and `scalar`, since none of them takes a domain; the distinction is observable because scenario 4 asserts value-equality on the carrier, where `nil` and `[]string{}` do not compare equal." — (NC; `SC-4`)

- [REQ-49] `0024:C3` "**The carry is value-preserving, not order-preserving, and the member list is SORTED.** … So `Domain` is the union sorted bytewise, and the partition grouping is not itself carried: it is fully recoverable from `Dispositions`" — (NC; `MC` fidelity row 1 lossy-exemption column)

- [REQ-50] `0024:C3` "that sweep does NOT reach this carrier: all three of its sub-tests are scoped to `table.Row` (A8, Refuted), so **Testing Strategy scenario 4 is the sole oracle for this clause** and Phase 2 owes no edit to `TestReq146`." — (NC; `PRE` A8 item; `PH2`; `SC-4`; `XC`) — **negative REQ**: `internal/table/roundtrip_test.go::TestReq146_EveryEmittedSequenceIsASortedSlice` is NOT edited.

- [REQ-51] `0024:C3` "the authored order of a flat-array domain is NOT preserved either, and nothing downstream reads it" — (NC)

- [REQ-52] `0024:C3` "**Value conformance is performed by constructing a throwaway `TagDecl{Kind: d.Kind, Domain: d.Domain}` per call and passing it to `internal/table/load.go::ConformValue(decl TagDecl, member string) error`.** That struct literal IS the adapter — no extraction, no new exported helper." — (NC; A4 Verified; `AUDIT` "REUSE AS-IS"; `PH1` "no `internal/guard` import") — **negative REQs**: no `internal/guard` import, no extraction, no new exported conformance helper.

- [REQ-53] `0024:C3` "`ConformValue` reads only `Kind`, `Domain`, `Min`, `Max`, and `Elements`; `Min`/`Max`/`Elements` stay nil, so the `int` arm's bounds and the `set` arm never fire." — (NC)

- [REQ-54] `0024:C3` "**A `scalar`-kinded key is short-circuited BEFORE the call** rather than relying on `conformKind`'s fall-through" — (NC) — an explicit skip, not an inherited fall-through.

- [REQ-55] `0024:C3` "**The reuse is safe-by-omission, so C2's value refusal DEPENDS on C1's arms having already fired.** … **if C1's arms are ever relaxed, `emit_value_out_of_domain` silently stops firing**" — (NC) — the C1-before-C2 ordering (REQ-25) is load-bearing for correctness, not only for the carrier.

- [REQ-56] `0024:C3` "The kernel (`internal/resolve`) continues to carry no declarations and no emit. This RDR adds exactly one reader of the carried declarations — C4's payload join — and `internal/graphlint` reads none of it here" — (NC; `TD` data-flow "This RDR adds exactly one downstream reader"; `MC` fidelity row 4, cell text "declarations are carried by neither" for the operation "Normalized model → kernel / dump"; `SC-4` "Kernel and dump surfaces carry none of it"; `XC` Concurrency) — **negative REQ**.

## D. Envelope: the `dispositions` payload field (C4)

- [REQ-57] `0024:C4` "ENVELOPE. The `flow resolve` success payload gains `dispositions`: a JSON object mapping emit key → the disposition token the declaration assigns the selected row's authored value, keys in byte order, present as `{}` — never `null`, never omitted — when no selected value carries one (undeclared model, non-enum kind, flat-array domain, or empty emit block alike)." — (NC; `SC-5`)

- [REQ-58] `0024:C4` "The field is `Dispositions map[string]string` with tag `json:\"dispositions\"` and no `omitempty`." — (NC; `XC` Canonical-form)

- [REQ-59] `0024:C4` "**Byte order is the marshaller's, not the join's**: `encoding/json` sorts map keys on marshal, so the join performs no sort and the carried value is an ordinary unordered Go map … An implementer must not add a sort over this map" — (NC; `XC` Canonical-form (ii)) — **negative REQ**.

- [REQ-60] `0024:C4` "**The map is keyed off the SELECTED ROW's authored emit, never off the declaration set.** An entry exists for key `k` exactly when the selected row authors `k` AND `k`'s declaration lists that authored value under a disposition. A declared key the selected row does not emit contributes NO entry — `dispositions` is never padded with nulls, empty strings, or absent-markers to the declared key set." — (NC; `SC-5` "a model declaring two keys whose selected row emits one yields exactly one `dispositions` entry")

- [REQ-61] `0024:C4` "It is inserted at one fixed position in `internal/cli/flow_resolve.go::resolvePayload`, immediately after `emit`" — (NC; `MVV` step 3; `SC-5`)

- [REQ-62] `0024:C4` "**The join itself is performed in `internal/cli`, reading `Model.EmitDecls`; `internal/table` gains no payload-shaped helper.** … a method on `*Model` returning the payload's map would shape the table package's API around the CLI's wire format" — (NC) — **negative REQ**: no `*Model` method returning the payload map.

- [REQ-63] `0024:C4` "**`omitempty` is available on this field and is deliberately not taken**" — (NC; `XC` Canonical-form) — **negative REQ**.

- [REQ-64] `0024:C4` "The token is surfaced verbatim; intrastate never interprets it." — (NC; `AP` "intrastate never interprets a disposition token — domains and dispositions are model-authored vocabulary, carried through and surfaced verbatim"; `MC` fidelity row 3)

- [REQ-65] `0024:C4` "Text mode renders through the generic payload renderer as `dispositions.<key>: <token>` / `dispositions: (none)` with no per-verb special case (the `0010:C4` clause)." — (NC; `SC-5`) — **negative REQ**: no per-verb text special case.

- [REQ-66] `0024:C4` "A plan rescued by an escape row joins `dispositions` from that escape row's OWN authored values — the same single `Plan.RuleID` join path as `emit` (`0010:C4`), so no defaulted or merged value can reach the payload unjoined" — (NC; `SC-5` "an escape-row rescue")

- [REQ-67] `0024:C4` "`flow next` carries no `dispositions`, for `0010:C4`'s reason: the answer is what `resolve` selects." — (NC; `SC-5`) — **negative REQ**.

- [REQ-68] `0024:C4` "The insertion takes `resolvePayload` from fourteen struct fields to **fifteen** and the wire from thirteen keys to **fourteen**, with `dispositions` at index 9 (`escape_class` is `omitempty` and absent on a non-escaped plan). Both values are normative fixtures" — (NC; `SC-5` carries the key list, REQ-83)

## E. Load-bearing decisions and naming

- [REQ-69] `0024:D-identity` "an emit declaration is identified by its emit key, byte-exact (the same identity emit keys already have, `0010:C3`); one declaration per key, duplicates refused by the TOML decoder as the tag table's are. A domain member's identity is its byte-exact string; one disposition per member." — (LBD)

- [REQ-70] `0024:D-naming` "the table is `[emit]`, mirroring `[tags]` / `[rule.emit]`" and "The new load categories are `malformed_emit_declaration`, `unknown_emit_key`, and `emit_value_out_of_domain`" and "The payload field is `dispositions`" — (LBD) — the three slugs and the field name are fixed literals.

- [REQ-71] `0024:D-naming` "Each slug also gets its `Cat*` constant in `internal/table/category.go` and an entry in `Categories()`, which is what a consumer branches on — the wire slug and the registered category are one decision, not two, and scenario 6 asserts the registration rather than trusting the append." — (LBD; `SC-6`)

- [REQ-72] `0024:D-selection-predicate` "dispositions are per **domain member**, not per key … When the selected row authors a value for a declared enum key, the disposition surfaced is the one the declaration lists that member under — exactly one exists by C1's duplicate refusal." — (LBD)

## F. Failure modes

- [REQ-73] `0024:F1` "a declared model with an undeclared key, out-of-domain value, or malformed declaration refuses at load — `intrastate lint` exits nonzero with one blocking finding — the first refusal the fail-fast pipeline reaches (category slug in the finding's `code`, offending file in `locator`) — and `flow resolve`/`flow next` refuse `flow-model-invalid` with the same finding." — (FM)

- [REQ-74] `0024:F2` "a `scalar`-declared key admits any value — declared-but-unvalidated is the documented escape hatch, visible in the model source. A model with zero declarations is today's world: nothing new fires." — (FM)

- [REQ-75] `0024:F3` "a model carrying `[emit]` refuses on an older binary as `unknown_schema_field` (strict decode) — loud, never a silent ignore of the declaration." — (FM; A1; `XC` Versioning) — satisfied by construction (`decodeStrict` at HEAD); no new code, but no code may weaken it.

## G. Implementation-plan obligations (named, non-optional edits)

- [REQ-76] `PH1` "Stamp all three categories through `atLine`, with the two keying strategies C2 fixes. … Write `emitHeaderLine(src []byte, ruleID string) int` as a rule-id-anchored, block-bounded forward scan: find the `id = \"<ruleID>\"` line (the id's VALUE — `[model]` also carries an `id` key), then scan to the emit block." — (PH1; REQ-30/31/32 are the contract legs)

- [REQ-77] `PH1` "**Amend three stale comments in the same change.** … `internal/table/model.go::EmitValue` (\"undeclared, uninterpreted, and compared by exact byte equality\"), `internal/table/normalize.go::emitSequence` (\"they are undeclared and uninterpreted\"), and `internal/table/dump.go::renderEmit` (\"an emit key has neither: it is undeclared, so there is no kind to consult\"). … only the \"no kind to consult\" and \"uninterpreted\" wordings narrow to \"no TAG kind\" and \"never parsed, canonicalized, or converted\"" — (PH1)

- [REQ-78] `PH1` "`TestReq27`'s ASSERTIONS need no edit, but its failure-message prose restates \"undeclared\" in the same absolute sense the three comments do — narrow those strings to \"not a tag key\" in this same change" — (PH1) — assertions unchanged, message strings narrowed.

- [REQ-79] `PH1` "`internal/table/emit_shape_0010_test.go::TestReq21_NoHandWrittenTypeCheckDiscriminatesTheEmitRefusal` … keeps passing … Do not \"fix\" TestReq21 when adding the kind arm; if it goes red, the kind check has been placed above the decoder, which C1 does not license." — (PH1) — **negative REQ**; TestReq21 is the tripwire.

- [REQ-80] `PH3` "(i) the 28 inline whole-payload assertions in `internal/cli/flow_demand_0011_test.go` — **all 28** pin the `\"emit\":{…},\"next\"` adjacency, so every one of them goes red and that is the licensed outcome" and "(ii) those 28 ARE the 3 `resolveGolden` byte-identity goldens plus the 25-entry `shippedResolveGoldens` sweep in the same file, not a further set." — (PH3; A2)

- [REQ-81] `PH3` "Two obligations ride with that authorization: the regenerated goldens must be captured from a build carrying ONLY the `dispositions` append (never from a tree with other pending work …), and the diff must be inspected to confirm every one of the 28 changed in exactly one way — the inserted `dispositions` key at index 9" — (PH3; `RISK` "**The residual risk is not diff width but evidentiary loss**")

- [REQ-82] `PH3` "(iii) `internal/cli/decision_table_0010_test.go::TestReq39_EmitSitsImmediatelyAfterGatesOnTheWire`, whose 13-key `slices.Equal` list becomes 14 keys and whose `NumField() != 14` guard becomes `!= 15` … the edit raises both it and the adjacent failure message (which names \"thirteen to fourteen\" in prose, and the sub-test title \"the struct declares fourteen fields\"). Those prose strings are part of this licensed edit" — (PH3; A2)

- [REQ-83] `PH3` "`TestReq43` (escaped between clear and escape_class) asserts relative indices and needs no edit. Any assertion outside this list going red is a signal the append did more than C4 licenses — investigate, do not update it." — (PH3; `SC-3` "A red outside that enumerated set is the real signal") — **negative REQ**: the licensed diff is closed.

- [REQ-84] `PH3` "**The §D1 registration obligation, whichever record lands second.** … `dispositions` is assigned PLAN there by name. That oracle does not exist at HEAD — 0023 is Final and unimplemented — so: if 0023 lands first, this change registers `dispositions` with the oracle it shipped; if 0024 lands first, 0023's implementation registers it when the oracle arrives." — (PH3; `0024:JC1`; `SC-7`)

- [REQ-85] `PH4` "Declare the pricing example's emit keys (`models/examples/pricing-decision-table.toml`), and extend `docs/cli-output-contract.md`'s resolve section with `dispositions` and the declaration grammar" — (PH4; `SC-8`)

- [REQ-86] `PH4` "**Declare the example against the values it actually authors**, not against the Illustrative Code's shapes: the model emits `plan` ∈ {`basic`, `pro`} and `dpa` ∈ {`required`, `none`} across its four rows. The Illustrative Code block shows `[emit.dpa]` with `domain = [\"required\", \"waived\"]` — a SHAPE illustration whose members are not this model's, and copying it refuses two of the four rows" — (PH4; `SC-8`; `IC` is non-normative)

- [REQ-87] `PH4` "this phase adds a second example — a small routing table whose one emit key partitions into `route` and `stop` members — under `models/examples/`, and `docs/cli-output-contract.md` documents `dispositions` against it." — (PH4; `SC-8` "Scenario 8 asserts both examples lint clean")

- [REQ-88] `PH4` "**Amend the two stale `docs/model-authoring.md` sentences** in the same change … The doc says emit keys are keys \"nothing declares\" (twice — in the reserved-key section and in the `[rule.emit]` walkthrough), and that `plan = \"<clear>\"` in an emit block \"loads and answers with that literal text\", which a declared domain excluding it now refuses. Narrow both to the same line C1 draws" — (PH4)

## H. Minimum Viable Validation (`0024:MVV`)

- [REQ-89] `0024:MVV` "Re-author the seed's two-rule adversarial table — first rule's emit value misspells a command, second rule types the key `nxet` — and confirm it still lints exit 0 with zero findings (the defect, reproduced)." — (MVV step 1; the negative control)

- [REQ-90] `0024:MVV` "Add an `[emit.next]` enum declaration with a route/stop-partitioned domain; run `intrastate lint`. **Expected**: exit nonzero with ONE blocking finding … Fix that rule and re-run: the second defect now surfaces, likewise as one finding. Across the two runs both categories are observed — `emit_value_out_of_domain` for the misspelled value and `unknown_emit_key` for the typo'd `nxet` key. Do NOT expect both in one run" — (MVV step 2)

- [REQ-91] `0024:MVV` "Fix both rules; lint exits 0. Run `flow resolve` selecting a rule whose value is listed under `stop`. **Expected**: payload carries the authored `emit` unchanged and `dispositions` mapping the key to `stop`, positioned immediately after `emit`." — (MVV step 3)

- [REQ-92] `0024:MVV` "Delete the `[emit]` table; re-run lint and resolve over the original (defective) table. **Expected**: exit 0, no findings, payload carries `dispositions: {}` — byte-identical behavior to today apart from the appended empty field." — (MVV step 4)

## I. Testing Strategy scenarios

- [REQ-93] `TS` "Done = every scenario below is green and the MVV run recorded; every C1 refusal arm and every C2 category maps to at least one scenario." — (TS)

- [REQ-94] `0024:S1` table-driven load of malformed declarations covering every C1 arm — "**Expected**: each case, loaded on its own, refuses `malformed_emit_declaration` at load with exactly one blocking finding, before any row is yielded. One defect per fixture — the pipeline is fail-fast, so a fixture carrying two defects proves nothing about the second." — (SC-1) — this scenario is A7's Pending closure gate; `PRE` states "Phase 1 does not close until it is green."

- [REQ-95] `0024:S1` "The no-usable-domain arm takes **three fixtures against one expectation** — no `domain` key, an empty `[emit.<key>.domain]` sub-table, and `domain = []` … Assert the shared category; do NOT assert distinct messages for the first two, which would pin a discrimination the decoder cannot make." — (SC-1; REQ-19)

- [REQ-96] `0024:S2` "**Expected**: `unknown_emit_key` on both rule classes; `emit_value_out_of_domain` for the enum/bool/int cases; `scalar` never refused — each from its own single-defect fixture, one blocking finding per run." — (SC-2)

- [REQ-97] `0024:S2` "A model carrying an emit defect AND a coexisting structural one asserts the fail-fast contract instead: exactly one finding, and the test does NOT assert WHICH — the pipeline order is unspecified and pinning it would fabricate a guarantee" — (SC-2; REQ-39)

- [REQ-98] `0024:S3` "**Expected**: full suite green with the loader change and no `models/*.toml` edited; no new refusal reachable; the only observable payload delta is C4's appended `dispositions: {}`." — (SC-3; A3's post-change leg)

- [REQ-99] `0024:S3` "**The pass criterion is LOAD-side equivalence, not payload byte-identity** … Green here means: the loader admits and refuses exactly the same models as before, and the only payload difference anywhere is the one appended field. \"Full suite green\" is measured AFTER Phase 3's licensed golden regeneration" — (SC-3)

- [REQ-100] `0024:S4` "**Expected**: read off `Model.EmitDecls[<key>]` — `Kind` equals the authored token, `Domain` equals the bytewise-sorted union of the authored members, and `Dispositions[<member>]` equals the token that member was listed under. Kernel and dump surfaces carry none of it." — (SC-4)

- [REQ-101] `0024:S4` "**Determinism is asserted here or nowhere**: loading the same partitioned-domain model repeatedly yields one `Domain` value, which is the leg that fails if the union is built from map iteration without the sort. … this assertion is the only thing standing between an unsorted union and a green suite." — (SC-4; REQ-50)

- [REQ-102] `0024:S5` "**Expected**: `dispositions` maps only member-disposition hits and is `{}` in every other case — never `null`, never omitted — sits immediately after `emit`, joins from the escape row's own authored values by `Plan.RuleID`, and renders in text mode through the generic payload renderer; `flow next` payloads carry no `dispositions`." — (SC-5)

- [REQ-103] `0024:S5` "The post-change wire key order is exactly these fourteen keys, in this order (the pre-change thirteen from `decision_table_0010_test.go::TestReq39`, with `dispositions` inserted at index 9)" — the record then gives the list in a fenced block, wrapped across two lines: `model, revision, observed, owned, readers, outcome, rule, gates, emit, dispositions, next, writes, clear, escaped` — "and `reflect.TypeOf(resolvePayload{}).NumField() == 15`" — (SC-5; the normative fixture C4 defers here; REQ-68)

- [REQ-104] `0024:S6` "**Expected**: `malformed_emit_declaration`, `unknown_emit_key` and `emit_value_out_of_domain` each appear in `table.Categories()` (`internal/table/category.go`) and each carries a checked-in `testdata/neg/*.toml` witness asserted by category, per the convention `internal/table/dump_test.go`'s REQ-119 witness map holds for every existing load category." — (SC-6; REQ-71)

- [REQ-105] `0024:S6` "**The convention has an exclusivity leg (REQ-131) this scenario must also carry**: each witness trips its own category *and no other*. Assert it per witness. … each witness must be single-defect by construction" — (SC-6)

- [REQ-106] `0024:S7` "**In the 0024-first order this scenario is not dischargeable by this record's own tests** … so it discharges instead as a written, checked-in note at the `dispositions` field declaration in `internal/cli/flow_resolve.go::resolvePayload` naming the PLAN assignment and `JDR 0002 §D1` … That note is the Done criterion in that order" — (SC-7; REQ-84)

- [REQ-107] `0024:S8` "**Expected**: `models/examples/pricing-decision-table.toml` lints exit 0 with its emit keys declared — the declared domains must cover the values the example actually authors, `plan` ∈ {`basic`, `pro`} and `dpa` ∈ {`required`, `none`} — and `docs/cli-output-contract.md` documents `dispositions`, asserted with the `readRepoFile` shape peers 0009 and 0011 already use to pin that same doc." — (SC-8; REQ-85/86/87)

- [REQ-108] `0024:S9` "**Expected**: an enum-declared emit value renders byte-identically to its undeclared rendering — `internal/table/dump.go::renderEmit` keeps bypassing `renderValue`. This is the negative control on C1's \"no value is parsed, canonicalized, or converted anywhere downstream\"" — (SC-9; REQ-7) — **negative REQ**: `renderEmit` MUST NOT be routed through `renderValue`.

## J. Cross-cutting concerns

- [REQ-109] `0024:G-cross-cutting` "**Concurrency model** — none introduced. Both new steps are `*loader` methods on the single-threaded load path, and the declaration carrier is written once during load and read-only thereafter" — (XC)

- [REQ-110] `0024:G-cross-cutting` "No hash, no content-addressed identity, and no replay-stable digest is claimed anywhere in this RDR" — (XC) — **negative REQ**.

---

## ASSUMPTIONS

- **ASSUMPTION-1 (REQ-13…REQ-18, REQ-94)** — C1's refusal list is read as **ten
  distinct arms** for test-coverage purposes: unknown `kind`; no usable domain
  (one arm, three fixtures per REQ-19/95); empty-string member; duplicate
  member; `domain` on a non-enum kind; empty-string disposition token;
  non-array/non-table `domain`; non-array disposition value; nesting below the
  disposition level; non-string member. This is the defensible reading because
  `TS` requires "every C1 refusal arm … maps to at least one scenario" and
  `SC-1` enumerates the same set in the same order. The record never numbers
  them; the numbering here is an artifact convention, not a contract claim.

- **ASSUMPTION-2 (REQ-26)** — the two step names and signatures
  (`loadEmitDecls`, `checkRuleEmit`) are read as **binding**, not illustrative,
  because C2 states them inside a `normative` fence with the reason (the step
  slice holds bound method values). An implementer who renames them satisfies
  the observable behaviour but not the fence; the defensible reading of a
  normative fence naming an exact signature is that it is normative.

- **ASSUMPTION-3 (REQ-32, REQ-76)** — `emitHeaderLine`'s exact signature comes
  from Phase 1 prose, not from a fence, and C2 explicitly says "Two helpers or
  one is an implementation choice; recovering BOTH lines is not". So the REQ is
  read as: **both lines are recoverable and stamped**; the helper's name and
  arity are the implementer's. REQ-76 records Phase 1's named shape as the
  default, not as a hard constraint.

- **ASSUMPTION-4 (REQ-30, REQ-73)** — "carries a source line" is asserted as a
  **non-zero, non-1 line number pointing at the offending block**, not as an
  exact line-number fixture, because `MVV` step 2 states the property as "the
  offending block's SOURCE LINE, not `:1`" and a literal-line fixture would
  break on any whitespace edit to the fixture model. The defensible reading is
  the property, not the integer.

- **ASSUMPTION-5 (REQ-80, REQ-81)** — the "28" count is treated as **evidence
  from A2's census, not a target to hit exactly**: the obligation is that every
  whole-payload assertion in `internal/cli/flow_demand_0011_test.go` pinning the
  `emit`/`next` adjacency is regenerated, and that the diff shows exactly one
  inserted key per site. If the census returns 27 or 29 at implementation time
  the obligation is unchanged; a materially different count is a deviation to
  record, not a REQ failure. Precedent: `0010`'s ASSUMPTION-8 treats its own
  103-fixture census the same way.

- **ASSUMPTION-6 (REQ-77, REQ-78, REQ-88)** — the comment/message/doc
  amendments are asserted by **reading the narrowed text back**, not by
  asserting exact new wording: the record fixes the *line to draw* ("emit keys
  are not TAGS … and MAY be declared") but not the sentence. Tests assert the
  absence of the falsified absolute claim where a test is warranted (`SC-8`
  already uses `readRepoFile` for `docs/cli-output-contract.md`); the code
  comments are reviewed, not asserted.

- **ASSUMPTION-7 (REQ-84, REQ-106)** — the run is assumed to be **0024-first**
  (0023 is Final and unimplemented at HEAD, per Phase 3's own statement), so
  `0024:S7` discharges via the checked-in note at the `dispositions` field
  declaration rather than via a reflective oracle. If 0023 has landed by the
  time Phase 3 runs, the alternative leg (register with the shipped oracle)
  applies and the note becomes redundant — that branch is checked at Phase 3,
  not assumed here.

- **ASSUMPTION-8 (REQ-16)** — "a `domain` on a non-enum kind" is read as
  refusing the *presence of the `domain` key* under `bool`/`int`/`scalar`, in
  either spelling, rather than only the flat-array spelling. C1's "None of the
  three takes a `domain`" is unqualified as to spelling, and the narrower
  reading would leave `[emit.<key>.domain]` under `kind = "bool"` silently
  admitted — which contradicts the arm's purpose.

- **ASSUMPTION-9 (REQ-29, REQ-96)** — "not a `bool` token" is read as
  `ConformValue`'s existing `bool` arm (RDR 0003's spellings), inherited by
  construction the same way REQ-9 inherits `strconv.Atoi`'s permissiveness for
  `int`. C1 says "`bool` fixes the implicit domain `true | false`" and C3 fixes
  that conformance runs through `ConformValue`; where the two could differ, the
  reuse clause wins, because C1's whole non-canonical-literal paragraph is an
  argument that reuse settles the lexical rules.

- **ASSUMPTION-10 (REQ-104)** — the `testdata/neg/*.toml` witness convention is
  followed as `internal/table/dump_test.go`'s REQ-119 map holds it at HEAD (one
  witness file per category, registered in that map). The record names the
  convention and the exclusivity leg but not the file layout; the defensible
  reading is "as the existing map does it", which is what `SC-6` cites.

---

## QUESTIONS

None blocking. Four clauses were candidates; each is settled by the record, its
evidence, or predecessor precedent, so none is escalated.

1. **Does `unknown_emit_key` fire per-rule or once per undeclared key?** C2
   says "any `[rule.emit]` key of any rule … not declared under `[emit]`",
   which reads either way. **Settled by fail-fast** (REQ-38): the load tier
   returns exactly one categorized error per run, so the distinction is
   unobservable — the first undeclared key reached refuses, and no accumulation
   is licensed (REQ-40). Recorded at REQ-28/REQ-38, no fork survives.

2. **Are `bool`/`int`/`scalar` keys eligible for `dispositions` entries?** C4's
   `{}`-cases list includes "non-enum kind", and C1 says "dispositions attach
   only to declared enum members" — so no. The competing reading (a non-enum
   key could carry a key-level disposition) is foreclosed by
   `0024:D-selection-predicate`: "dispositions are per **domain member**, not
   per key". Recorded at REQ-6/REQ-57/REQ-72.

3. **Does the `[emit]` grammar's `domain` sub-table admit a member listed under
   two dispositions?** C1 refuses "a duplicate member (duplicates across
   disposition lists included — one member, one disposition)" — explicit.
   The only residual reading was whether the same member repeated *within one*
   disposition list is also a duplicate; "a duplicate member" is unqualified and
   `0024:D-identity` fixes "one disposition per member" with member identity as
   the byte-exact string, so both cases refuse. Recorded at REQ-15/REQ-69.

4. **Must the `dispositions` join run for `flow next`?** REQ-67 is explicit
   ("`flow next` carries no `dispositions`"), and `SC-5` asserts it. The
   competing reading — that `next`'s payload should stay shape-parallel to
   `resolve`'s — is foreclosed by the cited `0010:C4` reason ("the answer is
   what `resolve` selects"), which already governs `emit` itself. Recorded at
   REQ-67.
