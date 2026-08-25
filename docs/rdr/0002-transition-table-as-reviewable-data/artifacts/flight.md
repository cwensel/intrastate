# Flight — `batch:rdr-0002`

`/kata-flight --label batch:rdr-0002 --drain`, run as the `§land-flight` tail of
`/rdr-implement-triage 2 --close-and-flight`.

```
flight: 13 shipped, 0 stopped, 0 skipped   (over 1 wave)
  shipped: rb25=2ad7a3e  ekpq=4a0bcc3  fmhr=bf686f4  g67s=737616a  meab=29bf3fb
           pwff=a9ab53e  q7jc=1b39c35  rjdy=e03ade6  sje2=c561c6e  v3yh=c0ec083
           xgt0=d37dc62  hjxk=f4e6b19  sp9e=e6d9f7b
  held:    06pe (kind:rdr-seed — routes to RDR authoring, not the drain)
```

Wave 1 resolved 13 eligible katas (no dependency edges; ordered by priority then
short_id). The scope-review gate passed all 13 as IN-SCOPE with embedded plans.
The `--drain` re-sweep returned 0 eligible, so the flight terminated normally.

## Scope-review gate

13 reviewed, 13 IN-SCOPE, 0 out-of-scope, 0 merged, 0 rdr-shaped, 0 surfaced.
Every kata had been grounded once already by triage, so none was dropped — but
the review still corrected five plans before any code was written (see
`_scope-review/batch-rdr-0002.md`, gitignored scratch).

## What the flight actually found

**Six of thirteen reviewed plans did not survive contact with source.** Each was
deviated from deliberately, with the structural reason recorded on the kata and
independently validated in refine. None was forced, and no oracle was weakened.

| kata | plan defect | substitution |
|---|---|---|
| `fmhr` | same-block `contains` pair is unauthorable (TOML rejects the duplicate table; contexts expose match blocks only, `load.go:408`) | tie built from two contexts authoring match `in`; oracle targets expansion suffix order |
| `g67s` | swapping the write target is insufficient — an owned tag needs `[read.*]`/`[write.*]` bindings + `[initial]` | added the loader-mandated scaffolding; minimality proved by dropping each block |
| `meab` | **the review's own "mutation-verified" RED recipe was unsound** — `NextTags`/`Writes` are two independent `cloneTagValues` calls, so the mutation aliases only across rows | substituted REQ-82's actual within-row aliasing |
| `rjdy` | spike fixtures are unloadable (all 71 carry `kind = "string"`, predating the `"scalar"` rename — deviations.md **D1**'s recorded residual) | per-fixture category pinning; the plan's aggregate census was provably vacuous |
| `sje2` | dropping the `kind` column refuses at load (`[dump].order` names it), so it fails for the wrong reason | mutation substituted with "kind renders empty" + a row-association mutant |
| `sp9e` | the same-rule-id unexpanded/expanded pair is **unreachable** (`expand` emits uniform suffix arity per rule; duplicate rule ids refused) — the plan's guard would fire on every input | second expanding atom so element 1 is the sole discriminator |

**Two production behavior changes landed**, both with fenced RDR-0003 authority:

- **`q7jc`** — `normalize.go` refuses a rule carrying no local `[rule.match]`
  block (`malformed_rule_shape`), firing *before* inheritance so a context cannot
  stand in for the rule's own selection. Deviation **D12** records the category
  reuse past REQ-106's parenthetical, citing D3/D9 precedent.
- **`xgt0`** — `load.go::tagDecl` refuses the `single_valued` marker on `set` and
  `scalar` kinds per `0003:1389`, keyed on **presence** (`*bool`) to match the
  sibling per-kind checks. `Categories()` still returns 25.

Both were swept against every tracked `.toml` **and** every inline Go raw-string
document to confirm no legal document newly refuses.

## Two oracles that were themselves vacuous

- **`sp9e`** — refine found the Phase-1 fix *still* could not witness suffix
  element 1: the Cartesian product already emits it pre-sorted within each
  element-0 group, so a comparator reading only element 0 passed silently. A
  comparator-level pin (`normalize_internal_test.go`) closes it. The initial
  whole-comparator inversion had caught the mutant for the wrong reason.
- **`sje2`** — refine found `cloneRows` left the unexported `setKeys` aliased and
  untestable with the fixture in use (`rdr` has 0 rows carrying it, `kata` has 2).
  Fixed via a package-internal `export_test.go` hook plus the `kata` leg.

## Notes

- No kata stopped, no rebase conflict, no worktree leaked. `main` advanced
  `a0a980b` → `e6d9f7b`; every branch ff-merged and was torn down.
- Two katas' claims were **corrected rather than parroted**: `ekpq` (dropping
  `escape`/`gate` fails the old oracle too — only `outcome` was genuinely blind)
  and `fmhr` (a genuine root key already refuses `unknown_schema_field`, so the
  proposed category change was unnecessary).
- `TestReq118`'s byte-for-byte fixture pin, landed mid-flight by `rjdy`, was
  spot-checked by `hjxk`'s refine — perturbing a genuinely pinned fixture failed
  with "differs from its approved spike content", proving the pin bites.
