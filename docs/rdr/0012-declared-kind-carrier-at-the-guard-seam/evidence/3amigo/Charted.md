Model: claude-opus-5[1m]

# Charted to successor — cli/0012 3amigo

One finding was real but net-new scope. Recorded here and dismissed
from this loop rather than absorbed into the RDR.

## CLI int canonicality (`--tag n=07`, `--write n=+1`)

- **Finding**: PM-1 / IMPL-5 / QA-7. C5 originally reached the CLI
  through `conformKind`, refusing non-canonical int spellings at
  `--tag` / `--write`. Narrowing C5 to the predicate ingress (the
  3amigo resolve) leaves the CLI accepting them, and under C2's parsed
  comparison a `--tag iter=07` against `iter eq 7` now MATCHES where
  today's raw-string seam prunes — a silent GuardFalse → GuardTrue
  flip, measured in `evidence/spikes/a4-typed-compare.md`.
- **Why out of scope here**: closing it at `conformKind` refuses
  literals that two BOUNDARY-marked tests in
  `internal/table/emit_grammar_0024_test.go` admit, encoding `0024`'s
  REQ-7 ("every kind check is a LEXICAL check … no value is parsed
  into a typed representation, canonicalized, or converted") and
  REQ-9 ("non-canonical literals settled in the permissive
  direction"). Retiring those is a cross-RDR contract change this RDR
  has no mandate for, and no measurement to justify: S6's zero counted
  GUARD ATOMS, never CLI, `[emit]`, `[initial]` or `[rule.write]`
  values. Both tests verified passing on `main` (afcc69d).
- **Suggested successor**: an RDR scoped to int-value canonicality
  ACROSS ingresses, whose job is to weigh `0024`'s lexical-check
  disposition against the typed-comparison agreement C2 wants, and to
  decide them together rather than from one side. It should carry the
  measurement this one lacks — the per-ingress census of non-canonical
  int spellings in the corpus and in persisted artifacts (see the new
  Pending assumption A6, which measures the held side).
- **Disposition in this loop**: dismissed as net-new scope; the
  residual is recorded honestly in §Failure Modes (F1) so the successor
  inherits a stated gap rather than a silent one.
