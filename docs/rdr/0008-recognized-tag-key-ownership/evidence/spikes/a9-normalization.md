Model: claude-opus-5[1m]

# Spike — A9: which normalized field does a tag predicate produce?

Stage 6 (reconcile). Settles A9: whether a tag-predicate reference to a
recognized-provenance declaration normalizes into a **view-read** row field
(`Row.Match` / the guard field) or into `Row.Outcome`. Under the second
reading the Phase 2 fixture rename would be a semantics change on a Final
peer's canonical evidence and would route to RDR 0002.

## Command

Re-ran RDR 0002's own normalizer spike — the peer's executable artifact, not
a stale captured output — against both committed fixtures:

```sh
cd docs/rdr/0002-transition-table-as-reviewable-data/evidence/spikes
GOCACHE=/private/tmp/intrastate-rdr0008-gocache GONOSUMDB='*' GOPROXY=off \
  go run . rdr-fixture.toml kata-fixture.toml
```

## Output (verbatim, reproduces `output.txt` byte-for-byte)

```
MODEL rdr rows=3 outcomes=round-clean,verdict-flapping,reconcile-block,finalized
rdr.continue-prelock kind=transition source=rdr:prelock match=[all:iter.lt=3; outcome.eq=round-clean; profile.in=[large foundational]; stage.eq=prelock; status.eq=Draft; unless:profile.eq=small] write=[iter=2; stage=prelock]
rdr.draft-no-match-escape kind=escape source=rdr:escape match=[outcome.eq=round-clean; status.eq=Draft] write=[escape=no_match]
rdr.reconcile-rewind kind=transition source=rdr:reconcile match=[outcome.eq=reconcile-block; status.eq=Draft] write=[prelock_lens=<clear>; rewind_scope=assumptions; stage=resolve; status=Draft]
MODEL kata rows=2 outcomes=accepted,needs-work,closed
kata.review-accepted kind=transition source=kata:review match=[all:owner.eq=current-session; outcome.eq=accepted; phase.eq=review; status.eq=open; unless:status.eq=closed] write=[phase=ship; status=accepted]
kata.review-needs-work kind=transition source=kata:review match=[outcome.eq=needs-work; phase.eq=review; status.eq=open] write=[phase=resolve; status=open]
```

## Verdict — A9's FIRST reading holds

**All three tag-predicate positions normalize into one predicate set; none
produces an outcome field.**

1. **The normalized row has no outcome field.** The spike's `Row` type
   (`0002-…/evidence/spikes/main.go::Row`) is
   `{Model, Rule, Kind, Source string; Match, Write []string}` — no `Outcome`
   member exists to receive a predicate.
2. **`[rule.match.<tag>]` lands in `match`.** `outcome.eq=round-clean`
   appears inside `match=[…]`, unprefixed and sorted alphabetically among
   `stage.eq`, `status.eq`, `profile.in` — handled as an ordinary tag
   predicate with no special-casing on the name `outcome`.
3. **Guard positions land in the same set, prefixed.**
   `0002-…/evidence/spikes/main.go::normalize` merges all three positions:
   `mergePredicates(predicates, "", rule.Match)`,
   `mergePredicates(predicates, "all:", rule.Guard.All)`,
   `mergePredicates(predicates, "unless:", rule.Guard.Unless)` — visible in
   the output as `all:iter.lt=3` and `unless:profile.eq=small`. So
   `[rule.guard.all.rewind_target]` (`0003-…/evidence/spikes/guard-fixture.toml`)
   normalizes to a view-read predicate too, not to a gate.
4. **The root `outcomes` alphabet is model metadata, printed separately
   (`outcomes=…` on the MODEL line) and never routed into a row** — so the
   alphabet is not what the row's predicate is bound to.

## Why `Row.Outcome` is not the destination

RDR 0002 never names `Row.Outcome` (zero `Row.` hits in the document); it is
a kernel field owned by RDR 0001, fed from `Input.Recognized`, not from the
table. The kernel gate is `row.Outcome != in.Recognized`
(`internal/resolve/resolve.go::Resolve`), comparing the input tuple's value —
while `::assemble` separately injects that same value as a tag under
`::recognizedTagKey` so a row may *also* match it as a tag (RDR 0001 REQ-17;
`deviations.md` D2: "a future RDR (0002 normalization) must spell the same key
for a table row to match on the recognized outcome **as a tag**").

The two paths are complementary, which is why both readings looked live: the
fixture's `[rule.match.outcome]` is a tag-level restatement of a value the
outcome gate reads independently.

## Consequence for this RDR

- The Phase 2 rename is **mechanical**: renaming the three declarations and
  carrying all six tag-predicate reference sites changes no gating semantics.
  No route to RDR 0002 is required; Phase 2 is unblocked.
- Noted rather than left implicit: the rename makes those fixtures *newly
  functional* against the kernel — the predicate key `outcome` cannot match
  `assemble`'s injected key `recognized` today, and after the rename it can.
  That is a strict improvement, not a regression, and nothing at HEAD depends
  on the current behavior: no normalizer exists under `internal/`, so no code
  path turns these fixtures into a `resolve.Table`.
