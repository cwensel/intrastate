# REQ List — RDR 0012 Declared-kind carrier at the guard value seam

Phase 0 audit artifact. Every clause below is a testable obligation drawn from
`docs/rdr/0012-declared-kind-carrier-at-the-guard-seam.md`. Quotes are
verbatim, copied from the projector (`recs inspect --select <id>`), never
transcribed; a line break inside a quoted source span is rendered as one space
and `…` marks an elision.

Element ids are carried where the REQ derives from a labelled contract
(`0012:C1` … `0012:C5`, `0012:MVV`) or a Testing Strategy scenario
(`0012:S1` … `0012:S9`). Source sections are abbreviated:

- `NC` = Proposed Solution / Technical Design / Normative Contracts (fenced)
- `NCP` = Normative Contracts, unfenced prose between C4 and C5
- `MC` = Technical Design / Pre-lock mini-checks
- `LBD` = Technical Design / Load-Bearing Decisions
- `FM` = Trade-offs / Failure Modes
- `MVV` = Implementation Plan / Minimum Viable Validation
- `IP` = Implementation Plan (phases)
- `TS` = Validation / Testing Strategy

**Three recorded deviations against predecessors are carried by the record
itself** (never amended here; listed so a later stage does not read them as
contradictions):

1. `0012:C3` amends `0007:REQ-71`'s parameter shape: `TestGuardEvaluatorContract`
   keeps its name and its kernel non-test home but takes a seam CONSTRUCTOR
   (`newSeam func(kinds map[string]string) GuardEvaluator`) instead of a
   constructed seam. The `var fn` pin in
   `internal/resolve/guard_mvv_test.go::TestReq71_GuardEvaluatorContractIsAnImportableCrossRDRSurface`
   is re-typed, not deleted (REQ-32).
2. `0012:A2` narrows `0003:REQ-34`'s implemented proxy — the zero-field
   reflection check in
   `internal/guard/guard_evaluator_0003_test.go::TestReq34_EvaluatorDecidesPresentValuesOnlyAndNeverReadsTheView`
   — to "no view-typed or runtime-valued state" (REQ-47). REQ-34's normative
   text ("MUST NOT read the tag view") is preserved.
3. `0012` §Approach narrows `0007:A17`'s headline ("nor tag declarations at
   evaluation time") to "needs no declaration it was not built over". The seam
   signature `0007:REQ-10` fixes is untouched (REQ-8).

---

## Carrier — construction, not signature (`0012:C1`)

- [REQ-1] "CARRIER. The declared kind travels to the value seam by CONSTRUCTION, not by signature, atom field, or upstream refusal. func NewEvaluator(kinds map[string]string) Evaluator" — (NC, `0012:C1`; LBD Naming "`NewEvaluator` (the Go-conventional constructor for the existing `Evaluator` type)")

- [REQ-2] "`kinds` maps tag key → declared kind token (the five-kind vocabulary `enum | bool | int | set | scalar`, RDR 0003's spelling). Its producer is EXPORTED and singular: func DeclaredKinds(m *table.Model) map[string]string" … "`DeclaredKinds` is homed in `internal/guard`" — (NC, `0012:C1`)

- [REQ-3] "It is TOTAL over `m.Tags`: every declared key gets an entry, and a key whose `Kind` is the empty string is OMITTED rather than mapped to `\"\"` — an empty-string entry would be indistinguishable at C2's arms from an absent key, collapsing two diagnostics into one." — (NC, `0012:C1`)

- [REQ-4] "Handed a nil model it returns an empty non-nil map rather than panicking" — (NC, `0012:C1`)

- [REQ-5] "The evaluator holds this mapping and NOTHING else: no `resolve.TagSet`, no runtime tag value, no view-typed state — \"never reads the tag view\" (`0003:C9`, REQ-34) remains a property of the type." — (NC, `0012:C1`; asserted by REQ-47 / `0012:S3`)

