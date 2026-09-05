Model: claude-sonnet-5 (delegated grounding pass)

# Source grounding for 3amigo findings — read at HEAD

## Q1 — duplicate vs empty-value precedence in `parseTags`

`internal/cli/flow_input.go::parseTags` (642-684), source order per entry:

- :649-650 grammar split (`strings.Cut`), refuse `codeTagInvalid` on malformed
- :654-665 `switch`: reserved (:655), owned (:659), duplicate `case seen[key]:` (:663)
- :667 `seen[key] = true` — AFTER the duplicate case, before the decl lookup
- :675 `decl := m.Tags[key]` (one-value, guarded)
- :677 `canonicalValue(key, value, decl, "tag")` — the empty-value check lives INSIDE, at :723
- :678 returns on first error

**Precedence at HEAD**: the duplicate check (:663) is strictly before the
empty-value site (:723). For a key already `seen`, duplicate wins; the
empty-value refusal is observable only on a key's FIRST appearance. `parseTags`
returns on the first error, so `--tag k= --tag k=` refuses empty-value on the
first entry and never reaches the second.

## Q2 — the empty-value check's site and callers

`canonicalValue` (714-752); the arm at :718-726:

    if !isSet {
        if looksArray { ... "is not set-valued; got the array literal " ... }
        if value == "" { ... "the tag `"+key+"` was given an empty value" ... }

Message code is `codeInvalidFor(flag)`: `"tag"` → `flow-tag-invalid`,
`"write"` → `flow-write-invalid` (:757-762, constants :36-41).

Two call sites only:
- `flow_input.go:677` — `parseTags`, flag `"tag"`
- `flow_state.go:470` — `parseWrites`'s `admitWrite`, flag `"write"`

**Confirmed**: `parseWrites` reaches `canonicalValue` only after
`writerFor(m, key, codeWriteUnbound)` succeeds (:464), so its `m.Tags[key]`
lookup is total — the key is declared by construction. The source comment at
:467-469 says so explicitly. This is the RDR's Technical Design claim, verified.

## Q3 — `codeInvalidFor`'s doc comment vs `--outcome`

Comment at :754-756 claims "`--tag` and `--outcome` take `flow-tag-invalid`".

`--outcome` is read at `flow_resolve.go:237` via `cmd.Flags().GetString` and
refuses inline at :239 with its own `userErr(codeTagInvalid, "outcome", ...)`.
It never enters `parseTags` or `canonicalValue`.

**Verdict**: the comment is accurate as to the resulting code VALUE
(`flow-tag-invalid` either way) but misleading as to path — `--outcome` does not
route through `codeInvalidFor`. Pre-existing, not introduced by 0020, and 0020
makes no claim about `--outcome`.
