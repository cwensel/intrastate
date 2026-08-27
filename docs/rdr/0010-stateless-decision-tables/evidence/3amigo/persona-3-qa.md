Model: claude-opus-5[1m]

# 3amigo — Persona 3 (QA / Tester), RDR 0010

Six findings: 1 High, 3 Medium, 2 Low. Scenarios S1–S6, MVV, and C1–C5 are
unusually well-instrumented for absence oracles — the Oracle table at
`0010:§oracle` already discriminates MVV step 3 against step 2, S4 carries a
match-only negative control and a stray-`terminal` control, and S5 carries
the `{}`-never-`null` row that the MVV table cannot exhibit. Findings below
are the residue after those controls are credited.

### P3-1 — C4's text-mode emit rendering contradicts S5, and neither matches the renderer the CLI actually has  [High]

**Anchor**: `0010:C4`, `0010:S5`, `0010:MVV` (step 4)

**Finding**: C4 states "Under `--as=text` the plan MUST render each emit pair
as one `key=value` line." S5's **Expected** states "text renders one
`key=value` line per pair, and **no line at all** for an unauthored block."
These are two different obligations, and the second is not derivable from the
first — C4 says nothing about the unauthored case in text mode, having spent
its `{}`-never-`null` clause on the JSON payload only.

