Model: claude-opus-5[1m]

# A12 — consumer-side verification for splitting `graph-dangling-edge` by arm

Scope: the mechanism half (both arms in `checkDanglingEdge`, distinguished by
`element`) is already proven at `evidence/spikes/cove-dt-terminal/result.md`.
This document verifies the CONSUMER half only: whether any golden fixture,
testdata expectation, exit-code path, or 0006 conformance test keys on the
finding code `graph-dangling-edge` AS A WHOLE such that class-keying the
missing-root arm alone would change an existing expectation.

Work was read-only on source. No source file was modified.

## 1. Baseline

```
$ go test ./internal/graphlint/...
ok  	github.com/newcoinc/intrastate/internal/graphlint	3.499s
```

PASS (exit 0). Green baseline established.

## 2. Enumeration of consumers

```
$ rg -n 'graph-dangling-edge|CodeDanglingEdge' internal/ testdata/
rg: testdata/: No such file or directory (os error 2)
internal/graphlint/taxonomy.go:25:	CodeDanglingEdge       = "graph-dangling-edge"
internal/graphlint/taxonomy.go:67:	CodeDanglingEdge,
internal/graphlint/invariants_0006_test.go:25:		graphlint.CodeDanglingEdge,
internal/graphlint/invariants_0006_test.go:70:	f := requireCode(t, r, graphlint.CodeDanglingEdge)
internal/graphlint/invariants_0006_test.go:82:	requireNoCode(t, clean, graphlint.CodeDanglingEdge)
internal/graphlint/invariants_0006_test.go:206:	f := requireCode(t, r, graphlint.CodeDanglingEdge)
internal/graphlint/invariants_0006_test.go:217:			graphlint.CodeDanglingEdge, render(r))
internal/graphlint/invariants_0006_test.go:893:	requireCode(t, r, graphlint.CodeDanglingEdge)
internal/graphlint/findings_0006_test.go:102:	f := requireCode(t, r, graphlint.CodeDanglingEdge)
internal/graphlint/findings_0006_test.go:138:			graphlint.CodeDanglingEdge, render(r))
internal/graphlint/findings_0006_test.go:251:// REQ-73: "The blocking code set is exactly: `graph-dangling-edge`,
internal/graphlint/findings_0006_test.go:261:		"graph-dangling-edge",
internal/graphlint/analysis.go:118:			Code:    CodeDanglingEdge,
internal/graphlint/analysis.go:139:				Code:     CodeDanglingEdge,
internal/cli/lint_fixtures_0006_test.go:112:// `graph-dangling-edge` naming the missing declaration.
internal/cli/lint_mvv_0006_test.go:105:// The named matrix is `graph-dangling-edge`, `graph-dead-end`,
internal/cli/lint_mvv_0006_test.go:128:			want: graphlint.CodeDanglingEdge},
internal/cli/lint_gate_0006_test.go:155:	// `graph-dangling-edge` and never to a clean empty reachable set.
```

Note: there is no top-level `testdata/` directory. The only `testdata/` tree in
the repo is `internal/table/testdata`, which holds loader fixtures; none of them
is consumed by a graph-lint expectation.

```
$ find . -path ./.git -prune -o -name '*.toml' -print | grep -v 'docs/rdr'
./.roborev.toml  ./.kata.toml  ./models/rdr.toml
./internal/table/testdata/*.toml  (loader fixtures, 27 files)
```

```
$ rg -n 'pos-no-initial' internal/
internal/table/testdata/... consumed only by:
internal/table/accessors_test.go:498:		m := mustLoad(t, "pos-no-initial.toml")
internal/table/roundtrip_test.go:1023:		m := mustLoad(t, "pos-no-initial.toml")
```

The one rootless checked-in fixture, `pos-no-initial.toml`, is consumed by
`internal/table` LOADER tests only (`mustLoad`), never by graph lint. It carries
no lint expectation to break.

## 3. Consumer classification

Legend: (a) asserts on code only · (b) asserts on code+element · (c) counts
findings · (d) other.

