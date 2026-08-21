# Resolve Research — RDR 0004 (scoped re-entry)

Model: claude-opus-5[1m]

Scope: the `Draft [revised from Final 2026-08-12]` re-entry. JDR 0001 §D3 added
one normative clause (read-accessor completeness). A1-A7 carry forward Verified;
this pass verifies only the added clause, which the re-entry qualifier declared
disturbed no assumption but which carries its own exactness claim.

## Reuse audit — does the code already do this?

Paths checked (`.rdr/env.md` reuse-audit set + the resolve kernel):
`internal/resolve/`, `internal/cli/`, `internal/cli/clierr/`,
`internal/cli/respond/`, `internal/cli/config/`.

**Verdict: NOT-PRESENT.** No reuse finding; the RDR introduces the capability
rather than re-specifying an existing one.

- No accessor executor exists under `internal/`. Every `accessor` hit in
  `internal/resolve/resolve.go` is a comment naming a layer not yet in code —
  `internal/resolve/resolve.go::Provenance` ("a tag read from a caller-provided
  artifact by the accessor layer"), `internal/resolve/resolve.go::Plan` ("the
  kernel describes writes; it never executes them").
- `internal/resolve/resolve.go::Input` declares `Owned []Tag`, and
  `internal/resolve/resolve.go::Tag` is `{Key, Value string}` — no error channel,
  no unread sentinel. `::Provenance` records where a tag came from, never whether
  the read succeeded.
- `internal/resolve/resolve.go::RefusalKinds` is a closed five-constant set
  (`no_match`, `ambiguous_match`, `owned_state_unavailable`, `guard_unevaluable`,
  `unmodeled_outcome`). No kind means partial or unreadable read.
- `internal/cli/clierr/clierr.go::CLIError` is the Go-error/exit-code path, and
  `::Refusal` states it "is not a CLI error and not the Go error path" — so it
  does not supply the missing channel either.

## Grounding — why the rule must sit at the accessor boundary

**Confirmed.** The kernel cannot recover the distinction after the fact, so
enforcement belongs upstream at the accessor seam, which is what the clause does.

`internal/resolve/resolve.go::assemble` flattens `in.Owned` into the `TagSet`
with `ProvenanceOwned`; no failure representation survives, because none entered.
`internal/resolve/resolve.go::gate` then calls `missingOwned` and, when the
result is non-empty, returns `Refusal{Kind: KindOwnedStateUnavailable, ...}`.

The deciding predicate in `internal/resolve/resolve.go::missingOwned`, iterating
`row.RequiresOwned`:

```go
if view.has(key, ProvenanceOwned) || seen[key] {
```

`internal/resolve/resolve.go::TagSet.has` is map presence plus provenance. So the
only question asked is "is this key in the snapshot map?" A key omitted because
the artifact genuinely lacks it and a key omitted because the read was truncated
are byte-identical inputs, and both yield `owned_state_unavailable`. That is
exactly the conflation JDR 0001 §D3 names.

## Negative results

- negative: read-accessor completeness — no corpus evidence sought in `DevRef`,
  `StateMachineLit`, `StateMachineRes`, `PapersFast`. Both questions are about
  THIS project's own code and the contract this RDR introduces, so the owning
  domain is project source, not external prior art. No external-behavior claim
  was made that would require a corpus pass.
