# Flight — `batch:rdr-0004`

`/kata-flight --label batch:rdr-0004 --drain`, run as `rdr-implement-triage`
Phase 4 (`lib-land-rdr §land-flight`) after RDR 0004 landed on `main`.

```
flight: 2 shipped, 0 stopped, 0 skipped [over 1 wave]
  shipped: bkf3=ba19b7b  742a=a30be0d
  stopped: (none)
  skipped: (none)
```

## Wave 1

Resolved `[bkf3, 742a]` — 2 eligible, 1 dropped (`p63c`, `kind:rdr-seed`:
its deliverable is an RDR draft, so it routes to RDR authoring and
`kata-ship` gate 4 refuses it by design). No intra-batch `blocks` edges,
so priority ordered the wave: `bkf3` (p1) then `742a` (p2).

**Review gate** (`/kata-scope-review`, default-on): both IN-SCOPE with
source-grounded plans embedded. No merges, no out-of-scope closes, no
rdr-shaped demotions, nothing surfaced. Seam-accretion count at
`internal/accessor` was 0 (fresh seam — no repeat point-fixes, so no
accretion route to an RDR). `§batch-correction`: both correctly carried
`batch:rdr-0004`; no change.

### `bkf3` → shipped `ba19b7b`

*Executor skips planned-owned-tag read-back when the pre-write baseline is
incomplete.* Two commits:

- `a3435ac` — the chartered fix. `Write` no longer short-circuits on an
  incomplete pre-write baseline; the re-read always runs, `0004:C12`'s
  baseline-independent planned-owned conjunct is evaluated, and a
  demonstrated inequality yields `read_back_mismatch`.
- `ba19b7b` — a **second-order defect the refine pass found**: the first
  commit moved the planned-owned conjunct ahead of the `baselineUnread`
  arm but left the *protected* conjunct behind it. With two or more
  protected keys, a partial snapshot can establish one key while failing
  another; a write clobbering the *established* key was still reported
  `read_back_incomplete`, naming the unrelated unread key. Same defect
  class as the kata, one conjunct deeper. Fixed by passing `before`
  instead of `nil` so one `verifyReadBack` call spans every evaluable
  conjunct.

**A plan correction worth recording.** The scope-review plan said to set
`applied = true` on the mismatch arm. The resolve agent declined, and was
right: `0004:C14` scopes the applied-but-unverified sense to exactly two
classes (`read_back_incomplete` and post-mutation `timeout`), while a
mismatch is verified-and-wrong. `TestReq67_MismatchAndIncompleteAreDistinct…`
and `TestMVV_…/8_post_mutation_reporting_and_no_compensation` both assert
`Applied() == false` on a mismatch as exactly that discriminator; setting it
would have required weakening two shipped tests, which §scope-discipline
bars. Regression tests ADV-6/7/8.

### `742a` → shipped `a30be0d`

*Strengthen RDR 0004 accessor suite oracles.* One test-only commit; the
production diff is empty. Five oracles that asserted a signal's *presence*
now assert its *identity*:

1. `planAdvancing` + a test proving only `plan.Writes` is applied, never
   `NextTags` — discharging deviation **D1**'s remaining half (RDR 0009's
   obligation, `0009:1515-1527`).
2. Timing bounds now `20 * fixtureTimeout` (1s) instead of `blockFor` (1h),
   with the binding still delaying an hour; read-back arm requires
   `ClassTimeout`.
3. REQ-101's simultaneity arm asserts `ob.sawUnreadable`.
4. `ctxOf` plants an unexported-key `context.WithValue` sentinel; the
   context-bound test asserts the binding observes it, so an executor
   discarding the caller context now fails.
5. ADV-1/1b `Refusal.Keys` assertions are exact set equality, not
   containment.

Every tightening was **mutation-verified** — each assertion was shown to
fail against its named mutant before landing. One documented limitation:
ADV-1c's single-key plan structurally cannot discriminate over-attribution;
ADV-1/1b carry that check, and the comment says so.

Refine returned zero findings on the first pass.

## Residual

`p63c` (`kind:rdr-seed`) stays open by design — *Reader cardinality per
artifact role is unbound at the RDR 0002/0004 seam*. It carries an
`## Open question` and needs RDR authoring, not a patch: the choice between
"a role admits at most one reader" (a new RDR 0002 binding validation) and
"read-back aggregates across same-role readers" (restating `0004:C12` over a
reader set) is a cross-RDR contract decision.
