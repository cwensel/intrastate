# Flight drain — `batch:rdr-0023`

`/kata-flight --label batch:rdr-0023 --drain`, run as the landing tail of
RDR 0023's implementation. One wave, one kata, no spin-offs.

```
flight: 1 shipped, 0 stopped, 0 skipped (over 1 wave)
  shipped: 0gm0=7541a69
```

## Wave 1

**Resolved:** `[0gm0]` (P3, no blockers — only `related` edges to `srz2`
and `rg0e`, which do not gate ordering). No drops.

### Review gate — `/kata-scope-review 0gm0` → IN-SCOPE

The review verified all three of the kata's claims **by mutation** rather
than by reading, and corrected the triage on two of them:

| Claim | Verdict | Evidence |
|---|---|---|
| 1. REQ-49 one-lexical-site check is vacuous | **REAL** | Injected a second `GetBool(planOnlyFlagName)` read into `runFlowResolve`; the oracle stayed GREEN. The conjunct is dead behind the `siteAt < 0` `t.Fatal`. |
| 2. exit-3 arm is not a reader-pass refusal | **Mislabel, not lost coverage** | The arm refuses at `[gate.unreachable]` via `runGates` (`flow_resolve.go:257`), code `flow-accessor-failed` — but it *does* discriminate mutation M-B. The finding overstated it. |
| 3. stale PLAN key leaves the suite green | **DOES NOT REPRODUCE** | `TestAdversarialF2_TheShippedDeclarationTableIsTotalOverThePayloadStruct` *is* the ADV-1 remedy; mutation M-A turns it red on both directions. The kata quoted that test's own status comment — describing the gap it closes — and read it as a live gap. Triage error, recorded rather than silently dropped. |

**Seam-accretion trail dismissed.** Four closed katas (`stn1` on RDR 0010;
`g1pa`, `yn1v`, `gq75` on RDR 0011) looked like prior point-fixes of one
defect class at one seam, which under the `§seam-accretion` rule would
route this to an RDR rather than a fifth patch. Grounding dismissed it:
those four are *fixture-reachability* defects (an oracle whose fixture
cannot express the clause, fixable only by authoring a discriminating
input), whereas this is a single dead boolean conjunct. Decisive
counter-evidence: RDR 0023 already ships the standing remedy an RDR would
have been asked to name — Phase 3b ran named mutations M-A/M-B against
each new oracle and recorded them in `verification.md`, both independently
re-verified during the review. A mutation gate already operates at this
seam.

### Ship

`resolve d872257` → `refine ffbc8dd, c5237a8` → squashed to **`7541a69`**,
ff-merged to `main`. Tests and lint green pre- and post-squash; the squash
tree is byte-identical to the un-squashed tip.

Refine ran 3 iterations over 4 findings: **3 fixed, 1 dismissed, 0 pushed.**

The material refine fix (`ffbc8dd`) widened the oracle beyond what the
resolve pass delivered: REQ-49's "exactly ONE lexical site" is a *package*
claim, but the Phase-1 oracle scoped only two files. `planOnlyReadSites0023`
now walks every non-test file in `internal/cli`, returning sorted
`<file>::<function>` sites, and asserts both that there is exactly one and
that it is `flow_projection.go::projectResolvePayload` — the
enclosing-function check the review's plan asked for. Verified by two new
mutations (a read in a third file; a read relocated into a helper), both
invisible to the earlier form.

The dismissal was itself mutation-tested: a reviewer proposed that
`Changed()`/aliased-name access could evade the AST oracle. Injecting that
exact evasion left the AST oracle green as predicted, but turned **23
shipped tests red** — including the REQ-47..53 upstream-read test and
ADV-2/ADV-3's differentials. The behaviour is covered; the finding was
refuted with evidence rather than waved off.

### Wave 2

Re-sweep of `batch:rdr-0023` resolved **empty** — no spin-offs were minted
during wave 1. Drain terminated normally.
