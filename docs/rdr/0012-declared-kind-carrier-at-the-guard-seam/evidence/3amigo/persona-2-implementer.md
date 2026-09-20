Model: claude-opus-5[1m]

# 3amigo — Persona 2: Implementer — cli/0012

## Widened

Yes, three times.

- C3 names a fixture the suite "publishes" but not how a caller gets it,
  and `TestGuardEvaluatorContract`'s signature is fixed source I had to
  read; that sent me to §implementation-plan (Phase 2) and to
  `internal/resolve/guardcontract.go`.
- C1's `NewEvaluator(kinds map[string]string)` names no producer of the
  map; §illustrative-code's `guard.DeclaredKinds(m)` is the only spelling
  of it anywhere, and it is stamped "not load-bearing". That sent me to
  §existing-infrastructure-audit (row 1, "behind a small kinds-extraction
  helper") and to `internal/guard/declaration.go`.
- C4's three-site enumeration is signature work, not a call-site edit, so
  I widened to A3 (§critical-assumptions) for the frame analysis and to
  F4/F6 (§failure-modes) for the migration disposition.

## Source-anchor check

Every `source-anchor` edge in `recs inspect --json` reports
`resolved: true`; I re-checked each named symbol against `main` anyway.

- `internal/guard/grammar.go::Evaluator` — EXISTS, `type Evaluator struct{}`
  (grammar.go:97), value receiver on `Evaluate` (grammar.go:104). Matches
  C1's "zero-value `struct{}`".
- `internal/guard/grammar.go::parseHeldSet` — EXISTS (grammar.go:179).
  Audit row is accurate: answers unevaluable without declarations.
- `internal/guard/declaration.go::DeclarationOf` — EXISTS
  (declaration.go:47), `func DeclarationOf(m *table.Model, key string) Declaration`.
  Audit's "takes `*table.Model`" limit is accurate.
- `internal/guard/declaration.go::Conforms` — EXISTS (declaration.go:280).
- `internal/guard/assignment.go::valueAssignments` — EXISTS
  (assignment.go:273); int arm renders via `strconv.Itoa(n)`, which is
  A4's claimed canonical render. Confirmed.
- `internal/guard/product.go::valueSatisfies` — EXISTS (product.go:541);
  constructs `seam := Evaluator{}` at product.go:542. C4 site 2 confirmed.
- `internal/cli/flow_resolve.go::guardSeam` — EXISTS (flow_resolve.go:561)
  but its signature is `func guardSeam() resolve.GuardEvaluator`, ZERO-ARG.
  §illustrative-code writes `func guardSeam(m *table.Model)`. Illustrative,
  so not a contract defect — but see IMPL-2.
- `internal/graphlint/reach.go::atomAdmitsValue` — EXISTS (reach.go:502);
  constructs `guard.Evaluator{}.Evaluate(...)` inline at reach.go:507. C4
  site 3 confirmed.
- `internal/graphlint/reach.go::matchSatisfiable` — EXISTS (reach.go:449),
  carries `m *table.Model`, as A3 states.
- `internal/cli/flow_next.go::probeRow` — EXISTS (flow_next.go:469); takes
  `row table.Row`, no model. A3's "widens `probeRow`'s signature" is
  accurate.
- `internal/table/load.go::conformKind` — EXISTS (load.go:1853); int arm is
  `strconv.Atoi` only. C5's stated starting point is accurate.
- `internal/table/load.go::valueMembers` — EXISTS (load.go:1748); float64
  arm is `strconv.FormatFloat(t, 'g', -1, 64)`. C5 says the TOML `-0.0`
  renders to `"-0"` there; `'g'` formatting of `-0.0` does yield `"-0"`.
  Confirmed.
- `internal/table/normalize.go::atom` — NO SUCH FREE FUNCTION. The
  atom-building code is in `normalize.go` and calls
  `conform(decl, operator, members)` at normalize.go:172, but there is no
  `func atom(` in the package (only `atomIdentity` at :219). C5 and the
  audit both spell it `normalize.go::atom`. See IMPL-6 (low).
- `internal/table/normalize.go::atomsFromBlock` — NO SUCH FUNCTION on
  `main`. `rg "func atomsFromBlock"` over `internal/` returns nothing. It
  appears only as a source-anchor edge, not in C1–C5 text. Low impact for
  this persona (no contract rests on it) but it is a stale anchor.
- `internal/table/normalize.go::renderWrites` — EXISTS as a METHOD,
  `func (l *loader) renderWrites(...)` (normalize.go:580), and it does
  already call `conform(decl, "eq", members)` (normalize.go:624), exactly
  as C4 says. Confirmed.
- `internal/cli/flow_input.go::canonicalValue` — EXISTS (flow_input.go:738),
  reaches `table.ConformValue` (:751 for scalars, :769 for set members).
  C5's reach claim confirmed.
- `internal/table/load.go::ConformValue` — EXISTS (load.go:1813), calls
  `conformKind` then `conformDomain`. Confirmed.
- `internal/resolve/guardcontract.go::TestGuardEvaluatorContract` — EXISTS
  (guardcontract.go:14), signature `func TestGuardEvaluatorContract(t *testing.T, seam GuardEvaluator)`.
  Every case hardcodes `Key: "subject"` — C3's "Today every case is built
  with `Key: \"subject\"`" is accurate.
- `internal/resolve/guard_mvv_test.go::conformingContractSeam` — EXISTS
  (guard_mvv_test.go:489), `type conformingContractSeam struct{}`,
  raw-string `eq`, driven at :462 and wrapped in `recordingSeam` at :256.
  C3's claim about both drivers is accurate.
- `internal/guard/guard_evaluator_0003_test.go::TestReq34_...` — EXISTS
  (:21); the zero-field assertion is at :46
  (`reflect.TypeOf(ev).NumField() != 0`). A2's "structural proxy in the
  test body" is accurate.
- `internal/resolve/resolve.go::assemble` — EXISTS (resolve.go:220).
- `internal/resolve/resolve.go::TagSet.matches` — present via `TagSet`;
  not load-bearing for C1–C5.
- `internal/accessor/model.go::ReadResult.OwnedSnapshot`,
  `internal/table/model.go::KernelTable` — present; not load-bearing here.

No high-severity anchor defect: every symbol C1–C5 rests on exists as
described. The two misses (`normalize.go::atom`, `atomsFromBlock`) are in
naming only and neither is a contract's load-bearing site.

## Findings

### IMPL-1 — C3 fixes the suite's fixture but not its delivery mechanism, and the suite's signature has no seat for it — severity: high
- anchor: 0012:C3 (widened to 0012:§phase-2-conformance-suite-extension,
  `internal/resolve/guardcontract.go:14`)
- question: How does a caller obtain the published fixture, and how does
  the suite get its per-case keys? C3 says the suite "publishes the
  declaration fixture its cases assume as a kernel-owned
  `map[string]string`" and "its doc states the caller obligation:
  construct the seam under test over that fixture before invoking the
  suite." That names an invariant the compiler cannot enforce. Concretely
  I cannot choose between: (a) a new exported
  `func GuardContractKinds() map[string]string` beside
  `TestGuardEvaluatorContract`, with the existing two-arg signature kept
  and the obligation left as prose; (b) widening to
  `TestGuardEvaluatorContract(t *testing.T, newSeam func(kinds map[string]string) GuardEvaluator)`
  so the suite constructs the seam itself and the obligation becomes
  structural; (c) an exported package var. C3's own wording ("construct
  the seam under test … before invoking") reads as (a) or (c), which
  means a caller that constructs over the WRONG map still compiles and
  the suite reports a semantic failure rather than a wiring one.
- blocks: the exported API shape of `internal/resolve/guardcontract.go` —
  a cross-package, non-test source file with two in-tree callers
  (`guard_mvv_test.go:462` and the 0003 implement stage). Choosing (b)
  breaks both call sites and is a different diff from (a). I would have
  to guess, and the guess is not reversible cheaply once 0003's
  implementer is driving it.

### IMPL-2 — C1 names the constructor's parameter but nothing names its producer — severity: high
- anchor: 0012:C1 (widened to 0012:§existing-infrastructure-audit row 1
  and 0012:§illustrative-code)
- question: What is the signature of the thing that turns a `*table.Model`
  into `map[string]string`, and which package owns it? C1 fixes
  `func NewEvaluator(kinds map[string]string) Evaluator`. The only
  spelling of a producer anywhere in the record is
  `guard.DeclaredKinds(m)` in §illustrative-code, explicitly stamped
  "Illustrative — shape only, not load-bearing"; the audit calls it "a
  small kinds-extraction helper" behind `DeclarationOf`. Open: is it
  exported (`guard.DeclaredKinds(m *table.Model) map[string]string`) or
  unexported? `internal/graphlint` needs it too (`atomAdmitsValue`,
  reach.go:507, is a different package), so an unexported helper does not
  serve all three C4 sites; that forces export, which the record never
  states. And does it key EVERY entry of `m.Tags`, or only tags with a
  non-empty `Kind`? Declared five-kind tokens are `enum|bool|int|set|scalar`
  (`internal/table/model.go:83`), but `TagDecl.Kind` is a bare `string` and
  a zero `TagDecl` has `Kind == ""` — an empty-string kind entering the map
  hits C2's `default` arm as an undeclared key, which is the same verdict
  as an ABSENT key, so the two are indistinguishable in a diagnostic.
- blocks: whether a new exported symbol lands in `internal/guard` (a public
  surface decision with its own doc and test obligations) versus three
  copies of a two-line loop, and whether the map is total over `m.Tags` or
  filtered. Both are first-hour choices and both are visible in the diff.

### IMPL-3 — C1's nil-map disposition is unreachable as written on a value receiver — severity: medium
- anchor: 0012:C1
- question: How is `Evaluator{}`'s LOUD disposition implemented, and does
  it survive the Go type system? C1 says "Go still admits a bare
  `Evaluator{}` literal — its disposition is fixed LOUD: a nil-mapping
  evaluator answers GuardUnevaluable for every `eq`/`in` atom". A nil map
  read in Go returns the zero value, so `ev.kinds[atom.Key]` on a nil map
  yields `""`, which falls to C2's "key absent from the mapping" arm →
  GuardUnevaluable. That works — for `eq`/`in`. But it does NOT hold for
  the other operators: `Evaluator{}` still answers `gte`/`contains`
  correctly under the operator-inferred parse, so a missed construction
  site with only integer-comparison or set guards migrates SILENTLY.
  F6 ("Zero-value evaluator survives migration") claims the miss "surfaces
  as loud `guard_unevaluable` refusals on previously-deciding guards" —
  true only of models whose guards use `eq`/`in`. Is that residual
  accepted, or is the intent a nil-check panic / a `lt/gte` arm that also
  consults the mapping?
- blocks: whether Phase 1 adds a nil-mapping guard clause at the top of
  `Evaluate` (changing `gte`/`contains` behavior for `Evaluator{}`, which
  the REQ-34 test at guard_evaluator_0003_test.go:26 and eight other test
  sites currently rely on) or leaves the arms untouched. The two produce
  different test-migration diffs across `internal/guard`'s eleven
  `Evaluator{}` literals.

### IMPL-4 — C2's `in` arm does not say whether an unparseable MEMBER poisons the whole list or only itself — severity: medium
- anchor: 0012:C2 (cross-checked against 0012:C3)
- question: For `in` over an `int` key with literal `["1","two","3"]` and
  held value `"1"`, is the answer GuardTrue (the matching member parsed
  fine) or GuardUnevaluable (one member did not)? C2's int bullet says
  "held value AND literal (each member, for `in`) must parse as integers"
  — that reads as all-members, i.e. GuardUnevaluable. C3 then asks for a
  case "an `in` list with a non-integer member on an `int` dimension
  (want GuardUnevaluable)" but does not say whether the held value matches
  a parseable member, so the case is authorable either way and does not
  disambiguate. Worse: C5 makes `"01"` a load refusal, so an authored
  literal can no longer be non-canonical — meaning the only way to reach
  this arm is a caller-built atom, which is exactly what
  `product.go::valueSatisfies` and `reach.go::atomAdmitsValue` do
  (they re-render literals through `renderSet`/`renderSetLiteral`).
- blocks: the loop structure of the `in` int arm — parse-all-then-compare
  versus parse-and-compare-per-member — and the exact expectation string
  in the new C3 case. Picking the permissive reading makes lint and the
  runtime disagree on a mixed list, which is the disagreement D-canonical-spelling-at-load
  exists to prevent; picking the strict one is a behavior change I cannot
  attribute to a written contract.

### IMPL-5 — C5 names four refusal categories but only three ingresses demonstrably reach `conformKind` with an int decl — severity: medium
- anchor: 0012:C5 (widened to 0012:§phase-4-canonical-int-spelling-at-load)
- question: Does the `[emit]` ingress actually reach C5's tightened check,
  and if so with what `TagDecl`? C5 lists `emit_value_out_of_domain` as
  one of the reused categories. The emit path (load.go:506) builds a
  THROWAWAY decl — `TagDecl{Kind: decl.Kind, Domain: decl.Domain}` — from
  an `EmitDecl`, whose documented kind vocabulary is four tokens
  (`enum | bool | int | scalar`, model.go:218), and short-circuits
  `scalar` before the call (load.go:497). So an int-kinded emit value DOES
  reach `conformKind` and WILL start refusing `"00"`. Is that intended?
  C5's rationale rests on S5's corpus measurement of zero `eq`/`in` atoms
  over `int` tags — a measurement about GUARD ATOMS, not about `[emit]`
  values, `[initial]` values, or `[rule.write]` values, which are three
  other populations the same tightening silently covers. Nothing in C5,
  S5, or S6 states those populations are also empty of non-canonical int
  spellings.
- blocks: whether Phase 4 is a one-line change to `conformKind`'s int arm
  (C5's stated home) or needs a per-ingress opt-in, and whether I should
  expect existing corpus/testdata models to start failing. I cannot tell
  from the record whether a green `make check` after Phase 4 is the
  expected outcome or evidence I scoped it wrong.

### IMPL-6 — C5 and the audit name `internal/table/normalize.go::atom`, which does not exist — severity: low
- anchor: 0012:C5, 0012:§existing-infrastructure-audit (last row)
- question: Which function is meant? C5 attributes
  `malformed_predicate_atom` to "guard and match atoms, via
  `internal/table/normalize.go::atom`", and the audit repeats "block-agnostic
  `normalize.go::atom`". There is no `func atom(` in `internal/table`. The
  behavior described is real — `normalize.go:172` calls
  `conform(decl, operator, members)` and the failure path is `badAtom(...)`
  — so the contract's SUBSTANCE holds; only the symbol name is wrong. The
  companion anchor `normalize.go::atomsFromBlock` likewise resolves to
  nothing on `main`.
- blocks: nothing structural — I would find the right call site in a
  minute. It costs a first-hour minute and it weakens my trust in the
  other anchors, which is why it is recorded rather than dropped.

### IMPL-7 — C4 fixes the three sites but not how the model reaches two of them, and A3 offers a CLI fork without choosing — severity: medium
- anchor: 0012:C4 (widened to 0012:A3, §critical-assumptions)
- question: For the CLI site, do I widen `probeRow`'s signature to take a
  `*table.Model`, or thread the already-constructed `resolve.GuardEvaluator`
  down? A3 says "the CLI additionally widens `probeRow`'s signature or
  passes the constructed evaluator down" — an explicit either/or that no
  contract resolves. The two differ observably: threading the evaluator
  means `guardSeam` is called ONCE per request and the same instance
  reaches all three call sites (flow_resolve.go:258, :445,
  flow_next.go:492), which is what "one model, one evaluator" in C4 most
  naturally means; widening `probeRow` means a fresh construction per row
  probe. `probeRow` is called in a loop, so the second is also a
  per-row map-build. Also unresolved: `guardSeam()` is referenced by four
  CLI TEST files (escape_shape_0009_test.go:83, :624,
  escape_shape_adv_0009_test.go:40, :159) that call it zero-arg — those
  migrate either way, but the migration differs.
- blocks: the shape of the `internal/cli` diff in Phase 3, and whether
  `guardSeam` survives as a helper at all.

### IMPL-8 — Phase 1 retires zero-value construction but nothing states the disposition of the eleven in-tree test literals — severity: low
- anchor: 0012:C1 ("The zero-value `Evaluator{}` construction form is
  retired; every non-test construction goes through `NewEvaluator`"),
  0012:§phase-3-producer-alignment ("test fixtures migrate to the
  constructor")
- question: Which is it for `internal/guard`'s test literals — Phase 1 or
  Phase 3? There are eleven `guard.Evaluator{}` literals in test files
  (guard_evaluator_0003_test.go ×7, guard_fixtures_0003_test.go ×1,
  guard_narrowing_0003_test.go ×3). C1 exempts tests from the constructor
  obligation ("every NON-TEST construction"); Phase 3 says "test fixtures
  migrate to the constructor". But Phase 1 ships the typed `eq`/`in` arms
  — so between Phase 1 and Phase 3 those eleven literals evaluate `eq`
  atoms against a nil mapping and answer GuardUnevaluable, and any test
  asserting GuardTrue/GuardFalse on an `eq` goes red. Does Phase 1 have to
  carry the test migration to stay green, or is a red intermediate
  accepted?
- blocks: whether Phase 1 is a self-contained commit that passes
  `make check` — which the project's stated convention ("After changing Go
  code, run `make check` before declaring done") makes a real question,
  not a bookkeeping one.