Worse, both are refuted by the surface that would have to produce them.
`flow resolve`'s text mode does not have a verb-specific renderer: the payload
goes through the generic `writeTextPayload` in
`/Users/cwensel/sandbox/newcoinc/intrastate/internal/cli/respond/text.go:42`,
which round-trips the payload's JSON and emits path-qualified leaf lines via
`flatten`/`label`/`join`. Over `resolvePayload` (declared at
`/Users/cwensel/sandbox/newcoinc/intrastate/internal/cli/flow_resolve.go:31`)
an emit block would render as `emit.next: propose` — path-qualified, `": "`
separated — not `next=propose`. And `flatten`'s `map[string]any` arm at
text.go:68 returns `label(path) + "(none)"` for an empty map, so an
unauthored block renders the line `emit: (none)` — precisely the "no line at
all" S5 forbids. That arm is load-bearing and documented as such at
text.go:61-64 ("a renderer that dropped empty containers would erase exactly
the distinction REQ-66 turns on"), so C4 is silently asking for a
special-case in a renderer whose no-special-cases property another RDR
depends on — and no contract in this RDR authorizes touching it.

The Technical Design at `0010:§technical-design` says only "`resolvePayload`
gains `emit`; text mode renders it," which is consistent with the generic
renderer doing the work and inconsistent with both C4 and S5.

**Prevents**: the text-mode arm of S5. I cannot write
`TestFlowResolveDecisionTableTextEmit` because I cannot decide what string to
assert. Three candidate assertions are each defensible and mutually
exclusive: (a) `stdout` contains the line `next=propose` (C4 literally);
(b) `stdout` contains `emit.next: propose` and, on the unauthored row,
`emit: (none)` (the renderer as built, no code change); (c) `stdout` contains
`next=propose` and, on the unauthored row, **no** line whose path begins
`emit` (S5 literally, requiring a special case in `writeTextPayload` that
this RDR nowhere authorizes and that text.go:61-64 argues against). Pick
wrong and the test pins the opposite of the intended behaviour. The
`--as=text` half of MVV step 4 is blocked for the same reason.

### P3-2 — S6 requires a pre-change binary, which the stated harness cannot produce  [Medium]

**Anchor**: `0010:S6`, `0010:§testing-strategy` (Done clause), `0010:A8`

**Finding**: S6's **Expected** ends "and a pre-change binary refuses the
fixture `unknown_schema_field` (A8)." But S6's own **Scenario** scopes the
sweep to "every checked-in model and fixture (`models/rdr.toml` included)
under `make check`," and the Testing Strategy's Done clause makes `make check`
the gate for every scenario. `make check` is `fmt-check vet lint build
graph-lint test` over the working tree
(`/Users/cwensel/sandbox/newcoinc/intrastate/Makefile`) — it builds exactly
one binary, the post-change one. There is no second toolchain, no vendored
prior binary, and no `git worktree`/`stash` step named anywhere in the four
implementation phases. A8's own evidence was captured out-of-band as a spike
(`evidence/spikes/a8-strict-decode.out`), which is the right method for a
one-time verification but is not a repeatable regression test, and S6
promotes the claim into the permanent suite without giving it a mechanism.

**Prevents**: the old-binary arm of `TestRegressionSweep`. I cannot write a
test asserting "a binary built from the parent commit exits 2 with category
`unknown_schema_field` over the new decision-table fixture," because no step
in the plan produces that binary and the RDR does not say whether the claim
is a permanent CI obligation (needing a build-two-binaries step, a `//go:build`
tag, or a scripted `git worktree` target) or a one-shot spike already
discharged by A8 and out of scope for the suite.

### P3-3 — "no existing golden output changes except for the `emit` payload field and dump column" names no artifact set to diff  [Medium]

**Anchor**: `0010:§testing-strategy` (Done clause), `0010:S6`

**Finding**: The Done clause reads "no existing golden output changes except
for the `emit` payload field and dump column," and S6 restates it as
"byte-identical lint, resolve, and dump output apart from the `emit`
field/column." This is the RDR's widest regression obligation and it is
stated as an adjective over an unnamed set. There is no golden-file
infrastructure to point at — `find internal -name '*.golden'` returns zero;
assertions live inline in Go test files and in `testdata` TOML/JSON fixtures.
So "existing golden output" has no denotation, and neither does the
exception clause: A4 establishes that **103** files under
`internal/table/testdata/` carry an explicit `[dump]` list and all 103 must
be edited to add `emit`. Those 103 edits change checked-in fixture bytes.
Are they "the dump column" exception, or are they violations of
"byte-identical"? Both readings are available and they classify 103 files in
opposite directions.

**Prevents**: the sweep's pass/fail decision. I cannot write a green/red
gate for `TestNoUnintendedOutputDrift` — I would need (i) an enumerated
baseline set (which files, which commands, which flags) and (ii) a rule that
mechanically separates a licensed diff from an unlicensed one, e.g. "a diff
is licensed iff every changed line differs only by the addition of an
` emit=[…]` token or an `"emit":` member." Without (ii), the reviewer of a
failing sweep has an adjective, not a criterion, and the natural failure mode
is to accept every diff by eye.

### P3-4 — S4's match-only negative control asserts "lints clean while incomplete" with no positive discriminator  [Medium]

**Anchor**: `0010:S4`, `0010:MVV` (Oracle row 1)

**Finding**: S4 requires "a **match-only negative control** — the same table
discriminating by `[rule.match.<key>]` instead of guard atoms" whose Expected
is "the match-only control lints clean **while incomplete**, pinning why the
fixtures must author guard atoms (C5)." The Oracle table's row 1 states the
same. This is an absence oracle and — unlike MVV step 3, which the Oracle
table explicitly discriminates against step 2 over *the same fixture* — it is
given no discriminator at all. `exit 0, findings []` over the match-only
fixture is exactly what a lint that crashed, short-circuited, or never
reached `checkGroups` would also produce. The control's whole purpose is to
prove that the *product* is empty (`internal/guard/product.go::Dimensions`
collecting from `guard.all`/`guard.unless` only), and emptiness of the
product is not observable from `lint`'s finding list: there is no CLI surface
that reports the scoped product's dimensions or cardinality.

The asymmetry matters because the risk register at `0010:§risks-and-mitigations`
names this exact path "the one path that produces it silently," and this
control is the only mitigation offered for it.

**Prevents**: `TestMatchOnlyControlIsVacuous`. I cannot write an assertion
that distinguishes "clean because the product is empty" from "clean because
the check did not run." What I need and cannot get from the record: either
(a) a positive companion assertion over the same match-only fixture — e.g.
introducing a *guard*-discriminated overlap into it and asserting
`graph-overlap` fires, proving lint reached the group machinery over that
fixture; or (b) a named non-CLI observable to assert against
(`guard.Dimensions(rows)` returning empty is package-visible and directly
assertable), which the RDR nowhere designates as the control's oracle.

### P3-5 — C5's "lint MUST NOT report the class as a finding of any severity" is an unbounded absence with no closed set to assert over  [Low]

**Anchor**: `0010:C5`, `0010:A10`, `0010:S4`

**Finding**: C5 ends "lint MUST NOT report the class as a finding of any
severity." A10 does the enumerable work — it partitions all fourteen 0006
codes into exercised-unchanged and provably-silent — and S4's Expected asks
for "exact finding lists against the full taxonomy (A10, A9)," which is
writable and strong. But C5's own clause is stated over a different quantity:
not "these codes do not fire" but "the class is not reported," which has no
enumerable extension. A finding "reporting the class" could be a new code, an
existing code with a class-mentioning `Message`, or a `Detail` naming
`decision-table`. `clierr.Finding` carries free-text `Message`, so the third
form is reachable without changing the taxonomy, and A10's code-partition
does not cover it.

**Prevents**: a direct test of C5's last sentence. `TestLintNeverReportsClass`
has no assertion I can write beyond the one S4 already gives (exact finding
lists over the fourteen-code taxonomy). If that is all C5 means, the clause is
redundant with A10 and should say so; if it means more — e.g. "no finding's
`Message` or `Detail` contains the substring `decision-table`" — that is a
writable assertion the record does not state, and the difference is exactly
what a test would have to encode.

### P3-6 — C1's "detail names the class and the offending count" fixes no form for the count  [Low]

**Anchor**: `0010:C1`, `0010:S1`

**Finding**: C1 requires the disagreement refusal's "detail names the class
and the offending count," and S1's Expected repeats "refuse `malformed model
declaration` with the class and count in the detail." `Failure.Detail`
(`/Users/cwensel/sandbox/newcoinc/intrastate/internal/table/category.go:79`)
is a free-text string, which is assertable by substring — but "the offending
count" is not pinned to a quantity. In the decision-table-with-owned-tags
direction the offending count is the owned tag count (≥1); in the
state-machine-with-none direction it is zero — and a detail asserting
"declares 0 owned tags" contains a digit that a substring test for "the
count" cannot distinguish from an incidental one. Neither C1 nor S1 gives a
token form, and 0002:C24's convention (cited in A8) is to assert on
*category*, never message text, which cuts against asserting the detail at
all.

**Prevents**: the detail half of S1's two disagreement rows. I can assert
`CatMalformedModelDeclaration` and exit status for all four S1 arms, but I
cannot write the `strings.Contains(detail, …)` assertion C1 obliges without
inventing the token myself — and an invented token pins the implementation to
whatever I guessed. Either C1 should fix a token form (as 0008 does for its
rule identifier, per `internal/table/reserved_key_0008_test.go`) or S1 should
drop the detail assertion and rest on the category.

## Widening

My owned set is `S`, `MVV`, and `C`. Three findings sent me outside it.

- **P3-1** — C4 and S5 disagree about text output, so I went to the CLI to
  see which one the existing surface produces. Neither: the answer was in
  `internal/cli/respond/text.go` and in `0010:§technical-design`, whose
  one-line "text mode renders it" is the third position.
- **P3-2, P3-3** — S6's expected outcomes are scoped to `make check`, which
  is defined in the Makefile, not the record; and the Done clause that binds
  every scenario lives in `0010:§testing-strategy` prose, outside any `S`
  element. The 103-file `[dump]` count that makes P3-3 bite came from
  `0010:A4`.
- **P3-4** — S4's control is stated as an absence, so I went to
  `0010:§risks-and-mitigations` to see what it was mitigating, and to
  `internal/guard/product.go` to see whether the emptiness it claims is
  observable anywhere.
