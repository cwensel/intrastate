# RDR 0009 — kata flight drain (`batch:rdr-0009`)

```
flight: 1 shipped, 0 stopped, 0 skipped  (over 2 waves)
  shipped: 0wc8=d39f68e
  stopped: —
  skipped: —
```

## Wave 1

**Resolved:** `0wc8` (1 eligible). `6cek` filtered out at resolution —
`kind:rdr-seed` is refused by the ship gate and routes to RDR authoring.

### Scope-review gate (ran by default)

`0wc8` → **IN-SCOPE**, `lifecycle:reviewed`. Two questions were settled against
source rather than against the kata body:

1. **The REQ-90 deferral did NOT block it — the kata's own "course of action"
   was wrong.** The kata said to wait for "RDR 0002's build". `req-list.md:793-802`
   records that stated reason as **stale** (RDR 0002 is `Implemented` at HEAD),
   and `:625-630` clarifies REQ-90's deferral is one of **ownership**, not
   readiness. REQ-98 is a shipped Phase 3 obligation whose test already executes
   against the real normalizer, so tightening it is pre-authorized as "a
   strengthening, not a scope breach". Deferring to the kata text would have left
   it blocked on work that already exists.
2. **Seam accretion: no.** `area:internal-table` carries 17 closed katas and five
   share this one's shape, but they patch five *different* oracles under five
   *different* REQs across four files (`rules_test.go`/REQ-44,
   `dump_test.go`/REQ-94, `roundtrip_test.go`/REQ-118,
   `format_test.go`+`rules_test.go`/REQ-9+34, `format_test.go`/REQ-18), mostly
   `batch:rdr-0002`. They share a review heuristic and a Go affordance, not a
   contract. REQ-98 has never been point-fixed; no closed kata touches
   `escape_shape_0009_test.go`. The accretion rule fires on repeated point-fixes
   to ONE contract — this is one-fix-each across a sweep.

### Ship

Steps 1–4 of the reviewed plan landed; **step 5 scoped out** (promoting
`EscapeShapeConformanceFixtures()` out of `table_test` may invert the existing
`table→resolve` dependency direction — a design fork, not kata work).

Both halves now read `ClearRuleID` / `ClearKey` from **one shared fixture
expectation** via a generic `onlyRowWithRuleID` selector, so row identity, tag
key, and value come from a single source and cannot drift. The whole-dump scan
and the `strings.Contains` substring match are both gone. `TestReq12_` now
asserts the refusal names its rule id. 4-bis also landed: the raw `"<clear>"`
literal in `internal/resolve/escape_shape_0009_test.go` became
`table.ClearSentinel` (no new import edge — that file is already `package
resolve_test` and already imports `internal/table`).

**Mutation-verified** (the point of the exercise — a strengthened oracle that
does not discriminate is worthless):

| Mutation | Result |
|---|---|
| clear renders multi-member `{"pending", ClearSentinel}` | **RED** — both halves fired |
| clear renders prefixed `{"x" + ClearSentinel}` | **RED** — both halves fired |
| clear moved to a different tag key | RED, but *for a load reason* — the base model declares one owned writer-served key, so a production-side key redirect is caught upstream by `malformed_accessor_binding`. Re-proved by perturbing the shared expectation (`ClearKey: "other"`), which fired the key assertion on both halves. Honestly recorded rather than claimed as a clean kill. |
| `+id+` dropped from the refusal message | **RED** — `TestReq12_` both subcases |

Production code untouched: `git diff main..HEAD -- internal/table/normalize.go
internal/table/model.go` empty; `normalize.go` restored byte-exact after each
mutation. Test-only diff, 2 files, +106/−34. Suite green, lint 0 issues.

## Wave 2 (drain re-sweep)

Resolved **empty** — normal termination. The one remaining open
`batch:rdr-0009` child is `6cek` (`kind:rdr-seed`), which is not flight-drainable
by construction: it exits to RDR authoring.

## Left open

| short_id | Why it is still open |
|---|---|
| `6cek` | `kind:rdr-seed` — JDR 0001 §JD-5 precondition precedence, an OPEN joint decision spanning RDR 0008 × 0009. Needs an RDR (or JDR amendment), not a ship. Carries an `## Open question`. |
