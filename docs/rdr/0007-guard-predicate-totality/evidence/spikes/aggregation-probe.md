Model: claude-opus-5[1m]

# RDR 0007 — Stage 4 aggregation spike

Probes the two resolution-level aggregation behaviors RDR 0007's normative
aggregation clause asserts, against the shipped kernel at `44b83e3`.

Both probes are throwaway tests compiled into `package resolve_test` (reusing
`fixtureGuards`, `noMatchInput`, `escapeRow`, `singleMatchTable` from
`fixtures_test.go`), run, then removed. The probe source is inlined under
**Probe source** below rather than left as a `.go` file beside this file: as a
standalone package under `docs/` it cannot compile (its fixtures live in
`internal/resolve`), which broke `go vet ./...` for the whole module.

## Command

```sh
# paste the Probe source below into the file, then:
go test ./internal/resolve/ -run 'TestProbe' -v
rm internal/resolve/zz_probe_test.go
```

## Output (verbatim)

```
    zz_probe_test.go:22: PROBE A: kind="guard_unevaluable" guard="unknown-predicate"  (RDR 0007 clause wants no_match)
--- PASS: TestProbeA_UnevaluableEscapeOverNoMatch (0.00s)
    zz_probe_test.go:47: PROBE B: kind="guard_unevaluable" guard="unknown-predicate"  (blocks true sibling: correct)
--- PASS: TestProbeB_UnevaluableBlocksTrueSibling (0.00s)
PASS
ok  	github.com/newcoinc/intrastate/internal/resolve	0.361s
```

## Probe A — unevaluable escape row over a `no_match` candidate set

Input: zero matching ordinary candidates (`no_match` condition) plus one
modeled `no_match` escape row whose guard the seam cannot decide.

Observed: `guard_unevaluable`, carrying the *escape row's* guard text.

RDR 0007's aggregation clause requires the opposite:

> an unevaluable escape row MUST NOT convert a candidate-set refusal into
> `guard_unevaluable`

**The shipped kernel contradicts the clause.** The path is
`resolve.go::escapeOrRefuse`, which returns the escape set's own blocking
refusal in place of the candidate-set refusal `r`:

```go
viable, blocked := gate(escapes, in.Guards, view)
if blocked != nil {
    return refuse(in, *blocked)   // replaces r; r survives only via case 0 below
}
```

This is **frozen deliberately**, not incidental —
`adversarial_test.go::TestAdv2_EscapeEdgeMustNotBypassTheGuardSeam`, subtest
`"guard UNEVALUABLE must not rescue"`, asserts exactly this kind, with a
written rationale ("an undecidable escape predicate is exactly the
guard_unevaluable condition"). `fixup_test.go::TestFixup1d_GuardedEscapeEdgeWithNilSeamMustNotRescue`
freezes the same conversion by the nil-seam route, and
`fixup_test.go::TestFixupGateIsUniformAcrossOrdinaryAndEscapeCandidates`
asserts the escape path returns the *same* refusal kind as the ordinary path.

Note the asymmetry that likely produced the RDR's error: on an escape row,
guard-FALSE and missing-owned-state leave the original refusal standing
(GuardFalse is pruned inside `gate`, yielding `len(viable)==0` → `case 0:
return refuse(in, r)`), whereas unevaluable exits through `blocked != nil`.
The clause reads as if it generalized the guard-FALSE behavior to
unevaluable, which the kernel does not do.

## Probe B — unevaluable candidate beside a decided-true sibling

Input: row `rdr.true` (guard decided TRUE) and row `rdr.unev` (guard
unevaluable), both matching.

Observed: `guard_unevaluable` naming the unevaluable row's guard — no plan.
The decided-true sibling is **not** selected past the unevaluable row.

This confirms the resolution-level half of RDR 0007's clause, and the MVV
scenario *unevaluable-blocks-true-sibling*. Mechanism in `resolve.go::gate`:
the `return nil, &Refusal{...}` on `len(undecidable) > 0` discards the
accumulated `selected` slice entirely.

**No shipped test freezes this.** `TestAdv3` has two rows but both are
unevaluable (it pins payload order-stability); `twoRowsOneGuardFalseTable`
pairs true with false. The behavior is correct but held by implementation
only — Phase 2's golden vector suite is where it should be pinned.

## Verdict

- Probe A: RDR 0007's escape half of the aggregation clause is **refuted** by
  frozen kernel behavior. The RDR text must change, or three shipped tests
  and RDR 0001 deviation D5 must be re-opened.
- Probe B: RDR 0007's candidate half is **confirmed**, currently untested.

## Probe source

Recreate as `internal/resolve/zz_probe_test.go` to re-run (see **Command**
above). Kept inline, not as a sibling `.go` file — under `docs/` it is a
package that cannot compile, and it broke `go vet ./...` module-wide until it
was removed.

```go
package resolve_test

import (
	"testing"

	"github.com/newcoinc/intrastate/internal/resolve"
)

// Probe A: unevaluable ESCAPE row over a no_match candidate set.
// RDR 0007's aggregation clause says this MUST NOT become guard_unevaluable.
func TestProbeA_UnevaluableEscapeOverNoMatch(t *testing.T) {
	esc := escapeRow("rdr.escape.unev", "flows/rdr.toml:90", resolve.KindNoMatch)
	esc.Guard = "unknown-predicate"
	in := noMatchInput()
	in.Table.Revision = "probe-a"
	in.Table.Rows = append(in.Table.Rows, esc)
	in.Guards = fixtureGuards{}
	got, _ := resolve.Resolve(in)
	if got.Refusal == nil {
		t.Fatalf("PROBE A: plan=%+v", got.Plan)
	}
	t.Logf("PROBE A: kind=%q guard=%q  (RDR 0007 clause wants no_match)", got.Refusal.Kind, got.Refusal.Guard)
}

// Probe B: unevaluable CANDIDATE row beside a decided-true sibling.
// RDR 0007 MVV scenario *unevaluable-blocks-true-sibling*.
func TestProbeB_UnevaluableBlocksTrueSibling(t *testing.T) {
	tbl := singleMatchTable()
	tbl.Revision = "probe-b"
	tbl.Rows[0].RuleID = "rdr.true"
	tbl.Rows[0].Guard = "always"
	unev := tbl.Rows[0]
	unev.RuleID = "rdr.unev"
	unev.SourceLocator = "flows/rdr.toml:20"
	unev.Guard = "unknown-predicate"
	tbl.Rows = append(tbl.Rows, unev)
	in := resolve.Input{
		Flow: "rdr", Table: tbl,
		Owned:      []resolve.Tag{{Key: "status", Value: "Draft"}},
		Recognized: "successful",
		Guards:     allGuardsTrue("always"),
	}
	got, _ := resolve.Resolve(in)
	if got.Refusal == nil {
		t.Fatalf("PROBE B: MASKING — plan=%q", got.Plan.RuleID)
	}
	t.Logf("PROBE B: kind=%q guard=%q  (blocks true sibling: correct)", got.Refusal.Kind, got.Refusal.Guard)
}
```
