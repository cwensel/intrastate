Model: claude-opus-5[1m]

# Grounding Step-0 — Codebase Claim Sweep, cli/0020

Scope: 24 `source-anchor` edges (all `resolved:true`; the 25th, `path::Symbol`
at :846, is TEMPLATE placeholder text inside a surviving guidance block — a
lint finding, not an anchor). Decision Rationale carries `Ground-sweep: clean
(17 anchors)` from propose, so the prose sweep scoped to the post-propose diff
(dc203d2..HEAD: the Resolve evidence layer, C1's empty-value hoist, the four
normative fixtures F1–F4) plus the inverse sibling-discriminator check.

18 unanchored prose claims read at source. CONFIRMED claims carry no finding.
Four REFUTED / NOT-FOUND below.

## G1 — NOT-FOUND. A1's "one asserting test" does not exist

A1 Evidence (`0020:A1`) states the "is not set-valued" message literal has
"one producer … and one asserting test —
`internal/cli/flow_input_0005_test.go`
`TestReq29And80_WrongKindAndEmptyTagValuesAreFlowTagInvalid`, subtest
`array-for-scalar-key`".

Found instead: `grep -rn "is not set-valued" --include="*_test.go"` returns
**zero** matches repo-wide. The string's only occurrence anywhere is the
emission site, `internal/cli/flow_input.go:721`. The named subtest exists
(`internal/cli/flow_input_0005_test.go:278`) and does reach that refusal arm
with `--tag profile=["mid"]`, but it asserts only the code and param:

```go
ce := requireRefusal(t, "flow-tag-invalid", 2, …, "--tag", `profile=["mid"]`, …)
if ce.Param != "profile" { … }
```

`requireRefusal` (`internal/cli/flow_harness_0005_test.go:306`) checks `Code`
and exit status; nothing inspects `ce.Message`. So the message text is
asserted by no test at all.

