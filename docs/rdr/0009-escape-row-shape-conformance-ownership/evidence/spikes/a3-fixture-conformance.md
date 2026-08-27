Model: claude-opus-5

# SPIKE A3 — escape-fixture conformance changes no frozen assertion outcome

RDR 0009, Stage 4 (Resolve). Assumption under test:

> Bringing escape fixtures into conformance (dropping `Writes` from
> `escapeRow`-built rows) changes no frozen assertion outcome: every
> ADV/MVV/boundary test still discriminates, so the entry precondition
> lands without weakening RDR 0001's frozen evidence.

Package under test: `internal/resolve` (Go). Date run: 2026-08-11.
Repo state at spike time: branch `via-claude`, HEAD `9678601`.

## Commands run

| Command | Output |
| --- | --- |
| `go test ./internal/resolve/... -v -count=1` (unmodified tree) | `a3-baseline.out` |
| `go test ./internal/resolve/... -v -count=1` (Writes stripped from escape rows) | `a3-conformed.out` |
| `go test ./internal/resolve/... -v -count=1` (conformed tree + mutated `escapeOrRefuse`) | `a3-conformed-mutant-gate-bypass.out` |
| `go test ./internal/resolve/... -v -count=1` (Writes **and** NextTags stripped) | `a3-conformed-plus-nexttags.out` |
| `git checkout -- internal/resolve/` then `git status --short` | tree clean of spike edits |

Outcome-set comparison used `grep -E '^\s*--- (PASS|FAIL|SKIP)'`, timing
normalized, sorted, then `diff`. Raw `-v` streams differ only in the
interleaving of map-ranged subtests (Go randomizes map iteration), which is
why the sorted-outcome comparison is the load-bearing one.

## What was edited (variant 1: conformance)

Three edits, all in `_test.go` files. `resolve.go` was not touched for the
conformance variants.

1. `internal/resolve/fixtures_test.go`, `escapeRow` — removed the line
   `Writes: []resolve.Tag{{Key: "status", Value: "Blocked"}},`.
2. `internal/resolve/adversarial_test.go:229` (ADV-2b) — removed the
   call-site override `escape.Writes = []resolve.Tag{{Key: "status", Value: "Escaped"}}`.
3. `internal/resolve/fixup_test.go:103` (Fixup-1e "missing owned state")
   — removed the identical call-site override.

Edits 2 and 3 matter: `escapeRow` is not the only place an escape row gets
`Writes`. Stripping only the builder would have left two escape rows still
carrying writes, so the conformance test would have been incomplete.

## Baseline vs conformed

| | Baseline | Conformed |
| --- | --- | --- |
| Exit code | 0 | 0 |
| `--- PASS` lines | 154 | 154 |
| `--- FAIL` lines | 0 | 0 |
| Sorted outcome set | — | **identical to baseline** |

No test changed outcome. No test was added, removed, or skipped.

## Discrimination analysis

"Still passes" is not the question; "still discriminates" is. Every
assertion in the escape suite was inspected for what it reads, and the
suite was then re-run against a deliberately broken kernel.

### a. Assertions that could have gone vacuous

The only `Plan.Writes` reads that touch an escape-built row are inside
`t.Fatalf` **failure-message format arguments**, never inside a condition:

- `adversarial_test.go:241` (ADV-2b) — `got.Plan.Writes` is an argument to
  the `t.Fatalf` that fires when `!got.Refused()`. The condition is
  `!got.Refused()`; the asserted property is `Refusal.Kind ==
  KindOwnedStateUnavailable` (line 243).
- `fixup_test.go:57` (Fixup-1d) — same shape; asserted properties are
  `Refusal.Kind == KindGuardUnevaluable` and `Refusal.Guard == escape.Guard`.
- `fixup_test.go:115` (Fixup-1e "missing owned state") — same shape;
  asserted property is `Refusal.Kind == KindOwnedStateUnavailable`.

