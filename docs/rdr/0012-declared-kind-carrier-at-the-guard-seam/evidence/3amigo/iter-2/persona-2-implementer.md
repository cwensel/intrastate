Model: claude-opus-5[1m]

# 3amigo iter-2 (delta) — Persona 2: Implementer — cli/0012

## Resolved by the rewrite

- **Producer identity (was: "who builds the map, and can three packages
  call it?").** C1 now names `func DeclaredKinds(m *table.Model)
  map[string]string`, states WHY it is exported (three C4 sites in three
  packages), fixes it as TOTAL over `m.Tags`, and fixes the empty-`Kind`
  disposition as OMIT rather than map-to-`""`. `§existing-infrastructure-audit`
  seats it over the existing `DeclarationOf`. Nothing left to ask.
- **`in`-over-`int` mixed-list disposition (was: "does one bad member
  poison the list?").** C2 now says parse-all-then-compare, one
  unparseable member ⇒ GuardUnevaluable even on an equal parsed member,
  with the lint/runtime-divergence rationale. Implementable as written.
- **`enum` vs `scalar` (was: "two arms or one?").** C2 explicitly shares
  one arm and states why the suite's per-kind-token requirement is still
  satisfiable.
- **`default` semantics (was: "is the switch exhaustive, and what does
  fallthrough mean?").** C2 fixes `default` = GuardUnevaluable and names
  the failure mode of the string-equality alternative.
- **Operator-inferred vs declared-kind competition (was: "does re-keying
  `contains` onto a `set` key make it GuardUnevaluable via the `set`
  arm?").** C2 now says lt/lte/gt/gte and `contains` never consult the
  kind, and the `set` arm governs `eq`/`in` ONLY. This was the single
  highest-risk ambiguity in iter-1 and it is now closed in the normative
  text, not the prose.
- **Suite delivery (was: "how does the suite get the right kinds without
  a prose obligation?").** C3 widens to a constructor parameter —
  structural, so miswiring does not compile.
- **The real implementer was never gated (was: "is `guard.Evaluator`
  actually run against the contract?").** C3 now ADDS
  `TestGuardEvaluatorContract(t, guard.NewEvaluator)`. Verified: on
  `main` it genuinely is not a caller.
- **CLI fork (was: "widen `probeRow` or thread the seam?").** C4 decides:
  thread the constructed evaluator down, with the per-row-rebuild
  argument. A3 left this an either/or; C4 now closes it.
- **C5 blast radius (was: "does tightening the int arm break `[emit]`/
  `[initial]`/`[rule.write]`?").** C5 narrowed to the predicate ingress,
  `conformKind` left shared and unchanged, the two 0024 BOUNDARY tests
  named. Verified below.
- **Phase-boundary redness (was: "Phase 1 lands typed arms; when do the
  `Evaluator{}` test literals migrate?").** Phase 1 now explicitly carries
  internal/guard's 11 literals and says why deferring to Phase 3 would
  leave `main` red.

## New-claim verification

- `DeclaredKinds` collides with nothing — `rg DeclaredKinds` over the
  whole repo returns ZERO hits. internal/guard's nearest neighbours are
  `Kinds() []string` (grammar.go:65, the 5-token vocabulary) and
  `DeclarationOf` / `DeclarationFields` (declaration.go:47/42). No
  shadowing, no near-miss that would confuse a reader. HOLDS.
- `DeclaredKinds(m *table.Model) map[string]string` is buildable as
  specified: `table.Model.Tags` is `map[string]TagDecl` (model.go:519)
  and `TagDecl` carries `Kind string` (the `conformKind` switch at
  load.go:1854 reads `decl.Kind`). A total loop with an empty-`Kind`
  skip is a five-line function. HOLDS.
- The widened `TestGuardEvaluatorContract` signature is compatible with
  its call sites — but there are THREE in-tree call sites, not two, and
  C3's own sentence is internally inconsistent about the count. See
  IMPL-D1. The three are guard_mvv_test.go:257 (`recordingSeam` wrapping
  `conformingContractSeam`), guard_mvv_test.go:462 (bare
  `conformingContractSeam{}`), and guard_evaluator_0003_test.go:67
  (`resolve.TestGuardEvaluatorContract(t, ev)`). PARTIALLY HOLDS.
- `conformingContractSeam` is `struct{}` (guard_mvv_test.go:489) with a
  value receiver on `Evaluate` (:491). Adding a `kinds map[string]string`
  field plus a constructor is mechanical; the value receiver still works
  (maps are reference types). `recordingSeam` (:473) holds
  `inner resolve.GuardEvaluator` and delegates at :479, so it needs no
  kinds of its own exactly as C3 says. HOLDS.
- The REQ-71 signature assertion at guard_mvv_test.go:222 —
  `var fn func(*testing.T, resolve.GuardEvaluator) = resolve.TestGuardEvaluatorContract`
  — is a COMPILE-TIME type assertion that the widening breaks. See
  IMPL-D2. C3's "breaks both in-tree call sites … and they migrate in
  the same phase" does not mention it, and it is not a call site: it is a
  normative signature pin from 0007. DOES NOT HOLD as stated.
- Threading the evaluator into `probeRow` works: `probeRow(row
  table.Row, view map[string]string, owned []resolve.Tag, observed
  []resolveTag, all bool) resolve.Result` (flow_next.go:469) has exactly
  one non-test caller, flow_next.go:259, inside `for _, row :=
  range req.model.Rows` — so `req.model` IS in scope one frame up,
  hoisting the construction above the loop is trivial, and C4's
  "rebuild per row" argument is factually the right one. Two test
  callers at flow_next_0011_test.go:806 and :861. HOLDS.
- `guardSeam()` is zero-arg today (flow_resolve.go:561) with three
  non-test call sites (flow_resolve.go:258, :445, flow_next.go:492) and
  four TEST call sites in escape_shape_0009_test.go (:83, :624) and
  escape_shape_adv_0009_test.go (:40, :159). C4/Phase 3 say
  `internal/cli` test fixtures migrate — those four are the migration.
  HOLDS, but the arity change ripples further than A3's prose implies
  (see IMPL-D3).
- `(*loader).atom` has decl/operator/members all in scope at
  normalize.go:172: the signature is `func (l *loader) atom(decl
  TagDecl, key, operator string, raw any, b Block, owner string) (Atom,
  error)` (normalize.go:43), `members` is derived above the operator
  switch, and :172 is literally `if err := conform(decl, operator,
  members); err != nil {`. `badAtom` is the in-scope refusal helper
  used throughout the function. Adding C5's `Itoa` round-trip beside
  :172 is a drop-in. HOLDS — the corrected method spelling is right,
  and `atomsFromBlock` is a method at :23 as stated.
- The 11 test literals are exactly where named: `guard_evaluator_0003_test.go`
  = 7, `guard_narrowing_0003_test.go` = 3, `guard_fixtures_0003_test.go`
  = 1. Total 11. HOLDS exactly.
- The three C4 non-test construction sites are exactly the three named:
  `rg 'Evaluator{}' -g '!*_test.go'` returns product.go:542
  (`valueSatisfies`), reach.go:507 (`atomAdmitsValue`),
  flow_resolve.go:561 (`guardSeam`). No fourth. HOLDS.
- C3's re-keying arithmetic: the suite has exactly 5 `gte` cases and
  exactly 4 `contains` cases. HOLDS exactly.
- C5's `conformKind` claim: load.go:1853-1865 is a two-arm switch whose
  `int` arm is `strconv.Atoi` only — permissive as stated — and it is
  reached from `conform` (load.go:1826), which has FOUR call sites:
  normalize.go:172 (predicate atoms), normalize.go:624 (`renderWrites`,
  the `[rule.write]` path), load.go:1701 (`[initial]`), and
  `ConformValue` (load.go:1801, the CLI path). So narrowing C5 to
  normalize.go:172 leaves the other three untouched — exactly the
  scoping C5 claims. HOLDS.
- C5's bare-float claim: `valueMembers` (load.go:1748) has a `case
  float64:` at :1756 returning `strconv.FormatFloat(t, 'g', -1, 64)`,
  which renders `-0.0` to `"-0"`. HOLDS.
- C2's `set`-arm-does-not-govern-`contains` guarantee is consistent with
  `valueSatisfies` (product.go:541) and `atomAdmitsValue` (reach.go:503),
  both of which route `contains` through the same `Evaluate` and would
  otherwise flip to GuardUnevaluable on every `set` key. C2's explicit
  carve-out is what keeps lint green. HOLDS.

## Widened

- C3's suite now binds three implementers, not two: the added
  `guard.NewEvaluator` caller means `internal/guard` gains a test-time
  import of `internal/resolve`'s contract function it did not have —
  guard_evaluator_0003_test.go:67 already does exactly this, so the
  import edge exists. No new dependency; the widening is real but free.

## Findings

### IMPL-D1 — C3 says the suite has two call sites; it has three — severity: medium
- anchor: 0012:C3
- passage: "on `main` the contract test has two call sites, both
  `conformingContractSeam` in `internal/resolve/guard_mvv_test.go` — the
  real implementer is never run against the contract it is said to
  satisfy."
- fact: `rg TestGuardEvaluatorContract` returns a THIRD invocation —
  `internal/guard/guard_evaluator_0003_test.go:67`,
  `resolve.TestGuardEvaluatorContract(t, ev)`. The real implementer IS
  already run against the contract; what it is not run against is the
  KIND-AWARE contract, because there are no kind-aware cases yet.
- question: is guard_evaluator_0003_test.go:67 the site C3's
  `TestGuardEvaluatorContract(t, guard.NewEvaluator)` REPLACES, or a
  fourth site added beside it? And does its local `ev` need the
  published kernel fixture, or its own model-derived kinds?
- blocks: Phase 2's migration list. C3's "both in-tree call sites … and
  they migrate in the same phase" undercounts the migration by one file
  in a different package, and the one it misses is the only one that
  crosses a package boundary. An implementer working the stated list
  leaves `internal/guard` red at the end of Phase 2.
- note: this does NOT weaken C3's conclusion — adding
  `TestGuardEvaluatorContract(t, guard.NewEvaluator)` as the standing
  cross-implementer gate is still the right call. Only the census is
  wrong, and the census is what drives the phase's work list.

### IMPL-D2 — the widening breaks a normative 0007 signature pin C3 does not name — severity: medium
- anchor: 0012:C3
- passage: "This widening breaks both in-tree call sites
  (`internal/resolve/guard_mvv_test.go` at the `conformingContractSeam`
  and `recordingSeam` drivers) and they migrate in the same phase."
- fact: `internal/resolve/guard_mvv_test.go:222` is not a call site. It
  is `var fn func(*testing.T, resolve.GuardEvaluator) =
  resolve.TestGuardEvaluatorContract` inside
  `TestReq71_GuardEvaluatorContractIsAnImportableCrossRDRSurface`, whose
  doc comment reads "The exact name and signature are normative surface"
  and "BOUNDARY". It is a compile-time assignment: the widening makes
  that file fail to COMPILE, and the test's stated purpose is to pin the
  very signature this RDR changes.
- question: does this RDR amend `0007:REQ-71`'s named signature, or is
  the pin narrowed to the name plus non-test-file location (the
  `kernelDeclaresFunc` half at :227) with the signature half dropped?
  Either way, which record carries the deviation — and does 0007 need a
  joint-check entry the way `0003:C9`/REQ-34 got one via A2?
- blocks: whether Phase 2 can land at all without a cross-RDR amendment.
  The RDR already handles the structurally identical case for 0003 —
  A2 + the `§existing-infrastructure-audit` "Amend (field check narrows)"
  row + a recorded deviation. The same treatment for 0007's REQ-71
  signature pin is missing, and an implementer who edits :222 without one
  has silently retired another record's BOUNDARY assertion — precisely
  the failure C5's own prose warns about for 0024.

### IMPL-D3 — `guardSeam`'s arity change hits four test call sites the plan does not budget — severity: low
- anchor: 0012:C4, 0012:§phase-3-producer-alignment
- passage: C4: "`guardSeam` takes the model and is called once per
  request". Phase 3: "`internal/cli` and `internal/graphlint` test
  fixtures migrate to the constructor."
- fact: `guardSeam()` has four TEST call sites beyond the three non-test
  ones — escape_shape_0009_test.go:83 and :624,
  escape_shape_adv_0009_test.go:40 and :159 — all four in RDR 0009's
  escape-shape tests, which build their own tables. Each needs a model
  (or a kinds map) it does not construct today.
- question: do those four get the loaded model they already have in
  scope, or does `guardSeam` keep a zero-arg test-only sibling? If the
  0009 fixtures have no `*table.Model` at all, does `NewEvaluator(nil)`
  in a test contradict C1's "the zero-value construction form is
  retired" (C1 retires it for NON-TEST construction, so a nil-mapping
  test seam is admissible — but then those 0009 `eq` assertions go
  GuardUnevaluable under the new loud disposition, exactly like the 11
  internal/guard literals Phase 1 had to migrate).
- blocks: Phase 3's size estimate, and whether Phase 3 can be green on
  its own. This is the same trap Phase 1 caught and documented for
  internal/guard; the CLI's version of it is one sentence ("test fixtures
  migrate") where Phase 1 got a paragraph and an exact count.

### IMPL-D4 — `DeclaredKinds` takes `*table.Model`, but two of the three C4 sites hold only a decls map — severity: low
- anchor: 0012:C1, 0012:C4, 0012:A3
- passage: C1: "func DeclaredKinds(m *table.Model) map[string]string".
  A3: "In all three the `*table.Model` is in scope at an ancestor frame,
  never at the immediate site".
- fact: A3 is right, and the ancestors do hold `*table.Model`
  (`Denotation(m *table.Model, …)`, `matchSatisfiable(m *table.Model,
  …)`). But `internal/guard` already has a precedent for the other
  shape: `CanRefuse(decls map[string]table.TagDecl, row table.Row)` at
  product.go:648 takes the decls map, not the model.
- question: should `DeclaredKinds` take `map[string]TagDecl` (matching
  `CanRefuse`) rather than `*table.Model`, so callers holding only decls
  need not carry a model? Or is the `*table.Model` form deliberate
  because every present caller has one?
- blocks: nothing on the critical path — both forms compile at all three
  sites today. Filed because C1 fixes the signature normatively and the
  package's one existing decls-taking helper disagrees with it, which is
  the kind of split a later reader has to re-litigate. A one-clause "the
  model form is chosen because all three sites hold one; `CanRefuse`'s
  decls-map form is not the pattern here" would close it permanently.