| # | Site | Class | Fixture / subject | Which arm produces it | Affected by class-keying the root arm? |
|---|---|---|---|---|---|
| 1 | `internal/graphlint/taxonomy.go:25` | (d) constant definition | — | — | No — the code string is unchanged; A12 changes no taxonomy member |
| 2 | `internal/graphlint/taxonomy.go:67` | (d) blocking-set membership | — | — | No — membership is unchanged |
| 3 | `internal/graphlint/invariants_0006_test.go:25` (`TestReq31`) | (d) asserts the code is a member of `BlockingCodes()` | no model at all | — | No — pure set-membership assertion, runs no lint |
| 4 | `internal/graphlint/invariants_0006_test.go:70` (`TestReq32`) | (a) code only | `noRoot` over `statusOnlyDecls` — **no `class` key** | root arm | No — fixture is state-machine class (no `class` declared) |
| 5 | `internal/graphlint/invariants_0006_test.go:82` (`TestReq32` control) | (a) code only, negative (`requireNoCode`) | `legalBody` over `twoStateDecls`, **has `[initial]`** | neither arm fires | No — a rooted state machine; the root arm was already silent here |
| 6 | `internal/graphlint/invariants_0006_test.go:206,217` (`TestReq129`) | (b) code+element/key — asserts a finding names `seen` / `done` | terminal predicate over the observed tag `seen`; fixture **has `[initial] status = "a"`** | **terminal arm only** | No — this is the arm A12 keeps unconditionally live, and the fixture has a root so the root arm never fired here even today |
| 7 | `internal/graphlint/invariants_0006_test.go:887,893` (`TestReq44`) | (c)+(a) — asserts `len(r.Blocking()) != 0`, then code | `noRoot` over `statusOnlyDecls` — no `class` | root arm | No — state-machine class; count assertion is a lower bound (`== 0` is the failure), not an exact total |
| 8 | `internal/graphlint/findings_0006_test.go:102,138` (`TestReq127`) | (b) code + identity-field discrimination | `noRoot` — no `class` | root arm | No — state-machine class. Note this test **specifically depends on the root arm firing**: line 136-140 fails if *every* dangling finding carries a rule id, because the missing-root finding is the one that exercises the span/element fallback. It is the strongest code-as-a-whole coupling found, and it is scoped to a class-omitted fixture |
| 9 | `internal/graphlint/findings_0006_test.go:261` (`TestReq73`) | (d) exact blocking-code-set equality | no model | — | No — asserts the ten-member code set; A12 adds/removes no code |
| 10 | `internal/cli/lint_mvv_0006_test.go:128` (`mvvIllegalMatrix`, case `dangling-edge/missing-root`) | (a) code only + exit-code path | `mvvNoRoot` (`internal/cli/lint_fixtures_0006_test.go:112`) — **no `class` key** | root arm | No — state-machine class |
| 11 | `internal/cli/lint_gate_0006_test.go:155` (`TestReq120`) | (d) exit/aggregate-code path via `dropInitialTable(models/rdr.toml)` | `models/rdr.toml` with `[initial]` stripped — **no `class` key** | root arm | No — the checked-in model declares no class; it is state-machine by the zero value |
| 12 | `internal/cli/lint_gate_0006_test.go` `TestReq122` (false-positive census) | (c) counts blocking findings over `models/rdr.toml` **as accepted** | `models/rdr.toml` intact, has `[initial]` at line 91 | neither arm | No — the model is rooted and clean; census expects zero |

Exit-code paths (items 10, 11) assert only `clierr.ErrorCode(err) ==
graphlint.AggregateCode` (`graph-lint-failed`), `GroupUserEnv`, and exit 2. They
never partition by the dangling code's arm.

Helper semantics confirm no exact-total coupling exists on this code:
`requireCode` (`internal/graphlint/fixtures_0006_test.go:225`) asserts
`len(withCode(...)) != 0`; `requireNoCode` (`:247`) asserts `countCode(...) == 0`.
`requireOneCode` (`:236`) does assert exactly one, but **no dangling-edge
consumer calls it** — the only `len(f) != 1` exact assertions in the suite are
`internal/graphlint/adversarial_0006_test.go:339` and
`internal/graphlint/coverage_0006_test.go:630`, both over `CodeOverlap`.

