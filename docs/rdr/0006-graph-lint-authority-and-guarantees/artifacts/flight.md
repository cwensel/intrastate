# Flight — batch:rdr-0006

`/kata-flight --label batch:rdr-0006 --drain`, run as the Phase-4 tail of
`/rdr-implement-triage 6 --close-and-flight`.

```
flight: 6 shipped, 0 stopped, 0 skipped  (over 2 waves)
  shipped: a5an=506ed0c  aqzt=593fe17  py0x=25e54e3
           ffgb=f3fad20  04xp=e20db68  dkaw=0e18bf7
  stopped: none
  skipped: none
```

Wave 1 resolved 6 eligible katas (3 `kind:rdr-seed` correctly excluded — they
route to RDR authoring, and kata-ship gate 4 refuses them). Wave 2 (the drain
re-sweep) resolved empty: normal termination.

## Scope review (gate, wave 1)

All 6 came back IN-SCOPE — no closes, no merges, no demotions. Detail in
`_scope-review/04xp-aqzt-ffgb-py0x-a5an-dkaw.md` (gitignored scratch).

The review's main product was **ship order**, derived from two real couplings
that a naive priority sort would have got wrong:

1. `aqzt → py0x` — shared fixture `lint_fixtures_0006_test.go::mvvMultiDefect`:
   aqzt rewrites its rows, py0x asserts against them.
2. `a5an → py0x` — py0x converts REQ-98 from substring containment to a
   line-oriented parse, which is exactly the property a5an's low-severity
   argument said was unpromised. a5an's escaping is a prerequisite for py0x's
   stricter parser to be sound.

Combined: **a5an → aqzt → py0x**, with `04xp`, `ffgb`, `dkaw` independent.

The review also tested and REJECTED a merge: `py0x` and `aqzt` share four files,
but mapping every claimed line to its enclosing test function showed **zero
test-function overlap**. They also have opposite verification protocols —
py0x's DoD is "suite still green", aqzt's is "suite goes red on demand".

## What shipped

| kata | sha | what |
|---|---|---|
| `a5an` | `506ed0c` | `clierr` quotes identity values; `sanitizeLine` folds CR/LF/VT/FF/NEL/LS/PS in Message, so one finding renders as one line |
| `aqzt` | `593fe17` | all 10 weak RDR 0006 fixtures rewritten to fail against a deliberately broken engine; each mutation-proved RED |
| `py0x` | `25e54e3` | 8 oracles tightened from containment to exact set/identity equality |
| `ffgb` | `f3fad20` | **production fix**: `heldValues` abstracts non-finite-domain tags to `OpaqueValue` (deviation D15) |
| `04xp` | `e20db68` | **production fix**: `escapeField` encodes composite keys injectively; members terminated not joined (deviations D16, D17) |
| `dkaw` | `0e18bf7` | regression test pinning the `in` atom's stable reason |

## Defects the pipeline caught that the katas did not name

- **aqzt**: the node-ceiling assertion sat inside `if len(Reach(...)) > ceiling`
  and had **never executed** — measured at 2 reachable nodes against a ceiling
  of 4096. `checkNodeCeiling` was deletable with the suite green.
- **aqzt**: REQ-88's ordering loop compared **zero** adjacent pairs, because the
  shared fixture put rule-scoped and element-scoped findings in disjoint
  `(Model, Code)` buckets.
- **04xp refine**: D16's own premortem was **wrong**. Backslash escaping is not
  sort-preserving — a sort compares against the neighbour's character, not the
  delimiter, so `["a,b"]` precedes `["a0"]` authored and follows it escaped.
  `compareAtoms` had been moved onto the escaped rendering, ordering atoms by an
  encoding artifact. Now sorts on canonical authored members. Recorded as D17.
- **04xp refine**: `escapeJoin` was not injective over cardinality — `[]` and
  `[""]` both rendered empty, and both are authorable. Members are now
  terminated rather than joined.
- **ffgb refine**: the new tests pinned both no-finite-domain arms but had no
  *finite* control, so a `heldValues` keyed on declared kind would have passed.
  A bounded-int control now pins that the gate keys on `guard.AssignmentCount`.
- **py0x refine**: `strings.Contains(line, "2048")` also matched `20480`.

## Spun off during the flight

- `8ea9` (`kind:rdr-seed`) — graphlint duplicates guard's unexported
  single-value operator set; unifying them means exporting new API from the
  package RDR 0003 owns while RDR 0006 is `Implemented`, so it needs its own
  scoping. Raised independently by both the scope review and dkaw's resolve.
- One finding routed `aqzt → py0x` mid-flight (assert that `mvvMultiDefect`
  yields one group with all four rule IDs) and shipped inside `py0x`.

## Still open (not drainable — route to RDR authoring)

`9en6`, `pz9z`, `8ea9` — all `kind:rdr-seed`. `pz9z` is deviation D12, the one
open author decision from the RDR 0006 launch.
