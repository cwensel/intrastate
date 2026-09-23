# Verification — RDR 0012 Declared-kind carrier at the guard value seam

## Phase 3a — CoVe

Independent pass: inputs taken from the record and `req-list.md` only. The
Phase 1 tests, `coverage.md` and `deviations.md` were not read. Every probe
ran against branch `worktree-rdr-0012` at `8bd2e85`. The binary was built
from that tree, and the baseline binary was built from `901ae3d`. Scratch
models and mutation copies live under /tmp.

FAIL-48 — the S4 zero-value backstop (REQ-48) does not catch several
zero-value constructions of `guard.Evaluator` in non-test code. REQ-48
requires "a typed check, not a text grep ... failing on any composite
literal, `var` declaration, `new`, or struct field of type
`guard.Evaluator` outside `NewEvaluator`'s own body". Failing inputs: each
one was a single new non-test file added to a scratch copy of the branch,
followed by `go test ./internal/guard/ -run 'TestReq48_0012|TestReq35_0012'`.
Observed behaviour: the result was `ok` (the check passed) for all of these:
  - composite literal through a type alias: `type zzT = guard.Evaluator; var zzG = zzT{}`
    (in `internal/cli`, and `type zzT = Evaluator; var zzA = zzT{}` inside `internal/guard`)
  - composite literal through a dot import: `import . ".../internal/guard"; var zzA = Evaluator{}`
  - composite literal with the element type elided: `var zzA = []guard.Evaluator{{}}`
  - zero-value declarations that are not a plain `var x guard.Evaluator`:
    a named result `func zzF() (e guard.Evaluator) { return }`, an array
    `var zzA [1]guard.Evaluator`, and `make([]guard.Evaluator, 1)[0]`
Each of these hands out a nil-mapping evaluator. The check catches the
plain forms: `guard.Evaluator{}`, `&guard.Evaluator{}`, a renamed import
(`g.Evaluator{}`), a package-scope or local `var`, `new(guard.Evaluator)`,
a named or embedded struct field, a site in `cmd/`, a site inside
`internal/guard`, and renaming `atomAdmitsValue`. So the check matches the
syntax of the type selector, not the type (go/types) that REQ-48 names.

Outcome: fixed:6c113a5 — the S4 check type-checks the module's non-test packages with go/types (dependencies from `go list -export` data) and matches guard.Evaluator by type through aliases, dot imports and elided literal types; zero-value vars include named results and arrays of it, plus `new`/`make`. Each failing input above is now a plant in `TestAdv1_0012_…`, and every plant fails the check.

Probed and holding (no violation), for the record:
- REQ-1/2/3/4/10: `NewEvaluator(map[string]string) guard.Evaluator`
  compiles as a concrete-typed value. `DeclaredKinds(nil)` and
  `DeclaredKinds(&table.Model{})` each return an empty non-nil map. A
  `Kind:""` tag is omitted, and `eq` over that key is Unevaluable. The
  loader separately refuses `kind = ""` and a missing kind with
  `malformed_tag_declaration`.
- REQ-5/47: the only field is `kinds map[string]string`. Adding a
  `resolve.TagSet`, `[]resolve.Tag`, `map[string][]string` or `string`
  field fails `TestReq34_…` and `TestReq47_0012_…`.
- REQ-7/22/49: a nil mapping built as `Evaluator{}`, `var`, `new`, an
  embedded field, or `NewEvaluator(nil)` gives Unevaluable for `eq`/`in`
  when the held value equals the literal or is a member. With that same
  mapping, `gte` and `contains` still decide True.
- REQ-11..18/20/21/23/24: a 48-case seam matrix. Results: int parsed
  (`07`,`+7` == 7), `many`/``/` 7`/`7.0`/overflow U, `in ["3","x"]` held 3
  U, bool `1`/`0`/`True`/`yes` U, bool `in ["true","yes"]` U,
  enum/scalar byte-equal, set `eq`/`in` U, set `contains` decided, absent
  key U, unknown token `decimal` U, `exists`/unknown/empty operator U, no
  panics.
- REQ-25..29/31/MVV-4: `TestGuardEvaluatorContract` is green over
  `guard.NewEvaluator`. `ContractKinds()` returns a fresh map and carries
  exactly the five tokens. Ten mutant seams each fail the suite: raw int
  `eq`, raw bool `eq`, permissive bool `in`, Unevaluable enum,
  Unevaluable scalar, raw set, raw absent key, raw unknown token, a
  per-member int `in`, and a constructor that ignores the kinds it is
  given. An int-arm mutant in `grammar.go` also fails the re-pointed
  `TestReq34_EvaluatorDecides…` call.
- REQ-MVV steps 2/3 and REQ-53: this is the owned door driven by the
  artifact, `{"iter":…}` with `flow resolve --artifact`. `many` refuses
  `flow-guard-unevaluable` with a finding "rule `guarded`: the atom on
  `iter` could not be decided (uncomparable)". The top-level keys are
  unchanged: code, message, schema_version and findings. `07` selects
  `guarded`, where the baseline binary selected `fallback`. `4` selects
  `fallback`. Owned bool `1`/`True`/`yes` refuses; `true` matches.
- REQ-38: the envelope is identical in shape to the baseline's existing
  `gte`-over-`many` refusal.
- REQ-54: `flow next` with owned `many` reports
  `unknown:[{key:iter,reason:uncomparable}]`. `07` is a candidate with no
  unknowns, and `4` is excluded.
