# A4 — Evidence derivation

Relocated from the RDR's A4 Evidence field at lock (Finalization Gate CHECK 9,
evidence-field budget on a `foundational` record). The field keeps the anchors
and points here; nothing below is summarized or dropped.

## HELD side

`internal/guard/assignment.go::valueAssignments` (called from
`internal/guard/product.go` with `m.Tags[key]`) renders the declared domain with
`strconv.Itoa` in its `int` arm, which emits no leading zero, `+`, or
whitespace, so only canonical spellings reach the lint side by THAT path (spike
`evidence/spikes/a4-render-path.md`).

That spike measured the product path, which under the narrowed scope is the
whole held side of this claim: `product.go::valueSatisfies` is reached from
`Denotation`, and `product.go::selectionOf` routes `BlockMatch` atoms away from
it. The other held-side path — `internal/graphlint/reach.go::heldValues` →
`canonicalValues` (no integer round-trip) — feeds `reach.go::atomAdmitsValue`
ONLY, which C4 excludes from typing, so lint and the kernel read those cells
identically and this claim does not quantify over them. Reconcile verified that
exclusion is structural, not incidental (both callers are `BlockMatch`-only;
derivation in `evidence/reconcile/report.md` §A4).

## LITERAL side

Today the literal is joined from authored bytes (`valueSatisfies`) and the
loader's kind check is `strconv.Atoi` (`internal/table/load.go::conformKind`),
which ACCEPTS `"00"`, `"01"`, `"+1"`; measured, an authored `n eq "00"` flips
lint from blocking `graph-coverage-gap` (exit 2) to clean (exit 0) — the
blocking→clean direction (fixture S7; spike `evidence/spikes/a4-lint-diff.md`).
C5 closes that leg by admitting only `strconv.Itoa(n) == authored`, so the
divergence the spike measured cannot be authored at all.

Bare floats are covered: a TOML `eq = -0.0` is rendered to the literal `"-0"` by
`internal/table/load.go::valueMembers` (`strconv.FormatFloat(t, 'g', -1, 64)`)
BEFORE the kind check, and `Atoi("-0")` succeeds while `Itoa(0)` is `"0"` — C5
refuses it, `conformKind` alone does not.

## Corpus coverage

The 123-model corpus shows no verdict flip, but it contains ZERO `eq`/`in` atoms
over an `int` tag, so that result measures coverage, not safety (S6).
