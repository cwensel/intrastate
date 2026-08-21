Model: claude-opus-5[1m]

# Cove Step 0 — grounding sweep (iteration 2, delta-scoped)

Scope: the seven assumptions named in the Status re-verify line (A3, A16, A17,
A18, A19, A21, A22) plus every Normative Contracts / Technical Design /
Implementation Plan / Testing Strategy passage citing them. Iteration 1's
findings (`../findings.md`) were read only to avoid duplicate reporting; every
verdict below was re-derived from source in this pass.

Sources read this pass: `internal/resolve/resolve.go`,
`internal/resolve/adversarial_test.go`, `internal/resolve/fixup_test.go`,
`internal/resolve/fixtures_test.go`, `internal/resolve/resolve_test.go`,
`internal/cli/root.go`, `internal/cli/clierr/clierr.go`,
`docs/jdr/0001-resolve-kernel-seam.md`, RDR 0002/0003/0005/0009,
`docs/rdr/0007-guard-predicate-totality/evidence/spikes/a3-reshape-probe.md`.

---

## Per-claim verdicts

### A3 — bounded kernel change; zero production importers

| # | Claim | Verdict | Cite |
| --- | --- | --- | --- |
| A3.1 | "the sweep was re-run this stage and is still CLEAN" — zero production importers of `internal/resolve` | **CONFIRMED (re-run by me this pass)** | `rg -n -t go --glob '!internal/resolve/**' 'newcoinc/intrastate/internal/resolve' .` → no output, exit=1. Broader sweep `rg --glob '!internal/resolve/**' --glob '!docs/**' 'internal/resolve' .` → no output. Repo-wide `.go` inventory: only `cmd/intrastate/main.go`, `internal/cli/{root,version,version_test,clierr,config,respond}`, `internal/version`, and `internal/resolve/*` itself. |
| A3.2 | "`internal/cli/root.go::NewRootCmd` registers only `newVersionCmd()`" | **CONFIRMED** | `internal/cli/root.go` `func NewRootCmd()`: single `cmd.AddCommand(newVersionCmd())`. No other `AddCommand` in the file. |
| A3.3 | Touch-list `Row`, `GuardEvaluator`, `evaluateGuard`, `Refusal` all resolve | **CONFIRMED** | `resolve.go:169 type Row struct` (with `Guard string` at :185, `RequiresOwned []string` at :182); `:91 type GuardEvaluator interface` with `Evaluate(guard string, view TagSet) GuardResult`; `:425 func evaluateGuard(seam GuardEvaluator, guard string, view TagSet) GuardResult`; `:255 type Refusal struct` with `Guard string` at :272. |
| A3.4 | "`gate`'s prune → owned-state → undecidable partition is byte-identical"; `gate` never reads `row.Guard` except to copy it into `Refusal.Guard` | **CONFIRMED** | `resolve.go:376 func gate(...)`: prune loop `if verdict == GuardFalse { continue }`; then `if missing := missingOwned(survivors, view); len(missing) > 0 { return nil, &Refusal{Kind: KindOwnedStateUnavailable, ...} }`; then the `undecidable` loop returning `&Refusal{Kind: KindGuardUnevaluable, Guard: lowest.Guard, Rows: rowRefs(undecidable)}`. `lowest.Guard` is the single `row.Guard` read, and it is exactly the `Refusal.Guard` copy. |
| A3.5 | "`escapeOrRefuse` needed ZERO changes" / uniform gating | **CONFIRMED (as a code fact)** | `resolve.go:473 func escapeOrRefuse(in Input, view TagSet, r Refusal) Result` — doc comment: "Escape candidates pass the same viability gate as ordinary candidates before the exact-one count is taken". The spike record's diff claim is not independently re-runnable here, but the shipped shape it asserts is present. |
| A3.6 | Six frozen tests named; all exist and pass | **CONFIRMED** | `adversarial_test.go:48 TestAdv1_GuardFalseRowsRequiresOwnedMustNotPoisonAnExactOneMatch`; `:114 TestAdv1b_AllGuardsFalseIsNoMatchNotOwnedStateUnavailable`; `:169 TestAdv2_EscapeEdgeMustNotBypassTheGuardSeam`; `:266 TestAdv3_GuardUnevaluableRefusalMustNotDependOnTableRowOrder`; `fixup_test.go:42 TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue`; `fixup_test.go:249 TestFixupGateIsUniformAcrossOrdinaryAndEscapeCandidates`. |
| A3.7 | "Baseline 154 PASS / 0 FAIL" | **CONFIRMED** | `go test ./internal/resolve/ -v` this pass: `grep -c -- '--- PASS'` = **154** (56 top-level + 83 subtests + nested), `--- FAIL` = **0**, package line `ok github.com/newcoinc/intrastate/internal/resolve`. |
| A3.8 | "the only two sites reading `Refusal.Guard`" are Fixup-1d (compares against the value the test set) and `TestReq33…` (non-emptiness check) | **CONFIRMED** | `resolve_test.go:1067 TestReq33_UnavailableOwnedStateAndUnevaluableGuardAreValueRefusals` exists; `fixup_test.go:42` exists. No third production/test site reads `Refusal.Guard` outside `resolve.go:413`'s construction. |
| A3.9 | Fixtures are keyed on guard *text* and discard the view | **CONFIRMED** | `fixtures_test.go:20 func (g fixtureGuards) Evaluate(guard string, _ resolve.TagSet) resolve.GuardResult` — the `TagSet` parameter is `_`. |
| A3.10 | New names `GuardAtom`, `evaluateAtom`, `undecidedAtoms`, `compareAtoms`, `copyAtoms` | **NOT-FOUND in shipped source (correctly — they are the RDR's proposed additions)** | `rg 'GuardAtom\|evaluateAtom\|undecidedAtoms\|compareAtoms\|copyAtoms' internal/` → no hits. `resolve.go` has `compareRefs` (:528) and `copyTags` (:569) as the existing analogues. No defect; recorded so a reader does not expect them today. |
| A3.11 | Mutation probe file `evidence/spikes/a3-reshape/spike_kleene_test.go.txt` exists | **CONFIRMED** | File present, 3584 bytes, dated 2026-08-21. |

### A16 — kernel recognizes the existence operator by one token

| # | Claim | Verdict | Cite |
| --- | --- | --- | --- |
| A16.1 | RDR 0003's matrix spells presence as exactly ONE token, the `exists` row, "Tests presence or absence, not value equality", literal shape `boolean` | **CONFIRMED verbatim** | `docs/rdr/0003-guard-predicate-exhaustiveness.md:239` — `| `exists` | optional scalar or optional set-valued tag | boolean | Tests presence or absence, not value equality. |` |
| A16.2 | "there is no `absent` counterpart token" | **CONFIRMED** | The matrix (0003 L235–240) has exactly five operator rows: `eq`, `in`, `lt/lte/gt/gte`, `exists`, `contains`. No `absent`. |
| A16.3 | The other four rows are stated purely as value comparisons | **CONFIRMED** | 0003 L236–238, L240: "Narrows the tag domain to one value" / "to the listed values" / "Narrows a bounded integer domain by comparison" / "Narrows the set-valued domain to assignments containing every listed element." |
| A16.4 | "RDR 0003's vocabulary clause closes the set at five" | **CONFIRMED** | 0003 L296–298 normative: "The initial operator vocabulary MUST be closed and typed: equality, membership, bounded integer comparison, existence, and set containment. Unknown operators MUST be rejected during parse or lint before resolution." |
| A16.5 | JDR 0001 §D1 fixes the atom as "(key, operator token, literal, block)" and says the kernel carries "cardinality, not meaning" | **CONFIRMED verbatim** | `docs/jdr/0001-resolve-kernel-seam.md:87` — "(key, operator token, literal, block) instead of a string. No reconstruction"; §D4 L178–179 quotes §D1 as having said the kernel "carries *cardinality*, not *meaning*." |
| A16.6 | §D4(b) prices the exception as "one grammar fact (the existence operator token and its two boolean literal forms), exported as constants 0002's normalizer MUST emit" | **CONFIRMED verbatim** | JDR L194–197: "the kernel learns one grammar fact (the existence operator token and its two boolean literal forms), exported as constants 0002's normalizer MUST emit." |
| A16.7 | "RDR 0003 and RDR 0002 each carry the identical joint-check line naming that kernel-exported token" | **CONFIRMED — byte-identical in both** | `0003:413` and `0002:444` carry the same sentence: "…§D4 settles the value-only per-atom `GuardEvaluator` seam, the kernel-exported existence operator token and boolean literal forms the normalizer MUST emit, exact key identity at the kernel, and the per-atom `guard_unevaluable` payload replacing `Refusal.Guard`. 0007 is the landing document; this RDR cites §D4 rather than restating it." |
| A16.8 | "their bodies restate it at their re-lock" — i.e. neither body carries it yet | **CONFIRMED (the negative holds)** | The joint-check line is the *only* hit for "kernel-exported / existence token / exported as constants" in either document; neither 0002's nor 0003's normative blocks mention the token. |

### A17 — value operator decidable from atom + present value alone

| # | Claim | Verdict | Cite |
| --- | --- | --- | --- |
| A17.1 | 0003's operator semantics quoted ("equality compares a tag value to one typed literal; membership checks a scalar tag against a typed literal set; bounded integer comparison uses `lt`, `lte`, `gt`, and `gte`; existence checks presence") | **CONFIRMED, with one wording drift** | `0003:358–361` Load-Bearing Decisions: "equality compares a tag value to one typed literal; membership checks a scalar tag against a typed literal set; bounded integer comparison uses `lt`, `lte`, `gt`, and `gte`; existence checks **presence of an optional tag value**; set containment checks declared set-valued tags against a typed element set." The RDR truncates at "existence checks presence"; substance unaffected. |
| A17.2 | "set containment 'Narrows the set-valued domain to assignments containing every listed element'" | **CONFIRMED verbatim** | `0003:240`. |
| A17.3 | "declared bounds are a LINT input, never an evaluation input" — quoting 0003 Technical Design | **CONFIRMED verbatim** | `0003:279–281`: "A guard may still compare an unbounded integer at runtime, but lint must report that it cannot prove exhaustive coverage for that dimension." |
| A17.4 | "RDR 0003's semantic identity tuple is `(tag, operator, literal)` with no kind field" | **CONFIRMED** | `0003:350–352`: "a guard predicate is identified by its source rule/context id plus its position within `all` or `unless`; semantic equality is the normalized tuple `(tag, operator, literal)`." No kind component. |
| A17.5 | "The declared kind needed to *parse* the value is known to the evaluator from the table it was built for, not from the kernel" | **CONFIRMED as consistent, but it is a design assertion, not a peer-document quote** | 0003's only declared-kind clause is `0003:302–304` ("A predicate whose literal cannot be parsed as the declared tag kind MUST be rejected before resolution") — about the authored literal, not runtime value plumbing. 0003 nowhere states where the evaluator obtains declared kinds. Recorded, not a refutation: the RDR's own "If wrong" branch already covers it. |

### A18 — seam MAY answer unevaluable for an uncomparable present value

| # | Claim | Verdict | Cite |
| --- | --- | --- | --- |
| A18.1 | "the nil-seam case `resolve.go::evaluateGuard` already answers `GuardUnevaluable`" | **CONFIRMED** | `resolve.go:425–433`: `if guard == "" { return GuardTrue }; if seam == nil { return GuardUnevaluable }; return seam.Evaluate(guard, view)`. |
| A18.2 | "`resolve.go::KindGuardUnevaluable` is documented 'the guard seam reported a predicate it could not decide'" | **CONFIRMED verbatim** | `resolve.go:56–58`: `// KindGuardUnevaluable: the guard seam reported a predicate it could // not decide.` |
| A18.3 | "RDR 0003's only parse-rejection clause is about the authored LITERAL" | **CONFIRMED** | `0003:302–304` normative: "Each operator MUST declare which tag value kinds it accepts. A predicate whose literal cannot be parsed as the declared tag kind MUST be rejected before resolution." Sweep of 0003 for `parsed`/`parse` finds no runtime-value clause. |
| A18.4 | "its Testing Strategy lists 'literal parse mismatch' among PARSE-time scenarios only" | **CONFIRMED verbatim** | `0003:748–749`: "**Scenario**: Parse malformed guard atoms: unknown tag, unknown operator, unsupported operator/tag-kind pair, and literal parse mismatch." Expected: "Each failure is rejected before resolution…" — parse-time. |
| A18.5 | "so 0003 does not make an unparseable RUNTIME value a load-time impossibility, and the gap this assumption anticipates is real" | **CONFIRMED (negative holds)** | No clause anywhere in 0003 governs a runtime tag value failing to parse as its declared kind. |
| A18.6 | The three-valued `GuardResult` is shipped | **CONFIRMED** | `resolve.go:77–87`: `GuardFalse`, `GuardTrue`, `GuardUnevaluable`. |

### A19 — Pending; blocked on an RDR 0005 / JDR §JD-8 envelope grant

| # | Claim | Verdict | Cite |
| --- | --- | --- | --- |
| A19.1 | "`internal/cli/clierr/clierr.go::CLIError` carries only flat strings — `Code`, `Message`, `Param`, `Detail`, `Hint`" | **CONFIRMED** | `clierr.go` `type CLIError struct`: `Code string`, `Message string`, `Param string`, `Detail string`, `Hint string`, plus non-serialized `Group ErrorGroup \`json:"-"\`` and `Cause error \`json:"-"\``. All wire fields are `string`. |
| A19.2 | "`Detail` is a `string` documented 'May be multi-line'" | **CONFIRMED verbatim** | `clierr.go`: `// Detail carries hard facts about the failure (the underlying // syscall reason, a parser diagnostic). May be multi-line.` |
| A19.3 | "`CLIError`'s own doc comment invites exactly this — 'Extend with new optional fields as needed — keep them `omitempty` so the envelope stays append-only and stable for tools'" | **CONFIRMED verbatim** | `clierr.go` type doc comment, immediately above `type CLIError struct`. |
| A19.4 | RDR 0005's normative split: "Failures MUST use the existing CLIError JSON/text envelope defined by `docs/cli-output-contract.md` and `internal/cli/clierr`" | **CONFIRMED verbatim** | `docs/rdr/0005-skill-integration-cli-contract.md:366–367`. |
| A19.5 | "`guard_unevaluable` already has a row in RDR 0005's stable-code table (`flow-guard-unevaluable`, `GroupUserEnv`) — though scoped there to 'supplied facts'" | **CONFIRMED verbatim** | `0005:611` — `| guard cannot be evaluated from supplied facts | \`flow-guard-unevaluable\` | \`GroupUserEnv\` |`. The "from supplied facts" scoping is exactly as the RDR reports. |
| A19.6 | "`owned_state_unavailable` has NO row" in 0005's code table | **CONFIRMED (negative holds)** | 0005's table L603–614 has no `owned_state_unavailable` row and no code resembling one; the nearest is `flow-fact-missing` ("required supplied fact is missing"). |
| A19.7 | "a gap §JD-8 names explicitly ('need `Code` values')" | **CONFIRMED verbatim** | JDR `:252–253`: "**JD-8 Refusal codes.** `owned_state_unavailable` and `reserved_tag_key` need `Code` values". |
| A19.8 | §JD-8 pre-authorization: "only require new `Code` constants or literals, **not new envelope fields or exit groups**" | **CONFIRMED verbatim** | JDR `:254–256`: "0005's own A-block pre-authorizes this: resolver-specific values 'only require new `Code` constants or literals, **not new envelope fields or exit groups**.'" |
| A19.9 | "RDR 0005 is `Final`" | **CONFIRMED** | `0005:9` — `- **Status**: Final [joint decision → JDR 0001 §JD-8, §JD-9]`. |
| A19.10 | The named plan's premise — that §JD-8 as written has NOT been amended to grant a structured field | **CONFIRMED, peer state is as described TODAY** | JDR §JD-8 (`:252–260`) is unchanged: it still says "not new envelope fields", still marks itself *(blank, except the exit-3 call)*, and carries no amendment. **However see problem #1** — §D4 L206–207 states the payload "reaches the CLI through §JD-8's `Detail`", which is a *decided* answer contradicting the RDR's "settled in DIRECTION … structured field, not flattened text." |

### A21 — escape rows carry empty `RequiresOwned` by composition

| # | Claim | Verdict | Cite |
| --- | --- | --- | --- |
| A21.1 | RDR 0009 Normative Contracts: "a Row with a non-empty Escape list MUST have an empty Writes slice" | **CONFIRMED verbatim (as a substring of the clause)** | `docs/rdr/0009-escape-row-shape-conformance-ownership.md:824–826`, inside a ```normative block: "…is a PRODUCER obligation on every constructor of resolve.Row values: a Row with a non-empty Escape list MUST have an empty Writes slice." |
| A21.2 | "RDR 0009 … is `Final`" | **CONFIRMED** | `0009:9` — `- **Status**: Final [joint decision → JDR 0001 §JD-5]`. |
| A21.3 | "note it mentions `RequiresOwned` ZERO times" | **CONFIRMED** | `grep -c 'RequiresOwned' docs/rdr/0009-*.md` → **0**. |
| A21.4 | "`resolve.go::missingOwned` ranges over `row.RequiresOwned`, so an empty slice never enters the loop and `missing` stays nil" | **CONFIRMED** | `resolve.go:444–458`: `for _, row := range candidates { for _, key := range row.RequiresOwned { … } }`; `var missing []string` stays nil if no iteration. `gate:388` guards on `len(missing) > 0`. |
| A21.5 | JDR §JD-3 names this the open composition; field "appears **zero** times in 0002"; calls the D1 `Row` reopening "the natural occasion to settle it" | **CONFIRMED verbatim** | JDR `:231–237`: "The field appears **zero** times in 0002, which owns the normalized row… On escape rows 0007 derives it from `Writes` while 0009 empties `Writes`. Not answered by D1, but D1 reopens the same type and is the natural occasion to settle it." |
| A21.6 | "§JD-3 remains NOT closed (unlike §JD-1/§JD-2/§JD-6/§JD-7/§JD-12)" | **CONFIRMED** | JD-3's entry carries no "**Closed by §Dn**" marker, unlike JD-1 (`:226` "**Closed by §D1**") and JD-12 (`:275` "**Closed by §D4**"). |
| A21.7 | "RDR 0002 carries the producer only as a re-entry Direction ('Name the `RequiresOwned` producer'), not yet a clause" | **CONFIRMED — and it is 0002's ONLY mention of the field** | `0002:863` (Refinement Context → Direction): "Name the `RequiresOwned` producer (JD-3) and rename the three canonical fixtures (JD-10)." `grep -c 'RequiresOwned' docs/rdr/0002-*.md` → **1**, that line. |

### A22 — exact-string key identity; canonicalization upstream

| # | Claim | Verdict | Cite |
| --- | --- | --- | --- |
| A22.1 | "`resolve.go::assemble` keys `view.tags` by `Tag.Key` verbatim" | **CONFIRMED** | `resolve.go:148–167`: `view.tags[t.Key] = taggedValue{...}` for observed and owned; `view.tags[recognizedTagKey]` for recognized. No transform on the key. |
| A22.2 | "`TagSet.Lookup` indexes the map by the string it is given" | **CONFIRMED** | `resolve.go:113–120`: `tv, ok := s.tags[key]`. |
| A22.3 | "the kernel has never canonicalized" | **CONFIRMED (negative holds)** | `rg -i 'canonical\|ToLower\|TrimSpace\|Fold' internal/resolve/resolve.go` → no hits. |
| A22.4 | "RDR 0002 has no canonicalization clause and no tag-name grammar or charset rule anywhere" | **CONFIRMED (negative holds)** | Sweep of 0002 for `canonical\|case-insensit\|case-sensit\|whitespace\|namespace\|charset\|tag name\|key spelling` returns only: L182 (a list of adjectives about the dump), L217 ("Tag declarations: tag name, provenance…" — a field list, not a grammar), L363–365 and L790–791 (canonical *examples* / canonical *semantic form* of the row, not key spelling), L444 (the joint-check line), L864 (Direction). No key-spelling rule. |
| A22.5 | 0002's declaration-completeness rule quoted: "The model MUST declare every tag it matches or writes, including each tag's provenance" | **CONFIRMED verbatim** | `0002:334`. |
| A22.6 | "The only supporting sentence is JDR 0001 §D4's landing note ('the normalizer … canonicalizes key spellings before a row exists')" | **CONFIRMED verbatim** | JDR `:214–215`: "**0002** cites it: the normalizer emits the kernel's existence constants and canonicalizes key spellings before a row exists." |
| A22.7 | Named plan's premise: "The duty is NOT currently on 0002's Refinement Context Direction list" | **CONFIRMED, peer state is as described TODAY** | `0002:861–864` Direction reads in full: "Restate the resolver flow as gate-then-count and qualify Scenario 4 so no sibling candidate is unevaluable. Add `row kind` and `escape failure classes` to the dump-derivation field list. Name the `RequiresOwned` producer (JD-3) and rename the three canonical fixtures (JD-10)." No canonicalization duty, and no existence-token duty either — **see problem #2**. |

### Delta-scoped passages outside the assumption blocks

| # | Passage | Verdict | Cite |
| --- | --- | --- | --- |
| N1 | Normative Contracts SEAM block: "the shape JDR 0001 §D1 fixes" | **CONFIRMED** | JDR `:87` "(key, operator token, literal, block) instead of a string. No reconstruction". |
| N2 | Normative Contracts: presence decided by "the `TagSet.Lookup` `ok` test" | **CONFIRMED** | `resolve.go:113` returns `(value, prov, ok)`; `ok` is map-presence, provenance-blind. |
| N3 | Normative Contracts (PRESENCE IS PROVENANCE-BLIND): "`missingOwned` tests owned provenance" | **CONFIRMED** | `resolve.go:452` `view.has(key, ProvenanceOwned)`; `has` at `:125–128` tests `tv.provenance == prov`. |
| N4 | Normative Contracts: "shipped: `evaluateGuard`'s empty-guard branch" (row with no atoms ⇒ TRUE) | **CONFIRMED** | `resolve.go:426–428` `if guard == "" { return GuardTrue }`. |
| N5 | Normative Contracts: "shipped: `escapeOrRefuse` returns the escape set's blocking refusal" | **CONFIRMED** | `resolve.go:473`+ doc comment and body; frozen by ADV-2 (`adversarial_test.go:192` subtest "guard UNEVALUABLE must not rescue") and Fixup-1d (`fixup_test.go:42`). |
| N6 | Technical Design: "`assemble`, `gate`, `missingOwned`, `escapeOrRefuse`, and `Resolve` are untouched in control flow" | **CONFIRMED as consistent with the A3 spike** | All five resolve in `resolve.go` at :148, :376, :444, :473, :318. |
| N7 | Testing Strategy header: "The frozen 154-test suite passes unchanged under the reshape" | **CONFIRMED (baseline half re-run)** | 154 PASS / 0 FAIL measured this pass. |
| N8 | Testing Strategy row 16: "no existing test contends the GUARD path (ADV-4/ADV-5 defend only the owned path)" | **CONFIRMED** | `adversarial_test.go:400 TestAdv4_ObservedTagsMustNotShadowTheOwnedSnapshot`, `:436 TestAdv5_ObservedTagCannotSatisfyAnOwnedStateRequirement` — both owned-snapshot assertions. No guard-path observed-substitution test exists. |
| N9 | Testing Strategy row 17: "`undecidedAtoms`/`compareAtoms`; frozen ADV-3 is the order-independence precedent" | **CONFIRMED for ADV-3; the two helpers are proposed, not shipped** | `adversarial_test.go:266`; `compareRefs` at `resolve.go:528` is the shipped analogue. |
| N10 | Implementation Plan Prerequisites bullet 3 (A19) and bullet 4 (A22) restate the two peer-document blockers | **CONFIRMED as accurate restatements of A19/A22** | Consistent with A19.9/A19.10 and A22.7 above; the §D4 `Detail` contradiction (problem #1) affects both the assumption and this bullet. |

---

## Inverse check

The RDR introduces five new discriminators. For each: does a sibling path — in
code or in a peer document's normative text — already decide it?

**1. Kernel-side existence-operator token recognition.**
- **Code: none.** `rg -n 'exists|existence' internal/ --type go` → **zero hits**
  anywhere under `internal/`. No token constant, no operator vocabulary, no
  grammar in the kernel. `GuardEvaluator.Evaluate(guard string, view TagSet)`
  (`resolve.go:91–95`) hands the whole predicate across as an opaque string.
- **Spec: a sibling DOES decide it — and the RDR cites it correctly.**
  `docs/jdr/0001-resolve-kernel-seam.md` §D4(b) L189–197 already resolves
  kernel-side token recognition, including the exact price ("one grammar fact …
  exported as constants 0002's normalizer MUST emit"), and §D4's landing note
  L211–215 names 0007 as the landing document. RDR 0002 (`:444`) and RDR 0003
  (`:413`) each carry a byte-identical joint-check line acknowledging it. This
  is the intended relationship, not a collision: 0007 is §D4's normative home.

**2. Absent-key ⇒ unevaluable derivation.**
- **Code: none.** `rg -n 'GuardUnevaluable' internal/ --type go`: every
  production occurrence is either the constant's declaration
  (`resolve.go:84–86`), the nil-seam branch (`:430`), the `gate` aggregation
  arm (`:401`, `:413`), or the refusal-kind list (`:70`). No code path derives
  undecidability from a key's absence — `gate` aggregates undecidability
  *already handed to it by the seam*. `fixtures_test.go:26` returns
  `GuardUnevaluable` for *unknown predicate text*, not for absent keys, and
  discards the view entirely (`:20`, `_ resolve.TagSet`).
- **Spec: JDR §D4 L179–180 states the rule** ("0007's domain rule — an atom
  over an absent key is unevaluable, never false") and L189–192 assigns
  enforcement to the kernel. Again the cited home, not a collision.

**3. Strong-Kleene combination across atoms.**
- **Code: none.** `rg -n -i 'kleene|three-valued' internal/` → zero production
  hits. `gate` has no combinator: it consumes one `GuardResult` per row.
- **Spec: JDR §D4(b) L192–193** ("combines under strong Kleene itself") decides
  it; no peer RDR states the tables. RDR 0003's row algebra (`0003:132`,
  "a row with `all` atoms denotes the intersection … its `unless` block denotes
  an excluded intersection that is subtracted") is set-theoretic over complete
  assignments — a two-valued lint-domain formulation, not a three-valued
  runtime combinator. The lift is genuinely 0007's. (Iteration 1 recorded this
  as F-7 against A2, which is outside this delta scope; re-confirmed here only
  as an inverse-check negative.)

**4. `RequiresOwned` narrowed to post-guard write-dependency keys.**
- **Code: the only meaning statement is the Go doc comment.**
  `resolve.go:182–184`: `// RequiresOwned names owned tag keys the row's
  evaluation needs. // A key absent from the owned snapshot yields
  owned_state_unavailable.` No peer RDR defines the field: 0009 mentions it
  zero times, 0002 mentions it once and only as an unfulfilled Direction item
  (`0002:863`).
- **Spec: NOT decided by any sibling — §JD-3 is explicitly open.** JDR
  `:231–237` records the 0007/0009 composition as unresolved and unassigned.
  A21's own "PRODUCER question stays open" framing matches the peer state.

**5. Exact-string tag-key identity.**
- **Code: the kernel is exact-string by construction** (`assemble` /
  `Lookup`, A22.1–A22.3) but states no rule — there is no comment or contract
  asserting exactness, only the map behavior.
- **Spec: JDR §D4 L214–215 files canonicalization on 0002; 0002 has not
  received it** (A22.4, A22.7). The joint-check line at `0002:444` names
  "exact key identity at the kernel" as settled by §D4, but 0002's body carries
  no clause. So a sibling *obligation* exists, unfulfilled — which is exactly
  what A22's Pending status and problem #2 below describe.

**Net: no code sibling makes any of the five decisions. The spec siblings are
JDR 0001 §D4 (which 0007 correctly names as its own source) and the two
unfulfilled obligations §D4 files on 0002. One §D4 sub-decision — the CLI
transport of the payload — contradicts the RDR (problem #1).**

---

## Refuted / not-found / sibling-exists

### 1. A19's premise that the envelope question is undecided is complicated by JDR §D4, which already decided it the *other* way — via `Detail`

RDR 0007, A19 Status:

> **Status**: Pending — blocked on a peer-document change, not on evidence. The
> kernel half is Verified; the envelope half is settled in DIRECTION by the
> author's round (2026-08-21: structured field, not flattened text) but cannot
> be stated by this RDR, because it contradicts §JD-8 as written. Named plan:
> raise the field at JDR 0001 as a §JD-8 amendment (or at RDR 0005's re-lock),
> then this assumption cites that decision and flips to Verified.

The RDR checks its plan against §JD-8 only. But **§D4 — the very decision this
RDR is the landing document for — already routed the payload, and routed it to
`Detail`**:

> `docs/jdr/0001-resolve-kernel-seam.md:204–207`
> "The `guard_unevaluable` refusal names what blocked the verdict per row and
> per atom — key, block, reason `absent` | `uncomparable` — replacing the
> single-valued `Refusal.Guard` text, which has no referent once guards are
> atoms; **it reaches the CLI through §JD-8's `Detail`.**"

So the peer-document state today is not "silent, pending a grant." It is a
*closed* joint decision (§D4, and §JD-12 is marked "**Closed by §D4**" at JDR
`:275`) whose CLI transport clause says flattened `Detail` text — the option
the author's round rejected. A19's named plan therefore understates the work:
it is not only a §JD-8 amendment adding a field, it is a reversal of a sentence
inside the closed §D4 that 0007 elsewhere cites as its mandate. The A19
Evidence block and the Implementation Plan Prerequisites bullet both quote
§JD-8 while omitting this §D4 sentence.

Everything else in A19 verifies exactly: `CLIError`'s five flat string fields,
`Detail`'s "May be multi-line" doc, the append-only `omitempty` invitation,
0005's failure-side split at `:366`, `flow-guard-unevaluable`'s "supplied
facts" scoping at `:611`, the missing `owned_state_unavailable` row, and
§JD-8's "not new envelope fields" at `:255`.

### 2. A16's obligation on RDR 0002 is in the same unfulfilled state A22 flags — but only A22 says so

RDR 0007, A16 Evidence:

> RDR 0003 and RDR 0002 each carry the identical joint-check line naming that
> kernel-exported token, so the obligation is acknowledged in both receiving
> documents; their bodies restate it at their re-lock.

The joint-check lines are verbatim-confirmed (`0002:444`, `0003:413`), and the
bodies indeed do not carry it. But A16 rests on 0002's normalizer actually
emitting the kernel constants, and — exactly as A22 discovered for
canonicalization — **0002's Refinement Context Direction list does not carry
that duty either**:

> `docs/rdr/0002-transition-table-as-reviewable-data.md:861–864`
> "**Direction.** Restate the resolver flow as gate-then-count and qualify
> Scenario 4 so no sibling candidate is unevaluable. Add `row kind` and `escape
> failure classes` to the dump-derivation field list. Name the `RequiresOwned`
> producer (JD-3) and rename the three canonical fixtures (JD-10)."

Four duties, none of them the existence-token emission. A22 is Pending
precisely because the canonicalization duty is missing from this same list, and
the Prerequisites bullet says so ("The duty is NOT currently on 0002's
Refinement Context Direction list — add it there, or A22 stays an inference").
A16 is marked **Verified** on the strength of the joint-check line alone, with
no equivalent note. The asymmetry is the finding: two obligations filed by the
same §D4 landing note on the same document sit in the same unfulfilled state,
and one is Verified while the other is Pending. Either A16's status should
carry the same caveat, or the Prerequisites bullet for 0002 should name both
duties.

### 3. A21's "verbatim as quoted" claim is accurate, but the quoted fragment is a clause tail whose head narrows it to producers

RDR 0007, A21 Evidence:

> RDR 0009 Normative Contracts: "a Row with a non-empty Escape list MUST have
> an empty Writes slice." … RDR 0009's clause is verbatim as quoted and is
> `Final`

The string is present verbatim (`0009:825–826`), but the RDR quotes it starting
mid-sentence. The full normative clause is:

> `docs/rdr/0009-escape-row-shape-conformance-ownership.md:822–834`
> "Escape-row shape conformance — an escape row carries no owned-state mutation
> — is a **PRODUCER obligation on every constructor of `resolve.Row` values**:
> a Row with a non-empty Escape list MUST have an empty Writes slice.
> (Authored clears normalize to `<clear>` writes per RDR 0002, so the Writes
> predicate carries both "no writes" and "no clears" at the kernel boundary.)
> The predicate is Writes-only and does NOT extend to NextTags…"

This matters for A21's composition. 0009 does not assert that escape rows *as
the kernel receives them* have empty `Writes`; it places the obligation on
producers — and JD-3 (`:231–237`) records that **no layer is currently obliged
to populate `RequiresOwned` at all**. A21's derivation therefore chains one
unfulfilled producer obligation (0009's, on a producer that does not exist yet)
onto another open producer question (§JD-3). The RDR's own hedge — "the
composition holds regardless of who populates the field; the deferral must
survive into 0002's re-lock or the derive-from-`Writes` definition has no
implementer" — is the right one, and the kernel half (`missingOwned` over an
empty slice, `resolve.go:444–458`) is exactly right. Recorded so a reader of
the truncated quote does not read 0009 as supplying a kernel-boundary
guarantee it scopes to producers.

### 4. Not-found (benign, recorded for the implementer): five A3 touch-list symbols do not exist in shipped source

RDR 0007, A3 Status:

> the reshape also introduces `GuardAtom`, `evaluateAtom` (where the
> presence/existence rule lands), and `undecidedAtoms`/`compareAtoms`/`copyAtoms`
> to keep the slice-valued payload tuple-deterministic under REQ-1.

`rg 'GuardAtom|evaluateAtom|undecidedAtoms|compareAtoms|copyAtoms' internal/`
returns **no hits**. This is correct — they are proposed additions, and the A3
text does say "introduces" — but Testing Strategy row 17 cites
"`undecidedAtoms`/`compareAtoms`" in the *Code path / evidence* column
alongside genuinely shipped anchors like `gate` and `escapeOrRefuse`, with no
marker distinguishing proposed from shipped. The shipped analogues are
`resolve.go:528 compareRefs` and `resolve.go:569 copyTags`. **NOT-FOUND, no
defect** — flagged only because the evidence column mixes the two categories.

### 5. Minor: A17's quotation of RDR 0003's existence semantics is truncated mid-clause

RDR 0007, A17 Evidence:

> RDR 0003's operator semantics — "equality compares a tag value to one typed
> literal; membership checks a scalar tag against a typed literal set; bounded
> integer comparison uses `lt`, `lte`, `gt`, and `gte`; existence checks
> presence"

The source reads:

> `docs/rdr/0003-guard-predicate-exhaustiveness.md:358–361`
> "…existence checks **presence of an optional tag value**; set containment
> checks declared set-valued tags against a typed element set."

The truncation drops "of an optional tag value." Substance is unaffected —
A17's claim (one value against one literal, no view needed) holds — but the
quote is presented inside quotation marks as if complete. Also recorded:
A17's supporting sentence "The declared kind needed to *parse* the value is
known to the evaluator from the table it was built for, not from the kernel"
has **no peer-document backing**; 0003 nowhere says where the evaluator obtains
declared kinds. It is a design assertion inside a **Method: Source Search**
assumption, and the "If wrong" branch already absorbs it.

---

## Summary

Of the seven re-verified assumptions, the **kernel-side** half of every one
verifies exactly against `internal/resolve/resolve.go` and the 154-test suite
(re-run: 154 PASS / 0 FAIL). **A3's re-verification trigger was re-run in this
pass and is CLEAN** — zero production importers, `NewRootCmd` still registers
only `newVersionCmd()`.

The problems are all at the peer-document boundary, and they cluster on
**RDR 0002's Refinement Context Direction list** and **JDR §D4's own text**:

- §D4 already routes the A19 payload to `Detail` (problem #1), so A19's plan is
  a reversal inside a closed decision, not just a §JD-8 addition;
- 0002's Direction list carries neither the A16 existence-token duty nor the
  A22 canonicalization duty (problem #2), yet A16 is Verified and A22 Pending;
- A21's 0009 quote is verbatim but scoped to producers, not to the kernel
  boundary (problem #3).

Two truncated-quotation notes (#4, #5) are recorded as precision issues, not
refutations.
