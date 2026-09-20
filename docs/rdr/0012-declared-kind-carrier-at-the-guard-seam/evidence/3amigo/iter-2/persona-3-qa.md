Model: claude-opus-5[1m]

# 3amigo iter-2 (delta) — Persona 3: QA / Tester — cli/0012

## Resolved by the rewrite

- **QA-1** (C5's `[emit]` ingress contradicted two live 0024 BOUNDARY
  tests) — RESOLVED, and resolved in the strongest available way. C5 is
  now narrowed to `internal/table/normalize.go::(*loader).atom`,
  `load.go::conformKind` is stated UNCHANGED, and C5 itself NAMES both
  colliding tests by function name as the reason. I re-verified both
  exist and are untouched: `internal/table/emit_grammar_0024_test.go:200`
  (`TestReq7_0024_EveryKindCheckIsLexicalAndTheAuthoredBytesSurvive`) and
  `:263` (`TestReq9_0024_NonCanonicalIntLiteralsAreAdmitted`). The
  `TestC5_EmitIntNonCanonicalRefuses` test I could not write is now
  correctly one I MUST NOT write — the `disposition` table's
  `[emit]`/`[initial]`/`[rule.write]` row says ADMITTED unchanged, and
  §Phase 4 gives the oracle for the narrowing itself ("a green
  `make check` … one that does [start failing] means the check landed on
  the shared arm"). That is a real, runnable negative control for the
  scope of C5.

- **QA-2** (nil-mapping LOUD disposition had no scenario) — RESOLVED by
  the new **0012:S5**. The scenario is writable as stated and it IS
  discriminating: held-equals-literal (`eq`, plus an `in` whose held
  value is a member) is exactly the input on which a surviving
  raw-string arm answers `GuardTrue` and the required arm answers
  `GuardUnevaluable`. The failing control is implicit but real and
  sharp — `internal/guard/grammar.go:97` is `type Evaluator struct{}`
  with a raw-string `eq` today, so the test is RED on `main` and GREEN
  only after C1's arm lands. S5 also correctly extends the assertion to
  the `var` / `new` zero values, which matters because those forms are
  in fact the majority on `main` (see QA-D1). Its **Scope note** fixing
  the residual to `lt/lte/gt/gte`/`contains` closes the last gap: I now
  know precisely which operators the test must NOT assert on.

- **QA-4** (`enum`/`scalar` indistinguishable, meta-check vacuous) —
  RESOLVED twice over. C2 now says plainly they "SHARE one arm — no
  input distinguishes them", so S2 no longer over-claims; and C2 adds
  the assertion that actually catches the regression I named: "Dispatch
  over the kind token is EXHAUSTIVE … the fallthrough `default` means
  GuardUnevaluable … so the suite asserts an unknown kind token answers
  GuardUnevaluable rather than comparing." That is a writable case with
  a real failing control (a `default: string-equal` implementation fails
  it), and §Testing Strategy's "Known coverage gaps" paragraph records
  the enum-vs-scalar half as accepted rather than claimed.

- **QA-5** (operator-inferred `gte`/`contains` re-keying unpinned) —
  RESOLVED. C2's closing paragraph now decides the two questions I could
  not answer: "their verdicts are unchanged — the kind lookup is not
  consulted for those operators at all, so `contains` over a `set` key
  does NOT become GuardUnevaluable under the `set` arm above (which
  governs `eq`/`in` only). Operator-inferred and declared-kind readings
  never compete." The before/after pin is now writable: the five `gte`
  cases (`internal/resolve/guardcontract.go:31-35`) and four `contains`
  cases (`:43-46`) keep their existing `want` values verbatim across the
  re-keying, and any drift is a test failure.

- **QA-6** (C4's "SAME loaded model" half untested) — RESOLVED as a
  DISCLOSURE, which is the correct disposition. §Testing Strategy's
  "Known coverage gaps" now states it in the place I asked for it, with
  the reason it has no oracle ("no observable — no panic, no refusal
  code — because both models are well-formed") and what substitutes
  (A3's audit; exactly three sites). I am no longer blocked on a test I
  was expected to write.

- **QA-7** (CLI canonicality leg unscenarioed, diagnostic envelope
  undecided) — RESOLVED by removal of the obligation. C5's "The CLI is
  OUT of scope: `--tag n=07` / `--write n=+1` stay accepted" plus the
  `disposition` row ("ADMITTED unchanged … this RDR adds no second
  visible break") means the CLI regression test I could not write is now
  one that must not exist. F1's "There is NO third break" and the
  accepted-residual paragraph make the silent GuardFalse→GuardTrue CLI
  flip an explicitly charted acceptance rather than an untested claim,
  and it is carried forward as new **A6** (Pending). A6 is the right
  shape for QA: its **If wrong** names the deliverable ("the
  silent-match cases must be enumerated").

## Cross-reference audit

Every S-number citation in the record, checked against what that
scenario now IS after the renumber (S1 suite / S2 dispatch / S3
reflection / S4 construction-site grep / S5 nil-mapping / S6 corpus
diff / S7 cover07 / S8 CLI `--tag` / S9 owned-reader).

| citing passage | cites | that scenario now is | correct? |
| --- | --- | --- | --- |
| §Problem Statement ("`--tag iter=many` … never reaches the seam, fixture S8") | S8 | CLI `--tag` fixture | YES |
| A4 Evidence ("blocking→clean direction (fixture S7…)") | S7 | cover07 `n eq "00"` | YES |
| A4 Evidence ("measures coverage, not safety (S6)") | S6 | corpus diff, 123 models | YES |
| A4 If-wrong ("caught by S6's before/after corpus diff and S7's regression case") | S6, S7 | corpus diff; cover07 | YES (both) |
| A6 Evidence ("S6 measured the AUTHORED side (123 models…)") | S6 | corpus diff | YES |
| C4 ("with Validation scenario 4 as its mechanical backstop") | scenario 4 | construction-site grep | YES |
| **C5 ("no measurement covers them (S5 counted guard atoms)")** | **S5** | **nil-mapping disposition** | **NO — see QA-D2** |
| C5 Evidence Record ("normative fixture S7 (the `cover07` model…)") | S7 | cover07 | YES |
| C5 Rationale ("zero `eq`/`in` atoms over an `int` tag (S6)") | S6 | corpus diff | YES |
| C5 Rationale ("S6's zero counts guard atoms only") | S6 | corpus diff | YES |
| **`authority` table, predicate-canonicality row ("S5 measured only this one")** | **S5** | **nil-mapping disposition** | **NO — see QA-D2** |
| `oracle` table, step 5 ("measured before-state (S7): exit 2 → exit 0 flip") | S7 | cover07 | YES |
| `trace` step 1 ("(S8/S9's `mvv-int`)") | S8, S9 | CLI `--tag`; owned-reader | YES (both share `mvv-int`) |
| §Existing Infrastructure Audit, CLI-conformance row ("Confirms the CLI is not the defect's ingress (S8)") | S8 | CLI `--tag` | YES |
| §Key Discoveries ("from blocking to clean (A4, S6)") | S6 | corpus diff | YES |
| F1 ("Zero committed models are affected (S6)") | S6 | corpus diff | YES |
| **F1 ("no measurement to justify (S5 counted guard atoms only)")** | **S5** | **nil-mapping disposition** | **NO — see QA-D2** |
| F6 ("Detection there falls to S4's grep") | S4 | construction-site grep | YES |
| **MVV step 1 ("the same `mvv-int` model and literal fixture S8 and S8 use")** | **S8, S8** | CLI `--tag` (twice) | **NO — see QA-D3** |
| MVV step 1 ("measured both before and after the change (fixture S8)") | S8 | CLI `--tag` | YES |
| MVV step 5 ("where today it loads and reports clean (fixture S7)") | S7 | cover07 | YES |
| **§Phase 4 ("Lands with S8's regression")** | **S8** | **CLI `--tag`, unchanged control** | **NO — see QA-D4** |
| §Testing Strategy gaps ("scenario 4 asserts only that no zero-value construction survives") | scenario 4 | construction-site grep | YES |
| S5 body ("the guarantee S4's grep explicitly defers to") | S4 | construction-site grep | YES |

Four dangling citations, all introduced by the renumber, all pointing at
an S-number whose occupant CHANGED. Filed below.

## Widened

Twice, both authorized by the brief's "you may read the code to check a
named test/fixture/symbol exists".

- **`internal/resolve/guardcontract.go` + every call site of
  `TestGuardEvaluatorContract`** — C3 makes a hard factual claim about
  the call-site set ("on `main` the contract test has two call sites,
  both `conformingContractSeam`"), and C3's new `newSeam` signature
  breaks every caller, so the caller census decides whether C3's
  migration list is complete. It is not. That produced QA-D1.
- **Every zero-value construction form of `guard.Evaluator` in
  `internal/guard`** — §Phase 1 enumerates "the eleven `Evaluator{}`
  literals" in three named files, and S4/S5 both extend the obligation
  to the `var` / `new` forms. I counted the `var` forms to check whether
  Phase 1's eleven is the whole migration. It is not; that also feeds
  QA-D1.

## Findings

### QA-D1 — C3 and S5 both assert `guard.Evaluator` is never driven against the conformance suite on `main`; it is, through a ZERO-VALUE evaluator, at a third call site neither C3's migration list nor Phase 1's eleven-literal census names — severity: high
- anchor: `0012:C3` (call-site census + "ADDED as a suite caller"), `0012:S5` (parenthetical), `0012:§phase-1-seam-constructor-and-typed-arms`
- finding: three passages of the rewrite rest on a call-site census that
  is wrong by one, and the missing site is the one that fails hardest.
  - C3: "on `main` the contract test has two call sites, both
    `conformingContractSeam` in `internal/resolve/guard_mvv_test.go` —
    the real implementer is never run against the contract it is said to
    satisfy … so this RDR lands
    `TestGuardEvaluatorContract(t, guard.NewEvaluator)` as the
    cross-implementer gate C3 exists to be."
  - S5 repeats it: "(and on `main` is never instantiated against
    `guard.Evaluator` at all — its two call sites are both
    `conformingContractSeam` …)".
  - C3's migration list: "This widening breaks both in-tree call sites
    (`internal/resolve/guard_mvv_test.go` at the `conformingContractSeam`
    and `recordingSeam` drivers) and they migrate in the same phase."

  There are THREE call sites on `main`, not two. The third is
  `internal/guard/guard_evaluator_0003_test.go:67`:
  `resolve.TestGuardEvaluatorContract(t, ev)` — and `ev` is declared at
  `:26` as `var ev guard.Evaluator`, i.e. the ZERO VALUE. The suite's
  own header comment says so in as many words: "RDR 0003's implement
  stage instantiates it against its evaluator, so it lives in non-test
  source and is called by name."

  Three consequences, each a different test problem:
  1. C3's factual premise for ADDING `guard.Evaluator` as a caller is
     false. The gate already exists; what is missing is that it runs
     against the *zero value*. The real obligation is to RE-POINT that
     call at a `NewEvaluator`-constructed seam, not to add a new one.
     As written I would land a second, duplicate call site.
  2. Under C3's new signature
     `newSeam func(kinds map[string]string) GuardEvaluator`, this site
     does not compile — it passes a value, not a constructor — so it is
     a third caller that MUST migrate in the same phase, and C3's
     "both in-tree call sites … migrate in the same phase" does not
     reach it (it is not in `internal/resolve`).
  3. Worse, it goes RED for a second, independent reason. Under C1/S5 a
     nil-mapping evaluator must answer `GuardUnevaluable` to every
     `eq`/`in`; the suite drives
     `{"eq/equal","eq","Draft","Draft",GuardTrue}` (guardcontract.go:25),
     `{"eq/unequal",…,GuardFalse}` (:26),
     `{"in/member",…,GuardTrue}` (:38) and
     `{"in/non-member",…,GuardFalse}` (:39). So S5's required disposition
     and this call site's existing assertions are in DIRECT conflict:
     the very test that proves S5 is the test S5 breaks.

  Phase 1's census misses it for the same reason. Phase 1 names "the
  eleven `Evaluator{}` literals in `guard_evaluator_0003_test.go` (7),
  `guard_narrowing_0003_test.go` (3) and `guard_fixtures_0003_test.go`
  (1)" — I confirmed those counts are exactly right for the
  composite-literal form. But Go's other zero-value forms, which S4 and
  S5 both explicitly include, add seven more in FOUR files, two of them
  unnamed by Phase 1: `guard_scope_0003_test.go:297,353`,
  `guard_grammar_0003_test.go:31,177,426`,
  `guard_evaluator_0003_test.go:26`, `guard_declaration_0003_test.go:267`.
  `guard_grammar_0003_test.go` alone carries 18 `eq`/`in` references, so
  it is not a cosmetic omission — Phase 1's stated purpose is "once this
  phase's typed arms land, those literals evaluate `eq` against a nil
  mapping and every test asserting GuardTrue/GuardFalse on an `eq` goes
  red", which is precisely what these seven do too.
- prevents: I cannot write Phase 1 as a single green commit, and I
  cannot write S5 at all without first knowing the ruling on the
  conflict. Concretely: S5's assertion
  (`var ev guard.Evaluator; ev.Evaluate({Key:"k",Operator:"eq",
  Literal:"x"}, "x") == GuardUnevaluable`) and
  `guard_evaluator_0003_test.go:67`'s existing
  `TestGuardEvaluatorContract(t, ev)` over the same `var ev` cannot both
  pass — one demands GuardUnevaluable on `eq/equal`, the other demands
  GuardTrue. The RDR gives no ruling on which wins, so there is no
  pass/fail criterion for S5 at the site that most needs one. It also
  leaves §Phase 1's "a green `make check`" unachievable as scoped: the
  migration list is short by seven zero-value sites across two unnamed
  files. Fix is mechanical and small — correct the census to three call
  sites, restate C3's obligation as RE-POINTING
  `guard_evaluator_0003_test.go:67` at a constructed seam (which is
  strictly stronger and makes C3's gate argument land properly), and
  widen Phase 1's enumeration from the eleven literals to all eighteen
  zero-value constructions — but until it is stated, the test is not
  writable.

### QA-D2 — three passages cite "S5" for the corpus measurement; S5 is now the nil-mapping scenario and the measurement is S6 — severity: medium
- anchor: `0012:C5` (CLI-out-of-scope paragraph), `0012:F1` (accepted-residual paragraph), `0012:§pre-lock-mini-checks` (`authority` table, predicate-canonicality row)
- finding: all three cite S5 as the guard-atom census, which is the
  renumbered S6. Verbatim:
  - C5: "no measurement covers them (S5 counted guard atoms)"
  - F1: "no measurement to justify (S5 counted guard atoms only)"
  - `authority` table: "deliberately NOT total over the other ingresses
    (S5 measured only this one)"

  S5 is now "Nil-mapping disposition (C1)" and counts nothing. The
  corpus scenario — "123 models, no diff … zero `eq`/`in` atoms over an
  `int` tag" — is S6. Note the record is INCONSISTENT with itself on
  this exact fact: C5's own Rationale paragraph, three paragraphs below
  the bad citation, cites it correctly twice ("zero `eq`/`in` atoms over
  an `int` tag (S6)", "S6's zero counts guard atoms only"), as do A4 and
  F1's own earlier sentence ("Zero committed models are affected (S6)").
  So the same claim is sourced to two different scenarios within one
  contract and within one failure mode.
- prevents: I cannot size the C5 scope-boundary regression from the
  record. The measurement is the ONLY justification offered for why C5
  stops at the predicate ingress, so a tester auditing that boundary
  follows the citation to S5, finds a unit test over a bare
  `Evaluator{}` that measures no population at all, and cannot tell
  whether the CLI exclusion rests on evidence or on nothing. It also
  breaks the A6 → S6 chain: A6's "S6 measured the AUTHORED side … the
  held side is a different population and is unmeasured" is the argument
  that makes A6's Pending spike scopeable, and two of the three places
  that state the authored-side measurement point somewhere else.

### QA-D3 — MVV step 1 cites "fixture S8 and S8", naming one scenario twice; the second should be S9 — severity: medium
- anchor: `0012:MVV` (step 1)
- finding: "the same `mvv-int` model and literal fixture S8 and S8 use,
  so every step below reads one guard." Two citations, same number —
  a renumber artifact (the pair was the old S7/S8, now S8/S9). The
  intended pair is unambiguous from the record: S8 is the CLI `--tag`
  fixture and S9 is the owned-reader fixture, and `trace` step 1 gets it
  right — "(S8/S9's `mvv-int`)". Both scenarios do in fact share the
  `mvv-int` model (`iter` int min 0 max 9, row `guarded` guarded
  `iter eq 7`, unguarded row `fallback`), so the claim is true; only the
  citation is broken.
- prevents: I cannot establish the MVV's fixture-sharing invariant from
  step 1 alone. The step's whole point is that one literal is read by
  every step ("so every step below reads one guard"), which is the
  property that makes MVV step 3 (`07` vs `eq 7`) and MVV step 2
  (`many`) comparable to each other and to the scenarios. As written it
  asserts S8 shares a fixture with itself — vacuously true, proving
  nothing — so the shared-fixture precondition for MVV steps 2/3/3b has
  no stated source. A tester building the fixture has to reverse-engineer
  the intended pair from `trace`.

### QA-D4 — Phase 4 says it "Lands with S8's regression"; S8 is the CLI `--tag` scenario Phase 4 explicitly leaves unchanged. The regression Phase 4 lands is S7 — severity: medium
- anchor: `0012:§phase-4-canonical-int-spelling-at-the-predicate-ingress`
- finding: Phase 4 implements C5's `strconv.Itoa(n) == authored` check at
  `normalize.go::(*loader).atom` and closes: "Lands with S8's regression
  and after Phase 1, so the literal side is closed before typed
  comparison ships — the ordering A4 depends on." S8 is now "CLI `--tag`
  conformance fires upstream of the seam", whose expected result is
  marked "unchanged before and after" — it is a CONTROL, not a
  regression, and Phase 4 in the same breath declares the CLI out of
  scope ("`load.go::conformKind` is left UNCHANGED"). Landing a phase
  with an unchanged-control scenario is a contradiction in terms.

  The regression Phase 4 actually lands is **S7** (cover07,
  `n eq "00"` → exit 2 with the canonical-rewrite diagnostic under
  `malformed_predicate_atom`). The record confirms this everywhere else:
  C5's Evidence Record says "normative fixture S7 (the `cover07` model
  and its refusal), plus MVV step 5"; MVV step 5 cites "(fixture S7)";
  the `oracle` table's step-5 row cites "measured before-state (S7)".
  This is the same off-by-one as QA-D2 — a pre-renumber S7 that is now
  S8 — and it happens to land on a scenario with the opposite polarity.
- prevents: I cannot write Phase 4's acceptance gate. Phase 4 already
  supplies a clean scope oracle ("A green `make check` after this phase
  is the expected outcome … one that [fails] means the check landed on
  the shared arm rather than the predicate site"), but a green
  `make check` is a NEGATIVE criterion only — it proves nothing was
  broken, not that the refusal fires. The positive criterion is the
  cited regression, and following the citation lands me on a fixture
  asserted to be byte-identical before and after, which can never be
  Phase 4's evidence. So the phase has a stated no-regression oracle and
  no stated pass oracle. One number.
