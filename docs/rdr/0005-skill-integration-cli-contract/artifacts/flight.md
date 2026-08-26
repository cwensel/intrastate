# Flight — batch:rdr-0005

`kata-flight --label batch:rdr-0005 --drain`, 1 wave (wave-2 re-sweep empty).
Scope-review gate ran first and graded all 7; 6 survived to the ship loop.

```
flight: 6 shipped, 0 stopped, 0 skipped
  shipped: 3cwy=b049a86  dkcf=26ed40c  r2ba=50506f9
           98n6=2715de8  8dg3=1301bef  g6vn=2665dcc
  routed:  sd16 → kind:rdr-seed (scope review; batch labels stripped)
```

## What shipped

| kata | fix | discrimination |
|---|---|---|
| `3cwy` | bare `flow` refuses through the gateway (`ValidateMode` then `command-error`, exit 2) | 3 behavioural arms failed pre-fix; 2 control arms passed both sides |
| `dkcf` | `Writer.Apply` clears `sealedKey` on a reachable write | stashed the fix: all 3 assertions failed; `TestReq104_AnUnreachableReadBackKeepsTheMutation` green both sides |
| `r2ba` | deny oracle asserts `result=="deny"` verbatim, repointed at `flowGateDenyModel` | coerce-to-allow mutation fails the new oracle at 3 sites and **passes the old one** |
| `98n6` | `canonicalValue` receives the `TagDecl`; `table.ConformValue` runs per scalar/member | 8 subtests failed pre-fix; each also asserts a conforming sibling is ACCEPTED, so refuse-everything fails too |
| `8dg3` | `writerFor` requires exactly one writer; loader guard widened to every declared tag | probed the hole first (two-writer fixture loaded clean); MVV zero-writer allowance preserved |
| `g6vn` | `flowRefusalInvocations` 16→21 entries, all five exit-3 classes | `envErr`→`GroupUserEnv` mutation fires the dead arm exactly 5×; `-count=40`, `-race`, zero flakes |

No production code was touched by `r2ba` or `g6vn` (test-only katas).

## Judgement calls worth recording

- **`98n6` deleted `flowbind.SetKeys`.** Its two callers were exactly the call
  sites 98n6 replaced with `decl.Kind == "set"`. Verified no other references;
  leaving it would be dead exported code whose doc comment named vanished call
  sites.
- **`8dg3` declined `testdata/neg/`.** That directory is the REQ-118 promoted
  spike set — `TestReq118` pins each member to the category it witnessed, so a
  new file would be an unpinned addition. Used in-test fixture mutation, the
  convention `TestReq24` already follows.
- **`dkcf`'s CLI oracle asserts CONVERGENCE, not one-shot repair.** The first
  corrected write still refuses honestly: the executor snapshots its
  pre-write baseline of the reader's protected keys while the artifact is
  still sealed, and `0004:C13` makes an unestablished baseline
  `read_back_incomplete`. Scoped to that invocation's own snapshot and out of
  scope here — pre-fix the loop never terminated at all.
- **One rebase conflict** (`g6vn` vs `98n6`, both appending to
  `flow_fixtures_0005_test.go`). Append-vs-append; both fixtures kept, full
  gate + flake check re-run on the resolved tip.

## `sd16` — routed out, not shipped

Every available fix names public surface RDR 0005 does not own: minting
`gate_allowed` is the additive-is-not-exempt violation DEV-7 already rejected
for `[model].revision`, and `Finding.Class` is assigned to **RDR 0006** by
REQ-18. A third reading is live — under REQ-17 (`message` must be
self-sufficient) the *doc comment* may be what is wrong, making the real fix a
comment correction with no behaviour change. A kata cannot adjudicate that.
