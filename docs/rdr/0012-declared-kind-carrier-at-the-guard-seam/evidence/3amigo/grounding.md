Model: claude-opus-5[1m]

Grounding pass — verdicts against `main` (afcc69d). Read-only.

## CLAIM 1 — the `[emit]` int-canonicality contradiction: **CONFIRMED (REAL)**

### 1a — `TestReq9_0024_NonCanonicalIntLiteralsAreAdmitted`: CONFIRMED
Anchor: `internal/table/emit_grammar_0024_test.go::TestReq9_0024_NonCanonicalIntLiteralsAreAdmitted` (L263; marker `// BOUNDARY` L262).
Loops `[]string{"03", "+5", "-0"}` — exact literals claimed — authoring `count = "<lit>"` under `[emit.count] kind = "int"`, and calls `loadSource` (the LOAD-succeeds helper), not `refuseSource`. Asserts admission. Comment: "`strconv.Atoi` admits all three, so reuse admits them too."

### 1b — `TestReq7_0024_EveryKindCheckIsLexicalAndTheAuthoredBytesSurvive`: CONFIRMED
Anchor: `internal/table/emit_grammar_0024_test.go::TestReq7_0024...` (L200; marker `// DOMAIN EDGE` L199).
Authors `code = "007"` under `[emit.code] kind = "int"`, loads, then asserts `e.Value == "007"` on the row's `Emit` — fails with "no value is parsed into a typed representation, canonicalized, or converted" otherwise. Authored bytes survive verbatim.

### 1c — does `[emit]` reach `conformKind` with an int decl?: **CONFIRMED** (trace exact)
Anchor: `internal/table/load.go::(*loader).checkRuleEmit` L497 + L506-507.
- L497: `if decl.Kind == "scalar" { continue }` — scalar short-circuited BEFORE the call, as claimed.
- L506: `ConformValue(TagDecl{Kind: decl.Kind, Domain: decl.Domain}, rule.Emit[key])` — throwaway `TagDecl` adapter, `Min`/`Max`/`Elements` left nil, as claimed.
- `ConformValue` (`load.go::ConformValue` L1813) calls `conformKind(decl, member)` L1814 unconditionally, then `conformDomain`.
- `conformKind` (`load.go::conformKind` L1853) int arm: `if _, err := strconv.Atoi(member); err != nil { … "is not an int" }`.
An int-kinded `[emit]` value DOES reach `conformKind`'s int arm. Only `scalar` escapes.

### 1d — same code path?: **SAME PATH. The contradiction is REAL, not apparent.**
Both tests author `kind = "int"` (not `scalar`), so L497 does not skip them; both flow through `checkRuleEmit` -> `ConformValue` -> `conformKind` int arm. The arm's admission test is `strconv.Atoi`, which accepts `03`, `+5`, `-0`, `007`. A rule requiring `strconv.Itoa(strconv.Atoi(s)) == s` would refuse all four literals and break both tests.
Verified empirically: both tests PASS on `main` (`go test ./internal/table -run 'TestReq9_0024_NonCanonicalIntLiteralsAreAdmitted|TestReq7_0024_...'` -> PASS, incl. subtests `03`, `+5`, `-0`). The proposed canonicality rule is refuted by live, passing tests on the same arm it would change.
Note: 1b compounds it — even if canonicality only gated admission, 1b independently pins that the AUTHORED bytes (`007`) ride onto the row uncanonicalized.

## CLAIM 2 — `guard.Evaluator` shape + `eq` arm: **CONFIRMED**
Anchor: `internal/guard/grammar.go::Evaluator` L97 `type Evaluator struct{}`; `grammar.go::Evaluator.Evaluate` L104 `func (Evaluator) Evaluate(atom resolve.GuardAtom, value string) resolve.GuardResult` — value receiver, unnamed.
`eq` arm L106-107: `return boolResult(value == atom.Literal)` — raw string comparison, no declaration consulted.

Zero-value `Evaluator{}` literals in `internal/guard`:
- TEST files: **11** across 3 files — `guard_evaluator_0003_test.go` (7: L84, L104, L124, L149, L164, L186, L220), `guard_narrowing_0003_test.go` (3: L142, L163, L568), `guard_fixtures_0003_test.go` (1: L294).
- NON-test files: **1** — `internal/guard/product.go:542` `seam := Evaluator{}`.
(Also outside the package: `internal/cli/flow_resolve.go:561` returns `guard.Evaluator{}`.)

## CLAIM 3 — `gte`/`lt`/`lte`/`gt` and `contains` are operator-inferred: **CONFIRMED**
Anchor: `internal/guard/grammar.go::Evaluator.Evaluate` L114-123 (comparisons), L124-140 (`contains`).
Comparisons: `strconv.Atoi(atom.Literal)` + `strconv.Atoi(value)`, unevaluable on either parse failure, else `compare(op, held, bound)`. `contains`: `parseSetLiteral(atom.Literal)` + `parseHeldSet(value)`, then membership. Neither consults any declaration, kinds map, or key — the operator alone fixes the reading. `Evaluator` is `struct{}`, so it holds no declaration state to consult.
A zero-value `Evaluator{}` with a nil kinds map: these five operators still answer correctly (they never read a map). `eq` (raw `==`) and `in` (`parseSetLiteral` + string `slices.Contains`) are likewise declaration-free TODAY, so nothing degrades on today's code — they would be the only arms that COULD degrade under a declared-kind design, since they are the arms whose reading a declared kind would change.