- [REQ-6] "The zero-value `Evaluator{}` construction form is retired; every non-test construction goes through `NewEvaluator`." — (NC, `0012:C1`; the one admitted exception is REQ-35's match site; mechanical backstop REQ-48 / `0012:S4`)

- [REQ-7] "a nil-mapping evaluator answers GuardUnevaluable for every `eq`/`in` atom (the no-declaration arm), surfacing as `guard_unevaluable` refusals — never silently reverting to raw-string comparison." — (NC, `0012:C1`; MC `disposition` row "Nil-mapping evaluator (missed C4 site)"; asserted by REQ-49 / `0012:S5`)

- [REQ-8] "The seam signature `Evaluate(atom GuardAtom, value string) GuardResult` (`0007:C1`, REQ-10) and the atom shape (JDR 0001 §D1: Key, Operator, Literal, Block) are UNCHANGED" — (NC, `0012:C1`) together with "the atom's semantic identity tuple stays `(tag, operator, literal)` (RDR 0003) with no kind member, and no kind is ever stamped into an atom, payload, or hash." — (LBD, Identity)

- [REQ-9] "the kernel's RESOLUTION PATH — `Resolve` and its callees — continues to read no declarations. The rule governs that path, not `internal/resolve` as a namespace" — (NC, `0012:C1`)

- [REQ-10] "`NewEvaluator`'s return type is NOT widened to the interface to avoid the wrapper: C1 fixes the concrete value shape deliberately" — (NC, `0012:C3`) `NewEvaluator` returns the concrete `guard.Evaluator`.

## Typed comparison on the guard path (`0012:C2`)

- [REQ-11] "TYPED COMPARISON APPLIES ON THE GUARD PATH. `eq` and `in` GUARD atoms resolve the atom's key against the constructed kind mapping and compare under the declared kind. MATCH atoms do NOT: the kernel byte-compares them at `internal/resolve/resolve.go::TagSet.matches`" — (NC, `0012:C2`; `0012:S2`)

- [REQ-12] "kind `int`: held value AND literal (each member, for `in`) must parse as integers; comparison is over the PARSED values. A present held value that does not parse is GuardUnevaluable — never GuardFalse." — (NC, `0012:C2`; `0012:S2`; parsed-equality leg REQ-53 / `0012:S9`)

- [REQ-13] "For `in`, ONE unparseable member poisons the whole list: GuardUnevaluable even when the held value equals a member that did parse. Parse-all-then-compare, not parse-and-compare-per-member" — (NC, `0012:C2`)

- [REQ-14] "The parse is `strconv.Atoi` … an out-of-range spelling is GuardUnevaluable on the target under test, and the suite pins the in-range cases only." — (NC, `0012:C2`; overflow leg of REQ-25)

- [REQ-15] "kind `bool`: held value and literal must be the boolean tokens (`true` | `false`, the `0007:A27` spellings); comparison is token equality; any other held value is GuardUnevaluable. The rejected set is every other spelling — `\"1\"`, `\"0\"`, `\"True\"`, `\"yes\"`" — (NC, `0012:C2`; `0012:S2`)

- [REQ-16] "kind `enum` | `scalar`: every string parses; comparison is exact string equality (unchanged behavior). These two SHARE one arm" — (NC, `0012:C2`; `0012:S2`)

- [REQ-17] "kind `set`: the operator/kind matrix does not admit `eq`/`in` over `set`; handed one anyway, the seam answers GuardUnevaluable" — (NC, `0012:C2`; `0012:S2`)

- [REQ-18] "key absent from the mapping: GuardUnevaluable — the seam cannot type the comparison." — (NC, `0012:C2`; `0012:S2`)

- [REQ-19] "For such a key this is the ordinary arm, not a defensive one: `eq`/`in` over it answers GuardUnevaluable where today it compares raw strings." … "so the conformance suite pins it as a named case rather than filing it under the unreachable defensive arms." — (NC, `0012:C2`; "such a key" = "a key that IS declared with `Kind` the empty string")

- [REQ-20] "unparseable LITERAL: GuardUnevaluable (defense in depth; load-time rejection per `0003:C8` … stays the primary guard, whose writer is the model loader)." — (NC, `0012:C2`)

- [REQ-21] "Dispatch over the kind token is EXHAUSTIVE across the five-kind vocabulary, and the fallthrough `default` means GuardUnevaluable" … "so the suite asserts an unknown kind token answers GuardUnevaluable rather than comparing." — (NC, `0012:C2`)

- [REQ-22] "`lt/lte/gt/gte` keep their operator-inferred integer parse (the matrix admits only `int`); `contains` and the §D13 set arms are unchanged; `exists` never reaches the seam (`0007:C1`)." … "the kind lookup is not consulted for those operators at all, so `contains` over a `set` key does NOT become GuardUnevaluable under the `set` arm above (which governs `eq`/`in` only)." — (NC, `0012:C2`) Pairs with the S5 scope note: "a nil-mapping evaluator still answers them normally" — (TS, `0012:S5`)

- [REQ-23] "an operator the seam does not recognize — `exists` included, which `0007:C1` keeps off this path — answers GuardUnevaluable like any other atom it cannot type. The seam NEVER panics" — (NC, `0012:C2`)

- [REQ-24] "\"Each member, for `in`\" above means the members decoded from that array, and parse-all-then-compare quantifies over them. No new members field and no delimiter convention is introduced." — (NC, `0012:C2`)

## Conformance suite (`0012:C3`)

- [REQ-25] "`eq`/`in` over an `int` and a `bool` dimension where the held value does not parse (want GuardUnevaluable) and where it parses but differs (want GuardFalse); plus the defensive arms — an `in` list with a non-integer member on an `int` dimension, an integer-overflow held value, `eq` over a `set`-kind key, and a key absent from the fixture (each want GuardUnevaluable)." — (NC, `0012:C3`; `0012:S1`)

- [REQ-26] "The suite carries at least one discriminating case per kind token, and a meta-check that the fixture's kind tokens are exactly the five-kind vocabulary." — (NC, `0012:C3`)

- [REQ-27] "The suite publishes the declaration fixture its cases assume as a kernel-owned `map[string]string` (key → kind token; no `internal/table` type crosses — A5). It is published as a FUNCTION returning a fresh map: func ContractKinds() map[string]string not a package-level `var`." … "A constructor hands each caller its own copy." — (NC, `0012:C3`)

- [REQ-28] "Delivery is STRUCTURAL, not a documented obligation — the suite takes a constructor and builds the seam itself: func TestGuardEvaluatorContract(t *testing.T, newSeam func(kinds map[string]string) GuardEvaluator)" — (NC, `0012:C3`)

- [REQ-29] "each case names its own key drawn from the fixture, chosen so the case's operator is one the matrix admits for that kind. Today every case is built with `Key: \"subject\"`; re-keying routes the existing `eq`/`in` cases onto an `enum`/`scalar` key, the five `gte` cases onto the `int` key, and the four `contains` cases onto the `set` key, with the new legs on the `int`, `bool` and `set` keys." — (NC, `0012:C3`) With "their verdicts are unchanged" for the re-keyed `gte`/`contains` cases — (NC, `0012:C2`)

- [REQ-30] "Both in-tree implementers are bound, not only `guard.Evaluator`" … "Concretely it gains a `kinds map[string]string` field and a matching constructor to satisfy the `newSeam` parameter above; `recordingSeam` keeps delegating to its `inner` and needs no kinds of its own." — (NC, `0012:C3`; "it" = `internal/resolve/guard_mvv_test.go::conformingContractSeam`)

- [REQ-31] "`guard.Evaluator`'s existing suite call is RE-POINTED, not duplicated." … "It therefore migrates to a call over the published fixture in the SAME phase as the signature widening" … "Each migrating site therefore supplies a one-line adapter: func(k map[string]string) resolve.GuardEvaluator { return guard.NewEvaluator(k) }" — (NC, `0012:C3`; the site is `internal/guard/guard_evaluator_0003_test.go:67`, today `var ev guard.Evaluator`)

- [REQ-32] "its importable-cross-RDR-surface guarantee is preserved — same name, still declared in the kernel's non-test sources (the `kernelDeclaresFunc` half at :227 is untouched) — and only the parameter shape changes" … "the `var fn` line is re-typed to the new signature rather than deleted, so the pin keeps pinning." — (NC, `0012:C3`)

## Construction and consumption (`0012:C4`)

- [REQ-33] "EVERY non-test evaluator construction site — present or later added — constructs over the declarations of the SAME loaded model whose rows it evaluates: one model, one evaluator, no mixed-model evaluation. The obligation is scoped to the GUARD path" … "On `main` that set is exactly two: `internal/cli/flow_resolve.go::guardSeam`, `internal/guard/product.go::valueSatisfies` (A3)." — (NC, `0012:C4`; `0012:S4`) Testable half: both sites construct through `NewEvaluator(DeclaredKinds(m))` over the model in play, witnessed end to end by a typed verdict at `flow resolve`, `flow next` and lint. The wrong-model half is EXCLUDED below.

- [REQ-34] "At the CLI the constructed evaluator is THREADED DOWN, not re-derived: `guardSeam` takes the model and is called once per request, and `internal/cli/flow_next.go::probeRow` receives the constructed `resolve.GuardEvaluator` rather than widening to take a `*table.Model`." — (NC, `0012:C4`)

- [REQ-35] "`internal/graphlint/reach.go::atomAdmitsValue` is NOT a construction site under this clause and keeps its `main` shape: it constructs `guard.Evaluator{}` with no mapping and compares raw bytes. That is required, not tolerated." … "`atomAdmitsValue` keeps its `bool` return and its name." — (NC, `0012:C4`; IP Phase 3 "deliberately NOT migrated")

- [REQ-36] "Every consumer of a seam verdict handles GuardUnevaluable as its OWN case, distinct from both GuardTrue and GuardFalse; a `!= GuardFalse` / `== GuardTrue` two-valued collapse at any consumer is a defect this clause forbids." — (NC, `0012:C4`; scoped to the guard path: "on the guard path both already hold on `main`, and this clause keeps them holding as the seam becomes typed")

- [REQ-37] "That split is the record's one asymmetry and it is named on purpose (C2 states it on the comparison side, and the user-facing docs carry it): one declared kind, two comparison rules, divided by whether the two deciders meet." — (NC, `0012:C4`; LBD "C2 and C4 state the split, and the user-facing docs carry it, so it is discoverable rather than surprising.")

## No new seam refusal surface (`NCP`)

- [REQ-38] "No new refusal kind, error code, or envelope field for the SEAM: a typed `GuardUnevaluable` reaches the user through the existing `guard_unevaluable` refusal and its payload (`0007:C8`), reason `uncomparable` (\"the key was present and its value was not compared to a verdict\") — the CLI already maps it to `flow-guard-unevaluable`." — (NCP; MC `disposition` rows 2–6)

## Canonical int spelling at the predicate ingress (`0012:C5`)

- [REQ-39] "CANONICAL INT SPELLING. A value authored against an `int`-declared tag is admitted only in its canonical decimal spelling: after `strconv.Atoi` succeeds, `strconv.Itoa(n)` MUST equal the authored string. Non-canonical spellings (`\"00\"`, `\"01\"`, `\"+1\"`, `\"-0\"`) are REFUSED at the ingress that admits them. The rule's scope is the PREDICATE ingress and only it: guard and match atom literals, the values a typed comparison reads." — (NC, `0012:C5`; IP Phase 4)

- [REQ-40] "It also covers the bare-float path there, where `internal/table/load.go::valueMembers` renders a TOML `eq = -0.0` to the literal `\"-0\"` via `strconv.FormatFloat` before any kind check." — (NC, `0012:C5`; `0012:S7`)

- [REQ-41] "The canonicality check is therefore applied at the atom-building site (`internal/table/normalize.go::(*loader).atom`, which already calls `conform(decl, operator, members)` at `normalize.go:172`), leaving `conformKind` itself unchanged." — (NC, `0012:C5`) Observable: `[emit]`/`[initial]`/`[rule.write]` non-canonical int values stay ADMITTED (MC `disposition` row "ADMITTED unchanged"), and the two BOUNDARY tests `internal/table/emit_grammar_0024_test.go::TestReq9_0024_NonCanonicalIntLiteralsAreAdmitted` and `::TestReq7_0024_EveryKindCheckIsLexicalAndTheAuthoredBytesSurvive` stay green.

- [REQ-42] "The CLI is OUT of scope: `--tag n=07` / `--write n=+1` stay accepted." — (NC, `0012:C5`; MC `disposition` row "Same, CLI … ADMITTED unchanged")

- [REQ-43] "The refusal is user-visible and DIAGNOSTIC: it names the offending site and the canonical rewrite, e.g. `<site>: \"00\" is not the canonical spelling of int tag n; write 0`. It travels as whatever error value `(*loader).atom`'s existing `conform` call already returns to its caller — the check is folded in beside that call and returns the same shape, adding no new error type and no new sentinel." — (NC, `0012:C5`)

- [REQ-44] "It reuses the one refusal code its ingress already carries — `malformed_predicate_atom` (guard and match atoms, via `internal/table/normalize.go::(*loader).atom`) — and adds no envelope field and no sixth refusal kind" — (NC, `0012:C5`)

- [REQ-45] "\"Canonical\" here is exactly the `strconv.Itoa` round-trip stated above and nothing wider — no case folding, whitespace or unicode normalization is claimed or performed." — (NC, `0012:C5`; MC `fidelity` row "\"Canonical\" scope")

## Minimum Viable Validation (`0012:MVV`)

- [REQ-MVV] "Author a model declaring `iter` as `int`, single-valued, with a row guarded `iter eq 7` and an unguarded fallback row on the same outcome, plus an accessor reading `iter` as an OWNED value" … "The reader is driven by the persisted artifact, not by a test double: steps 2 and 3 set `iter`'s owned value in the artifact the accessor reads (`many`, then `07`, then `4`) and re-run the command" … "2. With the reader returning `many` for `iter`: `flow resolve --outcome <o>` → refuses `flow-guard-unevaluable`; the payload names the guarded row and the `iter` atom with reason `uncomparable`." … "The envelope carries the existing `code`/`message`/`schema_version` fields and NO new top-level key" … "3. With the reader returning `07` against the guard `iter eq 7`: the guarded row MATCHES under parsed comparison where it was pruned under raw-string comparison" … "With the reader returning `4`: the guarded row prunes (GuardFalse) and the fallback plans" … "4. `TestGuardEvaluatorContract`, extended per C3 and driven against `guard.NewEvaluator` constructed over the published fixture kinds, is green — including the int/bool want-Unevaluable legs that fail against today's raw-string arms." … "5. `lint --model` on a model authoring `n eq \"00\"` over an `int` tag refuses with C5's canonical-spelling diagnostic naming the site and the rewrite" — (MVV, `0012:MVV`)

  The venue is load-bearing: "A seam-level unit test over `NewEvaluator` would satisfy the assertions while testing none of the ingress, which is the whole claim." Steps 2 and 3 run the real `OwnedSnapshot` → `resolve::assemble` path end to end; step 3's `4` leg is the negative control ("proving step 3 is not blanket-true").

## Testing Strategy obligations (outside `normative` fences)

- [REQ-46] "**Expected**: the C2 matrix verbatim." over "Per-kind dispatch at the seam — `int`, `bool`, `enum`, `scalar`, `set`, and key-absent — over `eq` and `in`." — (TS, `0012:S2`) One guard-package table driving REQ-12…REQ-18 through `NewEvaluator`.

- [REQ-47] "The zero-field reflection assertion in `guard_evaluator_0003_test.go` narrowed to \"no view-typed or runtime-valued state\". **Expected**: passes with the declaration mapping present, fails if a `resolve.TagSet` or runtime value is added." — (TS, `0012:S3`; IP Phase 1 "REQ-34 test's field check narrowed per A2")

- [REQ-48] "no zero-value construction of `guard.Evaluator` outside `NewEvaluator` survives in non-test code, with ONE allow-listed exception: `internal/graphlint/reach.go::atomAdmitsValue`" … "The assertion matches the composite-literal form (`Evaluator{}`) *and* the `var` / `new` / embedded-field zero values Go equally admits" … "The allow-list is ONE named function, asserted by name rather than by pattern, so a second untyped site cannot slip in under it and moving or renaming that function fails the check" … "Implement it as a typed check, not a text grep: a `go/ast` or `golang.org/x/tools/go/packages` pass over non-test files in the module, failing on any composite literal, `var` declaration, `new`, or struct field of type `guard.Evaluator` outside `NewEvaluator`'s own body." … "it is wired into `make check` rather than left as prose." — (TS, `0012:S4`)

- [REQ-49] "Construct a bare `guard.Evaluator{}` (and the `var` / `new` zero values) and evaluate an `eq` atom whose held value equals its literal, plus an `in` atom whose held value is a member. **Expected**: `GuardUnevaluable` for BOTH — never `GuardTrue`." — (TS, `0012:S5`)

- [REQ-50] "Lint/runtime agreement over the committed corpus — before/after verdict diff. **Expected**: no committed model's verdict flips." — (TS, `0012:S6`) With "A green `make check` after this phase is the expected outcome — no existing corpus or testdata model should start failing" — (IP, Phase 4)

- [REQ-51] "model `cover07` declares `[tags.n]` int min 0 max 1, row `r-zero` guarded `n eq \"00\"` emitting `zero`, row `r-one` guarded `n eq 1` emitting `one`; invoked `intrastate lint --model cover07.toml --as json`. **Expected** (normative fixture, C5): exit 2 with the canonical-spelling load refusal naming the site and the rewrite — `\"00\" is not the canonical spelling of int tag n; write 0` — under the ingress's existing `malformed_predicate_atom` code." … "Also covers the bare-float spelling `eq = -0.0`, which `valueMembers` renders to the literal `\"-0\"` (A4)." — (TS, `0012:S7`)

- [REQ-52] "same `mvv-int` model (`iter` int min 0 max 9, row `guarded` guarded `iter eq 7`, unguarded row `fallback`); invoked `intrastate flow resolve --model mvv-int.toml --outcome step --tag iter=many --plan-only --as json`." … "**Expected** (normative fixture), unchanged before and after: exit 2," — (TS, `0012:S8`) The expected envelope, verbatim from the record:

  ```text
  {"code":"flow-tag-invalid","message":"the value for `iter` does not
  conform to its declaration: \"many\" is not an int",
  "schema_version":"0.1","param":"iter"}
  ```

- [REQ-53] "Owned-reader parsed comparison — a reader returning `07` for `iter` against the guard `iter eq 7`, on the `mvv-int` model. **Expected**: the guarded row MATCHES (parsed 7 == 7) where today's raw-string seam prunes it." — (TS, `0012:S9`; A6 "S9's owned-reader fixture is the standing regression case for that direction")

- [REQ-54] "`flow next` derives its `unknown` pairs from the same probe seam (`internal/cli/flow_next.go::probeRow`), where the `guard_unevaluable` payload is the only source of `uncomparable`, so a malformed owned int is now reported `uncomparable` there too, and a non-canonical held value (`\"07\"` against `n eq 7`) moves a row between excluded and candidate." — (FM, Visible break; MC `disposition` row "Same, via `flow next`")

---

## EXCLUDED

- EXCLUDED: "First, C4's \"SAME loaded model\" half: scenario 4 asserts only that no zero-value construction survives, so a site constructing over the WRONG model's kinds still compiles and types comparisons silently (F4). There is no observable" — (TS, Known coverage gaps) — the record itself declares no oracle; guarded by review and the A3 site audit. REQ-33 keeps the testable half.
- EXCLUDED: "Second, the `enum`-vs-`scalar` distinction: they share one arm by construction (C2), so no input discriminates them" — (TS, Known coverage gaps) — no observable by construction; REQ-16 and REQ-21 pin what matters.
- EXCLUDED: "STEP ORDER. The unparseable-LITERAL arm is a per-kind obligation inside the `int` and `bool` arms, not a pre-dispatch gate above the kind switch." — (NC, `0012:C2`) — code-structure claim; every input it distinguishes already yields the verdict REQ-16/REQ-18/REQ-21 assert.
- EXCLUDED: "`GuardResult` ENCODING. … This RDR fixes no ordering and adds no verdict" — (NC, `0012:C2`) — a shape the clause reads and does not change; pinned by `0007:REQ-34`'s existing tests.
- EXCLUDED: "That encoding is those functions' to define and this clause makes no exactness claim of its own about it" — (NC, `0012:C2`, `in` MEMBERS) — peer-owned (§D13; `renderSet` / `renderSetLiteral`); pinned by the suite's existing `in/member` case.
- EXCLUDED: "Fixing the width is out of scope here" — (NC, `0012:C2`) — scope the record defers.
- EXCLUDED: "Undeclared-key ADMISSION policy — whether such a key may enter the view at all — is upstream and deliberately not decided here" — (NC, `0012:C2`) — deferred (JDR 0001 §JD-18 / RDR 0020).
- EXCLUDED: "This clause reaches the CONSTRUCTION and the CONSUMPTION of a verdict, never its AGGREGATION" — (NC, `0012:C4`) — aggregation is `0007:C6`/`C8`'s; nothing here changes it.
- EXCLUDED: "Dependency: `JDR 0004 §JD-3` obliges the loader-side admitted-cell shim RDR 0030's Phase 2 extracts to construct through `NewEvaluator`" — (NC, `0012:C4`) — the shim does not exist on `main`; RDR 0030's obligation. REQ-48's universal covers it once it lands.
- EXCLUDED: "The two excluded functions are also read by an open peer over the same seam; that overlap is adjudicated in §Load-Bearing Decisions and changes nothing in this clause." — (NC, `0012:C4`; LBD RDR 0022 entry) — coordination note, no observable.
- EXCLUDED: "Whether the comparison is factored into a private helper is implementation latitude" — (NC, `0012:C5`) — implementation latitude.
- EXCLUDED: "Whether the CLI should also demand canonical spellings is a separate question this RDR does not decide" — (NC, `0012:C5`) — deferred; REQ-42 pins the unchanged behavior.
- EXCLUDED: "A `--fix` affordance may be layered on later; it does not substitute for the refusal." — (NC, `0012:C5`) — deferred.
- EXCLUDED: "**Residual, stated not hidden.** A non-canonical int in an `[initial]`/`[rule.write]` cell still reaches `heldValues` → `canonicalValues` with no integer round-trip." — (LBD) — out of scope, carried by kata `intrastate#ch99`.
- EXCLUDED: "The narrower \"validate-then-byte-compare\" variant is the recorded fallback if A4's lint-agreement verification fails." — (LBD, Selection / predicate) — conditional on a failed assumption; A4 is Verified.
- EXCLUDED: "A6 is IN the gate, not an optional tail" — (IP, Prerequisites) — record-state precondition; A6 is Verified by census. Its standing regression is REQ-53.
- EXCLUDED: "It lives as a Go test under `internal/guard` — the package whose invariant it guards" — (TS, `0012:S4`) — code-siting claim; the implementer follows it, REQ-48 carries the observable.
- EXCLUDED: "Lands with S7's regression and after Phase 1, so the literal side is closed before typed comparison ships" — (IP, Phase 4) and "Phase 1 also carries `internal/guard`'s own test migration" — (IP, Phase 1) — phase sequencing, not behavior.
- EXCLUDED: "The PRODUCTION surface is small and counted … The TEST surface is not small — eighteen zero-value construction sites" — (Trade-offs, Consequences) — descriptive census.
- EXCLUDED: Alternatives 1–3 and Briefly Rejected (`0012:ALT1`–`ALT3`, `BR1`–`BR3`) — rejected designs. Their one negative consequence (no `Kind` on the atom) is REQ-8.
- EXCLUDED: `0012:JC1` joint-check and the JDR 0004 citations (§Decision Rationale) — cross-record coordination; no observable here.

---

## ASSUMPTIONS

- ASSUMPTION: MVV step 2's "the row, the atom and the reason are named in `message`" is met by the EXISTING per-atom finding messages (`findings[].message`, rendered by `internal/cli/flow_exec.go::undecidedFindings` as "rule `<id>`: the atom on `<key>` could not be decided (<reason>)", with flat `rule`/`key`/`operator`/`literal`/`block` fields), and the top-level `message` is unchanged. Grounded: the step's own constraint is "NO new top-level key", and `findings` is an existing key; NCP (REQ-38) routes a typed Unevaluable "through the existing `guard_unevaluable` refusal and its payload (`0007:C8`)"; RDR 0005's REQ-96 test (`internal/cli/flow_resolve_0005_test.go::TestReq96_GuardUnevaluableCarriesRowsAndFlatAtomFields`) fixes `findings[]` as the carrier. Rewriting the top-level message would be message churn the record neither needs nor asks for.
- ASSUMPTION: the suite's integer-overflow held value (REQ-25) is a spelling that overflows 64-bit `int` (so it overflows on every supported target), for example `"99999999999999999999"`. Grounded: REQ-14 lets the suite pin a verdict only where the target does not decide it.
- ASSUMPTION: REQ-21's unknown-kind-token case is built by passing `newSeam` a copy of `ContractKinds()` with one extra key mapped to a non-vocabulary token. The published fixture itself stays exactly the five tokens, so REQ-26's meta-check holds. Grounded: REQ-27 hands each caller a fresh copy, and REQ-28 has the suite build the seam itself.
- ASSUMPTION: REQ-19's empty-`Kind` population is pinned in two places. The suite has a distinctly named case whose key is absent from the mapping, because `DeclaredKinds` omits empty-`Kind` keys (REQ-3) and the seam therefore sees an absent key. An `internal/guard` test builds a `*table.Model` with an empty-`Kind` tag and asserts both halves: `DeclaredKinds` omits the key, and `NewEvaluator(DeclaredKinds(m))` answers GuardUnevaluable for `eq`/`in` over it. Grounded: A5 forbids the suite from naming `table` types.
- ASSUMPTION: C5 (REQ-39) covers every int-typed literal over an `int`-declared tag at `(*loader).atom`: the `eq` literal, each `in` member, and the `lt/lte/gt/gte` bound, in guard and match blocks alike. Grounded: "A value authored against an `int`-declared tag"; "guard and match atom literals"; MC `fidelity` row "Authored PREDICATE int literal". `(*loader).atom` is block-agnostic (A1).
- ASSUMPTION: for `in`, every member is checked. The diagnostic names the offending member and its rewrite, taking the first in authored order when the loader reports a single error. Grounded: REQ-43's example names one value, and the loader's existing `conform` path reports one error per atom.
- ASSUMPTION: `<site>` in the C5 diagnostic is the same location prefix that the loader's existing `malformed_predicate_atom` errors for that atom already carry. Grounded: REQ-43 says "returns the same shape" as `conform`'s error. S7 pins only the message tail after the site.
- ASSUMPTION: in the `bool` arm, an `in` literal is decoded as its §D13 array and every member must be a boolean token. One non-token member poisons the list, by the same parse-all-then-compare rule REQ-13 states for `int`. Grounded: REQ-15 says "held value and literal must be the boolean tokens", and REQ-24 says that for `in` the literal means its decoded members.
- ASSUMPTION: the narrowed REQ-34 reflection check (REQ-47) walks `guard.Evaluator`'s fields. It admits the declaration mapping (`map[string]string`) and fails on any field whose type is or contains `resolve.TagSet`, `resolve.Tag`, a view type, or other runtime-valued state. Grounded: REQ-5, "holds this mapping and NOTHING else".
- ASSUMPTION: "the user-facing docs" in REQ-37 means `docs/model-authoring.md`, the hand-authored guide that already documents the operator/kind matrix and the match/guard block split. It gains a statement of the split: guard `eq`/`in` compare under the declared kind, and match atoms byte-compare. `docs/cli-reference.md` and `llms.txt` change only through `make docs` if a command help body changes. Grounded: AGENTS.md §Docs.
- ASSUMPTION: `guardSeam` becomes `guardSeam(m *table.Model) resolve.GuardEvaluator`, following the Illustrative Code. Its callers in `flow_resolve.go` pass the request's model, including at `escapeClassOf`, which A3 says resolves a probe built from the same model. `flow next` builds the evaluator once per request and passes it to `probeRow`. The `internal/cli` test call sites that use `guardSeam()` (`escape_shape_0009_test.go`, `escape_shape_adv_0009_test.go`) migrate with the signature. Grounded: REQ-34 and A3.
- ASSUMPTION: at `internal/guard/product.go::valueSatisfies`, the evaluator is constructed from the `*table.Model` that `Denotation(m, …)` holds and reaches `valueSatisfies` as the "one added local parameter" A3 names. It is not rebuilt per atom. Grounded: A3 and C4's "one model, one evaluator".

---

## QUESTIONS

None.

Four clauses were examined as ambiguity candidates, and each one resolved against the record's own text or a predecessor artifact. None is recorded as a QUESTION:

1. **Where MVV step 2 names the row, atom and reason** (top-level `message` versus the existing `findings[]` messages). Resolved as the first ASSUMPTION, from the step's own "NO new top-level key", from NCP's "through the existing … payload", and from RDR 0005's REQ-96 carrier.
2. **Overflow case versus "the suite pins the in-range cases only"** (REQ-14 against REQ-25). These are reconciled by choosing a spelling that overflows on every target (second ASSUMPTION).
3. **Unknown-kind case versus the "exactly the five-kind vocabulary" meta-check** (REQ-21 against REQ-26). These are reconciled because the suite builds its seam over a caller-local copy (third ASSUMPTION).
4. **C5's reach across operators** (only `eq`/`in`, or every int literal). Resolved by C5's own subject, "a value authored against an `int`-declared tag", and by the `fidelity` row (fifth ASSUMPTION). The narrower reading would leave `gte = "01"` admitted, and the literal-side agreement A4 depends on only concerns `eq`/`in`. Either reading therefore keeps A4 intact, and the wider one is the record's wording.