Because these are diagnostics on the failure path and not predicates, no
assertion outcome depends on the fixture's `Writes` content. Stripping
`Writes` degrades the failure *message* (it would now print `[]` instead of
`[{status Escaped}]` if the kernel regressed) but changes no verdict.

The substantive `Writes`-equality assertions live entirely on
**non-escape** fixtures and are untouched by this edit:

- `resolve_test.go:671-673` (REQ-18) asserts `p.Writes` DeepEquals
  `[{status Final}]` — built from `singleMatchTable`, not `escapeRow`.
- `mvv_test.go:67-70` (MVV leg 1) asserts the same against
  `legalInput()` — `singleMatchTable`, not `escapeRow`. The MVV comment
  "The replayed plan must be substantive, or 'identical' is vacuous" names
  exactly the vacuity risk, and the fixture it guards is unaffected.
- `resolve_test.go:709`, `:731`, `:986` (REQ-19, REQ-29) — all
  `legalInput()`.
- `mvv_test.go:110`, `resolve_test.go:689`, `:1005`, `:1058` — these assert
  a *refusal* carries no plan/writes. Their guard is `got.Plan != nil` /
  `got.Refused()`; they were already reading a nil plan, so escape-row
  `Writes` never fed them either.

**No test's assertions become vacuous.** No test that previously read a
non-empty `Writes` in a *condition* now reads an empty one.

### b. Per-test verdict for escape-touching tests

| Test | Reads escape `Writes`? | Property still tested after strip |
| --- | --- | --- |
| `TestAdv1b_AllGuardsFalseIsNoMatchNotOwnedStateUnavailable` | no | escape rescues a guard-FALSE zero-match; `RuleID`+`Escaped` asserted — real |
| `TestAdv2_.../guard FALSE must not rescue` | no | `Refusal.Kind == no_match` — real |
| `TestAdv2_.../guard UNEVALUABLE must not rescue` | no | `Refusal.Kind == guard_unevaluable` — real |
| `TestAdv2b_EscapeEdgeMustNotBypassTheOwnedStateRequirement` | message only | `Refusal.Kind == owned_state_unavailable` — real |
| `TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue` | message only | `Kind == guard_unevaluable` **and** `Refusal.Guard` — real |
| `TestFixup1e_.../guard FALSE must not rescue an ambiguity` | no | `Kind == ambiguous_match` — real |
| `TestFixup1e_.../missing owned state must not rescue an ambiguity` | message only | `Kind == owned_state_unavailable` — real |
| `TestFixup1e_.../a viable ambiguous-class escape still rescues` | no | control: `RuleID`+`Escaped` — real |
| `TestFixup3c_DegradedEscapeAmbiguityPayloadMustNotDependOnRowOrder` | no | `Kind == ambiguous_match` + `DeepEqual` on the two refusal payloads — real |
| `TestFixupGateIsUniformAcrossOrdinaryAndEscapeCandidates` (3 subtests) | no | ordinary and escape refusal kinds must agree — real |
| `TestReq15_EscapeEdgeRescuesOnlyWhenItMatchesExactlyOnce` (3 subtests) | no | rescue exact-one, degrade-on-two, wrong-class-no-rescue — real |
| `TestReq16_.../multiple with a modeled ambiguous escape emits a plan` | no | `RuleID == rdr.escape.ambiguous` — real |
| `TestReq15_AllFourListedConditionsRefuseWithoutAnEscapeEdge` | no (asserts `len(row.Escape) == 0` on its fixtures) | unaffected — real |
| MVV suite (all 4 legs) | no escape fixtures used | unaffected — real |
| `boundary_test.go` | file holds AST helpers only; no escape fixtures | unaffected |

Note `boundary_test.go` in this package contains only package-boundary
inspection helpers (import/AST checks) and constructs no escape rows, so it
has no escape-fixture exposure at all.

### c. Mutation check — the tests still catch a real regression