Direction of the error is favourable — the load-bearing negative ("no
assertion, caller branch, or fixture requires `flow-tag-invalid` over an
UNDECLARED key") is *stronger* than stated, not weaker; A1's verdict stands.
But the count is wrong as written, and the same wrong count is repeated at
`0020:§key-discoveries` ("one producer and one asserting test") and leaned on
by the Failure-Modes/Risks argument.

## G2 — REFUTED. A3's "byte-equality lives in exactly one place"

A3 Evidence (`0020:A3`) states byte-equality "lives in exactly one place,
`internal/accessor/executor.go::verifyReadBack`".

Found instead: a second tag-value equality site,
`internal/cli/flow_exec.go::readBackFindings` (`internal/cli/flow_exec.go:632`):

```go
got, held := observed[want.Key]
if held && got == want.Value {
    continue
}
```

Its inputs are `accessor.Refusal`'s `Expected` / `Observed`
(`internal/accessor/model.go:330`), i.e. the same write-plan-vs-accessor-re-read
domain `verifyReadBack` already validated, re-rendered CLI-side as per-key
findings (REQ-103). It never reads the `--tag` observed slice.

So A3's conclusion holds — no byte-equality path consumes an undeclared
observed value — but its stated basis (sole site) is false. The accurate
claim is two sites, both confined to the write-plan / re-read domain.

## G3 — REFUTED (incomplete). A2's kernel-view reader enumeration omits a reader

A2 Evidence (`0020:A2`) states the kernel view "is read only through `matches`
(iterating `row.Match` — rule atoms only), `has`, and guard `Lookup`".

Found instead: a fourth reader of `TagSet.tags` —
`internal/resolve/resolve.go::conflicting` (`internal/resolve/resolve.go:154`):

```go
func (s TagSet) conflicting(key string) bool {
	return s.tags[key].conflicted
}
```

called from `internal/resolve/guard.go:204`. It is **value-blind** — it reads
only the `conflicted` bool, never `.value` — so A2's conclusion (an undeclared
key's VALUE is structurally unreadable) is unaffected. But a closed-world
enumeration that omits a member is not closed; the omission is the defect,
since the whole force of A2 is exhaustiveness.

## G4 — REFUTED. A1 mis-describes the fixture's declaration

A1 Evidence describes the `array-for-scalar-key` subtest as running "against a
key DECLARED at `internal/cli/flow_fixtures_0005_test.go` `[tags.profile]`" —
correct — but the surrounding argument treats `profile` as the declared-scalar
arm generically. The fixture (`internal/cli/flow_fixtures_0005_test.go:110`) is:

```toml
[tags.profile]
provenance = "observed"
kind = "enum"
domain = ["mid", "foundational"]
single_valued = true
```

`kind = "enum"`, not `"scalar"`. It is scalar-SHAPED (`single_valued = true`,
so `isSet` is false and the `!isSet && looksArray` arm fires as claimed), so
the evidence's mechanical point is sound. Minor, but "declared scalar" is not
what the fixture says, and C1 speaks in terms of the five kind tokens.

## Inverse sibling-discriminator check — PASS, no finding

C1's identity rule ("key ABSENT from the loaded model's normalized tag table,
the two-value `m.Tags[key]` lookup") mints no new discriminator. The same
two-value presence idiom already decides declared-vs-undeclared at six sites:

- `internal/table/load.go::accessorTable` (`_, ok := l.model.Tags[key]`, :979)
- `internal/table/load.go::loadInitial` (:1684)
- `internal/table/normalize.go::atomsFromBlock` (:26)
- `internal/table/normalize.go::renderWrites` write arm (:526) and clear arm (:584)
- `internal/graphlint/analysis.go` (:142, :204) and `coverage.go` (:363) — lint-only

No site anywhere uses a `Kind == ""` zero-decl sentinel as an admission-time
discriminator; the zero-`TagDecl` shape appears only downstream inside
`canonicalValue` / `ConformValue` as the shape-only fallback. `0020:D-identity`'s
"Reuse" verdict and the Existing Infrastructure Audit row are accurate.

## Confirmed (no finding owed)

All five `CatUnknownTag` refusal sites resolve exactly as A2 names them, and
the list is exhaustive at five. `parseTags`'s admission order is grammar →
reserved → owned → duplicate → `canonicalValue`, with `codeTagDuplicate`
(:663) strictly preceding the `canonicalValue` call (:677) — C1's
`D-selection-predicate` order claim confirmed. The empty-value arm does sit
inside `canonicalValue`'s `!isSet` branch (:723), and a declared SET key given
an empty value does take the "is set-valued and takes a JSON array literal;
got " arm (:735) — C1's hoist rationale confirmed on both legs.
`parseWrites` proves `writerFor` before `canonicalValue` and carries the
verbatim comment "the lookup here is total, unlike `parseTags`'s"
(`internal/cli/flow_state.go:461-470`). The loader refuses a kind outside the
five tokens at `internal/table/load.go::tagDecl` (:858) via `IsDeclaredKind`
(`internal/table/model.go:81`), tokens `enum, bool, int, set, scalar`.
`CheckInput` reads `Observed` only for the reserved key NAME
(`internal/resolve/precondition.go:32`). `merge` does compare a repeated key's
value (:235) and `parseTags` refuses the duplicate first (:663). The three
named non-comparing consumers (`observedTagMap`, `kernelTags`, the `groupEcho`
`observed` echo) are all verbatim copies with no comparison; no omitted
consumer found.

Two RDR quotations are faithful paraphrases rather than byte-exact source
text — `parseTags`'s guarded-lookup comment (source: "which conforms
everything and leaves the shape-only behaviour a caller already relies on
intact") and `EmitValue`'s doc (source: "never parsed, canonicalized, or
converted … compared by exact byte equality"). Both are materially accurate;
recorded here, no finding.