- REQ-11/35 (match path): `[rule.match.iter] eq = 7` with
  `--tag iter=07` does not match (it is byte-compared), and `in = [3,7]`
  with `07` does not match. `atomAdmitsValue` keeps its name, its `bool`
  return, and its sole `guard.Evaluator{}` site.
- REQ-39/40/41/43/44/45/51 (C5): `cover07` gives exit 2 with
  `malformed_predicate_atom: rule r-zero guard.all: n.eq: "00" is not the
  canonical spelling of int tag n; write 0`. These are refused with the
  rewrite named: `-0.0` (rendered `-0`), `"-0"`, `"+1"`, `"01"`, and `in`
  members (the first offender in authored order, `"00"` in
  `[0,"00","01"]`). The same holds for the `gte` bound, a
  `[rule.match.n]` literal, a `guard.unless` literal, and a
  `[context.*.match]` literal. These are admitted: `"0"`, `0`, `0.0`,
  `[rule.write] "08"`, `[initial] "03"`, and an int `[emit] "07"`. A
  `" 0"` or overflow literal is refused by the existing "is not an int",
  with no normalization.
- REQ-42/52: `--tag iter=07` and `--tag iter=+1` are accepted. `flow
  set-state --write iter=+1` is accepted. `--tag iter=many` gives exactly
  `{"code":"flow-tag-invalid","message":"the value for `iter` does not
  conform to its declaration: \"many\" is not an int",
  "schema_version":"0.1","param":"iter"}` with exit 2, the same before and
  after.
- REQ-33/34: `guardSeam(m)` is called once per request in `flow resolve`,
  where it is threaded to `escapeClassOf`, and once per request in `flow
  next`, where it is threaded to `probeRow`. There are no other kernel
  `Resolve` call sites.
- REQ-9: the non-test imports of `internal/resolve` are only errors, fmt,
  maps, slices, strings and testing. There is no declaration type.
- REQ-50: `lint --model --as json` was run over all 218 committed `.toml`
  files with the baseline and branch binaries. There were 0 differences in
  output or exit code. `go test ./... -count=1` on the branch is green.
- REQ-37: `docs/model-authoring.md` states the guard-typed / match-byte
  split and the C5 rule.

## Phase 3b — Adversarial

Three failure modes, anchored in the record's Failure Modes section. Each
was probed with a test that plants the failure and asserts the detection
the record relies on. Only a probe that fails now is recorded as an ADV
entry. Probes that passed were removed, because a passing probe catches
nothing.

ADV-1 — "Zero-value evaluator survives migration" (Failure Modes; S4 / REQ-48). BLOCKING.
The record says that a missed construction site whose guards use only
`lt/lte/gt/gte` or `contains` "migrates with no signal at all" and that
detection "falls to S4's typed construction check". The S4 check
(`TestReq48_0012_NoZeroValueEvaluatorSurvivesOutsideTheOneAllowListedFunction`)
catches the spellings it lists: `guard.Evaluator{}`, `var ev guard.Evaluator`,
`*new(guard.Evaluator)`, an embedded struct field, and an in-package
`Evaluator{}`. It PASSES, and so misses, four other zero-value forms Go
admits in non-test code:
  - a composite literal through a type alias (`type E = guard.Evaluator; E{}`);
  - an elided composite literal (`[]guard.Evaluator{{}}[0]`);
  - a named result (`func f() (ev guard.Evaluator) { return }`);
  - a make-allocated element (`make([]guard.Evaluator, 1)[0]`).
The first two are composite literals of type `guard.Evaluator`, which S4
names outright ("failing on any composite literal … of type
`guard.Evaluator`"). The check therefore matches the syntax, not the type,
and S4 asks for the type. The last two are the "zero values Go equally
admits" that S4 is there to close. Any of the four reopens the silent
migration path for ordering and `contains` guards.
Test: `internal/guard/guard_adv_0012_test.go::TestAdv1_0012_TheZeroValueCheckCatchesEveryZeroValueForm`.
It copies the module to a temp dir, confirms the S4 check passes on the
pristine copy, then plants each form alone as a compilable non-test file
and runs the S4 check there. The 5 control plants FAIL the check as they
should. The 4 forms above PASS it, so the test FAILS now (4 of 9
subtests).
Remedy direction: make the S4 pass resolve types with go/types over every
expression and declaration (`types.Identical(T, guard.Evaluator)` on
composite literals, including elided and alias-typed ones, on value
specs, named results, `new`/`make` calls, and struct fields), not
selector text.

Outcome: fixed:6c113a5 — `TestAdv1_0012_TheZeroValueCheckCatchesEveryZeroValueForm` is green: all 12 plants (5 controls, the 4 forms above, plus FAIL-48's dot import, in-package alias and array var) fail the S4 check and name their file.

Probed, no finding:
  - "Suite/fixture drift" (Failure Modes; C3). `TestGuardEvaluatorContract`
    was run in a subprocess over eight drifted constructors: all-scalar,
    implementer-own keys, int→scalar, bool→scalar, set→scalar, int→bool,
    scalar→int and enum→bool. The suite FAILS on every one, so the
    discriminating legs do what the record says. (enum vs scalar is
    EXCLUDED by construction.)
  - "Silent failure guarded against" / "Visible break" (Failure Modes). A
    loaded model with an int and a bool owned tag, guarded
    `unless iter eq 7`, `all iter in [7, 8]` and `unless flag eq true`,
    was resolved through `NewEvaluator(DeclaredKinds(m))` with the owned
    values `many`, `yes` and `1`. Every case refuses `guard_unevaluable`,
    and none misroutes to a plan through the `unless` block or the `in`
    list.
