# Flight — batch:rdr-0028

Drain of the five kata roborev triage filed against the RDR 0028
implementation. One wave shipped; the re-sweep resolved empty and the drain
terminated normally.

```
flight: 4 shipped, 0 stopped, 1 routed  [over 2 waves]
  shipped: j7jw=5781dd5  sr1e=94f6119  d8kp=0bc7199  aj7z=d315e95
  routed:  mqwx → kind:rdr-seed (seam accretion; batch label stripped)
```

## Scope review routed one kata out of the drain

`mqwx` collected three weak-oracle defects, each confirmed by a surviving
mutant. It did not ship. A backlog scan for closed kata whose title names a
weak, tautological, or non-discriminating test oracle returned **sixteen**,
spread over at least seven records: 1xpr, 1qgz, x5vq, 878e, 5hy1, wkqs,
0gm0, g1pa, yn1v, stn1, 742a, py0x, s1r1, 9yeq, hjxk, g67s. Each was filed,
patched and closed on its own; this would have been the seventeenth.

Two prior point-fixes at one seam is the threshold that says the gap is a
missing design decision rather than a missing patch. Sixteen is well past
it. What an RDR must decide, from this kata's own evidence: what
observability seam production code owes its tests. Each defect it names is
an assertion forced to observe weakly because no better seam exists —
cardinality standing in for identity, an endpoint inode comparison standing
in for a write count, a claimed failure case with no way to induce the
failure.

`aj7z` was judged against the same bar and shipped. Its oracle named a
surface that did not exist when it was written and was never revisited when
D4 shipped that surface; the design lesson was already drawn and recorded,
only the patch was missed. No fork for a seed to decide.

## What shipped

**j7jw — a refusal named the wrong thing.** Argv-precondition refusals were
reported through the line-edit envelope, naming a write accessor and an edit
writer that never ran, with a finding whose rule subject did not exist.
Neither the invoking phase nor the entry's carrier discriminates: the
read-back preconditions sit behind an edit-table check and are still about
the reader's argv. The origin does. A new sentinel wraps the existing
request marker, so the exit-group check still answers true and no refusal
class was added.

**sr1e — the output contract promised a subject shape it does not carry.**
Four shapes ship; the contract named one. The first amendment left the false
claim standing above the correction, which a top-down reader hits first, so
the rationale was rewritten to rest on cardinality — true across all four
shapes. Each shape is now driven live through the CLI, so the prose is
pinned to the carriers rather than to itself.

**d8kp — the acceptance walk validated a model the loader rejects.** The
fixture declared D9's two-key shape; the helper registered four accessors
under one key. Steps 2–8 never touched the fixture. The helper now derives
its registry from the fixture it ships with, so the two cannot diverge
again.

**aj7z — the dump oracle could not see the surface it named.** A substring
assertion satisfied by the fixture's model id. Now a full-line match per
declared rule, an emitted-line count, and canonical ordering.

## Three fixtures passed for the wrong reason

Caught only because each fix was required to fail before it was believed:

1. A tokenless-refusal fixture was refused at load for an unrelated defect
   that also carries no token, so the assertion never reached the arm under
   test.
2. A one-key collapse guard failed as malformed TOML before the arity rule
   was consulted, which a category-blind assertion would have passed while
   proving nothing.
3. An ordering fixture with two elements per dimension passed both sorting
   mutants on a single run, but caught them only **34 times in 40**: an
   unsorted range over a two-element map reproduces sorted order often
   enough that a real regression would survive CI about one run in seven.
   Four elements per dimension gives 40/40.

The third is the one worth carrying forward. Fixture *width* is
load-bearing for an ordering assertion, and a single-run mutant check hides
it.

## Notes

- Every ship rebased where main had moved. One fast-forward would have
  silently reverted two hook files; the scope check caught it, not the
  merge.
- A repo hook now rejects agent-attribution trailers at commit and push.
