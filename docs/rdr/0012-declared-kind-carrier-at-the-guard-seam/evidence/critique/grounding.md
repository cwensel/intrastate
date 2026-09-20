Model: claude-sonnet-5

# RDR 0012 critique grounding — read-only source verification on main

## 1. reach.go::atomAdmitsValue — CONFIRMED
`internal/graphlint/reach.go:513`:
```go
return verdict != resolve.GuardFalse
```
GuardUnevaluable and the zero value (GuardFalse=0 is excluded by name, but
GuardTrue and GuardUnevaluable both satisfy `!= GuardFalse`) both fold into
"admissible" — i.e. Unevaluable is treated as admitting the value.

## 2. Callers of atomAdmitsValue — CONFIRMED, exactly two, one positive one negated
- `internal/graphlint/reach.go:484` (positive): `if atomAdmitsValue(a, v) { return true }`
- `internal/graphlint/analysis.go:438` (negated), inside `nodeMeetsAll`:
  `if !atomAdmitsValue(atom, v) { return false }`

## 3. heldValues / canonicalValues — CONFIRMED
`canonicalValues` (reach.go:343-345):
```go
func canonicalValues(value []string) []string {
	return slices.Compact(slices.Sorted(slices.Values(value)))
}
```
Only sort+dedupe, no Atoi/int round-trip. `heldValues` (reach.go:436-441)
returns `canonicalValues(value)` for finite-domain tags or `[]string{OpaqueValue}`
otherwise. Its inputs come from exactly two call sites: the root
(`reach.go:120`, from `m.Initial` `[initial]` cells) and a row's write
(`reach.go:317`, from `w.Value` in `writesOf(row)`, i.e. `rule.write` cells).

## 4. resolve.go::TagSet.matches — CONFIRMED
`internal/resolve/resolve.go:175-183`:
```go
func (s TagSet) matches(want []Tag) bool {
	for _, w := range want {
		tv, ok := s.tags[w.Key]
		if !ok || tv.conflicted || tv.value != w.Value {
			return false
		}
	}
	return true
}
```
Byte/string compare (`tv.value != w.Value`), no Atoi anywhere in this function.

## 5. load.go::conformKind — CONFIRMED, admits "00"/"+1"/"-0"; NO default case; empty Kind falls through unvalidated
`internal/table/load.go:1853-1865`:
```go
func conformKind(decl TagDecl, member string) error {
	switch decl.Kind {
	case "int":
		if _, err := strconv.Atoi(member); err != nil {
			return fmt.Errorf("%s is not an int", strconv.Quote(member))
		}
	case "bool":
		if member != "true" && member != "false" {
			return fmt.Errorf("%s is not a bool", strconv.Quote(member))
		}
	}
	return nil
}
```
Bare `strconv.Atoi`, no canonical round-trip check. Verified empirically:
`strconv.Atoi("00")` -> 0,nil; `strconv.Atoi("+1")` -> 1,nil; `strconv.Atoi("-0")`
-> 0,nil — all three ADMITTED. No `default` case. For Kind == "" (empty
string), neither `case` matches, the switch falls through, and the function
returns `nil` — i.e. an empty-Kind tag's members pass conformKind with NO
validation at all (same true for conformDomain's switch, load.go:1870-1891,
which also has no default and no "" case, so domain checks are likewise
skipped for empty Kind).

## 6. DeclaredKinds-equivalent constructor for the guard evaluator — REFUTED
No such constructor/helper exists anywhere in `internal/guard` or
`internal/table`. The only "declaredKinds" artifact in the codebase is an
unrelated frozen vocabulary list, `internal/table/model.go:84`:
```go
var declaredKinds = []string{"enum", "bool", "int", "set", "scalar"}
```
— a static list of the five legal `Kind` token strings, not a per-model
map/constructor keyed by tag, and it does not omit or filter anything by
value (it has no notion of a Model instance at all). Exhaustive symbol
search of `internal/guard/*.go` (non-test) function list turned up no
map-building constructor filtered on `Kind != ""` (closest candidates —
`declarationOf`, `DeclarationOf`, `AssignmentCount` — build per-tag
declaration/count values, none omit empty-Kind keys).

## 7. load.go::valueMembers — CONFIRMED
`internal/table/load.go:1756-1757`, float64 arm:
```go
case float64:
	return []string{strconv.FormatFloat(t, 'g', -1, 64)}, nil
```
`FormatFloat(-0.0, 'g', -1, 64)` renders the literal string `"-0"`. This
call happens inside `valueMembers`, invoked at load.go:1688 — BEFORE the
kind/domain check `conform(decl, "eq", members)` at load.go:1701. So a TOML
`-0.0` becomes the string `"-0"` prior to any kind check, confirmed by call
order (1688 precedes 1701) inside `loadInitial`.

