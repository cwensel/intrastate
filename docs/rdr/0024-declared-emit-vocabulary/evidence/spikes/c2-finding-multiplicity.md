Model: claude-opus-5[1m]

# Spike — does one load refusal report one finding, or one per hit?

Question: `0024:C2` claims "one findings[] entry PER HIT, so an emit refusal
reports beside, not instead of, a coexisting structural one", and MVV step 2
expects "exactly two blocking findings" from a single `intrastate lint` run.
This spike tests whether the load pipeline reports more than one finding.

## Run 1 — a model with three independent malformed tag declarations

Built from `models/examples/pricing-decision-table.toml` by replacing every
`kind = "enum"` with an unknown kind token (3 substitutions, 3 independent
defects in 3 distinct declarations).

```
$ go run ./cmd/intrastate lint --model <model> --as=json
```

Output (findings[] extracted):

```
code   = model-invalid
findings count = 1
  - malformed_tag_declaration | <model>:31
```

**Result: 1 finding for 3 independent defects.** Only the first is reported.

## Run 2 — a model with defects in two different categories

A model carrying both a reserved-key violation and an unknown kind token
reported a single finding, and in fact reported a *third* defect that
preceded both:

```
findings count = 1
  - missing_recognized_outcome_alphabet | <model>:1
```

**Result: 1 finding. The pipeline stops at the first refusal reached.**

## Why — the source

`internal/table/load.go::Load` returns `nil, err` at every refusal arm
(21 such returns in the file); there is no accumulator. `loadFindings`
(`internal/cli/flow_input.go:185`) maps that ONE error to exactly ONE
`clierr.Finding`:

```go
func loadFindings(path string, err error) []clierr.Finding {
	category, ok := table.CategoryOf(err)
	...
	finding := clierr.Finding{Code: string(category), ...}
```

`internal/cli/lint.go` calls it on the single `err` from
`table.LoadWithAdvisories` and returns immediately.

This is not merely how it happens to be built — RDR 0008's spec test
`internal/table/reserved_key_0008_test.go` ACTIVELY FORBIDS multi-error
load reporting:

```go
if joined, ok := err.(interface{ Unwrap() []error }); ok {
	t.Errorf("the load reported %d failures; exactly one is required",
		len(joined.Unwrap()))
}
```

## Verdict

The load pipeline is **fail-fast by locked contract**: one refusal, one
finding, first hit only. `0005:C1`'s text does say "one findings[] entry per
category hit", but no load path ever produces more than one hit to map — the
contract text and the built behavior already diverge at HEAD, independently
of this RDR.

Consequences for 0024:
- C2's "reports beside, not instead of, a coexisting structural one" is
  FALSE against the built pipeline and against `0008`'s test.
- MVV step 2's "exit nonzero, exactly two blocking findings" is
  UNACHIEVABLE without new accumulation behavior C2 does not specify.
