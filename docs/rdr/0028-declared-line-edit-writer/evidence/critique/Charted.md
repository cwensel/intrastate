# Charted — net-new scope from the critique lens (not folded into 0028)

Each row is a real concern the lens raised that falls outside this RDR's fence.
Recorded here so it is durable, then dismissed from this loop.

- **Symlink resolution semantics (D-15, `0028:C1` C1.3 `target:`)** — the clause
  says "symlinks resolved" without stating final-component vs whole-path
  resolution, TOCTOU discipline, or any bound on where the derived staging
  directory may land. Out of scope: the `path` carrier has the identical
  property today and `§consequences` already disposes of caller-bound
  indirection ("a caller who binds a hook file has bound a hook file"), so this
  is a pre-existing seam-wide question, not one `edit` introduces. The security
  framing in the origin row (writing outside the repo as a new hazard) was
  refuted on that ground. Suggested successor: an RDR pinning artifact-path
  resolution semantics across ALL carriers, where the answer can be one rule
  instead of a per-carrier clause.

- **A literal-anchor dialect (D-11 residue, `0028:A3`)** — anchor uniqueness is
  verified against a corpus that grows with every RDR seeded, and the corpus is
  documentation about the format the anchor matches, so decoys accrue over time.
  A3's own "If wrong" names the cure: a literal-anchor mode alongside RE2. Out
  of scope for v1 — it moves C1.2's "RE2" pin, which is a grammar change, not a
  clarification. The residual risk is now stated honestly in
  `§risks-and-mitigations` rather than implied away. Suggested successor: an RDR
  adding a literal/anchored dialect once a second consumer exists to shape it.

- **A round-trip guard for the carrier set (from pass B's A6 row, REFUTED as a
  defect of A6 but real as a future hazard)** — A6 is correctly Verified: no
  accessor dump/normalize round-trip exists in `internal/`, so its stated
  failure mode has no live target today. The hazard is that a LATER round-trip
  (a model dump, a table migration) omits the `edit` carrier and silently drops
  it. Out of scope: gating a path that does not exist is future-proofing against
  an unwritten feature. Suggested successor: whichever RDR introduces the first
  dump/normalize path owns a carrier-completeness test at that time.
