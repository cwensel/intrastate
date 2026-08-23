Model: claude-opus-5 (fresh context; rdr-common §strong-consult)

# Strong consult — grouping predicate + no-false-green guarantee

Invoked under the resolve prompt's TIEBREAKER-REDUCTION GATE after the fork
failed to collapse on cluster evidence. One consult, queried once. Factored
brief: the two claims in tension plus the shipped kernel, with no prior verdicts
or justifying prose supplied.

## Packet

```text
verdict: NEEDS_DECISION
blocking: yes
next_action: Route the merged-node coverage question to a human fork decision,
  then rewrite 0006's "Source state is a tag-set" paragraph and its
  Load-Bearing false-green guarantee before 0006 finalizes; do not reopen 0003.
summary_50w: 0006's node keying is a genuinely different partition than 0003's,
  and 0006 asserts it while its own normative clause forbids it. The
  no-false-green claim holds for overlap but is backwards for coverage at
  merged nodes. 0006 must move; 0003 stays locked.
```

Evidence: `0003::Technical Design` (row-group definition), `0003::Normative
Contracts` ("Match keys are not product dimensions"), `0003::A10` + its cluster
gate venue, `0006::Technical Design` ("Source state is a tag-set"),
`0006::Normative Contracts` (MUST NOT define a second grouping predicate),
`0006::Load-Bearing Decisions` (union widening; no-false-green),
`internal/resolve/resolve.go::Resolve` (`view.matches(row.Match)` over one
concrete assembled view, then `gate`).

## Findings

**Q1 — different partition, a cross-cut (not a refinement).** 0003 maps each row
to exactly one group by authored match pattern. 0006 puts one row in *many*
groups (it says so: "a row participating in several reachable nodes is checked
once per node") and puts rows with *different* match patterns in one group
whenever both are satisfiable in that node. Neither coarsens nor refines the
other. Additionally, 0006 projects over **owned** tags only, dropping
observed/recognized match atoms that `Resolve` does filter on.

**Q2 — the guarantee splits by invariant.**

- **Overlap: sound.** A merged node admits a superset of concrete views, so any
  two rows co-enabled at some concrete state are co-enabled at the merged node.
  Spurious reports possible; missed ones not. 0006's claim holds here.
- **Coverage: unsound — it inverts.** At a merged node the union of
  row-accepted assignments is a *superset* of the union at each concrete state
  the node abstracts, because rows satisfiable only in some abstracted states
  still contribute assignments. A gap real at one concrete state gets filled by
  a row that is not a candidate there: the union closes, lint certifies green,
  and `Resolve` at that concrete state returns `no_match`. **False green in the
  check the RDR exists to provide.**

  Consult's stated uncertainty: only whether product widening sometimes offsets
  this. It does not reliably.

So `0006::Load-Bearing Decisions`' blanket "can only produce a false positive,
never a false green" is not sound for invariant 4. Corroborates critique pass B
C-1 over pass A's false-positive-only framing (diff.md divergence 1).

**Q3 — substantive, not a drafting fix.** Rewording cannot close it. Deleting
the offending paragraph leaves 0006 with no answer to "what is a source state" —
0003 names one but never defines its identity, since the cluster has no named
states. **0006 raised a real gap**; it just answered it in a way that breaks the
proof. Minimal reconciliation the consult offers:

1. Keep **0003's match-pattern keying** as group identity.
2. Demote 0006's node relation to a **reachability filter and witness** — prove
   a group at each reachable node satisfying its context, with the coverage
   union computed from that group's **authored rows only**, never pooled across
   contexts.
3. Certify coverage **only at nodes whose per-tag value sets are singletons for
   every participating match key**; emit `graph-unprovable-coverage` at merged
   nodes instead of green.
4. Overlap may keep the merged-node reading.

**Q4 — 0006's text moves; no route-back on 0003.** 0006's Technical Design
contradicts 0006's *own* Normative Contracts, so the Draft is internally
inconsistent regardless of which design wins — the burden is on 0006 either way.
0003 is Final/locked and its clauses agree with the shipped kernel. The A10
venue was already homed at "RDR 0006's refine" by the 2026-08-22 cluster gate,
so this resolves where the locked document already designated.

**A2's `Verified` is not earned** — its claim to close `0003::A10` is false as
written; it returns to `Pending`.

## Disposition

`NEEDS_DECISION` → escalated to the human per §strong-consult. The consult
collapsed the *technical* fork (coverage is a false green; 0006 moves, 0003
stays locked) but the remedy is a substantive redesign of this RDR's grouping
and guarantee, which is a design call, not a lens fix.
