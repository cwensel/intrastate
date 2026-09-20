Model: claude-opus-5[1m]

# 3amigo — Persona 3: QA / Tester — cli/0012

## Widened

Yes, four times.

- **§pre-lock-mini-checks** (`oracle` + `trace` tables) — the starting-set
  scenarios S1/S4 read "green" / "no construction survives", which is the
  documented vacuity smell. The `oracle` table is where the RDR states its
  negative controls, so it had to be read before I could call a scenario
  non-discriminating. It supplied real controls for MVV steps 2, 3 and 4
  and that removed three would-be findings.
- **§failure-modes (F1..F6)** — S4's expected-result text hands the real
  guarantee off to "C1's LOUD disposition … the zero-value-survives
  failure mode", so the pass criterion is only complete if F6 states one.
- **§critical-assumptions (A1..A5)** — A3 is the only place the three C4
  construction sites are enumerated with their call frames, which is what
  decides whether S4's assertion is writable at all; A4/A5 carry the
  spike measurements S5/S6 cite as their oracle.
- **Source repo** — the persona brief authorizes checking named
  tests/fixtures/symbols. I verified every file, function and test name
  the scenarios and contracts cite (all exist as described), and then
  swept `internal/table` for existing tests that assert the OPPOSITE of
  C5. That sweep produced QA-1.

## Contract coverage

- **C1** (CARRIER — `NewEvaluator`, zero-value retired, nil-mapping LOUD):
  exercised by `0012:S4` (no zero-value construction survives) and
  `0012:S1`/`0012:MVV` step 4 (suite driven against `NewEvaluator`). The
  nil-mapping LOUD disposition itself is asserted by NO scenario — see
  QA-2.
- **C2** (TYPED COMPARISON, six arms): `0012:S2` names all six arms and
  defers to "the C2 matrix verbatim"; `0012:S8` and `0012:MVV` step 3/3b
  pin the `int` arm end-to-end. The `enum`-vs-`scalar` split is
  untestable by construction — see QA-4. `lt/lte/gt/gte` and `contains`
  unchanged-behavior clauses have no scenario — see QA-5.
- **C3** (CONFORMANCE SUITE): `0012:S1`, `0012:S2`, `0012:MVV` step 4.
  The second bound implementer (`conformingContractSeam`) has no scenario
  of its own — see QA-3.
- **C4** (PRODUCER OBLIGATION): `0012:S4` (mechanical backstop),
  `0012:S5` (corpus verdict diff). The "SAME loaded model" half of the
  clause — the residual silent mode F4 names — is UNCOVERED by any
  scenario; S4's grep cannot see it. See QA-6.
- **C5** (CANONICAL INT SPELLING): `0012:S6` (`malformed_predicate_atom`
  ingress), `0012:S7` (CLI `--tag`, but as an UNCHANGED control),
  `0012:MVV` step 5. Three of C5's five named ingresses —
  `malformed_initial_declaration`, `malformed_tag_declaration`,
  `emit_value_out_of_domain` — and the `flow-write-invalid` CLI leg have
  NO scenario. See QA-1 and QA-7.

## Findings

### QA-1 — C5's `[emit]` ingress contradicts two live BOUNDARY tests from RDR 0024, and no scenario covers the collision — severity: high
- anchor: `0012:C5` (ingress list: `emit_value_out_of_domain`)
- finding: C5 states that a value authored against an `int`-declared tag
  is admitted only when `strconv.Itoa(strconv.Atoi(s)) == s`, and names
  `emit_value_out_of_domain` (`[emit]` values) as one of the ingresses
  that must enforce it. Two tests on `main` assert the exact opposite for
  that ingress and are marked as fenced obligations of RDR 0024:
  - `internal/table/emit_grammar_0024_test.go::TestReq9_0024_NonCanonicalIntLiteralsAreAdmitted`
    (marked `// BOUNDARY`) loads a model with `count = "03"`, `"+5"`,
    `"-0"` under `kind = "int"` and asserts each LOADS, quoting REQ-9:
    *"Reusing `ConformValue` (A4) settles non-canonical literals in the
    permissive direction: `03`, `+5` and `-0` are admitted for `int`."*
    `-0` is verbatim one of C5's four named refusals.
  - `internal/table/emit_grammar_0024_test.go::TestReq7_0024_EveryKindCheckIsLexicalAndTheAuthoredBytesSurvive`
    (marked `// DOMAIN EDGE`) authors `code = "007"` under `kind = "int"`
    and asserts the row carries `"007"`, quoting REQ-7: *"Every kind
    check is a LEXICAL check on the authored string: no value is parsed
    into a typed representation, canonicalized, or converted."*

  The RDR nowhere records this as a deviation. §Contradiction Check,
  §Cross-Cutting Concerns and the `authority` mini-check all treat
  `conformKind` as an unclaimed spelling authority ("Sibling arms: none").
  C5 also asserts "REQ-7's taxonomy is untouched" while REQ-7 here is a
  0024 requirement C5 directly negates for the `[emit]` arm.