## CLAIM 4 — `TestGuardEvaluatorContract`: **CONFIRMED**
Anchor: `internal/resolve/guardcontract.go::TestGuardEvaluatorContract` L14 — signature exactly `func TestGuardEvaluatorContract(t *testing.T, seam GuardEvaluator)`. Non-test source file (so it is exported for cross-package reuse).
`Key: "subject"` appears ONCE (L60), in the single shared `GuardAtom` literal built inside the case loop — so every case carries `Key: "subject"`. Claim holds; the hardcoding is one site, not per-case.
Case counts: `gte` = **5** (`above`, `equal`, `below`, `unparseable value`, `unparseable literal`). `contains` = **4** (`superset`, `missing member`, `unparseable value`, `null value`). (Also `eq` = 2, `in` = 2; 13 cases total.)
In-tree call sites (both in `internal/resolve/guard_mvv_test.go`): L462 `TestReq71_AConformingValueSeamSatisfiesTheContractTest` passes `conformingContractSeam{}`; L256 `TestReq71_ContractTestExercisesThePresentValueLiteralOperatorProduct` passes `rec` (a `*recordingSeam`). No `internal/guard` call site — the contract test is NOT instantiated against `guard.Evaluator`.

## CLAIM 5 — `conformingContractSeam` + `recordingSeam`: **CONFIRMED**
Anchor: `internal/resolve/guard_mvv_test.go::conformingContractSeam` L489 `type conformingContractSeam struct{}`; `Evaluate` L491, `opEq` arm L493-497 does `if value == atom.Literal` — raw string comparison, mirroring `guard.Evaluator`'s.
Wrapped: yes — `guard_mvv_test.go:256` `rec := &recordingSeam{inner: conformingContractSeam{}}`. `recordingSeam` (L473) delegates to `inner` and appends `calls`/`verdicts`.
Consumer: `TestReq71_ContractTestExercisesThePresentValueLiteralOperatorProduct` (L243) — it feeds `rec` to `resolve.TestGuardEvaluatorContract` (L257) and then asserts on `rec.calls`/`rec.verdicts`, i.e. it tests the CONTRACT TEST's own coverage (that it exercises present value x literal x operator, every call carrying a non-empty operator), using a conforming seam so the inner contract run still passes.

## CLAIM 6 — `atom` / `atomsFromBlock` in `internal/table`: **PARTIAL (refuted as stated, symbols exist with different form)**
- A FREE function `atom` in `normalize.go`: **REFUTED.** It is a METHOD: `internal/table/normalize.go::(*loader).atom` L43 — `func (l *loader) atom(decl TagDecl, key, operator string, raw any, b Block, owner string) (Atom, error)`.
- `atomsFromBlock` in `internal/`: **CONFIRMED, also a method** — `internal/table/normalize.go::(*loader).atomsFromBlock` L23.
The function calling `conform(decl, operator, members)` is `(*loader).atom` at `normalize.go:172`. Other `conform` call sites are `normalize.go:624` and `load.go:1701`, both `conform(decl, "eq", members)`. `conform` itself: `internal/table/load.go::conform` L1826.

## CLAIM 7 — `cli::guardSeam`: **CONFIRMED (zero-arity)**
Anchor: `internal/cli/flow_resolve.go::guardSeam` L561 — `func guardSeam() resolve.GuardEvaluator { return guard.Evaluator{} }`. Arity 0 in, 1 out.
Call sites: **6** total — non-test `flow_resolve.go:258`, `flow_resolve.go:445`, `flow_next.go:492`; test `escape_shape_0009_test.go:83`, `escape_shape_0009_test.go:624`, `escape_shape_adv_0009_test.go:40`, `escape_shape_adv_0009_test.go:159`. (That is 3 non-test + 4 test = 7 call sites; `flow_next.go:435` is a prose comment reference, not a call.)

## CLAIM 8 — `valueMembers` float64 arm: **CONFIRMED**
Anchor: `internal/table/load.go::valueMembers` L1748, float64 arm L1757 — `return []string{strconv.FormatFloat(t, 'g', -1, 64)}, nil`. Exact as claimed.
`-0.0` -> `"-0"`: **CONFIRMED empirically.** `pelletier/go-toml/v2` decodes `v = -0.0` to a `float64` negative zero, and `strconv.FormatFloat(-0.0, 'g', -1, 64)` yields `"-0"`. Both legs verified by execution.
Bearing on CLAIM 1: `"-0"` is exactly one of the literals `TestReq9` admits, so a float-authored `-0.0` and an int-kinded `"-0"` converge on the same non-canonical string a canonicality rule would refuse.