## 4. The decisive scoping fact: no model in the repo can be decision-table class

```
$ rg -n '"class"' internal/
(no match for a `class` layout key anywhere in internal/)
```

Every `Class` hit under `internal/` is either `clierr.Finding.Class` (the escape
failure-class field, RDR 0003/0009) or `accessor.Class*` refusal constants —
neither is a model class. Confirmed at the struct:

`internal/table/model.go:366-393` — `type Model struct` carries
`ID, Version, Description, Metadata, Outcomes, Initial, Terminal, Tags, Readers,
Writers, Gates, DumpOrder, Rows`. **There is no `Class` field.**

Consequences:

1. The `class` key does not exist in the loader, so no checked-in model, no test
   fixture, and no `internal/table/testdata` file can declare
   `class = "decision-table"`. 0010 introduces the field.
2. Therefore the class-keyed predicate C5 specifies for the root arm —
   `class == "decision-table" || len(Initial) > 0`, **augmenting** rather than
   replacing `len(Initial) == 0` — evaluates identically to today's
   `len(a.model.Initial) == 0` for every model that exists in the repo right
   now, because `Model.Class()`'s empty zero value reads as `state-machine`
   (0010:C1).
3. Every fixture in the table above is state-machine class by construction. The
   root arm keeps firing on all of them, unchanged.

This answers question 3 affirmatively: the class-keying IS scoped so
state-machine models keep both arms unchanged, and that scoping is enforced by
the augmenting form of the predicate, not merely by convention. Substituting a
replacing predicate (`class == "decision-table"` alone) would break items 4, 7,
8, 10, and 11 — which is exactly why C5 spells the predicate as an OR and
0006:C18 is listed as NOT overridden.

## 5. Answers to the posed questions

**Q1 — does any consumer key on the code as a whole such that class-keying only
the missing-root arm changes an existing expectation?**

No. Two consumers key on the code as a whole in a way that depends on the root
arm firing (items 8 and 7), and one exit-code path does (item 11), but all three
run over class-omitted fixtures, which remain state-machine and keep the root
arm live. No consumer keys on the code as a whole over a decision-table model,
because no such model can be authored today.

**Q2 — enumeration and classification.**

Done in §3. Twelve consumer sites: four (d) taxonomy/set-membership or
exit-path, five (a) code-only, two (b) code+element, two (c) counting — and the
counting ones are lower-bound (`== 0` is failure) or a zero-census over a rooted
clean model, never an exact total on the dangling code.

**Q3 — would keeping the terminal arm live while silencing the root arm for the
decision-table class break any of them?**

No. The only consumer that asserts on the terminal arm is item 6
(`TestReq129`), whose fixture declares `[initial] status = "a"` — it exercises
the terminal arm in isolation today, and A12 leaves the terminal arm
unconditional, so it is untouched. No consumer asserts that the root and
terminal arms co-occur, and no consumer asserts an exact total count of
`graph-dangling-edge` findings.

**Q4 — baseline.** PASS, §1.

## 6. Residual note (not blocking)

`internal/graphlint/findings_0006_test.go:136-140` (REQ-127) is the one place
where the ROOT arm is load-bearing for a property the code-as-a-whole carries:
it requires at least one `graph-dangling-edge` finding without a rule id, to
prove the span/element fallback is exercised. If a future decision-table fixture
were ever added to *that* test's table, silencing the root arm would make its
guard fire. It is not affected today (single class-omitted fixture, inline
const), but an implementer adding decision-table cases to REQ-127 should read
that assertion first. This is a note for the implementer, not a defect in A12.

## Verdict

**PASS.** A12's consumer half is verified. No golden fixture, testdata
expectation, exit-code path, or 0006 conformance test keys on
`graph-dangling-edge` as a whole over a model that could be decision-table
class. All twelve consumers run over state-machine-class models — necessarily,
since `table.Model` has no `Class` field and no `class` layout key exists — and
C5's augmenting predicate (`class == "decision-table" || len(Initial) > 0`)
leaves every one of their expectations bit-identical. The terminal arm stays
unconditional and its sole consumer (`TestReq129`) already exercises it over a
rooted model. Splitting the arms needs no taxonomy change and no new finding
code.