## 8. flow_input.go::canonicalValue — CONFIRMED (conforms before resolution; refuses non-true/false bool); admits "n=07" on int TODAY
`internal/cli/flow_input.go:738-776`. For non-set kinds it calls
`table.ConformValue(decl, value)` (line 751) which is `conformKind` +
`conformDomain` (load.go:1813-1818) — this runs at parse/CLI-input time,
before the value ever reaches `resolve.Resolve`. `conformKind`'s bool arm
(load.go:1859-1862) refuses any value that is not literally `"true"` or
`"false"`. For `--tag n=07` on an int-declared tag: `conformKind`'s int arm
is bare `strconv.Atoi("07")`, which returns `7, nil` (verified) — no
leading-zero rejection — so `"07"` is ADMITTED today, not refused.

## 9. Guard evaluator construction sites, non-test code — CORRECTED, three real sites but not the claimed three
Exact `guard.Evaluator{}` (or bare `Evaluator{}` within package guard)
literal construction in non-test code:
- `internal/cli/flow_resolve.go:561` — `func guardSeam() resolve.GuardEvaluator { return guard.Evaluator{} }`
- `internal/graphlint/reach.go:507` — inline in `atomAdmitsValue`
- `internal/guard/product.go:542` — inline in `valueSatisfies` (`seam := Evaluator{}`) — THIS ONE WAS NOT IN THE CLAIMED LIST

`internal/cli/flow_next.go::probeRow` does NOT itself construct a
`guard.Evaluator{}` — its own doc comment (flow_next.go:435) says it reaches
the evaluator "through the same `guardSeam()` `flow resolve` hands the
kernel," i.e. probeRow is a consumer of the `flow_resolve.go` construction
site, not an independent construction site. So the claimed list
{flow_next.go::probeRow, flow_resolve.go::guardSeam, reach.go::atomAdmitsValue}
is WRONG on one item (probeRow) and MISSING one item (guard/product.go:542).
The real count of non-test construction sites is three, but the membership
differs from the claim.

`internal/table/normalize.go::renderWrites` does NOT exist — no such
function/file content matched; `internal/table/normalize.go` was checked
directly and grep for `func renderWrites` returned nothing. This part of
the claim is REFUTED as stated (no such constructor).

## 10. Zero-value guard.Evaluator{} test-code sites and touched-package count — CONFIRMED counts
Exact literal `guard.Evaluator{}` occurrences in `internal/guard/*_test.go`:
11 total, across two files:
- `guard_evaluator_0003_test.go`: 7 (lines 84,104,124,149,164,186,220)
- `guard_fixtures_0003_test.go`: 1 (line 294)
- `guard_narrowing_0003_test.go`: 3 (lines 142,163,568)
(claim 11's test uses a separate non-literal `var ev guard.Evaluator` form,
not counted here.)

Packages importing `internal/guard` (non-test source): exactly 2 —
`internal/cli` and `internal/graphlint`. Including test-only importers adds
`internal/guard` itself (its own test suite) and one evidence-spike file
under `docs/rdr/0010.../evidence/spikes/a13-zerodim` (not production code).

## 11. guard_evaluator_0003_test.go::TestReq34... — CONFIRMED, stronger than "no TagSet field"
`internal/guard/guard_evaluator_0003_test.go:21-68`. The assertion is
compound, not merely "no field of type resolve.TagSet":
1. `fn.NumIn() != 3` — Evaluate must take exactly receiver+atom+value (no view param at all, by arity).
2. Loop over ALL params checking `fn.In(i) == viewType` where `viewType = reflect.TypeOf(resolve.TagSet{})` — no parameter of any position may be a TagSet.
3. `reflect.TypeOf(ev).NumField() != 0` — the receiver struct itself must carry ZERO fields of ANY type, not just no TagSet field ("a view-free evaluator holds no state to read one from").
4. Behavioral: `ev.Evaluate(exists atom, "anything")` must equal `GuardUnevaluable` (existence atoms are kernel-owned, not evaluator-decided).
5. Runs the full cross-RDR contract oracle: `resolve.TestGuardEvaluatorContract(t, ev)`.
So the claim "no field of type resolve.TagSet" understates the test; it actually pins zero fields total plus arity plus existence-atom behavior plus the full contract suite.

## 12. strconv.Atoi platform-width pin — REFUTED, no such pin exists
Exhaustive grep for width-related tokens (`int32`, `MaxInt32`, `MaxInt64`,
`bits.UintSize`, "platform"/"width") across `internal/table`,
`internal/guard`, `internal/resolve`, `internal/graphlint` found nothing.
Every `strconv.Atoi` call site in these packages (`internal/table/edit.go:163`,
`internal/table/load.go:1856,1880`, `internal/guard/grammar.go:115,119`) uses
Go's platform-native `int` with no documented or enforced width contract in
the conformance suite or the guard package.