- prevents: I cannot write the `[emit]` leg of the C5 regression suite. A
  test `TestC5_EmitIntNonCanonicalRefuses` with input `count = "03"` has
  two contradictory oracles on `main` — refuse (C5) vs admit
  (REQ-9_0024) — and whichever I write, `make test` goes red on the
  other. The RDR gives no ruling on which wins, so there is no pass/fail
  criterion for C5 at the `[emit]` ingress at all. The blast radius is
  not hypothetical: it is two named BOUNDARY-marked tests that CI runs
  today. Either C5 must narrow its ingress list to exclude `[emit]`, or
  the RDR must record an explicit deviation naming both test functions as
  ones the implementation retires.

### QA-2 — C1's nil-mapping LOUD disposition has no scenario; S4 explicitly disclaims covering it — severity: high
- anchor: `0012:C1` (nil-mapping arm), `0012:S4`
- finding: C1 fixes the disposition of a bare `Evaluator{}` literal LOUD:
  "a nil-mapping evaluator answers GuardUnevaluable for every `eq`/`in`
  atom (the no-declaration arm) … never silently reverting to raw-string
  comparison." This is the single clause the whole zero-value-survives
  failure mode (`0012:F6`) rests on. `0012:S4`'s expected result is a
  grep-shaped assertion ("no zero-value construction … survives in
  non-test code") and then explicitly hands the guarantee elsewhere:
  "The guarantee proper is C1's LOUD disposition, not this grep: a site
  the pattern misses still surfaces as `guard_unevaluable` refusals."
  No scenario S1..S8 and no MVV step constructs a bare `Evaluator{}` and
  asserts its verdict. The `disposition` mini-check table has a row for
  it ("Nil-mapping evaluator (missed C4 site) → `GuardUnevaluable` for
  every `eq`/`in`") but a table row is not a scenario and the `oracle`
  table has no entry for it.
- prevents: I cannot write the one-line unit test that is the actual
  proof of the failure mode —
  `var ev guard.Evaluator; ev.Evaluate(GuardAtom{Key:"k",Operator:"eq",
  Literal:"x"}, "x")` must return `GuardUnevaluable`, NOT `GuardTrue`.
  Without a scenario fixing that oracle, an implementation that keeps the
  raw-string arm as the nil-mapping fallback passes S1 through S8 in
  full: S4 greps for construction sites and finds none, and S1/S2 drive
  only `NewEvaluator`-constructed seams. The silent-revert mode the RDR
  says is impossible would ship untested. Note this is the arm most at
  risk: `internal/guard/grammar.go:97` is `type Evaluator struct{}` with
  raw-string `eq` today, so the "do nothing to that arm" implementation
  is the path of least resistance.

### QA-3 — C3 binds `conformingContractSeam` as a second implementer but no scenario states its pass/fail criterion — severity: medium
- anchor: `0012:C3` (both-implementers clause), `0012:S1`
- finding: C3's last paragraph is a hard obligation: "Both in-tree
  implementers are bound, not only `guard.Evaluator`: the kernel
  package's own reference seam
  `internal/resolve/guard_mvv_test.go::conformingContractSeam` … is a
  raw-string `struct{}` today, so it is made kind-aware over the same
  published fixture in the phase the new cases land, or those tests
  fail." I verified this on `main`: `internal/resolve/guard_mvv_test.go:489`
  is `type conformingContractSeam struct{}` and its `opEq` arm at :492 is
  `if value == atom.Literal` — raw-string, exactly as described. But
  `conformingContractSeam` is not a `guard.Evaluator`, so it cannot be
  "constructed over the published fixture" through `NewEvaluator`; it
  needs its own fixture-carrying shape, and the RDR does not say what
  that shape is. `0012:S1` names only "driven against
  `guard.NewEvaluator`". No scenario names the second implementer, and it
  is additionally wrapped in `recordingSeam` at
  `internal/resolve/guard_mvv_test.go:256` for 0007's REQ-71 tests, so
  changing its shape has a second consumer the RDR does not scope.
- prevents: I cannot write the second-implementer conformance run —
  `TestGuardEvaluatorContract(t, <kind-aware conformingContractSeam>)`
  — because the RDR does not say how that seam receives its kinds (a new
  field? a constructor? the same `map[string]string` by value?), nor
  what `recordingSeam`'s wrapping must do with them. The failure mode is
  the one F5 names as prevented — "an implementer constructs over its own
  kinds instead of the published fixture and the kind-aware cases pass
  vacuously" — landing on the very implementer F5's detection story
  assumes is fixed.

### QA-4 — C2's `enum` and `scalar` arms are behaviorally identical, so S2's "per-kind dispatch" cannot discriminate them and C3's per-kind-token meta-check is satisfiable vacuously — severity: medium
- anchor: `0012:C2` (`enum | scalar` arm), `0012:S2`, `0012:C3`
  (one-discriminating-case-per-kind-token clause)
- finding: C2 collapses two of the five kinds into one behavior: "kind
  `enum` | `scalar`: every string parses; comparison is exact string
  equality (unchanged behavior)." `0012:S2`'s scenario lists `enum` and
  `scalar` as separate dispatch arms and its expected result is "the C2
  matrix verbatim" — but the matrix gives them the same verdict for every
  input, so no input exists that distinguishes them. C3 then requires
  "at least one discriminating case per kind token" — for `enum` and
  `scalar` that case discriminates against `int`/`bool`, not against each
  other, and an implementation that maps `scalar` onto the `enum` branch,
  or that simply defaults every unknown kind token to string equality,
  passes both scenarios and the meta-check.
- prevents: I cannot write a negative control separating `enum` from
  `scalar` dispatch. Concretely: there is no assertion that catches an
  implementation whose `switch` has no `case "scalar"` arm at all and
  falls through to `default: string-equal`. That is a real regression
  risk because C2's key-absent arm and the `set` arm BOTH want
  `GuardUnevaluable` from a default branch — an implementation that makes
  `default` mean "string equality" satisfies `enum`/`scalar` and silently
  breaks the two defensive arms. Either S2 should state that dispatch is
  exhaustive over the five tokens with no fallthrough default, or C2
  should say plainly that `enum` and `scalar` share one arm so the suite
  does not claim to discriminate what it cannot.

### QA-5 — C2's "unchanged" clauses (`lt/lte/gt/gte`, `contains`, §D13 set arms) have no scenario or regression pin — severity: medium
- anchor: `0012:C2` (final paragraph), `0012:C3`
- finding: C2 closes with "`lt/lte/gt/gte` keep their operator-inferred
  integer parse (the matrix admits only `int`); `contains` and the §D13
  set arms are unchanged; `exists` never reaches the seam." These are
  no-regression claims, and C3's re-keying plan actively MOVES the tests
  that would catch a regression: "re-keying routes … the five `gte` cases
  onto the `int` key, and the four `contains` cases onto the `set` key."
  I verified those counts against `internal/resolve/guardcontract.go`
  (5 `gte` cases at lines 30-34, 4 `contains` cases at lines 41-44, all
  currently `Key: "subject"`) — so the plan is accurate, but it means
  every existing `gte`/`contains` case gets a new key AND a new kind at
  the same moment the seam becomes kind-aware. No scenario asserts the
  verdicts are unchanged across that move.
- prevents: I cannot write the before/after pin for the operator-inferred
  arms. The dangerous interaction is concrete and unstated: once
  `gte` cases are re-keyed onto an `int`-declared key, C2's `int` arm and
  the operator-inferred parse both apply to the same atom — and C2 does
  not say which wins, or whether they must agree. Worse for `contains`:
  its four cases move onto a `set` key, but C2's `set` arm says the seam
  answers `GuardUnevaluable` for `eq`/`in` over `set` — it is silent on
  whether the kind lookup is consulted at all for `contains`, so
  `{"contains/superset", … GuardTrue}` could legitimately become
  `GuardUnevaluable` under a kind-first dispatch and I would have no
  contract text saying that is wrong.

### QA-6 — C4's "SAME loaded model" half is UNCOVERED; S4's oracle only tests the negation of zero-value construction — severity: medium
- anchor: `0012:C4`, `0012:S4`, `0012:F4`
- finding: C4 carries two obligations: (a) every site constructs through
  `NewEvaluator` (no zero-value), and (b) each constructs "over the
  declarations of the SAME loaded model whose rows it evaluates: one
  model, one evaluator, no mixed-model evaluation." `0012:S4`'s expected
  result tests only (a) — a source-pattern assertion over
  `Evaluator{}` / `var` / `new` / embedded-field zero values. Obligation
  (b) is acknowledged as untested by `0012:F4`: "Residual silent mode: a
  construction site violating C4 (wrong model's kinds) can still type a
  comparison wrongly without refusing — guarded by the site audit (A3)
  and the conformance suite's fixture-bound cases, not by the kernel."
  A site audit is a human read, not a test, and the conformance suite
  constructs its own fixture so it cannot observe a caller's model
  mismatch.
- prevents: I cannot write a mixed-model test. There is no stated
  observable for "evaluator built over model A, rows from model B" — no
  panic, no refusal code, no assertion target. Given A3's finding that
  `valueSatisfies` and `atomAdmitsValue` have the `*table.Model` only at
  an ancestor frame (`Denotation(m *table.Model, …)` one frame up;
  `matchSatisfiable(m *table.Model, …)` two frames up) and that
  `flow_next.go::probeRow` takes only a `table.Row` with the model one
  frame further, the plumbing that must be threaded is exactly where a
  mismatch would be introduced — and it is the one part of C4 with no
  pass/fail criterion. Accepting this as untested is defensible, but it
  should be stated in the §Testing Strategy as a known coverage gap, not
  only inside a failure mode.

### QA-7 — C5's `flow-write-invalid` CLI leg and three of its five refusal ingresses have no scenario — severity: low
- anchor: `0012:C5` (refusal-code list), `0012:S6`, `0012:S7`
- finding: C5 names five refusal codes it must fire under:
  `malformed_predicate_atom`, `malformed_initial_declaration`,
  `malformed_tag_declaration`, `emit_value_out_of_domain`, and
  `flow-tag-invalid` / `flow-write-invalid`. Only
  `malformed_predicate_atom` gets a scenario (`0012:S6`, via `n eq "00"`
  in a guard). `0012:S7` covers the CLI but is expressly an UNCHANGED
  control (`--tag iter=many` is a kind failure, not a canonicality
  failure) — the NEW CLI break C5 and `0012:F1` both describe,
  "`--tag n=07` / `--write n=+1`, accepted today, now refuse", has no
  scenario at all. I confirmed both codes exist
  (`internal/cli/flow_input.go:36,40`) and that `canonicalValue` reaches
  `table.ConformValue` at `flow_input.go:751` (single-valued) and `:769`
  (set members), so the leg is real and reachable.
- prevents: I cannot write the CLI canonicality regression —
  `intrastate flow resolve --tag n=07 …` expecting `flow-tag-invalid`
  with the canonical-rewrite diagnostic, and `--write n=+1` expecting
  `flow-write-invalid`. §Failure Modes calls this a user-visible break
  ("a third" break in F1), which is precisely the class that needs a
  pinned normative fixture, and its diagnostic wording is unspecified for
  the CLI surface: C5 gives one example string shaped for a load site
  (`<site>: "00" is not the canonical spelling of int tag n; write 0`),
  while the CLI's existing wrapper prefixes "the value for `n` does not
  conform to its declaration: " (`flow_input.go:753`). Which of the two
  envelopes wins is undecided, so no exact-match assertion is writable.
  Low rather than medium only because `0012:S7`'s adjacent fixture
  establishes the exit code and envelope shape, leaving just the message
  text unpinned.