Passing after a fixture change can also mean the fixture change removed the
teeth. To rule that out, the conformed (Writes-stripped) tree was run
against a deliberately broken kernel: in `resolve.go`'s `escapeOrRefuse`,
the uniform viability gate was bypassed —

```go
viable := escapes // MUTANT: gate bypassed on the escape path
```

replacing `viable, blocked := gate(escapes, in.Guards, view)` and its
`blocked != nil` return. Result (`a3-conformed-mutant-gate-bypass.out`):

```
--- FAIL: TestAdv2_EscapeEdgeMustNotBypassTheGuardSeam
    --- FAIL: .../guard_FALSE_must_not_rescue
    --- FAIL: .../guard_UNEVALUABLE_must_not_rescue
--- FAIL: TestAdv2b_EscapeEdgeMustNotBypassTheOwnedStateRequirement
--- FAIL: TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue
--- FAIL: TestFixup1e_AmbiguousClassEscapeIsGatedLikeAnyOtherCandidate
    --- FAIL: .../guard_FALSE_must_not_rescue_an_ambiguity
    --- FAIL: .../missing_owned_state_must_not_rescue_an_ambiguity
--- FAIL: TestFixupGateIsUniformAcrossOrdinaryAndEscapeCandidates
    --- FAIL: .../missing_owned_state
    --- FAIL: .../undecidable_guard
    --- FAIL: .../absent_guard_seam
FAIL	github.com/cwensel/intrastate/internal/resolve
```

All eleven escape-gate assertions still fire on a conformed fixture set.
The frozen evidence is not weakened. `resolve.go` was restored from a
byte-copy immediately after and re-verified green (`git diff --stat --
internal/resolve/resolve.go` empty).

## Variant 2 — also stripping `NextTags` (the A4-contingent widening)

Applied on top of variant 1: removed
`NextTags: []resolve.Tag{{Key: "status", Value: "Blocked"}},` from
`escapeRow` as well, leaving escape rows carrying only `RuleID`,
`SourceLocator`, `Outcome`, `Match`, `RequiresOwned`, and `Escape`.

Result (`a3-conformed-plus-nexttags.out`): exit 0, 154 `--- PASS`, 0
`--- FAIL`, sorted outcome set **identical to baseline**.

No test in this package reads `Plan.NextTags` on an escape-derived plan.
The `NextTags` assertions (`resolve_test.go` REQ-18, `mvv_test.go` leg 1)
are all on `legalInput()` / `singleMatchTable`. So the A4-contingent
widening is also non-breaking *for this package's frozen tests*.

Caveat worth carrying forward: this variant is evidence only about
`internal/resolve` test outcomes. A plan emitted from a `NextTags`-less
escape row carries an empty next state, which is a question for whatever
consumes `Plan.NextTags` downstream, not a question these tests answer.
A3 as written scopes only to frozen assertion outcomes, so this caveat
does not bear on the A3 verdict.

## Verdict

**A3 HOLDS.** Dropping `Writes` from escape rows changes no frozen
assertion outcome (154 PASS before, 154 PASS after, identical sorted
outcome sets), makes no assertion vacuous (every escape-row `Writes` read
is a failure-message argument, never a predicate; all substantive
`Writes`-equality assertions are on non-escape fixtures), and leaves every
escape-gate test discriminating (verified by mutation: an escape-path gate
bypass still fails 11 assertions across ADV-2, ADV-2b, Fixup-1d, Fixup-1e,
and the uniform-gate test).

The entry precondition can land without weakening RDR 0001's frozen
evidence. Two call-site `Writes` overrides (`adversarial_test.go:229`,
`fixup_test.go:103`) must be stripped alongside the `escapeRow` builder or
the conformance is incomplete.

## Revert confirmation

`git checkout -- internal/resolve/` followed by `git status --short`
reported no modifications under `internal/resolve/`, and
`go test ./internal/resolve/... -count=1` returned `ok`. The only working
tree entries are untracked files under `docs/rdr/`. All spike edits are
gone; the evidence files in this directory remain.
