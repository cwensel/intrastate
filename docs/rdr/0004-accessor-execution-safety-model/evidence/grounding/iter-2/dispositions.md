Model: claude-opus-5[1m]

# Grounding Dispositions — iter-2

- **fixed** — origin: A8 Evidence cited three spike symbols
  (`completeRead`, `truncatedRead`, `absentKeyRead`) that do not resolve in
  `main.go` — sections touched: Critical Assumptions A8. Replaced with the
  three fixture builders that actually produce the cited dispositions, in the
  same order as the sentence describes them: `main.go::newArtifacts` (complete
  read), `main.go::newPartialArtifacts` (unreadable key →
  `incomplete_read`), `main.go::newSparseArtifacts` (genuine absence →
  `<absent>`). Every symbol in the A8 Evidence line now resolves, as do
  `internal/resolve/resolve.go::missingOwned` and `TagSet.has`.

Needs verification: None. The edit changes evidence citation form only. The
behavior A8 asserts was independently confirmed at `main.go::read` and at the
transcript (`output.txt:11-13`); no load-bearing claim, normative clause, or
exactness word was added or altered, so §amendment-sweep has no predicate token
to propagate.

Charted: None.

Tiebreakers: None.

## Mini-check cue read (first lens pass of this re-entry)

The draft carried no mini-check tables — it reached Final without the cue read
having been run — so this pass owes the full read. Cues taken from Normative
Contracts, the Round-Trip / Inverse Invariants block, the MVV, and Testing
Strategy; not from a file grep.

| Mini-check | Fired | Cue read |
| --- | --- | --- |
| disposition | **yes** | The contracts enumerate input classes and set an outcome for each: timeout, execution failure, incomplete read, gate denied/indeterminate, capability mismatch, unknown accessor, non-owned write, read-back mismatch. Table written as **Disposition Table**. |
| test-discriminability | **yes** | MVV Scenario 1 ("fail before resolution") and Scenario 5 ("without package prints") both pass by absence-of-error. Table written as **Oracle Discriminability**. |
| round-trip / fidelity | **yes** | The draft declares an explicit `write -> read` Round-Trip / Inverse Invariant, and the read contract adds a totality claim over the requested key set. Table written as **Fidelity Table**. |
| desk trace | **yes** | Eleven normative clauses bear on one output surface — the structured accessor result the MVV ends at. Table written as **Desk Trace**. No CONTRADICTION row. |
| source-authority census | no | The read-completeness rule is a single-arm branch split at one boundary, not a fallback, a derived/propagated output, or two competing sources of truth. The executor is the sole writer; there are no sibling arms to census. |

Determinacy trigger: **n/a** — the contract legislates execution safety and
refusal classing at an I/O boundary. It carries no parse/deparse,
import/export, compose/decompose, hashing, identity, or migration behavior, no
data-model field whose ownership could be read more than one way (tag ownership
is fixed by RDR 0001's owned/observed/recognized split), and the MVV rests on
per-call branch dispositions rather than multi-step transformation fidelity.

- **fixed** — origin: mini-check cue read, test-discriminability — sections
  touched: MVV Scenario 5. The scenario could pass vacuously by producing no
  output; it now asserts positively that the returned structured value carries
  the refusal class, and captures stdout/stderr around the call requiring both
  empty. §amendment-sweep: the other no-print sites (normative clause at
  Normative Contracts, Existing Infrastructure Audit's `respond::Fail` row,
  Testing Strategy's "without direct stdout/stderr output" line) were re-read
  and all agree; the edit strengthens the proof, it does not restate the rule.
