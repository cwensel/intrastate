# Flight — `batch:rdr-0019`

Drain of the bug children roborev triage filed for RDR 0019, run at the
end of the implementation's landing tail.

```
flight: 1 shipped, 0 stopped, 0 skipped, over 2 waves
  shipped: es2r=0179168
```

## Wave 1

Resolved `[es2r]`. `ptn7` was filtered before the ship queue: it carries
`kind:rdr-seed`, which exits to RDR authoring rather than draining as a
mechanical fix (kata-ship gate 4 refuses it by design).

**Scope-review gate — IN-SCOPE.** The review rebuilt the binary and
re-ran the reproduction independently rather than trusting the kata body,
and resolved the latent contract question: `0019:C1` REQ-33 (ordering:
carrier before the artifact-binding family) and REQ-44 (a no-match must
not mint a separate missing-reader code) are **orthogonal** — REQ-44
constrains the code minted, REQ-33 its ordinal position — so a two-pass
split satisfies both and is not a contract fork. Seam-accretion count at
`::initCarrierGate` was 0 (the verb shipped in this same batch), so there
was no missing-design-decision signal routing it to an RDR seed.

**Ship.** `internal/cli/flow_initstate.go::initCarrierGate` is now two
passes over the same `needed` slice: pass one answers every carrier
question and records reader-less roles; pass two returns the refusal for
the first such role in `needed` order, preserving the deterministic writer
choice and the message text verbatim. `::initReaderFor` was not touched
(REQ-46 forbids re-ordering `Definitions` — this re-orders `needed`-derived
refusals, never the registry slice), and both passes stay field-reads
only, so REQ-45 / S10's no-accessor-invoked property survives by
construction.

Red/green was proven by reverting only the source file: the new two-writer
row failed with `flow-artifact-missing` and passes as
`flow-init-carrier-unsupported`; an order-independence control (the
command-backed writer sorting first) passes both before and after, so the
fix is not a flipped bias. Named regressions REQ-33/36/37/44/45/46 and the
S10 case all stayed green.

Refine found no issues (per-commit and branch reviews both clean). Full
suite and `make check` green; ff-merge first try; 2 files, 289 insertions,
8 deletions — inside the scope floor.

## Wave 2

Re-sweep of the standing selector resolved empty. No spin-offs were minted
during wave 1, so the drain terminated normally.

## Residual

`ptn7` (`kind:rdr-seed`) remains open by design: RDR 0019's scenario S8
states a fixture recipe that cannot be built against the shipped seam, and
the correction is scenario TEXT in a Final record, which is never amended
in place. It routes to RDR authoring, not to a flight.
